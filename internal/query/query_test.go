package query

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, string) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create schema_migrations table
	_, err = db.Exec(`CREATE TABLE schema_migrations (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at INTEGER NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	// Create basic archives table for resolver
	_, err = db.Exec(`CREATE TABLE archives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT UNIQUE NOT NULL,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		size INTEGER NOT NULL,
		created INTEGER NOT NULL,
		checksum TEXT,
		profile TEXT,
		managed BOOLEAN DEFAULT true,
		status TEXT DEFAULT 'present',
		last_seen INTEGER,
		deleted_at INTEGER,
		original_path TEXT,
		uploaded BOOLEAN DEFAULT false,
		destination TEXT,
		uploaded_at INTEGER,
		metadata TEXT
	)`)
	if err != nil {
		t.Fatalf("Failed to create archives table: %v", err)
	}

	return db, dbPath
}

func TestQueryManager_SaveAndList(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	// Create registry and resolver for the query manager
	registry := &storage.Registry{}
	resolver := storage.NewResolver(registry)
	
	// Create query manager
	qm := NewQueryManager(db, resolver)

	// Test saving a query
	filters := map[string]string{
		"managed": "true",
		"profile": "documents",
	}

	err := qm.Save("test-query", filters)
	if err != nil {
		t.Fatalf("Failed to save query: %v", err)
	}

	// Test listing queries
	queries, err := qm.List()
	if err != nil {
		t.Fatalf("Failed to list queries: %v", err)
	}

	if len(queries) != 1 {
		t.Fatalf("Expected 1 query, got %d", len(queries))
	}

	query := queries[0]
	if query.Name != "test-query" {
		t.Errorf("Expected query name 'test-query', got '%s'", query.Name)
	}

	if len(query.Filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(query.Filters))
	}

	if query.Filters["managed"] != "true" {
		t.Errorf("Expected managed filter to be 'true', got '%s'", query.Filters["managed"])
	}

	if query.UseCount != 0 {
		t.Errorf("Expected use count 0, got %d", query.UseCount)
	}
}

func TestQueryManager_GetAndDelete(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	registry := &storage.Registry{}
	resolver := storage.NewResolver(registry)
	qm := NewQueryManager(db, resolver)

	// Save a query
	filters := map[string]string{"status": "present"}
	err := qm.Save("delete-test", filters)
	if err != nil {
		t.Fatalf("Failed to save query: %v", err)
	}

	// Test getting the query
	query, err := qm.Get("delete-test")
	if err != nil {
		t.Fatalf("Failed to get query: %v", err)
	}

	if query.Name != "delete-test" {
		t.Errorf("Expected query name 'delete-test', got '%s'", query.Name)
	}

	// Test deleting the query
	err = qm.Delete("delete-test")
	if err != nil {
		t.Fatalf("Failed to delete query: %v", err)
	}

	// Verify it's gone
	_, err = qm.Get("delete-test")
	if err == nil {
		t.Error("Expected error when getting deleted query")
	}
}

func TestQueryManager_SaveEmpty(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	registry := &storage.Registry{}
	resolver := storage.NewResolver(registry)
	qm := NewQueryManager(db, resolver)

	// Test saving empty query name
	err := qm.Save("", map[string]string{"test": "value"})
	if err == nil {
		t.Error("Expected error when saving query with empty name")
	}
}

func TestApplyQueryMigration(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	// Apply the migration
	err := ApplyQueryMigration(db)
	if err != nil {
		t.Fatalf("Failed to apply query migration: %v", err)
	}

	// Verify the queries table was created
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='queries'")
	var tableName string
	err = row.Scan(&tableName)
	if err != nil {
		t.Fatalf("Queries table was not created: %v", err)
	}

	if tableName != "queries" {
		t.Errorf("Expected table name 'queries', got '%s'", tableName)
	}

	// Verify the migration was recorded
	applied, err := IsQueryMigrationApplied(db)
	if err != nil {
		t.Fatalf("Failed to check migration status: %v", err)
	}

	if !applied {
		t.Error("Migration was not recorded as applied")
	}
}

func TestTimeHandling(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	registry := &storage.Registry{}
	resolver := storage.NewResolver(registry)
	qm := NewQueryManager(db, resolver)

	// Save a query
	filters := map[string]string{"test": "value"}
	err := qm.Save("time-test", filters)
	if err != nil {
		t.Fatalf("Failed to save query: %v", err)
	}

	// Get the query and verify time handling
	query, err := qm.Get("time-test")
	if err != nil {
		t.Fatalf("Failed to get query: %v", err)
	}

	// Check that created time is reasonable (within last minute)
	now := time.Now()
	if query.Created.After(now) || query.Created.Before(now.Add(-time.Minute)) {
		t.Errorf("Created time seems wrong: %v", query.Created)
	}

	// Check that LastUsed is nil for new query
	if query.LastUsed != nil {
		t.Errorf("Expected LastUsed to be nil for new query, got %v", query.LastUsed)
	}
}