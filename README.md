# Greenlight API

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![RESTful API](https://img.shields.io/badge/API-RESTful-007ACC?style=for-the-badge)
![Makefile](https://img.shields.io/badge/Makefile-000000?style=for-the-badge&logo=gnu-make&logoColor=white)

## 📚 Table of Contents

- [Greenlight API](#greenlight-api)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [Technologies Used](#technologies-used)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Setup](#setup)
    - [Running the API](#running-the-api)
    - [Database Migrations](#database-migrations)
  - [Makefile Commands](#makefile-commands)
  - [API Documentation](#api-documentation)
  - [Project Structure](#project-structure)
  - [Contributing](#contributing)

## 🌟 Overview

Greenlight is a robust and scalable RESTful API built with Go, designed to manage movie information, user authentication, and permissions. It provides a clean and efficient backend for a movie catalog application, featuring secure user management, token-based authentication, and a well-structured database.

## ✨ Features

-   **Movie Management:** CRUD operations for movie entries, including title, year, runtime, and genres.
-   **User Authentication:** Secure user registration, login, and password management using bcrypt hashing.
-   **Token-Based Authentication:** API key and activation token support for secure access and user verification.
-   **Permissions System:** Role-based access control for different user types (e.g., admin, regular user).
-   **Database Migrations:** Managed database schema evolution using `migrate`.
-   **Structured Logging:** JSON-based logging for better observability.
-   **Email Notifications:** Integration for sending welcome emails to new users.
-   **Health Checks:** Endpoint for monitoring API health.
-   **Rate Limiting:** Basic rate limiting to prevent abuse.
-   **CORS Support:** Configurable Cross-Origin Resource Sharing.

## 🛠️ Technologies Used

-   **Go (Golang):** The primary language for building the API.
-   **PostgreSQL:** Relational database for storing application data.
-   **`github.com/julienschmidt/httprouter`:** High-performance HTTP request router.
-   **`golang.org/x/crypto/bcrypt`:** For secure password hashing.
-   **`github.com/lib/pq`:** PostgreSQL driver for Go.
-   **`github.com/joho/godotenv`:** For loading environment variables from `.env` files.
-   **`gopkg.in/mail.v2`:** For sending emails.
-   **`github.com/tomasen/realip`:** For getting the real IP address of a client.
-   **`github.com/swaggo/swag`:** For generating OpenAPI (Swagger) documentation.
-   **`github.com/golang-migrate/migrate/v4`:** For database schema migrations.

## 🚀 Getting Started

### 📋 Prerequisites

Before you begin, ensure you have the following installed:

-   Go (version 1.22 or higher)
-   PostgreSQL
-   `migrate` command-line tool:
    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
-   `swag` command-line tool:
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```
-   `staticcheck` command-line tool:
    ```bash
    go install honnef.co/go/tools/cmd/staticcheck@latest
    ```

### 🛠️ Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/greenlight.git
    cd greenlight
    ```
2.  **Create a `.env` file:**
    Create a `.env` file in the root directory of the project and add your database connection string and other environment variables.
    ```env
    DB_DSN="postgres://user:password@localhost:5432/greenlight?sslmode=disable"
    SMTP_HOST="smtp.mailtrap.io"
    SMTP_PORT=2525
    SMTP_USERNAME="your_mailtrap_username"
    SMTP_PASSWORD="your_mailtrap_password"
    SMTP_SENDER="Greenlight <no-reply@greenlight.com>"
    CORS_TRUSTED_ORIGINS="http://localhost:4000,http://localhost:8080"
    ```
    *Replace with your actual database credentials and SMTP settings.*
3.  **Install dependencies and vendor them:**
    ```bash
    make vendor
    ```
4.  **Run database migrations:**
    ```bash
    make db/migrations/up
    ```

### ▶️ Running the API

To start the API server, use the following command:

```bash
make run/api
```

The API will typically run on `http://localhost:4000`.

### 🗄️ Database Migrations

The project uses `migrate` for database migrations. You can manage your database schema using the following Makefile commands:

-   **Create a new migration:**
    ```bash
    make db/migrations/new name=create_movies_table
    ```
    This will create two new files in the `migrations` directory: `[timestamp]_create_movies_table.up.sql` and `[timestamp]_create_movies_table.down.sql`.
-   **Apply all up migrations:**
    ```bash
    make db/migrations/up
    ```
-   **Rollback all migrations:**
    ```bash
    make db/migrations/down
    ```
-   **Migrate to a specific version:**
    ```bash
    make db/migrations/goto version=<version_number>
    ```
-   **Rollback a specific number of migrations:**
    ```bash
    make db/migrations/rollback version=<number_of_migrations>
    ```
-   **Force database schema to a specific version (use with caution):**
    ```bash
    make db/migrations/force version=<version_number>
    ```

## ⚙️ Makefile Commands

The `Makefile` provides several convenient commands for development, quality control, and building.

-   `make help`: Displays a list of all available Makefile commands with their descriptions.
-   `make audit`: Formats the Go code, runs `go vet`, `staticcheck`, and executes all tests with race detection.
-   `make vendor`: Tidies and verifies Go module dependencies, then vendors them into the `vendor/` directory.
-   `make run/api`: Runs the main Greenlight API server. Requires `DB_DSN` to be set in the `.env` file.
-   `make db/migrations/new name=...`: Creates a new database migration file. Replace `...` with the desired migration name.
-   `make db/migrations/up`: Applies all pending "up" database migrations.
-   `make db/migrations/goto version=...`: Migrates the database to a specific version. Replace `...` with the target version number.
-   `make db/migrations/down`: Applies all "down" database migrations, effectively rolling back the entire database schema.
-   `make db/migrations/rollback version=...`: Rolls back a specific number of migrations. Replace `...` with the number of migrations to rollback.
-   `make db/migrations/force version=...`: Forces the database schema to a specific version. Use with extreme caution as this can lead to data loss.
-   `make docs/gen`: Generates OpenAPI (Swagger) documentation for the API using `swag`.
-   `make build/api`: Builds the `cmd/api` application for the current operating system and also creates a Linux AMD64 executable in `bin/linux_amd64/api`.

## 📄 API Documentation

API documentation is generated using `swag`. To generate or update the documentation:

```bash
make docs/gen
```

Once generated, the documentation can typically be served by a Swagger UI instance, often available at `/swagger/index.html` when the API is running.

## 📂 Project Structure

```
.
├── cmd/api/             # Main API application code
│   ├── context.go       # Context handling for requests
│   ├── errors.go        # Custom error types and handlers
│   ├── healthcheck.go   # API health check endpoint
│   ├── helpers.go       # Utility functions for API handlers
│   ├── main.go          # Main entry point for the API server
│   ├── middleware.go    # HTTP middleware functions (e.g., authentication, rate limiting)
│   ├── movies.go        # Handlers for movie-related operations
│   ├── routes.go        # Defines API routes
│   ├── server.go        # HTTP server setup and configuration
│   ├── tokens.go        # Handlers for token-related operations (e.g., activation, authentication)
│   └── users.go         # Handlers for user-related operations
├── internal/            # Internal packages and business logic
│   ├── env/             # Environment variable loading
│   ├── jsonlog/         # Structured JSON logging
│   ├── mailer/          # Email sending functionality
│   │   └── templates/   # Email templates
│   ├── store/           # Database interaction and data models
│   │   ├── filters.go   # Filtering and pagination logic
│   │   ├── movies.go    # Movie data access
│   │   ├── permissions.go # User permissions management
│   │   ├── runtime.go   # Runtime configuration
│   │   ├── storage.go   # Database connection and setup
│   │   ├── tokens.go    # Token data access
│   │   └── users.go     # User data access
│   ├── validator/       # Data validation utilities
│   └── vcs/             # Version control system information
├── migrations/          # SQL migration files for database schema changes
├── vendor/              # Go module dependencies (vendored)
├── .air.toml            # Configuration for Air (live-reloading for Go apps)
├── .gitignore           # Git ignore rules
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Makefile             # Build and development automation scripts
└── README.md            # Project README
```

## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

