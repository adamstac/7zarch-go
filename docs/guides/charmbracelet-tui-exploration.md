# 🎨 Charmbracelet TUI Exploration Guide

**Purpose**: Comprehensive guide to explore what's possible with Charmbracelet's TUI ecosystem for 7zarch-go TUI evolution  
**Context**: Building on 7EP-0010 success, planning 7EP-0016 TUI-first interface evolution  
**Current Status**: Basic TUI implemented with 9 themes, simple navigation, and viewport architecture  

## 🌟 The Charmbracelet Ecosystem Overview

### **Core Framework Stack**
- **BubbleTea**: Main TUI framework (Elm Architecture)
- **Lipgloss**: CSS-like styling and layout
- **Bubbles**: Pre-built UI components 
- **Huh**: Interactive forms and prompts
- **Harmonica**: Physics-based animations

### **Architecture Philosophy: The Elm Architecture**
```go
type Model struct {
    // Your application state
}

func (m Model) Init() tea.Cmd {
    // Initial commands (I/O, timers, etc.)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages, update state, return commands
}

func (m Model) View() string {
    // Render current state as string
}
```

**Key Benefits:**
- Predictable state management
- Functional, immutable updates
- Clear separation of concerns
- Excellent debugging experience

## 🧱 BubbleTea Framework Capabilities

### **What BubbleTea Excels At:**
- **Interactive Applications**: Rich keyboard/mouse handling
- **State Management**: Complex application states with predictable updates
- **Multiple Views**: Switching between different interface modes
- **Real-time Updates**: Timers, network requests, file watching
- **Cross-Platform**: Works consistently across terminals
- **Performance**: Efficient rendering and framerate control

### **BubbleTea Patterns for 7zarch-go:**

#### **Multi-View Application Pattern**
```go
type ViewMode int

const (
    ListView ViewMode = iota
    DetailView
    UploadView
    SearchView
    SettingsView
)

type Model struct {
    currentView ViewMode
    listModel   list.Model
    detailModel DetailModel
    // ... other view models
}

func (m Model) View() string {
    switch m.currentView {
    case ListView:
        return m.listModel.View()
    case DetailView:
        return m.detailModel.View()
    // ... other views
    }
}
```

#### **Command-Driven Operations Pattern**
```go
type ArchiveOperationMsg struct {
    Operation string // "upload", "compress", "verify"
    ArchiveID string
    Result    error
}

func uploadArchiveCmd(id string) tea.Cmd {
    return func() tea.Msg {
        err := performUpload(id)
        return ArchiveOperationMsg{
            Operation: "upload",
            ArchiveID: id,
            Result:    err,
        }
    }
}
```

## 🎨 Lipgloss Styling Capabilities

### **What Lipgloss Enables:**
- **CSS-like Styling**: Familiar padding, margins, borders, colors
- **Adaptive Colors**: Automatic light/dark terminal detection
- **Layout System**: Horizontal/vertical joining, alignment, placement
- **Component Styling**: Tables, lists, trees with rich formatting
- **Responsive Design**: Terminal size-aware layouts

### **Lipgloss Patterns for Archive Management:**

#### **Dashboard Layout System**
```go
import "github.com/charmbracelet/lipgloss"

var (
    headerStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("12")).
        Border(lipgloss.RoundedBorder()).
        Padding(0, 1)
    
    sidebarStyle = lipgloss.NewStyle().
        Width(25).
        Border(lipgloss.NormalBorder(), false, true, false, false).
        Padding(1)
    
    contentStyle = lipgloss.NewStyle().
        Flex(1).
        Padding(1)
)

func renderDashboard(sidebar, content string) string {
    body := lipgloss.JoinHorizontal(
        lipgloss.Top,
        sidebarStyle.Render(sidebar),
        contentStyle.Render(content),
    )
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        headerStyle.Render("7ZARCH DASHBOARD"),
        body,
    )
}
```

#### **Archive Status Cards**
```go
func archiveCard(archive Archive) string {
    statusColor := lipgloss.Color("10") // Green
    if archive.Status == "missing" {
        statusColor = lipgloss.Color("9") // Red
    }
    
    cardStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Margin(0, 1, 1, 0).
        Width(30)
    
    statusBadge := lipgloss.NewStyle().
        Background(statusColor).
        Foreground(lipgloss.Color("0")).
        Padding(0, 1).
        Render(archive.Status)
    
    title := lipgloss.NewStyle().Bold(true).Render(archive.Name)
    size := lipgloss.NewStyle().Faint(true).Render(archive.Size)
    
    return cardStyle.Render(
        lipgloss.JoinVertical(
            lipgloss.Left,
            title,
            statusBadge,
            size,
            lipgloss.NewStyle().Faint(true).Render(archive.Created),
        ),
    )
}
```

#### **Progress Visualization**
```go
func renderUploadProgress(progress float64, filename string) string {
    progressBar := lipgloss.NewStyle().
        Width(40).
        Background(lipgloss.Color("236")).
        Render(strings.Repeat("█", int(progress*40)))
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        fmt.Sprintf("Uploading: %s", filename),
        progressBar,
        fmt.Sprintf("%.1f%% complete", progress*100),
    )
}
```

## 🧩 Bubbles Component Library

### **Available Components & Archive Tool Applications:**

#### **1. List Component - Perfect for Archive Browsing**
```go
import "github.com/charmbracelet/bubbles/list"

type ArchiveItem struct {
    title, desc, uid string
    size            string
    status          string
}

func (i ArchiveItem) FilterValue() string { return i.title }
func (i ArchiveItem) Title() string       { return i.title }
func (i ArchiveItem) Description() string { 
    return fmt.Sprintf("%s • %s • %s", i.size, i.status, i.uid[:8]) 
}

// Features:
// - Fuzzy filtering ("/podcast" finds "podcast-423.7z")
// - Pagination for large collections
// - Custom styling and themes
// - Built-in help system
// - Status messages and loading states
```

#### **2. Table Component - Detailed Archive Views**
```go
import "github.com/charmbracelet/bubbles/table"

func createArchiveTable(archives []Archive) table.Model {
    columns := []table.Column{
        {Title: "Name", Width: 30},
        {Title: "Size", Width: 10},
        {Title: "Status", Width: 10},
        {Title: "Profile", Width: 12},
        {Title: "Created", Width: 12},
    }
    
    rows := make([]table.Row, len(archives))
    for i, a := range archives {
        rows[i] = table.Row{
            a.Name,
            a.Size,
            statusWithIcon(a.Status),
            a.Profile,
            a.Created.Format("2006-01-02"),
        }
    }
    
    return table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(20),
    )
}
```

#### **3. Progress Component - Upload/Compression Feedback**
```go
import "github.com/charmbracelet/bubbles/progress"

type OperationModel struct {
    progress progress.Model
    operation string
    filename  string
}

func (m OperationModel) View() string {
    return fmt.Sprintf(
        "%s: %s\n%s", 
        m.operation, 
        m.filename,
        m.progress.View(),
    )
}
```

#### **4. Spinner Component - Background Operations**
```go
import "github.com/charmbracelet/bubbles/spinner"

func compressionSpinner() spinner.Model {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    return s
}
```

#### **5. Viewport Component - Content Scrolling**
```go
import "github.com/charmbracelet/bubbles/viewport"

// Perfect for:
// - Archive content previews
// - Log viewing
// - Configuration displays
// - Help documentation

func createLogViewer(content string) viewport.Model {
    vp := viewport.New(80, 20)
    vp.SetContent(content)
    vp.Style = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1)
    return vp
}
```

#### **6. Text Input Components - User Input**
```go
import (
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/textarea"
)

// Text Input for names, paths, queries
func createSearchInput() textinput.Model {
    ti := textinput.New()
    ti.Placeholder = "Search archives..."
    ti.Focus()
    ti.CharLimit = 100
    ti.Width = 50
    return ti
}

// Textarea for descriptions, metadata
func createDescriptionInput() textarea.Model {
    ta := textarea.New()
    ta.Placeholder = "Archive description..."
    ta.SetWidth(50)
    ta.SetHeight(5)
    return ta
}
```

## 📋 Huh Forms for Configuration & Setup

### **Interactive Setup Wizards**
```go
import "github.com/charmbracelet/huh"

func configurationWizard() error {
    var (
        storageLocation string
        trueNASEnabled  bool
        defaultProfile  string
        compressionLevel int
    )
    
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Default Storage Location").
                Options(
                    huh.NewOption("Managed (~/.7zarch-go/archives/)", "managed"),
                    huh.NewOption("Custom Directory", "custom"),
                ).
                Value(&storageLocation),
                
            huh.NewConfirm().
                Title("Enable TrueNAS Integration?").
                Description("Connect to TrueNAS for remote storage").
                Value(&trueNASEnabled),
        ),
        
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Default Compression Profile").
                Options(
                    huh.NewOption("Balanced (recommended)", "balanced"),
                    huh.NewOption("Media (photos/videos)", "media"),
                    huh.NewOption("Documents (text/code)", "documents"),
                    huh.NewOption("Maximum Compression", "maximum"),
                ).
                Value(&defaultProfile),
                
            huh.NewSelect[int]().
                Title("Compression Level (1-9)").
                Options(
                    huh.NewOption("Fast (1)", 1),
                    huh.NewOption("Balanced (5)", 5),
                    huh.NewOption("Maximum (9)", 9),
                ).
                Value(&compressionLevel),
        ),
    )
    
    return form.Run()
}
```

### **Archive Creation Forms**
```go
func archiveCreationForm() (*ArchiveConfig, error) {
    var config ArchiveConfig
    
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Archive Name").
                Value(&config.Name).
                Validate(func(str string) error {
                    if str == "" {
                        return errors.New("Name is required")
                    }
                    return nil
                }),
                
            huh.NewMultiSelect[string]().
                Title("Source Directories").
                Options(getAvailableDirectories()...).
                Value(&config.Sources),
                
            huh.NewSelect[string]().
                Title("Compression Profile").
                Options(
                    huh.NewOption("Media (photos, videos)", "media"),
                    huh.NewOption("Documents (text, code)", "documents"),
                    huh.NewOption("Balanced (mixed content)", "balanced"),
                ).
                Value(&config.Profile),
        ),
        
        huh.NewGroup(
            huh.NewConfirm().
                Title("Upload to TrueNAS after creation?").
                Value(&config.AutoUpload),
                
            huh.NewConfirm().
                Title("Verify archive integrity?").
                Value(&config.Verify),
        ),
    )
    
    err := form.Run()
    return &config, err
}
```

## 🎭 Advanced TUI Patterns for 7zarch-go

### **1. Multi-Panel Dashboard**
```go
type DashboardModel struct {
    activePanel    int
    archiveList    list.Model
    statusPanel    StatusModel
    operationPanel OperationModel
    helpPanel      help.Model
}

func (m DashboardModel) View() string {
    leftPanel := lipgloss.NewStyle().Width(40).Render(m.archiveList.View())
    rightTop := lipgloss.NewStyle().Height(10).Render(m.statusPanel.View())
    rightBottom := lipgloss.NewStyle().Render(m.operationPanel.View())
    
    rightPanel := lipgloss.JoinVertical(
        lipgloss.Left,
        rightTop,
        rightBottom,
    )
    
    body := lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftPanel,
        rightPanel,
    )
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        headerStyle.Render("7ZARCH DASHBOARD"),
        body,
        m.helpPanel.View(m.keyMap()),
    )
}
```

### **2. Command Palette (Vim-style)**
```go
type CommandPaletteModel struct {
    active   bool
    input    textinput.Model
    commands []Command
    filtered []Command
}

func (m CommandPaletteModel) Update(msg tea.Msg) (CommandPaletteModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == ":" && !m.active {
            m.active = true
            m.input.Focus()
            return m, textinput.Blink
        }
        if msg.String() == "esc" && m.active {
            m.active = false
            m.input.Blur()
            m.input.SetValue("")
            return m, nil
        }
        if m.active {
            m.input, _ = m.input.Update(msg)
            m.filterCommands()
        }
    }
    return m, nil
}

// Commands like:
// :create podcast-424 - Create new archive
// :upload 01JEX      - Upload archive by ID
// :search podcasts   - Search archives
// :config            - Open configuration
```

### **3. Real-time Status Updates**
```go
type StatusUpdates struct {
    uploadProgress   map[string]float64
    compressionJobs  []CompressionJob
    systemHealth     HealthStatus
}

func statusUpdateCmd() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return StatusUpdateMsg{
            timestamp: t,
            updates:   getLatestStatus(),
        }
    })
}
```

### **4. Interactive Archive Browser**
```go
type BrowserModel struct {
    currentPath   string
    entries       []ArchiveEntry
    selectedEntry int
    preview       viewport.Model
    breadcrumbs   []string
}

func (m BrowserModel) View() string {
    breadcrumb := strings.Join(m.breadcrumbs, " > ")
    
    entryList := make([]string, len(m.entries))
    for i, entry := range m.entries {
        icon := "📁"
        if !entry.IsDir {
            icon = "📄"
        }
        
        style := lipgloss.NewStyle()
        if i == m.selectedEntry {
            style = style.Background(lipgloss.Color("62"))
        }
        
        entryList[i] = style.Render(fmt.Sprintf("%s %s", icon, entry.Name))
    }
    
    leftPanel := lipgloss.JoinVertical(lipgloss.Left, entryList...)
    rightPanel := m.preview.View()
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        breadcrumbStyle.Render(breadcrumb),
        lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel),
    )
}
```

## 🚀 Implementation Roadmap for 7EP-0016

### **Phase 1: Enhanced Navigation (Building on 7EP-0010)**
- Multi-view switching (List ↔ Detail ↔ Settings)
- Command palette with `:` prefix
- Improved keyboard shortcuts
- Contextual help system

### **Phase 2: Rich Data Visualization** 
- Dashboard with multiple panels
- Storage usage charts and graphs
- Archive health indicators
- Tabular archive listings

### **Phase 3: Interactive Operations**
- Real-time upload/compression progress
- Interactive archive creation forms
- Live search with filtering
- Configuration wizards

### **Phase 4: Advanced Features**
- Remote TrueNAS browsing
- Responsive layouts for different terminal sizes
- Advanced themes and customization
- Analytics and insights

## 💡 Key Technical Considerations

### **Performance Optimization**
- Use `viewport` for large lists/content
- Implement lazy loading for archive metadata
- Optimize rendering with selective updates
- Cache expensive operations (file system access)

### **User Experience Principles**
- **Consistent Navigation**: Arrow keys, Enter, Esc patterns
- **Progressive Disclosure**: Simple → Advanced features
- **Immediate Feedback**: Loading states, progress indicators
- **Error Recovery**: Clear error messages with recovery options

### **Integration Points**
- **Registry Operations**: Real-time database updates
- **File System Watching**: Live status updates
- **Network Operations**: TrueNAS integration feedback
- **Background Jobs**: Compression, uploads with progress

## 🎯 Success Metrics for TUI Evolution

### **User Adoption**
- Users prefer TUI over CLI for archive management
- Session duration increases (users explore more)
- Feature discovery improves (forms guide users)
- Error rates decrease (better feedback)

### **Technical Excellence**
- Responsive performance (< 100ms interactions)
- Memory efficiency with large archive collections
- Cross-platform compatibility
- Graceful degradation on limited terminals

## 🔗 Learning Resources

### **Essential Documentation**
- [BubbleTea Tutorial](https://github.com/charmbracelet/bubbletea#tutorial)
- [Lipgloss Examples](https://github.com/charmbracelet/lipgloss#example)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Huh Forms Guide](https://github.com/charmbracelet/huh#tutorial)

### **Reference Implementations**
- [Glow](https://github.com/charmbracelet/glow) - Markdown viewer
- [Soft Serve](https://github.com/charmbracelet/soft-serve) - Git server TUI
- [Pop](https://github.com/charmbracelet/pop) - Email TUI
- [VHS](https://github.com/charmbracelet/vhs) - Terminal recording

This comprehensive exploration guide provides the foundation for understanding what's possible with Charmbracelet's TUI ecosystem and how it can transform 7zarch-go into a compelling terminal-first archive management experience.
