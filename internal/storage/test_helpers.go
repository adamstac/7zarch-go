package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
)

// Test helpers for MAS operations per 7EP-0004

// TestRegistry creates an in-memory registry for testing
func TestRegistry(t *testing.T) *Registry {
	t.Helper()
	
	// Create temp directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	
	reg, err := NewRegistry(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test registry: %v", err)
	}
	
	// Initialize schema
	if err := reg.initSchema(); err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}
	
	t.Cleanup(func() {
		reg.Close()
		os.RemoveAll(tmpDir)
	})
	
	return reg
}

// CreateTestArchive creates a test archive with predictable data
func CreateTestArchive(t *testing.T, reg *Registry, name string, opts ...TestArchiveOption) *Archive {
	t.Helper()
	
	cfg := &testArchiveConfig{
		name:     name,
		size:     1024,
		profile:  "balanced",
		managed:  true,
		status:   "present",
		created:  time.Now(),
	}
	
	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}
	
	// Generate ULID
	uid := generateUID()
	
	// Create archive
	archive := &Archive{
		UID:      uid,
		Name:     cfg.name,
		Path:     filepath.Join("/test/archives", cfg.name),
		Size:     cfg.size,
		Created:  cfg.created,
		Checksum: fmt.Sprintf("sha256:%x", []byte(cfg.name)),
		Profile:  cfg.profile,
		Managed:  cfg.managed,
		Status:   cfg.status,
	}
	
	// Register archive
	if err := reg.Register(archive); err != nil {
		t.Fatalf("Failed to register test archive: %v", err)
	}
	
	return archive
}

// TestArchiveOption configures test archive creation
type TestArchiveOption func(*testArchiveConfig)

type testArchiveConfig struct {
	name     string
	size     int64
	profile  string
	managed  bool
	status   string
	created  time.Time
}

// WithSize sets archive size
func WithSize(size int64) TestArchiveOption {
	return func(cfg *testArchiveConfig) {
		cfg.size = size
	}
}

// WithProfile sets compression profile
func WithProfile(profile string) TestArchiveOption {
	return func(cfg *testArchiveConfig) {
		cfg.profile = profile
	}
}

// WithManaged sets managed status
func WithManaged(managed bool) TestArchiveOption {
	return func(cfg *testArchiveConfig) {
		cfg.managed = managed
	}
}

// WithStatus sets archive status
func WithStatus(status string) TestArchiveOption {
	return func(cfg *testArchiveConfig) {
		cfg.status = status
	}
}

// WithCreated sets creation time
func WithCreated(created time.Time) TestArchiveOption {
	return func(cfg *testArchiveConfig) {
		cfg.created = created
	}
}

// CreateTestSet creates a standard set of test archives
func CreateTestSet(t *testing.T, reg *Registry) []*Archive {
	t.Helper()
	
	archives := []*Archive{
		CreateTestArchive(t, reg, "project-backup.7z", 
			WithSize(2*1024*1024), 
			WithProfile("documents")),
		
		CreateTestArchive(t, reg, "project-docs.7z", 
			WithSize(512*1024), 
			WithProfile("documents")),
		
		CreateTestArchive(t, reg, "media-files.7z", 
			WithSize(100*1024*1024), 
			WithProfile("media"),
			WithManaged(false)),
		
		CreateTestArchive(t, reg, "old-archive.7z", 
			WithCreated(time.Now().Add(-30*24*time.Hour)),
			WithStatus("deleted")),
		
		CreateTestArchive(t, reg, "missing-file.7z",
			WithStatus("missing")),
	}
	
	return archives
}

// AssertResolves verifies that an ID resolves to expected archive
func AssertResolves(t *testing.T, resolver *Resolver, input string, expected *Archive) {
	t.Helper()
	
	result, err := resolver.ResolveID(input)
	if err != nil {
		t.Errorf("Failed to resolve '%s': %v", input, err)
		return
	}
	
	if result.UID != expected.UID {
		t.Errorf("Resolved wrong archive for '%s': got %s, want %s", 
			input, result.UID, expected.UID)
	}
}

// AssertAmbiguous verifies that an ID is ambiguous
func AssertAmbiguous(t *testing.T, resolver *Resolver, input string, expectedCount int) {
	t.Helper()
	
	_, err := resolver.ResolveID(input)
	if err == nil {
		t.Errorf("Expected ambiguous error for '%s', got success", input)
		return
	}
	
	ambErr, ok := err.(*AmbiguousIDError)
	if !ok {
		t.Errorf("Expected AmbiguousIDError for '%s', got %T", input, err)
		return
	}
	
	if len(ambErr.Matches) != expectedCount {
		t.Errorf("Wrong match count for '%s': got %d, want %d",
			input, len(ambErr.Matches), expectedCount)
	}
}

// AssertNotFound verifies that an ID is not found
func AssertNotFound(t *testing.T, resolver *Resolver, input string) {
	t.Helper()
	
	_, err := resolver.ResolveID(input)
	if err == nil {
		t.Errorf("Expected not found error for '%s', got success", input)
		return
	}
	
	_, ok := err.(*ArchiveNotFoundError)
	if !ok {
		t.Errorf("Expected ArchiveNotFoundError for '%s', got %T", input, err)
	}
}

// BenchmarkResolver helps benchmark resolution performance
func BenchmarkResolver(b *testing.B, archiveCount int) {
	// Create registry with many archives
	reg := &Registry{} // Would use TestRegistry in real implementation
	resolver := NewResolver(reg)
	
	// Create test archives
	archives := make([]*Archive, archiveCount)
	for i := 0; i < archiveCount; i++ {
		archives[i] = &Archive{
			UID:  fmt.Sprintf("01K2E%06d", i),
			Name: fmt.Sprintf("archive-%d.7z", i),
		}
		reg.Register(archives[i])
	}
	
	// Benchmark resolution
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test various resolution types
		resolver.ResolveID(archives[i%archiveCount].UID[:8])
	}
}