package search

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func setupSearchTestDB(t *testing.T) (*sql.DB, *storage.Registry) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	// Use proper Registry constructor which handles all initialization
	registry, err := storage.NewRegistry(dbPath)
	if err != nil {
		t.Fatalf("Failed to create registry: %v", err)
	}

	// Insert test archives using the registry's database
	db := registry.DB()
	testArchives := []struct {
		uid, name, path, profile, status, metadata, checksum string
		size int64
		managed bool
	}{
		{"01K2E3BEJV6G", "project-backup", "/archives/project-backup.7z", "documents", "present", "Important project files", "sha256:abc123", 1024000, true},
		{"01K2E3CKJD8H", "media-photos", "/archives/media-photos.7z", "media", "present", "Family vacation photos", "sha256:def456", 5120000, true},
		{"01K2E3DMKF9J", "code-repository", "/archives/code-repository.7z", "documents", "present", "Source code backup", "sha256:ghi789", 512000, false},
		{"01K2E3GNLH2K", "video-project", "/archives/video-project.7z", "media", "missing", "Video editing project", "sha256:jkl012", 10240000, true},
	}

	for _, arch := range testArchives {
		_, err = db.Exec(`INSERT INTO archives (uid, name, path, size, created, checksum, profile, managed, status, original_path, uploaded, destination, metadata) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			arch.uid, arch.name, arch.path, arch.size, time.Now().Unix(), arch.checksum, arch.profile, arch.managed, arch.status, "", false, "", arch.metadata)
		if err != nil {
			t.Fatalf("Failed to insert test archive: %v", err)
		}
	}

	return db, registry
}

func TestSearchEngine_FullTextSearch(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	// Ensure search table exists
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test full-text search
	results, err := searchEngine.Search("project")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'project', got %d", len(results))
	}

	// Verify results contain expected archives
	found := make(map[string]bool)
	for _, result := range results {
		found[result.Name] = true
	}

	if !found["project-backup"] {
		t.Error("Expected to find 'project-backup' in search results")
	}
	if !found["video-project"] {
		t.Error("Expected to find 'video-project' in search results")
	}
}

func TestSearchEngine_FieldSpecificSearch(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test field-specific search
	results, err := searchEngine.SearchField("name", "backup")
	if err != nil {
		t.Fatalf("Field search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for name='backup', got %d", len(results))
	}

	if len(results) > 0 && results[0].Name != "project-backup" {
		t.Errorf("Expected 'project-backup', got '%s'", results[0].Name)
	}
}

func TestSearchEngine_RegexSearch(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test regex search
	results, err := searchEngine.SearchRegex("name", ".*-project$")
	if err != nil {
		t.Fatalf("Regex search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for regex '.*-project$', got %d", len(results))
	}

	if len(results) > 0 && results[0].Name != "video-project" {
		t.Errorf("Expected 'video-project', got '%s'", results[0].Name)
	}
}

func TestSearchEngine_SearchOptions(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test search with options
	opts := SearchOptions{
		Field:      "profile",
		MaxResults: 1,
	}

	results, err := searchEngine.SearchWithOptions("media", opts)
	if err != nil {
		t.Fatalf("Search with options failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result (limited), got %d", len(results))
	}
}

func TestSearchEngine_Performance(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test search performance
	start := time.Now()
	_, err := searchEngine.Search("project")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	duration := time.Since(start)

	// Should be well under 500ms target
	if duration > 100*time.Millisecond {
		t.Errorf("Search took %v, expected <100ms for small dataset", duration)
	}
}

func TestSearchEngine_Reindex(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test reindexing
	start := time.Now()
	err := searchEngine.Reindex()
	if err != nil {
		t.Fatalf("Reindex failed: %v", err)
	}
	duration := time.Since(start)

	// Reindexing should be fast
	if duration > 10*time.Millisecond {
		t.Errorf("Reindex took %v, expected <10ms for small dataset", duration)
	}

	// Verify index was updated
	if searchEngine.index.lastUpdate.IsZero() {
		t.Error("Index lastUpdate should be set after reindexing")
	}
}

func TestSearchEngine_EmptyQuery(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test empty query
	results, err := searchEngine.Search("")
	if err == nil {
		t.Error("Expected error for empty search query")
	}
	
	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty query, got %d", len(results))
	}
}

func TestSearchEngine_NoResults(t *testing.T) {
	db, registry := setupSearchTestDB(t)
	defer db.Close()

	searchEngine := NewSearchEngine(registry)
	
	if err := searchEngine.EnsureSearchTable(); err != nil {
		t.Fatalf("Failed to create search table: %v", err)
	}

	// Test query with no results
	results, err := searchEngine.Search("nonexistent")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", len(results))
	}
}

func TestLRUCache_Basic(t *testing.T) {
	cache := NewLRUCache(2)

	// Test cache miss
	result := cache.Get("key1")
	if result != nil {
		t.Error("Expected cache miss for 'key1'")
	}

	// Test cache set and get
	testArchives := []*storage.Archive{
		{Name: "test1"},
		{Name: "test2"},
	}
	
	cache.Set("key1", testArchives, 5*time.Minute)
	result = cache.Get("key1")
	if result == nil {
		t.Error("Expected cache hit for 'key1'")
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 cached archives, got %d", len(result))
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	cache := NewLRUCache(2)

	// Fill cache
	cache.Set("key1", []*storage.Archive{{Name: "test1"}}, 5*time.Minute)
	cache.Set("key2", []*storage.Archive{{Name: "test2"}}, 5*time.Minute)
	cache.Set("key3", []*storage.Archive{{Name: "test3"}}, 5*time.Minute)

	// key1 should be evicted
	if cache.Get("key1") != nil {
		t.Error("Expected key1 to be evicted")
	}
	
	// key2 and key3 should still be there
	if cache.Get("key2") == nil {
		t.Error("Expected key2 to still be cached")
	}
	if cache.Get("key3") == nil {
		t.Error("Expected key3 to still be cached")
	}
}

func TestLRUCache_Expiration(t *testing.T) {
	cache := NewLRUCache(10)

	// Set item with very short TTL
	cache.Set("key1", []*storage.Archive{{Name: "test1"}}, 1*time.Millisecond)
	
	// Wait for expiration
	time.Sleep(5 * time.Millisecond)
	
	// Should be expired
	result := cache.Get("key1")
	if result != nil {
		t.Error("Expected item to be expired")
	}
}