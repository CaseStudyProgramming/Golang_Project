.PHONY: build build-linux build-windows build-darwin clean test test-race lint format run dev help

# Build variables
BINARY_NAME=taskmanager-server
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application for current platform
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) main.go

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 main.go

build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe main.go

build-darwin: ## Build for macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 main.go

build-all: build-linux build-windows build-darwin ## Build for all platforms

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)*
	rm -f *.exe
	rm -f *.test
	rm -f *.out

test: ## Run tests
	$(GOTEST) -v ./...

test-race: ## Run tests with race detection
	$(GOTEST) -race -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

format: ## Format code
	gofmt -w .
	$(GOCMD) fmt ./...

format-check: ## Check if code is formatted
	gofmt -l .

run: ## Run the application
	$(GOBUILD) -o $(BINARY_NAME) main.go
	./$(BINARY_NAME)

dev: ## Run with live reload using air
	air

deps: ## Download dependencies
	$(GOGET) -u -t ./...
	$(GOMOD) tidy

deps-check: ## Check for dependency updates
	$(GOGET) -u ./...
	$(GOMOD) tidy

docker-build: ## Build Docker image
	docker build -t taskmanager-backend:$(VERSION) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 taskmanager-backend:$(VERSION)

security-scan: ## Run security scan
	gosec ./...

benchmark: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

profiling-cpu: ## Run CPU profiling
	$(GOBUILD) -o $(BINARY_NAME) main.go
	./$(BINARY_NAME) &
	$(GOCMD) tool pprof http://localhost:8080/debug/pprof/profile

profiling-memory: ## Run memory profiling
	$(GOBUILD) -o $(BINARY_NAME) main.go
	./$(BINARY_NAME) &
	$(GOCMD) tool pprof http://localhost:8080/debug/pprof/heap
