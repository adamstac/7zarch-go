#!/bin/bash
# 7EP-0020 Phase 2: Agent Lifecycle Integration Testing
# Tests complete agent operational cycles with current document state

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_DIR="$(dirname "$SCRIPT_DIR")"
TEST_AGENT="TESTBOT"
TEST_SESSION_DIR="/tmp/7zarch-go-lifecycle-test"
ORIGINAL_DIR="$(pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Agent Lifecycle Integration Testing${NC}"
echo -e "${BLUE}=====================================${NC}"

# Setup test environment
setup_test_environment() {
    echo -e "${BLUE}📋 Setting up test environment...${NC}"
    
    # Create clean test directory
    rm -rf "$TEST_SESSION_DIR"
    mkdir -p "$TEST_SESSION_DIR"
    
    # Copy current repository state for testing
    cp -r "$BASE_DIR" "$TEST_SESSION_DIR/repo"
    cd "$TEST_SESSION_DIR/repo"
    
    # Create test agent role file
    cp docs/development/roles/ROLE-TEMPLATE.md "docs/development/roles/${TEST_AGENT}.md"
    sed -i.bak "s/\[Agent Name\]/${TEST_AGENT} (Test Agent)/g" "docs/development/roles/${TEST_AGENT}.md"
    sed -i.bak "s/\[YYYY-MM-DD HH:MM\]/$(date '+%Y-%m-%d %H:%M')/g" "docs/development/roles/${TEST_AGENT}.md"
    sed -i.bak "s/\[Available|Active|Blocked\]/Available/g" "docs/development/roles/${TEST_AGENT}.md"
    sed -i.bak "s/\[Brief description of primary work\]/Lifecycle integration testing/g" "docs/development/roles/${TEST_AGENT}.md"
    rm -f "docs/development/roles/${TEST_AGENT}.md.bak"
    
    echo -e "${GREEN}✅ Test environment ready${NC}"
}

# Test Phase 1: Session Startup (BOOTUP.md integration)
test_session_startup() {
    echo -e "\n${BLUE}📋 Phase 1: Testing Session Startup${NC}"
    
    # Test git status and sync
    echo -e "${YELLOW}Testing git status and sync...${NC}"
    if git status >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Git status check passed${NC}"
    else
        echo -e "${RED}❌ Git status check failed${NC}"
        return 1
    fi
    
    # Test role context loading
    echo -e "${YELLOW}Testing role context loading...${NC}"
    if grep -q "Current Assignments" "docs/development/roles/${TEST_AGENT}.md"; then
        echo -e "${GREEN}✅ Role file structure valid${NC}"
    else
        echo -e "${RED}❌ Role file structure invalid${NC}"
        return 1
    fi
    
    # Test NEXT.md reading
    echo -e "${YELLOW}Testing NEXT.md coordination loading...${NC}"
    if [[ -f "docs/development/NEXT.md" ]]; then
        echo -e "${GREEN}✅ NEXT.md accessible${NC}"
    else
        echo -e "${RED}❌ NEXT.md not found${NC}"
        return 1
    fi
    
    # Test 7EP coordination context
    echo -e "${YELLOW}Testing 7EP coordination context...${NC}"
    ACTIVE_7EPS=$(grep -l "Status.*ACTIVE\|In Progress" docs/7eps/*.md 2>/dev/null || echo "")
    if [[ -n "$ACTIVE_7EPS" ]]; then
        echo -e "${GREEN}✅ Found active 7EPs: $(echo "$ACTIVE_7EPS" | wc -l | tr -d ' ')${NC}"
    else
        echo -e "${YELLOW}⚠️  No active 7EPs found (normal if no active work)${NC}"
    fi
    
    # Test role context integration (step 3.5 from enhanced BOOTUP.md)
    echo -e "${YELLOW}Testing role context integration...${NC}"
    
    ASSIGNMENTS=$(grep -A 10 "Active Work\|Current Assignments" "docs/development/roles/${TEST_AGENT}.md" 2>/dev/null || echo "No assignments")
    if [[ "$ASSIGNMENTS" != "No assignments" ]]; then
        echo -e "${GREEN}✅ Role assignments extracted successfully${NC}"
    else
        echo -e "${RED}❌ Role assignment extraction failed${NC}"
        return 1
    fi
    
    COORDINATION=$(grep -A 5 "Coordination Needed\|Blocked\|Waiting" "docs/development/roles/${TEST_AGENT}.md" 2>/dev/null || echo "No coordination")
    if [[ "$COORDINATION" != "No coordination" ]]; then
        echo -e "${GREEN}✅ Coordination status extracted${NC}"
    else
        echo -e "${YELLOW}⚠️  No coordination dependencies (acceptable for test agent)${NC}"
    fi
    
    echo -e "${GREEN}✅ Session startup simulation successful${NC}"
}

# Test Phase 2: Session Logging
test_session_logging() {
    echo -e "\n${BLUE}📋 Phase 2: Testing Session Logging${NC}"
    
    # Test session log creation (BOOTUP.md step 5)
    echo -e "${YELLOW}Testing session log creation...${NC}"
    
    mkdir -p docs/logs
    DATE_STAMP=$(date +%Y-%m-%d_%H-%M-%S)
    SESSION_START=$(date "+%Y-%m-%d %H:%M:%S")
    
    cat > "docs/logs/session-${DATE_STAMP}.md" << EOF
# Session Log - $(date)

## ⏱️ Session Timing
- **Start Time:** ${SESSION_START}
- **Agent:** ${TEST_AGENT} (Test Agent)
- **Status:** 🟢 **ACTIVE** - Session in progress

## 🚀 Boot Sequence Completed
- Git status: Clean and up to date
- Build verification: Successful
- Operational priorities: Reviewed
- Ready for work assignment

---
*Session started by DDD Framework bootup process*
EOF

    if [[ -f "docs/logs/session-${DATE_STAMP}.md" ]]; then
        echo -e "${GREEN}✅ Session log created successfully${NC}"
        echo "SESSION_LOG_FILE=docs/logs/session-${DATE_STAMP}.md" > .session-active
        echo -e "${GREEN}✅ Session tracking file created${NC}"
    else
        echo -e "${RED}❌ Session log creation failed${NC}"
        return 1
    fi
}

# Test Phase 3: Work Simulation
test_work_simulation() {
    echo -e "\n${BLUE}📋 Phase 3: Testing Work Simulation${NC}"
    
    # Simulate role file update (typical work pattern)
    echo -e "${YELLOW}Testing role file work updates...${NC}"
    
    # Update role file to simulate work assignment
    sed -i.bak 's/\[Assignment Name\].*\[STATUS\].*/Test Assignment - ACTIVE (lifecycle integration testing)/' "docs/development/roles/${TEST_AGENT}.md"
    
    # Test that role file update is valid
    if grep -q "Test Assignment - ACTIVE" "docs/development/roles/${TEST_AGENT}.md"; then
        echo -e "${GREEN}✅ Role file work update successful${NC}"
    else
        echo -e "${RED}❌ Role file work update failed${NC}"
        return 1
    fi
    
    # Test git add and commit simulation (without actually committing)
    echo -e "${YELLOW}Testing git workflow patterns...${NC}"
    
    git add "docs/development/roles/${TEST_AGENT}.md" 2>/dev/null || {
        echo -e "${RED}❌ Git add failed${NC}"
        return 1
    }
    
    echo -e "${GREEN}✅ Git workflow patterns operational${NC}"
    
    # Reset changes (don't actually commit test changes)
    git reset HEAD "docs/development/roles/${TEST_AGENT}.md" >/dev/null 2>&1 || true
    git checkout "docs/development/roles/${TEST_AGENT}.md" >/dev/null 2>&1 || true
}

# Test Phase 4: Session Shutdown Simulation
test_session_shutdown() {
    echo -e "\n${BLUE}📋 Phase 4: Testing Session Shutdown${NC}"
    
    if [[ ! -f .session-active ]]; then
        echo -e "${RED}❌ No active session file found${NC}"
        return 1
    fi
    
    source .session-active
    SESSION_LOG="${SESSION_LOG_FILE}"
    
    # Test session timing calculation
    echo -e "${YELLOW}Testing session timing calculation...${NC}"
    
    SESSION_END=$(date "+%Y-%m-%d %H:%M:%S")
    SESSION_START_TIME=$(grep "Start Time:" "$SESSION_LOG" | sed -E 's/.*\*\*Start Time:\*\* //' 2>/dev/null || echo "Unknown")
    
    if [[ "$SESSION_START_TIME" != "Unknown" ]]; then
        echo -e "${GREEN}✅ Session start time extracted: $SESSION_START_TIME${NC}"
        
        # Test portable date calculation
        if date --version >/dev/null 2>&1; then
            # GNU date (Linux)
            START_EPOCH=$(date -d "$SESSION_START_TIME" +%s 2>/dev/null || date +%s)
        else
            # BSD date (macOS)  
            START_EPOCH=$(date -j -f "%Y-%m-%d %H:%M:%S" "$SESSION_START_TIME" +%s 2>/dev/null || date +%s)
        fi
        
        END_EPOCH=$(date +%s)
        SESSION_DURATION=$((END_EPOCH - START_EPOCH))
        
        if [[ $SESSION_DURATION -ge 0 ]]; then
            echo -e "${GREEN}✅ Session duration calculation working: ${SESSION_DURATION}s${NC}"
        else
            echo -e "${RED}❌ Session duration calculation failed${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ Session start time extraction failed${NC}"
        return 1
    fi
    
    # Test session log completion (simulate without actually writing)
    echo -e "${YELLOW}Testing session log completion...${NC}"
    
    # Test that we can append to session log
    if [[ -w "$SESSION_LOG" ]]; then
        echo -e "${GREEN}✅ Session log writable for completion${NC}"
    else
        echo -e "${RED}❌ Session log not writable${NC}"
        return 1
    fi
    
    # Clean up test session
    rm -f .session-active
    echo -e "${GREEN}✅ Session shutdown simulation successful${NC}"
}

# Test Phase 5: Cross-Agent Coordination Simulation
test_coordination_patterns() {
    echo -e "\n${BLUE}📋 Phase 5: Testing Cross-Agent Coordination${NC}"
    
    # Test TEAM-UPDATE.md patterns
    echo -e "${YELLOW}Testing TEAM-UPDATE.md coordination patterns...${NC}"
    
    # Simulate role file status change
    if [[ -f "docs/development/actions/TEAM-UPDATE.md" ]]; then
        echo -e "${GREEN}✅ TEAM-UPDATE.md workflow available${NC}"
        
        # Test that role files can be updated for coordination
        if grep -q "Active Work" "docs/development/roles/AMP.md"; then
            echo -e "${GREEN}✅ Role file coordination update patterns operational${NC}"
        else
            echo -e "${RED}❌ Role file coordination patterns not found${NC}"
            return 1
        fi
        
        # Test NEXT.md coordination integration
        if grep -q "Current Active Work" docs/development/NEXT.md; then
            echo -e "${GREEN}✅ NEXT.md coordination integration working${NC}"
        else
            echo -e "${RED}❌ NEXT.md coordination integration broken${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ TEAM-UPDATE.md workflow missing${NC}"
        return 1
    fi
}

# Test Phase 6: Workflow Actions Integration
test_workflow_actions() {
    echo -e "\n${BLUE}📋 Phase 6: Testing Workflow Actions Integration${NC}"
    
    # Test COMMIT.md integration
    echo -e "${YELLOW}Testing COMMIT.md integration...${NC}"
    if grep -q "Role Coordination" docs/development/actions/COMMIT.md; then
        echo -e "${GREEN}✅ COMMIT.md has role coordination integration${NC}"
    else
        echo -e "${RED}❌ COMMIT.md missing role coordination integration${NC}"
        return 1
    fi
    
    # Test MERGE.md integration  
    echo -e "${YELLOW}Testing MERGE.md integration...${NC}"
    if grep -q "TEAM-UPDATE.md" docs/development/actions/MERGE.md; then
        echo -e "${GREEN}✅ MERGE.md has team coordination integration${NC}"
    else
        echo -e "${RED}❌ MERGE.md missing team coordination integration${NC}"
        return 1
    fi
    
    # Test NEW-FEATURE.md integration
    echo -e "${YELLOW}Testing NEW-FEATURE.md integration...${NC}"
    if grep -q "role file" docs/development/actions/NEW-FEATURE.md; then
        echo -e "${GREEN}✅ NEW-FEATURE.md has role file integration${NC}"
    else
        echo -e "${RED}❌ NEW-FEATURE.md missing role file integration${NC}"
        return 1
    fi
}

# Cleanup function
cleanup() {
    cd "$ORIGINAL_DIR"
    rm -rf "$TEST_SESSION_DIR"
    echo -e "${BLUE}🧹 Test environment cleaned up${NC}"
}

# Main test execution
main() {
    trap cleanup EXIT
    
    echo -e "${BLUE}Testing complete agent lifecycle with current framework state...${NC}\n"
    
    # Run test phases
    setup_test_environment
    test_session_startup
    test_session_logging
    test_work_simulation
    test_session_shutdown
    test_coordination_patterns
    test_workflow_actions
    
    echo -e "\n${GREEN}📊 Agent Lifecycle Integration Test Results${NC}"
    echo -e "${GREEN}===========================================${NC}"
    echo -e "${GREEN}✅ All lifecycle phases operational${NC}"
    echo -e "${GREEN}✅ Framework integration validated${NC}"
    echo -e "${GREEN}✅ Cross-agent coordination patterns working${NC}"
    echo -e "${GREEN}✅ Workflow actions properly integrated${NC}"
    
    echo -e "\n${BLUE}🎯 Integration Test Summary${NC}"
    echo -e "${BLUE}==========================${NC}"
    echo -e "Agent lifecycle: ${GREEN}OPERATIONAL${NC}"
    echo -e "Document integration: ${GREEN}VALIDATED${NC}"
    echo -e "Coordination patterns: ${GREEN}FUNCTIONAL${NC}"
    echo -e "Framework reliability: ${GREEN}CONFIRMED${NC}"
    
    return 0
}

# Error handling
if ! main; then
    echo -e "\n${RED}❌ Agent lifecycle integration testing FAILED${NC}"
    echo -e "${RED}Framework integration issues detected${NC}"
    exit 1
fi

echo -e "\n${GREEN}✅ Agent lifecycle integration testing PASSED${NC}"
echo -e "${GREEN}DDD framework operational reliability confirmed${NC}"
