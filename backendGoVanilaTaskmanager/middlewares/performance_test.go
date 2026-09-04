package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPerformanceMonitoring(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make a request
	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Check metrics
	metrics := GetMetrics()
	endpoint := "GET:/api/test"

	if metrics.RequestCounts[endpoint] != 1 {
		t.Errorf("Expected 1 request count, got %d", metrics.RequestCounts[endpoint])
	}

	if len(metrics.ResponseTimes[endpoint]) != 1 {
		t.Errorf("Expected 1 response time, got %d", len(metrics.ResponseTimes[endpoint]))
	}

	if metrics.ErrorCounts[endpoint] != 0 {
		t.Errorf("Expected 0 error count, got %d", metrics.ErrorCounts[endpoint])
	}
}

func TestPerformanceMonitoring_Error(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make a request
	req := httptest.NewRequest("GET", "/api/error", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check metrics
	metrics := GetMetrics()
	endpoint := "GET:/api/error"

	if metrics.ErrorCounts[endpoint] != 1 {
		t.Errorf("Expected 1 error count, got %d", metrics.ErrorCounts[endpoint])
	}
}

func TestPerformanceMonitoring_MultipleRequests(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make multiple requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/multiple", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// Check metrics
	metrics := GetMetrics()
	endpoint := "GET:/api/multiple"

	if metrics.RequestCounts[endpoint] != 5 {
		t.Errorf("Expected 5 request count, got %d", metrics.RequestCounts[endpoint])
	}

	if len(metrics.ResponseTimes[endpoint]) != 5 {
		t.Errorf("Expected 5 response times, got %d", len(metrics.ResponseTimes[endpoint]))
	}
}

func TestCalculateP95(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate varying response times
		time.Sleep(time.Duration(10) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make multiple requests
	for i := 0; i < 20; i++ {
		req := httptest.NewRequest("GET", "/api/p95", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// Calculate P95
	p95 := GlobalMetrics.CalculateP95("GET:/api/p95")

	if p95 == 0 {
		t.Error("Expected non-zero P95 value")
	}

	t.Logf("P95 response time: %v", p95)
}

func TestCalculateAverage(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond) // Add minimal delay
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make multiple requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/api/average", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// Calculate average
	avg := GlobalMetrics.CalculateAverage("GET:/api/average")

	if avg == 0 {
		t.Error("Expected non-zero average value")
	}

	t.Logf("Average response time: %v", avg)
}

func TestCalculateMinMax(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond) // Add minimal delay
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make multiple requests
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/api/minmax", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// Calculate min and max
	min := GlobalMetrics.CalculateMin("GET:/api/minmax")
	max := GlobalMetrics.CalculateMax("GET:/api/minmax")

	if min == 0 {
		t.Error("Expected non-zero min value")
	}

	if max == 0 {
		t.Error("Expected non-zero max value")
	}

	if min > max {
		t.Error("Expected min to be less than max")
	}

	t.Logf("Min response time: %v, Max response time: %v", min, max)
}

func TestResetMetrics(t *testing.T) {
	// Add some metrics
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	req := httptest.NewRequest("GET", "/api/reset", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Verify metrics exist
	metrics := GetMetrics()
	if metrics.RequestCounts["GET:/api/reset"] == 0 {
		t.Error("Expected metrics to exist before reset")
	}

	// Reset metrics
	ResetMetrics()

	// Verify metrics are cleared
	metrics = GetMetrics()
	if metrics.RequestCounts["GET:/api/reset"] != 0 {
		t.Error("Expected metrics to be cleared after reset")
	}

	if len(metrics.ResponseTimes) != 0 {
		t.Error("Expected response times to be cleared after reset")
	}
}

func TestGetMetrics_CopiesData(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make a request
	req := httptest.NewRequest("GET", "/api/copy", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get metrics copy
	metrics1 := GetMetrics()

	// Make another request
	req = httptest.NewRequest("GET", "/api/copy", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Get another metrics copy
	metrics2 := GetMetrics()

	// Verify the first copy didn't change
	if metrics1.RequestCounts["GET:/api/copy"] != 1 {
		t.Errorf("Expected first copy to have 1 request, got %d", metrics1.RequestCounts["GET:/api/copy"])
	}

	// Verify the second copy has the updated count
	if metrics2.RequestCounts["GET:/api/copy"] != 2 {
		t.Errorf("Expected second copy to have 2 requests, got %d", metrics2.RequestCounts["GET:/api/copy"])
	}
}

func TestPerformanceMonitoring_DifferentEndpoints(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := PerformanceMonitoring(nextHandler)

	// Make requests to different endpoints
	endpoints := []string{"/api/endpoint1", "/api/endpoint2", "/api/endpoint3"}
	for _, endpoint := range endpoints {
		req := httptest.NewRequest("GET", endpoint, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// Check metrics
	metrics := GetMetrics()

	for _, endpoint := range endpoints {
		fullEndpoint := "GET:" + endpoint
		if metrics.RequestCounts[fullEndpoint] != 1 {
			t.Errorf("Expected 1 request for %s, got %d", fullEndpoint, metrics.RequestCounts[fullEndpoint])
		}
	}
}
