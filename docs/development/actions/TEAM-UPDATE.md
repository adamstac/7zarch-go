# Team Coordination Update Protocol

**Purpose**: Standardized process for cross-agent coordination updates  
**Framework**: Document Driven Development (7EP-0017, 7EP-0019)  
**When to Use**: Role changes, coordination shifts, assignment updates, blocking resolution

## 🎯 Quick Update Process

### 1. Role File Status Update
```bash
# Update your role file with current status
# Replace [AGENT] with: CLAUDE, AMP, AUGMENT, ADAM

# Update Active Work section
vim docs/development/roles/[AGENT].md

# Focus areas to update:
# - Active Work status changes (ACTIVE → COMPLETED, TODO → ACTIVE)
# - Next Priorities reordering based on new information
# - Coordination Needed updates (new blockers, resolved dependencies)
# - Recently Completed additions
```

### 2. Team Coordination Update
```bash
# Update shared team status
vim docs/development/NEXT.md

# Key sections to update:
# - Current Active Work (agent status changes)
# - Next Priorities (sequential order changes)  
# - Active Coordination Points (new dependencies, resolved blockers)
# - Blocked/Waiting (add/remove items)
```

### 3. Cross-Agent Notifications
```bash
# When your work affects other agents, update their context
# Example: CC completes backend work that unblocks AC frontend work

# Update other agent's role file Coordination Needed section
vim docs/development/roles/[OTHER-AGENT].md
# Remove blocker or add new dependency

# Document the handoff in commit message
git add docs/development/roles/
git commit -m "coordination: [brief description of status change affecting team]"
```

## 🔄 Update Patterns by Scenario

### Assignment Completion
**Trigger**: Finished significant work, ready for next assignment

**Updates Required**:
1. **Your Role File**: Move item from Active Work → Recently Completed
2. **NEXT.md**: Update your Current Active Work status
3. **Affected Agents**: Remove any coordination blockers your completion resolves

**Commit Pattern**: `coordination: completed [work-description], unblocked [other-agents]`

### New Assignment Acceptance  
**Trigger**: Taking on new work, changing priorities

**Updates Required**:
1. **Your Role File**: Add to Active Work, reorder Next Priorities
2. **NEXT.md**: Update Current Active Work, adjust coordination points
3. **Coordination**: Add any new dependencies to relevant agent role files

**Commit Pattern**: `coordination: accepting [assignment], coordination with [other-agents]`

### Blocking Issue Discovery
**Trigger**: Discovered dependency or technical blocker

**Updates Required**:
1. **Your Role File**: Add to Coordination Needed with clear blocker description
2. **NEXT.md**: Add to Blocked/Waiting with owner assignment
3. **Blocking Agent**: Update their role file with new coordination request

**Commit Pattern**: `coordination: blocked on [issue], assigned to [agent/person]`

### Strategic Context Changes
**Trigger**: New information affecting team direction

**Updates Required**:
1. **All Relevant Role Files**: Update strategic context in Implementation Notes
2. **NEXT.md**: Update strategic coordination points and priorities
3. **7EP Updates**: Update affected 7EP status or dependencies

**Commit Pattern**: `coordination: strategic update - [brief description]`

## 📋 Quality Checklist

### Before Committing Coordination Updates
- [ ] **Role file accuracy**: Your status reflects actual current state
- [ ] **Cross-agent impact**: Other agents' blockers/dependencies updated
- [ ] **NEXT.md consistency**: Team view matches individual role files
- [ ] **Clear communication**: Coordination changes easy to understand
- [ ] **No orphaned references**: No broken dependencies or outdated blockers

### Validation Commands
```bash
# Check for consistency issues
grep -r "Available.*Assignment\|READY.*assignment" docs/development/roles/
grep -B 2 -A 2 "Blocked.*\|Waiting.*" docs/development/NEXT.md

# Verify no duplicate coordination points
grep -n -A 3 "Coordination.*Needed\|Active.*Coordination" docs/development/
```

## 🚨 Common Coordination Mistakes

**Mistake**: Updating only your role file without cross-agent coordination  
**Solution**: Always check if your changes affect other agents' dependencies

**Mistake**: Stale coordination references after work completion  
**Solution**: Remove coordination requests from other agents when you complete blocking work

**Mistake**: Vague coordination status updates  
**Solution**: Use specific work descriptions and clear owner assignments

**Mistake**: Forgetting NEXT.md updates after role file changes  
**Solution**: Always update both - role files are individual, NEXT.md is team coordination

## 🎯 Integration with Workflow Actions

### COMMIT.md Integration
```bash
# Standard commit enhanced with coordination context
git commit -m "[type]: [description] - affects [coordination-context]"

# Examples:
git commit -m "feat: search performance optimization - unblocks AC search UI work"
git commit -m "fix: database migration bug - removes CC deployment blocker"
```

### MERGE.md Integration
```bash
# Include coordination updates in merge descriptions
# Update role files before merging to reflect coordination changes
```

### NEW-FEATURE.md Integration
```bash
# Define coordination patterns upfront in feature planning
# Identify cross-agent dependencies during 7EP creation
```

---

**Framework Integration**: This coordination protocol integrates with BOOTUP.md and SHUTDOWN.md to create complete agent lifecycle coordination, ensuring team coordination state is preserved and communicated effectively across all work transitions.

**Next Steps**: Use this protocol whenever role changes occur → Update coordination documents → Commit with clear context → Continue with coordinated work execution
