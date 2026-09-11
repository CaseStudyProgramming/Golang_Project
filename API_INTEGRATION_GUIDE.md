# API Integration Guide

This guide provides comprehensive information for integrating with the Task Manager backend API, including authentication, endpoints, request/response formats, error handling, and code examples.

## Table of Contents

- [API Overview](#api-overview)
- [Authentication](#authentication)
- [Base URL](#base-url)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [API Endpoints](#api-endpoints)
- [Rate Limiting](#rate-limiting)
- [Webhooks](#webhooks)
- [SDK Examples](#sdk-examples)
- [Testing](#testing)

## API Overview

The Task Manager API is a RESTful API built with Go that provides endpoints for task management, user authentication, categories, tags, and analytics. All endpoints return JSON responses and follow REST conventions.

### Key Features

- JWT-based authentication
- Multi-tenancy support (user data isolation)
- Comprehensive CRUD operations
- Advanced filtering and pagination
- Activity logging and audit trail
- CSV export functionality
- Swagger documentation

## Authentication

### JWT Token-based Authentication

The API uses JWT (JSON Web Tokens) for authentication. You must include a valid JWT token in the `Authorization` header for protected endpoints.

#### Getting a Token

1. **Register a new user:**
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

2. **Login to get token:**
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "1",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

#### Using the Token

Include the token in the `Authorization` header:

```http
GET /tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Token Refresh

Tokens have an expiration time. When a token expires, you need to login again to get a new token.

### Authentication Headers

All protected endpoints require:

```http
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

## Base URL

### Development
```
http://localhost:8080
```

### Production
```
https://api.yourdomain.com
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "status": "success",
  "message": "Operation successful",
  "data": {
    // Response data here
  }
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
    "data": [
      // Array of items
    ],
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

## Error Handling

### HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required or invalid
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Common Error Responses

#### Authentication Error
```json
{
  "status": "error",
  "message": "Invalid or expired token",
  "data": null
}
```

#### Validation Error
```json
{
  "status": "error",
  "message": "Validation failed",
  "data": {
    "errors": [
      {
        "field": "title",
        "message": "Title is required"
      }
    ]
  }
}
```

#### Not Found Error
```json
{
  "status": "error",
  "message": "Task not found",
  "data": null
}
```

### Error Handling Best Practices

1. **Always check the status field** in responses
2. **Handle HTTP status codes** appropriately
3. **Display user-friendly error messages** based on the message field
4. **Implement retry logic** for transient failures
5. **Log errors** for debugging purposes
6. **Validate input** before sending requests

## API Endpoints

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "1",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "1",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

#### Get Current User
```http
GET /auth/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "id": "1",
    "name": "John Doe",
    "email": "john@example.com",
    "timezone": "UTC"
  }
}
```

### Task Endpoints

#### Create Task
```http
POST /tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Complete project documentation",
  "sub_title": "High priority task",
  "description": "Write comprehensive documentation for the project",
  "priority": "HIGH",
  "due_date": "2026-12-31T23:59:59Z",
  "category_id": "1"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Task created successfully",
  "data": {
    "id": "1",
    "title": "Complete project documentation",
    "sub_title": "High priority task",
    "description": "Write comprehensive documentation for the project",
    "priority": "HIGH",
    "completed": false,
    "due_date": "2026-12-31T23:59:59Z",
    "category_id": "1",
    "user_id": "1",
    "created_at": "2026-09-11T10:00:00Z",
    "updated_at": "2026-09-11T10:00:00Z"
  }
}
```

#### Get All Tasks
```http
GET /tasks?page=1&limit=10&completed=false&search=project&priority=HIGH&category_id=1&sort_by=due_date&sort_order=asc
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (integer) - Page number (default: 1)
- `limit` (integer) - Items per page (default: 10)
- `completed` (boolean) - Filter by completion status
- `search` (string) - Search in title and description
- `priority` (string) - Filter by priority (LOW, MEDIUM, HIGH, URGENT)
- `category_id` (string) - Filter by category
- `sort_by` (string) - Sort field (due_date, created_at, priority)
- `sort_order` (string) - Sort direction (asc, desc)

**Response:**
```json
{
  "status": "success",
  "message": "Tasks retrieved successfully",
  "data": {
    "data": [
      {
        "id": "1",
        "title": "Complete project documentation",
        "priority": "HIGH",
        "completed": false,
        "due_date": "2026-12-31T23:59:59Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_data": 25,
      "total_page": 3,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

#### Get Task by ID
```http
GET /tasks/{id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Task retrieved successfully",
  "data": {
    "id": "1",
    "title": "Complete project documentation",
    "description": "Write comprehensive documentation",
    "priority": "HIGH",
    "completed": false,
    "due_date": "2026-12-31T23:59:59Z",
    "category": {
      "id": "1",
      "name": "Work",
      "color_hex": "#3B82F6"
    },
    "tags": [
      {
        "id": "1",
        "name": "urgent",
        "color_hex": "#EF4444"
      }
    ],
    "subtasks": [
      {
        "id": "1",
        "title": "Write API documentation",
        "completed": false
      }
    ]
  }
}
```

#### Update Task
```http
PUT /tasks/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated task title",
  "description": "Updated description",
  "priority": "MEDIUM",
  "completed": true
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Task updated successfully",
  "data": {
    "id": "1",
    "title": "Updated task title",
    "description": "Updated description",
    "priority": "MEDIUM",
    "completed": true,
    "updated_at": "2026-09-11T11:00:00Z"
  }
}
```

#### Delete Task (Soft Delete)
```http
DELETE /tasks/{id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Task deleted successfully",
  "data": null
}
```

#### Mark Task as Completed
```http
PATCH /tasks/{id}/complete
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Task marked as completed",
  "data": {
    "id": "1",
    "completed": true,
    "updated_at": "2026-09-11T11:00:00Z"
  }
}
```

#### Mark Task as Uncompleted
```http
PATCH /tasks/{id}/uncomplete
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Task marked as uncompleted",
  "data": {
    "id": "1",
    "completed": false,
    "updated_at": "2026-09-11T11:00:00Z"
  }
}
```

#### Restore Deleted Task
```http
PATCH /tasks/{id}/restore
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Task restored successfully",
  "data": {
    "id": "1",
    "deleted_at": null,
    "updated_at": "2026-09-11T11:00:00Z"
  }
}
```

### Bulk Operations

#### Bulk Delete Tasks
```http
POST /tasks/bulk-delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "task_ids": [1, 2, 3, 4, 5]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Tasks deleted successfully",
  "data": {
    "deleted_count": 5
  }
}
```

#### Bulk Complete Tasks
```http
POST /tasks/bulk-complete
Authorization: Bearer <token>
Content-Type: application/json

{
  "task_ids": [1, 2, 3, 4, 5]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Tasks completed successfully",
  "data": {
    "completed_count": 5
  }
}
```

### Export Endpoints

#### Export Tasks to CSV
```http
GET /tasks/export/csv?completed=true&priority=HIGH&category_id=1
Authorization: Bearer <token>
```

**Query Parameters:**
- `completed` (boolean) - Filter by completion status
- `priority` (string) - Filter by priority
- `category_id` (string) - Filter by category

**Response:**
CSV file with columns: ID, Title, SubTitle, Description, Completed, DueDate, Priority, CategoryID, CreatedAt, UpdatedAt

### Category Endpoints

#### Create Category
```http
POST /categories
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Work",
  "color_hex": "#3B82F6"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Category created successfully",
  "data": {
    "id": "1",
    "name": "Work",
    "color_hex": "#3B82F6",
    "user_id": "1",
    "created_at": "2026-09-11T10:00:00Z"
  }
}
```

#### Get All Categories
```http
GET /categories
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": "1",
      "name": "Work",
      "color_hex": "#3B82F6",
      "task_count": 15
    },
    {
      "id": "2",
      "name": "Personal",
      "color_hex": "#10B981",
      "task_count": 8
    }
  ]
}
```

### Tag Endpoints

#### Create Tag
```http
POST /tags
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "urgent",
  "color_hex": "#EF4444"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Tag created successfully",
  "data": {
    "id": "1",
    "name": "urgent",
    "color_hex": "#EF4444",
    "user_id": "1",
    "created_at": "2026-09-11T10:00:00Z"
  }
}
```

#### Get All Tags
```http
GET /tags
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Tags retrieved successfully",
  "data": [
    {
      "id": "1",
      "name": "urgent",
      "color_hex": "#EF4444"
    },
    {
      "id": "2",
      "name": "important",
      "color_hex": "#F59E0B"
    }
  ]
}
```

#### Add Tag to Task
```http
POST /tasks/{id}/tags
Authorization: Bearer <token>
Content-Type: application/json

{
  "tag_id": "1"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Tag added to task successfully",
  "data": null
}
```

#### Remove Tag from Task
```http
DELETE /tasks/{id}/tags/{tag_id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Tag removed from task successfully",
  "data": null
}
```

### Subtask Endpoints

#### Create Subtask
```http
POST /tasks/{id}/subtasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Subtask 1"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Subtask created successfully",
  "data": {
    "id": "1",
    "title": "Subtask 1",
    "completed": false,
    "task_id": "1",
    "created_at": "2026-09-11T10:00:00Z"
  }
}
```

#### Toggle Subtask
```http
PATCH /subtasks/{id}/toggle
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Subtask toggled successfully",
  "data": {
    "id": "1",
    "completed": true,
    "updated_at": "2026-09-11T11:00:00Z"
  }
}
```

#### Delete Subtask
```http
DELETE /subtasks/{id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Subtask deleted successfully",
  "data": null
}
```

### Analytics Endpoints

#### Get Analytics Summary
```http
GET /tasks/analytics/summary
Authorization: Bearer <token>
```

**Response:**
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

### Activity Log Endpoints

#### Get User Activity Logs
```http
GET /activity-logs?page=1&limit=10
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Activity logs retrieved successfully",
  "data": {
    "data": [
      {
        "id": "1",
        "action": "CREATE",
        "entity_type": "task",
        "entity_id": "1",
        "description": "Created task 'Complete project documentation'",
        "user_id": "1",
        "created_at": "2026-09-11T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_data": 50,
      "total_page": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

#### Get Task Activity Logs
```http
GET /tasks/{id}/activity-logs
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Task activity logs retrieved successfully",
  "data": [
    {
      "id": "1",
      "action": "CREATE",
      "entity_type": "task",
      "entity_id": "1",
      "description": "Created task",
      "user_id": "1",
      "created_at": "2026-09-11T10:00:00Z"
    },
    {
      "id": "2",
      "action": "UPDATE",
      "entity_type": "task",
      "entity_id": "1",
      "description": "Updated task priority to HIGH",
      "user_id": "1",
      "created_at": "2026-09-11T10:30:00Z"
    }
  ]
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse. Current limits:

- **100 requests per minute** per IP address
- **1000 requests per hour** per IP address

### Rate Limit Headers

Responses include rate limit headers:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1631234567
```

### Rate Limit Error Response

```json
{
  "status": "error",
  "message": "Rate limit exceeded. Please try again later.",
  "data": null
}
```

HTTP Status: `429 Too Many Requests`

## Webhooks

Webhooks allow you to receive real-time notifications about task events.

### Setting Up Webhooks

Webhooks can be configured through the API (contact administrator for setup).

### Webhook Events

- `task.created` - When a task is created
- `task.updated` - When a task is updated
- `task.deleted` - When a task is deleted
- `task.completed` - When a task is marked as completed
- `user.registered` - When a new user registers

### Webhook Payload Example

```json
{
  "event": "task.created",
  "timestamp": "2026-09-11T10:00:00Z",
  "data": {
    "task": {
      "id": "1",
      "title": "Complete project documentation",
      "priority": "HIGH",
      "user_id": "1"
    }
  }
}
```

## SDK Examples

### JavaScript/TypeScript Example

```typescript
class TaskManagerAPI {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers
    };

    const response = await fetch(url, { ...options, headers });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Request failed');
    }

    return response.json();
  }

  // Authentication
  async register(data: { name: string; email: string; password: string }) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  async login(email: string, password: string) {
    const response = await this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    this.setToken(response.data.token);
    return response;
  }

  // Tasks
  async getTasks(params?: Record<string, any>) {
    const queryString = new URLSearchParams(params).toString();
    return this.request(`/tasks?${queryString}`);
  }

  async createTask(data: any) {
    return this.request('/tasks', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  async updateTask(id: string, data: any) {
    return this.request(`/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }

  async deleteTask(id: string) {
    return this.request(`/tasks/${id}`, {
      method: 'DELETE'
    });
  }

  // Categories
  async getCategories() {
    return this.request('/categories');
  }

  async createCategory(data: { name: string; color_hex: string }) {
    return this.request('/categories', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  // Tags
  async getTags() {
    return this.request('/tags');
  }

  async createTag(data: { name: string; color_hex: string }) {
    return this.request('/tags', {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }
}

// Usage
const api = new TaskManagerAPI('http://localhost:8080');

await api.login('john@example.com', 'password123');
const tasks = await api.getTasks({ page: 1, limit: 10 });
console.log(tasks);
```

### Python Example

```python
import requests
from typing import Dict, Any, Optional

class TaskManagerAPI:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.token: Optional[str] = None

    def set_token(self, token: str):
        self.token = token

    def _request(self, endpoint: str, method: str = 'GET', data: Optional[Dict] = None) -> Dict:
        url = f"{self.base_url}{endpoint}"
        headers = {
            'Content-Type': 'application/json',
            **({'Authorization': f'Bearer {self.token}'} if self.token else {})
        }
        
        response = requests.request(method, url, json=data, headers=headers)
        response.raise_for_status()
        return response.json()

    def register(self, name: str, email: str, password: str) -> Dict:
        return self._request('/auth/register', 'POST', {
            'name': name,
            'email': email,
            'password': password
        })

    def login(self, email: str, password: str) -> Dict:
        response = self._request('/auth/login', 'POST', {
            'email': email,
            'password': password
        })
        self.set_token(response['data']['token'])
        return response

    def get_tasks(self, params: Optional[Dict] = None) -> Dict:
        query_string = '&'.join(f"{k}={v}" for k, v in (params or {}).items())
        return self._request(f'/tasks?{query_string}')

    def create_task(self, data: Dict) -> Dict:
        return self._request('/tasks', 'POST', data)

    def update_task(self, task_id: str, data: Dict) -> Dict:
        return self._request(f'/tasks/{task_id}', 'PUT', data)

    def delete_task(self, task_id: str) -> Dict:
        return self._request(f'/tasks/{task_id}', 'DELETE')

# Usage
api = TaskManagerAPI('http://localhost:8080')
api.login('john@example.com', 'password123')
tasks = api.get_tasks({'page': 1, 'limit': 10})
print(tasks)
```

### Go Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type TaskManagerAPI struct {
    baseURL string
    token   string
    client  *http.Client
}

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func NewTaskManagerAPI(baseURL string) *TaskManagerAPI {
    return &TaskManagerAPI{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (api *TaskManagerAPI) SetToken(token string) {
    api.token = token
}

func (api *TaskManagerAPI) request(method, endpoint string, body interface{}) (*APIResponse, error) {
    url := api.baseURL + endpoint
    
    var bodyReader io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }
    
    req, err := http.NewRequest(method, url, bodyReader)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    if api.token != "" {
        req.Header.Set("Authorization", "Bearer "+api.token)
    }
    
    resp, err := api.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var apiResp APIResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }
    
    if resp.StatusCode >= 400 {
        return &apiResp, fmt.Errorf(apiResp.Message)
    }
    
    return &apiResp, nil
}

func (api *TaskManagerAPI) Login(email, password string) error {
    resp, err := api.request("POST", "/auth/login", map[string]string{
        "email":    email,
        "password": password,
    })
    if err != nil {
        return err
    }
    
    data := resp.Data.(map[string]interface{})
    api.SetToken(data["token"].(string))
    return nil
}

func (api *TaskManagerAPI) GetTasks(page, limit int) (*APIResponse, error) {
    return api.request("GET", fmt.Sprintf("/tasks?page=%d&limit=%d", page, limit), nil)
}

func main() {
    api := NewTaskManagerAPI("http://localhost:8080")
    
    err := api.Login("john@example.com", "password123")
    if err != nil {
        fmt.Println("Login failed:", err)
        return
    }
    
    tasks, err := api.GetTasks(1, 10)
    if err != nil {
        fmt.Println("Failed to get tasks:", err)
        return
    }
    
    fmt.Println("Tasks:", tasks)
}
```

## Testing

### Testing with cURL

```bash
# Register user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'

# Get tasks (replace TOKEN with actual token)
curl -X GET http://localhost:8080/tasks \
  -H "Authorization: Bearer TOKEN"

# Create task
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"New task","priority":"HIGH"}'
```

### Testing with Postman

1. Import the API collection from Swagger documentation
2. Set base URL to `http://localhost:8080`
3. Use the `/auth/login` endpoint to get a token
4. Add the token to the Authorization header for protected endpoints
5. Test each endpoint with sample data

### Interactive API Documentation

Interactive API documentation is available via Swagger UI:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/openapi.yaml`

This provides a complete interface for testing all API endpoints directly from your browser.

## Additional Resources

- [Backend README](./backendGoVanilaTaskmanager/readme.md) - More backend details
- [AGENTS.md](./backendGoVanilaTaskmanager/AGENTS.md) - Backend development guidelines
- [Swagger Documentation](http://localhost:8080/swagger/index.html) - Interactive API docs
- [Component Documentation](./COMPONENT_DOCUMENTATION.md) - Frontend component guide
