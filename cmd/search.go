package cmd

import (
	"fmt"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/search"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func SearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search archives by content with full-text and field-specific capabilities",
		Long: `Perform powerful full-text search across archive metadata with field-specific 
filtering and regex support.

The search engine provides sub-500ms performance on large archive collections
and supports both simple text queries and advanced pattern matching.`,
		Example: `  # Full-text search across all fields
  7zarch-go search "project backup 2024"
  
  # Field-specific search
  7zarch-go search --field=name "important"
  7zarch-go search --field=path "/Users/*/Documents"
  
  # Regex pattern matching
  7zarch-go search --field=name --regex ".*\.sql$"
  
  # Case-sensitive search with result limit
  7zarch-go search "Project" --case-sensitive --limit 10
  
  # Combined with filters (coming in Phase 4)
  7zarch-go search "backup" --profile=documents --managed`,
	}

	cmd.AddCommand(searchQueryCmd())
	cmd.AddCommand(searchReindexCmd())

	return cmd
}

func searchQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query <search-terms>",
		Short: "Search archives with specified query",
		Long: `Execute a search query against the archive metadata.

Supports full-text search across all fields or field-specific search
with optional regex pattern matching.`,
		Example: `  # Search for archives containing "backup" and "2024"
  7zarch-go search query "backup 2024"
  
  # Search only in archive names
  7zarch-go search query --field=name "project"
  
  # Use regex pattern
  7zarch-go search query --field=path --regex "/Users/.*/Documents/.*"
  
  # Case-sensitive search
  7zarch-go search query "Project" --case-sensitive`,
		Args: cobra.MinimumNArgs(1),
		RunE: runSearchQuery,
	}

	// Search options
	cmd.Flags().String("field", "", "Search specific field (name|path|metadata)")
	cmd.Flags().Bool("regex", false, "Use regex pattern matching")
	cmd.Flags().Bool("case-sensitive", false, "Case-sensitive search")
	cmd.Flags().Int("limit", 0, "Maximum number of results (0 = no limit)")

	// Output options
	cmd.Flags().String("output", "", "Output format: table|json|csv|yaml")
	cmd.Flags().Bool("details", false, "Show detailed information")

	// Display mode options
	cmd.Flags().Bool("table", false, "Use table display mode")
	cmd.Flags().Bool("compact", false, "Use compact display mode")
	cmd.Flags().Bool("card", false, "Use card display mode")
	cmd.Flags().Bool("tree", false, "Use tree display mode")
	cmd.Flags().Bool("dashboard", false, "Use dashboard display mode")

	return cmd
}

func searchReindexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reindex",
		Short: "Rebuild the search index",
		Long: `Rebuild the search index from current archive data.

This is useful when archives have been modified outside of 7zarch-go
or when search performance degrades due to index fragmentation.`,
		Example: `  # Rebuild search index
  7zarch-go search reindex`,
		RunE: runSearchReindex,
	}

	return cmd
}

func runSearchQuery(cmd *cobra.Command, args []string) error {
	query := args[0]
	if len(args) > 1 {
		// Join multiple arguments into single query
		for i := 1; i < len(args); i++ {
			query += " " + args[i]
		}
	}

	// Parse search options
	opts := search.SearchOptions{
		Field:         getString(cmd, "field"),
		UseRegex:      getBool(cmd, "regex"),
		CaseSensitive: getBool(cmd, "case-sensitive"),
		MaxResults:    getInt(cmd, "limit"),
	}

	// Initialize search engine
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage manager: %w", err)
	}
	defer storageManager.Close()

	searchEngine := search.NewSearchEngine(storageManager.Registry())

	// Ensure search index is available
	if err := searchEngine.EnsureSearchTable(); err != nil {
		return fmt.Errorf("failed to initialize search index: %w", err)
	}

	// Perform search with timing
	startTime := time.Now()
	results, err := searchEngine.SearchWithOptions(query, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}
	searchTime := time.Since(startTime)

	// Check for output format first
	outputFormat := getString(cmd, "output")
	if outputFormat != "" {
		switch outputFormat {
		case "json":
			return outputJSON(results)
		case "csv":
			return outputCSV(results)
		case "yaml":
			return outputYAML(results)
		default:
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}
	}

	// Display results
	if len(results) == 0 {
		fmt.Printf("🔍 Search '%s' - No archives found\n", query)
		if opts.Field != "" {
			fmt.Printf("   Field: %s\n", opts.Field)
		}
		if opts.UseRegex {
			fmt.Printf("   Pattern matching: regex\n")
		}
		fmt.Printf("   Search time: %v\n", searchTime)
		return nil
	}

	fmt.Printf("🔍 Search '%s' - %d archives found", query, len(results))
	if opts.Field != "" {
		fmt.Printf(" (field: %s)", opts.Field)
	}
	if opts.UseRegex {
		fmt.Printf(" (regex)")
	}
	if opts.MaxResults > 0 && len(results) >= opts.MaxResults {
		fmt.Printf(" (limited to %d)", opts.MaxResults)
	}
	fmt.Printf("\n")
	fmt.Printf("Search time: %v", searchTime)
	
	// Performance warning if over target
	if searchTime > 500*time.Millisecond {
		fmt.Printf(" ⚠️  (target: <500ms)")
	}
	fmt.Printf("\n\n")

	// Display results using existing archive display logic
	return printGroupedArchives(results, getBool(cmd, "details"))
}

func runSearchReindex(cmd *cobra.Command, args []string) error {
	// Initialize search engine
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	storageManager, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage manager: %w", err)
	}
	defer storageManager.Close()

	searchEngine := search.NewSearchEngine(storageManager.Registry())

	// Ensure search table exists
	if err := searchEngine.EnsureSearchTable(); err != nil {
		return fmt.Errorf("failed to initialize search table: %w", err)
	}

	fmt.Printf("🔄 Rebuilding search index...\n")
	
	startTime := time.Now()
	if err := searchEngine.Reindex(); err != nil {
		return fmt.Errorf("failed to rebuild search index: %w", err)
	}
	indexTime := time.Since(startTime)

	fmt.Printf("✅ Search index rebuilt in %v\n", indexTime)
	
	return nil
}

// Helper function for int flags
func getInt(cmd *cobra.Command, name string) int {
	v, _ := cmd.Flags().GetInt(name)
	return v
}