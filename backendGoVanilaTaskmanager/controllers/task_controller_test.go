package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"taskmanager/models"
	"testing"
	"time"
)

// MockTaskService is a mock implementation of TaskServiceInterface for testing
type MockTaskService struct {
	CreateFunc            func(task *models.Task) (*models.Task, error)
	GetAllFunc            func(completed *bool, page, limit int, search string, priority *models.Priority, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error)
	GetByIDFunc           func(id int64) (*models.Task, error)
	UpdateFunc            func(id int64, task *models.Task) (*models.Task, error)
	DeleteFunc            func(id int64) error
	RestoreFunc           func(id int64) error
	MarkAsCompletedFunc   func(id int64) error
	MarkAsUncompletedFunc func(id int64) error
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
	return req
}

func (m *MockTaskService) Create(task *models.Task) (*models.Task, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(task)
	}
	return nil, nil
}

func (m *MockTaskService) GetAll(completed *bool, page, limit int, search string, priority *models.Priority, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(completed, page, limit, search, priority, sortBy, sortOrder)
	}
	return nil, nil, nil
}

func (m *MockTaskService) GetByID(id int64) (*models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockTaskService) Update(id int64, task *models.Task) (*models.Task, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, task)
	}
	return nil, nil
}

func (m *MockTaskService) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockTaskService) Restore(id int64) error {
	if m.RestoreFunc != nil {
		return m.RestoreFunc(id)
	}
	return nil
}

func (m *MockTaskService) MarkAsCompleted(id int64) error {
	if m.MarkAsCompletedFunc != nil {
		return m.MarkAsCompletedFunc(id)
	}
	return nil
}

func (m *MockTaskService) MarkAsUncompleted(id int64) error {
	if m.MarkAsUncompletedFunc != nil {
		return m.MarkAsUncompletedFunc(id)
	}
	return nil
}

func TestCreateTaskHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
			task.ID = 1
			task.CreatedAt = time.Now()
			task.UpdatedAt = time.Now()
			return task, nil
		},
	}

	controller := NewTaskController(mockService)

	taskJSON := `{"title": "Test Task", "description": "Test Description"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
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
	w := httptest.NewRecorder()

	controller.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateTaskHandler_EmptyTitle(t *testing.T) {
	mockService := &MockTaskService{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
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
		GetAllFunc: func(completed *bool, page, limit int, search string, priority *models.Priority, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error) {
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
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return mockTask, nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", controller.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/1", nil)
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
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetTaskByIDHandler_NotFound(t *testing.T) {
	mockService := &MockTaskService{
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", controller.GetTaskByID)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
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
		UpdateFunc: func(id int64, task *models.Task) (*models.Task, error) {
			return mockTask, nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", controller.UpdateTask)

	taskJSON := `{"title": "Updated Task", "completed": true}`
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
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
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateTaskHandler_NotFound(t *testing.T) {
	mockService := &MockTaskService{
		UpdateFunc: func(id int64, task *models.Task) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", controller.UpdateTask)

	taskJSON := `{"title": "Updated Task"}`
	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteTaskHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		DeleteFunc: func(id int64) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
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
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteTaskHandler_NotFound(t *testing.T) {
	mockService := &MockTaskService{
		DeleteFunc: func(id int64) error {
			return sql.ErrNoRows
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", controller.DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestMarkTaskAsCompletedHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		MarkAsCompletedFunc: func(id int64) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/complete", controller.MarkTaskAsCompleted)

	req := httptest.NewRequest("PATCH", "/tasks/1/complete", nil)
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
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestMarkTaskAsUncompletedHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		MarkAsUncompletedFunc: func(id int64) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/uncomplete", controller.MarkTaskAsUncompleted)

	req := httptest.NewRequest("PATCH", "/tasks/1/uncomplete", nil)
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
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRestoreDeletedTaskHandler_Success(t *testing.T) {
	mockService := &MockTaskService{
		RestoreFunc: func(id int64) error {
			return nil
		},
	}

	controller := NewTaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /tasks/{id}/restore", controller.RestoreDeletedTask)

	req := httptest.NewRequest("PATCH", "/tasks/1/restore", nil)
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
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
