package modes

import (
	"fmt"
	"time"

	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// CompactDisplay provides terminal-friendly minimal output
type CompactDisplay struct{}

// NewCompactDisplay creates a new compact display mode
func NewCompactDisplay() *CompactDisplay {
	return &CompactDisplay{}
}

// Name returns the display mode name
func (cd *CompactDisplay) Name() string {
	return "compact"
}

// MinWidth returns the minimum terminal width for this display
func (cd *CompactDisplay) MinWidth() int {
	return 60
}

// Render displays archives in compact format
func (cd *CompactDisplay) Render(archives []*storage.Archive, opts display.Options) error {
	if len(archives) == 0 {
		if !opts.ShowHeaders {
			// Script-friendly: just return silently
			return nil
		}
		fmt.Printf("No archives found\n")
		return nil
	}

	// Group archives by status for summary
	var activeCount, missingCount, deletedCount int
	for _, a := range archives {
		switch a.Status {
		case "deleted":
			deletedCount++
		case "missing":
			missingCount++
		default:
			activeCount++
		}
	}

	// Print compact summary (unless headers disabled)
	if opts.ShowHeaders {
		fmt.Printf("%d archives", len(archives))
		if missingCount > 0 || deletedCount > 0 {
			fmt.Printf(" (%d active", activeCount)
			if missingCount > 0 {
				fmt.Printf(", %d missing", missingCount)
			}
			if deletedCount > 0 {
				fmt.Printf(", %d deleted", deletedCount)
			}
			fmt.Printf(")")
		}
		fmt.Println()
	}

	// Print each archive in compact format
	for _, archive := range archives {
		cd.printCompactArchive(archive, opts)
	}

	return nil
}

// printCompactArchive prints a single archive in compact format
func (cd *CompactDisplay) printCompactArchive(archive *storage.Archive, opts display.Options) {
	// Format: ID  name  size  age  status
	
	// Truncate ID to 12 chars for consistency with show command
	id := archive.UID
	if len(id) > 12 {
		id = id[:12]
	}

	// Truncate name for compact display
	name := archive.Name
	if len(name) > 25 {
		name = name[:24] + "…"
	}

	// Format size compactly
	size := display.FormatSize(archive.Size)
	
	// Format age compactly
	age := formatCompactAge(archive.Created)
	
	// Format status compactly
	status := cd.formatCompactStatus(archive)

	if opts.Details {
		// Detailed compact: add profile and managed indicator
		profile := archive.Profile
		if profile == "" {
			profile = "-"
		}
		if len(profile) > 8 {
			profile = profile[:7] + "…"
		}

		location := "EXTERNAL"
		if archive.Managed {
			location = "MANAGED"
		}

		fmt.Printf("%-12s  %-25s  %8s  %-8s  %3s  %-8s  %s\n", 
			id, name, size, profile, age, location, status)
	} else {
		// Basic compact format
		fmt.Printf("%-12s  %-25s  %8s  %3s  %s\n", 
			id, name, size, age, status)
	}
}

// formatCompactStatus returns a compact status indicator
func (cd *CompactDisplay) formatCompactStatus(archive *storage.Archive) string {
	return display.FormatStatus(archive.Status, false) // Use text format for compact
}

// formatCompactAge formats duration in a very compact way
func formatCompactAge(created time.Time) string {
	age := time.Since(created)
	
	if age < time.Minute {
		return "now"
	}
	if age < time.Hour {
		return fmt.Sprintf("%dm", int(age.Minutes()))
	}
	if age < 24*time.Hour {
		return fmt.Sprintf("%dh", int(age.Hours()))
	}
	if age < 7*24*time.Hour {
		return fmt.Sprintf("%dd", int(age.Hours()/24))
	}
	if age < 30*24*time.Hour {
		return fmt.Sprintf("%dw", int(age.Hours()/(24*7)))
	}
	if age < 365*24*time.Hour {
		return fmt.Sprintf("%dmo", int(age.Hours()/(24*30)))
	}
	return fmt.Sprintf("%dy", int(age.Hours()/(24*365)))
}