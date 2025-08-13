# 7EP-0010: Interactive TUI Application

**Status:** ✅ Implemented  
**Author(s):** Claude Code (CC), Adam Stacoviak, Amp (Sourcegraph)  
**Assignment:** CC (Full Implementation) + Amp (UX Design & Oversight)  
**Difficulty:** 3 (moderate - builds on existing PR21 foundation with theme system)  
**Created:** 2025-08-12  
**Updated:** 2025-08-13 (Amp UX design collaboration + foundation strategy)  
**Implementation:** 2025-08-13 (PR #21 merged - Simple TUI with theme system)  

## Executive Summary

Create a simple, elegant Terminal User Interface (TUI) focused on **podcast archival workflows** with clean list/details navigation and beautiful themes. Build on existing PR21 foundation while simplifying the interface to prioritize ease of viewing, creating, and managing large media archives for TrueNAS storage with future Tailscale integration.

## ✅ Implementation Summary (PR #21)

**Successfully Delivered:**
- **9 Beautiful Themes**: Blue Tech, Terminal Green, Purple, Cyan, Charmbracelet + 4 Dracula variations
- **Simple Navigation**: Arrow keys, Enter for details, Space for multi-select, single letters for actions
- **Multiple Entry Points**: `browse` (primary), `ui`, `i`, `tui` (compatibility)
- **Professional Architecture**: Viewport framework with proper margins and responsive design
- **Podcast Workflow Optimized**: Large media file friendly, staging → remote storage ready

**User Experience:**
- **List View**: Clean archive listing with themed highlighting and status indicators
- **Details View**: Essential metadata display with theme consistency
- **Multi-Select**: Space bar selection with visual checkboxes
- **Actions**: Delete, move, upload with confirmation dialogs
- **Theme Selection**: `--theme` flag with 9 gorgeous options

**Technical Foundation:**
- **Bubbletea Architecture**: Professional TUI framework with viewport component
- **Theme System**: Complete color specification system with runtime switching
- **Responsive Design**: Automatic terminal size adaptation
- **Safety Patterns**: Confirmation dialogs for destructive operations

**Future Evolution**: Architecture designed to support embedded command line (7EP-0016)

## 🎯 Amp (Sourcegraph) UX Design Vision

**Core Philosophy:** Simple lists with selection and action - optimized for podcast production archival workflow.

### **Primary Use Case: Podcast Archive Management**
**User:** Podcast producer managing large episode archives  
**Workflow:** Create → Stage → Upload to TrueNAS → Browse remote → Download when needed  
**Focus:** Large media files, long-term storage, simple navigation, bulletproof safety

### **UX Principles:**
- **List-first interface** - Simple archive lists with arrow key navigation
- **Immediate actions** - Single key commands for common operations  
- **Theme system** - Multiple color schemes including Dracula
- **Safety-first** - Soft deletes with clear confirmation patterns
- **Remote integration** - TrueNAS feels local via Tailscale (future)

## Evidence & Reasoning

**Current CLI limitations for complex workflows:**
- Multi-step operations require multiple command invocations
- No visual feedback for long-running operations (create, upload, move)
- Difficult to understand archive relationships and organization patterns
- Batch operations require scripting or multiple manual commands
- No interactive exploration of archive contents or metadata
- Path selection requires typing full paths or using separate tools

**TUI advantages for archive management:**
- **Visual overview**: Dashboard showing storage health, recent activity, problem areas
- **Interactive selection**: Multi-select archives for batch operations with visual feedback
- **Real-time operations**: Progress bars, status updates, live directory browsing
- **Exploration workflows**: Navigate archive hierarchies, preview contents, discover relationships
- **Reduced cognitive load**: Visual organization reduces need to remember commands and paths
- **Path assistance**: Built-in file browser with tab completion and directory navigation

**Why a comprehensive TUI matters:**
- Archive management is inherently visual (files, directories, sizes, relationships)
- Complex workflows benefit from persistent state and navigation
- Modern users expect rich terminal interfaces (see `lazygit`, `k9s`, `htop`)
- TUI complements CLI without replacing it for automation/scripting

## Use Cases

### Primary Use Case: Simple List Navigation for Podcast Archives
```bash
# Launch TUI application  
7zarch-go tui

# Simple list interface:
# - Arrow keys navigate archives
# - Enter for details view
# - Single letters for actions (d=delete, m=move, u=upload)
# - Multi-select with Space bar
# - Theme selection for personal preference
```

**Core Interface Pattern:**
```
7zarch-go                                    # Theme-colored header

Archives: 247 (12.4 TB)                     # Status line

> episode-423.7z        89 MB    2h ago    ✓ # Selected item highlighted
  episode-422.7z        92 MB    1d ago    ✓ # Normal items
  vacation-photos.7z   3.8 GB    1w ago    ✓
  
[Enter] Details  [d] Delete  [u] Upload  [q] Quit
```

**Simple List Interface Examples:**

**Blue Tech Theme:**
```
7zarch-go                                                     # Bright Blue

Archives: 247 (12.4 TB)                                       # White

> episode-423.7z                89 MB    2h ago    ✓           # Blue bg, White text
  episode-422.7z                92 MB    1d ago    ✓           # White, Cyan, Green
  vacation-photos.7z           3.8 GB    1w ago    ✓           # White, Cyan, Green
  old-project.7z               456 MB    1m ago    ?           # White, Cyan, Yellow

[Enter] Details  [d] Delete  [m] Move  [u] Upload  [q] Quit    # Blue
```

**Dracula Theme:**
```
7zarch-go                                                     # Purple (#BD93F9)

Archives: 247 (12.4 TB)                                       # Foreground (#F8F8F2)

> episode-423.7z                89 MB    2h ago    ✓           # Selection bg (#44475A), Purple text
  episode-422.7z                92 MB    1d ago    ✓           # Foreground, Cyan (#8BE9FD), Green (#50FA7B)
  vacation-photos.7z           3.8 GB    1w ago    ✓           # Foreground, Cyan (#8BE9FD), Green (#50FA7B)
  old-project.7z               456 MB    1m ago    ?           # Foreground, Cyan (#8BE9FD), Yellow (#F1FA8C)

[Enter] Details  [d] Delete  [m] Move  [u] Upload  [q] Quit    # Pink (#FF79C6)
```

### Secondary Use Cases

#### Interactive Archive Browsing
```
Archive Browser - 127 archives found                        [Filter: all] [Sort: date]
┌─ Filters ──────────────────┬─ Archive List ─────────────────────────────────────────────┐
│                            │                                                            │
│ Status:                    │ > 📦 project-backup-2024-08-12.7z          [01HF5K2M] ✅  │
│ ○ All archives (127)       │   💾 1.2 GB • 📊 balanced • ⏰ 2h ago                     │
│ ○ Active only (119)        │                                                            │
│ ○ Missing only (6)         │   📦 documents-sync-daily.7z               [01HF5K1N] ✅  │
│ ○ Deleted only (2)         │   💾 456 MB • 📊 documents • ⏰ 4h ago                    │
│                            │                                                            │
│ Profile:                   │   📦 media-vacation-photos.7z              [01HF5K0P] ✅  │
│ ○ All profiles             │   💾 3.8 GB • 📊 media • ⏰ 6h ago                        │
│ ○ Media (45)               │                                                            │
│ ○ Documents (38)           │ ⚠️ 📦 backup-external-drive.7z             [01HF5J9Q] ❌  │
│ ○ Balanced (44)            │   💾 2.1 GB • 📊 balanced • ⏰ 2d ago • Missing          │
│                            │                                                            │
│ Size Range:                │   📦 old-photos-archive.7z                 [01HF5J8R] ✅  │
│ Min: [______] MB           │   💾 892 MB • 📊 media • ⏰ 3d ago                        │
│ Max: [______] MB           │                                                            │
│                            │ ┌─ Selection Actions ─────────────────────────────────────┐ │
│ Age:                       │ │ Space: Select   Enter: Details   d: Delete             │ │
│ ○ All                      │ │ m: Move         u: Upload        c: Create similar     │ │
│ ○ Last 24h                 │ │ /: Search       f: Filter        Tab: Switch panes    │ │
│ ○ Last week                │ └─────────────────────────────────────────────────────────┘ │
│ ○ Last month               │                                                            │
│ ○ Older                    │ 📊 Showing 127 archives • 8.4 TB total                   │
│                            │                                                            │
│ 🔍 Search:                 │                                                            │
│ [________________]         │                                                            │
└────────────────────────────┴────────────────────────────────────────────────────────────┘
Status: project-backup-2024-08-12.7z selected | Use arrow keys, Space to select, Enter for details
```

#### Interactive Archive Creation with Path Browser
```
Create New Archive                                                     [Esc: Cancel]
┌─ Source Selection ────────────────────────────────────────────────────────────────────┐
│ 📁 Current Path: /Users/adam/Documents/Projects/                                      │
│                                                                     [Tab: Complete]   │
│ ┌─ Directory Browser ─────────────────────────────────────────────────────────────┐  │
│ │ 📁 Parent: /Users/adam/Documents/                                               │  │
│ │                                                                                 │  │
│ │ > 📁 Projects/          (current)                                              │  │
│ │   📁 client-work/       [Select: Space] [Enter: Browse]                       │  │
│ │   📁 open-source/       2.3 GB • 1,456 files                                  │  │
│ │   📁 personal/          156 MB • 89 files                                      │  │
│ │   📁 archive-old/       4.2 GB • 2,891 files                                  │  │
│ │   📄 README.md          4 KB                                                   │  │
│ │   📄 .gitignore         156 bytes                                              │  │
│ │                                                                                 │  │
│ │ Selected for archiving:                                                         │  │
│ │ ✓ client-work/ (2.3 GB)                                                       │  │
│ │ ✓ open-source/ (156 MB)                                                       │  │
│ │   Total: 2.5 GB • 1,545 files                                                  │  │
│ └─────────────────────────────────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────────────────────────────────┘
┌─ Archive Configuration ───────────────────────────────────────────────────────────────┐
│ Name: [client-work-backup-2024-08-12________________] (.7z will be added)              │
│                                                                                       │
│ Profile: ○ Media     ○ Documents     ● Balanced     ○ Custom                         │
│         (photos,     (office docs,   (mixed files,  (configure                       │
│          videos)     code, text)     general use)   manually)                        │
│                                                                                       │
│ Destination: ● Managed Storage (/managed/archives/2024/08/)                          │
│              ○ Custom Path: [________________________________] [Browse]              │
│                                                                                       │
│ Options: [✓] Create checksum     [✓] Add to registry     [ ] Upload after creation  │
│          [ ] Exclude hidden      [ ] Follow symlinks     [ ] Verify after creation  │
│                                                                                       │
│ 📊 Estimated: ~800MB compressed (68% reduction) • ~5-8 minutes to create            │
└───────────────────────────────────────────────────────────────────────────────────────┘
[Tab: Navigate fields] [Space: Toggle] [Enter: Create archive] [Esc: Cancel] [F1: Help]
```

#### Live Search with Real-Time Results
```
Live Archive Search                                           [Type to search] [Esc: Exit]
┌─ Search Query ────────────────────────────────────────────────────────────────────────┐
│ project backup 2024 photos█                                              [Clear: ^L] │
└───────────────────────────────────────────────────────────────────────────────────────┘
┌─ Filters ─────────────────┬─ Live Results (4 found) ────────────────────────────────────┐
│                           │                                                             │
│ Search in:                │ > 📦 project-backup-2024-08-12.7z         [01HF5K2M] ✅   │
│ ✓ Archive names           │   🔍 Match: name (project, backup, 2024)                   │
│ ✓ File paths              │   💾 1.2 GB • 📊 balanced • ⏰ 2h ago • Managed           │
│ ✓ Metadata                │   📂 /managed/archives/2024/08/project-backup-2024-08-...  │
│ ○ File contents           │                                                             │
│                           │   📦 client-project-photos-2024.7z        [01HF5J3K] ✅   │
│ Match type:               │   🔍 Match: name (project, photos, 2024)                   │
│ ○ Exact phrase            │   💾 3.2 GB • 📊 media • ⏰ 1w ago • External              │
│ ● All words               │   📂 /external/projects/photos/client-project-photos-...   │
│ ○ Any word                │                                                             │
│ ○ Regex                   │   📦 backup-photos-vacation-2024.7z       [01HF5H7L] ⚠️   │
│                           │   🔍 Match: name, path (backup, photos, 2024)              │
│ Status:                   │   💾 2.8 GB • 📊 media • ⏰ 3w ago • Missing               │
│ ○ Any status              │   📂 /external/vacation/backup-photos-vacation-2024.7z     │
│ ✓ Present archives        │                                                             │
│ ✓ Missing archives        │   📦 project-docs-backup.7z               [01HF5G9M] ✅   │
│ ○ Deleted archives        │   🔍 Match: metadata (contains "2024 project files")       │
│                           │   💾 456 MB • 📊 documents • ⏰ 1m ago • Managed          │
│ Profile:                  │   📂 /managed/archives/2024/08/project-docs-backup.7z      │
│ ○ All profiles            │                                                             │
│ ✓ Media                   │ ┌─ Search Actions ────────────────────────────────────────┐ │
│ ○ Documents               │ │ ↑/↓: Navigate    Enter: Show details    Space: Select  │ │
│ ✓ Balanced                │ │ Tab: Advanced    F3: Find next          /: New search  │ │
│                           │ └─────────────────────────────────────────────────────────┘ │
└───────────────────────────┴─────────────────────────────────────────────────────────────┘
Search updates as you type • Use Tab for advanced filters • Press Enter for details
```

#### Batch Operations with Visual Progress
```
Batch Move Operation - 5 archives selected                               [Esc: Cancel]
┌─ Selected Archives (5) ──────────────────────┬─ Destination Browser ─────────────────┐
│                                              │                                        │
│ ✓ project-backup-2024-08-12.7z              │ 📁 /external/backups/projects/         │
│   💾 1.2 GB • balanced • Currently: managed │ ┌─ Path Navigation ─────────────────┐  │
│                                              │ │ 📁 / (root)                       │  │
│ ✓ client-project-photos-2024.7z             │ │ 📁 external/                      │  │
│   💾 3.2 GB • media • Currently: external   │ │ 📁 backups/                       │  │
│                                              │ │ 📁 projects/ (current)            │  │
│ ✓ old-project-archive.7z                    │ └───────────────────────────────────┘  │
│   💾 892 MB • balanced • Currently: managed │                                        │
│                                              │ 📁 Subdirectories:                    │
│ ✓ project-docs-backup.7z                    │ > 📁 2024/                            │
│   💾 456 MB • documents • Currently: managed│   📁 2023/                            │
│                                              │   📁 client-work/                     │
│ ✓ media-project-files.7z                    │   📁 archive/                         │
│   💾 2.1 GB • media • Currently: external   │                                        │
│                                              │ 📝 Custom path:                       │
│ Total: 5 archives • 7.9 GB                  │ /external/backups/projects/2024/█     │
│                                              │ [Tab: Browse] [Enter: Confirm]        │
│ ┌─ Move Options ──────────────────────────┐  │                                        │
│ │ [✓] Update registry entries             │  │ 🎯 Final destination:                 │
│ │ [✓] Verify files after move             │  │ /external/backups/projects/2024/      │
│ │ [✓] Create destination if missing       │  │                                        │
│ │ [ ] Copy instead of move                │  │ ⚠️ Will change 3 archives from        │
│ │ [ ] Create date-based subdirectories    │  │    managed → external status          │
│ └─────────────────────────────────────────┘  │                                        │
└──────────────────────────────────────────────┴────────────────────────────────────────┘
[Tab: Switch panes] [Enter: Start move] [Space: Toggle options] [Esc: Cancel] [?: Help]

# During operation:
┌─ Moving Archives ────────────────────────────────────────────────────────────────────┐
│ Progress: 2 of 5 completed (40%)                                            [Cancel] │
│                                                                                      │
│ ✅ project-backup-2024-08-12.7z     → moved successfully (1.2 GB)                   │
│ ✅ client-project-photos-2024.7z    → moved successfully (3.2 GB)                   │
│ 🔄 old-project-archive.7z           → moving... ████████░░░░░░░░ 65% (580/892 MB)   │
│ ⏳ project-docs-backup.7z           → pending                                        │
│ ⏳ media-project-files.7z           → pending                                        │
│                                                                                      │
│ Elapsed: 2m 34s • Estimated remaining: 1m 15s • Speed: 45.2 MB/s                   │
└──────────────────────────────────────────────────────────────────────────────────────┘
```

#### Archive Details with Interactive Metadata
```
Archive Details - project-backup-2024-08-12.7z                    [01HF5K2M] [Esc: Back]
┌─ Overview ────────────────────────────────────────────────────────────────────────────┐
│ 📦 Name: project-backup-2024-08-12.7z                                                │
│ 📏 Size: 1.2 GB (compressed from 3.4 GB - 65% compression)                          │
│ 📊 Profile: balanced • 🏷️ Status: ✅ Present • 📍 Type: Managed                     │
│ ⏰ Created: 2024-08-12 14:30:15 (2 hours ago)                                        │
│ 📂 Path: /managed/archives/2024/08/project-backup-2024-08-12.7z                     │
│ 🔐 Checksum: verified ✅ • 📋 Files: 1,247 files in 89 directories                  │
└───────────────────────────────────────────────────────────────────────────────────────┘
┌─ Actions ─────────────────────────────────────────────────────────────────────────────┐
│ v: Verify integrity    e: Extract files     m: Move archive      d: Delete archive   │
│ u: Upload to cloud     c: Create copy       r: Rename archive    i: Show metadata    │
│ l: List contents       s: Search in files   h: Show history      t: Add tags         │
└───────────────────────────────────────────────────────────────────────────────────────┘
┌─ File Contents (10 of 1,247 shown) ─┬─ Metadata & History ──────────────────────────┐
│                                      │                                               │
│ > 📁 src/                           │ 📊 Compression Details:                      │
│   📁 components/                     │ • Algorithm: LZMA2                           │
│   📁 utils/                          │ • Dictionary: 32MB                            │
│   📁 tests/                          │ • Threads: 8                                  │
│   📄 main.go              125 KB     │ • Solid archive: Yes                         │
│   📄 config.yaml          2.4 KB     │                                               │
│                                      │ 🏷️ Tags:                                     │
│ > 📁 docs/                          │ • client-work                                │
│   📄 README.md            8.9 KB     │ • backup                                     │
│   📄 API.md               45 KB      │ • project                                    │
│                                      │ • 2024-q3                                    │
│ > 📁 assets/                        │                                               │
│   📄 logo.png             234 KB     │ 📈 History:                                  │
│   📄 banner.jpg           1.2 MB     │ • Created: 2024-08-12 14:30                 │
│                                      │ • Last verified: 2024-08-12 14:31           │
│ ┌─ Content Actions ─────────────────┐ │ • Registry updated: 2024-08-12 14:31        │
│ │ Enter: Browse folder              │ │ • No move operations                         │
│ │ Space: Select files               │ │ • No upload history                          │
│ │ e: Extract selected               │ │                                               │
│ │ /: Search in archive              │ │ 🔗 Related Archives:                         │
│ └───────────────────────────────────┘ │ • project-docs-backup.7z (same source)      │
│                                      │ • client-work-archive.7z (similar content)   │
└──────────────────────────────────────┴───────────────────────────────────────────────┘
Use arrow keys to navigate, Tab to switch panes, Enter to explore folders
```

## Technical Design

### TUI Architecture Overview

#### Core Application Structure
```go
// Main TUI application using Bubble Tea framework
type App struct {
    state      AppState
    manager    *storage.Manager
    navigation Navigation
    views      map[ViewType]View
    theme      Theme
}

type AppState struct {
    CurrentView    ViewType
    ViewHistory    []ViewType
    SelectedItems  []SelectedItem
    GlobalFilters  FilterState
    UserPrefs      UserPreferences
}

type ViewType int
const (
    DashboardView ViewType = iota
    BrowserView
    SearchView
    CreateView
    DetailsView
    BatchOpsView
    SettingsView
)
```

#### Navigation System
```go
type Navigation struct {
    keyMap      KeyMap
    breadcrumbs []NavigationCrumb
    shortcuts   map[rune]ViewType
}

type KeyMap struct {
    Global    GlobalKeys
    ViewLocal map[ViewType]LocalKeys
}

// Global keyboard shortcuts
type GlobalKeys struct {
    Quit        key.Binding  // q
    Help        key.Binding  // ?
    Dashboard   key.Binding  // h (home)
    Browse      key.Binding  // l (list)
    Search      key.Binding  // /
    Create      key.Binding  // c
    Settings    key.Binding  // ,
    Back        key.Binding  // Esc
    Forward     key.Binding  // Tab
}
```

#### View System with Bubble Tea Components
```go
// Base view interface
type View interface {
    Init() tea.Cmd
    Update(tea.Msg) (View, tea.Cmd)
    View() string
    HandleKey(tea.KeyMsg) tea.Cmd
    Focus() tea.Cmd
    Blur() tea.Cmd
}

// Dashboard view - entry point
type DashboardView struct {
    storageStats  StorageStatsModel
    recentFiles   RecentFilesModel
    quickActions  QuickActionsModel
    healthCheck   HealthCheckModel
    focused       int
}

// Browser view - main archive listing
type BrowserView struct {
    filterPanel   FilterPanelModel
    archiveList   ArchiveListModel
    detailPanel   DetailPanelModel
    pagination    PaginationModel
    selection     SelectionModel
}

// Search view - live search interface
type SearchView struct {
    searchInput   SearchInputModel
    filterPanel   FilterPanelModel
    resultsList   ResultsListModel
    previewPanel  PreviewPanelModel
}

// Create view - interactive archive creation
type CreateView struct {
    pathBrowser   PathBrowserModel
    configPanel   ConfigPanelModel
    previewPanel  PreviewPanelModel
    progressBar   ProgressModel
}
```

#### Path Browser with Tab Completion
```go
type PathBrowserModel struct {
    currentPath   string
    entries       []DirEntry
    selected      map[string]bool
    cursor        int
    completion    CompletionEngine
}

type CompletionEngine struct {
    cache         map[string][]string
    maxResults    int
    caseSensitive bool
}

func (pe *CompletionEngine) Complete(partial string) []string {
    // Implement tab completion for file paths
    // Support glob patterns, fuzzy matching
    // Cache results for performance
    // Handle permissions and access errors gracefully
}

// Path input with real-time completion
type PathInputModel struct {
    input         textinput.Model
    suggestions   []string
    showSuggestions bool
    selectedSuggestion int
}
```

#### Multi-Select and Batch Operations
```go
type SelectionModel struct {
    items         []SelectableItem
    selected      map[string]bool
    multiSelect   bool
    lastSelected  int
}

type BatchOperationModel struct {
    operation     BatchOperationType
    targets       []storage.Archive
    progress      ProgressTracker
    options       BatchOptions
    confirmation  ConfirmationModel
}

type BatchOperationType int
const (
    BatchMove BatchOperationType = iota
    BatchDelete
    BatchUpload
    BatchVerify
    BatchTag
)
```

### TUI Framework Analysis

#### 1. Bubble Tea (Recommended)
**Pros:**
- Modern reactive architecture based on Elm/Redux patterns
- Excellent composability for complex UIs
- Active development and strong community
- Built-in support for modern terminal features
- Comprehensive input handling and key bindings
- Good performance and memory management
- Extensive documentation and examples

**Cons:**
- Newer framework (less mature than tview)
- Steeper learning curve for complex applications
- Some widgets require custom implementation

**Best for:** Modern, complex TUI applications with rich interactions

#### 2. tview
**Pros:**
- Mature, stable framework with extensive widget library
- Rich set of built-in components (tables, forms, trees, modals)
- Simple programming model
- Good documentation and examples
- Proven in production applications

**Cons:**
- Less flexible for custom layouts
- Harder to achieve modern UI patterns
- Limited animation and transition support
- Monolithic design can limit customization

**Best for:** Traditional TUI applications with standard widgets

#### 3. termui
**Pros:**
- Excellent for dashboard-style interfaces
- Built-in charting and visualization widgets
- Good for data-heavy applications
- Simple API for basic use cases

**Cons:**
- Limited interactive capabilities
- Less suitable for form-heavy applications
- Smaller community and development activity
- Limited layout flexibility

**Best for:** Dashboard and monitoring applications

#### 4. tcell (Low-level)
**Pros:**
- Maximum control over terminal capabilities
- Excellent performance
- Direct access to all terminal features
- Small footprint

**Cons:**
- Requires building everything from scratch
- Significant development overhead
- Complex input handling
- No built-in widgets or layout system

**Best for:** High-performance applications with specific requirements

### Recommended Technology Stack

**Primary Framework: Bubble Tea**
- Modern architecture aligns with complex TUI requirements
- Excellent composability for our multi-view application
- Strong community support and active development
- Built-in support for animations and rich interactions

**Supporting Libraries:**
```go
// Core TUI framework
"github.com/charmbracelet/bubbletea"

// UI components and styling
"github.com/charmbracelet/lipgloss"    // Styling and layout
"github.com/charmbracelet/bubbles"     // Common UI components

// Enhanced input handling
"github.com/atotto/clipboard"          // Clipboard integration
"github.com/mattn/go-runewidth"        // Unicode width handling

// File system operations
"github.com/fsnotify/fsnotify"         // File system watching
"github.com/spf13/afero"               // File system abstraction

// Performance and utilities
"github.com/dustin/go-humanize"        // Human-readable sizes/times
"github.com/schollz/progressbar/v3"    // Progress bars
```

### User Experience Design

#### Dashboard-First Philosophy
**Entry Experience Design:**
1. **Immediate value**: Show storage overview, health, recent activity
2. **Problem awareness**: Highlight missing archives, space issues, maintenance needs
3. **Quick navigation**: Single-key shortcuts to common functions
4. **Visual hierarchy**: Use space, color, and typography to guide attention
5. **Progressive disclosure**: Start simple, reveal complexity as needed

#### Keyboard-First Navigation
**Primary Navigation:**
- `h` - Dashboard (home)
- `l` - Browse archives (list)
- `/` - Search
- `c` - Create archive
- `m` - Missing archives
- `t` - Trash management
- `u` - Upload operations
- `?` - Help
- `q` - Quit

**Context-Sensitive Keys:**
- `Space` - Select/deselect in lists
- `Enter` - Action/details
- `Tab` - Switch panes/focus
- `Esc` - Back/cancel
- Arrow keys - Navigate lists/trees

#### Visual Design Principles
**Information Hierarchy:**
1. **Critical info first**: Status, size, health indicators
2. **Contextual details**: Show relevant information based on current task
3. **Progressive depth**: Summary → details → actions
4. **Visual grouping**: Related information clustered visually

**Color and Typography:**
```go
// Theme system for consistent visual language
type Theme struct {
    Primary     lipgloss.Color  // #007ACC (blue)
    Success     lipgloss.Color  // #28A745 (green)
    Warning     lipgloss.Color  // #FFC107 (yellow)
    Danger      lipgloss.Color  // #DC3545 (red)
    Muted       lipgloss.Color  // #6C757D (gray)
    Background  lipgloss.Color  // Terminal default
    
    HeaderStyle lipgloss.Style
    InfoStyle   lipgloss.Style
    ErrorStyle  lipgloss.Style
    BorderStyle lipgloss.Style
}
```

### Integration with Existing CLI

#### Hybrid Architecture
```go
// Shared core between CLI and TUI modes
package core

type ArchiveManager struct {
    storage   *storage.Manager
    registry  *storage.Registry
    config    *config.Config
}

// CLI mode (existing)
func RunCLI(args []string) error {
    manager := NewArchiveManager()
    return cobra.Execute(manager, args)
}

// TUI mode (new)
func RunTUI() error {
    manager := NewArchiveManager()
    app := NewTUIApp(manager)
    return tea.NewProgram(app).Start()
}

// Entry point decides mode
func main() {
    if len(os.Args) > 1 && os.Args[1] == "tui" {
        RunTUI()
    } else {
        RunCLI(os.Args[1:])
    }
}
```

#### Command Compatibility
```bash
# CLI mode (existing, unchanged)
7zarch list --profile=media
7zarch create /path/to/source
7zarch show 01HF5K2M

# TUI mode (new)
7zarch tui                    # Full TUI application
7zarch tui --dashboard        # Start in dashboard view
7zarch tui --browse           # Start in browse view
7zarch tui --search=query     # Start with search

# Hybrid operations
7zarch list --interactive     # Interactive selection in CLI context
7zarch create --interactive   # Interactive creation with path browser
```

## Implementation Strategy

### Building on PR21 Foundation
**Existing Foundation (PR21 `feat/tui-actions-errors-overlay`):**
- ✅ Bubble Tea application structure working
- ✅ Table component with archive data integration
- ✅ Selection system and multi-select functionality
- ✅ Action overlays and confirmation dialogs
- ✅ Basic keyboard navigation patterns

### Implementation Plan

### Phase 1: Simplify & Theme System (CC)
- [ ] **Strip Complex UI** (CC)
  - [ ] Remove complex dashboard and overlays
  - [ ] Focus on simple list → details → actions flow
  - [ ] Keep working table and selection from PR21
  - [ ] Simplify navigation to arrow keys + single letter actions

- [ ] **Theme System Implementation** (CC)
  - [ ] Theme configuration structure
  - [ ] 5 base themes: Blue Tech, Terminal Green, Purple, Cyan, Charmbracelet
  - [ ] 4 Dracula variations: Classic, Warm, Cool, Minimal
  - [ ] Theme switching mechanism (`--theme` flag or config)
  - [ ] Color application to all UI elements

### Phase 2: Core Functionality (CC)
- [ ] **List View Enhancement** (CC)
  - [ ] Simple archive list with theme application
  - [ ] Arrow key navigation with visual selection
  - [ ] Multi-select with Space bar (checkboxes)
  - [ ] Integration with 7EP-0007 query system when ready

- [ ] **Details View** (CC)
  - [ ] Simple archive details display
  - [ ] Essential metadata (size, date, status, location)
  - [ ] Navigation back to list
  - [ ] Theme-consistent styling

### Phase 3: Actions & Safety (CC)
- [ ] **Simple Actions** (CC)
  - [ ] Single-key actions (d=delete, m=move, u=upload)
  - [ ] Multi-select batch operations
  - [ ] Safe confirmation dialogs
  - [ ] Integration with existing CLI commands

- [ ] **Upload Integration** (CC)
  - [ ] TrueNAS upload preparation
  - [ ] Progress tracking for large files
  - [ ] Remote storage status integration
  - [ ] Future Tailscale network transparency

### Phase 4: Polish & Future Features (CC)
- [ ] **Theme Management** (CC)
  - [ ] Runtime theme switching
  - [ ] User preference persistence
  - [ ] Custom theme creation capability

- [ ] **Remote Storage Integration** (CC)
  - [ ] TrueNAS browsing (when backend ready)
  - [ ] Download from remote archives
  - [ ] Integrity verification without download
  - [ ] Seamless local/remote navigation

## 🎨 Complete Theme System Specification

### Theme Architecture
```go
type Theme struct {
    Name        string
    Header      lipgloss.Color  // Title/app name
    Foreground  lipgloss.Color  // Normal text
    Selection   lipgloss.Color  // Selected item background
    SelText     lipgloss.Color  // Selected item text
    Metadata    lipgloss.Color  // Size/date info
    StatusOK    lipgloss.Color  // ✓ Present archives
    StatusMiss  lipgloss.Color  // ? Missing archives  
    StatusDel   lipgloss.Color  // X Deleted archives
    Commands    lipgloss.Color  // Command help text
}
```

### Complete Theme Definitions

#### **1. Blue Tech Theme**
- Header: `#00BFFF`, Foreground: `#FFFFFF`, Selection: `#1E3A8A`, SelText: `#FFFFFF`
- Metadata: `#22D3EE`, StatusOK: `#10B981`, StatusMiss: `#F59E0B`, Commands: `#3B82F6`

#### **2. Terminal Green Theme**  
- Header: `#00FF00`, Foreground: `#FFFFFF`, Selection: `#065F46`, SelText: `#00FF00`
- Metadata: `#86EFAC`, StatusOK: `#00FF00`, StatusMiss: `#FBBF24`, Commands: `#22C55E`

#### **3. Purple Gradient Theme**
- Header: `#FF00FF`, Foreground: `#FFFFFF`, Selection: `#7C3AED`, SelText: `#F472B6`
- Metadata: `#C4B5FD`, StatusOK: `#EC4899`, StatusMiss: `#FB923C`, Commands: `#8B5CF6`

#### **4. Neon Cyan Theme**
- Header: `#00FFFF`, Foreground: `#FFFFFF`, Selection: `#0F766E`, SelText: `#00FFFF`
- Metadata: `#FB923C`, StatusOK: `#10B981`, StatusMiss: `#F97316`, Commands: `#06B6D4`

#### **5. Charmbracelet Theme**
- Header: `#FF1493`, Foreground: `#FFFFFF`, Selection: `#EC4899`, SelText: `#FFFFFF`
- Metadata: `#FBCFE8`, StatusOK: `#32CD32`, StatusMiss: `#FFD700`, Commands: `#EC4899`

#### **6. Dracula Classic Theme**
- Header: `#BD93F9`, Foreground: `#F8F8F2`, Selection: `#44475A`, SelText: `#BD93F9`
- Metadata: `#8BE9FD`, StatusOK: `#50FA7B`, StatusMiss: `#F1FA8C`, Commands: `#FF79C6`

#### **7. Dracula Warm Theme**
- Header: `#BD93F9`, Foreground: `#F8F8F2`, Selection: `#44475A`, SelText: `#FFB86C`
- Metadata: `#FFB86C`, StatusOK: `#50FA7B`, StatusMiss: `#F1FA8C`, Commands: `#FF79C6`

#### **8. Dracula Cool Theme**  
- Header: `#FF79C6`, Foreground: `#F8F8F2`, Selection: `#44475A`, SelText: `#FF79C6`
- Metadata: `#BD93F9`, StatusOK: `#50FA7B`, StatusMiss: `#FF5555`, Commands: `#BD93F9`

#### **9. Dracula Minimal Theme**
- Header: `#BD93F9`, Foreground: `#F8F8F2`, Selection: `#44475A`, SelText: `#F8F8F2`
- Metadata: `#6272A4`, StatusOK: `#50FA7B`, StatusMiss: `#F1FA8C`, Commands: `#6272A4`

### Dependencies
- **7EP-0004**: MAS Foundation ✅ (completed) - provides storage and registry infrastructure
- **7EP-0014**: Critical Foundation Gaps ✅ (completed) - provides reliable foundation  
- **7EP-0007**: Enhanced MAS Operations 🔄 (CC implementing) - search and query integration
- **PR21**: TUI Foundation ✅ (existing) - Bubble Tea app, table, selection system

## Testing Strategy

### Acceptance Criteria
- [ ] TUI launches quickly (<500ms) and provides immediate value in dashboard
- [ ] All major workflows accessible within 3 keystrokes from dashboard
- [ ] Path browser with tab completion works across all major platforms
- [ ] Multi-select and batch operations handle 100+ archives efficiently
- [ ] TUI maintains full feature parity with CLI for core operations
- [ ] Responsive layout adapts to terminal sizes 80-200 columns
- [ ] Memory usage remains under 100MB for typical usage (1000+ archives)

### Test Scenarios

#### User Experience Testing
- **Entry experience**: Dashboard provides immediate orientation and value
- **Navigation efficiency**: Common workflows completable with minimal keystrokes
- **Discovery**: New users can understand capabilities without documentation
- **Accessibility**: Works well in various terminal environments and configurations

#### Functional Testing
- **Archive operations**: Create, move, delete, upload work correctly in TUI
- **Search and filtering**: All search capabilities available and performant
- **Batch operations**: Multi-select operations handle edge cases gracefully
- **Data integrity**: All operations maintain archive registry consistency

#### Performance Testing
- **Large datasets**: Smooth performance with 10,000+ archives
- **Real-time updates**: Live search and filtering remain responsive
- **Memory efficiency**: No memory leaks during long-running sessions
- **Terminal compatibility**: Works across major terminal emulators

#### Integration Testing
- **CLI compatibility**: TUI and CLI modes share configuration and data
- **Cross-platform**: Consistent behavior on macOS, Linux, Windows
- **Terminal variety**: Works in SSH, tmux, screen, various emulators
- **Theme compatibility**: Visual appearance adapts to terminal capabilities

### Performance Benchmarks
- **Application startup**: <500ms cold start, <200ms warm start
- **View transitions**: <100ms between major views
- **List rendering**: <200ms for 1000 archives
- **Search performance**: <50ms for real-time search results
- **Batch operations**: Progress updates every 100ms for responsive feedback

## Migration/Compatibility

### Breaking Changes
None - TUI mode is completely additive to existing CLI functionality.

### Upgrade Path
- Existing CLI commands continue working unchanged
- TUI mode available via `7zarch tui` command
- Gradual feature parity migration for interactive operations
- Shared configuration and data between CLI and TUI modes

### Backward Compatibility
Complete CLI compatibility maintained:
- All existing commands and flags continue working
- Configuration format unchanged
- Data format and storage unchanged
- Script integration unaffected

## Alternatives Considered

**Web-based interface**: Considered building a web UI instead of TUI, but terminal-native interface provides better integration with developer workflows and doesn't require browser/server setup.

**IDE plugin**: Evaluated creating VSCode/editor plugins, but TUI provides broader accessibility and doesn't tie users to specific development environments.

**Extending existing CLI with interactive flags**: Considered adding `--interactive` flags to existing commands, but comprehensive TUI provides much richer experience for complex workflows.

**Desktop GUI application**: Evaluated cross-platform desktop app, but TUI maintains the command-line focus and doesn't require GUI framework dependencies.

## CC Implementation Strategy

### CC (Claude Code) Responsibilities - Full TUI Implementation
- **TUI Simplification**: Strip PR21 complexity, focus on list → details → actions
- **Theme System**: 9 complete themes with runtime switching capability  
- **List Navigation**: Arrow keys, single-letter actions, multi-select with Space
- **Details View**: Simple metadata display with theme consistency
- **Action Integration**: Leverage existing CLI commands for operations
- **Safety Patterns**: Confirmation dialogs for destructive operations
- **Performance**: Efficient rendering for large archive collections

### Implementation Coordination with Amp
- **UX Design**: Amp provides UI mockups, flow guidance, and theme specifications
- **Architecture Review**: Amp validates approach against podcast archival workflow
- **Progress Oversight**: Amp monitors implementation against simple UI vision
- **Quality Assurance**: Amp ensures TUI maintains CLI functionality while improving UX

### PR21 Foundation Strategy
**Leverage Existing Work:**
- ✅ **Keep**: Bubble Tea architecture, table component, selection system
- 🔄 **Simplify**: Remove complex overlays, wizards, multi-pane layouts
- 🎨 **Enhance**: Add comprehensive theme system with 9 color schemes
- 🎯 **Focus**: List view + details view + simple actions

## Future Considerations

### Advanced TUI Features
- **Plugin System**: Extensible interface for custom operations and views
- **Remote Mode**: TUI interface for remote archive management over SSH
- **Integration APIs**: Hooks for external tools and automation
- **Multi-instance**: Support for managing multiple archive repositories

### Collaborative Features
- **Multi-user**: Shared archive management with conflict resolution
- **Activity Streams**: Real-time updates from other users/systems
- **Commenting**: Archive and operation annotation system
- **Approval Workflows**: Review and approval processes for sensitive operations

### Advanced Visualizations
- **Timeline Views**: Archive activity and evolution over time
- **Relationship Graphs**: Visual archive dependency and similarity mapping
- **Storage Analytics**: Advanced disk usage and optimization recommendations
- **Trend Analysis**: Predictive storage planning and capacity management

## References

- **Builds on**: 7EP-0004 MAS Foundation Implementation, 7EP-0009 Enhanced Display System
- **Integrates with**: 7EP-0001 Trash Management, 7EP-0007 Enhanced MAS Operations
- **Enables**: Rich interactive archive management beyond CLI limitations
- **Inspired by**: Modern TUI applications like `lazygit`, `k9s`, `htop`, `ncdu`, `ranger`