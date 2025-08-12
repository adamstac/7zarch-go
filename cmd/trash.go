package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
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
		flagJSON       bool
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

			if flagJSON {
				// Enrich with countdown
				type row struct {
					UID       string     `json:"uid"`
					Name      string     `json:"name"`
					Path      string     `json:"path"`
					DeletedAt *time.Time `json:"deleted_at"`
					PurgeDate string     `json:"purge_date"`
					DaysLeft  int        `json:"days_left"`
				}
				rows := make([]row, 0, len(out))
				for _, a := range out {
					var purgeStr string
					days := -1
					if a.DeletedAt != nil {
						purge := a.DeletedAt.Add(time.Duration(cfg.Storage.RetentionDays) * 24 * time.Hour)
						purgeStr = purge.Format("2006-01-02")
						days = int(time.Until(purge).Hours() / 24)
					}
					rows = append(rows, row{UID: a.UID, Name: a.Name, Path: a.Path, DeletedAt: a.DeletedAt, PurgeDate: purgeStr, DaysLeft: days})
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(rows)
			}

			// Text output
			fmt.Printf("🗑️  Deleted archives (%d)\n", len(out))
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
				fmt.Printf("- %s (%s)\n  deleted: %s | purge: %s (%s)\n", a.Name, a.UID[:8], delStr, purgeStr, countdown)
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&flagWithinDays, "within-days", 0, "Show items purging within N days (0=all)")
	cmd.Flags().StringVar(&flagBefore, "before", "", "Only show items deleted before YYYY-MM-DD")
	cmd.Flags().BoolVar(&flagJSON, "json", false, "Output as JSON")
	return cmd
}

func trashPurgeCmd() *cobra.Command {
	var (
		flagAll    bool
		flagForce  bool
		flagDryRun bool
		flagWithin int
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
				fmt.Println("Nothing to purge.")
				return nil
			}

			// Confirmation unless --force
			if !flagForce {
				fmt.Printf("About to purge %d archives. Proceed? [y/N]: ", len(eligible))
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(strings.ToLower(line))
				if line != "y" && line != "yes" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			// Dry-run summary
			if flagDryRun {
				fmt.Printf("Would purge %d archives:\n", len(eligible))
				for _, a := range eligible {
					fmt.Printf("- %s (%s)\n", a.Name, a.UID[:8])
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
			fmt.Printf("Purged %d archives.\n", len(eligible))
			return nil
		},
	}
	cmd.Flags().BoolVar(&flagAll, "all", false, "Purge all trashed archives (ignore retention)")
	cmd.Flags().BoolVar(&flagForce, "force", false, "Skip confirmation prompts")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show actions without making changes")
	cmd.Flags().IntVar(&flagWithin, "within-days", 0, "Only purge items purging within N days (0=all)")
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
