# Deployment Guide - Content Review API

## Production Deployment Checklist

- [ ] Environment configuration
- [ ] Database setup and migrations
- [ ] SSL/TLS certificates
- [ ] Reverse proxy configuration
- [ ] Monitoring and logging
- [ ] Backup strategy
- [ ] Security hardening
- [ ] Performance testing
- [ ] Health checks
- [ ] CI/CD pipeline

---

## Prerequisites

- Docker and Docker Compose (recommended)
- PostgreSQL 15 (if not using Docker)
- Go 1.21+ (if building from source)
- nginx or equivalent reverse proxy
- SSL certificate (Let's Encrypt)

---

## Option 1: Docker Deployment (Recommended)

### Step 1: Prepare Environment File

```bash
# Create production .env file
cat > .env.production << 'EOF'
# Application
ENVIRONMENT=production
APP_NAME=content-review-api
DEBUG=false

# Database (change these!)
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=generate_strong_password_here
DB_NAME=content_review
DB_SSLMODE=require
DB_MAX_CONN=50

# JWT (generate strong secrets)
JWT_SECRET=generate_very_strong_random_secret_64_chars_minimum
JWT_REFRESH_SECRET=generate_very_strong_random_secret_64_chars_minimum
JWT_EXPIRY_HOURS=24

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
EOF
```

### Step 2: Generate Secure Secrets

```bash
# Generate JWT secrets (Linux/Mac)
openssl rand -base64 64

# Or use Python
python3 -c "import secrets; print(secrets.token_urlsafe(64))"
```

### Step 3: Update docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: content_review_postgres
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - content-review-network

  api:
    build: .
    container_name: content_review_api
    restart: always
    env_file: .env.production
    expose:
      - "8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - content-review-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  nginx:
    image: nginx:alpine
    container_name: content_review_nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - api
    networks:
      - content-review-network

volumes:
  postgres_data:
    driver: local

networks:
  content-review-network:
    driver: bridge
```

### Step 4: Nginx Configuration

```bash
# Create nginx.conf
cat > nginx.conf << 'EOF'
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 10M;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript 
               application/json application/javascript application/xml+rss;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/s;
    limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=5r/m;

    # Upstream API
    upstream api_backend {
        server api:8080;
    }

    # HTTP to HTTPS redirect
    server {
        listen 80;
        server_name _;
        return 301 https://$host$request_uri;
    }

    # HTTPS server
    server {
        listen 443 ssl http2;
        server_name api.example.com;

        # SSL configuration
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers on;
        ssl_session_cache shared:SSL:10m;
        ssl_session_timeout 10m;

        # Security headers
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;

        # Health check endpoint
        location /health {
            access_log off;
            proxy_pass http://api_backend;
        }

        # API endpoints
        location /api/ {
            limit_req zone=api_limit burst=20 nodelay;
            
            proxy_pass http://api_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_read_timeout 90;
        }

        # Auth endpoints (stricter rate limiting)
        location /api/v1/auth/ {
            limit_req zone=auth_limit burst=3 nodelay;
            
            proxy_pass http://api_backend;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # 404 for unknown paths
        location / {
            return 404;
        }
    }
}
EOF
```

### Step 5: SSL Certificate Setup

```bash
# Using Let's Encrypt with Certbot
docker run --rm -it \
  -v /etc/letsencrypt:/etc/letsencrypt \
  -v /var/lib/letsencrypt:/var/lib/letsencrypt \
  -p 80:80 \
  certbot/certbot certonly --standalone \
  -d api.example.com

# Copy certificates
mkdir -p ssl
cp /etc/letsencrypt/live/api.example.com/fullchain.pem ssl/cert.pem
cp /etc/letsencrypt/live/api.example.com/privkey.pem ssl/key.pem
```

### Step 6: Deploy

```bash
# Start services
docker-compose -f docker-compose.yml up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f api

# Test API
curl https://api.example.com/health
```

---

## Option 2: Manual Linux Deployment

### Step 1: Install Dependencies

```bash
# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install PostgreSQL
sudo apt-get install postgresql postgresql-contrib -y

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Nginx
sudo apt-get install nginx -y

# Install supervisor for process management
sudo apt-get install supervisor -y
```

### Step 2: Setup Database

```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database and user
CREATE DATABASE content_review;
CREATE USER api_user WITH PASSWORD 'strong_password';
ALTER ROLE api_user SET client_encoding TO 'utf8';
ALTER ROLE api_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE api_user SET default_transaction_deferrable TO on;
ALTER ROLE api_user SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE content_review TO api_user;
\q

# Run migrations
psql -h localhost -U api_user -d content_review -f migrations/001_initial_schema.sql
```

### Step 3: Deploy Application

```bash
# Clone repository
cd /opt
sudo git clone https://github.com/yourusername/content-review-api.git
cd content-review-api

# Setup environment
sudo cp .env.example .env.production
sudo nano .env.production  # Edit with production values

# Build application
export PATH=$PATH:/usr/local/go/bin
go build -o bin/app main.go

# Test run
./bin/app
```

### Step 4: Setup Supervisor

```bash
# Create supervisor config
sudo tee /etc/supervisor/conf.d/content-review-api.conf > /dev/null << EOF
[program:content-review-api]
directory=/opt/content-review-api
command=/opt/content-review-api/bin/app
autostart=true
autorestart=true
stderr_logfile=/var/log/content-review-api/err.log
stdout_logfile=/var/log/content-review-api/out.log
environment=PATH="/usr/local/go/bin:/usr/bin",ENVIRONMENT="production"
EOF

# Create log directory
sudo mkdir -p /var/log/content-review-api
sudo chown nobody:nogroup /var/log/content-review-api

# Start supervisor
sudo systemctl restart supervisor
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start content-review-api
```

### Step 5: Configure Nginx

```bash
# Create Nginx config
sudo tee /etc/nginx/sites-available/content-review-api > /dev/null << 'EOF'
server {
    listen 80;
    server_name api.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/content-review-api /etc/nginx/sites-enabled/

# Test nginx
sudo nginx -t

# Restart nginx
sudo systemctl restart nginx
```

---

## Monitoring & Logging

### Setup Logging with ELK Stack

```bash
# Docker Compose with ELK
docker-compose -f docker-compose.monitoring.yml up -d
```

### Health Checks

```bash
# Manual health check
curl https://api.example.com/health

# Automated monitoring script
#!/bin/bash
while true; do
  response=$(curl -s -o /dev/null -w "%{http_code}" https://api.example.com/health)
  if [ $response != "200" ]; then
    echo "API is down! Status: $response"
    # Send alert (email, Slack, PagerDuty)
  fi
  sleep 30
done
```

---

## Backup Strategy

### Database Backups

```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/backups/database"
DATE=$(date +%Y%m%d_%H%M%S)

# Full backup
docker-compose exec -T postgres pg_dump -U postgres content_review > $BACKUP_DIR/backup_$DATE.sql

# Compress
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

# Upload to S3 (optional)
aws s3 cp $BACKUP_DIR/backup_$DATE.sql.gz s3://backups/content-review/

echo "Backup completed: $DATE"
```

### Restore Database

```bash
# From backup file
gunzip -c backup_20240115_100000.sql.gz | psql -h localhost -U api_user -d content_review
```

---

## Scaling Considerations

### Horizontal Scaling

```yaml
# Multiple API instances in docker-compose
services:
  api-1:
    build: .
    env_file: .env.production
    expose: ["8080"]
  
  api-2:
    build: .
    env_file: .env.production
    expose: ["8080"]
  
  api-3:
    build: .
    env_file: .env.production
    expose: ["8080"]
  
  nginx:
    # Configure upstream with all instances
```

### Database Optimization

```sql
-- Analyze query performance
EXPLAIN ANALYZE
SELECT * FROM contents
WHERE program_id = $1
AND status = 'approved'
ORDER BY created_at DESC;

-- Add indexes if needed
CREATE INDEX idx_contents_program_status ON contents(program_id, status);
```

---

## Security Hardening

### Database Security

```bash
# Restrict PostgreSQL access
sudo nano /etc/postgresql/15/main/pg_hba.conf
# Change: local   all             postgres                                peer
# To:     local   all             postgres                                md5

# Enable SSL in PostgreSQL
sudo nano /etc/postgresql/15/main/postgresql.conf
# Set: ssl = on

# Restart PostgreSQL
sudo systemctl restart postgresql
```

### Application Security

1. Rotate JWT secrets regularly
2. Enable HTTPS/TLS
3. Use security headers
4. Implement rate limiting
5. Regular dependency updates
6. Security scanning with OWASP
7. WAF (Web Application Firewall)

### Firewall Configuration

```bash
# UFW firewall
sudo ufw enable
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 5432/tcp  # PostgreSQL (internal only)
```

---

## Performance Tuning

### PostgreSQL Configuration

```bash
# Edit postgresql.conf
sudo nano /etc/postgresql/15/main/postgresql.conf

# Set for production:
max_connections = 200
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 16MB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
```

### Go Application Tuning

```go
// Set in main.go
runtime.GOMAXPROCS(runtime.NumCPU())

// Database connection pooling
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)
```

---

## Disaster Recovery Plan

### RTO/RPO Targets

- **RTO**: 1 hour (time to recover)
- **RPO**: 1 hour (data loss acceptable)

### Recovery Procedure

```bash
# 1. Check current status
docker-compose ps

# 2. If container failed, restart
docker-compose up -d api

# 3. If database corrupted, restore from backup
docker-compose down
# Restore database from backup
docker-compose up -d

# 4. Verify health
curl https://api.example.com/health

# 5. Test critical endpoints
curl -X GET https://api.example.com/api/v1/contents \
  -H "Authorization: Bearer <test_token>"
```

---

## Maintenance Schedule

| Task | Frequency | Window |
|------|-----------|--------|
| Security Updates | As needed | Off-peak |
| Database Maintenance | Weekly | Sunday 2-3 AM |
| Log Rotation | Daily | 3 AM |
| Certificate Renewal | 30 days before expiry | Automated |
| Backup Verification | Monthly | First Monday |
| Performance Review | Monthly | First Wednesday |
| Security Audit | Quarterly | Last week of quarter |

---

## Troubleshooting

### Common Issues

#### API won't start
```bash
# Check logs
docker-compose logs api

# Verify environment
docker-compose config

# Test database connection
psql -h postgres -U postgres -d content_review -c "SELECT 1"
```

#### High memory usage
```bash
# Check memory limits
docker stats

# Review database queries
SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;
```

#### Slow API response
```bash
# Check slow queries
docker-compose exec postgres psql -U postgres -d content_review \
  -c "ALTER DATABASE content_review SET log_min_duration_statement = 1000;"

# Monitor with pgAdmin
docker-compose exec postgres pgAdmin4
```

---

For more information, see [README.md](./README.md)
