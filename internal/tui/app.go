package tui

import (
	"fmt"
	"strings"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

// SimpleApp is the clean TUI implementation
type SimpleApp struct {
	archives      []*storage.Archive
	cursor        int
	width         int
	height        int
	currentView   ViewType
	selected      map[string]bool
	showConfirm   bool
	confirmMsg    string
	confirmAction func() tea.Cmd
	theme         Theme
	manager       *storage.Manager
	resolver      *storage.Resolver
	
	// Viewport for proper content management
	viewport      viewport.Model
}

type ViewType int

const (
	ListView ViewType = iota
	DetailView
)

// NewSimpleApp creates the TUI app
func NewSimpleApp(themeName string) *SimpleApp {
	cfg, _ := config.Load()
	manager, _ := storage.NewManager(cfg.Storage.ManagedPath)
	
	var resolver *storage.Resolver
	if manager != nil {
		resolver = storage.NewResolver(manager.Registry())
	}

	// Initialize viewport with margins
	vp := viewport.New(80, 24) // Default size, will be updated on window resize
	vp.Style = lipgloss.NewStyle().
		Margin(1, 2). // 1 line top/bottom, 2 chars left/right
		Border(lipgloss.HiddenBorder())

	return &SimpleApp{
		currentView: ListView,
		selected:    make(map[string]bool),
		theme:       GetTheme(themeName),
		manager:     manager,
		resolver:    resolver,
		viewport:    vp,
	}
}

// Init loads archives
func (a *SimpleApp) Init() tea.Cmd {
	return func() tea.Msg {
		if a.manager == nil {
			return archivesLoadedMsg{err: fmt.Errorf("no storage manager")}
		}
		archives, err := a.manager.List()
		return archivesLoadedMsg{archives: archives, err: err}
	}
}

// Update handles input
func (a *SimpleApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		if a.showConfirm {
			switch m.String() {
			case "enter", "y":
				a.showConfirm = false
				if a.confirmAction != nil {
					return a, a.confirmAction()
				}
			case "esc", "n":
				a.showConfirm = false
				a.confirmAction = nil
			}
			return a, nil
		}
		
		// Main navigation
		switch a.currentView {
		case ListView:
			return a.handleListKeys(m)
		case DetailView:
			return a.handleDetailKeys(m)
		}
		
	case tea.WindowSizeMsg:
		a.width, a.height = m.Width, m.Height
		// Update viewport size with proper margins
		a.viewport.Width = m.Width - 4  // 2 chars margin on each side
		a.viewport.Height = m.Height - 2 // 1 line margin top/bottom
		
	case archivesLoadedMsg:
		if m.err == nil {
			a.archives = m.archives
		}
	}
	
	return a, nil
}

// View renders the interface
func (a *SimpleApp) View() string {
	// Generate content based on current view
	var content string
	switch a.currentView {
	case ListView:
		content = a.renderListContent()
	case DetailView:
		content = a.renderDetailContent()
	}
	
	// Set viewport content
	a.viewport.SetContent(content)
	
	// Render viewport (handles margins properly)
	base := a.viewport.View()
	
	// Overlay confirmation dialog if needed
	if a.showConfirm {
		return base + "\n" + a.renderConfirm()
	}
	
	return base
}

// renderListContent generates content for viewport (no manual margins)
func (a *SimpleApp) renderListContent() string {
	var lines []string
	
	// Header (viewport handles margins)
	header := lipgloss.NewStyle().
		Foreground(a.theme.Header).
		Bold(true).
		Render("7zarch-go")
	lines = append(lines, header)
	lines = append(lines, "")
	
	// Summary
	summary := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render(fmt.Sprintf("Archives: %d", len(a.archives)))
	lines = append(lines, summary)
	lines = append(lines, "")
	
	// Archives (viewport handles margins)
	for i, archive := range a.archives {
		line := a.renderArchive(archive, i == a.cursor)
		lines = append(lines, line)
	}
	
	lines = append(lines, "")
	
	// Commands
	commands := lipgloss.NewStyle().
		Foreground(a.theme.Commands).
		Render("[Enter] Details  [Space] Select  [d] Delete  [m] Move  [u] Upload  [q] Quit")
	lines = append(lines, commands)
	
	return strings.Join(lines, "\n")
}

// Archive item renderer
func (a *SimpleApp) renderArchive(archive *storage.Archive, isSelected bool) string {
	selector := "  "
	if isSelected {
		selector = "> "
	}
	
	checkbox := "[ ]"
	if a.selected[archive.UID] {
		checkbox = "[✓]"
	}
	
	statusIcon := "✓"
	
	name := archive.Name
	if len(name) > 30 {
		name = name[:27] + "..."
	}
	
	size := humanize.Bytes(uint64(archive.Size))
	
	// Build content
	content := fmt.Sprintf("%s%s %-30s %10s %8s %s",
		selector, checkbox, name, size, "2h ago", statusIcon)
	
	if isSelected {
		return lipgloss.NewStyle().
			Background(a.theme.Selection).
			Foreground(a.theme.SelText).
			Render(content)
	}
	
	return lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render(content)
}

// renderDetailContent generates detail content for viewport
func (a *SimpleApp) renderDetailContent() string {
	if len(a.archives) == 0 || a.cursor >= len(a.archives) {
		return "No archive selected"
	}
	
	archive := a.archives[a.cursor]
	var lines []string
	
	// Header (viewport handles margins)
	header := lipgloss.NewStyle().
		Foreground(a.theme.Header).
		Bold(true).
		Render(archive.Name)
	lines = append(lines, header)
	lines = append(lines, "")
	
	// Details with proper formatting
	sizeText := "Size: " + lipgloss.NewStyle().Foreground(a.theme.Metadata).Render(humanize.Bytes(uint64(archive.Size)))
	lines = append(lines, sizeText)
	
	createdText := "Created: " + lipgloss.NewStyle().Foreground(a.theme.Metadata).Render(archive.Created.Format("January 2, 2006 3:04 PM"))
	lines = append(lines, createdText)
	
	statusText := "Status: " + a.getStatusDisplay(archive.Status)
	lines = append(lines, statusText)
	
	locationText := "Location: " + lipgloss.NewStyle().Foreground(a.theme.Foreground).Render(archive.Path)
	lines = append(lines, locationText)
	
	if archive.Checksum != "" {
		checksumText := "Checksum: " + lipgloss.NewStyle().Foreground(a.theme.Metadata).Render(archive.Checksum[:16]+"...")
		lines = append(lines, checksumText)
	}
	
	lines = append(lines, "")
	
	// Commands
	commands := lipgloss.NewStyle().
		Foreground(a.theme.Commands).
		Render("[Enter] Back  [d] Delete  [m] Move  [u] Upload  [q] Quit")
	lines = append(lines, commands)
	
	return strings.Join(lines, "\n")
}

// Confirmation dialog (rendered as overlay)
func (a *SimpleApp) renderConfirm() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.StatusMiss).
		Padding(1).
		Margin(2, 4). // Center with margins
		Foreground(a.theme.Foreground).
		Background(lipgloss.Color("#000000")) // Dark background for contrast
	
	content := fmt.Sprintf("%s\n\n[y] Yes  [n] No  [Esc] Cancel", a.confirmMsg)
	return box.Render(content)
}

// getStatusDisplay returns styled status text
func (a *SimpleApp) getStatusDisplay(status string) string {
	switch status {
	case "present":
		return lipgloss.NewStyle().Foreground(a.theme.StatusOK).Render("Present ✓")
	case "missing":
		return lipgloss.NewStyle().Foreground(a.theme.StatusMiss).Render("Missing ?")
	case "deleted":
		return lipgloss.NewStyle().Foreground(a.theme.StatusDel).Render("Deleted X")
	default:
		return lipgloss.NewStyle().Foreground(a.theme.Foreground).Render(status)
	}
}

// Key handlers
func (a *SimpleApp) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "up", "k":
		if a.cursor > 0 {
			a.cursor--
		}
	case "down", "j":
		if a.cursor < len(a.archives)-1 {
			a.cursor++
		}
	case "enter":
		a.currentView = DetailView
	case " ":
		if len(a.archives) > 0 && a.cursor < len(a.archives) {
			uid := a.archives[a.cursor].UID
			a.selected[uid] = !a.selected[uid]
		}
	case "d":
		a.confirmMsg = "Delete archive?"
		a.showConfirm = true
	}
	return a, nil
}

func (a *SimpleApp) handleDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "enter", "esc":
		a.currentView = ListView
	}
	return a, nil
}

// Messages
type archivesLoadedMsg struct {
	archives []*storage.Archive
	err      error
}
