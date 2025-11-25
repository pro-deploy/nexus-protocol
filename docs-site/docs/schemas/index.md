---
id: schemas-index
title: Схемы валидации
sidebar_label: Обзор
slug: /schemas
---

# JSON Schema валидации

Nexus Protocol использует JSON Schema для валидации всех сообщений и структур данных.

## 📋 Обзор

### Назначение схем

JSON Schema обеспечивают:

- ✅ **Валидацию структуры** - проверка формата сообщений
- ✅ **Типизацию данных** - строгое определение типов
- ✅ **Документацию** - самоописывающиеся структуры
- ✅ **Генерацию кода** - автоматическая генерация типов
- ✅ **Совместимость** - проверка версий протокола

### Структура схем

```
schemas/
├── message-schema.json    # Основная схема протокола
├── types/                 # Типы данных
├── examples/              # Примеры валидных данных
└── validation/            # Правила валидации
```

## 🔍 Основная схема

### Message Schema

Полная схема Nexus Protocol: [`message-schema.json`](../schemas/message-schema.json)

**Ключевые компоненты:**

#### RequestMetadata
```json
{
  "type": "object",
  "required": ["request_id", "protocol_version", "client_version"],
  "properties": {
    "request_id": {
      "$ref": "#/definitions/UUID"
    },
    "protocol_version": {
      "$ref": "#/definitions/Version"
    },
    "client_version": {
      "$ref": "#/definitions/Version"
    }
  }
}
```

#### ResponseMetadata
```json
{
  "type": "object",
  "required": ["request_id", "server_version", "protocol_version", "timestamp"],
  "properties": {
    "request_id": {
      "$ref": "#/definitions/UUID"
    },
    "server_version": {
      "$ref": "#/definitions/Version"
    }
  }
}
```

#### ErrorDetail
```json
{
  "type": "object",
  "required": ["code", "type", "message"],
  "properties": {
    "code": {
      "type": "string",
      "enum": [
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
    }
  }
}
```

## 🛠️ Использование схем

### Валидация в JavaScript

```javascript
const Ajv = require('ajv');
const ajv = new Ajv();

const schema = require('./schemas/message-schema.json');
const validate = ajv.compile(schema);

// Валидация сообщения
const message = {
  metadata: {
    request_id: "550e8400-e29b-41d4-a716-446655440000",
    protocol_version: "2.0.0",
    client_version: "2.0.0"
  },
  data: {
    query: "хочу борщ"
  }
};

const valid = validate(message);
if (!valid) {
  console.log('Validation errors:', validate.errors);
}
```

### Валидация в Python

```python
import json
import jsonschema

# Загрузка схемы
with open('schemas/message-schema.json', 'r') as f:
    schema = json.load(f)

# Валидация сообщения
message = {
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}

try:
    jsonschema.validate(instance=message, schema=schema)
    print("Message is valid")
except jsonschema.ValidationError as e:
    print(f"Validation error: {e.message}")
```

### Валидация в Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"

    "github.com/xeipuuv/gojsonschema"
)

func main() {
    // Загрузка схемы
    schemaBytes, err := ioutil.ReadFile("schemas/message-schema.json")
    if err != nil {
        panic(err)
    }

    schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

    // Сообщение для валидации
    message := map[string]interface{}{
        "metadata": map[string]interface{}{
            "request_id":      "550e8400-e29b-41d4-a716-446655440000",
            "protocol_version": "2.0.0",
            "client_version":   "2.0.0",
        },
        "data": map[string]interface{}{
            "query": "хочу борщ",
        },
    }

    documentLoader := gojsonschema.NewGoLoader(message)

    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        panic(err)
    }

    if result.Valid() {
        fmt.Println("Message is valid")
    } else {
        fmt.Println("Validation errors:")
        for _, desc := range result.Errors() {
            fmt.Printf("- %s\n", desc)
        }
    }
}
```

### CLI валидация

```bash
# Использование jq для проверки JSON
cat message.json | jq .

# Валидация с помощью jsonschema
jsonschema -i message.json schemas/message-schema.json

# Или с помощью ajv
ajv validate -s schemas/message-schema.json -d message.json
```

## 📋 Определения типов

### UUID
```json
{
  "UUID": {
    "type": "string",
    "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
    "description": "Universally Unique Identifier (UUID) format"
  }
}
```

### Version (Semantic Versioning)
```json
{
  "Version": {
    "type": "string",
    "pattern": "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9.-]+)?(\\+[a-zA-Z0-9.-]+)?$",
    "description": "Semantic version format (MAJOR.MINOR.PATCH)"
  }
}
```

### Timestamp
```json
{
  "Timestamp": {
    "type": "string",
    "format": "date-time",
    "description": "ISO 8601 timestamp"
  }
}
```

## 🔄 Расширение схем

### Кастомные свойства

Схемы поддерживают расширение через `additionalProperties`:

```json
{
  "type": "object",
  "properties": {
    "request_id": { "$ref": "#/definitions/UUID" },
    "protocol_version": { "$ref": "#/definitions/Version" }
  },
  "additionalProperties": true
}
```

### Version-specific схемы

Для разных версий протокола можно создавать отдельные схемы:

```
schemas/
├── v1.0/
│   └── message-schema.json
├── v2.0/
│   └── message-schema.json
└── current -> v2.0/
```

## 🧪 Тестирование схем

### Unit тесты

```javascript
const schema = require('./schemas/message-schema.json');
const testMessages = require('./test-messages.json');

describe('Message Schema Validation', () => {
  const ajv = new Ajv();
  const validate = ajv.compile(schema);

  testMessages.forEach((message, index) => {
    test(`validates message ${index}`, () => {
      const valid = validate(message);
      expect(valid).toBe(true);
      expect(validate.errors).toBeNull();
    });
  });
});
```

### Integration тесты

```go
func TestSchemaValidation(t *testing.T) {
    schemaBytes, err := ioutil.ReadFile("schemas/message-schema.json")
    require.NoError(t, err)

    schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

    testCases := []struct {
        name    string
        message map[string]interface{}
        valid   bool
    }{
        {
            name: "valid execute template request",
            message: map[string]interface{}{
                "metadata": map[string]interface{}{
                    "request_id":      "550e8400-e29b-41d4-a716-446655440000",
                    "protocol_version": "2.0.0",
                    "client_version":   "2.0.0",
                },
                "data": map[string]interface{}{
                    "query": "хочу борщ",
                },
            },
            valid: true,
        },
        {
            name: "invalid - missing required field",
            message: map[string]interface{}{
                "data": map[string]interface{}{
                    "query": "хочу борщ",
                },
            },
            valid: false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            documentLoader := gojsonschema.NewGoLoader(tc.message)
            result, err := gojsonschema.Validate(schemaLoader, documentLoader)
            require.NoError(t, err)

            if tc.valid {
                assert.True(t, result.Valid(), "Expected message to be valid")
            } else {
                assert.False(t, result.Valid(), "Expected message to be invalid")
            }
        })
    }
}
```

## 📚 Дополнительные ресурсы

- [JSON Schema Specification](https://json-schema.org/specification.html)
- [Understanding JSON Schema](https://json-schema.org/understanding-json-schema/)
- [AJV Documentation](https://ajv.js.org/)
- [gojsonschema](https://github.com/xeipuuv/gojsonschema)

## 🔗 Связанные документы

- [Формат сообщений](../protocol/message-format) - структура сообщений
- [Обработка ошибок](../protocol/error-handling) - формат ошибок
- [Версионирование](../protocol/versioning) - управление версиями схем
