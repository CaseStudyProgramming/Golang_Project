package models

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupActivityLogTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Berjuang#382 dbname=taskmanager_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: database not reachable")
	}

	// Clean up and create tables
	_, _ = db.Exec("DROP TABLE IF EXISTS activity_logs CASCADE")
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

	_, err = db.Exec(`CREATE TABLE activity_logs (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		task_id BIGINT,
		action VARCHAR(50) NOT NULL,
		entity_type VARCHAR(50) NOT NULL,
		entity_id BIGINT,
		details TEXT,
		ip_address VARCHAR(45),
		user_agent TEXT,
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("Failed to create activity_logs table: %v", err)
	}

	return db
}

func cleanupActivityLogTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestNewActivityLogModel(t *testing.T) {
	db := &sql.DB{}
	model := NewActivityLogModel(db)

	if model == nil {
		t.Fatal("NewActivityLogModel returned nil")
	}

	if model.DB != db {
		t.Error("NewActivityLogModel did not set DB correctly")
	}
}

func TestActivityLogModel_Create(t *testing.T) {
	db := setupActivityLogTestDB(t)
	defer cleanupActivityLogTestDB(db)

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

	model := NewActivityLogModel(db)

	log := &ActivityLog{
		UserID:     createdUser.ID,
		Action:     ActionCreate,
		EntityType: EntityTask,
		EntityID:   &[]int64{1}[0],
		Details:    "Created a test task",
		IPAddress:  "127.0.0.1",
		UserAgent:  "Test Agent",
	}

	createdLog, err := model.Create(log)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdLog.ID == 0 {
		t.Error("Create did not set log ID")
	}

	if createdLog.Action != ActionCreate {
		t.Errorf("Expected action %s, got %s", ActionCreate, createdLog.Action)
	}

	if createdLog.EntityType != EntityTask {
		t.Errorf("Expected entity type %s, got %s", EntityTask, createdLog.EntityType)
	}
}

func TestActivityLogModel_GetByUserID(t *testing.T) {
	db := setupActivityLogTestDB(t)
	defer cleanupActivityLogTestDB(db)

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

	model := NewActivityLogModel(db)

	// Create multiple activity logs
	for i := 0; i < 3; i++ {
		log := &ActivityLog{
			UserID:     createdUser.ID,
			Action:     ActionCreate,
			EntityType: EntityTask,
			EntityID:   &[]int64{int64(i + 1)}[0],
			Details:    "Test log entry",
		}
		_, err := model.Create(log)
		if err != nil {
			t.Fatalf("Failed to create activity log: %v", err)
		}
	}

	// Test GetByUserID
	logs, total, err := model.GetByUserID(createdUser.ID, 0, 10)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total count 3, got %d", total)
	}

	if len(logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}

	// Test pagination
	logs, total, err = model.GetByUserID(createdUser.ID, 0, 2)
	if err != nil {
		t.Fatalf("GetByUserID with pagination failed: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs with pagination, got %d", len(logs))
	}
}

func TestActivityLogModel_GetByTaskID(t *testing.T) {
	db := setupActivityLogTestDB(t)
	defer cleanupActivityLogTestDB(db)

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

	model := NewActivityLogModel(db)

	taskID := int64(123)

	// Create activity logs for the same task
	for i := 0; i < 2; i++ {
		log := &ActivityLog{
			UserID:     createdUser.ID,
			TaskID:     &taskID,
			Action:     ActionUpdate,
			EntityType: EntityTask,
			EntityID:   &taskID,
			Details:    "Updated task",
		}
		_, err := model.Create(log)
		if err != nil {
			t.Fatalf("Failed to create activity log: %v", err)
		}
	}

	// Test GetByTaskID
	logs, total, err := model.GetByTaskID(taskID, 0, 10)
	if err != nil {
		t.Fatalf("GetByTaskID failed: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected total count 2, got %d", total)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}
}

func TestActivityLogModel_GetByID(t *testing.T) {
	db := setupActivityLogTestDB(t)
	defer cleanupActivityLogTestDB(db)

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

	model := NewActivityLogModel(db)

	log := &ActivityLog{
		UserID:     createdUser.ID,
		Action:     ActionDelete,
		EntityType: EntityCategory,
		EntityID:   &[]int64{456}[0],
		Details:    "Deleted category",
	}

	createdLog, err := model.Create(log)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByID
	retrievedLog, err := model.GetByID(createdLog.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedLog.ID != createdLog.ID {
		t.Errorf("Expected ID %d, got %d", createdLog.ID, retrievedLog.ID)
	}

	if retrievedLog.Action != ActionDelete {
		t.Errorf("Expected action %s, got %s", ActionDelete, retrievedLog.Action)
	}

	// Test with non-existent ID
	_, err = model.GetByID(99999)
	if err == nil {
		t.Error("Expected error when getting non-existent activity log")
	}
}
