package tui

import (
	"fmt"
	"os/exec"

	"github.com/adamstac/7zarch-go/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// Action confirmation methods
func (a *SimpleApp) confirmDelete() (tea.Model, tea.Cmd) {
	selected := a.getSelectedArchives()
	if len(selected) == 0 {
		return a, nil
	}
	
	if len(selected) == 1 {
		a.confirmMsg = fmt.Sprintf("Delete %s?", selected[0].Name)
	} else {
		a.confirmMsg = fmt.Sprintf("Delete %d archives?", len(selected))
	}
	
	a.showConfirm = true
	a.confirmAction = func() tea.Cmd {
		return a.executeDelete(selected)
	}
	
	return a, nil
}

func (a *SimpleApp) confirmMove() (tea.Model, tea.Cmd) {
	selected := a.getSelectedArchives()
	if len(selected) == 0 {
		return a, nil
	}
	
	if len(selected) == 1 {
		a.confirmMsg = fmt.Sprintf("Move %s to managed storage?", selected[0].Name)
	} else {
		a.confirmMsg = fmt.Sprintf("Move %d archives to managed storage?", len(selected))
	}
	
	a.showConfirm = true
	a.confirmAction = func() tea.Cmd {
		return a.executeMove(selected)
	}
	
	return a, nil
}

func (a *SimpleApp) confirmUpload() (tea.Model, tea.Cmd) {
	selected := a.getSelectedArchives()
	if len(selected) == 0 {
		return a, nil
	}
	
	if len(selected) == 1 {
		a.confirmMsg = fmt.Sprintf("Upload %s to TrueNAS?", selected[0].Name)
	} else {
		a.confirmMsg = fmt.Sprintf("Upload %d archives to TrueNAS?", len(selected))
	}
	
	a.showConfirm = true
	a.confirmAction = func() tea.Cmd {
		return a.executeUpload(selected)
	}
	
	return a, nil
}

// Action execution methods
func (a *SimpleApp) executeDelete(archives []*storage.Archive) tea.Cmd {
	return func() tea.Msg {
		for _, archive := range archives {
			// Use CLI command integration for consistency
			cmd := exec.Command("7zarch-go", "delete", archive.UID)
			cmd.Run() // Simple execution, ignore output for now
		}
		
		// Clear selection after action
		a.selected = make(map[string]bool)
		
		// Reload archives
		return a.loadArchives()
	}
}

func (a *SimpleApp) executeMove(archives []*storage.Archive) tea.Cmd {
	return func() tea.Msg {
		for _, archive := range archives {
			// Use CLI command integration
			cmd := exec.Command("7zarch-go", "move", archive.UID)
			cmd.Run()
		}
		
		// Clear selection after action
		a.selected = make(map[string]bool)
		
		// Reload archives
		return a.loadArchives()
	}
}

func (a *SimpleApp) executeUpload(archives []*storage.Archive) tea.Cmd {
	return func() tea.Msg {
		for _, archive := range archives {
			// Use CLI command integration
			cmd := exec.Command("7zarch-go", "upload", archive.UID)
			cmd.Run()
		}
		
		// Clear selection after action
		a.selected = make(map[string]bool)
		
		// Reload archives
		return a.loadArchives()
	}
}

// getSelectedArchives returns selected archives or cursor archive if none selected
func (a *SimpleApp) getSelectedArchives() []*storage.Archive {
	var selected []*storage.Archive
	
	// Check if any archives are multi-selected
	hasSelection := false
	for _, isSelected := range a.selected {
		if isSelected {
			hasSelection = true
			break
		}
	}
	
	// If no multi-selection, use cursor position
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
