package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/adamstac/7zarch-go/internal/cmdutil"
	errs "github.com/adamstac/7zarch-go/internal/errors"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

// copyFile copies src to dst with mode preservation
func copyFile(src, dst string) error {
	// #nosec G304: src/dst validated via resolver and managed paths
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Use restrictive permissions for new file
	// #nosec G304: destination path validated/constructed earlier and rooted within managed base when applicable
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Sync(); err != nil {
		return err
	}
	if info, err := os.Stat(src); err == nil {
		_ = os.Chmod(dst, info.Mode())
	}
	return nil
}

func MasMoveCmd() *cobra.Command {
	var to string
	cmd := &cobra.Command{
		Use:   "move <id>",
		Short: "Move an archive (default to managed storage if --to omitted)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			_, mgr, cleanup, err := cmdutil.InitStorageManager()
			if err != nil {
				return err
			}
			defer cleanup()

			resolver := storage.NewResolver(mgr.Registry())
			arc, err := resolver.Resolve(id)
			if err != nil {
				return cmdutil.HandleResolverError(err, id)
			}
			
			// Check if archive can be moved
			if arc.Status == "deleted" {
				return &errs.InvalidOperationError{
					Operation: "move",
					Resource:  "archive",
					Reason:    "archive is deleted",
				}
			}

			dest := to
			if dest == "" {

				name := arc.Name
				if name == "" {
					name = filepath.Base(arc.Path)
				}
				dest = mgr.GetManagedPath(name)
			}

			// If dest is an existing directory, place the file under it by name
			if info, err := os.Stat(dest); err == nil && info.IsDir() {
				name := arc.Name
				if name == "" {
					name = filepath.Base(arc.Path)
				}
				dest = filepath.Join(dest, name)
			}

			// #nosec G301: restrict permissions on created directory
			if err := os.MkdirAll(filepath.Dir(dest), 0750); err != nil {
				return err
			}

			// Prevent accidental overwrite
			if info, err := os.Stat(dest); err == nil && !info.IsDir() {
				return fmt.Errorf("destination file already exists: %s", dest)
			}

			if err := os.Rename(arc.Path, dest); err != nil {
				// Handle cross-device rename (EXDEV)
				var linkErr *os.LinkError
				if errors.As(err, &linkErr) && errors.Is(linkErr.Err, syscall.EXDEV) {
					if err := copyFile(arc.Path, dest); err != nil {
						return fmt.Errorf("copy fallback failed: %w", err)
					}
					if err := os.Remove(arc.Path); err != nil {
						return fmt.Errorf("cleanup source failed after copy: %w", err)
					}
				} else {
					return err
				}
			}

			arc.Path = dest
			// More precise managed-path check
			rel, _ := filepath.Rel(mgr.GetBasePath(), dest)
			up := ".." + string(os.PathSeparator)
			arc.Managed = rel != ".." && !strings.HasPrefix(rel, up)

			return mgr.Registry().Update(arc)
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Destination path or managed default if omitted")
	return cmd
}
