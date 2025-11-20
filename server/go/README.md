# Nexus Protocol Server

**Nexus Protocol Server** - высокопроизводительный AI-сервер, реализующий Nexus Application Protocol v1.1.0 с enterprise-функциональностью.

## 🚀 Быстрый старт

### Использование Docker Compose

```bash
# Из директории server/go
cd server/go

# Запуск всех сервисов
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f nexus-api
```

### Локальный запуск

```bash
# Установка зависимостей
go mod download

# Настройка переменных окружения
export NEXUS_DATABASE_USER=nexus
export NEXUS_DATABASE_PASSWORD=nexus_password
export NEXUS_AUTH_JWT_SECRET=your-secret-key
export NEXUS_AI_API_KEY=your-openai-key

# Запуск сервера
go run ./cmd/main.go
```

### Проверка работоспособности

```bash
# Health check
curl http://localhost:8080/health

# Readiness check (enterprise)
curl http://localhost:8080/ready

# Version info
curl http://localhost:8080/api/v1/version

# Metrics (Prometheus)
curl http://localhost:9090/metrics
```

## 🏗️ Архитектура

### Основные компоненты

1. **API Gateway** - HTTP REST + gRPC + WebSocket
2. **Context-Aware Templates** - AI-powered обработка запросов
3. **Batch Operations** - Параллельное выполнение
4. **Webhook Service** - Асинхронные уведомления
5. **IAM (Auth)** - Аутентификация и авторизация
6. **Analytics** - Метрики и аналитика
7. **Conversations** - AI беседы

### Внешние сервисы

- **PostgreSQL** - основная БД
- **Redis** - кэш и сессии
- **Prometheus** - метрики
- **Grafana** - визуализация

## 📋 API Endpoints

### Health Checks
- `GET /health` - Базовая проверка здоровья
- `GET /ready` - Детальная проверка готовности (enterprise)

### Аутентификация (Keycloak)
- `POST /api/v1/auth/register` - Регистрация в Keycloak
- `POST /api/v1/auth/login` - Вход через Keycloak
- `POST /api/v1/auth/refresh` - Обновление токена

**Keycloak URLs:**
- Admin Console: `http://localhost:8081`
- User Account: `http://localhost:8081/realms/nexus/account`
- OpenID Config: `http://localhost:8081/realms/nexus/.well-known/openid-connect-configuration`

### Templates (AI)
- `POST /api/v1/templates/execute` - Выполнение шаблона
- `GET /api/v1/templates/status/{id}` - Статус выполнения
- `GET /api/v1/templates/stream/{id}` - Поток результатов (SSE)

### Batch Operations
- `POST /api/v1/batch/execute` - Выполнение батча
- `GET /api/v1/batch/status/{id}` - Статус батча

### Webhooks
- `POST /api/v1/webhooks` - Регистрация webhook
- `GET /api/v1/webhooks` - Список webhooks
- `DELETE /api/v1/webhooks/{id}` - Удаление webhook

### Conversations
- `POST /api/v1/conversations` - Создание беседы
- `POST /api/v1/conversations/{id}/messages` - Отправка сообщения

### Analytics
- `POST /api/v1/analytics/events` - Логирование события
- `GET /api/v1/analytics/stats` - Получение статистики

## 🔐 Аутентификация через Keycloak

Сервер поддерживает два режима аутентификации:
- **Keycloak** (рекомендуется для production)
- **JWT** (локальная аутентификация)

### Настройка Keycloak

1. **Запуск Keycloak:**
```bash
make keycloak
```

2. **Доступ к админке:**
   - URL: `http://localhost:8081`
   - Логин: `admin`
   - Пароль: `admin`

3. **Realm "nexus"** автоматически импортируется из `docker/keycloak-realm.json`

4. **Создание пользователей:**
   - Через Keycloak Admin Console → Users
   - Или через API: `POST /api/v1/auth/register`

5. **Тестирование логина:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user@example.com","password":"password"}'
```

6. **Получение токена:**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 300,
  "login_url": "http://localhost:8081/realms/nexus/account"
}
```

7. **Использование токена:**
```bash
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
     http://localhost:8080/api/v1/templates/execute
```

## ⚙️ Конфигурация

### Переменные окружения

```bash
# Server
NEXUS_SERVER_PORT=8080

# Database
NEXUS_DATABASE_HOST=localhost
NEXUS_DATABASE_USER=nexus
NEXUS_DATABASE_PASSWORD=password

# Redis
NEXUS_REDIS_HOST=localhost
NEXUS_REDIS_PORT=6379

# Auth
NEXUS_AUTH_JWT_SECRET=your-secret

# AI
NEXUS_AI_PROVIDER=openai
NEXUS_AI_API_KEY=your-key

# Enterprise features
NEXUS_ENABLE_ENTERPRISE_FEATURES=true
NEXUS_RATE_LIMIT_ENABLED=true
NEXUS_CACHE_ENABLED=true
NEXUS_WEBHOOK_ENABLED=true
```

### Файл конфигурации

```yaml
# config/config.yaml
server:
  port: 8080

database:
  host: localhost
  user: nexus
  password: password

# ... остальные настройки
```

## 🏭 Enterprise возможности

### Response Metadata
```json
{
  "metadata": {
    "rate_limit_info": {
      "limit": 1000,
      "remaining": 950,
      "reset_at": 1640996100
    },
    "cache_info": {
      "cache_hit": true,
      "cache_ttl": 300
    },
    "quota_info": {
      "quota_used": 50000,
      "quota_limit": 100000
    }
  }
}
```

### Batch Operations
```json
{
  "operations": [
    {"type": "execute_template", "data": {...}},
    {"type": "log_event", "data": {...}}
  ],
  "options": {
    "parallel": true,
    "stop_on_error": false
  }
}
```

### Webhooks
```json
{
  "url": "https://app.example.com/webhooks",
  "events": ["template.completed", "batch.finished"],
  "secret": "webhook-secret",
  "retry_policy": {
    "max_retries": 3,
    "initial_delay": 1000
  }
}
```

## 📊 Мониторинг

### Метрики
- **Application**: requests, errors, latency
- **Business**: conversions, user engagement
- **System**: CPU, memory, connections

### Health Checks
```bash
# Basic health
curl http://localhost:8080/health
# {"status":"healthy","version":"1.1.0"}

# Enterprise readiness
curl http://localhost:8080/ready
# {
#   "status": "ready",
#   "components": {
#     "database": {"status": "healthy"},
#     "redis": {"status": "healthy"},
#     "ai_service": {"status": "healthy"}
#   },
#   "capacity": {
#     "current_load": 0.75,
#     "active_connections": 7500
#   }
# }
```

## 🔧 Разработка

### Структура проекта
```
server/go/
├── cmd/           # Точка входа
├── internal/      # Внутренняя логика
│   ├── api/       # HTTP handlers
│   ├── auth/      # Аутентификация
│   ├── ai/        # AI сервисы
│   └── ...
├── pkg/           # Переиспользуемые пакеты
│   ├── config/    # Конфигурация
│   ├── types/     # Типы данных
│   └── middleware/# HTTP middleware
├── docker/        # Docker файлы
├── config/        # Конфигурационные файлы
└── migrations/    # БД миграции
```

### Добавление нового endpoint

1. **Создать handler** в `internal/api/handlers/`
2. **Добавить маршрут** в `internal/api/router.go`
3. **Реализовать бизнес-логику** в соответствующем сервисе

### Тестирование

```bash
# Unit тесты
make test

# Integration тесты
make test-integration

# Benchmarks
go test -bench=. ./...
```

### End-to-End тестирование с Keycloak

```bash
# 1. Запуск сервисов
make docker-run

# 2. Создание тестового пользователя
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'

# 3. Логин
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test@example.com",
    "password": "password123"
  }' | jq -r .access_token)

# 4. Выполнение запроса с токеном
curl -X POST http://localhost:8080/api/v1/templates/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "хочу рецепт борща",
    "language": "ru",
    "metadata": {
      "request_id": "test-123",
      "client_version": "1.0.0"
    }
  }'
```

## 🚀 Deployment

### Docker Compose (разработка)
```bash
docker-compose up -d
```

### Kubernetes (production)
```bash
kubectl apply -f kubernetes/
```

### Cloud (AWS/GCP/Azure)
- Используйте managed PostgreSQL и Redis
- Настройте auto-scaling
- Включите monitoring

## 📞 Поддержка

- **Документация**: [Nexus Protocol](../../README.md)
- **API Specs**: [OpenAPI](../../api/rest/openapi.yaml)
- **Enterprise Guide**: [ADVANCED.md](../../sdk/go/ADVANCED.md)

---

**Nexus Protocol Server v1.1.0** - Production-ready AI платформа! 🚀
