# AI Team Member Boot-up Guide

**Purpose**: Standardized startup sequence for all AI team members  
**Framework**: Document Driven Development (7EP-0017)  
**Team Context**: See `/docs/development/TEAM-CONTEXT.md` for shared project overview

## 🎯 Project Vision

**7zarch-go**: Advanced archive management for everyone who cares about their data. It's 7-Zip with a brain.

We're building the tool that makes advanced compression accessible to everyone:
- **Smart compression** - Automatically optimizes based on content type (media, code, documents)
- **Instant search** - Find any file across all archives in microseconds
- **Built-in verification** - Never worry about corrupt archives
- **Power when needed** - Simple defaults, advanced features available

## 🚀 Quick Start Sequence

Execute these commands in order to get project context:

### 1. Check Git Status
```bash
git status && git pull && git branch
```

### 2. Review Current State
```bash
# Check for active PRs and recent commits
gh pr list && git log --oneline -10
```

### 3. **CHECK OPERATIONAL PRIORITIES** (DDD Framework)
```bash
# Personal assignments (replace [ROLE] with: CLAUDE, AMP, AUGMENT, ADAM)
cat docs/development/[ROLE].md | head -20

# Shared team priorities and blockers  
cat docs/development/NEXT.md | head -30

# Active 7EP coordination context
grep -l "Status.*ACTIVE\|In Progress" docs/7eps/*.md | xargs ls -la
```

### 4. Test Build
```bash
make dev && ./7zarch-go list --dashboard
```

### 5. Log Session Start
```bash
# Create logs directory if needed
mkdir -p docs/logs

# Initialize session log with start timestamp
DATE_STAMP=$(date +%Y-%m-%d_%H-%M-%S)
SESSION_START=$(date "+%Y-%m-%d %H:%M:%S")
cat > docs/logs/session-${DATE_STAMP}.md << EOF
# Session Log - $(date)

## ⏱️ Session Timing
- **Start Time:** ${SESSION_START}
- **Agent:** CC (Claude Code)
- **Status:** 🟢 **ACTIVE** - Session in progress

## 🚀 Boot Sequence Completed
- Git status: Clean and up to date
- Build verification: Successful
- Operational priorities: Reviewed
- Ready for work assignment

---
*Session started by DDD Framework bootup process*
EOF

# Store session file for shutdown reference (local only)
echo "SESSION_LOG_FILE=docs/logs/session-${DATE_STAMP}.md" > .session-active

# Add initial log to git (NOT .session-active - it's local only)
git add docs/logs/session-${DATE_STAMP}.md
git commit -m "session: start new session $(date)"
```

## 📋 Information Priority Order

**Daily Operations (Check First)**:
1. `docs/development/NEXT.md` - What's happening now across all teams
2. `docs/development/[YOUR-ROLE].md` - Your current assignments and priorities
3. Active 7EPs - Sprint-level coordination requirements

**Reference Information (As Needed)**:
4. `/AGENT.md` - Technical build/test patterns and code style
5. `/docs/development/TEAM-CONTEXT.md` - Project structure and team overview
6. `docs/7eps/index.md` - Long-term feature planning

## 🎯 Current Project State

**Foundation Phase**: ✅ Complete - Production-ready CLI with TUI  
**Advanced Features Phase**: ✅ Complete - Full query/search/batch operations  
**Current Phase**: 🎯 Strategic Direction Planning

### Team Status
- **CC (Claude Code)**: Available for strategic assignment
- **AC (Augment Code)**: Available for strategic assignment  
- **Amp**: Strategic planning coordination and framework oversight
- **Adam**: Strategic priority decision needed

### Key Context
- **7EP-0007 Complete** - Enterprise query/search/batch transformation delivered
- **DDD Framework Operational** - All team coordination documents active and effective
- **Strategic Decision Pending** - Next focus: TUI evolution vs TrueNAS integration vs optimization
- **All Major Work Blocked** - Waiting for Adam's strategic direction

## 🔄 Role-Specific Boot-up

### For Claude Code (CC)
```bash
# Check CC-specific assignments
cat docs/development/CLAUDE.md

# Focus areas: Backend features, infrastructure, technical implementation
```

### For Augment Code (AC)  
```bash
# Check AC-specific assignments
cat docs/development/AUGMENT.md

# Focus areas: User experience, frontend features, quality assurance
```

### For Amp (Strategic/Technical Leadership)
```bash  
# Check Amp assignments
cat docs/development/AMP.md

# Focus areas: Strategic planning, architectural oversight, cross-team coordination
```

### For Adam
```bash
# Check team coordination and decisions needed
cat docs/development/NEXT.md

# Review pending strategic decisions
grep -B2 -A5 "Strategic\|Decision\|Waiting.*Adam" docs/development/NEXT.md

# Check for items awaiting approval
grep -B2 -A5 "Awaiting.*decision\|Approval" docs/development/*.md docs/7eps/*.md
```

**Key Decision Points for Adam**:
1. **Review NEXT.md** - What decisions are blocking the team?
2. **Check Draft 7EPs** - Any new proposals awaiting approval?
3. **Strategic Direction** - Set next major focus area when ready
4. **Resource Allocation** - Assign AI agents to specific work streams

## 📊 Success Checklist

After boot-up, you should clearly understand:
- [ ] **Current work status** - What you're actively working on
- [ ] **Next priorities** - What comes next in priority order
- [ ] **Coordination needs** - Who you're waiting on or coordinating with
- [ ] **Strategic context** - How your work fits into project direction
- [ ] **Immediate blockers** - What's preventing progress

## 🚨 Common Issues & Solutions

**Issue**: Documents reference outdated or completed work  
**Solution**: Update your role document (`docs/development/[ROLE].md`) with current status

**Issue**: Unclear coordination dependencies  
**Solution**: Check `NEXT.md` for cross-team coordination points

**Issue**: No clear next priorities  
**Solution**: Strategic decision needed from Adam - surface in coordination docs

**Issue**: Build or dependencies broken  
**Solution**: Check recent commits (`git log`), review any failing CI/tests

## 🎯 Ready State Confirmation

**When you're ready to work, you should be able to answer:**
1. What is my current active work?
2. What are my next 2-3 priorities?
3. Who am I coordinating with or waiting on?
4. What strategic context affects my work?
5. Are there any immediate blockers?

**Report Format**: Update your role document and coordinate via NEXT.md as needed.

---

**Framework Integration**: This boot-up sequence validates and uses the DDD operational framework for consistent team coordination.

**Next Steps**: Execute boot-up → Update role document → Begin assigned work → Coordinate via NEXT.md