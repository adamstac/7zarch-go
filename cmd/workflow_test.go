package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// TestCompleteTrashWorkflow tests the full delete -> restore -> purge lifecycle
func TestCompleteTrashWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "7zarch-workflow-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive file
	testFile := filepath.Join(tmpDir, "test-archive.7z")
	if err := os.WriteFile(testFile, []byte("test archive content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	managedPath := filepath.Join(tmpDir, "managed")
	cfg := &config.Config{
		Storage: config.StorageConfig{
			ManagedPath:   managedPath,
			RetentionDays: 7,
		},
	}

	// Initialize storage manager
	mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}
	defer mgr.Close()

	// Step 1: Create an archive
	stat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	err = mgr.Add("test-archive.7z", testFile, stat.Size(), "balanced", "", "", true)
	if err != nil {
		t.Fatalf("Failed to add archive: %v", err)
	}

	// Retrieve the archive
	archive, err := mgr.Registry().Get("test-archive.7z")
	if err != nil {
		t.Fatalf("Failed to get archive: %v", err)
	}

	if archive.Status != "present" {
		t.Errorf("Expected archive status 'present', got '%s'", archive.Status)
	}

	// Step 2: Delete the archive (soft delete)
	now := time.Now()
	archive.Status = "deleted"
	archive.DeletedAt = &now
	archive.OriginalPath = archive.Path

	// Move to trash
	trashDir := mgr.GetTrashPath()
	if err := os.MkdirAll(trashDir, 0750); err != nil {
		t.Fatalf("Failed to create trash directory: %v", err)
	}
	trashPath := filepath.Join(trashDir, filepath.Base(archive.Path))
	if err := os.Rename(archive.Path, trashPath); err != nil {
		t.Fatalf("Failed to move to trash: %v", err)
	}
	archive.Path = trashPath

	if err := mgr.Registry().Update(archive); err != nil {
		t.Fatalf("Failed to update archive: %v", err)
	}

	// Step 3: Verify archive is in deleted state
	resolver := storage.NewResolver(mgr.Registry())
	deletedArchive, err := resolver.Resolve(archive.UID)
	if err != nil {
		t.Fatalf("Failed to resolve deleted archive: %v", err)
	}

	if deletedArchive.Status != "deleted" {
		t.Errorf("Expected deleted archive status 'deleted', got '%s'", deletedArchive.Status)
	}

	if deletedArchive.DeletedAt == nil {
		t.Error("Expected deleted archive to have DeletedAt timestamp")
	}

	// Step 4: Test --deleted filter in list command
	archives, err := mgr.List()
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	opts := listFilters{onlyDeleted: true}
	filtered := applyAllFilters(archives, opts)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 deleted archive, got %d", len(filtered))
	}

	if filtered[0].UID != archive.UID {
		t.Errorf("Expected deleted archive UID %s, got %s", archive.UID, filtered[0].UID)
	}

	// Step 5: Test restore functionality
	// Simulate restore command logic
	target := deletedArchive.OriginalPath
	if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
		t.Fatalf("Failed to prepare restore destination: %v", err)
	}

	if err := os.Rename(deletedArchive.Path, target); err != nil {
		t.Fatalf("Failed to restore file: %v", err)
	}

	deletedArchive.Path = target
	deletedArchive.Status = "present"
	deletedArchive.DeletedAt = nil
	now = time.Now()
	deletedArchive.LastSeen = &now

	if err := mgr.Registry().Update(deletedArchive); err != nil {
		t.Fatalf("Failed to update restored archive: %v", err)
	}

	// Step 6: Verify restoration
	restoredArchive, err := resolver.Resolve(archive.UID)
	if err != nil {
		t.Fatalf("Failed to resolve restored archive: %v", err)
	}

	if restoredArchive.Status != "present" {
		t.Errorf("Expected restored archive status 'present', got '%s'", restoredArchive.Status)
	}

	if restoredArchive.DeletedAt != nil {
		t.Error("Expected restored archive to have nil DeletedAt")
	}

	// Step 7: Delete again and test purge
	now = time.Now()
	restoredArchive.Status = "deleted"
	restoredArchive.DeletedAt = &now
	restoredArchive.OriginalPath = restoredArchive.Path

	// Move back to trash
	trashPath = filepath.Join(trashDir, filepath.Base(restoredArchive.Path))
	if err := os.Rename(restoredArchive.Path, trashPath); err != nil {
		t.Fatalf("Failed to move to trash again: %v", err)
	}
	restoredArchive.Path = trashPath

	if err := mgr.Registry().Update(restoredArchive); err != nil {
		t.Fatalf("Failed to update re-deleted archive: %v", err)
	}

	// Step 8: Test purge functionality
	// Simulate immediate purge (--all flag behavior)
	if _, err := os.Stat(restoredArchive.Path); err == nil {
		if err := os.Remove(restoredArchive.Path); err != nil {
			t.Fatalf("Failed to remove file during purge: %v", err)
		}
	}

	if err := mgr.Registry().Delete(restoredArchive.Name); err != nil {
		t.Fatalf("Failed to remove from registry during purge: %v", err)
	}

	// Step 9: Verify purge
	_, err = resolver.Resolve(archive.UID)
	if err == nil {
		t.Error("Expected archive to be completely purged, but it still exists")
	}
}

// TestMachineReadableOutput tests JSON, CSV, and YAML output formats
func TestMachineReadableOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "7zarch-output-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	managedPath := filepath.Join(tmpDir, "managed")
	mgr, err := storage.NewManager(managedPath)
	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}
	defer mgr.Close()

	// Create test archive
	testFile := filepath.Join(tmpDir, "test.7z")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	err = mgr.Add("test.7z", testFile, stat.Size(), "balanced", "", "", true)
	if err != nil {
		t.Fatalf("Failed to add archive: %v", err)
	}

	archive, err := mgr.Registry().Get("test.7z")
	if err != nil {
		t.Fatalf("Failed to get archive: %v", err)
	}

	// Test JSON output
	t.Run("JSON Output", func(t *testing.T) {
		archives := []*storage.Archive{archive}
		err := outputJSON(archives)
		if err != nil {
			t.Errorf("JSON output failed: %v", err)
		}
	})

	// Test CSV output
	t.Run("CSV Output", func(t *testing.T) {
		archives := []*storage.Archive{archive}
		err := outputCSV(archives)
		if err != nil {
			t.Errorf("CSV output failed: %v", err)
		}
	})

	// Test YAML output
	t.Run("YAML Output", func(t *testing.T) {
		archives := []*storage.Archive{archive}
		err := outputYAML(archives)
		if err != nil {
			t.Errorf("YAML output failed: %v", err)
		}
	})
}

// TestTrashListOutput tests machine-readable output for trash list
func TestTrashListOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "7zarch-trash-output-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	managedPath := filepath.Join(tmpDir, "managed")
	mgr, err := storage.NewManager(managedPath)
	if err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}
	defer mgr.Close()

	// Create and delete an archive
	testFile := filepath.Join(tmpDir, "deleted.7z")
	if err := os.WriteFile(testFile, []byte("deleted content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	err = mgr.Add("deleted.7z", testFile, stat.Size(), "balanced", "", "", true)
	if err != nil {
		t.Fatalf("Failed to add archive: %v", err)
	}

	archive, err := mgr.Registry().Get("deleted.7z")
	if err != nil {
		t.Fatalf("Failed to get archive: %v", err)
	}

	// Delete the archive
	now := time.Now()
	archive.Status = "deleted"
	archive.DeletedAt = &now
	if err := mgr.Registry().Update(archive); err != nil {
		t.Fatalf("Failed to update deleted archive: %v", err)
	}

	// Test trash list JSON output
	t.Run("Trash List JSON", func(t *testing.T) {
		var output strings.Builder
		archives := []*storage.Archive{archive}
		err := outputTrashList(archives, "json", 7, &output)
		if err != nil {
			t.Errorf("Trash list JSON output failed: %v", err)
			return
		}

		// Verify JSON structure
		var rows []trashRow
		if err := json.Unmarshal([]byte(output.String()), &rows); err != nil {
			t.Errorf("Failed to parse JSON output: %v", err)
			return
		}

		if len(rows) != 1 {
			t.Errorf("Expected 1 trash row, got %d", len(rows))
			return
		}

		row := rows[0]
		if row.UID != archive.UID {
			t.Errorf("Expected UID %s, got %s", archive.UID, row.UID)
		}

		if row.Name != archive.Name {
			t.Errorf("Expected name %s, got %s", archive.Name, row.Name)
		}

		if row.DaysLeft < 6 || row.DaysLeft > 7 {
			t.Errorf("Expected 6-7 days left, got %d", row.DaysLeft)
		}
	})

	// Test trash list CSV output
	t.Run("Trash List CSV", func(t *testing.T) {
		var output strings.Builder
		archives := []*storage.Archive{archive}
		err := outputTrashList(archives, "csv", 7, &output)
		if err != nil {
			t.Errorf("Trash list CSV output failed: %v", err)
			return
		}

		lines := strings.Split(strings.TrimSpace(output.String()), "\n")
		if len(lines) != 2 { // header + data row
			t.Errorf("Expected 2 lines in CSV output, got %d", len(lines))
			return
		}

		header := lines[0]
		if !strings.Contains(header, "uid") || !strings.Contains(header, "name") {
			t.Errorf("CSV header missing expected fields: %s", header)
		}
	})
}

// TestOlderThanParsing tests the parseOlderThan helper function
func TestOlderThanParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"30d", 30 * 24 * time.Hour, false},
		{"1w", 7 * 24 * time.Hour, false},
		{"2h", 2 * time.Hour, false},
		{"30m", 30 * time.Minute, false},
		{"", 0, true},
		{"invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			cutoff, err := parseOlderThan(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error for input '%s', got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
				return
			}

			now := time.Now()
			expectedCutoff := now.Add(-tc.expected)
			diff := cutoff.Sub(expectedCutoff)

			// Allow 1 second tolerance for test execution time
			if diff > time.Second || diff < -time.Second {
				t.Errorf("Expected cutoff around %v, got %v (diff: %v)", expectedCutoff, cutoff, diff)
			}
		})
	}
}
