package services

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	"taskmanager/models"
)

// MockTaskModel is a mock implementation of TaskModelInterface for testing
type MockTaskModel struct {
	CreateFunc              func(task *models.Task) (*models.Task, error)
	GetAllFunc              func(completed *bool, offset, limit int, search string) ([]models.Task, int, error)
	GetByIDFunc             func(id int64) (*models.Task, error)
	UpdateFunc              func(task *models.Task) error
	CompleteFunc            func(id int64) error
	DeleteFunc              func(id int64, softDelete bool) error
	RestoreTaskFunc         func(id int64) error
	MarkTaskAsCompletedFunc func(id int64) error
	MarkTaskAsUncompletedFunc func(id int64) error
}

// Ensure MockTaskModel implements TaskModelInterface
var _ models.TaskModelInterface = (*MockTaskModel)(nil)

func (m *MockTaskModel) Create(task *models.Task) (*models.Task, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(task)
	}
	return nil, nil
}

func (m *MockTaskModel) GetAll(completed *bool, offset, limit int, search string) ([]models.Task, int, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(completed, offset, limit, search)
	}
	return nil, 0, nil
}

func (m *MockTaskModel) GetByID(id int64) (*models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockTaskModel) Update(task *models.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(task)
	}
	return nil
}

func (m *MockTaskModel) Complete(id int64) error {
	if m.CompleteFunc != nil {
		return m.CompleteFunc(id)
	}
	return nil
}

func (m *MockTaskModel) Delete(id int64, softDelete bool) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id, softDelete)
	}
	return nil
}

func (m *MockTaskModel) RestoreTask(id int64) error {
	if m.RestoreTaskFunc != nil {
		return m.RestoreTaskFunc(id)
	}
	return nil
}

func (m *MockTaskModel) MarkTaskAsCompleted(id int64) error {
	if m.MarkTaskAsCompletedFunc != nil {
		return m.MarkTaskAsCompletedFunc(id)
	}
	return nil
}

func (m *MockTaskModel) MarkTaskAsUncompleted(id int64) error {
	if m.MarkTaskAsUncompletedFunc != nil {
		return m.MarkTaskAsUncompletedFunc(id)
	}
	return nil
}

func TestCreateTask_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
			task.ID = 1
			task.CreatedAt = time.Now()
			task.UpdatedAt = time.Now()
			return task, nil
		},
	}

	service := NewTaskService(mockModel)
	
	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	result, err := service.Create(task)
	
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
	service := NewTaskService(mockModel)
	
	task := &models.Task{
		Title: "",
	}

	_, err := service.Create(task)
	
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
	
	if err.Error() != "Title tidak boleh kosong" {
		t.Errorf("Expected 'Title tidak boleh kosong' error, got %v", err)
	}
}

func TestCreateTask_PastDueDate(t *testing.T) {
	mockModel := &MockTaskModel{}
	service := NewTaskService(mockModel)
	
	pastTime := time.Now().Add(-24 * time.Hour)
	task := &models.Task{
		Title:   "Test Task",
		DueDate: &pastTime,
	}

	_, err := service.Create(task)
	
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
			task.CreatedAt = time.Now()
			task.UpdatedAt = time.Now()
			return task, nil
		},
	}

	service := NewTaskService(mockModel)
	
	futureTime := time.Now().Add(24 * time.Hour)
	task := &models.Task{
		Title:   "Test Task",
		DueDate: &futureTime,
	}

	result, err := service.Create(task)
	
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
		GetAllFunc: func(completed *bool, offset, limit int, search string) ([]models.Task, int, error) {
			return mockTasks, 2, nil
		},
	}

	service := NewTaskService(mockModel)
	
	tasks, meta, err := service.GetAll(nil, 1, 10, "")
	
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
		GetAllFunc: func(completedFilter *bool, offset, limit int, search string) ([]models.Task, int, error) {
			if completedFilter != nil && *completedFilter == true {
				return mockTasks, 1, nil
			}
			return []models.Task{}, 0, nil
		},
	}

	service := NewTaskService(mockModel)
	
	tasks, meta, err := service.GetAll(&completed, 1, 10, "")
	
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
		GetAllFunc: func(completed *bool, offset, limit int, search string) ([]models.Task, int, error) {
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

	service := NewTaskService(mockModel)
	
	// Request page 2 when only 1 page exists
	_, _, err := service.GetAll(nil, 2, 5, "")
	
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
		GetAllFunc: func(completed *bool, offset, limit int, search string) ([]models.Task, int, error) {
			return []models.Task{}, 0, nil
		},
	}

	service := NewTaskService(mockModel)
	
	_, _, err := service.GetAll(nil, 1, 10, "nonexistent")
	
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
		GetByIDFunc: func(id int64) (*models.Task, error) {
			if id == 1 {
				return mockTask, nil
			}
			return nil, sql.ErrNoRows
		},
	}

	service := NewTaskService(mockModel)
	
	task, err := service.GetByID(1)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewTaskService(mockModel)
	
	_, err := service.GetByID(999)
	
	if err != sql.ErrNoRows {
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
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return mockTask, nil
		},
		UpdateFunc: func(task *models.Task) error {
			return nil
		},
	}

	service := NewTaskService(mockModel)
	
	updatedTask := &models.Task{
		Title:     "Updated Task",
		Completed: true,
	}
	
	result, err := service.Update(1, updatedTask)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}
}

func TestUpdate_PastDueDate(t *testing.T) {
	mockModel := &MockTaskModel{}
	service := NewTaskService(mockModel)
	
	pastTime := time.Now().Add(-24 * time.Hour)
	updatedTask := &models.Task{
		Title:   "Updated Task",
		DueDate: &pastTime,
	}
	
	_, err := service.Update(1, updatedTask)
	
	if err == nil {
		t.Error("Expected error for past due date, got nil")
	}
	
	if err.Error() != "Due date must be equal or greater than current date" {
		t.Errorf("Expected due date validation error, got %v", err)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewTaskService(mockModel)
	
	updatedTask := &models.Task{
		Title: "Updated Task",
	}
	
	_, err := service.Update(999, updatedTask)
	
	if err != sql.ErrNoRows {
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
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return mockTask, nil
		},
		DeleteFunc: func(id int64, softDelete bool) error {
			return nil
		},
	}

	service := NewTaskService(mockModel)
	
	err := service.Delete(1)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	mockModel := &MockTaskModel{
		GetByIDFunc: func(id int64) (*models.Task, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := NewTaskService(mockModel)
	
	err := service.Delete(999)
	
	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}

func TestMarkAsCompleted_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		MarkTaskAsCompletedFunc: func(id int64) error {
			return nil
		},
	}

	service := NewTaskService(mockModel)
	
	err := service.MarkAsCompleted(1)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMarkAsUncompleted_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		MarkTaskAsUncompletedFunc: func(id int64) error {
			return nil
		},
	}

	service := NewTaskService(mockModel)
	
	err := service.MarkAsUncompleted(1)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRestore_Success(t *testing.T) {
	mockModel := &MockTaskModel{
		RestoreTaskFunc: func(id int64) error {
			return nil
		},
	}

	service := NewTaskService(mockModel)
	
	err := service.Restore(1)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRestore_Error(t *testing.T) {
	mockModel := &MockTaskModel{
		RestoreTaskFunc: func(id int64) error {
			return errors.New("restore failed")
		},
	}

	service := NewTaskService(mockModel)
	
	err := service.Restore(1)
	
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	if err.Error() != "restore failed" {
		t.Errorf("Expected 'restore failed' error, got %v", err)
	}
}