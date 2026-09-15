package models

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupSubtaskTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Berjuang#382 dbname=taskmanager_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: database not reachable")
	}

	// Clean up and create tables
	_, _ = db.Exec("DROP TABLE IF EXISTS subtasks CASCADE")
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
		user_id BIGINT NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		priority VARCHAR(20) DEFAULT 'medium',
		due_date BIGINT,
		is_completed BOOLEAN DEFAULT false,
		category_id BIGINT,
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("Failed to create tasks table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE subtasks (
		id SERIAL PRIMARY KEY,
		task_id BIGINT NOT NULL,
		title VARCHAR(255) NOT NULL,
		is_completed BOOLEAN DEFAULT false,
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("Failed to create subtasks table: %v", err)
	}

	return db
}

func cleanupSubtaskTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestNewSubtaskModel(t *testing.T) {
	db := &sql.DB{}
	model := NewSubtaskModel(db)

	if model == nil {
		t.Fatal("NewSubtaskModel returned nil")
	}

	if model.DB != db {
		t.Error("NewSubtaskModel did not set DB correctly")
	}
}

func TestSubtaskModel_Create(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	subtask := &Subtask{
		TaskID:      createdTask.ID,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	createdSubtask, err := model.Create(subtask)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdSubtask.ID == 0 {
		t.Error("Create did not set subtask ID")
	}

	if createdSubtask.Title != "Test Subtask" {
		t.Errorf("Expected title 'Test Subtask', got '%s'", createdSubtask.Title)
	}

	if createdSubtask.IsCompleted != false {
		t.Errorf("Expected is_completed false, got %v", createdSubtask.IsCompleted)
	}
}

func TestSubtaskModel_GetByTaskID(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	// Create multiple subtasks
	subtaskTitles := []string{"Subtask 1", "Subtask 2", "Subtask 3"}
	for _, title := range subtaskTitles {
		subtask := &Subtask{
			TaskID:      createdTask.ID,
			Title:       title,
			IsCompleted: false,
		}
		_, err := model.Create(subtask)
		if err != nil {
			t.Fatalf("Failed to create subtask: %v", err)
		}
	}

	// Test GetByTaskID
	subtasks, err := model.GetByTaskID(createdTask.ID)
	if err != nil {
		t.Fatalf("GetByTaskID failed: %v", err)
	}

	if len(subtasks) != 3 {
		t.Errorf("Expected 3 subtasks, got %d", len(subtasks))
	}
}

func TestSubtaskModel_GetByID(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	subtask := &Subtask{
		TaskID:      createdTask.ID,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	createdSubtask, err := model.Create(subtask)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByID
	retrievedSubtask, err := model.GetByID(createdSubtask.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedSubtask.ID != createdSubtask.ID {
		t.Errorf("Expected ID %d, got %d", createdSubtask.ID, retrievedSubtask.ID)
	}

	if retrievedSubtask.Title != "Test Subtask" {
		t.Errorf("Expected title 'Test Subtask', got '%s'", retrievedSubtask.Title)
	}

	// Test with non-existent ID
	_, err = model.GetByID(99999)
	if err == nil {
		t.Error("Expected error when getting non-existent subtask")
	}
}

func TestSubtaskModel_Update(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	subtask := &Subtask{
		TaskID:      createdTask.ID,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	createdSubtask, err := model.Create(subtask)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update subtask
	createdSubtask.Title = "Updated Subtask"
	createdSubtask.IsCompleted = true
	err = model.Update(createdSubtask)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrievedSubtask, err := model.GetByID(createdSubtask.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if retrievedSubtask.Title != "Updated Subtask" {
		t.Errorf("Expected title 'Updated Subtask', got '%s'", retrievedSubtask.Title)
	}

	if retrievedSubtask.IsCompleted != true {
		t.Errorf("Expected is_completed true, got %v", retrievedSubtask.IsCompleted)
	}
}

func TestSubtaskModel_Delete(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	subtask := &Subtask{
		TaskID:      createdTask.ID,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	createdSubtask, err := model.Create(subtask)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete subtask
	err = model.Delete(createdSubtask.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify subtask is deleted
	_, err = model.GetByID(createdSubtask.ID)
	if err == nil {
		t.Error("Expected error when getting deleted subtask")
	}
}

func TestSubtaskModel_Toggle(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	subtask := &Subtask{
		TaskID:      createdTask.ID,
		Title:       "Test Subtask",
		IsCompleted: false,
	}

	createdSubtask, err := model.Create(subtask)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Toggle subtask (false -> true)
	err = model.Toggle(createdSubtask.ID)
	if err != nil {
		t.Fatalf("Toggle failed: %v", err)
	}

	// Verify toggle
	retrievedSubtask, err := model.GetByID(createdSubtask.ID)
	if err != nil {
		t.Fatalf("GetByID after toggle failed: %v", err)
	}

	if retrievedSubtask.IsCompleted != true {
		t.Errorf("Expected is_completed true after toggle, got %v", retrievedSubtask.IsCompleted)
	}

	// Toggle again (true -> false)
	err = model.Toggle(createdSubtask.ID)
	if err != nil {
		t.Fatalf("Second toggle failed: %v", err)
	}

	retrievedSubtask, err = model.GetByID(createdSubtask.ID)
	if err != nil {
		t.Fatalf("GetByID after second toggle failed: %v", err)
	}

	if retrievedSubtask.IsCompleted != false {
		t.Errorf("Expected is_completed false after second toggle, got %v", retrievedSubtask.IsCompleted)
	}
}

func TestSubtaskModel_CalculateProgress(t *testing.T) {
	db := setupSubtaskTestDB(t)
	defer cleanupSubtaskTestDB(db)

	// Create a test user and task
	userModel := NewUserModel(db)
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	createdUser, err := userModel.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	taskModel := NewTaskModel(db)
	task := &Task{
		UserID:    createdUser.ID,
		Title:     "Test Task",
		Priority:  PriorityMedium,
		Completed: false,
	}
	createdTask, err := taskModel.Create(task)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	model := NewSubtaskModel(db)

	// Test with no subtasks
	progress, err := model.CalculateProgress(createdTask.ID)
	if err != nil {
		t.Fatalf("CalculateProgress failed with no subtasks: %v", err)
	}
	if progress != 0 {
		t.Errorf("Expected progress 0 with no subtasks, got %d", progress)
	}

	// Create 4 subtasks, 2 completed
	for i := 0; i < 4; i++ {
		subtask := &Subtask{
			TaskID:      createdTask.ID,
			Title:       "Test Subtask",
			IsCompleted: i < 2, // First 2 are completed
		}
		_, err := model.Create(subtask)
		if err != nil {
			t.Fatalf("Failed to create subtask: %v", err)
		}
	}

	// Test progress calculation
	progress, err = model.CalculateProgress(createdTask.ID)
	if err != nil {
		t.Fatalf("CalculateProgress failed: %v", err)
	}
	if progress != 50 {
		t.Errorf("Expected progress 50 (2/4), got %d", progress)
	}

	// Complete all subtasks
	subtasks, _ := model.GetByTaskID(createdTask.ID)
	for _, subtask := range subtasks {
		if !subtask.IsCompleted {
			model.Toggle(subtask.ID)
		}
	}

	// Test 100% progress
	progress, err = model.CalculateProgress(createdTask.ID)
	if err != nil {
		t.Fatalf("CalculateProgress failed: %v", err)
	}
	if progress != 100 {
		t.Errorf("Expected progress 100 (4/4), got %d", progress)
	}
}
