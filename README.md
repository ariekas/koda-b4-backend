# Project Name

> Brief description of what this service does and its purpose in the system architecture

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.9+-00ADD8?style=flat)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io/)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Overview

This is a production-grade RESTful API service built with Go (Golang) using the Gin web framework. The service provides [brief description of main functionality] and is designed with scalability, maintainability, and performance in mind.

## Tech Stack

- **Language**: Go 1.25+
- **Web Framework**: Gin v1.9+
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **ORM**: GORM / sqlx
- **Migration**: golang-migrate
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.25 or higher
- PostgreSQL 15 or higher
- Redis 7 or higher
- Docker & Docker Compose (optional, for containerized deployment)

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/your-organization/project-name.git
cd project-name
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration values.

## Running the Application

### Using Go directly

```bash
# Run database migrations (if using golang-migrate)
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up

# Start the server
go run cmd/api/main.go
```

### Using Docker Compose

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Development with Hot Reload (Optional)

If you want hot reload during development, install Air:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

The API will be available at `http://localhost:8080`

## API Documentation

Once the application is running, you can access the API documentation at:

- **Postman Documentasion**: https://.postman.co/workspace/My-Workspace~c2a1fd0c-0481-446f-80c2-62d586be5d52/collection/28747554-227eae38-7577-45dc-acb2-7d35267d91ab?action=share&creator=28747554&active-environment=28747554-34c6ebd3-99eb-4563-9e7f-08c9aa9a0ac4

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── controllers/             # HTTP handlers
│   ├── middlewares/             # Custom middlewares
│   ├── models/                  # Data models
│   ├── repositories/            # Data access layer
│   ├── routes/                  # Route definitions
|
├── migrations/                  # Database migrations
├── docs/                        # Swagger documentation
├── tests/                       # Test files
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker image definition
├── Makefile                     # Build and run commands
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── .env.example                 # Example environment variables
└── README.md                    # This file
```

## Testing

### Run all tests

```bash
go test ./...
```

### Run tests with verbose output

```bash
go test -v ./...
```

### Run tests with coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run specific test

```bash
go test -v ./internal/services/...
```

### Run tests for specific package

```bash
go test -v ./internal/controllers
```

## Deployment

### Building for Production

```bash
# Build binary
go build -o bin/app cmd/api/main.go

# Run the binary
./bin/app
```

### Build with specific OS/Architecture

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o bin/app-linux cmd/api/main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o bin/app.exe cmd/api/main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o bin/app-mac cmd/api/main.go
```

### Build Docker image

```bash
docker build -t project-name:latest .

# Run the Docker container
docker run -p 8080:8080 --env-file .env project-name:latest
```

## Database Migrations

### Using golang-migrate CLI

First, install golang-migrate:

```bash
# Install on macOS
brew install golang-migrate

# Install on Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Install on Windows (using scoop)
scoop install migrate
```

### Migration Commands

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_users_table

# Run all pending migrations
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" down 1

# Rollback all migrations
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" down

# Check migration version
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" version

# Force set version (use with caution)
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" force VERSION
```

### Alternative: Using GORM Auto-Migration

If you prefer GORM's auto-migration feature, add this to your main.go:

```go
db.AutoMigrate(&models.User{}, &models.Product{})
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code:
- Follows Go best practices and conventions
- Includes appropriate tests
- Updates documentation as needed
- Passes all CI checks