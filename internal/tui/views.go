package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

// renderListView renders the simple archive list
func (a *App) renderListView() string {
	if len(a.archives) == 0 {
		return a.renderEmptyState()
	}

	var lines []string
	
	// Header
	header := lipgloss.NewStyle().
		Foreground(a.theme.Header).
		Bold(true).
		Render("7zarch-go")
	lines = append(lines, header)
	lines = append(lines, "")

	// Summary line
	totalSize := int64(0)
	activeCount := 0
	missingCount := 0
	deletedCount := 0
	
	for _, archive := range a.archives {
		totalSize += archive.Size
		switch archive.Status {
		case "present":
			activeCount++
		case "missing":
			missingCount++
		case "deleted":
			deletedCount++
		}
	}
	
	summary := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render(fmt.Sprintf("Archives: %d (%s)", len(a.archives), humanize.Bytes(uint64(totalSize))))
	lines = append(lines, summary)
	lines = append(lines, "")

	// Archive list
	start := 0
	end := len(a.archives)
	if a.height > 0 {
		maxItems := a.height - 8 // Reserve space for header, summary, commands
		if maxItems > 0 && len(a.archives) > maxItems {
			start = a.cursor
			if start > len(a.archives)-maxItems {
				start = len(a.archives) - maxItems
			}
			end = start + maxItems
		}
	}

	for i := start; i < end && i < len(a.archives); i++ {
		archive := a.archives[i]
		line := a.renderArchiveLine(archive, i == a.cursor)
		lines = append(lines, line)
	}

	// Add some spacing
	lines = append(lines, "")
	
	// Commands
	commands := a.renderCommands()
	lines = append(lines, commands)

	return strings.Join(lines, "\n")
}

// renderArchiveLine renders a single archive line
func (a *App) renderArchiveLine(archive *storage.Archive, isSelected bool) string {
	// Selection indicator
	selector := "  "
	if isSelected {
		selector = "> "
	}
	
	// Multi-select checkbox
	checkbox := "[ ]"
	if a.selected[archive.UID] {
		checkbox = "[✓]"
	}
	
	// Status icon
	statusIcon := "✓"
	statusColor := a.theme.StatusOK
	switch archive.Status {
	case "missing":
		statusIcon = "?"
		statusColor = a.theme.StatusMiss
	case "deleted":
		statusIcon = "X"
		statusColor = a.theme.StatusDel
	}
	
	// Format components
	name := archive.Name
	if len(name) > 40 {
		name = name[:37] + "..."
	}
	
	size := humanize.Bytes(uint64(archive.Size))
	age := formatAge(archive.Created)
	
	// Style the line
	var style lipgloss.Style
	if isSelected {
		style = lipgloss.NewStyle().
			Background(a.theme.Selection).
			Foreground(a.theme.SelText).
			Width(a.width)
	} else {
		style = lipgloss.NewStyle().
			Foreground(a.theme.Foreground)
	}
	
	// Build line content
	nameStyle := style.Copy()
	metadataStyle := lipgloss.NewStyle().Foreground(a.theme.Metadata)
	statusStyle := lipgloss.NewStyle().Foreground(statusColor)
	
	if isSelected {
		metadataStyle = metadataStyle.Background(a.theme.Selection)
		statusStyle = statusStyle.Background(a.theme.Selection)
	}
	
	content := fmt.Sprintf("%s%s %-40s %10s %10s %s",
		selector,
		checkbox,
		name,
		metadataStyle.Render(size),
		metadataStyle.Render(age),
		statusStyle.Render(statusIcon))
	
	return style.Render(content)
}

// renderCommands renders the command help line
func (a *App) renderCommands() string {
	var commands []string
	
	if a.currentView == ListView {
		commands = []string{
			"[Enter] Details",
			"[Space] Select", 
			"[d] Delete",
			"[m] Move",
			"[u] Upload",
			"[/] Search",
			"[q] Quit",
		}
	} else {
		commands = []string{
			"[Enter] Back to list",
			"[d] Delete",
			"[m] Move", 
			"[u] Upload",
			"[q] Quit",
		}
	}
	
	return lipgloss.NewStyle().
		Foreground(a.theme.Commands).
		Render(strings.Join(commands, "  "))
}

// renderEmptyState renders when no archives found
func (a *App) renderEmptyState() string {
	var lines []string
	
	header := lipgloss.NewStyle().
		Foreground(a.theme.Header).
		Bold(true).
		Render("7zarch-go")
	lines = append(lines, header)
	lines = append(lines, "")
	
	message := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render("No archives found")
	lines = append(lines, message)
	lines = append(lines, "")
	
	commands := lipgloss.NewStyle().
		Foreground(a.theme.Commands).
		Render("[c] Create archive  [q] Quit")
	lines = append(lines, commands)
	
	return strings.Join(lines, "\n")
}

// renderDetailView renders the archive details view
func (a *App) renderDetailView() string {
	if len(a.archives) == 0 || a.cursor >= len(a.archives) {
		return a.renderEmptyState()
	}
	
	archive := a.archives[a.cursor]
	var lines []string
	
	// Header with archive name
	header := lipgloss.NewStyle().
		Foreground(a.theme.Header).
		Bold(true).
		Render(archive.Name)
	lines = append(lines, header)
	lines = append(lines, "")
	
	// Details
	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"Size:", humanize.Bytes(uint64(archive.Size)), a.theme.Metadata},
		{"Created:", archive.Created.Format("January 2, 2006 3:04 PM"), a.theme.Metadata},
		{"Status:", getStatusText(archive.Status), getStatusColor(archive.Status, a.theme)},
		{"Location:", archive.Path, a.theme.Foreground},
		{"Profile:", archive.Profile, a.theme.Metadata},
	}
	
	if archive.Checksum != "" {
		details = append(details, struct {
			label string
			value string
			color lipgloss.Color
		}{"Checksum:", archive.Checksum[:16] + "...", a.theme.Metadata})
	}
	
	for _, detail := range details {
		label := lipgloss.NewStyle().
			Foreground(a.theme.Foreground).
			Render(fmt.Sprintf("%-12s", detail.label))
		value := lipgloss.NewStyle().
			Foreground(detail.color).
			Render(detail.value)
		lines = append(lines, label+value)
	}
	
	lines = append(lines, "")
	
	// Commands
	commands := a.renderCommands()
	lines = append(lines, commands)
	
	return strings.Join(lines, "\n")
}

// renderConfirmDialog renders the confirmation dialog
func (a *App) renderConfirmDialog() string {
	// Create confirmation box
	confirmBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.StatusMiss).
		Padding(1).
		Background(lipgloss.Color("#000000")).
		Foreground(a.theme.Foreground).
		Width(50).
		Align(lipgloss.Center)
	
	content := fmt.Sprintf("%s\n\n[y] Yes  [n] No  [Enter] Confirm  [Esc] Cancel", a.confirmMsg)
	
	return confirmBox.Render(content)
}

// Helper functions
func getStatusText(status string) string {
	switch status {
	case "present":
		return "Present ✓"
	case "missing":
		return "Missing ?"
	case "deleted":
		return "Deleted X"
	default:
		return status
	}
}

func getStatusColor(status string, theme Theme) lipgloss.Color {
	switch status {
	case "present":
		return theme.StatusOK
	case "missing":
		return theme.StatusMiss
	case "deleted":
		return theme.StatusDel
	default:
		return theme.Foreground
	}
}

// formatAge formats time.Time into a human readable age string
func formatAge(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Hour {
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	} else if duration < 7*24*time.Hour {
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	} else if duration < 30*24*time.Hour {
		return fmt.Sprintf("%dw ago", int(duration.Hours()/(7*24)))
	} else {
		return fmt.Sprintf("%dm ago", int(duration.Hours()/(30*24)))
	}
}
