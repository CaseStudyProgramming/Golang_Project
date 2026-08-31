# Database Migrations

This directory contains SQL migration files for the Task Manager database.

## Recent Migrations for Authentication & Multi-Tenancy

### 004.create_users_table.up.sql
Creates the `users` table for user authentication:
- `id` (SERIAL PRIMARY KEY)
- `name` (VARCHAR(100) NOT NULL)
- `email` (VARCHAR(150) UNIQUE NOT NULL)
- `password_hash` (VARCHAR(255) NOT NULL)
- `created_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
- `updated_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP)

### 005.add_user_id_to_tasks.up.sql
Adds foreign key relationship to tasks table:
- `user_id` (INT NULL REFERENCES users(id) ON DELETE CASCADE)

## Running Migrations

### Using psql directly:
```bash
psql -h localhost -U postgres -d taskmanager -f migrations/004.create_users_table.up.sql
psql -h localhost -U postgres -d taskmanager -f migrations/005.add_user_id_to_tasks.up.sql
```

### Using the batch script (Windows):
```bash
migrations\run_migrations.bat
```

## Rolling Back Migrations

### Rollback users table:
```bash
psql -h localhost -U postgres -d taskmanager -f migrations/004.create_users_table.down.sql
```

### Rollback user_id column:
```bash
psql -h localhost -U postgres -d taskmanager -f migrations/005.add_user_id_to_tasks.down.sql
```

## Notes
- Make sure to run migrations in order (004 before 005)
- The `user_id` column is nullable initially to allow existing tasks to be updated
- ON DELETE CASCADE ensures tasks are deleted when users are deleted
