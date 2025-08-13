package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// TrashCmd groups trash-related subcommands.
func TrashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trash",
		Short: "Manage deleted archives in trash",
	}
	cmd.AddCommand(trashListCmd())
	cmd.AddCommand(trashPurgeCmd())
	return cmd
}

func trashListCmd() *cobra.Command {
	var (
		flagWithinDays int
		flagBefore     string
		flagOutput     string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List deleted archives and purge eligibility",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return fmt.Errorf("failed to init storage: %w", err)
			}
			defer mgr.Close()

			archives, err := mgr.List()
			if err != nil {
				return err
			}
			var out []*storage.Archive
			for _, a := range archives {
				if a.Status != "deleted" {
					continue
				}
				if flagBefore != "" && a.DeletedAt != nil {
					cutoff, err := parseYMD(flagBefore)
					if err != nil {
						return fmt.Errorf("invalid --before date: %w", err)
					}
					if a.DeletedAt.After(cutoff) {
						continue
					}
				}
				if flagWithinDays > 0 && a.DeletedAt != nil {
					purge := a.DeletedAt.Add(time.Duration(cfg.Storage.RetentionDays) * 24 * time.Hour)
					days := int(time.Until(purge).Hours() / 24)
					if days > flagWithinDays {
						continue
					}
				}
				out = append(out, a)
			}

			if flagOutput != "" {
				return outputTrashList(out, flagOutput, cfg.Storage.RetentionDays, cmd.OutOrStdout())
			}

			// Text output
			outWriter := cmd.OutOrStdout()
			fmt.Fprintf(outWriter, "🗑️  Deleted archives (%d)\n", len(out))
			for _, a := range out {
				var delStr, purgeStr, countdown string
				if a.DeletedAt != nil {
					delStr = a.DeletedAt.Format("2006-01-02")
					purge := a.DeletedAt.Add(time.Duration(cfg.Storage.RetentionDays) * 24 * time.Hour)
					purgeStr = purge.Format("2006-01-02")
					days := int(time.Until(purge).Hours() / 24)
					if days < 0 {
						countdown = "overdue"
					} else {
						countdown = fmt.Sprintf("%dd", days)
					}
				}
				fmt.Fprintf(outWriter, "- %s (%s)\n  deleted: %s | purge: %s (%s)\n", a.Name, a.UID[:8], delStr, purgeStr, countdown)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&flagWithinDays, "within-days", 0, "Show items purging within N days (0=all)")
	cmd.Flags().StringVar(&flagBefore, "before", "", "Only show items deleted before YYYY-MM-DD")
	cmd.Flags().StringVar(&flagOutput, "output", "", "Output format: json|csv|yaml (default: human-readable)")
	return cmd
}

func trashPurgeCmd() *cobra.Command {
	var (
		flagAll       bool
		flagForce     bool
		flagDryRun    bool
		flagWithin    int
		flagOlderThan string
	)
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Permanently delete trashed archives past retention",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return fmt.Errorf("failed to init storage: %w", err)
			}
			defer mgr.Close()

			archives, err := mgr.List()
			if err != nil {
				return err
			}
			var eligible []*storage.Archive
			now := time.Now()
			for _, a := range archives {
				if a.Status != "deleted" {
					continue
				}
				if a.DeletedAt == nil && !flagAll {
					continue
				}
				purgeDate := time.Time{}
				if a.DeletedAt != nil {
					purgeDate = a.DeletedAt.Add(time.Duration(cfg.Storage.RetentionDays) * 24 * time.Hour)
				}
				// Eligibility
				if flagAll {
					eligible = append(eligible, a)
					continue
				}
				if flagOlderThan != "" && a.DeletedAt != nil {
					cutoff, err := parseOlderThan(flagOlderThan)
					if err != nil {
						return fmt.Errorf("invalid --older-than format: %w", err)
					}
					if a.DeletedAt.Before(cutoff) {
						eligible = append(eligible, a)
					}
					continue
				}
				if flagWithin > 0 {
					days := int(purgeDate.Sub(now).Hours() / 24)
					if days <= flagWithin {
						eligible = append(eligible, a)
					}
					continue
				}
				if !purgeDate.IsZero() && !purgeDate.After(now) {
					eligible = append(eligible, a)
				}
			}

			if len(eligible) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Nothing to purge.")
				return nil
			}

			// Confirmation unless --force
			if !flagForce {
				fmt.Fprintf(cmd.OutOrStdout(), "About to purge %d archives. Proceed? [y/N]: ", len(eligible))
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(strings.ToLower(line))
				if line != "y" && line != "yes" {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
					return nil
				}
			}

			// Dry-run summary
			if flagDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would purge %d archives:\n", len(eligible))
				for _, a := range eligible {
					fmt.Fprintf(cmd.OutOrStdout(), "- %s (%s)\n", a.Name, a.UID[:8])
				}
				return nil
			}

			trashDir := mgr.GetTrashPath()
			for _, a := range eligible {
				// Remove file if it appears under trash (managed archives)
				if a.Managed && strings.HasPrefix(a.Path, trashDir+string(os.PathSeparator)) {
					_ = os.Remove(a.Path)
				}
				// Remove from registry
				_ = mgr.Registry().Delete(a.Name)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Purged %d archives.\n", len(eligible))
			return nil
		},
	}
	cmd.Flags().BoolVar(&flagAll, "all", false, "Purge all trashed archives (ignore retention)")
	cmd.Flags().BoolVar(&flagForce, "force", false, "Skip confirmation prompts")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show actions without making changes")
	cmd.Flags().IntVar(&flagWithin, "within-days", 0, "Only purge items purging within N days (0=all)")
	cmd.Flags().StringVar(&flagOlderThan, "older-than", "", "Purge archives deleted older than duration (e.g., '30d', '1w')")
	return cmd
}

// parseYMD parses YYYY-MM-DD into time at midnight local
func parseYMD(s string) (time.Time, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("expected YYYY-MM-DD")
	}
	y, m, d := parts[0], parts[1], parts[2]
	year, err := strconv.Atoi(y)
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(m)
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(d)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// parseOlderThan parses duration strings like "30d", "1w" and returns the cutoff time
func parseOlderThan(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty duration")
	}
	if strings.HasSuffix(s, "d") {
		n, err := strconv.ParseInt(strings.TrimSuffix(s, "d"), 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid days: %w", err)
		}
		return time.Now().Add(-time.Duration(n) * 24 * time.Hour), nil
	}
	if strings.HasSuffix(s, "w") {
		n, err := strconv.ParseInt(strings.TrimSuffix(s, "w"), 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid weeks: %w", err)
		}
		return time.Now().Add(-time.Duration(n) * 7 * 24 * time.Hour), nil
	}
	// Try standard time.ParseDuration
	dur, err := time.ParseDuration(s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid duration format: %w", err)
	}
	return time.Now().Add(-dur), nil
}

// trashRow represents a row in trash list output
type trashRow struct {
	UID       string     `json:"uid" yaml:"uid" csv:"uid"`
	Name      string     `json:"name" yaml:"name" csv:"name"`
	Path      string     `json:"path" yaml:"path" csv:"path"`
	DeletedAt *time.Time `json:"deleted_at" yaml:"deleted_at" csv:"deleted_at"`
	PurgeDate string     `json:"purge_date" yaml:"purge_date" csv:"purge_date"`
	DaysLeft  int        `json:"days_left" yaml:"days_left" csv:"days_left"`
}

// outputTrashList outputs trash list in machine-readable format
func outputTrashList(archives []*storage.Archive, format string, retentionDays int, out io.Writer) error {

	rows := make([]trashRow, 0, len(archives))
	for _, a := range archives {
		var purgeStr string
		days := -1
		if a.DeletedAt != nil {
			purge := a.DeletedAt.Add(time.Duration(retentionDays) * 24 * time.Hour)
			purgeStr = purge.Format("2006-01-02")
			days = int(time.Until(purge).Hours() / 24)
		}
		rows = append(rows, trashRow{
			UID:       a.UID,
			Name:      a.Name,
			Path:      a.Path,
			DeletedAt: a.DeletedAt,
			PurgeDate: purgeStr,
			DaysLeft:  days,
		})
	}

	switch format {
	case "json":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	case "csv":
		return outputTrashCSV(rows, out)
	case "yaml":
		enc := yaml.NewEncoder(out)
		defer enc.Close()
		return enc.Encode(rows)
	default:
		return fmt.Errorf("unsupported output format: %s (supported: json, csv, yaml)", format)
	}
}

func outputTrashCSV(rows []trashRow, out io.Writer) error {
	writer := csv.NewWriter(out)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{
		"uid", "name", "path", "deleted_at", "purge_date", "days_left",
	}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, row := range rows {
		var deletedAt string
		if row.DeletedAt != nil {
			deletedAt = row.DeletedAt.Format(time.RFC3339)
		}

		csvRow := []string{
			row.UID,
			row.Name,
			row.Path,
			deletedAt,
			row.PurgeDate,
			fmt.Sprintf("%d", row.DaysLeft),
		}

		if err := writer.Write(csvRow); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}
