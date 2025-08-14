package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adamstac/7zarch-go/internal/batch"
	"github.com/adamstac/7zarch-go/internal/query"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func BatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch <operation> [flags]",
		Short: "Perform operations on multiple archives",
		Long: `Perform operations on multiple archives using saved queries, stdin, or filter combinations.

Operations:
  move   - Move archives to a new location
  delete - Delete archives (with trash integration)

Selection methods:
  --query=<name>     Use saved query to select archives
  --stdin            Read archive UIDs from stdin (one per line)
  [filters...]       Use filter flags to select archives

Examples:
  # Move archives using saved query
  7zarch-go batch move --query=old-files --to=/archive/old/

  # Delete archives from stdin
  7zarch-go list --older-than=1y --output=json | jq -r '.[].uid' | 7zarch-go batch delete --stdin --confirm

  # Move with filters
  7zarch-go batch move --profile=documents --larger-than=100MB --to=/backup/docs/

  # Batch delete with confirmation
  7zarch-go batch delete --query=temp-files --confirm`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: batchOperationCompletion,
		RunE:              runBatch,
	}

	// Selection flags
	cmd.Flags().String("query", "", "Use saved query to select archives")
	cmd.Flags().Bool("stdin", false, "Read archive UIDs from stdin")
	cmd.Flags().Bool("all", false, "Process all archives (REQUIRED if no query/stdin specified)")

	// Operation-specific flags
	cmd.Flags().String("to", "", "Destination path for move operation")
	cmd.Flags().Bool("confirm", false, "Confirm destructive operations")
	cmd.Flags().Bool("dry-run", false, "Show what would be done without executing")

	// Progress and output flags
	cmd.Flags().Bool("progress", true, "Show progress during batch operations")
	cmd.Flags().Int("concurrent", 4, "Number of concurrent operations")
	cmd.Flags().String("output", "table", "Output format: table, json, csv, yaml")

	// Include all list filter flags for direct filtering
	addListFilterFlags(cmd)

	return cmd
}

func runBatch(cmd *cobra.Command, args []string) error {
	operation := args[0]

	// Validate operation
	validOps := []string{"move", "delete"}
	isValid := false
	for _, valid := range validOps {
		if operation == valid {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid operation '%s'. Valid operations: %s", operation, strings.Join(validOps, ", "))
	}

	// Get selection method
	queryName, _ := cmd.Flags().GetString("query")
	useStdin, _ := cmd.Flags().GetBool("stdin")
	useAll, _ := cmd.Flags().GetBool("all")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Validate selection method - exactly one must be specified
	selectionCount := 0
	if queryName != "" {
		selectionCount++
	}
	if useStdin {
		selectionCount++
	}
	if useAll {
		selectionCount++
	}

	if selectionCount == 0 {
		return fmt.Errorf("must specify exactly one selection method: --query, --stdin, or --all")
	}
	if selectionCount > 1 {
		return fmt.Errorf("cannot combine selection methods: use only one of --query, --stdin, or --all")
	}

	// Initialize storage
	manager, err := storage.NewManager("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	var archives []*storage.Archive

	if queryName != "" {
		// Use saved query
		resolver := storage.NewResolver(manager.Registry())
		queryManager := query.NewQueryManager(manager.Registry().DB(), resolver)
		archives, err = queryManager.Run(queryName)
		if err != nil {
			return fmt.Errorf("failed to run query '%s': %w", queryName, err)
		}
	} else if useStdin {
		// Read UIDs from stdin
		archives, err = readArchivesFromStdin(manager.Registry())
		if err != nil {
			return fmt.Errorf("failed to read archives from stdin: %w", err)
		}
	} else if useAll {
		// Use all archives (explicitly requested with --all flag)
		archives, err = manager.Registry().List()
		if err != nil {
			return fmt.Errorf("failed to list archives: %w", err)
		}
	} else {
		// This should never happen due to validation above
		return fmt.Errorf("internal error: no valid selection method")
	}

	if len(archives) == 0 {
		fmt.Println("No archives selected for batch operation")
		return nil
	}

	// Show what will be processed
	fmt.Printf("Selected %d archive(s) for %s operation:\n", len(archives), operation)
	for i, archive := range archives {
		if i < 10 { // Show first 10
			// Safe UID prefix extraction with bounds checking
			uidDisplay := archive.UID
			if len(archive.UID) > 12 {
				uidDisplay = archive.UID[:12]
			}
			fmt.Printf("  - %s (%s)\n", archive.Name, uidDisplay)
		} else if i == 10 {
			fmt.Printf("  ... and %d more\n", len(archives)-10)
			break
		}
	}
	fmt.Println()

	if dryRun {
		fmt.Println("Dry run - no operations performed")
		return nil
	}

	// Validate operation-specific requirements
	switch operation {
	case "move":
		to, _ := cmd.Flags().GetString("to")
		if to == "" {
			return fmt.Errorf("move operation requires --to flag")
		}
		return performMove(cmd, archives, to, manager)
	case "delete":
		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm {
			return fmt.Errorf("delete operation requires --confirm flag for safety")
		}
		return performDelete(cmd, archives, manager)
	}

	return nil
}

func readArchivesFromStdin(registry *storage.Registry) ([]*storage.Archive, error) {
	var uids []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			uids = append(uids, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading stdin: %w", err)
	}

	var archives []*storage.Archive
	for _, uid := range uids {
		archive, err := registry.Get(uid)
		if err != nil {
			return nil, fmt.Errorf("failed to get archive %s: %w", uid, err)
		}
		archives = append(archives, archive)
	}

	return archives, nil
}

func performMove(cmd *cobra.Command, archives []*storage.Archive, destPath string, manager *storage.Manager) error {
	concurrent, _ := cmd.Flags().GetInt("concurrent")
	showProgress, _ := cmd.Flags().GetBool("progress")

	processor := batch.NewProcessor(manager)
	processor.SetConcurrency(concurrent)

	var progressCallback batch.ProgressCallback
	if showProgress {
		progressCallback = func(update batch.ProgressUpdate) {
			percent := float64(update.Completed) / float64(update.Total) * 100
			fmt.Printf("\rProgress: %d/%d (%.1f%%) - %s", 
				update.Completed, update.Total, percent, update.Current)
			if update.Completed == update.Total {
				fmt.Printf(" - Completed in %v\n", update.Elapsed)
			}
		}
	}

	ctx := cmd.Context()
	err := processor.Move(ctx, archives, destPath, progressCallback)
	if err != nil {
		return fmt.Errorf("batch move failed: %w", err)
	}

	if showProgress {
		fmt.Printf("Successfully moved %d archives to %s\n", len(archives), destPath)
	}

	return nil
}

func performDelete(cmd *cobra.Command, archives []*storage.Archive, manager *storage.Manager) error {
	concurrent, _ := cmd.Flags().GetInt("concurrent")
	showProgress, _ := cmd.Flags().GetBool("progress")

	processor := batch.NewProcessor(manager)
	processor.SetConcurrency(concurrent)

	var progressCallback batch.ProgressCallback
	if showProgress {
		progressCallback = func(update batch.ProgressUpdate) {
			percent := float64(update.Completed) / float64(update.Total) * 100
			fmt.Printf("\rProgress: %d/%d (%.1f%%) - %s", 
				update.Completed, update.Total, percent, update.Current)
			if update.Completed == update.Total {
				elapsed := update.Elapsed
				if len(update.Errors) > 0 {
					fmt.Printf(" - Completed with %d errors in %v\n", len(update.Errors), elapsed)
					fmt.Println("\nErrors:")
					for _, err := range update.Errors {
						fmt.Printf("  - %v\n", err)
					}
				} else {
					fmt.Printf(" - Completed in %v\n", elapsed)
				}
			}
		}
	}

	ctx := cmd.Context()
	err := processor.Delete(ctx, archives, progressCallback)
	if err != nil {
		return fmt.Errorf("batch delete failed: %w", err)
	}

	if showProgress && len(archives) > 0 {
		fmt.Printf("Successfully deleted %d archives\n", len(archives))
	}

	return nil
}

func batchOperationCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return []string{"move", "delete"}, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}


func addListFilterFlags(cmd *cobra.Command) {
	cmd.Flags().String("profile", "", "Filter by compression profile")
	cmd.Flags().Bool("managed", false, "Show only managed archives")
	cmd.Flags().String("status", "", "Filter by status (ok, missing, error)")
	cmd.Flags().String("pattern", "", "Filter by name pattern (glob)")
	cmd.Flags().String("larger-than", "", "Filter by minimum size (e.g., 100MB, 1GB)")
	cmd.Flags().String("smaller-than", "", "Filter by maximum size")
	cmd.Flags().String("older-than", "", "Filter by age (e.g., 30d, 1y)")
	cmd.Flags().String("newer-than", "", "Filter by recency")
}