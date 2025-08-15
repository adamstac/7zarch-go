# Development Documentation Overview

**Document Driven Development (DDD) Framework**

## 📋 **Current Coordination Files**

### Team Assignment Docs (Current Work)
- **AMP.md** - Amp current assignments and architectural work
- **CLAUDE.md** - CC current assignments and implementation focus  
- **AUGMENT.md** - AC current assignments and availability
- **ADAM.md** - Project owner strategic priorities and decisions needed

### Shared Coordination
- **NEXT.md** - Shared priority hub - what's next for everyone

### Process Documentation (Keep)
- **emoji-usage-guidelines.md** - Standards for monospace-safe UI design
- **migration-best-practices.md** - Technical reference for database migrations

### Legacy Coordination (Update as needed)
- **pr-merge-roadmap.md** - Transitioned to framework model
- **sprint-planning-analysis.md** - Strategic analysis (reference)

## 🎯 **How to Use This System**

### Agent Lifecycle Framework (7EP-0019)
Complete operational framework for AI team members from session startup to shutdown:

1. **Session Startup** - Use [`actions/BOOTUP.md`](actions/BOOTUP.md)
   - Git sync and build verification
   - Load role context and assignments
   - Validate coordination needs
   - Ready state confirmation

2. **Daily Operations** - Role-driven work execution
   - Follow assignments in your role file (`roles/[AGENT].md`)
   - Use workflow actions (`COMMIT.md`, `MERGE.md`, `NEW-FEATURE.md`)
   - Update coordination via [`actions/TEAM-UPDATE.md`](actions/TEAM-UPDATE.md)
   - Maintain real-time status in role file and `NEXT.md`

3. **Session Shutdown** - Use [`actions/SHUTDOWN.md`](actions/SHUTDOWN.md)
   - Commit work with clear status
   - Update role file with session completion
   - Document next session priorities
   - Preserve coordination context

### For Team Members
1. **Start each session** with `actions/BOOTUP.md` sequence
2. **Check your role file** - `docs/development/roles/[AGENT].md` for current assignments
3. **Check shared priorities** - `docs/development/NEXT.md` for team coordination
4. **Update coordination** using `actions/TEAM-UPDATE.md` when status changes
5. **End each session** with `actions/SHUTDOWN.md` sequence

### For Project Owner (Adam)
1. **Strategic decisions** - Use [`STRATEGIC-DECISION-FRAMEWORK.md`](STRATEGIC-DECISION-FRAMEWORK.md)
2. **Team assignment** - Update agent role files with new priorities  
3. **Coordination oversight** - Monitor `NEXT.md` for team-wide blockers and progress

### For Cross-Team Coordination
1. **Role integration** - All agents follow standardized lifecycle patterns
2. **Real-time status** - Role files and `NEXT.md` maintain current coordination state
3. **Systematic handoffs** - BOOTUP/SHUTDOWN preserve context across sessions
4. **Quality validation** - Use `scripts/validate-roles.sh` to ensure compliance

## 📂 **Archive Location**

Historical development documentation moved to:
- **docs/archive/development/** - Tomorrow plans, old roadmaps, completed guides

## 🔄 **Framework Status**

**Implementation:** 
- 7EP-0017 Document Driven Development Framework ✅ Complete
- 7EP-0019 Agent Lifecycle & Coordination Standardization ✅ Complete

**Status:** Operational - Full agent lifecycle framework active  
**Quality:** All role files compliant, validation linter operational  
**Next:** Team adoption validation and framework refinement based on usage
