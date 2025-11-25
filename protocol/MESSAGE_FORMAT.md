# Формат сообщений Nexus Protocol

## Обзор

Nexus Application Protocol определяет единый формат сообщений для всех транспортных протоколов (HTTP REST, gRPC, WebSocket).

## Структура сообщения

### Базовый формат

Все сообщения следуют единой структуре:

```json
{
  "metadata": {
    "request_id": "uuid",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "string",
    "client_type": "web|mobile|sdk|api|desktop",
    "timestamp": 1640995200,
    "custom_headers": {}
  },
  "data": {
    // Payload зависит от операции
  }
}
```

## Типы сообщений

### 1. Request Message (Запрос)

**Структура:**
```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "web-app",
    "client_type": "web",
    "timestamp": 1640995200
  },
  "data": {
    // Операция-специфичные данные
  }
}
```

**Обязательные поля:**
- `metadata.request_id` - UUID запроса
- `metadata.protocol_version` - версия протокола
- `metadata.client_version` - версия клиента
- `metadata.timestamp` - Unix timestamp

**Опциональные поля:**
- `metadata.client_id` - идентификатор клиента
- `metadata.client_type` - тип клиента
- `metadata.custom_headers` - кастомные заголовки

### 2. Response Message (Ответ)

**Структура:**
```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "timestamp": 1640995235,
    "processing_time_ms": 3500
  },
  "data": {
    // Результат операции
  }
}
```

**Обязательные поля:**
- `metadata.request_id` - UUID запроса (из RequestMetadata)
- `metadata.protocol_version` - версия протокола
- `metadata.server_version` - версия сервера
- `metadata.timestamp` - Unix timestamp ответа
- `metadata.processing_time_ms` - время обработки в миллисекундах

### 3. Error Message (Ошибка)

**Структура:**
```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query",
    "details": "The query field is required",
    "metadata": {
      "request_id": "550e8400-e29b-41d4-a716-446655440000",
      "timestamp": 1640995200
    }
  }
}
```

[Подробнее об ошибках →](./ERROR_HANDLING.md)

## Форматы по транспортам

### HTTP REST

**Request:**
```http
POST /api/v1/templates/execute HTTP/1.1
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "query": "хочу борщ",
  "language": "ru",
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  }
}
```

**Response:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": {
    "execution_id": "exec-456",
    "status": "completed",
    "sections": [...]
  },
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "timestamp": 1640995235,
    "processing_time_ms": 3500
  }
}
```

### gRPC

**Request:**
```protobuf
message ExecuteTemplateRequest {
  string query = 1;
  string language = 2;
  RequestMetadata metadata = 3;
}

message RequestMetadata {
  string request_id = 1;
  string protocol_version = 2;
  string client_version = 3;
  string client_id = 4;
  string client_type = 5;
  int64 timestamp = 6;
}
```

**Response:**
```protobuf
message ExecuteTemplateResponse {
  string execution_id = 1;
  string status = 2;
  ResponseMetadata response_metadata = 3;
}

message ResponseMetadata {
  string request_id = 1;
  string protocol_version = 2;
  string server_version = 3;
  int64 timestamp = 4;
  int32 processing_time_ms = 5;
}
```

### WebSocket

**Request:**
```json
{
  "type": "context_aware_template",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {
    "query": "хочу борщ",
    "language": "ru"
  },
  "timestamp": "2025-01-18T10:00:00Z"
}
```

**Response:**
```json
{
  "type": "context_aware_template_result",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": true,
  "data": {
    "execution_id": "exec-456",
    "status": "completed"
  },
  "timestamp": "2025-01-18T10:00:15Z"
}
```

## Валидация сообщений

### JSON Schema

Все сообщения должны соответствовать JSON Schema:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["metadata"],
  "properties": {
    "metadata": {
      "$ref": "#/definitions/RequestMetadata"
    },
    "data": {
      "type": "object"
    }
  }
}
```

[Полная схема →](../schemas/message-schema.json)

### Валидация полей

#### request_id
- **Тип:** string (UUID)
- **Формат:** `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
- **Обязательно:** Да
- **Пример:** `"550e8400-e29b-41d4-a716-446655440000"`

#### protocol_version
- **Тип:** string
- **Формат:** Semantic Versioning `MAJOR.MINOR.PATCH`
- **Обязательно:** Да
- **Пример:** `"2.0.0"`

#### client_version
- **Тип:** string
- **Формат:** Semantic Versioning `MAJOR.MINOR.PATCH`
- **Обязательно:** Да
- **Пример:** `"2.0.0"`

#### client_type
- **Тип:** string
- **Значения:** `web`, `mobile`, `sdk`, `api`, `desktop`
- **Обязательно:** Нет
- **Пример:** `"web"`

#### timestamp
- **Тип:** integer (int64)
- **Формат:** Unix timestamp (секунды с 1970-01-01)
- **Обязательно:** Да
- **Пример:** `1640995200`

## Примеры

### Пример 1: ExecuteTemplate Request

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "web-app",
    "client_type": "web",
    "timestamp": 1640995200
  },
  "data": {
    "query": "хочу борщ",
    "language": "ru",
    "context": {
      "user_id": "user-123",
      "session_id": "session-456"
    },
    "options": {
      "timeout_ms": 30000,
      "max_results_per_domain": 5,
      "include_web_search": true
    }
  }
}
```

### Пример 2: ExecuteTemplate Response (Информационный запрос)

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "server_version": "1.1.3",
    "timestamp": 1640995235,
    "processing_time_ms": 3500
  },
  "data": {
    "execution_id": "exec-789",
    "intent_id": "intent-abc",
    "status": "completed",
    "query_type": "information_only",
    "sections": [
      {
        "domain_id": "recipes",
        "title": "Рецепты и кулинария",
        "status": "success",
        "results": [...]
      }
    ],
    "domain_analysis": {
      "selected_domains": [
        {
          "domain_id": "recipes",
          "name": "Рецепты и кулинария",
          "type": "recipes",
          "confidence": 0.87,
          "relevance": 0.92,
          "reason": "Высокая уверенность: найдены ключевые слова рецепта",
          "priority": 60
        }
      ],
      "rejected_domains": [
        {
          "domain_id": "commerce",
          "name": "Коммерция и покупки",
          "confidence": 0.23,
          "reason": "Низкая уверенность: минимальные признаки релевантности"
        }
      ],
      "confidence": 0.87,
      "analysis_algorithm": "hybrid_keyword_semantic"
    }
  }
}
```

### Пример 2a: ExecuteTemplate Response (Запрос с покупкой)

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440001",
    "protocol_version": "2.0.0",
    "server_version": "1.1.3",
    "timestamp": 1640995235,
    "processing_time_ms": 245
  },
  "data": {
    "execution_id": "exec-790",
    "intent_id": "intent-def",
    "status": "completed",
    "query_type": "with_purchases_services",
    "sections": [
      {
        "domain_id": "commerce",
        "title": "Коммерческие предложения",
        "status": "success",
        "response_time_ms": 200,
        "results": [
          {
            "id": "product-456",
            "type": "product_purchase",
            "title": "Coca-Cola 1л бутылка",
            "description": "Найдено в 3 магазинах рядом с вами",
            "data": {
              "price": "89 ₽",
              "availability": "в наличии",
              "stores_count": 3,
              "nearest_store": {
                "name": "Пятерочка",
                "distance": "200м",
                "address": "ул. Ленина, 15",
                "pickup_available": true,
                "work_hours": "Круглосуточно"
              }
            },
            "relevance": 0.95,
            "confidence": 0.88,
            "actions": [
              {
                "type": "reserve_product",
                "label": "Зарезервировать товар",
                "method": "POST",
                "url": "/api/v1/commerce/reserve"
              },
              {
                "type": "purchase",
                "label": "Купить сейчас",
                "method": "POST",
                "url": "/api/v1/commerce/purchase"
              }
            ]
          }
        ]
      }
    ],
    "ranking": {
      "items": [
        {
          "id": "product-456",
          "score": 0.92,
          "rank": 1
        }
      ],
      "algorithm": "weighted_relevance_confidence"
    }
  }
}
```

### Пример 3: Error Response

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query",
    "details": "The query field is required for template execution",
    "metadata": {
      "request_id": "550e8400-e29b-41d4-a716-446655440000",
      "timestamp": 1640995200
    }
  }
}
```

## Правила обработки

### 1. Обязательность метаданных

Все запросы **должны** включать `RequestMetadata`:
- `request_id` - обязательно
- `protocol_version` - обязательно
- `client_version` - обязательно
- `timestamp` - обязательно

Все ответы **должны** включать `ResponseMetadata`:
- `request_id` - обязательно (из RequestMetadata)
- `protocol_version` - обязательно
- `server_version` - обязательно
- `timestamp` - обязательно
- `processing_time_ms` - обязательно

### 2. Корреляция запросов

`request_id` из `RequestMetadata` должен совпадать с `request_id` в `ResponseMetadata` для корреляции запросов и ответов.

### 3. Версионирование

`protocol_version` должна соответствовать версии протокола, которую поддерживает клиент. Сервер проверяет совместимость версий.

### 4. Временные метки

- `timestamp` в запросе - время создания запроса клиентом
- `timestamp` в ответе - время создания ответа сервером
- `processing_time_ms` - время обработки запроса на сервере

## Расширение формата

### Кастомные заголовки

Можно добавлять кастомные заголовки в `metadata.custom_headers`:

```json
{
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "custom_headers": {
      "x-feature-flag": "new-ui",
      "x-experiment-id": "exp-456"
    }
  }
}
```

### Дополнительные поля в data

Поле `data` может содержать любые данные, специфичные для операции. Структура определяется API спецификацией для конкретной операции.

## Frontend Configuration

Клиенты могут получать активную конфигурацию визуала через публичный endpoint:

**Request:**
```http
GET /api/v1/frontend/config HTTP/1.1
```

**Response:**
```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "server_version": "1.1.3",
    "timestamp": 1640995235,
    "processing_time_ms": 5
  },
  "data": {
    "id": "frontend-config-001",
    "name": "Corporate Theme",
    "theme": "light",
    "colors": {
      "primary": "#0066CC",
      "secondary": "#00CC66",
      "accent": "#FF6600",
      "background": "#FFFFFF",
      "text": "#333333"
    },
    "layout": {
      "header": {
        "height": "64px",
        "sticky": true
      },
      "sidebar": {
        "width": "240px",
        "collapsible": true
      }
    },
    "components": {
      "button": {
        "border_radius": "8px",
        "padding": "12px 24px"
      }
    },
    "branding": {
      "logo": "https://cdn.example.com/logo.png",
      "name": "Nexus Protocol",
      "favicon": "https://cdn.example.com/favicon.ico"
    },
    "active": true
  }
}
```

**Особенности:**
- Публичный endpoint (не требует аутентификации)
- Возвращает только активную конфигурацию
- Используется клиентами для настройки UI

## См. также

- [Метаданные](./METADATA.md) - детали RequestMetadata/ResponseMetadata
- [Обработка ошибок](./ERROR_HANDLING.md) - формат ошибок
- [Версионирование](../versioning/README.md) - правила версий
- [JSON Schema](../schemas/message-schema.json) - схема валидации
