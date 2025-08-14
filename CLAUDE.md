# Claude Code (CC) Context Guide

## 📍 Essential References for Claude Code

**Technical Competency**: `/AGENT.md`
- Build/test/lint commands and development workflow
- Architecture overview and codebase structure
- Code style conventions and patterns specific to the project

**Role Context & Assignments**: `/docs/development/CLAUDE.md`
- Current assignments and priorities
- Team coordination context and communication patterns
- Role-specific responsibilities and focus areas

**Session Startup**: `/BOOTUP.md`
- Standardized startup sequence for AI team members
- Current project state and immediate priorities

---

**Usage**: Consult both AGENT.md (for technical patterns) and docs/development/CLAUDE.md (for coordination context) when starting work.

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

- **Amp-s** - Senior Strategic Architect (Renamed 2025-08-13) ⭐⭐⭐⭐⭐
  - Responsibilities: Product strategy, business impact, executive leadership, roadmap planning
  - Focus: Strategic vision, user value, competitive positioning, resource allocation
  - Activation: `Switching to Amp-s role for strategic planning.`
  - **PROVEN EXCELLENCE**: 7EP-0014 delivered exceptional foundation gap analysis
  - **Confidence Level**: **VERY HIGH** - Demonstrated strategic thinking and product vision

- **CR (CodeRabbit)** - Automated code reviewer
  - Triggers: On all PRs automatically
  - Purpose: Catches issues, suggests improvements
  - Config: `.coderabbit.yaml`

- **Amp-t** - Senior Technical Architect (NEW 2025-08-13)
  - Role: Technical review, architecture oversight, code quality governance
  - Focus: High-level review, documentation standards, process leadership
  - Activation: `Switching to Amp-t role for technical review.`
  - Scope: PR reviews, architecture evaluation, technical standards

## 📍 Key Locations

### Documentation
- `/docs/7eps/` - Enhancement proposals (our roadmap)
- `/docs/development/pr-merge-roadmap.md` - Current PR status and priorities
- `/docs/development/tomorrow-plan.md` - Daily planning
- `/docs/development/AMP.md` - Unified Amp-s (strategic) and Amp-p (technical) role documentation
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

2.5. **CHECK OPERATIONAL PRIORITIES** (DDD Framework)
   ```bash
   # Personal current assignments and coordination
   cat docs/development/CLAUDE.md | head -20
   
   # Shared team priorities and blockers  
   cat docs/development/NEXT.md | head -30
   
   # Active 7EP coordination context
   grep -l "Status.*ACTIVE\|In Progress" docs/7eps/*.md | xargs ls -la
   ```

3. **Understand today's priorities**
   - Check `/docs/development/NEXT.md` for shared team coordination
   - Review `/docs/development/CLAUDE.md` for personal assignments
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

### Search Engine Testing (NEW - 7EP-0007 Phase 2)
```bash
# Test search functionality
./7zarch-go search reindex                        # Rebuild index
./7zarch-go search query "test"                   # Full-text search
./7zarch-go search query --field=name "backup"    # Field-specific
./7zarch-go search query --field=profile "media"  # Profile search
./7zarch-go search query --regex ".*\\.7z$"      # Regex patterns

# Performance testing
time ./7zarch-go search query "performance test"  # Should be <1ms
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
| 0007 | Enhanced MAS Ops | 🚀 Phase 2 Complete | CC | Search engine ~60-100µs performance, PR #27 |
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
- **7EP-0007 Phase 2 Complete**: Search engine with exceptional ~60-100µs performance
- Search architecture: Inverted index + LRU cache + thread-safe design
- All searchable fields indexed: name, path, profile, metadata

### Gotchas
- Show command requires 12-char ULID minimum (not 8!)
- Display modes must handle narrow terminals (<80 cols)
- Status must be consistent across all displays
- Path display may need truncation in cards

### Session-Specific Context
<!-- Update this section with temporary context that might not persist -->
- 2025-08-13: **DUAL AMP ROLES CREATED** 🎯 - Clear strategic vs technical separation
- **AMP-S (Strategic)**: Product strategy, business impact, executive leadership, roadmap planning  
- **AMP-T (Technical)**: Code quality, architecture review, technical standards, implementation oversight
- **ACTIVATION**: See `AMP.md` for complete role documentation and activation examples
- **PURPOSE**: Clear role separation - strategic vision (Amp-s) + technical execution (Amp-t)
- 2025-08-14: **7EP-0007 FULLY COMPLETE** ✅ - All 3 phases merged to main
- **PHASE 2**: ✅ MERGED - Search engine (~60-100µs performance, 5000x faster than target)
- **PHASE 3**: ✅ MERGED - Batch operations with enterprise-grade concurrency and safety
- **FEATURES DELIVERED**: Complete Query → Search → Batch workflow operational
- **PR STATUS**: PR #27 & #28 both merged with dual leadership approval
- **TECHNICAL ACHIEVEMENT**: 4,338 total lines added across all phases
- **TRANSFORMATION COMPLETE**: Basic archive manager → Enterprise power user command center
- **PRODUCTION READY**: Full workflow from discovery to bulk operations
- **STRATEGIC IMPACT**: Complete enterprise archive management solution in production

---

**Remember**: You're CC (Claude Code). You build things. You ship features. You write clean code without unnecessary comments. You're direct and concise. And sometimes, when Adam says "that'll do pig," you know you've hit the sweet spot. 🐷

**Last Updated**: 2025-08-13 by CC (Phase 2 search engine completion)
