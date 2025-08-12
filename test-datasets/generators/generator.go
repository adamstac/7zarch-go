// Package generators provides comprehensive test dataset generation for 7zarch-go
// Built on 7EP-0006's proven metadata-only approach for efficient testing
package generators

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// ScenarioSpec defines what kind of test scenario to generate
// Simplified from DatasetSpec based on 7EP-0006 learnings
type ScenarioSpec struct {
	Name        string
	Count       int                   // Number of archives
	ULIDPattern ULIDPattern           // How to generate ULIDs
	Profiles    []ProfileDistribution // Profile mix
	TimeSpread  time.Duration         // Spread creation times over period
	SizePattern SizePattern           // Size distribution
	EdgeCases   []EdgeCase            // Special scenarios to include
}

// ULIDPattern defines how to generate ULID patterns for testing
type ULIDPattern string

const (
	ULIDUnique     ULIDPattern = "unique"     // All unique (7EP-0006 scaling tests)
	ULIDSimilar    ULIDPattern = "similar"    // Controlled similarity (7EP-0006 disambiguation)
	ULIDCollisions ULIDPattern = "collisions" // Intentional prefix overlaps
)

// SizePattern defines file size distribution patterns
type SizePattern string

const (
	SizeUniform   SizePattern = "uniform"   // Even distribution
	SizeRealistic SizePattern = "realistic" // Log-normal (most files small, few large)
	SizeLargeFiles SizePattern = "large"    // Focus on large file scenarios
)

// EdgeCase defines special scenarios to test
type EdgeCase string

const (
	UnicodeFilenames      EdgeCase = "unicode_filenames"
	EmojiFilenames        EdgeCase = "emoji_filenames"
	SpecialCharacters     EdgeCase = "special_characters"
	LongFilenames         EdgeCase = "long_filenames"
	MaxPathLength         EdgeCase = "max_path_length"
	MinMaxFileSizes       EdgeCase = "min_max_file_sizes"
	TimeBoundaries        EdgeCase = "time_boundaries"
	MixedManagedExternal  EdgeCase = "mixed_managed_external"
	TimeSequencing        EdgeCase = "time_sequencing"
	ManagedExternalMix    EdgeCase = "managed_external_mix"
	CrossProfileFiltering EdgeCase = "cross_profile_filtering"
)

// ProfileDistribution defines how to distribute archive profiles
type ProfileDistribution struct {
	Profile string
	Weight  float32 // 0.0-1.0
}

// Archive represents a test archive (simplified from storage.Archive)
type Archive struct {
	UID          string
	Name         string
	Path         string
	Size         int64
	Created      time.Time
	Profile      string
	Managed      bool
	Status       string
	Checksum     string
	Uploaded     bool
	UploadedAt   *time.Time
	DeletedAt    *time.Time
	OriginalPath string
}

// Generator creates reproducible test scenarios using 7EP-0006 patterns
type Generator struct {
	seed int64
	rng  *rand.Rand
}

// NewGenerator creates a new test data generator with fixed seed
func NewGenerator(seed int64) *Generator {
	return &Generator{
		seed: seed,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// GenerateScenario creates archives based on scenario spec
// Builds directly on 7EP-0006's successful generateTestArchives() pattern
func (g *Generator) GenerateScenario(tb testing.TB, spec ScenarioSpec) []*Archive {
	tb.Helper()

	archives := make([]*Archive, spec.Count)
	profiles := extractProfiles(spec.Profiles)
	sizes := generateSizeDistribution(spec.SizePattern)
	baseTime := time.Now().Add(-spec.TimeSpread)

	for i := 0; i < spec.Count; i++ {
		archive := &Archive{
			UID:      g.generateUID(i, spec.ULIDPattern),
			Name:     fmt.Sprintf("%s-%04d.7z", spec.Name, i),
			Path:     fmt.Sprintf("/tmp/test-%s-%04d.7z", spec.Name, i),
			Size:     selectSize(sizes, g.rng),
			Created:  g.generateCreationTime(baseTime, spec.TimeSpread, i, spec.Count),
			Profile:  selectProfile(profiles, g.rng),
			Managed:  g.rng.Float32() < 0.9, // 90% managed (7EP-0006 pattern)
			Status:   "present",
			Checksum: fmt.Sprintf("sha256:%064x", i),
		}

		// Apply edge case modifications
		g.applyEdgeCases(archive, spec.EdgeCases, i)

		archives[i] = archive
	}

	return archives
}

// generateUID creates a ULID based on the pattern
func (g *Generator) generateUID(i int, pattern ULIDPattern) string {
	switch pattern {
	case ULIDUnique:
		// Each ULID is completely unique - for scaling tests
		return fmt.Sprintf("01K%02d%05d%012d%05d",
			g.rng.Intn(26), i, g.rng.Int63n(999999999999), g.rng.Intn(99999))

	case ULIDSimilar:
		// Controlled similarity for disambiguation testing (7EP-0006 pattern)
		if i < 100 {
			// Group by tens for first 100: 01K2E00, 01K2E01, 01K2E02, etc.
			return fmt.Sprintf("01K2E%02d%012d%08d", i/10, i, i*17)
		}
		return fmt.Sprintf("01K2F%02d%012d%08d", (i-100)/100, i, i*23)

	case ULIDCollisions:
		// Intentional prefix collisions for stress testing resolution
		prefixCount := min(i/10+1, 5) // Group into 5 prefix buckets max
		return fmt.Sprintf("01K2G%02d%012d%08d", prefixCount, i, g.rng.Intn(99999999))

	default:
		return g.generateUID(i, ULIDUnique) // Safe default
	}
}

// generateCreationTime creates a time within the specified spread
func (g *Generator) generateCreationTime(baseTime time.Time, spread time.Duration, index, total int) time.Time {
	if total <= 1 {
		return baseTime
	}
	// Distribute evenly across the time spread
	offset := time.Duration(float64(spread) * float64(index) / float64(total-1))
	return baseTime.Add(offset)
}

// applyEdgeCases modifies archive based on edge case requirements
func (g *Generator) applyEdgeCases(archive *Archive, edgeCases []EdgeCase, index int) {
	for _, ec := range edgeCases {
		switch ec {
		case UnicodeFilenames:
			names := []string{"测试文件", "файл", "ファイル", "αρχείο", "파일"}
			archive.Name = fmt.Sprintf("%s-%04d.7z", names[index%len(names)], index)

		case EmojiFilenames:
			emojis := []string{"🚀", "📝", "🎉", "🔥", "💾"}
			archive.Name = fmt.Sprintf("%sproject-%04d.7z", emojis[index%len(emojis)], index)

		case SpecialCharacters:
			special := []string{"file with spaces", "file[brackets]", "file(parens)", "file@symbol", "file#hash"}
			archive.Name = fmt.Sprintf("%s-%04d.7z", special[index%len(special)], index)

		case LongFilenames:
			// Create a 255-character filename
			longName := fmt.Sprintf("%s", repeatString("a", 240))
			archive.Name = fmt.Sprintf("%s-%04d.7z", longName, index)

		case MaxPathLength:
			// Create deep directory hierarchy
			deepPath := ""
			for j := 0; j < 20; j++ {
				deepPath += fmt.Sprintf("/level%02d", j)
			}
			archive.Path = fmt.Sprintf("%s/archive-%04d.7z", deepPath, index)

		case MinMaxFileSizes:
			if index%3 == 0 {
				archive.Size = 0 // Zero-byte file
			} else if index%3 == 1 {
				archive.Size = 10 * 1024 * 1024 * 1024 // 10GB file
			}

		case TimeBoundaries:
			if index%3 == 0 {
				archive.Created = time.Unix(0, 0) // Unix epoch
			} else if index%3 == 1 {
				archive.Created = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC) // Far future
			}

		case MixedManagedExternal:
			archive.Managed = index%2 == 0 // Alternate between managed and external

		case TimeSequencing:
			// Already handled in creation time generation
			continue

		case ManagedExternalMix:
			// 70% managed, 30% external
			archive.Managed = g.rng.Float32() < 0.7

		case CrossProfileFiltering:
			// Ensure diverse profile distribution
			profiles := []string{"documents", "media", "balanced", "code", "data"}
			archive.Profile = profiles[index%len(profiles)]
		}
	}
}

// Helper functions

func extractProfiles(distributions []ProfileDistribution) []string {
	if len(distributions) == 0 {
		return []string{"balanced"}
	}

	profiles := make([]string, 0)
	for _, dist := range distributions {
		// Add profile proportionally to its weight
		count := int(dist.Weight * 100)
		for i := 0; i < count; i++ {
			profiles = append(profiles, dist.Profile)
		}
	}
	return profiles
}

func generateSizeDistribution(pattern SizePattern) []int64 {
	switch pattern {
	case SizeUniform:
		return []int64{1024, 100 * 1024, 10 * 1024 * 1024} // 1KB, 100KB, 10MB

	case SizeRealistic:
		// 70% small, 25% medium, 5% large (realistic distribution)
		sizes := make([]int64, 0, 100)
		for i := 0; i < 70; i++ {
			sizes = append(sizes, 1024+int64(i*100)) // 1-8KB
		}
		for i := 0; i < 25; i++ {
			sizes = append(sizes, 100*1024+int64(i*10000)) // 100KB-350KB
		}
		for i := 0; i < 5; i++ {
			sizes = append(sizes, 10*1024*1024+int64(i*1024*1024)) // 10MB-15MB
		}
		return sizes

	case SizeLargeFiles:
		return []int64{100 * 1024 * 1024, 500 * 1024 * 1024, 1024 * 1024 * 1024} // 100MB, 500MB, 1GB

	default:
		return []int64{1024, 100 * 1024, 10 * 1024 * 1024} // Safe default
	}
}

func selectSize(sizes []int64, rng *rand.Rand) int64 {
	if len(sizes) == 0 {
		return 1024 // Default 1KB
	}
	return sizes[rng.Intn(len(sizes))]
}

func selectProfile(profiles []string, rng *rand.Rand) string {
	if len(profiles) == 0 {
		return "balanced"
	}
	return profiles[rng.Intn(len(profiles))]
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}