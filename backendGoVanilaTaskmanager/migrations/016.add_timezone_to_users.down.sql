-- Rollback: Remove timezone column from users table

DROP INDEX IF EXISTS idx_users_timezone;

ALTER TABLE users 
DROP COLUMN IF EXISTS timezone;
