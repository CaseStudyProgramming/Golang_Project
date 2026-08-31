-- This migration assigns existing tasks to a default user or makes them accessible
-- Uncomment and modify as needed for your specific use case

-- Option 1: Assign all existing tasks to a specific user (replace USER_ID with actual user ID)
-- UPDATE tasks SET user_id = USER_ID WHERE user_id IS NULL;

-- Option 2: Delete tasks without user_id (clean slate approach)
-- DELETE FROM tasks WHERE user_id IS NULL;

-- Option 3: Create a default system user and assign tasks to it
-- INSERT INTO users (name, email, password_hash) VALUES 
-- ('System User', 'system@example.com', '$2a$10$placeholder_hash_here')
-- RETURNING id;
-- Then update tasks with that user_id
