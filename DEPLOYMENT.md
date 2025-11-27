# Production Deployment Guide

## Overview

This guide covers deploying Shurl API to production environments.

## Pre-Deployment Checklist

- [ ] Environment variables configured for production
- [ ] Database credentials secured and rotated
- [ ] TLS/SSL certificates obtained
- [ ] Application tested in staging environment
- [ ] Database backups configured
- [ ] Monitoring and logging setup
- [ ] Security audit completed
- [ ] Load balancing configured (if applicable)

## Environment Configuration

### Production `.env` File

```dotenv
# ==========================================
# Database Configuration
# ==========================================
POSTGRES_HOST=your-production-db-host
POSTGRES_USER=prod_user
POSTGRES_PASSWORD=strong_secure_password_here
POSTGRES_DB=shurl_production
POSTGRES_PORT=5432
POSTGRES_SSLMODE=require

# ==========================================
# Application Configuration
# ==========================================
PORT=8080
SECRET_KEY=your-very-long-secret-key-minimum-32-chars

# ==========================================
# Environment
# ==========================================
GIN_MODE=release

# ==========================================
# CORS Configuration
# ==========================================
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### Critical Production Settings

1. **GIN_MODE=release**: Disables debug logging and improves performance
2. **POSTGRES_SSLMODE=require**: Enforces SSL for database connections
3. **SECRET_KEY**: Use a strong random key (minimum 32 characters)
4. **TLS/HTTPS**: Configure reverse proxy with SSL certificates

## Docker Deployment

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o Shurl .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/Shurl .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./Shurl"]
```

### Docker Compose Example

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:16-alpine
    container_name: shurl_postgres
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - shurl_network

  shurl-api:
    build: .
    container_name: shurl_api
    ports:
      - "8080:8080"
    environment:
      POSTGRES_HOST: postgres
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_PORT: 5432
      PORT: 8080
      SECRET_KEY: ${SECRET_KEY}
      GIN_MODE: release
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--quiet",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 30s
      timeout: 3s
      retries: 3
    networks:
      - shurl_network
    restart: unless-stopped

volumes:
  postgres_data:

networks:
  shurl_network:
```

### Deploy with Docker Compose

```bash
# Set environment variables
export POSTGRES_USER=prod_user
export POSTGRES_PASSWORD=secure_password
export POSTGRES_DB=shurl_production
export SECRET_KEY=your-secret-key-here

# Build and start services
docker-compose -f docker-compose.yml up -d

# View logs
docker-compose logs -f shurl-api

# Stop services
docker-compose down
```

## Kubernetes Deployment

### Kubernetes Manifest Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: shurl-config
  namespace: default
data:
  GIN_MODE: "release"
  POSTGRES_SSLMODE: "require"
  POSTGRES_PORT: "5432"

---
apiVersion: v1
kind: Secret
metadata:
  name: shurl-secret
  namespace: default
type: Opaque
stringData:
  POSTGRES_HOST: your-db-host
  POSTGRES_USER: prod_user
  POSTGRES_PASSWORD: secure_password
  POSTGRES_DB: shurl_production
  SECRET_KEY: your-secret-key
  PORT: "8080"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shurl-api
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: shurl-api
  template:
    metadata:
      labels:
        app: shurl-api
    spec:
      containers:
        - name: shurl-api
          image: your-registry/shurl:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
              name: http
          envFrom:
            - configMapRef:
                name: shurl-config
            - secretRef:
                name: shurl-secret
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi

---
apiVersion: v1
kind: Service
metadata:
  name: shurl-service
  namespace: default
spec:
  type: LoadBalancer
  selector:
    app: shurl-api
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
      name: http
```

### Deploy to Kubernetes

```bash
# Create namespace
kubectl create namespace shurl

# Apply manifests
kubectl apply -f k8s/manifest.yaml

# Check deployment status
kubectl get deployments -n shurl
kubectl get pods -n shurl

# View logs
kubectl logs -n shurl -f deployment/shurl-api

# Port forward for testing
kubectl port-forward -n shurl service/shurl-service 8080:80
```

## Nginx Reverse Proxy Setup

### Nginx Configuration

```nginx
upstream shurl_backend {
    least_conn;
    server backend1.local:8080 max_fails=3 fail_timeout=30s;
    server backend2.local:8080 max_fails=3 fail_timeout=30s;
    server backend3.local:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Logging
    access_log /var/log/nginx/shurl_access.log;
    error_log /var/log/nginx/shurl_error.log;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # Proxy settings
    location / {
        proxy_pass http://shurl_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Swagger UI endpoint
    location /swagger/ {
        proxy_pass http://shurl_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## Database Setup

### PostgreSQL Configuration

```sql
-- Create production database
CREATE DATABASE shurl_production;

-- Create application user
CREATE USER shurl_user WITH PASSWORD 'secure_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE shurl_production TO shurl_user;
GRANT ALL PRIVILEGES ON SCHEMA public TO shurl_user;

-- Connect to database
\c shurl_production

-- Run migrations
-- This will be done by the application on startup
```

### Backup Strategy

```bash
#!/bin/bash
# Automated daily backup script

BACKUP_DIR="/backups/shurl"
DB_NAME="shurl_production"
DB_USER="shurl_user"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create backup
pg_dump -U $DB_USER -h localhost -d $DB_NAME | gzip > "$BACKUP_DIR/shurl_$TIMESTAMP.sql.gz"

# Keep only last 30 days of backups
find $BACKUP_DIR -name "shurl_*.sql.gz" -mtime +30 -delete

# Upload to S3 (optional)
# aws s3 cp "$BACKUP_DIR/shurl_$TIMESTAMP.sql.gz" s3://your-bucket/backups/
```

## Monitoring & Logging

### Prometheus Metrics

Add to your Go application:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Add metrics endpoint
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Log Aggregation Setup

Configure JSON logging for ELK/CloudWatch:

```bash
# View logs in production
docker logs shurl_api | tail -f

# Filter for errors
docker logs shurl_api | grep ERROR

# Monitor specific endpoints
docker logs shurl_api | grep "POST /links"
```

## Security Hardening

### Essential Security Measures

1. **Database**

   - Use strong passwords (minimum 16 characters)
   - Enable SSL/TLS for database connections
   - Use private subnets for databases
   - Enable audit logging

2. **Application**

   - Set `GIN_MODE=release`
   - Enable HTTPS with valid SSL certificate
   - Implement rate limiting
   - Add request logging

3. **Infrastructure**

   - Use Web Application Firewall (WAF)
   - Enable DDoS protection
   - Implement firewall rules
   - Regular security updates

4. **Secrets Management**
   - Use managed secrets service (AWS Secrets Manager, HashiCorp Vault)
   - Never hardcode secrets in code
   - Rotate secrets regularly
   - Implement secret versioning

## Scaling

### Horizontal Scaling

- Deploy multiple instances behind load balancer
- Use connection pooling for database
- Implement caching layer (Redis)
- Use CDN for static content

### Vertical Scaling

- Increase server resources (CPU, RAM)
- Optimize database indexes
- Enable query caching
- Tune Go runtime parameters

## Health Checks & Monitoring

### Application Health

```bash
# Check API status
curl https://yourdomain.com/health

# Expected response
{
  "status": "ok",
  "message": "API is running"
}
```

### Database Health

```bash
# Check database connection
psql -h host -U user -d database -c "SELECT 1"
```

## Disaster Recovery

### RTO and RPO Targets

- **RTO** (Recovery Time Objective): 1 hour
- **RPO** (Recovery Point Objective): 15 minutes

### Recovery Procedure

1. Restore latest database backup
2. Deploy application from latest stable commit
3. Verify application health checks
4. Monitor logs for errors
5. Perform smoke tests

## Maintenance

### Regular Tasks

- **Daily**: Monitor logs and metrics
- **Weekly**: Review security alerts, backup verification
- **Monthly**: Update dependencies, security patches
- **Quarterly**: Performance review, capacity planning

## Cost Optimization

- Use reserved instances for databases
- Implement auto-scaling based on load
- Use managed database services
- Optimize container image size
- Monitor resource utilization

## Further Reading

- [Go Production Deployment](https://golang.org/doc/tutorial/web-service-gin)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [Kubernetes Security Best Practices](https://kubernetes.io/docs/concepts/security/)
- [PostgreSQL Administration](https://www.postgresql.org/docs/current/admin.html)
