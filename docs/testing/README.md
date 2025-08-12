# Testing Documentation

This directory contains testing documentation, reports, and validation results for the 7zarch-go project.

## Directory Structure

### `/smoke-test-reports/`
Ongoing test execution reports from team members. These are living documents that get updated with each test run to track system evolution and communicate findings between AC and CC.

**Current Reports:**
- `e2e-smoke-test.md` - End-to-end workflow testing (create → delete → restore)
- `performance-benchmarks.md` - System performance testing and benchmarks
- `integration-testing.md` - Cross-component integration validation

### Formal Test Documentation
- `mas-foundation-integration.md` - 7EP-0004 MAS foundation integration tests
- Other formal test specifications and results

## Report Conventions

### Naming
- **Descriptive names**: Focus on what's tested, not when
- **Iterative updates**: Replace content but preserve history in collapsible sections
- **Consistent format**: Use the provided template for new reports

### Template for New Reports
```markdown
# {Test Type} Report

**Last Updated**: YYYY-MM-DD by {Tester}  
**Test Scope**: {Brief description}  
**Status**: ✅ PASS | 🔴 FAIL | ⚠️ PARTIAL  

## Latest Results ({Date})
[Current findings]

## Test Details
[Step-by-step results]

## Issues Found
[Problems discovered]

## Recommendations
[Suggested improvements]

## Previous Results
<details>
<summary>{Previous Date} Results</summary>
[Historical content]
</details>
```

## Communication Flow

1. **Test Execution**: AC/CC run tests and document findings
2. **Report Update**: Update relevant report in `/smoke-test-reports/`
3. **Cross-Reference**: Link findings in roadmap and 7EP documents
4. **Team Sync**: Reports serve as async communication between team members

## Quick Links

- [PR Merge Roadmap](../development/pr-merge-roadmap.md)
- [7EPs Index](../7eps/index.md)
- [MAS Foundation Tests](mas-foundation-integration.md)