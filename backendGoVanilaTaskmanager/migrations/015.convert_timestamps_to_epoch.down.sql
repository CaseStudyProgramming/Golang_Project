-- Rollback: Convert BIGINT columns back to TIMESTAMP
-- This rollback converts all epoch milliseconds back to TIMESTAMP format

-- Revert tasks table timestamps
ALTER TABLE tasks 
ALTER COLUMN created_at TYPE TIMESTAMP USING to_timestamp(created_at / 1000.0),
ALTER COLUMN updated_at TYPE TIMESTAMP USING to_timestamp(updated_at / 1000.0),
ALTER COLUMN due_date TYPE TIMESTAMP USING to_timestamp(due_date / 1000.0),
ALTER COLUMN deleted_at TYPE TIMESTAMP USING to_timestamp(deleted_at / 1000.0);

-- Revert default values for tasks table
ALTER TABLE tasks 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

-- Revert users table timestamps
ALTER TABLE users 
ALTER COLUMN created_at TYPE TIMESTAMP USING to_timestamp(created_at / 1000.0),
ALTER COLUMN updated_at TYPE TIMESTAMP USING to_timestamp(updated_at / 1000.0);

-- Revert default values for users table
ALTER TABLE users 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

-- Revert categories table timestamps
ALTER TABLE categories 
ALTER COLUMN created_at TYPE TIMESTAMP USING to_timestamp(created_at / 1000.0);

-- Revert default values for categories table
ALTER TABLE categories 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;

-- Revert tags table timestamps
ALTER TABLE tags 
ALTER COLUMN created_at TYPE TIMESTAMP USING to_timestamp(created_at / 1000.0);

-- Revert default values for tags table
ALTER TABLE tags 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;

-- Revert subtasks table timestamps
ALTER TABLE subtasks 
ALTER COLUMN created_at TYPE TIMESTAMP USING to_timestamp(created_at / 1000.0),
ALTER COLUMN updated_at TYPE TIMESTAMP USING to_timestamp(updated_at / 1000.0);

-- Revert default values for subtasks table
ALTER TABLE subtasks 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

-- Revert activity_logs table timestamps
ALTER TABLE activity_logs 
ALTER COLUMN created_at TYPE TIMESTAMP USING to_timestamp(created_at / 1000.0);

-- Revert default values for activity_logs table
ALTER TABLE activity_logs 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;
