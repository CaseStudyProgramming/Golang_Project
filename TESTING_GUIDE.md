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

The project uses GitHub Actions for continuous integration and deployment:

- **Backend CI** (`.github/workflows/backend-ci.yml`): Unit tests, integration tests, security scan, build
- **Frontend CI** (`.github/workflows/frontend-ci.yml`): Quality checks, unit tests, build, e2e tests with real backend
- **Main CI** (`.github/workflows/ci.yml`): Orchestrates backend and frontend CI, builds Docker images, deploys

### Key CI Features

- **MSW from OpenAPI**: Frontend unit tests use MSW handlers based on OpenAPI spec
- **Real Backend in E2E**: E2E tests use real backend with test database
- **Parallel Execution**: Tests run in parallel when possible
- **Artifact Upload**: Test results and coverage reports are uploaded
- **Docker Caching**: Build optimization with Docker layer caching

## Backend Testing (Go)

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
bun run test

# Run in watch mode
bun run test:ui

# Run with coverage
bun run test:coverage

# Run specific test file
bun test src/lib/features/tasks/stores/task.store.logic.test.ts
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
End-to-end tests using Playwright with real backend and frontend.

```bash
cd sveltekit-taskmanager-frontend

# Start test environment (backend + frontend + test DB)
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

**E2E Test Configuration:**
- Uses real backend server (Go)
- Uses real PostgreSQL test database
- Frontend runs on http://localhost:5173
- Backend runs on http://localhost:8080
- Tests run on Chromium, Firefox, WebKit, and mobile browsers

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
bun run test

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
cd ../sveltekit-taskmanager-frontend && bun run test
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

- **Backend CI** (`.github/workflows/backend-ci.yml`): Format → Lint → Unit tests → Integration tests → Security scan → Build
- **Frontend CI** (`.github/workflows/frontend-ci.yml`): API generation → Format → Lint → Type check → Unit tests → Build → E2E tests
- **Main CI** (`.github/workflows/ci.yml`): Orchestrates both CIs, builds Docker images, deploys

### CI Environment Setup

**Backend CI:**
- Uses PostgreSQL service container for integration tests
- Environment variables configured for test database
- Runs with Go 1.25

**Frontend CI:**
- Generates MSW handlers from OpenAPI spec automatically
- Uses PostgreSQL service container for e2e tests
- Starts backend server for e2e tests
- Runs with Bun and Node.js 24

**E2E Tests in CI:**
- Real backend server (compiled Go binary)
- Real PostgreSQL test database
- Playwright browsers installed with system dependencies
- Tests run on Chromium, Firefox, WebKit, and mobile browsers

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
bun test

# Test e2e locally with CI-like setup
cd ..
docker-compose -f docker-compose.ci.yml up -d
cd sveltekit-taskmanager-frontend
CI=true PUBLIC_API_BASE_URL=http://localhost:8080 bun test:e2e
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
bun run test
bun run build
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
- Check frontend is accessible: `curl http://localhost:5173`
- Verify test user can be created via API

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
   bun run test
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
