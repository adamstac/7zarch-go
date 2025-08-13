# 7EP-0010: TUI Implementation Guide - Simplified Charmbracelet Approach

**For**: AC overnight implementation  
**Philosophy**: Simple, beautiful, composable - the Charmbracelet way  
**Scope**: MVP that delights users immediately

## 🎯 Core Principle: Leverage What We Built

We already have 5 beautiful display modes. The TUI should be a **navigation wrapper** around these displays, not a reimplementation.

## 📐 Architecture: Three Simple Views

### View 1: Dashboard (Entry Point)
**Purpose**: Show the user their archive world at a glance  
**Implementation**: Wrap our existing `dashboard.go` display

```go
type DashboardView struct {
    content  string  // Rendered dashboard from display package
    viewport viewport.Model  // Bubbles viewport for scrolling
    help     help.Model      // Bubbles help component
}

func (v DashboardView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "l", "enter":
            return switchToListView()  // Navigate to list
        case "q", "ctrl+c":
            return v, tea.Quit
        case "?":
            v.help.ShowAll = !v.help.ShowAll
        }
    }
    // Update viewport for scrolling
    v.viewport, cmd = v.viewport.Update(msg)
    return v, cmd
}

func (v DashboardView) View() string {
    // Render our existing dashboard display into the viewport
    archives := storage.LoadArchives()
    var buf bytes.Buffer
    dashboardDisplay := modes.NewDashboardDisplay()
    dashboardDisplay.Render(archives, display.Options{})
    
    v.viewport.SetContent(buf.String())
    return v.viewport.View() + "\n" + v.help.View()
}
```

### View 2: List (Browse & Select)
**Purpose**: Interactive archive browsing with mode switching  
**Implementation**: Bubbles list + our display modes

```go
type ListView struct {
    list        list.Model       // Bubbles list component
    archives    []*storage.Archive
    displayMode display.Mode     // Current display mode
    selected    map[string]bool  // Multi-select tracking
}

func (v ListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "1":  // Switch display modes on the fly!
            v.displayMode = display.ModeTable
        case "2":
            v.displayMode = display.ModeCompact
        case "3":
            v.displayMode = display.ModeCard
        case "4":
            v.displayMode = display.ModeTree
        case "enter":
            selected := v.archives[v.list.Index()]
            return switchToDetailView(selected)
        case "space":
            // Toggle selection for batch operations
            id := v.archives[v.list.Index()].UID
            v.selected[id] = !v.selected[id]
        case "d":
            // Quick delete with confirmation
            return v, showConfirmation("Delete selected archives?")
        case "esc", "q":
            return switchToDashboardView()
        }
    }
}

// The magic: each item renders using our display modes!
func (v ListView) itemView(index int, archive *storage.Archive) string {
    // Use compact mode for list items, but allow switching
    switch v.displayMode {
    case display.ModeCompact:
        return compactLine(archive)
    case display.ModeCard:
        return cardView(archive)  // Mini card in list
    default:
        return compactLine(archive)
    }
}
```

### View 3: Detail (Inspect & Act)
**Purpose**: Full archive details with operations  
**Implementation**: Viewport + our card display + operation buttons

```go
type DetailView struct {
    archive  *storage.Archive
    viewport viewport.Model
    buttons  []string  // ["Delete", "Move", "Upload", "Verify"]
    selected int       // Which button is selected
}

func (v DetailView) View() string {
    // Render using our card display
    var buf bytes.Buffer
    cardDisplay := modes.NewCardDisplay()
    cardDisplay.Render([]*storage.Archive{v.archive}, display.Options{Details: true})
    
    // Add operation buttons using Lipgloss
    buttonBar := renderButtons(v.buttons, v.selected)
    
    return v.viewport.View() + "\n" + buttonBar
}
```

## 🎨 Charmbracelet Component Usage

### Essential Bubbles Components

```go
import (
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
)

// Pre-built components to use:
// - list: For archive browsing (View 2)
// - viewport: For scrolling content (all views)
// - textinput: For search/filter input
// - spinner: For loading states
// - help: For keyboard shortcuts (all views)
```

### Lipgloss Styling Patterns

```go
// Define a consistent style system
var (
    // Base styles - Charmbracelet philosophy: subtle, beautiful
    subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
    
    // Status styles matching our display system
    successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
    warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
    errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
    
    // Layout helpers
    borderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62"))
    
    // Button styles for operations
    buttonStyle = lipgloss.NewStyle().
        Padding(0, 2).
        Background(lipgloss.Color("62")).
        Foreground(lipgloss.Color("230"))
    
    activeButtonStyle = buttonStyle.Copy().
        Background(lipgloss.Color("42"))
)
```

## 🚀 Implementation Path for AC

### Hour 1-2: Foundation
```bash
# Create the TUI structure
internal/tui/
├── app.go           # Main app model and routing
├── views/
│   ├── dashboard.go # Dashboard view
│   ├── list.go      # List view
│   └── detail.go    # Detail view
├── styles/
│   └── theme.go     # Lipgloss styles
└── keys/
    └── bindings.go  # Key bindings using Bubbles key package
```

### Hour 3-4: Dashboard View
1. Create basic Bubble Tea app structure
2. Implement dashboard view using our `dashboard.Display`
3. Add viewport for scrolling
4. Add help component showing shortcuts
5. Test navigation to list view (just a placeholder)

### Hour 5-6: List View  
1. Implement list using Bubbles list component
2. Load archives and create list items
3. Add display mode switching (1,2,3,4 keys)
4. Add multi-select with space bar
5. Connect navigation to detail view

### Hour 7-8: Detail View & Polish
1. Implement detail view with full card display
2. Add operation buttons (styled with Lipgloss)
3. Add transitions between views
4. Polish keyboard navigation
5. Add loading spinners for operations

## 🎯 Key Simplifications

### 1. Read-Only MVP
Start with browsing and viewing only. Operations can be stubs that show "Coming soon" messages.

### 2. Reuse Display Logic
Don't rewrite our display modes! Render them into strings and display in viewports.

```go
// Example: Reusing table display in TUI
func renderTableView(archives []*storage.Archive) string {
    var buf bytes.Buffer
    tableDisplay := modes.NewTableDisplay()
    
    // Capture output
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    tableDisplay.Render(archives, display.Options{})
    
    w.Close()
    out, _ := io.ReadAll(r)
    os.Stdout = old
    
    return string(out)
}
```

### 3. Simple State Management
```go
type AppModel struct {
    currentView View        // dashboard, list, or detail
    archives    []*storage.Archive  // Loaded once
    selected    *storage.Archive    // Currently selected
    error       error               // Any error to display
}

type View int
const (
    DashboardView View = iota
    ListView
    DetailView
)
```

### 4. Minimal Navigation
Just three views with simple transitions:
- Dashboard → List (press 'l' or Enter)
- List → Detail (press Enter on item)
- Any → Dashboard (press Esc or 'h')

## 📋 Testing Strategy

### Manual Testing Checklist
```markdown
- [ ] Dashboard loads and displays correctly
- [ ] Can navigate to list view
- [ ] List shows all archives
- [ ] Can switch display modes in list (1,2,3,4)
- [ ] Can select archive and view details
- [ ] Esc returns to previous view
- [ ] Help shows correct shortcuts
- [ ] Scrolling works in all views
- [ ] Terminal resize handled gracefully
```

### Quick Iteration
```go
// Add debug mode for fast development
if os.Getenv("TUI_DEBUG") != "" {
    // Log to file while TUI runs
    logFile, _ := os.Create("/tmp/7zarch-tui.log")
    log.SetOutput(logFile)
}
```

## 🎨 UX Polish Details

### Smooth Transitions
```go
// Use Bubble Tea's commands for smooth updates
func switchToListView() tea.Cmd {
    return func() tea.Msg {
        // Can add slight delay for transition effect
        time.Sleep(50 * time.Millisecond)
        return switchViewMsg{view: ListView}
    }
}
```

### Status Line (Bottom Bar)
```go
// Consistent status line across all views
func statusLine(model AppModel) string {
    left := fmt.Sprintf(" %d archives", len(model.archives))
    center := model.currentView.String()
    right := "? help | q quit "
    
    width := 80 // Get actual terminal width
    padding := width - len(left) - len(center) - len(right)
    
    return lipgloss.JoinHorizontal(
        lipgloss.Left, left,
        lipgloss.Center, strings.Repeat(" ", padding/2) + center,
        lipgloss.Right, right,
    )
}
```

### Loading States
```go
// Use Bubbles spinner for any async operations
type LoadingModel struct {
    spinner  spinner.Model
    message  string
}

func (m LoadingModel) View() string {
    return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}
```

## 🚦 Success Criteria

### MVP Must-Haves
1. ✅ Dashboard view showing system overview
2. ✅ List view with our display modes
3. ✅ Detail view for single archive
4. ✅ Keyboard navigation between views
5. ✅ Help system showing shortcuts

### Nice-to-Haves (If Time Allows)
6. ⭐ Search/filter in list view
7. ⭐ Multi-select for batch operations
8. ⭐ Delete operation with confirmation
9. ⭐ Settings view for preferences
10. ⭐ Color themes

## 💡 Pro Tips for AC

### 1. Start Simple
Get navigation working first. You can always add features.

### 2. Use Bubbles Examples
Charmbracelet has excellent examples. Copy and adapt:
- https://github.com/charmbracelet/bubbletea/tree/master/examples
- https://github.com/charmbracelet/bubbles/tree/master/examples

### 3. Debug Visually
```go
// Add debug info to your views during development
if debug {
    debugInfo := fmt.Sprintf(
        "View: %s | Archives: %d | Selected: %s",
        v.currentView, len(v.archives), v.selected,
    )
    return mainContent + "\n" + subtle.Render(debugInfo)
}
```

### 4. Keep Models Small
Each view should have its own focused model. Don't put everything in AppModel.

### 5. Leverage Our Work
We built beautiful displays. The TUI just needs to show them and add navigation.

## 🎯 Final Architecture

```
┌─────────────────────────────────────────┐
│            TUI App (Bubble Tea)         │
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────┐  ┌──────────┐  ┌──────┐ │
│  │Dashboard │  │   List   │  │Detail│ │
│  │  View    │←→│   View   │←→│ View │ │
│  └──────────┘  └──────────┘  └──────┘ │
│       ↓             ↓            ↓     │
├─────────────────────────────────────────┤
│         Display Modes (Existing)        │
│  Table│Compact│Card│Tree│Dashboard      │
├─────────────────────────────────────────┤
│         Storage Layer (Existing)        │
│      Archive Registry & Operations      │
└─────────────────────────────────────────┘
```

The TUI is a **thin navigation layer** over our existing displays. This is the Charmbracelet way: compose simple, beautiful components.

---

## 🚀 AC: You've Got This!

Start with `cmd/tui.go`:
```go
package cmd

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/adamstac/7zarch-go/internal/tui"
)

var tuiCmd = &cobra.Command{
    Use:   "tui",
    Short: "Launch interactive TUI",
    RunE: func(cmd *cobra.Command, args []string) error {
        app := tui.NewApp()
        p := tea.NewProgram(app, tea.WithAltScreen())
        return p.Start()
    },
}
```

Then build outward. Dashboard first. List second. Detail third. Ship it! 🚢