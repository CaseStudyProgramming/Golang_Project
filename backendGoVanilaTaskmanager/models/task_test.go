package models

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Berjuang#382 dbname=taskmanager_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: database not reachable")
	}

	// Clean up and create tables
	_, _ = db.Exec("DROP TABLE IF EXISTS tasks CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS users CASCADE")

	_, err = db.Exec(`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(150) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		timezone VARCHAR(50) DEFAULT 'UTC',
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
	)`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE tasks (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		category_id INT NULL,
		title VARCHAR(255) NOT NULL,
		sub_title VARCHAR(255) NULL,
		description TEXT NULL,
		completed BOOLEAN DEFAULT false,
		due_date BIGINT NULL,
		priority VARCHAR(20) DEFAULT 'MEDIUM' CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')),
		progress_percentage INT DEFAULT 0,
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		deleted_at BIGINT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create tasks table: %v", err)
	}

	return db
}

func cleanupTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestNewTaskModel(t *testing.T) {
	db := &sql.DB{}
	model := NewTaskModel(db)

	if model == nil {
		t.Fatal("NewTaskModel returned nil")
	}

	if model.DB != db {
		t.Error("NewTaskModel did not set DB correctly")
	}
}

func TestTaskModel_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create a test user first
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:      userID,
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdTask.ID == 0 {
		t.Error("Create did not set task ID")
	}

	if createdTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", createdTask.Title)
	}

	if createdTask.ProgressPercentage != 0 {
		t.Errorf("Expected progress_percentage 0, got %d", createdTask.ProgressPercentage)
	}
}

func TestTaskModel_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:   userID,
		Title:    "Test Task",
		Priority: PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByID
	retrievedTask, err := model.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedTask.ID != createdTask.ID {
		t.Errorf("Expected ID %d, got %d", createdTask.ID, retrievedTask.ID)
	}

	if retrievedTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", retrievedTask.Title)
	}

	// Test GetByID with wrong user
	_, err = model.GetByID(userID+999, createdTask.ID)
	if err == nil {
		t.Error("Expected error when getting task with wrong user ID")
	}
}

func TestTaskModel_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	// Create multiple tasks
	for i := 0; i < 5; i++ {
		task := &Task{
			UserID:   userID,
			Title:    "Test Task",
			Priority: PriorityMedium,
		}
		_, err := model.Create(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
	}

	// Test GetAll
	tasks, total, err := model.GetAll(userID, nil, 0, 10, "", nil, nil, "created_at", "DESC")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(tasks) != 5 {
		t.Errorf("Expected 5 tasks, got %d", len(tasks))
	}

	// Test with completed filter
	completed := true
	tasks, total, err = model.GetAll(userID, &completed, 0, 10, "", nil, nil, "created_at", "DESC")
	if err != nil {
		t.Fatalf("GetAll with completed filter failed: %v", err)
	}

	if total != 0 {
		t.Errorf("Expected 0 completed tasks, got %d", total)
	}
}

func TestTaskModel_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:   userID,
		Title:    "Original Title",
		Priority: PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update task
	createdTask.Title = "Updated Title"
	createdTask.Description = "Updated Description"
	err = model.Update(userID, createdTask)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrievedTask, err := model.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if retrievedTask.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", retrievedTask.Title)
	}

	if retrievedTask.Description != "Updated Description" {
		t.Errorf("Expected description 'Updated Description', got '%s'", retrievedTask.Description)
	}
}

func TestTaskModel_MarkTaskAsCompleted(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:   userID,
		Title:    "Test Task",
		Priority: PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Mark as completed
	err = model.MarkTaskAsCompleted(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("MarkTaskAsCompleted failed: %v", err)
	}

	// Verify
	retrievedTask, err := model.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("GetByID after MarkTaskAsCompleted failed: %v", err)
	}

	if !retrievedTask.Completed {
		t.Error("Task was not marked as completed")
	}
}

func TestTaskModel_MarkTaskAsUncompleted(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:    userID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: true,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Mark as uncompleted
	err = model.MarkTaskAsUncompleted(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("MarkTaskAsUncompleted failed: %v", err)
	}

	// Verify
	retrievedTask, err := model.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("GetByID after MarkTaskAsUncompleted failed: %v", err)
	}

	if retrievedTask.Completed {
		t.Error("Task was not marked as uncompleted")
	}
}

func TestTaskModel_Delete_SoftDelete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:   userID,
		Title:    "Test Task",
		Priority: PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Soft delete
	err = model.Delete(userID, createdTask.ID, true)
	if err != nil {
		t.Fatalf("Delete (soft) failed: %v", err)
	}

	// Verify task is not found (soft deleted)
	_, err = model.GetByID(userID, createdTask.ID)
	if err == nil {
		t.Error("Expected error when getting soft-deleted task")
	}
}

func TestTaskModel_RestoreTask(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:   userID,
		Title:    "Test Task",
		Priority: PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Soft delete
	err = model.Delete(userID, createdTask.ID, true)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Restore
	err = model.RestoreTask(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("RestoreTask failed: %v", err)
	}

	// Verify task is found again
	retrievedTask, err := model.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("GetByID after restore failed: %v", err)
	}

	if retrievedTask.DeletedAt != nil {
		t.Error("Task still has deleted_at set after restore")
	}
}

func TestTaskModel_BulkDelete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	// Create multiple tasks
	var taskIDs []int64
	for i := 0; i < 3; i++ {
		task := &Task{
			UserID:   userID,
			Title:    "Test Task",
			Priority: PriorityMedium,
		}
		createdTask, err := model.Create(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
		taskIDs = append(taskIDs, createdTask.ID)
	}

	// Bulk delete
	err = model.BulkDelete(userID, taskIDs)
	if err != nil {
		t.Fatalf("BulkDelete failed: %v", err)
	}

	// Verify all tasks are deleted
	for _, id := range taskIDs {
		_, err := model.GetByID(userID, id)
		if err == nil {
			t.Errorf("Task %d should be deleted but was found", id)
		}
	}
}

func TestTaskModel_BulkComplete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	// Create multiple tasks
	var taskIDs []int64
	for i := 0; i < 3; i++ {
		task := &Task{
			UserID:   userID,
			Title:    "Test Task",
			Priority: PriorityMedium,
		}
		createdTask, err := model.Create(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
		taskIDs = append(taskIDs, createdTask.ID)
	}

	// Bulk complete
	err = model.BulkComplete(userID, taskIDs)
	if err != nil {
		t.Fatalf("BulkComplete failed: %v", err)
	}

	// Verify all tasks are completed
	for _, id := range taskIDs {
		task, err := model.GetByID(userID, id)
		if err != nil {
			t.Fatalf("Failed to get task %d: %v", id, err)
		}
		if !task.Completed {
			t.Errorf("Task %d should be completed but is not", id)
		}
	}
}

func TestTaskModel_GetAnalyticsSummary(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	model := NewTaskModel(db)

	// Create some tasks
	for i := 0; i < 3; i++ {
		task := &Task{
			UserID:    userID,
			Title:     "Test Task",
			Priority:  PriorityMedium,
			Completed: i == 0, // First task completed
		}
		_, err := model.Create(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
	}

	// Get analytics
	summary, err := model.GetAnalyticsSummary(userID)
	if err != nil {
		t.Fatalf("GetAnalyticsSummary failed: %v", err)
	}

	totalActive, ok := summary["total_active"].(int)
	if !ok {
		t.Error("total_active is not an int")
	}
	if totalActive != 3 {
		t.Errorf("Expected total_active 3, got %d", totalActive)
	}

	totalCompleted, ok := summary["total_completed"].(int)
	if !ok {
		t.Error("total_completed is not an int")
	}
	if totalCompleted != 1 {
		t.Errorf("Expected total_completed 1, got %d", totalCompleted)
	}

	priorityDist, ok := summary["priority_distribution"].(map[string]int)
	if !ok {
		t.Error("priority_distribution is not a map[string]int")
	}
	if priorityDist["MEDIUM"] != 3 {
		t.Errorf("Expected 3 MEDIUM priority tasks, got %d", priorityDist["MEDIUM"])
	}
}

func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		expected string
	}{
		{"Low Priority", PriorityLow, "LOW"},
		{"Medium Priority", PriorityMedium, "MEDIUM"},
		{"High Priority", PriorityHigh, "HIGH"},
		{"Urgent Priority", PriorityUrgent, "URGENT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.priority) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.priority)
			}
		})
	}
}

func TestTaskModel_UpdateProgress(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	// Create test user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
		"Test User", "test@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", "test@example.com").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	// Create subtasks table
	_, err = db.Exec(`CREATE TABLE subtasks (
		id SERIAL PRIMARY KEY,
		task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		is_completed BOOLEAN DEFAULT false,
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
	)`)
	if err != nil {
		t.Fatalf("Failed to create subtasks table: %v", err)
	}

	model := NewTaskModel(db)

	task := &Task{
		UserID:   userID,
		Title:    "Test Task",
		Priority: PriorityMedium,
	}

	createdTask, err := model.Create(task)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update progress (should work even without subtasks)
	err = model.UpdateProgress(createdTask.ID)
	if err != nil {
		t.Fatalf("UpdateProgress failed: %v", err)
	}

	// Verify progress was updated
	retrievedTask, err := model.GetByID(userID, createdTask.ID)
	if err != nil {
		t.Fatalf("GetByID after UpdateProgress failed: %v", err)
	}

	// Progress should be 0 since there are no subtasks
	if retrievedTask.ProgressPercentage != 0 {
		t.Errorf("Expected progress_percentage 0, got %d", retrievedTask.ProgressPercentage)
	}
}
