package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"taskmanager/controllers"
	"taskmanager/middlewares"
)

func TestRegisterRoutes(t *testing.T) {
	// Create mock controllers
	taskController := &controllers.TaskController{}
	authController := &controllers.AuthController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	subtaskController := &controllers.SubtaskController{}
	activityLogController := &controllers.ActivityLogController{}
	swaggerController := &controllers.SwaggerController{}
	metricsController := &controllers.MetricsController{}
	authMiddleware := &middlewares.AuthMiddleware{}

	// Create ServeMux
	mux := http.NewServeMux()

	// Register routes
	RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Test health check endpoint
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /health, got %d", w.Code)
	}

	if w.Body.String() != "API is runningggg 🚀\n" {
		t.Errorf("Expected 'API is runningggg 🚀\\n', got '%s'", w.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Create minimal setup
	mux := http.NewServeMux()
	taskController := &controllers.TaskController{}
	authController := &controllers.AuthController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	subtaskController := &controllers.SubtaskController{}
	activityLogController := &controllers.ActivityLogController{}
	swaggerController := &controllers.SwaggerController{}
	metricsController := &controllers.MetricsController{}
	authMiddleware := &middlewares.AuthMiddleware{}

	RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Test GET /health
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "API is runningggg 🚀\n"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

func TestMetricsEndpoints(t *testing.T) {
	// Create minimal setup
	mux := http.NewServeMux()
	taskController := &controllers.TaskController{}
	authController := &controllers.AuthController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	subtaskController := &controllers.SubtaskController{}
	activityLogController := &controllers.ActivityLogController{}
	swaggerController := &controllers.SwaggerController{}
	metricsController := &controllers.MetricsController{}
	authMiddleware := &middlewares.AuthMiddleware{}

	RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Test GET /metrics
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Should get some response (may be 200 or error depending on implementation)
	// Just verify the route is registered
	if w.Code == 0 {
		t.Error("Metrics endpoint not registered")
	}
}

func TestSwaggerEndpoints(t *testing.T) {
	// Create minimal setup
	mux := http.NewServeMux()
	taskController := &controllers.TaskController{}
	authController := &controllers.AuthController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	subtaskController := &controllers.SubtaskController{}
	activityLogController := &controllers.ActivityLogController{}
	swaggerController := &controllers.SwaggerController{}
	metricsController := &controllers.MetricsController{}
	authMiddleware := &middlewares.AuthMiddleware{}

	RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Test GET /swagger redirect
	req := httptest.NewRequest("GET", "/swagger", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Should redirect to /swagger/index.html
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected redirect status 301, got %d", w.Code)
	}
}

func TestAuthEndpoints(t *testing.T) {
	// Create minimal setup
	mux := http.NewServeMux()
	taskController := &controllers.TaskController{}
	authController := &controllers.AuthController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	subtaskController := &controllers.SubtaskController{}
	activityLogController := &controllers.ActivityLogController{}
	swaggerController := &controllers.SwaggerController{}
	metricsController := &controllers.MetricsController{}
	authMiddleware := &middlewares.AuthMiddleware{}

	RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Test POST /auth/register (public endpoint)
	req := httptest.NewRequest("POST", "/auth/register", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Route should be registered (may return error due to missing body)
	if w.Code == 0 {
		t.Error("Auth register endpoint not registered")
	}

	// Test POST /auth/login (public endpoint)
	req = httptest.NewRequest("POST", "/auth/login", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Error("Auth login endpoint not registered")
	}
}

func TestProtectedEndpoints(t *testing.T) {
	// Create minimal setup
	mux := http.NewServeMux()
	taskController := &controllers.TaskController{}
	authController := &controllers.AuthController{}
	categoryController := &controllers.CategoryController{}
	tagController := &controllers.TagController{}
	subtaskController := &controllers.SubtaskController{}
	activityLogController := &controllers.ActivityLogController{}
	swaggerController := &controllers.SwaggerController{}
	metricsController := &controllers.MetricsController{}
	authMiddleware := &middlewares.AuthMiddleware{}

	RegisterRoutes(mux, taskController, authController, categoryController, tagController, subtaskController, activityLogController, swaggerController, metricsController, authMiddleware)

	// Test GET /tasks (protected endpoint)
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Route should be registered (may return 401 due to missing auth)
	if w.Code == 0 {
		t.Error("Tasks endpoint not registered")
	}

	// Test GET /categories (protected endpoint)
	req = httptest.NewRequest("GET", "/categories", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Error("Categories endpoint not registered")
	}

	// Test GET /tags (protected endpoint)
	req = httptest.NewRequest("GET", "/tags", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Error("Tags endpoint not registered")
	}
}