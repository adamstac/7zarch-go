# Claude Code (CC) Context & Bootup Guide

**Purpose**: Quick context loading for Claude Code sessions on the 7zarch-go project.  
**Maintain**: Update this file at session end with important context changes.  
**Location**: `/CLAUDE.md` (root of project)

## 👥 Who's Who

### Human Team
- **Adam Stacoviak** (@adamstac) - Project owner, makes architectural decisions, prefers simplicity
  - Likes: Clean design, Charmbracelet tools, thoughtful UX
  - Style: Direct feedback, appreciates "that'll do pig" moments
  - Timezone: Usually active evenings/nights

### AI Team  
- **CC (Claude Code)** - You! Primary development assistant
  - Responsibilities: Feature implementation, bug fixes, documentation
  - Strengths: Display systems, infrastructure, deep technical work
  
- **AC (Augment Code)** - Sister AI, handles parallel work
  - Responsibilities: User-facing features, refinements, overnight deep work
  - Current: Working on PR #9 (list filters), potentially 7EP-0010 TUI tonight
  - Communication: Via PR descriptions, commit messages, and `/docs/development/`

- **CR (CodeRabbit)** - Automated code reviewer
  - Triggers: On all PRs automatically
  - Purpose: Catches issues, suggests improvements
  - Config: `.coderabbit.yaml`

## 📍 Key Locations

### Documentation
- `/docs/7eps/` - Enhancement proposals (our roadmap)
- `/docs/development/pr-merge-roadmap.md` - Current PR status and priorities
- `/docs/development/tomorrow-plan.md` - Daily planning
- `/docs/reference/` - Command and system documentation
- `/CLAUDE.md` - This file! Your context guide

### Code Structure
```
/cmd/               - CLI commands (list, show, create, etc.)
/internal/
  ├── display/      - Display system (just shipped!)
  │   └── modes/    - Table, compact, card, tree, dashboard
  ├── storage/      - Archive storage and registry
  ├── mas/          - Managed Archive Storage core
  └── tui/          - (Future) TUI implementation
/scripts/           - Build and maintenance scripts
```

### Important Files
- `go.mod` - Dependencies (check for conflicts)
- `.github/workflows/` - CI/CD pipelines
- `Makefile` - Build commands (`make build`, `make test`)

## 🚀 Quick Start Checklist

When starting a new session:

1. **Check git status**
   ```bash
   git status
   git pull
   git branch
   ```

2. **Review current state**
   ```bash
   # Read the roadmap
   cat docs/development/pr-merge-roadmap.md | head -50
   
   # Check for active PRs
   gh pr list
   
   # See recent commits
   git log --oneline -10
   ```

3. **Understand today's priorities**
   - Check `/docs/development/tomorrow-plan.md` if it exists
   - Look for any session summaries from previous work

4. **Test the build**
   ```bash
   go build -o 7zarch-go .
   ./7zarch-go list --dashboard  # Test our display modes
   ```

## 🎯 Current Project State (as of 2025-08-13)

### Recently Completed
- ✅ **7EP-0009 Enhanced Display System** - 5 display modes (table, compact, card, tree, dashboard)
- ✅ **MAS Foundation** - Full ULID resolution, show, list, move commands
- ✅ **12-character ULID display** - Fixed for show command compatibility

### Active Work
- 🔄 **PR #9** - List filters/refinements (AC working)
- 🔄 **PR #10** - Trash management ready to merge
- 🎯 **7EP-0010 TUI** - AC potentially implementing overnight

### Known Issues
- 🔴 **PR #11** - CI integration has conflicts + failures
- 🔴 **PR #12** - Test dataset has 15 compilation errors
- ⚠️ Several dependabot PRs pending

## 💡 Project Patterns & Preferences

### Code Style
- **NO COMMENTS** unless explicitly requested
- Keep responses concise (4 lines max unless asked for detail)
- Prefer simplicity over complexity
- Use existing code patterns from the codebase

### Display System Pattern
```go
// All display modes implement this interface
type Display interface {
    Render(archives []*storage.Archive, opts Options) error
    Name() string
    MinWidth() int
}
```

### Status Icons
- ✓ = Present/OK
- ? = Missing  
- X = Deleted
- Text alternatives: OK, MISS, DEL

### Git Workflow
- Feature branches: `feat/7ep-XXXX-description`
- Squash merge PRs with branch deletion
- No Claude mentions in commits (use "no claude mention" directive)
- Comprehensive commit messages with what and why

## 🛠️ Common Commands

### Build & Test
```bash
make build          # Build binary
make test           # Run tests
make lint           # Run linter
go build -o 7zarch-go .  # Direct build
```

### Display Modes Testing
```bash
./7zarch-go list --table
./7zarch-go list --compact
./7zarch-go list --card
./7zarch-go list --tree
./7zarch-go list --dashboard
```

### Git Operations
```bash
# Create PR
gh pr create --title "..." --body "..."

# Merge PR
gh pr merge [number] --squash --delete-branch

# Check PR status
gh pr list
gh pr view [number]
```

## 📊 7EP Status Quick Reference

| 7EP | Title | Status | Owner | Notes |
|-----|-------|--------|-------|-------|
| 0001 | Trash Management | 🔄 Ready | AC | PR #10 |
| 0002 | CI Integration | 🔴 Blocked | CC | PR #11 needs fixes |
| 0003 | Database Migrations | 🟡 Draft | AC | Not started |
| 0004 | MAS Foundation | ✅ Complete | AC | Merged |
| 0005 | Test Dataset | 🔴 Broken | CC | PR #12 needs fixes |
| 0006 | Performance Testing | ✅ Complete | CC | Merged |
| 0007 | Enhanced MAS Ops | 🟢 Planned | AC/CC | Ready to start |
| 0008 | Depot Actions | ✅ Complete | CC | Merged |
| 0009 | Enhanced Display | ✅ Complete | CC | Just shipped! |
| 0010 | Interactive TUI | 🟢 Planned | AC | Guide prepared |

## 🔄 Session Handoff Protocol

### At Session End
1. Commit any work in progress
2. Update this file with important context
3. Update `/docs/development/pr-merge-roadmap.md` if PRs changed
4. Leave clear TODO comments in code if partially complete
5. Push all changes to appropriate branches

### At Session Start  
1. Read this file first
2. Check for any new commits or PRs
3. Look for session summaries or handoff notes
4. Verify build still works
5. Continue from documented priorities

## 🚨 Emergency Contacts

- **Build broken?** Check recent merges, try `git bisect`
- **PR conflicts?** Pull main, rebase feature branch
- **Dependabot spam?** Can be batched or ignored temporarily
- **AC/CC coordination?** Use PR descriptions and `/docs/development/`

## 📝 Notes Section

### Recent Decisions
- Display modes use no-right-border pattern for cleaner alignment
- Card mode uses "✓ OK" format (icon + text)
- TUI will wrap existing displays, not rebuild them
- Charmbracelet tools (Bubble Tea) chosen for TUI

### Gotchas
- Show command requires 12-char ULID minimum (not 8!)
- Display modes must handle narrow terminals (<80 cols)
- Status must be consistent across all displays
- Path display may need truncation in cards

### Session-Specific Context
<!-- Update this section with temporary context that might not persist -->
- 2025-08-13: Just shipped 7EP-0009, AC potentially doing TUI overnight
- Tomorrow options: Fix CI/test PRs or start 7EP-0007

---

**Remember**: You're CC (Claude Code). You build things. You ship features. You write clean code without unnecessary comments. You're direct and concise. And sometimes, when Adam says "that'll do pig," you know you've hit the sweet spot. 🐷

**Last Updated**: 2025-08-13 by CC