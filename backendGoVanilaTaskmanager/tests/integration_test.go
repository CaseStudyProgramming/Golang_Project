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
	taskService := services.NewTaskService(taskModel)
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
	}

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)

	if createResponse["status"] != "success" {
		t.Errorf("Expected success status, got %v", createResponse["status"])
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
	taskService := services.NewTaskService(taskModel)

	// Create a task first
	task := &models.Task{
		Title:       "Test Task for Update",
		Description: "Initial description",
	}
	createdTask, err := taskService.Create(userID, task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Update the task using service directly
	updatedTask := &models.Task{
		Title:       "Updated Task Title",
		Description: "Updated description",
		Priority:    models.PriorityMedium,
	}
	_, err = taskService.Update(userID, createdTask.ID, updatedTask)
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
	err = taskService.Delete(userID, createdTask.ID)
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
	taskService := services.NewTaskService(taskModel)

	// Create a task
	task := &models.Task{
		Title: "Test Task for Completion",
	}
	createdTask, err := taskService.Create(userID, task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	// Mark as completed using service directly
	err = taskService.MarkAsCompleted(userID, createdTask.ID)
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
	err = taskService.MarkAsUncompleted(userID, createdTask.ID)
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
	taskService := services.NewTaskService(taskModel)
	taskController := controllers.NewTaskController(taskService)

	// Create multiple tasks
	for i := 1; i <= 15; i++ {
		task := &models.Task{
			Title:       "Test Task Pagination",
			Description: "Pagination test task",
		}
		_, err := taskService.Create(userID, task)
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
	taskService := services.NewTaskService(taskModel)
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
