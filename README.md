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

## Package Versions

- **Go**: see `go.mod` (recommended: Go 1.25+)
- **Echo**: v4.x ([github.com/labstack/echo/v4](https://github.com/labstack/echo))
- **MongoDB Go Driver**: v1.17.x ([go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver))
- **Air** (for live reload): v1.62.0 ([cosmtrek/air](https://github.com/air-verse/air))

See `go.mod` for the full list of dependencies and their versions.

## Project Structure

```
cmd/api/main.go         # Entry point
config/                 # Configuration files
internal/
  handlers/             # HTTP handlers
  middleware/           # Custom middleware
  models/               # Data models
  repository/           # Data access layer
    mongo/              # MongoDB-specific repositories
  routes/               # Route definitions
pkg/logger/             # Logger utility
```

## Getting Started

### 1. Clone the repository

```sh
git clone <your-repo-url>
cd base-api-nosql
```

### 2. Configure the application

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

## License

MIT

---

This project is intended as a base for new Go REST API projects. Feel free to fork and adapt for your needs.
