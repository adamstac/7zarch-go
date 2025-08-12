package archive

import "testing"

func Test_recommendProfileWithThresholds(t *testing.T) {
	stats := &ContentStats{TotalBytes: 100}

	// Media dominant
	stats.MediaBytes, stats.DocumentBytes, stats.CompressedBytes = 80, 10, 10
	if got := recommendProfileWithThresholds(stats, 70, 60); got.Name != "Media" {
		t.Errorf("media dominant: want Media, got %s", got.Name)
	}

	// Documents dominant
	stats.MediaBytes, stats.DocumentBytes, stats.CompressedBytes = 10, 70, 20
	if got := recommendProfileWithThresholds(stats, 70, 60); got.Name != "Documents" {
		t.Errorf("docs dominant: want Documents, got %s", got.Name)
	}

	// Compressed dominant
	stats.MediaBytes, stats.DocumentBytes, stats.CompressedBytes = 10, 30, 60
	if got := recommendProfileWithThresholds(stats, 70, 60); got.Name != "Media" {
		t.Errorf("compressed dominant: want Media (fast), got %s", got.Name)
	}

	// Balanced default
	stats.MediaBytes, stats.DocumentBytes, stats.CompressedBytes = 30, 30, 40
	if got := recommendProfileWithThresholds(stats, 70, 60); got.Name != "Balanced" {
		t.Errorf("balanced: want Balanced, got %s", got.Name)
	}
}
