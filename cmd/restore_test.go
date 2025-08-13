package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
)

func execRestore(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := RestoreCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	return buf.String(), cmd.Execute()
}

func TestRestoreManagedAndForce(t *testing.T) {
	base := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := config.DefaultConfig()
	cfg.Storage.ManagedPath = base
	data, _ := yaml.Marshal(cfg)
	_ = os.WriteFile(filepath.Join(home, ".7zarch-go-config"), data, 0644)

	mgr, err := storage.NewManager(base)
	if err != nil {
		t.Fatalf("mgr: %v", err)
	}
	defer mgr.Close()

	// Seed a managed deleted archive whose file is in trash
	name := "restore-me.7z"
	trashFile := filepath.Join(mgr.GetTrashPath(), name)
	_ = os.MkdirAll(mgr.GetTrashPath(), 0755)
	_ = os.WriteFile(trashFile, []byte("data"), 0644)

	now := time.Now()
	arc := &storage.Archive{
		UID:          "uid-restore",
		Name:         name,
		Path:         trashFile,
		Size:         4,
		Created:      now.Add(-24 * time.Hour),
		Managed:      true,
		Status:       "deleted",
		DeletedAt:    &now,
		OriginalPath: filepath.Join(mgr.GetArchivesPath(), name),
	}
	if err := mgr.Registry().Add(arc); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Case 1: restore normally
	out, err := execRestore(t, arc.UID)
	if err != nil {
		t.Fatalf("restore: %v; out=%s", err, out)
	}
	// File should be moved back to OriginalPath
	if _, err := os.Stat(arc.OriginalPath); err != nil {
		t.Fatalf("restored file missing: %v", err)
	}

	// Re-fetch record; status should be present
	updated, err := mgr.Registry().Get(name)
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if updated.Status != "present" || updated.DeletedAt != nil {
		t.Fatalf("wrong status after restore: %+v", updated)
	}

	// Case 2: existing target but --force (simulate by moving back to trash then re-creating file and restoring with force)
	// Move file back to trash to re-run
	_ = os.Rename(updated.Path, trashFile)
	updated.Path = trashFile
	updated.Status = "deleted"
	when := time.Now()
	updated.DeletedAt = &when
	_ = mgr.Registry().Update(updated)
	// Create conflicting file at target
	_ = os.WriteFile(updated.OriginalPath, []byte("conflict"), 0644)

	cmd := RestoreCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--force", updated.UID})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("restore --force: %v; out=%s", err, buf.String())
	}
}

func TestRestoreExternalSoftDeleted(t *testing.T) {
	base := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := config.DefaultConfig()
	cfg.Storage.ManagedPath = base
	data, _ := yaml.Marshal(cfg)
	_ = os.WriteFile(filepath.Join(home, ".7zarch-go-config"), data, 0644)

	mgr, err := storage.NewManager(base)
	if err != nil {
		t.Fatalf("mgr: %v", err)
	}
	defer mgr.Close()

	name := "external-rest.7z"
	orig := filepath.Join(base, "outside", name) // pretend external
	_ = os.MkdirAll(filepath.Dir(orig), 0755)
	_ = os.WriteFile(orig, []byte("data"), 0644)

	now := time.Now()
	arc := &storage.Archive{
		UID:          "uid-external",
		Name:         name,
		Path:         orig, // for external soft delete, file remains
		Created:      now.Add(-48 * time.Hour),
		Managed:      false,
		Status:       "deleted",
		DeletedAt:    &now,
		OriginalPath: orig,
	}
	if err := mgr.Registry().Add(arc); err != nil {
		t.Fatalf("seed: %v", err)
	}

	out, err := execRestore(t, arc.UID)
	if err != nil {
		t.Fatalf("restore external: %v; out=%s", err, out)
	}
	updated, err := mgr.Registry().Get(name)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if updated.Status != "present" || updated.DeletedAt != nil {
		t.Fatalf("bad status: %+v", updated)
	}
}
