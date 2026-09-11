# Environment Configuration Guide

This guide provides comprehensive documentation for configuring environment variables and settings for both the backend and frontend components of the Task Manager application.

## Table of Contents

- [Overview](#overview)
- [Backend Configuration](#backend-configuration)
- [Frontend Configuration](#frontend-configuration)
- [Secrets Management](#secrets-management)
- [Environment-Specific Configurations](#environment-specific-configurations)
- [Security Best Practices](#security-best-practices)
- [Validation and Testing](#validation-and-testing)

## Overview

The Task Manager application uses a layered configuration approach:

1. **Configuration Files**: YAML files for default settings
2. **Environment Variables**: Override file settings for environment-specific values
3. **Secrets Management**: Secure storage for sensitive data

### Configuration Priority

Environment variables override configuration file values, providing flexibility for different deployment environments.

## Backend Configuration

### Configuration File Structure

The backend uses `env/config.yaml` for configuration:

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password_here"
  dbname: "taskmanager"
  sslmode: "disable"

jwt:
  secret: "your_jwt_secret_key_here_change_in_production"

cors:
  allowed_origins:
    - "http://localhost:5173"
    - "http://localhost:3000"
    - "http://localhost:8080"

logging:
  level: "info"
  format: "json"

app:
  env: "development"
  rate_limit:
    enabled: true
    requests_per_minute: 100
    requests_per_hour: 1000
  session:
    timeout_hours: 24
```

### Environment Variables

#### Server Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | Server port | 8080 | No |
| `APP_ENV` | Application environment | development | No |

#### Database Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | localhost | Yes |
| `DB_PORT` | Database port | 5432 | No |
| `DB_USER` | Database user | postgres | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | taskmanager | Yes |
| `DB_SSLMODE` | SSL mode | disable | No |

#### JWT Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | JWT secret key | - | Yes |

#### CORS Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins | - | Yes |

#### Logging Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info | No |
| `LOG_FORMAT` | Log format (json, text) | json | No |

#### Rate Limiting Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `RATE_LIMIT_ENABLED` | Enable rate limiting | true | No |
| `RATE_LIMIT_RPM` | Requests per minute | 100 | No |
| `RATE_LIMIT_RPH` | Requests per hour | 1000 | No |

#### Session Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SESSION_TIMEOUT_HOURS` | Session timeout in hours | 24 | No |

### Environment-Specific Configuration Files

#### Development (`env/config.yaml`)

```yaml
server:
  port: "8080"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "taskmanager"
  sslmode: "disable"

jwt:
  secret: "dev-secret-key-change-in-production"

cors:
  allowed_origins:
    - "http://localhost:5173"
    - "http://localhost:3000"
    - "http://localhost:8080"

logging:
  level: "debug"
  format: "text"

app:
  env: "development"
  rate_limit:
    enabled: false
  session:
    timeout_hours: 168  # 1 week for development
```

#### Production (`env/config.production.yaml`)

```yaml
server:
  port: "8080"

database:
  host: "your-production-db-host"
  port: 5432
  user: "your-production-db-user"
  password: "your-secure-db-password"
  dbname: "taskmanager"
  sslmode: "require"

jwt:
  secret: "your-very-secure-jwt-secret-key-minimum-32-characters"

cors:
  allowed_origins:
    - "https://your-frontend-domain.com"
    - "https://www.your-frontend-domain.com"

logging:
  level: "info"
  format: "json"

app:
  env: "production"
  rate_limit:
    enabled: true
    requests_per_minute: 100
    requests_per_hour: 1000
  session:
    timeout_hours: 24
```

## Frontend Configuration

### Environment Variables

Frontend environment variables must start with `PUBLIC_` to be exposed to the client.

#### Required Variables

| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `PUBLIC_API_URL` | Backend API URL | https://api.example.com | Yes |

#### Optional Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PUBLIC_APP_NAME` | Application name | Task Manager | My App |
| `PUBLIC_APP_VERSION` | Application version | 1.0.0 | 2.0.0 |
| `PUBLIC_ENABLE_ANALYTICS` | Enable analytics feature | true | false |
| `PUBLIC_ENABLE_EXPORT` | Enable export feature | true | false |
| `PUBLIC_ENABLE_BULK_OPERATIONS` | Enable bulk operations | true | false |
| `PUBLIC_SENTRY_DSN` | Sentry DSN for error tracking | - | https://... |
| `PUBLIC_GOOGLE_ANALYTICS_ID` | Google Analytics ID | - | G-XXXXXXXXXX |
| `PUBLIC_DEFAULT_PAGE_SIZE` | Default pagination size | 20 | 50 |
| `PUBLIC_MAX_FILE_SIZE` | Maximum file size in bytes | 10485760 | 52428800 |
| `PUBLIC_SUPPORTED_IMAGE_TYPES` | Supported image types | image/jpeg,... | image/png,... |
| `PUBLIC_DEFAULT_TIMEZONE` | Default timezone | UTC | America/New_York |
| `PUBLIC_ENABLE_TIMEZONE_DETECTION` | Enable timezone detection | true | false |

### Environment-Specific Configuration Files

#### Development (`.env`)

```env
PUBLIC_API_URL=http://localhost:8080
PUBLIC_APP_NAME=Task Manager
PUBLIC_APP_VERSION=1.0.0
PUBLIC_ENABLE_ANALYTICS=true
PUBLIC_ENABLE_EXPORT=true
PUBLIC_ENABLE_BULK_OPERATIONS=true
PUBLIC_DEFAULT_PAGE_SIZE=20
PUBLIC_DEFAULT_TIMEZONE=UTC
PUBLIC_ENABLE_TIMEZONE_DETECTION=true
```

#### Production (`.env.production`)

```env
PUBLIC_API_URL=https://api.yourdomain.com
PUBLIC_APP_NAME=Task Manager
PUBLIC_APP_VERSION=1.0.0
PUBLIC_ENABLE_ANALYTICS=true
PUBLIC_ENABLE_EXPORT=true
PUBLIC_ENABLE_BULK_OPERATIONS=true
PUBLIC_SENTRY_DSN=https://your-sentry-dsn
PUBLIC_GOOGLE_ANALYTICS_ID=G-XXXXXXXXXX
PUBLIC_DEFAULT_PAGE_SIZE=20
PUBLIC_DEFAULT_TIMEZONE=UTC
PUBLIC_ENABLE_TIMEZONE_DETECTION=true
```

## Secrets Management

### Why Secrets Management?

Never commit secrets to version control. Use secure secret management systems for:

- Database passwords
- JWT secrets
- API keys
- Third-party service credentials

### Secrets Management Options

#### 1. Environment Variables

For simple deployments:

```bash
export DB_PASSWORD=your_secure_password
export JWT_SECRET=your_jwt_secret
```

#### 2. Cloud Provider Secret Managers

##### AWS Secrets Manager

```bash
# Store secret
aws secretsmanager create-secret \
  --name taskmanager/db-password \
  --secret-string "your_secure_password"

# Retrieve secret in application
aws secretsmanager get-secret-value \
  --secret-id taskmanager/db-password
```

##### Google Secret Manager

```bash
# Store secret
gcloud secrets create db-password --data-file=- <<< "your_secure_password"

# Retrieve secret
gcloud secrets versions access latest --secret="db-password"
```

##### Azure Key Vault

```bash
# Store secret
az keyvault secret set \
  --vault-name your-keyvault \
  --name db-password \
  --value "your_secure_password"

# Retrieve secret
az keyvault secret show \
  --vault-name your-keyvault \
  --name db-password
```

#### 3. Docker Secrets

```yaml
# docker-compose.yml
version: '3.8'
services:
  backend:
    secrets:
      - db_password
      - jwt_secret

secrets:
  db_password:
    file: ./secrets/db_password.txt
  jwt_secret:
    file: ./secrets/jwt_secret.txt
```

#### 4. Kubernetes Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: taskmanager-secrets
type: Opaque
stringData:
  db-password: your_secure_password
  jwt-secret: your_jwt_secret
```

### Generating Secure Secrets

#### JWT Secret

```bash
# Generate 32-byte random secret
openssl rand -base64 32

# Or use
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

#### Database Password

```bash
# Generate strong password
openssl rand -base64 24

# Or use password manager
# 1Password, LastPass, Bitwarden, etc.
```

## Environment-Specific Configurations

### Development Environment

**Purpose**: Local development and testing

**Characteristics**:
- Debug logging enabled
- No rate limiting
- Relaxed CORS settings
- Local database
- Longer session timeouts

**Configuration**:
- Use `env/config.yaml`
- Use `.env` for frontend
- Development tools enabled

### Staging Environment

**Purpose**: Pre-production testing

**Characteristics**:
- Production-like configuration
- Test data only
- Monitoring enabled
- Rate limiting enabled
- Separate database

**Configuration**:
- Use `env/config.staging.yaml`
- Use `.env.staging` for frontend
- Monitoring and error tracking enabled

### Production Environment

**Purpose**: Live production deployment

**Characteristics**:
- Maximum security
- Rate limiting enabled
- Strict CORS settings
- Production database
- Optimized logging
- Error tracking enabled

**Configuration**:
- Use `env/config.production.yaml`
- Use `.env.production` for frontend
- All security features enabled
- Secrets from secret manager

## Security Best Practices

### 1. Never Commit Secrets

```bash
# Add to .gitignore
.env
.env.local
.env.*.local
config.production.yaml
secrets/
```

### 2. Use Strong Secrets

- Minimum 32 characters for JWT secrets
- Use cryptographically secure random generators
- Rotate secrets regularly
- Never reuse secrets across environments

### 3. Principle of Least Privilege

- Database users with minimum required permissions
- API keys with specific scopes
- Limited access to secrets

### 4. Environment Isolation

- Separate configurations for each environment
- Never share secrets between environments
- Use environment-specific databases

### 5. Audit and Monitor

- Log access to secrets
- Monitor for unauthorized access
- Regular security audits
- Review access permissions

### 6. Secure Storage

- Encrypt secrets at rest
- Use TLS for secret transmission
- Implement secret rotation policies
- Use hardware security modules (HSM) for critical secrets

## Validation and Testing

### Configuration Validation

#### Backend Validation

```go
// In config/config.go
func ValidateConfig(cfg *Config) error {
    if cfg.Database.Host == "" {
        return errors.New("DB_HOST is required")
    }
    if cfg.Database.Password == "" {
        return errors.New("DB_PASSWORD is required")
    }
    if len(cfg.JWT.Secret) < 32 {
        return errors.New("JWT_SECRET must be at least 32 characters")
    }
    return nil
}
```

#### Frontend Validation

```typescript
// In src/lib/env.ts
import { z } from 'zod';

const envSchema = z.object({
  PUBLIC_API_URL: z.string().url(),
  PUBLIC_APP_NAME: z.string().optional(),
  PUBLIC_APP_VERSION: z.string().optional(),
});

export const env = envSchema.parse(import.meta.env);
```

### Testing Configuration

#### Unit Tests

```go
func TestConfigValidation(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: Config{
                Database: DatabaseConfig{
                    Host:     "localhost",
                    Password: "password",
                },
                JWT: JWTConfig{
                    Secret: "very-secure-secret-key-32-chars-long",
                },
            },
            wantErr: false,
        },
        {
            name: "missing db host",
            config: Config{
                Database: DatabaseConfig{
                    Host:     "",
                    Password: "password",
                },
            },
            wantErr: true,
        },
    }
    // ... test implementation
}
```

#### Integration Tests

```bash
# Test with different configurations
APP_ENV=test go test ./...
APP_ENV=production go test ./...
```

### Configuration Drift Detection

```bash
# Compare configurations
diff env/config.yaml env/config.production.yaml

# Or use configuration management tools
# Terraform, Ansible, etc.
```

## Troubleshooting

### Common Issues

#### 1. Configuration Not Loading

```bash
# Check file permissions
ls -la env/config.yaml

# Check file path
pwd
ls -la env/

# Verify YAML syntax
python3 -c "import yaml; yaml.safe_load(open('env/config.yaml'))"
```

#### 2. Environment Variables Not Working

```bash
# Check if variable is set
echo $DB_HOST

# Check process environment
ps aux | grep taskmanager
cat /proc/<pid>/environ

# Test in Go
fmt.Println(os.Getenv("DB_HOST"))
```

#### 3. Invalid Configuration

```bash
# Validate YAML
yamllint env/config.yaml

# Check for typos
grep -r "DB_HOST" env/

# Test configuration loading
go run main.go --config-test
```

### Debug Mode

Enable debug logging to troubleshoot configuration issues:

```bash
export LOG_LEVEL=debug
go run main.go
```

## Additional Resources

- [Deployment Guide](./DEPLOYMENT_GUIDE.md) - Deployment instructions
- [Build Optimization Guide](./BUILD_OPTIMIZATION.md) - Build configurations
- [AGENTS.md](./backendGoVanilaTaskmanager/AGENTS.md) - Backend guidelines
- [AGENTS_FrontEnd.md](./sveltekit-taskmanager-frontend/AGENTS_FrontEnd.md) - Frontend guidelines
- [API Integration Guide](./API_INTEGRATION_GUIDE.md) - API documentation
