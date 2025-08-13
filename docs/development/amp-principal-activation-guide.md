# Amp-Principal Activation Guide

Quick reference for activating the Amp-Principal replicant for technical review and architectural oversight.

## 🚀 Role Activation

### How to Invoke Amp-Principal

Use this prompt structure to activate Amp-Principal for reviews:

```
Switching to Amp-Principal role for technical review.

[Review Request Details]

Amp-Principal, please provide a comprehensive technical review of [specific item].
```

### Common Activation Scenarios

**1. PR Review Request:**
```
Switching to Amp-Principal role for technical review.

Please review PR #27 (7EP-0007 Phase 2 Search Engine) with focus on:
- Architecture quality and scalability
- Code patterns and maintainability  
- Performance implementation
- Documentation completeness

Amp-Principal, please provide comprehensive technical review.
```

**2. Architecture Design Review:**
```
Switching to Amp-Principal role for architectural review.

Please review the proposed architecture for 7EP-0007 Phase 3 Batch Operations:
- Multi-archive operation framework design
- Progress tracking implementation approach
- Error handling and rollback strategies
- Integration with existing MAS foundation

Amp-Principal, please evaluate architectural soundness and provide recommendations.
```

**3. Code Quality Assessment:**
```
Switching to Amp-Principal role for code quality review.

Please assess the current codebase quality in [specific area]:
- Adherence to Go best practices
- Error handling consistency
- Test coverage and quality
- Documentation standards compliance

Amp-Principal, please provide quality assessment with improvement recommendations.
```

## 📋 Review Request Templates

### Standard PR Review
```
## PR Review Request

**PR Details:**
- Number: #[number]
- Title: [title]
- Branch: [branch-name]
- Files Changed: [count] files, [lines] lines

**Review Scope:**
- [ ] Architecture & Design
- [ ] Code Quality & Patterns
- [ ] Performance Impact
- [ ] Documentation Completeness
- [ ] Test Coverage
- [ ] Security Considerations

**Specific Focus Areas:**
[List any specific concerns or areas requiring attention]

**Timeline:** [Required completion date if urgent]

Amp-Principal, please provide comprehensive technical review.
```

### Feature Architecture Review
```
## Architecture Review Request

**Feature:** [Feature name and 7EP reference]
**Scope:** [Brief description of architectural scope]

**Review Focus:**
- System design and component interactions
- Scalability and performance characteristics
- Integration with existing architecture
- Technical debt implications
- Security and reliability considerations

**Supporting Materials:**
- Design documents: [links]
- Previous architectural decisions: [context]
- Performance requirements: [specifications]

Amp-Principal, please evaluate architectural design and provide recommendations.
```

### Technical Standards Review
```
## Standards & Best Practices Review

**Area:** [Codebase area or technical domain]
**Current State:** [Brief assessment of current standards]

**Review Objectives:**
- Evaluate current patterns and practices
- Identify improvement opportunities
- Recommend standardization approaches
- Assess tool and process effectiveness

**Success Criteria:**
[Define what constitutes successful standards review]

Amp-Principal, please assess technical standards and recommend improvements.
```

## 🎯 Expected Review Outcomes

### Comprehensive Technical Feedback
Amp-Principal reviews will include:
- **Architectural Assessment**: System design evaluation with specific recommendations
- **Code Quality Analysis**: Pattern consistency, maintainability, and best practices compliance
- **Performance Evaluation**: Performance implications and optimization opportunities
- **Documentation Review**: Completeness and quality of technical documentation
- **Risk Assessment**: Technical risks and mitigation strategies
- **Action Items**: Prioritized list of required changes and improvements

### Review Deliverables Format
```
# Technical Review: [Item Name]

## 🏗️ Architecture Assessment
[Evaluation of system design, patterns, and architectural decisions]

## 📊 Code Quality Analysis  
[Code patterns, maintainability, and best practices evaluation]

## ⚡ Performance Evaluation
[Performance characteristics and optimization opportunities]

## 📚 Documentation Review
[Documentation completeness and quality assessment]

## 🚨 Risk Assessment
[Technical risks and recommended mitigation strategies]

## ✅ Approval Status
[Clear approval decision with any conditional requirements]

## 🎯 Required Actions
[Prioritized action items with specific recommendations]

## 💡 Future Recommendations
[Suggestions for long-term improvements and technical evolution]
```

## 🔄 Integration with Development Workflow

### When to Request Amp-Principal Review

**Required Reviews:**
- New feature implementations (7EP phases)
- Major architectural changes or refactoring
- Performance-critical code modifications
- Security-sensitive implementations
- Breaking changes to public APIs

**Recommended Reviews:**
- Complex bug fixes with architectural implications
- Third-party library integrations
- Database schema changes
- CI/CD pipeline modifications
- Documentation structure changes

**Optional Reviews:**
- Minor bug fixes (unless pattern-breaking)
- Documentation updates (unless architectural)
- Configuration changes
- Build script modifications

### Review Process Integration

1. **Pre-Implementation**: Architecture review for major features
2. **During Development**: Code review for complex implementations  
3. **Pre-Merge**: Final technical review before PR approval
4. **Post-Implementation**: Retrospective review for continuous improvement

## 📞 Communication Protocols

### Review Request Channels
- **Primary**: Direct prompt activation (as shown above)
- **Context Sharing**: Include relevant technical context and requirements
- **Timeline Communication**: Specify any urgency or deadline requirements
- **Follow-up Process**: Clear process for addressing feedback and re-review

### Response Expectations
- **Initial Acknowledgment**: Confirmation of review scope and timeline
- **Progress Updates**: Status updates for complex reviews requiring extended time
- **Final Deliverable**: Comprehensive review with clear action items
- **Follow-up Availability**: Availability for clarification and additional review cycles

## 🛠️ Best Practices for Effective Reviews

### Preparing for Review
1. **Complete Context**: Provide all relevant technical context and requirements
2. **Clear Scope**: Define specific areas requiring review attention
3. **Supporting Materials**: Include design docs, performance requirements, etc.
4. **Realistic Timeline**: Allow adequate time for thorough review

### During Review Process
1. **Responsive Communication**: Promptly address questions and provide clarification
2. **Open Collaboration**: Engage in technical discussion and alternative approaches
3. **Change Documentation**: Track and document all review-driven changes
4. **Quality Focus**: Prioritize thoroughness over speed for critical reviews

### Post-Review Actions
1. **Action Item Tracking**: Systematically address all required changes
2. **Re-review Process**: Request follow-up review after addressing feedback
3. **Documentation Updates**: Update documentation based on review recommendations
4. **Process Improvement**: Incorporate lessons learned into future development

## 📈 Continuous Improvement

### Review Process Evolution
- **Feedback Collection**: Regular assessment of review effectiveness
- **Process Refinement**: Continuous improvement of review standards and workflows
- **Tool Integration**: Leverage tools to improve review efficiency and quality
- **Knowledge Sharing**: Share insights and best practices across the team

### Standards Evolution
- **Technology Updates**: Adapt standards for new technologies and best practices
- **Project Growth**: Evolve standards as project complexity and scale increase
- **Team Learning**: Incorporate team learning and experience into standards
- **Industry Alignment**: Align with evolving industry best practices and patterns

---

**Activation Status**: 🟢 **Ready for Use**  
**Role Documentation**: [amp-principal-role-definition.md](./amp-principal-role-definition.md)  
**Last Updated**: 2025-08-13