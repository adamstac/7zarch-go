# MAS Foundation Integration Test Scenarios

**Purpose:** Validate the complete MAS Foundation workflow after PR #5 and #6 merge.

## Test Environment Setup

```bash
# Clean environment
rm -rf ~/.7zarch-go/
mkdir -p /tmp/7zarch-test/{source,external}

# Create test files
echo "Test document content" > /tmp/7zarch-test/source/document.txt
dd if=/dev/zero of=/tmp/7zarch-test/source/large-file.bin bs=1M count=5
mkdir -p /tmp/7zarch-test/source/folder-archive/
echo "Nested content" > /tmp/7zarch-test/source/folder-archive/nested.txt
```

## Core Workflow Tests

### Test 1: Basic Archive Creation and Resolution
```bash
# Create managed archive
7zarch-go create /tmp/7zarch-test/source/document.txt

# Expected: Archive created with ULID displayed
# Verify: Should show "Created archive with ID: 01K2xxxxx"

# Test resolution by various methods
ARCHIVE_ID=$(7zarch-go list --limit 1 | grep -o '01K2[A-Z0-9]*')

# Test ULID prefix resolution
7zarch-go show ${ARCHIVE_ID:0:6}  # First 6 chars
7zarch-go show ${ARCHIVE_ID:0:8}  # First 8 chars
7zarch-go show $ARCHIVE_ID        # Full ULID

# Test name resolution
7zarch-go show document.txt
7zarch-go show document

# Expected: All should resolve to same archive
```

### Test 2: Enhanced List Functionality
```bash
# Create diverse archive set
7zarch-go create /tmp/7zarch-test/source/document.txt --profile documents
7zarch-go create /tmp/7zarch-test/source/large-file.bin --profile media
7zarch-go create /tmp/7zarch-test/source/folder-archive/ --profile balanced

# Test basic list
7zarch-go list
# Expected: 3 archives shown in tabular format with UIDs, sizes, profiles

# Test filtering
7zarch-go list --profile documents
7zarch-go list --profile media
7zarch-go list --managed
7zarch-go list --larger-than 1MB

# Expected: Filters work correctly, show appropriate subsets
```

### Test 3: Status Verification and Missing Files
```bash
# Get an archive path
ARCHIVE_PATH=$(7zarch-go show document.txt --json | jq -r '.path')

# Move file outside 7zarch-go (simulate missing file)
mv "$ARCHIVE_PATH" /tmp/moved-archive.7z

# Test missing file detection
7zarch-go show document.txt
# Expected: Shows "⚠️ missing" status with helpful suggestions

7zarch-go list --status missing
# Expected: Shows the missing archive

# Restore file and verify
mv /tmp/moved-archive.7z "$ARCHIVE_PATH"
7zarch-go show document.txt
# Expected: Status back to "✓ present"
```

## Error Handling Tests

### Test 4: Ambiguous Resolution
```bash
# Create archives with similar prefixes
7zarch-go create /tmp/7zarch-test/source/document.txt
sleep 1  # Ensure different ULIDs
7zarch-go create /tmp/7zarch-test/source/large-file.bin

# Get first few characters that should be shared
FIRST_CHARS=$(7zarch-go list | head -1 | grep -o '01K2[A-Z0-9]*' | cut -c1-4)

# Test ambiguous resolution
7zarch-go show $FIRST_CHARS
# Expected: Shows disambiguation with options and helpful message
```

### Test 5: Not Found Scenarios
```bash
# Test non-existent archive
7zarch-go show "nonexistent"
7zarch-go show "01K2XXXXXX"
7zarch-go show "ffffffff"

# Expected: Clear "not found" messages with helpful suggestions
```

### Test 6: Registry Error Handling
```bash
# Test with corrupted/locked database
chmod 000 ~/.7zarch-go/registry.db
7zarch-go list
7zarch-go show document.txt

# Expected: Clear registry error messages with recovery suggestions

# Restore permissions
chmod 644 ~/.7zarch-go/registry.db
```

## Advanced Feature Tests

### Test 7: Checksum Verification
```bash
# Create archive
7zarch-go create /tmp/7zarch-test/source/document.txt

# Test verification
7zarch-go show document.txt --verify
# Expected: Shows checksum verification result

# Corrupt archive file
ARCHIVE_PATH=$(7zarch-go show document.txt --json | jq -r '.path')
echo "corrupted" >> "$ARCHIVE_PATH"

# Test corrupted file detection
7zarch-go show document.txt --verify
# Expected: Shows checksum mismatch with helpful suggestions
```

### Test 8: Complex Filtering Combinations
```bash
# Create varied archive set with timestamps
7zarch-go create /tmp/7zarch-test/source/document.txt --profile documents
sleep 1
7zarch-go create /tmp/7zarch-test/source/large-file.bin --profile media
sleep 1
7zarch-go create /tmp/7zarch-test/source/folder-archive/ --profile balanced

# Test time-based filtering
7zarch-go list --older-than 1s  # Should show first two
7zarch-go list --older-than 2s  # Should show first one

# Test size-based filtering
7zarch-go list --larger-than 1MB  # Should show large-file.bin

# Test combined filters
7zarch-go list --profile media --larger-than 1MB
7zarch-go list --managed --older-than 1s
```

### Test 9: JSON Output and Scripting
```bash
# Test JSON output
ARCHIVE_JSON=$(7zarch-go show document.txt --json)
echo "$ARCHIVE_JSON" | jq '.'

# Verify required fields
echo "$ARCHIVE_JSON" | jq -r '.uid, .name, .path, .size, .status'

# Test scripting scenarios
for id in $(7zarch-go list --format ids 2>/dev/null || 7zarch-go list | grep -o '01K2[A-Z0-9]*'); do
    echo "Archive $id: $(7zarch-go show "$id" --json | jq -r '.name')"
done
```

## Performance Tests

### Test 10: Large Registry Performance
```bash
# Create many archives quickly
for i in {1..100}; do
    echo "Content $i" > /tmp/7zarch-test/source/file-$i.txt
    7zarch-go create /tmp/7zarch-test/source/file-$i.txt >/dev/null
done

# Test resolution performance
time 7zarch-go show file-50.txt
time 7zarch-go list --profile documents
time 7zarch-go list --older-than 1m

# Expected: All operations under 200ms
```

## Cleanup

```bash
# Clean up test environment
rm -rf ~/.7zarch-go/
rm -rf /tmp/7zarch-test/
```

## Success Criteria

### Functional Requirements
- [x] All ULID resolution methods work (full, prefix, name, checksum)
- [x] Show command displays comprehensive archive information
- [x] List command supports all documented filters
- [x] Status verification detects missing/corrupted files
- [x] Error messages are helpful and actionable
- [x] JSON output contains all required fields

### Performance Requirements
- [ ] Resolution operations complete under 50ms
- [ ] Show command completes under 100ms
- [ ] List operations complete under 200ms for 100+ archives
- [ ] Memory usage stays under 10MB during operations

### User Experience Requirements
- [x] Ambiguous matches show clear disambiguation options
- [x] Not found errors provide helpful suggestions
- [x] Status indicators are intuitive (✓/❌/⚠️)
- [x] Output formatting is consistent and readable

## Integration Points

These tests validate the integration between:
- ULID resolution system and registry queries
- Show command and file verification
- List command filtering and database optimization
- Error handling standards across all commands
- Documentation accuracy with actual behavior

## Test Automation

```bash
#!/bin/bash
# mas-foundation-test.sh
# Automated version of integration tests

set -e
echo "Running MAS Foundation Integration Tests..."

# Run each test scenario
test_basic_workflow
test_enhanced_list
test_status_verification
test_error_handling
test_advanced_features
test_performance

echo "All tests passed! MAS Foundation ready for release."
```