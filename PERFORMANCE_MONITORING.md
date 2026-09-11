# Performance Monitoring Guide

This guide provides comprehensive documentation for implementing and using performance monitoring in the Task Manager application.

## Table of Contents

- [Overview](#overview)
- [Frontend Performance Monitoring](#frontend-performance-monitoring)
- [Backend Performance Monitoring](#backend-performance-monitoring)
- [Integration with Monitoring Services](#integration-with-monitoring-services)
- [Alerting and Notifications](#alerting-and-notifications)
- [Performance Optimization](#performance-optimization)
- [Best Practices](#best-practices)

## Overview

Performance monitoring is essential for maintaining application health, identifying bottlenecks, and ensuring optimal user experience. This guide covers both frontend and backend monitoring implementations.

### Key Metrics

#### Frontend Metrics
- **FCP** (First Contentful Paint): Time when first content is painted
- **LCP** (Largest Contentful Paint): Time when largest content is painted
- **FID** (First Input Delay): Time until first user interaction
- **CLS** (Cumulative Layout Shift): Visual stability score
- **TTFB** (Time to First Byte): Time to receive first byte
- **API Response Time**: Duration of API calls
- **Component Render Time**: Time to render components
- **Page Load Time**: Total page load duration

#### Backend Metrics
- **API Call Count**: Number of API calls per endpoint
- **API Duration**: Response time percentiles (p50, p95, p99)
- **Memory Usage**: Memory allocation and GC stats
- **Goroutine Count**: Number of active goroutines
- **Error Rate**: Percentage of failed requests
- **Request Rate**: Requests per second
- **Database Query Time**: Database operation duration

## Frontend Performance Monitoring

### Implementation

The frontend uses the `web-vitals` library for Core Web Vitals tracking and custom performance monitoring.

#### Setup

```typescript
// Install dependencies
bun add web-vitals
```

#### Usage

```typescript
import { getPerformanceMonitor } from '$lib/shared/performance-monitoring';

// Initialize performance monitor
const monitor = getPerformanceMonitor();

// Track API calls
const startTime = performance.now();
await fetch('/api/tasks');
monitor.trackApiCall('/api/tasks', startTime);

// Track component rendering
const componentStart = performance.now();
// ... component logic
monitor.trackComponentRender('TaskList', componentStart);

// Track page rendering
const pageStart = performance.now();
// ... page logic
monitor.trackPageRender('TasksPage', pageStart);
```

#### Configuration

Add to `.env.production`:

```env
PUBLIC_ENABLE_PERFORMANCE_MONITORING=true
PUBLIC_PERFORMANCE_ENDPOINT=https://your-metrics-endpoint.com/api/metrics
```

### Performance Data Structure

```typescript
interface PerformanceData {
  webVitals: {
    FCP?: number;    // First Contentful Paint (ms)
    LCP?: number;    // Largest Contentful Paint (ms)
    FID?: number;    // First Input Delay (ms)
    CLS?: number;    // Cumulative Layout Shift
    TTFB?: number;   // Time to First Byte (ms)
  };
  customMetrics: {
    apiResponseTime: number;
    renderTime: number;
    componentLoadTime: number;
  };
  timestamp: number;
  url: string;
  userAgent: string;
}
```

### Performance Monitoring in Components

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { getPerformanceMonitor } from '$lib/shared/performance-monitoring';

  const monitor = getPerformanceMonitor();
  let renderStart: number;

  onMount(() => {
    renderStart = performance.now();
  });

  function onRenderComplete() {
    if (renderStart) {
      monitor.trackComponentRender('TaskList', renderStart);
    }
  }
</script>

<!-- Component template -->
<div on:renderComplete>
  <!-- Your component content -->
</div>
```

### API Call Monitoring

```typescript
import { getPerformanceMonitor } from '$lib/shared/performance-monitoring';

const monitor = getPerformanceMonitor();

async function fetchTasks() {
  const startTime = performance.now();
  try {
    const response = await fetch('/api/tasks');
    const data = await response.json();
    monitor.trackApiCall('/api/tasks', startTime);
    return data;
  } catch (error) {
    monitor.trackApiCall('/api/tasks', startTime);
    throw error;
  }
}
```

## Backend Performance Monitoring

### Implementation

The backend uses a custom performance monitoring system built with Go.

#### Setup

```go
// In main.go
import "taskmanager/metrics"

func main() {
    // Initialize performance monitor
    metrics.InitPerformanceMonitor()
    
    // ... rest of application setup
    
    // Add metrics endpoint
    mux.HandleFunc("/metrics", GlobalMonitor.MetricsHandler())
    mux.HandleFunc("/health", GlobalMonitor.HealthHandler())
}
```

#### Middleware Integration

```go
// Wrap handlers with performance monitoring
mux.HandleFunc("/api/tasks", GlobalMonitor.Middleware("/api/tasks", taskController.GetTasks))
mux.HandleFunc("/api/tasks/create", GlobalMonitor.Middleware("/api/tasks/create", taskController.CreateTask))
```

#### Metrics Endpoint

```bash
# Get current metrics
curl http://localhost:8080/metrics

# Health check
curl http://localhost:8080/health
```

### Performance Data Structure

```json
{
  "api_metrics": {
    "call_count": {
      "/api/tasks": 150,
      "/api/tasks/create": 25
    },
    "avg_duration": {
      "/api/tasks": 45.5,
      "/api/tasks/create": 120.3
    },
    "p50_duration": {
      "/api/tasks": 40.0,
      "/api/tasks/create": 110.0
    },
    "p95_duration": {
      "/api/tasks": 80.0,
      "/api/tasks/create": 200.0
    },
    "p99_duration": {
      "/api/tasks": 150.0,
      "/api/tasks/create": 350.0
    },
    "total_requests": 175
  },
  "system_metrics": {
    "goroutines": 45,
    "memory_alloc": 8192000,
    "memory_total_alloc": 24576000,
    "memory_sys": 16384000,
    "memory_heap_alloc": 8192000,
    "memory_heap_sys": 12288000,
    "gc_pause_total": 2500000,
    "gc_pause_count": 12
  },
  "error_metrics": {
    "error_count": 5,
    "error_rate": 2.86
  },
  "timing": {
    "start_time": "2026-09-11T10:00:00Z",
    "last_updated": "2026-09-11T12:00:00Z",
    "uptime": "2h0m0s"
  }
}
```

### Custom Metrics Recording

```go
// Record API call with duration
startTime := time.Now()
// ... API logic
duration := time.Since(startTime)
GlobalMonitor.RecordAPICall("/api/custom", duration)

// Record error
if err != nil {
    GlobalMonitor.RecordError()
}

// Update system metrics
GlobalMonitor.UpdateSystemMetrics()
```

## Integration with Monitoring Services

### Prometheus Integration

#### Backend

```go
import "github.com/prometheus/client_golang/prometheus"

// Create custom metrics
var (
    apiCallsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "api_calls_total",
            Help: "Total number of API calls",
        },
        []string{"endpoint", "method"},
    )
    
    apiDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "api_duration_seconds",
            Help:    "API call duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"endpoint", "method"},
    )
)

func init() {
    prometheus.MustRegister(apiCallsTotal)
    prometheus.MustRegister(apiDuration)
}

// Record metrics
apiCallsTotal.WithLabelValues("/api/tasks", "GET").Inc()
apiDuration.WithLabelValues("/api/tasks", "GET").Observe(duration.Seconds())
```

#### Frontend

```typescript
// Send metrics to Prometheus Pushgateway
async function sendToPrometheus(metrics: PerformanceData) {
  const response = await fetch('https://your-prometheus-pushgateway/metrics/job/taskmanager', {
    method: 'POST',
    body: formatMetricsForPrometheus(metrics),
  });
}
```

### Grafana Dashboards

Create Grafana dashboards to visualize metrics:

#### Backend Dashboard Panels

1. **Request Rate**: Requests per second by endpoint
2. **Response Time**: P50, P95, P99 latency
3. **Error Rate**: Error percentage over time
4. **Memory Usage**: Heap allocation and GC stats
5. **Goroutine Count**: Active goroutines
6. **Database Performance**: Query duration

#### Frontend Dashboard Panels

1. **Web Vitals**: FCP, LCP, FID, CLS over time
2. **API Performance**: API response times
3. **Page Load Time**: Page load duration
4. **User Engagement**: Session duration, bounce rate
5. **Error Rate**: JavaScript errors

### Sentry Integration

#### Frontend

```typescript
import * as Sentry from "@sentry/svelte";

Sentry.init({
  dsn: import.meta.env.PUBLIC_SENTRY_DSN,
  environment: import.meta.env.MODE,
  tracesSampleRate: 1.0,
  integrations: [
    new Sentry.BrowserTracing(),
  ],
});

// Track performance
Sentry.startTransaction({
  op: 'navigation',
  name: 'Task List Page',
});
```

#### Backend

```go
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.Options{
  Dsn: "your-sentry-dsn",
  Environment: "production",
  TracesSampleRate: 1.0,
})

// Capture performance
hub := sentry.CurrentHub()
transaction := sentry.StartTransaction(hub, "api-call", "http")
transaction.SetData("endpoint", "/api/tasks")
transaction.Finish()
```

## Alerting and Notifications

### Alert Rules

#### Performance Alerts

```yaml
# Prometheus alert rules
groups:
  - name: performance_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(api_errors_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          
      - alert: SlowResponseTime
        expr: histogram_quantile(0.95, api_duration_seconds) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "95th percentile response time too high"
          
      - alert: HighMemoryUsage
        expr: memory_heap_alloc_bytes / memory_heap_sys_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Memory usage approaching limit"
```

#### Web Vitals Alerts

```yaml
  - alert: PoorLCP
    expr: web_vitals_lcp_seconds > 2.5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Largerst Contentful Paint exceeds threshold"
      
  - alert: HighCLS
    expr: web_vitals_cls > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Cumulative Layout Shift too high"
```

### Notification Channels

#### Slack Notifications

```yaml
# Alertmanager configuration
receivers:
  - name: 'slack-notifications'
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/WEBHOOK/URL'
        channel: '#performance-alerts'
        send_resolved: true
```

#### Email Notifications

```yaml
receivers:
  - name: 'email-notifications'
    email_configs:
      - to: 'team@example.com'
        send_resolved: true
```

## Performance Optimization

### Frontend Optimization

#### 1. Code Splitting

```typescript
// Lazy load components
const HeavyComponent = lazy(() => import('./HeavyComponent.svelte'));
```

#### 2. Image Optimization

```typescript
// Use WebP format with fallbacks
<picture>
  <source srcset="image.webp" type="image/webp">
  <img src="image.jpg" alt="Description">
</picture>
```

#### 3. Caching Strategy

```typescript
// Implement service worker caching
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});
```

### Backend Optimization

#### 1. Database Query Optimization

```go
// Use indexes and optimize queries
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_priority ON tasks(priority);
```

#### 2. Connection Pooling

```go
// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

#### 3. Caching

```go
// Implement Redis caching
func (s *TaskService) GetTasksWithCache(userID string) ([]Task, error) {
    cacheKey := fmt.Sprintf("tasks:%s", userID)
    
    // Try cache first
    if cached, err := redis.Get(cacheKey); err == nil {
        return decodeTasks(cached)
    }
    
    // Fallback to database
    tasks, err := s.taskModel.GetByUserID(userID)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    redis.Set(cacheKey, encodeTasks(tasks), 5*time.Minute)
    
    return tasks, nil
}
```

## Best Practices

### 1. Monitoring Strategy

- **Start Early**: Implement monitoring from the beginning
- **Monitor Key Metrics**: Focus on business-critical metrics
- **Set Baselines**: Establish normal performance baselines
- **Review Regularly**: Review metrics weekly/monthly

### 2. Alert Management

- **Meaningful Alerts**: Only alert on actionable issues
- **Avoid Alert Fatigue**: Don't over-alert
- **Escalation Rules**: Define clear escalation paths
- **Documentation**: Document alert procedures

### 3. Data Retention

- **Retain Historical Data**: Keep data for trend analysis
- **Aggregate Data**: Aggregate old data to save space
- **Backup Data**: Regular backups of monitoring data
- **Compliance**: Consider data retention policies

### 4. Performance Targets

#### Web Vitals Targets

- **FCP**: < 1.8 seconds
- **LCP**: < 2.5 seconds
- **FID**: < 100 milliseconds
- **CLS**: < 0.1
- **TTFB**: < 600 milliseconds

#### API Performance Targets

- **P50 Response Time**: < 200ms
- **P95 Response Time**: < 500ms
- **P99 Response Time**: < 1000ms
- **Error Rate**: < 1%
- **Availability**: > 99.9%

### 5. Testing Performance

#### Load Testing

```bash
# Use k6 for load testing
k6 run --vus 100 --duration 5m load-test.js
```

#### Performance Testing

```bash
# Lighthouse CI
npm install -g @lhci/cli
lhci autorun
```

## Troubleshooting

### Common Issues

#### 1. High Memory Usage

```bash
# Check memory metrics
curl http://localhost:8080/metrics

# Check for memory leaks
# Profile with pprof
go tool pprof http://localhost:8080/debug/pprof/heap
```

#### 2. Slow API Response

```bash
# Check API duration metrics
curl http://localhost:8080/metrics

# Profile CPU usage
go tool pprof http://localhost:8080/debug/pprof/profile
```

#### 3. High Error Rate

```bash
# Check error metrics
curl http://localhost:8080/metrics

# Review logs
journalctl -u taskmanager -f
```

### Debug Mode

Enable detailed logging for debugging:

```bash
export LOG_LEVEL=debug
export LOG_FORMAT=text
```

## Additional Resources

- [Web Vitals Documentation](https://web.dev/vitals/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Sentry Documentation](https://docs.sentry.io/)
- [Build Optimization Guide](./BUILD_OPTIMIZATION.md) - Build optimizations
- [Deployment Guide](./DEPLOYMENT_GUIDE.md) - Deployment monitoring
- [CI/CD Guide](./CICD_GUIDE.md) - Pipeline monitoring
