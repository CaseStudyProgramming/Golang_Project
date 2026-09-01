ALTER TABLE tasks ADD COLUMN progress_percentage INT DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100);
