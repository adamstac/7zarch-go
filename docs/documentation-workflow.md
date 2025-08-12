# Documentation Workflow Process

This document defines how we transition from development planning (7EPs) to comprehensive user documentation that welcomes power users.

## Documentation Philosophy

**Goal**: Welcome power users with documentation that explains the whats, whys, and hows with clear examples.

**Principles**:
- **What**: Clear description of functionality and capabilities
- **Why**: Context and reasoning behind design decisions  
- **How**: Step-by-step examples with real-world scenarios
- **When**: Appropriate use cases and timing guidance

## Documentation Lifecycle

### Phase 1: Development Planning (7EP)
**Input**: Feature requirements, user feedback  
**Output**: Structured 7EP with technical design  
**Artifacts**: `docs/7eps/7ep-XXXX-feature-name.md`

### Phase 2: Implementation Learnings
**Input**: Implementation experience, edge cases discovered  
**Output**: Updated 7EP with implementation notes  
**Artifacts**: Implementation notes added to existing 7EP

### Phase 3: User Documentation Creation
**Input**: Completed 7EP + implementation learnings  
**Output**: Comprehensive user documentation  
**Artifacts**: Updated README sections, new user guides

### Phase 4: Documentation Maintenance  
**Input**: User feedback, support questions  
**Output**: Refined documentation with better examples  
**Artifacts**: Iterative improvements to user docs

## Documentation Types & Responsibilities

### 7EPs (Development Documentation)
**Audience**: Developers, maintainers  
**Content**: Technical design, implementation details, alternatives considered  
**Lifecycle**: Permanent archive for historical reference  
**Owner**: Assigned developer (CC/AC)

### User Documentation (Public Documentation)
**Audience**: End users, power users, new adopters  
**Content**: Usage patterns, examples, best practices  
**Lifecycle**: Living documents updated with feature changes  
**Owner**: Feature implementer

### Reference Documentation (API/Command Reference)
**Audience**: Daily users, automation builders  
**Content**: Complete command syntax, all flags, return codes  
**Lifecycle**: Generated from code where possible  
**Owner**: Automated + manual review

## Transition Process: 7EP → User Docs

### When to Create User Documentation
- ✅ Feature implementation is complete and tested
- ✅ Edge cases and error conditions are well-understood  
- ✅ User testing has validated the approach
- ✅ Command interface is stable

### What to Extract from 7EPs

#### From "Use Cases" Section
**Extract**: Real-world scenarios  
**Transform to**: Step-by-step tutorials with context

#### From "Technical Design" Section  
**Extract**: User-facing behavior and concepts  
**Transform to**: Conceptual explanations and mental models

#### From "Implementation Plan" Section
**Extract**: Feature capabilities and limitations  
**Transform to**: Usage guidance and best practices

#### From "Testing Strategy" Section
**Extract**: Acceptance criteria and test scenarios  
**Transform to**: Examples and edge case documentation

### Documentation Enhancement Process

#### 1. Implementation Completion
```bash
# When feature is complete, update 7EP with learnings
# Add "Implementation Notes" section to 7EP
# Document unexpected challenges and solutions
```

#### 2. User Documentation Creation
```bash
# Create or update relevant README sections
# Add comprehensive examples with context
# Include troubleshooting for common issues
# Link back to 7EP for developers who want technical details
```

#### 3. Cross-Reference Linking
```bash
# 7EP references user documentation location
# User docs reference 7EP for technical background
# Maintain bidirectional traceability
```

## Documentation Standards

### User Documentation Requirements

#### Command Documentation Format
```markdown
### command-name

Brief description of what this command does and why you'd use it.

#### Basic Usage
`7zarch-go command-name [flags] <arguments>`

#### Common Scenarios

**Scenario 1: [Real-world use case]**
Context: When you need to...
```bash
7zarch-go command-name --flag value
```
Expected output: [Show actual output]
Next steps: [What to do after]

**Scenario 2: [Another use case]**
[Repeat pattern]

#### Advanced Usage
[Power user features, combining with other commands]

#### Troubleshooting
Common error: [Error message]
Solution: [How to fix]
```

#### Conceptual Documentation Format
```markdown
## Concept Name

### What It Is
[Clear definition without jargon]

### Why It Exists  
[The problem it solves, user benefit]

### How It Works
[Mental model, key behaviors]

### When to Use It
[Appropriate scenarios, alternatives]

### Examples
[Progressive examples from simple to complex]
```

## Documentation Quality Checklist

### Before Publishing User Documentation
- [ ] Tested all examples on clean environment
- [ ] Covers common error conditions and solutions
- [ ] Explains the "why" behind commands, not just "how"
- [ ] Includes context for when to use vs. alternatives
- [ ] Links to related commands and concepts
- [ ] Uses consistent terminology throughout
- [ ] Avoids implementation details users don't need

### Documentation Review Process
1. **Technical Review**: Accuracy and completeness
2. **User Experience Review**: Clarity and discoverability  
3. **Example Validation**: All code examples work as shown
4. **Link Verification**: All internal and external links function

## Example: Trash Management Documentation Plan

### From 7EP-0001 to User Documentation

**Current 7EP Content**: Technical design for restore, trash list, purge commands  
**User Documentation Needed**:

1. **README Section**: "Archive Recovery and Management"
   - Why: Explain accident recovery and cleanup workflows
   - What: Overview of trash system behavior
   - How: Step-by-step recovery scenarios

2. **Command Reference Updates**:
   - `restore` command with resolver examples
   - `trash list` with filtering options  
   - `trash purge` with safety considerations

3. **Workflow Documentation**:
   - "Recovering Accidentally Deleted Archives"
   - "Managing Archive Cleanup"
   - "Understanding Auto-Purge Policies"

### Implementation Timeline
- **Phase 1**: Complete trash management implementation
- **Phase 2**: Extract learnings and update 7EP-0001
- **Phase 3**: Create comprehensive user documentation
- **Phase 4**: Validate with user testing

## Future Documentation Enhancements

### Power User Documentation Ideas
- **Configuration Deep Dive**: Advanced config patterns and optimization
- **Automation Cookbook**: Scripts and workflows for common tasks
- **Integration Guide**: Using 7zarch-go in backup pipelines and CI/CD
- **Troubleshooting Playbook**: Comprehensive error resolution guide
- **Performance Tuning**: Optimization strategies for large-scale usage

### Documentation Automation Opportunities
- **Command Help Generation**: Extract from cobra command definitions
- **Example Testing**: Automated validation of documentation examples
- **Usage Analytics**: Track which documentation sections are most accessed
- **Feedback Integration**: Easy ways for users to suggest improvements

---

This documentation workflow ensures that every implemented feature results in high-quality user documentation that welcomes power users and provides the depth they need to be successful.