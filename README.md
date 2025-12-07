# Just Sruput

>This service is a backend system built using Golang and the Gin Web Framework to manage core operations of a Coffee Shop application. It provides a set of secure and scalable RESTful APIs that handle essential business processes, including user authentication, product management, cart operations, order transactions, payment methods, ratings, and history tracking.

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
This is a production-grade RESTful API service built with Go (Golang) using the Gin web framework. The service provides a complete backend system for a Coffee Shop application, including user authentication, product management, shopping cart operations, order processing, payment methods, ratings, and transaction history. It is designed with scalability, maintainability, and high performance in mind, serving as the core API layer that connects the frontend clients with the database.

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

## ERD Diagram
```mermaid
erDiagram

    profile ||--|| users : "profile_id"
    users ||--o{ carts : "users_id"
    users ||--o{ transactions : "users_id"
    users ||--o{ ratings : "users_id"

    category_products ||--o{ products : "category_products_id"
    discounts ||--o{ products : "discounts_id"

    products ||--o{ product_images : "products_id"
    products ||--o{ product_sizes : "products_id"
    products ||--o{ product_variants : "products_id"
    products ||--o{ carts : "products_id"
    products ||--o{ transaction_items : "products_id"
    products ||--o{ ratings : "products_id"

    size_products ||--o{ product_sizes : "size_products_id"
    variant_products ||--o{ product_variants : "variant_products_id"

    product_sizes ||--o{ carts : "product_sizes_id"
    product_sizes ||--o{ transaction_items : "product_sizes_id"

    product_variants ||--o{ carts : "product_variants_id"
    product_variants ||--o{ transaction_items : "product_variants_id"

    deliverys ||--o{ transactions : "deliverys_id"
    payment_methods ||--o{ transactions : "payment_methods_id"
    status_transactions ||--o{ transactions : "status_transactions_id"
    taxes ||--o{ transactions : "taxes_id"

    transactions ||--o{ transaction_items : "transactions_id"
    transactions ||--o{ ratings : "transactions_id"

    profile {
        int id
        text pic
        varchar phone
        text address
        timestamp created_at
        timestamp updated_at
    }

    users {
        int id
        varchar fullname
        varchar email
        varchar password
        varchar role
        int profile_id
        timestamp created_at
        timestamp updated_at
    }

    category_products {
        int id
        varchar name
        timestamp created_at
        timestamp updated_at
    }

    discounts {
        int id
        varchar name
        decimal discount_percentage
        timestamp created_at
        timestamp updated_at
    }

    products {
        int id
        int discounts_id
        varchar name
        decimal price
        text description
        int stock
        boolean is_flashsale
        boolean is_favorite_product
        int category_products_id
        int ratings
        timestamp created_at
        timestamp updated_at
    }

    product_images {
        int id
        text image
        int products_id
        boolean is_primary
        timestamp created_at
        timestamp updated_at
    }

    size_products {
        int id
        varchar name
        decimal additional_costs
        timestamp created_at
        timestamp updated_at
    }

    variant_products {
        int id
        varchar name
        decimal additional_costs
        timestamp created_at
        timestamp updated_at
    }

    product_sizes {
        int id
        int products_id
        int size_products_id
        timestamp created_at
        timestamp updated_at
    }

    product_variants {
        int id
        int products_id
        int variant_products_id
        timestamp created_at
        timestamp updated_at
    }

    taxes {
        int id
        varchar name
        float tax
        timestamp created_at
        timestamp updated_at
    }

    carts {
        int id
        int users_id
        int products_id
        int product_sizes_id
        int product_variants_id
        int quantity
        timestamp created_at
        timestamp updated_at
    }

    deliverys {
        int id
        varchar name
        decimal price
        timestamp created_at
        timestamp updated_at
    }

    payment_methods {
        int id
        varchar name
        text image_payment
        timestamp created_at
        timestamp updated_at
    }

    status_transactions {
        int id
        varchar status
        timestamp created_at
        timestamp updated_at
    }

    transactions {
        int id
        int users_id
        int deliverys_id
        int payment_methods_id
        int status_transactions_id
        int taxes_id
        varchar name_user
        text address_user
        varchar phone_user
        varchar email_user
        decimal subtotal
        decimal tax_amount
        decimal total
        varchar invoice_num
        timestamp created_at
        timestamp updated_at
    }

    transaction_items {
        int id
        int transactions_id
        int products_id
        int quantity
        decimal price_at_time
        decimal discount_at_time
        int product_sizes_id
        int product_variants_id
        timestamp created_at
        timestamp updated_at
    }

    ratings {
        int id
        int users_id
        int products_id
        int transactions_id
        int rating
        text review
        timestamp created_at
        timestamp updated_at
    }

```
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