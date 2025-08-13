# Database Migration Best Practices

This document outlines best practices for working with the database migration system in 7zarch-go.

## Overview

The migration system provides safe, transactional database schema evolution with automatic backups and rollback capability. It's implemented in the `internal/storage/migrations.go` file and accessed via the `7zarch-go db` CLI commands.

## CLI Commands

### Check Migration Status

```bash
# Show current schema version and pending migrations
7zarch-go db status
```

Output includes:
- Database file path and size
- Current schema version
- List of applied migrations with timestamps
- List of pending migrations

### Apply Migrations

```bash
# Apply all pending migrations with automatic backup
7zarch-go db migrate

# Dry run - show what would be applied without executing
7zarch-go db migrate --dry-run

# Create backup without applying migrations
7zarch-go db migrate --backup-only
```

### Create Manual Backups

```bash
# Create timestamped backup of the registry database
7zarch-go db backup
```

Backup files are created in the same directory as the database with format: `registry-YYYYMMDD-HHMMSS.bak`

## Migration Safety Features

### Automatic Backups
- Every migration run creates a timestamped backup before applying changes
- Backups are stored in the same directory as the database
- If migration fails, backup path is provided in the error message

### Transaction Safety
- Each migration runs in a database transaction
- Automatic rollback on any failure
- Schema validation after each migration

### Idempotent Operations
- Migrations can be run multiple times safely
- Schema introspection prevents duplicate column additions
- Migration tracking prevents re-applying completed migrations

## Writing New Migrations

### Migration Structure

Migrations are defined in the `MigrationRunner` struct and follow this pattern:

```go
const (
    migrationNewFeatureID   = "0004_new_feature"
    migrationNewFeatureName = "Add columns for new feature"
)
```

### Adding Migration Logic

1. **Add constants** for the new migration ID and name
2. **Update GetPendingMigrations()** to include the new migration
3. **Add migration logic** in `applyMigration()` method
4. **Test thoroughly** with both fresh and existing databases

### Best Practices for Schema Changes

#### Prefer Additive Changes
```go
// Good - additive changes are safer
ALTER TABLE archives ADD COLUMN new_field TEXT
CREATE INDEX IF NOT EXISTS idx_new_field ON archives(new_field)

// Avoid - dropping columns can cause data loss
ALTER TABLE archives DROP COLUMN old_field  // Don't do this
```

#### Use Safe Column Operations
```go
// Check if column exists before adding
if !columnExists(db, "archives", "new_column") {
    if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN new_column TEXT`); err != nil {
        _ = tx.Rollback()
        return fmt.Errorf("failed to add new_column: %w", err)
    }
}
```

#### Handle Complex Changes Carefully
For complex schema changes, use the copy-rebuild pattern:
1. Create new table with desired schema
2. Copy data from old table to new table
3. Drop old table and rename new table
4. Recreate indexes and triggers

## Testing Migrations

### Unit Tests
Every migration should have corresponding unit tests:

```go
func TestMigration_NewFeature(t *testing.T) {
    tmpDir := t.TempDir()
    dbPath := filepath.Join(tmpDir, "test.db")
    
    // Set up database with pre-migration state
    // Apply migration
    // Verify schema changes
    // Verify data integrity
}
```

### Integration Tests
Test the complete migration flow:
- Fresh database initialization
- Incremental migration from various starting points
- Rollback scenarios
- Large dataset performance

### Manual Testing
Before releasing a migration:
1. Test on a copy of production database
2. Verify backup/restore functionality
3. Test both fresh installation and upgrade paths
4. Performance test with realistic data volumes

## Recovery Procedures

### Migration Failure Recovery

If a migration fails:
1. The error message will include the backup path
2. Stop the application to prevent further damage
3. Restore from the backup:
   ```bash
   cp /path/to/backup.bak ~/.7zarch-go/registry.db
   ```
4. Investigate the failure cause
5. Fix the migration code and retry

### Database Corruption Recovery

If the database becomes corrupted:
1. Check for recent backups in the database directory
2. Restore from the most recent backup
3. If no backups available, you may need to reinitialize:
   ```bash
   # Backup current corrupted database first
   mv ~/.7zarch-go/registry.db ~/.7zarch-go/registry.db.corrupted
   
   # Let the application recreate the database
   7zarch-go list
   ```

### Manual Schema Fixes

If you need to manually fix schema issues:
1. Create a backup first: `7zarch-go db backup`
2. Use SQLite tools to inspect and modify:
   ```bash
   sqlite3 ~/.7zarch-go/registry.db
   .schema archives
   .quit
   ```
3. Update migration tracking if needed:
   ```sql
   INSERT INTO schema_migrations (id, name, applied_at) 
   VALUES ('manual_fix', 'Manual schema fix', datetime('now'));
   ```

## Performance Considerations

### Migration Performance
- Migrations should complete within 30 seconds (configurable timeout)
- For large datasets, consider batch processing
- Add progress indicators for long-running migrations
- Test performance with realistic data volumes

### Database Locking
- Migrations use `BEGIN IMMEDIATE` to acquire exclusive lock
- Other operations will wait during migration
- Keep migration duration minimal to reduce blocking time

### Backup Performance
- Backup time scales with database size
- Large databases (>100MB) may take several seconds
- Consider cleanup of old backups to manage disk space

## Development Workflow

### Adding a New Migration

1. **Design the change**
   - Document what tables/columns will be affected
   - Plan for data preservation and type conversion
   - Consider backward compatibility

2. **Implement the migration**
   - Add constants for the migration ID and name
   - Update `GetPendingMigrations()` to include it
   - Add logic to `applyMigration()`
   - Use proper transaction handling and error checking

3. **Test thoroughly**
   - Unit tests for the migration logic
   - Integration tests with sample data
   - Manual testing on development databases

4. **Document the change**
   - Update CHANGELOG.md with migration details
   - Add any special upgrade notes
   - Update version compatibility information

### Testing Against Production Data

1. **Create a copy** of production database
2. **Test migration** on the copy
3. **Verify results** - schema and data integrity
4. **Measure performance** - time and resource usage
5. **Test rollback** procedures

## Troubleshooting

### Common Issues

**"Migration already applied" error**
- Check `7zarch-go db status` to see applied migrations
- Verify the migration ID is unique and correct

**"Column already exists" error**
- Use `columnExists()` helper function
- Make migrations idempotent

**"Database is locked" error**
- Ensure no other 7zarch-go processes are running
- Check for long-running queries or connections

**Slow migration performance**
- Check database size and table row counts
- Consider adding progress indicators
- For large datasets, implement batch processing

### Debug Information

Enable verbose logging for migration debugging:
```bash
# Set debug environment variable
export ZARCH_LOG_LEVEL=debug
7zarch-go db migrate
```

Check SQLite database state:
```bash
sqlite3 ~/.7zarch-go/registry.db
.tables
.schema archives
SELECT * FROM schema_migrations;
.quit
```

## Future Enhancements

The migration system may be enhanced with:
- **Rollback migrations** - Down migrations for version downgrades
- **Migration validation** - Schema integrity checks and repair
- **Distributed coordination** - For multiple registry instances
- **Performance monitoring** - Migration timing and resource usage
- **Migration dependencies** - Complex migration ordering

## References

- [7EP-0003: Database Migrations & Schema Management](../7eps/7ep-0003-database-migrations.md)
- [7EP-0014: Critical Foundation Gaps](../7eps/7ep-0014-critical-foundation-gaps.md)
- [SQLite ALTER TABLE documentation](https://www.sqlite.org/lang_altertable.html)
- [SQLite transaction documentation](https://www.sqlite.org/lang_transaction.html)
