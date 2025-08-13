package cmd

import (
	"fmt"

	"github.com/adamstac/7zarch-go/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// TuiCmd starts the interactive TUI (7EP-0010)
func TuiCmd() *cobra.Command {
	var theme string
	
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Interactive TUI for archive management",
		Long: `Launch the interactive terminal UI to browse, inspect, and manage archives.

Simple list-based interface optimized for podcast archival workflows.
Navigate with arrow keys, select with space, and take actions with single letters.

Available themes: blue, green, purple, cyan, charmbracelet, dracula, dracula-warm, dracula-cool, dracula-minimal`,
		Example: `  # Launch with default Dracula theme
  7zarch-go tui
  
  # Launch with specific theme
  7zarch-go tui --theme blue
  7zarch-go tui --theme dracula-warm`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTui(theme)
		},
	}
	
	cmd.Flags().StringVar(&theme, "theme", "dracula", "Color theme (blue|green|purple|cyan|charmbracelet|dracula|dracula-warm|dracula-cool|dracula-minimal)")
	
	return cmd
}

func runTui(theme string) error {
	app := tui.NewSimpleApp(theme)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tui exited with error: %w", err)
	}
	return nil
}
