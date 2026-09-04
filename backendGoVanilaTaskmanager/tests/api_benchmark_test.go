package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"taskmanager/controllers"
	"taskmanager/middlewares"
	"taskmanager/models"
	"taskmanager/routes"
	"taskmanager/services"
	"taskmanager/utils"
)

// TestAPIPerformance_BenchmarkEndpoints benchmarks API endpoints to ensure P95 < 50ms
func TestAPIPerformance_BenchmarkEndpoints(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Setup test environment
	jwtManager := utils.NewJWTManager("test-secret-key")
	userModel := models.NewUserModel(db)
	userService := services.NewUserService(userModel, jwtManager)
	authController := controllers.NewAuthController(userService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtManager)

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	activityLogModel := models.NewActivityLogModel(db)
	activityLogService := services.NewActivityLogService(activityLogModel)
	taskService := services.NewTaskService(taskModel, tagModel, activityLogService)
	taskController := controllers.NewTaskController(taskService)

	categoryModel := models.NewCategoryModel(db)
	categoryService := services.NewCategoryService(categoryModel)
	categoryController := controllers.NewCategoryController(categoryService)

	tagService := services.NewTagService(tagModel)
	tagController := controllers.NewTagController(tagService)

	subtaskModel := models.NewSubtaskModel(db)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)
	subtaskController := controllers.NewSubtaskController(subtaskService)

	activityLogController := controllers.NewActivityLogController(activityLogService)

	metricsController := controllers.NewMetricsController()
	swaggerController := controllers.NewSwaggerController()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Apply middlewares
	handler := middlewares.Recovery(middlewares.Logger(middlewares.PerformanceMonitoring(middlewares.CORS([]string{"*"})(mux))))

	// Reset metrics before benchmarking
	middlewares.ResetMetrics()

	// Create test user and get token
	registerReq := map[string]interface{}{
		"name":     "benchmarkuser",
		"email":    "benchmark@example.com",
		"password": "benchmarkpass123",
		"timezone": "UTC",
	}
	registerBody, _ := json.Marshal(registerReq)

	registerReqHTTP := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()

	handler.ServeHTTP(registerResp, registerReqHTTP)

	if registerResp.Code != http.StatusCreated {
		t.Fatalf("Failed to register user: %s", registerResp.Body.String())
	}

	// Login to get token
	loginReq := map[string]interface{}{
		"email":    "benchmark@example.com",
		"password": "benchmarkpass123",
	}
	loginBody, _ := json.Marshal(loginReq)

	loginReqHTTP := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	handler.ServeHTTP(loginResp, loginReqHTTP)

	if loginResp.Code != http.StatusOK {
		t.Fatalf("Failed to login: %s", loginResp.Body.String())
	}

	var authResp map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &authResp)

	data, ok := authResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from response: %v", authResp)
	}

	token, ok := data["token"].(string)
	if !ok {
		t.Fatalf("Failed to extract token from response: %v", data)
	}

	// Benchmark endpoints
	benchmarkIterations := 100

	t.Run("POST_/auth/login", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"email":    "benchmark@example.com",
			"password": "benchmarkpass123",
		}
		loginBody, _ := json.Marshal(loginReq)

		durations := make([]time.Duration, benchmarkIterations)

		for i := 0; i < benchmarkIterations; i++ {
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(resp, req)
			durations[i] = time.Since(start)

			if resp.Code != http.StatusOK && i == 0 {
				t.Logf("First request failed: %s", resp.Body.String())
			}
		}

		// Auth endpoints have higher threshold due to bcrypt hashing for security
		analyzeBenchmarkResults(t, "POST /auth/login", durations, 500*time.Millisecond)
	})

	t.Run("GET_/tasks", func(t *testing.T) {
		// Create some tasks first
		for i := 0; i < 10; i++ {
			createTaskReq := map[string]interface{}{
				"title":       fmt.Sprintf("Benchmark Task %d", i),
				"description": "Task for benchmarking",
				"priority":    "medium",
				"status":      "pending",
			}
			createTaskBody, _ := json.Marshal(createTaskReq)

			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(createTaskBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)
		}

		durations := make([]time.Duration, benchmarkIterations)

		for i := 0; i < benchmarkIterations; i++ {
			req := httptest.NewRequest("GET", "/tasks", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(resp, req)
			durations[i] = time.Since(start)

			if resp.Code != http.StatusOK && i == 0 {
				t.Logf("First request failed: %s", resp.Body.String())
			}
		}

		analyzeBenchmarkResults(t, "GET /tasks", durations, 50*time.Millisecond)
	})

	t.Run("POST_/tasks", func(t *testing.T) {
		durations := make([]time.Duration, benchmarkIterations)

		for i := 0; i < benchmarkIterations; i++ {
			createTaskReq := map[string]interface{}{
				"title":       fmt.Sprintf("Benchmark Create Task %d", i),
				"description": "Task for benchmarking",
				"priority":    "medium",
				"status":      "pending",
			}
			createTaskBody, _ := json.Marshal(createTaskReq)

			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(createTaskBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(resp, req)
			durations[i] = time.Since(start)

			if resp.Code != http.StatusCreated && i == 0 {
				t.Logf("First request failed: %s", resp.Body.String())
			}
		}

		analyzeBenchmarkResults(t, "POST /tasks", durations, 50*time.Millisecond)
	})

	t.Run("GET_/categories", func(t *testing.T) {
		durations := make([]time.Duration, benchmarkIterations)

		for i := 0; i < benchmarkIterations; i++ {
			req := httptest.NewRequest("GET", "/categories", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(resp, req)
			durations[i] = time.Since(start)

			if resp.Code != http.StatusOK && i == 0 {
				t.Logf("First request failed: %s", resp.Body.String())
			}
		}

		analyzeBenchmarkResults(t, "GET /categories", durations, 50*time.Millisecond)
	})

	t.Run("GET_/tags", func(t *testing.T) {
		durations := make([]time.Duration, benchmarkIterations)

		for i := 0; i < benchmarkIterations; i++ {
			req := httptest.NewRequest("GET", "/tags", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(resp, req)
			durations[i] = time.Since(start)

			if resp.Code != http.StatusOK && i == 0 {
				t.Logf("First request failed: %s", resp.Body.String())
			}
		}

		analyzeBenchmarkResults(t, "GET /tags", durations, 50*time.Millisecond)
	})

	t.Run("GET_/activity-logs", func(t *testing.T) {
		durations := make([]time.Duration, benchmarkIterations)

		for i := 0; i < benchmarkIterations; i++ {
			req := httptest.NewRequest("GET", "/activity-logs", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(resp, req)
			durations[i] = time.Since(start)

			if resp.Code != http.StatusOK && i == 0 {
				t.Logf("First request failed: %s", resp.Body.String())
			}
		}

		analyzeBenchmarkResults(t, "GET /activity-logs", durations, 50*time.Millisecond)
	})

	// Check overall metrics
	metrics := middlewares.GetMetrics()
	t.Log("=== Overall Performance Metrics ===")
	t.Logf("Total Requests: %d", getTotalRequestCount(metrics))
	t.Logf("Total Errors: %d", getTotalErrorCount(metrics))

	for endpoint := range metrics.RequestCounts {
		p95 := metrics.CalculateP95(endpoint)
		avg := metrics.CalculateAverage(endpoint)
		t.Logf("Endpoint: %s | P95: %v | Average: %v | Requests: %d", endpoint, p95, avg, metrics.RequestCounts[endpoint])

		// Skip auth endpoints for 50ms threshold check (bcrypt hashing is intentionally slow for security)
		if endpoint == "POST:/auth/login" || endpoint == "POST:/auth/register" {
			if p95 > 500*time.Millisecond {
				t.Errorf("❌ FAIL: Auth endpoint %s has P95 response time %v > 500ms threshold", endpoint, p95)
			} else {
				t.Logf("✅ PASS: Auth endpoint %s meets P95 < 500ms requirement (accounting for bcrypt security)", endpoint)
			}
		} else {
			if p95 > 50*time.Millisecond {
				t.Errorf("❌ FAIL: Endpoint %s has P95 response time %v > 50ms threshold", endpoint, p95)
			} else {
				t.Logf("✅ PASS: Endpoint %s meets P95 < 50ms requirement", endpoint)
			}
		}
	}
}

// analyzeBenchmarkResults calculates and reports benchmark statistics
func analyzeBenchmarkResults(t *testing.T, endpoint string, durations []time.Duration, threshold time.Duration) {
	if len(durations) == 0 {
		t.Fatalf("No duration data for endpoint %s", endpoint)
	}

	// Calculate statistics
	var total time.Duration
	min := durations[0]
	max := durations[0]

	for _, d := range durations {
		total += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	average := total / time.Duration(len(durations))

	// Calculate P95
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	p95Index := int(float64(len(sorted)) * 0.95)
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	p95 := sorted[p95Index]

	t.Logf("=== Benchmark Results: %s ===", endpoint)
	t.Logf("Iterations: %d", len(durations))
	t.Logf("Average: %v", average)
	t.Logf("P95: %v", p95)
	t.Logf("Min: %v", min)
	t.Logf("Max: %v", max)

	if p95 > threshold {
		t.Errorf("❌ FAIL: P95 response time %v exceeds threshold %v", p95, threshold)
	} else {
		t.Logf("✅ PASS: P95 response time %v meets threshold %v", p95, threshold)
	}
}

// getTotalRequestCount calculates total request count from metrics
func getTotalRequestCount(metrics *middlewares.PerformanceMetrics) int64 {
	var total int64
	for _, count := range metrics.RequestCounts {
		total += count
	}
	return total
}

// getTotalErrorCount calculates total error count from metrics
func getTotalErrorCount(metrics *middlewares.PerformanceMetrics) int64 {
	var total int64
	for _, count := range metrics.ErrorCounts {
		total += count
	}
	return total
}
