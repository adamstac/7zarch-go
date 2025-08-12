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
- **7EP-0001**: Trash Management System (AC assigned, difficulty 3) - 🔷 Ready to implement
- **7EP-0002**: CI Integration & Automation (CC assigned, difficulty 2) - 🟡 Draft
- **7EP-0003**: Database Migrations & Schema Management (AC assigned, difficulty 2) - 🟡 Draft
- **7EP-0004**: MAS Foundation Implementation (AC assigned, difficulty 4) - ✅ Completed
- **7EP-0005**: Comprehensive Test Dataset System (CC assigned, difficulty 3) - 🟡 Draft
- **7EP-0006**: Minimal Performance Testing (CC assigned, difficulty 1) - ✅ Completed
- **7EP-0007**: Enhanced MAS Operations (AC/CC assigned, difficulty 3) - 🟡 Draft

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
- ✅ **MAS Foundation (7EP-0004)**: ULID resolution, show command, enhanced list
- ✅ **Performance Validation (7EP-0006)**: All requirements exceeded by 100-2,941x margins
- ✅ Enhanced list display with ULID support and status indicators
- ✅ Delete/move commands with soft delete and trash management
- ✅ Auto-purge countdown with configurable retention days
- ✅ 7EP system for formal enhancement proposals
- ✅ Documentation workflow process

### Active Development  
- 🔷 **Trash Management Commands** (AC): restore, trash list, trash purge (7EP-0001) - Ready to start
- 🟡 **Enhanced MAS Operations** (AC/CC): move, import, batch operations (7EP-0007) - Draft
- 🔄 **Database Migrations** (AC): Safe schema evolution system (7EP-0003) - Draft
- 🔄 **CI Integration** (CC): GitHub Actions workflows (7EP-0002) - Draft
- 🟡 **Comprehensive Test System** (CC): Edge cases and integration scenarios (7EP-0005) - Draft

## Development Roadmap

### **Current Sprint: Trash Management & Enhanced Operations**
**7EP-0001: Trash Management System** (AC assigned, difficulty 3)
- Status: 🔷 Ready for implementation - all dependencies complete
- Primary: restore, trash list, trash purge commands
- Foundation: 7EP-0004 MAS patterns proven and ready to use
- **Estimated time: 9-13 hours total**

**7EP-0007: Enhanced MAS Operations** (AC/CC assigned, difficulty 3)
- Status: 🟡 Draft - comprehensive advanced operations
- Primary: move/relocate, import workflows, batch operations
- Dependencies: 7EP-0004 complete, builds on proven patterns
- **Complements trash management for complete lifecycle**

### **Next Sprint: Development Infrastructure**
**7EP-0002: CI Integration & Automation** (CC assigned, difficulty 2)
- Status: 🟡 Draft - GitHub Actions workflows needed
- Dependencies: None (parallel work)
- Focus: Automated testing, quality gates, release automation

**7EP-0005: Comprehensive Test Dataset System** (CC assigned, difficulty 3)
- Status: 🟡 Draft - builds on 7EP-0006 learnings
- Dependencies: 7EP-0006 patterns proven successful
- Focus: Edge case coverage, integration scenarios

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

#### **Milestone 1: MAS Core Complete** ✅ (Completed 2025-08-12)
- ✅ ULID resolution working
- ✅ Show command operational  
- ✅ Enhanced list with filters
- ✅ Performance validated (100-2,941x faster than requirements)
- **Success Criteria**: Users can manage archives by ID without paths

#### **Milestone 2: Production Ready** (Target: 2-3 weeks from now)
- 🔄 Database migrations automated
- 🔷 Trash management complete (ready to start)
- 🔄 CI pipeline operational
- **Success Criteria**: Safe for daily use, automated testing

#### **Milestone 3: Advanced Workflows** (Target: 4-6 weeks from now)  
- 🔄 Enhanced MAS operations (move, import, batch)
- 🔄 Comprehensive test dataset system
- 🔄 Performance optimization for large collections
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
- **7EP-0008**: Multi-Location Storage Support (TrueNAS, cloud backends)
- **7EP-0009**: Archive Deduplication and Content Management  
- **7EP-0010**: Performance Optimization for Large Collections
- **7EP-0011**: Configuration Management and User Profiles

### CodeRabbit Integration Improvements (AC Focus)
- Enhanced PR review context
- Automated code pattern enforcement
- Integration with 7EP compliance checking
- Quality gate automation

---

**Note**: This file serves as shared context for all AI agents. For private development notes and half-formed ideas, use PRIVATE.md (git-ignored).