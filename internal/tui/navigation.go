package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleListViewKeys handles keyboard input for the list view
func (a *App) handleListViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case "home", "g":
		a.cursor = 0
		return a, nil
	case "end", "G":
		if len(a.archives) > 0 {
			a.cursor = len(a.archives) - 1
		}
		return a, nil
		
	// Selection
	case " ":
		if len(a.archives) > 0 && a.cursor < len(a.archives) {
			uid := a.archives[a.cursor].UID
			a.selected[uid] = !a.selected[uid]
		}
		return a, nil
		
	// Views
	case "enter":
		if len(a.archives) > 0 && a.cursor < len(a.archives) {
			a.currentView = DetailView
		}
		return a, nil
		
	// Actions
	case "d":
		return a.handleDeleteAction()
	case "m":
		return a.handleMoveAction()
	case "u":
		return a.handleUploadAction()
		
	// Utility
	case "r":
		return a, a.loadArchivesCmd()
	}
	
	return a, nil
}

// handleDetailViewKeys handles keyboard input for the details view
func (a *App) handleDetailViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "enter", "esc":
		a.currentView = ListView
		return a, nil
		
	// Actions from details view
	case "d":
		return a.handleDeleteAction()
	case "m":
		return a.handleMoveAction()
	case "u":
		return a.handleUploadAction()
	}
	
	return a, nil
}

// handleDeleteAction handles delete with confirmation
func (a *App) handleDeleteAction() (tea.Model, tea.Cmd) {
	selectedArchives := a.getSelectedArchives()
	if len(selectedArchives) == 0 {
		return a, nil
	}
	
	// Show confirmation
	if len(selectedArchives) == 1 {
		a.confirmMsg = "Delete " + selectedArchives[0].Name + "?"
	} else {
		a.confirmMsg = fmt.Sprintf("Delete %d archives?", len(selectedArchives))
	}
	
	a.showConfirm = true
	a.confirmAction = func() tea.Cmd {
		// Execute delete via CLI command integration
		return a.executeDelete(selectedArchives)
	}
	
	return a, nil
}

// handleMoveAction handles move action
func (a *App) handleMoveAction() (tea.Model, tea.Cmd) {
	selectedArchives := a.getSelectedArchives()
	if len(selectedArchives) == 0 {
		return a, nil
	}
	
	// For now, move to managed storage (simple implementation)
	a.confirmMsg = fmt.Sprintf("Move %d archives to managed storage?", len(selectedArchives))
	a.showConfirm = true
	a.confirmAction = func() tea.Cmd {
		return a.executeMove(selectedArchives)
	}
	
	return a, nil
}

// handleUploadAction handles upload action
func (a *App) handleUploadAction() (tea.Model, tea.Cmd) {
	selectedArchives := a.getSelectedArchives()
	if len(selectedArchives) == 0 {
		return a, nil
	}
	
	a.confirmMsg = fmt.Sprintf("Upload %d archives to TrueNAS?", len(selectedArchives))
	a.showConfirm = true
	a.confirmAction = func() tea.Cmd {
		return a.executeUpload(selectedArchives)
	}
	
	return a, nil
}

// getSelectedArchives returns currently selected archives or cursor archive
func (a *App) getSelectedArchives() []*storage.Archive {
	var selected []*storage.Archive
	
	// If nothing is multi-selected, use cursor position
	hasSelection := false
	for _, isSelected := range a.selected {
		if isSelected {
			hasSelection = true
			break
		}
	}
	
	if !hasSelection {
		if len(a.archives) > 0 && a.cursor < len(a.archives) {
			selected = append(selected, a.archives[a.cursor])
		}
		return selected
	}
	
	// Return all selected archives
	for _, archive := range a.archives {
		if a.selected[archive.UID] {
			selected = append(selected, archive)
		}
	}
	
	return selected
}
