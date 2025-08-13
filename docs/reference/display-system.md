# Display System Reference

The 7zarch-go display system provides multiple viewing modes for archive listing, each optimized for different user workflows and terminal environments.

## Quick Reference

### Display Modes

| Mode | Flag | Purpose | Min Width | Use Case |
|------|------|---------|-----------|----------|
| **Auto** | (default) | Context-aware selection | - | Smart defaults |
| **Table** | `--table` | High-density scanning | 80 cols | Power users, large collections |
| **Compact** | `--compact` | Minimal terminal output | 60 cols | SSH, mobile, scripting |
| **Card** | `--card` | Rich detail exploration | 80 cols | Small collections, detailed view |
| **Tree** | `--tree` | Hierarchical organization | 70 cols | Understanding structure |
| **Dashboard** | `--dashboard` | Management overview | 90 cols | Health monitoring, statistics |

### Common Usage

```bash
# Auto-detection (recommended)
7zarch list

# Specific modes
7zarch list --table
7zarch list --compact
7zarch list --card
7zarch list --tree  
7zarch list --dashboard

# With filters
7zarch list --compact --missing
7zarch list --dashboard --profile=media
```

## Implementation Details

### Status Icon System

Consistent status formatting across all display modes:

- **✓** Present archive (available and accessible)
- **?** Missing archive (referenced but file not found)  
- **X** Deleted archive (in trash, pending purge)

**Text Alternatives** (table/compact modes):
- **OK** Present
- **MISS** Missing
- **DEL** Deleted

### ULID Display Standards

All modes display **12-character ULID prefixes** for compatibility with the show command:

```bash
# List shows 12-char prefixes
7zarch list --compact
01K2E3BEJV6G  test-pod.7z  34.3KB  1d  ✓

# Show command accepts these prefixes
7zarch show 01K2E3BEJV6G
```

### Auto-Detection Logic

The system automatically selects the best display mode based on context:

| Condition | Selected Mode | Reason |
|-----------|---------------|--------|
| Terminal width < 80 | Compact | Narrow terminal optimization |
| Archive count > 50 | Table | High-density scanning needed |
| `--missing` filter | Compact | Problem identification workflow |
| Piped output | Compact | Script compatibility |
| Default | Table | General use optimization |

## Display Mode Details

### Table Mode (`--table`)

**Purpose**: High-density information scanning for large collections

**Features**:
- Bordered table with proper alignment
- 12-character ULID display
- Grouped by management status (MANAGED/EXTERNAL)
- Summary header with counts
- Text-based status indicators

**Best for**: Power users managing 50+ archives

### Compact Mode (`--compact`)

**Purpose**: Terminal-friendly minimal output

**Features**:
- Single line per archive
- Essential information only
- No visual decorations
- Script-friendly format
- Abbreviated status indicators

**Best for**: SSH sessions, mobile terminals, script integration

### Card Mode (`--card`)

**Purpose**: Rich information display with visual hierarchy

**Features**:
- Bordered cards with detailed metadata
- Visual grouping by status and location
- Icon-based status indicators
- Full path display
- Expandable details with `--details` flag

**Best for**: Exploring small collections (5-15 archives) in detail

### Tree Mode (`--tree`)

**Purpose**: Hierarchical organization view

**Features**:
- Directory-based grouping
- Visual tree structure with Unicode characters
- Status icons for each archive
- Group statistics
- Archive metadata in context

**Best for**: Understanding archive organization and relationships

### Dashboard Mode (`--dashboard`)

**Purpose**: Management overview and health monitoring

**Features**:
- Sectioned overview with statistics
- Health indicators and recommendations
- Storage breakdown and utilization
- Recent activity tracking
- Professional formatting with no-right-border design

**Best for**: Archive collection health monitoring and management planning

## Architecture

### Display Interface

```go
type Display interface {
    Render(archives []*storage.Archive, opts Options) error
    Name() string
    MinWidth() int
}
```

### Registration System

```go
// Display manager handles mode registration and selection
displayManager := display.NewManager()
displayManager.Register(display.ModeTable, modes.NewTableDisplay())
displayManager.Register(display.ModeCompact, modes.NewCompactDisplay())
// ... etc
```

### Key Design Patterns

#### 1. No-Right-Border Pattern (Dashboard)

Instead of complex padding calculations:
```go
// ❌ Complex and error-prone
fmt.Printf("│ %s%s │\n", content, strings.Repeat(" ", padding))

// ✅ Clean and maintainable  
fmt.Printf("│  %s\n", content)  // Note: extra space for header alignment
```

#### 2. Consistent Status Formatting

```go
// Centralized status formatting
status := display.FormatStatus(archive.Status, useIcons)
```

#### 3. Context-Aware Selection

```go
// Auto-detection based on environment
func (m *Manager) detectBestMode(opts Options) Mode {
    if ctx.OutputPiped { return ModeCompact }
    if ctx.TerminalWidth < 80 { return ModeCompact }
    if ctx.ArchiveCount > 50 { return ModeTable }
    return ModeTable // default
}
```

## Extending the System

### Adding a New Display Mode

1. **Create Mode Implementation**:
```go
type MyDisplay struct{}

func (md *MyDisplay) Name() string { return "my-mode" }
func (md *MyDisplay) MinWidth() int { return 60 }
func (md *MyDisplay) Render(archives []*storage.Archive, opts display.Options) error {
    // Implementation here
    return nil
}
```

2. **Register in List Command**:
```go
// Add flag
cmd.Flags().Bool("my-mode", false, "Use my display mode")

// Register mode
myDisplay := modes.NewMyDisplay()
displayManager.Register(display.ModeMyMode, myDisplay)
```

3. **Add to Auto-Detection** (optional):
```go
// Add detection rule for specific contexts
```

### Design Guidelines

1. **Consistency**: Use shared status formatting and ULID display standards
2. **Flexibility**: Support both `--details` and basic views
3. **Performance**: Render efficiently for large collections (1000+ archives)
4. **Accessibility**: Provide text alternatives to visual elements
5. **Maintainability**: Use simple, clear formatting patterns

## Troubleshooting

### Common Issues

**Display garbled in narrow terminal**:
- Use `--compact` mode explicitly
- Check terminal width with `echo $COLUMNS`

**Status icons not displaying**:
- Terminal may not support Unicode
- Text alternatives automatically used in table/compact modes

**Misaligned borders**:
- Ensure consistent use of border patterns
- Consider no-right-border approach for complex layouts

**Show command can't find archives from list**:
- Verify 12-character ULID display
- Check ULID prefix length in implementation

### Performance Optimization

- Large collections (1000+ archives): Use table mode for efficient scanning
- Memory-constrained environments: Compact mode uses minimal resources
- Terminal width auto-detection: Relies on `term.GetSize()` for optimal layout

## Related Documentation

- [7EP-0009: Enhanced Display System](../7eps/7ep-0009-enhanced-display-system.md) - Complete specification
- [List Command Reference](commands/list.md) - Full command documentation
- [Archive Management Guide](../guides/managed-storage.md) - Archive workflow context