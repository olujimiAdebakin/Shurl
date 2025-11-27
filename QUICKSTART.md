# Quick Start Guide

## 1-Minute Setup

### Prerequisites

- Go 1.16+
- PostgreSQL 12+

### Installation

```bash
# Clone repository
git clone https://github.com/olujimiAdebakin/Shurl.git
cd Shurl

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Edit .env with your database credentials
# nano .env  # or use your preferred editor
```

### Environment Setup

Edit `.env`:

```dotenv
POSTGRES_HOST=your-database-host
POSTGRES_USER=your-username
POSTGRES_PASSWORD=your-password
POSTGRES_DB=shurl_db
POSTGRES_PORT=5432
POSTGRES_SSLMODE=require
PORT=8080
SECRET_KEY=your-secret-key-here
GIN_MODE=debug
```

### Run the Application

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

---

## Testing the API

### 1. Sign Up

```bash
curl -X POST http://localhost:8080/api/v1/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePassword123"
  }'
```

**Response:**

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

**Save the token** - you'll need it for authenticated requests.

### 2. Create a Shortened Link

```bash
TOKEN="your-token-from-signup"

curl -X POST http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
    "shortCode": "shurl-repo"
  }'
```

**Response:**

```json
{
  "success": true,
  "data": {
    "shortCode": "shurl-repo",
    "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
    "clicks": 0,
    "favicon": null,
    "userId": 1
  }
}
```

### 3. Test the Redirect

```bash
# This will redirect to the original URL
curl -L http://localhost:8080/shurl-repo
```

### 4. Access Swagger UI

Open in your browser:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive interface to test all endpoints.

---

## API Endpoints Overview

### Authentication

- `POST /api/v1/users/signup` - Create account
- `POST /api/v1/users/login` - Login
- `GET /api/v1/users/validate` - Validate token (requires auth)

### Links (Authenticated)

- `POST /api/v1/links` - Create link
- `GET /api/v1/links` - Get your links
- `PATCH /api/v1/links/:shortCode` - Update link
- `DELETE /api/v1/links/:shortCode` - Delete link

### Links (Public)

- `GET /api/v1/links/:shortCode` - Get link info
- `GET /:shortCode` - Redirect to original URL

### System

- `GET /health` - Health check

---

## Common Tasks

### Login and Get Token

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123"
  }'
```

### Get All Your Links

```bash
TOKEN="your-token"

curl -X GET http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer $TOKEN"
```

### Update a Link

```bash
TOKEN="your-token"

curl -X PATCH http://localhost:8080/api/v1/links/shurl-repo \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "originalUrl": "https://new-url.com"
  }'
```

### Delete a Link

```bash
TOKEN="your-token"

curl -X DELETE http://localhost:8080/api/v1/links/shurl-repo \
  -H "Authorization: Bearer $TOKEN"
```

---

## Using Postman

### Import Swagger Spec

1. Open Postman
2. Click "Import"
3. Go to "Link" tab
4. Enter: `http://localhost:8080/swagger/swagger.json`
5. Click "Import"

All endpoints will be imported with pre-configured request templates.

---

## Using Swagger UI

1. Start the application
2. Open `http://localhost:8080/swagger/index.html`
3. Click "Authorize" button (top right)
4. Enter token as: `Bearer <your-token>`
5. Click on any endpoint to "Try it out"

---

## Building for Production

```bash
# Build executable
go build -o Shurl

# Run with release mode
GIN_MODE=release PORT=8080 ./Shurl
```

---

## Docker Deployment

```bash
# Build image
docker build -t shurl:latest .

# Run container
docker run -p 8080:8080 \
  -e POSTGRES_HOST=host.docker.internal \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=shurl_db \
  -e SECRET_KEY=your-secret \
  shurl:latest
```

---

## Troubleshooting

### Port Already in Use

```bash
# Use a different port
PORT=3000 go run main.go
```

### Database Connection Error

1. Verify PostgreSQL is running
2. Check `.env` credentials
3. Ensure database exists: `createdb shurl_db`
4. Test connection: `psql -h localhost -U postgres -d shurl_db`

### Swagger UI Not Loading

- Ensure `docs/` directory exists
- Regenerate docs: `swag init -g main.go`
- Restart the application

### Token Expired

Tokens expire after 30 days. Login again to get a new token.

---

## Next Steps

1. **Read Full Documentation**: See [README.md](README.md)
2. **Explore Swagger Docs**: Visit `http://localhost:8080/swagger/index.html`
3. **Production Deployment**: See [DEPLOYMENT.md](DEPLOYMENT.md)
4. **Swagger Configuration**: See [SWAGGER.md](SWAGGER.md)

---

## Support

- Documentation: [Full API Docs](README.md)
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Issues: [GitHub Issues](https://github.com/olujimiAdebakin/Shurl/issues)
