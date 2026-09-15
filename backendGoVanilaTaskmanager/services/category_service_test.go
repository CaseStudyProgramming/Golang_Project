package services

import (
	"errors"
	"taskmanager/models"
	"taskmanager/utils"
	"testing"
)

func TestNewCategoryService(t *testing.T) {
	mockModel := NewMockCategoryModel()
	service := NewCategoryService(mockModel)

	if service == nil {
		t.Fatal("NewCategoryService returned nil")
	}

	if service.model != mockModel {
		t.Error("NewCategoryService did not set model correctly")
	}
}

func TestCategoryService_Create(t *testing.T) {
	mockModel := NewMockCategoryModel()
	service := NewCategoryService(mockModel)

	tests := []struct {
		name        string
		userID      int64
		category    *models.Category
		expectError bool
		errorType   error
	}{
		{
			name:   "successful creation",
			userID: 1,
			category: &models.Category{
				Name:     "Work",
				ColorHex: "#FF5733",
			},
			expectError: false,
		},
		{
			name:   "successful creation with default color",
			userID: 1,
			category: &models.Category{
				Name: "Personal",
			},
			expectError: false,
		},
		{
			name:   "missing name",
			userID: 1,
			category: &models.Category{
				ColorHex: "#FF5733",
			},
			expectError: true,
			errorType:   utils.ErrMissingRequired,
		},
		{
			name:        "empty category",
			userID:      1,
			category:    &models.Category{},
			expectError: true,
			errorType:   utils.ErrMissingRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(tt.userID, tt.category)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected result but got nil")
				}
				if result != nil && result.ID == 0 {
					t.Error("Expected category ID to be set")
				}
				if result != nil && result.UserID != tt.userID {
					t.Errorf("Expected userID %d, got %d", tt.userID, result.UserID)
				}
				if tt.category.ColorHex == "" && result.ColorHex != "#3B82F6" {
					t.Errorf("Expected default color #3B82F6, got %s", result.ColorHex)
				}
			}
		})
	}
}

func TestCategoryService_Create_ModelError(t *testing.T) {
	mockModel := NewMockCategoryModel()
	mockModel.createError = errors.New("database error")
	service := NewCategoryService(mockModel)

	category := &models.Category{
		Name:     "Work",
		ColorHex: "#FF5733",
	}

	_, err := service.Create(1, category)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestCategoryService_GetAll(t *testing.T) {
	mockModel := NewMockCategoryModel()
	service := NewCategoryService(mockModel)

	// Setup test data
	userID := int64(1)
	category1 := &models.Category{UserID: userID, Name: "Work", ColorHex: "#FF5733"}
	category2 := &models.Category{UserID: userID, Name: "Personal", ColorHex: "#00FF00"}

	mockModel.Create(category1)
	mockModel.Create(category2)

	tests := []struct {
		name        string
		userID      int64
		expectError bool
		expectCount int
	}{
		{
			name:        "successful get all",
			userID:      userID,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "empty result for different user",
			userID:      999,
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetAll(tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectCount {
					t.Errorf("Expected %d categories, got %d", tt.expectCount, len(result))
				}
			}
		})
	}
}

func TestCategoryService_GetAll_ModelError(t *testing.T) {
	mockModel := NewMockCategoryModel()
	mockModel.getAllError = errors.New("database error")
	service := NewCategoryService(mockModel)

	_, err := service.GetAll(1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestCategoryService_GetByID(t *testing.T) {
	mockModel := NewMockCategoryModel()
	service := NewCategoryService(mockModel)

	// Setup test data
	userID := int64(1)
	category := &models.Category{UserID: userID, Name: "Work", ColorHex: "#FF5733"}
	createdCategory, _ := mockModel.Create(category)

	tests := []struct {
		name        string
		userID      int64
		id          int64
		expectError bool
	}{
		{
			name:        "successful get by id",
			userID:      userID,
			id:          createdCategory.ID,
			expectError: false,
		},
		{
			name:        "category not found",
			userID:      userID,
			id:          99999,
			expectError: true,
		},
		{
			name:        "wrong user",
			userID:      999,
			id:          createdCategory.ID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetByID(tt.userID, tt.id)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected result but got nil")
				}
				if result != nil && result.ID != tt.id {
					t.Errorf("Expected ID %d, got %d", tt.id, result.ID)
				}
			}
		})
	}
}

func TestCategoryService_Update(t *testing.T) {
	mockModel := NewMockCategoryModel()
	service := NewCategoryService(mockModel)

	// Setup test data
	userID := int64(1)
	category := &models.Category{UserID: userID, Name: "Work", ColorHex: "#FF5733"}
	createdCategory, _ := mockModel.Create(category)

	tests := []struct {
		name        string
		userID      int64
		id          int64
		category    *models.Category
		expectError bool
		errorType   error
	}{
		{
			name:   "successful update",
			userID: userID,
			id:     createdCategory.ID,
			category: &models.Category{
				Name:     "Updated Work",
				ColorHex: "#00FF00",
			},
			expectError: false,
		},
		{
			name:   "missing name",
			userID: userID,
			id:     createdCategory.ID,
			category: &models.Category{
				ColorHex: "#00FF00",
			},
			expectError: true,
			errorType:   utils.ErrMissingRequired,
		},
		{
			name:        "category not found",
			userID:      userID,
			id:          99999,
			category:    &models.Category{Name: "Test"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Update(tt.userID, tt.id, tt.category)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected result but got nil")
				}
				if result != nil && result.Name != tt.category.Name {
					t.Errorf("Expected name %s, got %s", tt.category.Name, result.Name)
				}
			}
		})
	}
}

func TestCategoryService_Delete(t *testing.T) {
	mockModel := NewMockCategoryModel()
	service := NewCategoryService(mockModel)

	// Setup test data
	userID := int64(1)
	category := &models.Category{UserID: userID, Name: "Work", ColorHex: "#FF5733"}
	createdCategory, _ := mockModel.Create(category)

	tests := []struct {
		name        string
		userID      int64
		id          int64
		expectError bool
	}{
		{
			name:        "successful delete",
			userID:      userID,
			id:          createdCategory.ID,
			expectError: false,
		},
		{
			name:        "category not found",
			userID:      userID,
			id:          99999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(tt.userID, tt.id)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Verify deletion
				_, err := mockModel.GetByID(tt.userID, tt.id)
				if err == nil {
					t.Error("Expected category to be deleted")
				}
			}
		})
	}
}
