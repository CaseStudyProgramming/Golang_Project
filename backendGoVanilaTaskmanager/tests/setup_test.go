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

// SetupTestDB creates a connection to the test database
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
	return db
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
	_, err := db.Exec("DELETE FROM tasks WHERE title LIKE 'Test%'")
	if err != nil {
		t.Logf("Warning: Failed to cleanup test tasks: %v", err)
	}
	_, err = db.Exec("DELETE FROM users WHERE name LIKE 'testuser%'")
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
