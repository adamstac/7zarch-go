# Role Files Directory

**Purpose**: Team member assignment and coordination documents  
**Framework**: Document Driven Development (7EP-0017)  
**Template**: Use `ROLE-TEMPLATE.md` as starting point for new role files

## 📁 Directory Contents

- **ROLE-TEMPLATE.md** - Clean template for creating new role files
- **ADAM.md** - Project lead strategic context and decisions
- **AMP.md** - Dual leadership roles (strategic + technical)
- **AUGMENT.md** - AC current assignments and coordination
- **CLAUDE.md** - CC current assignments and coordination

## 📋 Role File Standard

### Purpose & Scope
Role files track **current operational status only**:
- Current assignments and work in progress
- Coordination needs and blockers  
- Recent accomplishments with context
- Implementation insights from recent work
- Role-specific operational patterns

### Required Structure

#### Header
```markdown
# [Role Name] ([Agent Code]) Current Assignments

**Last Updated:** [Date]  
**Status:** [Available/Active/Blocked]  
**Current Focus:** [Brief current priority]
```

#### Core Sections
1. **🎯 Current Assignments** - Active work and next priorities
2. **🔗 Coordination Needed** - Blockers, dependencies, coordination points  
3. **✅ Recently Completed** - Recent achievements with context

#### Optional Sections
- **📝 Implementation Notes** - Technical insights and decisions (for implementation roles)
- **🎯 [Role] Identity & Approach** - Role personality and approach (when unique)

## 🚫 Content Boundaries

### What Belongs Here
✅ **Current work status and assignments**  
✅ **Active coordination needs and blockers**  
✅ **Recent accomplishments with context**  
✅ **Implementation insights from current work**  
✅ **Role-specific operational context**

### What Belongs Elsewhere
❌ **Technical guidance** → `/AGENT.md`  
❌ **Team context** → `/docs/development/TEAM-CONTEXT.md`  
❌ **Project overview** → `/docs/development/TEAM-CONTEXT.md`  
❌ **Workflow processes** → `/docs/development/actions/`  
❌ **Historical archives** → Keep only recent/relevant items

## 📏 Consistency Guidelines

### Section Naming Standards
- **Exact names**: "Current Assignments", "Coordination Needed", "Recently Completed"
- **Emoji prefixes**: 🎯 (assignments), 🔗 (coordination), ✅ (completed)
- **Sub-sections**: Can be role-specific but consistent within each role

### Content Style
- **Active Work**: Present tense, status in CAPS (READY, ACTIVE, BLOCKED)
- **Next Priorities**: Numbered list with brief descriptions
- **Coordination Needed**: Specific actionable items, not general statements
- **Recently Completed**: Past tense, focus on outcomes and context

### Update Frequency
- **Header date**: Every modification
- **Status**: Reflects current availability
- **Active Work**: Updated as work progresses
- **Recently Completed**: Keep last 3-5 significant achievements

## 🎯 Role-Specific Patterns

### Project Lead (ADAM.md)
**Focus**: Strategic decisions, team coordination, priority setting  
**Unique sections**: Strategic Context, Leadership Actions, Decision Framework  
**Content**: Decision matrices, recommended actions, strategic options

### Dual Leadership (AMP.md)
**Focus**: Strategic vs technical role separation  
**Unique sections**: Role Overview, Quick Activation  
**Content**: Role switching commands, dual responsibility coordination

### Implementation Roles (CLAUDE.md, AUGMENT.md)
**Focus**: Current assignments, technical decisions, coordination  
**Unique sections**: Implementation Notes, Technical Insights  
**Content**: Architecture decisions, design patterns, coordination insights

## 🔄 Maintenance Protocol

### Daily Operations
- Update status when work changes
- Move completed items to Recently Completed  
- Add coordination needs as they arise

### Session End Updates
- Reflect final session state in Active Work
- Update coordination based on session outcomes
- Document significant accomplishments

### Weekly Review
- Archive old Recently Completed items (keep last 3-5)
- Update Next Priorities based on strategic direction
- Ensure Implementation Notes remain current and valuable

## 🛠️ Creating New Role Files

1. **Copy template**: `cp ROLE-TEMPLATE.md NEW-ROLE.md`
2. **Customize header**: Update role name, agent code, initial status
3. **Fill core sections**: Add actual assignments, coordination needs, context
4. **Remove unused sections**: Delete optional sections not needed for this role
5. **Follow naming**: Use consistent section names and emoji prefixes

## ✅ Quality Standards

Every role file should:
- [ ] Have current header information (date, status, focus)
- [ ] Reflect actual current work state in Active Work
- [ ] List specific, actionable coordination needs
- [ ] Provide useful context in Recently Completed
- [ ] Focus on operational insights in Implementation Notes
- [ ] Follow content boundary rules (what belongs vs doesn't)
- [ ] Use standard section names and format
- [ ] Serve as single source of truth for role status

## 🚨 Common Issues

### Content Creep
**Problem**: Files become historical archives or reference guides  
**Solution**: Regular cleanup, focus on current operational needs only

### Inconsistent Structure  
**Problem**: Each role uses different section names or formats  
**Solution**: Follow standard exactly, customize content not structure

### Stale Information
**Problem**: Files not updated regularly, become misleading  
**Solution**: Update dates and status every session, weekly review cycle

### Scope Confusion
**Problem**: Technical guides mixed with assignments  
**Solution**: Use content boundaries - technical info goes in `/AGENT.md`

---

**Remember**: Role files are living operational documents, not historical records. Keep them current, focused, and directly useful for daily team coordination.