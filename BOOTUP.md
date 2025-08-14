# AI Team Member Boot-up Guide

**Purpose**: Standardized startup sequence for all AI team members  
**Framework**: Document Driven Development (7EP-0017)  
**Team Context**: See `/docs/development/TEAM-CONTEXT.md` for shared project overview

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