# =============================================================================
# Weather API - Build Pipeline
# =============================================================================
# Organized build stages for backend, frontend, and integration
# =============================================================================

.PHONY: all help
.PHONY: backend-test backend-build backend-run
.PHONY: frontend-deps frontend-test frontend-build
.PHONY: docker-build docker-up docker-down
.PHONY: e2e-test test-all clean

# Default target
.DEFAULT_GOAL := help

# =============================================================================
# Backend Build Stages
# =============================================================================

## Stage 1: Backend - Run tests
backend-test:
	@echo "🔍 Running backend tests..."
	go test -v -race -cover ./...

## Stage 2: Backend - Build binary
backend-build: backend-test
	@echo "🔨 Building backend..."
	go build -o weather-api .

## Stage 3: Backend - Run locally
backend-run: backend-build
	@echo "🚀 Starting backend server..."
	./weather-api

# =============================================================================
# Frontend Build Stages
# =============================================================================

## Stage 1: Frontend - Install dependencies
frontend-deps:
	@echo "📦 Installing frontend dependencies with Bun..."
	cd frontend && bun install

## Stage 2: Frontend - Run unit tests (uses Node.js for compatibility)
frontend-test: frontend-deps
	@echo "🧪 Running frontend unit tests..."
	cd frontend && npx vitest run

## Stage 3: Frontend - Build production bundle
frontend-build: frontend-deps
	@echo "📦 Building frontend production bundle..."
	cd frontend && bun run build

## Stage 4: Frontend - Run dev server
frontend-dev: frontend-deps
	@echo "🚀 Starting frontend dev server..."
	cd frontend && bun run dev

# =============================================================================
# Integration & E2E Testing
# =============================================================================

## Install Playwright browsers
e2e-setup: frontend-deps
	@echo "🎭 Installing Playwright browsers..."
	cd frontend && bunx playwright install chromium

## Run E2E tests
e2e-test: e2e-setup backend-build frontend-build
	@echo "🎭 Running E2E tests..."
	cd frontend && bun run test:e2e

## Run E2E tests with UI
e2e-test-ui: e2e-setup
	@echo "🎭 Running E2E tests with UI..."
	cd frontend && bun run test:e2e:ui

# =============================================================================
# Full Pipeline
# =============================================================================

## Run all tests (backend + frontend)
test-all: backend-test frontend-test
	@echo "✅ All tests passed!"

## Full production build
build-all: clean backend-build frontend-build
	@echo "✅ Full production build complete!"

## Complete CI pipeline
ci: test-all build-all
	@echo "✅ CI pipeline complete!"

# =============================================================================
# Docker Operations
# =============================================================================

## Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t weather-api:latest .

## Start all services with docker-compose
docker-up:
	@echo "🐳 Starting Docker services..."
	docker-compose up -d

## Start with development frontend
docker-up-dev:
	@echo "🐳 Starting Docker services with dev frontend..."
	docker-compose --profile dev up -d

## Stop all services
docker-down:
	@echo "🛑 Stopping Docker services..."
	docker-compose down

## View logs
docker-logs:
	@echo "📋 Viewing Docker logs..."
	docker-compose logs -f

# =============================================================================
# Development Utilities
# =============================================================================

## Format Go code
fmt:
	@echo "📝 Formatting Go code..."
	go fmt ./...

## Tidy Go modules
tidy:
	@echo "🧹 Tidying Go modules..."
	go mod tidy

## Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f weather-api
	rm -f coverage.out coverage.html
	rm -f weather_cache.db test_*.db
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -f frontend/bun.lockb

## Security scan
security:
	@echo "🔒 Running security scan..."
	gosec ./...

# =============================================================================
# Help
# =============================================================================

help:
	@echo "Weather API - Build Pipeline"
	@echo ""
	@echo "Backend Stages:"
	@echo "  make backend-test     - Run backend tests"
	@echo "  make backend-build    - Build backend binary (runs tests first)"
	@echo "  make backend-run      - Run backend locally (builds first)"
	@echo ""
	@echo "Frontend Stages:"
	@echo "  make frontend-deps    - Install frontend dependencies (Bun)"
	@echo "  make frontend-test    - Run frontend unit tests"
	@echo "  make frontend-build   - Build frontend production bundle"
	@echo "  make frontend-dev     - Start frontend dev server"
	@echo ""
	@echo "Integration & E2E:"
	@echo "  make e2e-setup        - Install Playwright browsers"
	@echo "  make e2e-test         - Run E2E tests (requires backend + frontend built)"
	@echo "  make e2e-test-ui      - Run E2E tests with UI mode"
	@echo ""
	@echo "Full Pipeline:"
	@echo "  make test-all         - Run all tests (backend + frontend)"
	@echo "  make build-all        - Full production build"
	@echo "  make ci               - Complete CI pipeline"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-up        - Start services with docker-compose"
	@echo "  make docker-up-dev    - Start with dev frontend"
	@echo "  make docker-down      - Stop services"
	@echo "  make docker-logs      - View logs"
	@echo ""
	@echo "Utilities:"
	@echo "  make fmt              - Format Go code"
	@echo "  make tidy             - Tidy Go modules"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make security         - Run security scan"
	@echo "  make help             - Show this help message"