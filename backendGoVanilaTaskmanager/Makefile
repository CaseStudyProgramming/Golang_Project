.PHONY: test test-unit test-integration test-performance test-all test-coverage build lint format clean help

# Test commands
test-unit:
	go test ./... -v

test-integration:
	go test ./tests/... -v

test-performance:
	go test ./tests/performance_test.go ./tests/api_benchmark_test.go -v

test-all: test-unit test-integration
	@echo "✅ All tests completed"

test-coverage:
	go test ./... -cover -v

test-race:
	go test ./... -race -v

# Build commands
build:
	go build -o taskmanager-server main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o taskmanager-server main.go

# Code quality
lint:
	golangci-lint run

format:
	gofmt -w .

format-check:
	gofmt -l .

vet:
	go vet ./...

# Database
db-migrate:
	@echo "Run migrations manually or use your migration tool"

# Development
dev:
	go run main.go

dev-test:
	TEST_DB_HOST=localhost TEST_DB_PORT=5433 go run main.go

# Clean
clean:
	rm -f taskmanager-server
	rm -rf test-results/

# Help
help:
	@echo "Available commands:"
	@echo "  test-unit          - Run unit tests"
	@echo "  test-integration   - Run integration tests (requires test DB)"
	@echo "  test-performance   - Run performance tests"
	@echo "  test-all           - Run all tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-race          - Run tests with race detection"
	@echo "  build              - Build the application"
	@echo "  build-linux        - Build for Linux"
	@echo "  lint               - Run linter"
	@echo "  format             - Format code"
	@echo "  format-check       - Check code formatting"
	@echo "  vet                - Run go vet"
	@echo "  dev                - Run development server"
	@echo "  dev-test           - Run development server with test DB"
	@echo "  clean              - Clean build artifacts"
	@echo "  help               - Show this help message"
