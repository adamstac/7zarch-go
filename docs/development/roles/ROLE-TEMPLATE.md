# [Agent Name] Context & Assignments

**Last Updated:** [YYYY-MM-DD HH:MM] ← Update each session  
**Status:** [Available|Active|Blocked] ← Current operational state  
**Current Focus:** [Brief description of primary work] ← What you're doing now

**Team Context**: See [`docs/development/TEAM-CONTEXT.md`](../TEAM-CONTEXT.md) for project structure and team overview

## 🎯 Current Assignments

### Active Work (This Week)
- **[Assignment Name]** - [STATUS] ([Brief description with context])
- **[Assignment Name]** - [STATUS] ([Brief description with context])

### Next Priorities
1. **[Priority 1]** - [Description with rationale]
2. **[Priority 2]** - [Description with rationale]
3. **[Priority 3]** - [Description with rationale]
4. **[Priority 4]** - [Description with rationale]

## 🔗 Coordination Needed
- **[Dependency/Blocker]** - [Description and who can resolve]
- **[Cross-agent work]** - [Description of coordination required]
- **[Strategic decision]** - [Description of decision needed]

## ✅ Recently Completed
- **[Recent work item]** - [Brief description and impact]
- **[Recent work item]** - [Brief description and impact]
- **[Recent work item]** - [Brief description and impact]

## 📝 Implementation Notes

### [Domain] Insights
- **[Key insight]** - [Description and implications]
- **[Technical pattern]** - [Description and reusability]
- **[Coordination pattern]** - [Description and effectiveness]

### [Specialization] Expertise
- **[Core skill/strength]** - [How this applies to project work]
- **[Domain knowledge]** - [How this applies to project work]
- **[Process expertise]** - [How this applies to project work]

---

## 📋 Role Template Usage Guide

### Required Sections
All role files MUST include these sections in this order:
1. **Header** (Last Updated, Status, Current Focus)
2. **Team Context Reference** (link to shared context)
3. **🎯 Current Assignments** (Active Work + Next Priorities)
4. **🔗 Coordination Needed** (dependencies and blockers)
5. **✅ Recently Completed** (achievements for context)
6. **📝 Implementation Notes** (domain insights and expertise)

### Header Field Standards
- **Last Updated**: ISO date format YYYY-MM-DD HH:MM
- **Status**: Available (ready for assignment) | Active (working on assignments) | Blocked (waiting on dependencies)
- **Current Focus**: One-line summary of primary work

### Status Indicators
Use these standard status indicators in Active Work:
- **ACTIVE** - Currently working on this
- **READY** - Prepared to begin when conditions met
- **BLOCKED** - Cannot proceed due to dependencies
- **WAITING** - Awaiting decision or input
- **ONGOING** - Continuous/maintenance work
- **COMPLETED** - Finished (move to Recently Completed)

### Content Boundaries
**Include in role files:**
- Your specific assignments and priorities
- Coordination needs specific to your work
- Technical insights from your domain expertise
- Work completion status and handoffs

**Do NOT include in role files:**
- Team structure information (use TEAM-CONTEXT.md)
- Strategic decision frameworks (use STRATEGIC-DECISION-FRAMEWORK.md)
- General project information (use TEAM-CONTEXT.md)
- Workflow instructions (use actions/ directory)

### Template Variations

#### Leadership Roles (Adam, Amp)
May include additional sections:
- **🎯 Strategic Context** (high-level project state)
- **📊 Decision Framework** (link to strategic tools)
- **🚀 Team Coordination** (cross-agent assignment patterns)

#### Implementation Roles (CC, AC)
Focus on:
- **Technical Implementation Notes** (architecture insights)
- **Feature Development Expertise** (domain-specific patterns)
- **Quality Assurance Patterns** (testing and validation approaches)

#### Dual-Role Patterns (AMP example)
- **Role Overview Table** (role switching matrix)
- **Quick Activation Commands** (role context switching)
- Preserve unique organizational value while maintaining standard structure

### Lifecycle Integration

#### Session Startup (BOOTUP.md integration)
Role files are read during bootup to:
- Load current assignments and priorities
- Identify coordination dependencies
- Validate assignment clarity

#### Daily Operations (TEAM-UPDATE.md integration)  
Role files are updated when:
- Assignment status changes (ACTIVE → COMPLETED)
- New blockers or dependencies discovered
- Strategic context shifts affect priorities
- Cross-agent coordination needs change

#### Session Shutdown (SHUTDOWN.md integration)
Role files are updated to:
- Preserve current work state
- Document completion status
- Set clear next-session priorities
- Maintain coordination handoffs

### Quality Checklist

#### Before Committing Role File Updates
- [ ] **Header current** - Last Updated reflects actual session time
- [ ] **Status accurate** - Status field matches actual availability
- [ ] **Assignments clear** - No ambiguous "Available for Assignment" states
- [ ] **Coordination updated** - Dependencies and blockers reflect current reality
- [ ] **Content boundaries** - No duplication with TEAM-CONTEXT.md or strategic docs
- [ ] **Standard structure** - All required sections present in correct order
- [ ] **Emoji prefixes** - All second-level headings use consistent emoji

#### Validation Commands
```bash
# Check header completeness
grep -E "Last Updated:|Status:|Current Focus:" docs/development/roles/[ROLE].md

# Verify no content boundary violations
grep -r "Adam Stacoviak\|Human Team\|AI Team" docs/development/roles/ || echo "✅ No team context duplication"

# Check for standard sections
grep -E "🎯 Current Assignments|🔗 Coordination|✅ Recently|📝 Implementation" docs/development/roles/[ROLE].md
```

### New Role Creation Process
1. **Copy template** - Start with this file as base structure
2. **Customize sections** - Adapt Implementation Notes to role expertise
3. **Set initial priorities** - Define first assignments and coordination needs  
4. **Link integrations** - Ensure bootup/shutdown processes reference role file
5. **Test lifecycle** - Run complete bootup → work → shutdown cycle
6. **Validate compliance** - Run quality checklist before first commit

---

**Framework Integration**: This template integrates with the complete agent lifecycle framework (7EP-0019) to ensure consistent agent operations from session startup through daily work execution to session shutdown.

**Next Steps**: Copy template → Customize for role → Test lifecycle → Begin productive work with full coordination
