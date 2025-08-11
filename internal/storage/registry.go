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
	db *sql.DB
}

// NewRegistry creates a new registry instance
func NewRegistry(dbPath string) (*Registry, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
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

	r := &Registry{db: db}

	// Initialize the schema
	if err := r.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return r, nil
}

// initSchema creates the database tables if they don't exist
func (r *Registry) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS archives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		size INTEGER NOT NULL,
		created TIMESTAMP NOT NULL,
		checksum TEXT,
		profile TEXT,
		uploaded BOOLEAN DEFAULT FALSE,
		destination TEXT,
		uploaded_at TIMESTAMP,
		metadata TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_archives_created ON archives(created);
	CREATE INDEX IF NOT EXISTS idx_archives_uploaded ON archives(uploaded);
	CREATE INDEX IF NOT EXISTS idx_archives_destination ON archives(destination);
	`

	_, err := r.db.Exec(query)
	return err
}

// Add inserts a new archive into the registry
func (r *Registry) Add(archive *Archive) error {
	query := `
	INSERT INTO archives (name, path, size, created, checksum, profile, uploaded, destination, uploaded_at, metadata)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		archive.Name,
		archive.Path,
		archive.Size,
		archive.Created,
		archive.Checksum,
		archive.Profile,
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
	SELECT id, name, path, size, created, checksum, profile, uploaded, destination, uploaded_at, metadata
	FROM archives
	WHERE name = ?
	`

	archive := &Archive{}
	err := r.db.QueryRow(query, name).Scan(
		&archive.ID,
		&archive.Name,
		&archive.Path,
		&archive.Size,
		&archive.Created,
		&archive.Checksum,
		&archive.Profile,
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

// List returns all archives
func (r *Registry) List() ([]*Archive, error) {
	query := `
	SELECT id, name, path, size, created, checksum, profile, uploaded, destination, uploaded_at, metadata
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
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
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
	SELECT id, name, path, size, created, checksum, profile, uploaded, destination, uploaded_at, metadata
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
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
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
	SELECT id, name, path, size, created, checksum, profile, uploaded, destination, uploaded_at, metadata
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
			&archive.Name,
			&archive.Path,
			&archive.Size,
			&archive.Created,
			&archive.Checksum,
			&archive.Profile,
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
	SET path = ?, size = ?, checksum = ?, profile = ?, uploaded = ?, destination = ?, uploaded_at = ?, metadata = ?
	WHERE id = ?
	`

	_, err := r.db.Exec(query,
		archive.Path,
		archive.Size,
		archive.Checksum,
		archive.Profile,
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

// Delete removes an archive from the registry
func (r *Registry) Delete(name string) error {
	query := `DELETE FROM archives WHERE name = ?`
	_, err := r.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete archive: %w", err)
	}
	return nil
}

// Close closes the database connection
func (r *Registry) Close() error {
	return r.db.Close()
}