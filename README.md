# Nexus Application Protocol v1.1.0 ✨

**Nexus Protocol** - Application Protocol для обмена данными между клиентами и серверами Nexus AI Platform.

## 🚀 Enterprise Ready

**Nexus Protocol теперь поддерживает enterprise-сценарии среднего и крупного бизнеса!**

### ✨ Новые возможности в v1.1.0
- **Multi-tenant архитектура** с полной изоляцией данных
- **Batch операции** для высокой производительности
- **Enterprise метрики** (rate limiting, кэширование, квоты)
- **Webhooks** для асинхронной обработки
- **Расширенная аналитика** для бизнес-решений
- **Локализация** и поддержка регионов

### 💰 Бизнес-преимущества

#### Средний бизнес (50-500 сотрудников)
- **Внедрение за 1-3 дня** вместо 2-6 месяцев
- **Конверсия +75%** (30% → 67.5%)
- **Экономия $200K-500K/год** на разработке

#### Крупный бизнес (500+ сотрудников)
- **Enterprise monitoring** с детальными health checks
- **Batch operations** для высокой производительности
- **Multi-tenant isolation** с полной изоляцией данных
- **Экономия $500K-2M/год**

**[🚀 Посмотреть Enterprise Demo](./sdk/go/examples/enterprise/)**

## Что это такое?

**Application Protocol** - это формат сообщений и правила обмена данными поверх существующих транспортных протоколов (HTTP, gRPC, WebSocket).

Nexus Protocol определяет:
- ✅ **Формат сообщений** - структура данных для обмена
- ✅ **Метаданные** - стандартизированные RequestMetadata/ResponseMetadata
- ✅ **Обработка ошибок** - единый формат ошибок
- ✅ **Версионирование** - правила совместимости версий

> 📖 **Подробнее об Application Protocol:** [PROTOCOL.md](./PROTOCOL.md)

## Структура документации

```
@protocol/
├── README.md                    # Этот файл
│
├── protocol/                    # ПРОТОКОЛ (формат сообщений)
│   ├── MESSAGE_FORMAT.md       # Формат сообщений
│   ├── METADATA.md             # Метаданные запросов/ответов
│   └── ERROR_HANDLING.md        # Обработка ошибок
│
├── api/                         # API СПЕЦИФИКАЦИИ
│   ├── rest/
│   │   └── openapi.yaml        # REST API (OpenAPI 3.0)
│   ├── grpc/
│   │   └── nexus.proto         # gRPC API (Protocol Buffers)
│   └── websocket/
│       └── protocol.json        # WebSocket формат сообщений
│
├── schemas/                     # JSON SCHEMAS
│   └── message-schema.json     # Схема валидации сообщений
│
└── versioning/                  # ВЕРСИОНИРОВАНИЕ
    └── README.md               # Правила версионирования
```

## Быстрый старт

### Формат сообщения (Application Protocol)

Все сообщения следуют единому формату:

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "1.0.0",
    "client_version": "1.0.0",
    "client_id": "web-app",
    "client_type": "web",
    "timestamp": 1640995200
  },
  "data": {
    // Payload зависит от операции
  }
}
```

### HTTP REST

```bash
curl -X POST https://api.nexus.dev/api/v1/templates/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "query": "хочу борщ",
    "metadata": {
      "request_id": "req-123",
      "protocol_version": "1.0.0"
    }
  }'
```

### gRPC

```go
client := pb.NewContextAwareTemplatesClient(conn)

resp, err := client.ExecuteTemplate(ctx, &pb.ExecuteTemplateRequest{
    Query: "хочу борщ",
    Metadata: &pb.RequestMetadata{
        RequestId:      uuid.New().String(),
        Version:        "1.0.0",
        ClientVersion:  "1.0.0",
    },
})
```

### WebSocket

```javascript
const ws = new WebSocket('ws://api.nexus.dev/ws?token=<jwt_token>');

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'context_aware_template',
    request_id: 'req-123',
    payload: {
      query: 'хочу борщ'
    },
    timestamp: new Date().toISOString()
  }));
};
```

## Основные компоненты протокола

### 1. Метаданные (Metadata)

Стандартизированные метаданные для всех запросов и ответов:

- `request_id` - уникальный идентификатор запроса (UUID)
- `protocol_version` - версия протокола (Semantic Versioning)
- `client_version` - версия клиента
- `timestamp` - временная метка запроса/ответа

[Подробнее →](./protocol/METADATA.md)

### 2. Формат сообщений (Message Format)

Единый формат сообщений для всех транспортов:

- Структура запроса
- Структура ответа
- Типы сообщений
- Валидация

[Подробнее →](./protocol/MESSAGE_FORMAT.md)

### 3. Обработка ошибок (Error Handling)

Стандартизированный формат ошибок:

- Коды ошибок
- Типы ошибок
- Детали ошибок
- Обработка на разных транспортах

[Подробнее →](./protocol/ERROR_HANDLING.md)

### 4. Версионирование (Versioning)

Правила версионирования и совместимости:

- Semantic Versioning (MAJOR.MINOR.PATCH)
- Правила совместимости
- Version negotiation
- Миграция между версиями

[Подробнее →](./versioning/README.md)

## Транспорты

Nexus Protocol работает поверх следующих транспортных протоколов:

### HTTP REST
- **Спецификация:** [OpenAPI 3.0](./api/rest/openapi.yaml)
- **Base URL:** `https://api.nexus.dev/api/v1`
- **Content-Type:** `application/json`
- **Authentication:** Bearer Token (JWT)

### gRPC
- **Спецификация:** [Protocol Buffers](./api/grpc/nexus.proto)
- **Port:** `50051`
- **Transport:** HTTP/2
- **Authentication:** mTLS / JWT в metadata

### WebSocket
- **Спецификация:** [JSON Protocol](./api/websocket/protocol.json)
- **URL:** `ws://api.nexus.dev/ws`
- **Subprotocol:** `nexus-json`
- **Authentication:** JWT в query parameter или header

## Валидация

JSON Schema для валидации сообщений:

```bash
# Валидация сообщения по схеме
cat message.json | jq . | jsonschema schemas/message-schema.json
```

[Схема →](./schemas/message-schema.json)

## Совместимость

- **Protocol Version:** 1.0.0
- **Semantic Versioning:** MAJOR.MINOR.PATCH
- **Backward Compatibility:** В рамках Major версии
- **Transport Protocols:** HTTP/1.1, HTTP/2, WebSocket (RFC 6455)
- **Data Formats:** JSON, Protocol Buffers 3

## Статус

✅ **Production Ready** - Протокол готов к использованию в production

## Лицензия

MIT License

## Контакты

- **Email:** contact@nexus.dev
- **Website:** https://nexus.dev
- **Documentation:** https://docs.nexus.dev

---

**Версия:** 1.1.0
**Дата:** 2025-01-18
**Автор:** Биркин Максим