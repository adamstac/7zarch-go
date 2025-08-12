package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TestRegistry creates a temporary registry for testing
func setupTestRegistry(t *testing.T) (*Registry, string) {
	tempDir, err := os.MkdirTemp("", "7zarch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	dbPath := filepath.Join(tempDir, "test.db")
	registry, err := NewRegistry(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test registry: %v", err)
	}
	
	// Clean up function
	t.Cleanup(func() {
		registry.Close()
		os.RemoveAll(tempDir)
	})
	
	return registry, tempDir
}

// seedTestArchives creates test archives with known ULIDs and checksums
func seedTestArchives(t *testing.T, registry *Registry) []*Archive {
	archives := []*Archive{
		{
			UID:          "01JEX4RT2N9K3M6P8Q7S5V4W2X",
			Name:         "project-backup.7z",
			Path:         "/managed/project-backup.7z",
			Size:         2097152, // 2MB
			Created:      time.Now().Add(-2 * 24 * time.Hour),
			Checksum:     "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
			Profile:      "media",
			Managed:      true,
			Status:       "present",
			LastSeen:     nil,
			DeletedAt:    nil,
			OriginalPath: "",
			Uploaded:     false,
			Destination:  "",
			UploadedAt:   nil,
			Metadata:     "",
		},
		{
			UID:          "01JEY5SU3O0L4N7Q9R8T6W5X3Y",
			Name:         "project-docs.7z", 
			Path:         "/external/project-docs.7z",
			Size:         524288, // 512KB
			Created:      time.Now().Add(-7 * 24 * time.Hour),
			Checksum:     "b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef234567",
			Profile:      "documents",
			Managed:      false,
			Status:       "present",
			LastSeen:     nil,
			DeletedAt:    nil,
			OriginalPath: "",
			Uploaded:     false,
			Destination:  "",
			UploadedAt:   nil,
			Metadata:     "",
		},
		{
			UID:          "01JEZ6TV4P1M5O8R0S9U7X6Y4Z",
			Name:         "podcast-103.7z",
			Path:         "/managed/podcast-103.7z", 
			Size:         163577856, // 156MB
			Created:      time.Now().Add(-7 * 24 * time.Hour),
			Checksum:     "c3d4e5f6789012345678901234567890abcdef1234567890abcdef345678",
			Profile:      "media",
			Managed:      true,
			Status:       "present",
			LastSeen:     nil,
			DeletedAt:    nil,
			OriginalPath: "",
			Uploaded:     false,
			Destination:  "",
			UploadedAt:   nil,
			Metadata:     "",
		},
	}
	
	for _, archive := range archives {
		err := registry.Add(archive)
		if err != nil {
			t.Fatalf("Failed to seed archive %s: %v", archive.Name, err)
		}
	}
	
	return archives
}

// Error types are now implemented by AC in resolver.go

// setupTestResolver creates a test resolver with AC's real implementation
func setupTestResolver(t *testing.T) (*Resolver, *Registry) {
	registry, _ := setupTestRegistry(t)
	resolver := NewResolver(registry)
	resolver.MinPrefixLength = 4 // Use shorter prefix for tests
	return resolver, registry
}

// Test exact ULID resolution
func TestResolveExactULID(t *testing.T) {
	resolver, registry := setupTestResolver(t)
	archives := seedTestArchives(t, registry)
	
	// Test exact ULID match
	result, err := resolver.Resolve("01JEX4RT2N9K3M6P8Q7S5V4W2X")
	if err != nil {
		t.Fatalf("Expected successful resolution, got error: %v", err)
	}
	
	if result.Name != "project-backup.7z" {
		t.Errorf("Expected 'project-backup.7z', got '%s'", result.Name)
	}
	
	if result.UID != archives[0].UID {
		t.Errorf("Expected UID %s, got %s", archives[0].UID, result.UID)
	}
}

// Test ULID prefix resolution
func TestResolveULIDPrefix(t *testing.T) {
	resolver, registry := setupTestResolver(t)
	seedTestArchives(t, registry)
	
	testCases := []struct {
		prefix   string
		expected string
	}{
		{"01JEX", "project-backup.7z"},
		{"01JEY", "project-docs.7z"},
		{"01JEZ", "podcast-103.7z"},
		{"01JEX4RT", "project-backup.7z"},
	}
	
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("prefix_%s", tc.prefix), func(t *testing.T) {
			result, err := resolver.Resolve(tc.prefix)
			if err != nil {
				t.Fatalf("Expected successful resolution for prefix '%s', got error: %v", tc.prefix, err)
			}
			
			if result.Name != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result.Name)
			}
		})
	}
}

// Test ambiguous ULID prefix resolution
func TestResolveAmbiguousULIDPrefix(t *testing.T) {
	resolver, registry := setupTestResolver(t)
	seedTestArchives(t, registry)
	
	// All test ULIDs start with "01JE" - should be ambiguous
	_, err := resolver.Resolve("01JE")
	
	if err == nil {
		t.Fatal("Expected ambiguous error, got nil")
	}
	
	ambiguousErr, ok := err.(*AmbiguousIDError)
	if !ok {
		t.Fatalf("Expected AmbiguousIDError, got %T", err)
	}
	
	if len(ambiguousErr.Matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(ambiguousErr.Matches))
	}
	
	if ambiguousErr.ID != "01JE" {
		t.Errorf("Expected ID '01JE', got '%s'", ambiguousErr.ID)
	}
}

// Test checksum prefix resolution
func TestResolveChecksumPrefix(t *testing.T) {
	resolver, registry := setupTestResolver(t)
	seedTestArchives(t, registry)
	
	testCases := []struct {
		prefix   string
		expected string
	}{
		{"a1b2c3d4", "project-backup.7z"},
		{"b2c3d4e5", "project-docs.7z"},
		{"c3d4e5f6", "podcast-103.7z"},
		{"a1b2c3d4e5f6", "project-backup.7z"},
	}
	
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("checksum_%s", tc.prefix), func(t *testing.T) {
			result, err := resolver.Resolve(tc.prefix)
			if err != nil {
				t.Fatalf("Expected successful resolution for checksum prefix '%s', got error: %v", tc.prefix, err)
			}
			
			if result.Name != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result.Name)
			}
		})
	}
}

// Test name resolution
func TestResolveName(t *testing.T) {
	resolver, registry := setupTestResolver(t)
	seedTestArchives(t, registry)
	
	result, err := resolver.Resolve("project-backup.7z")
	if err != nil {
		t.Fatalf("Expected successful name resolution, got error: %v", err)
	}
	
	if result.Name != "project-backup.7z" {
		t.Errorf("Expected 'project-backup.7z', got '%s'", result.Name)
	}
	
	if result.UID != "01JEX4RT2N9K3M6P8Q7S5V4W2X" {
		t.Errorf("Expected specific UID, got '%s'", result.UID)
	}
}

// Test archive not found
func TestResolveNotFound(t *testing.T) {
	resolver, registry := setupTestResolver(t)
	seedTestArchives(t, registry)
	
	_, err := resolver.Resolve("nonexistent")
	
	if err == nil {
		t.Fatal("Expected not found error, got nil")
	}
	
	notFoundErr, ok := err.(*ArchiveNotFoundError)
	if !ok {
		t.Fatalf("Expected ArchiveNotFoundError, got %T", err)
	}
	
	if notFoundErr.ID != "nonexistent" {
		t.Errorf("Expected ID 'nonexistent', got '%s'", notFoundErr.ID)
	}
}

// Test empty registry
func TestResolveEmptyRegistry(t *testing.T) {
	resolver, _ := setupTestResolver(t)
	
	_, err := resolver.Resolve("01JEX")
	
	if err == nil {
		t.Fatal("Expected not found error, got nil")
	}
	
	_, ok := err.(*ArchiveNotFoundError)
	if !ok {
		t.Fatalf("Expected ArchiveNotFoundError, got %T", err)
	}
}

// Benchmark resolver performance
func BenchmarkResolverULIDPrefix(b *testing.B) {
	registry, tempDir := setupTestRegistry(&testing.T{})
	defer func() {
		registry.Close()
		os.RemoveAll(tempDir)
	}()
	
	// Seed with more archives for realistic benchmarking
	for i := 0; i < 100; i++ {
		archive := &Archive{
			UID:      generateUID(),
			Name:     fmt.Sprintf("archive-%03d.7z", i),
			Path:     fmt.Sprintf("/test/archive-%03d.7z", i),
			Size:     int64(1024 * (i + 1)),
			Created:  time.Now(),
			Checksum: fmt.Sprintf("hash%096d", i),
			Profile:  "balanced",
			Managed:  true,
			Status:   "present",
		}
		registry.Add(archive)
	}
	
	resolver := NewResolver(registry)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test resolving a ULID prefix
		resolver.Resolve("01J")
	}
}