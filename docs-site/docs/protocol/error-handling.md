---
id: error-handling
title: Обработка ошибок
sidebar_label: Обработка ошибок
---

# Обработка ошибок Nexus Protocol

## Обзор

Nexus Application Protocol определяет стандартизированный формат обработки ошибок для всех транспортных протоколов. Единый формат ошибок обеспечивает:

- ✅ Консистентную обработку ошибок
- ✅ Легкую отладку и трассировку
- ✅ Автоматическую обработку на клиенте
- ✅ Мониторинг и аналитику ошибок

## Формат ошибки

### Базовая структура

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

### Поля ErrorDetail

#### code (обязательно)

**Тип:** string  
**Описание:** Машинно-читаемый код ошибки  
**Формат:** `UPPER_SNAKE_CASE`  
**Пример:** `"VALIDATION_FAILED"`

**Правила:**
- Уникальный идентификатор типа ошибки
- Используется для программной обработки
- Не должен изменяться между версиями протокола

#### type (обязательно)

**Тип:** string (enum)  
**Описание:** Категория ошибки  
**Значения:**
- `VALIDATION_ERROR` - ошибка валидации входных данных
- `AUTHENTICATION_ERROR` - ошибка аутентификации
- `AUTHORIZATION_ERROR` - ошибка авторизации
- `NOT_FOUND` - ресурс не найден
- `CONFLICT` - конфликт ресурсов
- `RATE_LIMIT_ERROR` - превышен лимит запросов
- `INTERNAL_ERROR` - внутренняя ошибка сервера
- `EXTERNAL_ERROR` - ошибка внешнего сервиса
- `PROTOCOL_VERSION_ERROR` - несовместимость версий

**Пример:** `"VALIDATION_ERROR"`

#### message (обязательно)

**Тип:** string  
**Описание:** Человеко-читаемое сообщение об ошибке  
**Максимальная длина:** 1000 символов  
**Пример:** `"Query cannot be empty"`

**Правила:**
- Описывает проблему понятным языком
- Может быть локализовано
- Не должно содержать чувствительную информацию

#### field (опционально)

**Тип:** string  
**Описание:** Поле, вызвавшее ошибку (для валидационных ошибок)  
**Максимальная длина:** 100 символов  
**Пример:** `"query"`, `"metadata.protocol_version"`

**Правила:**
- Указывает конкретное поле с ошибкой
- Используется для валидационных ошибок
- Может быть вложенным (через точку)

#### details (опционально)

**Тип:** string  
**Описание:** Детальная информация об ошибке  
**Максимальная длина:** 5000 символов  
**Пример:** `"The query field is required for template execution"`

**Правила:**
- Дополнительная информация для отладки
- Может содержать технические детали
- Не должно содержать чувствительную информацию

#### metadata (опционально)

**Тип:** object  
**Описание:** Дополнительные метаданные ошибки  
**Пример:**
```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": 1640995200,
  "trace_id": "trace-123",
  "span_id": "span-456"
}
```

## Коды ошибок

### VALIDATION_ERROR

**Коды:**
- `VALIDATION_FAILED` - общая ошибка валидации
- `INVALID_FORMAT` - неверный формат данных
- `MISSING_REQUIRED_FIELD` - отсутствует обязательное поле
- `INVALID_VALUE` - неверное значение поля
- `FIELD_TOO_LONG` - поле превышает максимальную длину
- `FIELD_TOO_SHORT` - поле меньше минимальной длины

**HTTP Status:** 400 Bad Request  
**gRPC Status:** INVALID_ARGUMENT

**Пример:**
```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query",
    "details": "The query field is required and must be at least 1 character"
  }
}
```

### AUTHENTICATION_ERROR

**Коды:**
- `AUTHENTICATION_FAILED` - ошибка аутентификации
- `INVALID_TOKEN` - неверный токен
- `TOKEN_EXPIRED` - токен истек
- `TOKEN_MALFORMED` - неверный формат токена

**HTTP Status:** 401 Unauthorized  
**gRPC Status:** UNAUTHENTICATED

**Пример:**
```json
{
  "error": {
    "code": "TOKEN_EXPIRED",
    "type": "AUTHENTICATION_ERROR",
    "message": "JWT token has expired",
    "details": "Token expired at 2025-01-18T10:00:00Z"
  }
}
```

### AUTHORIZATION_ERROR

**Коды:**
- `AUTHORIZATION_FAILED` - ошибка авторизации
- `INSUFFICIENT_PERMISSIONS` - недостаточно прав
- `FORBIDDEN_RESOURCE` - доступ к ресурсу запрещен

**HTTP Status:** 403 Forbidden  
**gRPC Status:** PERMISSION_DENIED

**Пример:**
```json
{
  "error": {
    "code": "INSUFFICIENT_PERMISSIONS",
    "type": "AUTHORIZATION_ERROR",
    "message": "User does not have permission to execute templates",
    "details": "Required permission: templates.execute"
  }
}
```

### NOT_FOUND

**Коды:**
- `RESOURCE_NOT_FOUND` - ресурс не найден
- `ENDPOINT_NOT_FOUND` - endpoint не найден
- `EXECUTION_NOT_FOUND` - выполнение не найдено

**HTTP Status:** 404 Not Found  
**gRPC Status:** NOT_FOUND

**Пример:**
```json
{
  "error": {
    "code": "EXECUTION_NOT_FOUND",
    "type": "NOT_FOUND",
    "message": "Execution with ID 'exec-123' not found",
    "details": "Execution may have been deleted or never existed"
  }
}
```

### CONFLICT

**Коды:**
- `RESOURCE_CONFLICT` - конфликт ресурсов
- `DUPLICATE_RESOURCE` - дублирующийся ресурс
- `CONCURRENT_MODIFICATION` - одновременная модификация

**HTTP Status:** 409 Conflict  
**gRPC Status:** ALREADY_EXISTS или ABORTED

**Пример:**
```json
{
  "error": {
    "code": "DUPLICATE_RESOURCE",
    "type": "CONFLICT",
    "message": "User with email 'user@example.com' already exists"
  }
}
```

### RATE_LIMIT_ERROR

**Коды:**
- `RATE_LIMIT_EXCEEDED` - превышен лимит запросов

**HTTP Status:** 429 Too Many Requests  
**gRPC Status:** RESOURCE_EXHAUSTED

**Пример:**
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "type": "RATE_LIMIT_ERROR",
    "message": "Rate limit exceeded",
    "details": "Limit: 1000 requests per minute. Try again in 60 seconds.",
    "metadata": {
      "limit": 1000,
      "remaining": 0,
      "reset_at": 1640996100
    }
  }
}
```

### INTERNAL_ERROR

**Коды:**
- `INTERNAL_ERROR` - внутренняя ошибка сервера
- `DATABASE_ERROR` - ошибка базы данных
- `PROCESSING_ERROR` - ошибка обработки

**HTTP Status:** 500 Internal Server Error  
**gRPC Status:** INTERNAL

**Пример:**
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "type": "INTERNAL_ERROR",
    "message": "An internal error occurred",
    "details": "Error ID: err-12345. Please contact support."
  }
}
```

### EXTERNAL_ERROR

**Коды:**
- `EXTERNAL_SERVICE_ERROR` - ошибка внешнего сервиса
- `SERVICE_UNAVAILABLE` - сервис недоступен
- `TIMEOUT` - таймаут запроса

**HTTP Status:** 502 Bad Gateway или 503 Service Unavailable  
**gRPC Status:** UNAVAILABLE

**Пример:**
```json
{
  "error": {
    "code": "EXTERNAL_SERVICE_ERROR",
    "type": "EXTERNAL_ERROR",
    "message": "External AI service is unavailable",
    "details": "Service 'ollama' returned error: connection timeout"
  }
}
```

### PROTOCOL_VERSION_ERROR

**Коды:**
- `PROTOCOL_VERSION_MISMATCH` - несовместимость версий

**HTTP Status:** 400 Bad Request  
**gRPC Status:** INVALID_ARGUMENT

**Пример:**
```json
{
  "error": {
    "code": "PROTOCOL_VERSION_MISMATCH",
    "type": "PROTOCOL_VERSION_ERROR",
    "message": "Protocol version mismatch",
    "details": "Client version 1.1.0 is not compatible with server version 1.0.0"
  }
}
```

## Форматы по транспортам

### HTTP REST

**Response:**
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query"
  }
}
```

**Response Headers:**
```http
X-Error-Code: VALIDATION_FAILED
X-Error-Type: VALIDATION_ERROR
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

### gRPC

**Status:**
```go
st := status.New(codes.InvalidArgument, "Validation failed")
st, _ = st.WithDetails(&pb.ErrorDetail{
    ErrorCode: "VALIDATION_FAILED",
    ErrorType: "VALIDATION_ERROR",
    Message:   "Query cannot be empty",
    Field:     "query",
    Details:   "The query field is required",
})
return nil, st.Err()
```

**Protocol Buffers:**
```protobuf
// Примечание: В Protocol Buffers используются имена error_code и error_type,
// но при сериализации в JSON они должны соответствовать полям "code" и "type"
// согласно Nexus Protocol v2.0.0
message ErrorDetail {
  string error_code = 1;  // В JSON: "code"
  string error_type = 2;   // В JSON: "type"
  string message = 3;
  string field = 4;
  string details = 5;
  map&lt;string, string&gt; metadata = 6;
}
```

### WebSocket

**Error Message:**
```json
{
  "type": "error",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": false,
  "error": "Query cannot be empty",
  "data": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "field": "query"
  },
  "timestamp": "2025-01-18T10:00:00Z"
}
```

## Обработка ошибок на клиенте

### JavaScript/TypeScript

```typescript
interface ErrorResponse {
  error: {
    code: string;
    type: string;
    message: string;
    field?: string;
    details?: string;
  };
}

async function handleError(response: Response): Promise<never> {
  const error: ErrorResponse = await response.json();
  
  switch (error.error.code) {
    case 'VALIDATION_FAILED':
      throw new ValidationError(error.error);
    case 'AUTHENTICATION_FAILED':
      throw new AuthenticationError(error.error);
    case 'RATE_LIMIT_EXCEEDED':
      throw new RateLimitError(error.error);
    default:
      throw new ProtocolError(error.error);
  }
}
```

### Go

```go
type ErrorDetail struct {
    Code    string            `json:"code"`
    Type    string            `json:"type"`
    Message string            `json:"message"`
    Field   string            `json:"field,omitempty"`
    Details string            `json:"details,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
}

func handleError(resp *http.Response) error {
    var errResp struct {
        Error ErrorDetail `json:"error"`
    }
    
    json.NewDecoder(resp.Body).Decode(&errResp)
    
    switch errResp.Error.Code {
    case "VALIDATION_FAILED":
        return &ValidationError{Detail: errResp.Error}
    case "AUTHENTICATION_FAILED":
        return &AuthenticationError{Detail: errResp.Error}
    default:
        return &ProtocolError{Detail: errResp.Error}
    }
}
```

## Best Practices

### 1. Всегда включайте request_id

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "metadata": {
      "request_id": "550e8400-e29b-41d4-a716-446655440000"
    }
  }
}
```

### 2. Используйте правильные HTTP статусы

```javascript
const statusMap = {
  'VALIDATION_ERROR': 400,
  'AUTHENTICATION_ERROR': 401,
  'AUTHORIZATION_ERROR': 403,
  'NOT_FOUND': 404,
  'CONFLICT': 409,
  'RATE_LIMIT_ERROR': 429,
  'INTERNAL_ERROR': 500,
  'EXTERNAL_ERROR': 502
};
```

### 3. Не раскрывайте внутренние детали

```json
// ❌ Плохо
{
  "error": {
    "message": "Database connection failed: Connection refused (127.0.0.1:5432)"
  }
}

// ✅ Хорошо
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An internal error occurred",
    "details": "Error ID: err-12345. Please contact support."
  }
}
```

### 4. Валидационные ошибки должны указывать поле

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query",
    "details": "The query field is required and must be at least 1 character"
  }
}
```

## См. также

- [Формат сообщений](./MESSAGE_FORMAT.md) - общая структура сообщений
- [Метаданные](./METADATA.md) - метаданные в ошибках
- [Версионирование](../versioning/README.md) - обработка ошибок версий
