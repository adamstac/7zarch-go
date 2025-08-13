package modes

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// DashboardDisplay provides management overview and statistics
type DashboardDisplay struct{}

// NewDashboardDisplay creates a new dashboard display mode
func NewDashboardDisplay() *DashboardDisplay {
	return &DashboardDisplay{}
}

// Name returns the display mode name
func (dd *DashboardDisplay) Name() string {
	return "dashboard"
}

// MinWidth returns the minimum terminal width for this display
func (dd *DashboardDisplay) MinWidth() int {
	return 90
}

// Render displays archives in dashboard format
func (dd *DashboardDisplay) Render(archives []*storage.Archive, opts display.Options) error {
	if len(archives) == 0 {
		dd.printEmptyDashboard()
		return nil
	}

	// Generate statistics
	stats := dd.generateStatistics(archives)

	// Print dashboard header
	dd.printHeader()

	// Print overview section
	dd.printOverview(stats)

	// Print storage breakdown
	dd.printStorageBreakdown(stats)

	// Print status summary
	dd.printStatusSummary(stats)

	// Print profile distribution
	dd.printProfileDistribution(stats)

	// Print recent activity
	dd.printRecentActivity(stats)

	// Print health indicators
	dd.printHealthIndicators(stats)

	return nil
}

// Statistics holds dashboard metrics
type Statistics struct {
	Total               int
	ManagedCount        int
	ExternalCount       int
	ActiveCount         int
	MissingCount        int
	DeletedCount        int
	TotalSize           int64
	ManagedSize         int64
	ExternalSize        int64
	ProfileDistribution map[string]int
	ProfileSizes        map[string]int64
	OldestArchive       *storage.Archive
	NewestArchive       *storage.Archive
	LargestArchive      *storage.Archive
	RecentArchives      []*storage.Archive
	DeletionCandidates  []*storage.Archive
	HealthScore         float64
}

// generateStatistics computes dashboard metrics
func (dd *DashboardDisplay) generateStatistics(archives []*storage.Archive) Statistics {
	stats := Statistics{
		ProfileDistribution: make(map[string]int),
		ProfileSizes:        make(map[string]int64),
	}

	var allActive []*storage.Archive

	for _, archive := range archives {
		stats.Total++
		stats.TotalSize += archive.Size

		// Count by status
		switch archive.Status {
		case "deleted":
			stats.DeletedCount++
		case "missing":
			stats.MissingCount++
			stats.ActiveCount++
		default:
			stats.ActiveCount++
			allActive = append(allActive, archive)
		}

		// Count by location
		if archive.Managed {
			stats.ManagedCount++
			stats.ManagedSize += archive.Size
		} else {
			stats.ExternalCount++
			stats.ExternalSize += archive.Size
		}

		// Profile distribution
		profile := archive.Profile
		if profile == "" {
			profile = "default"
		}
		stats.ProfileDistribution[profile]++
		stats.ProfileSizes[profile] += archive.Size

		// Track extremes
		if stats.OldestArchive == nil || archive.Created.Before(stats.OldestArchive.Created) {
			stats.OldestArchive = archive
		}
		if stats.NewestArchive == nil || archive.Created.After(stats.NewestArchive.Created) {
			stats.NewestArchive = archive
		}
		if stats.LargestArchive == nil || archive.Size > stats.LargestArchive.Size {
			stats.LargestArchive = archive
		}
	}

	// Sort for recent archives (last 5)
	sort.Slice(allActive, func(i, j int) bool {
		return allActive[i].Created.After(allActive[j].Created)
	})

	recentCount := 5
	if len(allActive) < recentCount {
		recentCount = len(allActive)
	}
	stats.RecentArchives = allActive[:recentCount]

	// Calculate health score
	stats.HealthScore = dd.calculateHealthScore(stats)

	return stats
}

// calculateHealthScore computes an overall health percentage
func (dd *DashboardDisplay) calculateHealthScore(stats Statistics) float64 {
	if stats.Total == 0 {
		return 100.0
	}

	score := 100.0

	// Deduct for missing archives
	missingPenalty := float64(stats.MissingCount) / float64(stats.Total) * 30.0
	score -= missingPenalty

	// Deduct for deleted archives (less severe)
	deletedPenalty := float64(stats.DeletedCount) / float64(stats.Total) * 10.0
	score -= deletedPenalty

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// printEmptyDashboard shows dashboard when no archives exist
func (dd *DashboardDisplay) printEmptyDashboard() {
	fmt.Printf("╔══════════════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                               7ZARCH DASHBOARD                               ║\n")
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                              ║\n")
	fmt.Printf("║                         📦 No archives found                                ║\n")
	fmt.Printf("║                                                                              ║\n")
	fmt.Printf("║                 Create archives with '7zarch-go create <path>'               ║\n")
	fmt.Printf("║                          to see them in this dashboard                      ║\n")
	fmt.Printf("║                                                                              ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════════════════════════════╝\n")
}

// printHeader prints the dashboard header
func (dd *DashboardDisplay) printHeader() {
	timestamp := fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05"))
	// Calculate padding for centering the timestamp (80 chars total width)
	padding := (80 - len(timestamp)) / 2
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("╔════════════════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("                               7ZARCH DASHBOARD\n")
	fmt.Printf("%s%s\n", strings.Repeat(" ", padding), timestamp)
	fmt.Printf("╚════════════════════════════════════════════════════════════════════════════════╝\n")
}

// printOverview prints the main statistics overview
func (dd *DashboardDisplay) printOverview(stats Statistics) {
	fmt.Printf("\n┌─ OVERVIEW ─────────────────────────────────────────────────────────────────────┐\n")

	// Align content with section header text
	fmt.Printf("│  Total Archives: %-10d  Storage Used: %-15s  Health: %.1f%%\n",
		stats.Total, display.FormatSize(stats.TotalSize), stats.HealthScore)
	fmt.Printf("│  Active: %-12d  Missing: %-10d  Deleted: %-10d\n",
		stats.ActiveCount, stats.MissingCount, stats.DeletedCount)

	fmt.Printf("└────────────────────────────────────────────────────────────────────────────────┘\n")
}

// printStorageBreakdown shows storage location distribution
func (dd *DashboardDisplay) printStorageBreakdown(stats Statistics) {
	fmt.Printf("\n┌─ STORAGE BREAKDOWN ────────────────────────────────────────────────────────────┐\n")

	managedPercent := 0.0
	externalPercent := 0.0
	if stats.TotalSize > 0 {
		managedPercent = float64(stats.ManagedSize) / float64(stats.TotalSize) * 100
		externalPercent = float64(stats.ExternalSize) / float64(stats.TotalSize) * 100
	}

	// Align content with section header text
	fmt.Printf("│  Managed Storage:  %3d archives  %15s  (%5.1f%%)\n",
		stats.ManagedCount, display.FormatSize(stats.ManagedSize), managedPercent)
	fmt.Printf("│  External Storage: %3d archives  %15s  (%5.1f%%)\n",
		stats.ExternalCount, display.FormatSize(stats.ExternalSize), externalPercent)

	fmt.Printf("└────────────────────────────────────────────────────────────────────────────────┘\n")
}

// printStatusSummary shows detailed status information
func (dd *DashboardDisplay) printStatusSummary(stats Statistics) {
	fmt.Printf("\n┌─ STATUS SUMMARY ───────────────────────────────────────────────────────────────┐\n")

	presentCount := stats.ActiveCount - stats.MissingCount
	presentIcon := display.FormatStatus("present", true)
	fmt.Printf("│  %s Present:  %3d archives\n", presentIcon, presentCount)

	if stats.MissingCount > 0 {
		missingIcon := display.FormatStatus("missing", true)
		fmt.Printf("│  %s Missing:  %3d archives  (requires attention)\n", missingIcon, stats.MissingCount)
	}

	if stats.DeletedCount > 0 {
		deletedIcon := display.FormatStatus("deleted", true)
		fmt.Printf("│  %s Deleted:  %3d archives  (auto-purge in 7 days)\n", deletedIcon, stats.DeletedCount)
	}

	fmt.Printf("└────────────────────────────────────────────────────────────────────────────────┘\n")
}

// printProfileDistribution shows profile usage statistics
func (dd *DashboardDisplay) printProfileDistribution(stats Statistics) {
	if len(stats.ProfileDistribution) == 0 {
		return
	}

	fmt.Printf("\n┌─ PROFILE DISTRIBUTION ─────────────────────────────────────────────────────────┐\n")

	// Sort profiles by count for consistent display
	type profileStat struct {
		name  string
		count int
		size  int64
	}

	var profiles []profileStat
	for name, count := range stats.ProfileDistribution {
		profiles = append(profiles, profileStat{
			name:  name,
			count: count,
			size:  stats.ProfileSizes[name],
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].count > profiles[j].count
	})

	for _, profile := range profiles {
		percent := float64(profile.count) / float64(stats.Total) * 100
		fmt.Printf("│  %-12s: %3d archives  %15s  (%5.1f%%)\n",
			profile.name, profile.count, display.FormatSize(profile.size), percent)
	}

	fmt.Printf("└────────────────────────────────────────────────────────────────────────────────┘\n")
}

// printRecentActivity shows recently created archives
func (dd *DashboardDisplay) printRecentActivity(stats Statistics) {
	if len(stats.RecentArchives) == 0 {
		return
	}

	fmt.Printf("\n┌─ RECENT ACTIVITY ──────────────────────────────────────────────────────────────┐\n")

	for _, archive := range stats.RecentArchives {
		age := dd.formatDashboardAge(archive.Created)
		size := display.FormatSize(archive.Size)
		name := archive.Name
		if len(name) > 35 {
			name = name[:34] + "…"
		}

		status := display.FormatStatus(archive.Status, true)

		// Align content with section header text
		fmt.Printf("│  %s %-35s  %8s  %s\n", status, name, size, age)
	}

	fmt.Printf("└────────────────────────────────────────────────────────────────────────────────┘\n")
}

// printHealthIndicators shows system health metrics
func (dd *DashboardDisplay) printHealthIndicators(stats Statistics) {
	fmt.Printf("\n┌─ HEALTH INDICATORS ────────────────────────────────────────────────────────────┐\n")

	// Overall health
	healthStatus := "Excellent"
	if stats.HealthScore < 90 {
		healthStatus = "Good"
	}
	if stats.HealthScore < 70 {
		healthStatus = "Fair"
	}
	if stats.HealthScore < 50 {
		healthStatus = "Poor"
	}

	fmt.Printf("│  Overall Health: %.1f%% (%s)\n", stats.HealthScore, healthStatus)

	// Size metrics
	if stats.LargestArchive != nil {
		avgSize := stats.TotalSize / int64(stats.Total)
		fmt.Printf("│  Average Size: %s    Largest: %s\n",
			display.FormatSize(avgSize), display.FormatSize(stats.LargestArchive.Size))
	}

	// Age metrics
	if stats.OldestArchive != nil && stats.NewestArchive != nil {
		oldestAge := dd.formatDashboardAge(stats.OldestArchive.Created)
		newestAge := dd.formatDashboardAge(stats.NewestArchive.Created)
		fmt.Printf("│  Archive Age Range: %s (oldest) to %s (newest)\n", oldestAge, newestAge)
	}

	fmt.Printf("└────────────────────────────────────────────────────────────────────────────────┘\n")
}

// formatDashboardAge formats duration for dashboard display
func (dd *DashboardDisplay) formatDashboardAge(created time.Time) string {
	age := time.Since(created)

	if age < time.Hour {
		return fmt.Sprintf("%dm ago", int(age.Minutes()))
	}
	if age < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(age.Hours()))
	}
	if age < 7*24*time.Hour {
		return fmt.Sprintf("%dd ago", int(age.Hours()/24))
	}
	if age < 30*24*time.Hour {
		return fmt.Sprintf("%dw ago", int(age.Hours()/(24*7)))
	}
	if age < 365*24*time.Hour {
		return fmt.Sprintf("%dmo ago", int(age.Hours()/(24*30)))
	}
	return fmt.Sprintf("%dy ago", int(age.Hours()/(24*365)))
}
