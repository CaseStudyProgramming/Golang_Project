# AGENTS.md

## Core Directives & Priority (Precedence Order)
1. **Security First**: OWASP guidelines override all formatting and style rules.
2. **Architecture Integrity**: Layered architecture rules override file-level preferences.
3. **Functionality > Style**: Working implementation takes priority over code style/refactoring.
4. **Circuit Breaker**: Stop and ask after 3 consecutive failed test/build attempts.
5. **When in Doubt, Ask**: Stop and request clarification if user instructions contradict these rules or if a requirement is ambiguous.

## Workflow & Task Execution
- **Branch-Based Development**:
  - Always work on a new branch for each sub-issue using git branch standard naming convention (e.g., `feat/add-task-filtering`).
- **Task Breakdown via Checklists / Issues**:
  - Read and parse task lists (`- [ ]`) or sub-issues sequentially.
  - Work on only ONE granular checklist item at a time. Complete it fully before moving to the next.
- **Atomic Commits & Verification**:
  - Run verification commands (`go test ./...`, `go build`) after completing each individual checklist item.
  - Immediately create a dedicated Git commit upon completing each item (e.g., `feat(scope): complete task X - description`).
  - Update or check off the task list item upon a successful commit.
- **Rollback Readiness**:
  - Each completed task MUST correspond to a clean, isolated Git commit to allow single-step rollbacks (`git revert`) without losing previous progress.
- **CI/CD Verification Workflow (Global Standard)**:
  - **Pre-commit Verification Order** (local development):
    1. Format check (`gofmt -l .`) - Fast-fail for basic formatting issues
    2. Lint (`golangci-lint run`) - Code quality and style checks
    3. Static analysis (`go vet ./...`) - Potential bugs and issues
    4. Tests (`go test ./...`) - Unit and integration tests
    5. Build (`go build`) - Compilation verification
  - **CI/CD Pipeline Order** (`.github/workflows/backend-ci.yml`):
    1. Format check - Fastest check, fail early
    2. Lint & Static analysis (PARALLEL) - Independent checks for speed
    3. Tests (with race detection & coverage) - Runtime verification
    4. Security scan (`govulncheck`) - Vulnerability check (can be parallel with tests)
    5. Build - Final compilation verification
  - **Fast-Fail Principle**: Always run fastest checks first (format → lint → test → build) to fail early and save CI/CD resources
  - **Parallel Execution**: Lint and static analysis checks can run in parallel as they are independent
  - **Security Scan**: Can be blocking (strict) or non-blocking (relaxed) depending on project velocity requirements

## Token Efficiency
- Skip recaps and conversational summaries unless the result is ambiguous or requires further input.

## Principles
- **OWASP Security Standard**: Validate all inputs, prevent SQL injection (use parameterized queries), prevent XSS, implement CSRF protection, follow OWASP Top 10.
- **Clarity and Consistency**: Clarity over cleverness. Match existing code patterns. Minimal changes unless refactoring is explicitly requested.
- **Modularity**: Keep functions under 50 lines. Break down when improves readability and structure. Each layer (controller, service, model) should have a single responsibility.
- **Type Safety**: Use Go's type system effectively. Use custom types for domain concepts. Avoid `interface{}` unless necessary.
- **Error Handling**: Always check errors. Use meaningful error messages. Use context.Context for cancellation and timeouts. Use wrapped errors with `fmt.Errorf` or `errors.Wrap` when relevant.
- **Concurrency**: Use goroutines and channels properly. Avoid race conditions. Use `sync` package or channels for synchronization. Use `-race` flag in tests.
- **Database Operations**: Use transactions for multi-step operations. Use connection pooling. Handle connection errors gracefully.
- **API Design**: Follow REST conventions. Use appropriate HTTP status codes. Implement consistent response format. Implement proper pagination.
- **Logging**: Use structured logging. Log critical operations, errors, and security events. Never log sensitive data (passwords, tokens, PII).
- **Testing**: Write unit tests for business logic. Write integration tests for API endpoints. Aim for >80% coverage. Use table-driven tests.
- **Documentation**: Keep Swagger/OpenAPI spec updated. Document public functions and types. Use godoc comments.
- **Debugging**: Hypothesis-driven debugging—formulate 1–3 most likely causes first, then validate incrementally.

## Commands (Go)
- Always use standard Go tooling:
  - `go run main.go` - Run the application
  - `go build -o server.exe` - Build for production
  - `go test ./...` - Run all tests
  - `go test -race ./...` - Run tests with race detection
  - `go test -cover ./...` - Run tests with coverage
  - `gofmt -w .` - Format all Go files
  - `gofmt -l .` - Check if code is properly formatted
  - `air` - Run with live reload (development)

## Git Commits
- **Conventional Commits**: Format as `type: summary without scope`.
- Summary must be a short, specific sentence explaining what changed and why.
- Valid types: `feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert`.
- Include `BREAKING CHANGE:` in the commit footer when applicable.

## Environment Variables & Configuration
- **Configuration Loading**: Load configuration from `env/config.yaml` using the existing config package.
- **No Direct Environment Access**: Do NOT use `os.Getenv` directly in application code. Use the centralized config structure.
- **Secrets Management**: Never commit secrets to repository. Use environment variables or secret management systems.
- **Validation**: Validate configuration at startup. Fail fast if required configuration is missing or invalid.
- **Database Configuration**: Use the existing `config.NewPostgresDB` function for database connections.

## Architecture & Layered Structure
- **Layered Architecture**:
  - **Controllers/**: HTTP request/response handlers. Handle input validation, call services, format responses.
  - **Services/**: Business logic layer. Implement domain rules, coordinate between models, handle transactions.
  - **Models/**: Database operations and data structures. Handle SQL queries, data mapping.
  - **Middlewares/**: Cross-cutting concerns (auth, logging, CORS, recovery, rate limiting).
  - **Routes/**: API endpoint definitions and route registration.
  - **Utils/**: Helper functions (JWT, response formatting, validation).
  - **Config/**: Configuration loading and validation.
- **Dependency Direction**:
  - Controllers → Services → Models
  - Controllers → Utils → Middlewares
  - Services can depend on multiple models and other services
  - Models should NOT depend on services or controllers
- **Error Handling Flow**:
  - Models return domain-specific errors
  - Services wrap and add context to errors
  - Controllers translate errors to appropriate HTTP responses
- **Transaction Management**:
  - Services should handle database transactions for multi-step operations
  - Use `db.Begin()`, `tx.Commit()`, `tx.Rollback()` appropriately
  - Ensure transactions are rolled back on error

## Go Best Practices
- **Package Organization**: One package per directory. Package names should be short, lowercase, single words.
- **Exported Names**: Exported names (functions, types, constants) must start with uppercase. Unexported names start with lowercase.
- **Interface Usage**: Use interfaces for dependency injection and testing. Keep interfaces small and focused (prefer composition over large interfaces).
- **Error Handling**: Always check errors immediately. Use named returns when needed. Use `defer` for cleanup.
- **Context Usage**: Pass `context.Context` as first parameter to functions that perform I/O or have long-running operations. Respect context cancellation.
- **Struct Initialization**: Use named struct literals for clarity. Use constructor functions (e.g., `NewUserService`) for complex initialization.
- **Constants**: Use `iota` for related constants. Use string constants for enum-like types.
- **Slices and Maps**: Initialize with capacity when size is known. Check for nil before appending.
- **Time Handling**: Use `time.Time` for timestamps. Use UTC for storage. Handle timezone conversion at presentation layer.

## Security
- **Input Validation**: Validate all inputs at controller layer. Use struct tags for validation rules.
- **SQL Injection Prevention**: Always use parameterized queries. Never concatenate strings into SQL queries.
- **Password Security**: Use bcrypt for password hashing with appropriate cost factor. Never store plain text passwords.
- **JWT Security**: Use strong secret keys. Set appropriate expiration times. Validate tokens on every protected request.
- **CORS Configuration**: Configure CORS properly. Limit allowed origins, methods, and headers.
- **Rate Limiting**: Implement rate limiting to prevent abuse. Use middleware or external services.
- **Sensitive Data**: Never log passwords, tokens, or PII. Use secure storage for secrets.
- **HTTPS**: Always use HTTPS in production. Redirect HTTP to HTTPS.
- **Dependencies**: Regularly update dependencies. Use `go get -u` and review security advisories.

## Database & SQL
- **Connection Pooling**: Use the existing connection pool configuration. Tune pool settings based on load.
- **Query Optimization**: Use appropriate indexes. Avoid N+1 queries. Use `EXPLAIN ANALYZE` for slow queries.
- **Transactions**: Use transactions for operations that modify multiple tables. Keep transactions short.
- **Migrations**: Use migration scripts in `migrations/` directory. Name migrations with sequential numbers (e.g., `001.init.up.sql`).
- **Error Handling**: Handle database connection errors gracefully. Implement retry logic for transient failures.
- **Data Consistency**: Use foreign key constraints. Use database constraints for data integrity.

## Testing
- **Unit Tests**: Test business logic in services layer. Mock dependencies (models, external services).
- **Integration Tests**: Test API endpoints with real database. Use test database configuration.
- **Table-Driven Tests**: Use table-driven tests for multiple test cases.
- **Race Detection**: Always run tests with `-race` flag before committing.
- **Coverage**: Aim for >80% code coverage. Use `go test -cover` to check coverage.
- **Test Structure**: Follow standard Go test structure. Use `Test` prefix for test functions. Use ` testify` for assertions if needed.

## API Design
- **REST Conventions**: Use appropriate HTTP methods (GET, POST, PUT, PATCH, DELETE). Use resource-based URLs.
- **HTTP Status Codes**: Use appropriate status codes (200, 201, 400, 401, 403, 404, 500, etc.).
- **Response Format**: Use consistent response format with `status`, `message`, and `data` fields.
- **Pagination**: Implement pagination for list endpoints. Return metadata (page, limit, total, has_next, has_prev).
- **Versioning**: Consider API versioning for breaking changes (e.g., `/api/v1/`).
- **Error Responses**: Return consistent error responses with clear error messages and error codes.
- **Validation**: Return validation errors with field-level details for client-side handling.

## Performance & Optimization
- **Database Queries**: Optimize queries. Use indexes. Avoid SELECT *. Use specific columns.
- **Connection Pooling**: Tune connection pool settings. Use connection pool for database.
- **Caching**: Consider caching for frequently accessed data. Use Redis or in-memory cache.
- **Goroutines**: Use goroutines for concurrent operations. Limit goroutine count to avoid resource exhaustion.
- **Profiling**: Use pprof for profiling. Profile CPU and memory usage.
- **Lazy Loading**: Load data only when needed. Use pagination for large datasets.

## System Quality & Reliability (Backend/Server-Side Scope)
- **System Observability & Incident Response**: Implement structured logging (JSON format), metrics collection (Prometheus), distributed tracing. Establish clear incident response procedures with defined runbooks and escalation paths. Ensure auditability of all critical operations.
- **High Availability & Fault Tolerance (HA/FT)**: Implement graceful degradation, health check endpoints (`/health`), circuit breakers for external dependencies, retry mechanisms with exponential backoff, and graceful shutdown.
- **Infrastructure & Database Resiliency**: Use connection pooling, implement retry logic for transient database failures, design for database failover, ensure backups and disaster recovery procedures are in place.
- **API Defensive Design (Defensive Programming) & Code Quality**: Validate all inputs and outputs, implement rate limiting middleware, use defensive coding practices (check for nil, bounds checking), follow SOLID principles, maintain high code coverage with meaningful tests.
- **Traffic Control & Throttling**: Implement rate limiting middleware, request queuing, backpressure handling, and load shedding to prevent system overload under high traffic conditions.
- **Data Integrity & Consistency**: Use database transactions for multi-step operations, implement idempotent operations (especially for non-GET requests), ensure data consistency across operations, validate data at boundaries (controller layer).
- **Performance & Concurrency**: Optimize database queries with indexes, implement caching strategies (Redis or in-memory), handle concurrent operations safely with proper synchronization (mutex, channels), use connection pooling, avoid N+1 query problems.
- **Security**: Follow OWASP Top 10 guidelines, implement JWT authentication and authorization, use secure communication (HTTPS), regularly update dependencies, never expose secrets or sensitive data in logs or responses, implement CORS properly.
- **Usability / Robustness**: Design intuitive error messages with clear error codes, implement graceful error handling with proper HTTP status codes, ensure the system provides helpful feedback for debugging, design for edge cases and unexpected input.
- **Data Privacy & Information Disclosure Protection**: Implement data minimization (only collect necessary data), encrypt sensitive data at rest (database encryption) and in transit (TLS), comply with privacy regulations (GDPR, CCPA), never log or expose sensitive information (passwords, tokens, PII), implement data retention policies.
- **High Availability & Fault Tolerance (HA/FT)**: Implement graceful degradation, health check endpoints (`/health`), circuit breakers for external dependencies, retry mechanisms with exponential backoff, and graceful shutdown (already implemented in main.go).
- **Infrastructure & Database Resiliency**: Use connection pooling (already implemented), implement retry logic for transient database failures, design for database failover, ensure backups and disaster recovery procedures are in place.
- **API Defensive Design (Defensive Programming) & Code Quality**: Validate all inputs and outputs, implement rate limiting middleware, use defensive coding practices (check for nil, bounds checking), follow SOLID principles, maintain high code coverage with meaningful tests.
- **Traffic Control & Throttling**: Implement rate limiting middleware, request queuing, backpressure handling, and load shedding to prevent system overload under high traffic conditions.
- **Data Integrity & Consistency**: Use database transactions for multi-step operations, implement idempotent operations (especially for non-GET requests), ensure data consistency across operations, validate data at boundaries (controller layer).
- **Performance & Concurrency**: Optimize database queries with indexes, implement caching strategies (Redis or in-memory), handle concurrent operations safely with proper synchronization (mutex, channels), use connection pooling, avoid N+1 query problems.
- **Security**: Follow OWASP Top 10 guidelines, implement JWT authentication and authorization, use secure communication (HTTPS), regularly update dependencies, never expose secrets or sensitive data in logs or responses, implement CORS properly.
- **Usability / Robustness**: Design intuitive error messages with clear error codes, implement graceful error handling with proper HTTP status codes, ensure the system provides helpful feedback for debugging, design for edge cases and unexpected input.
- **Data Privacy & Information Disclosure Protection**: Implement data minimization (only collect necessary data), encrypt sensitive data at rest (database encryption) and in transit (TLS), comply with privacy regulations (GDPR, CCPA), never log or expose sensitive information (passwords, tokens, PII), implement data retention policies.

## Development Tools
- **Live Reload**: Use `air` for development with live reload.
- **Code Formatting**: Use `gofmt` for consistent formatting. Use provided scripts (`check_format.sh` or `check_format.bat`).
- **Linting**: Consider using `golangci-lint` for additional linting rules.
- **Documentation**: Keep Swagger documentation updated. Use `godoc` for package documentation.
- **Migration Management**: Use migration scripts in `migrations/` directory. Execute migrations manually or with migration tool.

## Project-Specific Notes
- **Authentication**: JWT-based authentication is implemented. Use the existing `utils.JWTManager` for token generation and validation.
- **Multi-tenancy**: User data isolation is implemented. Always filter data by user_id in queries.
- **Soft Delete**: Tasks support soft delete with `deleted_at` field. Consider soft delete for other entities if needed.
- **Activity Logging**: Activity logging is implemented for audit trail. Log critical operations (create, update, delete).
- **Swagger**: Swagger documentation is available at `/swagger/index.html`. Keep it updated when adding new endpoints.
- **Graceful Shutdown**: Graceful shutdown is implemented in main.go. Maintain this pattern for any new background processes.
