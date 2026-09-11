# CI/CD Pipeline Guide

This guide provides comprehensive documentation for the CI/CD pipeline setup, workflows, and best practices for the Task Manager application.

## Table of Contents

- [Overview](#overview)
- [Pipeline Architecture](#pipeline-architecture)
- [GitHub Actions Workflows](#github-actions-workflows)
- [Pipeline Stages](#pipeline-stages)
- [Secrets Management](#secrets-management)
- [Deployment Strategies](#deployment-strategies)
- [Monitoring and Notifications](#monitoring-and-notifications)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

The CI/CD pipeline automates the build, test, and deployment process for both the backend (Go) and frontend (SvelteKit) components.

### Key Features

- **Automated Testing**: Unit tests, integration tests, and E2E tests
- **Code Quality**: Linting, formatting checks, and security scanning
- **Build Optimization**: Docker image building with caching
- **Automated Deployment**: Push to production on main branch
- **Monitoring**: Health checks and deployment notifications

### Pipeline Triggers

- **Push to main/develop**: Full pipeline execution
- **Pull requests**: CI checks only
- **Path filters**: Backend/frontend specific triggers

## Pipeline Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    GitHub Actions Trigger                  │
└──────────────────────┬────────────────────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        │                             │
┌───────▼────────┐           ┌────────▼────────┐
│  Backend CI    │           │  Frontend CI    │
└───────┬────────┘           └────────┬────────┘
        │                             │
        └──────────────┬──────────────┘
                       │
              ┌────────▼────────┐
              │  Build & Push   │
              │  Docker Images  │
              └────────┬────────┘
                       │
              ┌────────▼────────┐
              │    Deploy       │
              │  to Production  │
              └────────┬────────┘
                       │
              ┌────────▼────────┐
              │  Health Checks  │
              │  & Notifications│
              └─────────────────┘
```

## GitHub Actions Workflows

### 1. Backend CI Workflow

**File**: `.github/workflows/backend-ci.yml`

**Triggers**:
- Push to main/develop branches (backend changes only)
- Pull requests to main/develop (backend changes only)

**Stages**:
1. **Setup**: Checkout code, set up Go environment
2. **Dependencies**: Install Go modules with caching
3. **Code Quality**: 
   - Format check with `gofmt`
   - Security scan with `govulncheck`
4. **Testing**:
   - Unit tests (`go test`)
   - Race detection (`go test -race`)
   - Coverage report (`go test -cover`)
5. **Build**: Compile Go binary
6. **Docker**: Build Docker image

**Key Features**:
- Go version caching for faster builds
- Path-based triggers for efficiency
- Coverage upload to Codecov
- Security vulnerability scanning

### 2. Frontend CI Workflow

**File**: `.github/workflows/frontend-ci.yml`

**Triggers**:
- Push to main/develop branches (frontend changes only)
- Pull requests to main/develop (frontend changes only)

**Stages**:
1. **Setup**: Checkout code, set up Bun environment
2. **Dependencies**: Install dependencies with caching
3. **Code Quality**:
   - Type check (`bun run check`)
   - Linting (`bun run lint`)
   - Format check (`bun run format:check`)
4. **Testing**:
   - Unit tests (`bun test`)
   - Coverage report (`bun test:coverage`)
   - E2E tests with Playwright
5. **Build**: Production build (`bun run build`)
6. **Docker**: Build Docker image

**Key Features**:
- Bun for fast dependency management
- Playwright for E2E testing
- Coverage upload to Codecov
- Test artifact retention

### 3. Main CI/CD Pipeline

**File**: `.github/workflows/ci.yml`

**Triggers**:
- Push to main/develop branches (all changes)
- Pull requests to main/develop (all changes)

**Stages**:
1. **Parallel CI**: Run backend and frontend CI in parallel
2. **Build & Push**: Build and push Docker images
3. **Deploy**: Deploy to production
4. **Health Checks**: Verify deployment health
5. **Notifications**: Send deployment status

**Key Features**:
- Reusable workflows for CI
- Conditional deployment (main branch only)
- Docker Buildx for advanced builds
- Slack notifications for deployment status

## Pipeline Stages

### Stage 1: Continuous Integration

#### Backend CI

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Set up Go 1.22
      - Install dependencies
      - Check formatting
      - Run tests
      - Run race detection
      - Generate coverage
      - Security scan
      - Build binary
      - Build Docker image
```

#### Frontend CI

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Set up Bun
      - Install dependencies
      - Type check
      - Lint
      - Format check
      - Run unit tests
      - Generate coverage
      - Run E2E tests
      - Build application
      - Build Docker image
```

### Stage 2: Build and Push

```yaml
jobs:
  build-and-push:
    needs: [backend-ci, frontend-ci]
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Set up Docker Buildx
      - Login to Docker Hub
      - Build and push backend image
      - Build and push frontend image
```

**Docker Image Tags**:
- `latest`: Latest stable version
- `{commit-sha}`: Specific commit version
- `buildcache`: Build cache for faster builds

### Stage 3: Deployment

```yaml
jobs:
  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    steps:
      - Checkout code
      - Deploy to production
      - Run database migrations
      - Health checks
      - Send notifications
```

## Secrets Management

### Required GitHub Secrets

Configure these secrets in your GitHub repository settings:

#### Docker Hub
- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub password or access token

#### Deployment
- `DEPLOY_HOST`: Production server hostname
- `DEPLOY_USER`: SSH username for deployment
- `DEPLOY_KEY`: SSH private key for deployment
- `DEPLOY_PATH`: Deployment path on server

#### Monitoring
- `SLACK_WEBHOOK`: Slack webhook URL for notifications
- `SENTRY_DSN`: Sentry DSN for error tracking (optional)

#### Database
- `DB_HOST`: Production database host
- `DB_PASSWORD`: Production database password
- `JWT_SECRET`: JWT secret for authentication

### Using Secrets in Workflows

```yaml
- name: Deploy to production
  run: |
    ssh -i ${{ secrets.DEPLOY_KEY }} ${{ secrets.DEPLOY_USER }}@${{ secrets.DEPLOY_HOST }} \
      'cd ${{ secrets.DEPLOY_PATH }} && docker-compose pull && docker-compose up -d'
```

### Security Best Practices

1. **Never commit secrets** to version control
2. **Use GitHub Secrets** for sensitive data
3. **Rotate secrets regularly**
4. **Use least privilege** access
5. **Audit secret access** regularly

## Deployment Strategies

### 1. Blue-Green Deployment

```yaml
- name: Blue-Green Deployment
  run: |
    # Deploy to green environment
    kubectl apply -f k8s/green/
    
    # Wait for health checks
    kubectl wait --for=condition=ready pod -l app=taskmanager-green
    
    # Switch traffic
    kubectl patch service taskmanager -p '{"spec":{"selector":{"version":"green"}}}'
    
    # Cleanup blue environment
    kubectl delete -f k8s/blue/
```

### 2. Rolling Update

```yaml
- name: Rolling Update
  run: |
    # Update deployment with rolling strategy
    kubectl set image deployment/taskmanager \
      backend=${{ secrets.DOCKER_USERNAME }}/taskmanager-backend:${{ github.sha }}
    
    # Wait for rollout to complete
    kubectl rollout status deployment/taskmanager
```

### 3. Canary Deployment

```yaml
- name: Canary Deployment
  run: |
    # Deploy canary version
    kubectl apply -f k8s/canary/
    
    # Route 10% traffic to canary
    kubectl patch service taskmanager -p '{"spec":{"selector":{"version":"canary"}}}'
    
    # Monitor metrics (add monitoring step)
    
    # Gradual rollout or rollback based on metrics
```

### 4. Docker Compose Deployment

```yaml
- name: Docker Compose Deployment
  run: |
    ssh -i ${{ secrets.DEPLOY_KEY }} ${{ secrets.DEPLOY_USER }}@${{ secrets.DEPLOY_HOST }} \
      'cd ${{ secrets.DEPLOY_PATH }} && \
       docker-compose pull && \
       docker-compose up -d && \
       docker-compose ps'
```

## Monitoring and Notifications

### Health Checks

```yaml
- name: Health Check
  run: |
    # Backend health check
    curl -f https://api.yourdomain.com/health || exit 1
    
    # Frontend health check
    curl -f https://yourdomain.com/health || exit 1
    
    # Database health check
    PGPASSWORD=${{ secrets.DB_PASSWORD }} psql -h ${{ secrets.DB_HOST }} -U postgres -c "SELECT 1"
```

### Slack Notifications

```yaml
- name: Notify Deployment
  if: always()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: 'Deployment to production completed'
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

### Custom Notifications

```yaml
- name: Custom Notification
  run: |
    curl -X POST ${{ secrets.NOTIFICATION_WEBHOOK }} \
      -H 'Content-Type: application/json' \
      -d '{
        "status": "${{ job.status }}",
        "commit": "${{ github.sha }}",
        "branch": "${{ github.ref }}",
        "author": "${{ github.actor }}"
      }'
```

## Best Practices

### 1. Pipeline Optimization

#### Caching

```yaml
- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

#### Parallel Execution

```yaml
jobs:
  backend-ci:
    # ... backend CI steps
  frontend-ci:
    # ... frontend CI steps
  # These run in parallel
```

#### Matrix Builds

```yaml
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go-version: ['1.21', '1.22']
```

### 2. Quality Gates

#### Required Checks

```yaml
- name: Quality Gate
  run: |
    if [ ${{ steps.test.outcome }} != 'success' ]; then
      echo "Tests failed - blocking deployment"
      exit 1
    fi
```

#### Branch Protection

Configure branch protection rules in GitHub:
- Require status checks to pass
- Require branches to be up to date
- Require review from code owners

### 3. Security

#### Dependency Scanning

```yaml
- name: Dependency Scan
  uses: actions/dependency-review-action@v3
```

#### CodeQL Analysis

```yaml
- name: Initialize CodeQL
  uses: github/codeql-action/init@v2
  with:
    languages: go, javascript
```

#### Secret Scanning

```yaml
- name: Secret Scan
  uses: trufflesecurity/trufflehog@main
  with:
    path: ./
    base: ${{ github.event.repository.default_branch }}
```

### 4. Performance

#### Build Time Optimization

- Use build caching
- Parallelize independent jobs
- Use faster runners (GitHub-hosted vs self-hosted)
- Optimize Docker layer caching

#### Resource Optimization

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 30
    concurrency:
      group: ${{ github.workflow }}-${{ github.ref }}
      cancel-in-progress: true
```

### 5. Documentation

#### Pipeline Documentation

- Document each workflow
- Include trigger conditions
- Explain required secrets
- Provide troubleshooting steps

#### Change Logs

```yaml
- name: Generate Change Log
  run: |
    # Generate changelog from commits
    # Add to release notes
```

## Troubleshooting

### Common Issues

#### 1. Pipeline Fails on Tests

```bash
# Run tests locally
cd backendGoVanilaTaskmanager
go test ./...

cd ../sveltekit-taskmanager-frontend
bun test
```

#### 2. Docker Build Fails

```bash
# Test Docker build locally
cd backendGoVanilaTaskmanager
docker build -t test .

cd ../sveltekit-taskmanager-frontend
docker build -t test .
```

#### 3. Deployment Fails

```bash
# Check deployment logs
ssh user@server 'docker-compose logs taskmanager'

# Check service status
ssh user@server 'systemctl status taskmanager'
```

#### 4. Secrets Not Available

```yaml
# Debug secret availability
- name: Debug Secrets
  run: |
    echo "DOCKER_USERNAME available: ${{ secrets.DOCKER_USERNAME != '' }}"
    echo "DEPLOY_HOST available: ${{ secrets.DEPLOY_HOST != '' }}"
```

### Debug Mode

Enable debug logging in GitHub Actions:

```yaml
- name: Debug Step
  run: |
    echo "::debug::Detailed debug information"
    echo "::warning::Warning message"
    echo "::error::Error message"
```

### Retry Logic

```yaml
- name: Deploy with Retry
  uses: nick-fields/retry@v2
  with:
    timeout_minutes: 10
    max_attempts: 3
    command: |
      ssh -i ${{ secrets.DEPLOY_KEY }} ${{ secrets.DEPLOY_USER }}@${{ secrets.DEPLOY_HOST }} \
        'cd ${{ secrets.DEPLOY_PATH }} && docker-compose up -d'
```

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Codecov Documentation](https://docs.codecov.com/)
- [Slack Action Documentation](https://github.com/8398a7/action-slack)
- [Deployment Guide](./DEPLOYMENT_GUIDE.md) - Deployment instructions
- [Build Optimization Guide](./BUILD_OPTIMIZATION.md) - Build configurations
- [Environment Configuration Guide](./ENVIRONMENT_CONFIGURATION.md) - Environment setup
