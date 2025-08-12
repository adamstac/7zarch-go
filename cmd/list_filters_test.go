package cmd

import (
	"github.com/adamstac/7zarch-go/internal/storage"
	"testing"
)

func TestApplyFilters_Status(t *testing.T) {
	archives := []*storage.Archive{{Status: "present"}, {Status: "missing"}, {Status: "deleted"}}
	filtered := applyFilters(archives, struct {
		status, profile string
		largerThan      int64
	}{"missing", "", 0})
	if len(filtered) != 1 || filtered[0].Status != "missing" {
		t.Fatalf("expected 1 missing, got %d", len(filtered))
	}
}

func TestApplyFilters_Profile(t *testing.T) {
	archives := []*storage.Archive{{Profile: "A"}, {Profile: "B"}, {Profile: "A"}}
	filtered := applyFilters(archives, struct {
		status, profile string
		largerThan      int64
	}{"", "A", 0})
	if len(filtered) != 2 {
		t.Fatalf("expected 2 with profile A, got %d", len(filtered))
	}
}

func TestApplyFilters_LargerThan(t *testing.T) {
	archives := []*storage.Archive{{Size: 10}, {Size: 20}, {Size: 5}}
	filtered := applyFilters(archives, struct {
		status, profile string
		largerThan      int64
	}{"", "", 10})
	if len(filtered) != 1 || filtered[0].Size != 20 {
		t.Fatalf("expected only size 20, got %d entries", len(filtered))
	}
}
