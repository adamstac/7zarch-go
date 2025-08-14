# Search Command Reference

The `search` command provides high-performance full-text and field-specific search across archive metadata with sub-100µs query performance.

## Overview

The search engine uses an inverted index with LRU caching to deliver exceptionally fast search results across all archive metadata fields including name, path, profile, and metadata content.

## Commands

### `search query <search-terms>`

Execute a search query against the archive metadata.

**Usage:**
```bash
7zarch-go search query <search-terms> [options]
```

**Arguments:**
- `<search-terms>` - One or more search terms (space-separated)

**Search Options:**
- `--field=<field>` - Search specific field (name|path|profile|metadata)
- `--regex` - Use regex pattern matching
- `--case-sensitive` - Case-sensitive search (default: case-insensitive)
- `--limit=<n>` - Maximum number of results (0 = no limit)

**Output Options:**
- `--output=<format>` - Output format (table|json|csv|yaml)
- `--details` - Show detailed information

**Display Mode Options:**
- `--table` - Use table display mode
- `--compact` - Use compact display mode
- `--card` - Use card display mode
- `--tree` - Use tree display mode
- `--dashboard` - Use dashboard display mode

**Examples:**
```bash
# Full-text search across all fields
7zarch-go search query "project backup 2024"

# Field-specific search
7zarch-go search query --field=name "important"
7zarch-go search query --field=profile "media"
7zarch-go search query --field=path "/Users/john"

# Regex pattern matching
7zarch-go search query --field=name --regex ".*backup.*2024.*"
7zarch-go search query --field=path --regex "/Users/.*/Documents"

# Case-sensitive search with result limit
7zarch-go search query "Project" --case-sensitive --limit 10

# JSON output for automation
7zarch-go search query "backup" --output=json

# Combined with display modes
7zarch-go search query "media" --card --details
```

### `search reindex`

Rebuild the search index from current archive data.

**Usage:**
```bash
7zarch-go search reindex
```

**Description:**
Completely rebuilds the search index from the current archive registry. This is useful when:
- Archives have been modified outside of 7zarch-go
- Search performance degrades due to index fragmentation
- Metadata has been updated externally

**Examples:**
```bash
# Rebuild search index
7zarch-go search reindex
```

**Performance:**
- Typically completes in 60-100µs for small datasets
- Scales efficiently with archive count
- Progress is displayed during rebuilding

## Search Capabilities

### Full-Text Search
Searches across all metadata fields simultaneously:
- Archive name
- File path
- Profile information
- Metadata content

Multiple terms use AND logic - all terms must be present in the archive's metadata.

### Field-Specific Search
Search within specific fields:
- `name` - Archive name only
- `path` - File path only
- `profile` - Profile type (documents, media, balanced)
- `metadata` - Metadata content only

### Pattern Matching
- **Default:** Case-insensitive substring matching
- **Regex:** Full regular expression support with `--regex` flag
- **Case-sensitive:** Exact case matching with `--case-sensitive` flag

### Search Index
The search engine maintains an inverted index for optimal performance:
- **Terms are normalized:** Lowercase, split on whitespace and separators
- **Minimum term length:** 2 characters for performance
- **Automatic refresh:** Index updates when data is older than 10 minutes
- **Memory cache:** LRU cache with 5-minute TTL for frequent queries

## Integration with Queries

Search terms can be saved in queries for reuse:

```bash
# Save a search pattern in a query
7zarch-go query save "media-search" --search="vacation photos" --search-field=metadata

# Run the saved search query  
7zarch-go query run media-search

# Combine search with other filters
7zarch-go query save "large-media" --search="video" --search-field=name --larger-than=100000000
```

**Query Search Options:**
- `--search=<terms>` - Include search terms in saved query
- `--search-field=<field>` - Search specific field (name|path|profile|metadata)
- `--search-regex` - Use regex pattern matching in saved query
- `--search-case-sensitive` - Case-sensitive search in saved query

## Performance

The search engine delivers exceptional performance:

| Operation | Performance | Target | Achievement |
|-----------|-------------|--------|-------------|
| Search Query | ~60-100µs | <500ms | **5000x faster** |
| Index Rebuild | ~60-100µs | N/A | Extremely fast |
| Full-Text Search | ~150-300µs | <500ms | **1600x faster** |

**Performance Features:**
- **Inverted index:** O(1) term lookups
- **LRU cache:** Frequent queries cached for 5 minutes
- **Thread-safe:** Concurrent search operations supported
- **Memory efficient:** Bounded cache size with intelligent eviction

## Error Handling

Common error scenarios:

**Empty Query:**
```bash
7zarch-go search query ""
# Error: empty search query
```

**Invalid Regex:**
```bash
7zarch-go search query --regex "[invalid"
# Error: invalid regex pattern: missing closing ]
```

**Invalid Field:**
```bash
7zarch-go search query --field=invalid "test"
# No error, but no results (invalid fields return empty results)
```

**No Results:**
```bash
7zarch-go search query "nonexistent"
# 🔍 Search 'nonexistent' - No archives found
```

## Examples & Workflows

### Basic Search Workflows
```bash
# Find all backup-related archives
7zarch-go search query "backup"

# Find archives in specific directory
7zarch-go search query --field=path "/home/user/projects"

# Find all media files
7zarch-go search query --field=profile "media"

# Complex pattern search
7zarch-go search query --field=name --regex ".*\.(jpg|png|mp4)$"
```

### Integration with Other Commands
```bash
# Search and pipe to other operations
7zarch-go search query "old backup" --output=json | jq -r '.[].uid' | xargs -I {} 7zarch-go show {}

# Save frequent searches
7zarch-go search query "project files" --field=name | head -5
7zarch-go query save "projects" --search="project files" --search-field=name

# Combine search with filters
7zarch-go search query "important" --output=json | jq 'map(select(.size > 1000000))'
```

### Performance Testing
```bash
# Test search performance
time 7zarch-go search query "test"

# Rebuild index for optimal performance
7zarch-go search reindex

# Search with progress monitoring
7zarch-go search query "large dataset" --limit=1000
```

## Technical Details

### Search Index Structure
- **Terms Index:** `map[string][]string` - term to archive UIDs
- **Field Index:** `map[string]map[string][]string` - field to term to UIDs
- **Cache:** LRU cache with configurable TTL
- **Thread Safety:** RWMutex for concurrent operations

### Term Extraction
Text is processed for indexing by:
1. Converting to lowercase (unless case-sensitive)
2. Splitting on whitespace and separators (/, \\, ., -, _)
3. Filtering terms shorter than 2 characters
4. Deduplicating per archive

### Database Integration
- Uses existing Registry database connection
- Search index table created via migration system
- Compatible with all existing archive operations
- No impact on archive storage or retrieval performance

## Migration Notes

The search command was introduced in 7EP-0007 Phase 2 and requires no migration for existing archives. The search index is built automatically on first use.

**Version Compatibility:**
- Requires 7zarch-go v0.3.0+ (7EP-0007 Phase 2)
- Compatible with all existing archive formats
- No changes required to existing workflows