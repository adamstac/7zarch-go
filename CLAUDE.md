# Claude Code (CC) Context & Bootup Guide

**Purpose**: Quick context loading for Claude Code sessions on the 7zarch-go project.
**Maintain**: Update this file at session end with important context changes.
**Location**: `/CLAUDE.md` (root of project)

## 👥 Who's Who

### Human Team
- **Adam Stacoviak** (@adamstac) - Project owner, makes architectural decisions, prefers simplicity
  - Likes: Clean design, Charmbracelet tools, thoughtful UX
  - Style: Direct feedback, big ideas, a fan of document driven development
  - Timezone: n/a

### AI Team
- **CC (Claude Code)** - You! Primary development assistant
  - Responsibilities: Feature implementation, bug fixes, documentation
  - Strengths: Display systems, infrastructure, deep technical work

- **AC (Augment Code)** - Sibling AI, handles parallel work
  - Responsibilities: User-facing features, refinements, overnight deep work
  - Current: Working on PR #9 (list filters), potentially 7EP-0010 TUI tonight
  - Communication: Via PR descriptions, commit messages, and `/docs/development/`

- **Amp (Sourcegraph)** - Advanced architectural analysis ⭐⭐⭐⭐⭐
  - Responsibilities: Strategic analysis, performance optimization, architectural insights
  - **PROVEN EXCELLENCE**: 7EP-0014 delivered exceptional foundation gap analysis
  - **Confidence Level**: **VERY HIGH** - Demonstrated deep codebase understanding and strategic thinking
  - **Impact**: Identified critical path that unblocks 7EP-0007, 7EP-0010, and production readiness

- **CR (CodeRabbit)** - Automated code reviewer
  - Triggers: On all PRs automatically
  - Purpose: Catches issues, suggests improvements
  - Config: `.coderabbit.yaml`

- **Amp-Principal** - Senior Technical Architect (NEW 2025-08-13)
  - Role: Technical review, architecture oversight, code quality governance
  - Focus: High-level review, documentation standards, process leadership
  - Activation: Use prompt structure in `/docs/development/amp-principal-activation-guide.md`
  - Scope: PR reviews, architecture evaluation, technical standards

## 📍 Key Locations

### Documentation
- `/docs/7eps/` - Enhancement proposals (our roadmap)
- `/docs/development/pr-merge-roadmap.md` - Current PR status and priorities
- `/docs/development/tomorrow-plan.md` - Daily planning
- `/docs/development/amp-principal-role-definition.md` - Amp-Principal role specification
- `/docs/development/amp-principal-activation-guide.md` - How to activate Amp-Principal for reviews
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
- `Makefile` - Build commands (`make dev`, `make dist`, `make validate`) 
- `.goreleaser.yml` - Professional build pipeline with Level 2 reproducibility

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
   make dev            # Build with Goreleaser and install
   ~/bin/7zarch-go list --dashboard  # Test display modes
   ```

## 🎯 Current Project State (as of 2025-08-13)

### Recently Completed  
- ✅ **7EP-0015 Code Quality Foundation** - Comprehensive quality improvements (JUST COMPLETED!)
- ✅ **7EP-0013 Build Pipeline** - Goreleaser + Level 2 reproducibility (PR #20 merged)
- ✅ **7EP-0005 Test Dataset System** - Comprehensive test infrastructure merged (PR #12)
- ✅ **7EP-0011 Lint Tightening** - Improved code quality standards merged (PR #19)
- ✅ **7EP-0009 Enhanced Display System** - 5 display modes (table, compact, card, tree, dashboard)
- ✅ **MAS Foundation** - Full ULID resolution, show, list, move commands
- ✅ **CI/CD Infrastructure** - Fixed all workflow issues, updated dependencies
- ✅ **Dependabot Cleanup** - 3 PRs merged, 1 incompatible PR properly closed

### 🌟 LATEST ACHIEVEMENT: 7EP-0015 Code Quality Foundation
**Status**: 🔍 **REVIEW** - [PR #25](https://github.com/adamstac/7zarch-go/pull/25) created for review
- **Standardized Error Handling**: All MAS commands use consistent patterns with helpful suggestions
- **Debug System**: `--debug` flag provides performance metrics (query time, memory, DB size)
- **Code Quality**: Extracted common patterns into `internal/cmdutil` reducing duplication
- **Test Coverage**: Added comprehensive test suites (50-100% coverage on core packages)
- **User Experience**: Enhanced help text and troubleshooting documentation

### Available for Next Work
- 🔄 **7EP-0010 TUI** - Ready for implementation (AC potentially working overnight)
- 🔄 **7EP-0007 Enhanced MAS Ops** - Ready to start (build infrastructure resolved)
- 🎯 **New strategic work** - Foundation is solid, ready for advanced features

### Current Status
- ✨ **Professional build pipeline** - Goreleaser with reproducible builds
- ✨ **Quality codebase** - Standardized errors, debug system, comprehensive tests
- ✨ **Clean infrastructure** - All PRs resolved, CI functional, dependencies current
- ✨ **Team coordination** - Amp completed 7EP-0014, AC available for 7EP-0010
- ✨ **Ready for advanced features** - Solid foundation enables complex development

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
- **Always create feature branches for new work** - Never work directly on main
- **Branch from main** for all new features, not from other feature branches  
- **Keep branches focused** - One branch per 7EP or major feature
- **Clean working directory** before starting new feature work
- Squash merge PRs with branch deletion
- **GPG SIGNING REQUIRED**: All commits to the remote repo must be GPG signed. Any unsigned commits must be squashed with a GPG signed commit before pushing
- **NO SIGNATURES**: Don't add "🤖 Generated with Claude Code" or "Co-Authored-By" to commits
- No Claude mentions in commits (Adam will say "no claude mention" if needed)
- Comprehensive commit messages with what and why

## 🛠️ Common Commands

### Build & Test
```bash
# Goreleaser build system (Level 2 reproducible) - JUST IMPLEMENTED!
make dev            # Build and install to ~/bin
make dist           # Build for current platform  
make validate       # Validate Goreleaser config
make release        # Create release (CI only)

# Legacy build system (still available)
make build          # Build binary
make test           # Run tests
make lint           # Run linter
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
| 0001 | Trash Management | ✅ Complete | AC | Merged PR #10 |
| 0002 | CI Integration | ✅ Complete | CC | Merged PR #11 |
| 0003 | Database Migrations | 🟡 Draft | AC | Not started |
| 0004 | MAS Foundation | ✅ Complete | AC | Merged |
| 0005 | Test Dataset | ✅ Complete | CC | Merged PR #12 |
| 0006 | Performance Testing | ✅ Complete | CC | Merged |
| 0007 | Enhanced MAS Ops | 🟢 Ready | AC/CC | Foundation complete, ready to start |
| 0008 | Depot Actions | ✅ Complete | CC | Merged |
| 0009 | Enhanced Display | ✅ Complete | CC | Merged |
| 0010 | Interactive TUI | 🟢 Ready | AC | Guide prepared, ready for implementation |
| 0011 | Lint Tightening | ✅ Complete | CC | Merged PR #19 |
| 0013 | Build Pipeline | ✅ Complete | CC | Merged PR #20 - Goreleaser + reproducible builds |
| 0014 | Critical Foundation | ✅ Complete | Amp | Exceptional strategic analysis completed |
| 0015 | Code Quality | 🔍 Review | CC | PR #25 - comprehensive quality improvements |

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
- 2025-08-13: **AMP-PRINCIPAL ROLE CREATED** 🎯 - Senior Technical Architect replicant for technical review
- **NEW CAPABILITY**: High-level code review, architecture oversight, documentation standards
- **ACTIVATION**: Use prompt structure in `/docs/development/amp-principal-activation-guide.md`
- **PURPOSE**: Complement strategic Amp with senior technical leadership and quality governance
- **READY FOR**: Technical review of 7EP-0007 Phase 2 or architecture evaluation

---

**Remember**: You're CC (Claude Code). You build things. You ship features. You write clean code without unnecessary comments. You're direct and concise. And sometimes, when Adam says "that'll do pig," you know you've hit the sweet spot. 🐷

**Last Updated**: 2025-08-13 by CC
