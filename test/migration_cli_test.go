package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrationCLI_Commands(t *testing.T) {
	// Skip if binary doesn't exist
	if _, err := exec.LookPath("7zarch-go"); err != nil {
		t.Skip("7zarch-go binary not found in PATH")
	}

	// Create temporary directory for test database
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	testConfigDir := filepath.Join(tmpDir, ".7zarch-go")

	// Set up test environment
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create the config directory
	if err := os.MkdirAll(testConfigDir, 0755); err != nil {
		t.Fatalf("failed to create test config dir: %v", err)
	}

	t.Run("db_status", func(t *testing.T) {
		// First run should initialize the database
		cmd := exec.Command("7zarch-go", "db", "status")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("db status failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)

		// Should show database path
		if !strings.Contains(outputStr, "Database:") {
			t.Error("Expected database path in output")
		}

		// Should show schema version
		if !strings.Contains(outputStr, "Schema Version:") {
			t.Error("Expected schema version in output")
		}

		// Should show migration counts
		if !strings.Contains(outputStr, "Applied Migrations:") {
			t.Error("Expected applied migrations count in output")
		}

		if !strings.Contains(outputStr, "Pending Migrations:") {
			t.Error("Expected pending migrations count in output")
		}
	})

	t.Run("db_backup", func(t *testing.T) {
		cmd := exec.Command("7zarch-go", "db", "backup")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("db backup failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)

		// Should show backup created message
		if !strings.Contains(outputStr, "Backup created:") {
			t.Error("Expected backup created message in output")
		}

		// Should contain timestamp format
		if !strings.Contains(outputStr, "registry-") || !strings.Contains(outputStr, ".bak") {
			t.Error("Expected timestamped backup filename in output")
		}
	})

	t.Run("db_migrate_dry_run", func(t *testing.T) {
		cmd := exec.Command("7zarch-go", "db", "migrate", "--dry-run")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("db migrate --dry-run failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)

		// Should show either dry run results or no pending migrations
		if !strings.Contains(outputStr, "No pending migrations") &&
			!strings.Contains(outputStr, "Dry run:") {
			t.Error("Expected dry run or no pending migrations message")
		}
	})

	t.Run("db_migrate", func(t *testing.T) {
		cmd := exec.Command("7zarch-go", "db", "migrate")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("db migrate failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)

		// Should show either successful migration or no pending migrations
		if !strings.Contains(outputStr, "No pending migrations") &&
			!strings.Contains(outputStr, "completed successfully") {
			t.Error("Expected migration success or no pending migrations message")
		}
	})

	t.Run("db_migrate_backup_only", func(t *testing.T) {
		cmd := exec.Command("7zarch-go", "db", "migrate", "--backup-only")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("db migrate --backup-only failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)

		// Should show backup created message
		if !strings.Contains(outputStr, "Backup created:") {
			t.Error("Expected backup created message in output")
		}
	})
}
