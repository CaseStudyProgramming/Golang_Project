@echo off
echo Running database migrations...

psql -h localhost -U postgres -d taskmanager -f migrations/004.create_users_table.up.sql
if errorlevel 1 (
    echo Failed to create users table
    exit /b 1
)

psql -h localhost -U postgres -d taskmanager -f migrations/005.add_user_id_to_tasks.up.sql
if errorlevel 1 (
    echo Failed to add user_id to tasks table
    exit /b 1
)

echo Migrations completed successfully!
