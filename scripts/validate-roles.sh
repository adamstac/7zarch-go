#!/bin/bash
# 7EP-0019 Role File Quality Validation
# Checks compliance with agent lifecycle framework standards

set -euo pipefail

ROLE_DIR="docs/development/roles"
FAILED=0
TOTAL_ROLES=0

echo "🔍 7EP-0019 Role File Validation"
echo "=================================="

# Function to check individual role file
validate_role() {
    local role_file="$1"
    local role_name=$(basename "$role_file" .md)
    local errors=0
    
    echo "📋 Checking $role_name..."
    
    # Check required header fields
    if ! grep -q "^\*\*Last Updated:\*\*" "$role_file"; then
        echo "  ❌ Missing 'Last Updated' header field"
        ((errors++))
    fi
    
    if ! grep -q "^\*\*Status:\*\*" "$role_file"; then
        echo "  ❌ Missing 'Status' header field" 
        ((errors++))
    fi
    
    if ! grep -q "^\*\*Current Focus:\*\*" "$role_file"; then
        echo "  ❌ Missing 'Current Focus' header field"
        ((errors++))
    fi
    
    # Check for required sections
    if ! grep -q "## 🎯 Current Assignments" "$role_file"; then
        echo "  ❌ Missing '🎯 Current Assignments' section"
        ((errors++))
    fi
    
    if ! grep -q "## 🔗 Coordination Needed" "$role_file"; then
        echo "  ❌ Missing '🔗 Coordination Needed' section"
        ((errors++))
    fi
    
    if ! grep -q "## ✅ Recently Completed" "$role_file"; then
        echo "  ❌ Missing '✅ Recently Completed' section"
        ((errors++))
    fi
    
    if ! grep -q "## 📝 Implementation Notes" "$role_file"; then
        echo "  ❌ Missing '📝 Implementation Notes' section"
        ((errors++))
    fi
    
    # Check for content boundary violations
    if grep -q "Adam Stacoviak.*@adamstac.*Project owner" "$role_file"; then
        echo "  ❌ Team context duplication (should reference TEAM-CONTEXT.md)"
        ((errors++))
    fi
    
    # Check for standard subsections in Current Assignments
    if ! grep -q "### Active Work" "$role_file"; then
        echo "  ❌ Missing 'Active Work' subsection in Current Assignments"
        ((errors++))
    fi
    
    if ! grep -q "### Next Priorities" "$role_file"; then
        echo "  ❌ Missing 'Next Priorities' subsection in Current Assignments"
        ((errors++))
    fi
    
    # Check for team context reference (except ADAM.md which may have strategic context)
    if [[ "$role_name" != "ADAM" ]] && ! grep -q "TEAM-CONTEXT.md" "$role_file"; then
        echo "  ⚠️  No reference to TEAM-CONTEXT.md (recommended)"
    fi
    
    # Success message or error count
    if [[ $errors -eq 0 ]]; then
        echo "  ✅ $role_name compliant"
    else
        echo "  ❌ $role_name has $errors compliance issues"
        ((FAILED += errors))
    fi
    
    return $errors
}

# Main validation loop
for role_file in "$ROLE_DIR"/*.md; do
    # Skip template and readme files
    filename=$(basename "$role_file")
    if [[ "$filename" == "ROLE-TEMPLATE.md" || "$filename" == "README.md" ]]; then
        continue
    fi
    
    if [[ -f "$role_file" ]]; then
        ((TOTAL_ROLES++))
        validate_role "$role_file"
        echo ""
    fi
done

# Summary
echo "📊 Validation Summary"
echo "===================="
echo "Roles checked: $TOTAL_ROLES"
echo "Total issues: $FAILED"

if [[ $FAILED -eq 0 ]]; then
    echo "✅ All role files are compliant with 7EP-0019 standards!"
    exit 0
else
    echo "❌ Found $FAILED compliance issues"
    echo ""
    echo "💡 Fix Issues:"
    echo "- Add missing header fields (Last Updated, Status, Current Focus)"
    echo "- Add missing required sections (🎯 🔗 ✅ 📝)" 
    echo "- Remove team context duplication (reference TEAM-CONTEXT.md instead)"
    echo "- Follow ROLE-TEMPLATE.md structure"
    exit 1
fi
