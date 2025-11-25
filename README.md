# Nexus Application Protocol v2.0.0 ✨

**Nexus Protocol** - Application Protocol для обмена данными между клиентами и серверами Nexus AI Platform.

## 🚀 Enterprise Ready

**Nexus Protocol теперь поддерживает enterprise-сценарии среднего и крупного бизнеса!**

### ✨ Новые возможности в v2.0.0
- **Multi-tenant архитектура** с полной изоляцией данных
- **Batch операции** для высокой производительности
- **Enterprise метрики** (rate limiting, кэширование, квоты)
- **Webhooks** для асинхронной обработки
- **Расширенная аналитика** для бизнес-решений
- **Локализация** и поддержка регионов

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
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
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

#### Пример 1: Информационный запрос
```bash
curl -X POST https://api.nexus.dev/api/v1/templates/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "query": "хочу борщ",
    "language": "ru",
    "metadata": {
      "request_id": "req-123",
      "protocol_version": "2.0.0"
    }
  }'
```

#### Пример 2: Запрос с покупкой и геолокацией
```bash
curl -X POST https://api.nexus.dev/api/v1/templates/execute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "query": "Найди где рядом продается кокакола и купи литровую бутылку колы заберу самостоятельно",
    "language": "ru",
    "context": {
      "user_id": "user-123",
      "location": {
        "latitude": 55.7558,
        "longitude": 37.6173,
        "accuracy": 50
      },
      "locale": "ru-RU",
      "currency": "RUB",
      "region": "RU"
    },
    "metadata": {
      "request_id": "req-456",
      "protocol_version": "2.0.0",
      "client_version": "2.0.0"
    }
  }'
```

**Ответ:**
```json
{
  "data": {
    "execution_id": "exec-789",
    "status": "completed",
    "query_type": "with_purchases_services",
    "sections": [
      {
        "domain_id": "commerce",
        "title": "Коммерческие предложения",
        "status": "success",
        "results": [
          {
            "id": "product-456",
            "type": "product_purchase",
            "title": "Coca-Cola 1л бутылка",
            "description": "Найдено в 3 магазинах рядом",
            "data": {
              "price": "89 ₽",
              "stores": [
                {
                  "name": "Пятерочка",
                  "distance": "200м",
                  "address": "ул. Ленина, 15",
                  "pickup_available": true,
                  "work_hours": "Круглосуточно"
                }
              ]
            },
            "relevance": 0.95,
            "actions": [
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
    ]
  },
  "metadata": {
    "request_id": "req-456",
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 245
  }
}
```

### gRPC

```go
client := pb.NewContextAwareTemplatesClient(conn)

resp, err := client.ExecuteTemplate(ctx, &pb.ExecuteTemplateRequest{
    Query: "хочу борщ",
    Metadata: &pb.RequestMetadata{
        RequestId:      uuid.New().String(),
        Version:        "2.0.0", // Nexus Protocol v2.0.0
        ClientVersion:  "2.0.0",
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

## 🎯 Примеры использования

### Простой информационный запрос
```bash
curl -X POST https://api.nexus.dev/api/v1/templates/execute \
  -H "Authorization: Bearer <token>" \
  -d '{"query": "хочу борщ", "language": "ru"}'
```

### Комплексный многошаговый сценарий
```bash
curl -X POST https://api.nexus.dev/api/v1/templates/execute \
  -H "Authorization: Bearer <token>" \
  -d '{
    "query": "закажи в макдоналдсе карточку фри, оплати, введи адрес доставки, и напоминай когда курьер выедет с заказом выпить таблетки, и через два часа выпить еще одни таблетки",
    "language": "ru",
    "context": {
      "user_id": "user-123",
      "location": {"latitude": 55.7558, "longitude": 37.6173}
    }
  }'
```

Система автоматически:
- ✅ Определяет несколько доменов (commerce, payment, delivery, notifications)
- ✅ Создает workflow с зависимостями между шагами
- ✅ Обрабатывает последовательность действий
- ✅ Создает напоминания с правильными триггерами

[📖 Подробные примеры →](./PURCHASE_EXAMPLES.md)

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

## 🎨 Frontend Configuration

Клиенты могут получать активную конфигурацию визуала (тема, цвета, layout, брендинг) через публичный endpoint:

```bash
GET /api/v1/frontend/config
```

**Пример ответа:**
```json
{
  "data": {
    "id": "frontend-config-001",
    "name": "Corporate Theme",
    "theme": "light",
    "colors": {
      "primary": "#0066CC",
      "secondary": "#00CC66",
      "accent": "#FF6600"
    },
    "branding": {
      "logo": "https://cdn.example.com/logo.png",
      "name": "Nexus Protocol"
    }
  }
}
```

**Использование в SDK:**
```go
config, err := client.GetFrontendConfig(ctx)
// Применить конфигурацию в UI
```

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

- **Protocol Version:** 2.0.0
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

**Версия:** 2.0.0
**Дата:** 2025-01-18
**Автор:** Биркин Максим