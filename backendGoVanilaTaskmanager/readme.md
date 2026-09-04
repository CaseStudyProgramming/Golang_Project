# Go Vanilla Task Manager API

A RESTful API task management backend built with pure Go (standard library) without external web frameworks.

## Architecture

This project follows a **Layered Architecture** pattern with clear separation of concerns:

```
backendGoVanilaTaskmanager/
├── config/             # Database configuration, environment variables
├── controllers/        # HTTP request/response handlers (Input validation)
├── models/             # Database schema / ORM operations
├── routes/             # API endpoint definitions
├── services/           # Business logic layer
├── middlewares/        # Middleware functions (Auth, Logging, CORS, Recovery)
├── utils/              # Helper functions (Response formatting, JWT)
├── tests/              # Unit and integration tests
├── migrations/         # Database migration scripts
├── swagger/            # OpenAPI/Swagger documentation
├── main.go             # Application entry point
└── readme.md           # This file
```

## Features

### Authentication & Authorization
- User registration and login with JWT tokens
- Password hashing with bcrypt
- Protected routes with authentication middleware
- Multi-tenancy support (user data isolation)

### Task Management
- CRUD operations for tasks
- Task completion status management
- Soft delete with restore functionality
- Advanced filtering and pagination
- Search functionality
- Priority levels (LOW, MEDIUM, HIGH, URGENT)
- Category support
- Tagging system
- Subtasks/checklist items
- Activity logging and audit trail

### Analytics & Reporting
- Task analytics summary endpoint
- CSV export functionality
- Bulk operations (delete and complete)

### API Documentation
- Swagger UI available at `/swagger/index.html`
- OpenAPI specification at `/swagger/openapi.yaml`

## Getting Started

### Prerequisites
- Go 1.22 or higher
- PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/CaseStudyProgramming/Golang_Project.git
cd backendGoVanilaTaskmanager
```

2. Install dependencies:
```bash
go mod download
```

3. Configure database:
```bash
# Copy example config
cp env/config.example.yaml env/config.yaml

# Edit database credentials in env/config.yaml
```

4. Run database migrations:
```bash
# Execute migration scripts in migrations/ folder
psql -U your_user -d taskmanager -f migrations/001.init.up.sql
psql -U your_user -d taskmanager -f migrations/002.add_deleted_at.up.sql
```

5. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Get Current User
```http
GET /auth/me
Authorization: Bearer <jwt_token>
```

### Tasks

#### Create Task
```http
POST /tasks
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Complete project documentation",
  "sub_title": "High priority task",
  "description": "Write comprehensive documentation for the project",
  "priority": "HIGH",
  "due_date": "2026-12-31T23:59:59Z",
  "category_id": 1
}
```

#### Get All Tasks
```http
GET /tasks?page=1&limit=10&completed=false&search=project&priority=HIGH&category_id=1&sort_by=due_date&sort_order=asc
Authorization: Bearer <jwt_token>
```

#### Get Task by ID
```http
GET /tasks/{id}
Authorization: Bearer <jwt_token>
```

#### Update Task
```http
PUT /tasks/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Updated task title",
  "completed": true
}
```

#### Delete Task (Soft Delete)
```http
DELETE /tasks/{id}
Authorization: Bearer <jwt_token>
```

#### Mark Task as Completed
```http
PATCH /tasks/{id}/complete
Authorization: Bearer <jwt_token>
```

#### Mark Task as Uncompleted
```http
PATCH /tasks/{id}/uncomplete
Authorization: Bearer <jwt_token>
```

#### Restore Deleted Task
```http
PATCH /tasks/{id}/restore
Authorization: Bearer <jwt_token>
```

### Analytics

#### Get Analytics Summary
```http
GET /tasks/analytics/summary
Authorization: Bearer <jwt_token>
```

Response:
```json
{
  "status": "success",
  "message": "Analytics summary retrieved successfully",
  "data": {
    "total_active": 10,
    "total_completed": 5,
    "total_overdue": 2,
    "completion_percentage": 50.0,
    "priority_distribution": {
      "LOW": 2,
      "MEDIUM": 4,
      "HIGH": 3,
      "URGENT": 1
    }
  }
}
```

### Bulk Operations

#### Bulk Delete Tasks
```http
POST /tasks/bulk-delete
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "task_ids": [1, 2, 3, 4, 5]
}
```

#### Bulk Complete Tasks
```http
POST /tasks/bulk-complete
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "task_ids": [1, 2, 3, 4, 5]
}
```

### Export

#### Export Tasks to CSV
```http
GET /tasks/export/csv?completed=true&priority=HIGH&category_id=1
Authorization: Bearer <jwt_token>
```

Response will be a CSV file with the following columns:
- ID, Title, SubTitle, Description, Completed, DueDate, Priority, CategoryID, CreatedAt, UpdatedAt

### Categories

#### Create Category
```http
POST /categories
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Work",
  "color_hex": "#3B82F6"
}
```

#### Get All Categories
```http
GET /categories
Authorization: Bearer <jwt_token>
```

### Tags

#### Create Tag
```http
POST /tags
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "urgent",
  "color_hex": "#EF4444"
}
```

#### Add Tag to Task
```http
POST /tasks/{id}/tags
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "tag_id": 1
}
```

### Subtasks

#### Create Subtask
```http
POST /tasks/{id}/subtasks
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Subtask 1"
}
```

#### Toggle Subtask
```http
PATCH /subtasks/{id}/toggle
Authorization: Bearer <jwt_token>
```

### Activity Logs

#### Get User Activity Logs
```http
GET /activity-logs?page=1&limit=10
Authorization: Bearer <jwt_token>
```

#### Get Task Activity Logs
```http
GET /tasks/{id}/activity-logs
Authorization: Bearer <jwt_token>
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run tests with race detection:
```bash
go test -race ./...
```

Or use the provided scripts:
```bash
# On Linux/Mac
./test_with_race.sh

# On Windows
test_with_race.bat
```

Run specific package tests:
```bash
go test ./controllers
go test ./services
go test ./models
```

## API Documentation

Interactive API documentation is available via Swagger UI:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/openapi.yaml`

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "status": "success",
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description",
  "data": null
}
```

### Paginated Response
```json
{
  "status": "success",
  "message": "success",
  "data": {
    "data": [ ... ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_data": 100,
      "total_page": 10,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

## Development

### Build
```bash
go build -o server.exe
```

### Run with live reload (using Air)
```bash
air
```

### Code Formatting
Check if code is properly formatted:
```bash
# On Linux/Mac
./check_format.sh

# On Windows
check_format.bat

# Or manually
gofmt -l .
```

Format all Go files:
```bash
gofmt -w .
```

### Environment Variables

Key environment variables (can be set in `env/config.yaml`):

- `server.port`: Server port (default: 8080)
- `database.host`: Database host
- `database.port`: Database port
- `database.user`: Database user
- `database.password`: Database password
- `database.dbname`: Database name
- `database.sslmode`: SSL mode
- `jwt.secret`: JWT secret key

## Security Features

- Password hashing with bcrypt
- JWT token authentication
- SQL injection prevention (parameterized queries)
- CORS configuration
- Request rate limiting (planned)
- Input validation

## Technology Stack

- **Language**: Go 1.22+
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **HTTP Server**: Go standard library `net/http`
- **Routing**: Go 1.22+ standard `http.ServeMux`
- **Testing**: Go testing framework

## Frontend Integration

This API is designed to work with a Svelte frontend. The API provides:

- RESTful endpoints following REST conventions
- Consistent JSON response format
- Proper HTTP status codes
- CORS support for frontend integration
- Swagger documentation for frontend developers
- CSV export for data analysis

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the MIT License.

## Future Enhancements

- WebSocket support for real-time updates
- File upload support for task attachments
- Email notifications for task reminders
- Advanced search with filters
- Task templates
- Team collaboration features
- Performance optimization and caching
- Rate limiting and DDoS protection