# 7EP-0003: Database Migrations & Schema Management

**Status:** Draft  
**Author(s):** Claude Code (CC), Augment Code (AC)  
**Assignment:** AC  
**Difficulty:** 2 (moderate - safety-critical but well-understood pattern)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  

## Executive Summary

Implement a lightweight, built-in migration system for the MAS registry to safely evolve schema over time with automatic backups, transactional safety, and clear observability.

## Evidence & Reasoning

**User feedback/pain points:**
- Database schema conflicts when upgrading versions
- Fear of data loss during schema changes
- Need manual intervention for registry updates
- Unclear state when migrations fail

**Current limitations:**
- No formal migration system for registry schema
- Schema changes require manual database deletion
- No backup/recovery mechanism for failed updates
- No visibility into current schema version

**Why now:**
- Registry is becoming central to all operations
- Multiple features (trash, ULID, status tracking) require schema changes
- User data safety is critical for adoption
- Foundation needed before more complex features

## Use Cases

### Primary Use Case: Safe Version Upgrades
```bash
# User upgrades 7zarch-go version with schema changes
7zarch-go list
# Output: "Applying 2 pending migrations... ✓ Complete"
# Registry automatically updated without data loss
```

### Secondary Use Cases
- **Development iterations**: Schema changes during feature development
- **Recovery scenarios**: Rollback failed migrations with backups
- **Observability**: Check current schema version and migration status
- **Manual control**: Run migrations explicitly when needed

## Technical Design

### Overview
Lightweight migration system with safety, idempotence, transactional updates, observability, and minimal user friction.

### Migration Strategy
- Maintain `schema_migrations` table (id, name, applied_at)
- Run pending migrations in order at startup (auto-migrate) or via CLI
- Each migration is transactional; rollback on failure
- Prefer additive changes (ALTER TABLE ADD COLUMN, CREATE INDEX IF NOT EXISTS)
- For complex changes, use copy-rebuild pattern (new table, copy/transform, swap)

### Safety Mechanisms
- Take timestamped backup before migrating
- Lock DB with `BEGIN IMMEDIATE` during migrations
- Validate migration integrity before applying
- Automatic rollback on any failure

### Migration File Structure
```
internal/storage/migrations/
├── 0001_baseline.sql           # Initial archives schema + indexes
├── 0002_identity_and_status.sql # Add uid, managed, status, last_seen
├── 0003_trash_management.sql   # Add deleted_at, original_path
└── migration.go               # Migration runner logic
```

### CLI Commands

#### `7zarch-go db status`
```bash
# Shows current migration state
Database: ~/.7zarch-go/registry.db
Schema Version: 0002_identity_and_status
Applied Migrations: 2/3
Pending Migrations:
  - 0003_trash_management (adds deleted_at, original_path columns)
Database Size: 42.1 KB
Last Backup: 2025-08-12 10:30:22
```

#### `7zarch-go db migrate [--dry-run] [--backup-only]`
```bash
# Apply pending migrations
7zarch-go db migrate --dry-run
# Shows what would be applied without executing

7zarch-go db migrate
# Creates backup, applies migrations, reports success
```

#### `7zarch-go db backup`
```bash
# Create timestamped backup
7zarch-go db backup
# Output: Backup created: ~/.7zarch-go/backups/registry-20250812-103022.db
```

### Configuration Integration
```yaml
storage:
  auto_migrate: true              # Run migrations on startup (default)
  backup_before_migrate: true     # Create backup before migrations
  migration_timeout: 30s          # Timeout for individual migrations
  
# Environment override for automation
# ZARCH_DB_AUTOMIGRATE=false
```

### Migration Examples

#### Migration 0001: Baseline
```sql
-- 0001_baseline.sql
CREATE TABLE IF NOT EXISTS archives (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    size INTEGER NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checksum TEXT,
    profile TEXT,
    uploaded BOOLEAN DEFAULT FALSE,
    destination TEXT,
    uploaded_at TIMESTAMP,
    metadata TEXT
);

CREATE INDEX IF NOT EXISTS idx_archives_name ON archives(name);
CREATE INDEX IF NOT EXISTS idx_archives_checksum ON archives(checksum);
```

#### Migration 0002: Identity and Status
```sql
-- 0002_identity_and_status.sql
ALTER TABLE archives ADD COLUMN uid TEXT;
ALTER TABLE archives ADD COLUMN managed BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE archives ADD COLUMN status TEXT NOT NULL DEFAULT 'present';
ALTER TABLE archives ADD COLUMN last_seen TIMESTAMP;

-- Create unique index after backfilling UIDs
CREATE UNIQUE INDEX IF NOT EXISTS idx_archives_uid ON archives(uid);
CREATE INDEX IF NOT EXISTS idx_archives_status ON archives(status);

-- Backfill UIDs for existing records (handled in Go migration code)
```

### Implementation Plan

#### Phase 1: Migration Infrastructure
- [ ] Create migration table and runner
- [ ] Implement backup system with timestamped files
- [ ] Add CLI commands for status, migrate, backup
- [ ] Integrate with startup sequence

#### Phase 2: Safety & Validation
- [ ] Add transaction safety and rollback
- [ ] Implement migration validation
- [ ] Add timeout and error handling
- [ ] Create comprehensive logging

#### Phase 3: Testing & Edge Cases
- [ ] Test with corrupted databases
- [ ] Validate rollback scenarios
- [ ] Test large database migrations
- [ ] Cross-platform compatibility

### Error Handling

#### Migration Failure Recovery
```bash
# If migration fails:
Error: Migration 0003_trash_management failed: column already exists
Backup preserved: ~/.7zarch-go/backups/registry-20250812-103022.db
To recover: cp ~/.7zarch-go/backups/registry-20250812-103022.db ~/.7zarch-go/registry.db
```

#### Corruption Detection
```bash
# Database corruption handling:
7zarch-go db status
# Output: Warning: Database corruption detected
#         Run '7zarch-go db recover' to attempt repair
```

## Testing Strategy

### Acceptance Criteria
- [ ] Migrations apply successfully on clean database
- [ ] Existing data preserved through schema changes
- [ ] Failed migrations rollback without data loss
- [ ] Backups can restore complete registry state
- [ ] Auto-migration works on startup
- [ ] Manual migration control available

### Test Scenarios
- Fresh installation (baseline migration)
- Incremental upgrades (apply pending migrations)
- Failed migration rollback
- Corrupt database recovery
- Large database performance
- Cross-platform compatibility

### Testing With Fixtures
```go
// Unit tests with versioned fixtures
func TestMigration_0001_to_0002(t *testing.T) {
    // Load 0001 database fixture
    db := loadFixture("0001_baseline.db")
    
    // Apply 0002 migration
    err := applyMigration(db, "0002_identity_and_status")
    require.NoError(t, err)
    
    // Verify schema and data integrity
    verifyColumns(t, db, []string{"uid", "managed", "status"})
    verifyIndexes(t, db, []string{"idx_archives_uid", "idx_archives_status"})
    verifyDataIntegrity(t, db)
}
```

## Migration/Compatibility

### Breaking Changes
None - migrations are additive and backward-compatible.

### Upgrade Path
Automatic migration on first run of new version.

### Backward Compatibility
New columns use sensible defaults; existing queries continue working.

## Alternatives Considered

**External migration tools**: Considered tools like golang-migrate but decided on built-in system for simpler deployment and registry-specific optimizations.

**Version-specific databases**: Evaluated separate DB files per version but decided unified approach with migrations provides better user experience.

**Manual migration scripts**: Considered requiring manual schema updates but automatic migration reduces user friction significantly.

## Future Considerations

- **Schema validation**: Automated schema integrity checks
- **Migration rollback**: Reverse migrations for downgrades
- **Distributed migrations**: Coordination across multiple registry instances
- **Performance optimization**: Lazy migrations for large datasets

## References

- Related: 7EP-0001 Trash Management (requires deleted_at, original_path columns)
- Related: ULID implementation (requires uid column and indexes)
- Database backup strategy documented in DEVELOPMENT.md