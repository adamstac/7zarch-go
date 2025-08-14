package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/debug"
	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/display/modes"
	"github.com/adamstac/7zarch-go/internal/query"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// parseHumanDuration supports 'd' (days) and 'w' (weeks) in addition to time.ParseDuration units
// Kept local to this file to avoid changing behavior elsewhere
func parseHumanDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}
	if strings.HasSuffix(s, "d") {
		n, err := strconv.ParseInt(strings.TrimSuffix(s, "d"), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid days: %w", err)
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	if strings.HasSuffix(s, "w") {
		n, err := strconv.ParseInt(strings.TrimSuffix(s, "w"), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid weeks: %w", err)
		}
		return time.Duration(n) * 7 * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

// listFilters collects flags for registry listing
type listFilters struct {
	details      bool
	notUploaded  bool
	pattern      string
	olderThan    string
	onlyManaged  bool
	onlyExternal bool
	onlyMissing  bool
	onlyDeleted  bool
	status       string
	profile      string
	largerThan   int64
	debug        bool
}

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List archives in the registry with various display and filtering options",
		Long: `List all archives tracked in the registry with rich display modes and filtering.

The list command provides multiple ways to view and filter your archives:
- Display modes: table, compact, card, tree, dashboard
- Filters: by size, age, status, profile, location
- Output formats: human-readable or machine-readable (JSON/CSV)

Display mode is auto-detected based on terminal width if not specified.`,
		Example: `  # List all archives with auto-detected display
  7zarch-go list

  # Use specific display modes
  7zarch-go list --table            # High-density table view
  7zarch-go list --dashboard        # Management overview
  7zarch-go list --card             # Detailed cards for each archive

  # Filter archives
  7zarch-go list --missing          # Only missing archives
  7zarch-go list --managed          # Only managed archives
  7zarch-go list --older-than 30d   # Archives older than 30 days
  7zarch-go list --larger-than 100M # Archives larger than 100MB
  
  # Machine-readable output
  7zarch-go list --output json      # JSON format for scripting
  7zarch-go list --output csv       # CSV format for spreadsheets`,
		RunE:  runList,
	}

	// Add flags
	cmd.Flags().String("directory", "", "List archives in specific directory instead of managed storage")
	cmd.Flags().Bool("details", false, "Show detailed information")
	cmd.Flags().Bool("not-uploaded", false, "Show only archives that haven't been uploaded")
	cmd.Flags().String("pattern", "", "Filter archives by name pattern")
	cmd.Flags().String("older-than", "", "Show archives older than duration (e.g., '7d', '1h')")
	cmd.Flags().Bool("managed", false, "Only managed archives")
	cmd.Flags().Bool("external", false, "Only external archives")
	cmd.Flags().Bool("missing", false, "Only missing archives")
	cmd.Flags().String("status", "", "Filter by status (present|missing|deleted)")
	cmd.Flags().String("profile", "", "Filter by profile (media|documents|balanced)")
	cmd.Flags().Int64("larger-than", 0, "Filter by size larger than bytes (e.g., 1048576)")
	cmd.Flags().Bool("deleted", false, "Show only deleted archives")
	cmd.Flags().String("output", "", "Output format: table|json|csv|yaml (default: table)")
	
	// Query integration flags
	cmd.Flags().String("save-query", "", "Save current filters as a named query")
	cmd.Flags().String("query", "", "Use a saved query instead of specifying filters")

	// Display mode flags
	cmd.Flags().Bool("table", false, "Use table display mode (enhanced)")
	cmd.Flags().Bool("compact", false, "Use compact display mode")
	cmd.Flags().Bool("card", false, "Use card display mode")
	cmd.Flags().Bool("tree", false, "Use tree display mode")
	cmd.Flags().Bool("dashboard", false, "Use dashboard display mode")
	
	// Debug flag
	cmd.Flags().Bool("debug", false, "Show performance and debug information")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	directory, _ := cmd.Flags().GetString("directory")
	
	// Check for query integration flags
	queryName := getString(cmd, "query")
	saveQueryName := getString(cmd, "save-query")
	
	// If using a saved query, load it instead of collecting individual flags
	if queryName != "" {
		return runListWithSavedQuery(cmd, queryName, saveQueryName)
	}
	
	opts := listFilters{
		details:      getBool(cmd, "details"),
		notUploaded:  getBool(cmd, "not-uploaded"),
		pattern:      getString(cmd, "pattern"),
		olderThan:    getString(cmd, "older-than"),
		onlyManaged:  getBool(cmd, "managed"),
		onlyExternal: getBool(cmd, "external"),
		onlyMissing:  getBool(cmd, "missing"),
		onlyDeleted:  getBool(cmd, "deleted"),
		status:       getString(cmd, "status"),
		profile:      getString(cmd, "profile"),
		largerThan:   getInt64(cmd, "larger-than"),
		debug:        getBool(cmd, "debug"),
	}
	
	// If save-query flag is set, save the current filters
	if saveQueryName != "" {
		if err := saveCurrentFiltersAsQuery(opts, saveQueryName); err != nil {
			return fmt.Errorf("failed to save query: %w", err)
		}
		fmt.Printf("✅ Query '%s' saved successfully\n", saveQueryName)
	}
	
	// Initialize metrics if debug mode
	var metrics *debug.Metrics
	if opts.debug {
		metrics = debug.NewMetrics()
	}

	// Check for output format
	outputFormat := getString(cmd, "output")
	if outputFormat != "" {
		return listRegistryArchivesWithOutput(opts, outputFormat)
	}

	// Determine display mode
	displayMode := determineDisplayMode(cmd)

	if directory != "" {
		// List archives in a specific directory
		return listDirectory(directory, opts.details, opts.pattern)
	}

	// List registry-tracked archives
	return listRegistryArchivesWithDisplay(opts, displayMode, metrics)
}

// flag helpers
func getBool(cmd *cobra.Command, name string) bool     { v, _ := cmd.Flags().GetBool(name); return v }
func getString(cmd *cobra.Command, name string) string { v, _ := cmd.Flags().GetString(name); return v }
func getInt64(cmd *cobra.Command, name string) int64   { v, _ := cmd.Flags().GetInt64(name); return v }

// determineDisplayMode selects the display mode based on flags
func determineDisplayMode(cmd *cobra.Command) display.Mode {
	if getBool(cmd, "table") {
		return display.ModeTable
	}
	if getBool(cmd, "compact") {
		return display.ModeCompact
	}
	if getBool(cmd, "card") {
		return display.ModeCard
	}
	if getBool(cmd, "tree") {
		return display.ModeTree
	}
	if getBool(cmd, "dashboard") {
		return display.ModeDashboard
	}
	// Default to auto-detection
	return display.ModeAuto
}

// listRegistryArchivesWithDisplay uses the new display system
func listRegistryArchivesWithDisplay(opts listFilters, mode display.Mode, metrics *debug.Metrics) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	// Initialize storage manager
	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize managed storage: %w", err)
	}
	defer storageManager.Close()

	// Get archives based on filters
	var archives []*storage.Archive
	if opts.notUploaded {
		archives, err = storageManager.ListNotUploaded()
	} else if opts.olderThan != "" {
		dur, parseErr := parseHumanDuration(opts.olderThan)
		if parseErr != nil {
			return fmt.Errorf("invalid duration format: %w", parseErr)
		}
		archives, err = storageManager.ListOlderThan(dur)
	} else {
		archives, err = storageManager.List()
	}
	if err != nil {
		return fmt.Errorf("failed to list archives: %w", err)
	}

	// Record query completion if metrics enabled
	if metrics != nil {
		metrics.RecordQueryTime()
		metrics.SetResultCount(len(archives))
		
		// Get database size if available
		if dbInfo, statErr := os.Stat(filepath.Join(cfg.Storage.ManagedPath, "registry.db")); statErr == nil {
			metrics.SetDatabaseSize(dbInfo.Size())
		}
	}

	// Apply filters
	archives = applyAllFilters(archives, opts)

	// Use enhanced display system for supported modes
	if mode == display.ModeTable || mode == display.ModeCompact || mode == display.ModeCard || mode == display.ModeTree || mode == display.ModeDashboard {
		// Initialize display manager
		displayManager := display.NewManager()

		// Register available display modes
		tableDisplay := modes.NewTableDisplay()
		compactDisplay := modes.NewCompactDisplay()
		cardDisplay := modes.NewCardDisplay()
		treeDisplay := modes.NewTreeDisplay()
		dashboardDisplay := modes.NewDashboardDisplay()
		displayManager.Register(display.ModeTable, tableDisplay)
		displayManager.Register(display.ModeCompact, compactDisplay)
		displayManager.Register(display.ModeCard, cardDisplay)
		displayManager.Register(display.ModeTree, treeDisplay)
		displayManager.Register(display.ModeDashboard, dashboardDisplay)

		// Configure display options
		displayOpts := display.Options{
			Mode:        mode,
			Details:     opts.details,
			ShowHeaders: mode != display.ModeCompact, // No headers for compact by default
		}

		// Render using the display system
		err := displayManager.Render(archives, displayOpts)
		
		// Record render time and show debug output if enabled
		if metrics != nil {
			metrics.RecordRenderTime()
			if opts.debug {
				fmt.Printf("\n%s\n", metrics.String())
			}
		}
		
		return err
	}

	// Fall back to original display for now (other modes not yet implemented)
	err = displayArchivesOriginal(archives, opts)
	
	// Show debug output for fallback display too
	if metrics != nil && opts.debug {
		fmt.Printf("\n%s\n", metrics.String())
	}
	
	return err
}

// applyAllFilters applies all configured filters to the archive list
func applyAllFilters(archives []*storage.Archive, opts listFilters) []*storage.Archive {
	// Apply pattern filter
	if opts.pattern != "" {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if matched, _ := filepath.Match(opts.pattern, a.Name); matched {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}

	// Apply managed/external filter
	if opts.onlyManaged || opts.onlyExternal {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if opts.onlyManaged && a.Managed {
				filtered = append(filtered, a)
			}
			if opts.onlyExternal && !a.Managed {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}

	// Apply missing filter
	if opts.onlyMissing {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if a.Status == "missing" {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}

	// Apply deleted filter
	if opts.onlyDeleted {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if a.Status == "deleted" {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}

	// Apply status/profile/larger-than filters
	archives = applyFilters(archives, struct {
		status, profile string
		largerThan      int64
	}{opts.status, opts.profile, opts.largerThan})

	return archives
}

// displayArchivesOriginal is the original display function (fallback)
func displayArchivesOriginal(archives []*storage.Archive, opts listFilters) error {
	if len(archives) == 0 {
		fmt.Printf("No archives found.\n")
		fmt.Printf("💡 Tip: Create archives with '7zarch-go create <path>' to see them here.\n")
		return nil
	}

	// Group and summarize
	var managedCount, externalCount, missingCount, deletedCount int
	var activeManaged, activeExternal, deletedArchives []*storage.Archive

	for _, a := range archives {
		if a.Status == "deleted" {
			deletedCount++
			deletedArchives = append(deletedArchives, a)
		} else if a.Managed {
			managedCount++
			activeManaged = append(activeManaged, a)
		} else {
			externalCount++
			activeExternal = append(activeExternal, a)
		}
		if a.Status == "missing" {
			missingCount++
		}
	}

	fmt.Printf("📦 Archives (%d found)\n", len(archives))
	fmt.Printf("Active: %d (Managed: %d, External: %d) | Missing: %d | Deleted: %d\n\n",
		managedCount+externalCount, managedCount, externalCount, missingCount, deletedCount)

	// Delegate to existing printer
	return printGroupedArchives(archives, opts.details)
}

func listRegistryArchives(opts listFilters) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	// Initialize storage manager
	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize managed storage: %w", err)
	}
	defer storageManager.Close()

	// Get archives based on filters
	var archives []*storage.Archive
	if opts.notUploaded {
		archives, err = storageManager.ListNotUploaded()
	} else if opts.olderThan != "" {
		dur, parseErr := parseHumanDuration(opts.olderThan)
		if parseErr != nil {
			return fmt.Errorf("invalid duration format: %w", parseErr)
		}
		archives, err = storageManager.ListOlderThan(dur)
	} else {
		archives, err = storageManager.List()
	}
	if err != nil {
		return fmt.Errorf("failed to list archives: %w", err)
	}

	// Apply pattern filter
	if opts.pattern != "" {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if matched, _ := filepath.Match(opts.pattern, a.Name); matched {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}
	// Apply managed/external filter
	if opts.onlyManaged || opts.onlyExternal {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if opts.onlyManaged && a.Managed {
				filtered = append(filtered, a)
			}
			if opts.onlyExternal && !a.Managed {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}
	// Apply missing filter
	if opts.onlyMissing {
		filtered := make([]*storage.Archive, 0)
		for _, a := range archives {
			if a.Status == "missing" {
				filtered = append(filtered, a)
			}
		}
		archives = filtered
	}

	// Simple filter block (status/profile/larger-than)
	archives = applyFilters(archives, struct {
		status, profile string
		largerThan      int64
	}{opts.status, opts.profile, opts.largerThan})

	if len(archives) == 0 {
		fmt.Printf("No archives found.\n")
		fmt.Printf("💡 Tip: Create archives with '7zarch-go create <path>' to see them here.\n")
		return nil
	}

	// Group and summarize
	var managedCount, externalCount, missingCount, deletedCount int
	var activeManaged, activeExternal, deletedArchives []*storage.Archive

	for _, a := range archives {
		if a.Status == "deleted" {
			deletedCount++
			deletedArchives = append(deletedArchives, a)
		} else if a.Managed {
			managedCount++
			activeManaged = append(activeManaged, a)
		} else {
			externalCount++
			activeExternal = append(activeExternal, a)
		}
		if a.Status == "missing" {
			missingCount++
		}
	}

	fmt.Printf("📦 Archives (%d found)\n", len(archives))
	fmt.Printf("Active: %d (Managed: %d, External: %d) | Missing: %d | Deleted: %d\n\n",
		managedCount+externalCount, managedCount, externalCount, missingCount, deletedCount)

	// Delegate to existing printer to keep behavior identical
	return printGroupedArchives(archives, opts.details)
}

// applyFilters applies status/profile/largerThan filters in sequence
func applyFilters(archives []*storage.Archive, filters struct {
	status, profile string
	largerThan      int64
}) []*storage.Archive {
	result := archives
	if filters.status != "" {
		filtered := make([]*storage.Archive, 0)
		for _, a := range result {
			if a.Status == filters.status {
				filtered = append(filtered, a)
			}
		}
		result = filtered
	}
	if filters.profile != "" {
		filtered := make([]*storage.Archive, 0)
		for _, a := range result {
			if a.Profile == filters.profile {
				filtered = append(filtered, a)
			}
		}
		result = filtered
	}
	if filters.largerThan > 0 {
		filtered := make([]*storage.Archive, 0)
		for _, a := range result {
			if a.Size > filters.largerThan {
				filtered = append(filtered, a)
			}
		}
		result = filtered
	}
	return result
}

// printGroupedArchives prints groups and summary (same behavior as before)
func printGroupedArchives(archives []*storage.Archive, details bool) error {
	// Group and summarize
	var managedCount, externalCount, missingCount, deletedCount int
	var activeManaged, activeExternal, deletedArchives []*storage.Archive

	for _, a := range archives {
		if a.Status == "deleted" {
			deletedCount++
			deletedArchives = append(deletedArchives, a)
		} else if a.Managed {
			managedCount++
			activeManaged = append(activeManaged, a)
		} else {
			externalCount++
			activeExternal = append(activeExternal, a)
		}
		if a.Status == "missing" {
			missingCount++
		}
	}

	fmt.Printf("📦 Archives (%d found)\n", len(archives))
	fmt.Printf("Active: %d (Managed: %d, External: %d) | Missing: %d | Deleted: %d\n\n",
		managedCount+externalCount, managedCount, externalCount, missingCount, deletedCount)

	// Print active archives
	if len(activeManaged) > 0 {
		fmt.Printf("ACTIVE - MANAGED\n")
		printArchiveTable(activeManaged, details)
	}
	if len(activeExternal) > 0 {
		fmt.Printf("ACTIVE - EXTERNAL\n")
		printArchiveTable(activeExternal, details)
	}

	// Print deleted archives
	if len(deletedArchives) > 0 {
		// Load config to get actual retention days
		cfg, _ := config.Load()
		retentionDays := 7 // default fallback
		if cfg != nil && cfg.Storage.RetentionDays > 0 {
			retentionDays = cfg.Storage.RetentionDays

		}
		fmt.Printf("DELETED (auto-purge older than %d days)\n", retentionDays)
		for _, a := range deletedArchives {
			displayDeletedArchive(a, details)
		}
	}

	return nil
}

func printArchiveTable(archives []*storage.Archive, details bool) {
	// Headers
	if details {
		fmt.Printf("%-12s  %-30s  %8s  %-10s  %-19s  %-7s\n", "ID", "Name", "Size", "Profile", "Created", "Status")
	} else {
		fmt.Printf("%-12s  %-30s  %8s  %-7s\n", "ID", "Name", "Size", "Status")
	}
	for _, a := range archives {
		id := a.UID
		if len(id) > 12 {
			id = id[:12]
		}
		sizeMB := fmt.Sprintf("%.1f MB", float64(a.Size)/(1024*1024))
		status := a.Status
		if status == "present" {
			status = "✓"
		} else if status == "missing" {
			status = "⚠️"
		}
		if details {
			created := a.Created.Format("2006-01-02 15:04:05")
			name := a.Name
			if len(name) > 30 {
				name = name[:29] + "…"
			}
			fmt.Printf("%-12s  %-30s  %8s  %-10s  %-19s  %-7s\n", id, name, sizeMB, a.Profile, created, status)
		} else {
			name := a.Name
			if len(name) > 30 {
				name = name[:29] + "…"
			}
			fmt.Printf("%-12s  %-30s  %8s  %-7s\n", id, name, sizeMB, status)
		}
	}
	fmt.Println()
}

func listManagedArchives(details, notUploaded bool, pattern, olderThan string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize storage manager
	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize managed storage: %w", err)
	}
	defer storageManager.Close()

	// Get archives based on filters
	var archives []*storage.Archive
	if notUploaded {
		archives, err = storageManager.ListNotUploaded()
	} else if olderThan != "" {
		duration, parseErr := parseHumanDuration(olderThan)
		if parseErr != nil {
			return fmt.Errorf("invalid duration format: %w", parseErr)
		}
		archives, err = storageManager.ListOlderThan(duration)
	} else {
		archives, err = storageManager.List()
	}

	if err != nil {
		return fmt.Errorf("failed to list archives: %w", err)
	}

	// Apply pattern filter if specified
	if pattern != "" {
		filtered := make([]*storage.Archive, 0)
		for _, archive := range archives {
			if matched, _ := filepath.Match(pattern, archive.Name); matched {
				filtered = append(filtered, archive)
			}
		}
		archives = filtered
	}

	if len(archives) == 0 {
		fmt.Printf("No archives found in managed storage.\n")
		fmt.Printf("💡 Tip: Create archives with '7zarch-go create <path>' to see them here.\n")
		return nil
	}

	// Display results
	fmt.Printf("📦 Managed Archives (%d found)\n", len(archives))
	fmt.Printf("Location: %s\n\n", storageManager.GetArchivesPath())

	for _, archive := range archives {
		displayArchive(archive, details)
	}

	return nil
}

func listDirectory(directory string, details bool, pattern string) error {
	fmt.Printf("📁 Listing .7z files in: %s\n\n", directory)

	// Find .7z files in directory
	matches, err := filepath.Glob(filepath.Join(directory, "*.7z"))
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	// Apply pattern filter if specified
	if pattern != "" {
		filtered := make([]string, 0)
		for _, match := range matches {
			if matched, _ := filepath.Match(pattern, filepath.Base(match)); matched {
				filtered = append(filtered, match)
			}
		}
		matches = filtered
	}

	if len(matches) == 0 {
		fmt.Printf("No .7z files found.\n")
		return nil
	}

	fmt.Printf("Found %d archive(s):\n\n", len(matches))

	for _, archivePath := range matches {
		info, err := os.Stat(archivePath)
		if err != nil {
			continue
		}

		fmt.Printf("📦 %s\n", filepath.Base(archivePath))
		if details {
			fmt.Printf("   Path: %s\n", archivePath)
			fmt.Printf("   Size: %.2f MB\n", float64(info.Size())/(1024*1024))
			fmt.Printf("   Created: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	}

	return nil
}

func displayArchive(archive *storage.Archive, details bool) {
	// Show upload status
	status := "📤 Not uploaded"
	if archive.Uploaded {
		status = "✅ Uploaded"
		if archive.Destination != "" {
			status += fmt.Sprintf(" to %s", archive.Destination)
		}
	}

	fmt.Printf("📦 %s - %s\n", archive.Name, status)

	if details {
		if archive.UID != "" {
			fmt.Printf("   ID: %s\n", archive.UID)
		}
		fmt.Printf("   Path: %s\n", archive.Path)
		fmt.Printf("   Size: %.2f MB\n", float64(archive.Size)/(1024*1024))
		fmt.Printf("   Created: %s\n", archive.Created.Format("2006-01-02 15:04:05"))
		if archive.Profile != "" {
			fmt.Printf("   Profile: %s\n", archive.Profile)
		}
		if archive.Uploaded && archive.UploadedAt != nil {
			fmt.Printf("   Uploaded: %s\n", archive.UploadedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("   Age: %s\n", formatDuration(archive.Age()))
	}
	fmt.Println()
}

func displayDeletedArchive(archive *storage.Archive, details bool) {
	// Show deleted status with trash emoji
	deleteTime := "unknown"
	if archive.DeletedAt != nil {
		deleteTime = archive.DeletedAt.Format("2006-01-02 15:04:05")
	}
	fmt.Printf("🗑️  %s - Deleted %s\n", archive.Name, deleteTime)

	if details {
		if archive.UID != "" {
			fmt.Printf("   ID: %s\n", archive.UID)
		}

		// Calculate days until auto-purge
		if archive.DeletedAt != nil {
			cfg, _ := config.Load()
			retentionDays := 7 // default fallback
			if cfg != nil && cfg.Storage.RetentionDays > 0 {
				retentionDays = cfg.Storage.RetentionDays
			}

			purgeDate := archive.DeletedAt.AddDate(0, 0, retentionDays)
			daysLeft := int(time.Until(purgeDate).Hours() / 24)

			if daysLeft > 1 {
				fmt.Printf("   Auto-purge: %d days (%s)\n", daysLeft, purgeDate.Format("2006-01-02"))
			} else if daysLeft == 1 {
				fmt.Printf("   Auto-purge: 1 day (%s)\n", purgeDate.Format("2006-01-02"))
			} else if daysLeft == 0 {
				fmt.Printf("   Auto-purge: today\n")
			} else {
				fmt.Printf("   Auto-purge: overdue by %d days\n", -daysLeft)
			}
		}

		if archive.OriginalPath != "" {
			fmt.Printf("   Original: %s\n", archive.OriginalPath)
		}
		fmt.Printf("   Trash: %s\n", archive.Path)
		fmt.Printf("   Size: %.2f MB\n", float64(archive.Size)/(1024*1024))
		if archive.Profile != "" {
			fmt.Printf("   Profile: %s\n", archive.Profile)
		}
	}
	fmt.Println()
}

func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// listRegistryArchivesWithOutput handles machine-readable output formats
func listRegistryArchivesWithOutput(opts listFilters, format string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize storage manager
	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize managed storage: %w", err)
	}
	defer storageManager.Close()

	// Get archives based on filters
	var archives []*storage.Archive
	if opts.notUploaded {
		archives, err = storageManager.ListNotUploaded()
	} else if opts.olderThan != "" {
		dur, parseErr := parseHumanDuration(opts.olderThan)
		if parseErr != nil {
			return fmt.Errorf("invalid duration format: %w", parseErr)
		}
		archives, err = storageManager.ListOlderThan(dur)
	} else {
		archives, err = storageManager.List()
	}
	if err != nil {
		return fmt.Errorf("failed to list archives: %w", err)
	}

	// Apply filters
	archives = applyAllFilters(archives, opts)

	// Output in requested format
	switch format {
	case "json":
		return outputJSON(archives)
	case "csv":
		return outputCSV(archives)
	case "yaml":
		return outputYAML(archives)
	case "table":
		// Fall back to display system
		return listRegistryArchivesWithDisplay(opts, display.ModeTable, nil)
	default:
		return fmt.Errorf("unsupported output format: %s (supported: json, csv, yaml, table)", format)
	}
}

func outputJSON(archives []*storage.Archive) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(archives)
}

func outputCSV(archives []*storage.Archive) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{
		"uid", "name", "path", "size", "created", "checksum", "profile",
		"managed", "status", "last_seen", "deleted_at", "original_path",
		"uploaded", "destination", "uploaded_at",
	}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, a := range archives {
		var lastSeen, deletedAt, uploadedAt string
		if a.LastSeen != nil {
			lastSeen = a.LastSeen.Format(time.RFC3339)
		}
		if a.DeletedAt != nil {
			deletedAt = a.DeletedAt.Format(time.RFC3339)
		}
		if a.UploadedAt != nil {
			uploadedAt = a.UploadedAt.Format(time.RFC3339)
		}

		row := []string{
			a.UID,
			a.Name,
			a.Path,
			fmt.Sprintf("%d", a.Size),
			a.Created.Format(time.RFC3339),
			a.Checksum,
			a.Profile,
			fmt.Sprintf("%t", a.Managed),
			a.Status,
			lastSeen,
			deletedAt,
			a.OriginalPath,
			fmt.Sprintf("%t", a.Uploaded),
			a.Destination,
			uploadedAt,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func outputYAML(archives []*storage.Archive) error {
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	return enc.Encode(archives)
}

// runListWithSavedQuery executes list using a saved query
func runListWithSavedQuery(cmd *cobra.Command, queryName, saveQueryName string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize storage manager
	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage manager: %w", err)
	}
	defer storageManager.Close()

	// Initialize query manager
	resolver := storage.NewResolver(storageManager.Registry())
	queryManager := query.NewQueryManager(storageManager.Registry().DB(), resolver)

	// Run the saved query
	archives, err := queryManager.Run(queryName)
	if err != nil {
		return fmt.Errorf("failed to run query '%s': %w", queryName, err)
	}

	// If save-query flag is also set, this doesn't make sense with --query, so warn
	if saveQueryName != "" {
		fmt.Printf("⚠️  Warning: --save-query ignored when using --query\n")
	}

	// Check for output format first
	outputFormat := getString(cmd, "output")
	if outputFormat != "" {
		switch outputFormat {
		case "json":
			return outputJSON(archives)
		case "csv":
			return outputCSV(archives)
		case "yaml":
			return outputYAML(archives)
		default:
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}
	}

	// Display results using same logic as normal list
	if len(archives) == 0 {
		fmt.Printf("📋 Query '%s' - No archives found\n", queryName)
		return nil
	}

	fmt.Printf("📋 Query '%s' - %d archives found\n\n", queryName, len(archives))

	// Use simplified display for now
	return printGroupedArchives(archives, getBool(cmd, "details"))
}

// saveCurrentFiltersAsQuery converts listFilters to a map and saves as query
func saveCurrentFiltersAsQuery(opts listFilters, queryName string) error {
	// Convert listFilters to map format expected by query system
	filters := make(map[string]string)

	if opts.notUploaded {
		filters["not-uploaded"] = "true"
	}
	if opts.pattern != "" {
		filters["pattern"] = opts.pattern
	}
	if opts.olderThan != "" {
		filters["older-than"] = opts.olderThan
	}
	if opts.onlyManaged {
		filters["managed"] = "true"
	}
	if opts.onlyExternal {
		filters["external"] = "true"
	}
	if opts.onlyMissing {
		filters["missing"] = "true"
	}
	if opts.onlyDeleted {
		filters["deleted"] = "true"
	}
	if opts.status != "" {
		filters["status"] = opts.status
	}
	if opts.profile != "" {
		filters["profile"] = opts.profile
	}
	if opts.largerThan > 0 {
		filters["larger-than"] = fmt.Sprintf("%d", opts.largerThan)
	}

	if len(filters) == 0 {
		return fmt.Errorf("no filters specified - cannot save empty query")
	}

	// Load configuration and initialize query manager
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage manager: %w", err)
	}
	defer storageManager.Close()

	resolver := storage.NewResolver(storageManager.Registry())
	queryManager := query.NewQueryManager(storageManager.Registry().DB(), resolver)

	// Save the query
	return queryManager.Save(queryName, filters)
}
