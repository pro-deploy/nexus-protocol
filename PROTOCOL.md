# Nexus Application Protocol

## Что такое Application Protocol?

**Application Protocol** - это формат сообщений и правила обмена данными, которые работают **поверх** транспортных протоколов (HTTP, gRPC, WebSocket).

Application Protocol определяет:
- ✅ **Структуру сообщений** - единый формат для всех запросов и ответов
- ✅ **Метаданные** - стандартизированные RequestMetadata/ResponseMetadata
- ✅ **Обработку ошибок** - единый формат ошибок
- ✅ **Версионирование** - правила совместимости версий

## Формат Application Protocol

### Базовый формат сообщения

Все сообщения в Application Protocol следуют единой структуре:

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "1.1.0",
    "client_version": "1.0.0",
    "timestamp": 1640995200
  },
  "data": {
    // Payload зависит от операции
  }
}
```

### RequestMessage (Запрос)

```json
{
  "metadata": {
    "request_id": "uuid",
    "protocol_version": "1.1.0",
    "client_version": "1.0.0",
    "timestamp": 1640995200
  },
  "data": {
    "query": "хочу борщ",
    "language": "ru"
  }
}
```

### ResponseMessage (Ответ)

```json
{
  "metadata": {
    "request_id": "uuid",
    "protocol_version": "1.1.0",
    "server_version": "1.0.2",
    "timestamp": 1640995235,
    "processing_time_ms": 3500
  },
  "data": {
    "execution_id": "exec-123",
    "status": "completed"
  }
}
```

## Application Protocol vs Transport Protocol

### Application Protocol (уровень приложения)
- Формат сообщений: `{metadata, data}`
- Метаданные: RequestMetadata/ResponseMetadata
- Обработка ошибок: ErrorDetail
- Версионирование: Semantic Versioning

### Transport Protocol (уровень транспорта)
- HTTP REST: JSON в теле запроса
- gRPC: Protocol Buffers
- WebSocket: JSON сообщения

## Применение в REST API

В REST API Application Protocol применяется следующим образом:

### Запрос
Метаданные могут быть:
1. В теле запроса (если требуется строгий формат)
2. В заголовках HTTP (для удобства REST API)
3. Автоматически добавлены SDK

**Текущая реализация SDK:**
- Метаданные автоматически добавляются в RequestMetadata (Application Protocol уровень)
- Запрос отправляется с метаданными в структуре запроса (Transport Protocol уровень)
- Это удобно для REST API и сохраняет совместимость
- Для строгого формата доступны типы `protocol.RequestMessage`/`protocol.ResponseMessage`

### Ответ
Ответ **всегда** приходит в формате Application Protocol:
```json
{
  "metadata": { ... },
  "data": { ... }
}
```

## Версионирование

Application Protocol использует Semantic Versioning:
- **MAJOR** - несовместимые изменения формата
- **MINOR** - обратно совместимые новые функции
- **PATCH** - исправления ошибок

Текущая версия: **1.1.0**

## См. также

- [Формат сообщений](./protocol/MESSAGE_FORMAT.md)
- [Метаданные](./protocol/METADATA.md)
- [Обработка ошибок](./protocol/ERROR_HANDLING.md)
- [Версионирование](./protocol/VERSIONING.md)

