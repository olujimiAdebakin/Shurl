# Swagger/OpenAPI Documentation Guide

## Overview

Shurl includes comprehensive Swagger/OpenAPI 3.0 documentation that is automatically generated from Go code annotations.

## Accessing Swagger UI

Once the application is running, visit:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive interface to:

- View all API endpoints
- Read endpoint descriptions and parameters
- Try out requests directly in the browser
- View request/response schemas
- See authentication requirements

## Swagger Files

The Swagger documentation is generated in the `docs/` directory:

```
docs/
├── docs.go         # Generated Go code with embedded swagger spec
├── swagger.json    # OpenAPI 3.0 specification in JSON
└── swagger.yaml    # OpenAPI 3.0 specification in YAML
```

## Regenerating Swagger Docs

If you modify any endpoint annotations or handlers, regenerate the documentation:

```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Regenerate documentation
swag init -g main.go
```

## Annotation Examples

### Function-Level Annotations

```go
// @Summary Sign Up
// @Description Create a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dtos.CreateUserRequest true "Signup credentials"
// @Success 201 {object} dtos.LoginResponse "User created with token"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /users/signup [post]
func SignUpWithToken(c *gin.Context) {
    // implementation
}
```

### Common Annotation Tags

| Tag            | Purpose                                      |
| -------------- | -------------------------------------------- |
| `@Summary`     | One-line endpoint summary                    |
| `@Description` | Detailed endpoint description                |
| `@Tags`        | Group endpoints by category                  |
| `@Accept`      | Content types accepted (json, xml, etc.)     |
| `@Produce`     | Content types produced (json, xml, etc.)     |
| `@Param`       | Define request parameters                    |
| `@Success`     | Define successful response with status code  |
| `@Failure`     | Define error responses with status codes     |
| `@Router`      | Endpoint path and HTTP method                |
| `@Security`    | Specify security scheme (e.g., Bearer token) |

## API Endpoints Overview

### Authentication (`/users`)

| Method | Endpoint          | Auth | Description                |
| ------ | ----------------- | ---- | -------------------------- |
| POST   | `/users/signup`   | ❌   | Create new account         |
| POST   | `/users/login`    | ❌   | Login and get token        |
| GET    | `/users/validate` | ✅   | Validate and inspect token |

### Links (`/links`)

| Method | Endpoint            | Auth | Description           |
| ------ | ------------------- | ---- | --------------------- |
| POST   | `/links`            | ✅   | Create shortened link |
| GET    | `/links`            | ✅   | Get user's links      |
| GET    | `/links/:shortCode` | ❌   | Get link info         |
| PATCH  | `/links/:shortCode` | ✅   | Update link (owner)   |
| DELETE | `/links/:shortCode` | ✅   | Delete link (owner)   |

### Redirect

| Method | Endpoint      | Auth | Description              |
| ------ | ------------- | ---- | ------------------------ |
| GET    | `/:shortCode` | ❌   | Redirect to original URL |

### System

| Method | Endpoint  | Auth | Description      |
| ------ | --------- | ---- | ---------------- |
| GET    | `/health` | ❌   | Check API status |

## Testing with Swagger UI

### Step 1: Sign Up

1. Click "Try it out" on the `/users/signup` endpoint
2. Enter sample credentials:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

3. Click "Execute"
4. Copy the returned `token` value

### Step 2: Authorize

1. Click the "Authorize" button at the top
2. Paste the token in the format: `Bearer <token>`
3. Click "Authorize"

### Step 3: Create a Link

1. Navigate to the `POST /links` endpoint
2. Click "Try it out"
3. Enter:

```json
{
  "originalUrl": "https://github.com/olujimiAdebakin/Shurl",
  "shortCode": "shurl-repo"
}
```

4. Click "Execute"

### Step 4: Test Redirect

1. Navigate to the `GET /:shortCode` endpoint
2. Enter `shortCode`: `shurl-repo`
3. Click "Execute"
4. Observe the redirect response

## Integration with CI/CD

### Generate Docs in CI Pipeline

```yaml
# Example GitHub Actions workflow
- name: Generate Swagger Docs
  run: |
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g main.go

- name: Commit Updated Docs
  run: |
    git add docs/
    git commit -m "chore: update swagger docs"
```

## Hosting Swagger Docs

### Option 1: Built-in Swagger UI (Recommended for Development)

```
http://localhost:8080/swagger/index.html
```

### Option 2: Swagger Editor

Use the online editor to view/edit swagger files:

- Go to: https://editor.swagger.io
- Select "File" → "Import URL"
- Enter: `http://your-domain:8080/swagger/swagger.json` (when deployed)

### Option 3: ReDoc

For a cleaner documentation view:

```html
<!-- Include in your documentation page -->
<redoc spec-url="http://your-domain:8080/swagger/swagger.json"></redoc>
<script src="https://cdn.jsdelivr.net/npm/redoc@latest/bundles/redoc.standalone.js"></script>
```

## Swagger Configuration

The API metadata is defined in `main.go`:

```go
// @title Shurl API
// @version 1.0
// @description A modern URL shortener API built with Go and Gin
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@shurl.dev
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @basePath /api/v1
// @schemes http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
```

## Common Issues

### Swagger UI Shows "404 Not Found"

- Ensure the application is running
- Check that port 8080 is correct (or your configured port)
- Verify `gin-swagger` is imported: `_ "github.com/olujimiAdebakin/Shurl/docs"`

### Annotations Not Appearing

- Regenerate docs: `swag init -g main.go`
- Rebuild and restart the application
- Check that annotations are in comments directly above function definitions

### Swagger UI Not Authenticating

- Use the "Authorize" button at the top of Swagger UI
- Format token as: `Bearer <your_jwt_token>`
- Tokens expire after 30 days

## Best Practices

1. **Keep Annotations Updated**: Update Swagger annotations whenever you modify endpoints
2. **Use Descriptive Summaries**: Make summaries clear and concise (under 200 characters)
3. **Document Error Codes**: Always document possible error responses (400, 401, 403, 404, 500)
4. **Use Tags**: Group related endpoints with tags for better organization
5. **Version Your API**: Use `@version` to track API versions
6. **Document Parameters**: Clearly describe all request parameters
7. **Include Examples**: Provide example request/response bodies in docs

## Troubleshooting

### Update Swagger Dependencies

```bash
go get -u github.com/swaggo/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### Clear Generated Docs

```bash
rm -rf docs/
swag init -g main.go
```

### View Generated Swagger.json

```bash
# Linux/Mac
cat docs/swagger.json | jq

# Windows
type docs\swagger.json
```

## Further Reading

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI 3.0 Specification](https://spec.openapis.org/oas/v3.0.0)
- [Swagger/OpenAPI Tutorial](https://swagger.io/tools/swagger-ui/)
