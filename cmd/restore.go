package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

// RestoreCmd returns the `restore` command.
// Usage: 7zarch-go restore <id>
func RestoreCmd() *cobra.Command {
	var (
		flagForce  bool
		flagDryRun bool
	)

	cmd := &cobra.Command{
		Use:   "restore <id>",
		Short: "Restore a deleted archive from trash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return fmt.Errorf("failed to init storage: %w", err)
			}
			defer mgr.Close()

			resolver := storage.NewResolver(mgr.Registry())
			arc, err := resolver.Resolve(id)
			if err != nil {
				return err
			}
			if arc.Status != "deleted" {
				return fmt.Errorf("archive '%s' is not deleted (status=%s)", arc.Name, arc.Status)
			}

			target := arc.OriginalPath
			if target == "" {
				name := arc.Name
				if name == "" {
					name = filepath.Base(arc.Path)
				}
				target = mgr.GetManagedPath(name)
			}

			// Plan
			if flagDryRun {
				cmd.Printf("Would restore %s -> %s\n", arc.Path, target)
				return nil
			}

			// Ensure parent dir exists
			// #nosec G301: restrict permissions on restored directory
			if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
				return fmt.Errorf("failed to prepare destination: %w", err)
			}

			// If managed archive, file lives in trash and needs moving back
			if arc.Managed {
				if err := moveOrCopy(arc.Path, target); err != nil {
					return fmt.Errorf("failed to restore file: %w", err)
				}
				arc.Path = target
			} else {
				// External soft delete: file likely remained in place; just flip status
				arc.Path = target
			}

			// Update registry state
			arc.Status = "present"
			arc.DeletedAt = nil
			now := time.Now()
			arc.LastSeen = &now
			if err := mgr.Registry().Update(arc); err != nil {
				return fmt.Errorf("failed to update registry: %w", err)
			}
			cmd.Printf("✅ Restored %s to %s\n", arc.Name, target)
			return nil
		},
	}

	cmd.Flags().BoolVar(&flagForce, "force", false, "Overwrite existing file if it exists at original location")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show actions without making changes")

	return cmd
}
