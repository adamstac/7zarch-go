package generators

import (
	"fmt"
	"time"
)

// PredefinedScenarios contains all available test scenarios
var PredefinedScenarios = map[string]ScenarioSpec{
	// Performance Testing Scenarios
	"disambiguation-stress": {
		Name:        "disambiguation-stress",
		Count:       1000,
		ULIDPattern: ULIDSimilar, // Creates controlled similarity like 7EP-0006
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 0.7},
			{Profile: "media", Weight: 0.2},
			{Profile: "balanced", Weight: 0.1},
		},
		SizePattern: SizeUniform,
		TimeSpread:  30 * 24 * time.Hour, // 30 days
	},

	"scaling-validation": {
		Name:        "scaling-validation",
		Count:       10000,
		ULIDPattern: ULIDUnique, // Each unique for pure scaling test
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 0.4},
			{Profile: "media", Weight: 0.3},
			{Profile: "balanced", Weight: 0.3},
		},
		SizePattern: SizeRealistic,        // Most small, few large
		TimeSpread:  365 * 24 * time.Hour, // Full year spread
	},

	"resolution-stress": {
		Name:        "resolution-stress",
		Count:       5000,
		ULIDPattern: ULIDCollisions, // Intentional prefix collisions
		Profiles: []ProfileDistribution{
			{Profile: "balanced", Weight: 1.0},
		},
		SizePattern: SizeUniform,
		TimeSpread:  90 * 24 * time.Hour, // 3 months
	},

	// Integration Testing Scenarios
	"create-list-show-delete": {
		Name:        "create-list-show-delete",
		Count:       50,
		ULIDPattern: ULIDUnique,
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 0.5},
			{Profile: "media", Weight: 0.3},
			{Profile: "balanced", Weight: 0.2},
		},
		SizePattern: SizeRealistic,
		TimeSpread:  7 * 24 * time.Hour, // Week timeline
		EdgeCases:   []EdgeCase{MixedManagedExternal, TimeSequencing},
	},

	"mixed-storage-scenario": {
		Name:        "mixed-storage-scenario",
		Count:       100,
		ULIDPattern: ULIDSimilar,
		Profiles: []ProfileDistribution{
			{Profile: "balanced", Weight: 1.0},
		},
		SizePattern: SizeUniform,
		TimeSpread:  30 * 24 * time.Hour,
		EdgeCases:   []EdgeCase{ManagedExternalMix, CrossProfileFiltering},
	},

	"time-series-archives": {
		Name:        "time-series-archives",
		Count:       365, // One per day for a year
		ULIDPattern: ULIDUnique,
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 0.6},
			{Profile: "media", Weight: 0.3},
			{Profile: "data", Weight: 0.1},
		},
		SizePattern: SizeRealistic,
		TimeSpread:  365 * 24 * time.Hour,
		EdgeCases:   []EdgeCase{TimeSequencing},
	},

	// Edge Case Testing Scenarios
	"unicode-names": {
		Name:        "unicode-names",
		Count:       25,
		ULIDPattern: ULIDUnique,
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 1.0},
		},
		SizePattern: SizeUniform,
		TimeSpread:  24 * time.Hour,
		EdgeCases: []EdgeCase{
			UnicodeFilenames,
			EmojiFilenames,
			SpecialCharacters,
			LongFilenames,
		},
	},

	"boundary-conditions": {
		Name:        "boundary-conditions",
		Count:       30,
		ULIDPattern: ULIDCollisions, // Stress test resolution
		Profiles: []ProfileDistribution{
			{Profile: "balanced", Weight: 1.0},
		},
		SizePattern: SizeLargeFiles, // Large file edge cases
		TimeSpread:  24 * time.Hour,
		EdgeCases: []EdgeCase{
			MaxPathLength,
			MinMaxFileSizes,
			TimeBoundaries,
		},
	},

	// Small test scenarios for unit tests
	"small-test": {
		Name:        "small-test",
		Count:       10,
		ULIDPattern: ULIDUnique,
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 1.0},
		},
		SizePattern: SizeUniform,
		TimeSpread:  24 * time.Hour,
	},

	"medium-test": {
		Name:        "medium-test",
		Count:       100,
		ULIDPattern: ULIDSimilar,
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 0.5},
			{Profile: "media", Weight: 0.5},
		},
		SizePattern: SizeRealistic,
		TimeSpread:  7 * 24 * time.Hour,
	},

	"large-test": {
		Name:        "large-test",
		Count:       1000,
		ULIDPattern: ULIDUnique,
		Profiles: []ProfileDistribution{
			{Profile: "documents", Weight: 0.4},
			{Profile: "media", Weight: 0.3},
			{Profile: "balanced", Weight: 0.3},
		},
		SizePattern: SizeRealistic,
		TimeSpread:  30 * 24 * time.Hour,
	},
}

// GetScenario retrieves a predefined scenario by name
func GetScenario(name string) (ScenarioSpec, error) {
	spec, exists := PredefinedScenarios[name]
	if !exists {
		return ScenarioSpec{}, fmt.Errorf("scenario %q not found", name)
	}
	return spec, nil
}

// ListScenarios returns all available scenario names
func ListScenarios() []string {
	names := make([]string, 0, len(PredefinedScenarios))
	for name := range PredefinedScenarios {
		names = append(names, name)
	}
	return names
}

// ScenarioCategories groups scenarios by their purpose
var ScenarioCategories = map[string][]string{
	"performance": {
		"disambiguation-stress",
		"scaling-validation",
		"resolution-stress",
	},
	"integration": {
		"create-list-show-delete",
		"mixed-storage-scenario",
		"time-series-archives",
	},
	"edge-cases": {
		"unicode-names",
		"boundary-conditions",
	},
	"unit-tests": {
		"small-test",
		"medium-test",
		"large-test",
	},
}
