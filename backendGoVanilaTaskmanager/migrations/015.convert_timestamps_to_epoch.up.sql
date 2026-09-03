-- Migration: Convert TIMESTAMP columns to BIGINT (epoch milliseconds)
-- This migration converts all timestamp columns from TIMESTAMP to BIGINT to store epoch milliseconds

-- Convert tasks table timestamps
ALTER TABLE tasks 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT,
ALTER COLUMN updated_at TYPE BIGINT USING (EXTRACT(EPOCH FROM updated_at) * 1000)::BIGINT,
ALTER COLUMN due_date TYPE BIGINT USING (EXTRACT(EPOCH FROM due_date) * 1000)::BIGINT,
ALTER COLUMN deleted_at TYPE BIGINT USING (EXTRACT(EPOCH FROM deleted_at) * 1000)::BIGINT;

-- Update default values for tasks table
ALTER TABLE tasks 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
ALTER COLUMN updated_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert users table timestamps
ALTER TABLE users 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT,
ALTER COLUMN updated_at TYPE BIGINT USING (EXTRACT(EPOCH FROM updated_at) * 1000)::BIGINT;

-- Update default values for users table
ALTER TABLE users 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
ALTER COLUMN updated_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert categories table timestamps
ALTER TABLE categories 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Update default values for categories table
ALTER TABLE categories 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert tags table timestamps
ALTER TABLE tags 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Update default values for tags table
ALTER TABLE tags 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert subtasks table timestamps
ALTER TABLE subtasks 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT,
ALTER COLUMN updated_at TYPE BIGINT USING (EXTRACT(EPOCH FROM updated_at) * 1000)::BIGINT;

-- Update default values for subtasks table
ALTER TABLE subtasks 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT,
ALTER COLUMN updated_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;

-- Convert activity_logs table timestamps
ALTER TABLE activity_logs 
ALTER COLUMN created_at TYPE BIGINT USING (EXTRACT(EPOCH FROM created_at) * 1000)::BIGINT;

-- Update default values for activity_logs table
ALTER TABLE activity_logs 
ALTER COLUMN created_at SET DEFAULT (EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000)::BIGINT;
