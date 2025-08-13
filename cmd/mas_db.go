package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func MasDbCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "db", Short: "Database operations (status, migrate, backup)"}
	cmd.AddCommand(masDbStatusCmd())
	cmd.AddCommand(masDbMigrateCmd())
	cmd.AddCommand(masDbBackupCmd())
	return cmd
}

func masDbStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show database version and applied migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return err
			}
			defer mgr.Close()

			if err := mgr.Registry().EnsureMigrationsTable(); err != nil {
				return err
			}
			// For now, we show whether baseline/identity are marked
			baseApplied, _ := mgr.Registry().IsMigrationApplied("0001_baseline")
			idApplied, _ := mgr.Registry().IsMigrationApplied("0002_identity_and_status")
			fmt.Printf("DB Path: %s\n", mgr.Registry().Path())
			fmt.Printf("Baseline: %v\n", baseApplied)
			fmt.Printf("Identity/Status: %v\n", idApplied)
			return nil
		},
	}
}

func masDbMigrateCmd() *cobra.Command {
	var dryRun bool
	var backupOnly bool
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Apply pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return err
			}
			defer mgr.Close()

			if backupOnly {
				return createDbBackup(mgr)
			}
			if dryRun {
				fmt.Println("Dry run: no migrations applied")
				return nil
			}
			// For now, schema is current; future pending migrations will be run here
			fmt.Println("No pending migrations")
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done")
	cmd.Flags().BoolVar(&backupOnly, "backup-only", false, "Create a backup without migrating")
	return cmd
}

func masDbBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Create a timestamped backup of the registry database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return err
			}
			defer mgr.Close()
			return createDbBackup(mgr)
		},
	}
}

func createDbBackup(mgr *storage.Manager) error {
	path := mgr.Registry().Path()
	if path == "" {
		return fmt.Errorf("registry path unknown")
	}
	backupDir := filepath.Dir(path)
	stamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("registry.%s.bak", stamp))

	// #nosec G304: path comes from validated config; used for local backup
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()
	// Use restrictive permissions for backup file
	// #nosec G304: backupPath is created under the same directory as the DB
	dst, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	fmt.Printf("Backup created: %s\n", backupPath)
	return nil
}
