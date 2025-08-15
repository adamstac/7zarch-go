# Role File Standard Template

**Purpose**: Standardized format for all team member role files in `/docs/development/roles/`  
**Framework**: Document Driven Development (7EP-0017)  
**Scope**: Current assignments, coordination needs, and operational context only

## 📋 Standard Structure

### Header (Required)
```markdown
# [Role Name] ([Agent Code]) Current Assignments

**Last Updated:** [Date]  
**Status:** [Available/Active/Blocked]  
**Current Focus:** [Brief current priority]
```

### Core Sections (Required)

#### 1. Current Assignments
```markdown
## 🎯 Current Assignments

### Active Work
- **[Work Item]** - [STATUS] ([brief context])
- **[Work Item]** - [STATUS] ([brief context])

### Next Priorities  
1. **[Priority 1]** - [Brief description]
2. **[Priority 2]** - [Brief description]
3. **[Priority 3]** - [Brief description]
```

#### 2. Coordination Needed
```markdown
## 🔗 Coordination Needed
- **[Item]:** [What coordination is needed and with whom]
- **[Blocker]:** [What's blocking progress and who can resolve]
- **[Dependency]:** [What this role is waiting on]
```

#### 3. Recently Completed
```markdown
## ✅ Recently Completed
- **[Achievement]** - [Brief outcome and context]
- **[Achievement]** - [Brief outcome and context]
```

### Optional Sections (Use When Relevant)

#### Implementation Notes (For Technical Roles)
```markdown
## 📝 Implementation Notes

### [Category] Insights
- [Technical insight or decision]
- [Architecture pattern or lesson learned]

### [Category] Decisions  
- [Recent technical or process decisions]
- [Rationale for choices made]
```

#### Role-Specific Context (When Unique)
```markdown
## 🎯 [Role] Identity & Approach
[Brief description of role personality, approach, or unique characteristics]
```

## 🚫 Content Rules

### What Belongs in Role Files
✅ **Current assignments and work status**  
✅ **Coordination needs and blockers**  
✅ **Recent accomplishments with context**  
✅ **Implementation insights from recent work**  
✅ **Role-specific operational context**

### What Does NOT Belong
❌ **Technical guidance** (belongs in `/AGENT.md`)  
❌ **Team context** (belongs in `/docs/development/TEAM-CONTEXT.md`)  
❌ **Project overview** (belongs in `/docs/development/TEAM-CONTEXT.md`)  
❌ **Historical context** (keep only recent/relevant items)  
❌ **General documentation** (belongs in appropriate docs)

## 📏 Consistency Guidelines

### Section Naming
- Use **exact standard names**: "Current Assignments", "Coordination Needed", "Recently Completed"
- Use consistent emoji prefixes: 🎯 (assignments), 🔗 (coordination), ✅ (completed)
- Sub-sections can be role-specific but should be consistent within the role

### Content Style
- **Active Work**: Present tense, status in CAPS (READY, ACTIVE, BLOCKED)
- **Next Priorities**: Numbered list, brief descriptions
- **Coordination Needed**: Specific items, not general statements
- **Recently Completed**: Past tense, focus on outcomes and context

### Update Frequency
- **Header date**: Update every time file is modified
- **Status**: Reflect current availability (Available/Active/Blocked)
- **Active Work**: Update as work progresses or completes
- **Recently Completed**: Keep last 3-5 significant achievements

## 🎯 Role-Specific Variations

### Project Lead (ADAM.md)
- **Additional sections**: Strategic Context, Leadership Actions, Decision Framework
- **Focus**: Strategic decisions, team coordination, priority setting
- **Unique elements**: Decision matrices, recommended actions

### Dual Role (AMP.md) 
- **Additional sections**: Role Overview, Quick Activation
- **Focus**: Strategic vs technical role separation
- **Unique elements**: Activation commands for role switching

### Implementation Roles (CLAUDE.md, AUGMENT.md)
- **Additional sections**: Implementation Notes, Technical Insights
- **Focus**: Current assignments, technical decisions, coordination
- **Unique elements**: Architecture insights, design decisions

## ✅ Quality Checklist

Before updating a role file:
- [ ] Header information is current (date, status, focus)
- [ ] Active work reflects actual current state
- [ ] Coordination needs are specific and actionable
- [ ] Recently completed items provide useful context
- [ ] Implementation notes focus on insights, not just facts
- [ ] Content follows "what belongs" rules
- [ ] Section names match standard format
- [ ] File serves as single source of truth for role status

## 🔄 Maintenance Protocol

### Daily Updates
- Update status and current focus when work changes
- Move completed items from Active Work to Recently Completed
- Add new coordination needs as they arise

### Session End Updates  
- Reflect final session state in Active Work
- Update coordination needs based on session outcomes
- Add significant accomplishments to Recently Completed

### Weekly Review
- Archive old Recently Completed items (keep last 3-5)
- Review and update Next Priorities based on strategic direction
- Ensure Implementation Notes remain relevant and valuable

---

**Remember**: Role files are operational documents, not historical archives. Keep them current, focused, and directly useful for team coordination.