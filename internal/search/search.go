package search

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
)

const (
	migrationSearchID   = "0005_search_index"
	migrationSearchName = "Add search_index table for full-text search support"
)

// SearchEngine provides full-text search across archive metadata
type SearchEngine struct {
	registry *storage.Registry
	index    *InvertedIndex
	mu       sync.RWMutex
}

// InvertedIndex stores term -> archive UID mappings for fast text search
type InvertedIndex struct {
	// terms maps search terms to archive UIDs that contain them
	terms map[string][]string
	// fieldTerms maps field names to their term mappings for field-specific search
	fieldTerms map[string]map[string][]string
	// cache provides LRU caching for frequent search terms
	cache *LRUCache
	// lastUpdate tracks when the index was last rebuilt
	lastUpdate time.Time
}

// LRUCache provides memory-efficient caching for search results
type LRUCache struct {
	capacity int
	items    map[string]*cacheItem
	order    *cacheList
	mu       sync.RWMutex
}

type cacheItem struct {
	key     string
	value   []*storage.Archive
	expires time.Time
	prev    *cacheItem
	next    *cacheItem
}

type cacheList struct {
	head *cacheItem
	tail *cacheItem
}

// SearchOptions configures search behavior
type SearchOptions struct {
	Field        string // Specific field to search (name, path, metadata)
	UseRegex     bool   // Enable regex pattern matching
	CaseSensitive bool   // Case-sensitive search
	MaxResults   int    // Limit number of results (0 = no limit)
}

// NewSearchEngine creates a new search engine instance
func NewSearchEngine(registry *storage.Registry) *SearchEngine {
	return &SearchEngine{
		registry: registry,
		index: &InvertedIndex{
			terms:      make(map[string][]string),
			fieldTerms: make(map[string]map[string][]string),
			cache:      NewLRUCache(1000), // Cache up to 1000 search results
		},
	}
}

// NewLRUCache creates a new LRU cache with specified capacity
func NewLRUCache(capacity int) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		items:    make(map[string]*cacheItem),
		order: &cacheList{
			head: &cacheItem{},
			tail: &cacheItem{},
		},
	}
	// Initialize the doubly-linked list
	cache.order.head.next = cache.order.tail
	cache.order.tail.prev = cache.order.head
	return cache
}

// Search performs full-text search across all archive metadata
func (se *SearchEngine) Search(query string) ([]*storage.Archive, error) {
	return se.SearchWithOptions(query, SearchOptions{})
}

// SearchField performs field-specific search with optional regex support
func (se *SearchEngine) SearchField(field, query string) ([]*storage.Archive, error) {
	return se.SearchWithOptions(query, SearchOptions{
		Field: field,
	})
}

// SearchRegex performs regex-based search on specified field
func (se *SearchEngine) SearchRegex(field, pattern string) ([]*storage.Archive, error) {
	return se.SearchWithOptions(pattern, SearchOptions{
		Field:    field,
		UseRegex: true,
	})
}

// SearchWithOptions performs search with full configuration control
func (se *SearchEngine) SearchWithOptions(query string, opts SearchOptions) ([]*storage.Archive, error) {
	startTime := time.Now()

	// Check cache first for performance
	cacheKey := fmt.Sprintf("%s:%s:%v:%v", query, opts.Field, opts.UseRegex, opts.CaseSensitive)
	if cached := se.index.cache.Get(cacheKey); cached != nil {
		return cached, nil
	}

	// Ensure index is up to date (check without lock first to avoid deadlock)
	if err := se.ensureIndexCurrent(); err != nil {
		return nil, fmt.Errorf("failed to update search index: %w", err)
	}
	
	// Now acquire read lock for search operations
	se.mu.RLock()
	defer se.mu.RUnlock()

	var results []*storage.Archive
	var err error

	if opts.UseRegex {
		results, err = se.searchRegex(query, opts)
	} else if opts.Field != "" {
		results, err = se.searchField(opts.Field, query, opts)
	} else {
		results, err = se.searchFullText(query, opts)
	}

	if err != nil {
		return nil, err
	}

	// Apply result limit if specified
	if opts.MaxResults > 0 && len(results) > opts.MaxResults {
		results = results[:opts.MaxResults]
	}

	// Cache results for future queries (expires in 5 minutes)
	se.index.cache.Set(cacheKey, results, 5*time.Minute)

	// Log performance for optimization (in production, use proper logging)
	queryTime := time.Since(startTime)
	if queryTime > 500*time.Millisecond {
		// This should trigger monitoring alerts in production
		fmt.Printf("PERF WARNING: Search query took %v (target: <500ms)\n", queryTime)
	}

	return results, nil
}

// searchFullText performs cross-field text search
func (se *SearchEngine) searchFullText(query string, opts SearchOptions) ([]*storage.Archive, error) {
	query = se.normalizeQuery(query, opts)
	terms := strings.Fields(query)
	
	if len(terms) == 0 {
		return nil, fmt.Errorf("empty search query")
	}

	// Find archives that contain all search terms (AND logic)
	var candidateUIDs []string
	for i, term := range terms {
		termUIDs := se.index.terms[term]
		if len(termUIDs) == 0 {
			// No archives contain this term, so no results possible
			return []*storage.Archive{}, nil
		}
		
		if i == 0 {
			candidateUIDs = termUIDs
		} else {
			// Intersect with previous results
			candidateUIDs = intersectStringSlices(candidateUIDs, termUIDs)
		}
		
		if len(candidateUIDs) == 0 {
			break // Early termination if no intersection
		}
	}

	return se.resolveArchiveUIDs(candidateUIDs)
}

// searchField performs field-specific search
func (se *SearchEngine) searchField(field, query string, opts SearchOptions) ([]*storage.Archive, error) {
	query = se.normalizeQuery(query, opts)
	
	fieldIndex, exists := se.index.fieldTerms[field]
	if !exists {
		return []*storage.Archive{}, nil
	}

	terms := strings.Fields(query)
	var candidateUIDs []string
	
	for i, term := range terms {
		termUIDs := fieldIndex[term]
		if len(termUIDs) == 0 {
			return []*storage.Archive{}, nil
		}
		
		if i == 0 {
			candidateUIDs = termUIDs
		} else {
			candidateUIDs = intersectStringSlices(candidateUIDs, termUIDs)
		}
		
		if len(candidateUIDs) == 0 {
			break
		}
	}

	return se.resolveArchiveUIDs(candidateUIDs)
}

// searchRegex performs regex-based search
func (se *SearchEngine) searchRegex(pattern string, opts SearchOptions) ([]*storage.Archive, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Get all archives and filter by regex
	archives, err := se.registry.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list archives: %w", err)
	}

	var results []*storage.Archive
	for _, archive := range archives {
		var searchText string
		
		switch opts.Field {
		case "name":
			searchText = archive.Name
		case "path":
			searchText = archive.Path
		case "metadata":
			searchText = archive.Metadata
		default:
			// Search all fields
			searchText = fmt.Sprintf("%s %s %s", archive.Name, archive.Path, archive.Metadata)
		}
		
		if !opts.CaseSensitive {
			searchText = strings.ToLower(searchText)
		}
		
		if regex.MatchString(searchText) {
			results = append(results, archive)
		}
	}

	return results, nil
}

// Reindex rebuilds the search index from current archive data
func (se *SearchEngine) Reindex() error {
	se.mu.Lock()
	defer se.mu.Unlock()

	startTime := time.Now()

	// Get all archives
	archives, err := se.registry.List()
	if err != nil {
		return fmt.Errorf("failed to list archives for indexing: %w", err)
	}

	// Clear existing index
	se.index.terms = make(map[string][]string)
	se.index.fieldTerms = make(map[string]map[string][]string)
	
	// Initialize field indices
	se.index.fieldTerms["name"] = make(map[string][]string)
	se.index.fieldTerms["path"] = make(map[string][]string)
	se.index.fieldTerms["profile"] = make(map[string][]string)
	se.index.fieldTerms["metadata"] = make(map[string][]string)

	// Build index from archives
	for _, archive := range archives {
		se.indexArchive(archive)
	}

	se.index.lastUpdate = time.Now()
	
	// Clear cache since index changed
	se.index.cache = NewLRUCache(1000)

	indexTime := time.Since(startTime)
	fmt.Printf("Search index rebuilt: %d archives indexed in %v\n", len(archives), indexTime)

	return nil
}

// indexArchive adds a single archive to the search index
func (se *SearchEngine) indexArchive(archive *storage.Archive) {
	// Index individual fields
	se.indexFieldTerms("name", archive.Name, archive.UID)
	se.indexFieldTerms("path", archive.Path, archive.UID)
	se.indexFieldTerms("profile", archive.Profile, archive.UID)
	se.indexFieldTerms("metadata", archive.Metadata, archive.UID)

	// Index all text for cross-field search
	allText := fmt.Sprintf("%s %s %s %s", archive.Name, archive.Path, archive.Profile, archive.Metadata)
	se.indexTerms(allText, archive.UID)
}

// indexFieldTerms adds terms from a specific field to the field-specific index
func (se *SearchEngine) indexFieldTerms(field, text, uid string) {
	terms := se.extractTerms(text)
	fieldIndex := se.index.fieldTerms[field]
	
	for _, term := range terms {
		if !stringSliceContains(fieldIndex[term], uid) {
			fieldIndex[term] = append(fieldIndex[term], uid)
		}
	}
}

// indexTerms adds terms to the main cross-field index
func (se *SearchEngine) indexTerms(text, uid string) {
	terms := se.extractTerms(text)
	
	for _, term := range terms {
		if !stringSliceContains(se.index.terms[term], uid) {
			se.index.terms[term] = append(se.index.terms[term], uid)
		}
	}
}

// extractTerms breaks text into searchable terms
func (se *SearchEngine) extractTerms(text string) []string {
	// Normalize to lowercase
	text = strings.ToLower(text)
	
	// Split on whitespace and common separators
	terms := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '/' || r == '\\' || r == '.' || r == '-' || r == '_'
	})
	
	// Filter out empty terms and very short terms
	var filtered []string
	for _, term := range terms {
		if len(term) >= 2 { // Minimum term length for performance
			filtered = append(filtered, term)
		}
	}
	
	return filtered
}

// normalizeQuery prepares query text for searching
func (se *SearchEngine) normalizeQuery(query string, opts SearchOptions) string {
	if !opts.CaseSensitive {
		query = strings.ToLower(query)
	}
	return strings.TrimSpace(query)
}

// ensureIndexCurrent checks if index needs rebuilding
func (se *SearchEngine) ensureIndexCurrent() error {
	// If index has never been built, build it now
	if se.index.lastUpdate.IsZero() {
		return se.Reindex()
	}
	// Simple time-based check - in production would check archive modification times
	if time.Since(se.index.lastUpdate) > 10*time.Minute {
		return se.Reindex()
	}
	return nil
}

// resolveArchiveUIDs converts UIDs to Archive objects
func (se *SearchEngine) resolveArchiveUIDs(uids []string) ([]*storage.Archive, error) {
	var archives []*storage.Archive
	
	for _, uid := range uids {
		archive, err := se.registry.GetByUID(uid)
		if err != nil {
			// Archive may have been deleted, skip it
			continue
		}
		archives = append(archives, archive)
	}
	
	return archives, nil
}

// LRU Cache implementation

// Get retrieves an item from cache
func (c *LRUCache) Get(key string) []*storage.Archive {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.expires) {
		if exists {
			c.removeItem(item)
		}
		return nil
	}

	// Move to front (most recently used)
	c.moveToFront(item)
	return item.value
}

// Set stores an item in cache with expiration
func (c *LRUCache) Set(key string, value []*storage.Archive, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove existing item if present
	if existing, exists := c.items[key]; exists {
		c.removeItem(existing)
	}

	// Create new item
	item := &cacheItem{
		key:     key,
		value:   value,
		expires: time.Now().Add(ttl),
	}

	// Add to front
	c.addToFront(item)
	c.items[key] = item

	// Evict oldest items if over capacity
	for len(c.items) > c.capacity {
		oldest := c.order.tail.prev
		if oldest == c.order.head {
			// Sanity check to prevent infinite loop
			break
		}
		c.removeItem(oldest)
	}
}

// moveToFront moves item to front of LRU list
func (c *LRUCache) moveToFront(item *cacheItem) {
	c.removeFromList(item)
	c.addToFront(item)
}

// addToFront adds item to front of LRU list
func (c *LRUCache) addToFront(item *cacheItem) {
	item.prev = c.order.head
	item.next = c.order.head.next
	c.order.head.next.prev = item
	c.order.head.next = item
}

// removeItem removes item from cache
func (c *LRUCache) removeItem(item *cacheItem) {
	delete(c.items, item.key)
	c.removeFromList(item)
}

// removeFromList removes item from LRU list
func (c *LRUCache) removeFromList(item *cacheItem) {
	item.prev.next = item.next
	item.next.prev = item.prev
}

// Utility functions

// intersectStringSlices returns elements present in both slices
func intersectStringSlices(a, b []string) []string {
	set := make(map[string]bool)
	for _, item := range a {
		set[item] = true
	}

	var result []string
	for _, item := range b {
		if set[item] {
			result = append(result, item)
		}
	}

	return result
}

// stringSliceContains checks if slice contains string
func stringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Migration support for search index table

// EnsureSearchTable creates search index table if it doesn't exist
func (se *SearchEngine) EnsureSearchTable() error {
	_, err := se.registry.DB().Exec(`
		CREATE TABLE IF NOT EXISTS search_index (
			term TEXT,
			archive_uid TEXT,
			field TEXT,
			PRIMARY KEY (term, archive_uid, field)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create search_index table: %w", err)
	}

	// Create index for performance
	_, err = se.registry.DB().Exec(`
		CREATE INDEX IF NOT EXISTS idx_search_term ON search_index(term)
	`)
	if err != nil {
		return fmt.Errorf("failed to create search index: %w", err)
	}

	return nil
}

// IsSearchMigrationApplied checks if search migration has been applied
func IsSearchMigrationApplied(db *sql.DB) (bool, error) {
	row := db.QueryRow(`SELECT 1 FROM schema_migrations WHERE id = ?`, migrationSearchID)
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

// ApplySearchMigration applies the search system migration
func ApplySearchMigration(db *sql.DB) error {
	// Create search_index table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS search_index (
			term TEXT,
			archive_uid TEXT,
			field TEXT,
			PRIMARY KEY (term, archive_uid, field)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create search_index table: %w", err)
	}

	// Create performance index
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_search_term ON search_index(term)
	`)
	if err != nil {
		return fmt.Errorf("failed to create search index: %w", err)
	}

	// Mark migration as applied
	_, err = db.Exec(`
		INSERT OR REPLACE INTO schema_migrations (id, name, applied_at) 
		VALUES (?, ?, ?)
	`, migrationSearchID, migrationSearchName, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to mark search migration as applied: %w", err)
	}

	return nil
}