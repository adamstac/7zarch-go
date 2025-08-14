# 7EP-0016: TUI-First Interface Evolution

**Status:** Draft  
**Author(s):** Amp (Sourcegraph), Adam Stacoviak  
**Assignment:** Future Planning  
**Difficulty:** 4 (architectural - interface paradigm shift)  
**Created:** 2025-08-13  
**Updated:** 2025-08-13  

## Executive Summary

Evolve 7zarch-go into a **staging-to-remote workflow tool** with simple TUI-first interface and embedded command line. Focus on **podcast archival workflow**: MAS local as 24-48h staging → intelligent compression analysis → push to TrueNAS remote storage → simple remote browsing and pull-back capabilities.

## Evidence & Reasoning

**Current Success:**
- ✅ TUI experiment successful - beautiful themes, simple navigation, user-tested
- ✅ Multiple entry points working (`browse`, `ui`, `i`, `tui`)
- ✅ Solid foundation with viewport framework and proper Bubbletea architecture
- ✅ User feedback positive - "very good and simple start, great base to build on"

**User Vision:**
- **MAS Local = Staging** - 24-48 hour staging area before remote storage
- **Remote Browsing** - View and traverse TrueNAS archive collections  
- **Push/Pull Workflow** - Git-like model (pull to staging, work, push remote)
- **Compression Intelligence** - Tool analyzes and suggests recompression
- **Simple TUI** - Terminal-friendly, minimal borders and visual clutter
- **Embedded Command Line** - Vim-style `:` commands for advanced operations

**Strategic Opportunity:**
- Differentiate from all other CLI archive tools
- Create compelling visual demo and user experience
- Natural evolution path from current success

## Use Cases

### Primary Use Case: Staging-to-Remote Workflow
```bash
# Default launch shows staging + remote status
7zarch-go

# Simple interface with staging awareness
7zarch-go                                       Local Staging | Remote (TrueNAS)

Staging: 3 episodes (267 MB) - Ready for remote
> episode-424.7z        89 MB    2h ago    ✓ media profile
  episode-423.7z        92 MB    1d ago    ✓ needs upload  
  vacation-photos.7z   3.8 GB    1w ago    ✓ compress better?

Remote: 1,247 episodes (14.2 TB) via Tailscale

:push recent                                   <-- Simple command line
:browse remote                                 <-- View TrueNAS archives
:pull episode-420                              <-- Download for editing
:analyze compression vacation-photos.7z        <-- Check compression efficiency
```

### Secondary Use Cases

#### Remote Archive Analysis
```bash
# Browse remote storage via TUI
:browse remote

# Analyze compression efficiency
:analyze episode-420.7z
→ "Compressed with documents profile, could save 40% with media profile"
→ "Original: 2.1GB → Current: 89MB (96%) → Optimal: 54MB (97.4%)"

# View Tailscale URL for sharing
:details episode-423.7z
→ "Tailscale URL: http://truenas.tail-net.ts.net/archives/2024/episode-423.7z"
→ "Press 'c' to copy URL to clipboard"
```

#### Batch Recompression Workflows  
```bash
# Find archives that could benefit from recompression
:analyze old --profile-mismatch
→ "Found 23 archives from 2022-2023 using suboptimal compression"
→ "Potential savings: 3.2GB (18% improvement)"

# Batch recompress with new smart profiles
:recompress 2022 --profile media --confirm
→ "Will recompress 47 episodes from 2022 with media profile"
→ "Estimated time: 2h 15m, savings: 1.8GB"
```

#### Push/Pull Operations
```bash
# Git-like push/pull model
:push staging                    # Upload everything in staging to remote
:pull episode-400 episode-401    # Download specific episodes to staging  
:pull year:2023 guest:alice      # Bulk pull based on criteria
:status                          # Show staging vs remote sync status
```

## Technical Design

### Architecture Overview - Simple & Focused
```go
type TUIApp struct {
    // Current components (keep)
    archives     []*storage.Archive
    viewport     viewport.Model
    theme        Theme
    
    // Staging/Remote awareness
    localStaging  []*storage.Archive    // MAS local (24-48h staging)
    remoteArchives []*storage.Archive   // TrueNAS remote storage
    viewMode      ViewMode              // Local | Remote
    
    // Simple command line
    commandMode   bool                  // ':' toggles command mode
    commandInput  textinput.Model       // Simple command input
    commandHist   []string             // Command history
    
    // Remote integration
    tailscaleURL  string               // Base Tailscale URL for remote files
    syncStatus    map[string]SyncState // Track push/pull status
}

type ViewMode int
const (
    LocalView ViewMode = iota  // Staging area view
    RemoteView                 // TrueNAS remote storage view
)

type SyncState int
const (
    InSync SyncState = iota    // Local and remote match
    NeedsUpload               // Local newer, needs push
    NeedsDownload             // Remote newer, can pull
    ConflictState             // Both modified, needs resolution
)
```

### Core Features

#### 1. **Staging/Remote View Toggle**
```bash
# Toggle between local staging and remote archives
[Tab] Switch Local ⟷ Remote
[l] Local staging view
[r] Remote TrueNAS view
```

#### 2. **Compression Analysis Engine**
```bash
:analyze episode-420.7z
→ "Profile used: documents (suboptimal for audio)"
→ "Compression: 2.1GB → 89MB (95.7%)"  
→ "Optimal: media profile → 54MB (97.4%)"
→ "Tool recommendation: Recompress with media profile"
```

#### 3. **Remote File Details with Tailscale URLs**
```bash
# Details view for remote files
episode-423.7z (Remote)

Size: 89 MB
Created: August 13, 2024 2:30 PM
Profile: media (optimal ✓)
Tailscale URL: http://truenas.tail-net.ts.net/archives/2024/08/episode-423.7z

[c] Copy URL  [p] Pull to staging  [a] Analyze compression
```

#### 4. **Push/Pull Commands**
```bash
:push                          # Push all staging to remote
:push episode-424.7z          # Push specific file
:pull episode-420             # Pull specific episode to staging
:pull recent 5                # Pull last 5 episodes  
:status                       # Show sync status (like git status)
```

#### 5. **Batch Recompression**
```bash
:recompress --analyze         # Find suboptimal compression
:recompress 2022 --profile media --dry-run
:recompress selected --confirm
```

### Command Integration Strategy
```go
// Reuse existing CLI commands
func (app *TUIApp) executeEmbeddedCommand(cmdLine string) tea.Cmd {
    // Parse command line
    parts := strings.Fields(cmdLine)
    if len(parts) == 0 {
        return nil
    }
    
    // Use same command handlers as CLI
    return app.executor.Execute(parts[0], parts[1:])
}
```

## Implementation Plan

### Phase 1: Command Line Infrastructure
- [ ] **Embedded Command Line** 
  - [ ] Add textinput component at bottom of viewport
  - [ ] Vim-style `:` key to enter command mode
  - [ ] Enter to execute, Esc to return to navigation
  - [ ] Command history with up/down arrows

- [ ] **Command Parser Integration**
  - [ ] Reuse existing cobra command structure
  - [ ] Route embedded commands to CLI handlers
  - [ ] Capture output and display in TUI context

### Phase 2: Enhanced Command Experience  
- [ ] **Tab Completion**
  - [ ] Integrate with existing shell completion system
  - [ ] Context-aware suggestions based on current view
  - [ ] Archive ID completion for selected items

- [ ] **Command Feedback**
  - [ ] Display command results in TUI context
  - [ ] Progress bars for long-running operations
  - [ ] Error handling with visual feedback

### Phase 3: Default Interface Transition
- [ ] **Smart Default Detection**
  - [ ] Launch TUI by default for interactive terminals
  - [ ] Use CLI mode for pipes and scripts
  - [ ] `--cli` flag for explicit CLI mode

- [ ] **Migration Documentation**
  - [ ] User guide for TUI-first workflow
  - [ ] Script compatibility documentation
  - [ ] Migration path for existing users

### Dependencies
- **7EP-0010**: Interactive TUI Application ✅ (completed) - provides visual foundation
- **7EP-0007**: Enhanced MAS Operations 🔄 (in progress) - provides command backend
- **Current CLI Implementation** ✅ - provides command handlers to reuse

## Testing Strategy

### Acceptance Criteria
- [ ] All CLI commands accessible via embedded command line
- [ ] Tab completion works in command mode
- [ ] Command execution updates visual interface immediately
- [ ] Performance remains responsive with embedded command processing
- [ ] Script compatibility maintained with `--cli` flag
- [ ] TUI launches by default for interactive sessions

### Test Scenarios
- **Command Mode Testing**: All CLI commands work in embedded mode
- **Integration Testing**: Visual interface updates after command execution
- **Compatibility Testing**: Scripts continue working with `--cli` flag
- **Performance Testing**: Command line processing doesn't slow visual interface

## Migration/Compatibility

### Breaking Changes
**Phase 1-2:** None - TUI is additive feature  
**Phase 3:** Default behavior changes (TUI first) but `--cli` preserves compatibility

### Upgrade Path
- Current TUI becomes foundation for embedded command line
- CLI commands become embeddable without modification
- Gradual transition to TUI-first experience

### Backward Compatibility
- All existing CLI commands preserved via `--cli` flag
- Script integration maintained
- Configuration and data formats unchanged

## 🎨 Simple Interface Design Principles

### Terminal-Friendly Aesthetics
- **Minimal borders** - Clean lines, no heavy box drawing
- **Text-based indicators** - Status with simple characters (✓ ? X)
- **Readable layouts** - Spacious but not wasteful  
- **Theme consistency** - Colors enhance, don't overwhelm

### Example Simple Interface
```
7zarch-go                                              Staging | Remote (TrueNAS)

Local Staging (3 files, 267 MB)
> episode-424.7z          89 MB   2h ago   ✓ ready → upload
  episode-423.7z          92 MB   1d ago   ✓ needs upload
  vacation-photos.7z     3.8 GB   1w ago   ? recompress?

Remote Storage (1,247 files, 14.2 TB) via Tailscale
  episode-422.7z          87 MB   2d ago   ✓ synced
  episode-421.7z          91 MB   3d ago   ✓ synced  
  episode-420.7z          89 MB   4d ago   ✓ synced

[Tab] Switch view  [Enter] Details  [:] Command mode  [q] Quit
:
```

## Future Considerations

### Compression Intelligence
- **Profile analysis** - Detect suboptimal compression choices
- **Efficiency scoring** - Rate compression effectiveness  
- **Batch optimization** - Recompress old archives with new profiles
- **Savings estimation** - Show potential space/time savings

### Remote Storage Integration  
- **Tailscale URL management** - Easy URL copying and sharing
- **Bandwidth optimization** - Smart push/pull based on connection
- **Integrity verification** - Checksum validation without download
- **Directory traversal** - Browse remote folder structures

## References

- **Builds on**: 7EP-0010 Interactive TUI Application (completed foundation)
- **Integrates with**: 7EP-0007 Enhanced MAS Operations (command backend)
- **Enables**: TUI-first archive management experience
- **Inspired by**: Vim command mode, modern TUI applications with embedded CLIs
