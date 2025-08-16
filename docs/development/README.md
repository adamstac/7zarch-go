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
4. **Quality validation** - Use `make validate-framework` for comprehensive compliance
5. **Framework monitoring** - Use `make framework-health` for continuous health assessment
6. **CI enforcement** - GitHub Actions automatically validate framework compliance on PRs

## 🛠️ **Framework Validation & Maintenance**

### Daily Validation Commands
```bash
# Quick role file validation
make validate-framework-roles

# Check cross-document consistency
make validate-framework-consistency  

# Complete framework validation suite
make validate-framework

# Framework health dashboard
make framework-health

# Complete integration testing
make validate-framework-integration
```

### Framework Health Monitoring
- **Structure Validation**: All documents follow standardized patterns
- **Consistency Checking**: Role files ↔ NEXT.md coordination synchronized
- **Integration Testing**: BOOTUP/SHUTDOWN/TEAM-UPDATE workflows operational
- **Git Pattern Compliance**: Session logs, coordination commits, branch naming
- **CI Automation**: GitHub Actions enforce compliance on all PRs
- **Auto-fix Capabilities**: Missing headers and common issues resolved automatically

### New Agent Onboarding Process
1. **Copy template**: Use `docs/development/roles/ROLE-TEMPLATE.md`
2. **Customize role**: Adapt sections for agent expertise and focus areas
3. **Set initial assignments**: Define first priorities and coordination needs
4. **Test lifecycle**: Run `BOOTUP.md` → work simulation → `SHUTDOWN.md`
5. **Validate compliance**: Run `make validate-framework-roles`
6. **Integration validation**: Confirm role works with team coordination patterns

### Framework Troubleshooting
**Issue**: Role file validation fails  
**Solution**: Run `make validate-framework-roles` for specific error details

**Issue**: Cross-document inconsistency detected  
**Solution**: Run `make validate-framework-consistency` and sync role files with NEXT.md

**Issue**: Workflow scripts fail during lifecycle testing  
**Solution**: Run `make validate-framework-integration` for detailed diagnostic

**Issue**: Framework health declining  
**Solution**: Run `make framework-health` for metrics and recommendations

**Issue**: CI blocking PRs due to framework violations  
**Solution**: Auto-fix will resolve common issues, manual fixes for complex violations

## 📂 **Archive Location**

Historical development documentation moved to:
- **docs/archive/development/** - Tomorrow plans, old roadmaps, completed guides

## 🔄 **Framework Status**

**Implementation:** 
- 7EP-0017 Document Driven Development Framework ✅ Complete
- 7EP-0019 Agent Lifecycle & Coordination Standardization ✅ Complete  
- 7EP-0020 DDD Framework Validation & Compliance Suite ✅ Complete

**Status:** Production Ready - Complete validation and CI enforcement active  
**Quality:** 100% framework health, comprehensive validation coverage, CI automated  
**Validation:** `make validate-framework` - Complete validation suite operational  
**Health:** `make framework-health` - Continuous monitoring dashboard
