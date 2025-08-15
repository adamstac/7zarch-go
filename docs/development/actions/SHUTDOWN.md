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
vim docs/development/roles/CLAUDE.md

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

### 5. Complete Session Log
```bash
# Get session log file from bootup
if [ -f .session-active ]; then
    source .session-active
    SESSION_LOG=${SESSION_LOG_FILE}
else
    # Fallback if no active session file
    DATE_STAMP=$(date +%Y-%m-%d_%H-%M-%S)
    SESSION_LOG="docs/logs/session-${DATE_STAMP}.md"
fi

# Calculate session timing
SESSION_END=$(date "+%Y-%m-%d %H:%M:%S")
SESSION_START_TIME=$(grep "Start Time:" "$SESSION_LOG" | sed 's/.*Start Time:** //' || echo "Unknown")

# Calculate duration if we can parse start time  
if [[ "$SESSION_START_TIME" != "Unknown" ]]; then
    START_EPOCH=$(date -j -f "%Y-%m-%d %H:%M:%S" "$SESSION_START_TIME" +%s 2>/dev/null || date +%s)
    END_EPOCH=$(date +%s)
    SESSION_DURATION=$((END_EPOCH - START_EPOCH))
    SESSION_HOURS=$((SESSION_DURATION / 3600))
    SESSION_MINUTES=$(((SESSION_DURATION % 3600) / 60))
    DURATION_TEXT="${SESSION_HOURS}h ${SESSION_MINUTES}m"
else
    DURATION_TEXT="Unknown"
fi

# Find session commits - look for commits since session start
# Get the commit that started this session (contains "session: start")
SESSION_START_COMMIT=$(git log --oneline --grep="session: start" -1 --format="%H" 2>/dev/null)

# If we found a start commit, get all commits since then, otherwise get recent commits
if [[ -n "$SESSION_START_COMMIT" ]]; then
    SESSION_COMMITS=$(git log --oneline --format="- [\`%h\`](../../commit/%H) %s" ${SESSION_START_COMMIT}..HEAD 2>/dev/null)
    COMMIT_COUNT=$(git rev-list --count ${SESSION_START_COMMIT}..HEAD 2>/dev/null || echo "0")
else
    # Fallback: get last 5 commits
    SESSION_COMMITS=$(git log -5 --oneline --format="- [\`%h\`](../../commit/%H) %s" 2>/dev/null)
    COMMIT_COUNT=$(echo "$SESSION_COMMITS" | wc -l | tr -d ' ')
fi

# Append shutdown summary to existing session log
cat >> "$SESSION_LOG" << EOF

## 🛑 Session Shutdown - $(date)

### ⏱️ Final Timing
- **End Time:** ${SESSION_END}
- **Duration:** ${DURATION_TEXT}
- **Status:** 🔴 **OFFLINE** - Session terminated successfully

### 📋 Session Deliverables
- [Key accomplishment 1]
- [Key accomplishment 2]

### 📁 Files Changed
$(git diff --name-only HEAD~1 HEAD 2>/dev/null || echo "- No changes detected")

### 🎯 Next Session Priorities  
- [What should happen next]

### 🔄 Role State Updates
$(echo "Agent Role File Updates:")
$(grep -A 3 "Active Work\|Current Assignments" docs/development/roles/[AGENT].md | head -5)

### 🤝 Team Coordination Changes
$(echo "NEXT.md Coordination Updates:")
$(git diff HEAD~1 HEAD docs/development/NEXT.md | grep "^+\|^-" | head -5 || echo "- No coordination changes")

### 📊 Session Stats
- **Commits This Session:** ${COMMIT_COUNT}
- **Files Modified:** $(git diff --name-only HEAD~1 HEAD 2>/dev/null | wc -l | tr -d ' ')

### 🔗 Session Commits
${SESSION_COMMITS:-"- No commits found for this session"}

---
*Session completed by DDD Framework shutdown process*
EOF

# Clean up session tracking
rm -f .session-active

# Add final log to git
git add "$SESSION_LOG"
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