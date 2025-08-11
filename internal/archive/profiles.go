package archive

import (
	"os"
	"path/filepath"
	"strings"
)

// CompressionProfile defines optimal 7z parameters for different content types
type CompressionProfile struct {
	Name           string
	Description    string
	Level          int    // -mx parameter
	DictionarySize string // -md parameter
	FastBytes      int    // -mfb parameter
	SolidMode      bool   // -ms parameter
	Algorithm      string // compression algorithm
}

// Predefined compression profiles
var profiles = map[string]CompressionProfile{
	"media": {
		Name:           "Media",
		Description:    "Optimized for video, audio, and images (faster, good ratio for pre-compressed data)",
		Level:          3,
		DictionarySize: "16m",
		FastBytes:      32,
		SolidMode:      false,
		Algorithm:      "lzma2",
	},
	"documents": {
		Name:           "Documents",
		Description:    "Optimized for text, code, and office files (maximum compression)",
		Level:          9,
		DictionarySize: "64m",
		FastBytes:      273,
		SolidMode:      true,
		Algorithm:      "lzma2",
	},
	"balanced": {
		Name:           "Balanced",
		Description:    "Good compression with reasonable speed for mixed content",
		Level:          7,
		DictionarySize: "32m",
		FastBytes:      64,
		SolidMode:      true,
		Algorithm:      "lzma2",
	},
}

// ContentStats holds analysis results of directory contents
type ContentStats struct {
	TotalBytes      int64
	TotalFiles      int
	MediaBytes      int64
	MediaFiles      int
	DocumentBytes   int64
	DocumentFiles   int
	CompressedBytes int64
	CompressedFiles int
	OtherBytes      int64
	OtherFiles      int
}

// GetProfile returns a compression profile by name
func GetProfile(name string) (CompressionProfile, bool) {
	profile, exists := profiles[name]
	return profile, exists
}

// ListProfiles returns all available profiles
func ListProfiles() []CompressionProfile {
	result := make([]CompressionProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, profile)
	}
	return result
}

// AnalyzeContent examines directory contents and recommends optimal compression profile
func AnalyzeContent(sourcePath string) (*ContentStats, CompressionProfile, error) {
	return AnalyzeContentWithThresholds(sourcePath, 70, 60)
}

// AnalyzeContentWithThresholds allows custom thresholds for media/docs percentages
func AnalyzeContentWithThresholds(sourcePath string, mediaThreshold int, docsThreshold int) (*ContentStats, CompressionProfile, error) {
	stats := &ContentStats{}

	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		size := info.Size()
		stats.TotalBytes += size
		stats.TotalFiles++
		switch {
		case isMediaFile(ext):
			stats.MediaBytes += size
			stats.MediaFiles++
		case isDocumentFile(ext):
			stats.DocumentBytes += size
			stats.DocumentFiles++
		case isCompressedFile(ext):
			stats.CompressedBytes += size
			stats.CompressedFiles++
		default:
			stats.OtherBytes += size
			stats.OtherFiles++
		}
		return nil
	})
	if err != nil {
		return nil, CompressionProfile{}, err
	}
	// Recommend profile based on content analysis and custom thresholds
	recommended := recommendProfileWithThresholds(stats, mediaThreshold, docsThreshold)
	return stats, recommended, nil
}

// isMediaFile checks if file extension indicates media content
func isMediaFile(ext string) bool {
	mediaExts := map[string]bool{
		// Video
		".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".wmv": true,
		".flv": true, ".webm": true, ".m4v": true, ".3gp": true, ".mpg": true,
		".mpeg": true, ".ts": true, ".mts": true,

		// Audio
		".mp3": true, ".wav": true, ".flac": true, ".aac": true, ".ogg": true,
		".wma": true, ".m4a": true, ".opus": true, ".aiff": true,

		// Images
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
		".tiff": true, ".webp": true, ".svg": true, ".ico": true, ".raw": true,
		".cr2": true, ".nef": true, ".arw": true, ".dng": true,
	}
	return mediaExts[ext]
}

// isDocumentFile checks if file extension indicates document/text content
func isDocumentFile(ext string) bool {
	docExts := map[string]bool{
		// Text/Code
		".txt": true, ".md": true, ".rst": true, ".rtf": true,
		".go": true, ".py": true, ".js": true, ".ts": true, ".java": true,
		".c": true, ".cpp": true, ".h": true, ".hpp": true, ".cs": true,
		".php": true, ".rb": true, ".rs": true, ".swift": true, ".kt": true,
		".html": true, ".css": true, ".scss": true, ".less": true,
		".xml": true, ".json": true, ".yaml": true, ".yml": true, ".toml": true,
		".sh": true, ".bat": true, ".ps1": true, ".sql": true,

		// Office Documents
		".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true, ".odt": true, ".ods": true, ".odp": true,
		".pdf": true, ".tex": true, ".epub": true, ".mobi": true,

		// Configuration
		".conf": true, ".cfg": true, ".ini": true, ".properties": true,
		".env": true, ".gitignore": true, ".dockerignore": true,
	}
	return docExts[ext]
}

// isCompressedFile checks if file is already compressed
func isCompressedFile(ext string) bool {
	compressedExts := map[string]bool{
		".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
		".bz2": true, ".xz": true, ".lz": true, ".lzma": true, ".zst": true,
		".jar": true, ".war": true, ".ear": true, ".apk": true, ".ipa": true,
	}
	return compressedExts[ext]
}

// recommendProfile analyzes content statistics and recommends best profile
func recommendProfile(stats *ContentStats) CompressionProfile {
	return recommendProfileWithThresholds(stats, 70, 60)
}

// recommendProfileWithThresholds uses configurable thresholds (percent values)
func recommendProfileWithThresholds(stats *ContentStats, mediaThreshold int, docsThreshold int) CompressionProfile {
	if stats.TotalBytes == 0 {
		return profiles["balanced"] // Default fallback
	}
	mediaPercent := float64(stats.MediaBytes) / float64(stats.TotalBytes) * 100
	docPercent := float64(stats.DocumentBytes) / float64(stats.TotalBytes) * 100
	compressedPercent := float64(stats.CompressedBytes) / float64(stats.TotalBytes) * 100
	if mediaPercent >= float64(mediaThreshold) {
		return profiles["media"]
	}
	if docPercent >= float64(docsThreshold) {
		return profiles["documents"]
	}
	if compressedPercent >= 50.0 {
		return profiles["media"] // Use media settings for already compressed data
	}
	return profiles["balanced"]
}
