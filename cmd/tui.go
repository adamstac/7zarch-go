package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/display/modes"
	"github.com/adamstac/7zarch-go/internal/storage"
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

type contentMsg struct {
	content string
	err     error
}

type tuiModel struct {
	content string
	width   int
	height  int
	err     error
	mode    display.Mode
}

func initialModel() tuiModel {
	return tuiModel{mode: display.ModeDashboard}
}

func (m tuiModel) Init() tea.Cmd {
	return fetchContent(m.mode)
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			return m, fetchContent(m.mode)
		case "1":
			m.mode = display.ModeTable
			return m, fetchContent(m.mode)
		case "2":
			m.mode = display.ModeCompact
			return m, fetchContent(m.mode)
		case "3":
			m.mode = display.ModeCard
			return m, fetchContent(m.mode)
		case "4":
			m.mode = display.ModeTree
			return m, fetchContent(m.mode)
		case "5":
			m.mode = display.ModeDashboard
			return m, fetchContent(m.mode)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case contentMsg:
		m.err = msg.err
		if msg.err == nil {
			m.content = msg.content
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("7zarch-go TUI — %s (q quit, r refresh, 1-5 switch modes)\n", m.mode))
	b.WriteString("────────────────────────────────────────────────────\n")
	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else if m.content == "" {
		b.WriteString("Loading…\n")
	} else {
		b.WriteString(m.content)
	}
	return b.String()
}

func runTui(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tui exited with error: %w", err)
	}
	return nil
}

func fetchContent(mode display.Mode) tea.Cmd {
	return func() tea.Msg {
		content, err := generateContent(mode)
		return contentMsg{content: content, err: err}
	}
}

// generateContent renders the selected display mode and captures its output
func generateContent(mode display.Mode) (string, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	// Initialize storage manager
	mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return "", fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer mgr.Close()

	// List all registry archives
	archives, err := mgr.List()
	if err != nil {
		return "", fmt.Errorf("failed to list archives: %w", err)
	}

	// Set up display manager and register modes
	dm := display.NewManager()
	dm.Register(display.ModeTable, modes.NewTableDisplay())
	dm.Register(display.ModeCompact, modes.NewCompactDisplay())
	dm.Register(display.ModeCard, modes.NewCardDisplay())
	dm.Register(display.ModeTree, modes.NewTreeDisplay())
	dm.Register(display.ModeDashboard, modes.NewDashboardDisplay())

	// Capture stdout while rendering
	var buf bytes.Buffer
	err = withCapturedStdout(func() error {
		return dm.Render(archives, display.Options{Mode: mode, ShowHeaders: mode != display.ModeCompact})
	}, &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// withCapturedStdout temporarily redirects stdout to a buffer while fn runs
func withCapturedStdout(fn func() error, out *bytes.Buffer) error {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	origStdout := os.Stdout
	// Replace stdout with the write end of the pipe
	os.Stdout = w
	// Ensure we restore stdout
	defer func() {
		_ = w.Close()
		os.Stdout = origStdout
	}()

	// Reader goroutine to copy pipe output into buffer
	done := make(chan error, 1)
	go func() {
		_, copyErr := out.ReadFrom(r)
		done <- copyErr
	}()

	// Run function
	callErr := fn()
	// Close writer to signal reader EOF
	_ = w.Close()
	// Wait for reader to finish
	readErr := <-done
	_ = r.Close()
	if callErr != nil {
		return callErr
	}
	return readErr
}
