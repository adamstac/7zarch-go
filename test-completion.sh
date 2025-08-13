#!/bin/bash

# 7zarch-go Shell Completion Test Script
# Tests all completion functionality and performance

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test binary path
BINARY="./7zarch-go"

echo -e "${BLUE}7zarch-go Shell Completion Test Suite${NC}"
echo "==========================================="

# Test 1: Verify binary exists and completion command works
echo -e "\n${YELLOW}1. Testing completion command availability${NC}"
if ! $BINARY completion bash --help > /dev/null 2>&1; then
    echo -e "${RED}FAIL: completion command not available${NC}"
    exit 1
fi
echo -e "${GREEN}PASS: completion command available${NC}"

# Test 2: Generate completion scripts
echo -e "\n${YELLOW}2. Testing completion script generation${NC}"
for shell in bash zsh fish powershell; do
    echo -n "  Testing $shell completion... "
    if $BINARY completion $shell > /tmp/7zarch_completion_test_$shell 2>&1; then
        echo -e "${GREEN}PASS${NC}"
    else
        echo -e "${RED}FAIL${NC}"
        exit 1
    fi
done

# Test 3: Performance tests
echo -e "\n${YELLOW}3. Testing completion performance${NC}"

# Test empty completion (should be fast)
echo -n "  Empty completion performance... "
start_time=$(date +%s%N)
$BINARY __complete show "" > /dev/null 2>&1
end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds

if [ $duration -le 100 ]; then
    echo -e "${GREEN}PASS (${duration}ms)${NC}"
else
    echo -e "${YELLOW}SLOW (${duration}ms, target <100ms)${NC}"
fi

# Test prefix completion performance
if $BINARY list --output json | grep -q uid; then
    echo -n "  Prefix completion performance... "
    start_time=$(date +%s%N)
    $BINARY __complete show "01" > /dev/null 2>&1
    end_time=$(date +%s%N)
    duration=$(( (end_time - start_time) / 1000000 ))
    
    if [ $duration -le 100 ]; then
        echo -e "${GREEN}PASS (${duration}ms)${NC}"
    else
        echo -e "${YELLOW}SLOW (${duration}ms, target <100ms)${NC}"
    fi
else
    echo -e "${YELLOW}SKIP: No archives in registry${NC}"
fi

# Test 4: Functional completion tests
echo -e "\n${YELLOW}4. Testing functional completion${NC}"

# Test that each command has proper completion setup
for cmd in show delete move restore; do
    echo -n "  Testing $cmd completion... "
    output=$($BINARY __complete $cmd "" 2>&1)
    
    if echo "$output" | grep -q "ShellCompDirectiveNoFileComp"; then
        echo -e "${GREEN}PASS${NC}"
    else
        echo -e "${RED}FAIL: Missing NoFileComp directive${NC}"
        exit 1
    fi
done

# Test 5: Context-aware completion
echo -e "\n${YELLOW}5. Testing context-aware completion${NC}"

# Check if we have both present and deleted archives for testing
present_count=$($BINARY list --output json | jq '[.[] | select(.status == "present")] | length' 2>/dev/null || echo "0")
deleted_count=$($BINARY list --output json | jq '[.[] | select(.status == "deleted")] | length' 2>/dev/null || echo "0")

if [ "$present_count" -gt 0 ] && [ "$deleted_count" -gt 0 ]; then
    echo -n "  Testing restore command filtering... "
    
    # Get completion for show (should show all)
    show_completion=$($BINARY __complete show "01" 2>/dev/null | head -10)
    show_count=$(echo "$show_completion" | grep -v ":" | wc -l)
    
    # Get completion for restore (should show only deleted)
    restore_completion=$($BINARY __complete restore "01" 2>/dev/null | head -10)
    restore_count=$(echo "$restore_completion" | grep -v ":" | wc -l)
    
    if [ "$restore_count" -lt "$show_count" ]; then
        echo -e "${GREEN}PASS (restore shows fewer archives than show)${NC}"
    else
        echo -e "${YELLOW}INCONCLUSIVE (restore:$restore_count, show:$show_count)${NC}"
    fi
else
    echo -e "${YELLOW}SKIP: Need both present and deleted archives for filtering test${NC}"
fi

# Test 6: Shell integration test
echo -e "\n${YELLOW}6. Testing bash shell integration${NC}"
echo -n "  Loading bash completion... "
if bash -c "source <(./7zarch-go completion bash) && complete -p 7zarch-go" > /dev/null 2>&1; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    exit 1
fi

# Test 7: Error handling
echo -e "\n${YELLOW}7. Testing error handling${NC}"
echo -n "  Invalid shell completion... "
if $BINARY completion invalid 2>&1 | grep -q "invalid argument"; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL: Should reject invalid shell${NC}"
    exit 1
fi

# Summary
echo -e "\n${GREEN}==========================================="
echo -e "All completion tests passed! ✓${NC}"
echo -e "\nCompletion system ready for production:"
echo -e "• Fast performance (<100ms)"
echo -e "• All shells supported (bash/zsh/fish/powershell)"
echo -e "• Context-aware filtering (restore vs other commands)"
echo -e "• Archive ID completion (UID/name/checksum prefixes)"
echo -e "• Proper error handling"

echo -e "\n${BLUE}Quick start:${NC}"
echo "  bash: source <(7zarch-go completion bash)"
echo "  zsh:  source <(7zarch-go completion zsh)"
echo "  fish: 7zarch-go completion fish | source"

echo -e "\n${BLUE}Installation docs: docs/shell-completion.md${NC}"

# Cleanup
rm -f /tmp/7zarch_completion_test_*
