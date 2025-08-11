package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List managed archives",
		Long:  `List all archives in managed storage.`,
		RunE:  runList,
	}

	// Add flags
	cmd.Flags().String("directory", "", "List archives in specific directory instead of managed storage")
	cmd.Flags().Bool("details", false, "Show detailed information")
	cmd.Flags().Bool("not-uploaded", false, "Show only archives that haven't been uploaded")
	cmd.Flags().String("pattern", "", "Filter archives by name pattern")
	cmd.Flags().String("older-than", "", "Show archives older than duration (e.g., '7d', '1h')")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	directory, _ := cmd.Flags().GetString("directory")
	details, _ := cmd.Flags().GetBool("details")
	notUploaded, _ := cmd.Flags().GetBool("not-uploaded")
	pattern, _ := cmd.Flags().GetString("pattern")
	olderThan, _ := cmd.Flags().GetString("older-than")

	if directory != "" {
		// List archives in a specific directory
		return listDirectory(directory, details, pattern)
	}

	// List managed archives
	return listManagedArchives(details, notUploaded, pattern, olderThan)
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
		duration, parseErr := time.ParseDuration(olderThan)
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

func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}