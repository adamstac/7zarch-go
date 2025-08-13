package modes

import (
	"fmt"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// CardDisplay provides rich information display with visual hierarchy
type CardDisplay struct{}

// NewCardDisplay creates a new card display mode
func NewCardDisplay() *CardDisplay {
	return &CardDisplay{}
}

// Name returns the display mode name
func (cd *CardDisplay) Name() string {
	return "card"
}

// MinWidth returns the minimum terminal width for this display
func (cd *CardDisplay) MinWidth() int {
	return 80
}

// Render displays archives in card format
func (cd *CardDisplay) Render(archives []*storage.Archive, opts display.Options) error {
	if len(archives) == 0 {
		fmt.Printf("No archives found\n")
		fmt.Printf("Create archives with '7zarch-go create <path>' to see them here.\n")
		return nil
	}

	// Group archives by status and location
	var managedActive, externalActive, deleted []*storage.Archive
	var missingCount int

	for _, a := range archives {
		if a.Status == "deleted" {
			deleted = append(deleted, a)
		} else if a.Managed {
			managedActive = append(managedActive, a)
		} else {
			externalActive = append(externalActive, a)
		}
		if a.Status == "missing" {
			missingCount++
		}
	}

	// Print summary header
	fmt.Printf("Archive Collection (%d archives found)\n", len(archives))
	fmt.Printf("Active: %d (Managed: %d, External: %d) | Missing: %d | Deleted: %d\n\n",
		len(managedActive)+len(externalActive), len(managedActive), len(externalActive), missingCount, len(deleted))

	// Print managed archives
	if len(managedActive) > 0 {
		cd.printGroupHeader("MANAGED ARCHIVES", len(managedActive))
		for i, archive := range managedActive {
			cd.printCard(archive, opts)
			if i < len(managedActive)-1 {
				fmt.Println()
			}
		}
		fmt.Println()
	}

	// Print external archives
	if len(externalActive) > 0 {
		cd.printGroupHeader("EXTERNAL ARCHIVES", len(externalActive))
		for i, archive := range externalActive {
			cd.printCard(archive, opts)
			if i < len(externalActive)-1 {
				fmt.Println()
			}
		}
		fmt.Println()
	}

	// Print deleted archives
	if len(deleted) > 0 {
		cd.printGroupHeader("DELETED ARCHIVES", len(deleted))
		for i, archive := range deleted {
			cd.printCard(archive, opts)
			if i < len(deleted)-1 {
				fmt.Println()
			}
		}
		fmt.Println()
	}

	return nil
}

// printGroupHeader prints a section header
func (cd *CardDisplay) printGroupHeader(title string, count int) {
	fmt.Printf("%s\n", title)
}

// printCard prints a single archive as a card
func (cd *CardDisplay) printCard(archive *storage.Archive, opts display.Options) {
	// Card border
	cardWidth := 70
	if opts.Width > 0 && opts.Width < 80 {
		cardWidth = opts.Width - 10
	}

	// Top border
	fmt.Printf("┌%s┐\n", strings.Repeat("─", cardWidth-2))

	// Archive name and ID
	name := archive.Name
	id := archive.UID
	if len(id) > 12 {
		id = id[:12]
	}
	
	// Calculate the full content: "name [id]"
	fullContent := fmt.Sprintf("%s [%s]", name, id)
	maxContentLen := cardWidth - 4 // 4 for "│ " + " │"
	
	// Truncate name if the full content is too long
	if len(fullContent) > maxContentLen {
		availableNameLen := maxContentLen - len(id) - 3 // 3 for " [" + "]"
		if availableNameLen > 3 {
			name = name[:availableNameLen-3] + "..."
		}
	}
	
	// Recalculate with potentially truncated name
	nameAndId := fmt.Sprintf("%s [%s]", name, id)
	titlePadding := cardWidth - 4 - len(nameAndId) // 4 for "│ " + " │"
	if titlePadding < 0 {
		titlePadding = 0
	}
	
	fmt.Printf("│ %s%s │\n", nameAndId, strings.Repeat(" ", titlePadding))

	// Separator
	fmt.Printf("│%s│\n", strings.Repeat("─", cardWidth-2))

	// Status and location
	status := cd.formatCardStatus(archive)
	location := "External Storage"
	if archive.Managed {
		location = "Managed Storage"
	}
	
	// Calculate padding for this row
	statusPart := fmt.Sprintf("Status: %s", status)
	locationPart := fmt.Sprintf("Location: %s", location)
	contentLen1 := len(statusPart) + len(locationPart) + 4 // 4 for "│ " + "    " + " │"
	padding1 := cardWidth - contentLen1
	if padding1 < 0 {
		padding1 = 1
	}
	
	fmt.Printf("│ %s%s%s │\n", statusPart, strings.Repeat(" ", padding1), locationPart)

	// Size and profile
	size := display.FormatSize(archive.Size)
	profile := archive.Profile
	if profile == "" {
		profile = "default"
	}
	
	// Calculate padding for this row
	sizePart := fmt.Sprintf("Size: %s", size)
	profilePart := fmt.Sprintf("Profile: %s", profile)
	contentLen2 := len(sizePart) + len(profilePart) + 4 // 4 for "│ " + "    " + " │"
	padding2 := cardWidth - contentLen2
	if padding2 < 0 {
		padding2 = 1
	}
	
	fmt.Printf("│ %s%s%s │\n", sizePart, strings.Repeat(" ", padding2), profilePart)

	// Created and age
	created := archive.Created.Format("2006-01-02 15:04:05")
	age := cd.formatCardAge(archive.Created)
	
	// Calculate padding for this row
	createdPart := fmt.Sprintf("Created: %s", created)
	agePart := fmt.Sprintf("Age: %s", age)
	contentLen3 := len(createdPart) + len(agePart) + 4 // 4 for "│ " + "    " + " │"
	padding3 := cardWidth - contentLen3
	if padding3 < 0 {
		padding3 = 1
	}
	
	fmt.Printf("│ %s%s%s │\n", createdPart, strings.Repeat(" ", padding3), agePart)

	// Path
	path := archive.Path
	pathPart := fmt.Sprintf("Path: %s", path)
	
	// Truncate path if too long
	maxPathPartLen := cardWidth - 4 // 4 for "│ " + " │"
	if len(pathPart) > maxPathPartLen {
		availablePathLen := maxPathPartLen - 6 // 6 for "Path: "
		if availablePathLen > 10 {
			truncatedPath := "..." + path[len(path)-(availablePathLen-3):]
			pathPart = fmt.Sprintf("Path: %s", truncatedPath)
		}
	}
	
	// Calculate padding for path row
	padding4 := cardWidth - 4 - len(pathPart) // 4 for "│ " + " │"
	if padding4 < 0 {
		padding4 = 0
	}
	
	fmt.Printf("│ %s%s │\n", pathPart, strings.Repeat(" ", padding4))

	// Additional details if requested
	if opts.Details {
		fmt.Printf("│%s│\n", strings.Repeat("─", cardWidth-2))
		
		// Checksum (if available)
		if archive.Checksum != "" {
			checksum := archive.Checksum
			checksumPart := fmt.Sprintf("Checksum: %s", checksum)
			
			// Truncate checksum if too long
			maxChecksumPartLen := cardWidth - 4 // 4 for "│ " + " │"
			if len(checksumPart) > maxChecksumPartLen {
				availableChecksumLen := maxChecksumPartLen - 10 // 10 for "Checksum: "
				if availableChecksumLen > 10 {
					truncatedChecksum := checksum[:availableChecksumLen-3] + "..."
					checksumPart = fmt.Sprintf("Checksum: %s", truncatedChecksum)
				}
			}
			
			// Calculate padding for checksum row
			padding5 := cardWidth - 4 - len(checksumPart) // 4 for "│ " + " │"
			if padding5 < 0 {
				padding5 = 0
			}
			
			fmt.Printf("│ %s%s │\n", checksumPart, strings.Repeat(" ", padding5))
		}

		// Additional metadata for deleted archives
		if archive.Status == "deleted" && archive.DeletedAt != nil {
			deletedTime := archive.DeletedAt.Format("2006-01-02 15:04:05")
			deletedPart := fmt.Sprintf("Deleted: %s", deletedTime)
			
			// Calculate padding for deleted row
			padding6 := cardWidth - 4 - len(deletedPart) // 4 for "│ " + " │"
			if padding6 < 0 {
				padding6 = 0
			}
			
			fmt.Printf("│ %s%s │\n", deletedPart, strings.Repeat(" ", padding6))
			
			// Days until auto-purge (assuming 7 day retention)
			purgeDate := archive.DeletedAt.AddDate(0, 0, 7)
			daysLeft := int(time.Until(purgeDate).Hours() / 24)
			var purgePart string
			if daysLeft > 0 {
				purgePart = fmt.Sprintf("Auto-purge in: %d days", daysLeft)
			} else {
				purgePart = "Auto-purge: overdue"
			}
			
			// Calculate padding for purge row
			padding7 := cardWidth - 4 - len(purgePart) // 4 for "│ " + " │"
			if padding7 < 0 {
				padding7 = 0
			}
			
			fmt.Printf("│ %s%s │\n", purgePart, strings.Repeat(" ", padding7))
		}

		// Original path for deleted archives
		if archive.OriginalPath != "" && archive.OriginalPath != archive.Path {
			origPath := archive.OriginalPath
			originalPart := fmt.Sprintf("Original: %s", origPath)
			
			// Truncate original path if too long
			maxOriginalPartLen := cardWidth - 4 // 4 for "│ " + " │"
			if len(originalPart) > maxOriginalPartLen {
				availableOrigPathLen := maxOriginalPartLen - 10 // 10 for "Original: "
				if availableOrigPathLen > 10 {
					truncatedOrigPath := "..." + origPath[len(origPath)-(availableOrigPathLen-3):]
					originalPart = fmt.Sprintf("Original: %s", truncatedOrigPath)
				}
			}
			
			// Calculate padding for original path row
			padding8 := cardWidth - 4 - len(originalPart) // 4 for "│ " + " │"
			if padding8 < 0 {
				padding8 = 0
			}
			
			fmt.Printf("│ %s%s │\n", originalPart, strings.Repeat(" ", padding8))
		}
	}

	// Bottom border
	fmt.Printf("└%s┘\n", strings.Repeat("─", cardWidth-2))
}

// formatCardStatus returns a formatted status for cards
func (cd *CardDisplay) formatCardStatus(archive *storage.Archive) string {
	icon := display.FormatStatus(archive.Status, true)
	switch archive.Status {
	case "present":
		return fmt.Sprintf("%s Present", icon)
	case "missing":
		return fmt.Sprintf("%s Missing", icon)
	case "deleted":
		return fmt.Sprintf("%s Deleted", icon)
	default:
		return fmt.Sprintf("%s %s", icon, archive.Status)
	}
}

// formatCardAge formats duration since creation for cards
func (cd *CardDisplay) formatCardAge(created time.Time) string {
	age := time.Since(created)
	
	if age < time.Hour {
		mins := int(age.Minutes())
		return fmt.Sprintf("%d minutes ago", mins)
	}
	if age < 24*time.Hour {
		hours := int(age.Hours())
		return fmt.Sprintf("%d hours ago", hours)
	}
	if age < 7*24*time.Hour {
		days := int(age.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
	if age < 30*24*time.Hour {
		weeks := int(age.Hours() / (24 * 7))
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	if age < 365*24*time.Hour {
		months := int(age.Hours() / (24 * 30))
		return fmt.Sprintf("%d months ago", months)
	}
	years := int(age.Hours() / (24 * 365))
	return fmt.Sprintf("%d years ago", years)
}