# Testing Guide - Pyramid Testing Strategy

This project follows a pyramid testing strategy with three layers of testing:
- **Unit Tests**: Fast, isolated tests with mocked dependencies
- **Integration Tests**: Medium-speed tests with real database
- **E2E Tests**: Slow tests with real backend and frontend

## Test Pyramid Structure

```
        /\
       /E2E\        ← Playwright, real backend + frontend (slowest)
      /------\
     /Integration\  ← Go tests, real DB (medium)
    /----------\
   /  Unit Tests \  ← Vitest, MSW mocks from OpenAPI (fastest)
  /______________\
```

## CI/CD Overview

The project uses GitHub Actions for continuous integration and deployment with a **Docker-based approach** for cross-platform consistency:

- **Backend CI** (`.github/workflows/backend-ci.yml`): Format → Lint → Unit tests → Integration tests → Security scan → Build
- **Frontend CI** (`.github/workflows/frontend-ci.yml`): API generation → Format → Lint → Type check → Unit tests → Build → Docker build → E2E tests with Docker containers
- **Main CI** (`.github/workflows/ci.yml`): Orchestrates backend and frontend CI, builds Docker images, deploys

### Key CI Features

- **Docker-based E2E**: E2E tests use Docker containers for backend (production parity)
- **MSW from OpenAPI**: Frontend unit tests use MSW handlers based on OpenAPI spec
- **Cross-platform Support**: Docker/Podman compatible for macOS, Linux, Windows teams
- **Environment-only Config**: Backend supports environment variables without config files
- **Parallel Execution**: Tests run in parallel when possible
- **Artifact Upload**: Test results and coverage reports are uploaded
- **Docker Caching**: Build optimization with Docker layer caching

## Backend Testing (Go)

### Type Checking in Go
**Important**: Go does not require explicit type checking in CI because:

- **Go is statically-typed**: Type checking happens at compile time
- **`go build` performs type checking**: Build fails if there are type errors
- **`go test` also type-checks**: Tests won't run if code has type errors
- **No separate type check needed**: Type safety is built into the language

**Type Safety in Go vs TypeScript:**
- **Go**: Compile-time type checking (happens automatically during build)
- **TypeScript**: Runtime language, needs explicit type checking tools (`tsc`, `svelte-check`)

This is why `backend-ci.yml` does not have a separate type check step - the `go build` command already handles it.

### E2E Tests in Backend
**Backend does not have E2E tests** because:

- **Backend is an API server**: No UI to test with browser automation
- **Integration tests are sufficient**: Real database testing covers API functionality
- **Performance overhead**: E2E tests are slower and unnecessary for pure API
- **Frontend handles E2E**: Frontend CI tests the complete system including backend

Backend uses **integration tests** instead:
- Real PostgreSQL database
- Complete request/response cycle testing
- API contract validation
- Business logic verification

### Unit Tests
Fast tests for individual functions and business logic with mocked dependencies.

```bash
# Run all unit tests
cd backendGoVanilaTaskmanager
go test ./... -v

# Run specific package tests
go test ./services/... -v
go test ./controllers/... -v

# Run with race detection
go test ./... -race -v

# Run with coverage
go test ./... -cover -v
```

### Integration Tests
Tests that use a real PostgreSQL database to test the complete request/response cycle.

```bash
# Start test database (optional - uses existing PostgreSQL by default)
docker-compose -f docker-compose.test.yml up postgres-test

# Run integration tests
cd backendGoVanilaTaskmanager
go test ./tests/... -v

# Run specific integration test
go test ./tests/integration_test.go -v -run TestTaskAPIIntegration_CreateAndGet

# With environment variables for custom DB
TEST_DB_HOST=localhost TEST_DB_PORT=5433 go test ./tests/... -v
```

### Performance Tests
Tests that measure API response times and database performance.

```bash
# Run performance benchmarks
cd backendGoVanilaTaskmanager
go test ./tests/performance_test.go -v

# Run API benchmarks
go test ./tests/api_benchmark_test.go -v
```

### Test Database Setup
Integration tests use a separate test database (`taskmanager_test`) to avoid affecting development data.

**Environment Variables:**
- `TEST_DB_HOST`: Database host (default: localhost)
- `TEST_DB_PORT`: Database port (default: 5432)
- `TEST_DB_USER`: Database user (default: postgres)
- `TEST_DB_PASSWORD`: Database password (default: Berjuang#382)
- `TEST_DB_NAME`: Database name (default: taskmanager_test)
- `TEST_DB_SSLMODE`: SSL mode (default: disable)

## Frontend Testing (SvelteKit)

### Unit Tests
Fast component and logic tests using Vitest with MSW for API mocking.

```bash
cd sveltekit-taskmanager-frontend

# Generate MSW handlers from OpenAPI spec (run this first after OpenAPI changes)
bun run api:generate

# Run all unit tests
bun run validate

# Run in watch mode
bun run test:ui

# Run with coverage
bun run test:coverage

# Run specific test file
bun run validate src/lib/features/tasks/stores/task.store.logic.test.ts
```

**Unit Test Configuration:**
- Uses MSW (Mock Service Worker) for API mocking
- MSW handlers are generated from OpenAPI spec using Orval
- Tests run in jsdom environment
- Excludes e2e tests and playwright tests
- Coverage threshold: 80% for lines, functions, branches, statements

**MSW from OpenAPI with Orval:**
- MSW handlers are automatically generated from `backendGoVanilaTaskmanager/swagger/openapi.yaml`
- Run `bun run api:generate` after OpenAPI spec changes
- Generated files:
  - `src/lib/api/index.ts` - API client with TypeScript types
  - `src/lib/api/index.msw.ts` - MSW handlers with Faker.js integration
  - `src/lib/api/index.faker.ts` - Mock data factories
  - `src/lib/api/handlers.ts` - Wrapper to export all handlers for MSW setup
- Orval also generates realistic mock data using Faker.js
- This ensures frontend mocks match backend API contract with minimal manual effort

### E2E Tests
End-to-end tests using Playwright with real backend via Docker containers.

```bash
cd sveltekit-taskmanager-frontend

# Start test environment (backend + frontend + test DB via Docker)
cd ..
docker-compose -f docker-compose.test.yml up -d

# Run e2e tests
cd sveltekit-taskmanager-frontend
bun run test:e2e

# Run e2e tests with UI
bun run test:e2e:ui

# Run e2e tests in debug mode
bun run test:e2e:debug

# Run specific e2e test
bun run test:e2e task-management.e2e.ts
```

**Docker-based E2E Testing (Alternative for macOS/Podman):**
```bash
# Build backend Docker image
cd backendGoVanilaTaskmanager
docker build -t taskmanager-backend:test .

# Run backend container with test database
docker run -d --name taskmanager-backend-e2e --network host \
  -e DB_HOST=localhost -e DB_PORT=5432 \
  -e DB_USER=postgres -e DB_PASSWORD=Berjuang#382 \
  -e DB_NAME=taskmanager_test \
  -e JWT_SECRET=test-secret \
  taskmanager-backend:test

# Run e2e tests
cd ../sveltekit-taskmanager-frontend
PUBLIC_API_BASE_URL=http://localhost:8080 bun test:e2e

# Cleanup
docker stop taskmanager-backend-e2e
docker rm taskmanager-backend-e2e
```

**E2E Test Configuration:**
- Uses real backend server via Docker container (production parity)
- Uses real PostgreSQL test database via service container
- Backend runs on http://localhost:8080 (Docker --network host)
- Tests run on Chromium, Firefox, WebKit, and mobile browsers
- Docker-based for cross-platform consistency (macOS, Linux, Windows)

## Running All Tests

### Complete Test Suite
Run all tests in the correct order (fast to slow):

```bash
# 1. Backend unit tests
cd backendGoVanilaTaskmanager
go test ./... -v

# 2. Backend integration tests (requires test DB)
go test ./tests/... -v

# 3. Frontend unit tests
cd ../sveltekit-taskmanager-frontend
bun run validate

# 4. E2E tests (requires test environment)
cd ..
docker-compose -f docker-compose.test.yml up -d
cd sveltekit-taskmanager-frontend
bun run test:e2e
```

### Quick Development Cycle
For fast development feedback:

```bash
# Run only unit tests (fastest)
cd backendGoVanilaTaskmanager && go test ./... -v
cd ../sveltekit-taskmanager-frontend && bun run validate
```

### Pre-Commit Validation
Frontend has a pre-commit validation script:

```bash
cd sveltekit-taskmanager-frontend
bun run validate
```

This runs: format check → lint → type check → unit tests → build

## Test Environment Management

### Local Development Environment
```bash
# Start test database and backend (for local e2e testing)
docker-compose -f docker-compose.test.yml up -d

# Check services are running
docker-compose -f docker-compose.test.yml ps
```

### CI Environment
```bash
# Start CI-like environment (for testing CI workflows locally)
docker-compose -f docker-compose.ci.yml up -d

# Check services are running
docker-compose -f docker-compose.ci.yml ps
```

### Stop Test Environment
```bash
# Stop local test environment
docker-compose -f docker-compose.test.yml down

# Stop CI environment
docker-compose -f docker-compose.ci.yml down

# Stop and remove with volumes
docker-compose -f docker-compose.test.yml down -v
docker-compose -f docker-compose.ci.yml down -v
```

### View Logs
```bash
# View local test environment logs
docker-compose -f docker-compose.test.yml logs
docker-compose -f docker-compose.test.yml logs backend-test
docker-compose -f docker-compose.test.yml logs postgres-test

# View CI environment logs
docker-compose -f docker-compose.ci.yml logs
docker-compose -f docker-compose.ci.yml logs backend-ci
docker-compose -f docker-compose.ci.yml logs postgres-ci
```

## CI/CD Integration

### GitHub Actions
The project includes CI/CD workflows that automatically run tests:

- **Backend CI** (`.github/workflows/backend-ci.yml`): Format → Lint → Unit tests → Integration tests → Security scan → Build → Docker build
  - **No E2E tests**: Backend uses integration tests (API server doesn't need UI testing)
  - **No explicit type check**: Go is statically-typed, type checking happens during `go build`

- **Frontend CI** (`.github/workflows/frontend-ci.yml`): Quality checks → Unit tests → Build → **E2E tests with Docker**
  - **E2E tests are here**: Playwright tests with real backend Docker container
  - **Type check included**: Explicit TypeScript/Svelte type checking with `svelte-check`
  - **Docker-based backend**: Builds and runs backend in container for E2E tests

- **Main CI** (`.github/workflows/ci.yml`): Orchestrates backend and frontend CI, builds Docker images, deploys
  - **No direct tests**: Delegates testing to individual workflows
  - **No E2E here**: E2E runs in frontend-ci.yml which is called by this workflow

### E2E Test Distribution Across Workflows

**Where are E2E tests run?**

| Workflow | Has E2E? | Why? |
|----------|---------|------|
| `backend-ci.yml` | ❌ No | Backend is API server, integration tests are sufficient |
| `frontend-ci.yml` | ✅ Yes | Frontend needs browser automation and user journey testing |
| `ci.yml` | ❌ No | Orchestrator workflow, delegates to individual workflows |

**E2E Test Location Strategy:**

1. **Backend CI (`backend-ci.yml`)**:
   - Focuses on API functionality
   - Uses integration tests with real database
   - No UI to test, so no E2E needed
   - `go test ./tests/...` covers API testing

2. **Frontend CI (`frontend-ci.yml`)**:
   - Has dedicated `e2e` job
   - Uses Playwright for browser automation
   - Tests complete user journeys (login, create task, etc.)
   - Requires real backend server (runs via Docker container)
   - Job structure: `quality-checks` → `test` → `build` → `e2e`

3. **Main CI (`ci.yml`)**:
   - Orchestrates both backend and frontend CI
   - Calls `frontend-ci.yml` which includes E2E
   - Focuses on build, push, and deployment
   - No direct test execution

**Why This Distribution?**

- **Separation of Concerns**: Each workflow has a focused responsibility
- **Efficiency**: Backend doesn't need browser automation overhead
- **Relevance**: Frontend is the only component that needs UI testing
- **Maintainability**: E2E tests live where they're most relevant

### CI Environment Setup

**Backend CI:**
- Uses PostgreSQL service container for integration tests
- Environment variables configured for test database
- Runs with Go 1.25
- Format → Lint → Unit tests → Integration tests → Security scan → Build → Docker build
- **No E2E tests**: Integration tests are sufficient for API server
- **No explicit type check**: Go's static typing handles this during build

**Frontend CI:**
- Has **4 separate jobs**: quality-checks, test, build, e2e
- Generates MSW handlers from OpenAPI spec automatically
- Quality checks: Format → Lint → Type check → Unit tests
- Build job: Builds frontend Docker image
- **E2E job**: Uses PostgreSQL service container + Docker backend container
- **Builds Docker image** for backend (production parity)
- **Starts backend via Docker container** for e2e tests
- Runs with Bun and Node.js 24
- Database initialization via SQL scripts
- Docker-based for cross-platform consistency

**E2E Tests in CI:**
- **Real backend server via Docker container** (not compiled binary)
- Real PostgreSQL test database via service container
- Docker --network host for localhost access
- Playwright browsers installed with system dependencies
- Tests run on Chromium, Firefox, WebKit, and mobile browsers
- Environment-only configuration (no config files needed)

### Running Tests in CI Context

To test CI workflows locally:

```bash
# Test backend CI locally
cd backendGoVanilaTaskmanager
docker-compose -f ../docker-compose.ci.yml up postgres-ci
TEST_DB_HOST=localhost TEST_DB_PORT=5432 go test ./tests/... -v

# Test frontend CI locally
cd sveltekit-taskmanager-frontend
bun run api:generate
bun run validate

# Test e2e locally with CI-like Docker setup
cd ..
docker-compose -f docker-compose.ci.yml up -d
cd sveltekit-taskmanager-frontend
CI=true PUBLIC_API_BASE_URL=http://localhost:8080 bun test:e2e

# Alternative: Manual Docker testing (for macOS/Podman)
cd backendGoVanilaTaskmanager
docker build -t taskmanager-backend:test .
docker run -d --name test-backend --network host \
  -e DB_HOST=localhost -e DB_PORT=5432 \
  -e DB_USER=postgres -e DB_PASSWORD=Berjuang#382 \
  -e DB_NAME=taskmanager_test \
  -e JWT_SECRET=test-secret \
  taskmanager-backend:test
cd ../sveltekit-taskmanager-frontend
PUBLIC_API_BASE_URL=http://localhost:8080 bun test:e2e
docker stop test-backend && docker rm test-backend
```

### Local CI Simulation
```bash
# Backend CI simulation
cd backendGoVanilaTaskmanager
gofmt -l .
golangci-lint run
go vet ./...
go test ./... -race
go build

# Frontend CI simulation
cd sveltekit-taskmanager-frontend
bun run format:check
bun run lint
bun run check
bun run validate
bun run build

# Docker build simulation (matches CI)
cd backendGoVanilaTaskmanager
docker build -t taskmanager-backend:test .
```

## Test Data Management

### E2E Test User
E2E tests use a dedicated test user:
- Email: `e2e-test@example.com`
- Password: `testpassword123`
- Name: `E2E Test User`

The test user is automatically registered before tests run and cleaned up after.

### Test Database Cleanup
Integration tests automatically clean up test data after each test run:
- Deletes test tasks (titles starting with "Test")
- Deletes test users (names starting with "testuser")
- Resets database state between tests

## Troubleshooting

### Backend Tests Fail
- Ensure PostgreSQL is running: `docker ps`
- Check test database exists: `docker exec -it taskmanager-postgres-test psql -U postgres -l`
- Verify connection settings in environment variables

### Frontend Unit Tests Fail
- Ensure dependencies are installed: `bun install`
- Check MSW setup in `tests/setup.ts`
- Verify vitest configuration in `vitest.config.ts`

### E2E Tests Fail
- Ensure test environment is running: `docker-compose -f docker-compose.test.yml ps`
- Check backend health: `curl http://localhost:8080/health`
- Verify Docker containers are running: `docker ps`
- Check backend container logs: `docker logs taskmanager-backend-e2e`
- Ensure database schema is initialized
- Verify environment variables are set correctly

### Docker Build Failures
- Ensure Docker/Podman is installed and running
- Check Go version in Dockerfile matches go.mod (currently 1.25)
- Verify Dockerfile syntax: `docker build -t test .`
- Check for sufficient disk space
- Ensure base images can be pulled: `docker pull golang:1.25-alpine`

### Port Conflicts
If ports are already in use:
- Change ports in `docker-compose.test.yml`
- Update environment variables accordingly
- Kill processes using the ports: `lsof -ti:8080 | xargs kill`

## Best Practices

1. **Write unit tests first** - They're fast and catch most bugs
2. **Keep unit tests isolated** - Mock external dependencies with MSW
3. **Use integration tests sparingly** - Only for critical paths
4. **Minimize e2e tests** - Focus on user journeys, not implementation details
5. **Run tests locally before committing** - Catch issues early
6. **Keep tests fast** - Slow tests discourage running them
7. **Test behavior, not implementation** - Focus on what, not how
8. **Use descriptive test names** - Explain what is being tested
9. **Clean up test data** - Don't leave pollution behind
10. **Monitor test coverage** - Aim for >80% coverage
11. **Regenerate MSW handlers after OpenAPI changes** - Run `bun run api:generate`
12. **Test CI workflows locally** - Use docker-compose.ci.yml to simulate CI environment
13. **Keep OpenAPI spec updated** - Single source of truth for API contract
14. **Use environment-specific configurations** - Local vs CI environments
15. **Leverage Docker for consistency** - Same containers for local, CI, and production
16. **Support environment-only config** - Backend should work without config files using env vars
17. **Test on multiple platforms** - Docker/Podman compatibility for macOS, Linux, Windows teams
18. **Clean up Docker resources** - Remove containers and volumes after testing

## Docker Environment Support

The project uses **Docker-based environments** for cross-platform consistency, especially important for teams with macOS members using Docker Desktop or Podman.

### Docker Configuration Files
- `docker-compose.yml` - Main development environment
- `docker-compose.test.yml` - Testing environment with test database
- `docker-compose.ci.yml` - CI simulation environment
- `Dockerfile` - Backend container build (Go 1.25-alpine)
- `env/config.docker.yaml` - Default Docker configuration

### Docker Configuration Benefits
- **Cross-platform**: Works on macOS, Linux, Windows with Docker/Podman
- **Production parity**: Same containers in development, CI, and production
- **Isolated dependencies**: No host system dependencies
- **Team consistency**: Same environment for all team members
- **Environment-only config**: Backend supports running without config files

### Backend Configuration System
The backend supports **environment-only configuration** for Docker environments:

**Priority Order:**
1. Environment variables (highest priority)
2. Config file values (if exists)
3. Default values (lowest priority)

**Environment Variables:**
- `SERVER_PORT` - Server port (default: 8080)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: taskmanager)
- `DB_SSLMODE` - SSL mode (default: disable)
- `JWT_SECRET` - JWT secret key (required)
- `APP_ENV` - Application environment (default: development)
- `CORS_ALLOWED_ORIGINS` - Comma-separated CORS origins

**Config Files:**
- `env/config.yaml` - Local development config (gitignored)
- `env/config.example.yaml` - Example config for reference
- `env/config.docker.yaml` - Default Docker config (copied in Dockerfile)

### Running Tests with Docker
```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run tests
cd backendGoVanilaTaskmanager
go test ./... -v

# Stop environment
docker-compose -f docker-compose.test.yml down -v
```

### Manual Docker Testing (Alternative Approach)
```bash
# Build backend image
cd backendGoVanilaTaskmanager
docker build -t taskmanager-backend:test .

# Run with environment variables (no config file needed)
docker run -d --name test-backend --network host \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=Berjuang#382 \
  -e DB_NAME=taskmanager_test \
  -e JWT_SECRET=test-secret \
  taskmanager-backend:test

# Health check
curl http://localhost:8080/health

# Cleanup
docker stop test-backend
docker rm test-backend
```

## OpenAPI & MSW Integration

### Why MSW from OpenAPI with Orval?

- **Contract Testing**: Ensures frontend mocks match backend API
- **Single Source of Truth**: OpenAPI spec drives both backend and frontend
- **Automatic Generation**: Orval generates MSW handlers from OpenAPI
- **Type Safety**: Generated TypeScript types from OpenAPI schema
- **Mock Data Factories**: Faker.js integration for realistic test data
- **Consistency**: Frontend and backend always in sync
- **Minimal Maintenance**: One command to regenerate all handlers

### What Orval Generates

Orval generates the following files from your OpenAPI spec:

1. **API Client** (`src/lib/api/index.ts`):
   - TypeScript interfaces for all API types
   - Axios-based API client functions
   - Type-safe request/response definitions
   - Example: `getTasks()`, `createTask()`, `updateTask()`

2. **MSW Handlers** (`src/lib/api/index.msw.ts`):
   - MSW request handlers for all API endpoints
   - Wildcard URL matching (`*/health`, `*/auth/login`, etc.)
   - Override support for custom responses
   - Example: `getGetTasksMockHandler()`, `getPostAuthLoginMockHandler()`

3. **Mock Data Factories** (`src/lib/api/index.faker.ts`):
   - Faker.js-based mock data generators
   - Realistic random data for testing
   - Configurable override support
   - Example: `getGetTasksResponseMock()`, `getPostAuthLoginResponseMock()`

4. **Handlers Wrapper** (`src/lib/api/handlers.ts`):
   - Exports all MSW handlers as an array
   - Used in test setup to configure MSW
   - Maintains list of all API endpoints

5. **Model Directory** (`src/lib/api/model/`):
   - Auto-generated TypeScript type definitions
   - Used by other generated files
   - Should not be edited manually

### Excluding Orval-Generated Files from Checks

**IMPORTANT**: All Orval-generated files must be excluded from type checking, linting, formatting, and coverage. These files are auto-generated and should not be manually edited or validated.

**Current Exclusions Status:**

1. **Vitest Coverage** (`vitest.config.ts`): ✅ **Already configured**
   ```typescript
   exclude: [
     'src/lib/api/index.ts',        // Auto-generated by Orval
     'src/lib/api/index.msw.ts',    // Auto-generated by Orval
     'src/lib/api/index.faker.ts',  // Auto-generated by Orval
     'src/lib/api/model/',          // Auto-generated types
   ]
   ```

2. **TypeScript Type Check** (`tsconfig.json`): ✅ **Configured**
   - Excludes `src/lib/api/` directory from type checking
   - Auto-generated files have `// @ts-nocheck` directive added

3. **Biome Linting/Formatting** (`biome.json`): ✅ **Configured**
   - Excludes `src/lib/api` directory from linting and formatting
   - Uses proper Biome pattern: `!**/src/lib/api` (without trailing `/**`)

4. **Git Ignore** (`.gitignore`): ✅ **Configured**
   - Added comments explaining Orval file strategy
   - Recommends committing generated files for team consistency

**Implemented Configurations:**

**TypeScript Configuration** (`tsconfig.json`):
```json
{
  "exclude": [
    "node_modules",
    ".svelte-kit",
    "build",
    "dist",
    "src/lib/api/"  // Exclude all Orval-generated files
  ]
}
```

**Important Note on svelte-check:**
Even with `tsconfig.json` exclusions, `svelte-check` may still check generated files. To handle this automatically:

**✅ Automated Solution Implemented:**

**1. Custom Script** (`scripts/add-ts-nocheck.mjs`):
- Automatically adds `// @ts-nocheck` directive to Orval-generated TypeScript files
- Processes known Orval output files: `index.ts`, `index.msw.ts`, `index.faker.ts`
- Checks if directive already exists to avoid duplicates
- Provides clear logging of processed files

**2. Updated Package.json Script**:
```json
"api:generate": "orval --mock -i ../backendGoVanilaTaskmanager/swagger/openapi.yaml -o src/lib/api/index.ts && node scripts/add-ts-nocheck.mjs"
```

**3. Usage**:
```bash
# Generate API files with automatic @ts-nocheck directive
cd sveltekit-taskmanager-frontend
bun run api:generate
```

**Benefits of This Implementation**:
- ✅ Automatic: No manual intervention needed after generation
- ✅ Reliable: Works consistently across all environments
- ✅ Maintainsable: Single script handles all Orval output files
- ✅ Idempotent: Safe to run multiple times without side effects
- ✅ Clear: Provides feedback on which files were modified

**How It Works**:
1. Orval generates files from OpenAPI spec
2. Custom script runs immediately after generation
3. Script adds `// @ts-nocheck` directive to all generated TypeScript files
4. TypeScript compiler skips type checking for these files
5. Development workflow remains clean and automated

**Biome Configuration** (`biome.json`):
```json
{
  "files": {
    "ignoreUnknown": false,
    "includes": ["**"],
    "excludes": [
      "**/node_modules",
      "**/dist",
      "**/.svelte-kit",
      "**/build",
      "**/coverage",
      "**/*.config.js",
      "**/*.config.ts",
      "**/playwright-report",
      "**/test-results",
      "**/public/mockServiceWorker.js",
      "**/src/lib/api"  // Exclude all Orval-generated files
    ]
  }
}
```

**Git Ignore** (`.gitignore`):
```gitignore
# Orval-generated files
# Note: We commit these files for team consistency and CI/CD
# Regenerate with: bun run api:generate
# If you need to ignore them locally, uncomment the lines below:
# src/lib/api/index.ts
# src/lib/api/index.msw.ts
# src/lib/api/index.faker.ts
# src/lib/api/handlers.ts
# src/lib/api/model/
```

**Note on Git Strategy:**
- **Commit generated files**: Recommended for team consistency and CI/CD
- **Ignore generated files**: Alternative if team regenerates locally (not recommended for this project)

**Verification Commands:**

```bash
# Verify TypeScript excludes Orval files
cd sveltekit-taskmanager-frontend
bun run check  # ✅ Should not report errors in src/lib/api/

# Verify Biome excludes Orval files
bun run lint   # ✅ Should not report errors in src/lib/api/
bun run format:check  # ✅ Should not report errors in src/lib/api/

# Verify Vitest excludes Orval files from coverage
bun run test:coverage  # ✅ src/lib/api/ files should not appear in coverage report
# Only handlers.ts may appear (wrapper file that's not excluded)
```

**Verification Results (Latest Run):**
- ✅ `bun run check`: No errors in Orval-generated files
- ✅ `bun run lint`: No errors in Orval-generated files
- ✅ `bun run format:check`: No formatting issues in Orval-generated files
- ✅ `bun run validate`: All checks pass (format → lint → type check → tests → build)
- ✅ Coverage report: `index.ts`, `index.msw.ts`, `index.faker.ts` excluded from coverage
- ✅ Only `handlers.ts` appears in coverage (intentional wrapper file)
- ✅ Automated `@ts-nocheck` directive addition via custom script
- ✅ Complete automation - no manual intervention needed

**Why This Exclusion is Critical:**

1. **False Positives**: Linting errors in generated files distract from real issues
2. **Wasted Time**: Fixing auto-generated code is pointless
3. **Team Friction**: Different developers may regenerate with different Orval versions
4. **CI Failures**: Generated file changes can cause unnecessary CI failures
5. **Maintainability**: Only OpenAPI spec should be the source of truth

### Updating MSW Handlers

When the backend API changes:

1. Update `backendGoVanilaTaskmanager/swagger/openapi.yaml`
2. Regenerate MSW handlers:
   ```bash
   cd sveltekit-taskmanager-frontend
   bun run api:generate
   ```
3. Review generated files in `src/lib/api/`
4. Customize mock data if needed (override generated defaults in tests)
5. Update unit tests if needed
6. Commit both OpenAPI changes and generated files

### OpenAPI Validation

The backend should validate against the OpenAPI spec. This can be added:
- `openapi-validator` middleware in Go
- Contract tests to verify API responses match spec
- Automated validation in CI pipeline

## Test Coverage Goals

- **Backend Unit Tests**: >80% coverage
- **Backend Integration Tests**: Critical API paths
- **Frontend Unit Tests**: >80% coverage
- **E2E Tests**: Critical user journeys (login, create task, update, delete)

## Test Coverage Strategy

When improving test coverage, follow this pragmatic approach based on industry best practices:

### Step-by-Step Coverage Improvement

1. **Run coverage report to see actual data**
   ```bash
   cd sveltekit-taskmanager-frontend
   bun run test:coverage
   ```

2. **Analyze which files have low coverage**
   - Review the coverage report output
   - Identify files with <80% coverage
   - Prioritize critical business logic files

3. **Determine approach based on file type**
   - **Type definitions/error pages/boilerplate** → Exclude from coverage requirements
   - **Business logic** → Add comprehensive tests
   - **Mixed cases** → Consider appropriate thresholds

### File Type Classification

**Exclude from Coverage Requirements:**
- Type definition files (`.d.ts`, interfaces, types)
- Error pages and error handling UI
- Boilerplate and configuration files
- Auto-generated files (Orval-generated API handlers) - see "Excluding Orval-Generated Files from Checks" section

**Focus Testing Efforts On:**
- Business logic functions
- Data validation schemas
- State management logic
- API integration layers
- Core utility functions with complex logic

**Threshold Considerations:**
- Simple utilities: 60-70% acceptable
- Complex business logic: 80-90% required
- Critical paths: 90-100% required
- Mixed complexity: Adjust based on risk assessment

### Pragmatic Coverage Goals

Instead of blindly pursuing 100% coverage, focus on:
- Testing critical business paths thoroughly
- Ensuring high coverage for complex, error-prone code
- Maintaining reasonable coverage for simple utilities
- Excluding files that don't benefit from testing

## Summary

This pyramid testing strategy ensures:
- **Fast feedback** from unit tests during development
- **Confidence** from integration tests that components work together
- **Validation** from e2e tests that the complete system works for users
- **No mock vs real API mismatch** by using MSW generated from OpenAPI
- **Clear separation** between test types with appropriate tooling
- **CI/CD ready** with GitHub Actions workflows
- **Contract compliance** through OpenAPI-driven MSW handlers

## Development Workflow

### Typical Development Cycle

1. **Backend Changes**:
   ```bash
   cd backendGoVanilaTaskmanager
   make test-unit
   make test-integration
   make build
   ```

2. **Frontend Changes**:
   ```bash
   cd sveltekit-taskmanager-frontend
   bun run api:generate  # If API changed
   bun run validate
   bun run build
   ```

3. **E2E Testing**:
   ```bash
   cd ..
   docker-compose -f docker-compose.test.yml up -d
   cd sveltekit-taskmanager-frontend
   bun test:e2e
   ```

4. **Pre-commit**:
   ```bash
   # Backend
   cd backendGoVanilaTaskmanager
   make test-all

   # Frontend
   cd ../sveltekit-taskmanager-frontend
   bun run validate
   ```

### CI/CD Workflow

1. **Push to feature branch**: Runs relevant CI based on changed files
2. **Pull Request**: Full CI suite including e2e tests
3. **Merge to main**: Builds Docker images and deploys
4. **Automated Testing**: Every commit triggers appropriate tests

This setup ensures quality at every level while maintaining fast development cycles.
