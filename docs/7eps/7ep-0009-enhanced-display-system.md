# 7EP-0009: Enhanced Display System

**Status:** Draft  
**Author(s):** Claude Code (CC), Adam Stacoviak  
**Assignment:** CC Lead (UI/UX), AC Support (Integration)  
**Difficulty:** 2 (moderate - UI/UX focused with clear patterns)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  

## Executive Summary

Implement a comprehensive display system for archive listing that provides multiple visual presentation modes optimized for different user workflows. Transform the current single-format list output into a flexible, theme-able, and context-aware display system with 5 distinct modes: Table, Card, Compact, Tree, and Dashboard.

## Evidence & Reasoning

**Current limitations:**
- Single fixed-width table format for all archive listing
- Poor readability in narrow terminals (SSH, mobile)
- Limited information hierarchy for large archive collections
- No overview/summary capabilities for archive management
- Difficult to identify patterns and relationships between archives

**User workflow pain points:**
- **Power users**: Need high-density scanning for large collections (100+ archives)
- **Casual users**: Want rich detail view for small collections (5-15 archives)
- **Mobile/SSH users**: Need compact output for narrow terminals
- **Managers**: Want overview statistics and health monitoring
- **Organizers**: Need hierarchical grouping to understand structure

**Why enhanced displays matter:**
- Different archive management tasks require different information presentation
- Terminal environment varies widely (width, capabilities, user preferences)
- Archive collections grow from dozens to thousands requiring scalable UX
- Modern CLI tools provide rich, adaptive user experiences

## Use Cases

### Primary Use Case: Context-Aware Display Selection
```bash
# Automatic smart defaults based on context
7zarch list                         # Auto-detects best display for terminal/count
7zarch list --missing              # Auto-switches to compact (problem scanning)
7zarch list --profile=media        # Auto-switches to tree (grouping context)
7zarch list | head -5              # Auto-switches to compact (piped output)

# Explicit display mode selection
7zarch list --table                # High-density scanning
7zarch list --card                 # Rich detail exploration
7zarch list --compact              # Terminal-friendly minimal
7zarch list --tree                 # Hierarchical grouping
7zarch list --dashboard            # Management overview
```

### Secondary Use Cases

#### Power User Workflows
```bash
# High-density archive scanning
7zarch list --table --sort-by=size,age --columns="id,name,size,age,status"

# Quick problem identification
7zarch list --missing --compact | grep "important"

# Batch operation preparation
7zarch list --table --profile=media --larger-than=1GB --save-query=large-media
```

#### Management & Organization
```bash
# Storage overview and health
7zarch list --dashboard --period=month

# Archive organization analysis
7zarch list --tree --group-by=profile
7zarch list --tree --group-by=date --detail=summary

# Project-specific exploration
7zarch list --card --pattern="*project*" --detail=full
```

#### Terminal Environment Adaptation
```bash
# Narrow terminal (SSH, mobile)
COLUMNS=80 7zarch list              # Auto-compact mode
7zarch list --compact --width=80    # Force compact

# Wide terminal (desktop)
COLUMNS=140 7zarch list             # Auto-table mode with full columns

# Script integration
7zarch list --compact --no-headers | awk '{print $1}' | xargs -I {} 7zarch show {}
```

## Technical Design

### Display Mode Architecture

#### 1. Display Interface (`internal/display/`)
```go
// Core display system interface
type DisplayMode interface {
    Render(archives []*storage.Archive, opts DisplayOptions) error
    Name() string
    AutoDetect(context DisplayContext) bool
}

type DisplayOptions struct {
    Details      bool
    Width        int
    Theme        Theme
    Columns      []string
    GroupBy      string
    SortBy       []string
    ShowHeaders  bool
}

type DisplayContext struct {
    TerminalWidth int
    ArchiveCount  int
    FilterContext string
    OutputPiped   bool
}
```

#### 2. Display Modes (`internal/display/modes/`)

**Table Mode** - High-density information scanning
```go
type TableDisplay struct {
    columns    []Column
    sortable   bool
    pagination bool
}

func (td *TableDisplay) Render(archives []*storage.Archive, opts DisplayOptions) error {
    // Adaptive column sizing based on terminal width
    // Sortable headers with visual indicators
    // Pagination for large datasets
    // Color coding for status (present/missing/deleted)
}
```

**Card Mode** - Rich information exploration
```go
type CardDisplay struct {
    detailLevel DetailLevel
    showMeta    bool
    expandable  bool
}

func (cd *CardDisplay) Render(archives []*storage.Archive, opts DisplayOptions) error {
    // Visual cards with hierarchical information
    // Emoji indicators and visual spacing
    // Expandable metadata sections
    // Related archive suggestions
}
```

**Compact Mode** - Terminal-friendly minimal output
```go
type CompactDisplay struct {
    singleLine bool
    essential  bool
}

func (cd *CompactDisplay) Render(archives []*storage.Archive, opts DisplayOptions) error {
    // Single line per archive
    // Essential information only
    // Pipe-friendly output format
    // Abbreviated status indicators
}
```

**Tree Mode** - Hierarchical grouping and organization
```go
type TreeDisplay struct {
    groupBy     GroupingStrategy
    collapsible bool
    showStats   bool
}

func (td *TreeDisplay) Render(archives []*storage.Archive, opts DisplayOptions) error {
    // Hierarchical tree structure
    // Group by profile, date, size, location
    // Collapsible sections
    // Group statistics and summaries
}
```

**Dashboard Mode** - Overview and management
```go
type DashboardDisplay struct {
    period    TimePeriod
    showGraph bool
    quickActions bool
}

func (dd *DashboardDisplay) Render(archives []*storage.Archive, opts DisplayOptions) error {
    // Storage statistics and health metrics
    // ASCII graphs for trends
    // Quick action suggestions
    // Problem identification and recommendations
}
```

#### 3. Theme System (`internal/display/themes/`)
```go
type Theme struct {
    Name        string
    Colors      ColorScheme
    Emojis      EmojiSet
    Borders     BorderStyle
    Spacing     SpacingRules
}

// Built-in themes
var (
    RichTheme       = Theme{/* full colors, emojis, borders */}
    MinimalTheme    = Theme{/* reduced visual elements */}
    MonochromeTheme = Theme{/* no colors, plain text */}
    TerminalTheme   = Theme{/* terminal-safe, universal */}
)
```

#### 4. Auto-Detection System (`internal/display/detector/`)
```go
type DisplayDetector struct {
    rules []DetectionRule
}

type DetectionRule struct {
    Condition func(DisplayContext) bool
    Mode      string
    Priority  int
}

func (dd *DisplayDetector) DetectBest(context DisplayContext) string {
    // Terminal width analysis
    // Archive count optimization
    // Filter context evaluation
    // User preference integration
}
```

### Command Integration

#### Enhanced List Command
```go
func ListCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List archives with enhanced display options",
        RunE:  runListWithDisplay,
    }
    
    // Display mode flags
    cmd.Flags().Bool("table", false, "Use table display mode")
    cmd.Flags().Bool("card", false, "Use card display mode")
    cmd.Flags().Bool("compact", false, "Use compact display mode")
    cmd.Flags().Bool("tree", false, "Use tree display mode")
    cmd.Flags().Bool("dashboard", false, "Use dashboard display mode")
    
    // Display options
    cmd.Flags().String("theme", "", "Display theme (rich|minimal|monochrome)")
    cmd.Flags().StringSlice("columns", nil, "Columns to show (table mode)")
    cmd.Flags().String("group-by", "profile", "Grouping strategy (tree mode)")
    cmd.Flags().String("detail", "summary", "Detail level (card mode)")
    cmd.Flags().Int("width", 0, "Force terminal width")
    cmd.Flags().Bool("no-headers", false, "Hide headers (compact mode)")
    
    return cmd
}
```

#### Configuration Integration
```go
// Configuration file support
type DisplayConfig struct {
    DefaultMode  string            `yaml:"default_mode"`
    AutoDetect   bool             `yaml:"auto_detect"`
    Theme        string           `yaml:"theme"`
    ContextModes map[string]string `yaml:"context_modes"`
    TableOptions TableConfig      `yaml:"table"`
    TreeOptions  TreeConfig       `yaml:"tree"`
}

// User configuration examples
var defaultDisplayConfig = DisplayConfig{
    DefaultMode: "auto",
    AutoDetect:  true,
    Theme:      "rich",
    ContextModes: map[string]string{
        "missing":        "compact",
        "piped":          "compact",
        "narrow_terminal": "compact",
        "profile_filter": "tree",
    },
}
```

## Display Mode Specifications

### Table Mode - Information Density Focus
**Purpose**: High-density scanning for power users with large collections

**Features**:
- Adaptive column sizing based on terminal width
- Sortable columns with visual indicators (↑↓)
- Pagination for datasets >50 archives
- Color-coded status indicators
- Column selection and reordering

**Example Output**:
```
┌─────────┬────────────────────────────────────┬──────────┬─────────┬────────────┬────────┐
│ ID ↓    │ Name                               │ Size     │ Profile │ Age        │ Status │
├─────────┼────────────────────────────────────┼──────────┼─────────┼────────────┼────────┤
│ 01K2E3F │ project-backup-2024-03-15.7z       │ 156.7 MB │ media   │ 5d         │ ✓ M    │
│ 01K2E3G │ documents-archive-large.7z         │  2.3 GB  │ docs    │ 12d        │ ✓ M    │
│ 01K2E3H │ external-backup.7z                 │ 945.2 MB │ balance │ 3d         │ ⚠️ E   │
└─────────┴────────────────────────────────────┴──────────┴─────────┴────────────┴────────┘
```

### Card Mode - Rich Information Display
**Purpose**: Detailed exploration of smaller archive sets

**Features**:
- Visual hierarchy with emojis and spacing
- Expandable metadata sections
- Full path display
- Related archive suggestions
- Visual status indicators

**Example Output**:
```
┌─ Archive Details ─────────────────────────────────────────────────────────────────────┐
│                                                                                       │
│  📁 project-backup-2024-03-15.7z                                      [01K2E3F4]     │
│     💾 156.7 MB • 📊 media profile • ⏰ 5 days ago • ✅ Present                       │
│     📂 /managed/archives/2024/03/project-backup-2024-03-15.7z                        │
│     🔍 Contains: 1,247 files • Compression: 68% • Checksum: verified                 │
│                                                                                       │
│  📁 documents-archive-large.7z                                       [01K2E3G5]     │
│     💾 2.3 GB • 📊 documents profile • ⏰ 12 days ago • ✅ Present                   │
│     📂 /managed/archives/2024/02/documents-archive-large.7z                          │
│     🔍 Contains: 3,891 files • Compression: 45% • Checksum: verified                 │
│                                                                                       │
└───────────────────────────────────────────────────────────────────────────────────────┘
```

### Compact Mode - Terminal-Friendly Minimal
**Purpose**: Narrow terminals, SSH sessions, script integration

**Features**:
- Single line per archive
- Essential information only
- Pipe-friendly output
- No visual decorations
- Abbreviated indicators

**Example Output**:
```
01K2E3F4  project-backup-2024-03-15.7z              156.7MB  media     5d   ✓M
01K2E3G5  documents-archive-large.7z                  2.3GB  docs     12d   ✓M  
01K2E3H6  external-backup.7z                        945.2MB  balance   3d   ⚠️E
01K2E3I7  photos-vacation-2024.7z                     1.8GB  media     1d   ✓M
```

### Tree Mode - Hierarchical Organization
**Purpose**: Understanding archive relationships and organization

**Features**:
- Multiple grouping strategies (profile, date, size, location)
- Collapsible sections
- Group statistics
- Visual tree structure
- Size visualization

**Example Output**:
```
📊 Archives by Profile (15 total, 2.1 TB)

📱 media (6 archives, 1.2 TB)
├─ ✅ photos-vacation-2024.7z          [01K2E3I6]  1.8 GB  managed    1d
├─ ✅ project-backup-2024-03-15.7z     [01K2E3F4]  156.7 MB  managed  5d  
├─ ✅ video-projects-archive.7z        [01K2E3N1]  892.1 MB  managed  15d
└─ ⚠️ missing-media-backup.7z          [01K2E3O2]  1.1 GB  external   8d

📄 documents (5 archives, 3.4 GB)  
├─ ✅ documents-archive-large.7z       [01K2E3G5]  2.3 GB  managed   12d
├─ ✅ code-repositories-Q1.7z          [01K2E3J7]  834.5 MB  managed  8d
└─ ✅ shared-drive-sync.7z             [01K2E3K8]  1.2 GB  external   6d

⚖️ balanced (4 archives, 2.8 GB)
├─ ⚠️ external-backup.7z               [01K2E3H6]  945.2 MB  external  3d
└─ ✅ mixed-content-archive.7z         [01K2E3P3]  1.9 GB  managed    9d
```

### Dashboard Mode - Management Overview
**Purpose**: High-level overview and maintenance planning

**Features**:
- Storage statistics and utilization
- Health metrics and recommendations
- Growth trends (ASCII graphs)
- Quick action suggestions
- Problem identification

**Example Output**:
```
7zarch-go Archive Dashboard
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 Storage Overview                     🏷️  Profiles                    ⚠️  Attention Needed
┌─────────────────────────────────┐    ┌─────────────────────────────┐   ┌─────────────────────┐
│ Total:    15 archives (2.1 TB) │    │ 📱 media:      6 (1.2 TB)  │   │ Missing:        2   │
│ Active:   12 archives (1.8 TB) │    │ 📄 documents:  5 (645 MB)  │   │ Deleted:        3   │  
│ Missing:   2 archives (1.1 GB) │    │ ⚖️ balanced:    4 (2.8 GB)  │   │ Old (>6mo):     4   │
│ Deleted:   3 archives (546 MB) │    └─────────────────────────────┘   └─────────────────────┘
└─────────────────────────────────┘    

📈 Growth Trend (Last 30 Days)          🎯 Quick Actions
Size:  ████████████▓▓▓▓ 75% growth      • list --missing (check 2 missing archives)
Count: █████████▓▓▓▓▓▓▓ 60% growth      • trash list (review 3 deleted archives)  
                                        • list --older-than=6m (cleanup candidates)

Recent Archives (5 most recent)
┌─────────┬─────────────────────────────────┬──────────┬─────────┬────────────────────┐
│ ID      │ Name                            │ Size     │ Status  │ Created            │  
├─────────┼─────────────────────────────────┼──────────┼─────────┼────────────────────┤
│ 01K2E3I │ photos-vacation-2024.7z         │  1.8 GB  │ ✅ M    │ 2024-08-11 14:30   │
│ 01K2E3H │ external-backup.7z              │ 945.2 MB │ ⚠️ E    │ 2024-08-09 09:15   │
│ 01K2E3G │ documents-archive-large.7z      │  2.3 GB  │ ✅ M    │ 2024-07-31 16:45   │
└─────────┴─────────────────────────────────┴──────────┴─────────┴────────────────────┘
```

## Auto-Detection & Smart Defaults

### Detection Rules
```go
var detectionRules = []DetectionRule{
    // Terminal width-based
    {
        Condition: func(ctx DisplayContext) bool { return ctx.TerminalWidth < 80 },
        Mode: "compact",
        Priority: 100,
    },
    
    // Archive count-based
    {
        Condition: func(ctx DisplayContext) bool { return ctx.ArchiveCount > 50 },
        Mode: "table",
        Priority: 80,
    },
    
    // Filter context-based
    {
        Condition: func(ctx DisplayContext) bool { return ctx.FilterContext == "missing" },
        Mode: "compact",
        Priority: 90,
    },
    
    // Piped output
    {
        Condition: func(ctx DisplayContext) bool { return ctx.OutputPiped },
        Mode: "compact",
        Priority: 95,
    },
}
```

### Context-Aware Behavior
- **Narrow terminals** (<80 cols): Auto-compact mode
- **Large collections** (>50 archives): Auto-table mode for scanning
- **Filter contexts**: Missing archives → compact, profile filters → tree
- **Piped output**: Always compact for script compatibility
- **User overrides**: Explicit flags always take precedence

## Implementation Plan

### Phase 1: Core Display Infrastructure (CC Lead)
- [ ] **Display Interface & Manager** (CC)
  - [ ] Create display mode interface and registration system
  - [ ] Implement display context detection
  - [ ] Build auto-detection rule engine
  - [ ] Add display configuration support

- [ ] **Basic Display Modes** (CC)
  - [ ] Enhanced table mode with adaptive columns
  - [ ] Compact mode for narrow terminals
  - [ ] Theme system foundation (rich, minimal, monochrome)

### Phase 2: Advanced Display Modes (CC Lead, AC Integration)
- [ ] **Rich Display Modes** (CC)
  - [ ] Card mode with visual hierarchy
  - [ ] Tree mode with grouping strategies
  - [ ] Dashboard mode with statistics

- [ ] **Command Integration** (AC)
  - [ ] Add display flags to list command
  - [ ] Integrate with existing filter system
  - [ ] Update help and documentation

### Phase 3: User Experience Polish (Shared)
- [ ] **Advanced Features** (CC)
  - [ ] Column sorting and selection (table mode)
  - [ ] Expandable sections (card mode)
  - [ ] Collapsible tree sections
  - [ ] ASCII graph rendering (dashboard)

- [ ] **Configuration & Customization** (AC)
  - [ ] User preference system
  - [ ] Context-specific mode configuration
  - [ ] Custom theme creation
  - [ ] Performance optimization for large datasets

### Dependencies
- **7EP-0004**: MAS Foundation (completed) - provides archive listing infrastructure
- **7EP-0001**: Trash Management (completed) - dashboard mode shows trash statistics
- **7EP-0007**: Enhanced MAS Operations (planned) - search integration with display modes

## Testing Strategy

### Acceptance Criteria
- [ ] All 5 display modes render correctly across terminal sizes (60-200 columns)
- [ ] Auto-detection selects appropriate mode for context 95% of the time
- [ ] Display performance <100ms for collections up to 1000 archives
- [ ] All display modes integrate seamlessly with existing filter system
- [ ] Configuration system allows persistent user preferences
- [ ] Pipe-friendly output maintains script compatibility

### Test Scenarios

#### Display Mode Testing
- Render each mode with various archive counts (1, 10, 100, 1000)
- Terminal width testing (60, 80, 120, 200 columns)
- Theme testing across different terminal capabilities
- Color/emoji rendering in various terminal environments

#### Auto-Detection Testing
- Verify correct mode selection for different contexts
- Test override behavior with explicit flags
- Configuration precedence testing
- Edge case handling (empty results, terminal detection failures)

#### Integration Testing
- Filter integration across all display modes
- Configuration loading and persistence
- Performance testing with large datasets
- Cross-platform compatibility (macOS, Linux, Windows)

### Performance Benchmarks
- **Display rendering**: <100ms for 1000 archives
- **Auto-detection**: <10ms for context analysis
- **Configuration loading**: <50ms for complex config files
- **Memory usage**: <50MB for 10,000 archive display

## Migration/Compatibility

### Breaking Changes
None - all enhancements are additive to existing list command.

### Upgrade Path
- Existing `7zarch list` behavior unchanged (auto-detection chooses best mode)
- New display flags opt-in to specific modes
- Configuration system starts with sensible defaults
- Legacy table format available via `--table` flag

### Backward Compatibility
Complete backward compatibility maintained:
- All existing list command flags continue working
- Output format remains identical without display flags
- Script integration unaffected unless using new flags

## Alternatives Considered

**Single enhanced table**: Considered improving only the existing table format, but different use cases require fundamentally different information layouts.

**External display tools**: Evaluated integration with tools like `fzf` or `bat`, but native display modes provide better integration and don't require external dependencies.

**Configuration-only approach**: Considered making everything configurable rather than distinct modes, but explicit mode selection provides clearer user mental models.

**TUI-first approach**: Evaluated building interactive TUI instead of enhanced CLI displays, but CLI-first maintains compatibility and broader accessibility.

## AC/CC Implementation Split

### CC (Claude Code) Responsibilities - Display Infrastructure
- **Display System Architecture**: Interface design, mode registration, auto-detection
- **Core Display Modes**: Table, compact, card rendering implementations
- **Theme System**: Color schemes, emoji sets, border styles
- **Performance Optimization**: Rendering efficiency, large dataset handling
- **Testing Infrastructure**: Display testing framework, performance benchmarks

### AC (Augment Code) Responsibilities - User Experience
- **Command Integration**: List command flag integration, help text
- **Configuration System**: User preferences, context-specific defaults
- **Advanced Display Features**: Column selection, grouping options, detail levels
- **User Workflow Design**: Mode selection logic, smart defaults
- **Documentation**: User guides, display mode examples, configuration reference

### Shared Responsibilities
- **API Design**: Display option interfaces and flag naming (AC leads, CC implements)
- **Integration Testing**: Cross-component validation and user workflow testing
- **Theme Design**: Visual appearance and accessibility (AC designs, CC implements)
- **Auto-Detection Logic**: Context analysis rules (AC defines, CC implements)

### Coordination Points
1. **Display Interface Design**: How modes integrate with existing list infrastructure
2. **Configuration Schema**: User preference format and storage approach
3. **Auto-Detection Rules**: Context analysis and mode selection logic
4. **Performance Requirements**: Rendering speed and memory usage targets

## Future Considerations

### Enhanced Interactivity
- **Interactive mode**: Arrow key navigation, real-time filtering
- **TUI integration**: Full-screen interface for complex operations
- **Multi-select**: Batch operations with visual selection
- **Live updates**: Real-time refresh for changing archive status

### Advanced Visualization
- **Charts and graphs**: Storage trends, archive aging analysis
- **Heatmaps**: Activity patterns, storage utilization
- **Export formats**: HTML, PDF reports for sharing
- **Custom layouts**: User-defined display templates

### Integration Opportunities
- **Search integration**: Enhanced display for search results (7EP-0007)
- **Batch operations**: Visual progress tracking for bulk actions
- **Cloud storage**: Provider-specific display enhancements
- **Monitoring**: Real-time status updates and notifications

## References

- **Builds on**: 7EP-0004 MAS Foundation Implementation
- **Integrates with**: 7EP-0001 Trash Management System, 7EP-0007 Enhanced MAS Operations
- **Enables**: Rich user experience for archive management workflows
- **Inspired by**: Modern CLI tools like `kubectl`, `docker`, `gh`, `exa`, `bat`