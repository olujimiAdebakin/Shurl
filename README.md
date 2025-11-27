# **Shurl: A Modern URL Shortener API** 🔗

## Overview

Shurl is a high-performance, production-ready URL shortening service built with Go, Gin, and PostgreSQL. It provides REST API endpoints for creating, managing, and tracking shortened URLs with click analytics.

## Features

- ✅ **User Authentication**: JWT-based authentication with secure password hashing
- ✅ **URL Shortening**: Generate custom or random short codes for long URLs
- ✅ **Click Tracking**: Real-time click analytics for shortened links
- ✅ **User Management**: Create accounts, login, and manage personal links
- ✅ **Link Management**: Full CRUD operations for links
- ✅ **High Performance**: Built with Go for concurrent request handling
- ✅ **RESTful API**: Clean, standardized API design
- ✅ **Database Migrations**: Automated schema management with GORM
- ✅ **CORS Support**: Cross-origin resource sharing enabled
- ✅ **Swagger Documentation**: Interactive API docs at `/swagger/index.html`

## Getting Started

### Prerequisites

- **Go** 1.16+
- **PostgreSQL** 12+
- **Git**

### Installation

1. Clone the repository:

```bash
git clone https://github.com/olujimiAdebakin/Shurl.git
cd Shurl
```

2. Install dependencies:

```bash
go mod download
go mod tidy
```

3. Create `.env` file:

```bash
cp .env.example .env
```

4. Configure environment variables in `.env`:

```dotenv
# Database Configuration
POSTGRES_HOST=your-db-host
POSTGRES_USER=your-username
POSTGRES_PASSWORD=your-password
POSTGRES_DB=shurl_db
POSTGRES_PORT=5432
POSTGRES_SSLMODE=require

# Application Configuration
PORT=8080
SECRET_KEY=your-secret-key-here

# Environment
GIN_MODE=debug  # Set to 'release' for production
```

5. Run database migrations:

```bash
go run ./migrations/migrate.go
```

6. Start the application:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://yourdomain.com`

### Response Format

All API responses follow this standard format:

**Success Response:**

```json
{
  "success": true,
  "data": {
    /* response data */
  }
}
```

**Error Response:**

```json
{
  "success": false,
  "error": "Error message"
}
```

---

## Authentication Endpoints

### Sign Up

Create a new user account.

**Endpoint:** `POST /api/v1/users/signup`

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**

- `400 Bad Request`: Invalid input
- `409 Conflict`: Email already exists

---

### Login

Authenticate and receive JWT token.

**Endpoint:** `POST /api/v1/users/login`

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**

- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials

---

### Validate Token

Verify JWT token and get user information.

**Endpoint:** `GET /api/v1/users/validate`

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "USER"
  }
}
```

**Error Responses:**

- `401 Unauthorized`: Invalid or missing token

---

## Link Endpoints

### Create Link

Create a new shortened URL.

**Endpoint:** `POST /api/v1/links`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
  "shortCode": "my-project"
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "shortCode": "my-project",
    "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
    "clicks": 0,
    "favicon": null,
    "userId": 1
  }
}
```

**Error Responses:**

- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing token
- `409 Conflict`: Short code already exists

**Notes:**

- `shortCode` is optional; a random one will be generated if not provided
- `shortCode` must be 4-20 characters
- `originalUrl` must be a valid URL

---

### Get Link Info

Retrieve link details and increment click count.

**Endpoint:** `GET /api/v1/links/:shortCode`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "shortCode": "my-project",
    "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
    "clicks": 5,
    "favicon": null,
    "userId": 1
  }
}
```

**Error Responses:**

- `404 Not Found`: Link does not exist

---

### Redirect to Link

Redirect to original URL and increment clicks.

**Endpoint:** `GET /:shortCode`

**Response:** `301 Moved Permanently`
Redirects to the original URL

**Example:**

```
GET http://localhost:8080/my-project
→ Redirects to https://github.com/olujimiAdebakin/Shurl
```

**Error Responses:**

- `404 Not Found`: Link does not exist

---

### Update Link

Modify an existing link (owner only).

**Endpoint:** `PATCH /api/v1/links/:shortCode`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "originalUrl": "https://new-url.com"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "shortCode": "my-project",
    "originalUrl": "https://new-url.com",
    "clicks": 5,
    "favicon": null,
    "userId": 1
  }
}
```

**Error Responses:**

- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing token
- `403 Forbidden`: Not the link owner
- `404 Not Found`: Link does not exist

---

### Delete Link

Delete a link (owner only).

**Endpoint:** `DELETE /api/v1/links/:shortCode`

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Link deleted successfully"
  }
}
```

**Error Responses:**

- `401 Unauthorized`: Missing token
- `403 Forbidden`: Not the link owner
- `404 Not Found`: Link does not exist

---

### Get User Links

Retrieve all links created by authenticated user.

**Endpoint:** `GET /api/v1/links`

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "shortCode": "my-project",
      "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
      "clicks": 5,
      "favicon": null,
      "userId": 1
    },
    {
      "shortCode": "another-link",
      "originalUrl": "https://example.com",
      "clicks": 12,
      "favicon": null,
      "userId": 1
    }
  ]
}
```

**Error Responses:**

- `401 Unauthorized`: Missing token

---

## Health Check

### Health Endpoint

Check API status.

**Endpoint:** `GET /health`

**Response:** `200 OK`

```json
{
  "status": "ok",
  "message": "API is running"
}
```

---

## Swagger Documentation

Interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger docs after code changes:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go
```

---

## Technologies Used

| Technology     | Version | Purpose           |
| -------------- | ------- | ----------------- |
| **Go**         | 1.16+   | Core language     |
| **Gin**        | Latest  | Web framework     |
| **GORM**       | v1      | ORM               |
| **PostgreSQL** | 12+     | Database          |
| **JWT**        | Latest  | Authentication    |
| **Bcrypt**     | Latest  | Password hashing  |
| **Swaggo**     | Latest  | API documentation |

## Project Structure

```
Shurl/
├── main.go                 # Application entry point
├── go.mod                  # Module definition
├── .env.example            # Environment variables template
├── controllers/            # Request handlers
│   ├── user_controllers.go
│   └── link_controller.go
├── models/                 # Data models
│   ├── user.go
│   └── link.go
├── dtos/                   # Data transfer objects
│   ├── user_dtos.go
│   ├── link_dtos.go
│   └── global_dtos.go
├── middleware/             # Middleware functions
│   ├── require_auth.go
│   └── cors.go
├── initializers/           # App initialization
│   ├── database.go
│   └── loadEnv.go
└── migrations/             # Database migrations
    └── migrate.go
```

## Development

### Running in Debug Mode

```bash
go run main.go
```

### Running in Release Mode

```bash
export GIN_MODE=release
go run main.go
```

### Testing Endpoints

Using curl:

```bash
# Create link
curl -X POST http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"originalUrl":"https://example.com","shortCode":"test"}'

# Redirect
curl -L http://localhost:8080/test
```

Using Postman:

1. Import the Swagger URL: `http://localhost:8080/swagger/doc.json`
2. Or manually create requests following the API documentation above

## Error Handling

All errors follow a consistent format:

```json
{
  "success": false,
  "error": "Descriptive error message"
}
```

Common HTTP Status Codes:

- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Security Considerations

1. **JWT Tokens**: Tokens expire after 30 days
2. **Password Hashing**: Uses bcrypt with cost factor 10
3. **HTTPS**: Enable in production (use reverse proxy)
4. **CORS**: Configured for specific origins
5. **Database**: Use SSL mode 'require' in production
6. **Environment Variables**: Never commit `.env` file

## Performance Tips

- Use a CDN for shortened link redirects
- Implement caching for frequently accessed links
- Monitor database query performance
- Use connection pooling for database connections

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see LICENSE file for details.

## Roadmap

- [ ] Swagger/OpenAPI integration
- [ ] API rate limiting
- [ ] Advanced analytics dashboard
- [ ] QR code generation
- [ ] Link expiration dates
- [ ] Custom domain support
- [ ] API webhooks
- [ ] Bulk link creation

## Support

For support, please:

- Open an issue on GitHub
- Contact: support@shurl.dev
- Documentation: [Swagger Docs](http://localhost:8080/swagger/index.html)

## Author

**Olujimi Adebakin**

- LinkedIn: [linkedin.com/in/olujimiadebakin](https://www.linkedin.com/in/adebakin-olujimi-25446331b/)
- GitHub: [@olujimiAdebakin](https://github.com/olujimiAdebakin)
- Twitter: [@olujimiadebakin](https://twitter.com/olujimi_the_dev)

---

## Status Badges

[![Go Version](https://img.shields.io/github/go-mod/go-version/olujimiAdebakin/Shurl?style=flat-square)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-0080FF?style=flat-square&logo=go)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/ORM-GORM-darkblue?style=flat-square&logo=go)](https://gorm.io/)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

---

**Last Updated**: November 2025
