package storage

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	migrationBaselineID   = "0001_baseline"
	migrationBaselineName = "Baseline schema"

	migrationIdentityID   = "0002_identity_and_status"
	migrationIdentityName = "Add uid/managed/status/last_seen and indexes"
)

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
	if err := r.EnsureMigrationsTable(); err != nil { return err }
	// If archives table exists, we can mark baseline applied
	if tableExists(r.db, "archives") {
		applied, err := r.IsMigrationApplied(migrationBaselineID)
		if err != nil { return err }
		if !applied {
			if err := r.MarkMigrationApplied(migrationBaselineID, migrationBaselineName); err != nil { return err }
		}
	}
	// If identity columns exist, mark identity migration applied
	if columnExists(r.db, "archives", "uid") && columnExists(r.db, "archives", "managed") && columnExists(r.db, "archives", "status") && columnExists(r.db, "archives", "last_seen") {
		applied, err := r.IsMigrationApplied(migrationIdentityID)
		if err != nil { return err }
		if !applied {
			if err := r.MarkMigrationApplied(migrationIdentityID, migrationIdentityName); err != nil { return err }
		}
	}
	return nil
}

const (
	migrationTrashID   = "0003_trash_fields"
	migrationTrashName = "Add deleted_at and original_path for trash support"
)

// ApplyPendingMigrations runs known migrations that haven't been marked applied yet
func (r *Registry) ApplyPendingMigrations() error {
	if err := r.EnsureMigrationsTable(); err != nil { return err }
	// 0001 & 0002 are effectively applied via initSchema + RecordCurrentSchemaAsApplied
	// 0003: trash support
	applied, err := r.IsMigrationApplied(migrationTrashID)
	if err != nil { return err }
	if !applied {
		tx, err := r.db.Begin()
		if err != nil { return err }
		// Add columns if missing
		if !columnExists(r.db, "archives", "deleted_at") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN deleted_at TIMESTAMP`); err != nil { tx.Rollback(); return err }
		}
		if !columnExists(r.db, "archives", "original_path") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN original_path TEXT`); err != nil { tx.Rollback(); return err }
		}
		if err := tx.Commit(); err != nil { return err }
		if err := r.MarkMigrationApplied(migrationTrashID, migrationTrashName); err != nil { return err }
	}
	return nil
}

func tableExists(db *sql.DB, table string) bool {
	row := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, table)
	var name string
	if err := row.Scan(&name); err != nil { return false }
	return name == table
}

func columnExists(db *sql.DB, table, column string) bool {
	rows, err := db.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, table))
	if err != nil { return false }
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil { return false }
		if name == column { return true }
	}
	return false
}

