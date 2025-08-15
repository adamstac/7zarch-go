# Merge Workflow: "Push to Main", "Submit PR", "Merge to Master"

**Purpose**: Standardized process for pushing local commits to remote and managing PRs  
**Framework**: Document Driven Development (7EP-0017)  
**Scope**: REMOTE OPERATIONS - Pushes to remote, creates PRs, merges branches  
**Trigger Phrases**: "Push to main", "Submit PR", "Create pull request", "Merge to master", "Ready for review", "Ship it"

**🚨 IMPORTANT**: This workflow handles remote operations. Use COMMIT.md workflow first for local commits.

## 🎯 Quick Decision Tree

```
"Push to main" / "Submit PR" / "Merge to master"
         |
    [Check Status]
         |
    Ready for PR?
    /         \
   YES         NO
    |           |
[Submit PR]  [Prepare First]
    |           |
[Wait Review] [Then Submit]
```

## 📋 Merge Scenarios

### Scenario A: Simple Push to Main ✅
**Status**: Local commits ready, just need to push to remote main

```bash
git status  # Should show "Your branch is ahead of 'origin/main' by X commits"
```

**Action**:
```bash
# Verify we're on main branch
git branch --show-current

# Push to remote
git push origin main
```

**Response**: "✅ Pushed to main: [X commits pushed]"  
**Role Coordination**: Update role files and NEXT.md if work completion affects other agents (use TEAM-UPDATE.md)

---

### Scenario B: Ready for PR ✅
**Status**: Work complete, tests passing, clean history

```bash
git status  # Clean working tree
make test   # All tests pass
git log --oneline -5  # Clean commit history
```

**Action**:
```bash
# Ensure we're on feature branch
git branch --show-current

# Final sync with main
git fetch origin main
git rebase origin/main

# Push feature branch
git push origin HEAD

# Create PR
gh pr create --title "[7EP-XXXX] Brief description" --body "$(cat <<'EOF'
## Summary
- [Key change 1]
- [Key change 2]
- [Key change 3]

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## 7EP Context
Addresses: 7EP-XXXX [7EP Name]
Phase: [1/2/3] of [total phases]

## Review Notes
[Any specific areas needing attention]
EOF
)"
```

**Response**: "✅ PR created: [PR URL] - Ready for review"

---

### Scenario C: Work Complete, Needs Cleanup 🔧
**Status**: Work done but messy commits or failing tests

**Action**:
1. **Clean up first**:
   ```bash
   # Fix any test failures
   make test
   
   # Interactive rebase to clean commits
   git rebase -i HEAD~[N]  # Squash/reword as needed
   
   # Then proceed with Scenario A
   ```

**Response**: "🔧 Cleaned up commits and tests, now submitting PR..."

---

### Scenario D: Work in Progress 🚧
**Status**: Partial work, not ready for final review

**Action**:
1. **Create draft PR for visibility**:
   ```bash
   git push origin HEAD
   gh pr create --draft --title "[WIP] [7EP-XXXX] Work description" --body "$(cat <<'EOF'
## 🚧 Work in Progress

### Completed
- [x] [Completed item 1]
- [x] [Completed item 2]

### In Progress  
- [ ] [Current work item]
- [ ] [Remaining item 1]
- [ ] [Remaining item 2]

### Status
Current focus: [what you're working on]
Next session: [what's planned next]
Blocker: [any blockers or dependencies]

**Note**: This is a WIP PR for visibility. Will mark ready for review when complete.
EOF
   )"
   ```

**Response**: "🚧 Created draft PR for work visibility: [PR URL]"

---

### Scenario E: Direct to Main (Hotfix) 🚨
**Status**: Critical fix needed immediately on main

**Action**:
```bash
# Ensure on main and up to date
git checkout main && git pull origin main

# Create hotfix branch
git checkout -b hotfix/critical-fix-$(date +%Y%m%d)

# Make minimal fix
[make changes]

# Commit with clear message
git add -A && git commit -m "fix: [critical issue description]

Problem: [what was broken]  
Solution: [minimal fix applied]
Impact: [who this affects]

Hotfix: Applied directly to main"

# Push and create PR
git push origin HEAD
gh pr create --title "🚨 HOTFIX: [issue]" --body "Critical fix - needs immediate review"
```

**Response**: "🚨 Hotfix PR created - flagged for immediate review"

---

## 🎨 PR Templates

### Standard Feature PR
```markdown
## Summary
Brief description of what this PR accomplishes

- [Key change 1]: [why this matters]
- [Key change 2]: [impact on users]  
- [Key change 3]: [technical improvement]

## Testing
- [ ] Unit tests pass (`make test`)
- [ ] Integration tests pass 
- [ ] Manual testing completed
- [ ] Performance impact assessed

## 7EP Context
**Addresses**: 7EP-XXXX [7EP Name]  
**Phase**: [X] of [Y] phases  
**Depends on**: [any dependencies]  

## Files Changed
- `file1.go`: [what changed]
- `file2.go`: [what changed]
- `docs/`: [documentation updates]

## Review Focus
Please pay attention to:
- [Specific area 1]: [why it needs attention]
- [Specific area 2]: [design decision rationale]
```

### WIP/Draft PR
```markdown  
## 🚧 Work in Progress

### Completed ✅
- [x] [Finished item with details]
- [x] [Another completed item]

### In Progress 🔄
- [ ] [Current focus item]
- [ ] [Next priority item]

### Planned 📋
- [ ] [Future work item]
- [ ] [Final integration step]

### Status Update
**Current Session**: [what you're working on now]  
**Next Session**: [what's planned next]  
**Estimated Completion**: [timeline if known]  
**Blockers**: [any dependencies or issues]

**Review Status**: Draft - not ready for final review
```

## 🚀 Branch Management

### Feature Branches
```bash
# Create feature branch from main
git checkout main && git pull
git checkout -b feat/7ep-XXXX-description

# Work on feature...

# Keep up to date with main
git fetch origin main
git rebase origin/main  # or merge if complex

# When ready, push and PR
git push origin HEAD
```

### Branch Naming Conventions
- `feat/7ep-XXXX-description` - New features
- `fix/issue-description` - Bug fixes  
- `docs/documentation-update` - Documentation only
- `refactor/code-improvement` - Code cleanup
- `hotfix/critical-fix-YYYYMMDD` - Emergency fixes

## 🔄 PR Lifecycle

### 1. Create PR
- Choose appropriate template
- Link to relevant 7EP
- Request specific reviewers
- Add labels if available

### 2. Address Feedback  
```bash
# Make requested changes
[edit files]

# Commit changes
git add -A
git commit -m "review: address feedback on [specific area]"

# Push updates
git push origin HEAD
```

### 3. Merge Process
```bash
# Once approved, merge via GitHub UI
# OR use gh cli:
gh pr merge [PR-number] --squash --delete-branch

# Clean up local branch
git checkout main
git pull origin main
git branch -d feat/7ep-XXXX-description
```

## 🎯 Quality Gates

### Before Creating PR
- [ ] All tests pass locally
- [ ] Code follows project conventions  
- [ ] Commit history is clean
- [ ] Documentation updated if needed
- [ ] No sensitive data in commits

### PR Description Requirements
- [ ] Clear summary of changes
- [ ] Testing checklist completed
- [ ] 7EP context provided
- [ ] Review focus areas identified

### Ready to Merge
- [ ] All CI checks pass
- [ ] Required approvals obtained
- [ ] No merge conflicts
- [ ] Documentation builds successfully

---

## 🚨 Emergency Protocols

### Broken Main Branch
1. **Identify the issue**: `git bisect` or check recent merges
2. **Create hotfix branch** immediately  
3. **Minimal fix** only - don't add features
4. **Fast-track review** - ping team immediately
5. **Deploy fix** as soon as approved

### Blocked by Dependencies
1. **Document the blocker** in PR description
2. **Create draft PR** for visibility  
3. **Link dependencies** in comments
4. **Update team coordination** in NEXT.md

---

**Remember**: When Adam says "submit PR" or "merge to master," use this workflow to handle the complexity while keeping the process smooth and predictable.