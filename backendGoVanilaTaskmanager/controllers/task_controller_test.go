package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"taskmanager/models"
	"taskmanager/utils"
	"testing"
)

// MockTaskService is a mock implementation of TaskServiceInterface for testing
type MockTaskService struct {
	CreateFunc              func(userID int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error)
	GetAllFunc              func(userID int64, completed *bool, page, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error)
	GetByIDFunc             func(userID int64, id int64) (*models.Task, error)
	UpdateFunc              func(userID int64, id int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error)
	DeleteFunc              func(userID int64, id int64, ipAddress string, userAgent string) error
	RestoreFunc             func(userID int64, id int64) error
	MarkAsCompletedFunc     func(userID int64, id int64, ipAddress string, userAgent string) error
	MarkAsUncompletedFunc   func(userID int64, id int64, ipAddress string, userAgent string) error
	AddTagToTaskFunc        func(userID int64, taskID int64, tagID int64) error
	RemoveTagFromTaskFunc   func(userID int64, taskID int64, tagID int64) error
	GetAnalyticsSummaryFunc func(userID int64) (map[string]interface{}, error)
	BulkDeleteFunc          func(userID int64, taskIDs []int64, ipAddress string, userAgent string) error
	BulkCompleteFunc        func(userID int64, taskIDs []int64, ipAddress string, userAgent string) error
	ExportToCSVFunc         func(userID int64, completed *bool, search string, priority *models.Priority, categoryID *int64, timezone string) ([]byte, error)
}

// Ensure MockTaskService implements TaskServiceInterface
var _ TaskServiceInterface = (*MockTaskService)(nil)

// Helper function to create request with PathValue support
func createRequest(method, path string, body *bytes.Buffer) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body.Bytes())
	}
	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Add user_id to context for authentication
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	// Add timezone to context for timezone-aware responses
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	return req
}

func (m *MockTaskService) Create(userID int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(userID, task, ipAddress, userAgent)
	}
	return nil, nil
}

func (m *MockTaskService) GetAll(userID int64, completed *bool, page, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(userID, completed, page, limit, search, priority, categoryID, sortBy, sortOrder)
	}
	return nil, nil, nil
}

func (m *MockTaskService) GetByID(userID int64, id int64) (*models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(userID, id)
	}
	return nil, nil
}

func (m *MockTaskService) Update(userID int64, id int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(userID, id, task, ipAddress, userAgent)
	}
	return nil, nil
}

func (m *MockTaskService) Delete(userID int64, id int64, ipAddress string, userAgent string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(userID, id, ipAddress, userAgent)
	}
	return nil
}

func (m *MockTaskService) Restore(userID int64, id int64) error {
	if m.RestoreFunc != nil {
		return m.RestoreFunc(userID, id)
	}
	return nil
}

func (m *MockTaskService) MarkAsCompleted(userID int64, id int64, ipAddress string, userAgent string) error {
	if m.MarkAsCompletedFunc != nil {
		return m.MarkAsCompletedFunc(userID, id, ipAddress, userAgent)
	}
	return nil
}

func (m *MockTaskService) MarkAsUncompleted(userID int64, id int64, ipAddress string, userAgent string) error {
	if m.MarkAsUncompletedFunc != nil {
		return m.MarkAsUncompletedFunc(userID, id, ipAddress, userAgent)
	}
	return nil
}

func (m *MockTaskService) AddTagToTask(userID int64, taskID int64, tagID int64) error {
	if m.AddTagToTaskFunc != nil {
		return m.AddTagToTaskFunc(userID, taskID, tagID)
	}
	return nil
}

func (m *MockTaskService) RemoveTagFromTask(userID int64, taskID int64, tagID int64) error {
	if m.RemoveTagFromTaskFunc != nil {
		return m.RemoveTagFromTaskFunc(userID, taskID, tagID)
	}
	return nil
}

func (m *MockTaskService) GetAnalyticsSummary(userID int64) (map[string]interface{}, error) {
	if m.GetAnalyticsSummaryFunc != nil {
		return m.GetAnalyticsSummaryFunc(userID)
	}
	return nil, nil
}

func (m *MockTaskService) BulkDelete(userID int64, taskIDs []int64, ipAddress string, userAgent string) error {
	if m.BulkDeleteFunc != nil {
		return m.BulkDeleteFunc(userID, taskIDs, ipAddress, userAgent)
	}
	return nil
}

func (m *MockTaskService) BulkComplete(userID int64, taskIDs []int64, ipAddress string, userAgent string) error {
	if m.BulkCompleteFunc != nil {
		return m.BulkCompleteFunc(userID, taskIDs, ipAddress, userAgent)
	}
	return nil
}

func (m *MockTaskService) ExportToCSV(userID int64, completed *bool, search string, priority *models.Priority, categoryID *int64, timezone string) ([]byte, error) {
	if m.ExportToCSVFunc != nil {
		return m.ExportToCSVFunc(userID, completed, search, priority, categoryID, timezone)
	}
	return nil, nil
}

func TestCreateTaskHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		CreateFunc: func(userID int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
			task.ID = 1
			task.CreatedAt = utils.CurrentEpochMillis()
			task.UpdatedAt = utils.CurrentEpochMillis()
			return task, nil
		},
	}

	controller := NewTaskController(mockService)

	taskJSON := `{"title": "Test Task", "description": "Test Description"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestCreateTaskHandler_InvalidJSON(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	invalidJSON := `{"title": "Test Task", "description": invalid}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateTaskHandler_EmptyTitle(t *testing.T) {
	mockService := &MockTaskService{
		CreateFunc: func(userID int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
			// Service layer should validate empty title and return error
			return nil, errors.New("Title tidak boleh kosong")
		},
	}

	controller := NewTaskController(mockService)

	taskJSON := `{"title": "", "description": "Test Description"}`
	req := createRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	w := httptest.NewRecorder()

	controller.CreateTask(w, req)

	// The service layer validates empty title, so we expect 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAllTasksHandler_Success(t *testing.T) {
	mockTasks := []models.Task{
		{ID: 1, Title: "Task 1", Completed: false},
		{ID: 2, Title: "Task 2", Completed: true},
	}

	mockService := &MockTaskService{
		GetAllFunc: func(userID int64, completed *bool, page, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error) {
			meta := map[string]interface{}{
				"page":       page,
				"limit":      limit,
				"total_data": 2,
				"total_page": 1,
				"has_next":   false,
				"has_prev":   false,
			}
			return mockTasks, meta, nil
		},
	}

	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks?page=1&limit=10", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.GetAllTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestGetAllTasksHandler_InvalidPage(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks?page=0", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.GetAllTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAllTasksHandler_InvalidLimit(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks?limit=0", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.GetAllTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAllTasksHandler_InvalidCompleted(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks?completed=invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.GetAllTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAllTasksHandler_EmptySearch(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks?search=", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.GetAllTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetTaskByIDHandler_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:        1,
		Title:     "Test Task",
		Completed: false,
	}

	mockService := &MockTaskService{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return mockTask, nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", controller.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestGetTaskByIDHandler_InvalidID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", controller.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetTaskByIDHandler_NotFound(t *testing.T) {
	mockService := &MockTaskService{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", controller.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateTaskHandler_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:        1,
		Title:     "Updated Task",
		Completed: true,
	}

	mockService := &MockTaskService{
		UpdateFunc: func(userID int64, id int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
			return mockTask, nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", controller.UpdateTask)

	taskJSON := `{"title": "Updated Task", "completed": true}`
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestUpdateTaskHandler_InvalidJSON(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", controller.UpdateTask)

	invalidJSON := `{"title": "Updated Task", "completed": invalid}`
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTaskHandler_InvalidID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", controller.UpdateTask)

	taskJSON := `{"title": "Updated Task"}`
	req := httptest.NewRequest("PUT", "/tasks/invalid", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTaskHandler_NotFound(t *testing.T) {
	mockService := &MockTaskService{
		UpdateFunc: func(userID int64, id int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", controller.UpdateTask)

	taskJSON := `{"title": "Updated Task"}`
	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteTaskHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		DeleteFunc: func(userID int64, id int64, ipAddress string, userAgent string) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestDeleteTaskHandler_InvalidID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteTaskHandler_NotFound(t *testing.T) {
	mockService := &MockTaskService{
		DeleteFunc: func(userID int64, id int64, ipAddress string, userAgent string) error {
			return sql.ErrNoRows
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestMarkTaskAsCompletedHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		MarkAsCompletedFunc: func(userID int64, id int64, ipAddress string, userAgent string) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/complete", controller.MarkTaskAsCompleted)

	req := httptest.NewRequest("PATCH", "/tasks/1/complete", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestMarkTaskAsCompletedHandler_InvalidID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/complete", controller.MarkTaskAsCompleted)

	req := httptest.NewRequest("PATCH", "/tasks/invalid/complete", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestMarkTaskAsUncompletedHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		MarkAsUncompletedFunc: func(userID int64, id int64, ipAddress string, userAgent string) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/uncomplete", controller.MarkTaskAsUncompleted)

	req := httptest.NewRequest("PATCH", "/tasks/1/uncomplete", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestMarkTaskAsUncompletedHandler_InvalidID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/uncomplete", controller.MarkTaskAsUncompleted)

	req := httptest.NewRequest("PATCH", "/tasks/invalid/uncomplete", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRestoreDeletedTaskHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		RestoreFunc: func(userID int64, id int64) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/restore", controller.RestoreDeletedTask)

	req := httptest.NewRequest("PATCH", "/tasks/1/restore", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestRestoreDeletedTaskHandler_InvalidID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/restore", controller.RestoreDeletedTask)

	req := httptest.NewRequest("PATCH", "/tasks/invalid/restore", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetAnalyticsSummaryHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		GetAnalyticsSummaryFunc: func(userID int64) (map[string]interface{}, error) {
			return map[string]interface{}{
				"total_active":          10,
				"total_completed":       5,
				"total_overdue":         2,
				"completion_percentage": 50.0,
				"priority_distribution": map[string]int{
					"LOW":    2,
					"MEDIUM": 4,
					"HIGH":   3,
					"URGENT": 1,
				},
			}, nil
		},
	}

	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks/analytics/summary", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	controller.GetAnalyticsSummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestBulkDeleteTasksHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		BulkDeleteFunc: func(userID int64, taskIDs []int64, ipAddress string, userAgent string) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	requestJSON := `{"task_ids": [1, 2, 3]}`
	req := httptest.NewRequest("POST", "/tasks/bulk-delete", bytes.NewBufferString(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	controller.BulkDeleteTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestBulkDeleteTasksHandler_EmptyTaskIDs(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	requestJSON := `{"task_ids": []}`
	req := httptest.NewRequest("POST", "/tasks/bulk-delete", bytes.NewBufferString(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	controller.BulkDeleteTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestBulkDeleteTasksHandler_TooManyTasks(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	// Create a request with 101 task IDs
	taskIDs := make([]int64, 101)
	for i := 0; i < 101; i++ {
		taskIDs[i] = int64(i + 1)
	}
	requestJSON, _ := json.Marshal(map[string]interface{}{"task_ids": taskIDs})
	req := httptest.NewRequest("POST", "/tasks/bulk-delete", bytes.NewBuffer(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	controller.BulkDeleteTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestBulkCompleteTasksHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		BulkCompleteFunc: func(userID int64, taskIDs []int64, ipAddress string, userAgent string) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	requestJSON := `{"task_ids": [1, 2, 3]}`
	req := httptest.NewRequest("POST", "/tasks/bulk-complete", bytes.NewBufferString(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	controller.BulkCompleteTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestBulkCompleteTasksHandler_EmptyTaskIDs(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	requestJSON := `{"task_ids": []}`
	req := httptest.NewRequest("POST", "/tasks/bulk-complete", bytes.NewBufferString(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	w := httptest.NewRecorder()

	controller.BulkCompleteTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestExportTasksToCSVHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		ExportToCSVFunc: func(userID int64, completed *bool, search string, priority *models.Priority, categoryID *int64, timezone string) ([]byte, error) {
			csvContent := "ID,Title,Completed\n1,Test Task,false\n"
			return []byte(csvContent), nil
		},
	}

	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks/export/csv", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.ExportTasksToCSV(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("Expected Content-Type text/csv, got %s", contentType)
	}

	contentDisposition := w.Header().Get("Content-Disposition")
	if contentDisposition != "attachment; filename=tasks_export.csv" {
		t.Errorf("Expected Content-Disposition header, got %s", contentDisposition)
	}
}

func TestExportTasksToCSVHandler_InvalidPriority(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks/export/csv?priority=INVALID", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.ExportTasksToCSV(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestExportTasksToCSVHandler_InvalidCategoryID(t *testing.T) {
	mockService := &MockTaskService{}
	controller := NewTaskController(mockService)

	req := httptest.NewRequest("GET", "/tasks/export/csv?category_id=invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.ExportTasksToCSV(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
