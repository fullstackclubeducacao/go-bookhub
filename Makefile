.PHONY: all build build-mongo build-all run run-mongo test test-integration test-all test-coverage test-coverage-all clean generate generate-api generate-sqlc generate-mocks migrate-up migrate-down migrate-create docker-build docker-build-mongo docker-build-all docker-run docker-run-mongo docker-compose docker-compose-down docker-postgres docker-mongodb lint deps help setup

# Variables
APP_NAME=bookhub
MAIN_FILE=cmd/api/main.go
MAIN_FILE_MONGO=cmd/api-mongo/main.go
BINARY_NAME=bin/$(APP_NAME)
BINARY_NAME_MONGO=bin/$(APP_NAME)-mongo
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_IMAGE_MONGO=$(APP_NAME)-mongo:latest

# Database - PostgreSQL
DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=postgres
DB_PASSWORD?=postgres
DB_NAME?=bookhub
DB_SSL_MODE?=disable
DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

# Database - MongoDB
MONGO_URI?=mongodb://mongo:mongo@localhost:27017
MONGO_DATABASE?=bookhub

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOMOD=$(GO) mod

# Tools
OAPI_CODEGEN=oapi-codegen
SQLC=sqlc
MIGRATE=migrate
MOCKGEN=mockgen

# Mocks directory
MOCKS_DIR=internal/mocks

all: generate build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all                  Generate code and build the application"
	@echo "  build                Build the PostgreSQL application binary"
	@echo "  build-mongo          Build the MongoDB application binary"
	@echo "  build-all            Build both PostgreSQL and MongoDB binaries"
	@echo "  run                  Run the PostgreSQL application"
	@echo "  run-mongo            Run the MongoDB application"
	@echo "  test                 Run unit tests"
	@echo "  test-integration     Run integration tests (requires Docker)"
	@echo "  test-all             Run all tests (unit + integration)"
	@echo "  test-coverage        Run unit tests with coverage report"
	@echo "  test-coverage-all    Run all tests with coverage report"
	@echo "  clean                Remove build artifacts"
	@echo "  generate             Generate all code (OpenAPI + SQLC)"
	@echo "  generate-api         Generate OpenAPI handlers with oapi-codegen"
	@echo "  generate-sqlc        Generate SQLC code"
	@echo "  generate-mocks       Generate mock files with mockgen"
	@echo "  migrate-up           Run database migrations"
	@echo "  migrate-down         Rollback database migrations"
	@echo "  migrate-create       Create a new migration file"
	@echo "  docker-build         Build PostgreSQL Docker image"
	@echo "  docker-build-mongo   Build MongoDB Docker image"
	@echo "  docker-build-all     Build both Docker images"
	@echo "  docker-run           Run PostgreSQL application in Docker"
	@echo "  docker-run-mongo     Run MongoDB application in Docker"
	@echo "  docker-compose       Run with docker compose"
	@echo "  docker-compose-down  Stop docker compose"
	@echo "  docker-postgres      Run PostgreSQL container only"
	@echo "  docker-mongodb       Run MongoDB container only"
	@echo "  lint                 Run linter"
	@echo "  deps                 Download dependencies"
	@echo "  setup                Initial project setup (tools, deps, code generation, .env)"

## build: Build the PostgreSQL application binary
build:
	@echo "Building $(APP_NAME) (PostgreSQL)..."
	@mkdir -p bin
	$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_NAME)"

## build-mongo: Build the MongoDB application binary
build-mongo:
	@echo "Building $(APP_NAME)-mongo (MongoDB)..."
	@mkdir -p bin
	$(GOBUILD) -ldflags="-s -w" -o $(BINARY_NAME_MONGO) $(MAIN_FILE_MONGO)
	@echo "Build complete: $(BINARY_NAME_MONGO)"

## build-all: Build both PostgreSQL and MongoDB binaries
build-all: build build-mongo

## run: Run the PostgreSQL application
run:
	@echo "Running $(APP_NAME) (PostgreSQL)..."
	$(GO) run $(MAIN_FILE)

## run-mongo: Run the MongoDB application
run-mongo:
	@echo "Running $(APP_NAME)-mongo (MongoDB)..."
	MONGO_URI=$(MONGO_URI) MONGO_DATABASE=$(MONGO_DATABASE) $(GO) run $(MAIN_FILE_MONGO)

## test: Run unit tests (excludes integration tests)
test:
	@echo "Running unit tests..."
	$(GOTEST) -v ./...

## test-integration: Run integration tests with testcontainers (requires Docker)
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./internal/infrastructure/repository/... -timeout 5m

## test-all: Run all tests (unit + integration)
test-all: test test-integration

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-coverage-all: Run all tests with coverage (unit + integration)
test-coverage-all:
	@echo "Running all tests with coverage..."
	$(GOTEST) -v -tags=integration -coverprofile=coverage.out ./... -timeout 5m
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## generate: Generate all code
generate: generate-api generate-sqlc

## generate-api: Generate OpenAPI handlers
generate-api:
	@echo "Generating OpenAPI handlers..."
	@mkdir -p api/generated
	$(OAPI_CODEGEN) -generate types,gin,spec -package generated -o api/generated/openapi.gen.go api/openapi/schema.yaml
	@echo "OpenAPI generation complete"

## generate-sqlc: Generate SQLC code
generate-sqlc:
	@echo "Generating SQLC code..."
	@cd internal/infrastructure/database/sqlc && $(SQLC) generate
	@echo "SQLC generation complete"

## generate-mocks: Generate mock files with mockgen
generate-mocks:
	@echo "Generating mocks..."
	@mkdir -p $(MOCKS_DIR)
	$(MOCKGEN) -source=internal/usecase/user_usecase.go -destination=$(MOCKS_DIR)/mock_user_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/usecase/book_usecase.go -destination=$(MOCKS_DIR)/mock_book_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/usecase/loan_usecase.go -destination=$(MOCKS_DIR)/mock_loan_usecase.go -package=mocks
	$(MOCKGEN) -source=internal/infrastructure/auth/jwt.go -destination=$(MOCKS_DIR)/mock_jwt_service.go -package=mocks
	@echo "Mocks generation complete"

## migrate-up: Run database migrations
migrate-up:
	@echo "Running migrations..."
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" up
	@echo "Migrations complete"

## migrate-down: Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down 1
	@echo "Rollback complete"

## migrate-create: Create a new migration (usage: make migrate-create name=migration_name)
migrate-create:
	@echo "Creating migration: $(name)"
	$(MIGRATE) create -ext sql -dir migrations -seq $(name)

## docker-build: Build PostgreSQL Docker image
docker-build:
	@echo "Building Docker image (PostgreSQL)..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

## docker-build-mongo: Build MongoDB Docker image
docker-build-mongo:
	@echo "Building Docker image (MongoDB)..."
	docker build -t $(DOCKER_IMAGE_MONGO) -f Dockerfile.mongo .
	@echo "Docker image built: $(DOCKER_IMAGE_MONGO)"

## docker-build-all: Build both Docker images
docker-build-all: docker-build docker-build-mongo

## docker-run: Run PostgreSQL application in Docker
docker-run:
	@echo "Running Docker container (PostgreSQL)..."
	docker run --rm -p 8080:8080 \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=$(DB_PORT) \
		-e DB_USER=$(DB_USER) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_NAME=$(DB_NAME) \
		$(DOCKER_IMAGE)

## docker-run-mongo: Run MongoDB application in Docker
docker-run-mongo:
	@echo "Running Docker container (MongoDB)..."
	docker run --rm -p 8080:8080 \
		-e MONGO_URI=mongodb://host.docker.internal:27017 \
		-e MONGO_DATABASE=$(MONGO_DATABASE) \
		$(DOCKER_IMAGE_MONGO)

## docker-compose: Run with docker compose
docker-compose:
	@echo "Starting with docker compose..."
	docker compose up --build

## docker-compose-down: Stop docker compose
docker-compose-down:
	@echo "Stopping docker compose..."
	docker compose down -v

## docker-postgres: Run PostgreSQL container only
docker-postgres:
	@echo "Starting PostgreSQL container..."
	docker compose up -d postgres
	@echo "PostgreSQL is running on port 5432"

## docker-mongodb: Run MongoDB container only
docker-mongodb:
	@echo "Starting MongoDB container..."
	docker compose up -d mongodb
	@echo "MongoDB is running on port 27017"

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies downloaded"

## setup: Initial project setup (installs tools, deps, generates code, creates .env)
setup:
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh
