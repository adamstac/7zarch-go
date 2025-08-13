# 7zarch-go Development Guide

This document provides shared context for all AI agents and human contributors working on 7zarch-go.

## AI Agent Roles & Coordination

### Claude Code (CC)
**Primary Responsibilities:**
- Infrastructure and tooling (CI/CD, testing frameworks)
- Code quality standards and automation
- Documentation workflow and processes
- Technical architecture and system design
- Build system and development tooling

**Current Focus:** CI integration (7EP-0002), testing automation

### Augment Code (AC)
**Primary Responsibilities:**
- User-facing command implementations
- Feature development and UX improvements
- Command interface design
- User workflow optimization

**Current Focus:**
- Trash management system implementation (7EP-0001: restore, trash list, trash purge)
- Deeper CodeRabbit (CR) integration workflows

### CodeRabbit (CR)
**Primary Responsibilities:**
- Automated code review on all PRs
- Code quality enforcement
- Security and performance analysis
- Consistency checking across codebase

**Integration Status:** AC is working on deeper CR integration patterns

## Development Patterns & Standards

### Command Implementation Pattern
```go
// Standard command registration in main.go
rootCmd.AddCommand(cmd.MasShowCmd())    // Top-level command access
rootCmd.AddCommand(cmd.MasDeleteCmd())  // Direct command registration
```

### Registry Integration Pattern
```go
// Archive status tracking
archive.Status = "deleted"
archive.DeletedAt = &now
archive.OriginalPath = archive.Path  // Preserve for restore
archive.Path = trashPath             // Move to trash location
```

### Error Handling Standard
```go
if err != nil {
    return fmt.Errorf("descriptive context: %w", err)
}
```

## Architecture Decisions

### Managed Archive Storage (MAS)
- **Database**: SQLite with 0600 permissions
- **Location**: `~/.7zarch-go/`
- **Registry**: Tracks metadata, status, relationships
- **Migration**: Automatic schema updates

### Command Design Philosophy
- **Top-level commands** preferred over deep subcommand nesting
- **ULID resolution** for user-friendly ID references
- **Status indicators** in list outputs for clear visual feedback
- **Configuration integration** - no hardcoded values

### Display Standards
```bash
# Status formatting
📦 archive-name - ✅ Uploaded to destination
🗑️  deleted-name - Deleted 2025-08-11 14:23:01
   Auto-purge: 29 days (2025-09-10)

# ID display for copy/paste
   ID: 01K2E33XW4HTX7RVPS9Y6CRGDY
```

## 7EP Enhancement Proposal Process

### Workflow
1. **Draft**: Create 7EP using template for significant features
2. **Assignment**: CC handles infrastructure, AC handles user features
3. **Implementation**: Follow 7EP technical design
4. **Documentation**: Transition 7EP learnings to user docs
5. **Archive**: 7EPs remain for historical reference

### Current Active 7EPs
- **7EP-0001**: Trash Management System (AC assigned, difficulty 3)
- **7EP-0002**: CI Integration & Automation (CC assigned, difficulty 2)
- **7EP-0003**: Database Migrations & Schema Management (AC assigned, difficulty 2)
- **7EP-0004**: MAS Foundation Implementation (AC assigned, difficulty 4)

## Code Quality Requirements

### Testing Standards
- Unit tests for all new commands
- Integration tests for registry operations
- Build verification across platforms (Linux, macOS, Windows)
- Lint and format checks must pass

### Documentation Standards
- Update README.md for user-facing changes
- Reference 7EPs in commit messages
- Include examples in all command documentation
- Follow documentation workflow for feature completion

### Commit Message Format
```
feat: implement trash list command

References 7EP-0001 for trash management system.
Adds filtering and status display for deleted archives.

🤖 Generated with [Claude Code](https://claude.ai/code)
Co-Authored-By: Claude <noreply@anthropic.com>
```

## Current Development Status

### Completed Features
- ✅ Managed Archive Storage (MAS) with SQLite registry
- ✅ Enhanced list display with ULID support and status indicators
- ✅ Delete/move commands with soft delete and trash management
- ✅ Auto-purge countdown with configurable retention days
- ✅ 7EP system for formal enhancement proposals
- ✅ Documentation workflow process

### Active Development
- 🔄 **Trash Management Commands** (AC): restore, trash list, trash purge (7EP-0001)
- 🔄 **Database Migrations** (AC): Safe schema evolution system (7EP-0003)
- 🔄 **CI Integration** (CC): GitHub Actions workflows (7EP-0002)
- 🔄 **CodeRabbit Integration** (AC): Enhanced PR review workflows

## Development Roadmap

### **Current Sprint: MAS Foundation** (AC Priority)
**7EP-0004: MAS Foundation Implementation** (AC assigned, difficulty 4)
- Status: 📋 7EP Complete, ready for implementation
- Primary: ULID resolution, show command, enhanced list
- Supporting: Error handling standards, testing, documentation
- **See 7EP-0004 for detailed implementation checklist**

### **Next Sprint: Data Safety & Lifecycle**
**7EP-0003: Database Migration System** (AC assigned, difficulty 2)
- Status: 📋 7EP Complete, ready for implementation
- Blockers: None (can run parallel with 7EP-0004)
- Focus: Safe schema evolution with backup/recovery

**7EP-0001: Trash Management System** (AC assigned, difficulty 3)
- Status: 📋 7EP Complete, schema ready
- Blockers: Show command patterns from 7EP-0004
- Focus: Restore, trash list, purge operations

### **Future Sprints: Advanced Features**
**Import Command** - Bulk registration of existing archives
- Status: 💭 Design needed (future 7EP)
- Blockers: Core MAS stability (7EP-0004 completion)

**Multi-Location Support** - Network and cloud storage
- Status: 💭 Concept documented in guides/
- Blockers: Core operations solid

**7EP-0002: CI Integration** (CC assigned, difficulty 2)
- Status: 📋 7EP Complete, ready for implementation
- Blockers: None (parallel work)
- Focus: GitHub Actions, automated testing

### **Development Milestones**

#### **Milestone 1: MAS Core Complete** (Target: 2-3 weeks)
- ✅ ULID resolution working
- ✅ Show command operational
- ✅ Enhanced list with filters
- **Success Criteria**: Users can manage archives by ID without paths

#### **Milestone 2: Production Ready** (Target: 4-6 weeks)
- ✅ Database migrations automated
- ✅ Trash management complete
- ✅ CI pipeline operational
- **Success Criteria**: Safe for daily use, automated testing

#### **Milestone 3: Advanced Workflows** (Target: 8-12 weeks)
- ✅ Import command for existing archives
- ✅ Multi-location support
- ✅ Performance optimization
- **Success Criteria**: Handles large archive collections efficiently

### **Risk Management**

#### **Technical Risks**
- **ULID Resolution Complexity**: Disambiguation UX might be complex
  - *Mitigation*: Start simple, iterate based on user feedback
- **Database Migration Safety**: Schema changes risk data loss
  - *Mitigation*: Comprehensive backup strategy, extensive testing
- **Performance**: Large registries might slow operations
  - *Mitigation*: Index optimization, lazy loading patterns

#### **Coordination Risks**
- **CC/AC Overlap**: Parallel development might create conflicts
  - *Mitigation*: Clear responsibility boundaries, regular sync
- **Feature Creep**: Advanced features might delay core functionality
  - *Mitigation*: Strict milestone focus, defer non-essential features

### **Success Metrics**

#### **Development Velocity**
- **7EP completion rate**: Target 1 major 7EP per 2-3 weeks
- **Implementation lag**: <1 week from 7EP acceptance to start
- **Bug escape rate**: <5% of features require post-release fixes

#### **Quality Indicators**
- **Test coverage**: >80% for core registry operations
- **Performance**: <100ms for common operations (list, show)
- **User experience**: No breaking changes to existing workflows

### Known Technical Debt
- Database migration edge cases (resolved via fresh install approach)
- Command discoverability (resolved via top-level registration)
- Configuration consistency (resolved via centralized config loading)

## Inter-Agent Coordination Protocol

### Feature Development
1. **CC** creates 7EP with technical design and assigns to appropriate agent
2. **AC** implements user-facing features following 7EP specification
3. **CR** reviews all PRs against established patterns and quality standards
4. **CC** handles CI/testing infrastructure and documentation processes

### Roadmap Management
When user requests "latest roadmap," synthesize from:
- **DEVELOPMENT.md**: Development priorities, milestones, technical roadmap
- **PRIVATE.md**: Strategic vision, business considerations, pivot points

Provide unified view with source attribution, then update appropriate locations as discussions evolve.

### Communication Channels
- **7EPs**: Formal feature specification and design documents
- **GitHub Issues**: Bug reports and enhancement discussions
- **PR Descriptions**: Implementation context for CR review
- **DEVELOPMENT.md**: Shared patterns and coordination (this file)

### Handoff Patterns
- Reference 7EP numbers in commits and PRs
- Update 7EP status as implementation progresses
- Document implementation learnings in 7EP completion notes
- Transition completed features to user documentation

## Build & Development Environment

### Makefile Targets
```bash
make build          # Build main binary
make test           # Run unit tests
make integration    # Integration tests
make lint           # Code linting and formatting
make build-all      # Multi-platform builds
```

### User Installation Pattern
```bash
# Development symlink approach
ln -sf $(pwd)/7zarch-go /usr/local/bin/7zarch-go
```

## Future Enhancement Areas

### Planned 7EPs
- **External Archive Registration**: Track non-managed archives
- **Backup Integration**: Cloud storage and rsync workflows
- **Performance Optimization**: Large-scale archive handling
- **Configuration Enhancement**: Advanced config management

### CodeRabbit Integration Improvements (AC Focus)
- Enhanced PR review context
- Automated code pattern enforcement
- Integration with 7EP compliance checking
- Quality gate automation

---

**Note**: This file serves as shared context for all AI agents. For private development notes and half-formed ideas, use PRIVATE.md (git-ignored).

## Resume Context Snapshot — 2025-08-12

This section summarizes the current development state, active PRs, and near-term next steps, suitable for handoff or quick onboarding.

- Overview
  - Focus area: 7EP-0004 (MAS foundation: resolver, show, list) with related move UX improvements
  - Secondary: 7EP-0002 (CI) prep and small refactors to improve code quality and testability

- Pull Requests
  - PR #5 — 7EP-0004: MAS foundation increments (move fix, list filters)
    - Review: CodeRabbit completed review; status SUCCESS with actionable notes
    - Applied locally: move command fixes
      - Base-name fallback when arc.Name is empty (uses filepath.Base(arc.Path))
      - If --to points to a directory, append archive name
      - Prevent overwrite; handle cross-device rename (EXDEV) with copy+remove fallback
      - Robust managed-path detection using filepath.Rel to the managed base
    - State: Changes build and tests pass locally; awaiting final push/merge after aligning with user’s latest mas_move.go adjustments
    - Next: Confirm desired final behavior vs. user edits, push update to feature/7ep-0004-mas-foundation, respond to remaining CR comments

  - PR #8 — prep(7EP-0004): resolver benchmark + table truncation
    - CI: Failed in GH Actions earlier
    - Cause (from logs): duplicate error type declarations in resolver vs. central errors; outdated helper method calls (e.g., ResolveID vs Resolve, Register vs Add)
    - Plan: align resolver to use central error types (internal/storage/errors.go); update helpers and tests to current method names; re-run CI

  - PR #9 — refactor(list): listFilters + applyFilters + tests
    - Scope: structural refactor only; no output/behavior changes
      - Introduced listFilters struct and refactored listRegistryArchives to accept it
      - Extracted applyFilters for status/profile/larger-than (simple sequential filters)
      - Added parseHumanDuration with support for d (days) and w (weeks)
    - Tests: added small unit tests (applyFilters, parseHumanDuration); later removed per user preference
    - Build: green locally and in repository tests

- Key Code Changes (landed locally and/or in PRs)
  - cmd/mas_move.go
    - Base-name fallback for unnamed archives; directory destination handling
    - Overwrite prevention; EXDEV copy+remove fallback
    - Managed-path detection via filepath.Rel for portability
  - cmd/list.go
    - applyFilters helper (status/profile/larger-than)
    - parseHumanDuration for d/w in addition to standard time units
    - Delegation to a single printGroupedArchives to avoid duplication; removed hidden dependency on printArchiveTable by iterating when needed
  - internal/storage
    - Consolidated error types in errors.go and removed duplicates from resolver.go
    - test_helpers: updated to use Registry.Add and Resolver.Resolve, adjusted imports

- Documentation
  - Resolved conflicts in docs/7eps/index.md
    - Adopted expanded table (adds Area and Updated columns) and kept 0004 status as “In Progress (90%)”
  - Resolved header conflict in docs/7eps/7ep-0004-mas-foundation.md, keeping the up-to-date “Current Status” and “Update Notes” sections

- Validation
  - go build ./... and go test ./... pass locally after refactors and move fixes
  - CodeRabbit review signals SUCCESS on PR #5; CI on PR #8 needs a follow-up patch as noted above

- Immediate Next Steps
  - PR #5: finalize which mas_move.go variant to keep (user’s latest vs. CR-refined), push changes, and address any final CodeRabbit comments
  - PR #8: push small fix aligning resolver/test helpers and re-run CI to green the pipeline
  - PR #9: proceed with review; minimal risk since no behavior changes were introduced

- Coordination Notes
  - CC (infrastructure) continues CI and documentation flow
  - AC (features) owns MAS commands and list refactors
  - CR (CodeRabbit) remains the automated review gate; actionable comments already integrated into move improvements

## Session Summary — 2025-08-12 (AC)

- Goals
  - Advance 7EP-0004 (MAS foundation) with small, low-risk refactors and address move command UX issues highlighted by CodeRabbit.
  - Check in on PR #5 and PR #8, unblock CI and reviews, and capture context for coordination.

- What I did
  - List command refactor (no behavior change):
    - Introduced listFilters struct; refactored listRegistryArchives to accept it.
    - Added applyFilters helper (status/profile/larger-than) to de-duplicate filtering patterns.
    - Added parseHumanDuration with support for d (days) and w (weeks) alongside standard units.
    - Reconciled output by delegating to printGroupedArchives; where table helper wasn’t available, iterated items explicitly to preserve output.
  - Tests:
    - Added small unit tests for applyFilters and parseHumanDuration; validated green locally.
  - Storage/resolver cleanup:
    - Removed duplicate error type definitions from resolver.go and used centralized types from errors.go.
    - Updated test helpers to use Registry.Add and Resolver.Resolve.
  - Move command (mas move) improvements per CR comments:
    - Base-name fallback when arc.Name is empty.
    - Directory destination handling (append name if --to is a directory).
    - Overwrite prevention and EXDEV copy+remove fallback.
    - Robust managed-path detection via filepath.Rel.
  - PR mechanics and docs:
    - Opened PR #9 for the list refactor and tests.
    - Reviewed PR #5 CodeRabbit comments; implemented move fixes and posted status.
    - Investigated PR #8 CI failure; captured root causes from logs.
    - Resolved markdown merge conflicts in 7EP docs (index and 7ep-0004) and preserved the latest status/columns.

- How it went
  - Local builds/tests: green after iterative fixes.
  - CodeRabbit on PR #5: SUCCESS with actionable feedback that improved safety and portability of move.
  - CI for PR #8: still failing upstream; cause is clear and fix is straightforward (align resolver/test helpers to central API and errors).

- Challenges encountered
  - Initial compile errors due to duplicate error type declarations (resolver vs errors) and outdated helper names in tests; addressed locally by unifying against errors.go and updating helper calls.
  - A few awkward intermediate diffs while refactoring list.go (formatting/placement); cleaned up to ensure no behavior changes.
  - Git conflicts in docs during rebase; resolved by keeping expanded columns and current status blocks.

- Concerns / risks
  - PR #8 CI remains red until we push the resolver/helper alignment; low effort but needs a targeted patch.
  - Drift on mas_move.go between CR-refined version and latest edits; we should reconcile intent before finalizing PR #5.
  - Small tests added for list refactor were later removed; decide whether to keep lightweight unit coverage for helpers.

- Learnings / patterns
  - Centralizing error types in internal/storage/errors.go avoids duplication and CI drift.
  - For CLI filters, a tiny applyFilters helper improves readability and testability without changing semantics.
  - Move semantics benefit from explicit safeguards (dir vs file dest, overwrite checks, EXDEV fallback) and a portable managed-path check.

- Recommendations for CC
  - PR #5: Review the move command changes for safety/UX; confirm preferred final variant (user’s latest vs CR-refined), then merge.
  - PR #8: Allow AC to push the resolver/helper alignment patch and re-run CI; should go green quickly.
  - Docs: The 7EP index with Area/Updated columns reads well—keep as canonical.
  - Tests: Consider keeping small unit tests for helper functions (filters/duration) to prevent regressions.

- Next steps (short)
  - Finalize mas_move.go reconciliation and push to PR #5.
  - Patch PR #8 to align resolver/test helpers, re-run CI, and confirm green.
  - Proceed with PR #9 review; minimal risk since no behavior changes were introduced.

### Active PR Status & Action Items for AC

See section 302-426 above for current PR status and action items.

---

**Note**: This file serves as shared context for all AI agents. For private development notes and half-formed ideas, use PRIVATE.md (git-ignored).
