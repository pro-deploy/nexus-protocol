.PHONY: build test examples clean deps fmt vet lint

# Переменные
GO=go
GOFMT=gofmt
GOLINT=golangci-lint
BINARY_NAME=nexus-sdk

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

all: deps fmt vet test build

# Установка зависимостей
deps:
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

# Форматирование кода
fmt:
	@echo "$(GREEN)Форматирование кода...$(NC)"
	$(GOFMT) -s -w .

# Проверка кода
vet:
	@echo "$(GREEN)Проверка кода...$(NC)"
	$(GO) vet ./...

# Линтинг
lint:
	@echo "$(GREEN)Линтинг кода...$(NC)"
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
	else \
		echo "$(YELLOW)golangci-lint не установлен, пропускаем...$(NC)"; \
	fi

# Тестирование
test:
	@echo "$(GREEN)Запуск тестов...$(NC)"
	$(GO) test -v ./...

# Integration тесты (требуют NEXUS_API_URL и NEXUS_API_TOKEN)
test-integration:
	@echo "$(GREEN)Запуск integration тестов...$(NC)"
	$(GO) test -tags=integration -v ./client/...

# Сборка примеров
examples:
	@echo "$(GREEN)Сборка примеров...$(NC)"
	$(GO) build -o examples/basic/basic ./examples/basic
	$(GO) build -o examples/error_handling/error_handling ./examples/error_handling
	$(GO) build -o examples/iam/iam ./examples/iam
	$(GO) build -o examples/conversations/conversations ./examples/conversations
	$(GO) build -o examples/analytics/analytics ./examples/analytics
	$(GO) build -o examples/retry/retry ./examples/retry
	$(GO) build -o examples/interceptors/interceptors ./examples/interceptors
	$(GO) build -o examples/metrics/metrics ./examples/metrics

# Запуск базового примера
run-basic:
	@echo "$(GREEN)Запуск базового примера...$(NC)"
	$(GO) run ./examples/basic

# Запуск примера обработки ошибок
run-error:
	@echo "$(GREEN)Запуск примера обработки ошибок...$(NC)"
	$(GO) run ./examples/error_handling

# Запуск примера IAM
run-iam:
	@echo "$(GREEN)Запуск примера IAM...$(NC)"
	$(GO) run ./examples/iam

# Запуск примера Conversations
run-conversations:
	@echo "$(GREEN)Запуск примера Conversations...$(NC)"
	$(GO) run ./examples/conversations

# Запуск примера Analytics
run-analytics:
	@echo "$(GREEN)Запуск примера Analytics...$(NC)"
	$(GO) run ./examples/analytics

# Запуск примера Retry
run-retry:
	@echo "$(GREEN)Запуск примера Retry...$(NC)"
	$(GO) run ./examples/retry

# Запуск примера Interceptors
run-interceptors:
	@echo "$(GREEN)Запуск примера Interceptors...$(NC)"
	$(GO) run ./examples/interceptors

# Запуск примера Metrics
run-metrics:
	@echo "$(GREEN)Запуск примера Metrics...$(NC)"
	$(GO) run ./examples/metrics

# Очистка
clean:
	@echo "$(GREEN)Очистка...$(NC)"
	$(GO) clean
	rm -f examples/basic/basic
	rm -f examples/error_handling/error_handling

# Помощь
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "  make deps      - Установить зависимости"
	@echo "  make fmt       - Форматировать код"
	@echo "  make vet        - Проверить код"
	@echo "  make lint      - Запустить линтер"
	@echo "  make test      - Запустить тесты"
	@echo "  make examples  - Собрать примеры"
	@echo "  make run-basic - Запустить базовый пример"
	@echo "  make run-error - Запустить пример обработки ошибок"
	@echo "  make clean     - Очистить артефакты сборки"
	@echo "  make help      - Показать эту справку"

