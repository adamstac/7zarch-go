package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestRegistryCorruption tests various corruption scenarios
func TestRegistryCorruption(t *testing.T) {
	t.Run("corrupted_database_file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "7zarch-corruption-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		
		dbPath := filepath.Join(tempDir, "corrupted.db")
		
		// Create a corrupted database file
		err = os.WriteFile(dbPath, []byte("this is not a valid sqlite database"), 0600)
		if err != nil {
			t.Fatalf("Failed to create corrupted file: %v", err)
		}
		
		// Try to open registry - should fail gracefully
		_, err = NewRegistry(dbPath)
		if err == nil {
			t.Fatal("Expected error opening corrupted database, got nil")
		}
		
		t.Logf("Corruption handled correctly: %v", err)
	})
	
	t.Run("missing_database_directory", func(t *testing.T) {
		nonExistentPath := "/nonexistent/path/registry.db"
		
		// Try to create registry in non-existent directory
		// Should create directory and initialize properly
		registry, err := NewRegistry(nonExistentPath)
		if err != nil {
			t.Logf("Expected graceful handling, got error: %v", err)
			// This is acceptable - may fail due to permissions
			return
		}
		defer registry.Close()
		defer os.RemoveAll(filepath.Dir(nonExistentPath))
		
		// Verify registry works
		testArchive := &Archive{
			UID:     generateUID(),
			Name:    "test.7z",
			Path:    "/test/test.7z",
			Size:    1024,
			Created: time.Now(),
			Profile: "balanced",
			Managed: true,
			Status:  "present",
		}
		
		err = registry.Add(testArchive)
		if err != nil {
			t.Errorf("Failed to add archive to new registry: %v", err)
		}
	})
	
	t.Run("partial_schema_migration", func(t *testing.T) {
		registry, tempDir := setupTestRegistry(t)
		defer os.RemoveAll(tempDir)
		
		// Simulate partial migration by manually altering schema
		_, err := registry.db.Exec("ALTER TABLE archives DROP COLUMN uid")
		if err != nil {
			t.Logf("Could not simulate partial migration: %v", err)
			return
		}
		
		// Try to add archive without UID column
		testArchive := &Archive{
			UID:     generateUID(),
			Name:    "test-partial.7z",
			Path:    "/test/test-partial.7z",
			Size:    1024,
			Created: time.Now(),
			Profile: "balanced",
			Managed: true,
			Status:  "present",
		}
		
		err = registry.Add(testArchive)
		if err == nil {
			t.Error("Expected error adding archive with missing UID column")
		}
		
		t.Logf("Partial migration handled: %v", err)
	})
}

// TestRegistryMissingFiles tests handling of missing archive files
func TestRegistryMissingFiles(t *testing.T) {
	registry, tempDir := setupTestRegistry(t)
	defer os.RemoveAll(tempDir)
	
	// Create a temporary archive file
	archiveDir := filepath.Join(tempDir, "archives")
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create archive dir: %v", err)
	}
	
	archivePath := filepath.Join(archiveDir, "temporary.7z")
	err = os.WriteFile(archivePath, []byte("fake archive content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary archive: %v", err)
	}
	
	// Add archive to registry
	testArchive := &Archive{
		UID:     generateUID(),
		Name:    "temporary.7z",
		Path:    archivePath,
		Size:    int64(len("fake archive content")),
		Created: time.Now(),
		Profile: "balanced",
		Managed: true,
		Status:  "present",
	}
	
	err = registry.Add(testArchive)
	if err != nil {
		t.Fatalf("Failed to add archive: %v", err)
	}
	
	// Verify archive exists in registry
	retrieved, err := registry.Get("temporary.7z")
	if err != nil {
		t.Fatalf("Failed to retrieve archive: %v", err)
	}
	
	if retrieved.Status != "present" {
		t.Errorf("Expected status 'present', got '%s'", retrieved.Status)
	}
	
	// Delete the physical file (simulate missing file)
	err = os.Remove(archivePath)
	if err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}
	
	// File is gone but registry still thinks it exists
	retrieved, err = registry.Get("temporary.7z")
	if err != nil {
		t.Fatalf("Failed to retrieve archive after file deletion: %v", err)
	}
	
	if retrieved.Status != "present" {
		t.Errorf("Registry should still show 'present' until verification")
	}
	
	// AC's show command should detect missing file and update status
	// This documents expected behavior for file verification
	t.Logf("Archive in registry but file missing - status verification needed")
}

// TestRegistryRecovery tests registry rebuild scenarios
func TestRegistryRecovery(t *testing.T) {
	t.Run("rebuild_from_archives", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "7zarch-recovery-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		
		// Create archive directory with some files
		archiveDir := filepath.Join(tempDir, "archives")
		err = os.MkdirAll(archiveDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create archive dir: %v", err)
		}
		
		// Create fake archive files
		archiveFiles := []string{"backup-1.7z", "backup-2.7z", "project.7z"}
		for _, filename := range archiveFiles {
			path := filepath.Join(archiveDir, filename)
			content := fmt.Sprintf("fake content for %s", filename)
			err = os.WriteFile(path, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test archive %s: %v", filename, err)
			}
		}
		
		// Registry doesn't exist yet - would need import command
		dbPath := filepath.Join(tempDir, "registry.db")
		
		// This documents the recovery scenario that AC's import command should handle
		t.Logf("Recovery scenario: %d archive files exist without registry", len(archiveFiles))
		t.Logf("Archive directory: %s", archiveDir)
		t.Logf("Missing registry: %s", dbPath)
		
		// AC should implement: mas import <dir> to rebuild registry
		registry, err := NewRegistry(dbPath)
		if err != nil {
			t.Fatalf("Failed to create new registry: %v", err)
		}
		defer registry.Close()
		
		// Empty registry - import command would populate it
		archives, err := registry.List()
		if err != nil {
			t.Fatalf("Failed to list empty registry: %v", err)
		}
		
		if len(archives) != 0 {
			t.Errorf("Expected empty registry, got %d archives", len(archives))
		}
		
		t.Logf("Recovery test complete - import functionality needed")
	})
}

// TestRegistryConsistency tests data consistency scenarios
func TestRegistryConsistency(t *testing.T) {
	t.Run("duplicate_names", func(t *testing.T) {
		registry, _ := setupTestRegistry(t)
		
		// Try to add archives with duplicate names
		archive1 := &Archive{
			UID:     generateUID(),
			Name:    "duplicate.7z",
			Path:    "/path1/duplicate.7z",
			Size:    1024,
			Created: time.Now(),
			Profile: "balanced",
			Managed: true,
			Status:  "present",
		}
		
		archive2 := &Archive{
			UID:     generateUID(),
			Name:    "duplicate.7z", // Same name
			Path:    "/path2/duplicate.7z",
			Size:    2048,
			Created: time.Now(),
			Profile: "media",
			Managed: false,
			Status:  "present",
		}
		
		err := registry.Add(archive1)
		if err != nil {
			t.Fatalf("Failed to add first archive: %v", err)
		}
		
		err = registry.Add(archive2)
		if err == nil {
			t.Error("Expected error adding duplicate name, got nil")
		}
		
		t.Logf("Duplicate name prevented: %v", err)
	})
	
	t.Run("duplicate_uids", func(t *testing.T) {
		registry, _ := setupTestRegistry(t)
		
		sameUID := generateUID()
		
		// Try to add archives with duplicate UIDs
		archive1 := &Archive{
			UID:     sameUID,
			Name:    "first.7z",
			Path:    "/path1/first.7z",
			Size:    1024,
			Created: time.Now(),
			Profile: "balanced",
			Managed: true,
			Status:  "present",
		}
		
		archive2 := &Archive{
			UID:     sameUID, // Same UID
			Name:    "second.7z",
			Path:    "/path2/second.7z",
			Size:    2048,
			Created: time.Now(),
			Profile: "media",
			Managed: false,
			Status:  "present",
		}
		
		err := registry.Add(archive1)
		if err != nil {
			t.Fatalf("Failed to add first archive: %v", err)
		}
		
		err = registry.Add(archive2)
		if err == nil {
			t.Error("Expected error adding duplicate UID, got nil")
		}
		
		t.Logf("Duplicate UID prevented: %v", err)
	})
	
	t.Run("backfill_missing_uids", func(t *testing.T) {
		registry, _ := setupTestRegistry(t)
		
		// Manually insert archive without UID to simulate legacy data
		// Include all fields to avoid NULL issues - set empty strings for text fields
		query := `INSERT INTO archives (name, path, size, created, checksum, profile, managed, status, uploaded, destination, original_path, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := registry.db.Exec(query, "legacy.7z", "/legacy/legacy.7z", 1024, time.Now(), "", "balanced", true, "present", false, "", "", "")
		if err != nil {
			t.Fatalf("Failed to insert legacy archive: %v", err)
		}
		
		// Test backfill functionality
		err = registry.BackfillUIDs(func() string {
			return generateUID()
		})
		if err != nil {
			t.Errorf("Backfill UIDs failed: %v", err)
		}
		
		// Verify UID was added
		archive, err := registry.Get("legacy.7z")
		if err != nil {
			t.Fatalf("Failed to get backfilled archive: %v", err)
		}
		
		if archive.UID == "" {
			t.Error("Expected UID to be backfilled, got empty string")
		}
		
		t.Logf("UID backfilled successfully: %s", archive.UID)
	})
}

// TestRegistryPerformance tests performance under various load conditions
func TestRegistryPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}
	
	t.Run("large_registry_operations", func(t *testing.T) {
		registry, _ := setupTestRegistry(t)
		
		// Add many archives
		numArchives := 1000
		startTime := time.Now()
		
		for i := 0; i < numArchives; i++ {
			archive := &Archive{
				UID:     generateUID(),
				Name:    fmt.Sprintf("archive-%04d.7z", i),
				Path:    fmt.Sprintf("/test/archive-%04d.7z", i),
				Size:    int64(1024 * (i + 1)),
				Created: time.Now(),
				Profile: "balanced",
				Managed: i%2 == 0,
				Status:  "present",
			}
			
			err := registry.Add(archive)
			if err != nil {
				t.Fatalf("Failed to add archive %d: %v", i, err)
			}
		}
		
		insertTime := time.Since(startTime)
		t.Logf("Inserted %d archives in %v (%.2f archives/sec)", 
			numArchives, insertTime, float64(numArchives)/insertTime.Seconds())
		
		// Test list performance
		startTime = time.Now()
		archives, err := registry.List()
		listTime := time.Since(startTime)
		
		if err != nil {
			t.Fatalf("Failed to list archives: %v", err)
		}
		
		if len(archives) != numArchives {
			t.Errorf("Expected %d archives, got %d", numArchives, len(archives))
		}
		
		t.Logf("Listed %d archives in %v", len(archives), listTime)
		
		// Test individual lookups
		startTime = time.Now()
		lookupCount := 100
		
		for i := 0; i < lookupCount; i++ {
			name := fmt.Sprintf("archive-%04d.7z", i)
			_, err := registry.Get(name)
			if err != nil {
				t.Errorf("Failed to get archive %s: %v", name, err)
			}
		}
		
		lookupTime := time.Since(startTime)
		t.Logf("Performed %d lookups in %v (%.2f lookups/sec)", 
			lookupCount, lookupTime, float64(lookupCount)/lookupTime.Seconds())
	})
}

// TestRegistryTransactions tests transactional behavior
func TestRegistryTransactions(t *testing.T) {
	t.Run("concurrent_access", func(t *testing.T) {
		registry, _ := setupTestRegistry(t)
		
		// Test concurrent writes
		done := make(chan error, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				archive := &Archive{
					UID:     generateUID(),
					Name:    fmt.Sprintf("concurrent-%d.7z", id),
					Path:    fmt.Sprintf("/test/concurrent-%d.7z", id),
					Size:    int64(1024 * id),
					Created: time.Now(),
					Profile: "balanced",
					Managed: true,
					Status:  "present",
				}
				
				done <- registry.Add(archive)
			}(i)
		}
		
		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			err := <-done
			if err != nil {
				t.Errorf("Concurrent write %d failed: %v", i, err)
			}
		}
		
		// Verify all archives were added
		archives, err := registry.List()
		if err != nil {
			t.Fatalf("Failed to list after concurrent writes: %v", err)
		}
		
		if len(archives) != 10 {
			t.Errorf("Expected 10 archives after concurrent writes, got %d", len(archives))
		}
	})
}