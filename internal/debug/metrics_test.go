package debug

import (
	"strings"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	metrics := NewMetrics()
	
	if metrics.StartTime.IsZero() {
		t.Error("Expected StartTime to be set")
	}
	if metrics.MemoryBefore == 0 {
		t.Error("Expected MemoryBefore to be set")
	}
}

func TestMetrics_RecordQueryTime(t *testing.T) {
	metrics := NewMetrics()
	
	// Simulate some work
	time.Sleep(10 * time.Millisecond)
	metrics.RecordQueryTime()
	
	if metrics.QueryTime <= 0 {
		t.Error("Expected QueryTime to be greater than 0")
	}
	if metrics.QueryTime < 10*time.Millisecond {
		t.Error("Expected QueryTime to be at least 10ms")
	}
}

func TestMetrics_RecordRenderTime(t *testing.T) {
	metrics := NewMetrics()
	
	// Record query time first
	time.Sleep(5 * time.Millisecond)
	metrics.RecordQueryTime()
	
	// Simulate render work
	time.Sleep(10 * time.Millisecond)
	metrics.RecordRenderTime()
	
	if metrics.RenderTime <= 0 {
		t.Error("Expected RenderTime to be greater than 0")
	}
}

func TestMetrics_SetResultCount(t *testing.T) {
	metrics := NewMetrics()
	
	metrics.SetResultCount(42)
	
	if metrics.ResultCount != 42 {
		t.Errorf("Expected ResultCount 42, got %d", metrics.ResultCount)
	}
}

func TestMetrics_SetDatabaseSize(t *testing.T) {
	metrics := NewMetrics()
	
	metrics.SetDatabaseSize(1024)
	
	if metrics.DatabaseSize != 1024 {
		t.Errorf("Expected DatabaseSize 1024, got %d", metrics.DatabaseSize)
	}
}

func TestMetrics_String(t *testing.T) {
	metrics := NewMetrics()
	metrics.SetResultCount(5)
	metrics.SetDatabaseSize(2048)
	
	// Record some times
	time.Sleep(5 * time.Millisecond)
	metrics.RecordQueryTime()
	time.Sleep(3 * time.Millisecond)
	metrics.RecordRenderTime()
	
	output := metrics.String()
	
	// Check that output contains expected parts
	if !strings.Contains(output, "DEBUG:") {
		t.Error("Expected output to contain 'DEBUG:'")
	}
	if !strings.Contains(output, "Query:") {
		t.Error("Expected output to contain 'Query:'")
	}
	if !strings.Contains(output, "Render:") {
		t.Error("Expected output to contain 'Render:'")
	}
	if !strings.Contains(output, "Results: 5") {
		t.Error("Expected output to contain 'Results: 5'")
	}
	if !strings.Contains(output, "DB: 2.0 kB") {
		t.Error("Expected output to contain 'DB: 2.0 kB'")
	}
}

func TestMetrics_Summary(t *testing.T) {
	metrics := NewMetrics()
	metrics.SetResultCount(10)
	metrics.SetDatabaseSize(4096)
	
	// Record some times
	time.Sleep(2 * time.Millisecond)
	metrics.RecordQueryTime()
	time.Sleep(3 * time.Millisecond)
	metrics.RecordRenderTime()
	
	output := metrics.Summary()
	
	// Check that output contains expected sections
	if !strings.Contains(output, "Performance Metrics:") {
		t.Error("Expected output to contain 'Performance Metrics:'")
	}
	if !strings.Contains(output, "Query Time:") {
		t.Error("Expected output to contain 'Query Time:'")
	}
	if !strings.Contains(output, "Render Time:") {
		t.Error("Expected output to contain 'Render Time:'")
	}
	if !strings.Contains(output, "Total Time:") {
		t.Error("Expected output to contain 'Total Time:'")
	}
	if !strings.Contains(output, "Result Count:  10") {
		t.Error("Expected output to contain 'Result Count:  10'")
	}
	if !strings.Contains(output, "Database Size: 4.1 kB") {
		t.Error("Expected output to contain 'Database Size: 4.1 kB'")
	}
}

func TestMetrics_Finish(t *testing.T) {
	metrics := NewMetrics()
	
	// MemoryAfter should be 0 initially
	if metrics.MemoryAfter != 0 {
		t.Error("Expected MemoryAfter to be 0 before Finish()")
	}
	
	metrics.Finish()
	
	// MemoryAfter should be set after Finish()
	if metrics.MemoryAfter == 0 {
		t.Error("Expected MemoryAfter to be set after Finish()")
	}
}