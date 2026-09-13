# Multi-stage build for optimized production image
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=${VERSION:-dev} -X main.BuildTime=${BUILD_TIME:-unknown}" \
    -o taskmanager-server main.go

# Runtime stage - minimal Alpine image
FROM alpine:latest

# Install ca-certificates and wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 taskmanager && \
    adduser -D -u 1000 -G taskmanager taskmanager

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=taskmanager:taskmanager /app/taskmanager-server .

# Copy migrations
COPY --from=builder --chown=taskmanager:taskmanager /app/migrations ./migrations

# Copy default Docker config
COPY --from=builder --chown=taskmanager:taskmanager /app/env/config.docker.yaml ./env/config.yaml

# Switch to non-root user
USER taskmanager

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./taskmanager-server"]
