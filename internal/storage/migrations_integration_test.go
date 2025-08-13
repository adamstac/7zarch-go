package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrationRunner_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a database with only the baseline schema
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Create baseline schema manually (simulating old database)
	_, err = db.Exec(`CREATE TABLE archives (
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
	)`)
	if err != nil {
		t.Fatalf("failed to create baseline schema: %v", err)
	}

	// Add some test data
	_, err = db.Exec(`INSERT INTO archives (name, path, size) VALUES ('test1.7z', '/path/test1.7z', 1024)`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	runner := NewMigrationRunner(db, dbPath)
	registry := &Registry{db: db}

	// Set up migrations table and record what should be "applied" based on existing schema
	if err := registry.EnsureMigrationsTable(); err != nil {
		t.Fatalf("failed to ensure migrations table: %v", err)
	}

	// Mark baseline as applied (since the table already exists)
	if err := registry.MarkMigrationApplied(migrationBaselineID, migrationBaselineName); err != nil {
		t.Fatalf("failed to mark baseline migration as applied: %v", err)
	}

	// Check initial pending migrations
	pending, err := runner.GetPendingMigrations()
	if err != nil {
		t.Fatalf("failed to get pending migrations: %v", err)
	}

	// Should have identity and trash migrations pending
	if len(pending) < 1 {
		t.Fatalf("expected at least 1 pending migration, got %d", len(pending))
	}

	// Apply pending migrations
	backupPath, err := runner.CreateBackup(dbPath)
	if err != nil {
		t.Fatalf("failed to create backup: %v", err)
	}

	// Apply the trash migration manually to test the process
	trashApplied, err := registry.IsMigrationApplied(migrationTrashID)
	if err != nil {
		t.Fatalf("failed to check if trash migration applied: %v", err)
	}

	if !trashApplied {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}

		if !columnExists(db, "archives", "deleted_at") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN deleted_at TIMESTAMP`); err != nil {
				_ = tx.Rollback()
				t.Fatalf("failed to add deleted_at column: %v", err)
			}
		}

		if !columnExists(db, "archives", "original_path") {
			if _, err := tx.Exec(`ALTER TABLE archives ADD COLUMN original_path TEXT`); err != nil {
				_ = tx.Rollback()
				t.Fatalf("failed to add original_path column: %v", err)
			}
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit migration: %v", err)
		}

		if err := registry.MarkMigrationApplied(migrationTrashID, migrationTrashName); err != nil {
			t.Fatalf("failed to mark trash migration as applied: %v", err)
		}
	}

	// Verify columns were added
	if !columnExists(db, "archives", "deleted_at") {
		t.Fatal("deleted_at column not found after migration")
	}

	if !columnExists(db, "archives", "original_path") {
		t.Fatal("original_path column not found after migration")
	}

	// Verify data was preserved
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM archives`).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count archives after migration: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 archive record after migration, got %d", count)
	}

	// Verify backup was created
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatalf("backup file does not exist: %s", backupPath)
	}

	// Check that no migrations are pending now
	pending, err = runner.GetPendingMigrations()
	if err != nil {
		t.Fatalf("failed to get pending migrations after applying: %v", err)
	}

	if len(pending) != 0 {
		t.Fatalf("expected no pending migrations, got %d", len(pending))
	}

	// Verify applied migrations
	applied, err := runner.GetAppliedMigrations()
	if err != nil {
		t.Fatalf("failed to get applied migrations: %v", err)
	}

	expectedMigrations := []string{migrationBaselineID, migrationTrashID}
	if len(applied) < len(expectedMigrations) {
		t.Fatalf("expected at least %d applied migrations, got %d", len(expectedMigrations), len(applied))
	}

	// Check that expected migrations are applied
	appliedIds := make(map[string]bool)
	for _, migration := range applied {
		appliedIds[migration.ID] = true
	}

	for _, expectedId := range expectedMigrations {
		if !appliedIds[expectedId] {
			t.Fatalf("expected migration %s to be applied", expectedId)
		}
	}
}
