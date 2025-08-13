package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Registry manages the SQLite database for archive metadata
type Registry struct {
	db     *sql.DB
	dbPath string
}

// NewRegistry creates a new registry instance
func NewRegistry(dbPath string) (*Registry, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	// #nosec G301: create registry directory with restricted permissions
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create registry directory: %w", err)
	}

	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set file permissions
	if err := os.Chmod(dbPath, 0600); err != nil && !os.IsNotExist(err) {
		db.Close()
		return nil, fmt.Errorf("failed to set database permissions: %w", err)
	}

	r := &Registry{db: db, dbPath: dbPath}

	// Initialize the schema
	if err := r.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return r, nil
}

// initSchema creates the database tables if they don't exist and applies migrations
func (r *Registry) initSchema() error {
	// Create full modern schema for new installations
	query := `
	CREATE TABLE IF NOT EXISTS archives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uid TEXT UNIQUE,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		size INTEGER NOT NULL,
		created TIMESTAMP NOT NULL,
		checksum TEXT,
		profile TEXT,
		managed BOOLEAN DEFAULT FALSE,
		status TEXT NOT NULL DEFAULT 'present',
		last_seen TIMESTAMP,
		deleted_at TIMESTAMP,
		original_path TEXT,
		uploaded BOOLEAN DEFAULT FALSE,
		destination TEXT,
		uploaded_at TIMESTAMP,
		metadata TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_archives_created ON archives(created);
	CREATE INDEX IF NOT EXISTS idx_archives_uploaded ON archives(uploaded);
	CREATE INDEX IF NOT EXISTS idx_archives_destination ON archives(destination);
	CREATE INDEX IF NOT EXISTS idx_archives_checksum ON archives(checksum);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_archives_uid ON archives(uid);
	`

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	// Apply proper migration system
	if err := r.RecordCurrentSchemaAsApplied(); err != nil {
		return fmt.Errorf("failed to record schema state: %w", err)
	}

	if err := r.ApplyPendingMigrations(); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Add inserts a new archive into the registry
func (r *Registry) Add(archive *Archive) error {
	query := `
	INSERT INTO archives (uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		archive.UID,
		archive.Name,
		archive.Path,
		archive.Size,
		archive.Created,
		archive.Checksum,
		archive.Profile,
		archive.Managed,
		archive.Status,
		archive.LastSeen,
		archive.DeletedAt,
		archive.OriginalPath,
		archive.Uploaded,
		archive.Destination,
		archive.UploadedAt,
		archive.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to add archive: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	archive.ID = id
	return nil
}

// Get retrieves an archive by name
func (r *Registry) Get(name string) (*Archive, error) {
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE name = ?
	`

	archive := &Archive{}
	err := r.db.QueryRow(query, name).Scan(
		&archive.ID,
		&archive.UID,
		&archive.Name,
		&archive.Path,
		&archive.Size,
		&archive.Created,
		&archive.Checksum,
		&archive.Profile,
		&archive.Managed,
		&archive.Status,
		&archive.LastSeen,
		&archive.DeletedAt,
		&archive.OriginalPath,
		&archive.Uploaded,
		&archive.Destination,
		&archive.UploadedAt,
		&archive.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("archive not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get archive: %w", err)
	}

	return archive, nil
}

// Exists reports whether an archive with the given name exists.
// Returns (true, nil) when found; (false, nil) when not found; (false, err) on other errors.
func (r *Registry) Exists(name string) (bool, error) {
	_, err := r.Get(name)
	if err != nil {
		// TODO: differentiate not-found via sentinel; for now, treat any error as not found
		return false, err
	}
	return true, nil
}

// List returns all archives
func (r *Registry) List() ([]*Archive, error) {
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	ORDER BY created DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list archives: %w", err)
	}
	defer rows.Close()

	var archives []*Archive
	for rows.Next() {
		archive := &Archive{}
		err := rows.Scan(
			&archive.ID,
			&archive.UID,
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
			&archive.Managed,
			&archive.Status,
			&archive.LastSeen,
			&archive.DeletedAt,
			&archive.OriginalPath,
			&archive.Uploaded,
			&archive.Destination,
			&archive.UploadedAt,
			&archive.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan archive: %w", err)
		}
		archives = append(archives, archive)
	}

	return archives, rows.Err()
}

// ListNotUploaded returns archives that haven't been uploaded
func (r *Registry) ListNotUploaded() ([]*Archive, error) {
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE uploaded = FALSE
	ORDER BY created DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list not uploaded archives: %w", err)
	}
	defer rows.Close()

	var archives []*Archive
	for rows.Next() {
		archive := &Archive{}
		err := rows.Scan(
			&archive.ID,
			&archive.UID,
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
			&archive.Managed,
			&archive.Status,
			&archive.LastSeen,
			&archive.DeletedAt,
			&archive.OriginalPath,
			&archive.Uploaded,
			&archive.Destination,
			&archive.UploadedAt,
			&archive.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan archive: %w", err)
		}
		archives = append(archives, archive)
	}

	return archives, rows.Err()
}

// ListOlderThan returns archives older than the specified duration
func (r *Registry) ListOlderThan(duration time.Duration) ([]*Archive, error) {
	cutoff := time.Now().Add(-duration)
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE created < ?
	ORDER BY created DESC
	`

	rows, err := r.db.Query(query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to list older archives: %w", err)
	}
	defer rows.Close()

	var archives []*Archive
	for rows.Next() {
		archive := &Archive{}
		err := rows.Scan(
			&archive.ID,
			&archive.UID,
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
			&archive.Managed,
			&archive.Status,
			&archive.LastSeen,
			&archive.DeletedAt,
			&archive.OriginalPath,
			&archive.Uploaded,
			&archive.Destination,
			&archive.UploadedAt,
			&archive.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan archive: %w", err)
		}
		archives = append(archives, archive)
	}

	return archives, rows.Err()
}

// Update updates an existing archive
func (r *Registry) Update(archive *Archive) error {
	query := `
	UPDATE archives
	SET uid = ?, path = ?, size = ?, checksum = ?, profile = ?, managed = ?, status = ?, last_seen = ?, deleted_at = ?, original_path = ?, uploaded = ?, destination = ?, uploaded_at = ?, metadata = ?
	WHERE id = ?
	`

	_, err := r.db.Exec(query,
		archive.UID,
		archive.Path,
		archive.Size,
		archive.Checksum,
		archive.Profile,
		archive.Managed,
		archive.Status,
		archive.LastSeen,
		archive.DeletedAt,
		archive.OriginalPath,
		archive.Uploaded,
		archive.Destination,
		archive.UploadedAt,
		archive.Metadata,
		archive.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update archive: %w", err)
	}

	return nil
}

// BackfillUIDs sets a UID for rows missing it
func (r *Registry) BackfillUIDs(gen func() string) error {
	rows, err := r.db.Query(`SELECT id FROM archives WHERE uid IS NULL OR uid = ''`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	for _, id := range ids {
		if _, err := r.db.Exec(`UPDATE archives SET uid = ? WHERE id = ?`, gen(), id); err != nil {
			return err
		}
	}
	return nil
}

// GetByID retrieves an archive by numeric id
func (r *Registry) GetByID(id int64) (*Archive, error) {
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE id = ?`
	archive := &Archive{}
	err := r.db.QueryRow(query, id).Scan(
		&archive.ID,
		&archive.UID,
		&archive.Name,
		&archive.Path,
		&archive.Size,
		&archive.Created,
		&archive.Checksum,
		&archive.Profile,
		&archive.Managed,
		&archive.Status,
		&archive.LastSeen,
		&archive.DeletedAt,
		&archive.OriginalPath,
		&archive.Uploaded,
		&archive.Destination,
		&archive.UploadedAt,
		&archive.Metadata,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("archive not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get archive by id: %w", err)
	}
	return archive, nil
}

// GetByUID retrieves an archive by exact UID
func (r *Registry) GetByUID(uid string) (*Archive, error) {
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE uid = ?`
	archive := &Archive{}
	err := r.db.QueryRow(query, uid).Scan(
		&archive.ID,
		&archive.UID,
		&archive.Name,
		&archive.Path,
		&archive.Size,
		&archive.Created,
		&archive.Checksum,
		&archive.Profile,
		&archive.Managed,
		&archive.Status,
		&archive.LastSeen,
		&archive.DeletedAt,
		&archive.OriginalPath,
		&archive.Uploaded,
		&archive.Destination,
		&archive.UploadedAt,
		&archive.Metadata,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("archive not found: %s", uid)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get archive by uid: %w", err)
	}
	return archive, nil
}

// FindByUIDPrefix returns archives whose UID starts with prefix
func (r *Registry) FindByUIDPrefix(prefix string, limit int) ([]*Archive, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE uid LIKE ?
	ORDER BY created DESC
	LIMIT ?`
	rows, err := r.db.Query(query, prefix+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query by uid prefix: %w", err)
	}
	defer rows.Close()
	var out []*Archive
	for rows.Next() {
		archive := &Archive{}
		if err := rows.Scan(
			&archive.ID,
			&archive.UID,
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
			&archive.Managed,
			&archive.Status,
			&archive.LastSeen,
			&archive.DeletedAt,
			&archive.OriginalPath,
			&archive.Uploaded,
			&archive.Destination,
			&archive.UploadedAt,
			&archive.Metadata,
		); err != nil {
			return nil, err
		}
		out = append(out, archive)
	}
	return out, rows.Err()
}

// FindByChecksumPrefix returns archives whose checksum starts with prefix
func (r *Registry) FindByChecksumPrefix(prefix string, limit int) ([]*Archive, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
	SELECT id, uid, name, path, size, created, checksum, profile, managed, status, last_seen, deleted_at, original_path, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE checksum LIKE ?
	ORDER BY created DESC
	LIMIT ?`
	rows, err := r.db.Query(query, prefix+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query by checksum prefix: %w", err)
	}
	defer rows.Close()
	var out []*Archive
	for rows.Next() {
		archive := &Archive{}
		if err := rows.Scan(
			&archive.ID,
			&archive.UID,
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
			&archive.Managed,
			&archive.Status,
			&archive.LastSeen,
			&archive.DeletedAt,
			&archive.OriginalPath,
			&archive.Uploaded,
			&archive.Destination,
			&archive.UploadedAt,
			&archive.Metadata,
		); err != nil {
			return nil, err
		}
		out = append(out, archive)
	}
	return out, rows.Err()
}

// Delete removes an archive from the registry
func (r *Registry) Delete(name string) error {
	query := `DELETE FROM archives WHERE name = ?`
	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete archive: %w", err)
	}
	return nil
}

// Path returns the underlying database file path (for backups)
func (r *Registry) Path() string { return r.dbPath }

// Close closes the database connection
func (r *Registry) Close() error {
	return r.db.Close()
}
