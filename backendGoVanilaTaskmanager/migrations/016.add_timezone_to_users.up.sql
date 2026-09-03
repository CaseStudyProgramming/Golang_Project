-- Migration: Add timezone column to users table
-- This adds a timezone column to store user timezone preferences

ALTER TABLE users 
ADD COLUMN timezone VARCHAR(50) DEFAULT 'UTC';

-- Add index for timezone column for potential queries
CREATE INDEX idx_users_timezone ON users(timezone);
