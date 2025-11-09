.PHONY: run docker-up docker-down docker-logs test build clean help

# Variables
APP_NAME=geekery-api
DOCKER_COMPOSE=cd deploy && docker-compose

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "╔═══════════════════════════════════════════════════════════╗"
	@echo "║              GEEKERY API - MAKEFILE COMMANDS              ║"
	@echo "╚═══════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Quick Start:"
	@echo "  make setup    - Complete initial setup (first time only)"
	@echo "  make dev      - Start everything (PostgreSQL + API)"
	@echo "  make run      - Run the API (PostgreSQL must be running)"
	@echo ""
	@echo "Documentation:"
	@echo "  make swagger       - Generate Swagger docs"
	@echo "  make open-swagger  - Open Swagger UI in browser"
	@echo "  📘 Swagger URL: http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "Database:"
	@echo "  make db-reset - Reset database (removes all data)"
	@echo "  make db-verify - Test database connection"
	@echo ""
	@echo "All commands:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | column -t -s ':'
	@echo ""

## run: Run the application
run:
	@echo "🚀 Starting Geekery API..."
	@go run cmd/main.go

## build: Build the application binary
build:
	@echo "🔨 Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/main.go
	@echo "✓ Build complete: bin/$(APP_NAME)"

## docker-up: Start PostgreSQL container
docker-up:
	@echo "🐳 Starting PostgreSQL container..."
	@$(DOCKER_COMPOSE) up -d
	@echo "✓ PostgreSQL is running on port 5432"

## docker-down: Stop PostgreSQL container
docker-down:
	@echo "🛑 Stopping PostgreSQL container..."
	@$(DOCKER_COMPOSE) down
	@echo "✓ PostgreSQL stopped"

## docker-logs: View PostgreSQL logs
docker-logs:
	@$(DOCKER_COMPOSE) logs -f postgres

## docker-restart: Restart PostgreSQL container
docker-restart: docker-down docker-up
	@echo "✓ PostgreSQL restarted"

## test: Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

## deps: Download Go module dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod verify
	@echo "✓ Dependencies downloaded"

## tidy: Tidy Go module dependencies
tidy:
	@echo "🔧 Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

## fmt: Format Go code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
	fi

## dev: Start development environment (PostgreSQL + API)
dev: docker-up
	@$(MAKE) db-setup
	@$(MAKE) run

## prod: Build and run production binary
prod: build
	@echo "🚀 Starting production server..."
	@./bin/$(APP_NAME)

## db-setup: Setup database user and grants
db-setup:
	@echo "🗄️  Setting up database..."
	@sleep 2
	@echo "✓ Database is already configured via Docker environment variables"
	@echo "  - User: geekery"
	@echo "  - Database: geekery_db"
	@echo "  - Password: configured in .env"

## db-reset: Reset database (removes volume and recreates everything)
db-reset:
	@echo "🔄 Resetting database (this will delete all data)..."
	@$(MAKE) docker-down
	@docker volume rm deploy_postgres_data 2>/dev/null || echo "Volume already removed or doesn't exist"
	@echo "✓ Old data removed"
	@$(MAKE) docker-up
	@echo "⏳ Waiting for PostgreSQL to initialize..."
	@sleep 5
	@docker exec geekery-postgres psql -U geekery -d geekery_db -c "SELECT 'Database ready!' as status;" && \
		echo "✓ Database reset complete!" || \
		echo "⚠️  Database is still initializing, wait a few seconds and run 'make run'"

## db-verify: Verify database connection
db-verify:
	@echo "🔍 Verifying database connection..."
	@docker exec geekery-postgres psql -U geekery -d geekery_db -c "SELECT current_database(), current_user, version();" && \
		echo "✓ Database is accessible!" || \
		echo "❌ Cannot connect to database. Run 'make db-reset' to fix."

## setup: Initial setup (install deps, start db, configure database)
setup:
	@echo "⚙️  Setting up Geekery API..."
	@$(MAKE) deps
	@$(MAKE) db-reset
	@echo "✓ Setup complete! Run 'make run' to start the API"

## seed: Seed database with sample data
seed:
	@echo "🌱 Seeding database with sample data..."
	@go run cmd/seed/main.go
	@echo "✓ Seeding complete!"

## swagger: Generate Swagger documentation
swagger:
	@echo "📚 Generating Swagger documentation..."
	@swag init -g cmd/main.go -o docs --parseDependency --parseInternal
	@echo "✓ Swagger docs generated! Access at http://localhost:8080/swagger/index.html"

## open-swagger: Open Swagger UI
open-swagger:
	@echo "🔗 Opening Swagger UI..."
	@open http://localhost:8080/swagger/index.html
	@echo "✓ Swagger UI opened"
