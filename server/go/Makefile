.PHONY: build run test clean docker-build docker-run deps migrate

# Переменные
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker-compose
BINARY_NAME=nexus-server
VERSION=1.1.0

# Пути
CMD_DIR=./cmd
BUILD_DIR=./build
CONFIG_DIR=./config

all: deps build

# Установка зависимостей
deps:
	@echo "📦 Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Сборка приложения
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go


# Запуск в development режиме
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	$(GO) run $(CMD_DIR)/main.go


# Запуск с Docker Compose
docker-run:
	@echo "🐳 Starting with Docker Compose..."
	$(DOCKER_COMPOSE) up -d

# Остановка Docker Compose
docker-stop:
	@echo "🛑 Stopping Docker Compose..."
	$(DOCKER_COMPOSE) down

# Сборка Docker образа
docker-build:
	@echo "🏗️ Building Docker image..."
	$(DOCKER) build -f docker/Dockerfile -t nexus-protocol/server:$(VERSION) .

# Тестирование
test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...

# Integration тесты
test-integration:
	@echo "🔗 Running integration tests..."
	$(GO) test -tags=integration -v ./...

# Benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Линтинг
lint:
	@echo "🔍 Running linter..."
	golangci-lint run ./...

# Форматирование кода
fmt:
	@echo "💅 Formatting code..."
	$(GO) fmt ./...

# Валидация кода
vet:
	@echo "🔬 Running go vet..."
	$(GO) vet ./...

# Очистка
clean:
	@echo "🧹 Cleaning up..."
	$(GO) clean
	rm -rf $(BUILD_DIR)
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Миграции БД
migrate-up:
	@echo "🗃️ Running database migrations..."
	# Implement migration logic

migrate-down:
	@echo "⬇️ Rolling back database migrations..."
	# Implement rollback logic

# Генерация кода (protobuf, etc.)
generate:
	@echo "🤖 Generating code..."
	$(GO) generate ./...

# Health check
health:
	@echo "❤️ Checking health..."
	curl -f http://localhost:8080/health || echo "Health check failed"

# Logs
logs:
	@echo "📋 Showing logs..."
	$(DOCKER_COMPOSE) logs -f nexus-api

# Development setup
dev-setup: deps
	@echo "🚀 Setting up development environment..."
	# Create config file if not exists
	@if [ ! -f $(CONFIG_DIR)/config.yaml ]; then \
		cp config/config.yaml.example $(CONFIG_DIR)/config.yaml; \
		echo "⚠️ Please update $(CONFIG_DIR)/config.yaml with your settings"; \
	fi

# Start Keycloak only
keycloak:
	@echo "🔐 Starting Keycloak..."
	$(DOCKER_COMPOSE) up keycloak -d
	@echo "Keycloak Admin: http://localhost:8081 (admin/admin)"

# Start all services except Keycloak
services:
	@echo "🚀 Starting services (without Keycloak)..."
	$(DOCKER_COMPOSE) up postgres redis nexus-api -d

# Start everything
all-services: docker-run

# Stop Keycloak
stop-keycloak:
	@echo "🛑 Stopping Keycloak..."
	$(DOCKER_COMPOSE) stop keycloak

# View Keycloak logs
logs-keycloak:
	@echo "📋 Keycloak logs..."
	$(DOCKER_COMPOSE) logs -f keycloak

# Production build
prod-build: clean deps fmt vet lint test build

# Help
help:
	@echo "🤖 Available commands:"
	@echo "  make deps          - Install dependencies"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run in development mode"
	@echo "  make docker-run    - Run all services with Docker Compose"
	@echo "  make keycloak      - Start Keycloak only"
	@echo "  make services      - Start services without Keycloak"
	@echo "  make docker-stop   - Stop Docker Compose"
	@echo "  make stop-keycloak - Stop Keycloak"
	@echo "  make test          - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make clean         - Clean up build artifacts"
	@echo "  make health        - Check service health"
	@echo "  make logs          - Show all logs"
	@echo "  make logs-keycloak - Show Keycloak logs"
	@echo "  make dev-setup     - Setup development environment"
	@echo "  make prod-build    - Build for production"
	@echo "  make help          - Show this help"

# Default target
.DEFAULT_GOAL := help
