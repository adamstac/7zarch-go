package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func MasMoveCmd() *cobra.Command {
	var to string
	cmd := &cobra.Command{
		Use:   "move <id>",
		Short: "Move an archive (default to managed storage if --to omitted)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil { return fmt.Errorf("failed to init storage: %w", err) }
			defer mgr.Close()

			resolver := storage.NewResolver(mgr.Registry())
			arc, err := resolver.Resolve(id)
			if err != nil { return err }

			dest := to
			if dest == "" {
				dest = mgr.GetManagedPath(arc.Name)
			}
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil { return err }
			if err := os.Rename(arc.Path, dest); err != nil { return err }

			arc.Path = dest
			arc.Managed = filepath.HasPrefix(dest, mgr.GetBasePath())
			return mgr.Registry().Update(arc)
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Destination path or managed default if omitted")
	return cmd
}

