package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrationRunner_CreateBackup(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Create a simple table with data
	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO test (name) VALUES ('test data')`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	runner := NewMigrationRunner(db, dbPath)
	backupPath, err := runner.CreateBackup(dbPath)
	if err != nil {
		t.Fatalf("backup failed: %v", err)
	}

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatalf("backup file does not exist: %s", backupPath)
	}

	// Verify backup has correct content
	backupDB, err := sql.Open("sqlite3", backupPath)
	if err != nil {
		t.Fatalf("failed to open backup database: %v", err)
	}
	defer backupDB.Close()

	var count int
	err = backupDB.QueryRow(`SELECT COUNT(*) FROM test`).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query backup database: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 row in backup, got %d", count)
	}
}

func TestMigrationRunner_GetPendingMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	runner := NewMigrationRunner(db, dbPath)
	registry := &Registry{db: db}

	// Initialize migrations table
	if err := registry.EnsureMigrationsTable(); err != nil {
		t.Fatalf("failed to ensure migrations table: %v", err)
	}

	pending, err := runner.GetPendingMigrations()
	if err != nil {
		t.Fatalf("failed to get pending migrations: %v", err)
	}

	// Should have trash migration pending
	found := false
	for _, migration := range pending {
		if migration.ID == migrationTrashID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected trash migration to be pending")
	}

	// Apply the migration
	if err := registry.MarkMigrationApplied(migrationTrashID, migrationTrashName); err != nil {
		t.Fatalf("failed to mark migration as applied: %v", err)
	}

	// Check pending migrations again
	pending, err = runner.GetPendingMigrations()
	if err != nil {
		t.Fatalf("failed to get pending migrations after applying: %v", err)
	}

	// Should no longer have trash migration pending
	for _, migration := range pending {
		if migration.ID == migrationTrashID {
			t.Fatalf("trash migration should not be pending after applying")
		}
	}
}

func TestMigrationRunner_GetAppliedMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	runner := NewMigrationRunner(db, dbPath)
	registry := &Registry{db: db}

	// Initialize migrations table
	if err := registry.EnsureMigrationsTable(); err != nil {
		t.Fatalf("failed to ensure migrations table: %v", err)
	}

	// Apply a test migration
	if err := registry.MarkMigrationApplied("test_001", "Test migration"); err != nil {
		t.Fatalf("failed to mark migration as applied: %v", err)
	}

	applied, err := runner.GetAppliedMigrations()
	if err != nil {
		t.Fatalf("failed to get applied migrations: %v", err)
	}

	if len(applied) != 1 {
		t.Fatalf("expected 1 applied migration, got %d", len(applied))
	}

	migration := applied[0]
	if migration.ID != "test_001" {
		t.Fatalf("expected migration ID 'test_001', got '%s'", migration.ID)
	}

	if migration.Name != "Test migration" {
		t.Fatalf("expected migration name 'Test migration', got '%s'", migration.Name)
	}

	// Verify timestamp is recent
	if time.Since(migration.AppliedAt) > time.Minute {
		t.Fatalf("migration timestamp seems too old: %v", migration.AppliedAt)
	}
}
