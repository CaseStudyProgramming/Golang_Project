# Deployment Guide

This guide provides comprehensive instructions for deploying the Task Manager application to production environments, including backend, frontend, database, and infrastructure setup.

## Table of Contents

- [Deployment Overview](#deployment-overview)
- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Database Deployment](#database-deployment)
- [Backend Deployment](#backend-deployment)
- [Frontend Deployment](#frontend-deployment)
- [Infrastructure Options](#infrastructure-options)
- [Monitoring and Logging](#monitoring-and-logging)
- [Security Considerations](#security-considerations)
- [Backup and Recovery](#backup-and-recovery)
- [Scaling Considerations](#scaling-considerations)
- [Troubleshooting](#troubleshooting)

## Deployment Overview

The Task Manager application consists of:

- **Backend**: Go API server
- **Frontend**: SvelteKit application
- **Database**: PostgreSQL
- **Infrastructure**: Can be deployed to various cloud providers

### Architecture

```
┌─────────────┐
│   CDN/WAF   │
└──────┬──────┘
       │
┌──────▼──────┐
│  Frontend   │ (SvelteKit)
│  (Static)   │
└──────┬──────┘
       │
┌──────▼──────┐
│   Backend   │ (Go API)
│   Server    │
└──────┬──────┘
       │
┌──────▼──────┐
│  PostgreSQL │
│  Database   │
└─────────────┘
```

## Prerequisites

### Required Tools

- **Go**: 1.22 or higher
- **Node.js**: 18+ or Bun
- **PostgreSQL**: 15+
- **Domain name** (for production)
- **SSL certificate** (for HTTPS)
- **Cloud provider account** (AWS, GCP, Azure, DigitalOcean, etc.)

### System Requirements

- **Backend**: Minimum 1 CPU, 2GB RAM
- **Frontend**: Static hosting (CDN)
- **Database**: Minimum 2 CPU, 4GB RAM, 20GB storage
- **Recommended**: 2 CPU, 4GB RAM for backend, 4 CPU, 8GB RAM for database

## Environment Configuration

### Backend Environment Variables

Create a production environment file for the backend:

```bash
# backendGoVanilaTaskmanager/.env.production
APP_ENV=production
SERVER_PORT=8080

# Database Configuration
DB_HOST=your-production-db-host
DB_PORT=5432
DB_USER=your-production-db-user
DB_PASSWORD=your-secure-db-password
DB_NAME=taskmanager
DB_SSLMODE=require

# JWT Configuration
JWT_SECRET=your-very-secure-jwt-secret-key-min-32-chars

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com,https://www.your-frontend-domain.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Frontend Environment Variables

Create a production environment file for the frontend:

```bash
# sveltekit-taskmanager-frontend/.env.production
PUBLIC_API_URL=https://your-api-domain.com
```

### Security Best Practices

1. **Never commit secrets** to version control
2. **Use strong passwords** (minimum 32 characters for JWT secret)
3. **Enable SSL** for database connections
4. **Use environment-specific** configurations
5. **Rotate secrets** regularly
6. **Use secret management** services (AWS Secrets Manager, etc.)

## Database Deployment

### Option 1: Managed PostgreSQL Service

#### AWS RDS

```bash
# Create RDS instance via AWS CLI
aws rds create-db-instance \
  --db-instance-identifier taskmanager-db \
  --db-instance-class db.t3.medium \
  --engine postgres \
  --engine-version 15.4 \
  --master-username taskmanager \
  --master-user-password your-secure-password \
  --allocated-storage 20 \
  --storage-type gp2 \
  --publicly-accessible false \
  --vpc-security-group-ids sg-xxxxxxxx \
  --db-subnet-group-name default-vpc-xxxxxxxx
```

#### DigitalOcean Managed Database

```bash
# Create via doctl
doctl databases create taskmanager-db \
  --engine pg \
  --version 15 \
  --num-nodes 1 \
  --size db-s-2vcpu-4gb \
  --region nyc1
```

#### Google Cloud SQL

```bash
# Create Cloud SQL instance
gcloud sql instances create taskmanager-db \
  --database-version POSTGRES_15 \
  --tier db-n1-standard-2 \
  --region us-central1 \
  --storage-size 20GB
```

### Option 2: Self-Hosted PostgreSQL

#### Using Docker

```bash
# Create docker-compose.prod.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    container_name: taskmanager-postgres
    environment:
      POSTGRES_USER: taskmanager
      POSTGRES_PASSWORD: your-secure-password
      POSTGRES_DB: taskmanager
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - taskmanager-network

volumes:
  postgres_data:

networks:
  taskmanager-network:
    driver: bridge
```

```bash
# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

### Database Setup

1. **Run migrations**:
```bash
cd backendGoVanilaTaskmanager

# Set production environment variables
export DB_HOST=your-db-host
export DB_USER=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=taskmanager

# Run migrations
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/001.init.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/002.add_deleted_at.up.sql
# ... run all migrations in order
```

2. **Create database user** (if not exists):
```sql
CREATE USER taskmanager WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;
```

3. **Verify connection**:
```bash
psql -h your-db-host -U taskmanager -d taskmanager -c "SELECT version();"
```

## Backend Deployment

### Option 1: Traditional VPS Deployment

#### 1. Build the Application

```bash
cd backendGoVanilaTaskmanager

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o taskmanager-server main.go

# Or build for specific target
GOOS=linux GOARCH=amd64 go build -o taskmanager-server main.go
```

#### 2. Create Systemd Service

Create `/etc/systemd/system/taskmanager.service`:

```ini
[Unit]
Description=Task Manager Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=taskmanager
WorkingDirectory=/opt/taskmanager
ExecStart=/opt/taskmanager/taskmanager-server
Restart=always
RestartSec=10
Environment="APP_ENV=production"
Environment="DB_HOST=localhost"
Environment="DB_PORT=5432"
Environment="DB_USER=taskmanager"
Environment="DB_PASSWORD=your-secure-password"
Environment="DB_NAME=taskmanager"
Environment="DB_SSLMODE=require"
Environment="JWT_SECRET=your-jwt-secret"
Environment="CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com"

[Install]
WantedBy=multi-user.target
```

#### 3. Deploy Files

```bash
# Create deployment directory
sudo mkdir -p /opt/taskmanager
sudo chown taskmanager:taskmanager /opt/taskmanager

# Copy binary
sudo cp taskmanager-server /opt/taskmanager/
sudo chmod +x /opt/taskmanager/taskmanager-server

# Copy migrations
sudo cp -r migrations /opt/taskmanager/
sudo chown -R taskmanager:taskmanager /opt/taskmanager
```

#### 4. Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable taskmanager

# Start service
sudo systemctl start taskmanager

# Check status
sudo systemctl status taskmanager

# View logs
sudo journalctl -u taskmanager -f
```

### Option 2: Docker Deployment

#### 1. Create Dockerfile

Create `backendGoVanilaTaskmanager/Dockerfile`:

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o taskmanager-server main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/taskmanager-server .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run application
CMD ["./taskmanager-server"]
```

#### 2. Create docker-compose.prod.yml

```yaml
version: '3.8'
services:
  backend:
    build: .
    container_name: taskmanager-backend
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=taskmanager
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=taskmanager
      - DB_SSLMODE=disable
      - JWT_SECRET=${JWT_SECRET}
      - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS}
    depends_on:
      - postgres
    restart: unless-stopped
    networks:
      - taskmanager-network

  postgres:
    image: postgres:15
    container_name: taskmanager-postgres
    environment:
      - POSTGRES_USER=taskmanager
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=taskmanager
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - taskmanager-network

volumes:
  postgres_data:

networks:
  taskmanager-network:
    driver: bridge
```

#### 3. Deploy

```bash
# Create .env file
cat > .env << EOF
DB_PASSWORD=your-secure-password
JWT_SECRET=your-jwt-secret
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com
EOF

# Build and start
docker-compose -f docker-compose.prod.yml up -d --build

# Run migrations
docker-compose -f docker-compose.prod.yml exec backend psql -h postgres -U taskmanager -d taskmanager -f migrations/001.init.up.sql
```

### Option 3: Cloud Platform Deployment

#### AWS EC2 + ECS

```bash
# Build and push Docker image
docker build -t taskmanager-backend .
docker tag taskmanager-backend:latest your-ecr-repo/taskmanager-backend:latest
docker push your-ecr-repo/taskmanager-backend:latest

# Deploy to ECS
aws ecs update-service --cluster taskmanager-cluster --service taskmanager-service --force-new-deployment
```

#### Google Cloud Run

```bash
# Build and push to Google Container Registry
gcloud builds submit --tag gcr.io/your-project/taskmanager-backend

# Deploy to Cloud Run
gcloud run deploy taskmanager-backend \
  --image gcr.io/your-project/taskmanager-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

#### Heroku

```bash
# Create Procfile
echo "web: ./taskmanager-server" > Procfile

# Deploy
git push heroku main
```

### Backend Health Check

Add health check endpoint to verify deployment:

```bash
# Test health endpoint
curl https://your-api-domain.com/health

# Expected response
{"status":"healthy","timestamp":"2026-09-11T10:00:00Z"}
```

## Frontend Deployment

### Option 1: Static Hosting (Vercel, Netlify)

#### Vercel Deployment

```bash
cd sveltekit-taskmanager-frontend

# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod

# Set environment variables
vercel env add PUBLIC_API_URL production
```

#### Netlify Deployment

```bash
# Install Netlify CLI
npm i -g netlify-cli

# Build
bun run build

# Deploy
netlify deploy --prod --dir=build
```

### Option 2: Cloud Storage + CDN

#### AWS S3 + CloudFront

```bash
# Build application
cd sveltekit-taskmanager-frontend
bun run build

# Upload to S3
aws s3 sync build/ s3://your-bucket-name --delete

# Configure CloudFront distribution
aws cloudfront create-distribution \
  --origin-domain-name your-bucket-name.s3.amazonaws.com \
  --default-root-object index.html \
  --default-cache-behavior TargetOriginId=your-origin-id,ViewerProtocolPolicy=redirect-to-https
```

#### Google Cloud Storage + Cloud CDN

```bash
# Build application
bun run build

# Upload to GCS
gsutil rsync -R build/ gs://your-bucket-name

# Configure Cloud CDN
gcloud compute backend-buckets create taskmanager-frontend \
  --gcs-bucket-name=your-bucket-name \
  --enable-cdn
```

### Option 3: Traditional Web Server

#### Nginx Configuration

```nginx
server {
    listen 80;
    server_name your-frontend-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-frontend-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    root /var/www/taskmanager-frontend;
    index index.html;

    # SvelteKit static files
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

#### Deploy Files

```bash
# Build application
cd sveltekit-taskmanager-frontend
bun run build

# Copy to web server
sudo cp -r build/* /var/www/taskmanager-frontend/

# Set permissions
sudo chown -R www-data:www-data /var/www/taskmanager-frontend
sudo chmod -R 755 /var/www/taskmanager-frontend
```

### Frontend Build Optimization

Ensure production build is optimized:

```bash
# Build with analysis
bun run build

# Check bundle size
# Consider code splitting and lazy loading
```

## Infrastructure Options

### AWS Architecture

```
                    ┌─────────────┐
                    │   Route53   │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  CloudFront │
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              │                         │
        ┌─────▼─────┐           ┌─────▼─────┐
        │   S3      │           │   ALB     │
        │ (Frontend)│           └─────┬─────┘
        └───────────┘                 │
                              ┌──────▼──────┐
                              │  ECS Fargate │
                              │  (Backend)   │
                              └──────┬──────┘
                                     │
                              ┌──────▼──────┐
                              │  RDS Postgres│
                              └─────────────┘
```

### Google Cloud Architecture

```
                    ┌─────────────┐
                    │  Cloud DNS  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Cloud CDN   │
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              │                         │
        ┌─────▼─────┐           ┌─────▼─────┐
        │  GCS      │           │ Cloud Run │
        │(Frontend) │           │ (Backend) │
        └───────────┘           └─────┬─────┘
                                     │
                              ┌──────▼──────┐
                              │ Cloud SQL   │
                              │  Postgres   │
                              └─────────────┘
```

### DigitalOcean Architecture

```
                    ┌─────────────┐
                    │  Cloud DNS  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   Spaces    │
                    │ (Frontend)  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Load Balancer│
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              │                         │
        ┌─────▼─────┐           ┌─────▼─────┐
        │  Droplet  │           │ Managed DB │
        │ (Backend) │           │  Postgres  │
        └───────────┘           └────────────┘
```

## Monitoring and Logging

### Application Monitoring

#### Prometheus + Grafana

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'taskmanager'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

#### Health Checks

```bash
# Setup health check monitoring
curl -f https://your-api-domain.com/health || alert
```

### Logging

#### Structured Logging

Backend uses structured JSON logging:

```json
{
  "level": "info",
  "timestamp": "2026-09-11T10:00:00Z",
  "message": "Task created successfully",
  "user_id": "1",
  "task_id": "1",
  "ip_address": "192.168.1.1"
}
```

#### Log Aggregation

- **CloudWatch Logs** (AWS)
- **Cloud Logging** (GCP)
- **Papertrail**
- **Loggly**
- **Datadog**

### Error Tracking

#### Sentry Integration

```go
// Add to backend
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.Options{
    Dsn: "your-sentry-dsn",
    Environment: "production",
})
```

## Security Considerations

### SSL/TLS Configuration

1. **Use HTTPS** for all endpoints
2. **Enable HSTS** headers
3. **Use strong SSL ciphers**
4. **Keep certificates updated**

### Firewall Configuration

```bash
# UFW example
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

### Security Headers

Add security headers to backend:

```go
// Add middleware
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
```

### Database Security

1. **Use strong passwords**
2. **Enable SSL connections**
3. **Restrict network access**
4. **Regular security updates**
5. **Backup encryption**

### API Security

1. **Rate limiting**
2. **Input validation**
3. **SQL injection prevention**
4. **XSS protection**
5. **CSRF protection**

## Backup and Recovery

### Database Backups

#### Automated Backups

```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -h your-db-host -U taskmanager taskmanager > $BACKUP_DIR/taskmanager_$DATE.sql

# Keep last 7 days
find $BACKUP_DIR -name "taskmanager_*.sql" -mtime +7 -delete
```

#### Cloud Provider Backups

- **AWS RDS**: Automated backups enabled
- **Google Cloud SQL**: Automated backups enabled
- **DigitalOcean**: Automated backups enabled

### Disaster Recovery

1. **Document recovery procedures**
2. **Test restore process regularly**
3. **Keep offsite backups**
4. **Have standby infrastructure ready**

## Scaling Considerations

### Horizontal Scaling

#### Load Balancing

```nginx
# Nginx load balancer configuration
upstream backend {
    server backend1.example.com:8080;
    server backend2.example.com:8080;
    server backend3.example.com:8080;
}

server {
    location /api/ {
        proxy_pass http://backend;
    }
}
```

#### Database Scaling

- **Read replicas** for read-heavy workloads
- **Connection pooling** (PgBouncer)
- **Database sharding** for large datasets

### Caching Strategy

#### Redis Cache

```bash
# Add Redis for caching
docker run -d -p 6379:6379 redis:alpine
```

Cache frequently accessed data:
- User sessions
- Task lists
- Category lists
- Analytics data

### CDN Configuration

- **Static assets** served from CDN
- **API responses** cached where appropriate
- **Geographic distribution** for better performance

## Troubleshooting

### Common Issues

#### Backend Won't Start

```bash
# Check logs
sudo journalctl -u taskmanager -n 50

# Check port availability
sudo netstat -tlnp | grep 8080

# Check database connection
psql -h your-db-host -U taskmanager -d taskmanager -c "SELECT 1;"
```

#### Database Connection Issues

```bash
# Test connection
psql -h your-db-host -U taskmanager -d taskmanager

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

#### Frontend Build Issues

```bash
# Clear cache
rm -rf node_modules .svelte-kit
bun install

# Check environment variables
echo $PUBLIC_API_URL
```

### Performance Issues

#### Database Performance

```bash
# Check slow queries
SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;

# Add indexes if needed
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_priority ON tasks(priority);
```

#### Backend Performance

```bash
# Profile with pprof
go tool pprof http://localhost:8080/debug/pprof/profile
```

### Monitoring Setup

```bash
# Setup basic monitoring
# CPU, memory, disk usage monitoring
# Response time monitoring
# Error rate monitoring
```

## Post-Deployment Checklist

- [ ] Verify database connectivity
- [ ] Test all API endpoints
- [ ] Verify frontend functionality
- [ ] Test authentication flow
- [ ] Verify SSL certificate
- [ ] Setup monitoring and alerting
- [ ] Configure backups
- [ ] Test disaster recovery
- [ ] Security audit
- [ ] Performance testing
- [ ] Load testing
- [ ] Documentation update

## Additional Resources

- [README.md](./README.md) - Project overview
- [API Integration Guide](./API_INTEGRATION_GUIDE.md) - API documentation
- [Component Documentation](./COMPONENT_DOCUMENTATION.md) - Frontend components
- [AGENTS.md](./backendGoVanilaTaskmanager/AGENTS.md) - Backend guidelines
- [AGENTS_FrontEnd.md](./sveltekit-taskmanager-frontend/AGENTS_FrontEnd.md) - Frontend guidelines
