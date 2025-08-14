package batch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
)

type ProgressUpdate struct {
	Total     int
	Completed int
	Errors    []error
	Current   string
	Elapsed   time.Duration
}

type ProgressCallback func(ProgressUpdate)

// ManagerInterface defines the interface needed by the batch processor
type ManagerInterface interface {
	Delete(uid string) error
	Registry() *storage.Registry
	GetBasePath() string
}

type Processor struct {
	manager    ManagerInterface
	concurrent int
	mu         sync.RWMutex
}

func NewProcessor(manager ManagerInterface) *Processor {
	return &Processor{
		manager:    manager,
		concurrent: 4, // Default concurrency
	}
}

func (p *Processor) SetConcurrency(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if n > 0 {
		p.concurrent = n
	}
}

func (p *Processor) Move(ctx context.Context, archives []*storage.Archive, destPath string, callback ProgressCallback) error {
	return p.processWithProgress(ctx, archives, callback, func(archive *storage.Archive) error {
		return p.moveArchive(archive, destPath)
	})
}

func (p *Processor) Delete(ctx context.Context, archives []*storage.Archive, callback ProgressCallback) error {
	return p.processWithProgress(ctx, archives, callback, func(archive *storage.Archive) error {
		return p.manager.Delete(archive.UID)
	})
}

// moveArchive implements the move logic from cmd/mas_move.go
func (p *Processor) moveArchive(archive *storage.Archive, destBasePath string) error {
	// Check if archive can be moved
	if archive.Status == "deleted" {
		return fmt.Errorf("archive %s is deleted", archive.UID)
	}

	// Determine destination path
	dest := destBasePath
	name := archive.Name
	if name == "" {
		name = filepath.Base(archive.Path)
	}

	// If dest is a directory, place the file under it by name
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		dest = filepath.Join(dest, name)
	}

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dest), 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(dest), err)
	}

	// Prevent accidental overwrite
	if info, err := os.Stat(dest); err == nil && !info.IsDir() {
		return fmt.Errorf("destination file already exists: %s", dest)
	}

	// Attempt rename first (fastest for same filesystem)
	if err := os.Rename(archive.Path, dest); err != nil {
		// Handle cross-device rename (EXDEV) with copy+remove fallback
		var linkErr *os.LinkError
		if errors.As(err, &linkErr) && errors.Is(linkErr.Err, syscall.EXDEV) {
			if err := p.copyFile(archive.Path, dest); err != nil {
				return fmt.Errorf("copy fallback failed: %w", err)
			}
			if err := os.Remove(archive.Path); err != nil {
				return fmt.Errorf("cleanup source failed after copy: %w", err)
			}
		} else {
			return fmt.Errorf("failed to move %s to %s: %w", archive.Path, dest, err)
		}
	}

	// Update registry with new path and managed status
	archive.Path = dest
	// Check if destination is under managed storage
	rel, _ := filepath.Rel(p.manager.GetBasePath(), dest)
	up := ".." + string(os.PathSeparator)
	archive.Managed = rel != ".." && !strings.HasPrefix(rel, up)

	return p.manager.Registry().Update(archive)
}

// copyFile copies src to dst with mode preservation
func (p *Processor) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

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

func (p *Processor) processWithProgress(ctx context.Context, archives []*storage.Archive, callback ProgressCallback, operation func(*storage.Archive) error) error {
	total := len(archives)
	if total == 0 {
		return nil
	}

	startTime := time.Now()
	completed := 0
	var errors []error
	var mu sync.Mutex

	updateProgress := func(current string) {
		if callback != nil {
			mu.Lock()
			update := ProgressUpdate{
				Total:     total,
				Completed: completed,
				Errors:    append([]error(nil), errors...),
				Current:   current,
				Elapsed:   time.Since(startTime),
			}
			mu.Unlock()
			callback(update)
		}
	}

	// Create worker pool
	p.mu.RLock()
	concurrency := p.concurrent
	p.mu.RUnlock()

	archiveChan := make(chan *storage.Archive, total)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for archive := range archiveChan {
				select {
				case <-ctx.Done():
					return
				default:
					current := fmt.Sprintf("Processing %s", archive.Name)
					updateProgress(current)

					if err := operation(archive); err != nil {
						mu.Lock()
						errors = append(errors, fmt.Errorf("failed to process %s: %w", archive.Name, err))
						mu.Unlock()
					}

					mu.Lock()
					completed++
					mu.Unlock()
				}
			}
		}()
	}

	// Send work to workers
	go func() {
		defer close(archiveChan)
		for _, archive := range archives {
			select {
			case <-ctx.Done():
				return
			case archiveChan <- archive:
			}
		}
	}()

	// Progress reporting goroutine
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan struct{})

	go func() {
		defer close(done)
		wg.Wait()
	}()

	// Wait for completion or context cancellation
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			updateProgress("In progress...")
		case <-done:
			// Final progress update
			updateProgress("Complete")
			if len(errors) > 0 {
				return fmt.Errorf("batch operation completed with %d errors (see progress for details)", len(errors))
			}
			return nil
		}
	}
}