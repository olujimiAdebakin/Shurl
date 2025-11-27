# Shurl API - Complete Documentation Index

## Welcome to Shurl! 🔗

Shurl is a production-ready URL shortener API built with Go, Gin, and PostgreSQL.

---

## 📚 Documentation Guide

### For First-Time Users

**Start here** → [`QUICKSTART.md`](QUICKSTART.md)

- 1-minute setup
- First curl commands
- Test with Swagger UI
- Common tasks

### For API Developers

**Main reference** → [`README.md`](README.md)

- Full API documentation
- All endpoints with examples
- Request/response formats
- Error handling
- Technologies used

### For API Response Details

**Response reference** → [`API_RESPONSES.md`](API_RESPONSES.md)

- All response formats
- Error codes explained
- Validation rules
- Example workflows
- Response data types

### For Swagger/OpenAPI

**Swagger guide** → [`SWAGGER.md`](SWAGGER.md)

- Accessing Swagger UI
- Annotation examples
- Testing endpoints
- Regenerating docs
- Integration guides

### For Production Deployment

**Deployment guide** → [`DEPLOYMENT.md`](DEPLOYMENT.md)

- Docker setup
- Kubernetes manifests
- Nginx configuration
- Database setup
- Security hardening
- Monitoring setup
- Scaling strategies

### For Production Readiness

**Summary** → [`PRODUCTION_READINESS.md`](PRODUCTION_READINESS.md)

- Changes made
- Features completed
- Verification steps
- Deployment paths

---

## 🚀 Quick Start (30 seconds)

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Database

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 3. Run Application

```bash
go run main.go
```

### 4. Access Swagger UI

```
http://localhost:8080/swagger/index.html
```

---

## 📋 API Endpoints Overview

### Authentication (3)

- `POST /api/v1/users/signup` - Create account
- `POST /api/v1/users/login` - Login & get token
- `GET /api/v1/users/validate` - Validate token

### Links (6)

- `POST /api/v1/links` - Create shortened link
- `GET /api/v1/links` - Get your links
- `GET /api/v1/links/:shortCode` - Get link details
- `PATCH /api/v1/links/:shortCode` - Update link
- `DELETE /api/v1/links/:shortCode` - Delete link
- `GET /:shortCode` - Redirect to original URL

### System (1)

- `GET /health` - API health check

**Total: 10 endpoints, all documented and production-ready**

---

## 🔧 Technology Stack

| Component         | Technology      | Version |
| ----------------- | --------------- | ------- |
| Language          | Go              | 1.25.3  |
| Web Framework     | Gin             | Latest  |
| ORM               | GORM            | v1      |
| Database          | PostgreSQL      | 12+     |
| Authentication    | JWT             | RS256   |
| Password Hashing  | Bcrypt          | 10 cost |
| API Documentation | Swagger/OpenAPI | 2.0     |

---

## 📁 File Structure

```
Shurl/
├── README.md                  # Main API documentation
├── QUICKSTART.md              # Quick start guide
├── DEPLOYMENT.md              # Production deployment
├── SWAGGER.md                 # Swagger configuration
├── API_RESPONSES.md           # Response reference
├── PRODUCTION_READINESS.md    # Readiness summary
├── .env.example               # Environment template
│
├── main.go                    # Application entry point
├── go.mod / go.sum            # Go dependencies
│
├── controllers/               # Request handlers
│   ├── user_controllers.go
│   └── link_controller.go
│
├── models/                    # Data structures
│   ├── user.go
│   └── link.go
│
├── dtos/                      # Data transfer objects
│   ├── global_dtos.go
│   ├── user_dtos.go
│   └── link_dtos.go
│
├── middleware/                # Middleware functions
│   ├── require_auth.go
│   └── cors.go
│
├── initializers/              # Initialization
│   ├── database.go
│   └── loadEnv.go
│
├── migrations/                # Database migrations
│   └── migrate.go
│
└── docs/                      # Generated Swagger docs
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
```

---

## ✅ Feature Checklist

### Core Features

- ✅ User authentication (signup/login)
- ✅ JWT token generation (30-day expiration)
- ✅ Password hashing with bcrypt
- ✅ URL shortening with custom codes
- ✅ Click tracking/analytics
- ✅ Link management (CRUD)
- ✅ User-owned link access control

### API Features

- ✅ RESTful endpoint design
- ✅ Consistent response format
- ✅ Comprehensive error handling
- ✅ Request validation
- ✅ CORS support
- ✅ Health check endpoint

### Documentation Features

- ✅ Complete README
- ✅ Quick start guide
- ✅ API response reference
- ✅ Swagger/OpenAPI spec
- ✅ Deployment guide
- ✅ Production readiness guide

### Security Features

- ✅ JWT authentication
- ✅ Bcrypt password hashing
- ✅ SQL injection prevention
- ✅ CORS configured
- ✅ Environment variable management
- ✅ SSL/TLS ready

### Production Features

- ✅ Release mode support
- ✅ Database connection pooling
- ✅ Error logging
- ✅ Health checks
- ✅ Backup strategy
- ✅ Monitoring ready

---

## 🧪 Testing the API

### Using Swagger UI (Easiest)

1. Start the application: `go run main.go`
2. Open: `http://localhost:8080/swagger/index.html`
3. Click "Try it out" on any endpoint

### Using curl (CLI)

```bash
# Sign up
curl -X POST http://localhost:8080/api/v1/users/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secure"}'

# Create link
curl -X POST http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer <token>" \
  -d '{"originalUrl":"https://example.com","shortCode":"test"}'

# Redirect
curl -L http://localhost:8080/test
```

### Using Postman

1. Open Postman
2. Click Import → Link
3. Enter: `http://localhost:8080/swagger/swagger.json`
4. All endpoints imported with examples

---

## 🚀 Deployment Options

### Development

```bash
GIN_MODE=debug go run main.go
```

### Production (Standalone)

```bash
GIN_MODE=release PORT=8080 ./Shurl
```

### Docker

```bash
docker build -t shurl:latest .
docker-compose -f docker-compose.yml up -d
```

### Kubernetes

```bash
kubectl apply -f k8s/manifest.yaml
```

See [`DEPLOYMENT.md`](DEPLOYMENT.md) for detailed instructions.

---

## 🔐 Security

### Implemented

- JWT token authentication
- Bcrypt password hashing (cost 10)
- CORS middleware
- SQL injection prevention
- Environment-based configuration

### Recommended for Production

- HTTPS/TLS termination
- Web Application Firewall (WAF)
- Rate limiting
- DDoS protection
- Database encryption
- Regular security audits

See [`DEPLOYMENT.md`](DEPLOYMENT.md) for security hardening guide.

---

## 📊 Monitoring & Logging

### Health Check

```bash
curl http://localhost:8080/health
```

### Logging

- Application logs to stdout
- Structured logging ready
- Error tracking capability
- Debug mode available

### Metrics Ready For

- Prometheus integration
- CloudWatch integration
- ELK stack integration
- Custom logging solutions

---

## 🆘 Troubleshooting

### Port Already in Use

```bash
PORT=3000 go run main.go
```

### Database Connection Error

1. Verify PostgreSQL is running
2. Check `.env` credentials
3. Ensure database exists

### Swagger UI Not Loading

```bash
swag init -g main.go
go build -o Shurl.exe
go run main.go
```

### Token Expired

Tokens expire after 30 days. Login again to get a new token.

See [`QUICKSTART.md`](QUICKSTART.md) for more troubleshooting.

---

## 📞 Support & Resources

| Resource      | Link                                            |
| ------------- | ----------------------------------------------- |
| Quick Start   | [`QUICKSTART.md`](QUICKSTART.md)                |
| Full Docs     | [`README.md`](README.md)                        |
| API Responses | [`API_RESPONSES.md`](API_RESPONSES.md)          |
| Swagger UI    | `http://localhost:8080/swagger/index.html`      |
| Deployment    | [`DEPLOYMENT.md`](DEPLOYMENT.md)                |
| GitHub Issues | https://github.com/olujimiAdebakin/Shurl/issues |

---

## 📝 API Versions

- **Current Version**: 1.0
- **Base Path**: `/api/v1`
- **Status**: Production-Ready
- **Documentation**: Complete

---

## 🎯 Next Steps

### For Development

1. Read [`QUICKSTART.md`](QUICKSTART.md)
2. Start the application
3. Test endpoints via Swagger UI
4. Read full [`README.md`](README.md)

### For Deployment

1. Read [`DEPLOYMENT.md`](DEPLOYMENT.md)
2. Choose deployment option
3. Configure production environment
4. Deploy and monitor

### For Integration

1. Check [`API_RESPONSES.md`](API_RESPONSES.md)
2. Review [`README.md`](README.md) examples
3. Use Swagger spec for code generation
4. Test with Postman collection

---

## 📈 Project Status

| Aspect           | Status         |
| ---------------- | -------------- |
| Build            | ✅ Passing     |
| Tests            | ✅ Ready       |
| Documentation    | ✅ Complete    |
| Swagger UI       | ✅ Active      |
| Security         | ✅ Implemented |
| Production Ready | ✅ Yes         |

---

## 👤 Author

**Olujimi Adebakin**

- GitHub: [@olujimiAdebakin](https://github.com/olujimiAdebakin)
- LinkedIn: [linkedin.com/in/olujimiadebakin](https://linkedin.com/in/olujimiadebakin)
- Twitter: [@olujimiadebakin](https://twitter.com/olujimiadebakin)

---

## 📄 License

This project is licensed under the MIT License.

---

**Last Updated**: November 27, 2025

**Shurl API v1.0 - Production Ready** ✅
