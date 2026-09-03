package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskmanager/controllers"
	"taskmanager/models"
	"taskmanager/services"
	"taskmanager/utils"
	"testing"
	"time"
)

// Integration tests for the Task API endpoints
// These tests use a real database connection to test the complete request/response cycle

func TestTaskAPIIntegration_CreateAndGet(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "testuser", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "testuser").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	// Setup the application layers
	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	taskController := controllers.NewTaskController(taskService)

	// Create a test task
	taskJSON := `{"title": "Test Task Integration", "description": "Integration test task"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	w := httptest.NewRecorder()

	taskController.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	if createResponse["status"] != "success" {
		t.Errorf("Expected success status, got %v", createResponse["status"])
		return
	}

	// Extract the created task ID
	data := createResponse["data"].(map[string]interface{})
	taskID := int64(data["id"].(float64))

	// Get the created task using service directly (skip controller test for PathValue)
	getTask, err := taskService.GetByID(userID, taskID)
	if err != nil {
		t.Errorf("Failed to get created task: %v", err)
	}

	if getTask.Title != "Test Task Integration" {
		t.Errorf("Expected title 'Test Task Integration', got %v", getTask.Title)
	}
}

func TestTaskAPIIntegration_UpdateAndDelete(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "testuser2", "test2@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "testuser2").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)

	// Create a task first
	task := &models.Task{
		Title:       "Test Task for Update",
		Description: "Initial description",
	}
	createdTask, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Update the task using service directly
	updatedTask := &models.Task{
		Title:       "Updated Task Title",
		Description: "Updated description",
		Priority:    models.PriorityMedium,
	}
	_, err = taskService.Update(userID, createdTask.ID, updatedTask, "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("Failed to update task: %v", err)
	}

	// Verify the update
	retrievedTask, err := taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Errorf("Failed to get updated task: %v", err)
	}

	if retrievedTask.Title != "Updated Task Title" {
		t.Errorf("Expected title 'Updated Task Title', got %v", retrievedTask.Title)
	}

	// Delete the task using service directly
	err = taskService.Delete(userID, createdTask.ID, "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("Failed to delete task: %v", err)
	}

	// Verify the task is soft-deleted
	_, err = taskService.GetByID(userID, createdTask.ID)
	if err == nil {
		t.Error("Expected error when getting deleted task, got nil")
	}
}

func TestTaskAPIIntegration_CompleteAndUncomplete(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "testuser3", "test3@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "testuser3").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)

	// Create a task
	task := &models.Task{
		Title: "Test Task for Completion",
	}
	createdTask, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Mark as completed using service directly
	err = taskService.MarkAsCompleted(userID, createdTask.ID, "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("Failed to mark task as completed: %v", err)
	}

	// Verify task is completed
	updatedTask, err := taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get updated task: %v", err)
	}

	if !updatedTask.Completed {
		t.Error("Expected task to be completed, but it's not")
	}

	// Mark as uncompleted using service directly
	err = taskService.MarkAsUncompleted(userID, createdTask.ID, "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("Failed to mark task as uncompleted: %v", err)
	}

	// Verify task is uncompleted
	finalTask, err := taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get final task: %v", err)
	}

	if finalTask.Completed {
		t.Error("Expected task to be uncompleted, but it's still completed")
	}
}

func TestTaskAPIIntegration_ListWithPagination(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	// Clean up all test data before this test
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("Failed to cleanup all test data: %v", err)
	}
	_, err = db.Exec("DELETE FROM users WHERE name LIKE 'testuser%' OR name LIKE 'timezoneuser%' OR name LIKE 'tzupdateuser%' OR name LIKE 'invalidtzuser%' OR name LIKE 'eastuser' OR name LIKE 'westuser' OR name LIKE 'subtestuser%' OR name LIKE 'user%'")
	if err != nil {
		t.Fatalf("Failed to cleanup test users: %v", err)
	}
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err = db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "testuser4", "test4@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "testuser4").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	taskController := controllers.NewTaskController(taskService)

	// Create multiple tasks
	for i := 1; i <= 15; i++ {
		task := &models.Task{
			Title:       "Test Task Pagination",
			Description: "Pagination test task",
		}
		_, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
		if err != nil {
			t.Fatalf("Failed to create test task %d: %v", i, err)
		}
	}

	// Test first page
	req := httptest.NewRequest("GET", "/tasks?page=1&limit=10", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	w := httptest.NewRecorder()

	taskController.GetAllTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	meta := data["meta"].(map[string]interface{})

	if meta["total_data"].(float64) != 15 {
		t.Errorf("Expected total_data 15, got %v", meta["total_data"])
	}

	if meta["total_page"].(float64) != 2 {
		t.Errorf("Expected total_page 2, got %v", meta["total_page"])
	}

	if meta["has_next"].(bool) != true {
		t.Error("Expected has_next to be true")
	}
}

func TestTaskAPIIntegration_WithDueDate(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "testuser5", "test5@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "testuser5").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)

	// Create a task with future due date using epoch milliseconds
	futureTime := time.Now().Add(24 * time.Hour)
	futureEpoch := futureTime.UnixMilli()
	task := &models.Task{
		Title:       "Test Task with Due Date",
		Description: "Task with future due date",
		DueDate:     &futureEpoch,
	}
	createdTask, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create task with future due date: %v", err)
	}

	if createdTask.DueDate == nil {
		t.Error("Expected due date to be set")
	} else if *createdTask.DueDate != futureEpoch {
		t.Errorf("Expected due date epoch %d, got %d", futureEpoch, *createdTask.DueDate)
	}

	// Try to create a task with past due date
	pastTime := time.Now().Add(-24 * time.Hour)
	pastEpoch := pastTime.UnixMilli()
	pastTask := &models.Task{
		Title:       "Test Task with Past Due Date",
		Description: "Task with past due date",
		DueDate:     &pastEpoch,
	}
	_, err = taskService.Create(userID, pastTask, "127.0.0.1", "test-agent")
	// Note: The service layer might not validate past due dates, so this might succeed
	// This test mainly checks that epoch timestamps work correctly
	if err != nil {
		t.Logf("Note: Past due date validation: %v", err)
	}
}

// Helper function to create a request with PathValue support
func createPathValueRequest(req *http.Request, idStr string) *http.Request {
	// For integration tests, we need to skip these since they require the actual router
	// or a more sophisticated mocking approach for PathValue
	return req
}

// Integration tests for Subtask API endpoints
func TestSubtaskAPIIntegration_CreateAndGet(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "subtestuser", "subtest@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "subtestuser").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	// Setup the application layers
	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	subtaskModel := models.NewSubtaskModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)
	taskController := controllers.NewTaskController(taskService)

	// Create a test task first
	taskJSON := `{"title": "Test Task for Subtasks", "description": "Task with subtasks"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	w := httptest.NewRecorder()

	taskController.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	// Extract the created task ID
	data := createResponse["data"].(map[string]interface{})
	taskID := int64(data["id"].(float64))

	// Create a subtask using service directly (skip controller test for PathValue)
	subtask := &models.Subtask{
		Title: "Test Subtask",
	}
	createdSubtask, err := subtaskService.Create(userID, taskID, subtask)
	if err != nil {
		t.Fatalf("Failed to create subtask: %v", err)
	}

	if createdSubtask.Title != "Test Subtask" {
		t.Errorf("Expected title 'Test Subtask', got %v", createdSubtask.Title)
	}

	// Get subtasks by task ID
	subtasks, err := subtaskService.GetByTaskID(userID, taskID)
	if err != nil {
		t.Fatalf("Failed to get subtasks: %v", err)
	}

	if len(subtasks) != 1 {
		t.Errorf("Expected 1 subtask, got %d", len(subtasks))
	}

	// Verify task progress was updated
	task, err := taskService.GetByID(userID, taskID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ProgressPercentage != 0 {
		t.Errorf("Expected progress 0, got %d", task.ProgressPercentage)
	}
}

func TestSubtaskAPIIntegration_UpdateAndToggle(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "subtestuser2", "subtest2@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "subtestuser2").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	subtaskModel := models.NewSubtaskModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)

	// Create a task
	task := &models.Task{
		Title: "Test Task for Subtask Update",
	}
	createdTask, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Create a subtask
	subtask := &models.Subtask{
		Title: "Original Subtask",
	}
	createdSubtask, err := subtaskService.Create(userID, createdTask.ID, subtask)
	if err != nil {
		t.Fatalf("Failed to create subtask: %v", err)
	}

	// Update the subtask
	updatedSubtask := &models.Subtask{
		Title: "Updated Subtask",
	}
	_, err = subtaskService.Update(userID, createdSubtask.ID, updatedSubtask)
	if err != nil {
		t.Errorf("Failed to update subtask: %v", err)
	}

	// Verify the update
	retrievedSubtask, err := subtaskService.GetByID(userID, createdSubtask.ID)
	if err != nil {
		t.Errorf("Failed to get updated subtask: %v", err)
	}

	if retrievedSubtask.Title != "Updated Subtask" {
		t.Errorf("Expected title 'Updated Subtask', got %v", retrievedSubtask.Title)
	}

	// Toggle the subtask to completed
	err = subtaskService.Toggle(userID, createdSubtask.ID)
	if err != nil {
		t.Errorf("Failed to toggle subtask: %v", err)
	}

	// Verify subtask is completed
	toggledSubtask, err := subtaskService.GetByID(userID, createdSubtask.ID)
	if err != nil {
		t.Fatalf("Failed to get toggled subtask: %v", err)
	}

	if !toggledSubtask.IsCompleted {
		t.Error("Expected subtask to be completed, but it's not")
	}

	// Verify task progress was updated to 100%
	task, err = taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ProgressPercentage != 100 {
		t.Errorf("Expected progress 100, got %d", task.ProgressPercentage)
	}
}

func TestSubtaskAPIIntegration_Delete(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "subtestuser3", "subtest3@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "subtestuser3").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	subtaskModel := models.NewSubtaskModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)

	// Create a task
	task := &models.Task{
		Title: "Test Task for Subtask Delete",
	}
	createdTask, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Create a subtask
	subtask := &models.Subtask{
		Title: "Subtask to Delete",
	}
	createdSubtask, err := subtaskService.Create(userID, createdTask.ID, subtask)
	if err != nil {
		t.Fatalf("Failed to create subtask: %v", err)
	}

	// Delete the subtask
	err = subtaskService.Delete(userID, createdSubtask.ID)
	if err != nil {
		t.Errorf("Failed to delete subtask: %v", err)
	}

	// Verify the subtask is deleted
	_, err = subtaskService.GetByID(userID, createdSubtask.ID)
	if err == nil {
		t.Error("Expected error when getting deleted subtask, got nil")
	}

	// Verify task progress was updated back to 0%
	task, err = taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ProgressPercentage != 0 {
		t.Errorf("Expected progress 0, got %d", task.ProgressPercentage)
	}
}

func TestSubtaskAPIIntegration_ProgressCalculation(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "subtestuser4", "subtest4@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get the user ID
	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "subtestuser4").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	subtaskModel := models.NewSubtaskModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)

	// Create a task
	task := &models.Task{
		Title: "Test Task for Progress Calculation",
	}
	createdTask, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Create 4 subtasks
	for i := 1; i <= 4; i++ {
		subtask := &models.Subtask{
			Title: "Subtask",
		}
		_, err := subtaskService.Create(userID, createdTask.ID, subtask)
		if err != nil {
			t.Fatalf("Failed to create subtask %d: %v", i, err)
		}
	}

	// Check initial progress (0%)
	task, err = taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ProgressPercentage != 0 {
		t.Errorf("Expected initial progress 0, got %d", task.ProgressPercentage)
	}

	// Get all subtasks
	subtasks, err := subtaskService.GetByTaskID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get subtasks: %v", err)
	}

	// Complete 2 out of 4 subtasks (50%)
	err = subtaskService.Toggle(userID, subtasks[0].ID)
	if err != nil {
		t.Fatalf("Failed to toggle first subtask: %v", err)
	}

	err = subtaskService.Toggle(userID, subtasks[1].ID)
	if err != nil {
		t.Fatalf("Failed to toggle second subtask: %v", err)
	}

	// Check progress (50%)
	task, err = taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ProgressPercentage != 50 {
		t.Errorf("Expected progress 50, got %d", task.ProgressPercentage)
	}

	// Complete all subtasks (100%)
	err = subtaskService.Toggle(userID, subtasks[2].ID)
	if err != nil {
		t.Fatalf("Failed to toggle third subtask: %v", err)
	}

	err = subtaskService.Toggle(userID, subtasks[3].ID)
	if err != nil {
		t.Fatalf("Failed to toggle fourth subtask: %v", err)
	}

	// Check final progress (100%)
	task, err = taskService.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ProgressPercentage != 100 {
		t.Errorf("Expected final progress 100, got %d", task.ProgressPercentage)
	}
}

func TestSubtaskAPIIntegration_AccessControl(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create two test users
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "user1", "user1@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)", "user2", "user2@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	// Get the user IDs
	var userID1, userID2 int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "user1").Scan(&userID1)
	if err != nil {
		t.Fatalf("Failed to get test user 1 ID: %v", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "user2").Scan(&userID2)
	if err != nil {
		t.Fatalf("Failed to get test user 2 ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	subtaskModel := models.NewSubtaskModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)
	subtaskService := services.NewSubtaskService(subtaskModel, taskModel)

	// Create a task for user1
	task := &models.Task{
		Title: "User1 Task",
	}
	createdTask, err := taskService.Create(userID1, task, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Create a subtask for user1's task
	subtask := &models.Subtask{
		Title: "User1 Subtask",
	}
	createdSubtask, err := subtaskService.Create(userID1, createdTask.ID, subtask)
	if err != nil {
		t.Fatalf("Failed to create subtask: %v", err)
	}

	// Try to access the subtask as user2 (should fail)
	_, err = subtaskService.GetByID(userID2, createdSubtask.ID)
	if err == nil {
		t.Error("Expected error when user2 tries to access user1's subtask, got nil")
	}

	// Try to update the subtask as user2 (should fail)
	updatedSubtask := &models.Subtask{
		Title: "Malicious Update",
	}
	_, err = subtaskService.Update(userID2, createdSubtask.ID, updatedSubtask)
	if err == nil {
		t.Error("Expected error when user2 tries to update user1's subtask, got nil")
	}

	// Try to delete the subtask as user2 (should fail)
	err = subtaskService.Delete(userID2, createdSubtask.ID)
	if err == nil {
		t.Error("Expected error when user2 tries to delete user1's subtask, got nil")
	}

	// Try to toggle the subtask as user2 (should fail)
	err = subtaskService.Toggle(userID2, createdSubtask.ID)
	if err == nil {
		t.Error("Expected error when user2 tries to toggle user1's subtask, got nil")
	}
}

// Integration tests for timezone functionality
func TestTimezoneIntegration_MultipleTimezones(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create users with different timezones
	timezones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
	}

	var userIDs []int64
	for i, tz := range timezones {
		userName := "timezoneuser" + string(rune('1'+i))
		_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
			userName, "timezone"+string(rune('1'+i))+"@example.com", "hashedpassword", tz)
		if err != nil {
			t.Fatalf("Failed to create test user for timezone %s: %v", tz, err)
		}

		var userID int64
		err = db.QueryRow("SELECT id FROM users WHERE name = $1", userName).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to get test user ID for timezone %s: %v", tz, err)
		}
		userIDs = append(userIDs, userID)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)

	// Create tasks for each user with the same epoch timestamp
	futureTime := time.Now().Add(24 * time.Hour)
	futureEpoch := futureTime.UnixMilli()

	for i, userID := range userIDs {
		task := &models.Task{
			Title:       "Timezone Test Task",
			Description: "Task for timezone testing",
			DueDate:     &futureEpoch,
		}
		_, err := taskService.Create(userID, task, "127.0.0.1", "test-agent")
		if err != nil {
			t.Fatalf("Failed to create task for user in timezone %s: %v", timezones[i], err)
		}
	}

	// Verify that all tasks were created with the same epoch timestamp
	for i, userID := range userIDs {
		tasks, _, err := taskService.GetAll(userID, nil, 1, 10, "", nil, nil, "created_at", "DESC")
		if err != nil {
			t.Fatalf("Failed to get tasks for user in timezone %s: %v", timezones[i], err)
		}

		if len(tasks) != 1 {
			t.Errorf("Expected 1 task for user in timezone %s, got %d", timezones[i], len(tasks))
		}

		// Verify the due date epoch is the same
		if tasks[0].DueDate == nil {
			t.Errorf("Expected due date to be set for user in timezone %s", timezones[i])
		} else if *tasks[0].DueDate != futureEpoch {
			t.Errorf("Expected due date epoch %d for user in timezone %s, got %d", futureEpoch, timezones[i], *tasks[0].DueDate)
		}
	}

	// Test timezone conversion utilities
	for _, tz := range timezones {
		// Test converting the epoch to each timezone
		convertedTime, err := utils.ConvertToTimezone(futureEpoch, tz)
		if err != nil {
			t.Errorf("Failed to convert epoch to timezone %s: %v", tz, err)
		}

		// Verify the conversion works
		if convertedTime.IsZero() {
			t.Errorf("Converted time for timezone %s should not be zero", tz)
		}

		// Test formatting in each timezone
		formatted, err := utils.FormatEpochMillis(futureEpoch, tz, "2006-01-02 15:04:05")
		if err != nil {
			t.Errorf("Failed to format epoch for timezone %s: %v", tz, err)
		}

		if formatted == "" {
			t.Errorf("Formatted time for timezone %s should not be empty", tz)
		}

		// Test timezone-aware timestamp
		tzAware, err := utils.ConvertToTimezoneAware(futureEpoch, tz)
		if err != nil {
			t.Errorf("Failed to create timezone-aware timestamp for %s: %v", tz, err)
		}

		if tzAware == nil {
			t.Errorf("Timezone-aware timestamp for %s should not be nil", tz)
		} else {
			if tzAware.Timezone != tz {
				t.Errorf("Expected timezone %s, got %s", tz, tzAware.Timezone)
			}
			if tzAware.EpochMillis != futureEpoch {
				t.Errorf("Expected epoch %d, got %d", futureEpoch, tzAware.EpochMillis)
			}
		}
	}
}

func TestTimezoneIntegration_UserTimezoneUpdates(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user with UTC timezone
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"tzupdateuser", "tzupdate@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "tzupdateuser").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	userModel := models.NewUserModel(db)

	// Update user timezone to America/New_York
	err = userModel.UpdateTimezone(userID, "America/New_York")
	if err != nil {
		t.Fatalf("Failed to update user timezone: %v", err)
	}

	// Verify the timezone was updated
	user, err := userModel.GetByID(userID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if user.Timezone != "America/New_York" {
		t.Errorf("Expected timezone America/New_York, got %s", user.Timezone)
	}

	// Update user timezone to Asia/Tokyo
	err = userModel.UpdateTimezone(userID, "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Failed to update user timezone to Asia/Tokyo: %v", err)
	}

	// Verify the timezone was updated again
	user, err = userModel.GetByID(userID)
	if err != nil {
		t.Fatalf("Failed to get updated user again: %v", err)
	}

	if user.Timezone != "Asia/Tokyo" {
		t.Errorf("Expected timezone Asia/Tokyo, got %s", user.Timezone)
	}
}

func TestTimezoneIntegration_TaskTimestampsAcrossTimezones(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create two users in different timezones
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"eastuser", "east@example.com", "hashedpassword", "America/New_York")
	if err != nil {
		t.Fatalf("Failed to create east coast user: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"westuser", "west@example.com", "hashedpassword", "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Failed to create west coast user: %v", err)
	}

	var eastUserID, westUserID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "eastuser").Scan(&eastUserID)
	if err != nil {
		t.Fatalf("Failed to get east user ID: %v", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "westuser").Scan(&westUserID)
	if err != nil {
		t.Fatalf("Failed to get west user ID: %v", err)
	}

	taskModel := models.NewTaskModel(db)
	tagModel := models.NewTagModel(db)
	taskService := services.NewTaskService(taskModel, tagModel, nil)

	// Create tasks at the same moment
	eastTask := &models.Task{
		Title:       "East Coast Task",
		Description: "Task created from east coast",
	}
	createdEastTask, err := taskService.Create(eastUserID, eastTask, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create east task: %v", err)
	}

	westTask := &models.Task{
		Title:       "West Coast Task",
		Description: "Task created from west coast",
	}
	createdWestTask, err := taskService.Create(westUserID, westTask, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("Failed to create west task: %v", err)
	}

	// Both tasks should have similar epoch timestamps (within a few seconds)
	timeDiff := createdEastTask.CreatedAt - createdWestTask.CreatedAt
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}

	// Allow for 2 seconds difference due to execution time
	if timeDiff > 2000 {
		t.Errorf("Tasks created at same time should have similar timestamps. Difference: %d ms", timeDiff)
	}

	// Convert both timestamps to their respective timezones
	eastTime, err := utils.ConvertToTimezone(createdEastTask.CreatedAt, "America/New_York")
	if err != nil {
		t.Fatalf("Failed to convert east task time: %v", err)
	}

	westTime, err := utils.ConvertToTimezone(createdWestTask.CreatedAt, "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Failed to convert west task time: %v", err)
	}

	// The formatted times should be different due to timezone offset
	eastFormatted := eastTime.Format("2006-01-02 15:04:05")
	westFormatted := westTime.Format("2006-01-02 15:04:05")

	if eastFormatted == westFormatted {
		t.Log("Note: Formatted times are the same, which is possible depending on the exact time created")
	}

	// But the epoch milliseconds should be very close
	epochDiff := createdEastTask.CreatedAt - createdWestTask.CreatedAt
	if epochDiff < 0 {
		epochDiff = -epochDiff
	}

	if epochDiff > 2000 {
		t.Errorf("Epoch timestamps should be very close. Difference: %d ms", epochDiff)
	}
}

func TestTimezoneIntegration_InvalidTimezoneHandling(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"invalidtzuser", "invalidtz@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "invalidtzuser").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	userModel := models.NewUserModel(db)

	// Try to update with invalid timezone (should fail in real application, but database might accept it)
	// The validation should happen at application level
	err = userModel.UpdateTimezone(userID, "Invalid/Timezone")
	if err != nil {
		t.Logf("Expected: Update with invalid timezone failed as expected: %v", err)
	}

	// Test timezone validation utility
	if utils.IsValidTimezone("Invalid/Timezone") {
		t.Error("IsValidTimezone should return false for invalid timezone")
	}

	if !utils.IsValidTimezone("UTC") {
		t.Error("IsValidTimezone should return true for valid timezone")
	}

	if !utils.IsValidTimezone("America/New_York") {
		t.Error("IsValidTimezone should return true for valid timezone")
	}
}
