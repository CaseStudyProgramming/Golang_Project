# Build Optimization Guide

This guide covers build optimization strategies for both the backend (Go) and frontend (SvelteKit) components of the Task Manager application.

## Table of Contents

- [Backend Build Optimization](#backend-build-optimization)
- [Frontend Build Optimization](#frontend-build-optimization)
- [Docker Optimization](#docker-optimization)
- [CI/CD Build Optimization](#cicd-build-optimization)
- [Performance Monitoring](#performance-monitoring)
- [Best Practices](#best-practices)

## Backend Build Optimization

### Go Build Optimizations

#### 1. Use Build Flags

The Makefile includes optimized build flags:

```makefile
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"
```

- `-s`: Remove symbol table
- `-w`: Remove DWARF debug information
- Reduces binary size by ~30-40%

#### 2. Cross-Platform Builds

Build for multiple platforms:

```bash
make build-linux    # Linux AMD64
make build-windows  # Windows AMD64
make build-darwin   # macOS AMD64
make build-all      # All platforms
```

#### 3. Dependency Optimization

```bash
# Vendor dependencies for reproducible builds
go mod vendor

# Use specific Go version for consistency
GOVERSION=1.22 go build
```

#### 4. Build Caching

```bash
# Use Go build cache
export GOCACHE=/path/to/cache
go build

# Clean cache if needed
go clean -cache
```

### Performance Optimizations

#### 1. Database Connection Pooling

Configure optimal connection pool settings:

```go
// In config/database.go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

#### 2. Enable Go Profiling

```go
// Add to main.go
import _ "net/http/pprof"

// Then access /debug/pprof/ for profiling
```

#### 3. Use Efficient Data Structures

- Use slices instead of maps for small datasets
- Pre-allocate slice capacity when size is known
- Use sync.Pool for object reuse

### Memory Optimization

#### 1. Reduce Memory Footprint

```go
// Use smaller data types
int8, int16 instead of int32, int64 when possible
float32 instead of float64 when precision allows

// Use value types instead of pointers for small structs
```

#### 2. Garbage Collection Tuning

```bash
# Set GOGC environment variable
export GOGC=100  # Default
export GOGC=50   # More aggressive GC
```

## Frontend Build Optimization

### Vite Configuration

The optimized `vite.config.ts` includes:

#### 1. Code Splitting

```typescript
rollupOptions: {
  output: {
    manualChunks: {
      'chart-vendor': ['chart.js', 'svelte-chartjs'],
      'zod': ['zod']
    }
  }
}
```

Benefits:
- Smaller initial bundle
- Better caching
- Parallel loading

#### 2. Tree Shaking

```typescript
minify: 'terser',
terserOptions: {
  compress: {
    drop_console: true,    // Remove console.log
    drop_debugger: true    // Remove debugger statements
  }
}
```

#### 3. Source Maps

```typescript
sourcemap: true  // Enable for production debugging
```

### Bundle Size Optimization

#### 1. Analyze Bundle Size

```bash
# Install bundle analyzer
bun add -D rollup-plugin-visualizer

# Add to vite.config.ts
import { visualizer } from 'rollup-plugin-visualizer';

plugins: [
  // ... other plugins
  visualizer({
    filename: './stats.html',
    open: true,
    gzipSize: true
  })
]
```

#### 2. Lazy Loading

```svelte
// Lazy load components
const HeavyComponent = lazy(() => import('./HeavyComponent.svelte'));

// Use in template
{#await HeavyComponent then Component}
  <Component />
{/await}
```

#### 3. Dynamic Imports

```typescript
// Load libraries only when needed
const loadChart = () => import('chart.js');

async function showChart() {
  const Chart = await loadChart();
  // Use Chart
}
```

### Asset Optimization

#### 1. Image Optimization

```bash
# Install image optimization tools
bun add -D vite-plugin-imagemin sharp

# Configure in vite.config.ts
import imagemin from 'vite-plugin-imagemin';

plugins: [
  imagemin({
    gifsicle: { optimizationLevel: 7 },
    optipng: { optimizationLevel: 7 },
    mozjpeg: { quality: 80 },
    pngquant: { quality: [0.8, 0.9] }
  })
]
```

#### 2. Font Optimization

```css
/* Use font-display: swap for faster rendering */
@font-face {
  font-family: 'Custom Font';
  src: url('./font.woff2') format('woff2');
  font-display: swap;
}
```

#### 3. CSS Optimization

```typescript
// Purge unused CSS
import { default as vitePurgeCss } from '@modular-css/vite';

plugins: [
  vitePurgeCss({
    content: ['./src/**/*.svelte', './src/**/*.ts']
  })
]
```

### Runtime Optimizations

#### 1. Code Splitting by Route

```typescript
// SvelteKit automatically code-splits by route
// Additional optimization: preload critical routes
```

#### 2. Prefetching

```svelte
<!-- Prefetch next page -->
<svelte:component 
  this={nextPage} 
  data-sveltekit-preload-data="hover"
/>
```

#### 3. Caching Strategy

```typescript
// Service worker for offline caching
// Use workbox for efficient caching
```

## Docker Optimization

### Multi-Stage Builds

Both backend and frontend use multi-stage builds:

#### Backend Dockerfile

```dockerfile
# Builder stage
FROM golang:1.22-alpine AS builder
# ... build steps ...

# Runtime stage
FROM alpine:latest
# ... minimal runtime ...
```

Benefits:
- Smaller final image
- No build tools in runtime
- Faster deployment

#### Frontend Dockerfile

```dockerfile
# Builder stage
FROM node:18-alpine AS builder
# ... build steps ...

# Runtime stage
FROM nginx:alpine
# ... serve static files ...
```

### Image Size Optimization

#### 1. Use Alpine Base

```dockerfile
FROM alpine:latest  # ~5MB
# Instead of
FROM ubuntu:latest  # ~70MB
```

#### 2. Minimize Layers

```dockerfile
# Combine RUN commands
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 1000 taskmanager && \
    adduser -D -u 1000 -G taskmanager taskmanager
```

#### 3. Clean Up

```dockerfile
# Remove unnecessary files
RUN rm -rf /var/cache/apk/* /tmp/*
```

### Build Cache Optimization

#### 1. Layer Caching

```dockerfile
# Copy dependencies first for better caching
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Then copy source code
COPY . .
```

#### 2. BuildKit Cache

```bash
# Use BuildKit for better caching
DOCKER_BUILDKIT=1 docker build .
```

#### 3. Registry Caching

```bash
# Use build cache from registry
docker build --cache-from=taskmanager-backend:latest .
```

## CI/CD Build Optimization

### Parallel Builds

```yaml
# GitHub Actions example
jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - run: make build
```

### Dependency Caching

```yaml
# Cache Go modules
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

# Cache node modules
- uses: actions/cache@v3
  with:
    path: node_modules
    key: ${{ runner.os }}-node-${{ hashFiles('**/bun.lock') }}
```

### Incremental Builds

```yaml
# Only build changed packages
- run: |
    if git diff --name-only HEAD~1 HEAD | grep -q "backend"; then
      cd backendGoVanilaTaskmanager && make build
    fi
```

### Artifact Caching

```yaml
# Cache build artifacts
- uses: actions/cache@v3
  with:
    path: |
      backendGoVanilaTaskmanager/taskmanager-server
      sveltekit-taskmanager-frontend/build
    key: build-${{ github.sha }}
```

## Performance Monitoring

### Build Time Monitoring

#### 1. Measure Build Times

```bash
# Backend build time
time make build

# Frontend build time
time bun run build
```

#### 2. Track Build Metrics

```yaml
# Add to CI/CD pipeline
- name: Report build metrics
  run: |
    echo "Backend build: ${BACKEND_BUILD_TIME}s"
    echo "Frontend build: ${FRONTEND_BUILD_TIME}s"
```

### Runtime Performance

#### 1. Backend Metrics

```go
// Add metrics endpoint
func metricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := map[string]interface{}{
        "memory": getMemoryStats(),
        "goroutines": runtime.NumGoroutine(),
        "requests": getRequestCount(),
    }
    json.NewEncoder(w).Encode(metrics)
}
```

#### 2. Frontend Performance

```typescript
// Measure page load time
window.addEventListener('load', () => {
  const perfData = performance.getEntriesByType('navigation')[0];
  console.log('Page load time:', perfData.loadEventEnd - perfData.fetchStart);
});
```

### Bundle Size Monitoring

#### 1. Track Bundle Size

```bash
# Add to package.json
"scripts": {
  "size-check": "bundlesize"
}

# .bundlesize.json
{
  "files": [
    {
      "path": "./build/index.js",
      "maxSize": "200 kB"
    }
  ]
}
```

#### 2. Lighthouse CI

```yaml
# Add to CI/CD pipeline
- name: Run Lighthouse CI
  run: |
    npm install -g @lhci/cli
    lhci autorun
```

## Best Practices

### 1. Build Strategy

- **Development**: Fast builds, source maps, hot reload
- **Staging**: Optimized builds, source maps, monitoring
- **Production**: Highly optimized, no source maps, minified

### 2. Version Management

- Use semantic versioning
- Tag releases
- Maintain changelog

### 3. Dependency Management

- Regular dependency updates
- Security audits
- License compliance

### 4. Testing Integration

- Run tests before builds
- Measure test coverage
- Integrate with CI/CD

### 5. Documentation

- Document build process
- Maintain build scripts
- Update deployment guides

### 6. Monitoring

- Monitor build times
- Track bundle sizes
- Measure runtime performance

### 7. Security

- Scan dependencies for vulnerabilities
- Use non-root containers
- Implement security headers

### 8. Scalability

- Optimize for scale
- Implement caching
- Use CDNs

## Optimization Checklist

### Backend
- [ ] Use build flags for size reduction
- [ ] Implement connection pooling
- [ ] Enable profiling for debugging
- [ ] Optimize database queries
- [ ] Use efficient data structures
- [ ] Implement caching strategy
- [ ] Monitor memory usage

### Frontend
- [ ] Enable code splitting
- [ ] Implement tree shaking
- [ ] Optimize bundle size
- [ ] Lazy load components
- [ ] Optimize images and assets
- [ ] Implement caching
- [ ] Use CDN for static assets

### Docker
- [ ] Use multi-stage builds
- [ ] Minimize image size
- [ ] Implement layer caching
- [ ] Use non-root users
- [ ] Implement health checks
- [ ] Optimize Dockerfile

### CI/CD
- [ ] Implement parallel builds
- [ ] Cache dependencies
- [ ] Use incremental builds
- [ ] Monitor build times
- [ ] Implement quality gates
- [ ] Automate deployments

## Troubleshooting

### Build Issues

#### Backend Build Fails

```bash
# Check Go version
go version

# Clean build cache
go clean -cache

# Verify dependencies
go mod verify
```

#### Frontend Build Fails

```bash
# Clear cache
rm -rf node_modules .svelte-kit build
bun install

# Check TypeScript errors
bun run check
```

### Performance Issues

#### Slow Build Times

```bash
# Check for cache misses
# Optimize Docker layer order
# Use BuildKit
```

#### Large Bundle Size

```bash
# Analyze bundle
bun run build:analyze

# Check for unused dependencies
# Implement code splitting
# Remove unused code
```

## Additional Resources

- [Go Build Optimization](https://golang.org/doc/gc_guide)
- [Vite Build Optimization](https://vitejs.dev/guide/build.html)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [SvelteKit Performance](https://kit.svelte.dev/docs/performance)
- [AGENTS.md](./backendGoVanilaTaskmanager/AGENTS.md) - Backend guidelines
- [AGENTS_FrontEnd.md](./sveltekit-taskmanager-frontend/AGENTS_FrontEnd.md) - Frontend guidelines
