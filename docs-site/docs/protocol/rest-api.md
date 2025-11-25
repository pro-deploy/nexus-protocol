---
id: rest-api
title: REST API
sidebar_label: REST API
---

# REST API

Nexus Protocol поддерживает HTTP REST API для взаимодействия с сервером через стандартные HTTP методы и JSON.

## 🌐 Базовая информация

### Базовый URL
```
https://api.nexus.dev/api/v1
```

### Поддерживаемые версии
- **Protocol Version**: 2.0.0
- **API Version**: v1

### Формат данных
- **Content-Type**: `application/json`
- **Encoding**: UTF-8

## 🔐 Аутентификация

Все API запросы требуют аутентификации через JWT токен:

```
Authorization: Bearer <jwt_token>
```

### Получение токена

```bash
POST /auth/login
Content-Type: application/json

{
  "username": "your-username",
  "password": "your-password"
}
```

**Ответ:**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

## 📋 Формат сообщений

### Структура запроса

Все запросы следуют единому формату Nexus Protocol:

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "web-app",
    "client_type": "web",
    "timestamp": 1640995200,
    "custom_headers": {
      "x-request-priority": "high"
    }
  },
  "data": {
    // Операция-специфичные данные
  }
}
```

### Структура ответа

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "timestamp": 1640995235,
    "processing_time_ms": 3500,
    "rate_limit_info": {
      "limit": 1000,
      "remaining": 999,
      "reset_at": 1640996100
    }
  },
  "data": {
    // Результат операции
  }
}
```

## 🚀 Основные эндпоинты

### Health Check

#### GET /health
Простая проверка доступности сервиса.

**Пример:**
```bash
curl -X GET https://api.nexus.dev/api/v1/health
```

**Ответ:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-18T10:00:00Z",
  "version": "2.0.0"
}
```

#### GET /ready
Детальная проверка готовности (для Kubernetes).

**Ответ:**
```json
{
  "status": "ready",
  "timestamp": "2025-01-18T10:00:00Z",
  "checks": {
    "database": "ok",
    "redis": "ok",
    "ai_services": "ok"
  }
}
```

### Основные операции

#### POST /templates/execute
Выполнение шаблона с AI.

**Пример запроса:**
```bash
curl -X POST https://api.nexus.dev/api/v1/templates/execute \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "request_id": "req-123",
      "protocol_version": "2.0.0",
      "client_version": "2.0.0"
    },
    "data": {
      "query": "хочу борщ",
      "language": "ru",
      "context": {
        "user_id": "user-123",
        "location": {
          "latitude": 55.7558,
          "longitude": 37.6173
        },
        "locale": "ru-RU",
        "currency": "RUB"
      },
      "options": {
        "timeout_ms": 30000,
        "max_results_per_domain": 5
      }
    }
  }'
```

**Ответ:**
```json
{
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 3500
  },
  "data": {
    "execution_id": "exec-456",
    "status": "completed",
    "query_type": "information_only",
    "sections": [
      {
        "domain_id": "recipes",
        "title": "Рецепты и кулинария",
        "status": "success",
        "results": [...]
      }
    ]
  }
}
```

#### GET /executions/\{execution_id\}
Получение статуса выполнения.

**Пример:**
```bash
curl -X GET https://api.nexus.dev/api/v1/executions/exec-456 \
  -H "Authorization: Bearer <token>"
```

#### GET /executions/\{execution_id\}/status
Получение статуса выполнения (короткий формат).

**Ответ:**
```json
{
  "execution_id": "exec-456",
  "status": "completed",
  "progress": 100,
  "created_at": "2025-01-18T10:00:00Z",
  "updated_at": "2025-01-18T10:00:15Z"
}
```

## 📦 Enterprise API (v2.0.0)

### Batch Operations

#### POST /batch/execute
Пакетное выполнение нескольких операций.

**Пример:**
```bash
curl -X POST https://api.nexus.dev/api/v1/batch/execute \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "request_id": "batch-123",
      "protocol_version": "2.0.0"
    },
    "data": {
      "requests": [
        {
          "query": "хочу борщ",
          "language": "ru"
        },
        {
          "query": "find pizza near me",
          "language": "en"
        }
      ],
      "options": {
        "parallel_execution": true,
        "max_concurrency": 5
      }
    }
  }'
```

### Webhooks

#### POST /webhooks
Создание вебхука.

```json
{
  "data": {
    "url": "https://your-app.com/webhook",
    "events": ["template.completed", "template.failed"],
    "secret": "webhook-secret",
    "active": true
  }
}
```

#### GET /webhooks
Список вебхуков.

#### DELETE /webhooks/\{webhook_id\}
Удаление вебхука.

### Analytics

#### GET /analytics/summary
Получение аналитических данных.

**Параметры:**
- `period`: `day|week|month|year`
- `start_date`: ISO 8601 date
- `end_date`: ISO 8601 date

**Пример:**
```bash
curl -X GET "https://api.nexus.dev/api/v1/analytics/summary?period=week" \
  -H "Authorization: Bearer <token>"
```

### Admin API

#### GET /admin/ai-models
Список AI моделей.

#### POST /admin/ai-models
Создание AI модели.

#### GET /admin/domains
Список доменов.

#### POST /admin/domains
Создание домена.

## ⚡ Rate Limiting

API использует rate limiting для контроля нагрузки:

- **Headers в ответе:**
  ```
  X-RateLimit-Limit: 1000
  X-RateLimit-Remaining: 999
  X-RateLimit-Reset: 1640996100
  ```

- **Информация в metadata:**
  ```json
  "rate_limit_info": {
    "limit": 1000,
    "remaining": 999,
    "reset_at": 1640996100
  }
  ```

## 🛡️ Error Handling

Все ошибки возвращаются в стандартизированном формате:

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query",
    "details": "The query field is required",
    "metadata": {
      "request_id": "req-123",
      "timestamp": 1640995200
    }
  }
}
```

### HTTP Status Codes

- `200` - Успешный запрос
- `400` - Ошибка валидации (VALIDATION_ERROR)
- `401` - Ошибка аутентификации (AUTHENTICATION_ERROR)
- `403` - Ошибка авторизации (AUTHORIZATION_ERROR)
- `404` - Ресурс не найден (NOT_FOUND)
- `409` - Конфликт ресурсов (CONFLICT)
- `429` - Превышен лимит запросов (RATE_LIMIT_ERROR)
- `500` - Внутренняя ошибка сервера (INTERNAL_ERROR)
- `502` - Ошибка внешнего сервиса (EXTERNAL_ERROR)

## 📊 Monitoring

### Metrics Endpoints

#### GET /metrics
Prometheus метрики (требует специального доступа).

#### GET /version
Информация о версиях.

**Ответ:**
```json
{
  "data": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "build_info": {
      "git_commit": "abc123",
      "build_time": "2025-01-18T10:00:00Z",
      "go_version": "1.21.0"
    }
  }
}
```

## 🔧 Frontend Configuration

#### GET /frontend/config
Получение конфигурации UI (публичный эндпоинт).

**Ответ:**
```json
{
  "data": {
    "id": "frontend-config-001",
    "theme": "light",
    "colors": {
      "primary": "#0066CC",
      "secondary": "#00CC66",
      "accent": "#FF6600"
    },
    "branding": {
      "logo": "https://cdn.example.com/logo.png",
      "name": "Nexus Protocol"
    },
    "active": true
  }
}
```

Ответы всегда приходят в формате Application Protocol:

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 3500
  },
  "data": {
    // Результат операции
  }
}
```

## Основные эндпоинты

### Выполнение шаблона

```bash
POST /api/v1/templates/execute
```

### Получение статуса выполнения

```bash
GET /api/v1/templates/status/\{execution_id\}
```

### Health Check

```bash
GET /api/v1/health
```

## Спецификация

Полная спецификация доступна в файле [OpenAPI 3.0](../../api/rest/openapi.yaml).

