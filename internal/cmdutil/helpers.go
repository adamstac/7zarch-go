package cmdutil

import (
	"fmt"

	"github.com/adamstac/7zarch-go/internal/config"
	errs "github.com/adamstac/7zarch-go/internal/errors"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// InitStorageManager is a helper that loads config and creates storage manager
// Returns config, storage manager, and cleanup function
func InitStorageManager() (*config.Config, *storage.Manager, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return cfg, nil, nil, fmt.Errorf("failed to initialize storage manager: %w", err)
	}

	cleanup := func() {
		mgr.Close()
	}

	return cfg, mgr, cleanup, nil
}

// LoadConfigOrDefault attempts to load config but falls back to defaults
// Used in create command and similar cases where config errors are non-fatal
func LoadConfigOrDefault() *config.Config {
	if cfg, err := config.Load(); err == nil {
		return cfg
	}
	return config.DefaultConfig()
}

// HandleResolverError converts storage resolver errors to consistent user-friendly errors
func HandleResolverError(err error, id string) error {
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