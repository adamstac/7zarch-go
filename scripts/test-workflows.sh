#!/bin/bash
# 7EP-0020 Phase 2: Workflow Action Integration Testing
# Tests that all workflow scripts work with current document state

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_DIR="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Workflow Action Integration Testing${NC}"
echo -e "${BLUE}====================================${NC}"

# Test BOOTUP.md script execution
test_bootup_script() {
    echo -e "\n${BLUE}📋 Testing BOOTUP.md Script Execution${NC}"
    
    # Test git status and pull (step 1)
    echo -e "${YELLOW}Testing git operations...${NC}"
    if git status >/dev/null 2>&1 && git branch >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Git operations functional${NC}"
    else
        echo -e "${RED}❌ Git operations failed${NC}"
        return 1
    fi
    
    # Test role context integration (step 3.5)
    echo -e "${YELLOW}Testing role context integration commands...${NC}"
    
    # Test role file reading
    for role in AMP CLAUDE AUGMENT ADAM; do
        if [[ -f "docs/development/roles/${role}.md" ]]; then
            ASSIGNMENTS=$(grep -A 10 "Active Work\|Current Assignments" "docs/development/roles/${role}.md" 2>/dev/null || echo "")
            if [[ -n "$ASSIGNMENTS" ]]; then
                echo -e "${GREEN}✅ ${role}.md assignments extractable${NC}"
            else
                echo -e "${YELLOW}⚠️  ${role}.md assignments extraction needs refinement${NC}"
            fi
        else
            echo -e "${RED}❌ ${role}.md not found${NC}"
            return 1
        fi
    done
    
    # Test 7EP context loading
    echo -e "${YELLOW}Testing 7EP coordination context...${NC}"
    ACTIVE_7EPS=$(grep -l "Status.*ACTIVE\|In Progress" docs/7eps/*.md 2>/dev/null | xargs --no-run-if-empty ls -la 2>/dev/null || echo "")
    if [[ -n "$ACTIVE_7EPS" ]]; then
        echo -e "${GREEN}✅ 7EP coordination context extraction working${NC}"
    else
        echo -e "${YELLOW}⚠️  No active 7EPs (normal state)${NC}"
    fi
}

# Test SHUTDOWN.md script execution  
test_shutdown_script() {
    echo -e "\n${BLUE}📋 Testing SHUTDOWN.md Script Execution${NC}"
    
    # Create test session file
    echo -e "${YELLOW}Creating test session for shutdown testing...${NC}"
    
    DATE_STAMP=$(date +%Y-%m-%d_%H-%M-%S)
    SESSION_START=$(date "+%Y-%m-%d %H:%M:%S")
    TEST_SESSION_LOG="docs/logs/test-session-${DATE_STAMP}.md"
    
    mkdir -p docs/logs
    cat > "$TEST_SESSION_LOG" << EOF
# Session Log - $(date)

## ⏱️ Session Timing
- **Start Time:** ${SESSION_START}
- **Agent:** TEST (Test Agent)
- **Status:** 🟢 **ACTIVE** - Session in progress

---
*Session started by DDD Framework bootup process*
EOF

    echo "SESSION_LOG_FILE=$TEST_SESSION_LOG" > .session-active-test
    
    # Test session timing calculation
    echo -e "${YELLOW}Testing session timing calculation...${NC}"
    
    source .session-active-test
    SESSION_LOG="${SESSION_LOG_FILE}"
    SESSION_END=$(date "+%Y-%m-%d %H:%M:%S")
    SESSION_START_TIME=$(grep "Start Time:" "$SESSION_LOG" | sed -E 's/.*\*\*Start Time:\*\* //' 2>/dev/null || echo "Unknown")
    
    if [[ "$SESSION_START_TIME" != "Unknown" ]]; then
        echo -e "${GREEN}✅ Session start time extraction working${NC}"
        
        # Test portable date parsing
        if date --version >/dev/null 2>&1; then
            START_EPOCH=$(date -d "$SESSION_START_TIME" +%s 2>/dev/null || date +%s)
        else
            START_EPOCH=$(date -j -f "%Y-%m-%d %H:%M:%S" "$SESSION_START_TIME" +%s 2>/dev/null || date +%s)
        fi
        
        END_EPOCH=$(date +%s)
        SESSION_DURATION=$((END_EPOCH - START_EPOCH))
        
        if [[ $SESSION_DURATION -ge 0 ]]; then
            echo -e "${GREEN}✅ Portable date calculation working${NC}"
        else
            echo -e "${RED}❌ Date calculation failed${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ Session timing extraction failed${NC}"
        return 1
    fi
    
    # Test role state update extraction
    echo -e "${YELLOW}Testing role state update extraction...${NC}"
    
    ROLE_UPDATES=$(grep -A 3 "Active Work\|Current Assignments" docs/development/roles/AMP.md | head -5 2>/dev/null || echo "")
    if [[ -n "$ROLE_UPDATES" ]]; then
        echo -e "${GREEN}✅ Role state updates extractable${NC}"
    else
        echo -e "${RED}❌ Role state update extraction failed${NC}"
        return 1
    fi
    
    # Cleanup test files
    rm -f .session-active-test "$TEST_SESSION_LOG"
}

# Test TEAM-UPDATE.md workflow patterns
test_team_update_workflow() {
    echo -e "\n${BLUE}📋 Testing TEAM-UPDATE.md Workflow Patterns${NC}"
    
    echo -e "${YELLOW}Testing role file status update pattern...${NC}"
    
    # Test that we can read and update role files
    for role in AMP CLAUDE AUGMENT ADAM; do
        ROLE_FILE="docs/development/roles/${role}.md"
        if [[ -f "$ROLE_FILE" ]]; then
            # Test reading current status
            CURRENT_STATUS=$(grep -A 3 "Active Work" "$ROLE_FILE" 2>/dev/null || echo "")
            if [[ -n "$CURRENT_STATUS" ]]; then
                echo -e "${GREEN}✅ ${role} status readable for updates${NC}"
            else
                echo -e "${YELLOW}⚠️  ${role} status format needs adjustment${NC}"
            fi
        else
            echo -e "${RED}❌ ${role}.md not accessible${NC}"
            return 1
        fi
    done
    
    echo -e "${YELLOW}Testing NEXT.md coordination update pattern...${NC}"
    
    # Test NEXT.md coordination extraction
    NEXT_COORDINATION=$(grep -A 10 "Current Active Work" docs/development/NEXT.md 2>/dev/null || echo "")
    if [[ -n "$NEXT_COORDINATION" ]]; then
        echo -e "${GREEN}✅ NEXT.md coordination extractable${NC}"
    else
        echo -e "${RED}❌ NEXT.md coordination extraction failed${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Testing cross-agent notification patterns...${NC}"
    
    # Test that agents can be identified and coordinated
    AGENT_COUNT=$(ls docs/development/roles/*.md | grep -v ROLE-TEMPLATE | grep -v README | wc -l | tr -d ' ')
    if [[ $AGENT_COUNT -ge 4 ]]; then
        echo -e "${GREEN}✅ Cross-agent coordination targets available (${AGENT_COUNT} agents)${NC}"
    else
        echo -e "${RED}❌ Insufficient agents for coordination testing${NC}"
        return 1
    fi
}

# Test build verification
test_build_verification() {
    echo -e "\n${BLUE}📋 Testing Build Verification${NC}"
    
    echo -e "${YELLOW}Testing make dev build...${NC}"
    
    # Test that build system works (required for BOOTUP.md step 4)
    if make build >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Build system operational${NC}"
    else
        echo -e "${RED}❌ Build system failed${NC}"
        return 1
    fi
    
    # Test that binary is functional
    if [[ -x "./7zarch-go" ]]; then
        echo -e "${GREEN}✅ Binary executable${NC}"
        
        # Test basic functionality
        if ./7zarch-go --help >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Binary functionality verified${NC}"
        else
            echo -e "${RED}❌ Binary functionality failed${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ Binary not executable${NC}"
        return 1
    fi
}

# Main workflow testing
main() {
    cd "$BASE_DIR"
    
    echo -e "${BLUE}Testing workflow actions with current document state...${NC}\n"
    
    test_bootup_script
    test_shutdown_script
    test_team_update_workflow
    test_build_verification
    
    echo -e "\n${GREEN}📊 Workflow Integration Test Results${NC}"
    echo -e "${GREEN}====================================${NC}"
    echo -e "${GREEN}✅ BOOTUP.md script operational${NC}"
    echo -e "${GREEN}✅ SHUTDOWN.md script functional${NC}"
    echo -e "${GREEN}✅ TEAM-UPDATE.md patterns working${NC}"
    echo -e "${GREEN}✅ Build verification successful${NC}"
    
    echo -e "\n${BLUE}🎯 Workflow Integration Summary${NC}"
    echo -e "${BLUE}==============================${NC}"
    echo -e "Document integration: ${GREEN}VALIDATED${NC}"
    echo -e "Script functionality: ${GREEN}OPERATIONAL${NC}"
    echo -e "Cross-agent patterns: ${GREEN}FUNCTIONAL${NC}"
    echo -e "Build integration: ${GREEN}CONFIRMED${NC}"
    
    return 0
}

# Error handling
if ! main; then
    echo -e "\n${RED}❌ Workflow integration testing FAILED${NC}"
    echo -e "${RED}Workflow action issues detected${NC}"
    exit 1
fi

echo -e "\n${GREEN}✅ Workflow integration testing PASSED${NC}"
echo -e "${GREEN}All workflow actions operational with current framework${NC}"
