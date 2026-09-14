package services

import (
	"errors"
	"taskmanager/models"
	"taskmanager/utils"
	"testing"
)

func TestNewTagService(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	if service == nil {
		t.Fatal("NewTagService returned nil")
	}

	if service.model != mockModel {
		t.Error("NewTagService did not set model correctly")
	}
}

func TestTagService_Create(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	tests := []struct {
		name        string
		userID      int64
		tag         *models.Tag
		expectError bool
		errorType   error
	}{
		{
			name:   "successful creation",
			userID: 1,
			tag: &models.Tag{
				Name:     "Urgent",
				ColorHex: "#FF5733",
			},
			expectError: false,
		},
		{
			name:   "successful creation with default color",
			userID: 1,
			tag: &models.Tag{
				Name: "Important",
			},
			expectError: false,
		},
		{
			name:   "missing name",
			userID: 1,
			tag: &models.Tag{
				ColorHex: "#FF5733",
			},
			expectError: true,
			errorType:   utils.ErrMissingRequired,
		},
		{
			name:        "empty tag",
			userID:      1,
			tag:         &models.Tag{},
			expectError: true,
			errorType:   utils.ErrMissingRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(tt.userID, tt.tag)

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
					t.Error("Expected tag ID to be set")
				}
				if result != nil && result.UserID != tt.userID {
					t.Errorf("Expected userID %d, got %d", tt.userID, result.UserID)
				}
				if tt.tag.ColorHex == "" && result.ColorHex != "#10B981" {
					t.Errorf("Expected default color #10B981, got %s", result.ColorHex)
				}
			}
		})
	}
}

func TestTagService_Create_ModelError(t *testing.T) {
	mockModel := NewMockTagModel()
	mockModel.createError = errors.New("database error")
	service := NewTagService(mockModel)

	tag := &models.Tag{
		Name:     "Urgent",
		ColorHex: "#FF5733",
	}

	_, err := service.Create(1, tag)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestTagService_GetAll(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	// Setup test data
	userID := int64(1)
	tag1 := &models.Tag{UserID: userID, Name: "Urgent", ColorHex: "#FF5733"}
	tag2 := &models.Tag{UserID: userID, Name: "Important", ColorHex: "#00FF00"}

	mockModel.Create(tag1)
	mockModel.Create(tag2)

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
					t.Errorf("Expected %d tags, got %d", tt.expectCount, len(result))
				}
			}
		})
	}
}

func TestTagService_GetAll_ModelError(t *testing.T) {
	mockModel := NewMockTagModel()
	mockModel.getAllError = errors.New("database error")
	service := NewTagService(mockModel)

	_, err := service.GetAll(1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestTagService_GetByID(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	// Setup test data
	userID := int64(1)
	tag := &models.Tag{UserID: userID, Name: "Urgent", ColorHex: "#FF5733"}
	createdTag, _ := mockModel.Create(tag)

	tests := []struct {
		name        string
		userID      int64
		id          int64
		expectError bool
	}{
		{
			name:        "successful get by id",
			userID:      userID,
			id:          createdTag.ID,
			expectError: false,
		},
		{
			name:        "tag not found",
			userID:      userID,
			id:          99999,
			expectError: true,
		},
		{
			name:        "wrong user",
			userID:      999,
			id:          createdTag.ID,
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

func TestTagService_Update(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	// Setup test data
	userID := int64(1)
	tag := &models.Tag{UserID: userID, Name: "Urgent", ColorHex: "#FF5733"}
	createdTag, _ := mockModel.Create(tag)

	tests := []struct {
		name        string
		userID      int64
		id          int64
		tag         *models.Tag
		expectError bool
		errorType   error
	}{
		{
			name:   "successful update",
			userID: userID,
			id:     createdTag.ID,
			tag: &models.Tag{
				Name:     "Updated Urgent",
				ColorHex: "#00FF00",
			},
			expectError: false,
		},
		{
			name:   "missing name",
			userID: userID,
			id:     createdTag.ID,
			tag: &models.Tag{
				ColorHex: "#00FF00",
			},
			expectError: true,
			errorType:   utils.ErrMissingRequired,
		},
		{
			name:        "tag not found",
			userID:      userID,
			id:          99999,
			tag:         &models.Tag{Name: "Test"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Update(tt.userID, tt.id, tt.tag)

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
				if result != nil && result.Name != tt.tag.Name {
					t.Errorf("Expected name %s, got %s", tt.tag.Name, result.Name)
				}
			}
		})
	}
}

func TestTagService_Delete(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	// Setup test data
	userID := int64(1)
	tag := &models.Tag{UserID: userID, Name: "Urgent", ColorHex: "#FF5733"}
	createdTag, _ := mockModel.Create(tag)

	tests := []struct {
		name        string
		userID      int64
		id          int64
		expectError bool
	}{
		{
			name:        "successful delete",
			userID:      userID,
			id:          createdTag.ID,
			expectError: false,
		},
		{
			name:        "tag not found",
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
					t.Error("Expected tag to be deleted")
				}
			}
		})
	}
}

func TestTagService_AddTagToTask(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	tests := []struct {
		name        string
		taskID      int64
		tagID       int64
		expectError bool
	}{
		{
			name:        "successful add tag to task",
			taskID:      1,
			tagID:       1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddTagToTask(tt.taskID, tt.tagID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestTagService_AddTagToTask_ModelError(t *testing.T) {
	mockModel := NewMockTagModel()
	mockModel.addTagError = errors.New("database error")
	service := NewTagService(mockModel)

	err := service.AddTagToTask(1, 1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestTagService_RemoveTagFromTask(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	// Setup test data
	taskID := int64(1)
	tagID := int64(1)
	mockModel.AddTagToTask(taskID, tagID)

	tests := []struct {
		name        string
		taskID      int64
		tagID       int64
		expectError bool
	}{
		{
			name:        "successful remove tag from task",
			taskID:      taskID,
			tagID:       tagID,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RemoveTagFromTask(tt.taskID, tt.tagID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestTagService_RemoveTagFromTask_ModelError(t *testing.T) {
	mockModel := NewMockTagModel()
	mockModel.removeTagError = errors.New("database error")
	service := NewTagService(mockModel)

	err := service.RemoveTagFromTask(1, 1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestTagService_GetTagsByTaskID(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	// Setup test data
	userID := int64(1)
	taskID := int64(1)
	tag1 := &models.Tag{UserID: userID, Name: "Urgent", ColorHex: "#FF5733"}
	tag2 := &models.Tag{UserID: userID, Name: "Important", ColorHex: "#00FF00"}

	createdTag1, _ := mockModel.Create(tag1)
	createdTag2, _ := mockModel.Create(tag2)

	mockModel.AddTagToTask(taskID, createdTag1.ID)
	mockModel.AddTagToTask(taskID, createdTag2.ID)

	tests := []struct {
		name        string
		taskID      int64
		expectError bool
		expectCount int
	}{
		{
			name:        "successful get tags by task id",
			taskID:      taskID,
			expectError: false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetTagsByTaskID(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectCount {
					t.Errorf("Expected %d tags, got %d", tt.expectCount, len(result))
				}
			}
		})
	}
}

func TestTagService_GetTagsByTaskID_ModelError(t *testing.T) {
	mockModel := NewMockTagModel()
	mockModel.getTagsError = errors.New("database error")
	service := NewTagService(mockModel)

	_, err := service.GetTagsByTaskID(1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestTagService_GetTasksByTagID(t *testing.T) {
	mockModel := NewMockTagModel()
	service := NewTagService(mockModel)

	tests := []struct {
		name        string
		tagID       int64
		expectError bool
	}{
		{
			name:        "successful get tasks by tag id",
			tagID:       1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetTasksByTagID(tt.tagID)

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
			}
		})
	}
}

func TestTagService_GetTasksByTagID_ModelError(t *testing.T) {
	mockModel := NewMockTagModel()
	mockModel.getTasksError = errors.New("database error")
	service := NewTagService(mockModel)

	_, err := service.GetTasksByTagID(1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}
