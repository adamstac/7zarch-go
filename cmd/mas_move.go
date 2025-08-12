package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			// If dest is an existing directory, place the file under it by name
			if info, err := os.Stat(dest); err == nil && info.IsDir() {
				dest = filepath.Join(dest, arc.Name)
			}
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil { return err }
			if err := os.Rename(arc.Path, dest); err != nil { return err }

			arc.Path = dest
			rel, err := filepath.Rel(mgr.GetBasePath(), dest)
			if err != nil { return err }
			arc.Managed = !strings.HasPrefix(rel, "..")
			return mgr.Registry().Update(arc)
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Destination path or managed default if omitted")
	return cmd
}

