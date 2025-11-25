---
id: validation-examples
title: Примеры валидации
sidebar_label: Примеры валидации
---

# Примеры валидации JSON Schema

Примеры валидных и невалидных сообщений с объяснениями.

## ✅ Валидные сообщения

### 1. ExecuteTemplate Request

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
      "x-feature-flag": "new-ui"
    }
  },
  "data": {
    "query": "хочу борщ",
    "language": "ru",
    "context": {
      "user_id": "user-123",
      "session_id": "session-456",
      "location": {
        "latitude": 55.7558,
        "longitude": 37.6173,
        "accuracy": 50
      },
      "locale": "ru-RU",
      "currency": "RUB",
      "region": "RU"
    },
    "options": {
      "timeout_ms": 30000,
      "max_results_per_domain": 5,
      "parallel_execution": true,
      "include_web_search": true
    }
  }
}
```

**Почему валидно:**
- ✅ `request_id` - корректный UUID
- ✅ `protocol_version` и `client_version` - semantic versioning
- ✅ Все обязательные поля присутствуют
- ✅ Типы данных соответствуют схеме

### 2. ExecuteTemplate Response

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
    "execution_id": "exec-456",
    "status": "completed",
    "query_type": "information_only",
    "sections": [
      {
        "domain_id": "recipes",
        "title": "Рецепты и кулинария",
        "status": "success",
        "results": [
          {
            "id": "recipe-123",
            "type": "recipe",
            "title": "Борщ украинский",
            "description": "Классический рецепт украинского борща",
            "relevance": 0.95,
            "confidence": 0.88
          }
        ]
      }
    ]
  }
}
```

**Почему валидно:**
- ✅ Корректная структура ResponseMetadata
- ✅ Enterprise поля (rate_limit_info) опциональны
- ✅ Массивы results правильно структурированы

### 3. Error Response

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

**Почему валидно:**
- ✅ `code` из разрешенного enum
- ✅ `type` соответствует коду ошибки
- ✅ Все обязательные поля присутствуют

## ❌ Невалидные сообщения

### 1. Отсутствует обязательное поле

```json
{
  "metadata": {
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}
```

**Почему невалидно:**
- ❌ Отсутствует обязательное поле `request_id` в `metadata`

**Ошибка валидации:**
```json
[
  {
    "keyword": "required",
    "dataPath": ".metadata",
    "schemaPath": "#/properties/metadata/required",
    "params": {
      "missingProperty": "request_id"
    },
    "message": "should have required property 'request_id'"
  }
]
```

### 2. Некорректный UUID

```json
{
  "metadata": {
    "request_id": "not-a-uuid",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}
```

**Почему невалидно:**
- ❌ `request_id` не соответствует паттерну UUID

**Ошибка валидации:**
```json
[
  {
    "keyword": "pattern",
    "dataPath": ".metadata.request_id",
    "schemaPath": "#/definitions/UUID/pattern",
    "params": {
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
    },
    "message": "should match pattern \"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$\""
  }
]
```

### 3. Некорректный enum

```json
{
  "error": {
    "code": "INVALID_ERROR_CODE",
    "type": "VALIDATION_ERROR",
    "message": "Test error"
  }
}
```

**Почему невалидно:**
- ❌ `code` не из разрешенного списка enum значений

**Ошибка валидации:**
```json
[
  {
    "keyword": "enum",
    "dataPath": ".error.code",
    "schemaPath": "#/definitions/ErrorDetail/properties/code/enum",
    "params": {
      "allowedValues": [
        "VALIDATION_FAILED",
        "AUTHENTICATION_FAILED",
        "AUTHORIZATION_FAILED",
        "NOT_FOUND",
        "CONFLICT",
        "RATE_LIMIT_ERROR",
        "INTERNAL_ERROR",
        "EXTERNAL_ERROR",
        "PROTOCOL_VERSION_ERROR"
      ]
    },
    "message": "should be equal to one of the allowed values"
  }
]
```

### 4. Некорректный тип version

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "latest",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}
```

**Почему невалидно:**
- ❌ `protocol_version` не соответствует паттерну semantic versioning

**Ошибка валидации:**
```json
[
  {
    "keyword": "pattern",
    "dataPath": ".metadata.protocol_version",
    "schemaPath": "#/definitions/Version/pattern",
    "params": {
      "pattern": "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9.-]+)?(\\+[a-zA-Z0-9.-]+)?$"
    },
    "message": "should match pattern \"^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9.-]+)?(\\+[a-zA-Z0-9.-]+)?$\""
  }
]
```

### 5. Превышение максимальной длины

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "a".repeat(101)  // 101 символ
  },
  "data": {
    "query": "хочу борщ"
  }
}
```

**Почему невалидно:**
- ❌ `client_id` превышает максимальную длину 100 символов

**Ошибка валидации:**
```json
[
  {
    "keyword": "maxLength",
    "dataPath": ".metadata.client_id",
    "schemaPath": "#/properties/metadata/properties/client_id/maxLength",
    "params": {
      "limit": 100
    },
    "message": "should NOT be longer than 100 characters"
  }
]
```

## 🔄 Частично валидные сообщения

### 1. Дополнительные поля (расширение)

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "custom_field": "allowed"
  },
  "data": {
    "query": "хочу борщ",
    "extra_data": {
      "custom_property": "also allowed"
    }
  }
}
```

**Почему валидно:**
- ✅ `additionalProperties: true` позволяет расширения
- ✅ Схема не ограничивает дополнительные поля

### 2. Опциональные поля отсутствуют

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}
```

**Почему валидно:**
- ✅ Опциональные поля `client_id`, `client_type`, `timestamp` могут отсутствовать
- ✅ Только `request_id`, `protocol_version`, `client_version` обязательны

## 🧪 Тестовые сценарии

### Batch Operations (Enterprise)

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440001",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "requests": [
      {
        "query": "хочу борщ",
        "language": "ru"
      },
      {
        "query": "find pizza",
        "language": "en"
      }
    ],
    "options": {
      "parallel_execution": true,
      "max_concurrency": 5,
      "timeout_ms": 60000
    }
  }
}
```

### Webhook Configuration

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440002",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "url": "https://api.example.com/webhooks/nexus",
    "events": ["template.completed", "template.failed"],
    "secret": "webhook-secret-123",
    "active": true,
    "headers": {
      "X-API-Key": "custom-api-key"
    }
  }
}
```

### Analytics Request

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440003",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "period": {
      "start_date": "2025-01-01T00:00:00Z",
      "end_date": "2025-01-18T23:59:59Z"
    },
    "metrics": [
      "requests_total",
      "success_rate",
      "avg_response_time",
      "error_rate"
    ],
    "filters": {
      "client_type": "web",
      "domain": "recipes"
    }
  }
}
```

## 📋 Сводка правил валидации

### Обязательные поля (required)

#### RequestMetadata
- `request_id` (UUID)
- `protocol_version` (Version)
- `client_version` (Version)

#### ResponseMetadata
- `request_id` (UUID)
- `protocol_version` (Version)
- `server_version` (Version)
- `timestamp` (integer)
- `processing_time_ms` (integer)

#### ErrorDetail
- `code` (enum)
- `type` (string)
- `message` (string)

### Ограничения типов

#### Строки
- `request_id`: UUID pattern
- `*_version`: Semantic version pattern
- `client_id`: maxLength 100
- `message`: maxLength 1000
- `details`: maxLength 5000

#### Числа
- `timestamp`: int64 Unix timestamp
- `processing_time_ms`: int32 ≥ 0
- `latitude`/`longitude`: double
- `accuracy`: double ≥ 0
- `relevance`/`confidence`: float 0.0-1.0

#### Массивы
- `roles`: array of strings
- `actions`: array of Action objects
- `results`: array of ResultItem objects

### Расширяемость

- `custom_headers`: object с additionalProperties
- `metadata`: additionalProperties: true
- `data`: без строгих ограничений
- `context`: extensible для будущих полей

## 🛠️ Инструменты валидации

### Онлайн валидаторы

1. **JSON Schema Validator**: https://www.jsonschemavalidator.net/
2. **AJV Playground**: https://ajv.js.org/
3. **JSON Schema Lint**: https://jsonschemalint.com/

### CLI инструменты

```bash
# ajv-cli
npm install -g ajv-cli
ajv validate -s schema.json -d data.json

# jsonschema
pip install jsonschema
python -c "import jsonschema; jsonschema.validate(data, schema)"

# gojsonschema
go install github.com/atombender/go-jsonschema/cmd/gojsonschema@latest
gojsonschema -p main schema.json
```

### IDE интеграция

- **VS Code**: JSON Schema Store
- **IntelliJ IDEA**: JSON Schema mappings
- **GoLand**: struct tags validation
