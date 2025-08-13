package debug

import (
	"fmt"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
)

// Metrics holds performance information for debug output
type Metrics struct {
	StartTime    time.Time
	QueryTime    time.Duration
	RenderTime   time.Duration
	ResultCount  int
	DatabaseSize int64
	MemoryBefore uint64
	MemoryAfter  uint64
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return &Metrics{
		StartTime:    time.Now(),
		MemoryBefore: m.Alloc,
	}
}

// RecordQueryTime marks the query completion
func (m *Metrics) RecordQueryTime() {
	m.QueryTime = time.Since(m.StartTime)
}

// RecordRenderTime marks the render completion
func (m *Metrics) RecordRenderTime() {
	m.RenderTime = time.Since(m.StartTime) - m.QueryTime
}

// SetResultCount sets the number of results
func (m *Metrics) SetResultCount(count int) {
	m.ResultCount = count
}

// SetDatabaseSize sets the database size
func (m *Metrics) SetDatabaseSize(size int64) {
	m.DatabaseSize = size
}

// Finish collects final memory stats
func (m *Metrics) Finish() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.MemoryAfter = ms.Alloc
}

// String formats metrics for output
func (m *Metrics) String() string {
	m.Finish()
	memUsed := uint64(0)
	if m.MemoryAfter > m.MemoryBefore {
		memUsed = m.MemoryAfter - m.MemoryBefore
	}
	
	return fmt.Sprintf("DEBUG: Query: %v, Render: %v, Results: %d, DB: %s, Memory: %s",
		m.QueryTime.Round(time.Millisecond),
		m.RenderTime.Round(time.Millisecond),
		m.ResultCount,
		humanize.Bytes(uint64(m.DatabaseSize)),
		humanize.Bytes(memUsed))
}

// Summary provides a multi-line debug summary
func (m *Metrics) Summary() string {
	m.Finish()
	memUsed := uint64(0)
	if m.MemoryAfter > m.MemoryBefore {
		memUsed = m.MemoryAfter - m.MemoryBefore
	}
	
	total := time.Since(m.StartTime)
	
	return fmt.Sprintf(`
Performance Metrics:
  Query Time:    %v
  Render Time:   %v
  Total Time:    %v
  Result Count:  %d
  Database Size: %s
  Memory Used:   %s`,
		m.QueryTime.Round(time.Millisecond),
		m.RenderTime.Round(time.Millisecond),
		total.Round(time.Millisecond),
		m.ResultCount,
		humanize.Bytes(uint64(m.DatabaseSize)),
		humanize.Bytes(memUsed))
}