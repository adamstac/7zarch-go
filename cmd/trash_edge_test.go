package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/pflag"
)

// Helpers
func runListWithFlags(t *testing.T, setFlags func(*pflag.FlagSet)) (string, error) {
	cmd := trashListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if setFlags != nil {
		setFlags(cmd.Flags())
	}
	if cmd.RunE != nil {
		_ = cmd.RunE(cmd, []string{})
	}
	return buf.String(), nil
}

func TestEdgeCases_TrashAndRestore(t *testing.T) {
	base := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := config.DefaultConfig()
	cfg.Storage.ManagedPath = base
	cfg.Storage.RetentionDays = 7
	data, _ := yaml.Marshal(cfg)
	_ = os.WriteFile(filepath.Join(home, ".7zarch-go-config"), data, 0644)

	mgr, err := storage.NewManager(base)
	if err != nil {
		t.Fatalf("mgr: %v", err)
	}
	defer mgr.Close()

	now := time.Now()

	// 1) Countdown boundary at 0 days (purge today)
	zeroName := "purge-today.7z"
	when := now.Add(-7 * 24 * time.Hour) // retention=7 -> 0 days left
	arcZero := &storage.Archive{
		UID:          "uid-zero",
		Name:         zeroName,
		Path:         filepath.Join(mgr.GetTrashPath(), zeroName),
		Created:      now.Add(-10 * 24 * time.Hour),
		Managed:      true,
		Status:       "deleted",
		DeletedAt:    &when,
		OriginalPath: filepath.Join(mgr.GetArchivesPath(), zeroName),
	}
	_ = mgr.Registry().Add(arcZero)
	_ = os.MkdirAll(mgr.GetTrashPath(), 0755)
	_ = os.WriteFile(arcZero.Path, []byte("x"), 0644)
	out, _ := runListWithFlags(t, func(f *pflag.FlagSet) { _ = f.Set("within-days", "0") })
	if !strings.Contains(out, zeroName) {
		t.Fatalf("expected 0-day item listed: %s", out)
	}

	// 2) Purge --all includes items without DeletedAt
	noDel := &storage.Archive{
		UID:          "uid-nodel",
		Name:         "no-deleted-at.7z",
		Path:         filepath.Join(mgr.GetTrashPath(), "no-deleted-at.7z"),
		Created:      now.Add(-10 * 24 * time.Hour),
		Managed:      true,
		Status:       "deleted",
		DeletedAt:    nil,
		OriginalPath: filepath.Join(mgr.GetArchivesPath(), "no-deleted-at.7z"),
	}
	_ = mgr.Registry().Add(noDel)
	outP, err := runPurgeWithFlags(t, func(f *pflag.FlagSet) {
		_ = f.Set("all", "true")
		_ = f.Set("dry-run", "true")
		_ = f.Set("force", "true")
	})
	if err != nil {
		t.Fatalf("purge --all: %v", err)
	}
	if !strings.Contains(outP, "no-deleted-at.7z") {
		t.Fatalf("--all should include missing DeletedAt item: %s", outP)
	}

	// 3a) Restore fallback to managed when OriginalPath empty and file exists in trash
	fallbackName := "fallback-managed.7z"
	fallbackArc := &storage.Archive{
		UID:          "uid-fallback",
		Name:         fallbackName,
		Path:         filepath.Join(mgr.GetTrashPath(), fallbackName),
		Created:      now.Add(-2 * 24 * time.Hour),
		Managed:      true,
		Status:       "deleted",
		DeletedAt:    &now,
		OriginalPath: "",
	}
	_ = mgr.Registry().Add(fallbackArc)
	_ = os.WriteFile(fallbackArc.Path, []byte("data"), 0644)
	outRF, err := execRestore(t, fallbackArc.UID)
	if err != nil {
		t.Fatalf("restore fallback: %v; out=%s", err, outRF)
	}
	updatedF, _ := mgr.Registry().Get(fallbackName)
	if updatedF.Status != "present" {
		t.Fatalf("expected present after restore fallback: %+v", updatedF)
	}
	if expected := mgr.GetManagedPath(fallbackName); updatedF.Path != expected {
		t.Fatalf("expected managed path %s, got %s", expected, updatedF.Path)
	}

	// 3b) Missing file in trash -> expect failure
	missingName := "missing-in-trash.7z"
	missingArc := &storage.Archive{
		UID:          "uid-missing",
		Name:         missingName,
		Path:         filepath.Join(mgr.GetTrashPath(), missingName),
		Created:      now.Add(-2 * 24 * time.Hour),
		Managed:      true,
		Status:       "deleted",
		DeletedAt:    &now,
		OriginalPath: "",
	}
	_ = mgr.Registry().Add(missingArc)
	outR, err := execRestore(t, missingArc.UID)
	if err == nil {
		t.Fatalf("expected error when restoring missing trash file; out=%s", outR)
	}
}
