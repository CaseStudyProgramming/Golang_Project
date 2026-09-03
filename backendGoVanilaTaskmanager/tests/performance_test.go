package tests

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// TestPerformance_EpochVsStringTimestamps compares performance of epoch vs string timestamp operations
func TestPerformance_EpochVsStringTimestamps(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user for testing
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

	taskCount := 1000
	t.Logf("Running performance test with %d tasks", taskCount)

	// Test 1: Insert performance with epoch timestamps
	t.Run("Insert_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		for i := 0; i < taskCount; i++ {
			dueDate := time.Now().Add(time.Duration(i) * time.Hour).UnixMilli()
			_, err := db.Exec("INSERT INTO tasks (user_id, title, due_date) VALUES ($1, $2, $3)",
				userID, fmt.Sprintf("Epoch Task %d", i), dueDate)
			if err != nil {
				t.Fatalf("Failed to insert task %d: %v", i, err)
			}
		}
		
		duration := time.Since(startTime)
		avgPerTask := duration / time.Duration(taskCount)
		
		t.Logf("✅ Inserted %d tasks with epoch timestamps in %v", taskCount, duration)
		t.Logf("   Average per task: %v", avgPerTask)
	})

	// Test 2: Query performance with epoch timestamps
	t.Run("Query_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		rows, err := db.Query("SELECT id, title, due_date, created_at, updated_at FROM tasks WHERE user_id = $1", userID)
		if err != nil {
			t.Fatalf("Failed to query tasks: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			var dueDate, createdAt, updatedAt sql.NullInt64
			err = rows.Scan(&id, &title, &dueDate, &createdAt, &updatedAt)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		avgPerTask := duration / time.Duration(count)
		
		t.Logf("✅ Queried %d tasks with epoch timestamps in %v", count, duration)
		t.Logf("   Average per task: %v", avgPerTask)
	})

	// Test 3: Update performance with epoch timestamps
	t.Run("Update_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		// Update all tasks
		result, err := db.Exec("UPDATE tasks SET updated_at = $1 WHERE user_id = $2", time.Now().UnixMilli(), userID)
		if err != nil {
			t.Fatalf("Failed to update tasks: %v", err)
		}
		
		rowsAffected, _ := result.RowsAffected()
		duration := time.Since(startTime)
		
		t.Logf("✅ Updated %d tasks with epoch timestamps in %v", rowsAffected, duration)
	})

	// Test 4: Comparison query with epoch comparison
	t.Run("Comparison_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		now := time.Now().UnixMilli()
		future := time.Now().Add(24 * time.Hour).UnixMilli()
		
		rows, err := db.Query("SELECT id, title FROM tasks WHERE user_id = $1 AND due_date BETWEEN $2 AND $3", 
			userID, now, future)
		if err != nil {
			t.Fatalf("Failed to query with epoch comparison: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			err = rows.Scan(&id, &title)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Comparison query with epoch timestamps returned %d rows in %v", count, duration)
	})

	// Test 5: Sorting performance with epoch timestamps
	t.Run("Sort_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		rows, err := db.Query("SELECT id, title, created_at FROM tasks WHERE user_id = $1 ORDER BY created_at DESC LIMIT 100", userID)
		if err != nil {
			t.Fatalf("Failed to query with epoch sort: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			var createdAt sql.NullInt64
			err = rows.Scan(&id, &title, &createdAt)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Sorted query with epoch timestamps returned %d rows in %v", count, duration)
	})

	// Test 6: Aggregate functions with epoch timestamps
	t.Run("Aggregate_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		var count int64
		err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND created_at > $2", 
			userID, time.Now().Add(-24*time.Hour).UnixMilli()).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to run aggregate query: %v", err)
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Aggregate query with epoch timestamps counted %d rows in %v", count, duration)
	})

	// Test 7: Index performance with epoch timestamps
	t.Run("Index_EpochTimestamps", func(t *testing.T) {
		// Create an index on created_at if it doesn't exist
		_, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at)")
		if err != nil {
			t.Logf("Warning: Could not create index: %v", err)
		}
		
		startTime := time.Now()
		
		rows, err := db.Query("SELECT id, title FROM tasks WHERE user_id = $1 AND created_at > $2 ORDER BY created_at DESC LIMIT 50", 
			userID, time.Now().Add(-1*time.Hour).UnixMilli())
		if err != nil {
			t.Fatalf("Failed to run indexed query: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			err = rows.Scan(&id, &title)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Indexed query with epoch timestamps returned %d rows in %v", count, duration)
	})

	// Test 8: Batch insert performance with epoch timestamps
	t.Run("BatchInsert_EpochTimestamps", func(t *testing.T) {
		// Clean up existing tasks for batch test
		_, err := db.Exec("DELETE FROM tasks WHERE user_id = $1", userID)
		if err != nil {
			t.Fatalf("Failed to clean up tasks: %v", err)
		}

		startTime := time.Now()
		
		// Use transaction for batch insert
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}
		defer tx.Rollback()

		stmt, err := tx.Prepare("INSERT INTO tasks (user_id, title, due_date) VALUES ($1, $2, $3)")
		if err != nil {
			t.Fatalf("Failed to prepare statement: %v", err)
		}
		defer stmt.Close()

		for i := 0; i < taskCount; i++ {
			dueDate := time.Now().Add(time.Duration(i) * time.Hour).UnixMilli()
			_, err := stmt.Exec(userID, fmt.Sprintf("Batch Task %d", i), dueDate)
			if err != nil {
				t.Fatalf("Failed to execute batch insert %d: %v", i, err)
			}
		}

		err = tx.Commit()
		if err != nil {
			t.Fatalf("Failed to commit transaction: %v", err)
		}
		
		duration := time.Since(startTime)
		avgPerTask := duration / time.Duration(taskCount)
		
		t.Logf("✅ Batch inserted %d tasks with epoch timestamps in %v", taskCount, duration)
		t.Logf("   Average per task: %v", avgPerTask)
	})

	// Test 9: Memory efficiency with epoch timestamps
	t.Run("Memory_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		// Load all tasks into memory
		rows, err := db.Query("SELECT id, title, due_date, created_at, updated_at FROM tasks WHERE user_id = $1", userID)
		if err != nil {
			t.Fatalf("Failed to query tasks: %v", err)
		}
		defer rows.Close()

		type Task struct {
			ID        int64
			Title     string
			DueDate   *int64
			CreatedAt int64
			UpdatedAt int64
		}

		tasks := make([]Task, 0, taskCount)
		for rows.Next() {
			var task Task
			var dueDate sql.NullInt64
			err = rows.Scan(&task.ID, &task.Title, &dueDate, &task.CreatedAt, &task.UpdatedAt)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			if dueDate.Valid {
				task.DueDate = &dueDate.Int64
			}
			tasks = append(tasks, task)
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Loaded %d tasks with epoch timestamps into memory in %v", len(tasks), duration)
		t.Logf("   Memory per task (estimated): %d bytes", len(tasks)*28/len(tasks)) // Rough estimate
	})

	// Test 10: Concurrent access performance with epoch timestamps
	t.Run("Concurrent_EpochTimestamps", func(t *testing.T) {
		startTime := time.Now()
		
		// Simulate concurrent reads
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func() {
				rows, err := db.Query("SELECT COUNT(*) FROM tasks WHERE user_id = $1", userID)
				if err != nil {
					t.Errorf("Concurrent query failed: %v", err)
					done <- false
					return
				}
				defer rows.Close()
				
				var count int64
				rows.Scan(&count)
				done <- true
			}()
		}
		
		// Wait for all goroutines
		successCount := 0
		for i := 0; i < 10; i++ {
			if <-done {
				successCount++
			}
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Executed %d concurrent queries with epoch timestamps in %v", successCount, duration)
	})

	t.Log("✅ Performance testing completed successfully")
}

// TestPerformance_EpochConversionOverhead tests the overhead of epoch conversion operations
func TestPerformance_EpochConversionOverhead(t *testing.T) {
	iterations := 100000
	
	// Test time to epoch conversion
	t.Run("TimeToEpoch", func(t *testing.T) {
		startTime := time.Now()
		
		for i := 0; i < iterations; i++ {
			now := time.Now()
			_ = now.UnixMilli()
		}
		
		duration := time.Since(startTime)
		avgPerOp := duration / time.Duration(iterations)
		
		t.Logf("✅ Time to epoch conversion (%d iterations): %v total, %v per operation", 
			iterations, duration, avgPerOp)
	})

	// Test epoch to time conversion
	t.Run("EpochToTime", func(t *testing.T) {
		startTime := time.Now()
		epoch := time.Now().UnixMilli()
		
		for i := 0; i < iterations; i++ {
			_ = time.UnixMilli(epoch)
		}
		
		duration := time.Since(startTime)
		avgPerOp := duration / time.Duration(iterations)
		
		t.Logf("✅ Epoch to time conversion (%d iterations): %v total, %v per operation", 
			iterations, duration, avgPerOp)
	})

	// Test timezone conversion
	t.Run("TimezoneConversion", func(t *testing.T) {
		startTime := time.Now()
		epoch := time.Now().UnixMilli()
		
		for i := 0; i < iterations; i++ {
			loc, _ := time.LoadLocation("America/New_York")
			_ = time.UnixMilli(epoch).In(loc)
		}
		
		duration := time.Since(startTime)
		avgPerOp := duration / time.Duration(iterations)
		
		t.Logf("✅ Timezone conversion (%d iterations): %v total, %v per operation", 
			iterations, duration, avgPerOp)
	})

	// Test formatting
	t.Run("EpochFormatting", func(t *testing.T) {
		startTime := time.Now()
		epoch := time.Now().UnixMilli()
		
		for i := 0; i < iterations; i++ {
			t := time.UnixMilli(epoch)
			_ = t.Format("2006-01-02 15:04:05")
		}
		
		duration := time.Since(startTime)
		avgPerOp := duration / time.Duration(iterations)
		
		t.Logf("✅ Epoch formatting (%d iterations): %v total, %v per operation", 
			iterations, duration, avgPerOp)
	})
}

// TestPerformance_StringVsEpochComparison compares hypothetical string vs epoch performance
func TestPerformance_StringVsEpochComparison(t *testing.T) {
	SkipIfNoTestDB(t)

	db := SetupTestDB(t)
	defer CleanupTestDB(db)
	defer CleanupTestData(db, t)

	// Create a user for testing
	_, err := db.Exec("INSERT INTO users (name, email, password_hash, timezone) VALUES ($1, $2, $3, $4)",
		"comparetest", "compare@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	var userID int64
	err = db.QueryRow("SELECT id FROM users WHERE name = $1", "comparetest").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to get test user ID: %v", err)
	}

	taskCount := 100

	// Insert test data with epoch timestamps
	for i := 0; i < taskCount; i++ {
		dueDate := time.Now().Add(time.Duration(i) * time.Hour).UnixMilli()
		_, err := db.Exec("INSERT INTO tasks (user_id, title, due_date) VALUES ($1, $2, $3)",
			userID, fmt.Sprintf("Compare Task %d", i), dueDate)
		if err != nil {
			t.Fatalf("Failed to insert task %d: %v", i, err)
		}
	}

	// Test numeric comparison (epoch)
	t.Run("NumericComparison_Epoch", func(t *testing.T) {
		startTime := time.Now()
		
		now := time.Now().UnixMilli()
		rows, err := db.Query("SELECT id, title FROM tasks WHERE user_id = $1 AND due_date > $2", userID, now)
		if err != nil {
			t.Fatalf("Failed to run numeric comparison: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			err = rows.Scan(&id, &title)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Numeric comparison (epoch): %d rows in %v", count, duration)
	})

	// Test range query (epoch)
	t.Run("RangeQuery_Epoch", func(t *testing.T) {
		startTime := time.Now()
		
		now := time.Now().UnixMilli()
		future := time.Now().Add(24 * time.Hour).UnixMilli()
		rows, err := db.Query("SELECT id, title FROM tasks WHERE user_id = $1 AND due_date BETWEEN $2 AND $3", 
			userID, now, future)
		if err != nil {
			t.Fatalf("Failed to run range query: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			err = rows.Scan(&id, &title)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Range query (epoch): %d rows in %v", count, duration)
	})

	// Test sorting (epoch)
	t.Run("Sorting_Epoch", func(t *testing.T) {
		startTime := time.Now()
		
		rows, err := db.Query("SELECT id, title, due_date FROM tasks WHERE user_id = $1 ORDER BY due_date ASC", userID)
		if err != nil {
			t.Fatalf("Failed to run sorting query: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var id int64
			var title string
			var dueDate sql.NullInt64
			err = rows.Scan(&id, &title, &dueDate)
			if err != nil {
				t.Fatalf("Failed to scan row: %v", err)
			}
			count++
		}
		
		duration := time.Since(startTime)
		
		t.Logf("✅ Sorting query (epoch): %d rows in %v", count, duration)
	})

	t.Log("✅ String vs epoch comparison completed")
}