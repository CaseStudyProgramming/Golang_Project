package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// TestDBConfig holds the test database configuration
type TestDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetTestDBConfig returns the test database configuration from environment variables or defaults
func GetTestDBConfig() TestDBConfig {
	return TestDBConfig{
		Host:     getEnv("TEST_DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("TEST_DB_USER", "postgres"),
		Password: getEnv("TEST_DB_PASSWORD", "Berjuang#382"),
		DBName:   getEnv("TEST_DB_NAME", "taskmanager_test"),
		SSLMode:  getEnv("TEST_DB_SSLMODE", "disable"),
	}
}

// SetupTestDB creates a connection to the test database and runs migrations
func SetupTestDB(t *testing.T) *sql.DB {
	cfg := GetTestDBConfig()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	log.Println("✅ Connected to test database")

	// Run migrations for test database
	RunTestMigrations(db, t)

	return db
}

// RunTestMigrations runs the database migrations for testing
func RunTestMigrations(db *sql.DB, t *testing.T) {
	// Drop existing tables to ensure clean state
	tables := []string{"activity_logs", "subtasks", "task_tags", "tasks", "tags", "categories", "users"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: Failed to drop table %s: %v", table, err)
		}
	}

	migrations := []string{
		// Users table with timezone column and epoch timestamps
		`CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(150) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			timezone VARCHAR(50) DEFAULT 'UTC',
			created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
			updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
		)`,
		// Categories table with epoch timestamps
		`CREATE TABLE categories (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			color_hex VARCHAR(7) DEFAULT '#3B82F6',
			created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
		)`,
		// Tags table with epoch timestamps
		`CREATE TABLE tags (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			color_hex VARCHAR(7) DEFAULT '#3B82F6',
			created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
		)`,
		// Tasks table with epoch timestamps
		`CREATE TABLE tasks (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			category_id INT NULL REFERENCES categories(id) ON DELETE SET NULL,
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
		)`,
		// Task-tags junction table with epoch timestamps
		`CREATE TABLE task_tags (
			task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
			created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
			PRIMARY KEY (task_id, tag_id)
		)`,
		// Subtasks table with epoch timestamps
		`CREATE TABLE subtasks (
			id SERIAL PRIMARY KEY,
			task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			is_completed BOOLEAN DEFAULT false,
			created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
			updated_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
		)`,
		// Activity logs table with epoch timestamps
		`CREATE TABLE activity_logs (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			task_id INT NULL REFERENCES tasks(id) ON DELETE SET NULL,
			action VARCHAR(20) NOT NULL,
			entity_type VARCHAR(20) NOT NULL,
			entity_id INT NULL,
			details TEXT NULL,
			ip_address VARCHAR(45) NULL,
			user_agent TEXT NULL,
			created_at BIGINT DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT
		)`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			t.Fatalf("Migration failed: %v", err)
		}
	}

	log.Println("✅ Test database migrations completed")
}

// CleanupTestDB closes the test database connection
func CleanupTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("✅ Test database connection closed")
	}
}

// CleanupTestData removes all test data from the database
func CleanupTestData(db *sql.DB, t *testing.T) {
	_, err := db.Exec("DELETE FROM task_tags")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test task_tags: %v", err)
	}
	_, err = db.Exec("DELETE FROM tags WHERE name LIKE 'Test%'")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test tags: %v", err)
	}
	_, err = db.Exec("DELETE FROM categories WHERE name LIKE 'Test%'")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test categories: %v", err)
	}
	_, err = db.Exec("DELETE FROM tasks WHERE title LIKE 'Test%' OR title LIKE 'Epoch%' OR title LIKE 'Batch%' OR title LIKE 'Compare%' OR title LIKE 'Performance%'")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test tasks: %v", err)
	}
	_, err = db.Exec("DELETE FROM users WHERE name LIKE 'testuser%' OR name LIKE 'timezoneuser%' OR name LIKE 'tzupdateuser%' OR name LIKE 'invalidtzuser%' OR name LIKE 'eastuser' OR name LIKE 'westuser' OR name LIKE 'subtestuser%' OR name LIKE 'user%' OR name LIKE 'perftest' OR name LIKE 'defaulttest' OR name LIKE 'integritytest' OR name LIKE 'nulltest' OR name LIKE 'tzvalidation%' OR name LIKE 'comparetest'")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test users: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SkipIfNoTestDB skips the test if test database is not available
func SkipIfNoTestDB(t *testing.T) {
	cfg := GetTestDBConfig()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		t.Skip("Skipping test: test database not available")
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping test: test database not reachable")
	}
}
