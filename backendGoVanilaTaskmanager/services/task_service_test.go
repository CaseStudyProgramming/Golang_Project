package services

import (
	"database/sql"
	"errors"
	"strings"
	"taskmanager/models"
	"taskmanager/utils"
	"testing"
)

// MockTaskModel is a mock implementation of TaskModelInterface for testing
type MockTaskModel struct {
	CreateFunc                func(task *models.Task) (*models.Task, error)
	GetAllFunc                func(userID int64, completed *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error)
	GetByIDFunc               func(userID int64, id int64) (*models.Task, error)
	UpdateFunc                func(userID int64, task *models.Task) error
	CompleteFunc              func(userID int64, id int64) error
	DeleteFunc                func(userID int64, id int64, softDelete bool) error
	RestoreTaskFunc           func(userID int64, id int64) error
	MarkTaskAsCompletedFunc   func(userID int64, id int64) error
	MarkTaskAsUncompletedFunc func(userID int64, id int64) error
	AddTagToTaskFunc          func(taskID int64, tagID int64) error
	RemoveTagFromTaskFunc     func(taskID int64, tagID int64) error
	LoadTagsFunc              func(task *models.Task) error
	LoadSubtasksFunc          func(task *models.Task) error
	UpdateProgressFunc        func(taskID int64) error
	GetAnalyticsSummaryFunc   func(userID int64) (map[string]interface{}, error)
	BulkDeleteFunc            func(userID int64, taskIDs []int64) error
	BulkCompleteFunc          func(userID int64, taskIDs []int64) error
}

// Ensure MockTaskModel implements TaskModelInterface
var _ models.TaskModelInterface = (*MockTaskModel)(nil)

func (m *MockTaskModel) Create(task *models.Task) (*models.Task, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(task)
	}
	return nil, nil
}

func (m *MockTaskModel) GetAll(userID int64, completed *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(userID, completed, offset, limit, search, priority, categoryID, sortBy, sortOrder)
	}
	return nil, 0, nil
}

func (m *MockTaskModel) GetByID(userID int64, id int64) (*models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(userID, id)
	}
	return nil, nil
}

func (m *MockTaskModel) Update(userID int64, task *models.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(userID, task)
	}
	return nil
}

func (m *MockTaskModel) Complete(userID int64, id int64) error {
	if m.CompleteFunc != nil {
		return m.CompleteFunc(userID, id)
	}
	return nil
}

func (m *MockTaskModel) Delete(userID int64, id int64, softDelete bool) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(userID, id, softDelete)
	}
	return nil
}

func (m *MockTaskModel) RestoreTask(userID int64, id int64) error {
	if m.RestoreTaskFunc != nil {
		return m.RestoreTaskFunc(userID, id)
	}
	return nil
}

func (m *MockTaskModel) MarkTaskAsCompleted(userID int64, id int64) error {
	if m.MarkTaskAsCompletedFunc != nil {
		return m.MarkTaskAsCompletedFunc(userID, id)
	}
	return nil
}

func (m *MockTaskModel) MarkTaskAsUncompleted(userID int64, id int64) error {
	if m.MarkTaskAsUncompletedFunc != nil {
		return m.MarkTaskAsUncompletedFunc(userID, id)
	}
	return nil
}

func (m *MockTaskModel) AddTagToTask(taskID int64, tagID int64) error {
	if m.AddTagToTaskFunc != nil {
		return m.AddTagToTaskFunc(taskID, tagID)
	}
	return nil
}

func (m *MockTaskModel) RemoveTagFromTask(taskID int64, tagID int64) error {
	if m.RemoveTagFromTaskFunc != nil {
		return m.RemoveTagFromTaskFunc(taskID, tagID)
	}
	return nil
}

func (m *MockTaskModel) LoadTags(task *models.Task) error {
	if m.LoadTagsFunc != nil {
		return m.LoadTagsFunc(task)
	}
	return nil
}

func (m *MockTaskModel) LoadSubtasks(task *models.Task) error {
	if m.LoadSubtasksFunc != nil {
		return m.LoadSubtasksFunc(task)
	}
	return nil
}

func (m *MockTaskModel) UpdateProgress(taskID int64) error {
	if m.UpdateProgressFunc != nil {
		return m.UpdateProgressFunc(taskID)
	}
	return nil
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

func TestCreateTask_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
			task.ID = 1
			task.CreatedAt = utils.CurrentEpochMillis()
			task.UpdatedAt = utils.CurrentEpochMillis()
			return task, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	result, err := service.Create(1, task, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}

	if result.Completed != false {
		t.Errorf("Expected Completed to be false, got %v", result.Completed)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	mockModel := &MockTaskModel{}
	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	task := &models.Task{
		Title: "",
	}

	_, err := service.Create(1, task, "127.0.0.1", "test-agent")

	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}

	if err != utils.ErrMissingRequired {
		t.Errorf("Expected ErrMissingRequired error, got %v", err)
	}
}

func TestCreateTask_PastDueDate(t *testing.T) {
	mockModel := &MockTaskModel{}
	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	pastTime := utils.CurrentEpochMillis() - (24 * 60 * 60 * 1000) // 24 hours ago in milliseconds
	task := &models.Task{
		Title:   "Test Task",
		DueDate: &pastTime,
	}

	_, err := service.Create(1, task, "127.0.0.1", "test-agent")

	if err == nil {
		t.Error("Expected error for past due date, got nil")
	}

	if err.Error() != "Due date must be equal or greater than current date" {
		t.Errorf("Expected due date validation error, got %v", err)
	}
}

func TestCreateTask_FutureDueDate(t *testing.T) {
	mockModel := &MockTaskModel{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
			task.ID = 1
			task.CreatedAt = utils.CurrentEpochMillis()
			task.UpdatedAt = utils.CurrentEpochMillis()
			return task, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	futureTime := utils.CurrentEpochMillis() + (24 * 60 * 60 * 1000) // 24 hours in the future in milliseconds
	task := &models.Task{
		Title:   "Test Task",
		DueDate: &futureTime,
	}

	result, err := service.Create(1, task, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error for future due date, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}
}

func TestGetAllTasks_Success(t *testing.T) {
	mockTasks := []models.Task{
		{ID: 1, Title: "Task 1", Completed: false},
		{ID: 2, Title: "Task 2", Completed: true},
	}

	mockModel := &MockTaskModel{
		GetAllFunc: func(userID int64, completed *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error) {
			return mockTasks, 2, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	tasks, meta, err := service.GetAll(1, nil, 1, 10, "", nil, nil, "", "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if meta["total_data"] != 2 {
		t.Errorf("Expected total_data 2, got %v", meta["total_data"])
	}
}

func TestGetAllTasks_WithFilter(t *testing.T) {
	completed := true
	mockTasks := []models.Task{
		{ID: 1, Title: "Task 1", Completed: true},
	}

	mockModel := &MockTaskModel{
		GetAllFunc: func(userID int64, completedFilter *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error) {
			if completedFilter != nil && *completedFilter == true {
				return mockTasks, 1, nil
			}
			return []models.Task{}, 0, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	tasks, meta, err := service.GetAll(1, &completed, 1, 10, "", nil, nil, "", "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	if meta["total_data"] != 1 {
		t.Errorf("Expected total_data 1, got %v", meta["total_data"])
	}
}

func TestGetAllTasks_InvalidPage(t *testing.T) {
	mockModel := &MockTaskModel{
		GetAllFunc: func(userID int64, completed *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error) {
			// Return 5 total items with limit 5, so only 1 page exists
			mockTasks := []models.Task{
				{ID: 1, Title: "Task 1", Completed: false},
				{ID: 2, Title: "Task 2", Completed: false},
				{ID: 3, Title: "Task 3", Completed: false},
				{ID: 4, Title: "Task 4", Completed: false},
				{ID: 5, Title: "Task 5", Completed: false},
			}
			return mockTasks, 5, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	// Request page 2 when only 1 page exists
	_, _, err := service.GetAll(1, nil, 2, 5, "", nil, nil, "", "")

	if err == nil {
		t.Error("Expected error for invalid page, got nil")
	}

	expectedError := "page must be less than or equal to 1"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected '%s' error, got %v", expectedError, err)
	}
}

func TestGetAllTasks_NoResults(t *testing.T) {
	mockModel := &MockTaskModel{
		GetAllFunc: func(userID int64, completed *bool, offset, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, int, error) {
			return []models.Task{}, 0, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	_, _, err := service.GetAll(1, nil, 1, 10, "nonexistent", nil, nil, "", "")

	if err == nil {
		t.Error("Expected error for no results, got nil")
	}

	if err.Error() != "no tasks found" {
		t.Errorf("Expected 'no tasks found' error, got %v", err)
	}
}

func TestGetByID_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:        1,
		Title:     "Test Task",
		Completed: false,
	}

	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			if id == 1 {
				return mockTask, nil
			}
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	task, err := service.GetByID(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	_, err := service.GetByID(1, 999)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:        1,
		Title:     "Updated Task",
		Completed: true,
	}

	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return mockTask, nil
		},
		UpdateFunc: func(userID int64, task *models.Task) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	updatedTask := &models.Task{
		Title:     "Updated Task",
		Completed: true,
	}

	result, err := service.Update(1, 1, updatedTask, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}
}

func TestUpdate_PastDueDate(t *testing.T) {
	mockModel := &MockTaskModel{}
	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	pastTime := utils.CurrentEpochMillis() - (24 * 60 * 60 * 1000) // 24 hours ago in milliseconds
	updatedTask := &models.Task{
		Title:   "Updated Task",
		DueDate: &pastTime,
	}

	_, err := service.Update(1, 1, updatedTask, "127.0.0.1", "test-agent")

	if err == nil {
		t.Error("Expected error for past due date, got nil")
	}

	if err.Error() != "Due date must be equal or greater than current date" {
		t.Errorf("Expected due date validation error, got %v", err)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	updatedTask := &models.Task{
		Title: "Updated Task",
	}

	_, err := service.Update(1, 999, updatedTask, "127.0.0.1", "test-agent")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:        1,
		Title:     "Test Task",
		Completed: false,
	}

	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return mockTask, nil
		},
		DeleteFunc: func(userID int64, id int64, softDelete bool) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.Delete(1, 1, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.Delete(1, 999, "127.0.0.1", "test-agent")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestMarkAsCompleted_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		MarkTaskAsCompletedFunc: func(userID int64, id int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.MarkAsCompleted(1, 1, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMarkAsUncompleted_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		MarkTaskAsUncompletedFunc: func(userID int64, id int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.MarkAsUncompleted(1, 1, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRestore_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		RestoreTaskFunc: func(userID int64, id int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.Restore(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRestore_Error(t *testing.T) {
	mockModel := &MockTaskModel{
		RestoreTaskFunc: func(userID int64, id int64) error {
			return errors.New("restore failed")
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.Restore(1, 1)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Check if error contains the expected message (accounting for error wrapping)
	if !strings.Contains(err.Error(), "restore failed") {
		t.Errorf("Expected error to contain 'restore failed', got %v", err)
	}
}

func TestAddTagToTask_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:    1,
		Title: "Test Task",
	}

	mockTag := &models.Tag{
		ID:     1,
		UserID: 1,
		Name:   "Test Tag",
	}

	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			if id == 1 {
				return mockTask, nil
			}
			return nil, sql.ErrNoRows
		},
		AddTagToTaskFunc: func(taskID int64, tagID int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	// Create the tag in the mock with the correct user ID
	mockTagModel.Create(mockTag)

	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.AddTagToTask(1, 1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAddTagToTask_TaskNotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.AddTagToTask(1, 999, 1)

	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}
}

func TestAddTagToTask_TagNotFound(t *testing.T) {
	mockTask := &models.Task{
		ID:    1,
		Title: "Test Task",
	}

	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			if id == 1 {
				return mockTask, nil
			}
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	// Set the mock to return error for any GetByID call
	mockTagModel.setGetByIDError(sql.ErrNoRows)

	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.AddTagToTask(1, 1, 999)

	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}
}

func TestRemoveTagFromTask_Success(t *testing.T) {
	mockTask := &models.Task{
		ID:    1,
		Title: "Test Task",
	}

	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			if id == 1 {
				return mockTask, nil
			}
			return nil, sql.ErrNoRows
		},
		RemoveTagFromTaskFunc: func(taskID int64, tagID int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.RemoveTagFromTask(1, 1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRemoveTagFromTask_TaskNotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(userID int64, id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.RemoveTagFromTask(1, 999, 1)

	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}
}

func TestGetAnalyticsSummary_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		GetAnalyticsSummaryFunc: func(userID int64) (map[string]interface{}, error) {
			return map[string]interface{}{
				"total_active":    5,
				"total_completed": 3,
				"priority_distribution": map[string]int{
					"HIGH":   2,
					"MEDIUM": 3,
				},
			}, nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	summary, err := service.GetAnalyticsSummary(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if summary["total_active"] != 5 {
		t.Errorf("Expected total_active 5, got %v", summary["total_active"])
	}
}

func TestGetAnalyticsSummary_Error(t *testing.T) {
	mockModel := &MockTaskModel{
		GetAnalyticsSummaryFunc: func(userID int64) (map[string]interface{}, error) {
			return nil, errors.New("analytics error")
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	_, err := service.GetAnalyticsSummary(1)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestBulkDelete_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		BulkDeleteFunc: func(userID int64, taskIDs []int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	taskIDs := []int64{1, 2, 3}
	err := service.BulkDelete(1, taskIDs, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBulkDelete_EmptyTaskIDs(t *testing.T) {
	mockModel := &MockTaskModel{}
	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.BulkDelete(1, []int64{}, "127.0.0.1", "test-agent")

	// Bulk operations with empty arrays should succeed (no-op)
	if err != nil {
		t.Errorf("Expected no error for empty task IDs, got %v", err)
	}
}

func TestBulkComplete_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		BulkCompleteFunc: func(userID int64, taskIDs []int64) error {
			return nil
		},
	}

	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	taskIDs := []int64{1, 2, 3}
	err := service.BulkComplete(1, taskIDs, "127.0.0.1", "test-agent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBulkComplete_EmptyTaskIDs(t *testing.T) {
	mockModel := &MockTaskModel{}
	mockTagModel := NewMockTagModel()
	service := NewTaskService(mockModel, mockTagModel, nil)

	err := service.BulkComplete(1, []int64{}, "127.0.0.1", "test-agent")

	// Bulk operations with empty arrays should succeed (no-op)
	if err != nil {
		t.Errorf("Expected no error for empty task IDs, got %v", err)
	}
}
