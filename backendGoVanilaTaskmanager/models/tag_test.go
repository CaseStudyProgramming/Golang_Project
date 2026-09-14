package models

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupTagTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Berjuang#382 dbname=taskmanager_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: database not reachable")
	}

	// Clean up and create tables
	_, _ = db.Exec("DROP TABLE IF EXISTS task_tags CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS tags CASCADE")
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

	_, err = db.Exec(`CREATE TABLE tags (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		name VARCHAR(100) NOT NULL,
		color_hex VARCHAR(7),
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("Failed to create tags table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE task_tags (
		task_id BIGINT NOT NULL,
		tag_id BIGINT NOT NULL,
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		PRIMARY KEY (task_id, tag_id),
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("Failed to create task_tags table: %v", err)
	}

	return db
}

func cleanupTagTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestNewTagModel(t *testing.T) {
	db := &sql.DB{}
	model := NewTagModel(db)

	if model == nil {
		t.Fatal("NewTagModel returned nil")
	}

	if model.DB != db {
		t.Error("NewTagModel did not set DB correctly")
	}
}

func TestTagModel_Create(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	model := NewTagModel(db)

	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}

	createdTag, err := model.Create(tag)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdTag.ID == 0 {
		t.Error("Create did not set tag ID")
	}

	if createdTag.Name != "Urgent" {
		t.Errorf("Expected name 'Urgent', got '%s'", createdTag.Name)
	}

	if createdTag.ColorHex != "#FF0000" {
		t.Errorf("Expected color hex '#FF0000', got '%s'", createdTag.ColorHex)
	}
}

func TestTagModel_GetAll(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	model := NewTagModel(db)

	// Create multiple tags
	tagNames := []string{"Urgent", "Important", "Work"}
	for _, name := range tagNames {
		tag := &Tag{
			UserID:   createdUser.ID,
			Name:     name,
			ColorHex: "#FF0000",
		}
		_, err := model.Create(tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}
	}

	// Test GetAll
	tags, err := model.GetAll(createdUser.ID)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
}

func TestTagModel_GetByID(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	model := NewTagModel(db)

	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}

	createdTag, err := model.Create(tag)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByID
	retrievedTag, err := model.GetByID(createdUser.ID, createdTag.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedTag.ID != createdTag.ID {
		t.Errorf("Expected ID %d, got %d", createdTag.ID, retrievedTag.ID)
	}

	if retrievedTag.Name != "Urgent" {
		t.Errorf("Expected name 'Urgent', got '%s'", retrievedTag.Name)
	}

	// Test with non-existent ID
	_, err = model.GetByID(createdUser.ID, 99999)
	if err == nil {
		t.Error("Expected error when getting non-existent tag")
	}
}

func TestTagModel_Update(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	model := NewTagModel(db)

	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}

	createdTag, err := model.Create(tag)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update tag
	createdTag.Name = "Very Urgent"
	createdTag.ColorHex = "#00FF00"
	err = model.Update(createdUser.ID, createdTag)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrievedTag, err := model.GetByID(createdUser.ID, createdTag.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if retrievedTag.Name != "Very Urgent" {
		t.Errorf("Expected name 'Very Urgent', got '%s'", retrievedTag.Name)
	}

	if retrievedTag.ColorHex != "#00FF00" {
		t.Errorf("Expected color hex '#00FF00', got '%s'", retrievedTag.ColorHex)
	}
}

func TestTagModel_Delete(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	model := NewTagModel(db)

	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}

	createdTag, err := model.Create(tag)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete tag
	err = model.Delete(createdUser.ID, createdTag.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify tag is deleted
	_, err = model.GetByID(createdUser.ID, createdTag.ID)
	if err == nil {
		t.Error("Expected error when getting deleted tag")
	}
}

func TestTagModel_AddTagToTask(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	// Create a test task
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

	// Create a test tag
	tagModel := NewTagModel(db)
	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}
	createdTag, err := tagModel.Create(tag)
	if err != nil {
		t.Fatalf("Failed to create test tag: %v", err)
	}

	// Add tag to task
	err = tagModel.AddTagToTask(createdTask.ID, createdTag.ID)
	if err != nil {
		t.Fatalf("AddTagToTask failed: %v", err)
	}

	// Verify the relationship
	tags, err := tagModel.GetTagsByTaskID(createdTask.ID)
	if err != nil {
		t.Fatalf("GetTagsByTaskID failed: %v", err)
	}

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}

	if tags[0].ID != createdTag.ID {
		t.Errorf("Expected tag ID %d, got %d", createdTag.ID, tags[0].ID)
	}
}

func TestTagModel_RemoveTagFromTask(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	// Create a test task
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

	// Create a test tag
	tagModel := NewTagModel(db)
	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}
	createdTag, err := tagModel.Create(tag)
	if err != nil {
		t.Fatalf("Failed to create test tag: %v", err)
	}

	// Add tag to task
	err = tagModel.AddTagToTask(createdTask.ID, createdTag.ID)
	if err != nil {
		t.Fatalf("AddTagToTask failed: %v", err)
	}

	// Remove tag from task
	err = tagModel.RemoveTagFromTask(createdTask.ID, createdTag.ID)
	if err != nil {
		t.Fatalf("RemoveTagFromTask failed: %v", err)
	}

	// Verify the relationship is removed
	tags, err := tagModel.GetTagsByTaskID(createdTask.ID)
	if err != nil {
		t.Fatalf("GetTagsByTaskID failed: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected 0 tags after removal, got %d", len(tags))
	}
}

func TestTagModel_GetTagsByTaskID(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	// Create a test task
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

	// Create multiple tags
	tagModel := NewTagModel(db)
	tagNames := []string{"Urgent", "Important"}
	var createdTags []Tag
	for _, name := range tagNames {
		tag := &Tag{
			UserID:   createdUser.ID,
			Name:     name,
			ColorHex: "#FF0000",
		}
		createdTag, err := tagModel.Create(tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}
		createdTags = append(createdTags, *createdTag)
	}

	// Add tags to task
	for _, tag := range createdTags {
		err = tagModel.AddTagToTask(createdTask.ID, tag.ID)
		if err != nil {
			t.Fatalf("AddTagToTask failed: %v", err)
		}
	}

	// Test GetTagsByTaskID
	tags, err := tagModel.GetTagsByTaskID(createdTask.ID)
	if err != nil {
		t.Fatalf("GetTagsByTaskID failed: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestTagModel_GetTasksByTagID(t *testing.T) {
	db := setupTagTestDB(t)
	defer cleanupTagTestDB(db)

	// Create a test user
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

	// Create multiple tasks
	taskModel := NewTaskModel(db)
	taskTitles := []string{"Task 1", "Task 2"}
	var createdTasks []Task
	for _, title := range taskTitles {
		task := &Task{
			UserID:    createdUser.ID,
			Title:     title,
			Priority:  PriorityMedium,
			Completed: false,
		}
		createdTask, err := taskModel.Create(task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
		createdTasks = append(createdTasks, *createdTask)
	}

	// Create a tag
	tagModel := NewTagModel(db)
	tag := &Tag{
		UserID:   createdUser.ID,
		Name:     "Urgent",
		ColorHex: "#FF0000",
	}
	createdTag, err := tagModel.Create(tag)
	if err != nil {
		t.Fatalf("Failed to create tag: %v", err)
	}

	// Add tag to tasks
	for _, task := range createdTasks {
		err = tagModel.AddTagToTask(task.ID, createdTag.ID)
		if err != nil {
			t.Fatalf("AddTagToTask failed: %v", err)
		}
	}

	// Test GetTasksByTagID
	tasks, err := tagModel.GetTasksByTagID(createdTag.ID)
	if err != nil {
		t.Fatalf("GetTasksByTagID failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}
