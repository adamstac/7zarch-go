# Commit Workflow: "Commit Those Changes"

**Purpose**: Standardized process for when Adam (or team lead) says "commit those changes"  
**Framework**: Document Driven Development (7EP-0017)  
**Scope**: LOCAL COMMITS ONLY - Does not push to remote  
**Trigger Phrases**: "Commit those changes", "Let's commit", "Save that work"

**🚨 IMPORTANT**: This workflow creates LOCAL commits only. It does NOT push to remote/main. Use MERGE.md workflow for pushing to remote.

## 🎯 Quick Decision Tree

```
"Commit those changes"
         |
    [Check Status]
         |
    Is it clear?
    /         \
   YES         NO
    |           |
[Scenario A]  [More Complex]
              /     |     \
           Mixed  Messy  Unclear
             |      |       |
        [Scen B] [Scen C] [Scen D]
```

## 📋 Scenario Workflows

### Scenario A: Clean & Clear ✅
**Status**: Only intended changes, no surprises
```bash
git status  # All changes are expected
```

**Action**:
```bash
git add -A
git commit -m "feat/fix/docs: [clear description]

- [Bullet point of key change 1]
- [Bullet point of key change 2]"
```

**Response**: "✅ Committed locally: [brief summary]"  
**Next**: Use MERGE.md workflow if you want to push to remote

---

### Scenario B: Mixed Changes 🔀
**Status**: Some work complete, some WIP, some experimental
```bash
git status  # Shows completed work + experiments + WIP
```

**Action**:
1. **First, clarify**:
   ```
   "I see:
   - ✅ Completed: [files related to X]
   - 🔧 WIP: [files for Y]
   - 🧪 Experimental: [test files Z]
   
   Should I commit:
   a) Everything as-is
   b) Only completed work
   c) Completed + WIP (not experimental)"
   ```

2. **Then commit based on response**:
   ```bash
   # Option b example:
   git add [specific completed files]
   git commit -m "feat: [completed feature]
   
   - [What was accomplished]
   - WIP and experimental changes not included"
   ```

---

### Scenario C: Messy State 🚧
**Status**: Broken tests, incomplete refactoring, or mixed concerns
```bash
git status      # Many files changed
git diff        # Changes span multiple features
make test       # Tests failing
```

**Action**:
1. **Report the situation**:
   ```
   "⚠️ Current state:
   - Changed files: X across Y directories
   - Test status: [passing/failing]
   - Changes span: [multiple features/single feature]
   
   Options:
   a) Commit as 'wip: [description]' - preserves current state
   b) Let me organize into logical commits
   c) Stash and clean up first"
   ```

2. **Handle based on response**:
   ```bash
   # Option a - WIP commit
   git add -A
   git commit -m "wip: [current work state]
   
   Status:
   - [What's working]
   - [What's broken]
   - [What's next]
   
   Tests: [status]
   Next: [what needs to happen]"
   ```

---

### Scenario D: Unclear Intent ❓
**Status**: Not sure what "those changes" refers to
```bash
git status          # Multiple unrelated changes
git log --oneline -3  # Recent commits for context
```

**Action**:
1. **Seek clarification**:
   ```
   "I see several different changes:
   
   Recent work (last 3 commits):
   - [commit 1]
   - [commit 2] 
   - [commit 3]
   
   Current uncommitted:
   - [Group 1]: [description] (X files)
   - [Group 2]: [description] (Y files)
   - [Group 3]: [description] (Z files)
   
   Which should I commit?"
   ```

2. **Wait for clarification, then proceed with Scenario A/B/C**

---

## 🎨 Commit Message Templates

### Feature Complete
```
feat: [what it does]

- [Implementation detail 1]
- [Implementation detail 2]
- [Files affected or approach]

[Optional: Closes #issue]
```

### Bug Fix
```
fix: [what was broken]

Problem: [what was happening]
Solution: [what you did]
Result: [what now works]
```

### Documentation
```
docs: [what was documented]

- [Section/file 1]: [what was added/changed]
- [Section/file 2]: [what was added/changed]
```

### Work in Progress
```
wip: [what you're working on]

Progress:
- [x] Completed part
- [ ] In progress part
- [ ] Not started part

Status: [why stopping now]
Next: [what happens next session]
```

### Mixed/Multiple
```
feat: [primary change] + misc improvements

Primary:
- [Main feature/fix description]

Additional:
- fix: [small fix included]
- docs: [documentation update]
- refactor: [code cleanup]
```

---

## 🚀 Quick Rules

### Always Do:
1. **Run `git status` first** - Understand what's there
2. **Check for surprises** - Files you didn't expect
3. **Verify the branch** - Ensure you're on main (or correct branch)
4. **Write meaningful messages** - Future you will thank you
5. **Push after commit** - Unless explicitly told not to

### Never Do:
1. **Commit passwords/secrets** - Check for .env, keys, tokens
2. **Commit broken builds** - Unless explicitly WIP
3. **Mix unrelated changes** - Unless explicitly told to
4. **Commit .session-active** - It's local only
5. **Force push to main** - Ever

### When Uncertain:
1. **Show the status** - Let Adam see what's there
2. **Group the changes** - Organize by feature/purpose
3. **Ask for clarification** - Better safe than confused
4. **Suggest options** - Give 2-3 clear paths forward

---

## 📝 Response Format

After committing, always report:

```
✅ Committed and pushed: [commit hash short]

Summary: [one line what was done]
Files: [number] files changed
Commit: "[commit message first line]"
Status: Clean working tree / [or current state]
```

---

## 🔄 Special Cases

### "Ship it"
- Implies: Commit, push, and celebrate
- Add a bit of personality to the success message

### "Save that work"  
- Implies: Might be WIP, preservation important
- Use 'wip:' prefix if incomplete

### "Commit before we continue"
- Implies: Checkpoint commit, more work coming
- Be extra clear about current state

### "Let's checkpoint"
- Implies: Create a checkpoint, might revert later
- Use descriptive message about current state

---

## 🎯 Team Coordination

When multiple agents are working:
1. **Check for others' work**: `git fetch && git status`
2. **Pull before committing**: Avoid conflicts
3. **Mention agent in commit**: "CC: fix: [change]" or "AC: feat: [change]"
4. **Update NEXT.md**: If commit affects coordination

---

**Remember**: When Adam says "commit those changes," he trusts you to handle it intelligently. Use this workflow to make the right decision quickly.