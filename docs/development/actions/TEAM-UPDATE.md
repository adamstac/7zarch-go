# Team Update Workflow: "Update the Team", "Let Everyone Know"

**Purpose**: Standardized process for coordinating status changes across team members  
**Framework**: Document Driven Development (7EP-0017) + Agent Role Lifecycle (7EP-0019)  
**Scope**: Team coordination and cross-agent communication  
**Trigger Phrases**: "Update the team", "Let everyone know", "Coordinate this change", "Status update needed"

## 🎯 Quick Decision Tree

```
"Update the team" / "Let everyone know"
         |
    [What Changed?]
         |
   Status/Priority/Blocker
    /         |         \
Personal    Cross-Team  Strategic
   |           |          |
[Update Role] [Update All] [Update + Escalate]
```

## 📋 Update Scenarios

### Scenario A: Personal Status Change ✅
**Trigger**: Work status changed, new assignment, completion
**Status**: Individual work update, affects personal coordination

**What Changed Examples**:
- Task completed → COMPLETED
- New blocker encountered → BLOCKED  
- Assignment received → ACTIVE
- Work handed off → HANDOFF-READY

**Action**:
```bash
# 1. Update personal role file
# docs/development/roles/CLAUDE.md - Current Assignments section
# Change status indicators and coordination needs

# 2. Update NEXT.md if affects team priorities
# Only if creates/resolves blockers or changes sequences

# 3. Brief status report
"✅ Team updated: [brief change description]
Personal status: [new status]
Team impact: [minimal/coordination needed/priority shift]"
```

---

### Scenario B: Cross-Team Coordination Change 🔗
**Trigger**: Handoffs needed, dependencies changed, blocking relationships
**Status**: Affects multiple team members or work sequences

**What Changed Examples**:
- Work ready for handoff to another agent
- Blocker resolved, unblocking others
- New dependency identified
- Cross-team work assignment

**Action**:
```bash
# 1. Update personal role file (same as Scenario A)
# docs/development/roles/CLAUDE.md

# 2. Update NEXT.md coordination sections
# - Active Coordination Points
# - Blocked/Waiting sections
# - Next Priorities if sequence changes

# 3. Update relevant role files if handoff involved
# Mention in other agent's Coordination Needed section

# 4. Coordination report
"🔗 Team coordination updated: [change description]
Affects: [list of team members/roles]
Next actions: [what needs to happen]
Timeline: [when handoff/unblocking occurs]"
```

---

### Scenario C: Strategic Priority Change 🎯
**Trigger**: Strategic decisions, major priority shifts, project direction changes
**Status**: Affects overall team direction and resource allocation

**What Changed Examples**:
- Strategic direction decision from Adam
- Major feature priority shift
- New 7EP activation or completion
- Resource allocation changes

**Action**:
```bash
# 1. Update all relevant role files
# Adjust Current Assignments and Next Priorities

# 2. Update NEXT.md comprehensively
# - Current Active Work
# - Next Priorities (Sequential)
# - Active Coordination Points
# - Strategic context sections

# 3. Update active 7EPs if implementation affected
# Coordination sections and timeline impacts

# 4. Strategic update report
"🎯 Strategic team update: [major change description]
Priority impact: [how priorities shifted]
Team assignments: [new/changed assignments]
Coordination: [new coordination needs]
Timeline: [impact on delivery schedules]"
```

---

### Scenario D: Blocker Resolution/Creation 🚧
**Trigger**: Blockers encountered or resolved that affect team flow
**Status**: Immediate coordination needed to maintain team momentum

**What Changed Examples**:
- Critical dependency resolved
- New technical blocker discovered  
- External decision needed
- Resource availability change

**Action**:
```bash
# 1. Immediate NEXT.md update
# Blocked/Waiting section - add/remove items
# Next Priorities - adjust sequence if unblocked

# 2. Update affected role files
# Status changes for blocked/unblocked work
# Coordination Needed adjustments

# 3. Escalation if needed
# If blocker requires strategic decision → mention Adam
# If blocker requires architectural guidance → mention Amp

# 4. Blocker status report
"🚧 Blocker update: [blocker description]
Status: [RESOLVED/NEW/ESCALATED]
Affects: [team members/work streams]
Action needed: [who needs to act]
Urgency: [timeline for resolution]"
```

---

## 🎨 Update Templates

### Personal Status Update
```markdown
**Role File Update** (docs/development/roles/CLAUDE.md):
- **[Work Item]** - [OLD STATUS] → [NEW STATUS] ([reason/context])
- **Coordination Needed**: [new/changed coordination requirements]

**NEXT.md Impact**: [None/Minor/Coordination/Priority]
```

### Cross-Team Coordination Update  
```markdown
**Multi-Role Update**:
- **CLAUDE.md**: [status change]
- **NEXT.md**: [coordination section updated]
- **[Other Role].md**: [handoff/dependency noted]

**Team Impact**: [description of cross-team effects]
**Next Actions**: [what each role needs to do]
```

### Strategic Priority Update
```markdown
**Strategic Change**: [decision/priority shift description]

**Role File Updates**:
- **CLAUDE.md**: [assignment changes]
- **[Role].md**: [priority adjustments]

**NEXT.md Updates**:
- Current Active Work: [changes]
- Next Priorities: [new sequence]
- Coordination Points: [new dependencies]

**7EP Impact**: [active 7EP coordination changes]
```

## 🚀 Update Patterns

### Status Indicators (Consistent Across All Docs)
- **ACTIVE** - Currently working on this task
- **COMPLETED** - Task finished, ready for handoff or next phase
- **BLOCKED** - Cannot proceed due to dependency
- **READY** - Ready to begin when prioritized  
- **HANDOFF-READY** - Complete and ready to hand off to another team member
- **WAITING** - Waiting for external input or decision
- **PAUSED** - Temporarily suspended (lower priority)

### Update Frequency Guidelines
- **Real-time**: Blockers, completions, urgent handoffs
- **Session-based**: Status changes, assignment updates
- **Daily**: Coordination adjustments, priority fine-tuning
- **Strategic**: Major direction changes, 7EP phase transitions

### Cross-Reference Patterns
```markdown
# In personal role file
## 🔗 Coordination Needed
- **Handoff to [ROLE]:** [specific work] (updated in NEXT.md)

# In NEXT.md
## 🔗 Active Coordination Points  
- **CC → AMP:** [specific handoff] (CC status: HANDOFF-READY)
```

## 📊 Integration with 7EP-0019

### Agent Lifecycle Integration
**Daily Operations**: Team updates are core part of role-driven execution patterns
- Personal assignment management with team visibility
- Cross-agent coordination through standardized update patterns
- Strategic alignment through systematic priority communication

**Session Transitions**: BOOTUP and SHUTDOWN workflows enhanced with team update patterns
- Session startup includes team coordination context loading
- Session shutdown includes status preservation for team visibility

**Workflow Integration**: All action workflows (COMMIT, MERGE, NEW-FEATURE) include team update decision points
- When to update team coordination during development workflow
- How updates integrate with git-based change management
- Cross-workflow coordination through team update patterns

### Framework Enhancement
**7EP-0017 Extension**: Team updates operationalize DDD coordination patterns
**7EP-0019 Component**: Critical piece of complete agent lifecycle framework  
**Operational Effectiveness**: Systematic team communication reduces coordination overhead

## 🔄 Response Format

### Update Confirmation
```
✅ Team updated: [brief summary]

Updated Documents:
- [Role].md: [what changed]
- NEXT.md: [coordination changes]
- [Other docs]: [additional updates]

Team Impact: [None/Coordination/Priority/Strategic]
Next Actions: [what happens next]
```

### Update with Coordination Needs
```
🔗 Team coordination updated: [change description]

Cross-Team Effects:
- [Role 1]: [what they need to know/do]
- [Role 2]: [coordination requirements]

Timeline: [when coordination happens]
Blockers: [any new dependencies created]
```

## 🎯 Quality Guidelines

### Before Updating
- [ ] Identify what actually changed (status, priority, dependency)
- [ ] Determine who is affected (individual, cross-team, strategic)
- [ ] Choose appropriate update scope (role file, NEXT.md, 7EPs)
- [ ] Consider timing (immediate, session-based, strategic)

### During Update
- [ ] Use consistent status indicators across all documents
- [ ] Update coordination needs and dependencies
- [ ] Cross-reference related changes in other documents
- [ ] Maintain clear timeline and sequence information

### After Update
- [ ] Verify no coordination information lost or duplicated
- [ ] Confirm affected team members have visibility
- [ ] Check that priorities and sequences remain clear
- [ ] Validate integration with ongoing workflows

---

**Remember**: Team updates are not just status reports - they're coordination actions that keep team momentum flowing and prevent coordination overhead from scaling with team size.