package middlewares

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// PerformanceMetrics tracks API performance metrics
type PerformanceMetrics struct {
	mu sync.RWMutex

	// Request counts by endpoint
	RequestCounts map[string]int64

	// Response times by endpoint (in nanoseconds)
	ResponseTimes map[string][]time.Duration

	// Error counts by endpoint
	ErrorCounts map[string]int64

	// Last request time
	LastRequestTime time.Time
}

// Global metrics instance
var GlobalMetrics = &PerformanceMetrics{
	RequestCounts:   make(map[string]int64),
	ResponseTimes:   make(map[string][]time.Duration),
	ErrorCounts:     make(map[string]int64),
	LastRequestTime: time.Now(),
}

// PerformanceMonitoring middleware tracks API performance metrics
func PerformanceMonitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(rw, r)

		duration := time.Since(startTime)
		endpoint := r.Method + ":" + r.URL.Path

		// Update metrics
		GlobalMetrics.mu.Lock()
		GlobalMetrics.RequestCounts[endpoint]++
		GlobalMetrics.ResponseTimes[endpoint] = append(GlobalMetrics.ResponseTimes[endpoint], duration)
		GlobalMetrics.LastRequestTime = time.Now()

		if rw.statusCode >= 400 {
			GlobalMetrics.ErrorCounts[endpoint]++
		}
		GlobalMetrics.mu.Unlock()

		// Log slow requests (> 50ms)
		if duration > 50*time.Millisecond {
			log.Printf("⚠️  SLOW REQUEST: %s %s took %v (status: %d)", r.Method, r.URL.Path, duration, rw.statusCode)
		}
	})
}

// GetMetrics returns current performance metrics
func GetMetrics() *PerformanceMetrics {
	GlobalMetrics.mu.RLock()
	defer GlobalMetrics.mu.RUnlock()

	// Create a copy to avoid race conditions
	metrics := &PerformanceMetrics{
		RequestCounts:   make(map[string]int64),
		ResponseTimes:   make(map[string][]time.Duration),
		ErrorCounts:     make(map[string]int64),
		LastRequestTime: GlobalMetrics.LastRequestTime,
	}

	for k, v := range GlobalMetrics.RequestCounts {
		metrics.RequestCounts[k] = v
	}

	for k, v := range GlobalMetrics.ResponseTimes {
		metrics.ResponseTimes[k] = make([]time.Duration, len(v))
		copy(metrics.ResponseTimes[k], v)
	}

	for k, v := range GlobalMetrics.ErrorCounts {
		metrics.ErrorCounts[k] = v
	}

	return metrics
}

// CalculateP95 calculates the 95th percentile response time for an endpoint
func (pm *PerformanceMetrics) CalculateP95(endpoint string) time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	times, exists := pm.ResponseTimes[endpoint]
	if !exists || len(times) == 0 {
		return 0
	}

	// Sort the response times
	sorted := make([]time.Duration, len(times))
	copy(sorted, times)

	// Simple bubble sort (for small datasets)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Calculate 95th percentile
	index := int(float64(len(sorted)) * 0.95)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// CalculateAverage calculates the average response time for an endpoint
func (pm *PerformanceMetrics) CalculateAverage(endpoint string) time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	times, exists := pm.ResponseTimes[endpoint]
	if !exists || len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	return total / time.Duration(len(times))
}

// CalculateMin calculates the minimum response time for an endpoint
func (pm *PerformanceMetrics) CalculateMin(endpoint string) time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	times, exists := pm.ResponseTimes[endpoint]
	if !exists || len(times) == 0 {
		return 0
	}

	min := times[0]
	for _, t := range times {
		if t < min {
			min = t
		}
	}

	return min
}

// CalculateMax calculates the maximum response time for an endpoint
func (pm *PerformanceMetrics) CalculateMax(endpoint string) time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	times, exists := pm.ResponseTimes[endpoint]
	if !exists || len(times) == 0 {
		return 0
	}

	max := times[0]
	for _, t := range times {
		if t > max {
			max = t
		}
	}

	return max
}

// ResetMetrics clears all performance metrics
func ResetMetrics() {
	GlobalMetrics.mu.Lock()
	defer GlobalMetrics.mu.Unlock()

	GlobalMetrics.RequestCounts = make(map[string]int64)
	GlobalMetrics.ResponseTimes = make(map[string][]time.Duration)
	GlobalMetrics.ErrorCounts = make(map[string]int64)
	GlobalMetrics.LastRequestTime = time.Now()
}
