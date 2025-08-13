package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/query"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func QueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Manage saved queries for archive filtering",
		Long: `Save, list, run, and delete queries that encapsulate filter combinations.

Saved queries allow you to store complex filter combinations and reuse them easily.
This is particularly useful for frequently used searches or complex filtering criteria.`,
		Example: `  # Save current list filters as a query
  7zarch-go query save "my-docs" --profile=documents --managed
  
  # List all saved queries
  7zarch-go query list
  
  # Run a saved query
  7zarch-go query run my-docs
  
  # Delete a saved query
  7zarch-go query delete my-docs`,
	}

	cmd.AddCommand(querySaveCmd())
	cmd.AddCommand(queryListCmd())
	cmd.AddCommand(queryRunCmd())
	cmd.AddCommand(queryDeleteCmd())
	cmd.AddCommand(queryShowCmd())

	return cmd
}

func querySaveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save <name> [filters...]",
		Short: "Save a new query with specified filters",
		Long: `Save a query with the specified name and filter combination.

The filters follow the same syntax as the list command flags.
This allows you to save complex filter combinations for reuse.`,
		Example: `  # Save a query for managed documents
  7zarch-go query save "my-docs" --profile=documents --managed
  
  # Save a query for large media files
  7zarch-go query save "big-media" --profile=media --larger-than=100000000
  
  # Save a query for old unuploaded archives
  7zarch-go query save "old-unuploaded" --not-uploaded --older-than=30d`,
		Args: cobra.MinimumNArgs(1),
		RunE: runQuerySave,
	}

	// Add the same filter flags as the list command
	cmd.Flags().Bool("not-uploaded", false, "Filter for archives that haven't been uploaded")
	cmd.Flags().String("pattern", "", "Filter archives by name pattern")
	cmd.Flags().String("older-than", "", "Filter archives older than duration (e.g., '7d', '1h')")
	cmd.Flags().Bool("managed", false, "Filter for managed archives only")
	cmd.Flags().Bool("external", false, "Filter for external archives only")
	cmd.Flags().Bool("missing", false, "Filter for missing archives only")
	cmd.Flags().String("status", "", "Filter by status (present|missing|deleted)")
	cmd.Flags().String("profile", "", "Filter by profile (media|documents|balanced)")
	cmd.Flags().Int64("larger-than", 0, "Filter by size larger than bytes")
	cmd.Flags().Bool("deleted", false, "Filter for deleted archives only")

	return cmd
}

func queryListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved queries",
		Long: `Display all saved queries with their usage statistics.

Shows query names, creation dates, last used dates, and usage counts.
Queries are sorted by most recently used first.`,
		Example: `  # List all saved queries
  7zarch-go query list
  
  # List with JSON output for scripting
  7zarch-go query list --output json`,
		RunE: runQueryList,
	}

	cmd.Flags().String("output", "", "Output format: table|json (default: table)")

	return cmd
}

func queryRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <name>",
		Short: "Execute a saved query",
		Long: `Run a saved query and display the matching archives.

This executes the stored filter combination and shows results using
the standard archive display format.`,
		Example: `  # Run a saved query
  7zarch-go query run my-docs
  
  # Run query with specific display mode
  7zarch-go query run my-docs --table
  
  # Run query with JSON output
  7zarch-go query run my-docs --output json`,
		Args: cobra.ExactArgs(1),
		RunE: runQueryRun,
	}

	// Display and output options
	cmd.Flags().Bool("table", false, "Use table display mode")
	cmd.Flags().Bool("compact", false, "Use compact display mode")
	cmd.Flags().Bool("card", false, "Use card display mode")
	cmd.Flags().Bool("tree", false, "Use tree display mode")
	cmd.Flags().Bool("dashboard", false, "Use dashboard display mode")
	cmd.Flags().String("output", "", "Output format: table|json|csv|yaml")
	cmd.Flags().Bool("details", false, "Show detailed information")

	return cmd
}

func queryDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a saved query",
		Long: `Remove a saved query permanently.

This action cannot be undone. The query will be completely removed
from the saved queries database.`,
		Example: `  # Delete a saved query
  7zarch-go query delete my-docs`,
		Args: cobra.ExactArgs(1),
		RunE: runQueryDelete,
	}

	return cmd
}

func queryShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show details of a saved query",
		Long: `Display detailed information about a specific saved query.

Shows the query name, filters, creation date, usage statistics,
and the actual filter configuration that will be applied.`,
		Example: `  # Show query details
  7zarch-go query show my-docs
  
  # Show with JSON output
  7zarch-go query show my-docs --output json`,
		Args: cobra.ExactArgs(1),
		RunE: runQueryShow,
	}

	cmd.Flags().String("output", "", "Output format: table|json (default: table)")

	return cmd
}

func runQuerySave(cmd *cobra.Command, args []string) error {
	queryName := args[0]

	// Collect filters from flags
	filters := make(map[string]string)

	if getBool(cmd, "not-uploaded") {
		filters["not-uploaded"] = "true"
	}
	if pattern := getString(cmd, "pattern"); pattern != "" {
		filters["pattern"] = pattern
	}
	if olderThan := getString(cmd, "older-than"); olderThan != "" {
		filters["older-than"] = olderThan
	}
	if getBool(cmd, "managed") {
		filters["managed"] = "true"
	}
	if getBool(cmd, "external") {
		filters["external"] = "true"
	}
	if getBool(cmd, "missing") {
		filters["missing"] = "true"
	}
	if status := getString(cmd, "status"); status != "" {
		filters["status"] = status
	}
	if profile := getString(cmd, "profile"); profile != "" {
		filters["profile"] = profile
	}
	if largerThan := getInt64(cmd, "larger-than"); largerThan > 0 {
		filters["larger-than"] = fmt.Sprintf("%d", largerThan)
	}
	if getBool(cmd, "deleted") {
		filters["deleted"] = "true"
	}

	if len(filters) == 0 {
		return fmt.Errorf("no filters specified - provide at least one filter flag")
	}

	// Initialize storage and query manager
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
	if err := queryManager.Save(queryName, filters); err != nil {
		return fmt.Errorf("failed to save query: %w", err)
	}

	fmt.Printf("✅ Query '%s' saved successfully\n", queryName)
	fmt.Printf("Filters: %s\n", formatFilters(filters))

	return nil
}

func runQueryList(cmd *cobra.Command, args []string) error {
	// Initialize storage and query manager
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

	// Get all queries
	queries, err := queryManager.List()
	if err != nil {
		return fmt.Errorf("failed to list queries: %w", err)
	}

	// Check output format
	outputFormat := getString(cmd, "output")
	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(queries)
	}

	// Display in table format
	if len(queries) == 0 {
		fmt.Printf("No saved queries found.\n")
		fmt.Printf("💡 Tip: Save queries with '7zarch-go query save <name> [filters...]'\n")
		return nil
	}

	fmt.Printf("📋 Saved Queries (%d found)\n\n", len(queries))

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tCREATED\tLAST USED\tUSE COUNT\tFILTERS\n")

	for _, q := range queries {
		lastUsed := "never"
		if q.LastUsed != nil {
			lastUsed = q.LastUsed.Format("2006-01-02 15:04")
		}

		filtersStr := formatFilters(q.Filters)
		if len(filtersStr) > 50 {
			filtersStr = filtersStr[:47] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			q.Name,
			q.Created.Format("2006-01-02 15:04"),
			lastUsed,
			q.UseCount,
			filtersStr,
		)
	}

	return w.Flush()
}

func runQueryRun(cmd *cobra.Command, args []string) error {
	queryName := args[0]

	// Initialize storage and query manager
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

	// Execute the query
	archives, err := queryManager.Run(queryName)
	if err != nil {
		return fmt.Errorf("failed to run query: %w", err)
	}

	// Check for output format first
	outputFormat := getString(cmd, "output")
	if outputFormat != "" {
		switch outputFormat {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(archives)
		case "csv":
			return outputCSV(archives)
		case "yaml":
			return outputYAML(archives)
		default:
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}
	}

	// Display results using same logic as list command
	if len(archives) == 0 {
		fmt.Printf("📋 Query '%s' - No archives found\n", queryName)
		return nil
	}

	fmt.Printf("📋 Query '%s' - %d archives found\n\n", queryName, len(archives))

	// Use simplified display for now - can integrate with display system later
	printArchiveTable(archives, getBool(cmd, "details"))
	return nil
}

func runQueryDelete(cmd *cobra.Command, args []string) error {
	queryName := args[0]

	// Initialize storage and query manager
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

	// Delete the query
	if err := queryManager.Delete(queryName); err != nil {
		return fmt.Errorf("failed to delete query: %w", err)
	}

	fmt.Printf("✅ Query '%s' deleted successfully\n", queryName)

	return nil
}

func runQueryShow(cmd *cobra.Command, args []string) error {
	queryName := args[0]

	// Initialize storage and query manager
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

	// Get the query
	query, err := queryManager.Get(queryName)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	// Check output format
	outputFormat := getString(cmd, "output")
	if outputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(query)
	}

	// Display in table format
	fmt.Printf("📋 Query Details: %s\n\n", query.Name)
	fmt.Printf("Created: %s\n", query.Created.Format("2006-01-02 15:04:05"))
	
	if query.LastUsed != nil {
		fmt.Printf("Last Used: %s\n", query.LastUsed.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("Last Used: never\n")
	}
	
	fmt.Printf("Use Count: %d\n", query.UseCount)
	fmt.Printf("Filters: %s\n", formatFilters(query.Filters))

	return nil
}

// formatFilters creates a readable string representation of filter map
func formatFilters(filters map[string]string) string {
	if len(filters) == 0 {
		return "none"
	}

	var parts []string
	for key, value := range filters {
		if value == "true" {
			parts = append(parts, fmt.Sprintf("--%s", key))
		} else {
			parts = append(parts, fmt.Sprintf("--%s=%s", key, value))
		}
	}

	return strings.Join(parts, " ")
}