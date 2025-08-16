# Simple Production Deployment Guide

## 🚀 **For Projects Without Staging Environment**

This guide is for teams that deploy directly from development to production.

## 📋 **Pre-Deployment Checklist**

### 1. **Environment Variables**
- [ ] Copy `env.example` to `.env.production`
- [ ] Update production values (hosts, database, JWT secret)
- [ ] Ensure `.env.production` is NOT committed to git
- [ ] Set `ENVIRONMENT=production`

### 2. **Security**
- [ ] Generate a strong, unique JWT secret
- [ ] Use HTTPS in production
- [ ] Set up proper CORS origins
- [ ] Use production database credentials

### 3. **Database**
- [ ] Use production PostgreSQL instance
- [ ] Set up database backups
- [ ] Test database migrations

## 🌐 **Environment Configuration**

### **Development**
```bash
ENVIRONMENT=development
DEV_HOST=localhost:8082
PORT=8082
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

## 📚 **Alternative to Staging**

If you want to test production-like environments locally:

### **Local Production Testing**
```bash
# Use a different port to simulate production
ENVIRONMENT=production \
PROD_HOST=localhost:8080 \
PORT=8080 \
go run cmd/main.go
```

### **Docker Compose for Local Testing**
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=production
      - PROD_HOST=localhost:8080
      - DB_URI=postgresql://user:pass@postgres:5432/music_lib
      - JWT_SECRET=test-secret
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=music_lib
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    ports:
      - "5432:5432"
```

## 🎯 **Summary**

- **Skip staging** if you don't need it
- **Deploy directly** from development to production
- **Test locally** with production-like settings
- **Use Docker** for consistent environments
- **Keep it simple** - fewer moving parts means fewer things to go wrong
