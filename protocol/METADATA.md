# Метаданные Nexus Protocol

## Обзор

Метаданные (Metadata) - это стандартизированная информация, которая включается во все запросы и ответы протокола Nexus. Метаданные обеспечивают:

- ✅ Корреляцию запросов и ответов
- ✅ Версионирование и совместимость
- ✅ Трассировку и отладку
- ✅ Мониторинг и аналитику

## RequestMetadata (Метаданные запроса)

### Структура

```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "protocol_version": "1.0.0",
  "client_version": "1.0.0",
  "client_id": "web-app",
  "client_type": "web",
  "timestamp": 1640995200,
  "custom_headers": {
    "x-feature-flag": "new-ui"
  }
}
```

### Поля

#### request_id (обязательно)

**Тип:** string (UUID v4)  
**Описание:** Уникальный идентификатор запроса для корреляции с ответом  
**Формат:** `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`  
**Пример:** `"550e8400-e29b-41d4-a716-446655440000"`

**Правила:**
- Должен быть уникальным для каждого запроса
- Используется для корреляции запросов и ответов
- Может использоваться для трассировки и логирования

#### protocol_version (обязательно)

**Тип:** string  
**Описание:** Версия протокола Nexus, которую использует клиент  
**Формат:** Semantic Versioning `MAJOR.MINOR.PATCH`  
**Пример:** `"1.0.0"`

**Правила:**
- Должна соответствовать версии протокола, которую поддерживает клиент
- Сервер проверяет совместимость версий
- Формат: `MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`

**Примеры:**
- `"1.0.0"` - стабильная версия
- `"1.1.0"` - новая минорная версия
- `"1.0.0-alpha.1"` - пререлизная версия

#### client_version (обязательно)

**Тип:** string  
**Описание:** Версия клиентского приложения или SDK  
**Формат:** Semantic Versioning `MAJOR.MINOR.PATCH`  
**Пример:** `"1.0.0"`

**Правила:**
- Версия клиентского приложения или SDK
- Используется для аналитики и отладки
- Может отличаться от protocol_version

#### client_id (опционально)

**Тип:** string  
**Описание:** Идентификатор клиентского приложения  
**Максимальная длина:** 100 символов  
**Пример:** `"web-app"`, `"mobile-ios"`, `"sdk-python"`

**Правила:**
- Уникальный идентификатор клиентского приложения
- Используется для аналитики и мониторинга
- Может быть использован для rate limiting

#### client_type (опционально)

**Тип:** string (enum)  
**Описание:** Тип клиентского приложения  
**Значения:** `web`, `mobile`, `sdk`, `api`, `desktop`  
**Пример:** `"web"`

**Правила:**
- Категоризация типа клиента
- Используется для аналитики и оптимизации
- Может влиять на поведение сервера

#### timestamp (обязательно)

**Тип:** integer (int64)  
**Описание:** Unix timestamp создания запроса (секунды с 1970-01-01 00:00:00 UTC)  
**Пример:** `1640995200`

**Правила:**
- Время создания запроса на стороне клиента
- Используется для расчета задержек и таймаутов
- Должен быть в UTC

#### custom_headers (опционально)

**Тип:** object (map<string, string>)  
**Описание:** Кастомные заголовки для расширения функциональности  
**Пример:**
```json
{
  "x-feature-flag": "new-ui",
  "x-experiment-id": "exp-456",
  "x-debug-mode": "true",
  "x-priority": "high",
  "x-cache-control": "cache-first",
  "x-cache-ttl": "300"
}
```

**Стандартизированные заголовки:**

**Приоритет запросов:**
- `x-priority` (string) - приоритет запроса
  - Значения: `low`, `normal`, `high`, `critical`
  - По умолчанию: `normal`

- `x-request-source` (string) - источник запроса
  - Значения: `user`, `system`, `batch`, `webhook`

**Кэширование:**
- `x-cache-control` (string) - контроль кэширования
  - Значения: `no-cache`, `cache-only`, `cache-first`, `network-first`

- `x-cache-ttl` (int32) - TTL кэша в секундах
- `x-cache-key` (string) - кастомный ключ кэша

**A/B тестирование:**
- `x-experiment-id` (string) - ID эксперимента
- `x-feature-{name}` (string) - feature flag

**Правила:**
- Произвольные ключ-значение пары
- Используются для feature flags, экспериментов, отладки
- Префикс `x-` для кастомных заголовков
- Не должны содержать чувствительную информацию

## ResponseMetadata (Метаданные ответа)

### Структура

```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "protocol_version": "1.0.0",
  "server_version": "1.0.2",
  "timestamp": 1640995235,
  "processing_time_ms": 3500,
  "rate_limit_info": {
    "limit": 1000,
    "remaining": 950,
    "reset_at": 1640996100
  },
  "cache_info": {
    "cache_hit": true,
    "cache_key": "template:query:hash",
    "cache_ttl": 300
  },
  "quota_info": {
    "quota_used": 50000,
    "quota_limit": 100000,
    "quota_type": "requests"
  }
}
```

### Поля

#### request_id (обязательно)

**Тип:** string (UUID v4)  
**Описание:** Идентификатор запроса из RequestMetadata  
**Правила:**
- Должен совпадать с `request_id` из RequestMetadata
- Используется для корреляции запросов и ответов

#### protocol_version (обязательно)

**Тип:** string  
**Описание:** Версия протокола, которую использует сервер  
**Правила:**
- Версия протокола, поддерживаемая сервером
- Может отличаться от client_version, если версии совместимы

#### server_version (обязательно)

**Тип:** string  
**Описание:** Версия серверного приложения  
**Формат:** Semantic Versioning `MAJOR.MINOR.PATCH`  
**Пример:** `"1.0.2"`

**Правила:**
- Версия серверного приложения
- Используется для отладки и мониторинга
- Может отличаться от protocol_version

#### timestamp (обязательно)

**Тип:** integer (int64)  
**Описание:** Unix timestamp создания ответа (секунды с 1970-01-01 00:00:00 UTC)  
**Пример:** `1640995235`

**Правила:**
- Время создания ответа на стороне сервера
- Используется для расчета задержек
- Должен быть в UTC

#### processing_time_ms (обязательно)

**Тип:** integer (int32)  
**Описание:** Время обработки запроса на сервере в миллисекундах  
**Минимум:** 0  
**Пример:** `3500`

**Правила:**
- Время от получения запроса до отправки ответа
- Включает время обработки на всех уровнях
- Используется для мониторинга производительности

#### rate_limit_info (опционально)

**Тип:** object (RateLimitInfo)  
**Описание:** Информация о rate limiting для текущего запроса  

**Структура:**
```json
{
  "limit": 1000,        // лимит запросов
  "remaining": 950,     // оставшиеся запросы
  "reset_at": 1640996100 // время сброса лимита (Unix timestamp)
}
```

#### cache_info (опционально)

**Тип:** object (CacheInfo)  
**Описание:** Информация о кэшировании для текущего запроса  

**Структура:**
```json
{
  "cache_hit": true,         // был ли кэш
  "cache_key": "template:query:hash", // ключ кэша
  "cache_ttl": 300           // TTL кэша в секундах
}
```

#### quota_info (опционально)

**Тип:** object (QuotaInfo)  
**Описание:** Информация о квотах для текущего клиента  

**Структура:**
```json
{
  "quota_used": 50000,   // использовано квоты
  "quota_limit": 100000, // лимит квоты
  "quota_type": "requests" // тип квоты (requests, data, storage, bandwidth)
}
```

## Использование метаданных

### Корреляция запросов и ответов

```javascript
// Клиент отправляет запрос
const requestId = uuid.v4();
const request = {
  metadata: {
    request_id: requestId,
    protocol_version: "1.0.0",
    client_version: "1.0.0",
    timestamp: Date.now() / 1000
  },
  data: { query: "хочу борщ" }
};

// Сервер возвращает ответ с тем же request_id
const response = {
  metadata: {
    request_id: requestId, // Тот же ID
    protocol_version: "1.0.0",
    server_version: "1.0.2",
    timestamp: Date.now() / 1000,
    processing_time_ms: 3500
  },
  data: { /* результат */ }
};
```

### Версионирование

```javascript
// Клиент указывает версию протокола
const metadata = {
  request_id: uuid.v4(),
  protocol_version: "1.0.0", // Версия протокола клиента
  client_version: "1.0.0",   // Версия клиентского приложения
  timestamp: Date.now() / 1000
};

// Сервер проверяет совместимость
if (!isCompatible(metadata.protocol_version, serverProtocolVersion)) {
  return {
    error: {
      code: "PROTOCOL_VERSION_MISMATCH",
      message: "Protocol version mismatch"
    }
  };
}
```

### Трассировка

```javascript
// Все логи включают request_id
logger.info("Processing request", {
  request_id: metadata.request_id,
  client_version: metadata.client_version,
  processing_time_ms: response.metadata.processing_time_ms
});
```

### Мониторинг

```javascript
// Метрики включают метаданные
metrics.record({
  operation: "execute_template",
  client_type: metadata.client_type,
  client_version: metadata.client_version,
  processing_time_ms: response.metadata.processing_time_ms,
  protocol_version: metadata.protocol_version
});
```

## Форматы по транспортам

### HTTP REST

**Request Headers:**
```http
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
X-Protocol-Version: 1.0.0
X-Client-Version: 1.0.0
```

**Request Body:**
```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "1.0.0",
    "client_version": "1.0.0",
    "timestamp": 1640995200
  },
  "data": { /* payload */ }
}
```

**Response Headers:**
```http
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
X-Protocol-Version: 1.0.0
X-Server-Version: 1.0.2
X-Processing-Time: 3500
```

### gRPC

**Protocol Buffers:**
```protobuf
message RequestMetadata {
  string request_id = 1;
  string protocol_version = 2;
  string client_version = 3;
  string client_id = 4;
  string client_type = 5;
  int64 timestamp = 6;
  map<string, string> custom_headers = 7;
}

message ResponseMetadata {
  string request_id = 1;
  string protocol_version = 2;
  string server_version = 3;
  int64 timestamp = 4;
  int32 processing_time_ms = 5;
}
```

**gRPC Metadata:**
```go
md := metadata.New(map[string]string{
    "request-id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol-version": "1.0.0",
    "client-version": "1.0.0",
})
ctx := metadata.NewOutgoingContext(ctx, md)
```

### WebSocket

**Message Format:**
```json
{
  "type": "context_aware_template",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {
    "query": "хочу борщ"
  },
  "timestamp": "2025-01-18T10:00:00Z"
}
```

## Валидация

### JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "request_id": {
      "type": "string",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
    },
    "protocol_version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9.-]+)?(\\+[a-zA-Z0-9.-]+)?$"
    },
    "client_version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9.-]+)?(\\+[a-zA-Z0-9.-]+)?$"
    },
    "timestamp": {
      "type": "integer",
      "minimum": 0
    }
  },
  "required": ["request_id", "protocol_version", "client_version", "timestamp"]
}
```

## Best Practices

### 1. Генерация request_id

```javascript
// Используйте UUID v4
const requestId = uuid.v4();

// Или используйте библиотеку для генерации
const requestId = crypto.randomUUID();
```

### 2. Временные метки

```javascript
// Используйте Unix timestamp (секунды)
const timestamp = Math.floor(Date.now() / 1000);

// Или используйте ISO 8601 для WebSocket
const timestamp = new Date().toISOString();
```

### 3. Версионирование

```javascript
// Всегда указывайте версию протокола
const metadata = {
  protocol_version: "1.0.0", // Версия протокола
  client_version: package.version, // Версия клиента
};
```

### 4. Кастомные заголовки

```javascript
// Используйте префикс x- для кастомных заголовков
const customHeaders = {
  "x-feature-flag": "new-ui",
  "x-experiment-id": "exp-456",
  "x-debug-mode": "true"
};
```

## См. также

- [Формат сообщений](./MESSAGE_FORMAT.md) - общая структура сообщений
- [Обработка ошибок](./ERROR_HANDLING.md) - метаданные в ошибках
- [Версионирование](../versioning/README.md) - правила версий
