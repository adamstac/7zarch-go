# Query Command Reference

The `query` command manages saved filter combinations for reuse, enabling power users to save complex filtering logic and recall it with simple names.

## Overview

Saved queries store filter combinations (including search terms) as named configurations that can be executed repeatedly. This eliminates the need to retype complex filter combinations and enables workflow automation.

## Commands

### `query save <name> [filters...]`

Save a new query with specified filters.

**Usage:**
```bash
7zarch-go query save <name> [filter-flags...]
```

**Arguments:**
- `<name>` - Unique name for the query

**Filter Flags (same as list command):**
- `--not-uploaded` - Archives that haven't been uploaded
- `--pattern=<pattern>` - Filter by name pattern (glob syntax)
- `--older-than=<duration>` - Archives older than duration (e.g., '7d', '1h')
- `--managed` - Managed archives only
- `--external` - External archives only  
- `--missing` - Missing archives only
- `--status=<status>` - Filter by status (present|missing|deleted)
- `--profile=<profile>` - Filter by profile (media|documents|balanced)
- `--larger-than=<size>` - Archives larger than size
- `--deleted` - Deleted archives only

**Search Integration Flags (7EP-0007 Phase 2):**
- `--search=<terms>` - Include search terms in saved query
- `--search-field=<field>` - Search specific field (name|path|profile|metadata)
- `--search-regex` - Use regex pattern matching for search
- `--search-case-sensitive` - Case-sensitive search

**Examples:**
```bash
# Save filter combinations
7zarch-go query save "my-docs" --profile=documents --managed
7zarch-go query save "large-media" --profile=media --larger-than=100000000
7zarch-go query save "old-unuploaded" --not-uploaded --older-than=30d

# Save queries with search terms
7zarch-go query save "backup-files" --search="backup" --search-field=name
7zarch-go query save "project-docs" --search="project" --profile=documents --managed

# Complex combinations
7zarch-go query save "cleanup-candidates" --external --older-than=6m --not-uploaded
```

### `query list`

List all saved queries with usage statistics.

**Usage:**
```bash
7zarch-go query list [options]
```

**Options:**
- `--output=<format>` - Output format (table|json)

**Examples:**
```bash
# List all saved queries
7zarch-go query list

# JSON output for scripting
7zarch-go query list --output json
```

**Output Format:**
```
📋 Saved Queries (3 found)

NAME             CREATED          LAST USED        USE COUNT  FILTERS
my-docs          2025-08-13 14:30 2025-08-13 15:45 5          --profile=documents --managed
large-media      2025-08-13 14:32 never            0          --profile=media --larger-than=100000000
backup-files     2025-08-13 14:35 2025-08-13 15:30 2          --search=backup --search-field=name
```

### `query run <name>`

Execute a saved query and display matching archives.

**Usage:**
```bash
7zarch-go query run <name> [display-options...]
```

**Arguments:**
- `<name>` - Name of saved query to execute

**Display Options:**
- `--table` - Use table display mode
- `--compact` - Use compact display mode
- `--card` - Use card display mode
- `--tree` - Use tree display mode
- `--dashboard` - Use dashboard display mode
- `--output=<format>` - Output format (table|json|csv|yaml)
- `--details` - Show detailed information

**Examples:**
```bash
# Run a saved query
7zarch-go query run my-docs

# Run with specific display mode
7zarch-go query run large-media --card

# JSON output for automation
7zarch-go query run backup-files --output json

# Combine with display modes
7zarch-go query run old-unuploaded --table --details
```

### `query show <name>`

Show details of a specific saved query.

**Usage:**
```bash
7zarch-go query show <name> [options]
```

**Arguments:**
- `<name>` - Name of query to display

**Options:**
- `--output=<format>` - Output format (table|json)

**Examples:**
```bash
# Show query details
7zarch-go query show my-docs

# JSON output
7zarch-go query show my-docs --output json
```

**Output Format:**
```
📋 Query Details: my-docs

Created: 2025-08-13 14:30:45
Last Used: 2025-08-13 15:45:12
Use Count: 5
Filters: --profile=documents --managed
```

### `query delete <name>`

Delete a saved query permanently.

**Usage:**
```bash
7zarch-go query delete <name>
```

**Arguments:**
- `<name>` - Name of query to delete

**Examples:**
```bash
# Delete a saved query
7zarch-go query delete old-query

# Confirmation message
✅ Query 'old-query' deleted successfully
```

## Integration with List Command

The `query` system integrates with the `list` command for seamless workflow:

### Saving Queries from List
```bash
# Test filters with list command
7zarch-go list --profile=media --larger-than=100MB --not-uploaded

# Save the filters once satisfied
7zarch-go list --profile=media --larger-than=100MB --not-uploaded --save-query large-media-pending
```

### Using Queries in List
```bash
# Use saved query in list command
7zarch-go list --query large-media-pending

# Combine query with additional filters
7zarch-go list --query large-media-pending --older-than=7d
```

## Integration with Search (7EP-0007 Phase 2)

Queries can include search terms for powerful combinations:

### Saving Search Queries
```bash
# Save search patterns
7zarch-go query save "photo-search" --search="vacation photos" --search-field=metadata
7zarch-go query save "code-backups" --search="backup" --search-field=name --profile=documents

# Save complex search + filter combinations
7zarch-go query save "old-media-search" --search="video" --search-field=name --profile=media --older-than=30d
```

### Executing Search Queries
```bash
# Run search-enabled queries
7zarch-go query run photo-search
7zarch-go query run code-backups --table

# Search queries work with all display modes
7zarch-go query run old-media-search --card --details
```

## Workflow Examples

### Daily Maintenance Workflows
```bash
# Set up maintenance queries
7zarch-go query save "daily-cleanup" --older-than=1y --not-uploaded
7zarch-go query save "upload-pending" --not-uploaded --larger-than=10MB
7zarch-go query save "health-check" --status=missing

# Daily execution
7zarch-go query run daily-cleanup
7zarch-go query run upload-pending
7zarch-go query run health-check
```

### Project-Specific Workflows
```bash
# Project organization
7zarch-go query save "current-project" --search="project-2024" --managed
7zarch-go query save "archived-projects" --search="project" --older-than=6m

# Review workflows
7zarch-go query run current-project --card
7zarch-go query run archived-projects --output json | jq length
```

### Batch Operation Preparation
```bash
# Prepare for batch operations
7zarch-go query save "batch-upload" --not-uploaded --profile=media
7zarch-go query save "batch-cleanup" --external --older-than=1y

# Use with future batch commands (Phase 3)
7zarch-go query run batch-upload --output json  # Prepare UIDs for batch processing
```

## Query Storage

Queries are stored in the SQLite database with the following information:

**Storage Schema:**
- `name` - Unique query name (primary key)
- `filters` - JSON-encoded filter configuration
- `created` - Creation timestamp
- `last_used` - Last execution timestamp
- `use_count` - Number of times executed

**Performance:**
- Query storage and retrieval is optimized for instant access
- Queries are sorted by most recently used first
- No limit on number of saved queries
- Compatible with database migration system

## Error Handling

Common error scenarios:

**Empty Query Name:**
```bash
7zarch-go query save "" --managed
# Error: query name cannot be empty
```

**Duplicate Query Name:**
```bash
7zarch-go query save "existing" --managed
# Overwrites existing query with same name
```

**Query Not Found:**
```bash
7zarch-go query run nonexistent
# Error: query not found: nonexistent
```

**No Filters Provided:**
```bash
7zarch-go query save "empty"
# Error: no filters specified - provide at least one filter flag
```

## Technical Implementation

The query system is implemented as part of 7EP-0007 Enhanced MAS Operations:

**Architecture:**
- SQLite storage using migration system (0004_query_system)
- JSON serialization of filter configurations
- Integration with existing resolver and display systems
- Thread-safe operations with proper error handling

**Database Migration:**
- Automatic schema creation on first use
- Compatible with existing database structure
- No impact on archive storage or retrieval

**Performance:**
- Query execution leverages existing list filtering logic
- Search integration uses high-performance search engine
- Memory-efficient storage and retrieval
- Optimized for frequent query execution

## Migration Notes

The query command was introduced in 7EP-0007 Phase 1 and requires no migration for existing installations. The query table is created automatically via the database migration system.

**Version Compatibility:**
- Requires 7zarch-go v0.3.0+ (7EP-0007 Phase 1)
- Search integration requires v0.3.1+ (7EP-0007 Phase 2)
- Compatible with all existing commands and workflows
- No changes required to existing data or configurations