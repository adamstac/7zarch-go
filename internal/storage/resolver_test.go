package storage

import (
	"database/sql"
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
			UID:      "01JEX4RT2N9K3M6P8Q7S5V4W2X",
			Name:     "project-backup.7z",
			Path:     "/managed/project-backup.7z",
			Size:     2097152, // 2MB
			Created:  time.Now().Add(-2 * 24 * time.Hour),
			Checksum: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
			Profile:  "media",
			Managed:  true,
			Status:   "present",
		},
		{
			UID:      "01JEY5SU3O0L4N7Q9R8T6W5X3Y",
			Name:     "project-docs.7z", 
			Path:     "/external/project-docs.7z",
			Size:     524288, // 512KB
			Created:  time.Now().Add(-7 * 24 * time.Hour),
			Checksum: "b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef234567",
			Profile:  "documents",
			Managed:  false,
			Status:   "present",
		},
		{
			UID:      "01JEZ6TV4P1M5O8R0S9U7X6Y4Z",
			Name:     "podcast-103.7z",
			Path:     "/managed/podcast-103.7z", 
			Size:     163577856, // 156MB
			Created:  time.Now().Add(-7 * 24 * time.Hour),
			Checksum: "c3d4e5f6789012345678901234567890abcdef1234567890abcdef345678",
			Profile:  "media",
			Managed:  true,
			Status:   "present",
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

// Resolver error types that AC will implement
type ArchiveNotFoundError struct {
	ID string
}

func (e *ArchiveNotFoundError) Error() string {
	return fmt.Sprintf("archive not found: %s", e.ID)
}

type AmbiguousIDError struct {
	ID      string
	Matches []*Archive
}

func (e *AmbiguousIDError) Error() string {
	return fmt.Sprintf("ambiguous ID '%s': %d matches found", e.ID, len(e.Matches))
}

// Mock Resolver implementation for testing - AC will replace this
type MockResolver struct {
	registry *Registry
}

func NewMockResolver(registry *Registry) *MockResolver {
	return &MockResolver{registry: registry}
}

func (r *MockResolver) ResolveID(input string) (*Archive, error) {
	// This is a simplified mock - AC will implement the full algorithm
	
	// Try exact UID match first
	archives, err := r.registry.List()
	if err != nil {
		return nil, err
	}
	
	var uidMatches []*Archive
	var checksumMatches []*Archive
	var nameMatches []*Archive
	
	for _, archive := range archives {
		// Exact UID match
		if archive.UID == input {
			return archive, nil
		}
		
		// UID prefix match
		if len(input) >= 4 && len(archive.UID) >= len(input) && 
		   archive.UID[:len(input)] == input {
			uidMatches = append(uidMatches, archive)
		}
		
		// Checksum prefix match
		if len(input) >= 8 && len(archive.Checksum) >= len(input) &&
		   archive.Checksum[:len(input)] == input {
			checksumMatches = append(checksumMatches, archive)
		}
		
		// Name match
		if archive.Name == input {
			nameMatches = append(nameMatches, archive)
		}
	}
	
	// Check UID prefix matches
	if len(uidMatches) == 1 {
		return uidMatches[0], nil
	} else if len(uidMatches) > 1 {
		return nil, &AmbiguousIDError{ID: input, Matches: uidMatches}
	}
	
	// Check checksum prefix matches  
	if len(checksumMatches) == 1 {
		return checksumMatches[0], nil
	} else if len(checksumMatches) > 1 {
		return nil, &AmbiguousIDError{ID: input, Matches: checksumMatches}
	}
	
	// Check name matches
	if len(nameMatches) == 1 {
		return nameMatches[0], nil
	} else if len(nameMatches) > 1 {
		return nil, &AmbiguousIDError{ID: input, Matches: nameMatches}
	}
	
	return nil, &ArchiveNotFoundError{ID: input}
}

// Test exact ULID resolution
func TestResolveExactULID(t *testing.T) {
	registry, _ := setupTestRegistry(t)
	archives := seedTestArchives(t, registry)
	resolver := NewMockResolver(registry)
	
	// Test exact ULID match
	result, err := resolver.ResolveID("01JEX4RT2N9K3M6P8Q7S5V4W2X")
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
	registry, _ := setupTestRegistry(t)
	seedTestArchives(t, registry)
	resolver := NewMockResolver(registry)
	
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
			result, err := resolver.ResolveID(tc.prefix)
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
	registry, _ := setupTestRegistry(t)
	seedTestArchives(t, registry)
	resolver := NewMockResolver(registry)
	
	// All test ULIDs start with "01JE" - should be ambiguous
	_, err := resolver.ResolveID("01JE")
	
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
	registry, _ := setupTestRegistry(t)
	seedTestArchives(t, registry)
	resolver := NewMockResolver(registry)
	
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
			result, err := resolver.ResolveID(tc.prefix)
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
	registry, _ := setupTestRegistry(t)
	seedTestArchives(t, registry)
	resolver := NewMockResolver(registry)
	
	result, err := resolver.ResolveID("project-backup.7z")
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
	registry, _ := setupTestRegistry(t)
	seedTestArchives(t, registry)
	resolver := NewMockResolver(registry)
	
	_, err := resolver.ResolveID("nonexistent")
	
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
	registry, _ := setupTestRegistry(t)
	resolver := NewMockResolver(registry)
	
	_, err := resolver.ResolveID("01JEX")
	
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
	
	resolver := NewMockResolver(registry)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test resolving a ULID prefix
		resolver.ResolveID("01J")
	}
}