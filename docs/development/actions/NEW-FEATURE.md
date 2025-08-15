# New Feature Workflow: "Start New Feature", "Begin 7EP Work"

**Purpose**: Standardized process for starting new feature development  
**Framework**: Document Driven Development (7EP-0017)  
**Trigger Phrases**: "Start new feature", "Begin 7EP-XXXX", "New feature branch", "Clean slate"

## 🎯 Quick Decision Tree

```
"Start new feature"
         |
    [Check Context]
         |
    What kind of work?
    /       |        \
  7EP    Feature    Experiment
   |       |          |
[7EP Flow] [Feature Flow] [Experiment Flow]
```

## 📋 Feature Start Scenarios

### Scenario A: 7EP Implementation 🎯
**Context**: Starting work on approved 7EP

**Pre-work Check**:
```bash
# Verify 7EP exists and status
ls docs/7eps/ | grep 7EP-XXXX
cat docs/7eps/7EP-XXXX-*.md | head -20

# Check if already assigned
grep -l "7EP-XXXX" docs/development/roles/*.md
```

**Action**:
```bash
# Ensure clean starting point
git checkout main && git pull origin main
git status  # Should be clean

# Create 7EP feature branch
7EP_NUM="XXXX"  # e.g., "0019"
7EP_DESC="short-description"  # e.g., "tui-enhancement"
git checkout -b "feat/7ep-${7EP_NUM}-${7EP_DESC}"

# Update role assignment
echo "## 🎯 Current Assignment: 7EP-${7EP_NUM}" >> docs/development/roles/CLAUDE.md
echo "**Status**: ACTIVE - Implementation in progress" >> docs/development/roles/CLAUDE.md
echo "**Branch**: feat/7ep-${7EP_NUM}-${7EP_DESC}" >> docs/development/roles/CLAUDE.md

# Commit role assignment
git add docs/development/roles/CLAUDE.md
git commit -m "assign: start 7EP-${7EP_NUM} implementation

Agent: CC (Claude Code)
Branch: feat/7ep-${7EP_NUM}-${7EP_DESC}
Status: Beginning implementation"

# Update coordination
echo "- **CC**: Started 7EP-${7EP_NUM} implementation on $(date +%Y-%m-%d)" >> docs/development/NEXT.md
```

**Response**: "🎯 Started 7EP-${7EP_NUM} on branch feat/7ep-${7EP_NUM}-${7EP_DESC}"

---

### Scenario B: General Feature 🔧
**Context**: New feature not tied to specific 7EP

**Action**:
```bash
# Clean starting point
git checkout main && git pull origin main

# Create descriptive feature branch
FEATURE_NAME="descriptive-name"  # e.g., "archive-validation"
git checkout -b "feat/${FEATURE_NAME}"

# Create basic feature structure
mkdir -p docs/features/
cat > docs/features/${FEATURE_NAME}.md << EOF
# Feature: ${FEATURE_NAME}

## Overview
[Brief description of what this feature does]

## Implementation Plan
- [ ] [Step 1]
- [ ] [Step 2] 
- [ ] [Step 3]

## Acceptance Criteria
- [ ] [Criteria 1]
- [ ] [Criteria 2]

## Testing Plan
- [ ] Unit tests for [component]
- [ ] Integration tests for [workflow]
- [ ] Manual testing scenarios
EOF

# Initial commit
git add docs/features/${FEATURE_NAME}.md
git commit -m "feat: start ${FEATURE_NAME} development

- Added feature specification
- Created implementation roadmap
- Ready for development"
```

**Response**: "🔧 Started feature '${FEATURE_NAME}' with planning document"

---

### Scenario C: Experiment/Spike 🧪
**Context**: Exploratory work or proof of concept

**Action**:
```bash
# Create experiment branch
EXPERIMENT_NAME="spike-description"  # e.g., "spike-search-optimization"
git checkout -b "experiment/${EXPERIMENT_NAME}"

# Create experiment log
mkdir -p docs/experiments/
cat > docs/experiments/${EXPERIMENT_NAME}.md << EOF
# Experiment: ${EXPERIMENT_NAME}

## Hypothesis
[What you're trying to prove or explore]

## Goals
- [ ] [Explore approach A]
- [ ] [Test performance of B]
- [ ] [Validate assumption C]

## Success Criteria
[What would make this experiment successful]

## Timeline
**Started**: $(date +%Y-%m-%d)
**Expected Duration**: [timeframe]

## Log
### $(date +%Y-%m-%d)
- Started experiment
- [Initial observations]

---
*Note: This is experimental work - may be discarded*
EOF

# Commit experiment setup
git add docs/experiments/${EXPERIMENT_NAME}.md
git commit -m "experiment: start ${EXPERIMENT_NAME}

Exploring: [brief description]
Timeline: [expected duration]
May be discarded after evaluation"
```

**Response**: "🧪 Started experiment '${EXPERIMENT_NAME}' - exploratory work"

---

### Scenario D: Converting WIP to Feature 🔄  
**Context**: Taking messy work and making it into a proper feature

**Action**:
```bash
# Assess current state
git status
git log --oneline -10
git stash  # Save any uncommitted work

# Create clean feature branch
FEATURE_NAME="cleaned-feature-name"
git checkout main && git pull origin main
git checkout -b "feat/${FEATURE_NAME}"

# Apply valuable work from stash/commits
git stash pop  # If there was stashed work
# OR cherry-pick specific commits:
# git cherry-pick [commit-hash]

# Clean up and organize
[organize code, clean commits, add tests]

# Document the cleanup
cat > CLEANUP-NOTES.md << EOF
# Feature Cleanup: ${FEATURE_NAME}

## Source Work
- Previous branch: [old-branch-name]
- Commits preserved: [list]
- Work discarded: [what was removed and why]

## Improvements Made
- [Cleanup action 1]
- [Cleanup action 2]
- [Added missing tests/docs]

## Ready For
- [What this feature now accomplishes]
- [How it integrates with existing code]
EOF

git add -A
git commit -m "feat: ${FEATURE_NAME} from WIP cleanup

- Organized previous exploratory work
- Added proper tests and documentation
- Ready for review and integration"
```

**Response**: "🔄 Converted WIP to clean feature: ${FEATURE_NAME}"

---

## 🎨 Feature Planning Templates

### 7EP Feature Plan
```markdown
# 7EP-XXXX Implementation Plan

## 7EP Context
**Title**: [7EP full title]
**Phase**: [X] of [Y] phases  
**Dependencies**: [other 7EPs or work]

## Implementation Approach
### Phase 1: Foundation
- [ ] [Core component 1]
- [ ] [Core component 2]
- [ ] [Basic integration]

### Phase 2: Features  
- [ ] [Feature 1]
- [ ] [Feature 2]
- [ ] [User interface]

### Phase 3: Polish
- [ ] [Performance optimization]
- [ ] [Error handling]
- [ ] [Documentation]

## Technical Decisions
- **Architecture**: [approach chosen]
- **Dependencies**: [new packages needed]
- **Testing Strategy**: [how to test]
- **Integration Points**: [where it connects]

## Success Metrics
- [ ] [Measurable outcome 1]
- [ ] [Measurable outcome 2]
- [ ] [User experience goal]
```

### General Feature Plan
```markdown
# Feature: [Feature Name]

## Problem Statement
[What problem does this solve]

## Proposed Solution
[How this feature addresses the problem]

## Implementation Plan
- [ ] [Development step 1]
- [ ] [Development step 2]
- [ ] [Integration step]
- [ ] [Testing step]

## Acceptance Criteria
- [ ] [Specific testable requirement]
- [ ] [User-facing behavior]
- [ ] [Performance requirement]

## Out of Scope
[What this feature explicitly does NOT do]
```

## 🚀 Best Practices

### Branch Hygiene
- **One feature per branch** - Keep branches focused
- **Branch from main** - Always start from latest main
- **Regular rebasing** - Keep history clean with main
- **Descriptive names** - Make branch purpose obvious

### Commit Practices
- **Atomic commits** - One logical change per commit
- **Clear messages** - Explain what and why
- **Test each commit** - Ensure each commit builds/works
- **No WIP commits** - In final feature branches

### Documentation
- **Plan before coding** - Write the plan first
- **Update as you go** - Keep docs current with implementation
- **Include examples** - Show how to use the feature
- **Testing instructions** - How others can verify it works

## 🎯 Quality Gates

### Before Starting Development
- [ ] Clear understanding of requirements
- [ ] Implementation approach decided
- [ ] Dependencies identified
- [ ] Timeline estimated

### During Development
- [ ] Regular commits with clear messages
- [ ] Tests written alongside code
- [ ] Documentation updated
- [ ] Integration points tested

### Before Submitting PR
- [ ] All tests pass
- [ ] Documentation complete
- [ ] Code review ready
- [ ] Integration verified

---

## 🚨 Common Pitfalls

### Starting Without Planning
**Problem**: Jumping into code without clear direction  
**Solution**: Always create planning document first

### Feature Creep
**Problem**: Adding unrelated functionality  
**Solution**: Keep strict scope, create new features for extras

### Broken History
**Problem**: Messy commits make review difficult  
**Solution**: Use interactive rebase to clean up before PR

### Missing Context
**Problem**: Others can't understand the feature purpose  
**Solution**: Clear documentation and commit messages

---

**Remember**: Starting well is half the battle. Take time to set up properly, and the development process will be much smoother.