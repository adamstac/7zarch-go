package cmd

import (
	"fmt"

	"github.com/adamstac/7zarch-go/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// TuiCmd starts the interactive TUI (7EP-0010)
func TuiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Interactive TUI for archive management",
		Long:  "Launch the interactive terminal UI to browse, inspect, and manage archives.",
		RunE:  runTui,
	}
	return cmd
}

func runTui(cmd *cobra.Command, args []string) error {
	app := tui.NewApp()
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tui exited with error: %w", err)
	}
	return nil
}
