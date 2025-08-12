package archive

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Archive represents a 7z archive with metadata
type Archive struct {
	Path         string
	Size         int64
	FileCount    int
	Created      time.Time
	Checksum     string
	Metadata     *Metadata
	OriginalSize int64
	Profile      CompressionProfile // Profile used for compression
}

// Metadata contains archive metadata
type Metadata struct {
	Files       []FileInfo
	Created     time.Time
	Compression string
	Notes       string
}

// FileInfo represents a file in the archive
type FileInfo struct {
	Path     string
	Size     int64
	Modified time.Time
	Mode     os.FileMode
}

// Manager handles archive operations
type Manager struct {
	// Can add configuration here
}

// NewManager creates a new archive manager
func NewManager() *Manager {
	return &Manager{}
}

// CreateOptions contains options for creating archives
type CreateOptions struct {
	Source           string
	Output           string
	CompressionLevel int
	Threads          int
	Exclude          []string
	Profile          string // Compression profile name
	SmartCompression bool   // Auto-detect optimal profile (deprecated - now default)
	Comprehensive    bool   // Create log and checksums
	Force            bool   // Overwrite existing files
	// Config-driven thresholds (percent values); 0 means use defaults
	MediaThreshold int
	DocsThreshold  int
}

// Create creates a new archive
func (m *Manager) Create(ctx context.Context, opts CreateOptions) (*Archive, error) {
	var profile CompressionProfile
	var err error

	// Always analyze content to educate the user, with config-driven thresholds
	mediaTh := opts.MediaThreshold
	docsTh := opts.DocsThreshold
	if mediaTh <= 0 {
		mediaTh = 70
	}
	if docsTh <= 0 {
		docsTh = 60
	}
	stats, recommended, analyzeErr := AnalyzeContentWithThresholds(opts.Source, mediaTh, docsTh)
	if analyzeErr != nil {
		// Don't fail on analysis error, just skip the educational output
		fmt.Printf("⚠️  Content analysis unavailable: %v\n\n", analyzeErr)
	} else {
		// Show content breakdown to educate user
		fmt.Printf("📊 Content Analysis:\n")
		fmt.Printf("  Total: %d files, %.1f MB\n", stats.TotalFiles, float64(stats.TotalBytes)/(1024*1024))
		if stats.MediaFiles > 0 {
			mediaPercent := float64(stats.MediaBytes) / float64(stats.TotalBytes) * 100
			fmt.Printf("  Media: %d files (%.1f%%), %.1f MB\n", stats.MediaFiles, mediaPercent, float64(stats.MediaBytes)/(1024*1024))
		}
		if stats.DocumentFiles > 0 {
			docPercent := float64(stats.DocumentBytes) / float64(stats.TotalBytes) * 100
			fmt.Printf("  Documents: %d files (%.1f%%), %.1f MB\n", stats.DocumentFiles, docPercent, float64(stats.DocumentBytes)/(1024*1024))
		}
		if stats.CompressedFiles > 0 {
			compPercent := float64(stats.CompressedBytes) / float64(stats.TotalBytes) * 100
			fmt.Printf("  Compressed: %d files (%.1f%%), %.1f MB\n", stats.CompressedFiles, compPercent, float64(stats.CompressedBytes)/(1024*1024))
		}
		if stats.OtherFiles > 0 {
			otherPercent := float64(stats.OtherBytes) / float64(stats.TotalBytes) * 100
			fmt.Printf("  Other: %d files (%.1f%%), %.1f MB\n", stats.OtherFiles, otherPercent, float64(stats.OtherBytes)/(1024*1024))
		}
		fmt.Printf("\n")
	}

	// Determine which compression profile to use
	if opts.Profile != "" {
		// Use specified profile
		var exists bool
		profile, exists = GetProfile(opts.Profile)
		if !exists {
			return nil, fmt.Errorf("unknown compression profile: %s", opts.Profile)
		}
		fmt.Printf("🎯 Using Profile: %s\n", profile.Name)
		fmt.Printf("   %s\n", profile.Description)
		fmt.Printf("   Settings: Level %d, Dictionary %s, Fast bytes %d\n\n",
			profile.Level, profile.DictionarySize, profile.FastBytes)
	} else if opts.CompressionLevel > 0 {
		// Manual compression level specified - use traditional mode
		profile = CompressionProfile{
			Level:          opts.CompressionLevel,
			DictionarySize: "32m",
			FastBytes:      64,
			SolidMode:      true,
			Algorithm:      "lzma2",
		}

		// Educational message about available optimizations
		if analyzeErr == nil {
			fmt.Printf("💡 Optimization Tip: Based on your content, --profile %s might be faster\n",
				strings.ToLower(recommended.Name))
			fmt.Printf("   Run '7zarch-go profiles' to see all available profiles\n\n")
		}
	} else {
		// Smart compression by default - use recommended profile
		if analyzeErr != nil {
			// Fallback to balanced if analysis failed
			profile, _ = GetProfile("balanced")
		} else {
			fmt.Printf("🎯 Using Smart Profile: %s\n", recommended.Name)
			fmt.Printf("   %s\n", recommended.Description)
			fmt.Printf("   Settings: Level %d, Dictionary %s, Fast bytes %d\n\n",
				recommended.Level, recommended.DictionarySize, recommended.FastBytes)

			profile = recommended
		}
	}

	// Build 7z command
	args := []string{"a"}

	// Add output file
	args = append(args, opts.Output)

	// Force overwrite without prompting
	args = append(args, "-y")

	// Apply compression profile parameters
	args = append(args, "-t7z")
	args = append(args, fmt.Sprintf("-m0=%s", profile.Algorithm))
	args = append(args, fmt.Sprintf("-mx=%d", profile.Level))
	args = append(args, fmt.Sprintf("-mfb=%d", profile.FastBytes))
	args = append(args, fmt.Sprintf("-md=%s", profile.DictionarySize))
	if profile.SolidMode {
		args = append(args, "-ms=on")
	} else {
		args = append(args, "-ms=off")
	}

	// Add thread count if specified
	if opts.Threads > 0 {
		args = append(args, fmt.Sprintf("-mmt=%d", opts.Threads))
	}

	// Add source
	args = append(args, opts.Source)

	// Add excludes
	for _, exclude := range opts.Exclude {
		args = append(args, fmt.Sprintf("-x!%s", exclude))
	}

	// Execute 7z command
	cmd := exec.CommandContext(ctx, "7z", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("7z failed: %w\nOutput: %s", err, string(output))
	}

	// Get archive info
	info, err := os.Stat(opts.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to stat archive: %w", err)
	}

	// Calculate checksum
	checksum, err := calculateFileChecksum(opts.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Get file count from output
	fileCount := extractFileCount(string(output))

	archive := &Archive{
		Path:      opts.Output,
		Size:      info.Size(),
		FileCount: fileCount,
		Created:   time.Now(),
		Checksum:  checksum,
		Profile:   profile,
	}

	// Use analysis totals as original size to avoid a second directory walk
	if analyzeErr == nil && stats != nil {
		archive.OriginalSize = stats.TotalBytes
	}

	// Handle comprehensive mode (create log and checksums)
	if opts.Comprehensive {
		// Create log file
		logPath := archive.Path + ".log"
		if err := CreateLogFile(logPath, archive, opts.Source); err != nil {
			fmt.Printf("Warning: Failed to create log: %v\n", err)
		} else {
			fmt.Printf("Log created: %s\n", logPath)
		}

		// Create checksum file
		checksumPath := archive.Path + ".sha256"
		if err := CreateChecksumFile(checksumPath, archive); err != nil {
			fmt.Printf("Warning: Failed to create checksum: %v\n", err)
		} else {
			fmt.Printf("Checksum created: %s\n", checksumPath)
		}
	}

	return archive, nil
}

// TestResult contains the results of archive testing
type TestResult struct {
	Passed        bool
	ChecksumValid bool
	MetadataValid bool
	FilesVerified int
	Errors        []string
	Duration      time.Duration
}

// Test verifies archive integrity
func (m *Manager) Test(ctx context.Context, archivePath string) (*TestResult, error) {
	startTime := time.Now()
	result := &TestResult{
		Passed: true,
		Errors: []string{},
	}

	// Test 1: Archive structure integrity
	if err := m.testArchiveIntegrity(ctx, archivePath); err != nil {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("Archive integrity: %v", err))
	}

	// Test 2: Checksum verification (if .sha256 exists)
	checksumFile := archivePath + ".sha256"
	if _, err := os.Stat(checksumFile); err == nil {
		if err := m.verifyChecksum(archivePath, checksumFile); err != nil {
			result.Passed = false
			result.ChecksumValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Checksum: %v", err))
		} else {
			result.ChecksumValid = true
		}
	}

	// Test 3: Metadata validation (if .log exists)
	metadataFile := archivePath + ".log"
	if _, err := os.Stat(metadataFile); err == nil {
		if err := m.validateMetadata(ctx, archivePath, metadataFile); err != nil {
			result.Passed = false
			result.MetadataValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Metadata: %v", err))
		} else {
			result.MetadataValid = true
		}
	}

	// Test 4: List files (quick extraction test)
	fileCount, err := m.listArchiveFiles(ctx, archivePath)
	if err != nil {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("File listing: %v", err))
	} else {
		result.FilesVerified = fileCount
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// testArchiveIntegrity runs 7z test command and relies on exit code for success
func (m *Manager) testArchiveIntegrity(ctx context.Context, archivePath string) error {
	cmd := exec.CommandContext(ctx, "7z", "t", archivePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("integrity test failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// countPathsInSlt counts file entries by scanning for "Path = " lines in -slt output
func countPathsInSlt(output string) int {
	lines := strings.Split(output, "\n")
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Path = ") {
			count++
		}
	}
	return count
}

// verifyChecksum compares archive checksum with stored value
func (m *Manager) verifyChecksum(archivePath, checksumFile string) error {
	// Read expected checksum
	data, err := os.ReadFile(checksumFile)
	if err != nil {
		return fmt.Errorf("failed to read checksum file: %w", err)
	}

	// Parse checksum (format: "hash  filename")
	parts := strings.Fields(string(data))
	if len(parts) < 1 {
		return fmt.Errorf("invalid checksum file format")
	}
	expectedChecksum := parts[0]

	// Calculate actual checksum
	actualChecksum, err := calculateFileChecksum(archivePath)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Compare
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// validateMetadata checks metadata consistency
func (m *Manager) validateMetadata(ctx context.Context, archivePath, metadataFile string) error {
	// This would parse the .log file and validate against archive contents
	// For now, just check if we can read it
	_, err := os.ReadFile(metadataFile)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	// TODO: Parse and validate metadata against actual archive contents

	return nil
}

// listArchiveFiles counts files in archive using structured listing (-slt)
func (m *Manager) listArchiveFiles(ctx context.Context, archivePath string) (int, error) {
	cmd := exec.CommandContext(ctx, "7z", "l", "-slt", "-scsUTF-8", archivePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to list files: %w\nOutput: %s", err, string(output))
	}

	// In -slt output, each file section starts with "Path = ...". Count these.
	lines := strings.Split(string(output), "\n")
	fileCount := 0
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Path = ") {
			fileCount++
		}
	}
	return fileCount, nil
}

// Helper functions

func calculateFileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func calculateDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func extractFileCount(output string) int {
	// Try to parse totals from standard output summary first
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)
		// Example patterns observed:
		// "Files: 15"
		// "15 files, 2048576 bytes"
		if strings.HasPrefix(l, "Files:") {
			var label string
			var count int
			if _, err := fmt.Sscanf(l, "%s %d", &label, &count); err == nil {
				return count
			}
		}
		if strings.Contains(l, " files") {
			parts := strings.Fields(l)
			if len(parts) > 0 {
				var count int
				if _, err := fmt.Sscanf(parts[0], "%d", &count); err == nil {
					return count
				}
			}
		}
	}
	return 0
}

// CreateLogFile creates a metadata log file for the archive
func CreateLogFile(logPath string, archive *Archive, sourcePath string) error {
	file, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write metadata
	fmt.Fprintf(file, "Archive: %s\n", archive.Path)
	fmt.Fprintf(file, "Created: %s\n", archive.Created.Format(time.RFC3339))
	fmt.Fprintf(file, "Source: %s\n", sourcePath)
	fmt.Fprintf(file, "Size: %d bytes\n", archive.Size)
	fmt.Fprintf(file, "Files: %d\n", archive.FileCount)
	fmt.Fprintf(file, "Checksum: %s\n", archive.Checksum)

	if archive.OriginalSize > 0 {
		ratio := float64(archive.Size) / float64(archive.OriginalSize) * 100
		fmt.Fprintf(file, "Compression: %.1f%%\n", ratio)
	}

	// TODO: Add file listing with details

	return nil
}

// CreateChecksumFile creates a SHA256 checksum file
func CreateChecksumFile(checksumPath string, archive *Archive) error {
	file, err := os.Create(checksumPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write checksum in standard format: "hash  filename"
	fmt.Fprintf(file, "%s  %s\n", archive.Checksum, filepath.Base(archive.Path))

	return nil
}
