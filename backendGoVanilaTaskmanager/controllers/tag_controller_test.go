package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskmanager/models"
	"testing"
)

// MockTagService is a mock implementation of TagServiceInterface for testing
type MockTagService struct {
	createFunc        func(userID int64, tag *models.Tag) (*models.Tag, error)
	getAllFunc        func(userID int64) ([]models.Tag, error)
	getByIDFunc       func(userID int64, id int64) (*models.Tag, error)
	updateFunc        func(userID int64, id int64, tag *models.Tag) (*models.Tag, error)
	deleteFunc        func(userID int64, id int64) error
	getTagsByTaskFunc func(taskID int64) ([]models.Tag, error)
	getTasksByTagFunc func(tagID int64) ([]models.Task, error)
}

func (m *MockTagService) Create(userID int64, tag *models.Tag) (*models.Tag, error) {
	if m.createFunc != nil {
		return m.createFunc(userID, tag)
	}
	return &models.Tag{ID: 1, Name: tag.Name, UserID: userID}, nil
}

func (m *MockTagService) GetAll(userID int64) ([]models.Tag, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(userID)
	}
	return []models.Tag{}, nil
}

func (m *MockTagService) GetByID(userID int64, id int64) (*models.Tag, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(userID, id)
	}
	return &models.Tag{ID: id, Name: "Test", UserID: userID}, nil
}

func (m *MockTagService) Update(userID int64, id int64, tag *models.Tag) (*models.Tag, error) {
	if m.updateFunc != nil {
		return m.updateFunc(userID, id, tag)
	}
	return &models.Tag{ID: id, Name: tag.Name, UserID: userID}, nil
}

func (m *MockTagService) Delete(userID int64, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(userID, id)
	}
	return nil
}

func (m *MockTagService) GetTagsByTaskID(taskID int64) ([]models.Tag, error) {
	if m.getTagsByTaskFunc != nil {
		return m.getTagsByTaskFunc(taskID)
	}
	return []models.Tag{}, nil
}

func (m *MockTagService) GetTasksByTagID(tagID int64) ([]models.Task, error) {
	if m.getTasksByTagFunc != nil {
		return m.getTasksByTagFunc(tagID)
	}
	return []models.Task{}, nil
}

func TestNewTagController(t *testing.T) {
	mockService := &MockTagService{}
	controller := NewTagController(mockService)

	if controller == nil {
		t.Fatal("NewTagController returned nil")
	}

	if controller.service != mockService {
		t.Error("NewTagController did not set service correctly")
	}
}

func TestTagController_GetAllTags_Success(t *testing.T) {
	mockService := &MockTagService{}
	controller := NewTagController(mockService)

	mockService.getAllFunc = func(userID int64) ([]models.Tag, error) {
		return []models.Tag{{ID: 1, Name: "Urgent"}}, nil
	}

	req := httptest.NewRequest("GET", "/tags", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetAllTags(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTagController_GetAllTags_Error(t *testing.T) {
	mockService := &MockTagService{}
	controller := NewTagController(mockService)

	mockService.getAllFunc = func(userID int64) ([]models.Tag, error) {
		return nil, errors.New("service error")
	}

	req := httptest.NewRequest("GET", "/tags", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetAllTags(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for service error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTagController_GetTagByID_Success(t *testing.T) {
	mockService := &MockTagService{
		getByIDFunc: func(userID int64, id int64) (*models.Tag, error) {
			return &models.Tag{ID: id, Name: "Test Tag", UserID: userID}, nil
		},
	}
	controller := NewTagController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tags/{id}", controller.GetTagByID)

	req := httptest.NewRequest("GET", "/tags/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTagController_GetTagByID_NotFound(t *testing.T) {
	mockService := &MockTagService{
		getByIDFunc: func(userID int64, id int64) (*models.Tag, error) {
			return nil, sql.ErrNoRows
		},
	}
	controller := NewTagController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tags/{id}", controller.GetTagByID)

	req := httptest.NewRequest("GET", "/tags/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTagController_DeleteTag_Success(t *testing.T) {
	mockService := &MockTagService{
		deleteFunc: func(userID int64, id int64) error {
			return nil
		},
	}
	controller := NewTagController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tags/{id}", controller.DeleteTag)

	req := httptest.NewRequest("DELETE", "/tags/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, w.Code)
	}
}

func TestTagController_DeleteTag_NotFound(t *testing.T) {
	mockService := &MockTagService{
		deleteFunc: func(userID int64, id int64) error {
			return sql.ErrNoRows
		},
	}
	controller := NewTagController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tags/{id}", controller.DeleteTag)

	req := httptest.NewRequest("DELETE", "/tags/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTagController_GetTagsByTaskID_Success(t *testing.T) {
	mockService := &MockTagService{
		getTagsByTaskFunc: func(taskID int64) ([]models.Tag, error) {
			return []models.Tag{{ID: 1, Name: "Urgent"}}, nil
		},
	}
	controller := NewTagController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}/tags", controller.GetTagsByTaskID)

	req := httptest.NewRequest("GET", "/tasks/1/tags", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTagController_GetTasksByTagID_Success(t *testing.T) {
	mockService := &MockTagService{
		getTasksByTagFunc: func(tagID int64) ([]models.Task, error) {
			return []models.Task{{ID: 1, Title: "Test Task"}}, nil
		},
	}
	controller := NewTagController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tags/{id}/tasks", controller.GetTasksByTagID)

	req := httptest.NewRequest("GET", "/tags/1/tasks", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTagController_CreateTag_Success(t *testing.T) {
	mockService := &MockTagService{
		createFunc: func(userID int64, tag *models.Tag) (*models.Tag, error) {
			return &models.Tag{ID: 1, Name: tag.Name, UserID: userID}, nil
		},
	}
	controller := NewTagController(mockService)

	tagJSON := `{"name": "Urgent"}`
	req := httptest.NewRequest("POST", "/tags", bytes.NewBufferString(tagJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.CreateTag(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestTagController_CreateTag_InvalidRequest(t *testing.T) {
	mockService := &MockTagService{}
	controller := NewTagController(mockService)

	invalidJSON := `{"name": invalid}`
	req := httptest.NewRequest("POST", "/tags", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.CreateTag(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid request, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTagController_CreateTag_ServiceError(t *testing.T) {
	mockService := &MockTagService{
		createFunc: func(userID int64, tag *models.Tag) (*models.Tag, error) {
			return nil, errors.New("service error")
		},
	}
	controller := NewTagController(mockService)

	tagJSON := `{"name": "Urgent"}`
	req := httptest.NewRequest("POST", "/tags", bytes.NewBufferString(tagJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	controller.CreateTag(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for service error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTagController_UpdateTag_Success(t *testing.T) {
	mockService := &MockTagService{
		updateFunc: func(userID int64, id int64, tag *models.Tag) (*models.Tag, error) {
			return &models.Tag{ID: id, Name: tag.Name, UserID: userID}, nil
		},
	}
	controller := NewTagController(mockService)

	tagJSON := `{"name": "Updated Tag"}`
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tags/{id}", controller.UpdateTag)

	req := httptest.NewRequest("PUT", "/tags/1", bytes.NewBufferString(tagJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTagController_UpdateTag_NotFound(t *testing.T) {
	mockService := &MockTagService{
		updateFunc: func(userID int64, id int64, tag *models.Tag) (*models.Tag, error) {
			return nil, sql.ErrNoRows
		},
	}
	controller := NewTagController(mockService)

	tagJSON := `{"name": "Updated Tag"}`
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tags/{id}", controller.UpdateTag)

	req := httptest.NewRequest("PUT", "/tags/999", bytes.NewBufferString(tagJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
