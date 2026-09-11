# Task Manager Application

A full-stack task management application with a Go backend and SvelteKit frontend.

## Overview

This project consists of:
- **Backend**: Go Vanilla Task Manager API (RESTful API with PostgreSQL)
- **Frontend**: SvelteKit + TypeScript + Tailwind CSS

## Features

### Backend Features
- User authentication with JWT
- Task CRUD operations with advanced filtering
- Priority levels (LOW, MEDIUM, HIGH, URGENT)
- Categories and tags
- Subtasks/checklist items
- Activity logging and audit trail
- Analytics and reporting
- CSV export functionality
- Soft delete with restore
- Bulk operations
- Swagger API documentation

### Frontend Features
- Modern Svelte 5 with runes
- TypeScript type safety
- Tailwind CSS for styling
- Responsive design
- Real-time task management
- Analytics dashboard
- Data visualization with Chart.js

## Prerequisites

- **Go**: 1.22 or higher
- **Node.js**: 18+ (or Bun as package manager)
- **PostgreSQL**: 15+ (via Podman/Docker)
- **Git**: Latest version

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/CaseStudyProgramming/Golang_Project.git
cd taskmanager
```

### 2. Start PostgreSQL with Podman

```bash
podman-compose up -d
```

Or with Docker:
```bash
docker-compose up -d
```

### 3. Configure Backend

Create or update `backendGoVanilaTaskmanager/env/config.yaml`:

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "taskmanager"
  sslmode: "disable"

jwt:
  secret: "your_jwt_secret_key_here_change_in_production"

cors:
  allowed_origins:
    - "http://localhost:5173"
    - "http://localhost:3000"
    - "http://localhost:8080"
```

### 4. Run Database Migrations

```bash
cd backendGoVanilaTaskmanager
psql -U postgres -d taskmanager -f migrations/001.init.up.sql
psql -U postgres -d taskmanager -f migrations/002.add_deleted_at.up.sql
```

### 5. Start Backend

```bash
cd backendGoVanilaTaskmanager
go run main.go
```

Backend will run at `http://localhost:8080`

### 6. Configure Frontend

Create or update `sveltekit-taskmanager-frontend/.env`:

```env
PUBLIC_API_URL=http://localhost:8080
```

### 7. Install Frontend Dependencies

```bash
cd sveltekit-taskmanager-frontend
bun install
```

Or with npm:
```bash
npm install
```

### 8. Start Frontend

```bash
bun run dev
```

Or with npm:
```bash
npm run dev
```

Frontend will run at `http://localhost:5173`

## Project Structure

```
taskmanager/
├── backendGoVanilaTaskmanager/     # Go backend API
│   ├── config/                     # Configuration files
│   ├── controllers/                # HTTP handlers
│   ├── models/                     # Database models
│   ├── services/                   # Business logic
│   ├── middlewares/                # Middleware functions
│   ├── routes/                     # API routes
│   ├── utils/                      # Helper functions
│   ├── tests/                      # Test files
│   ├── migrations/                 # Database migrations
│   ├── swagger/                    # API documentation
│   └── main.go                     # Application entry point
├── sveltekit-taskmanager-frontend/ # SvelteKit frontend
│   ├── src/                        # Source code
│   │   ├── lib/                    # Shared utilities
│   │   │   ├── features/           # Feature modules
│   │   │   ├── shared/             # Shared components
│   │   │   └── server/             # Server-side code
│   │   └── routes/                 # SvelteKit routes
│   ├── static/                     # Static assets
│   ├── tests/                      # Test files
│   └── playwright-tests/           # E2E tests
├── docker-compose.yml              # PostgreSQL configuration
└── README.md                       # This file
```

## Documentation

- [Backend README](./backendGoVanilaTaskmanager/readme.md) - Backend API documentation
- [Frontend README](./sveltekit-taskmanager-frontend/README.md) - Frontend documentation
- [Setup Guide](./SETUP_GUIDE.md) - Detailed setup instructions
- [Podman Setup](./PODMAN_SETUP.md) - Podman-specific setup
- [AGENTS.md](./backendGoVanilaTaskmanager/AGENTS.md) - Backend development guidelines
- [AGENTS_FrontEnd.md](./sveltekit-taskmanager-frontend/AGENTS_FrontEnd.md) - Frontend development guidelines

## API Documentation

Interactive API documentation is available via Swagger UI:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/openapi.yaml`

## Development

### Backend Development

```bash
cd backendGoVanilaTaskmanager

# Run with live reload
air

# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Format code
gofmt -w .

# Build for production
go build -o server.exe
```

### Frontend Development

```bash
cd sveltekit-taskmanager-frontend

# Start development server
bun run dev

# Type checking
bun run check

# Run tests
bun test

# Run tests with UI
bun test:ui

# Run tests with coverage
bun test:coverage

# Run E2E tests
bun test:e2e

# Build for production
bun run build

# Preview production build
bun run preview

# Lint code
bun run lint

# Format code
bun run format
```

## Environment Variables

### Backend Environment Variables

These can be set in `backendGoVanilaTaskmanager/env/config.yaml`:

- `server.port`: Server port (default: 8080)
- `database.host`: Database host
- `database.port`: Database port
- `database.user`: Database user
- `database.password`: Database password
- `database.dbname`: Database name
- `database.sslmode`: SSL mode
- `jwt.secret`: JWT secret key
- `cors.allowed_origins`: CORS allowed origins

### Frontend Environment Variables

These should be set in `sveltekit-taskmanager-frontend/.env`:

- `PUBLIC_API_URL`: Backend API URL (required, must start with PUBLIC_)

## Testing

### Backend Testing

```bash
cd backendGoVanilaTaskmanager

# Run all tests
go test ./...

# Run specific package tests
go test ./controllers
go test ./services
go test ./models

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...
```

### Frontend Testing

```bash
cd sveltekit-taskmanager-frontend

# Run unit tests
bun test

# Run tests with UI
bun test:ui

# Run tests with coverage
bun test:coverage

# Run E2E tests
bun test:e2e

# Run E2E tests with UI
bun test:e2e:ui

# Run E2E tests in debug mode
bun test:e2e:debug
```

## Production Deployment

### Backend Deployment

1. **Build the application**:
```bash
cd backendGoVanilaTaskmanager
go build -o server
```

2. **Set production environment variables**:
```bash
export DB_HOST=your_production_db_host
export DB_PORT=5432
export DB_USER=your_db_user
export DB_PASSWORD=your_secure_password
export DB_NAME=taskmanager
export DB_SSLMODE=require
export JWT_SECRET=your_secure_jwt_secret
export APP_ENV=production
```

3. **Run the application**:
```bash
./server
```

### Frontend Deployment

1. **Build the application**:
```bash
cd sveltekit-taskmanager-frontend
bun run build
```

2. **Set production environment variables**:
```bash
export PUBLIC_API_URL=https://your-api-domain.com
```

3. **Deploy the build**:
The `build/` directory contains the production-ready files. Deploy to your hosting service (Vercel, Netlify, etc.)

## Troubleshooting

### PostgreSQL Connection Issues

1. Check if Podman container is running:
```bash
podman ps
```

2. Check logs:
```bash
podman logs taskmanager-postgres
```

3. Test connection:
```bash
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d taskmanager -c "\l"
```

### Port Already in Use

```bash
# Find process using the port
lsof -i :8080
lsof -i :5173

# Kill the process
kill -9 <PID>
```

### Backend Can't Connect to Database

1. Verify PostgreSQL is running
2. Check config.yaml credentials
3. Ensure database exists: `psql -U postgres -c "\l"`
4. Check firewall settings

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on GitHub.
