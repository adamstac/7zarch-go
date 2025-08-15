# Team Context & Project Overview

**Purpose**: Shared context for all AI team members about project structure, team, and current state  
**Audience**: CC, AC, Amp, and any new AI team members  
**Framework**: Document Driven Development (7EP-0017)

## 👥 Team Structure

### Human Team
- **Adam Stacoviak** (@adamstac) - Project owner, makes architectural decisions, prefers simplicity
  - Likes: Clean design, Charmbracelet tools, thoughtful UX
  - Style: Direct feedback, big ideas, document driven development

### AI Team
- **CC (Claude Code)** - Backend features, infrastructure, technical implementation
- **AC (Augment Code)** - User experience, frontend features, quality assurance  
- **Amp-s** - Strategic planning, business impact, executive leadership
- **Amp-t** - Technical architecture, code quality governance
- **CR (CodeRabbit)** - Automated code review (triggers on all PRs)

## 📍 Project Architecture

### Code Structure
```
/cmd/               - CLI commands (list, show, create, etc.)
/internal/
  ├── storage/      - SQLite registry, archive management
  ├── tui/          - Bubbletea interface (9 themes)
  ├── display/      - 5 display modes (table, compact, card, tree, dashboard)
  ├── query/        - Saved queries and search engine
  ├── batch/        - Multi-archive operations
  └── config/       - Configuration management
/docs/
  ├── 7eps/         - Enhancement proposals (roadmap)
  ├── development/  - Team coordination and assignments
  └── reference/    - Command documentation
```

### Database
- **SQLite registry**: `~/.7zarch-go/registry.db`
- **ULID-based**: User-facing archive IDs (01JEX...)
- **Managed vs External**: Tracks storage location type
- **Soft deletes**: Status field (present/missing/deleted)

## 🎯 Current Project State

### Completed Phases
- **Foundation Phase**: ✅ Production-ready CLI with TUI
- **Advanced Features Phase**: ✅ Query/search/batch operations

### Current Phase
- **Strategic Direction Planning** - Awaiting Adam's next focus decision

### Team Availability
- **CC**: Available for strategic assignment
- **AC**: Available for strategic assignment  
- **Amp**: Framework oversight and strategic coordination
- **Adam**: Strategic priority decision needed

## 📋 Key Documentation Locations

### Daily Operations
1. `docs/development/NEXT.md` - Current cross-team priorities
2. `docs/development/[ROLE].md` - Role-specific assignments
3. Active 7EPs - Sprint coordination context

### Reference
4. `docs/7eps/index.md` - Long-term feature roadmap
5. `/AGENT.md` - Technical build/test/style patterns
6. `docs/development/README.md` - DDD framework usage

## 🔄 Workflow Patterns

### Communication
- PR descriptions for major updates
- Commit messages reference 7EP numbers
- Update role documents for coordination needs
- Use NEXT.md for cross-team dependencies

### Development
- Feature branches: `feat/7ep-XXXX-description`
- GPG signed commits required
- Squash merge PRs with branch deletion
- Clean working directory before new work

## 📊 7EP Status Quick Reference

| 7EP | Title | Status | Owner | Notes |
|-----|-------|--------|-------|-------|
| 0001 | Trash Management | ✅ Complete | AC | Merged PR #10 |
| 0002 | CI Integration | ✅ Complete | CC | Merged PR #11 |
| 0003 | Database Migrations | 🟡 Draft | AC | Not started |
| 0004 | MAS Foundation | ✅ Complete | AC | Merged |
| 0005 | Test Dataset | ✅ Complete | CC | Merged PR #12 |
| 0006 | Performance Testing | ✅ Complete | CC | Merged |
| 0007 | Enhanced MAS Ops | ✅ Complete | CC | All 3 phases complete |
| 0008 | Depot Actions | ✅ Complete | CC | Merged |
| 0009 | Enhanced Display | ✅ Complete | CC | Merged |
| 0010 | Interactive TUI | 🟢 Ready | AC | Guide prepared, ready for implementation |
| 0011 | Lint Tightening | ✅ Complete | CC | Merged PR #19 |
| 0013 | Build Pipeline | ✅ Complete | CC | Merged PR #20 - Goreleaser + reproducible builds |
| 0014 | Critical Foundation | ✅ Complete | Amp | Exceptional strategic analysis completed |
| 0015 | Code Quality | ✅ Complete | CC | Merged - comprehensive quality improvements |
| 0017 | DDD Framework | ✅ Complete | CC | All operational documents and structure complete |
| 0018 | Static Blog Generator | 🟡 Draft | CC | Awaiting Adam's decision |

## 🚨 Team Troubleshooting

### Common Issues
- **Build broken?** Check recent merges, try `git bisect`
- **PR conflicts?** Pull main, rebase feature branch  
- **Dependabot spam?** Can be batched or ignored temporarily
- **Team coordination issues?** Update NEXT.md and role documents
- **Can't find technical commands?** Check `AGENT.md` for build/test/lint patterns

### Workflow Support
- **Session startup:** Use `docs/development/actions/BOOTUP.md`
- **Session shutdown:** Use `docs/development/actions/SHUTDOWN.md`
- **Commit workflows:** Use `docs/development/actions/COMMIT.md`
- **PR workflows:** Use `docs/development/actions/MERGE.md`
- **New features:** Use `docs/development/actions/NEW-FEATURE.md`

---

This shared context eliminates duplication while providing essential team and project knowledge for all AI agents.
