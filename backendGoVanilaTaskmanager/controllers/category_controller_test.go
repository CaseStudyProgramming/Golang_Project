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

// MockCategoryService is a mock implementation of CategoryServiceInterface for testing
type MockCategoryService struct {
	createFunc  func(userID int64, category *models.Category) (*models.Category, error)
	getAllFunc  func(userID int64) ([]models.Category, error)
	getByIDFunc func(userID int64, id int64) (*models.Category, error)
	updateFunc  func(userID int64, id int64, category *models.Category) (*models.Category, error)
	deleteFunc  func(userID int64, id int64) error
}

func (m *MockCategoryService) Create(userID int64, category *models.Category) (*models.Category, error) {
	if m.createFunc != nil {
		return m.createFunc(userID, category)
	}
	return &models.Category{ID: 1, Name: category.Name, UserID: userID}, nil
}

func (m *MockCategoryService) GetAll(userID int64) ([]models.Category, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(userID)
	}
	return []models.Category{}, nil
}

func (m *MockCategoryService) GetByID(userID int64, id int64) (*models.Category, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(userID, id)
	}
	return &models.Category{ID: id, Name: "Test", UserID: userID}, nil
}

func (m *MockCategoryService) Update(userID int64, id int64, category *models.Category) (*models.Category, error) {
	if m.updateFunc != nil {
		return m.updateFunc(userID, id, category)
	}
	return &models.Category{ID: id, Name: category.Name, UserID: userID}, nil
}

func (m *MockCategoryService) Delete(userID int64, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(userID, id)
	}
	return nil
}

func TestNewCategoryController(t *testing.T) {
	mockService := &MockCategoryService{}
	controller := NewCategoryController(mockService)

	if controller == nil {
		t.Fatal("NewCategoryController returned nil")
	}

	if controller.service != mockService {
		t.Error("NewCategoryController did not set service correctly")
	}
}

func TestCategoryController_CreateCategory_Basic(t *testing.T) {
	mockService := &MockCategoryService{
		createFunc: func(userID int64, category *models.Category) (*models.Category, error) {
			return &models.Category{ID: 1, Name: category.Name, UserID: userID}, nil
		},
	}
	controller := NewCategoryController(mockService)

	categoryJSON := `{"name": "Work"}`
	req := httptest.NewRequest("POST", "/categories", bytes.NewBufferString(categoryJSON))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.CreateCategory(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected success status, got %v", response["status"])
	}
}

func TestCategoryController_GetAllCategories_Success(t *testing.T) {
	mockService := &MockCategoryService{}
	controller := NewCategoryController(mockService)

	mockService.getAllFunc = func(userID int64) ([]models.Category, error) {
		return []models.Category{{ID: 1, Name: "Work"}}, nil
	}

	req := httptest.NewRequest("GET", "/categories", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetAllCategories(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryController_GetAllCategories_Error(t *testing.T) {
	mockService := &MockCategoryService{}
	controller := NewCategoryController(mockService)

	mockService.getAllFunc = func(userID int64) ([]models.Category, error) {
		return nil, errors.New("service error")
	}

	req := httptest.NewRequest("GET", "/categories", nil)
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	ctx = context.WithValue(ctx, "timezone", "UTC")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	controller.GetAllCategories(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for service error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCategoryController_GetCategoryByID_Success(t *testing.T) {
	mockService := &MockCategoryService{
		getByIDFunc: func(userID int64, id int64) (*models.Category, error) {
			return &models.Category{ID: id, Name: "Work", UserID: userID}, nil
		},
	}
	controller := NewCategoryController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{id}", controller.GetCategoryByID)

	req := httptest.NewRequest("GET", "/categories/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryController_GetCategoryByID_NotFound(t *testing.T) {
	mockService := &MockCategoryService{
		getByIDFunc: func(userID int64, id int64) (*models.Category, error) {
			return nil, sql.ErrNoRows
		},
	}
	controller := NewCategoryController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{id}", controller.GetCategoryByID)

	req := httptest.NewRequest("GET", "/categories/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCategoryController_UpdateCategory_Success(t *testing.T) {
	mockService := &MockCategoryService{
		updateFunc: func(userID int64, id int64, category *models.Category) (*models.Category, error) {
			return &models.Category{ID: id, Name: category.Name, UserID: userID}, nil
		},
	}
	controller := NewCategoryController(mockService)

	categoryJSON := `{"name": "Updated Category"}`
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", controller.UpdateCategory)

	req := httptest.NewRequest("PUT", "/categories/1", bytes.NewBufferString(categoryJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryController_UpdateCategory_NotFound(t *testing.T) {
	mockService := &MockCategoryService{
		updateFunc: func(userID int64, id int64, category *models.Category) (*models.Category, error) {
			return nil, sql.ErrNoRows
		},
	}
	controller := NewCategoryController(mockService)

	categoryJSON := `{"name": "Updated Category"}`
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /categories/{id}", controller.UpdateCategory)

	req := httptest.NewRequest("PUT", "/categories/999", bytes.NewBufferString(categoryJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCategoryController_DeleteCategory_Success(t *testing.T) {
	mockService := &MockCategoryService{
		deleteFunc: func(userID int64, id int64) error {
			return nil
		},
	}
	controller := NewCategoryController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /categories/{id}", controller.DeleteCategory)

	req := httptest.NewRequest("DELETE", "/categories/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status %d, got %d", http.StatusAccepted, w.Code)
	}
}

func TestCategoryController_DeleteCategory_NotFound(t *testing.T) {
	mockService := &MockCategoryService{
		deleteFunc: func(userID int64, id int64) error {
			return sql.ErrNoRows
		},
	}
	controller := NewCategoryController(mockService)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /categories/{id}", controller.DeleteCategory)

	req := httptest.NewRequest("DELETE", "/categories/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", int64(1)))
	req = req.WithContext(context.WithValue(req.Context(), "timezone", "UTC"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
