#!/bin/bash

# BookHub Setup Script
# This script sets up the development environment for the BookHub project

set -e

echo "=========================================="
echo "  BookHub - Development Setup"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print status
print_status() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check Go installation
echo ""
echo "Checking Go installation..."
if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    print_status "Go is installed: $GO_VERSION"
else
    print_error "Go is not installed. Please install Go 1.21+ from https://golang.org/dl/"
    exit 1
fi

# Check PostgreSQL
echo ""
echo "Checking PostgreSQL..."
if command_exists psql; then
    PSQL_VERSION=$(psql --version | awk '{print $3}')
    print_status "PostgreSQL client is installed: $PSQL_VERSION"
else
    print_warning "PostgreSQL client (psql) is not installed"
    print_warning "You can still use Docker for PostgreSQL"
fi

# Check MongoDB
echo ""
echo "Checking MongoDB..."
if command_exists mongosh; then
    MONGO_VERSION=$(mongosh --version | head -1)
    print_status "MongoDB Shell is installed: $MONGO_VERSION"
else
    print_warning "MongoDB Shell (mongosh) is not installed"
    print_warning "You can still use Docker for MongoDB"
fi

# Install Go tools
echo ""
echo "Installing Go tools..."

echo "  - Installing oapi-codegen..."
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
print_status "oapi-codegen installed"

echo "  - Installing sqlc..."
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
print_status "sqlc installed"

echo "  - Installing migrate..."
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
print_status "migrate installed"

echo "  - Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
print_status "golangci-lint installed"

echo "  - Installing mockgen (uber-go/mock)..."
go install go.uber.org/mock/mockgen@latest
print_status "mockgen installed"

# Download dependencies
echo ""
echo "Downloading Go dependencies..."
go mod download
go mod tidy
print_status "Dependencies downloaded"

# Generate code
echo ""
echo "Generating code..."

echo "  - Generating OpenAPI handlers..."
mkdir -p api/generated
oapi-codegen -generate types,gin,spec -package generated -o api/generated/openapi.gen.go api/openapi/openapi.yaml
print_status "OpenAPI handlers generated"

echo "  - Generating SQLC code..."
cd internal/infrastructure/database/sqlc && sqlc generate && cd -
print_status "SQLC code generated"

echo "  - Generating mocks..."
mkdir -p internal/mocks
mockgen -source=internal/usecase/user_usecase.go -destination=internal/mocks/mock_user_usecase.go -package=mocks
mockgen -source=internal/usecase/book_usecase.go -destination=internal/mocks/mock_book_usecase.go -package=mocks
mockgen -source=internal/usecase/loan_usecase.go -destination=internal/mocks/mock_loan_usecase.go -package=mocks
mockgen -source=internal/infrastructure/auth/jwt.go -destination=internal/mocks/mock_jwt_service.go -package=mocks
print_status "Mocks generated"

# Swagger UI info
echo ""
echo "Swagger UI..."
print_status "Swagger UI is embedded via github.com/flowchartsman/swaggerui (no manual download needed)"
print_status "Access Swagger UI at http://localhost:8080/docs after starting the server"

# Create .env file if not exists
echo ""
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cat > .env << 'EOF'
# Server Configuration
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

# PostgreSQL Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bookhub
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_LIFETIME=5m

# MongoDB Database Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=bookhub
MONGO_MAX_POOL_SIZE=100
MONGO_MIN_POOL_SIZE=10
MONGO_MAX_IDLE_TIME=5m

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-key-change-in-production
JWT_TOKEN_DURATION=24h
JWT_ISSUER=bookhub
EOF
    print_status ".env file created"
else
    print_warning ".env file already exists, skipping"
fi

# Summary
echo ""
echo "=========================================="
echo "  Setup Complete!"
echo "=========================================="
echo ""
echo "Next steps for PostgreSQL:"
echo "  1. Start PostgreSQL (or use: docker-compose up -d postgres)"
echo "  2. Create database: createdb bookhub"
echo "  3. Run migrations: make migrate-up"
echo "  4. Start the server: make run"
echo ""
echo "Next steps for MongoDB:"
echo "  1. Start MongoDB (or use: docker-compose up -d mongodb)"
echo "  2. Start the server: make run-mongo"
echo ""
echo "Available commands:"
echo "  make help              - Show all available commands"
echo "  make run               - Run the PostgreSQL application"
echo "  make run-mongo         - Run the MongoDB application"
echo "  make build-all         - Build both PostgreSQL and MongoDB binaries"
echo "  make test              - Run tests"
echo "  make docker-compose    - Run with Docker Compose (both APIs)"
echo "  make docker-build-all  - Build both Docker images"
echo ""
