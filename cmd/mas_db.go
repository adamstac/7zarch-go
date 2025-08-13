package cmd

import (
	"fmt"
	"os"

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

			runner := mgr.NewMigrationRunner()

			applied, err := runner.GetAppliedMigrations()
			if err != nil {
				return fmt.Errorf("failed to get applied migrations: %w", err)
			}

			pending, err := runner.GetPendingMigrations()
			if err != nil {
				return fmt.Errorf("failed to get pending migrations: %w", err)
			}

			fmt.Printf("Database: %s\n", mgr.Registry().Path())
			
			if len(applied) > 0 {
				latestMigration := applied[len(applied)-1]
				fmt.Printf("Schema Version: %s\n", latestMigration.ID)
			} else {
				fmt.Printf("Schema Version: (none applied)\n")
			}

			fmt.Printf("Applied Migrations: %d\n", len(applied))
			for _, migration := range applied {
				fmt.Printf("  ✓ %s: %s (applied %s)\n", 
					migration.ID, 
					migration.Name, 
					migration.AppliedAt.Format("2006-01-02 15:04:05"))
			}

			if len(pending) > 0 {
				fmt.Printf("Pending Migrations: %d\n", len(pending))
				for _, migration := range pending {
					fmt.Printf("  - %s: %s\n", migration.ID, migration.Description)
				}
			} else {
				fmt.Printf("Pending Migrations: 0\n")
			}

			// Get database size
			if stat, err := os.Stat(mgr.Registry().Path()); err == nil {
				fmt.Printf("Database Size: %.1f KB\n", float64(stat.Size())/1024)
			}

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

			runner := mgr.NewMigrationRunner()

			if backupOnly {
				backupPath, err := runner.CreateBackup(mgr.Registry().Path())
				if err != nil {
					return fmt.Errorf("backup failed: %w", err)
				}
				fmt.Printf("Backup created: %s\n", backupPath)
				return nil
			}

			pending, err := runner.GetPendingMigrations()
			if err != nil {
				return fmt.Errorf("failed to get pending migrations: %w", err)
			}

			if len(pending) == 0 {
				fmt.Println("No pending migrations")
				return nil
			}

			if dryRun {
				fmt.Printf("Dry run: would apply %d migration(s)\n", len(pending))
				for _, migration := range pending {
					fmt.Printf("  - %s: %s\n", migration.ID, migration.Description)
				}
				return nil
			}

			fmt.Printf("Applying %d pending migration(s)...\n", len(pending))

			if err := runner.ApplyPending(mgr.Registry().Path()); err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}

			fmt.Println("✓ Migrations completed successfully")
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

			runner := mgr.NewMigrationRunner()
			backupPath, err := runner.CreateBackup(mgr.Registry().Path())
			if err != nil {
				return fmt.Errorf("backup failed: %w", err)
			}
			fmt.Printf("Backup created: %s\n", backupPath)
			return nil
		},
	}
}
