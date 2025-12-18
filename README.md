# Device API

A REST API for managing device resources, built with Go, Gin, and Gorm.

## Architecture

This project follows DDD (Domain-Driven Design) and Clean Architecture principles.

- **Domain Layer**: `internal/domain` (Entities, Interfaces)
- **Service Layer**: `internal/service` (Business Logic)
- **Repository Layer**: `internal/repository` (Data Access)
- **Handler Layer**: `internal/handler` (HTTP Transport)

## Prerequisites

- Go 1.23+
- Docker & Docker Compose

## Configuration

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```
2. Update the values in `.env` as needed.

## Running the Application

### Using Docker Compose (Recommended)

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

## API Documentation

The API documentation is available via:

- **Swagger UI (Internal)**: `http://localhost:8080/swagger/index.html`
- **Swagger UI (Standalone)**: `http://localhost:8081`
- **Postman Collection**: `postman_collection.json` (Import into Postman)

### Local Development

1. Set `DATABASE_URL` environment variable.
2. Run `go run cmd/api/main.go`.

## API Endpoints

- `POST /api/v1/devices`: Create a new device.
- `GET /api/v1/devices/:id`: Get a device by ID.
- `GET /api/v1/devices`: List all devices (supports `?brand=X` or `?state=Y` filters).
- `PUT/PATCH /api/v1/devices/:id`: Update a device (details or state).
- `DELETE /api/v1/devices/:id`: Delete a device.

## Testing

Run unit and integration tests:

```bash
go test ./...
```
