# Task Manager Setup Guide

## Prerequisites

- Podman (recommended) or Docker
- Go 1.20+ (for backend)
- Node.js 18+ (for frontend)

## Quick Start

### 1. Start PostgreSQL with Podman

```bash
cd /path/to/taskmanager
podman-compose up -d
```

### 2. Configure Backend

Create/update `backendGoVanilaTaskmanager/env/config.yaml`:

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

Or use environment variables (recommended for team collaboration):

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=taskmanager
export DB_SSLMODE=disable
```

### 3. Run Backend

```bash
cd backendGoVanilaTaskmanager
go run main.go
```

Backend will run at `http://localhost:8080`

### 4. Run Frontend

```bash
cd sveltekit-taskmanager-frontend
npm install
npm run dev
```

Frontend will run at `http://localhost:5173`

## Podman Commands

```bash
# Start PostgreSQL
podman-compose up -d

# Stop PostgreSQL
podman-compose down

# View status
podman ps

# View logs
podman logs taskmanager-postgres

# Connect to database
podman exec -it taskmanager-postgres psql -U postgres

# Database operations
podman exec -it taskmanager-postgres psql -U postgres -c "\l"  # List databases
podman exec -it taskmanager-postgres psql -U postgres -d taskmanager  # Connect to taskmanager DB
```

## Docker Commands (for team members using Docker)

Same commands work with Docker:

```bash
docker-compose up -d
docker-compose down
docker ps
docker logs taskmanager-postgres
```

## Environment Variables

The backend supports environment variables that override config file values:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: taskmanager)
- `DB_SSLMODE` - SSL mode (default: disable)
- `JWT_SECRET` - JWT secret key
- `APP_ENV` - Application environment (development/production)
- `CORS_ALLOWED_ORIGINS` - CORS allowed origins

## Database Setup

The PostgreSQL container will automatically:
- Create database named `taskmanager`
- Set default user `postgres` with password `postgres`
- Be accessible on port 5432

## Project Structure

```
taskmanager/
├── docker-compose.yml              # PostgreSQL configuration
├── PODMAN_SETUP.md                 # Detailed Podman setup
├── SETUP_GUIDE.md                  # This file
├── backendGoVanilaTaskmanager/     # Go backend
│   ├── config/                     # Configuration
│   ├── env/                        # Environment files
│   └── main.go                     # Main application
└── sveltekit-taskmanager-frontend/ # SvelteKit frontend
```

## Troubleshooting

### PostgreSQL connection failed
1. Check if Podman container is running: `podman ps`
2. Check logs: `podman logs taskmanager-postgres`
3. Verify port 5432 is not used by native PostgreSQL
4. Stop native PostgreSQL if needed: `sudo systemctl stop postgresql`

### Port already in use
```bash
# Stop native PostgreSQL
sudo systemctl stop postgresql@15
sudo systemctl stop postgresql@18

# Or kill process
sudo pkill -9 postgres
```

### Backend can't connect to database
1. Verify PostgreSQL is running: `podman ps`
2. Test connection: `PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d taskmanager -c "\l"`
3. Check config.yaml or environment variables are correct

## Team Collaboration

This setup ensures:
- ✅ Same PostgreSQL version (15) across all team members
- ✅ Consistent environment configuration
- ✅ Cross-platform compatibility (Mac, Windows, Linux)
- ✅ Easy database backup and restore
- ✅ Production-ready secrets management

## Production Deployment

For production:
1. Change database password in docker-compose.yml
2. Use environment variables for secrets
3. Set proper CORS origins
4. Use strong JWT secret
5. Enable SSL for database connections
