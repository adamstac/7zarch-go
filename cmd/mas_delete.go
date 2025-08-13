package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	errs "github.com/adamstac/7zarch-go/internal/errors"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func MasDeleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an archive (soft by default; --force to remove file)",
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
				if _, ok := err.(*storage.AmbiguousIDError); ok {
					return &errs.ValidationError{
						Field:   "archive ID",
						Value:   id,
						Message: "matches multiple archives. Use a longer prefix or full UID",
					}
				}
				if _, ok := err.(*storage.ArchiveNotFoundError); ok {
					return errs.NewArchiveNotFound(id)
				}
				return err
			}
			
			// Check if already deleted
			if arc.Status == "deleted" {
				return &errs.InvalidOperationError{
					Operation: "delete",
					Resource:  "archive",
					Reason:    "archive is already deleted",
				}
			}

			now := time.Now()
			orig := arc.Path

			if force {
				// Physically remove file if present
				_ = os.Remove(arc.Path)
				arc.Status = "deleted"
				arc.DeletedAt = &now
				if arc.OriginalPath == "" {
					arc.OriginalPath = orig
				}
				return mgr.Registry().Update(arc)
			}

			// Soft delete
			if arc.Managed {
				// Move to managed trash directory
				trashDir := mgr.GetTrashPath()
				// #nosec G301: restrict permissions on created trash directory
				if err := os.MkdirAll(trashDir, 0750); err != nil {
					return fmt.Errorf("failed to create trash: %w", err)
				}
				trashPath := filepath.Join(trashDir, filepath.Base(arc.Path))
				if err := moveOrCopy(arc.Path, trashPath); err != nil {
					return fmt.Errorf("failed to move to trash: %w", err)
				}
				arc.Path = trashPath
			} else {
				// External: default DB-only delete (do not touch file)
			}
			arc.Status = "deleted"
			arc.DeletedAt = &now
			if arc.OriginalPath == "" {
				arc.OriginalPath = orig
			}
			return mgr.Registry().Update(arc)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Physically remove file instead of soft delete")
	return cmd
}

// moveOrCopy tries to rename; if it fails (e.g., cross-device), it copies then removes
func moveOrCopy(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	// fallback to copy + remove
	// #nosec G304: src and dst are derived from managed paths within storage; validated upstream
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	// Use restrictive permissions for destination file
	// #nosec G304: destination is inside managed trash or archives path
	dstF, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer dstF.Close()
	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}
	return os.Remove(src)
}
