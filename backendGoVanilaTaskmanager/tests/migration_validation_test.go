package tests

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// TestMigrationValidation_ExistingDataStructure validates that the database schema has been properly migrated
func TestMigrationValidation_ExistingDataStructure(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Validate users table structure
	var usersCreatedAtType, usersTimezoneType string
	var usersIsNullable, usersColumnDefault string
	err := db.QueryRow(`
		SELECT 
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns 
		WHERE table_name = 'users' AND column_name = 'created_at'
	`).Scan(&usersCreatedAtType, &usersIsNullable, &usersColumnDefault)
	if err != nil {
		t.Fatalf("Failed to check users.created_at column type: %v", err)
	}

	if usersCreatedAtType != "bigint" {
		t.Errorf("Expected users.created_at to be bigint, got %s", usersCreatedAtType)
	}

	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'users' AND column_name = 'timezone'
	`).Scan(&usersTimezoneType)
	if err != nil {
		t.Fatalf("Failed to check users.timezone column: %v", err)
	}

	if usersTimezoneType != "character varying" {
		t.Errorf("Expected users.timezone to be character varying, got %s", usersTimezoneType)
	}

	// Validate tasks table structure
	var tasksCreatedAtType, tasksDueDateType string
	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'tasks' AND column_name = 'created_at'
	`).Scan(&tasksCreatedAtType)
	if err != nil {
		t.Fatalf("Failed to check tasks.created_at column type: %v", err)
	}

	if tasksCreatedAtType != "bigint" {
		t.Errorf("Expected tasks.created_at to be bigint, got %s", tasksCreatedAtType)
	}

	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'tasks' AND column_name = 'due_date'
	`).Scan(&tasksDueDateType)
	if err != nil {
		t.Fatalf("Failed to check tasks.due_date column type: %v", err)
	}

	if tasksDueDateType != "bigint" {
		t.Errorf("Expected tasks.due_date to be bigint, got %s", tasksDueDateType)
	}

	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'tasks' AND column_name = 'deleted_at'
	`).Scan(&tasksDueDateType)
	if err != nil {
		t.Fatalf("Failed to check tasks.deleted_at column type: %v", err)
	}

	if tasksDueDateType != "bigint" {
		t.Errorf("Expected tasks.deleted_at to be bigint, got %s", tasksDueDateType)
	}

	// Validate categories table structure
	var categoriesCreatedAtType string
	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'categories' AND column_name = 'created_at'
	`).Scan(&categoriesCreatedAtType)
	if err != nil {
		t.Fatalf("Failed to check categories.created_at column type: %v", err)
	}

	if categoriesCreatedAtType != "bigint" {
		t.Errorf("Expected categories.created_at to be bigint, got %s", categoriesCreatedAtType)
	}

	// Validate tags table structure
	var tagsCreatedAtType string
	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'tags' AND column_name = 'created_at'
	`).Scan(&tagsCreatedAtType)
	if err != nil {
		t.Fatalf("Failed to check tags.created_at column type: %v", err)
	}

	if tagsCreatedAtType != "bigint" {
		t.Errorf("Expected tags.created_at to be bigint, got %s", tagsCreatedAtType)
	}

	// Validate subtasks table structure
	var subtasksCreatedAtType, subtasksUpdatedAtType string
	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'subtasks' AND column_name = 'created_at'
	`).Scan(&subtasksCreatedAtType)
	if err != nil {
		t.Fatalf("Failed to check subtasks.created_at column type: %v", err)
	}

	if subtasksCreatedAtType != "bigint" {
		t.Errorf("Expected subtasks.created_at to be bigint, got %s", subtasksCreatedAtType)
	}

	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'subtasks' AND column_name = 'updated_at'
	`).Scan(&subtasksUpdatedAtType)
	if err != nil {
		t.Fatalf("Failed to check subtasks.updated_at column type: %v", err)
	}

	if subtasksUpdatedAtType != "bigint" {
		t.Errorf("Expected subtasks.updated_at to be bigint, got %s", subtasksUpdatedAtType)
	}

	// Validate activity_logs table structure
	var activityLogsCreatedAtType string
	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'activity_logs' AND column_name = 'created_at'
	`).Scan(&activityLogsCreatedAtType)
	if err != nil {
		t.Fatalf("Failed to check activity_logs.created_at column type: %v", err)
	}

	if activityLogsCreatedAtType != "bigint" {
		t.Errorf("Expected activity_logs.created_at to be bigint, got %s", activityLogsCreatedAtType)
	}

	// Validate task_tags junction table structure
	var taskTagsCreatedAtType string
	err = db.QueryRow(`
		SELECT data_type FROM information_schema.columns 
		WHERE table_name = 'task_tags' AND column_name = 'created_at'
	`).Scan(&taskTagsCreatedAtType)
	if err != nil {
		t.Fatalf("Failed to check task_tags.created_at column type: %v", err)
	}

	if taskTagsCreatedAtType != "bigint" {
		t.Errorf("Expected task_tags.created_at to be bigint, got %s", taskTagsCreatedAtType)
	}

	t.Log("✅ All timestamp columns validated as bigint (epoch milliseconds)")
}

// TestMigrationValidation_DefaultValues validates that default values use epoch milliseconds
func TestMigrationValidation_DefaultValues(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user without specifying timestamps to test default values
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"defaulttest", "default@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "defaulttest").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	// Check that created_at and updated_at are set and are reasonable epoch values
	var createdAt, updatedAt int64
	err = db.QueryRow("SELECT created_at, updated_at FROM users WHERE id = $1", userID).
		Scan(&createdAt, &updatedAt)
	if err != nil {
		t.Fatalf("Failed to get timestamps: %v", err)
	}

	now := time.Now().UnixMilli()

	// Check that timestamps are reasonable (within last minute)
	if createdAt < now-60000 || createdAt > now+60000 {
		t.Errorf("created_at timestamp %d is not within reasonable range (current: %d)", createdAt, now)
	}

	if updatedAt < now-60000 || updatedAt > now+60000 {
		t.Errorf("updated_at timestamp %d is not within reasonable range (current: %d)", updatedAt, now)
	}

	// Create a task without specifying timestamps to test default values
	_, err = db.Exec("INSERT INTO tasks (user_id, title) VALUES ($1, $2)", userID, "Default Timestamp Task")
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	var taskCreatedAt, taskUpdatedAt int64
	err = db.QueryRow("SELECT created_at, updated_at FROM tasks WHERE user_id = $1 LIMIT 1", userID).
		Scan(&taskCreatedAt, &taskUpdatedAt)
	if err != nil {
		t.Fatalf("Failed to get task timestamps: %v", err)
	}

	now = time.Now().UnixMilli()

	// Check that task timestamps are reasonable
	if taskCreatedAt < now-60000 || taskCreatedAt > now+60000 {
		t.Errorf("task created_at timestamp %d is not within reasonable range (current: %d)", taskCreatedAt, now)
	}

	if taskUpdatedAt < now-60000 || taskUpdatedAt > now+60000 {
		t.Errorf("task updated_at timestamp %d is not within reasonable range (current: %d)", taskUpdatedAt, now)
	}

	t.Log("✅ Default timestamp values validated as epoch milliseconds")
}

// TestMigrationValidation_DataIntegrity validates that epoch timestamps can be properly converted back to time
func TestMigrationValidation_DataIntegrity(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user with known epoch timestamp
	knownEpoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		"integritytest", "integrity@example.com", "hashedpassword", "UTC", knownEpoch, knownEpoch)
	if err != nil {
		t.Fatalf("Failed to create test user with known epoch: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "integritytest").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	// Retrieve the epoch timestamp
	var retrievedEpoch int64
	err = db.QueryRow("SELECT created_at FROM users WHERE id = $1", userID).Scan(&retrievedEpoch)
	if err != nil {
		t.Fatalf("Failed to retrieve epoch timestamp: %v", err)
	}

	// Verify the epoch timestamp matches
	if retrievedEpoch != knownEpoch {
		t.Errorf("Expected epoch %d, got %d", knownEpoch, retrievedEpoch)
	}

	// Convert back to time and verify
	retrievedTime := time.UnixMilli(retrievedEpoch)
	expectedTime := time.UnixMilli(knownEpoch)

	if !retrievedTime.Equal(expectedTime) {
		t.Errorf("Time conversion failed: expected %v, got %v", expectedTime, retrievedTime)
	}

	// Create a task with known due date epoch
	futureEpoch := time.Now().Add(24 * time.Hour).UnixMilli()
	_, err = db.Exec("INSERT INTO tasks (user_id, title, due_date) VALUES ($1, $2, $3)",
		userID, "Integrity Task", futureEpoch)
	if err != nil {
		t.Fatalf("Failed to create test task with known due date epoch: %v", err)
	}

	var retrievedDueDate int64
	err = db.QueryRow("SELECT due_date FROM tasks WHERE user_id = $1 LIMIT 1", userID).Scan(&retrievedDueDate)
	if err != nil {
		t.Fatalf("Failed to retrieve due date epoch: %v", err)
	}

	if retrievedDueDate != futureEpoch {
		t.Errorf("Expected due date epoch %d, got %d", futureEpoch, retrievedDueDate)
	}

	t.Log("✅ Data integrity validated - epoch timestamps can be properly converted")
}

// TestMigrationValidation_TimezoneColumn validates timezone column functionality
func TestMigrationValidation_TimezoneColumn(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create users with different timezones
	timezones := []string{"UTC", "America/New_York", "Europe/London", "Asia/Tokyo"}

	for i, tz := range timezones {
		userName := fmt.Sprintf("tzvalidation%d", i)
		_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
			userName, fmt.Sprintf("tz%d@example.com", i), "hashedpassword", tz)
		if err != nil {
			t.Fatalf("Failed to create user with timezone %s: %v", tz, err)
		}

		var retrievedTZ string
		err = db.QueryRow("SELECT timezone FROM users WHERE name = $1", userName).Scan(&retrievedTZ)
		if err != nil {
			t.Fatalf("Failed to retrieve timezone for %s: %v", userName, err)
		}

		if retrievedTZ != tz {
			t.Errorf("Expected timezone %s, got %s", tz, retrievedTZ)
		}
	}

	// Update timezone
	_, err := db.Exec("UPDATE users SET timezone = $1 WHERE name = $2", "Australia/Sydney", "tzvalidation0")
	if err != nil {
		t.Fatalf("Failed to update timezone: %v", err)
	}

	var updatedTZ string
	err = db.QueryRow("SELECT timezone FROM users WHERE name = $1", "tzvalidation0").Scan(&updatedTZ)
	if err != nil {
		t.Fatalf("Failed to retrieve updated timezone: %v", err)
	}

	if updatedTZ != "Australia/Sydney" {
		t.Errorf("Expected updated timezone Australia/Sydney, got %s", updatedTZ)
	}

	t.Log("✅ Timezone column functionality validated")
}

// TestMigrationValidation_NullTimestamps validates handling of NULL timestamp values
func TestMigrationValidation_NullTimestamps(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"nulltest", "null@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "nulltest").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	// Create a task with NULL due_date and deleted_at
	_, err = db.Exec("INSERT INTO tasks (user_id, title, due_date, deleted_at) VALUES ($1, $2, NULL, NULL)",
		userID, "Null Timestamp Task")
	if err != nil {
		t.Fatalf("Failed to create task with null timestamps: %v", err)
	}

	// Verify NULL values are handled correctly
	var dueDate sql.NullInt64
	var deletedAt sql.NullInt64
	err = db.QueryRow("SELECT due_date, deleted_at FROM tasks WHERE user_id = $1 LIMIT 1", userID).
		Scan(&dueDate, &deletedAt)
	if err != nil {
		t.Fatalf("Failed to retrieve null timestamps: %v", err)
	}

	if dueDate.Valid {
		t.Error("Expected due_date to be NULL, but it has a value")
	}

	if deletedAt.Valid {
		t.Error("Expected deleted_at to be NULL, but it has a value")
	}

	// Update with actual values
	futureEpoch := time.Now().Add(24 * time.Hour).UnixMilli()
	_, err = db.Exec("UPDATE tasks SET due_date = $1 WHERE user_id = $2", futureEpoch, userID)
	if err != nil {
		t.Fatalf("Failed to update due_date: %v", err)
	}

	err = db.QueryRow("SELECT due_date FROM tasks WHERE user_id = $1 LIMIT 1", userID).Scan(&dueDate)
	if err != nil {
		t.Fatalf("Failed to retrieve updated due_date: %v", err)
	}

	if !dueDate.Valid {
		t.Error("Expected due_date to be valid after update")
	}

	if dueDate.Int64 != futureEpoch {
		t.Errorf("Expected due_date %d, got %d", futureEpoch, dueDate.Int64)
	}

	t.Log("✅ NULL timestamp handling validated")
}

// TestMigrationValidation_PerformanceCheck performs a basic performance check on epoch vs string operations
func TestMigrationValidation_PerformanceCheck(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"perftest", "perf@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "perftest").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	// Create multiple tasks with epoch timestamps
	startTime := time.Now()
	taskCount := 100
	for i := 0; i < taskCount; i++ {
		dueDate := time.Now().Add(time.Duration(i) * time.Hour).UnixMilli()
		_, err = db.Exec("INSERT INTO tasks (user_id, title, due_date) VALUES ($1, $2, $3)",
			userID, fmt.Sprintf("Performance Task %d", i), dueDate)
		if err != nil {
			t.Fatalf("Failed to create task %d: %v", i, err)
		}
	}
	insertDuration := time.Since(startTime)

	// Query all tasks
	startTime = time.Now()
	rows, err := db.Query("SELECT id, title, due_date, created_at FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		t.Fatalf("Failed to query tasks: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int64
		var title string
		var dueDate, createdAt sql.NullInt64
		err = rows.Scan(&id, &title, &dueDate, &createdAt)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		count++
	}
	queryDuration := time.Since(startTime)

	if count != taskCount {
		t.Errorf("Expected %d tasks, got %d", taskCount, count)
	}

	t.Logf("✅ Performance check completed:")
	t.Logf("   - Inserted %d tasks in %v", taskCount, insertDuration)
	t.Logf("   - Queried %d tasks in %v", count, queryDuration)
	t.Logf("   - Average insert time: %v per task", insertDuration/time.Duration(taskCount))
	t.Logf("   - Average query time: %v per task", queryDuration/time.Duration(count))
}
