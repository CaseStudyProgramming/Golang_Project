package metrics

import (
	"testing"
	"time"
)

func TestNewPerformanceMonitor(t *testing.T) {
	pm := NewPerformanceMonitor()

	if pm == nil {
		t.Fatal("NewPerformanceMonitor returned nil")
	}

	if pm.metrics == nil {
		t.Fatal("PerformanceMonitor metrics is nil")
	}

	if pm.metrics.RequestCount != 0 {
		t.Errorf("Expected initial request count 0, got %d", pm.metrics.RequestCount)
	}

	if pm.metrics.ErrorCount != 0 {
		t.Errorf("Expected initial error count 0, got %d", pm.metrics.ErrorCount)
	}
}

func TestPerformanceMonitor_RecordAPICall(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Record an API call
	pm.RecordAPICall("GET /test", 100*time.Millisecond)

	if pm.metrics.RequestCount != 1 {
		t.Errorf("Expected request count 1, got %d", pm.metrics.RequestCount)
	}

	if pm.metrics.ErrorCount != 0 {
		t.Errorf("Expected error count 0, got %d", pm.metrics.ErrorCount)
	}

	if pm.metrics.APICallCount["GET /test"] != 1 {
		t.Errorf("Expected API call count 1, got %d", pm.metrics.APICallCount["GET /test"])
	}
}

func TestPerformanceMonitor_RecordError(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Record an error
	pm.RecordError()

	if pm.metrics.ErrorCount != 1 {
		t.Errorf("Expected error count 1, got %d", pm.metrics.ErrorCount)
	}
}

func TestPerformanceMonitor_GetMetrics(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Record some API calls
	pm.RecordAPICall("GET /test1", 100*time.Millisecond)
	pm.RecordAPICall("POST /test2", 150*time.Millisecond)
	pm.RecordError()

	metrics := pm.GetMetrics()

	if metrics == nil {
		t.Fatal("GetMetrics returned nil")
	}

	if pm.metrics.RequestCount != 2 {
		t.Errorf("Expected request count 2, got %d", pm.metrics.RequestCount)
	}

	if pm.metrics.ErrorCount != 1 {
		t.Errorf("Expected error count 1, got %d", pm.metrics.ErrorCount)
	}
}

func TestPerformanceMonitor_ResetMetrics(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Record some data
	pm.RecordAPICall("GET /test", 100*time.Millisecond)
	pm.RecordError()

	// Reset
	pm.ResetMetrics()

	if pm.metrics.RequestCount != 0 {
		t.Errorf("Expected request count 0 after reset, got %d", pm.metrics.RequestCount)
	}

	if pm.metrics.ErrorCount != 0 {
		t.Errorf("Expected error count 0 after reset, got %d", pm.metrics.ErrorCount)
	}

	if len(pm.metrics.APICallCount) != 0 {
		t.Errorf("Expected 0 API call counts after reset, got %d", len(pm.metrics.APICallCount))
	}
}

func TestPerformanceMonitor_UpdateSystemMetrics(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Update system metrics
	pm.UpdateSystemMetrics()

	if pm.metrics.GoroutineCount == 0 {
		t.Error("Expected positive goroutine count")
	}

	if pm.metrics.MemoryStats.Alloc == 0 {
		t.Error("Expected positive memory allocation")
	}
}

func TestPerformanceMonitor_Percentile(t *testing.T) {
	pm := NewPerformanceMonitor()

	// Record some durations
	durations := []time.Duration{
		100 * time.Millisecond,
		150 * time.Millisecond,
		200 * time.Millisecond,
		250 * time.Millisecond,
		300 * time.Millisecond,
	}

	// Test percentile calculation
	p50 := pm.percentile(durations, 50)
	p95 := pm.percentile(durations, 95)
	p99 := pm.percentile(durations, 99)

	if p50 <= 0 {
		t.Errorf("Expected positive P50, got %f", p50)
	}

	if p95 <= 0 {
		t.Errorf("Expected positive P95, got %f", p95)
	}

	if p99 <= 0 {
		t.Errorf("Expected positive P99, got %f", p99)
	}
}

func TestPerformanceMonitor_InitPerformanceMonitor(t *testing.T) {
	// Initialize global monitor
	InitPerformanceMonitor()

	if GlobalMonitor == nil {
		t.Fatal("GlobalMonitor is nil after initialization")
	}
}
