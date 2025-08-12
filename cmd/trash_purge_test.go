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

func runPurgeWithFlags(t *testing.T, setFlags func(*pflag.FlagSet)) (string, error) {
	t.Helper()
	cmd := trashPurgeCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if setFlags != nil {
		setFlags(cmd.Flags())
	}
	if cmd.RunE != nil {
		err := cmd.RunE(cmd, []string{})
		return buf.String(), err
	}
	return buf.String(), nil
}

func TestTrashPurgeEligibilityAndDryRun(t *testing.T) {
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
	mk := func(name string, delAgeDays int) {
		when := now.Add(time.Duration(-delAgeDays) * 24 * time.Hour)
		arc := &storage.Archive{
			UID:          "uid-" + name,
			Name:         name,
			Path:         filepath.Join(mgr.GetTrashPath(), name),
			Created:      now.Add(-30 * 24 * time.Hour),
			Managed:      true,
			Status:       "deleted",
			DeletedAt:    &when,
			OriginalPath: filepath.Join(mgr.GetArchivesPath(), name),
		}
		_ = mgr.Registry().Add(arc)
	}
	_ = os.MkdirAll(mgr.GetTrashPath(), 0755)
	mk("soon.7z", 6)    // 1 day left
	mk("overdue.7z", 9) // overdue

	// within-days=2 should target "soon" and "overdue" (<=2 includes overdue)
	out, err := runPurgeWithFlags(t, func(f *pflag.FlagSet) {
		_ = f.Set("within-days", "2")
		_ = f.Set("dry-run", "true")
		_ = f.Set("force", "true")
	})
	if err != nil {
		t.Fatalf("purge within dry-run: %v", err)
	}
	if !strings.Contains(out, "Would purge") || !strings.Contains(out, "soon.7z") || !strings.Contains(out, "overdue.7z") {
		t.Fatalf("dry-run summary missing items: %s", out)
	}

	// default purge (no flags) should purge overdue only
	out2, err := runPurgeWithFlags(t, func(f *pflag.FlagSet) { _ = f.Set("force", "true") })
	if err != nil {
		t.Fatalf("purge default: %v", err)
	}
	if !strings.Contains(out2, "Purged 1 archives.") {
		t.Fatalf("expected purge count 1, got: %s", out2)
	}
}
