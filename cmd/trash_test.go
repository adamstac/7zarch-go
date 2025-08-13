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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// helper to invoke RunE directly with flags
func runEWithFlags(t *testing.T, cmd *cobra.Command, setFlags func(f *pflag.FlagSet)) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if setFlags != nil {
		setFlags(cmd.Flags())
	}
	// Directly execute RunE to avoid root/subcommand wiring
	if cmd.RunE != nil {
		err := cmd.RunE(cmd, []string{})
		return buf.String(), err
	}
	return buf.String(), nil
}

func TestParseYMD(t *testing.T) {
	got, err := parseYMD("2025-01-15")
	if err != nil {
		t.Fatalf("parseYMD error: %v", err)
	}
	if got.Year() != 2025 || int(got.Month()) != 1 || got.Day() != 15 {
		t.Fatalf("parseYMD wrong date: %v", got)
	}
	if _, err := parseYMD("2025/01/15"); err == nil {
		t.Fatalf("expected error for bad format")
	}
}

func TestTrashListFiltersAndJSON(t *testing.T) {
	// Set up a temp managed store and config
	base := t.TempDir()
	tempDir := base
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := config.DefaultConfig()
	cfg.Storage.ManagedPath = tempDir
	cfg.Storage.RetentionDays = 7

	// Write YAML config at $HOME/.7zarch-go-config so config.Load picks it up
	cfgPath := filepath.Join(home, ".7zarch-go-config")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("yaml marshal: %v", err)
	}
	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		t.Fatalf("write cfg: %v", err)
	}

	// Create manager directly for seeding data
	mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		t.Fatalf("manager error: %v", err)
	}
	defer mgr.Close()

	now := time.Now()
	mk := func(name string, delAgeDays int, status string) *storage.Archive {
		arc := &storage.Archive{
			UID:          "uid-" + name,
			Name:         name,
			Path:         filepath.Join(mgr.GetTrashPath(), name),
			Size:         100,
			Created:      now.Add(-30 * 24 * time.Hour),
			Managed:      true,
			Status:       status,
			DeletedAt:    nil,
			OriginalPath: filepath.Join(mgr.GetArchivesPath(), name),
		}
		if delAgeDays >= 0 {
			when := now.Add(time.Duration(-delAgeDays) * 24 * time.Hour)
			arc.DeletedAt = &when
		}
		if err := mgr.Registry().Add(arc); err != nil {
			t.Fatalf("seed add: %v", err)
		}
		return arc
	}

	// Seed: one purging soon (in 3 days), one overdue, one far future, one not-deleted
	_ = os.MkdirAll(mgr.GetTrashPath(), 0755)
	mk("soon.7z", 4, "deleted")     // deleted 4 days ago -> purge at 7 -> 3 days left
	mk("overdue.7z", 10, "deleted") // deleted 10 days ago -> overdue
	mk("future.7z", 0, "deleted")   // deleted today -> 7 days left
	mk("present.7z", -1, "present") // not deleted

	// within-days=3 should include 'soon' and 'overdue'
	listCmd := trashListCmd()
	out, err := runEWithFlags(t, listCmd, func(f *pflag.FlagSet) { _ = f.Set("within-days", "3") })
	if err != nil {
		t.Fatalf("list within-days error: %v", err)
	}
	// The CLI prints header lines; just assert names appear
	if !containsAll(out, []string{"soon.7z", "overdue.7z"}) {
		t.Fatalf("within-days output missing expected items:\n%s", out)
	}
	if strings.Contains(out, "future.7z") {
		t.Fatalf("within-days output should not contain future.7z:\n%s", out)
	}

	// before filter: only items deleted before cutoff (overdue only)
	listCmd2 := trashListCmd()
	out2, err := runEWithFlags(t, listCmd2, func(f *pflag.FlagSet) { _ = f.Set("before", now.Add(-5*24*time.Hour).Format("2006-01-02")) })
	if err != nil {
		t.Fatalf("list before error: %v", err)
	}
	if !containsAll(out2, []string{"overdue.7z"}) || strings.Contains(out2, "soon.7z") || strings.Contains(out2, "future.7z") {
		t.Fatalf("before output mismatch:\n%s", out2)
	}

	// JSON mode must include DaysLeft and PurgeDate
	listCmd3 := trashListCmd()
	out3, err := runEWithFlags(t, listCmd3, func(f *pflag.FlagSet) { _ = f.Set("output", "json") })
	if err != nil {
		t.Fatalf("list json error: %v", err)
	}
	if !strings.Contains(out3, "\"days_left\"") || !strings.Contains(out3, "\"purge_date\"") {
		t.Fatalf("json output missing fields: %s", out3)
	}
}

// helpers
func containsAll(hay string, needles []string) bool {
	for _, n := range needles {
		if !strings.Contains(hay, n) {
			return false
		}
	}
	return true
}
