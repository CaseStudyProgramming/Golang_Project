# Observability Status & Migration Guide

## Current Implementation Status

### ✅ Implemented
- **Basic HTTP Logging**: Standard `log.Printf` for HTTP request logging (method, path, status, duration, IP)
- **Custom In-Memory Metrics**: Performance metrics tracking in `middlewares/performance.go`
  - Request counts per endpoint
  - Response times (average, P95, min, max)
  - Error counts per endpoint
  - Slow endpoint tracking
  - Metrics endpoint: `/metrics`
- **Activity Logging**: Database-based activity logging for audit trail
- **Graceful Shutdown**: Implemented in main.go for clean server shutdown

### ❌ Not Implemented (Production Gaps)
- **Structured Logging**: No JSON format, log levels, or structured logging library
- **Prometheus Metrics**: Custom metrics are not Prometheus-compatible
- **Distributed Tracing**: No OpenTelemetry or distributed tracing
- **Health Check Endpoints**: No `/health` or `/ready` endpoints for load balancer/orchestration
- **Error Tracking Service**: No integration with error tracking services (Sentry, etc.)
- **Log Aggregation**: No integration with log aggregation systems (ELK, Loki, etc.)
- **Alerting**: No automated alerting for critical metrics or errors

## Migration Guide (Recommended for Production)

### Phase 1: Structured Logging
**Priority**: High
**Estimated Effort**: 2-3 days

**Goals**:
- Replace `log.Printf` with structured logging library
- Implement log levels (debug, info, warn, error)
- Add contextual fields to logs (request_id, user_id, endpoint, duration)
- Enable JSON format for production

**Implementation Steps**:
1. Choose logging library:
   - `log/slog` (Go 1.21+, built-in, recommended)
   - `zap` (Uber, high performance)
   - `zerolog` (zero allocation, fast)

2. Replace existing logging:
   ```go
   // Old
   log.Printf("[%s] %s %s %d %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, statusCode, duration, r.RemoteAddr)

   // New with slog
   slog.Info("HTTP request",
       "method", r.Method,
       "path", r.URL.Path,
       "status", statusCode,
       "duration", duration,
       "remote_addr", r.RemoteAddr,
       "request_id", requestID,
   )
   ```

3. Add request_id middleware for trace correlation:
   ```go
   func RequestIDMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           requestID := r.Header.Get("X-Request-ID")
           if requestID == "" {
               requestID = uuid.New().String()
           }
           ctx := context.WithValue(r.Context(), "request_id", requestID)
           w.Header().Set("X-Request-ID", requestID)
           next.ServeHTTP(w, r.WithContext(ctx))
       })
   }
   ```

4. Configure log levels based on environment:
   ```go
   var logLevel slog.Level
   if os.Getenv("ENV") == "production" {
       logLevel = slog.LevelInfo
   } else {
       logLevel = slog.LevelDebug
   }

   opts := &slog.HandlerOptions{
       Level: logLevel,
   }

   var handler slog.Handler
   if os.Getenv("ENV") == "production" {
       handler = slog.NewJSONHandler(os.Stdout, opts)
   } else {
       handler = slog.NewTextHandler(os.Stdout, opts)
   }

   logger := slog.New(handler)
   slog.SetDefault(logger)
   ```

**Files to Modify**:
- `middlewares/logger.go`
- `middlewares/performance.go`
- `main.go` (for logger initialization)
- All controllers and services (for contextual logging)

### Phase 2: Prometheus Metrics
**Priority**: High
**Estimated Effort**: 3-4 days

**Goals**:
- Replace custom metrics with Prometheus-compatible metrics
- Expose metrics at `/metrics` endpoint
- Add metrics for HTTP requests, database operations, and business metrics

**Implementation Steps**:
1. Add Prometheus dependency:
   ```bash
   go get github.com/prometheus/client_golang
   go get github.com/prometheus/client_golang/prometheus/promhttp
   ```

2. Replace custom metrics with Prometheus metrics:
   ```go
   import (
       "github.com/prometheus/client_golang/prometheus"
       "github.com/prometheus/client_golang/prometheus/promauto"
   )

   var (
       httpRequestsTotal = promauto.NewCounterVec(
           prometheus.CounterOpts{
               Name: "http_requests_total",
               Help: "Total number of HTTP requests",
           },
           []string{"method", "endpoint", "status"},
       )

       httpRequestDuration = promauto.NewHistogramVec(
           prometheus.HistogramOpts{
               Name:    "http_request_duration_seconds",
               Help:    "HTTP request duration in seconds",
               Buckets: prometheus.DefBuckets,
           },
           []string{"method", "endpoint"},
       )

       dbConnectionsActive = promauto.NewGauge(
           prometheus.GaugeOpts{
               Name: "db_connections_active",
               Help: "Number of active database connections",
           },
       )
   )
   ```

3. Update middleware to record Prometheus metrics:
   ```go
   func PrometheusMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           start := time.Now()

           rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
           next.ServeHTTP(rw, r)

           duration := time.Since(start).Seconds()

           httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.statusCode)).Inc()
           httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
       })
   }
   ```

4. Add `/metrics` endpoint:
   ```go
   import "github.com/prometheus/client_golang/prometheus/promhttp"

   // In routes registration
   mux.Handle("/metrics", promhttp.Handler())
   ```

5. Add database metrics:
   ```go
   // Track connection pool metrics
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(25)
   dbConnectionsActive.Set(float64(db.Stats().OpenConnections))
   ```

**Files to Modify**:
- `go.mod` (add dependency)
- `middlewares/performance.go` (replace custom metrics)
- `routes/routes.go` (add /metrics endpoint)
- `config/database.go` (add connection metrics)

### Phase 3: Health Check Endpoints
**Priority**: Medium
**Estimated Effort**: 1-2 days

**Goals**:
- Implement `/health` endpoint for liveness probes
- Implement `/ready` endpoint for readiness probes
- Check database connectivity and external dependencies

**Implementation Steps**:
1. Create health check controller:
   ```go
   // controllers/health_controller.go
   type HealthController struct {
       db *sql.DB
   }

   func NewHealthController(db *sql.DB) *HealthController {
       return &HealthController{db: db}
   }

   func (c *HealthController) LivenessHandler(w http.ResponseWriter, r *http.Request) {
       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(map[string]string{
           "status": "healthy",
       })
   }

   func (c *HealthController) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
       // Check database connectivity
       if err := c.db.Ping(); err != nil {
           w.WriteHeader(http.StatusServiceUnavailable)
           json.NewEncoder(w).Encode(map[string]string{
               "status": "unhealthy",
               "error":  "database connection failed",
           })
           return
       }

       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(map[string]string{
           "status": "ready",
       })
   }
   ```

2. Register health check routes:
   ```go
   mux.HandleFunc("/health", healthController.LivenessHandler)
   mux.HandleFunc("/ready", healthController.ReadinessHandler)
   ```

3. Configure Kubernetes probes (if using Kubernetes):
   ```yaml
   livenessProbe:
     httpGet:
       path: /health
       port: 8080
     initialDelaySeconds: 30
     periodSeconds: 10

   readinessProbe:
     httpGet:
       path: /ready
       port: 8080
     initialDelaySeconds: 5
     periodSeconds: 5
   ```

**Files to Create**:
- `controllers/health_controller.go`

**Files to Modify**:
- `routes/routes.go`
- `main.go` (initialize health controller)

### Phase 4: Distributed Tracing (Optional)
**Priority**: Low (for single-service), High (for microservices)
**Estimated Effort**: 5-7 days

**Goals**:
- Implement OpenTelemetry instrumentation
- Configure Jaeger or Tempo as tracing backend
- Add tracing to HTTP handlers and database operations

**Implementation Steps**:
1. Add OpenTelemetry dependencies:
   ```bash
   go get go.opentelemetry.io/otel
   go get go.opentelemetry.io/otel/trace
   go get go.opentelemetry.io/otel/exporters/jaeger
   go get go.opentelemetry.io/otel/sdk/trace
   ```

2. Initialize tracer:
   ```go
   import (
       "go.opentelemetry.io/otel"
       "go.opentelemetry.io/otel/exporters/jaeger"
       "go.opentelemetry.io/otel/sdk/trace"
   )

   func initTracer(serviceName string) error {
       exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
       if err != nil {
           return err
       }

       tp := trace.NewTracerProvider(
           trace.WithBatcher(exporter),
           trace.WithResource(resources.NewWithAttributes(
               semconv.SchemaURL,
               semconv.ServiceNameKey.String(serviceName),
           )),
       )

       otel.SetTracerProvider(tp)
       return nil
   }
   ```

3. Add tracing middleware:
   ```go
   func TracingMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           tracer := otel.Tracer("http-server")
           ctx, span := tracer.Start(r.Context(), r.URL.Path,
               trace.WithAttributes(
                   semconv.HTTPMethodKey.String(r.Method),
                   semconv.HTTPURLKey.String(r.URL.String()),
               ),
           )
           defer span.End()

           next.ServeHTTP(w, r.WithContext(ctx))
       })
   }
   ```

4. Add database tracing:
   ```go
   func QueryWithTracing(ctx context.Context, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
       tracer := otel.Tracer("database")
       ctx, span := tracer.Start(ctx, "database.query",
           trace.WithAttributes(semconv.DBStatementKey.String(query)),
       )
       defer span.End()

       return db.QueryContext(ctx, query, args...)
   }
   ```

**Files to Modify**:
- `go.mod` (add dependencies)
- `main.go` (initialize tracer)
- `middlewares/` (add tracing middleware)
- `models/` (add database tracing)

## Production Checklist

Before deploying to production, ensure:

- [ ] Structured logging implemented with JSON format
- [ ] Log levels configured (debug for dev, info for prod)
- [ ] Request ID correlation implemented
- [ ] Prometheus metrics implemented and exposed
- [ ] Metrics include: HTTP requests, response times, errors, DB operations
- [ ] Health check endpoints implemented (`/health`, `/ready`)
- [ ] Alerting configured for critical metrics
- [ ] Log aggregation system configured (ELK, Loki, etc.)
- [ ] Error tracking service integrated (Sentry, etc.)
- [ ] Distributed tracing implemented (if microservices)
- [ ] Performance budgets defined and monitored
- [ ] Incident response procedures documented
- [ ] On-call rotation established
- [ ] Runbooks created for common incidents

## Monitoring & Alerting Recommendations

### Key Metrics to Monitor
- **Availability**: Uptime percentage
- **Response Time**: P50, P95, P99 response times
- **Error Rate**: HTTP 4xx and 5xx rates
- **Throughput**: Requests per second
- **Database**: Connection pool usage, query performance
- **System**: CPU, memory, disk I/O, network I/O

### Alert Thresholds (Example)
- Error rate > 1% for 5 minutes
- P95 response time > 500ms for 5 minutes
- Database connection pool > 80% utilization
- CPU > 80% for 10 minutes
- Memory > 80% for 10 minutes
- Any 5xx error

### Recommended Tools
- **Metrics**: Prometheus + Grafana
- **Logging**: Loki + Grafana, or ELK Stack
- **Tracing**: Jaeger, Tempo, or Honeycomb
- **Error Tracking**: Sentry, Rollbar, or Bugsnag
- **APM**: Datadog, New Relic, or Dynatrace
