# Ayo Cuci Backend (Clean Architecture)

Backend system for Ayo Cuci laundry management app, built with Golang and Gin Framework.

## Tech Stack

- **Language:** Go (1.22+)
- **Framework:** Gin Gonic
- **ORM:** GORM
- **Database:** SQLite (Development) / PostgreSQL (Production ready)
- **Auth:** JWT (JSON Web Token)

## Architecture

This project follows **Clean Architecture** principles:

- `cmd/`: Application entry points.
- `internal/module/`: Domain logic split by modules (Auth, Outlet, Employee, etc.).
- `middleware/`: Custom Gin middlewares (Auth, Role check).

## How to Run

1. Clone the repository.
2. Run `go mod download`.
3. Set up your `.env` file based on environment variables.
4. Run `go run ./cmd/server/main.go`.
