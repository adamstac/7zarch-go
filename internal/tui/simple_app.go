package tui

import (
	"fmt"
	"strings"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SimpleApp is the new simplified TUI implementation
type SimpleApp struct {
	// Core data
	archives []*storage.Archive
	cursor   int
	width    int
	height   int
	
	// Views
	currentView ViewType
	
	// State
	selected      map[string]bool
	showConfirm   bool
	confirmMsg    string
	confirmAction func() tea.Cmd
	
	// Theme
	theme Theme
	
	// Storage
	manager  *storage.Manager
	resolver *storage.Resolver
}

// NewSimpleApp creates a new simplified TUI app
func NewSimpleApp(themeName string) *SimpleApp {
	// Initialize storage
	cfg, _ := config.Load()
	manager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		manager = nil
	}
	
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

// Init initializes the app
func (a *SimpleApp) Init() tea.Cmd {
	return a.loadArchives
}

// Update handles messages
func (a *SimpleApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		// Handle confirmation dialog
		if a.showConfirm {
			return a.handleConfirmKeys(m)
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
		return a, nil
		
	case archivesLoadedMsg:
		if m.err != nil {
			// Handle error
			return a, nil
		}
		a.archives = m.archives
		return a, nil
	}
	
	return a, nil
}

// View renders the current view
func (a *SimpleApp) View() string {
	if a.showConfirm {
		return a.renderListView() + "\n\n" + a.renderConfirmDialog()
	}
	
	switch a.currentView {
	case ListView:
		return a.renderListView()
	case DetailView:
		return a.renderDetailView()
	default:
		return a.renderListView()
	}
}

// Messages
type archivesLoadedMsg struct {
	archives []*storage.Archive
	err      error
}

// loadArchives loads the archive list
func (a *SimpleApp) loadArchives() tea.Msg {
	if a.manager == nil {
		return archivesLoadedMsg{err: fmt.Errorf("storage manager not available")}
	}
	
	archives, err := a.manager.List()
	return archivesLoadedMsg{archives: archives, err: err}
}

// Key handlers
func (a *SimpleApp) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "y":
		a.showConfirm = false
		if a.confirmAction != nil {
			cmd := a.confirmAction()
			a.confirmAction = nil
			return a, cmd
		}
		return a, nil
	case "esc", "n":
		a.showConfirm = false
		a.confirmAction = nil
		return a, nil
	}
	return a, nil
}

func (a *SimpleApp) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
		
	// Navigation
	case "up", "k":
		if a.cursor > 0 {
			a.cursor--
		}
		return a, nil
	case "down", "j":
		if a.cursor < len(a.archives)-1 {
			a.cursor++
		}
		return a, nil
		
	// Selection
	case " ":
		if len(a.archives) > 0 && a.cursor < len(a.archives) {
			uid := a.archives[a.cursor].UID
			a.selected[uid] = !a.selected[uid]
		}
		return a, nil
		
	// View navigation
	case "enter":
		if len(a.archives) > 0 && a.cursor < len(a.archives) {
			a.currentView = DetailView
		}
		return a, nil
		
	// Actions
	case "d":
		return a.confirmDelete()
	case "m":
		return a.confirmMove()
	case "u":
		return a.confirmUpload()
		
	// Utility
	case "r":
		return a, a.loadArchives
	}
	
	return a, nil
}

func (a *SimpleApp) handleDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "enter", "esc":
		a.currentView = ListView
		return a, nil
		
	// Actions from detail view
	case "d":
		return a.confirmDelete()
	case "m":
		return a.confirmMove()
	case "u":
		return a.confirmUpload()
	}
	
	return a, nil
}
