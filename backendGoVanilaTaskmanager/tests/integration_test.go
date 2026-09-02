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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "testuser", "test@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "testuser2", "test2@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "testuser3", "test3@example.com", "hashedpassword")
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
	_, err = db.Exec("DELETE FROM users WHERE name LIKE 'testuser%'")
	if err != nil {
		t.Fatalf("Failed to cleanup test users: %v", err)
	}
	defer CleanupTestData(db, t)

	// Create a test user first
	_, err = db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "testuser4", "test4@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "testuser5", "test5@example.com", "hashedpassword")
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
	taskController := controllers.NewTaskController(taskService)

	// Create a task with future due date
	futureTime := time.Now().Add(24 * time.Hour)
	taskJSON := `{"title": "Test Task with Due Date", "due_date": "` + futureTime.Format(time.RFC3339) + `"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	w := httptest.NewRecorder()

	taskController.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Try to create a task with past due date (should fail)
	pastTime := time.Now().Add(-24 * time.Hour)
	pastTaskJSON := `{"title": "Test Task with Past Due Date", "due_date": "` + pastTime.Format(time.RFC3339) + `"}`
	pastReq := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(pastTaskJSON))
	pastReq.Header.Set("Content-Type", "application/json")
	pastReq = pastReq.WithContext(context.WithValue(pastReq.Context(), "user_id", userID))
	pastW := httptest.NewRecorder()

	taskController.CreateTask(pastW, pastReq)

	if pastW.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for past due date, got %d. Body: %s", pastW.Code, pastW.Body.String())
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "subtestuser", "subtest@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "subtestuser2", "subtest2@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "subtestuser3", "subtest3@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "subtestuser4", "subtest4@example.com", "hashedpassword")
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
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "user1", "user1@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", "user2", "user2@example.com", "hashedpassword")
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
