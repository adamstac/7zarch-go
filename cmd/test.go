package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/adamstac/7zarch-go/internal/archive"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	testRemote    bool
	testDirectory bool
	maxConcurrent int
)

func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <archive|directory>",
		Short: "Test archive integrity",
		Long: `Test the integrity of archives by verifying structure, checksums, and metadata.
Can test single archives or entire directories concurrently.`,
		Args: cobra.ExactArgs(1),
		RunE: runTest,
	}

	// Add flags
	cmd.Flags().BoolVar(&testRemote, "remote", false, "Run tests on TrueNAS server")
	cmd.Flags().BoolVarP(&testDirectory, "directory", "d", false, "Test all archives in directory")
	cmd.Flags().IntVar(&maxConcurrent, "concurrent", 10, "Max concurrent tests")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be tested")

	return cmd
}

func runTest(cmd *cobra.Command, args []string) error {
	target := args[0]

	if dryRun {
		return runTestDryRun(target)
	}

	if testDirectory {
		return runTestDirectory(target)
	}

	return runTestSingle(target)
}

func runTestDryRun(target string) error {
	fmt.Printf("DRY RUN MODE - No tests will be executed\n\n")

	if testDirectory {
		// Find archives in directory
		archives, err := findArchives(target)
		if err != nil {
			return fmt.Errorf("failed to find archives: %w", err)
		}

		fmt.Printf("Would test %d archives in %s:\n", len(archives), target)
		for _, arch := range archives {
			fmt.Printf("  - %s\n", filepath.Base(arch))
		}
		fmt.Printf("\nTests to run:\n")
		fmt.Printf("  ✓ Archive structure integrity\n")
		fmt.Printf("  ✓ Checksum verification\n")
		fmt.Printf("  ✓ Metadata validation\n")
		fmt.Printf("  ✓ Extraction test\n")
		fmt.Printf("\nMax concurrent tests: %d\n", maxConcurrent)
	} else {
		fmt.Printf("Would test archive: %s\n", target)
		fmt.Printf("\nTests to run:\n")
		fmt.Printf("  ✓ Archive structure integrity\n")
		fmt.Printf("  ✓ Checksum verification\n")
		fmt.Printf("  ✓ Metadata validation\n")
		fmt.Printf("  ✓ Extraction test\n")
	}

	if testRemote {
		fmt.Printf("\nExecution mode: Remote (on TrueNAS)\n")
	} else {
		fmt.Printf("\nExecution mode: Local\n")
	}

	return nil
}

func runTestSingle(archivePath string) error {
	fmt.Printf("Testing archive: %s\n\n", filepath.Base(archivePath))

	manager := archive.NewManager()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Run tests
	result, err := manager.Test(ctx, archivePath)
	if err != nil {
		return fmt.Errorf("test failed: %w", err)
	}

	// Display results
	printTestResult(archivePath, result)

	if !result.Passed {
		return fmt.Errorf("archive verification failed")
	}

	return nil
}

func runTestDirectory(dir string) error {
	// Find all archives
	archives, err := findArchives(dir)
	if err != nil {
		return fmt.Errorf("failed to find archives: %w", err)
	}

	if len(archives) == 0 {
		fmt.Printf("No archives found in %s\n", dir)
		return nil
	}

	fmt.Printf("Testing %d archives in %s\n\n", len(archives), dir)

	// Create progress bar
	bar := progressbar.Default(int64(len(archives)))

	// Results storage
	results := make([]*archive.TestResult, len(archives))
	var resultsMu sync.Mutex

	// Run tests concurrently
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(maxConcurrent)

	for i, archivePath := range archives {
		i, archivePath := i, archivePath // Capture loop variables

		g.Go(func() error {
			if err := ctx.Err(); err != nil {
				return err
			}

<<<<<<< HEAD
			// Test archive (per-archive timeout for parity with single mode)
=======
			// Test archive
>>>>>>> origin/main
			manager := archive.NewManager()
			ctxArchive, cancel := context.WithTimeout(ctx, 10*time.Minute)
			defer cancel()
			result, err := manager.Test(ctxArchive, archivePath)
			if err != nil {
				result = &archive.TestResult{
					Passed: false,
					Errors: []string{err.Error()},
				}
			}

			// Store result
			resultsMu.Lock()
			results[i] = result
			bar.Add(1)
			resultsMu.Unlock()

			return nil
		})
	}

	// Wait for all tests to complete
	if err := g.Wait(); err != nil {
		return fmt.Errorf("testing failed: %w", err)
	}

	bar.Finish()
	fmt.Printf("\n")

	// Print summary
	printBatchSummary(archives, results)

	// Check if any failed
	failedCount := 0
	for _, result := range results {
		if !result.Passed {
			failedCount++
		}
	}

	if failedCount > 0 {
		return fmt.Errorf("%d archives failed verification", failedCount)
	}

	return nil
}

func findArchives(dir string) ([]string, error) {
	var archives []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(path), ".7z") {
			archives = append(archives, path)
		}
		return nil
	})
	return archives, err
}

func printTestResult(path string, result *archive.TestResult) {
	if result.Passed {
		fmt.Printf("✅ PASS: Archive integrity verified\n")
		fmt.Printf("  Archive structure: VALID\n")
		if result.ChecksumValid {
			fmt.Printf("  Checksums: ALL MATCH (%d files verified)\n", result.FilesVerified)
		}
		if result.MetadataValid {
			fmt.Printf("  Metadata: CONSISTENT\n")
		}
		fmt.Printf("  Extraction: SUCCESS\n")
		fmt.Printf("  Completeness: ALL ARTIFACTS PRESENT\n")
	} else {
		fmt.Printf("❌ FAIL: Archive verification failed\n")
		if len(result.Errors) > 0 {
			fmt.Printf("  Errors:\n")
			for _, err := range result.Errors {
				fmt.Printf("    - %s\n", err)
			}
		}
	}
	fmt.Printf("\n")
}

func printBatchSummary(archives []string, results []*archive.TestResult) {
	passed := 0
	failed := 0
	totalFiles := 0

	for _, result := range results {
		if result.Passed {
			passed++
		} else {
			failed++
		}
		totalFiles += result.FilesVerified
	}

	fmt.Printf("Batch Summary:\n")
	fmt.Printf("- Total archives tested: %d\n", len(archives))
	fmt.Printf("- Passed: %d (%.1f%%)\n", passed, float64(passed)/float64(len(archives))*100)
	if failed > 0 {
		fmt.Printf("- Failed: %d (%.1f%%)\n", failed, float64(failed)/float64(len(archives))*100)

		// List failed archives
		fmt.Printf("\nFailed archives:\n")
		for i, result := range results {
			if !result.Passed {
				fmt.Printf("  ❌ %s\n", filepath.Base(archives[i]))
			}
		}
	}
	fmt.Printf("- Total files verified: %d\n", totalFiles)

	if passed == len(archives) {
		fmt.Printf("\n✅ All archives passed verification!\n")
	}
}
