-- Migration: Convert TIMESTAMP columns to BIGINT (epoch milliseconds)
-- This migration converts all timestamp columns from TIMESTAMP to BIGINT to store epoch milliseconds
-- This is important for multi-region support as epoch time is timezone-independent

-- Convert tasks table timestamps
-- Step 1: Remove default values first
ALTER TABLE tasks 
ALTER COLUMN created_at DROP DEFAULT,
ALTER COLUMN updated_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE tasks 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT,
ALTER COLUMN updated_at TYPE BIGINT USING (EXTRACT(EPOCH FROM updated_at) * 1000)::BIGINT,
ALTER COLUMN due_date TYPE BIGINT USING (EXTRACT(EPOCH FROM due_date) * 1000)::BIGINT,
ALTER COLUMN deleted_at TYPE BIGINT USING (EXTRACT(EPOCH FROM deleted_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE tasks 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
ALTER COLUMN updated_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert users table timestamps
-- Step 1: Remove default values first
ALTER TABLE users 
ALTER COLUMN created_at DROP DEFAULT,
ALTER COLUMN updated_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE users 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT,
ALTER COLUMN updated_at TYPE BIGINT USING (EXTRACT(EPOCH FROM updated_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE users 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
ALTER COLUMN updated_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert categories table timestamps
-- Step 1: Remove default values first
ALTER TABLE categories 
ALTER COLUMN created_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE categories 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE categories 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert tags table timestamps
-- Step 1: Remove default values first
ALTER TABLE tags 
ALTER COLUMN created_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE tags 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE tags 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert subtasks table timestamps
-- Step 1: Remove default values first
ALTER TABLE subtasks 
ALTER COLUMN created_at DROP DEFAULT,
ALTER COLUMN updated_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE subtasks 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT,
ALTER COLUMN updated_at TYPE BIGINT USING (EXTRACT(EPOCH FROM updated_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE subtasks 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
ALTER COLUMN updated_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert activity_logs table timestamps
-- Step 1: Remove default values first
ALTER TABLE activity_logs 
ALTER COLUMN created_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE activity_logs 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE activity_logs 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert task_tags junction table timestamps
-- Step 1: Remove default values first
ALTER TABLE task_tags 
ALTER COLUMN created_at DROP DEFAULT;

-- Step 2: Convert column types to BIGINT
ALTER TABLE task_tags 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Step 3: Set new default values with epoch time
ALTER TABLE task_tags 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;
