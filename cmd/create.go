package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/archive"
	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	compressionLevel int
	threads         int
	verbose         bool
	dryRun          bool
	outputPath      string
	comprehensive   bool
	createLog       bool
	createChecksums bool
	forceOverwrite  bool
	profileName     string
	presetName      string
	noManaged       bool
)

func CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <path>",
		Short: "Create a new archive",
		Long:  `Create a new 7z archive from the specified path with optional compression settings.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runCreate,
	}

	// Add flags
	cmd.Flags().IntVarP(&compressionLevel, "compression", "c", 0, "Compression level (0-9, 0=smart default)")
	cmd.Flags().IntVarP(&threads, "threads", "t", 0, "Number of threads (0=auto)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without doing it")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for archive (default: managed storage)")
	cmd.Flags().BoolVar(&comprehensive, "comprehensive", false, "Create archive with log and checksums")
	cmd.Flags().BoolVar(&createLog, "log", false, "Create metadata log file")
	cmd.Flags().BoolVar(&createChecksums, "checksums", false, "Create SHA256 checksum file")
	cmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Overwrite existing archive")
	cmd.Flags().StringVar(&profileName, "profile", "", "Compression profile (media, documents, balanced)")
	cmd.Flags().StringVar(&presetName, "preset", "", "Use predefined settings preset")
	cmd.Flags().BoolVar(&noManaged, "no-managed", false, "Don't use managed storage (use current directory)")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("⚠️  Config loading failed, using defaults: %v\n", err)
		cfg = config.DefaultConfig()
	}
	
	// Apply preset if specified
	if presetName != "" {
		preset, exists := cfg.Presets[presetName]
		if !exists {
			return fmt.Errorf("unknown preset: %s", presetName)
		}
		
		// Apply preset values (CLI flags override presets)
		if profileName == "" && preset.Profile != "" {
			profileName = preset.Profile
		}
		if !comprehensive && preset.Comprehensive {
			comprehensive = preset.Comprehensive
		}
		if !forceOverwrite && preset.Force {
			forceOverwrite = preset.Force
		}
		if outputPath == "" && preset.Output != "" {
			// Expand tilde in output path
			if strings.HasPrefix(preset.Output, "~/") {
				home, _ := os.UserHomeDir()
				outputPath = filepath.Join(home, preset.Output[2:])
			} else {
				outputPath = preset.Output
			}
		}
		if threads == 0 && preset.Threads > 0 {
			threads = preset.Threads
		}
		
		fmt.Printf("📋 Using preset: %s\n", presetName)
	}
	
	// Apply config defaults (CLI flags and presets override config)
	if !comprehensive && cfg.Defaults.Create.Comprehensive {
		comprehensive = cfg.Defaults.Create.Comprehensive
	}
	if !forceOverwrite && cfg.Defaults.Create.Force {
		forceOverwrite = cfg.Defaults.Create.Force
	}
	if threads == 0 && cfg.Defaults.Create.Threads > 0 {
		threads = cfg.Defaults.Create.Threads
	}
	
	// Resolve absolute path
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if source exists
	_, err = os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("source path does not exist: %w", err)
	}

	// Initialize storage manager if using managed storage
	var storageManager *storage.Manager
	var useManaged bool
	
	// Determine if we should use managed storage
	if outputPath == "" && !noManaged && cfg.Storage.UseManagedDefault {
		// Use managed storage
		useManaged = true
		storageManager, err = storage.NewManager(cfg.Storage.ManagedPath)
		if err != nil {
			return fmt.Errorf("failed to initialize managed storage: %w", err)
		}
		defer storageManager.Close()
	}

	// Determine archive name and path
	var archiveName string
	baseName := filepath.Base(absPath) + ".7z"
	
	if outputPath != "" {
		// Explicit output path specified
		if filepath.Ext(outputPath) == ".7z" {
			archiveName = outputPath
		} else {
			// If directory specified, add filename
			archiveName = filepath.Join(outputPath, baseName)
		}
	} else if useManaged {
		// Use managed storage
		archiveName = storageManager.GetManagedPath(baseName)
	} else {
		// Default to current directory
		archiveName = baseName
	}

	// Enable log and checksums if comprehensive mode (handled in archive.Manager)
	// We keep flags for display only; artifact creation is centralized in internal/archive
	if comprehensive {
		createLog = true
		createChecksums = true
	}

	// Check if archive already exists
	if _, err := os.Stat(archiveName); err == nil && !forceOverwrite {
		// File exists and force not specified
		fmt.Printf("❌ Archive already exists: %s\n", archiveName)
		fmt.Printf("\nOptions:\n")
		fmt.Printf("  • Use --force to overwrite\n")
		fmt.Printf("  • Use a different --output path\n")
		fmt.Printf("  • Delete the existing file first\n")
		return fmt.Errorf("archive already exists (use --force to overwrite)")
	}

	if dryRun {
		fmt.Printf("DRY RUN MODE - No files will be created\n\n")
		fmt.Printf("Would create archive: %s\n", archiveName)
		fmt.Printf("Source: %s\n", absPath)
		fmt.Printf("Compression level: %d\n", compressionLevel)
		if threads > 0 {
			fmt.Printf("Threads: %d\n", threads)
		} else {
			fmt.Printf("Threads: auto\n")
		}
		if createLog {
			fmt.Printf("Would create log: %s.log\n", archiveName)
		}
		if createChecksums {
			fmt.Printf("Would create checksum: %s.sha256\n", archiveName)
		}
		return nil
	}

	// Show meaningful start message (after profile is determined)
	fmt.Printf("Creating archive: %s\n", filepath.Base(archiveName))
	fmt.Printf("Source: %s\n", absPath)
	// Note: Compression level will be shown after profile determination
	if threads > 0 {
		fmt.Printf("Threads: %d\n", threads)
	} else {
		fmt.Printf("Threads: auto\n")
	}
	fmt.Printf("\n")

	// Create a spinner that shows we're working
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Compressing"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(200*time.Millisecond),
		progressbar.OptionSpinnerType(14),
	)

	// Create archive manager
	manager := archive.NewManager()
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Start the spinner
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				bar.Add(1)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Handle excludes from preset
	var excludes []string
	if presetName != "" {
		if preset, exists := cfg.Presets[presetName]; exists {
			excludes = append(excludes, preset.Exclude...)
		}
	}

	// Create the archive
	opts := archive.CreateOptions{
		Source:           absPath,
		Output:           archiveName,
		CompressionLevel: compressionLevel,
		Threads:          threads,
		Profile:          profileName,
		Comprehensive:    comprehensive,
		Force:            forceOverwrite,
		Exclude:          excludes,
	}

	startTime := time.Now()
	result, err := manager.Create(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	bar.Finish()
	duration := time.Since(startTime)

	// Artifact creation (log/checksum) is handled inside archive.Manager when --comprehensive is used.
	// If the user explicitly requested only one artifact without --comprehensive, we could support that here.
	// For now, we centralize to avoid duplication.

	// Register in managed storage if applicable
	if useManaged && storageManager != nil {
		if err := storageManager.Add(
			filepath.Base(result.Path),
			result.Path,
			result.Size,
			result.Profile.Name,
			result.Checksum,
			"",
		); err != nil {
			// Non-fatal error - archive was created successfully
			fmt.Printf("⚠️  Warning: Failed to register archive in managed storage: %v\n", err)
		}
	}

	// Print results
	fmt.Printf("\n✅ Archive created successfully!\n")
	if useManaged {
		fmt.Printf("📦 Stored in managed storage: %s\n", filepath.Base(result.Path))
	} else {
		fmt.Printf("Archive: %s\n", result.Path)
	}
	fmt.Printf("Size: %.2f MB\n", float64(result.Size)/(1024*1024))
	fmt.Printf("Files: %d\n", result.FileCount)
	fmt.Printf("Compression: Level %d (%s profile)\n", result.Profile.Level, result.Profile.Name)
	fmt.Printf("Duration: %s\n", duration.Round(time.Second))
	
	if result.Size > 0 && result.OriginalSize > 0 {
		ratio := float64(result.Size) / float64(result.OriginalSize) * 100
		fmt.Printf("Size reduction: %.1f%%\n", 100-ratio)
	}
	
	if useManaged {
		fmt.Printf("\n💡 Tip: Use '7zarch-go list' to see all managed archives\n")
	}

	return nil
}