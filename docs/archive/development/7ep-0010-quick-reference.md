# 7EP-0010 TUI Quick Reference Card

## 🎯 Your Mission
Build a 3-view TUI that wraps our existing display modes with navigation.

## 📦 Dependencies to Add
```go
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
```

## 🏗️ File Structure
```
internal/tui/
├── app.go          # Main model: type App struct { view View; archives []*storage.Archive }
├── dashboard.go    # Dashboard view: wraps modes.DashboardDisplay
├── list.go         # List view: bubbles/list + our display modes  
├── detail.go       # Detail view: viewport + modes.CardDisplay
└── keys.go         # Key bindings: type KeyMap struct { Quit, Help, Enter key.Binding }
```

## ⌨️ Essential Key Bindings
```go
// Global
"q", "ctrl+c" → Quit
"?" → Toggle help
"esc" → Go back

// Dashboard
"l", "enter" → Go to list
"/" → Search (future)

// List  
"enter" → View details
"space" → Multi-select
"1-5" → Switch display modes
"d" → Delete (with confirm)

// Detail
"d" → Delete
"m" → Move
"u" → Upload
"←/→" → Navigate buttons
```

## 🎨 Core Pattern: Wrap Existing Displays
```go
// DON'T rebuild the display logic
// DO capture and show existing output

func renderDashboard(archives []*storage.Archive) string {
    var buf bytes.Buffer
    
    // Create a writer that captures output
    writer := bufio.NewWriter(&buf)
    
    // Use existing display
    display := modes.NewDashboardDisplay()
    display.RenderToWriter(archives, opts, writer)
    
    writer.Flush()
    return buf.String()
}
```

## 🔄 View Switching Pattern
```go
type viewType int
const (
    viewDashboard viewType = iota
    viewList
    viewDetail
)

type App struct {
    view     viewType
    // ... other fields
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case switchViewMsg:
        a.view = msg.view
        return a, nil
    }
    
    // Delegate to current view
    switch a.view {
    case viewDashboard:
        return a.updateDashboard(msg)
    case viewList:
        return a.updateList(msg)
    case viewDetail:
        return a.updateDetail(msg)
    }
}
```

## 📊 Bubbles Components Cheat Sheet

### List Component
```go
import "github.com/charmbracelet/bubbles/list"

// Setup
items := []list.Item{}
for _, archive := range archives {
    items = append(items, archiveItem{archive})
}

l := list.New(items, list.NewDefaultDelegate(), 0, 0)
l.Title = "Archives"
l.SetShowHelp(false)  // We handle help globally
```

### Viewport Component  
```go
import "github.com/charmbracelet/bubbles/viewport"

// Setup
vp := viewport.New(80, 24)
vp.SetContent(contentString)

// Update loop
vp, cmd = vp.Update(msg)

// Keys it handles
// ↑/↓, PgUp/PgDn, Home/End
```

### Help Component
```go
import "github.com/charmbracelet/bubbles/help"

// Setup
h := help.New()
h.ShortSeparator = " • "

// Define keys
keys := []key.Binding{
    key.NewBinding(
        key.WithKeys("q"),
        key.WithHelp("q", "quit"),
    ),
}

// Render
h.View(keys)
```

## 💅 Lipgloss Quick Styles
```go
var (
    // Text styles
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("42"))
    
    // Box with border
    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2)
    
    // Status colors (match our icons)
    okStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))    // Green
    warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Yellow  
    errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))  // Red
    
    // Layout helpers
    centerStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

// Join layouts
lipgloss.JoinHorizontal(lipgloss.Top, left, right)
lipgloss.JoinVertical(lipgloss.Left, top, bottom)
```

## 🚀 Hour-by-Hour Goals

### Hours 1-2: Foundation ✓
- [ ] Basic app structure with 3 view types
- [ ] View switching logic
- [ ] Load archives once on startup

### Hours 3-4: Dashboard ✓  
- [ ] Render existing dashboard display
- [ ] Add viewport for scrolling
- [ ] Navigation to list view

### Hours 5-6: List View ✓
- [ ] Bubbles list with archive items
- [ ] Display mode switching (1-5 keys)
- [ ] Navigation to detail view

### Hours 7-8: Detail View ✓
- [ ] Show full card display
- [ ] Operation buttons (can be stubs)
- [ ] Back navigation

### If Time Allows:
- [ ] Search in list view
- [ ] Multi-select with space
- [ ] Delete confirmation dialog
- [ ] Loading spinners

## 🐛 Debug Tips
```go
// Log to file while TUI runs
func init() {
    if debug := os.Getenv("DEBUG"); debug != "" {
        f, _ := os.Create("/tmp/tui-debug.log")
        log.SetOutput(f)
    }
}

// Then in your code
log.Printf("Current view: %v, Archives: %d", app.view, len(app.archives))
```

## ✅ Definition of Done
1. `7zarch tui` launches dashboard
2. Can navigate: Dashboard → List → Detail → Dashboard
3. List view can switch display modes (1-5)
4. Keyboard shortcuts work (q, ?, esc, enter)
5. Scrolling works in all views
6. No crashes on resize

## 🎯 Remember
- **Simple > Complex**: Get navigation working first
- **Reuse > Rebuild**: Use our existing displays
- **Ship > Perfect**: MVP that works beats perfect that doesn't exist

You got this! 🚀 The Charmbracelet tools are incredible - let them do the heavy lifting.