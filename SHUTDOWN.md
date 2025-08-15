# AI Team Member Shutdown Guide

**Purpose**: Standardized shutdown sequence for all AI team members  
**Framework**: Document Driven Development (7EP-0017)  
**Counterpart**: See `BOOTUP.md` for session initialization

## 🤖 Shutdown Protocol

### Requesting Shutdown Status
**User says**: "Can we shutdown?" / "Ready to shutdown?" / "Shutdown status?"
**AI responds**: Pre-shutdown status check with current state

### Initiating Shutdown  
**User says**: "Shutdown. Do it." / "Shutdown now" / "Execute shutdown"
**AI responds**: "Are you sure you want to shutdown? (Yes/No)"
**User confirms**: "Yes" / "Confirmed" / "Do it"
**AI executes**: Full shutdown sequence

## 🛑 Shutdown Sequence

Execute these steps in order when ending a session:

### 1. Update Role Documents
```bash
# Update your current status and progress
vim docs/development/CLAUDE.md

# Key updates needed:
# - Mark completed tasks as ✅ Complete
# - Update "Active Work" section with current state
# - Note any blockers or handoff points
# - Update "Next Priorities" for next session
```

### 2. Update Team Coordination  
```bash
# Update shared team status
vim docs/development/NEXT.md

# Key updates needed:
# - Move completed items to "Recently Completed"
# - Update active work status and % complete
# - Add any new blockers or dependencies
# - Clear any resolved blockers
```

### 3. Handle Work in Progress
```bash
# Check git status
git status

# For completed work:
git add [files]
git commit -m "descriptive commit message"

# For work in progress:
# Option A: Commit WIP with clear message
git add [files] 
git commit -m "wip: [description] - [what's left to do]"

# Option B: Stash if not ready to commit
git stash push -m "session-end: [description]"
```

### 4. Push Changes
```bash
# Push committed work
git push origin main

# Verify push succeeded
git status
```

### 5. Document Session Context
```bash
# If session produced significant work, add brief summary
echo "## Session Summary $(date +%Y-%m-%d)" >> SESSION-NOTES.md
echo "- [Key accomplishment 1]" >> SESSION-NOTES.md  
echo "- [Key accomplishment 2]" >> SESSION-NOTES.md
echo "- [Next session should focus on...]" >> SESSION-NOTES.md
echo "" >> SESSION-NOTES.md
```

## 🔄 Handoff Protocol

### For Completed Work
- [ ] All changes committed and pushed
- [ ] Role document updated with completed status
- [ ] NEXT.md reflects new state
- [ ] Any blocking dependencies resolved or documented

### For Work in Progress  
- [ ] Clear WIP commit or detailed stash message
- [ ] Role document shows current state and next steps
- [ ] NEXT.md updated with progress percentage
- [ ] Clear TODO comments in code for continuation points

### For Blocked Work
- [ ] Blocker clearly documented in NEXT.md
- [ ] Role document indicates who/what is blocking progress
- [ ] Any research or attempted solutions documented
- [ ] Alternative approaches noted for next session

## 📋 Pre-Shutdown Checklist

**Code State**:
- [ ] No broken builds left behind
- [ ] No obvious syntax errors in modified files
- [ ] Tests pass (if modified test-related code)

**Documentation State**:
- [ ] Role document accurately reflects current state
- [ ] NEXT.md shows real coordination status
- [ ] Any new patterns or decisions documented

**Repository State**:
- [ ] All important work committed
- [ ] Commits have descriptive messages
- [ ] Changes pushed to remote
- [ ] No sensitive information committed

## 🚨 Emergency Shutdown

If you need to shut down quickly:

```bash
# Minimal sequence
git add -A
git commit -m "session-end: emergency shutdown - [brief status]"
git push origin main

# Update coordination (critical)
echo "Emergency shutdown $(date)" >> docs/development/NEXT.md
echo "Last work: [brief description]" >> docs/development/NEXT.md
```

## 🎯 Success Metrics

**Good Shutdown**:
- Next session can resume in <2 minutes
- No context loss about current work
- Team coordination is accurate
- Repository state is clean

**Poor Shutdown**:
- Next session requires >10 minutes to understand state
- Work in progress is unclear
- Team coordination is out of sync
- Repository has uncommitted changes

## 📝 Session Patterns

### Daily Work Sessions
- Focus on updating role docs and NEXT.md
- Commit frequently during session
- Minimal bootdown overhead

### Major Feature Work
- Ensure progress is documented in 7EP
- Update architecture docs if needed
- Add session summary for complex changes

### Collaboration Handoffs
- Update role docs for both agents involved
- Clear documentation of handoff points
- Test that next agent can continue immediately

---

**Framework Integration**: This bootdown sequence maintains the DDD operational framework by ensuring all coordination documents reflect actual project state at session end.

**Next Steps**: Execute bootdown → Update role documents → Push changes → Ready for next session startup via BOOTUP.md