# Contributing to 7zarch-go

Thanks for your interest in contributing to 7zarch-go! This document provides guidelines for contributing to our systematic document-driven development approach.

## 🎯 Development Framework

**7zarch-go uses Document Driven Development (DDD)** - a systematic coordination framework that ensures quality and team scalability. All contributions must follow DDD framework compliance.

### **Framework Validation Required**
Before submitting any PR, run our validation suite:
```bash
# Complete framework validation
make validate-framework

# Quick role file check (if contributing to coordination docs)
make validate-framework-roles

# Framework health assessment
make framework-health
```

## 📋 Contribution Types

### **Code Contributions**
- Follow existing patterns in `AGENT.md` for technical conventions
- Run tests: `go test ./...` and `golangci-lint run`
- Ensure `make build` succeeds on your platform
- Add tests for new functionality

### **Documentation Contributions**
- **Role files**: Must follow `docs/development/roles/ROLE-TEMPLATE.md` structure
- **7EPs**: Use `docs/7eps/template.md` for new enhancement proposals
- **Framework docs**: Validate with `make validate-framework-consistency`
- **Session logs**: Follow `docs/development/actions/BOOTUP.md` and `SHUTDOWN.md` patterns

### **Framework Contributions**
- **Validation tools**: Extend `scripts/validate-*.go` for new validation needs
- **Workflow actions**: Follow patterns in `docs/development/actions/`
- **CI improvements**: Update `.github/workflows/ddd-validation.yml`
- **Role patterns**: Preserve unique value while maintaining standardization

## 🚀 Getting Started

### **Setup Development Environment**
```bash
# Clone and build
git clone https://github.com/adamstac/7zarch-go.git
cd 7zarch-go
make dev

# Validate framework compliance  
make validate-framework

# Check framework health
make framework-health
```

### **Before Making Changes**
1. **Read framework docs**: `docs/development/README.md`
2. **Check current coordination**: `docs/development/NEXT.md`
3. **Understand role patterns**: `docs/development/roles/ROLE-TEMPLATE.md`
4. **Run validation baseline**: `make validate-framework`

### **During Development**
- **Follow branching strategy**: Use `feature/description` or `feat/7ep-XXXX-name` branches
- **Document changes**: Update relevant role files and coordination docs
- **Test continuously**: `make validate-framework` catches issues early
- **Use coordination patterns**: Follow `docs/development/actions/TEAM-UPDATE.md`

### **Before Submitting PR**
```bash
# Complete validation (required)
make validate-framework

# Integration testing  
make validate-framework-integration

# Build verification
make build && go test ./...

# Framework health check
make framework-health
```

## 📏 Quality Standards

### **Framework Compliance (Required)**
- **Role files**: 100% compliance with 7EP-0019 standards
- **Cross-document consistency**: No coordination mismatches
- **Git patterns**: Follow session log and coordination commit patterns
- **Workflow integration**: BOOTUP/SHUTDOWN/TEAM-UPDATE patterns operational

### **Code Quality (Required)**  
- **Tests pass**: `go test ./...` 
- **Linting clean**: `golangci-lint run`
- **Build succeeds**: `make build` on your platform
- **No security violations**: `gosec` clean (checked by CI)

### **Documentation Quality (Recommended)**
- **Clear structure**: Follow existing document patterns
- **Content boundaries**: No duplication across document types
- **Team context**: Reference `docs/development/TEAM-CONTEXT.md` vs embedding
- **Strategic frameworks**: Use `docs/development/STRATEGIC-DECISION-FRAMEWORK.md`

## 🤝 Team Coordination

### **Multi-Agent Coordination**
7zarch-go uses AI agents for development. Contributors should:
- **Follow coordination patterns**: Update role files and `NEXT.md` for assignment changes
- **Use framework workflows**: BOOTUP → Work → SHUTDOWN patterns for session management
- **Respect assignment boundaries**: Check agent assignments before taking on coordination work
- **Communicate systematically**: Use `TEAM-UPDATE.md` patterns for cross-agent work

### **Framework Adoption**
- **New contributors**: Follow new agent onboarding process (<30 minutes to productivity)
- **Framework evolution**: Contribute validation improvements and coordination pattern enhancements
- **Team scaling**: Framework patterns support growth without coordination overhead increase

## 🚨 Common Issues & Solutions

### **Framework Validation Failures**
```bash
# Issue: Role file validation fails
# Solution: Check against ROLE-TEMPLATE.md structure
make validate-framework-roles

# Issue: Cross-document inconsistency
# Solution: Sync role files with NEXT.md coordination
make validate-framework-consistency

# Issue: Workflow integration broken
# Solution: Test complete lifecycle patterns
make validate-framework-integration
```

### **Development Environment Issues**
```bash
# Issue: Build fails
# Solution: Clean dependencies and rebuild
go mod tidy && make clean && make build

# Issue: Tests fail
# Solution: Verify test data and run specific packages
go test -v ./internal/storage/

# Issue: Framework health declining
# Solution: Run diagnostics and follow recommendations
make framework-health
```

## 📊 Success Metrics

### **Contribution Quality Indicators**
- **Framework compliance**: 100% validation pass rate
- **Integration success**: Complete lifecycle testing passes  
- **Team coordination**: Clear assignment updates and handoffs
- **Documentation quality**: Framework boundaries respected

### **Team Impact Measures**
- **Coordination efficiency**: Framework patterns reduce communication overhead
- **Knowledge preservation**: Context maintained across sessions and handoffs
- **Framework adoption**: Contributors prefer DDD patterns over ad-hoc coordination
- **Team scaling**: Framework supports new contributor integration in <30 minutes

## 🎯 Framework Philosophy

**Document Driven Development**: All coordination, assignments, and strategic decisions are captured in git-tracked documents with systematic validation.

**Agent Lifecycle Integration**: Contributors follow standardized patterns from session startup through work execution to session shutdown.

**Systematic Coordination**: Cross-team dependencies and handoffs use validated patterns that scale from individual work to complex multi-contributor projects.

**Quality Through Automation**: Framework validation prevents coordination degradation and ensures consistent team effectiveness.

---

**Questions?** Check `docs/development/README.md` for framework usage or create an issue for contribution guidance.

**Framework Health**: Run `make framework-health` to see current coordination patterns and framework effectiveness metrics.
