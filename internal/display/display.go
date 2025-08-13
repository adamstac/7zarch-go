package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/adamstac/7zarch-go/internal/storage"
	"golang.org/x/term"
)

// Mode represents a display mode type
type Mode string

const (
	ModeAuto      Mode = "auto"
	ModeTable     Mode = "table"
	ModeCompact   Mode = "compact"
	ModeCard      Mode = "card"
	ModeTree      Mode = "tree"
	ModeDashboard Mode = "dashboard"
)

// Options configures display behavior
type Options struct {
	Mode         Mode
	Details      bool
	Width        int
	ShowHeaders  bool
	Columns      []string
	GroupBy      string
	SortBy       []string
	ColorEnabled bool
}

// Context provides environmental information for display decisions
type Context struct {
	TerminalWidth int
	ArchiveCount  int
	FilterContext string
	OutputPiped   bool
}

// Display interface that all display modes must implement
type Display interface {
	Render(archives []*storage.Archive, opts Options) error
	Name() string
	MinWidth() int
}

// Manager handles display mode selection and rendering
type Manager struct {
	displays map[Mode]Display
	context  Context
}

// NewManager creates a new display manager
func NewManager() *Manager {
	return &Manager{
		displays: make(map[Mode]Display),
		context:  detectContext(),
	}
}

// Register adds a display mode to the manager
func (m *Manager) Register(mode Mode, display Display) {
	m.displays[mode] = display
}

// Render displays archives using the specified or auto-detected mode
func (m *Manager) Render(archives []*storage.Archive, opts Options) error {
	mode := opts.Mode
	if mode == ModeAuto || mode == "" {
		mode = m.detectBestMode(opts)
	}

	display, exists := m.displays[mode]
	if !exists {
		return fmt.Errorf("unknown display mode: %s", mode)
	}

	// Update options with terminal width if not specified
	if opts.Width == 0 {
		opts.Width = m.context.TerminalWidth
	}

	return display.Render(archives, opts)
}

// detectBestMode determines the optimal display mode based on context
func (m *Manager) detectBestMode(opts Options) Mode {
	ctx := m.context

	// Piped output always uses compact
	if ctx.OutputPiped {
		return ModeCompact
	}

	// Narrow terminals use compact
	if ctx.TerminalWidth < 80 {
		return ModeCompact
	}

	// Filter-specific defaults
	if ctx.FilterContext == "missing" {
		return ModeCompact
	}

	// Large collections use table for scanning
	if ctx.ArchiveCount > 50 {
		return ModeTable
	}

	// Default to table for normal use
	return ModeTable
}

// detectContext gathers environmental information
func detectContext() Context {
	ctx := Context{
		TerminalWidth: getTerminalWidth(),
		OutputPiped:   !isTerminal(),
	}
	return ctx
}

// getTerminalWidth returns the current terminal width
func getTerminalWidth() int {
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return width
	}
	return 80 // default fallback
}

// isTerminal checks if output is going to a terminal
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Helper functions for formatting

// TruncateString truncates a string to maxLen with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}

// PadRight pads a string to the specified width
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// FormatSize formats bytes as human-readable size
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatStatus formats archive status with consistent icons
func FormatStatus(status string, useIcons bool) string {
	if useIcons {
		switch status {
		case "present":
			return "✓"
		case "missing":
			return "?"
		case "deleted":
			return "X"
		default:
			return "?"
		}
	} else {
		switch status {
		case "present":
			return "OK"
		case "missing":
			return "MISS"
		case "deleted":
			return "DEL"
		default:
			return status
		}
	}
}
