package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskmanager/models"
	"testing"
)

// MockActivityLogService is a mock implementation of ActivityLogServiceInterface for testing
type MockActivityLogService struct {
	logActivityFunc func(userID int64, taskID *int64, action string, entityType string, entityID *int64, details string, ipAddress string, userAgent string) error
	getUserLogsFunc func(userID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error)
	getTaskLogsFunc func(taskID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error)
	getByIDFunc     func(id int64) (*models.ActivityLog, error)
}

func (m *MockActivityLogService) LogActivity(userID int64, taskID *int64, action string, entityType string, entityID *int64, details string, ipAddress string, userAgent string) error {
	if m.logActivityFunc != nil {
		return m.logActivityFunc(userID, taskID, action, entityType, entityID, details, ipAddress, userAgent)
	}
	return nil
}

func (m *MockActivityLogService) GetUserActivityLogs(userID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
	if m.getUserLogsFunc != nil {
		return m.getUserLogsFunc(userID, page, limit)
	}
	return []models.ActivityLog{}, map[string]interface{}{}, nil
}

func (m *MockActivityLogService) GetTaskActivityLogs(taskID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
	if m.getTaskLogsFunc != nil {
		return m.getTaskLogsFunc(taskID, page, limit)
	}
	return []models.ActivityLog{}, map[string]interface{}{}, nil
}

func (m *MockActivityLogService) GetActivityLogByID(id int64) (*models.ActivityLog, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return &models.ActivityLog{ID: id, Action: "create"}, nil
}

func TestNewActivityLogController(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	if controller == nil {
		t.Fatal("NewActivityLogController returned nil")
	}

	if controller.service != mockService {
		t.Error("NewActivityLogController did not set service correctly")
	}
}

func TestActivityLogController_GetUserActivityLogs_Success(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	mockService.getUserLogsFunc = func(userID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
		return []models.ActivityLog{{ID: 1, Action: "create"}}, map[string]interface{}{"page": 1, "limit": 10}, nil
	}

	req := httptest.NewRequest("GET", "/activity-logs?page=1&limit=10", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetUserActivityLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestActivityLogController_GetUserActivityLogs_InvalidPage(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	req := httptest.NewRequest("GET", "/activity-logs?page=0&limit=10", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetUserActivityLogs(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestActivityLogController_GetUserActivityLogs_InvalidLimit(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	req := httptest.NewRequest("GET", "/activity-logs?page=1&limit=0", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetUserActivityLogs(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestActivityLogController_GetUserActivityLogs_Error(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	mockService.getUserLogsFunc = func(userID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
		return nil, nil, errors.New("service error")
	}

	req := httptest.NewRequest("GET", "/activity-logs?page=1&limit=10", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetUserActivityLogs(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for service error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestActivityLogController_GetTaskActivityLogs_Success(t *testing.T) {
	mockService := &MockActivityLogService{
		getTaskLogsFunc: func(taskID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
			return []models.ActivityLog{{ID: 1, Action: "update"}}, map[string]interface{}{"page": 1, "limit": 10}, nil
		},
	}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/activity-logs", controller.GetTaskActivityLogs)

	req := httptest.NewRequest("GET", "/tasks/1/activity-logs?page=1&limit=10", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestActivityLogController_GetTaskActivityLogs_InvalidTaskID(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/activity-logs", controller.GetTaskActivityLogs)

	req := httptest.NewRequest("GET", "/tasks/invalid/activity-logs?page=1&limit=10", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestActivityLogController_GetTaskActivityLogs_Error(t *testing.T) {
	mockService := &MockActivityLogService{
		getTaskLogsFunc: func(taskID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
			return nil, nil, errors.New("service error")
		},
	}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/activity-logs", controller.GetTaskActivityLogs)

	req := httptest.NewRequest("GET", "/tasks/1/activity-logs?page=1&limit=10", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestActivityLogController_GetActivityLogByID_Success(t *testing.T) {
	mockService := &MockActivityLogService{
		getByIDFunc: func(id int64) (*models.ActivityLog, error) {
			return &models.ActivityLog{ID: id, Action: "create", UserID: 1}, nil
		},
	}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /activity-logs/{id}", controller.GetActivityLogByID)

	req := httptest.NewRequest("GET", "/activity-logs/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestActivityLogController_GetActivityLogByID_NotFound(t *testing.T) {
	mockService := &MockActivityLogService{
		getByIDFunc: func(id int64) (*models.ActivityLog, error) {
			return nil, sql.ErrNoRows
		},
	}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /activity-logs/{id}", controller.GetActivityLogByID)

	req := httptest.NewRequest("GET", "/activity-logs/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestActivityLogController_GetActivityLogByID_Forbidden(t *testing.T) {
	mockService := &MockActivityLogService{
		getByIDFunc: func(id int64) (*models.ActivityLog, error) {
			return &models.ActivityLog{ID: id, Action: "create", UserID: 2}, nil
		},
	}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /activity-logs/{id}", controller.GetActivityLogByID)

	req := httptest.NewRequest("GET", "/activity-logs/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestActivityLogController_GetActivityLogByID_InvalidID(t *testing.T) {
	mockService := &MockActivityLogService{}
	controller := NewActivityLogController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /activity-logs/{id}", controller.GetActivityLogByID)

	req := httptest.NewRequest("GET", "/activity-logs/invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
