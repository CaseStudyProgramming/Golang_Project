package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskmanager/models"
	"taskmanager/utils"
	"testing"
)

// MockSubtaskService is a mock implementation of SubtaskServiceInterface for testing
type MockSubtaskService struct {
	CreateFunc      func(userID int64, taskID int64, subtask *models.Subtask) (*models.Subtask, error)
	GetByTaskIDFunc func(userID int64, taskID int64) ([]models.Subtask, error)
	GetByIDFunc     func(userID int64, id int64) (*models.Subtask, error)
	UpdateFunc      func(userID int64, id int64, subtask *models.Subtask) (*models.Subtask, error)
	DeleteFunc      func(userID int64, id int64) error
	ToggleFunc      func(userID int64, id int64) error
}

// Ensure MockSubtaskService implements SubtaskServiceInterface
var _ SubtaskServiceInterface = (*MockSubtaskService)(nil)

func (m *MockSubtaskService) Create(userID int64, taskID int64, subtask *models.Subtask) (*models.Subtask, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(userID, taskID, subtask)
	}
	return nil, nil
}

func (m *MockSubtaskService) GetByTaskID(userID int64, taskID int64) ([]models.Subtask, error) {
	if m.GetByTaskIDFunc != nil {
		return m.GetByTaskIDFunc(userID, taskID)
	}
	return nil, nil
}

func (m *MockSubtaskService) GetByID(userID int64, id int64) (*models.Subtask, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(userID, id)
	}
	return nil, nil
}

func (m *MockSubtaskService) Update(userID int64, id int64, subtask *models.Subtask) (*models.Subtask, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(userID, id, subtask)
	}
	return nil, nil
}

func (m *MockSubtaskService) Delete(userID int64, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(userID, id)
	}
	return nil
}

func (m *MockSubtaskService) Toggle(userID int64, id int64) error {
	if m.ToggleFunc != nil {
		return m.ToggleFunc(userID, id)
	}
	return nil
}

// Helper function to create request with context
func createSubtaskRequest(method, path string, body *bytes.Buffer) *http.Request {
	var bodyReader = bytes.NewReader(body.Bytes())
	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	return req
}

func TestNewSubtaskController(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	if controller == nil {
		t.Error("Expected non-nil controller")
	}

	if controller.service != mockService {
		t.Error("Expected service to be set")
	}
}

func TestCreateSubtask_Success(t *testing.T) {
	mockService := &MockSubtaskService{
		CreateFunc: func(userID int64, taskID int64, subtask *models.Subtask) (*models.Subtask, error) {
			subtask.ID = 1
			subtask.CreatedAt = utils.CurrentEpochMillis()
			subtask.UpdatedAt = utils.CurrentEpochMillis()
			return subtask, nil
		},
	}

	controller := NewSubtaskController(mockService)

	subtaskJSON := `{"title": "Test Subtask"}`
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/{id}/subtasks", controller.CreateSubtask)

	req := httptest.NewRequest("POST", "/tasks/1/subtasks", bytes.NewBufferString(subtaskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestCreateSubtask_InvalidJSON(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/{id}/subtasks", controller.CreateSubtask)

	invalidJSON := `{"title": "Test Subtask", invalid}`
	req := httptest.NewRequest("POST", "/tasks/1/subtasks", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateSubtask_InvalidTaskID(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/{id}/subtasks", controller.CreateSubtask)

	subtaskJSON := `{"title": "Test Subtask"}`
	req := httptest.NewRequest("POST", "/tasks/invalid/subtasks", bytes.NewBufferString(subtaskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateSubtask_ServiceError(t *testing.T) {
	mockService := &MockSubtaskService{
		CreateFunc: func(userID int64, taskID int64, subtask *models.Subtask) (*models.Subtask, error) {
			return nil, utils.ErrMissingRequired
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/{id}/subtasks", controller.CreateSubtask)

	subtaskJSON := `{"title": ""}`
	req := httptest.NewRequest("POST", "/tasks/1/subtasks", bytes.NewBufferString(subtaskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetSubtasksByTaskID_Success(t *testing.T) {
	mockSubtasks := []models.Subtask{
		{ID: 1, Title: "Subtask 1", TaskID: 1, IsCompleted: false},
		{ID: 2, Title: "Subtask 2", TaskID: 1, IsCompleted: true},
	}

	mockService := &MockSubtaskService{
		GetByTaskIDFunc: func(userID int64, taskID int64) ([]models.Subtask, error) {
			return mockSubtasks, nil
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/subtasks", controller.GetSubtasksByTaskID)

	req := httptest.NewRequest("GET", "/tasks/1/subtasks", nil)
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

func TestGetSubtasksByTaskID_InvalidTaskID(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/subtasks", controller.GetSubtasksByTaskID)

	req := httptest.NewRequest("GET", "/tasks/invalid/subtasks", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetSubtasksByTaskID_ServiceError(t *testing.T) {
	mockService := &MockSubtaskService{
		GetByTaskIDFunc: func(userID int64, taskID int64) ([]models.Subtask, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/subtasks", controller.GetSubtasksByTaskID)

	req := httptest.NewRequest("GET", "/tasks/999/subtasks", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetSubtaskByID_Success(t *testing.T) {
	mockSubtask := &models.Subtask{
		ID:          1,
		Title:       "Test Subtask",
		TaskID:      1,
		IsCompleted: false,
	}

	mockService := &MockSubtaskService{
		GetByIDFunc: func(userID int64, id int64) (*models.Subtask, error) {
			return mockSubtask, nil
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /subtasks/{id}", controller.GetSubtaskByID)

	req := httptest.NewRequest("GET", "/subtasks/1", nil)
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

func TestGetSubtaskByID_InvalidID(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /subtasks/{id}", controller.GetSubtaskByID)

	req := httptest.NewRequest("GET", "/subtasks/invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetSubtaskByID_NotFound(t *testing.T) {
	mockService := &MockSubtaskService{
		GetByIDFunc: func(userID int64, id int64) (*models.Subtask, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /subtasks/{id}", controller.GetSubtaskByID)

	req := httptest.NewRequest("GET", "/subtasks/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateSubtask_Success(t *testing.T) {
	mockSubtask := &models.Subtask{
		ID:          1,
		Title:       "Updated Subtask",
		TaskID:      1,
		IsCompleted: true,
	}

	mockService := &MockSubtaskService{
		UpdateFunc: func(userID int64, id int64, subtask *models.Subtask) (*models.Subtask, error) {
			return mockSubtask, nil
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subtasks/{id}", controller.UpdateSubtask)

	subtaskJSON := `{"title": "Updated Subtask", "completed": true}`
	req := httptest.NewRequest("PUT", "/subtasks/1", bytes.NewBufferString(subtaskJSON))
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

func TestUpdateSubtask_InvalidJSON(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subtasks/{id}", controller.UpdateSubtask)

	invalidJSON := `{"title": "Updated Subtask", "completed": invalid}`
	req := httptest.NewRequest("PUT", "/subtasks/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateSubtask_InvalidID(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subtasks/{id}", controller.UpdateSubtask)

	subtaskJSON := `{"title": "Updated Subtask"}`
	req := httptest.NewRequest("PUT", "/subtasks/invalid", bytes.NewBufferString(subtaskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUpdateSubtask_NotFound(t *testing.T) {
	mockService := &MockSubtaskService{
		UpdateFunc: func(userID int64, id int64, subtask *models.Subtask) (*models.Subtask, error) {
			return nil, sql.ErrNoRows
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /subtasks/{id}", controller.UpdateSubtask)

	subtaskJSON := `{"title": "Updated Subtask"}`
	req := httptest.NewRequest("PUT", "/subtasks/999", bytes.NewBufferString(subtaskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestDeleteSubtask_Success(t *testing.T) {
	mockService := &MockSubtaskService{
		DeleteFunc: func(userID int64, id int64) error {
			return nil
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /subtasks/{id}", controller.DeleteSubtask)

	req := httptest.NewRequest("DELETE", "/subtasks/1", nil)
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

func TestDeleteSubtask_InvalidID(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /subtasks/{id}", controller.DeleteSubtask)

	req := httptest.NewRequest("DELETE", "/subtasks/invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteSubtask_NotFound(t *testing.T) {
	mockService := &MockSubtaskService{
		DeleteFunc: func(userID int64, id int64) error {
			return sql.ErrNoRows
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /subtasks/{id}", controller.DeleteSubtask)

	req := httptest.NewRequest("DELETE", "/subtasks/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestToggleSubtask_Success(t *testing.T) {
	mockService := &MockSubtaskService{
		ToggleFunc: func(userID int64, id int64) error {
			return nil
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /subtasks/{id}/toggle", controller.ToggleSubtask)

	req := httptest.NewRequest("PATCH", "/subtasks/1/toggle", nil)
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

func TestToggleSubtask_InvalidID(t *testing.T) {
	mockService := &MockSubtaskService{}
	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /subtasks/{id}/toggle", controller.ToggleSubtask)

	req := httptest.NewRequest("PATCH", "/subtasks/invalid/toggle", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestToggleSubtask_NotFound(t *testing.T) {
	mockService := &MockSubtaskService{
		ToggleFunc: func(userID int64, id int64) error {
			return sql.ErrNoRows
		},
	}

	controller := NewSubtaskController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /subtasks/{id}/toggle", controller.ToggleSubtask)

	req := httptest.NewRequest("PATCH", "/subtasks/999/toggle", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
