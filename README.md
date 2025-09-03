# Go REST API Server With Echo and MongoDB

This project is a **Base API** server written in Go, using the [Echo](https://echo.labstack.com/) web framework and MongoDB as the database. It is designed as a starting point for new Go REST API projects, providing a clean structure and best practices for scalable API development.

## Features

- RESTful API structure
- Echo web framework
- MongoDB integration
- Modular project layout (handlers, middleware, models, repository, routes)
- Configurable via YAML
- Built-in logger
- Hot reload with [Air](https://github.com/air-verse/air)
- Standardized API response template for all endpoints
- Request validation for incoming data
- Basic authentication using JWT
- Authorization with user roles

## Package Versions

- **Go**: see `go.mod` (recommended: Go 1.25+)
- **Echo**: v4.x ([github.com/labstack/echo/v4](https://github.com/labstack/echo))
- **MongoDB Go Driver**: v1.17.x ([go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver))
- **Air** (for live reload): v1.62.0 ([air-verse/air](https://github.com/air-verse/air))

See `go.mod` for the full list of dependencies and their versions.

## Project Structure

```
go.mod                  # Go module definition
go.sum                  # Go module checksums
cmd/
  api/
    main.go             # Application entry point
config/
  config.go             # Configuration loader (reads config file and env)
  config.example.yaml   # Template configuration file (copy to config.yaml for setup)
internal/
  auth/
    auth.go             # Authentication logic (login, token generation, password hashing)
    middleware.go       # Auth-related middleware (JWT validation, role checks)
  handlers/
    handlers.go         # General handlers (base handler functions)
    user_handler.go     # User-related handlers (user endpoints: CRUD, profile)
    auth_handler.go     # Auth endpoints (login, register, refresh token)
    role_handler.go     # Role endpoints (role management)
  middleware/
    middleware.go       # Custom middleware (request logging, error handling, CORS, etc.)
  models/
    auth.go             # Auth-related data models (JWT claims, login/register structs)
    role.go             # Role data model (role struct, permissions)
    user.go             # User data model (user struct, validation)
  repository/
    repository.go       # Repository interfaces (data access abstraction)
    mongo/
      auth_repo.go      # MongoDB auth repository implementation (login, register)
      role_repo.go      # MongoDB role repository implementation (role CRUD)
      user_repo.go      # MongoDB user repository implementation (user CRUD)
  routes/
    routes.go           # Route definitions and registration (Echo router)
pkg/
  logger/
    logger.go           # Logger utility (structured logging)
  response/
    response.go         # Standardized API response template (success/error responses)
  validation/
    validation.go       # Request validation logic (struct validation, custom rules)
```

## Getting Started

### 1. Clone the repository

```sh
git clone https://github.com/madhiyono/go-echo_base-api-no-sql.git
cd go-echo_base-api-no-sql
```

### 2. Configure the application

Rename `config/config.example.yaml` to `config/config.yaml`
Edit `config/config.yaml` to set your MongoDB URI and other settings.

### 3. Install dependencies

```sh
go mod tidy
```

### 4. Run the server with Air (hot reload)

Install Air if you haven't:

```sh
go install github.com/air-verse/air@latest
```

Run The Air Init:

```sh
air init
```

Change `.air.toml` Pointing the cmd:

```sh
cmd = "go build -o ./tmp/main.exe cmd/api/main.go"
```

Save and Run the server:

```sh
air
```

Or, run directly:

```sh
go run cmd/api/main.go
```

## How to Add New Features

### Add a New Route

1. Create a new handler in `internal/handlers/` (e.g., `product_handler.go`).
2. Register the route in `internal/routes/routes.go` using Echo's router.
3. Implement any required business logic or validation in the handler.

### Add a New Repository

1. Create a new repository file in `internal/repository/` or `internal/repository/mongo/` (e.g., `product_repo.go`).
2. Define the interface and implementation for data access.
3. Inject the repository into your handler as needed.

### Add a New Model

1. Define your struct in `internal/models/` (e.g., `product.go`).
2. Use the model in your handler and repository.

### Add Middleware

1. Create a new middleware in `internal/middleware/`.
2. Register it in `cmd/api/main.go` or in the router as needed.

---

This project is intended as a base for new Go REST API projects. Feel free to fork and adapt for your needs.

## Updates

### New Features

- Added new features (see CHANGELOG.md for details)

### Changes

- Various improvements and changes have been made to enhance stability and usability.

For a complete list of updates, see [CHANGELOG.md](./CHANGELOG.md).
