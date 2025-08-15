# 7EP-0017: Document Driven Development Framework

**Status:** Active  
**Author(s):** Amp (Sourcegraph), Adam Stacoviak  
**Assignment:** Framework Design  
**Difficulty:** 3 (moderate - process standardization with systematic coordination)  
**Created:** 2025-08-13  
**Updated:** 2025-08-13  

## Executive Summary

Establish a systematic Document Driven Development (DDD) framework that separates "what to build" (7EPs), "who builds it" (team identity docs), and "current priorities" (team assignment docs + shared NEXT.md), creating clear coordination patterns that scale from individual tasks to cross-team sprints.

## Evidence & Reasoning

**Current Documentation Success:**
- ✅ **7EP system** working excellently for feature specification and technical design
- ✅ **Team identity docs** (AMP.md, CLAUDE.md, AUGMENT.md, ADAM.md) provide role clarity
- ✅ **Recent coordination success** - 7EP-0014 foundation sprint (2 days vs 4-6 target)
- ✅ **Cross-team collaboration** - Amp TUI + CC backend integration working well

**Coordination Gaps Identified:**
- **Priority confusion** - No clear "what's next for everyone" visibility
- **Assignment ambiguity** - 7EPs have assignments but no current task tracking
- **Handoff friction** - Success patterns exist but not systematized
- **Cross-dependencies** - Limited visibility into blocking relationships

**Strategic Need:**
- **Scale coordination** - Framework that works for 2-10 team members
- **Maintain momentum** - Clear priorities prevent idle time and duplicate work
- **Document everything** - All decisions and priorities captured in git-tracked docs
- **Enable autonomy** - Team members can see their priorities and dependencies clearly

## Operational Activation Protocol

**Strategic Recognition**: Framework effectiveness requires **immediate operational activation** when coordination friction is identified.

### **Activation Triggers**
- **Coordination Friction**: Team members unclear on current priorities or next tasks
- **Cross-Team Dependencies**: Blocking relationships visible but not systematically tracked
- **Handoff Inefficiency**: Success patterns exist but require manual discovery each time
- **Strategic Misalignment**: Implementation work disconnected from strategic priorities

### **Same-Day Activation Process**
1. **Immediate Document Creation**
   - Create operational docs using framework templates (docs/development/)
   - Populate with current real project state (not example content)
   - Establish NEXT.md as primary coordination hub

2. **Boot-up Integration** 
   - Update existing team startup documents (/CLAUDE.md, etc.)
   - Add operational document sequence to team onboarding
   - Document operational framework usage patterns

3. **Live Validation**
   - Apply framework to current coordination challenges (same session)
   - Use real project handoffs to test coordination patterns
   - Refine templates based on immediate operational feedback

### **Operational Success Metrics**
- **<30 Second Priority Discovery**: Team members find current work priorities quickly
- **Visible Cross-Dependencies**: Blocking relationships clear in NEXT.md
- **Reduced Coordination Overhead**: Less time spent clarifying "what should I work on next"
- **Preserved Team Momentum**: Framework enhances existing patterns without disruption

## Use Cases

### Primary Use Case: Clear Priority Coordination
```bash
# Team member wants to know what to work on
1. Check personal assignment doc: docs/development/roles/CLAUDE.md
2. Check shared priorities: docs/development/NEXT.md  
3. Check active 7EP coordination: docs/7eps/7ep-XXXX.md

# Result: Clear understanding of current work + next priorities + coordination needs
```

## Boot-up Integration Specification

**Strategic Integration**: Framework must connect seamlessly with existing team startup sequences to ensure adoption and operational effectiveness.

### **Primary Boot-up Path Enhancement**

#### **Target Integration Point**: /CLAUDE.md Quick Start Checklist
**Current Boot-up Sequence (Enhanced)**:
```bash
1. Check git status
   git status && git pull && git branch

2. Review current state  
   cat docs/development/pr-merge-roadmap.md | head -50
   gh pr list && git log --oneline -10

2.5. **CHECK OPERATIONAL PRIORITIES** (NEW)
     # Personal current assignments and coordination
     cat docs/development/roles/CLAUDE.md | head -20
     
     # Shared team priorities and blockers  
     cat docs/development/NEXT.md | head -30
     
     # Active 7EP coordination context
     grep -l "Status.*ACTIVE\|In Progress" docs/7eps/*.md | xargs ls -la

3. Understand today's priorities
   cat docs/development/tomorrow-plan.md  # (deprecated → use NEXT.md)

4. Test the build
   go build -o 7zarch-go . && ./7zarch-go list --dashboard
```

### **Day 1 Team Onboarding Sequence**
**For New Team Members or Role Transitions**:
1. **Identity Understanding**
   - Read `/CLAUDE.md` (role definition and capabilities)
   - Read `/AMP.md` (leadership roles and activation)
   - Understand team member capabilities and authorities

2. **Current Work Context**
   - Read `docs/development/roles/CLAUDE.md` (current assignments)  
   - Read `docs/development/NEXT.md` (shared coordination hub)
   - Review active 7EPs for coordination requirements

3. **Operational Integration**
   - Execute enhanced boot-up checklist (above)
   - Identify immediate coordination needs or blockers
   - Update personal assignment doc with availability/capacity

### **Framework Document Reading Sequence**
**Optimized Information Discovery**:
```markdown
## Operational Priority (Daily)
1. docs/development/NEXT.md       # What's happening now across all teams
2. docs/development/roles/CLAUDE.md     # My current work and next priorities  
3. Active 7EPs                    # Sprint-level coordination context

## Reference Information (As Needed)
4. docs/7eps/index.md            # Long-term feature planning
5. /CLAUDE.md                    # Role definition and context
6. docs/development/README.md     # Framework usage patterns
7. docs/development/actions/       # Standardized workflow processes (COMMIT.md, MERGE.md, NEW-FEATURE.md)
```

### **Integration Success Metrics**
- **Boot-up Time Reduction**: From "unclear priorities" to "ready to work" in <3 minutes
- **Coordination Clarity**: Cross-team dependencies visible within boot-up sequence
- **Framework Adoption**: Natural integration enhances (doesn't replace) existing successful patterns
- **Context Preservation**: All coordination context accessible through familiar startup process

### Secondary Use Cases

#### Cross-Team Sprint Coordination
```markdown
# Example: 7EP-0014 Foundation Sprint pattern
docs/7eps/7ep-0014.md:
- Phase assignments with clear deliverables
- Success metrics and timeline
- Coordination points between team members

docs/development/NEXT.md:
- Sprint status: "Phase 2 active, Phase 3 waiting"
- Blockers: "PR #11 CI fixes needed before Phase 2 completion"
- Next actions: "CC: Fix CI, Amp: Review Phase 2 when complete"
```

#### Individual Task Management  
```markdown
# Example: CC current work clarity
docs/development/roles/CLAUDE.md:
## Current Assignments (Updated: 2025-08-13)
- **7EP-0007 Phase 3** - ACTIVE (implementing batch operations, ETA: 3-4 days)
- **Performance review** - BLOCKED (waiting for Amp architectural feedback)
- **Code review queue** - READY (2 PRs need review)

## Next Priorities
1. Complete 7EP-0007 Phase 3 batch operations
2. Hand off to Amp for architectural review  
3. Begin TrueNAS backend integration (7EP-0018)
```

## Technical Design

### Document Architecture

#### **Layer 1: Identity & Roles** (Root Level)
```
/AMP.md           # Amp capabilities and activation protocols
/CLAUDE.md        # CC role definition and context  
/AUGMENT.md       # AC role definition and context
/ADAM.md          # Project owner perspective and priorities
```

**Purpose:** Who are the team members, what are their capabilities, how to activate specific roles

#### **Layer 2: Specifications** (docs/7eps/)
```
/docs/7eps/7ep-XXXX.md    # Feature specifications (WHAT to build)
/docs/7eps/index.md       # Status tracking and overview
/docs/7eps/template.md    # Standardized 7EP format
```

**Purpose:** Long-term planning, technical specifications, no current assignments

#### **Layer 3: Current Work** (docs/development/)
```
/docs/development/roles/AMP.md      # Amp current assignments and priorities
/docs/development/roles/CLAUDE.md   # CC current work and next tasks
/docs/development/roles/AUGMENT.md  # AC current assignments and coordination
/docs/development/roles/ADAM.md     # Project owner current focus and decisions
/docs/development/NEXT.md     # Shared coordination hub for everyone
/docs/development/actions/    # Workflow processes (COMMIT.md, MERGE.md, NEW-FEATURE.md)
```

**Purpose:** Current task tracking, priority coordination, immediate planning

## Documentation Hierarchy Clarification

**Strategic Hierarchy**: Clear operational documentation relationship for effective team coordination and priority management.

### **Primary Operational Hub**: docs/development/NEXT.md

#### **Hub Designation Rationale**
- **Single Source of Truth** for "what's happening now" across all team members
- **Cross-Team Visibility** - All coordination points visible in one document
- **Real-Time Updates** - Updated when blockers occur or completions happen
- **Strategic Alignment** - Connects individual work to project priorities

#### **Operational Document Relationship**
```markdown
## Primary Coordination Flow
docs/development/NEXT.md              # Team coordination hub
├── docs/development/roles/CLAUDE.md        # CC current assignments (feeds into NEXT.md)
├── docs/development/roles/AUGMENT.md       # AC current assignments (feeds into NEXT.md)  
├── docs/development/roles/AMP.md           # Amp current assignments (feeds into NEXT.md)
├── docs/development/roles/ADAM.md          # Adam priorities (drives NEXT.md priorities)
└── docs/development/actions/               # Standardized workflow processes

## Supporting Reference
docs/7eps/index.md                    # Long-term planning reference
docs/7eps/7ep-XXXX.md                 # Active sprint coordination (linked from NEXT.md)
```

### **Documentation Type Classifications**

#### **Operational Documents** (Updated Frequently)
- **docs/development/NEXT.md** - Updated when: blockers occur, phases complete, priorities shift
- **docs/development/[TEAM].md** - Updated when: status changes, new assignments, handoffs occur
- **Active 7EPs** - Updated during: implementation phases, architectural changes

#### **Reference Documents** (Updated Infrequently)  
- **Root identity docs** (/AMP.md, /CLAUDE.md) - Role definitions and capabilities
- **Completed 7EPs** - Historical specification reference
- **Process documentation** - Framework standards and patterns

#### **Strategic Documents** (Updated Based on Decisions)
- **docs/development/roles/ADAM.md** - Adam's current strategic priorities and decisions needed
- **Active 7EP coordination sections** - Sprint-level strategic direction

### **Update Responsibility Matrix**

| Document | Primary Owner | Update Trigger | Coordination Responsibility |
|----------|---------------|----------------|----------------------------|
| **NEXT.md** | Shared (whoever encounters blocker/completion) | Real-time | Cross-team visibility |
| **CLAUDE.md** | CC | Status changes | Personal coordination with team |
| **AUGMENT.md** | AC | Status changes | Personal coordination with team |
| **AMP.md** | Amp | Status changes | Leadership coordination |
| **ADAM.md** | Adam | Strategic decisions | Strategic direction setting |
| **Active 7EPs** | Assigned implementer | Implementation progress | Sprint coordination |

### **Coordination Effectiveness Patterns**

#### **Daily Operations**
1. **Team members** update personal assignment docs when status changes
2. **NEXT.md** updated by whoever encounters blockers or completes major phases  
3. **Active 7EPs** updated during implementation with technical progress
4. **Cross-team** handoffs documented in NEXT.md coordination sections

#### **Strategic Coordination**
1. **Adam** updates strategic priorities in ADAM.md and NEXT.md
2. **Amp strategic/technical** reviews reflected in NEXT.md and AMP.md
3. **Implementation teams** align personal assignment docs with strategic direction
4. **Project phases** coordinated through NEXT.md priority sequencing

### **Operational Hierarchy Success Metrics**
- **Information Discovery Speed**: Team members find current priorities within NEXT.md quickly
- **Cross-Team Visibility**: Dependencies and blockers visible across team boundaries
- **Strategic Alignment**: Individual work clearly connected to project strategic direction  
- **Coordination Overhead Reduction**: Less time spent in meetings clarifying "what's next"

#### **Layer 4: Process Documentation** (docs/development/)
```
/docs/development/handoff-protocols.md
/docs/development/coordination-patterns.md
/docs/development/documentation-standards.md
```

**Purpose:** How coordination works, standardized patterns, process guidance

### NEXT.md Coordination Hub

#### **Structure Template:**
```markdown
# What's Next for Everyone

**Last Updated:** 2025-08-13 14:30  
**Project Phase:** Advanced Features Development

## 🔄 Current Active Work
**7EP-0007 Phase 3:** CC implementing batch operations (ETA: 3-4 days)
**TUI Polish:** Amp reviewing feedback and planning command line integration

## 📋 Next Priorities (In Priority Order)
1. **CC:** Complete 7EP-0007 Phase 3 → Hand off for Amp review
2. **Amp:** Architectural review of batch operations → Provide optimization feedback  
3. **CC:** Begin TrueNAS backend integration (7EP-0018) when 7EP-0007 complete
4. **Amp:** Implement 7EP-0016 TUI command line when backend features ready

## 🔗 Coordination Points
- **CC → Amp:** Batch operations architecture review needed (Phase 3 completion)
- **Amp → CC:** Remote storage patterns for TUI integration feedback
- **Adam decision needed:** TrueNAS backend vs advanced TUI features priority

## 🚫 Blocked/Waiting
- **7EP-0016 implementation:** Waiting for 7EP-0007 completion
- **TrueNAS integration:** Waiting for Adam priority decision
- **Performance optimization:** Ready but lower priority than feature completion

## 🎯 Success Metrics This Week
- [ ] 7EP-0007 Phase 3 complete with batch operations working
- [ ] TUI feedback incorporated and command line architecture planned
- [ ] Next sprint priorities clarified based on Adam strategic direction
```

### Team Assignment Doc Template

#### **Structure for docs/development/[TEAM].md:**
```markdown
# [TEAM NAME] Current Assignments

**Last Updated:** 2025-08-13  
**Status:** Active  
**Current Focus:** [Primary work stream]

## 🎯 Current Assignments
### Active Work (This Week)
- **[7EP-XXXX Phase Y]** - ACTIVE (brief description, ETA)
- **[Task]** - BLOCKED (reason, waiting for what)
- **[Review]** - READY (ready for action)

### Next Priorities (Priority Order)
1. **[Next task]** - [Dependencies and context]
2. **[Following task]** - [When it becomes available]
3. **[Future task]** - [Longer-term planning]

## 🔗 Coordination Needed
- **Handoff to [TEAM]:** [What needs handoff and when]
- **Waiting from [TEAM]:** [What blocking dependencies]
- **Adam decision:** [Decisions needed from project owner]

## ✅ Recently Completed
- **[Completed work]** - [Impact and handoff status]
- **[Finished task]** - [Results and next steps]

## 📝 Implementation Notes
[Context, learnings, and technical notes relevant for coordination]
```

### 7EP Assignment Integration

#### **Enhanced 7EP Assignment Tracking:**
```markdown
## Implementation Coordination (When Active)

### Current Phase: Phase 3 - Batch Operations
**Assigned:** CC (Primary implementation)  
**Support:** Amp (Architectural review)  
**Timeline:** 3-4 days

#### CC Tasks (Active)
- [ ] Multi-archive operation framework - IN PROGRESS
- [ ] Progress tracking implementation - NEXT  
- [ ] CLI integration - HANDOFF-READY (for Amp review)

#### Amp Tasks (Review)  
- [ ] Architecture review - WAITING (for CC Phase 3 completion)
- [ ] Performance validation - READY (when implementation complete)
- [ ] Integration guidance - ON-DEMAND (as needed)

#### Adam Decisions Needed
- [ ] TrueNAS backend priority vs advanced TUI features
- [ ] Performance targets for batch operations
```

## Implementation Plan

### Phase 1: Framework Design & Templates
- [ ] **Document Structure Definition**
  - [ ] Finalize Layer 3 (current work) structure
  - [ ] Create standardized templates for team assignment docs
  - [ ] Design NEXT.md coordination patterns
  - [ ] Document handoff status indicators

- [ ] **Integration with Existing 7EP Process**
  - [ ] Enhance 7EP template with coordination sections
  - [ ] Define assignment patterns for active 7EPs
  - [ ] Establish handoff protocols and status tracking

### Phase 2: Template Creation & Examples  
- [ ] **Team Assignment Templates**
  - [ ] Create docs/development/roles/AMP.md template and example
  - [ ] Create docs/development/roles/CLAUDE.md template with current work
  - [ ] Create docs/development/roles/AUGMENT.md template
  - [ ] Create docs/development/roles/ADAM.md template for project owner priorities
  - [ ] Create docs/development/actions/COMMIT.md workflow
  - [ ] Create docs/development/actions/MERGE.md workflow
  - [ ] Create docs/development/actions/NEW-FEATURE.md workflow

- [ ] **NEXT.md Coordination Hub**
  - [ ] Design shared priority coordination structure
  - [ ] Create example with current project state
  - [ ] Define update cadence and ownership
  - [ ] Document coordination patterns

### Phase 3: Process Documentation & Standards
- [ ] **Coordination Standards**
  - [ ] Document handoff status indicators (ACTIVE, BLOCKED, READY, HANDOFF-READY)
  - [ ] Create coordination pattern examples
  - [ ] Establish update responsibilities and cadence
  - [ ] Define escalation paths for priority conflicts

- [ ] **Integration Guidelines**
  - [ ] How to transition 7EPs to active implementation with assignments
  - [ ] When and how to use NEXT.md for coordination
  - [ ] Standards for team assignment doc maintenance
  - [ ] Cross-team communication protocols

### Phase 4: Documentation Cleanup & Migration
- [ ] **Audit Existing docs/development/ Files**
  - [ ] Review tomorrow-plan.md → Migrate relevant content to NEXT.md or archive
  - [ ] Review pr-merge-roadmap.md → Update for new coordination patterns or archive
  - [ ] Review 7ep-0010-*.md → Migrate to reference docs or archive as historical
  - [ ] Review sprint-planning-analysis.md → Integrate with new priority framework

- [ ] **Content Migration Strategy**
  - [ ] Identify content to migrate to new team assignment docs
  - [ ] Archive historical content that's no longer relevant
  - [ ] Update remaining docs to work with new framework
  - [ ] Remove duplicate or outdated coordination files

- [ ] **Final Documentation Structure**
  - [ ] Clean docs/development/ with only current framework files
  - [ ] Ensure all coordination flows through new system
  - [ ] Validate no critical information lost in migration
  - [ ] Document what was archived and why

## Framework Benefits

### **Individual Clarity**
- **Clear current work** - Each team member knows their priorities
- **Visible dependencies** - Understand what they're waiting for
- **Next task visibility** - No idle time, clear progression path
- **Context preservation** - All work context documented and preserved

### **Team Coordination**
- **Shared visibility** - NEXT.md shows everyone's priorities and blockers
- **Dependency management** - Clear handoff points and coordination needs
- **Priority alignment** - Adam's strategic priorities visible to all team members
- **Conflict resolution** - Systematic approach to priority and resource conflicts

### **Project Management**
- **Systematic planning** - All work flows through documented frameworks
- **Progress tracking** - Multiple levels of progress visibility
- **Quality assurance** - Documented handoff and review patterns
- **Knowledge preservation** - All decisions and context captured in git

### **Scalability**
- **Add team members** - Framework patterns work for expanding teams
- **Complex projects** - Coordination patterns scale to larger initiatives
- **Long-term planning** - Strategic vision connects to tactical execution
- **Process improvement** - Framework itself can evolve based on usage patterns

## Testing Strategy

### Acceptance Criteria
- [ ] Team members can quickly understand their current priorities
- [ ] Cross-team dependencies are clearly visible and trackable
- [ ] NEXT.md provides accurate shared coordination view
- [ ] Handoff protocols reduce coordination friction
- [ ] Adam can see and adjust strategic priorities effectively
- [ ] Framework scales from individual tasks to cross-team sprints

### Validation Approach
- **Real-world testing** - Use framework for next 7EP implementation cycle
- **Team feedback** - Gather input from CC, AC, and Adam on coordination effectiveness
- **Iteration tracking** - Document what works well and what needs adjustment
- **Process optimization** - Refine templates and patterns based on usage

### Success Metrics
- **Coordination time reduction** - Less time spent clarifying priorities and dependencies
- **Implementation velocity** - Faster progression from planning to execution
- **Quality maintenance** - Systematic coordination without quality degradation
- **Team satisfaction** - Clear priorities and reduced coordination friction

## Migration/Compatibility

### Breaking Changes
None - framework enhances existing documentation, doesn't replace it.

### Upgrade Path
- Create team assignment docs in docs/development/ for existing team members
- Populate NEXT.md with current project priorities
- Enhance active 7EPs with coordination sections
- Gradually adopt framework patterns for new work

### Backward Compatibility
- All existing 7EPs continue working unchanged
- Team identity docs (AMP.md, CLAUDE.md, etc.) remain as role definitions
- Existing coordination patterns continue while new framework adds structure

## Implementation Examples

### Example: docs/development/roles/AMP.md
```markdown
# Amp Current Assignments

**Last Updated:** 2025-08-13 15:30  
**Status:** Active  
**Current Focus:** TUI architecture and 7EP-0007 coordination

## 🎯 Current Assignments
### Active Work (This Week)
- **TUI Polish & Feedback** - ACTIVE (incorporating user feedback, planning command line)
- **7EP-0007 Coordination** - READY (architectural oversight for CC's batch operations)
- **7EP-0016 Planning** - ONGOING (TUI-first interface evolution design)

### Next Priorities
1. **7EP-0007 Phase 3 Review** - Provide architectural feedback when CC completes batch operations
2. **7EP-0016 Implementation** - Begin TUI command line when 7EP-0007 backend complete
3. **Performance Optimization** - Review and optimize completed features

## 🔗 Coordination Needed
- **Handoff from CC:** 7EP-0007 Phase 3 completion for architectural review
- **Coordination with Adam:** Strategic priority between TrueNAS backend vs TUI evolution
- **Support to CC:** Ongoing architectural guidance as needed

## ✅ Recently Completed
- **7EP-0010 TUI Implementation** - Simple themed interface with 9 color schemes
- **7EP-0014 Foundation Sprint** - 2-day coordination of critical foundation gaps
- **Documentation Framework** - Current DDD framework design and proposal
```

### Example: docs/development/NEXT.md
```markdown
# What's Next for Everyone

**Last Updated:** 2025-08-13 15:30  
**Project Phase:** Advanced Features Development  
**Sprint Status:** 7EP-0007 Phase 3 active, TUI foundation complete

## 🔄 Current Active Work
**CC:** 7EP-0007 Phase 3 batch operations (ETA: 3-4 days)  
**Amp:** TUI feedback integration and 7EP-0016 planning  
**Adam:** Strategic priorities and TrueNAS backend decision  

## 📋 Next Priorities (Sequential)
1. **CC completes 7EP-0007 Phase 3** → Batch operations working
2. **Amp reviews architecture** → Provide optimization feedback
3. **CC hands off to Amp** → Begin 7EP-0016 or TrueNAS backend
4. **Adam strategic decision** → TrueNAS integration vs TUI evolution priority

## 🔗 Active Coordination Points
- **CC → Amp:** Phase 3 completion handoff for architectural review
- **Amp → CC:** TUI command line integration patterns (when ready)
- **Adam → Team:** Strategic priority decision (backend vs frontend features)

## 🚫 Blocked/Waiting
- **7EP-0016 implementation:** Waiting for 7EP-0007 completion + Adam priority decision
- **TrueNAS integration:** Waiting for Adam backend priority confirmation
- **Advanced TUI features:** Dependent on backend feature completion

## 🎯 Success Metrics This Week
- [ ] 7EP-0007 Phase 3 complete with working batch operations
- [ ] Architectural review completed with optimization feedback
- [ ] Next sprint priorities clarified and assigned
- [ ] Framework effectiveness validated through real usage
```

## Framework Standards

### Status Indicators
**For Team Assignment Docs:**
- **ACTIVE** - Currently working on this task
- **NEXT** - Immediate next priority when current work complete
- **BLOCKED** - Cannot proceed due to dependency
- **READY** - Ready to begin when prioritized
- **HANDOFF-READY** - Complete and ready to hand off to another team member
- **WAITING** - Waiting for external input or decision

### Update Responsibilities  
- **Team assignment docs** - Updated by team member when status changes
- **NEXT.md** - Updated by whoever makes significant progress or encounters blockers
- **Active 7EPs** - Updated by assigned implementer during active work
- **Framework itself** - Updated based on usage patterns and team feedback

### Coordination Patterns
- **Daily:** Team members update their assignment docs when status changes
- **Milestone:** NEXT.md updated when major phases complete or blockers encountered
- **Strategic:** Adam reviews and adjusts priorities in NEXT.md as needed
- **Cross-team:** Coordination sections in active 7EPs used for sprint-level coordination

## Migration/Compatibility

### Breaking Changes
None - framework enhances existing coordination without breaking current patterns.

### Upgrade Path
1. **Create team assignment docs** for current team members
2. **Populate NEXT.md** with current project state
3. **Enhance active 7EPs** with coordination sections
4. **Clean up docs/development/** - Migrate, archive, or update existing files
5. **Begin using patterns** for new work while existing work continues

### Documentation Cleanup Analysis

#### **Files to Review in docs/development/:**
- **tomorrow-plan.md** - Daily planning → Migrate to NEXT.md or archive
- **pr-merge-roadmap.md** - Outdated PR tracking → Update with current coordination or archive
- **7ep-0010-tui-implementation-guide.md** - Implementation guide → Move to reference docs or archive
- **7ep-0010-quick-reference.md** - Quick reference → Move to reference docs or archive  
- **sprint-planning-analysis.md** - Strategic planning → Integrate with NEXT.md or keep as reference
- **migration-best-practices.md** - Keep (technical reference)
- **emoji-usage-guidelines.md** - Keep (technical standards)

#### **Migration Strategy:**
1. **Content audit** - Identify what's still relevant vs historical
2. **Reference migration** - Move technical guides to docs/reference/
3. **Priority integration** - Migrate planning content to NEXT.md and team assignment docs
4. **Historical archiving** - Move outdated content to docs/archive/ 
5. **Final structure** - docs/development/ contains only active framework files

### Backward Compatibility
- Existing 7EP process unchanged
- Team identity docs remain as role definitions  
- Current coordination patterns continue while framework adds structure
- All content preserved through migration process

## Future Considerations

### Advanced Coordination Features
- **Dependency tracking** - Automated detection of blocking relationships
- **Timeline management** - Integration with project timeline planning
- **Metrics collection** - Track coordination effectiveness and team velocity
- **External integration** - Hooks for project management tools if needed

### Process Evolution
- **Framework refinement** - Continuous improvement based on usage patterns
- **Template standardization** - Standardized templates based on successful patterns
- **Automation opportunities** - Potential for automated status tracking and updates
- **Scale preparation** - Framework patterns that work for larger teams

## AI Agent Technical Competency Framework

### **Document Architecture for AI Agents**

The DDD framework integrates a systematic approach to AI agent technical competency and coordination:

**Technical Competency Layer**: `/AGENT.md` (universal for all AI agents)
- Build/test/lint commands and development workflow
- Architecture overview and package structure
- Code style conventions and patterns specific to the project
- Key technical concepts unique to the codebase

**Role Coordination Layer**: `/docs/development/[ROLE].md` (role-specific)
- Current assignments and priorities
- Team coordination context and communication patterns  
- Role-specific responsibilities and focus areas
- Cross-team dependencies and blocking relationships

**Operational Procedure Layer**: `/BOOTUP.md` (session startup)
- Standardized startup sequence for quick context loading
- References to both technical and role-specific documentation
- Current project state and team status
- Immediate priorities and coordination needs

### **Root File Strategy**

**Pointer Files in Root** (for agent accessibility):
- `/CLAUDE.md` → References both AGENT.md and docs/development/roles/CLAUDE.md
- `/AMP.md` → References both AGENT.md and docs/development/roles/AMP.md  
- `/AUGMENT.md` → References both AGENT.md and docs/development/roles/AUGMENT.md

**Single Source of Truth** (eliminates duplication):
- Technical patterns: AGENT.md only
- Role coordination: docs/development/[ROLE].md only
- Root files become simple navigation aids

### **Integration with DDD Framework**

This technical competency layer complements existing DDD components:
1. **Strategic Planning** (7EPs) - What to build
2. **Role Definition** (docs/development/) - Who builds it  
3. **Priority Coordination** (NEXT.md) - Current focus
4. **Technical Execution** (AGENT.md) - How to build effectively

**Result**: AI agents get both coordination context and technical competency from a clear, non-duplicated document structure that scales with team growth.

## References

- **Builds on:** Current 7EP process, team identity docs, existing coordination patterns
- **Enables:** Systematic priority management and cross-team coordination
- **Validates:** 7EP-0014 foundation sprint coordination patterns
- **Integrates with:** Existing documentation architecture and git-based workflow
