package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"taskmanager/middlewares"
)

func TestNewMetricsController(t *testing.T) {
	controller := NewMetricsController()
	if controller == nil {
		t.Fatal("NewMetricsController returned nil")
	}
}

func TestMetricsController_GetMetricsHandler(t *testing.T) {
	// Reset metrics before test
	middlewares.ResetMetrics()

	// Record some sample metrics
	middlewares.RecordAPICall("/api/tasks", 100*time.Millisecond, false)
	middlewares.RecordAPICall("/api/tasks", 150*time.Millisecond, false)
	middlewares.RecordAPICall("/api/tasks", 200*time.Millisecond, true) // error
	middlewares.RecordAPICall("/api/users", 50*time.Millisecond, false)
	middlewares.RecordAPICall("/api/users", 75*time.Millisecond, false)

	controller := NewMetricsController()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	controller.GetMetricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MetricsResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify totals
	if response.TotalRequests != 5 {
		t.Errorf("Expected total requests 5, got %d", response.TotalRequests)
	}

	if response.TotalErrors != 1 {
		t.Errorf("Expected total errors 1, got %d", response.TotalErrors)
	}

	// Verify endpoint metrics
	if len(response.EndpointMetrics) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(response.EndpointMetrics))
	}

	// Verify sorting (should be sorted by request count descending)
	if len(response.EndpointMetrics) > 1 {
		if response.EndpointMetrics[0].RequestCount < response.EndpointMetrics[1].RequestCount {
			t.Error("Endpoint metrics not sorted by request count descending")
		}
	}

	// Verify slow endpoints tracking
	// P95 for /api/tasks should be > 50ms (max is 200ms)
	hasSlowEndpoints := false
	for _, endpoint := range response.EndpointMetrics {
		if endpoint.P95Time > 50*time.Millisecond {
			hasSlowEndpoints = true
			break
		}
	}
	if !hasSlowEndpoints {
		t.Error("Expected slow endpoints to be tracked")
	}
}

func TestMetricsController_GetMetricsHandler_EmptyMetrics(t *testing.T) {
	// Reset metrics to ensure empty state
	middlewares.ResetMetrics()

	controller := NewMetricsController()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	controller.GetMetricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MetricsResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify empty state
	if response.TotalRequests != 0 {
		t.Errorf("Expected total requests 0, got %d", response.TotalRequests)
	}

	if response.TotalErrors != 0 {
		t.Errorf("Expected total errors 0, got %d", response.TotalErrors)
	}

	if len(response.EndpointMetrics) != 0 {
		t.Errorf("Expected 0 endpoints, got %d", len(response.EndpointMetrics))
	}
}

func TestMetricsController_ResetMetricsHandler(t *testing.T) {
	// Record some metrics first
	middlewares.RecordAPICall("/api/tasks", 100*time.Millisecond, false)
	middlewares.RecordAPICall("/api/users", 50*time.Millisecond, true)

	controller := NewMetricsController()
	req := httptest.NewRequest("POST", "/metrics/reset", nil)
	w := httptest.NewRecorder()

	controller.ResetMetricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "Metrics reset successfully" {
		t.Errorf("Expected message 'Metrics reset successfully', got '%s'", response["message"])
	}

	// Verify metrics are actually reset
	metrics := middlewares.GetMetrics()
	if len(metrics.RequestCounts) != 0 {
		t.Error("Metrics were not reset - RequestCounts not empty")
	}

	if len(metrics.ErrorCounts) != 0 {
		t.Error("Metrics were not reset - ErrorCounts not empty")
	}
}

func TestMetricsController_GetMetricsHandler_SlowEndpointsSorting(t *testing.T) {
	// Reset metrics
	middlewares.ResetMetrics()

	// Record metrics with varying P95 times
	middlewares.RecordAPICall("/api/slow1", 300*time.Millisecond, false)
	middlewares.RecordAPICall("/api/slow1", 350*time.Millisecond, false)
	middlewares.RecordAPICall("/api/slow2", 150*time.Millisecond, false)
	middlewares.RecordAPICall("/api/slow2", 180*time.Millisecond, false)
	middlewares.RecordAPICall("/api/fast", 30*time.Millisecond, false)

	controller := NewMetricsController()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	controller.GetMetricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response MetricsResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify slow endpoints are sorted by P95 time descending
	if len(response.SlowEndpoints) > 1 {
		for i := 0; i < len(response.SlowEndpoints)-1; i++ {
			if response.SlowEndpoints[i].P95Time < response.SlowEndpoints[i+1].P95Time {
				t.Error("Slow endpoints not sorted by P95 time descending")
			}
		}
	}
}
