package services

import (
	"database/sql"
	"errors"
	"taskmanager/models"
	"testing"
	"time"
)

// MockSubtaskModel is a mock implementation of SubtaskModelInterface for testing
type MockSubtaskModel struct {
	CreateFunc            func(subtask *models.Subtask) (*models.Subtask, error)
	GetByTaskIDFunc       func(taskID int64) ([]models.Subtask, error)
	GetByIDFunc           func(id int64) (*models.Subtask, error)
	UpdateFunc            func(subtask *models.Subtask) error
	DeleteFunc            func(id int64) error
	ToggleFunc            func(id int64) error
	CalculateProgressFunc func(taskID int64) (int, error)
}

// MockTaskModel is a mock implementation of TaskModelInterface for testing (minimal version for subtask tests)
type MockTaskModel struct {
	GetByIDFunc             func(userID int64, id int64) (*models.Task, error)
	UpdateProgressFunc      func(taskID int64) error
	GetAnalyticsSummaryFunc func(userID int64) (map[string]interface{}, error)
	BulkDeleteFunc          func(userID int64, taskIDs []int64) error
	BulkCompleteFunc        func(userID int64, taskIDs []int64) error
}

func (m *MockTaskModel) GetAnalyticsSummary(userID int64) (map[string]interface{}, error) {
	if m.GetAnalyticsSummaryFunc != nil {
		return m.GetAnalyticsSummaryFunc(userID)
	}
	return nil, nil
}

func (m *MockTaskModel) BulkDelete(userID int64, taskIDs []int64) error {
	if m.BulkDeleteFunc != nil {
		return m.BulkDeleteFunc(userID, taskIDs)
	}
	return nil
}

func (m *MockTaskModel) BulkComplete(userID int64, taskIDs []int64) error {
	if m.BulkCompleteFunc != nil {
		return m.BulkCompleteFunc(userID, taskIDs)
	}
	return nil
}

// Implement remaining required methods for TaskModelInterface
func (m *MockTaskModel) Create(task *models.Task) (*models.Task, error) { return nil, nil }
func (m *MockTaskModel) GetAll(userID int64, completed *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error) {
	return nil, 0, nil
}
func (m *MockTaskModel) Update(userID int64, task *models.Task) error         { return nil }
func (m *MockTaskModel) Complete(userID int64, id int64) error                { return nil }
func (m *MockTaskModel) Delete(userID int64, id int64, softDelete bool) error { return nil }
func (m *MockTaskModel) RestoreTask(userID int64, id int64) error             { return nil }
func (m *MockTaskModel) MarkTaskAsCompleted(userID int64, id int64) error     { return nil }
func (m *MockTaskModel) MarkTaskAsUncompleted(userID int64, id int64) error   { return nil }
func (m *MockTaskModel) AddTagToTask(taskID int64, tagID int64) error         { return nil }
func (m *MockTaskModel) RemoveTagFromTask(taskID int64, tagID int64) error    { return nil }
func (m *MockTaskModel) LoadTags(task *models.Task) error                     { return nil }
func (m *MockTaskModel) LoadSubtasks(task *models.Task) error                 { return nil }
func (m *MockTaskModel) UpdateProgress(taskID int64) error {
	if m.UpdateProgressFunc != nil {
		return m.UpdateProgressFunc(taskID)
	}
	return nil
}
func (m *MockTaskModel) GetByID(userID int64, id int64) (*models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(userID, id)
	}
	return nil, nil
}

func (m *MockSubtaskModel) Create(subtask *models.Subtask) (*models.Subtask, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(subtask)
	}
	return nil, nil
}

func (m *MockSubtaskModel) GetByTaskID(taskID int64) ([]models.Subtask, error) {
	if m.GetByTaskIDFunc != nil {
		return m.GetByTaskIDFunc(taskID)
	}
	return nil, nil
}

func (m *MockSubtaskModel) GetByID(id int64) (*models.Subtask, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockSubtaskModel) Update(subtask *models.Subtask) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(subtask)
	}
	return nil
}

func (m *MockSubtaskModel) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockSubtaskModel) Toggle(id int64) error {
	if m.ToggleFunc != nil {
		return m.ToggleFunc(id)
	}
	return nil
}

func (m *MockSubtaskModel) CalculateProgress(taskID int64) (int, error) {
	if m.CalculateProgressFunc != nil {
		return m.CalculateProgressFunc(taskID)
	}
	return 0, nil
}

func TestCreateSubtask_Success(t *testing.T) {
	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
		UpdateProgressFunc: func(taskID int64) error {
			return nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		CreateFunc: func(subtask *models.Subtask) (*models.Subtask, error) {
			subtask.ID = 1
			subtask.CreatedAt = time.Now()
			subtask.UpdatedAt = time.Now()
			return subtask, nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	subtask := &models.Subtask{
		Title: "Test Subtask",
	}

	result, err := service.Create(1, 1, subtask)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}

	if result.IsCompleted != false {
		t.Errorf("Expected IsCompleted to be false, got %v", result.IsCompleted)
	}
}

func TestCreateSubtask_EmptyTitle(t *testing.T) {
	mockTaskModel := &MockTaskModel{}
	mockSubtaskModel := &MockSubtaskModel{}
	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	subtask := &models.Subtask{
		Title: "",
	}

	_, err := service.Create(1, 1, subtask)

	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}

	if err.Error() != "Subtask title cannot be empty" {
		t.Errorf("Expected 'Subtask title cannot be empty' error, got %v", err)
	}
}

func TestCreateSubtask_TaskNotFound(t *testing.T) {
	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockSubtaskModel := &MockSubtaskModel{}
	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	subtask := &models.Subtask{
		Title: "Test Subtask",
	}

	_, err := service.Create(1, 999, subtask)

	if err == nil {
		t.Error("Expected error for task not found, got nil")
	}

	if err.Error() != "Task not found or access denied" {
		t.Errorf("Expected 'Task not found or access denied' error, got %v", err)
	}
}

func TestGetSubtasksByTaskID_Success(t *testing.T) {
	mockSubtasks := []models.Subtask{
		{ID: 1, TaskID: 1, Title: "Subtask 1", IsCompleted: false},
		{ID: 2, TaskID: 1, Title: "Subtask 2", IsCompleted: true},
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByTaskIDFunc: func(taskID int64) ([]models.Subtask, error) {
			return mockSubtasks, nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	subtasks, err := service.GetByTaskID(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(subtasks) != 2 {
		t.Errorf("Expected 2 subtasks, got %d", len(subtasks))
	}
}

func TestGetSubtasksByTaskID_TaskNotFound(t *testing.T) {
	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockSubtaskModel := &MockSubtaskModel{}
	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	_, err := service.GetByTaskID(1, 999)

	if err == nil {
		t.Error("Expected error for task not found, got nil")
	}

	if err.Error() != "Task not found or access denied" {
		t.Errorf("Expected 'Task not found or access denied' error, got %v", err)
	}
}

func TestGetSubtaskByID_Success(t *testing.T) {
	mockSubtask := &models.Subtask{
		ID:          1,
		TaskID:      1,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return mockSubtask, nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	subtask, err := service.GetByID(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if subtask.ID != 1 {
		t.Errorf("Expected ID 1, got %d", subtask.ID)
	}
}

func TestGetSubtaskByID_AccessDenied(t *testing.T) {
	mockSubtask := &models.Subtask{
		ID:     1,
		TaskID: 1,
		Title:  "Test Subtask",
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return mockSubtask, nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	_, err := service.GetByID(2, 1)

	if err == nil {
		t.Error("Expected error for access denied, got nil")
	}

	if err.Error() != "Access denied to subtask" {
		t.Errorf("Expected 'Access denied to subtask' error, got %v", err)
	}
}

func TestUpdateSubtask_Success(t *testing.T) {
	existingSubtask := &models.Subtask{
		ID:     1,
		TaskID: 1,
		Title:  "Old Title",
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
		UpdateProgressFunc: func(taskID int64) error {
			return nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return existingSubtask, nil
		},
		UpdateFunc: func(subtask *models.Subtask) error {
			return nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	updatedSubtask := &models.Subtask{
		Title: "Updated Title",
	}

	result, err := service.Update(1, 1, updatedSubtask)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", result.Title)
	}
}

func TestUpdateSubtask_EmptyTitle(t *testing.T) {
	mockTaskModel := &MockTaskModel{}
	mockSubtaskModel := &MockSubtaskModel{}
	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	updatedSubtask := &models.Subtask{
		Title: "",
	}

	_, err := service.Update(1, 1, updatedSubtask)

	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}

	if err.Error() != "Subtask title cannot be empty" {
		t.Errorf("Expected 'Subtask title cannot be empty' error, got %v", err)
	}
}

func TestDeleteSubtask_Success(t *testing.T) {
	existingSubtask := &models.Subtask{
		ID:     1,
		TaskID: 1,
		Title:  "Test Subtask",
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
		UpdateProgressFunc: func(taskID int64) error {
			return nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return existingSubtask, nil
		},
		DeleteFunc: func(id int64) error {
			return nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	err := service.Delete(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteSubtask_AccessDenied(t *testing.T) {
	existingSubtask := &models.Subtask{
		ID:     1,
		TaskID: 1,
		Title:  "Test Subtask",
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return existingSubtask, nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	err := service.Delete(2, 1)

	if err == nil {
		t.Error("Expected error for access denied, got nil")
	}

	if err.Error() != "Access denied to subtask" {
		t.Errorf("Expected 'Access denied to subtask' error, got %v", err)
	}
}

func TestToggleSubtask_Success(t *testing.T) {
	existingSubtask := &models.Subtask{
		ID:          1,
		TaskID:      1,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
		UpdateProgressFunc: func(taskID int64) error {
			return nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return existingSubtask, nil
		},
		ToggleFunc: func(id int64) error {
			return nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	err := service.Toggle(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestToggleSubtask_AccessDenied(t *testing.T) {
	existingSubtask := &models.Subtask{
		ID:     1,
		TaskID: 1,
		Title:  "Test Subtask",
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return existingSubtask, nil
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	err := service.Toggle(2, 1)

	if err == nil {
		t.Error("Expected error for access denied, got nil")
	}

	if err.Error() != "Access denied to subtask" {
		t.Errorf("Expected 'Access denied to subtask' error, got %v", err)
	}
}

func TestToggleSubtask_ServiceError(t *testing.T) {
	existingSubtask := &models.Subtask{
		ID:     1,
		TaskID: 1,
		Title:  "Test Subtask",
	}

	mockTaskModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return &models.Task{ID: 1, UserID: 1}, nil
		},
	}

	mockSubtaskModel := &MockSubtaskModel{
		GetByIDFunc: func(id int64) (*models.Subtask, error) {
			return existingSubtask, nil
		},
		ToggleFunc: func(id int64) error {
			return errors.New("toggle failed")
		},
	}

	service := NewSubtaskService(mockSubtaskModel, mockTaskModel)

	err := service.Toggle(1, 1)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "toggle failed" {
		t.Errorf("Expected 'toggle failed' error, got %v", err)
	}
}
