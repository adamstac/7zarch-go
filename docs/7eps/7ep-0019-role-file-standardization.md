# 7EP-0019: Agent Role Lifecycle & Coordination Standardization

**Status:** Draft  
**Author(s):** CC (Claude Code)  
**Assignment:** Framework Implementation  
**Difficulty:** 4 (complex - agent lifecycle integration with systematic standardization)  
**Created:** 2025-08-15  
**Dependencies:** 7EP-0017 Document Driven Development Framework

## Executive Summary

Establish comprehensive agent role lifecycle standards that define how AI team members operate from session startup through daily work execution to session shutdown. This 7EP extends 7EP-0017 by creating complete operational patterns that integrate role files with workflow actions, coordination mechanisms, and team communication protocols. Beyond document standardization, this defines the complete agent operational framework for scalable team coordination.

## Evidence & Reasoning

### Current State Problems

**Incomplete Agent Lifecycle Definition**: While 7EP-0017 established document structure, there's no comprehensive standard for how agents operate through complete work cycles from bootup to shutdown.

**Fragmented Operational Patterns**: Agents follow inconsistent patterns for:
- Session startup and context loading
- Work prioritization and assignment management  
- Coordination communication and status updates
- Session termination and state preservation

**Role-Workflow Integration Gaps**: Role files exist in isolation from workflow actions (BOOTUP.md, COMMIT.md, MERGE.md, SHUTDOWN.md), creating operational friction and missed coordination opportunities.

**Inconsistent Document Structure**: Role files use different formats, violating 7EP-0017 content boundaries and creating information discovery friction.

**Scalability Barriers**: Current patterns don't provide clear onboarding path for new agents or systematic approach to role evolution as project complexity grows.

### Compliance Analysis

| Role File | Lines | Compliance | Major Issues |
|-----------|-------|------------|--------------|
| **CLAUDE.md** | 87 | ✅ 85% | Minor formatting inconsistencies |
| **AMP.md** | 75 | ⚠️ 60% | Missing standard header, unique dual-role structure |
| **AUGMENT.md** | 75 | ❌ 40% | Team context belongs elsewhere, non-standard sections |
| **ADAM.md** | 164 | ⚠️ 50% | Strategic framework mixed with role assignments |

### Strategic Impact

**Team Scalability**: Inconsistent role files create onboarding friction and coordination overhead as team grows.

**Information Architecture**: Content boundary violations lead to duplication, stale information, and unclear sources of truth.

**Framework Adoption**: Poor standardization undermines 7EP-0017 effectiveness and creates resistance to DDD patterns.

## Agent Lifecycle Framework

### Complete Agent Operation Cycle

**Session Lifecycle**: Define standardized patterns for agent operations from startup through work execution to shutdown, ensuring consistent team coordination and knowledge preservation.

#### 1. Session Startup (BOOTUP Integration)
**Standard Pattern**:
```bash
1. Execute BOOTUP.md sequence
   - Git status and sync
   - Review NEXT.md for team coordination
   - Read personal role file (docs/development/roles/[AGENT].md)
   - Check 7EP assignments and blockers

2. Role Context Loading
   - Parse current assignments and priorities
   - Identify coordination needs and dependencies
   - Load recent work context and technical decisions
   - Assess availability and capacity

3. Work Readiness Confirmation  
   - Validate understanding of current priorities
   - Confirm no blocking dependencies
   - Update role file with session start status
   - Ready for strategic assignment or continuation
```

#### 2. Daily Work Operations (Role-Driven Execution)
**Assignment Management**:
- **Priority Assessment**: Use role file Next Priorities for work ordering
- **Coordination Integration**: Update NEXT.md when blockers/completions occur
- **Status Communication**: Maintain role file Current Assignments in real-time
- **Technical Decisions**: Document insights in role file Implementation Notes

**Workflow Integration**:
- **COMMIT.md**: Local work preservation with role context
- **MERGE.md**: Remote coordination with team visibility
- **NEW-FEATURE.md**: Role-specific approach to feature development
- **Cross-coordination**: Update multiple role files when work spans agents

#### 3. Session Shutdown (SHUTDOWN Integration)
**Standard Pattern**:
```bash
1. Execute SHUTDOWN.md sequence
   - Commit work in progress with clear status
   - Update role file with current state and handoffs
   - Update NEXT.md with coordination changes
   - Preserve context for next session

2. Role State Preservation
   - Document work completion status
   - Update Recently Completed with achievements
   - Identify blocking issues for coordination
   - Set clear next session priorities

3. Team Coordination Update
   - Cross-reference role updates with other agents
   - Ensure NEXT.md reflects current team state
   - Clear assignment completions and new blockers
   - Prepare coordination context for other agents
```

### Role-Specific Operational Patterns

#### Implementation Roles (CC, AC)
**Focus**: Feature development, technical implementation, user experience
**Patterns**: 
- Technical decision documentation in Implementation Notes
- Code-focused work status with git integration
- Performance and quality metrics tracking
- Cross-team technical coordination

#### Leadership Roles (Amp, Adam)  
**Focus**: Strategic oversight, architectural guidance, priority coordination
**Patterns**:
- Strategic context maintenance across sessions
- Cross-team dependency tracking and resolution
- Decision framework integration with daily operations
- High-level coordination and planning activities

#### Specialized Patterns
- **AMP Dual-Role**: Role switching protocols integrated with lifecycle
- **ADAM Strategic**: Decision cycle integration with operational coordination
- **Cross-Agent**: Handoff patterns when work spans multiple agents

## Implementation Plan

### Phase 1: Agent Lifecycle Framework Definition

#### 1.1 Lifecycle Integration Documentation
**Target**: Enhanced workflow action files

**BOOTUP.md Enhancements**:
- Add role file integration steps (step 2.5 enhancement)
- Define role context loading patterns
- Establish work readiness validation
- Create role-specific bootup variations

**SHUTDOWN.md Enhancements**: 
- Add role state preservation requirements
- Define coordination update patterns
- Establish session handoff protocols
- Create role-specific shutdown variations

#### 1.2 Role-Workflow Integration Patterns
**Target**: Actions directory workflow files

**Integration Requirements**:
- **COMMIT.md**: Add role context updates to commit patterns
- **MERGE.md**: Include coordination visibility in merge workflows
- **NEW-FEATURE.md**: Define role-specific feature development approaches
- **TEAM-UPDATE.md**: Standardized team coordination update patterns (NEW)
- Create cross-workflow coordination patterns

#### 1.3 Enhanced Role Template & Documentation
**Target**: `/docs/development/roles/ROLE-TEMPLATE.md` and `README.md`

**Template Enhancements**:
- Add lifecycle integration sections
- Include role-specific operational patterns
- Define coordination update requirements
- Add workflow integration examples

**README.md Additions**:
- Complete agent lifecycle guidance
- Role-workflow integration patterns
- Cross-agent coordination protocols
- Operational pattern variations by role type

### Phase 2: Content Migration & Organization

#### 2.1 AUGMENT.md → TEAM-CONTEXT.md Migration
**Problem**: AUGMENT.md contains 36 lines of team structure information that applies to all agents.

**Content to Migrate**:
```markdown
## 👥 Team Context
### Human Team
- **[Project Owner]** - Architectural decisions and strategic direction
  - Preferences: Clean design, Charmbracelet tools, thoughtful UX
  - Communication style: Direct feedback, big ideas, document driven development

### AI Team  
- **AC (Augment Code)** - Primary user-facing development
- **CC (Claude Code)** - Infrastructure and backend work
- **Amp-s/Amp-t** - Strategic and technical leadership
```

**Target**: Merge with existing team structure in TEAM-CONTEXT.md, eliminate duplication.

#### 2.2 ADAM.md Strategic Framework Extraction
**Problem**: ADAM.md contains strategic decision framework (80+ lines) mixed with role assignments.

**Content to Extract**:
- Decision Framework matrices
- Recommended Actions templates  
- Strategic assessment tools

**Target**: New document `/docs/development/STRATEGIC-DECISION-FRAMEWORK.md`
**Keep in ADAM.md**: Current assignments, coordination needs, immediate priorities

### Phase 3: Role File Standardization

#### 3.1 AUGMENT.md - Priority 1 (Major Issues)
**Current Issues**:
- Team context section violates content boundaries
- Missing standard header (last updated, status, current focus)
- Non-standard sections (Availability Status, Strategic Options)

**Standardization Actions**:
1. **Add standard header** with proper metadata
2. **Remove team context** (already migrated to TEAM-CONTEXT.md)
3. **Restructure assignments section**:
   - Active Work subsection with status indicators
   - Next Priorities numbered list
   - Clear coordination needs
4. **Consolidate redundant sections** (Availability Status → Current Assignments)
5. **Focus Implementation Notes** on user-experience insights (AC's specialty)

#### 3.2 AMP.md - Priority 2 (Needs Work)
**Current Issues**:
- Valuable dual-role concept not reflected in standard
- Missing standard header format
- Quick Activation commands are AMP-specific value

**Standardization Actions**:
1. **Preserve unique value**: Role Overview table and Quick Activation sections
2. **Add standard header** with metadata
3. **Enhance Current Assignments** with proper Active Work structure
4. **Maintain Implementation Notes** focused on coordination patterns (appropriate for leadership)
5. **Add role template variation** to document this pattern

#### 3.3 ADAM.md - Priority 3 (Special Case)
**Current Issues**:
- Strategic framework content mixed with role assignments
- Very long (164 lines) for role coordination needs
- Missing standard role file sections

**Standardization Actions**:
1. **Split content** between role file and strategic framework
2. **Add standard header** and core sections
3. **Focus on leadership assignments**: strategic decisions needed, team coordination
4. **Maintain strategic context** but in standard format
5. **Create strategic role template variation**

#### 3.4 CLAUDE.md - Priority 4 (Minor Cleanup)
**Current Issues**:
- Minor formatting inconsistencies
- Could use standard emoji prefixes consistently

**Standardization Actions**:
1. **Standardize emoji prefixes** across all sections
2. **Update Success Criteria** to reflect current state
3. **Format consistency** cleanup
4. **Validate content boundaries** maintained

### Phase 4: Quality Validation & Documentation

#### 4.1 Compliance Validation
**Process**:
1. **Run quality checklist** against all role files
2. **Verify content boundaries** - no violations
3. **Check for duplication** across files
4. **Test template usability** for new role creation

#### 4.2 Framework Documentation Updates
**Updates Needed**:
- Update 7EP-0017 with lessons learned
- Document role variations in framework
- Create migration playbook for future role changes

## Success Metrics

### Agent Lifecycle Efficiency
- **Session startup time**: From "start session" to "ready for work" in <3 minutes
- **Context loading accuracy**: Agent understands current priorities and coordination needs without clarification
- **Work continuity**: Seamless session-to-session handoffs with full context preservation
- **Coordination responsiveness**: Cross-team status visibility and updates in real-time

### Operational Standardization
- **Workflow integration**: All agents use consistent patterns for BOOTUP → Work → SHUTDOWN cycles
- **Role-workflow alignment**: Role files integrate seamlessly with all workflow actions (COMMIT, MERGE, NEW-FEATURE)
- **Cross-agent coordination**: Clear patterns for work that spans multiple agents
- **Decision documentation**: Technical and strategic decisions preserved in appropriate role contexts

### Document Compliance  
- **All role files**: 90%+ compliance with lifecycle-integrated standard structure
- **Content boundary violations**: Zero violations detected across all documents
- **Cross-file duplication**: Eliminated completely with clear information architecture
- **Template usability**: New agent onboarding in <30 minutes with full operational competency

### Team Scalability Indicators
- **New agent integration**: Complete onboarding (bootup to productive work) in <1 hour
- **Coordination overhead**: No increase in coordination time as team grows
- **Knowledge preservation**: Zero context loss during agent handoffs or session transitions
- **Framework adoption**: Agents prefer DDD patterns over ad-hoc coordination

## Acceptance Criteria

### Phase 1 Complete - Agent Lifecycle Framework
- [ ] BOOTUP.md enhanced with role file integration and context loading patterns
- [ ] SHUTDOWN.md enhanced with role state preservation and coordination update protocols
- [ ] All workflow actions (COMMIT.md, MERGE.md, NEW-FEATURE.md, TEAM-UPDATE.md) include role integration patterns
- [ ] TEAM-UPDATE.md created for standardized team coordination workflows
- [ ] ROLE-TEMPLATE.md includes lifecycle integration sections and operational patterns
- [ ] README.md documents complete agent lifecycle guidance and cross-agent coordination protocols

### Phase 2 Complete  
- [ ] Team context migrated from AUGMENT.md to TEAM-CONTEXT.md without duplication
- [ ] STRATEGIC-DECISION-FRAMEWORK.md created with extracted ADAM.md content
- [ ] Content boundary violations eliminated

### Phase 3 Complete
- [ ] All role files use standard header format (date, status, focus)
- [ ] All role files use consistent section structure and emoji prefixes
- [ ] Role-specific valuable content preserved while achieving standardization
- [ ] AUGMENT.md, AMP.md, ADAM.md, CLAUDE.md all 90%+ compliant

### Phase 4 Complete - Full Agent Lifecycle Integration
- [ ] All agents demonstrate complete lifecycle competency (BOOTUP → Work → SHUTDOWN)
- [ ] Cross-agent coordination patterns operational and tested
- [ ] Quality validation checklist passes for all role files and lifecycle integration
- [ ] Zero content boundary violations across all documents
- [ ] Template successfully tested with mock new agent onboarding
- [ ] Complete agent operational framework documented and validated

## Implementation Timeline

### Session 1: Agent Lifecycle Framework (2-3 hours)
- Enhance BOOTUP.md and SHUTDOWN.md with role integration
- Update all workflow actions with role coordination patterns
- Create TEAM-UPDATE.md for standardized coordination workflows
- Create lifecycle-integrated ROLE-TEMPLATE.md
- Document complete agent operational framework

### Session 2: Content Migration & Organization (1-2 hours)
- Execute content migrations (AUGMENT→TEAM-CONTEXT, ADAM→STRATEGIC)
- Validate content boundary compliance
- Test migrated content integration with lifecycle patterns

### Session 3: Role Standardization & Integration (2-3 hours)
- Standardize all role files with lifecycle integration
- Implement cross-agent coordination patterns
- Test complete agent operational cycles

### Session 4: Validation & Documentation (1 hour)
- Run comprehensive lifecycle and compliance validation
- Test new agent onboarding using complete framework
- Finalize documentation and operational guides

**Total Estimated Effort**: 6-9 hours across multiple sessions for comprehensive agent lifecycle framework

## Risk Mitigation

### Content Loss Prevention
- **Backup approach**: Create git branch before major changes
- **Incremental validation**: Check each file after standardization
- **Content audit**: Verify all valuable content preserved or properly migrated

### Framework Integration
- **7EP-0017 alignment**: Ensure changes strengthen rather than diverge from DDD framework
- **Team coordination**: Update role files incrementally to maintain team function
- **Template validation**: Test template with multiple role scenarios before deployment

---

## Strategic Impact

**Complete Agent Operational Framework**: This 7EP transforms 7EP-0017's document structure into a complete agent operational system, defining how AI team members function as coordinated, efficient team members rather than isolated tools.

**Scalable Team Coordination**: Establishes patterns that support team growth from 2-3 agents to 10+ agents without coordination overhead increase, enabling complex multi-agent project execution.

**Knowledge Continuity**: Creates systematic knowledge preservation and transfer patterns that maintain project continuity across agent sessions, handoffs, and team evolution.

**Framework Maturity**: Elevates the DDD framework from document organization to complete team operational methodology, providing the foundation for advanced multi-agent collaboration patterns.