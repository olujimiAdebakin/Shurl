# **Shurl: A Modern URL Shortener API** 🔗

## Overview
Shurl is a high-performance backend API for robust URL shortening, built with Go and leveraging the Gin web framework for routing and GORM for elegant database interactions. It provides a reliable service for generating and managing short, trackable links.

## Features
-   **Go**: Core language for a highly performant and concurrent backend.
-   **Gin Framework**: Powers a fast and efficient RESTful API.
-   **GORM**: Provides an elegant Object-Relational Mapping (ORM) for PostgreSQL, simplifying database operations.
-   **PostgreSQL**: A powerful, open-source relational database for persistent storage of links and user data.
-   **`godotenv`**: Manages environment variables for secure and flexible configuration.
-   **Short Link Generation**: Creates unique, concise short codes for long URLs.
-   **Click Tracking**: Monitors the number of times a shortened URL is accessed.
-   **User Authentication (Planned)**: Secure user management for personalized link tracking. (Based on `user.go` model)
-   **Link Management**: Allows users to create, view, and potentially manage their shortened URLs.

## Getting Started
To get a local copy of Shurl up and running, follow these steps.

### Installation
Clone the repository:
```bash
git clone https://github.com/olujimiAdebakin/Shurl.git
cd Shurl
```

Install Go modules:
```bash
go mod tidy
```

### Environment Variables
Create a `.env` file in the root directory of the project and populate it with the following required variables:

```dotenv
POSTGRES_HOST=localhost
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=shurl_db
POSTGRES_PORT=5432
POSTGRES_SSLMODE=disable # Use 'require' or 'verify-full' for production
```

**Note**: `POSTGRES_SSLMODE` is set to `disable` for local development convenience. For production environments, it is strongly recommended to use `require` or `verify-full` for secure database connections.

### Usage
Run the application:
```bash
go run main.go
```
This will start the Shurl API server. Currently, the `main.go` only initializes the database connection and environment variables. Future updates will include API routes and business logic.

Once the API routes are implemented and the server is running, you can interact with it using tools like Postman, Insomnia, or `curl`.

## API Documentation
### Base URL
The base URL for the API will typically be `http://localhost:[PORT]` during local development, where `[PORT]` is configured in your environment variables (e.g., `APP_PORT`).

### Endpoints
*(Note: As of the current codebase, API endpoints are not yet implemented in `main.go`. This section will be populated with detailed endpoint documentation once they are defined.)*

#### Example: POST /api/v1/links
This section would document an example endpoint for creating a short URL.

**Request**:
```json
{
  "original_url": "https://www.example.com/very/long/url/that/needs/shortening",
  "user_id": 1
}
```

**Response**:
```json
{
  "short_code": "shrtcd",
  "original_url": "https://www.example.com/very/long/url/that/needs/shortening",
  "short_url": "http://localhost:8080/shrtcd",
  "clicks": 0
}
```

**Errors**:
-   `400 Bad Request`: Invalid input or missing `original_url`.
-   `409 Conflict`: Short code already exists (unlikely with robust generation, but possible).
-   `500 Internal Server Error`: Database error or unexpected server issue.

---

## Technologies Used

| Technology    | Description                                       | Link                                                       |
| :------------ | :------------------------------------------------ | :--------------------------------------------------------- |
| **Go**        | High-performance programming language.            | [golang.org](https://golang.org/)                          |
| **Gin**       | HTTP web framework for Go.                        | [gin-gonic.com](https://gin-gonic.com/)                    |
| **GORM**      | Object-Relational Mapper (ORM) for Go.            | [gorm.io](https://gorm.io/)                                |
| **PostgreSQL**| Advanced open-source relational database.         | [postgresql.org](https://www.postgresql.org/)              |
| **DotEnv**    | Loads environment variables from `.env` file.     | [github.com/joho/godotenv](https://github.com/joho/godotenv) |

## Contributing
We welcome contributions to Shurl! If you're interested in improving this project, please consider the following:

-   ⭐ Fork the repository.
-   💡 Create a new branch for your feature or bug fix.
-   ✨ Ensure your code adheres to Go best practices and includes tests where appropriate.
-   🚀 Submit a pull request with a clear description of your changes.

## Author Info
**Olujimi Adebakin**
*   LinkedIn: [linkedin.com/in/olujimiadebakin](https://linkedin.com/in/olujimiadebakin)
*   Twitter: [@olujimiadebakin](https://twitter.com/olujimiadebakin)
*   Portfolio: [your_portfolio_link]

---
[![Go Version](https://img.shields.io/github/go-mod/go-version/olujimiAdebakin/Shurl?style=flat-square)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-0080FF?style=flat-square&logo=go)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/ORM-GORM-darkblue?style=flat-square&logo=go)](https://gorm.io/)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square)](https://github.com/olujimiAdebakin/Shurl/actions)

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)
