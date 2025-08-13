package tui

import (
	"fmt"
	"strings"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
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

	return &SimpleApp{
		currentView: ListView,
		selected:    make(map[string]bool),
		theme:       GetTheme(themeName),
		manager:     manager,
		resolver:    resolver,
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
		
	case archivesLoadedMsg:
		if m.err == nil {
			a.archives = m.archives
		}
	}
	
	return a, nil
}

// View renders the interface
func (a *SimpleApp) View() string {
	base := ""
	switch a.currentView {
	case ListView:
		base = a.renderList()
	case DetailView:
		base = a.renderDetail()
	}
	
	if a.showConfirm {
		return base + "\n\n" + a.renderConfirm()
	}
	
	return base
}

// List view renderer
func (a *SimpleApp) renderList() string {
	var lines []string
	
	// Header
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
	
	// Archives
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

// Detail view renderer
func (a *SimpleApp) renderDetail() string {
	if len(a.archives) == 0 || a.cursor >= len(a.archives) {
		return "No archive selected"
	}
	
	archive := a.archives[a.cursor]
	var lines []string
	
	// Header
	header := lipgloss.NewStyle().
		Foreground(a.theme.Header).
		Bold(true).
		Render(archive.Name)
	lines = append(lines, header)
	lines = append(lines, "")
	
	// Details
	lines = append(lines, "Size: "+humanize.Bytes(uint64(archive.Size)))
	lines = append(lines, "Status: Present ✓")
	lines = append(lines, "")
	
	// Commands
	commands := lipgloss.NewStyle().
		Foreground(a.theme.Commands).
		Render("[Enter] Back  [d] Delete  [q] Quit")
	lines = append(lines, commands)
	
	return strings.Join(lines, "\n")
}

// Confirmation dialog
func (a *SimpleApp) renderConfirm() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.StatusMiss).
		Padding(1).
		Foreground(a.theme.Foreground)
	
	content := fmt.Sprintf("%s\n\n[y] Yes  [n] No", a.confirmMsg)
	return box.Render(content)
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
