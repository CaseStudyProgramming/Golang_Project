package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// PerformanceMetrics holds performance monitoring data
type PerformanceMetrics struct {
	mu sync.RWMutex

	// API metrics
	APICallCount    map[string]int64
	APIDuration    map[string][]time.Duration
	APIDurationSum map[string]time.Duration

	// System metrics
	MemoryStats     runtime.MemStats
	GoroutineCount  int
	RequestCount    int64
	ErrorCount      int64

	// Timing
	StartTime       time.Time
	LastUpdated     time.Time
}

// PerformanceMonitor manages performance monitoring
type PerformanceMonitor struct {
	metrics *PerformanceMetrics
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics: &PerformanceMetrics{
			APICallCount:    make(map[string]int64),
			APIDuration:    make(map[string][]time.Duration),
			APIDurationSum: make(map[string]time.Duration),
			StartTime:      time.Now(),
			LastUpdated:    time.Now(),
		},
	}
}

// RecordAPICall records an API call with its duration
func (pm *PerformanceMonitor) RecordAPICall(endpoint string, duration time.Duration) {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.APICallCount[endpoint]++
	pm.metrics.APIDuration[endpoint] = append(pm.metrics.APIDuration[endpoint], duration)
	pm.metrics.APIDurationSum[endpoint] += duration
	pm.metrics.RequestCount++
	pm.metrics.LastUpdated = time.Now()

	// Keep only last 1000 durations per endpoint to prevent memory bloat
	if len(pm.metrics.APIDuration[endpoint]) > 1000 {
		pm.metrics.APIDuration[endpoint] = pm.metrics.APIDuration[endpoint][1:]
	}
}

// RecordError records an error occurrence
func (pm *PerformanceMonitor) RecordError() {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.ErrorCount++
	pm.metrics.LastUpdated = time.Now()
}

// UpdateSystemMetrics updates system-level metrics
func (pm *PerformanceMonitor) UpdateSystemMetrics() {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	runtime.ReadMemStats(&pm.metrics.MemoryStats)
	pm.metrics.GoroutineCount = runtime.NumGoroutine()
	pm.metrics.LastUpdated = time.Now()
}

// GetMetrics returns current performance metrics
func (pm *PerformanceMonitor) GetMetrics() map[string]interface{} {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()

	// Calculate average durations
	avgDurations := make(map[string]float64)
	for endpoint, sum := range pm.metrics.APIDurationSum {
		count := pm.metrics.APICallCount[endpoint]
		if count > 0 {
			avgDurations[endpoint] = float64(sum.Milliseconds()) / float64(count)
		}
	}

	// Calculate percentiles
	p50Durations := make(map[string]float64)
	p95Durations := make(map[string]float64)
	p99Durations := make(map[string]float64)

	for endpoint, durations := range pm.metrics.APIDuration {
		if len(durations) > 0 {
			p50Durations[endpoint] = pm.percentile(durations, 50)
			p95Durations[endpoint] = pm.percentile(durations, 95)
			p99Durations[endpoint] = pm.percentile(durations, 99)
		}
	}

	return map[string]interface{}{
		"api_metrics": map[string]interface{}{
			"call_count":     pm.metrics.APICallCount,
			"avg_duration":   avgDurations,
			"p50_duration":   p50Durations,
			"p95_duration":   p95Durations,
			"p99_duration":   p99Durations,
			"total_requests": pm.metrics.RequestCount,
		},
		"system_metrics": map[string]interface{}{
			"goroutines":       pm.metrics.GoroutineCount,
			"memory_alloc":     pm.metrics.MemoryStats.Alloc,
			"memory_total_alloc": pm.metrics.MemoryStats.TotalAlloc,
			"memory_sys":       pm.metrics.MemoryStats.Sys,
			"memory_heap_alloc": pm.metrics.MemoryStats.HeapAlloc,
			"memory_heap_sys":  pm.metrics.MemoryStats.HeapSys,
			"gc_pause_total":   pm.metrics.MemoryStats.PauseTotalNs,
			"gc_pause_count":   pm.metrics.MemoryStats.NumGC,
		},
		"error_metrics": map[string]interface{}{
			"error_count": pm.metrics.ErrorCount,
			"error_rate":  float64(pm.metrics.ErrorCount) / float64(pm.metrics.RequestCount) * 100,
		},
		"timing": map[string]interface{}{
			"start_time":  pm.metrics.StartTime,
			"last_updated": pm.metrics.LastUpdated,
			"uptime":      time.Since(pm.metrics.StartTime).String(),
		},
	}
}

// percentile calculates the percentile value from a slice of durations
func (pm *PerformanceMonitor) percentile(durations []time.Duration, p int) float64 {
	if len(durations) == 0 {
		return 0
	}

	// Make a copy to avoid modifying original
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

	// Simple sort (for production, use more efficient algorithm)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := (len(sorted) * p) / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return float64(sorted[index].Milliseconds())
}

// ResetMetrics resets all metrics
func (pm *PerformanceMonitor) ResetMetrics() {
	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.APICallCount = make(map[string]int64)
	pm.metrics.APIDuration = make(map[string][]time.Duration)
	pm.metrics.APIDurationSum = make(map[string]time.Duration)
	pm.metrics.RequestCount = 0
	pm.metrics.ErrorCount = 0
	pm.metrics.StartTime = time.Now()
	pm.metrics.LastUpdated = time.Now()
}

// MetricsHandler returns an HTTP handler for metrics endpoint
func (pm *PerformanceMonitor) MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pm.UpdateSystemMetrics()
		
		metrics := pm.GetMetrics()
		
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(metrics); err != nil {
			http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
			return
		}
	}
}

// HealthHandler returns an HTTP handler for health check endpoint
func (pm *PerformanceMonitor) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pm.UpdateSystemMetrics()
		
		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now(),
			"uptime":    time.Since(pm.metrics.StartTime).String(),
			"metrics": map[string]interface{}{
				"goroutines": pm.metrics.GoroutineCount,
				"memory":    pm.metrics.MemoryStats.Alloc,
				"requests":  pm.metrics.RequestCount,
				"errors":    pm.metrics.ErrorCount,
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, "Failed to encode health status", http.StatusInternalServerError)
			return
		}
	}
}

// Middleware wraps an HTTP handler with performance monitoring
func (pm *PerformanceMonitor) Middleware(endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next(rw, r)
		
		duration := time.Since(start)
		pm.RecordAPICall(endpoint, duration)
		
		if rw.statusCode >= 400 {
			pm.RecordError()
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Global performance monitor instance
var GlobalMonitor *PerformanceMonitor

// InitPerformanceMonitor initializes the global performance monitor
func InitPerformanceMonitor() {
	GlobalMonitor = NewPerformanceMonitor()
}
