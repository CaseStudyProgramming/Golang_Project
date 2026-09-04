package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"taskmanager/config"
	"taskmanager/controllers"
	"taskmanager/middlewares"
	"taskmanager/models"
	"taskmanager/routes"
	"taskmanager/services"
	"taskmanager/utils"
)

// TestAPIPerformance_BasicOperations benchmarks basic API operations
func TestAPIPerformance_BasicOperations(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Initialize components
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

	// Initialize other controllers
	categoryModel := models.NewCategoryModel(db)
	categoryService := services.NewCategoryService(categoryModel)
	categoryController := controllers.NewCategoryController(categoryService)

	tagService := services.NewTagService(tagModel)
	tagController := controllers.NewTagController(tagService)

	subtaskModel := models.NewSubtaskModel(db)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)
	subtaskController := controllers.NewSubtaskController(subtaskService)

	activityLogController := controllers.NewActivityLogController(activityLogService)

	swaggerController := controllers.NewSwaggerController()
	metricsController := controllers.NewMetricsController()

	// Reset metrics
	middlewares.ResetMetrics()

	// Setup router
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Apply middlewares
	handler := middlewares.Recovery(middlewares.Logger(middlewares.PerformanceMonitoring(mux)))

	// Register a test user
	registerReq := map[string]string{
		"name":     "perftest",
		"email":    "perftest@example.com",
		"password": "password123",
		"timezone": "UTC",
	}
	reqBody, _ := json.Marshal(registerReq)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to register user: %d", rr.Code)
	}

	var authResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &authResponse)
	token := authResponse["token"].(string)

	t.Log("✅ Registered test user")

	// Benchmark task creation
	t.Run("CreateTask", func(t *testing.T) {
		iterations := 50
		var totalDuration time.Duration

		for i := 0; i < iterations; i++ {
			taskReq := map[string]interface{}{
				"title":       fmt.Sprintf("Performance Task %d", i),
				"description": "Test task for performance benchmarking",
				"priority":    "medium",
			}
			reqBody, _ := json.Marshal(taskReq)

			startTime := time.Now()
			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			duration := time.Since(startTime)
			totalDuration += duration

			if rr.Code != http.StatusCreated {
				t.Errorf("Failed to create task: %d", rr.Code)
			}
		}

		avgDuration := totalDuration / time.Duration(iterations)
		t.Logf("✅ Created %d tasks in %v (avg: %v per task)", iterations, totalDuration, avgDuration)

		if avgDuration > 50*time.Millisecond {
			t.Logf("⚠️  Warning: Average response time exceeds 50ms threshold: %v", avgDuration)
		}
	})

	// Benchmark task listing
	t.Run("ListTasks", func(t *testing.T) {
		iterations := 50
		var totalDuration time.Duration

		for i := 0; i < iterations; i++ {
			startTime := time.Now()
			req := httptest.NewRequest("GET", "/tasks?page=1&limit=10", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			duration := time.Since(startTime)
			totalDuration += duration

			if rr.Code != http.StatusOK {
				t.Errorf("Failed to list tasks: %d", rr.Code)
			}
		}

		avgDuration := totalDuration / time.Duration(iterations)
		t.Logf("✅ Listed tasks %d times in %v (avg: %v per request)", iterations, totalDuration, avgDuration)

		if avgDuration > 50*time.Millisecond {
			t.Logf("⚠️  Warning: Average response time exceeds 50ms threshold: %v", avgDuration)
		}
	})

	// Benchmark task update
	t.Run("UpdateTask", func(t *testing.T) {
		// First create a task
		taskReq := map[string]interface{}{
			"title":       "Task to update",
			"description": "Test task",
		}
		reqBody, _ := json.Marshal(taskReq)
		req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var createdTask map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &createdTask)
		taskID := int64(createdTask["id"].(float64))

		iterations := 50
		var totalDuration time.Duration

		for i := 0; i < iterations; i++ {
			updateReq := map[string]interface{}{
				"title": fmt.Sprintf("Updated Task %d", i),
			}
			reqBody, _ := json.Marshal(updateReq)

			startTime := time.Now()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/tasks/%d", taskID), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			duration := time.Since(startTime)
			totalDuration += duration

			if rr.Code != http.StatusOK {
				t.Errorf("Failed to update task: %d", rr.Code)
			}
		}

		avgDuration := totalDuration / time.Duration(iterations)
		t.Logf("✅ Updated task %d times in %v (avg: %v per request)", iterations, totalDuration, avgDuration)

		if avgDuration > 50*time.Millisecond {
			t.Logf("⚠️  Warning: Average response time exceeds 50ms threshold: %v", avgDuration)
		}
	})

	// Check metrics
	t.Run("CheckMetrics", func(t *testing.T) {
		metrics := middlewares.GetMetrics()

		t.Logf("📊 Performance Metrics Summary:")
		t.Logf("   Total Requests: %d", len(metrics.RequestCounts))

		for endpoint, count := range metrics.RequestCounts {
			avgTime := metrics.CalculateAverage(endpoint)
			p95Time := metrics.CalculateP95(endpoint)
			t.Logf("   %s: %d requests, avg: %v, p95: %v", endpoint, count, avgTime, p95Time)

			if p95Time > 50*time.Millisecond {
				t.Logf("   ⚠️  %s exceeds 50ms p95 threshold: %v", endpoint, p95Time)
			}
		}
	})
}

// TestAPIPerformance_ConcurrentRequests tests concurrent request handling
func TestAPIPerformance_ConcurrentRequests(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Initialize components
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

	// Initialize other controllers
	categoryModel := models.NewCategoryModel(db)
	categoryService := services.NewCategoryService(categoryModel)
	categoryController := controllers.NewCategoryController(categoryService)

	tagService := services.NewTagService(tagModel)
	tagController := controllers.NewTagController(tagService)

	subtaskModel := models.NewSubtaskModel(db)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)
	subtaskController := controllers.NewSubtaskController(subtaskService)

	activityLogController := controllers.NewActivityLogController(activityLogService)

	swaggerController := controllers.NewSwaggerController()
	metricsController := controllers.NewMetricsController()

	// Reset metrics
	middlewares.ResetMetrics()

	// Setup router
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Apply middlewares
	handler := middlewares.Recovery(middlewares.Logger(middlewares.PerformanceMonitoring(mux)))

	// Register a test user
	registerReq := map[string]string{
		"name":     "concurrenttest",
		"email":    "concurrent@example.com",
		"password": "password123",
		"timezone": "UTC",
	}
	reqBody, _ := json.Marshal(registerReq)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to register user: %d", rr.Code)
	}

	var authResponse map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &authResponse)
	token := authResponse["token"].(string)

	// Test concurrent requests
	concurrentRequests := 20
	done := make(chan bool, concurrentRequests)
	startTime := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			defer func() { done <- true }()

			taskReq := map[string]interface{}{
				"title":       fmt.Sprintf("Concurrent Task %d", index),
				"description": "Test task for concurrent performance",
			}
			reqBody, _ := json.Marshal(taskReq)

			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				t.Errorf("Concurrent request %d failed: %d", index, rr.Code)
			}
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}

	totalDuration := time.Since(startTime)
	avgDuration := totalDuration / time.Duration(concurrentRequests)

	t.Logf("✅ Completed %d concurrent requests in %v (avg: %v per request)", concurrentRequests, totalDuration, avgDuration)

	// Check metrics
	metrics := middlewares.GetMetrics()
	t.Logf("📊 Concurrent Request Metrics:")
	for endpoint, count := range metrics.RequestCounts {
		if count > 0 {
			avgTime := metrics.CalculateAverage(endpoint)
			p95Time := metrics.CalculateP95(endpoint)
			t.Logf("   %s: %d requests, avg: %v, p95: %v", endpoint, count, avgTime, p95Time)
		}
	}
}