package query

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/adamstac/7zarch-go/internal/search"
	"github.com/adamstac/7zarch-go/internal/storage"
)

const (
	migrationQueryID   = "0004_query_system"
	migrationQueryName = "Add queries table for saved query support"
)

// Query represents a saved filter configuration
type Query struct {
	Name     string            `json:"name"`
	Filters  map[string]string `json:"filters"`
	Created  time.Time         `json:"created"`
	LastUsed *time.Time        `json:"last_used,omitempty"`
	UseCount int               `json:"use_count"`
}

// QueryManager handles saved query operations
type QueryManager struct {
	db           *sql.DB
	resolver     *storage.Resolver
	searchEngine *search.SearchEngine
}

// NewQueryManager creates a new query manager instance
func NewQueryManager(db *sql.DB, resolver *storage.Resolver) *QueryManager {
	searchEngine := search.NewSearchEngine(resolver.Registry())
	return &QueryManager{
		db:           db,
		resolver:     resolver,
		searchEngine: searchEngine,
	}
}

// Save stores a new saved query
func (qm *QueryManager) Save(name string, filters map[string]string) error {
	if name == "" {
		return fmt.Errorf("query name cannot be empty")
	}

	// Ensure query table exists
	if err := qm.ensureQueryTable(); err != nil {
		return fmt.Errorf("failed to ensure query table: %w", err)
	}

	// Serialize filters to JSON
	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		return fmt.Errorf("failed to serialize filters: %w", err)
	}

	// Insert or replace the query
	_, err = qm.db.Exec(`
		INSERT OR REPLACE INTO queries (name, filters, created, last_used, use_count)
		VALUES (?, ?, ?, NULL, 0)
	`, name, string(filtersJSON), time.Now().Unix())

	if err != nil {
		return fmt.Errorf("failed to save query: %w", err)
	}

	return nil
}

// List returns all saved queries
func (qm *QueryManager) List() ([]*Query, error) {
	if err := qm.ensureQueryTable(); err != nil {
		return nil, fmt.Errorf("failed to ensure query table: %w", err)
	}

	rows, err := qm.db.Query(`
		SELECT name, filters, created, last_used, use_count 
		FROM queries 
		ORDER BY last_used DESC, created DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query saved queries: %w", err)
	}
	defer rows.Close()

	var queries []*Query
	for rows.Next() {
		var query Query
		var filtersJSON string
		var createdUnix int64
		var lastUsedUnix sql.NullInt64

		if err := rows.Scan(&query.Name, &filtersJSON, &createdUnix, &lastUsedUnix, &query.UseCount); err != nil {
			return nil, fmt.Errorf("failed to scan query row: %w", err)
		}

		// Convert Unix timestamps to time.Time
		query.Created = time.Unix(createdUnix, 0)
		if lastUsedUnix.Valid {
			lastUsedTime := time.Unix(lastUsedUnix.Int64, 0)
			query.LastUsed = &lastUsedTime
		}

		// Parse filters JSON
		if err := json.Unmarshal([]byte(filtersJSON), &query.Filters); err != nil {
			return nil, fmt.Errorf("failed to parse filters for query %s: %w", query.Name, err)
		}

		queries = append(queries, &query)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating query rows: %w", err)
	}

	return queries, nil
}

// Run executes a saved query and returns matching archives
func (qm *QueryManager) Run(name string) ([]*storage.Archive, error) {
	if err := qm.ensureQueryTable(); err != nil {
		return nil, fmt.Errorf("failed to ensure query table: %w", err)
	}

	// Get the query
	var filtersJSON string
	err := qm.db.QueryRow(`SELECT filters FROM queries WHERE name = ?`, name).Scan(&filtersJSON)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("query not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// Parse filters
	var filters map[string]string
	if err := json.Unmarshal([]byte(filtersJSON), &filters); err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	// Update usage statistics
	_, err = qm.db.Exec(`
		UPDATE queries 
		SET last_used = ?, use_count = use_count + 1 
		WHERE name = ?
	`, time.Now().Unix(), name)
	if err != nil {
		// Don't fail the query execution for statistics update failure
		// Log this in production
	}

	// Execute the filters using the resolver
	return qm.executeFilters(filters)
}

// Delete removes a saved query
func (qm *QueryManager) Delete(name string) error {
	if err := qm.ensureQueryTable(); err != nil {
		return fmt.Errorf("failed to ensure query table: %w", err)
	}

	result, err := qm.db.Exec(`DELETE FROM queries WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("query not found: %s", name)
	}

	return nil
}

// Get retrieves a specific saved query
func (qm *QueryManager) Get(name string) (*Query, error) {
	if err := qm.ensureQueryTable(); err != nil {
		return nil, fmt.Errorf("failed to ensure query table: %w", err)
	}

	var query Query
	var filtersJSON string
	var createdUnix int64
	var lastUsedUnix sql.NullInt64

	err := qm.db.QueryRow(`
		SELECT name, filters, created, last_used, use_count 
		FROM queries 
		WHERE name = ?
	`, name).Scan(&query.Name, &filtersJSON, &createdUnix, &lastUsedUnix, &query.UseCount)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("query not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	// Convert Unix timestamps to time.Time
	query.Created = time.Unix(createdUnix, 0)
	if lastUsedUnix.Valid {
		lastUsedTime := time.Unix(lastUsedUnix.Int64, 0)
		query.LastUsed = &lastUsedTime
	}

	// Parse filters JSON
	if err := json.Unmarshal([]byte(filtersJSON), &query.Filters); err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	return &query, nil
}

// executeFilters converts filter map to archive list using the resolver and search engine
func (qm *QueryManager) executeFilters(filters map[string]string) ([]*storage.Archive, error) {
	var archives []*storage.Archive
	var err error

	// Check if this query includes search terms
	if searchTerm, hasSearch := filters["search"]; hasSearch {
		// Use search engine for the base archive set
		searchOpts := search.SearchOptions{}
		
		// Apply search options from filters
		if field, ok := filters["search-field"]; ok {
			searchOpts.Field = field
		}
		if _, ok := filters["search-regex"]; ok {
			searchOpts.UseRegex = true
		}
		if _, ok := filters["search-case-sensitive"]; ok {
			searchOpts.CaseSensitive = true
		}

		// Ensure search table exists
		if err := qm.searchEngine.EnsureSearchTable(); err != nil {
			return nil, fmt.Errorf("failed to initialize search: %w", err)
		}

		archives, err = qm.searchEngine.SearchWithOptions(searchTerm, searchOpts)
		if err != nil {
			return nil, fmt.Errorf("search failed: %w", err)
		}
	} else {
		// No search terms, get all archives from registry
		registry := qm.resolver.Registry()
		if registry == nil {
			return nil, fmt.Errorf("registry not available")
		}

		archives, err = registry.List()
		if err != nil {
			return nil, fmt.Errorf("failed to list archives: %w", err)
		}
	}

	// Apply additional filters (non-search filters)
	filtered := make([]*storage.Archive, 0)
	for _, archive := range archives {
		if qm.matchesFilters(archive, filters) {
			filtered = append(filtered, archive)
		}
	}

	return filtered, nil
}

// matchesFilters checks if an archive matches the saved filter criteria
func (qm *QueryManager) matchesFilters(archive *storage.Archive, filters map[string]string) bool {
	for key, value := range filters {
		switch key {
		case "status":
			if archive.Status != value {
				return false
			}
		case "profile":
			if archive.Profile != value {
				return false
			}
		case "managed":
			if value == "true" && !archive.Managed {
				return false
			}
			if value == "false" && archive.Managed {
				return false
			}
		case "uploaded":
			if value == "true" && !archive.Uploaded {
				return false
			}
			if value == "false" && archive.Uploaded {
				return false
			}
		case "missing":
			if value == "true" && archive.Status != "missing" {
				return false
			}
		case "deleted":
			if value == "true" && archive.Status != "deleted" {
				return false
			}
		case "larger-than":
			// Parse the size value
			if threshold, err := strconv.ParseInt(value, 10, 64); err == nil {
				if archive.Size <= threshold {
					return false
				}
			}
		case "not-uploaded":
			if value == "true" && archive.Uploaded {
				return false
			}
		// Skip search-specific filters - they're handled separately
		case "search", "search-field", "search-regex", "search-case-sensitive":
			continue
		// Note: More complex filters like pattern, older-than
		// can be implemented as needed
		}
	}
	return true
}

// ensureQueryTable creates the queries table if it doesn't exist
func (qm *QueryManager) ensureQueryTable() error {
	_, err := qm.db.Exec(`
		CREATE TABLE IF NOT EXISTS queries (
			name TEXT PRIMARY KEY,
			filters TEXT NOT NULL,
			created INTEGER NOT NULL,
			last_used INTEGER,
			use_count INTEGER DEFAULT 0
		)
	`)
	return err
}

// IsQueryMigrationApplied checks if the query system migration has been applied
func IsQueryMigrationApplied(db *sql.DB) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM schema_migrations WHERE id = ?`, migrationQueryID)
	var one int
	switch err := row.Scan(&one); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// ApplyQueryMigration applies the query system migration
func ApplyQueryMigration(db *sql.DB) error {
	// Create queries table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS queries (
			name TEXT PRIMARY KEY,
			filters TEXT NOT NULL,
			created INTEGER NOT NULL,
			last_used INTEGER,
			use_count INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create queries table: %w", err)
	}

	// Mark migration as applied
	_, err = db.Exec(`
		INSERT OR REPLACE INTO schema_migrations (id, name, applied_at) 
		VALUES (?, ?, ?)
	`, migrationQueryID, migrationQueryName, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to mark query migration as applied: %w", err)
	}

	return nil
}