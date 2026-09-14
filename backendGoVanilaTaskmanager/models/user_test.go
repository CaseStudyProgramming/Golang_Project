package models

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupUserTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Berjuang#382 dbname=taskmanager_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: database not reachable")
	}

	// Clean up and create tables
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

	return db
}

func cleanupUserTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestNewUserModel(t *testing.T) {
	db := &sql.DB{}
	model := NewUserModel(db)
	
	if model == nil {
		t.Fatal("NewUserModel returned nil")
	}
	
	if model.DB != db {
		t.Error("NewUserModel did not set DB correctly")
	}
}

func TestUserModel_Create(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	model := NewUserModel(db)
	
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Timezone:     "UTC",
	}

	createdUser, err := model.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdUser.ID == 0 {
		t.Error("Create did not set user ID")
	}

	if createdUser.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", createdUser.Name)
	}

	if createdUser.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", createdUser.Email)
	}

	if createdUser.Timezone != "UTC" {
		t.Errorf("Expected timezone 'UTC', got '%s'", createdUser.Timezone)
	}
}

func TestUserModel_GetByEmail(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	model := NewUserModel(db)
	
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	_, err := model.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByEmail
	retrievedUser, err := model.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}

	if retrievedUser.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", retrievedUser.Email)
	}

	if retrievedUser.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", retrievedUser.Name)
	}

	// Test with non-existent email
	_, err = model.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestUserModel_GetByID(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	model := NewUserModel(db)
	
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	createdUser, err := model.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByID
	retrievedUser, err := model.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected ID %d, got %d", createdUser.ID, retrievedUser.ID)
	}

	if retrievedUser.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", retrievedUser.Email)
	}

	// Test with non-existent ID
	_, err = model.GetByID(99999)
	if err == nil {
		t.Error("Expected error when getting non-existent user by ID")
	}
}

func TestUserModel_Update(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	model := NewUserModel(db)
	
	user := &User{
		Name:         "Original Name",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Timezone:     "UTC",
	}

	createdUser, err := model.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update user
	createdUser.Name = "Updated Name"
	createdUser.Timezone = "Asia/Jakarta"
	err = model.Update(createdUser)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrievedUser, err := model.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if retrievedUser.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", retrievedUser.Name)
	}

	if retrievedUser.Timezone != "Asia/Jakarta" {
		t.Errorf("Expected timezone 'Asia/Jakarta', got '%s'", retrievedUser.Timezone)
	}
}

func TestUserModel_Delete(t *testing.T) {
	db := setupUserTestDB(t)
	defer cleanupUserTestDB(db)

	model := NewUserModel(db)
	
	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	createdUser, err := model.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete user
	err = model.Delete(createdUser.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify user is deleted
	_, err = model.GetByID(createdUser.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}