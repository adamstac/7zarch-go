package storage

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	migrationBaselineID   = "0001_baseline"
	migrationBaselineName = "Baseline schema"

	migrationIdentityID   = "0002_identity_and_status"
	migrationIdentityName = "Add uid/managed/status/last_seen and indexes"

	migrationTrashID   = "0003_trash_fields"
	migrationTrashName = "Add deleted_at and original_path for trash support"

	migrationQueryID   = "0004_query_system"
	migrationQueryName = "Add queries table for saved query support"

	migrationSearchID   = "0005_search_index"
	migrationSearchName = "Add search_index table for full-text search support"
)

type MigrationRunner struct {
	db         *sql.DB
	backupPath string
	timeout    time.Duration
}

func NewMigrationRunner(db *sql.DB, dbPath string) *MigrationRunner {
	backupDir := filepath.Dir(dbPath)
	return &MigrationRunner{
		db:         db,
		backupPath: backupDir,
		timeout:    30 * time.Second,
	}
}

type PendingMigration struct {
	ID          string
	Name        string
	Description string
}

func (mr *MigrationRunner) GetPendingMigrations() ([]PendingMigration, error) {
	registry := &Registry{db: mr.db}
	if err := registry.EnsureMigrationsTable(); err != nil {
		return nil, err
	}

	var pending []PendingMigration

	applied, err := registry.IsMigrationApplied(migrationTrashID)
	if err != nil {
		return nil, err
	}
	if !applied {
		pending = append(pending, PendingMigration{
			ID:          migrationTrashID,
			Name:        migrationTrashName,
			Description: "Adds deleted_at and original_path columns for trash functionality",
		})
	}

	applied, err = registry.IsMigrationApplied(migrationQueryID)
	if err != nil {
		return nil, err
	}
	if !applied {
		pending = append(pending, PendingMigration{
			ID:          migrationQueryID,
			Name:        migrationQueryName,
			Description: "Adds queries table for saved query functionality",
		})
	}

	applied, err = registry.IsMigrationApplied(migrationSearchID)
	if err != nil {
		return nil, err
	}
	if !applied {
		pending = append(pending, PendingMigration{
			ID:          migrationSearchID,
			Name:        migrationSearchName,
			Description: "Adds search_index table for full-text search functionality",
		})
	}

	return pending, nil
}

func (mr *MigrationRunner) CreateBackup(dbPath string) (string, error) {
	if dbPath == "" {
		return "", fmt.Errorf("database path is required for backup")
	}

	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("registry-%s.bak", timestamp)
	backupPath := filepath.Join(mr.backupPath, backupName)

	src, err := os.Open(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source database: %w", err)
	}
	defer src.Close()

	if err := os.MkdirAll(mr.backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	dst, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy database: %w", err)
	}

	return backupPath, nil
}

func (mr *MigrationRunner) ApplyPending(dbPath string) error {
	pending, err := mr.GetPendingMigrations()
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(pending) == 0 {
		return nil
	}

	backupPath, err := mr.CreateBackup(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create backup before migration: %w", err)
	}

	registry := &Registry{db: mr.db}

	for _, migration := range pending {
		if err := mr.applyMigration(registry, migration); err != nil {
			return fmt.Errorf("migration %s failed: %w\nBackup preserved at: %s", migration.ID, err, backupPath)
		}
	}

	return nil
}

func (mr *MigrationRunner) applyMigration(registry *Registry, migration PendingMigration) error {
	tx, err := mr.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	switch migration.ID {
	case migrationTrashID:
		if !columnExists(mr.db, "archives", "deleted_at") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN deleted_at TIMESTAMP`); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to add deleted_at column: %w", err)
			}
		}
		if !columnExists(mr.db, "archives", "original_path") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN original_path TEXT`); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to add original_path column: %w", err)
			}
		}
	case migrationQueryID:
		if !tableExists(mr.db, "queries") {
			if _, err := tx.Exec(`
				CREATE TABLE queries (
					name TEXT PRIMARY KEY,
					filters TEXT NOT NULL,
					created INTEGER NOT NULL,
					last_used INTEGER,
					use_count INTEGER DEFAULT 0
				)
			`); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to create queries table: %w", err)
			}
		}
	case migrationSearchID:
		if !tableExists(mr.db, "search_index") {
			if _, err := tx.Exec(`
				CREATE TABLE search_index (
					term TEXT,
					archive_uid TEXT,
					field TEXT,
					PRIMARY KEY (term, archive_uid, field)
				)
			`); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to create search_index table: %w", err)
			}
			// Create performance index
			if _, err := tx.Exec(`
				CREATE INDEX idx_search_term ON search_index(term)
			`); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to create search index: %w", err)
			}
		}
	default:
		_ = tx.Rollback()
		return fmt.Errorf("unknown migration: %s", migration.ID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	if err := registry.MarkMigrationApplied(migration.ID, migration.Name); err != nil {
		return fmt.Errorf("failed to mark migration as applied: %w", err)
	}

	return nil
}

type AppliedMigration struct {
	ID        string
	Name      string
	AppliedAt time.Time
}

func (mr *MigrationRunner) GetAppliedMigrations() ([]AppliedMigration, error) {
	registry := &Registry{db: mr.db}
	if err := registry.EnsureMigrationsTable(); err != nil {
		return nil, err
	}

	rows, err := mr.db.Query(`SELECT id, name, applied_at FROM schema_migrations ORDER BY applied_at`)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var applied []AppliedMigration
	for rows.Next() {
		var migration AppliedMigration
		if err := rows.Scan(&migration.ID, &migration.Name, &migration.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		applied = append(applied, migration)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migration rows: %w", err)
	}

	return applied, nil
}

func (m *Manager) NewMigrationRunner() *MigrationRunner {
	return NewMigrationRunner(m.registry.db, m.registry.Path())
}

// EnsureMigrationsTable creates the schema_migrations table if missing
func (r *Registry) EnsureMigrationsTable() error {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL
	)`)
	return err
}

// IsMigrationApplied checks if a migration id is recorded
func (r *Registry) IsMigrationApplied(id string) (bool, error) {
	row := r.db.QueryRow(`SELECT 1 FROM schema_migrations WHERE id = ?`, id)
	var one int
	switch err := row.Scan(&one); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// MarkMigrationApplied records a migration as applied
func (r *Registry) MarkMigrationApplied(id, name string) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO schema_migrations (id, name, applied_at) VALUES (?, ?, ?)`, id, name, time.Now())
	return err
}

// RecordCurrentSchemaAsApplied marks baseline and identity migrations as applied
// if the current schema already contains their effects (idempotent).
func (r *Registry) RecordCurrentSchemaAsApplied() error {
	if err := r.EnsureMigrationsTable(); err != nil {
		return err
	}
	// If archives table exists, we can mark baseline applied
	if tableExists(r.db, "archives") {
		applied, err := r.IsMigrationApplied(migrationBaselineID)
		if err != nil {
			return err
		}
		if !applied {
			if err := r.MarkMigrationApplied(migrationBaselineID, migrationBaselineName); err != nil {
				return err
			}
		}
	}
	// If identity columns exist, mark identity migration applied
	if columnExists(r.db, "archives", "uid") && columnExists(r.db, "archives", "managed") && columnExists(r.db, "archives", "status") && columnExists(r.db, "archives", "last_seen") {
		applied, err := r.IsMigrationApplied(migrationIdentityID)
		if err != nil {
			return err
		}
		if !applied {
			if err := r.MarkMigrationApplied(migrationIdentityID, migrationIdentityName); err != nil {
				return err
			}
		}
	}
	return nil
}

// ApplyPendingMigrations runs known migrations that haven't been marked applied yet
func (r *Registry) ApplyPendingMigrations() error {
	if err := r.EnsureMigrationsTable(); err != nil {
		return err
	}
	// 0001 & 0002 are effectively applied via initSchema + RecordCurrentSchemaAsApplied
	// 0003: trash support
	applied, err := r.IsMigrationApplied(migrationTrashID)
	if err != nil {
		return err
	}
	if !applied {
		tx, err := r.db.Begin()
		if err != nil {
			return err
		}
		// Add columns if missing
		if !columnExists(r.db, "archives", "deleted_at") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN deleted_at TIMESTAMP`); err != nil {
				_ = tx.Rollback() // best-effort rollback
				return err
			}
		}
		if !columnExists(r.db, "archives", "original_path") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN original_path TEXT`); err != nil {
				_ = tx.Rollback() // best-effort rollback
				return err
			}
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		if err := r.MarkMigrationApplied(migrationTrashID, migrationTrashName); err != nil {
			return err
		}
	}
	
	// 0004: query system support
	applied, err = r.IsMigrationApplied(migrationQueryID)
	if err != nil {
		return err
	}
	if !applied {
		// Create queries table if missing
		if !tableExists(r.db, "queries") {
			if _, err := r.db.Exec(`
				CREATE TABLE queries (
					name TEXT PRIMARY KEY,
					filters TEXT NOT NULL,
					created INTEGER NOT NULL,
					last_used INTEGER,
					use_count INTEGER DEFAULT 0
				)
			`); err != nil {
				return err
			}
		}
		if err := r.MarkMigrationApplied(migrationQueryID, migrationQueryName); err != nil {
			return err
		}
	}
	
	// 0005: search index support
	applied, err = r.IsMigrationApplied(migrationSearchID)
	if err != nil {
		return err
	}
	if !applied {
		// Create search_index table if missing
		if !tableExists(r.db, "search_index") {
			if _, err := r.db.Exec(`
				CREATE TABLE search_index (
					term TEXT,
					archive_uid TEXT,
					field TEXT,
					PRIMARY KEY (term, archive_uid, field)
				)
			`); err != nil {
				return err
			}
			// Create performance index
			if _, err := r.db.Exec(`
				CREATE INDEX idx_search_term ON search_index(term)
			`); err != nil {
				return err
			}
		}
		if err := r.MarkMigrationApplied(migrationSearchID, migrationSearchName); err != nil {
			return err
		}
	}
	return nil
}

func tableExists(db *sql.DB, table string) bool {
	row := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table)
	var name string
	if err := row.Scan(&name); err != nil {
		return false
	}
	return name == table
}

func columnExists(db *sql.DB, table, column string) bool {
	rows, err := db.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, table))
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return false
		}
		if name == column {
			return true
		}
	}
	return false
}
