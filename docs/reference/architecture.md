# Architecture Decisions

## Managed Archive Storage (MAS)

- **Database**: SQLite with 0600 permissions
- **Location**: `~/.7zarch-go/`
- **Registry**: Tracks metadata, status, relationships
- **Migration**: Automatic schema updates

## Command Design Philosophy

- **Top-level commands** preferred over deep subcommand nesting
- **ULID resolution** for user-friendly ID references
- **Status indicators** in list outputs for clear visual feedback
- **Configuration integration** - no hardcoded values

## Display Standards

```bash
# Status formatting
📦 archive-name - ✅ Uploaded to destination
🗑️  deleted-name - Deleted 2025-08-11 14:23:01
   Auto-purge: 29 days (2025-09-10)

# ID display for copy/paste
   ID: 01K2E33XW4HTX7RVPS9Y6CRGDY
```

## Core Principles

### Simplicity First
- Prefer straightforward implementations over complex abstractions
- Commands should be discoverable and intuitive
- Error messages should be helpful and actionable

### User Experience
- Fast response times (<100ms for common operations)
- Clear visual feedback for all operations
- Progressive disclosure of complexity (simple defaults, advanced options)

### Data Safety
- Never delete data without explicit user confirmation
- Soft deletes with recovery period (trash system)
- Automatic backups before migrations

### Extensibility
- Plugin-friendly architecture for future extensions
- Clear interfaces between components
- Configuration-driven behavior