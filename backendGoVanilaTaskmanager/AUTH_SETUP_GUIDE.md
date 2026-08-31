# Authentication & Multi-Tenancy Setup Guide

## Prerequisites
- PostgreSQL database running locally
- Database `taskmanager` created
- User `postgres` with appropriate permissions

## Setup Steps

### 1. Run Database Migrations

Run the new migrations to create the users table and add user_id to tasks:

```bash
# Option 1: Using psql directly
psql -h localhost -U postgres -d taskmanager -f migrations/004.create_users_table.up.sql
psql -h localhost -U postgres -d taskmanager -f migrations/005.add_user_id_to_tasks.up.sql

# Option 2: Using the batch script (Windows)
migrations\run_migrations.bat
```

### 2. Update Configuration

Make sure your `env/config.yaml` includes the JWT secret:

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password_here"
  dbname: "taskmanager"
  sslmode: "disable"

jwt:
  secret: "your_jwt_secret_key_here_change_in_production"
```

Or set it via environment variable:

```bash
$env:JWT_SECRET="your-secret-key"
```

### 3. Start the Server

```bash
go run main.go
```

### 4. Test Authentication Flow

#### Register a new user:
```powershell
$body = @{
    name = "Test User"
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/auth/register" -Method POST -Body $body -ContentType "application/json"
```

#### Login:
```powershell
$body = @{
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri "http://localhost:8080/auth/login" -Method POST -Body $body -ContentType "application/json"
$token = ($response.Content | ConvertFrom-Json).data.token
```

#### Access protected endpoint:
```powershell
$headers = @{
    "Authorization" = "Bearer $token"
}

Invoke-WebRequest -Uri "http://localhost:8080/auth/me" -Method GET -Headers $headers
```

#### Create a task (protected):
```powershell
$taskBody = @{
    title = "Test Task"
    description = "This is a test task"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -Body $taskBody -ContentType "application/json" -Headers $headers
```

## API Endpoints

### Public Endpoints
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user
- `POST /auth/logout` - Logout user
- `GET /health` - Health check

### Protected Endpoints (require JWT token)
- `GET /auth/me` - Get current user info
- `GET /tasks` - Get all tasks for current user
- `POST /tasks` - Create new task
- `GET /tasks/{id}` - Get specific task
- `PUT /tasks/{id}` - Update task
- `DELETE /tasks/{id}` - Delete task
- `PATCH /tasks/{id}/complete` - Mark task as completed
- `PATCH /tasks/{id}/uncomplete` - Mark task as uncompleted
- `PATCH /tasks/{id}/restore` - Restore deleted task

## Security Features Implemented

1. **Password Hashing**: Uses bcrypt with default cost factor
2. **JWT Authentication**: Stateless token-based authentication
3. **Multi-Tenancy**: All task operations are scoped to the authenticated user
4. **IDOR Prevention**: Users can only access their own tasks
5. **Token Validation**: Middleware validates JWT tokens on protected routes

## Testing Notes

- All task operations now require authentication
- Each user can only access their own tasks
- Existing tasks without user_id will need to be updated or assigned
- JWT tokens expire after 24 hours by default
