package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"taskmanager/middlewares"
)

type MetricsController struct{}

func NewMetricsController() *MetricsController {
	return &MetricsController{}
}

// MetricStats represents statistics for a single endpoint
type MetricStats struct {
	Endpoint     string        `json:"endpoint"`
	RequestCount int64         `json:"request_count"`
	ErrorCount   int64         `json:"error_count"`
	AverageTime  time.Duration `json:"average_time"`
	P95Time      time.Duration `json:"p95_time"`
	MinTime      time.Duration `json:"min_time"`
	MaxTime      time.Duration `json:"max_time"`
}

// MetricsResponse represents the complete metrics response
type MetricsResponse struct {
	TotalRequests    int64            `json:"total_requests"`
	TotalErrors      int64            `json:"total_errors"`
	LastRequestTime  time.Time        `json:"last_request_time"`
	EndpointMetrics  []MetricStats    `json:"endpoint_metrics"`
	SlowEndpoints    []MetricStats    `json:"slow_endpoints"` // Endpoints with P95 > 50ms
}

// GetMetricsHandler returns performance metrics
func (c *MetricsController) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := middlewares.GetMetrics()

	response := &MetricsResponse{
		LastRequestTime: metrics.LastRequestTime,
	}

	// Calculate totals
	var totalRequests int64
	var totalErrors int64

	// Calculate stats for each endpoint
	for endpoint := range metrics.RequestCounts {
		requestCount := metrics.RequestCounts[endpoint]
		errorCount := metrics.ErrorCounts[endpoint]
		averageTime := metrics.CalculateAverage(endpoint)
		p95Time := metrics.CalculateP95(endpoint)
		minTime := metrics.CalculateMin(endpoint)
		maxTime := metrics.CalculateMax(endpoint)

		stats := MetricStats{
			Endpoint:     endpoint,
			RequestCount: requestCount,
			ErrorCount:   errorCount,
			AverageTime:  averageTime,
			P95Time:      p95Time,
			MinTime:      minTime,
			MaxTime:      maxTime,
		}

		response.EndpointMetrics = append(response.EndpointMetrics, stats)

		totalRequests += requestCount
		totalErrors += errorCount

		// Track slow endpoints (P95 > 50ms)
		if p95Time > 50*time.Millisecond {
			response.SlowEndpoints = append(response.SlowEndpoints, stats)
		}
	}

	response.TotalRequests = totalRequests
	response.TotalErrors = totalErrors

	// Sort endpoints by request count (descending)
	sort.Slice(response.EndpointMetrics, func(i, j int) bool {
		return response.EndpointMetrics[i].RequestCount > response.EndpointMetrics[j].RequestCount
	})

	// Sort slow endpoints by P95 time (descending)
	sort.Slice(response.SlowEndpoints, func(i, j int) bool {
		return response.SlowEndpoints[i].P95Time > response.SlowEndpoints[j].P95Time
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResetMetricsHandler resets all performance metrics
func (c *MetricsController) ResetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	middlewares.ResetMetrics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Metrics reset successfully",
	})
}