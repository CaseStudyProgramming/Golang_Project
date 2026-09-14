package services

import (
	"errors"
	"taskmanager/models"
	"testing"
)

func TestNewActivityLogService(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	service := NewActivityLogService(mockModel)

	if service == nil {
		t.Fatal("NewActivityLogService returned nil")
	}
}

func TestActivityLogService_LogActivity(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	service := NewActivityLogService(mockModel)

	tests := []struct {
		name        string
		userID      int64
		taskID      *int64
		action      string
		entityType  string
		entityID    *int64
		details     string
		ipAddress   string
		userAgent   string
		expectError bool
	}{
		{
			name:        "successful log activity",
			userID:      1,
			taskID:      func() *int64 { id := int64(1); return &id }(),
			action:      "create",
			entityType:  "task",
			entityID:    func() *int64 { id := int64(1); return &id }(),
			details:     "Created new task",
			ipAddress:   "127.0.0.1",
			userAgent:   "Mozilla/5.0",
			expectError: false,
		},
		{
			name:        "log activity without task",
			userID:      1,
			taskID:      nil,
			action:      "login",
			entityType:  "user",
			entityID:    func() *int64 { id := int64(1); return &id }(),
			details:     "User logged in",
			ipAddress:   "127.0.0.1",
			userAgent:   "Mozilla/5.0",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.LogActivity(tt.userID, tt.taskID, tt.action, tt.entityType, tt.entityID, tt.details, tt.ipAddress, tt.userAgent)

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

func TestActivityLogService_LogActivity_ModelError(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	mockModel.createError = errors.New("database error")
	service := NewActivityLogService(mockModel)

	err := service.LogActivity(1, nil, "login", "user", nil, "test", "127.0.0.1", "Mozilla")
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestActivityLogService_GetUserActivityLogs(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	service := NewActivityLogService(mockModel)

	// Setup test data
	userID := int64(1)
	taskID := int64(1)

	for i := 0; i < 5; i++ {
		log := &models.ActivityLog{
			UserID:     userID,
			TaskID:     &taskID,
			Action:     models.ActionType("create"),
			EntityType: models.EntityType("task"),
			EntityID:   &taskID,
			Details:    "Test log",
			IPAddress:  "127.0.0.1",
			UserAgent:  "Mozilla",
		}
		mockModel.Create(log)
	}

	tests := []struct {
		name        string
		userID      int64
		page        int
		limit       int
		expectError bool
		expectCount int
	}{
		{
			name:        "successful get user logs",
			userID:      userID,
			page:        1,
			limit:       10,
			expectError: false,
			expectCount: 5,
		},
		{
			name:        "pagination",
			userID:      userID,
			page:        1,
			limit:       2,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "invalid page",
			userID:      userID,
			page:        0,
			limit:       10,
			expectError: true,
		},
		{
			name:        "invalid limit",
			userID:      userID,
			page:        1,
			limit:       0,
			expectError: true,
		},
		{
			name:        "page exceeds total",
			userID:      userID,
			page:        10,
			limit:       10,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, meta, err := service.GetUserActivityLogs(tt.userID, tt.page, tt.limit)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectCount {
					t.Errorf("Expected %d logs, got %d", tt.expectCount, len(result))
				}
				if meta == nil {
					t.Error("Expected metadata but got nil")
				}
				if meta != nil {
					if meta["page"] != tt.page {
						t.Errorf("Expected page %d, got %v", tt.page, meta["page"])
					}
					if meta["limit"] != tt.limit {
						t.Errorf("Expected limit %d, got %v", tt.limit, meta["limit"])
					}
				}
			}
		})
	}
}

func TestActivityLogService_GetUserActivityLogs_ModelError(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	mockModel.getUserError = errors.New("database error")
	service := NewActivityLogService(mockModel)

	_, _, err := service.GetUserActivityLogs(1, 1, 10)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestActivityLogService_GetTaskActivityLogs(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	service := NewActivityLogService(mockModel)

	// Setup test data
	userID := int64(1)
	taskID := int64(1)

	for i := 0; i < 3; i++ {
		log := &models.ActivityLog{
			UserID:     userID,
			TaskID:     &taskID,
			Action:     models.ActionType("update"),
			EntityType: models.EntityType("task"),
			EntityID:   &taskID,
			Details:    "Test log",
			IPAddress:  "127.0.0.1",
			UserAgent:  "Mozilla",
		}
		mockModel.Create(log)
	}

	tests := []struct {
		name        string
		taskID      int64
		page        int
		limit       int
		expectError bool
		expectCount int
	}{
		{
			name:        "successful get task logs",
			taskID:      taskID,
			page:        1,
			limit:       10,
			expectError: false,
			expectCount: 3,
		},
		{
			name:        "pagination",
			taskID:      taskID,
			page:        1,
			limit:       2,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "invalid page",
			taskID:      taskID,
			page:        0,
			limit:       10,
			expectError: true,
		},
		{
			name:        "invalid limit",
			taskID:      taskID,
			page:        1,
			limit:       0,
			expectError: true,
		},
		{
			name:        "page exceeds total",
			taskID:      taskID,
			page:        10,
			limit:       10,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, meta, err := service.GetTaskActivityLogs(tt.taskID, tt.page, tt.limit)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expectCount {
					t.Errorf("Expected %d logs, got %d", tt.expectCount, len(result))
				}
				if meta == nil {
					t.Error("Expected metadata but got nil")
				}
				if meta != nil {
					if meta["page"] != tt.page {
						t.Errorf("Expected page %d, got %v", tt.page, meta["page"])
					}
					if meta["limit"] != tt.limit {
						t.Errorf("Expected limit %d, got %v", tt.limit, meta["limit"])
					}
				}
			}
		})
	}
}

func TestActivityLogService_GetTaskActivityLogs_ModelError(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	mockModel.getTaskError = errors.New("database error")
	service := NewActivityLogService(mockModel)

	_, _, err := service.GetTaskActivityLogs(1, 1, 10)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}

func TestActivityLogService_GetActivityLogByID(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	service := NewActivityLogService(mockModel)

	// Setup test data
	log := &models.ActivityLog{
		UserID:     1,
		Action:     models.ActionType("create"),
		EntityType: models.EntityType("task"),
		Details:    "Test log",
		IPAddress:  "127.0.0.1",
		UserAgent:  "Mozilla",
	}
	createdLog, _ := mockModel.Create(log)

	tests := []struct {
		name        string
		id          int64
		expectError bool
	}{
		{
			name:        "successful get by id",
			id:          createdLog.ID,
			expectError: false,
		},
		{
			name:        "log not found",
			id:          99999,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetActivityLogByID(tt.id)

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

func TestActivityLogService_GetActivityLogByID_ModelError(t *testing.T) {
	mockModel := NewMockActivityLogModel()
	mockModel.getByIDError = errors.New("database error")
	service := NewActivityLogService(mockModel)

	_, err := service.GetActivityLogByID(1)
	if err == nil {
		t.Error("Expected error from model failure")
	}
}
