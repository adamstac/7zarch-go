package cmd

import (
	"fmt"

	"github.com/adamstac/7zarch-go/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// BrowseCmd launches the interactive archive browser (primary TUI interface)
func BrowseCmd() *cobra.Command {
	var theme string
	
	cmd := &cobra.Command{
		Use:   "browse",
		Short: "Browse archives with interactive visual interface",
		Long: `Launch the interactive archive browser with beautiful themes and simple navigation.

Browse your archive collection with a clean, visual interface optimized for 
podcast archival workflows. Navigate with arrow keys, select with space, 
and take actions with single letters.

Available themes: blue, green, purple, cyan, charmbracelet, dracula, dracula-warm, dracula-cool, dracula-minimal`,
		Example: `  # Browse with default Dracula theme
  7zarch-go browse
  
  # Browse with specific theme
  7zarch-go browse --theme blue
  7zarch-go browse --theme dracula-warm
  
  # Quick access aliases
  7zarch-go ui --theme purple
  7zarch-go i --theme charmbracelet`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInteractiveBrowser(theme)
		},
	}
	
	cmd.Flags().StringVar(&theme, "theme", "dracula", "Color theme (blue|green|purple|cyan|charmbracelet|dracula|dracula-warm|dracula-cool|dracula-minimal)")
	
	return cmd
}

// UICmd is an alias for browse (shorter name)
func UICmd() *cobra.Command {
	var theme string
	
	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Interactive UI for archive management (alias for browse)",
		Long:  "Alias for the browse command. Launch interactive archive browser.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInteractiveBrowser(theme)
		},
	}
	
	cmd.Flags().StringVar(&theme, "theme", "dracula", "Color theme")
	
	return cmd
}

// InteractiveCmd is a single-letter alias for browse
func InteractiveCmd() *cobra.Command {
	var theme string
	
	cmd := &cobra.Command{
		Use:   "i",
		Short: "Interactive browser (alias for browse)",
		Long:  "Single-letter alias for the browse command.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInteractiveBrowser(theme)
		},
	}
	
	cmd.Flags().StringVar(&theme, "theme", "dracula", "Color theme")
	
	return cmd
}

// runInteractiveBrowser launches the TUI with specified theme
func runInteractiveBrowser(theme string) error {
	app := tui.NewSimpleApp(theme)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("interactive browser exited with error: %w", err)
	}
	return nil
}
