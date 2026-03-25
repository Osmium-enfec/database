.PHONY: help build run test docker-build docker-up docker-down db-migrate db-reset clean lint

help:
	@echo "Content Review API - Development Commands"
	@echo "=========================================="
	@echo "make build          - Build the application"
	@echo "make run            - Run the application"
	@echo "make test           - Run tests"
	@echo "make lint           - Run linter"
	@echo "make docker-build   - Build Docker image"
	@echo "make docker-up      - Start Docker containers"
	@echo "make docker-down    - Stop Docker containers"
	@echo "make docker-logs    - View Docker logs"
	@echo "make db-migrate     - Run database migrations"
	@echo "make db-reset       - Reset database"
	@echo "make clean          - Clean build artifacts"
	@echo "make fmt            - Format code"
	@echo "make deps           - Download dependencies"

# Build commands
build:
	@echo "Building application..."
	go build -o bin/app main.go

run: build
	@echo "Running application..."
	./bin/app

# Testing commands
test:
	@echo "Running tests..."
	go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Code quality commands
lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t content-review-api:latest .

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f api

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker rmi content-review-api:latest

# Database commands
db-migrate:
	@echo "Running database migrations..."
	psql -h localhost -U postgres -d content_review -f migrations/001_initial_schema.sql

db-reset:
	@echo "Resetting database..."
	dropdb -h localhost -U postgres content_review || true
	createdb -h localhost -U postgres content_review
	psql -h localhost -U postgres -d content_review -f migrations/001_initial_schema.sql

db-backup:
	@echo "Backing up database..."
	pg_dump -h localhost -U postgres content_review > backup_$$(date +%Y%m%d_%H%M%S).sql

# Dependency commands
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Cleanup commands
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean
	rm -f coverage.out

# Development setup
setup:
	@echo "Setting up development environment..."
	cp .env.example .env
	go mod download
	go mod tidy
	@echo "Setup complete! Please configure .env file"

# Run all checks before commit
pre-commit: fmt lint test
	@echo "Pre-commit checks passed!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init

.DEFAULT_GOAL := help
