# API Response Reference

## Response Format

All API responses follow a consistent JSON format:

### Success Response

```json
{
  "success": true,
  "data": {
    /* response data */
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

---

## Authentication Responses

### Sign Up (POST /api/v1/users/signup)

**Request:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Success Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsImV4cCI6MTczNjA5MjAyNH0.signature"
  }
}
```

**Error Response (400 Bad Request):**

```json
{
  "success": false,
  "error": "Invalid input: Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

**Error Response (409 Conflict):**

```json
{
  "success": false,
  "error": "User with this email already exists"
}
```

---

### Login (POST /api/v1/users/login)

**Request:**

```json
{
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsImV4cCI6MTczNjA5MjAyNH0.signature"
  }
}
```

**Error Response (401 Unauthorized):**

```json
{
  "success": false,
  "error": "Invalid email or password"
}
```

---

### Validate Token (GET /api/v1/users/validate)

**Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200 OK):**

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

**Error Response (401 Unauthorized):**

```json
{
  "success": false,
  "error": "User not found in context"
}
```

---

## Link Endpoint Responses

### Create Link (POST /api/v1/links)

**Request:**

```json
{
  "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
  "shortCode": "shurl-repo"
}
```

**Success Response (201 Created):**

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

**Error Response (400 Bad Request):**

```json
{
  "success": false,
  "error": "Invalid input: Key: 'CreateLinkRequest.OriginalURL' Error:Field validation for 'OriginalURL' failed on the 'url' tag"
}
```

**Error Response (409 Conflict):**

```json
{
  "success": false,
  "error": "Short code already exists"
}
```

**Error Response (401 Unauthorized):**

```json
{
  "success": false,
  "error": "Unauthorized - User not found"
}
```

---

### Get Link (GET /api/v1/links/:shortCode)

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "shortCode": "shurl-repo",
    "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
    "clicks": 5,
    "favicon": null,
    "userId": 1
  }
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "error": "Link not found"
}
```

---

### Get User Links (GET /api/v1/links)

**Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "shortCode": "shurl-repo",
      "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
      "clicks": 5,
      "favicon": null,
      "userId": 1
    },
    {
      "shortCode": "my-blog",
      "originalUrl": "https://example.com",
      "clicks": 12,
      "favicon": null,
      "userId": 1
    }
  ]
}
```

**Error Response (401 Unauthorized):**

```json
{
  "success": false,
  "error": "User not found in context"
}
```

---

### Update Link (PATCH /api/v1/links/:shortCode)

**Request:**

```json
{
  "originalUrl": "https://new-url.com"
}
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "shortCode": "shurl-repo",
    "originalUrl": "https://new-url.com",
    "clicks": 5,
    "favicon": null,
    "userId": 1
  }
}
```

**Error Response (400 Bad Request):**

```json
{
  "success": false,
  "error": "Invalid input: Key: 'LinkUpdateRequest.OriginalURL' Error:Field validation for 'OriginalURL' failed on the 'url' tag"
}
```

**Error Response (403 Forbidden):**

```json
{
  "success": false,
  "error": "You can only update your own links"
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "error": "Link not found"
}
```

---

### Delete Link (DELETE /api/v1/links/:shortCode)

**Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "message": "Link deleted successfully"
  }
}
```

**Error Response (403 Forbidden):**

```json
{
  "success": false,
  "error": "You can only delete your own links"
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "error": "Link not found"
}
```

---

## Redirect Response

### Redirect (GET /:shortCode)

**Success Response (301 Moved Permanently):**

```
Location: https://github.com/olujimiAdebakin/Shurl
```

The client will automatically follow the redirect to the original URL.

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "error": "Link not found"
}
```

---

## System Responses

### Health Check (GET /health)

**Success Response (200 OK):**

```json
{
  "status": "ok",
  "message": "API is running"
}
```

---

## HTTP Status Codes Reference

| Code | Meaning               | Common Causes                                        |
| ---- | --------------------- | ---------------------------------------------------- |
| 200  | OK                    | Request successful                                   |
| 201  | Created               | Resource successfully created                        |
| 301  | Moved Permanently     | Redirect to original URL                             |
| 400  | Bad Request           | Invalid input, validation failed                     |
| 401  | Unauthorized          | Missing or invalid authentication token              |
| 403  | Forbidden             | Authenticated but not authorized (e.g., not owner)   |
| 404  | Not Found             | Resource doesn't exist                               |
| 409  | Conflict              | Resource already exists (e.g., duplicate short code) |
| 500  | Internal Server Error | Unexpected server error                              |

---

## Error Categories

### Authentication Errors (401)

- Missing Authorization header
- Invalid token format
- Expired token
- Invalid token signature
- User not found

### Validation Errors (400)

- Missing required fields
- Invalid field format (e.g., email, URL)
- Invalid JSON
- Field value constraints violated

### Authorization Errors (403)

- Not the link owner
- Permission denied
- Insufficient role

### Not Found Errors (404)

- Link doesn't exist
- User doesn't exist
- Resource deleted

### Conflict Errors (409)

- Email already exists
- Short code already taken
- Duplicate resource

### Server Errors (500)

- Database connection failed
- Unexpected exception
- Token generation failed

---

## Field Validation Rules

### CreateUserRequest

```
name:     required, min 2 chars
email:    required, valid email format
password: required, min 6 chars
```

### LoginUserRequest

```
email:    required, valid email format
password: required, min 6 chars
```

### CreateLinkRequest

```
originalUrl: required, valid URL format
shortCode:   optional, min 4 chars, max 20 chars
```

### LinkUpdateRequest

```
originalUrl: optional, valid URL format if provided
```

---

## Response Data Types

### User Object

```typescript
{
  id: number,
  name: string,
  email: string,
  token?: string,
  role?: string
}
```

### Link Object

```typescript
{
  shortCode: string,
  originalUrl: string,
  clicks: number,
  favicon: string | null,
  userId: number
}
```

### Error Object

```typescript
{
  success: false,
  error: string
}
```

### Success Object

```typescript
{
  success: true,
  data: any
}
```

---

## Example Workflows

### Complete User Signup → Create Link → Redirect Workflow

```bash
# 1. Sign up
RESPONSE=$(curl -X POST http://localhost:8080/api/v1/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePassword123"
  }')

TOKEN=$(echo $RESPONSE | jq -r '.data.token')

# 2. Create link
LINK_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "originalUrl": "https://example.com",
    "shortCode": "test-link"
  }')

SHORT_CODE=$(echo $LINK_RESPONSE | jq -r '.data.shortCode')

# 3. Redirect via short code
curl -L http://localhost:8080/$SHORT_CODE
# Will redirect to https://example.com
```

---

## Rate Limiting (Future)

Currently no rate limiting is implemented. For production, consider:

- Per-user request limits
- IP-based throttling
- Token-based rate limiting
- Burst allowances

---

## Pagination (Future)

Currently no pagination is implemented for list endpoints. All user links are returned at once.

Future implementation may include:

- `?page=1&limit=10` parameters
- Cursor-based pagination
- Offset/limit pagination

---

## API Versioning

- **Current Version**: 1.0
- **Base Path**: `/api/v1`
- **Stability**: Production-ready

Future versions will be available at:

- `/api/v2` (when applicable)
- Backward compatibility will be maintained

---

## Additional Resources

- [Full API Documentation](README.md)
- [Quick Start Guide](QUICKSTART.md)
- [Swagger/OpenAPI Docs](http://localhost:8080/swagger/index.html)
- [Deployment Guide](DEPLOYMENT.md)
