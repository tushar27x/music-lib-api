# Production Deployment Guide

## 🚀 **Pre-Deployment Checklist**

### 1. **Environment Variables**
- [ ] Copy `env.example` to `.env.production`
- [ ] Update all production values (hosts, database, JWT secret)
- [ ] Ensure `.env.production` is NOT committed to git
- [ ] Set `ENVIRONMENT=production`

### 2. **Security**
- [ ] Generate a strong, unique JWT secret
- [ ] Use HTTPS in production
- [ ] Set up proper CORS origins
- [ ] Configure rate limiting
- [ ] Use environment-specific database credentials

### 3. **Database**
- [ ] Use production PostgreSQL instance
- [ ] Set up proper database backups
- [ ] Configure connection pooling
- [ ] Test database migrations

## 🌐 **Environment Configuration**

> **Note:** If you don't have a staging domain, you can:
> - Skip staging environment entirely and go directly from development to production
> - Use a subdomain like `staging.yourdomain.com` instead of `staging-api.yourdomain.com`
> - Use localhost with different ports for staging (e.g., `localhost:8081`)

### **Development**
```bash
ENVIRONMENT=development
DEV_HOST=localhost:8082
PORT=8082
```

### **Staging (Optional)**
```bash
ENVIRONMENT=staging
# STAGING_HOST=staging-api.yourdomain.com  # Optional - if not set, will use PROD_HOST
PORT=8080
```

### **Production**
```bash
ENVIRONMENT=production
PROD_HOST=api.yourdomain.com
PORT=8080
```

## 📝 **Generating Production Swagger Docs**

### **Using Scripts**
```bash
# Linux/Mac
./scripts/generate-swagger.sh production

# Windows
scripts\generate-swagger.bat production
```

### **Manual Generation**
```bash
# Set environment variables
export ENVIRONMENT=production
export PROD_HOST=api.yourdomain.com

# Generate docs
swag init -g cmd/main.go
```

## 🐳 **Docker Deployment**

### **Dockerfile (already exists)**
```dockerfile
FROM golang:1.23.5-alpine
# ... existing configuration
```

### **Docker Compose for Production**
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=production
      - PROD_HOST=api.yourdomain.com
      - DB_URI=${DB_URI}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

## 🔒 **Security Best Practices**

### **JWT Configuration**
```bash
# Generate a strong secret (32+ characters)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### **CORS Configuration**
```bash
# Only allow your frontend domain
CORS_ORIGIN=https://yourdomain.com
```

### **Rate Limiting**
```bash
RATE_LIMIT=100
RATE_LIMIT_WINDOW=1m
```

## 📊 **Monitoring & Logging**

### **Health Check Endpoint**
```bash
GET /api/ping
```

### **Log Levels**
- `development`: debug
- `staging`: info
- `production`: warn/error

## 🚀 **Deployment Commands**

### **1. Build for Production**
```bash
go build -ldflags="-s -w" -o music-api cmd/main.go
```

### **2. Run with Production Config**
```bash
ENVIRONMENT=production PROD_HOST=api.yourdomain.com ./music-api
```

### **3. Using Docker**
```bash
docker build -t music-api .
docker run -d \
  -e ENVIRONMENT=production \
  -e PROD_HOST=api.yourdomain.com \
  -e DB_URI="your-db-uri" \
  -e JWT_SECRET="your-jwt-secret" \
  -p 8080:8080 \
  music-api
```

## 🔍 **Post-Deployment Verification**

### **1. Health Check**
```bash
curl https://api.yourdomain.com/api/ping
```

### **2. Swagger Documentation**
```bash
# Should show production host
https://api.yourdomain.com/swagger/index.html
```

### **3. Database Connection**
- Check logs for successful database connection
- Verify tables are created/migrated

### **4. Authentication**
- Test user registration
- Test user login
- Verify JWT tokens work

## 🚨 **Common Issues & Solutions**

### **Swagger Host Not Updated**
```bash
# Regenerate docs for production
./scripts/generate-swagger.sh production
```

### **Environment Variables Not Loading**
- Ensure `.env.production` exists
- Check file permissions
- Verify variable names match code

### **Database Connection Issues**
- Verify database credentials
- Check network connectivity
- Ensure database is running

## 📚 **Additional Resources**

- [Go Production Best Practices](https://golang.org/doc/effective_go.html)
- [Gin Production Guidelines](https://gin-gonic.com/docs/examples/graceful-restart-or-stop/)
- [GORM Production Tips](https://gorm.io/docs/performance.html)
- [Swagger Production Setup](https://swaggo.github.io/swaggo.io/)
