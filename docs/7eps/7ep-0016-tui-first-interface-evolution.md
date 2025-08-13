# 7EP-0016: TUI-First Interface Evolution

**Status:** Draft  
**Author(s):** Amp (Sourcegraph), Adam Stacoviak  
**Assignment:** Future Planning  
**Difficulty:** 4 (architectural - interface paradigm shift)  
**Created:** 2025-08-13  
**Updated:** 2025-08-13  

## Executive Summary

Evolve 7zarch-go from CLI-with-TUI to TUI-first application with embedded Vim-style command line, making the visual interface the primary user experience while maintaining full CLI compatibility for automation and scripting.

## Evidence & Reasoning

**Current Success:**
- ✅ TUI experiment successful - beautiful themes, simple navigation, user-tested
- ✅ Multiple entry points working (`browse`, `ui`, `i`, `tui`)
- ✅ Solid foundation with viewport framework and proper Bubbletea architecture
- ✅ User feedback positive - "very good and simple start, great base to build on"

**User Vision:**
- TUI becomes the entire application interface
- Vim-style command line embedded within TUI
- All CLI commands accessible from within visual interface
- Seamless workflow for podcast archival management

**Strategic Opportunity:**
- Differentiate from all other CLI archive tools
- Create compelling visual demo and user experience
- Natural evolution path from current success

## Use Cases

### Primary Use Case: TUI-First Experience
```bash
# Default launch goes to visual interface
7zarch-go

# Visual interface with embedded command line
┌─ Archive Browser ─────────────────────────────┐
│ Archives: 247 (12.4 TB)                       │
│                                               │ 
│ > episode-423.7z        89 MB    2h ago    ✓  │
│   episode-422.7z        92 MB    1d ago    ✓  │
│   vacation-photos.7z   3.8 GB    1w ago    ✓  │
│                                               │
│ [Visual navigation and content]               │
└───────────────────────────────────────────────┘
│ :create episode-424 --profile media          │  <-- Vim-style command line
└───────────────────────────────────────────────┘

# CLI mode still available for automation
7zarch-go --cli list --output json
```

### Secondary Use Cases
- **Command discovery** - Tab completion in embedded command line
- **Complex operations** - Command line for advanced workflows  
- **Hybrid workflows** - Visual browsing + command execution
- **Script integration** - CLI mode preserves automation compatibility

## Technical Design

### Architecture Overview
```go
type TUIApp struct {
    // Current components (keep)
    archives     []*storage.Archive
    viewport     viewport.Model
    theme        Theme
    
    // New command line components
    commandMode  bool           // Toggle between nav and command mode
    commandLine  textinput.Model // Vim-style command input
    commandHist  []string       // Command history
    suggestions  []string       // Tab completion suggestions
    
    // Command execution
    executor     CommandExecutor // Shared with CLI mode
}

type CommandExecutor interface {
    Execute(cmd string, args []string) tea.Cmd
    Complete(partial string) []string
    Validate(cmd string) error
}
```

### Interface Modes
1. **Navigation Mode** (default)
   - Arrow keys navigate archives
   - Single letters trigger actions
   - Visual feedback and selection

2. **Command Mode** (`:` key like Vim)
   - Embedded command line at bottom
   - Full CLI command access
   - Tab completion and history
   - Execute and return to navigation

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

## Future Considerations

### Advanced Command Integration
- **Visual command builder** - Build commands by selecting archives in TUI
- **Command templates** - Save frequently used command patterns
- **Macro support** - Record and replay command sequences
- **Multi-pane interface** - Command output in separate pane

### User Experience Enhancements
- **Command suggestion** - AI-like suggestions based on context
- **Visual command feedback** - Progress and results integrated with interface
- **Keyboard shortcuts** - Quick access to common commands without typing
- **Customizable interface** - User-defined layouts and command mappings

## References

- **Builds on**: 7EP-0010 Interactive TUI Application (completed foundation)
- **Integrates with**: 7EP-0007 Enhanced MAS Operations (command backend)
- **Enables**: TUI-first archive management experience
- **Inspired by**: Vim command mode, modern TUI applications with embedded CLIs
