package models

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupCategoryTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=Berjuang#382 dbname=taskmanager_test sslmode=disable")
	if err != nil {
		t.Skip("Skipping test: database not available")
	}

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: database not reachable")
	}

	// Clean up and create tables
	_, _ = db.Exec("DROP TABLE IF EXISTS categories CASCADE")
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

	_, err = db.Exec(`CREATE TABLE categories (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		name VARCHAR(100) NOT NULL,
		color_hex VARCHAR(7),
		created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		t.Fatalf("Failed to create categories table: %v", err)
	}

	return db
}

func cleanupCategoryTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestNewCategoryModel(t *testing.T) {
	db := &sql.DB{}
	model := NewCategoryModel(db)

	if model == nil {
		t.Fatal("NewCategoryModel returned nil")
	}

	if model.DB != db {
		t.Error("NewCategoryModel did not set DB correctly")
	}
}

func TestCategoryModel_Create(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer cleanupCategoryTestDB(db)

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

	model := NewCategoryModel(db)

	category := &Category{
		UserID:   createdUser.ID,
		Name:     "Work",
		ColorHex: "#FF5733",
	}

	createdCategory, err := model.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdCategory.ID == 0 {
		t.Error("Create did not set category ID")
	}

	if createdCategory.Name != "Work" {
		t.Errorf("Expected name 'Work', got '%s'", createdCategory.Name)
	}

	if createdCategory.ColorHex != "#FF5733" {
		t.Errorf("Expected color hex '#FF5733', got '%s'", createdCategory.ColorHex)
	}
}

func TestCategoryModel_GetAll(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer cleanupCategoryTestDB(db)

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

	model := NewCategoryModel(db)

	// Create multiple categories
	categoryNames := []string{"Work", "Personal", "Shopping"}
	for _, name := range categoryNames {
		category := &Category{
			UserID:   createdUser.ID,
			Name:     name,
			ColorHex: "#FF5733",
		}
		_, err := model.Create(category)
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}
	}

	// Test GetAll
	categories, err := model.GetAll(createdUser.ID)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(categories))
	}
}

func TestCategoryModel_GetByID(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer cleanupCategoryTestDB(db)

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

	model := NewCategoryModel(db)

	category := &Category{
		UserID:   createdUser.ID,
		Name:     "Work",
		ColorHex: "#FF5733",
	}

	createdCategory, err := model.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Test GetByID
	retrievedCategory, err := model.GetByID(createdUser.ID, createdCategory.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if retrievedCategory.ID != createdCategory.ID {
		t.Errorf("Expected ID %d, got %d", createdCategory.ID, retrievedCategory.ID)
	}

	if retrievedCategory.Name != "Work" {
		t.Errorf("Expected name 'Work', got '%s'", retrievedCategory.Name)
	}

	// Test with non-existent ID
	_, err = model.GetByID(createdUser.ID, 99999)
	if err == nil {
		t.Error("Expected error when getting non-existent category")
	}
}

func TestCategoryModel_Update(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer cleanupCategoryTestDB(db)

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

	model := NewCategoryModel(db)

	category := &Category{
		UserID:   createdUser.ID,
		Name:     "Work",
		ColorHex: "#FF5733",
	}

	createdCategory, err := model.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update category
	createdCategory.Name = "Updated Work"
	createdCategory.ColorHex = "#00FF00"
	err = model.Update(createdUser.ID, createdCategory)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrievedCategory, err := model.GetByID(createdUser.ID, createdCategory.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if retrievedCategory.Name != "Updated Work" {
		t.Errorf("Expected name 'Updated Work', got '%s'", retrievedCategory.Name)
	}

	if retrievedCategory.ColorHex != "#00FF00" {
		t.Errorf("Expected color hex '#00FF00', got '%s'", retrievedCategory.ColorHex)
	}
}

func TestCategoryModel_Delete(t *testing.T) {
	db := setupCategoryTestDB(t)
	defer cleanupCategoryTestDB(db)

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

	model := NewCategoryModel(db)

	category := &Category{
		UserID:   createdUser.ID,
		Name:     "Work",
		ColorHex: "#FF5733",
	}

	createdCategory, err := model.Create(category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete category
	err = model.Delete(createdUser.ID, createdCategory.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify category is deleted
	_, err = model.GetByID(createdUser.ID, createdCategory.ID)
	if err == nil {
		t.Error("Expected error when getting deleted category")
	}
}
