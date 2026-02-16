# Integration and Deployment Guide

## Overview

This document provides comprehensive guidance for integrating, testing, and deploying the Financial Planning Calculator application.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Browser    │  │    Mobile    │  │   Desktop    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTPS
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (Next.js)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Pages     │  │  Components  │  │   Contexts   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ REST API
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend (Go/Echo)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Controllers  │  │  Use Cases   │  │   Domain     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ SQL
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Database (PostgreSQL)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Financial    │  │    Goals     │  │   Reports    │      │
│  │    Data      │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## System Requirements

### Development Environment

- **Node.js**: 18.x or higher
- **Go**: 1.21 or higher
- **PostgreSQL**: 15.x or higher
- **Docker**: 24.x or higher (optional)
- **Git**: 2.x or higher

### Production Environment

- **CPU**: 2+ cores recommended
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 20GB minimum
- **Network**: Stable internet connection
- **OS**: Linux (Ubuntu 22.04 LTS recommended)

## Local Development Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd financial-planning-calculator
```

### 2. Setup Database

```bash
# Using Docker
docker-compose up -d postgres

# Or install PostgreSQL locally
# Then create database
createdb financial_planning
```

### 3. Setup Backend

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your settings
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=password
# DB_NAME=financial_planning

# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go

# Seed database (optional)
go run cmd/seed/main.go

# Start server
go run main.go
```

Backend will be available at `http://localhost:8080`

### 4. Setup Frontend

```bash
cd frontend

# Copy environment file
cp .env.example .env.local

# Edit .env.local
# NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Install dependencies
npm install

# Start development server
npm run dev
```

Frontend will be available at `http://localhost:3000`

### 5. Verify Integration

```bash
# Run integration test script
./scripts/test-integration.sh
```

## Testing

### Unit Tests

#### Backend

```bash
cd backend
go test ./... -v
go test ./... -cover
```

#### Frontend

```bash
cd frontend
npm test
npm run test:coverage
```

### Integration Tests

```bash
# Ensure both servers are running
./scripts/test-integration.sh
```

### E2E Tests

```bash
cd e2e

# Install dependencies
npm install
npm run install

# Run all tests
npm test

# Run specific test suite
npm run test:financial
npm run test:goals
npm run test:api

# Run with UI
npm run test:ui

# Run in debug mode
npm run test:debug
```

## API Documentation

### Accessing Swagger UI

Once the backend is running, access the API documentation at:

```
http://localhost:8080/swagger/index.html
```

### Key Endpoints

#### Health Check

```bash
GET /health
GET /health/detailed
GET /ready
```

#### Financial Data

```bash
POST   /api/financial-data
GET    /api/financial-data?user_id={id}
PUT    /api/financial-data/{user_id}/profile
PUT    /api/financial-data/{user_id}/retirement
PUT    /api/financial-data/{user_id}/emergency-fund
DELETE /api/financial-data/{user_id}
```

#### Calculations

```bash
POST /api/calculations/asset-projection
POST /api/calculations/retirement
POST /api/calculations/emergency-fund
POST /api/calculations/goal-projection
```

#### Goals

```bash
POST   /api/goals
GET    /api/goals?user_id={id}
GET    /api/goals/{id}?user_id={id}
PUT    /api/goals/{id}?user_id={id}
PUT    /api/goals/{id}/progress?user_id={id}
DELETE /api/goals/{id}?user_id={id}
```

## Performance Optimization

### Frontend Optimizations

1. **Code Splitting**: Automatic route-based splitting
2. **Image Optimization**: Next.js Image component
3. **Caching**: API response caching (5 min TTL)
4. **Memoization**: React.memo, useMemo, useCallback
5. **Bundle Size**: Tree shaking and minification

See [frontend/PERFORMANCE.md](frontend/PERFORMANCE.md) for details.

### Backend Optimizations

1. **Connection Pooling**: 25 max connections
2. **Query Optimization**: Indexed queries
3. **Caching**: In-memory cache for calculations
4. **Compression**: Gzip compression enabled
5. **Rate Limiting**: 100 requests/second

See [backend/infrastructure/database/optimization.md](backend/infrastructure/database/optimization.md) for details.

## Deployment

### Using Docker Compose

```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop services
docker-compose -f docker-compose.prod.yml down
```

### Manual Deployment

#### Backend

```bash
cd backend

# Build binary
go build -o server main.go

# Run with production settings
PORT=8080 \
DEBUG=false \
DB_HOST=your-db-host \
DB_PORT=5432 \
DB_USER=your-db-user \
DB_PASSWORD=your-db-password \
DB_NAME=financial_planning \
ALLOWED_ORIGINS=https://your-domain.com \
./server
```

#### Frontend

```bash
cd frontend

# Build for production
npm run build

# Start production server
npm run start
```

### Environment Variables

#### Backend

```bash
# Server
PORT=8080
DEBUG=false

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=financial_planning
DB_SSLMODE=require

# CORS
ALLOWED_ORIGINS=https://your-domain.com

# Performance (IP-based rate limiting)
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=50
REQUEST_TIMEOUT=30s
MAX_REQUEST_SIZE=10M
ENABLE_GZIP=true
```

#### Frontend

```bash
# API
NEXT_PUBLIC_API_URL=https://api.your-domain.com/api

# Environment
NODE_ENV=production
```

## Monitoring

### Health Checks

```bash
# Basic health
curl http://localhost:8080/health

# Detailed health with component status
curl http://localhost:8080/health/detailed

# Readiness check
curl http://localhost:8080/ready
```

### Performance Metrics

Backend exposes performance metrics:

```bash
curl http://localhost:8080/api/metrics
```

### Logging

#### Backend Logs

```bash
# View logs
tail -f backend/logs/app.log

# Filter errors
grep ERROR backend/logs/app.log
```

#### Frontend Logs

```bash
# Development
npm run dev

# Production
pm2 logs frontend
```

## Troubleshooting

### Common Issues

#### Backend won't start

1. Check database connection
2. Verify environment variables
3. Check port availability
4. Review logs for errors

```bash
# Test database connection
psql -h localhost -U postgres -d financial_planning

# Check if port is in use
lsof -i :8080
```

#### Frontend can't connect to backend

1. Verify backend is running
2. Check CORS configuration
3. Verify API URL in .env.local
4. Check network connectivity

```bash
# Test backend connectivity
curl http://localhost:8080/health

# Check CORS
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -X OPTIONS \
     http://localhost:8080/api/financial-data
```

#### Database connection errors

1. Verify PostgreSQL is running
2. Check credentials
3. Verify database exists
4. Check network/firewall

```bash
# Test connection
psql -h localhost -U postgres -d financial_planning -c "SELECT 1"

# Check PostgreSQL status
systemctl status postgresql
```

#### E2E tests failing

1. Ensure servers are running
2. Check test configuration
3. Verify test data
4. Review screenshots/videos

```bash
# Run with debug
cd e2e
npm run test:debug

# View test report
npm run report
```

## Security Considerations

### Backend Security

1. **Input Validation**: All inputs validated
2. **SQL Injection**: Prepared statements used
3. **CORS**: Configured for specific origins
4. **Rate Limiting**: Prevents abuse
5. **HTTPS**: Required in production

### Frontend Security

1. **XSS Protection**: React escapes by default
2. **CSRF**: Token-based protection
3. **Content Security Policy**: Configured headers
4. **Secure Cookies**: HttpOnly, Secure flags

### Database Security

1. **Encryption**: SSL/TLS for connections
2. **Access Control**: Limited user permissions
3. **Backups**: Regular automated backups
4. **Audit Logging**: Track data changes

## Backup and Recovery

### Database Backup

```bash
# Create backup
pg_dump -h localhost -U postgres financial_planning > backup.sql

# Restore backup
psql -h localhost -U postgres financial_planning < backup.sql
```

### Automated Backups

```bash
# Add to crontab
0 2 * * * pg_dump -h localhost -U postgres financial_planning > /backups/$(date +\%Y\%m\%d).sql
```

## Support and Resources

### Documentation

- [Frontend Performance Guide](frontend/PERFORMANCE.md)
- [Backend Optimization Guide](backend/infrastructure/database/optimization.md)
- [E2E Testing Guide](e2e/README.md)
- [API Documentation](http://localhost:8080/swagger/index.html)

### Tools

- [Swagger UI](http://localhost:8080/swagger/index.html) - API documentation
- [Playwright](https://playwright.dev/) - E2E testing
- [Go Documentation](https://pkg.go.dev/) - Go packages
- [Next.js Documentation](https://nextjs.org/docs) - Frontend framework

### Getting Help

1. Check documentation
2. Review logs
3. Run integration tests
4. Check GitHub issues
5. Contact development team

## Changelog

### Version 1.0.0

- Initial release
- Complete financial planning features
- E2E test coverage
- Performance optimizations
- Production-ready deployment
