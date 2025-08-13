package modes

import (
	"fmt"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// TableDisplay provides high-density information scanning
type TableDisplay struct{}

// NewTableDisplay creates a new table display mode
func NewTableDisplay() *TableDisplay {
	return &TableDisplay{}
}

// Name returns the display mode name
func (td *TableDisplay) Name() string {
	return "table"
}

// MinWidth returns the minimum terminal width for this display
func (td *TableDisplay) MinWidth() int {
	return 80
}

// Render displays archives in an enhanced table format
func (td *TableDisplay) Render(archives []*storage.Archive, opts display.Options) error {
	if len(archives) == 0 {
		fmt.Printf("📦 No archives found\n")
		fmt.Printf("💡 Tip: Create archives with '7zarch-go create <path>' to see them here.\n")
		return nil
	}

	// Group archives by status
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
	td.printSummary(len(archives), len(managedActive), len(externalActive), missingCount, len(deleted))

	// Configure columns based on terminal width and options
	columns := td.selectColumns(opts)

	// Print active archives
	if len(managedActive) > 0 {
		fmt.Printf("\nACTIVE - MANAGED\n")
		td.printTable(managedActive, columns, opts)
	}

	if len(externalActive) > 0 {
		fmt.Printf("\nACTIVE - EXTERNAL\n")
		td.printTable(externalActive, columns, opts)
	}

	// Print deleted archives
	if len(deleted) > 0 {
		fmt.Printf("\nDELETED (auto-purge in 7 days)\n")
		td.printTable(deleted, columns, opts)
	}

	return nil
}

// printSummary prints the archive summary header
func (td *TableDisplay) printSummary(total, managed, external, missing, deleted int) {
	fmt.Printf("📦 Archives (%d found)\n", total)
	fmt.Printf("Active: %d (Managed: %d, External: %d) | Missing: %d | Deleted: %d\n",
		managed+external, managed, external, missing, deleted)
}

// Column represents a table column
type Column struct {
	Name   string
	Width  int
	Format func(*storage.Archive) string
}

// selectColumns determines which columns to show based on options
func (td *TableDisplay) selectColumns(opts display.Options) []Column {
	// Default columns for basic mode
	columns := []Column{
		{
			Name:  "ID",
			Width: 13,
			Format: func(a *storage.Archive) string {
				id := a.UID
				if len(id) > 12 {
					id = id[:12]
				}
				return id
			},
		},
		{
			Name:  "Name",
			Width: 30,
			Format: func(a *storage.Archive) string {
				return display.TruncateString(a.Name, 30)
			},
		},
		{
			Name:  "Size",
			Width: 9,
			Format: func(a *storage.Archive) string {
				return display.FormatSize(a.Size)
			},
		},
	}

	// Add detailed columns if requested
	if opts.Details {
		columns = append(columns,
			Column{
				Name:  "Profile",
				Width: 10,
				Format: func(a *storage.Archive) string {
					if a.Profile == "" {
						return "-"
					}
					return display.TruncateString(a.Profile, 10)
				},
			},
			Column{
				Name:  "Created",
				Width: 20,
				Format: func(a *storage.Archive) string {
					return a.Created.Format("2006-01-02 15:04:05")
				},
			},
			Column{
				Name:  "Age",
				Width: 7,
				Format: func(a *storage.Archive) string {
					return formatAge(a.Created)
				},
			},
		)
	}

	// Status column always last
	columns = append(columns, Column{
		Name:  "Status",
		Width: 7,
		Format: func(a *storage.Archive) string {
			return display.FormatStatus(a.Status, false) // Use text format for table
		},
	})

	// Adjust column widths for terminal width if needed
	if opts.Width > 0 {
		td.adjustColumnWidths(columns, opts.Width)
	}

	return columns
}

// adjustColumnWidths adapts column widths to fit terminal
func (td *TableDisplay) adjustColumnWidths(columns []Column, termWidth int) {
	// Calculate total width needed
	totalWidth := 0
	for _, col := range columns {
		totalWidth += col.Width + 2 // +2 for spacing
	}

	// If we're over terminal width, adjust the Name column
	if totalWidth > termWidth && len(columns) > 1 {
		excess := totalWidth - termWidth
		for i := range columns {
			if columns[i].Name == "Name" {
				columns[i].Width = max(15, columns[i].Width-excess)
				break
			}
		}
	}
}

// printTable prints archives in table format
func (td *TableDisplay) printTable(archives []*storage.Archive, columns []Column, opts display.Options) {
	// Print header with borders
	if opts.ShowHeaders {
		td.printBorder(columns, "top")
		td.printHeader(columns)
		td.printBorder(columns, "middle")
	}

	// Print rows
	for _, archive := range archives {
		td.printRow(archive, columns)
	}

	// Always print bottom border for clean look
	td.printBorder(columns, "bottom")

	fmt.Println()
}

// printBorder prints table borders
func (td *TableDisplay) printBorder(columns []Column, position string) {
	var left, middle, right, horizontal string

	switch position {
	case "top":
		left, middle, right, horizontal = "┌", "┬", "┐", "─"
	case "middle":
		left, middle, right, horizontal = "├", "┼", "┤", "─"
	case "bottom":
		left, middle, right, horizontal = "└", "┴", "┘", "─"
	}

	fmt.Print(left)
	for i, col := range columns {
		// Add 1 for left padding space
		fmt.Print(strings.Repeat(horizontal, col.Width+1))
		if i < len(columns)-1 {
			fmt.Print(middle)
		}
	}
	fmt.Println(right)
}

// printHeader prints the table header
func (td *TableDisplay) printHeader(columns []Column) {
	fmt.Print("│")
	for _, col := range columns {
		fmt.Printf(" %-*s│", col.Width, col.Name)
	}
	fmt.Println()
}

// printRow prints a single archive row
func (td *TableDisplay) printRow(archive *storage.Archive, columns []Column) {
	fmt.Print("│")
	for _, col := range columns {
		value := col.Format(archive)
		fmt.Printf(" %-*s│", col.Width, value)
	}
	fmt.Println()
}

// formatAge formats duration since creation
func formatAge(created time.Time) string {
	age := time.Since(created)

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

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
