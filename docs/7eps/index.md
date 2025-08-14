# 7zarch Enhancement Proposals (7EPs)

This directory contains **7zarch Enhancement Proposals (7EPs)** - structured documents that describe new features, significant changes, or improvements to the 7zarch-go project.

## What is a 7EP?

A 7EP is a design document that provides:
- **Clear problem definition** with evidence and reasoning
- **Technical design** and implementation approach  
- **Use cases** and acceptance criteria
- **Assignment** to development teams (CC/AC)
- **Difficulty assessment** for resource planning

## 7EP Process

1. **Draft**: Create 7EP using the [template](template.md)
2. **Discussion**: Open GitHub issue for community feedback  
3. **Review**: Formal review by maintainers
4. **Decision**: Accept/Reject with documented rationale
5. **Implementation**: Assigned developer(s) implement
6. **Completion**: Mark as Implemented when merged

## Status Legend

- 🟡 **Draft** - Initial proposal, seeking feedback
- 🔵 **Under Review** - Formal review in progress  
- 🟢 **Accepted** - Approved for implementation
- 🔷 **Ready** - Ready to start implementation (dependencies complete)
- 🔄 **In Progress** - Implementation underway
- 🔴 **Rejected** - Not moving forward (with rationale)
- ✅ **Implemented** - Completed and merged

## Current 7EPs

| 7EP | Title | Status | Assignee | Difficulty | Created |
|-----|-------|--------|----------|------------|---------|
| [0001](7ep-0001-trash-management.md) | Trash Management System | ✅ **Implemented** | AC | 3 | 2025-08-11 |
| [0002](7ep-0002-ci-integration.md) | CI Integration & Automation | ✅ **Implemented** | CC | 2 | 2025-08-12 |
| [0003](7ep-0003-database-migrations.md) | Database Migrations & Schema Management | ✅ **Implemented** | AC | 2 | 2025-08-12 |
| [0004](7ep-0004-mas-foundation.md) | MAS Foundation Implementation | ✅ **Implemented** | AC | 4 | 2025-08-12 |
| [0005](7ep-0005-test-dataset-system.md) | Comprehensive Test Dataset System | 🟡 Draft | CC | 3 | 2025-08-12 |
| [0006](7ep-0006-minimal-performance-testing.md) | Minimal Performance Testing for 7EP-0004 | ✅ **Implemented** | CC | 1 | 2025-08-12 |
| [0007](7ep-0007-enhanced-mas-operations.md) | Enhanced MAS Operations | 🔄 **Phase 3 Complete** | CC | 3 | 2025-08-12 |
| [0010](7ep-0010-interactive-tui-application.md) | Interactive TUI Application | ✅ **Implemented** | Amp/CC | 3 | 2025-08-12 |
| [0011](7ep-0011-re-tighten-golangci-lint.md) | Code Quality Linting Strategy | ✅ **Resolved** | CC | 2 | 2025-08-13 |
| [0013](7ep-0013-robust-build-pipeline.md) | Robust Build Pipeline | ✅ **Implemented** | CC | 3 | 2025-08-13 |
| [0014](7ep-0014-critical-foundation-gaps.md) | Critical Foundation Gaps | ✅ **Implemented** | Amp | 4 | 2025-08-13 |
| [0015](7ep-0015-code-quality-foundation.md) | Code Quality Foundation | ✅ **Implemented** | CC | 3 | 2025-08-13 |
| [0016](7ep-0016-tui-first-interface-evolution.md) | TUI-First Interface Evolution | 🟡 **Draft** | Future | 4 | 2025-08-13 |

## How to Contribute

### Proposing a New 7EP

1. **Check existing 7EPs** to avoid duplication
2. **Copy the [template](template.md)** to `7ep-XXXX-title.md`
3. **Fill out all sections** thoroughly
4. **Open a GitHub issue** linking to your 7EP for discussion
5. **Iterate based on feedback** until ready for formal review

### Implementation Guidelines

- **Follow the technical design** outlined in the accepted 7EP
- **Update the 7EP status** as implementation progresses  
- **Add implementation notes** to the 7EP when complete
- **Reference the 7EP** in commit messages and PRs

## Templates

- **[7EP Template](template.md)** - Standard format for new proposals

## Archive

Completed 7EPs remain in this index for historical reference and to inform future enhancements.

---

💡 **Questions?** Open an issue or discussion for guidance on the 7EP process.