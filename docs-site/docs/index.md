---
id: index
title: Nexus Protocol - Документация
sidebar_label: Главная
slug: /
---

# Nexus Application Protocol v2.0.0

## Предыстория и задачи

### Проблема

В современной разработке приложений с AI и интеграцией различных сервисов разработчики сталкиваются с множеством вызовов:

- **Фрагментация протоколов** - каждый сервис использует свой формат данных и API
- **Сложность интеграции** - необходимо писать отдельный код для каждого сервиса
- **Отсутствие стандартизации** - нет единого подхода к обработке ошибок, метаданных, версионированию
- **Высокая стоимость разработки** - каждый новый интеграционный проект требует месяцев разработки
- **Проблемы масштабирования** - сложно масштабировать систему при росте нагрузки
- **Отсутствие контекста** - сервисы не понимают контекст пользователя (геолокация, предпочтения, история)

### Решение

**Nexus Application Protocol** - это единый протокол обмена данными, который решает все эти проблемы:

✅ **Единый формат сообщений** для всех транспортных протоколов (HTTP, gRPC, WebSocket)  
✅ **Стандартизированные метаданные** для трассировки, версионирования и аналитики  
✅ **Унифицированная обработка ошибок** с понятными кодами и типами  
✅ **Контекстно-зависимые запросы** с поддержкой геолокации, локализации, валюты  
✅ **Enterprise возможности** - batch операции, webhooks, аналитика, multi-tenant  
✅ **Быстрое внедрение** - SDK для популярных языков программирования  

### Какие задачи решает Nexus Protocol?

#### 1. Унификация интеграций

**Проблема:** Каждый сервис требует свой формат запросов и обработки ответов.

**Решение:** Единый формат Application Protocol работает поверх любых транспортных протоколов. Один раз интегрируете SDK - работаете со всеми сервисами одинаково.

```go
// Один и тот же код для всех сервисов
req := &types.ExecuteTemplateRequest{
    Query: "хочу борщ",
    Language: "ru",
}
result, err := client.ExecuteTemplate(ctx, req)
```

#### 2. Контекстно-зависимые запросы

**Проблема:** Сервисы не знают контекст пользователя (где он находится, какая валюта, язык, предпочтения).

**Решение:** Nexus Protocol передает полный контекст пользователя в каждом запросе:

```go
Context: &types.UserContext{
    UserID:    "user-123",
    Location:  &types.UserLocation{Latitude: 55.7558, Longitude: 37.6173},
    Locale:    "ru-RU",
    Currency:  "RUB",
    Region:    "RU",
}
```

#### 3. Многошаговые сценарии (Workflow)

**Проблема:** Сложные задачи требуют координации нескольких сервисов (заказ → оплата → доставка → уведомления).

**Решение:** Nexus Protocol автоматически создает workflow с зависимостями между шагами:

```go
Query: "закажи в макдоналдсе карточку фри, оплати, введи адрес доставки, 
        и напоминай когда курьер выедет с заказом"
// Система автоматически создает workflow:
// 1. Заказ еды (commerce)
// 2. Оплата (payment) - зависит от шага 1
// 3. Доставка (delivery) - зависит от шага 2
// 4. Напоминания (notifications) - зависит от шага 3
```

#### 4. Enterprise масштабирование

**Проблема:** При росте нагрузки система не справляется, нет метрик и мониторинга.

**Решение:** Встроенные enterprise возможности:

- **Batch операции** - параллельное выполнение множественных запросов
- **Webhooks** - асинхронная обработка результатов
- **Rate limiting** - контроль нагрузки
- **Кэширование** - оптимизация производительности
- **Аналитика** - метрики конверсии и производительности
- **Multi-tenant** - изоляция данных по клиентам

#### 5. Быстрое внедрение

**Проблема:** Интеграция нового сервиса занимает месяцы разработки.

**Решение:** 
- **Внедрение за 1-3 дня** вместо 2-6 месяцев
- **Готовые SDK** для популярных языков
- **Автоматическая генерация метаданных**
- **Встроенная обработка ошибок**

#### 6. Версионирование и совместимость

**Проблема:** Обновление API ломает существующие интеграции.

**Решение:** Semantic Versioning с правилами совместимости:
- MAJOR версия - несовместимые изменения
- MINOR версия - обратно совместимые новые функции
- PATCH версия - исправления ошибок

#### 7. Единая обработка ошибок

**Проблема:** Каждый сервис возвращает ошибки в своем формате.

**Решение:** Стандартизированный формат ошибок с типами и кодами:

```go
if errDetail.IsValidationError() {
    // Ошибка валидации
} else if errDetail.IsRateLimitError() {
    // Превышен лимит запросов
}
```

---

# Протокол


## Формат сообщений

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

---


## Метаданные

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
  "protocol_version": "2.0.0",
  "client_version": "2.0.0",
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

**Тип:** object (map&lt;string, string&gt;)  
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
  "protocol_version": "2.0.0",
  "server_version": "2.0.0",
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
    protocol_version: "2.0.0",
    client_version: "2.0.0",
    timestamp: Date.now() / 1000
  },
  data: { query: "хочу борщ" }
};

// Сервер возвращает ответ с тем же request_id
const response = {
  metadata: {
    request_id: requestId, // Тот же ID
    protocol_version: "2.0.0",
    server_version: "2.0.0",
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
  protocol_version: "2.0.0", // Версия протокола клиента
  client_version: "2.0.0",   // Версия клиентского приложения
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
X-Protocol-Version: 2.0.0
X-Client-Version: 2.0.0
```

**Request Body:**
```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "timestamp": 1640995200
  },
  "data": { /* payload */ }
}
```

**Response Headers:**
```http
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
X-Protocol-Version: 2.0.0
X-Server-Version: 2.0.0
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
    "protocol-version": "2.0.0",
    "client-version": "2.0.0",
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
  protocol_version: "2.0.0", // Версия протокола
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

---


## Обработка ошибок протокола

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

---


## Версионирование

## Обзор

Nexus Application Protocol использует **Semantic Versioning** (SemVer) для управления версиями протокола. Версионирование обеспечивает:

- ✅ Предсказуемое развитие протокола
- ✅ Обратную совместимость
- ✅ Понятные правила обновления
- ✅ Безопасную миграцию между версиями

## Формат версии

### Semantic Versioning

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
│     │    │        │           │
│     │    │        │           └─ Build metadata (опционально)
│     │    │        └───────────── Prerelease identifier (опционально)
│     │    └────────────────────── Patch: обратно совместимые исправления
│     └────────────────────────── Minor: обратно совместимые новые функции
└──────────────────────────────── Major: несовместимые изменения
```

### Примеры версий

- `1.0.0` - первая стабильная версия
- `1.1.0` - добавлена новая функция (обратно совместимая)
- `1.1.1` - исправлена ошибка (обратно совместимая)
- `2.0.0` - несовместимые изменения
- `1.0.0-alpha.1` - альфа-версия
- `1.0.0-rc.1+build.123` - release candidate

## Текущая версия

**Protocol Version:** `2.0.0` ✨ (Advanced Enterprise Edition)

## Правила версионирования

### MAJOR версия (несовместимые изменения)

Увеличивается при:
- ❌ Изменении формата сообщений (RequestMetadata, ResponseMetadata)
- ❌ Удалении полей из сообщений
- ❌ Изменении кодов ошибок
- ❌ Изменении обязательности полей
- ❌ Изменении поведения существующих функций

**Примеры:**
- `1.0.0` → `2.0.0`: Изменена структура RequestMetadata
- `1.0.0` → `2.0.0`: Удалено поле `client_id` из RequestMetadata
- `1.0.0` → `2.0.0`: Изменен формат ErrorDetail

### MINOR версия (обратно совместимые новые функции)

Увеличивается при:
- ✅ Добавлении новых опциональных полей
- ✅ Добавлении новых типов сообщений
- ✅ Добавлении новых кодов ошибок
- ✅ Расширении существующих enum значений

**Примеры:**
- `1.0.0` → `1.1.0`: Добавлено опциональное поле `custom_headers` в RequestMetadata
- `1.0.0` → `1.1.0`: Добавлен новый тип сообщения `subscribe`
- `1.0.0` → `1.1.0`: Добавлен новый код ошибки `TIMEOUT`

### PATCH версия (обратно совместимые исправления)

Увеличивается при:
- ✅ Исправлении ошибок в документации
- ✅ Уточнении описаний полей
- ✅ Исправлении примеров
- ✅ Улучшении валидации (без изменения поведения)

**Примеры:**
- `1.0.0` → `1.0.1`: Исправлена опечатка в документации
- `1.0.0` → `1.0.1`: Уточнено описание поля `timestamp`
- `1.0.0` → `1.0.1`: Улучшена валидация UUID

## Совместимость версий

### Правила совместимости

Версии совместимы если:

1. **Major версии совпадают** - несовместимые изменения только в Major
2. **Minor версия клиента ≤ Minor версии сервера** - клиент не может быть новее сервера

### Матрица совместимости

| Client Version | Server Version | Compatible | Notes |
|----------------|----------------|------------|-------|
| 1.0.0         | 1.0.0         | ✅        | Полная совместимость |
| 1.0.0         | 1.1.0         | ✅        | Server имеет дополнительные функции |
| 1.1.0         | 1.0.0         | ❌        | Client ожидает функций, которых нет |
| 1.0.0         | 2.0.0         | ❌        | Несовместимые изменения |
| 2.0.0         | 1.0.0         | ❌        | Client несовместим |

### Проверка совместимости

```javascript
function isCompatible(clientVersion, serverVersion) {
  const [clientMajor, clientMinor] = clientVersion.split('.').map(Number);
  const [serverMajor, serverMinor] = serverVersion.split('.').map(Number);
  
  // Major версии должны совпадать
  if (clientMajor !== serverMajor) {
    return false;
  }
  
  // Minor версия клиента не должна быть больше сервера
  if (clientMinor > serverMinor) {
    return false;
  }
  
  return true;
}
```

## Version Negotiation

### Request Metadata

Клиент указывает версию протокола в RequestMetadata:

```json
{
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "1.0.0",
    "client_version": "1.0.0",
    "timestamp": 1640995200
  }
}
```

### Response Metadata

Сервер возвращает версию протокола в ResponseMetadata:

```json
{
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "1.0.0",
    "server_version": "1.0.2",
    "timestamp": 1640995235,
    "processing_time_ms": 3500
  }
}
```

### Проверка на сервере

```go
func validateProtocolVersion(clientVersion string) error {
    serverVersion := "1.0.0"
    
    if !isCompatible(clientVersion, serverVersion) {
        return &ProtocolVersionError{
            ClientVersion: clientVersion,
            ServerVersion: serverVersion,
            Message: "Protocol version mismatch",
        }
    }
    
    return nil
}
```

## Миграция между версиями

### Миграция 1.0.0 → 1.1.0 (Minor)

**Изменения:**
- Добавлено опциональное поле `custom_headers` в RequestMetadata
- Добавлен новый тип сообщения `subscribe`

**Миграция:**
1. Обновить клиент до версии 1.1.0
2. Использовать новые функции (опционально)
3. Обратная совместимость гарантирована

### Миграция 1.x.x → 2.0.0 (Major)

**Изменения:**
- Изменена структура RequestMetadata
- Удалено поле `client_id`
- Изменен формат ErrorDetail

**Миграция:**
1. **Audit:** Проверить использование устаревших полей
2. **Update:** Обновить код для новой структуры
3. **Test:** Протестировать в staging среде
4. **Deploy:** Развернуть с rollback plan

## Deprecation Policy

### Процесс устаревания

1. **Announcement:** Поле/функция помечается как deprecated в релизе
2. **Documentation:** Обновляется документация с предупреждением
3. **Alternative:** Предоставляется альтернативный способ
4. **Removal:** Удаляется через 2 major версии

### Пример

```json
// v1.0.0
{
  "metadata": {
    "request_id": "req-123",
    "client_id": "web-app"  // Актуально
  }
}

// v1.1.0 - deprecated
{
  "metadata": {
    "request_id": "req-123",
    "client_id": "web-app",  // deprecated, используйте client_type
    "client_type": "web"     // Новое поле
  }
}

// v2.0.0 - удалено
{
  "metadata": {
    "request_id": "req-123",
    "client_type": "web"     // client_id удален
  }
}
```

## Version Endpoints

### GET /api/v1/version

**Response:**
```json
{
  "protocol_version": "1.0.0",
  "server_version": "1.0.2",
  "api_version": "v1",
  "build_info": {
    "git_commit": "abc123def456",
    "build_time": "2025-01-18T10:00:00Z"
  }
}
```

## Best Practices

### 1. Всегда указывайте protocol_version

```json
{
  "metadata": {
    "protocol_version": "1.0.0"  // Обязательно
  }
}
```

### 2. Проверяйте совместимость на клиенте

```javascript
const clientVersion = "1.0.0";
const serverVersion = response.metadata.protocol_version;

if (!isCompatible(clientVersion, serverVersion)) {
  console.error("Protocol version mismatch");
  // Обработка несовместимости
}
```

### 3. Используйте Semantic Versioning

```javascript
// ✅ Правильно
const version = "1.0.0";

// ❌ Неправильно
const version = "1.0";
const version = "v1.0.0";
const version = "1.0.0-beta";
```

## Изменения в версии 1.1.0 (Enterprise Edition)

**Дата релиза:** 2024-12-01
**Тип:** MINOR (обратная совместимость сохранена)

### 🎯 Enterprise фичи

#### Метрики и мониторинг
- ✅ **ResponseMetadata** расширена полями `rate_limit_info`, `cache_info`, `quota_info`
- ✅ **ReadinessResponse** с детальным статусом компонентов и емкостью системы
- ✅ **Performance metrics** в аналитике (перцентили, throughput, error rates)

#### Batch операции и производительность
- ✅ **Batch operations** для параллельного выполнения множественных запросов
- ✅ **Webhook integration** для асинхронной обработки результатов
- ✅ **Priority management** через custom_headers (low, normal, high, critical)

#### Расширенный поиск и фильтрация
- ✅ **Advanced filters** в ExecuteTemplate (домены, релевантность, дата, сортировка)
- ✅ **Pagination support** с курсорами и метаданными страниц
- ✅ **Cache control** через custom_headers (cache-first, network-first, etc.)

#### Локализация и мультитенантность
- ✅ **UserContext** расширен полями `locale`, `timezone`, `currency`, `region`
- ✅ **Multi-tenant isolation** через `tenant_id` в контексте
- ✅ **A/B testing** поддержка через `x-experiment-id` header

#### Аналитика для бизнеса
- ✅ **Conversion metrics** (поиск→результат, результат→действие, удержание)
- ✅ **Domain breakdown** с метриками по каждому домену
- ✅ **Business intelligence** метрики для принятия решений

### Миграция на 1.1.0

**Для существующих клиентов:**
1. Обновите SDK до версии 1.1.0
2. Установите `protocol_version: "1.1.0"` в RequestMetadata
3. Используйте новые enterprise фичи (опционально)
4. **Обратная совместимость гарантирована**

**Пример обновления клиента:**
```go
cfg := client.Config{
    BaseURL:         "https://api.company.com",
    ProtocolVersion: "1.1.0", // Обновлено с 1.0.0
    // ... остальные настройки
}
```

### Бизнес-преимущества 1.1.0

#### Средний бизнес (50-500 сотрудников)
- **Внедрение:** 1-3 дня вместо 2-6 месяцев
- **Конверсия:** +75% (30% → 67.5%)
- **Экономия:** $200K-500K/год

#### Крупный бизнес (500+ сотрудников)
- **Multi-tenant:** полная изоляция данных
- **Enterprise monitoring:** детальные health checks
- **Batch operations:** высокая производительность
- **Экономия:** $500K-2M/год

## См. также

- [Формат сообщений](./MESSAGE_FORMAT.md) - структура сообщений
- [Метаданные](./METADATA.md) - protocol_version в метаданных
- [Обработка ошибок](./ERROR_HANDLING.md) - PROTOCOL_VERSION_ERROR
- [Enterprise Demo](../../sdk/go/examples/enterprise/) - примеры использования enterprise фич

---


## REST API

Nexus Protocol поддерживает HTTP REST API для взаимодействия с сервером.

## Базовый URL

```
https://api.nexus.dev/api/v1
```

## Аутентификация

Все запросы требуют JWT токен в заголовке:

```
Authorization: Bearer <jwt_token>
```

## Content-Type

Все запросы должны использовать:

```
Content-Type: application/json
```

## Формат запросов

Запросы следуют формату Application Protocol:

```json
{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    // Операция-специфичные данные
  }
}
```

## Формат ответов

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

---


## gRPC API

Nexus Protocol поддерживает gRPC для высокопроизводительного взаимодействия.

## Подключение

```
api.nexus.dev:50051
```

## Protocol Buffers

Спецификация gRPC API определена в файле [nexus.proto](../../api/grpc/nexus.proto).

## Аутентификация

gRPC поддерживает два метода аутентификации:

1. **mTLS** - взаимная аутентификация с сертификатами
2. **JWT в metadata** - передача JWT токена в gRPC metadata

### Пример с JWT:

```go
import (
    "google.golang.org/grpc/metadata"
)

md := metadata.New(map[string]string{
    "authorization": "Bearer <jwt_token>",
})
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

## Пример использования

```go
import (
    pb "github.com/nexus-protocol/api/grpc"
)

conn, err := grpc.Dial("api.nexus.dev:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewContextAwareTemplatesClient(conn)

resp, err := client.ExecuteTemplate(ctx, &pb.ExecuteTemplateRequest{
    Query: "хочу борщ",
    Metadata: &pb.RequestMetadata{
        RequestId:      uuid.New().String(),
        ProtocolVersion: "2.0.0",
        ClientVersion:  "2.0.0",
    },
})
```

## Преимущества gRPC

- Высокая производительность (HTTP/2, бинарный протокол)
- Типобезопасность (Protocol Buffers)
- Поддержка стриминга
- Автоматическая генерация клиентов

---


## WebSocket API

Nexus Protocol поддерживает WebSocket для двусторонней связи в реальном времени.

## Подключение

```
ws://api.nexus.dev/ws?token=<jwt_token>
```

или

```
wss://api.nexus.dev/ws?token=<jwt_token>
```

## Subprotocol

Используйте subprotocol `nexus-json`:

```javascript
const ws = new WebSocket('wss://api.nexus.dev/ws?token=<jwt_token>', ['nexus-json']);
```

## Формат сообщений

Все сообщения следуют формату Application Protocol:

### Запрос

```json
{
  "type": "execute_template",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "payload": {
    "query": "хочу борщ",
    "language": "ru"
  },
  "timestamp": 1640995200
}
```

### Ответ

```json
{
  "type": "execute_template_response",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 3500
  },
  "payload": {
    "execution_id": "exec-123",
    "status": "completed"
  },
  "timestamp": 1640995235
}
```

## Пример использования

```javascript
const ws = new WebSocket('wss://api.nexus.dev/ws?token=<jwt_token>', ['nexus-json']);

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'execute_template',
    request_id: 'req-123',
    metadata: {
      protocol_version: '2.0.0',
      client_version: '2.0.0'
    },
    payload: {
      query: 'хочу борщ',
      language: 'ru'
    },
    timestamp: Date.now() / 1000
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket closed');
};
```

## Типы сообщений

- `execute_template` - выполнение шаблона
- `execute_template_response` - ответ на выполнение шаблона
- `status_update` - обновление статуса выполнения
- `error` - сообщение об ошибке

## Спецификация

Полная спецификация доступна в файле [protocol.json](../../api/websocket/protocol.json).

---


# SDK


## Установка SDK

## Требования

- Go 1.18 или выше
- Доступ к Nexus Protocol API

## Установка через go get

```bash
go get github.com/pro-deploy/nexus-protocol/sdk/go
```

## Установка через go.mod

Добавьте в ваш `go.mod`:

```go
module your-module

go 1.18

require (
    github.com/pro-deploy/nexus-protocol/sdk/go v2.0.0
)
```

Затем выполните:

```bash
go mod download
go mod tidy
```

## Импорт

```go
import (
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)
```

## Проверка установки

Создайте простой тестовый файл:

```go
package main

import (
    "fmt"
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
)

func main() {
    cfg := client.Config{
        BaseURL: "http://localhost:8080",
    }
    client := client.NewClient(cfg)
    fmt.Println("SDK установлен успешно!")
}
```

Запустите:

```bash
go run main.go
```

## Зависимости

SDK использует следующие зависимости:

- `github.com/google/uuid` - генерация UUID
- `github.com/xeipuuv/gojsonschema` - валидация JSON Schema (опционально)

Все зависимости устанавливаются автоматически при установке SDK.

---


## Быстрый старт SDK

Это руководство поможет вам быстро начать работу с Nexus Protocol SDK.

## Шаг 1: Установка

```bash
go get github.com/pro-deploy/nexus-protocol/sdk/go
```

## Шаг 2: Создание клиента

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
    // Создаем клиент
    cfg := client.Config{
        BaseURL:         "https://api.nexus.dev",
        Token:           "your-jwt-token",
        ProtocolVersion: "2.0.0",
        ClientVersion:   "2.0.0",
    }
    
    nexusClient := client.NewClient(cfg)
    ctx := context.Background()
    
    // Выполняем шаблон
    req := &types.ExecuteTemplateRequest{
        Query:    "хочу борщ",
        Language: "ru",
    }
    
    result, err := nexusClient.ExecuteTemplate(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Execution ID: %s\n", result.ExecutionID)
    fmt.Printf("Status: %s\n", result.Status)
}
```

## Шаг 3: Запуск

```bash
go run main.go
```

## Следующие шаги

- [Базовое использование](./basic-usage) - подробнее о создании клиента
- [Руководство по использованию](./usage-guide) - полное руководство
- [Примеры](./examples) - больше примеров кода

---


## Базовое использование SDK

## Создание клиента

### Минимальная конфигурация

```go
import (
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
)

cfg := client.Config{
    BaseURL: "https://api.nexus.dev",
    Token:   "your-jwt-token",
}

nexusClient := client.NewClient(cfg)
```

### Полная конфигурация

```go
import (
    "time"
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
)

cfg := client.Config{
    BaseURL:         "https://api.nexus.dev",
    Token:           "jwt-token",
    Timeout:         30 * time.Second,
    ProtocolVersion: "2.0.0",
    ClientVersion:   "2.0.0",
    ClientID:        "my-application",
    ClientType:      "web", // web, mobile, sdk, api, desktop
}

nexusClient := client.NewClient(cfg)
```

## Выполнение шаблона

### Простой запрос

```go
import (
    "context"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

ctx := context.Background()

req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
}

result, err := nexusClient.ExecuteTemplate(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Execution ID: %s\n", result.ExecutionID)
```

### С контекстом пользователя

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
    Context: &types.UserContext{
        UserID:    "user-123",
        SessionID: "session-456",
        Locale:    "ru-RU",
        Currency:  "RUB",
        Region:    "RU",
    },
}

result, err := nexusClient.ExecuteTemplate(ctx, req)
```

## Получение статуса выполнения

```go
status, err := nexusClient.GetExecutionStatus(ctx, "execution-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status)
```

## Обработка ошибок

```go
result, err := nexusClient.ExecuteTemplate(ctx, req)
if err != nil {
    if errDetail, ok := err.(*types.ErrorDetail); ok {
        fmt.Printf("Error: %s (%s)\n", errDetail.Message, errDetail.Code)
    } else {
        log.Printf("Unexpected error: %v", err)
    }
    return
}
```

## Изменение токена

```go
nexusClient.SetToken("new-token")
```

## Получение конфигурации фронтенда

```go
config, err := nexusClient.GetFrontendConfig(ctx)
if err == nil {
    fmt.Printf("Theme: %s\n", config.Theme)
    fmt.Printf("Primary Color: %s\n", config.Colors["primary"])
}
```

---


## Руководство по использованию SDK

## Быстрый старт

### 1. Установка

```bash
go get github.com/pro-deploy/nexus-protocol/sdk/go
```

### 2. Базовое использование

```go
package main

import (
    "fmt"
    "log"
    
    nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
    // Создаем клиент
    client := nexus.NewClient(nexus.Config{
        BaseURL: "http://localhost:8080",
        Token:   "your-jwt-token",
    })
    
    // Выполняем запрос
    result, err := client.ExecuteTemplate(&types.ExecuteTemplateRequest{
        Query:    "хочу борщ",
        Language: "ru",
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Execution ID: %s\n", result.ExecutionID)
}
```

## Конфигурация клиента

### Полная конфигурация

```go
cfg := nexus.Config{
    BaseURL:         "https://api.nexus.dev",
    Token:           "jwt-token",
    Timeout:         30 * time.Second,
    ProtocolVersion: "2.0.0", // Nexus Protocol v2.0.0
    ClientVersion:   "2.0.0",
    ClientID:        "my-application",
    ClientType:      "web", // web, mobile, sdk, api, desktop
}

client := nexus.NewClient(cfg)
```

### Изменение токена

```go
client.SetToken("new-token")
```

## Выполнение шаблонов

### Простой запрос

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
}

result, err := client.ExecuteTemplate(req)
```

### С контекстом пользователя

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
    Context: &types.UserContext{
        UserID:    "user-123",
        SessionID: "session-456",
        TenantID:  "tenant-789",
        Location: &types.UserLocation{
            Latitude:  55.7558,
            Longitude: 37.6173,
        },
    },
}
```

### Запрос с покупкой и геолокацией

```go
req := &types.ExecuteTemplateRequest{
    Query:    "Найди где рядом продается кокакола и купи литровую бутылку колы заберу самостоятельно",
    Language: "ru",
    Context: &types.UserContext{
        UserID:    "user-123",
        SessionID: "session-456",
        Location: &types.UserLocation{
            Latitude:  55.7558,  // Москва
            Longitude: 37.6173,
            Accuracy:  50,      // точность 50 метров
        },
        Locale:    "ru-RU",
        Currency:  "RUB",
        Region:    "RU",
    },
    Options: &types.ExecuteOptions{
        TimeoutMS:           30000,
        MaxResultsPerDomain: 10,
        ParallelExecution:   true,
    },
}

result, err := client.ExecuteTemplate(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Проверка типа запроса
if result.QueryType == "with_purchases_services" {
    fmt.Println("Запрос с возможностью покупки")
    
        // Обработка результатов из commerce домена
        for _, section := range result.Sections {
            if section.DomainID == "commerce" {
                for _, item := range section.Results {
                    fmt.Printf("Товар: %s\n", item.Title)
                    fmt.Printf("Релевантность: %.2f\n", item.Relevance)
                    
                    // Обработка данных товара (цена, магазины и т.д.)
                    if item.Data != nil {
                        if price, ok := item.Data["price"].(string); ok {
                            fmt.Printf("Цена: %s\n", price)
                        }
                        if stores, ok := item.Data["stores"]; ok {
                            fmt.Printf("Магазины: %v\n", stores)
                        }
                    }
                    
                    // Обработка действий (покупка, резервирование)
                    for _, action := range item.Actions {
                        fmt.Printf("Действие: %s - %s\n", action.Type, action.Label)
                    }
                }
            }
        }
}
```

### Многошаговый сценарий (заказ еды + оплата + доставка + напоминания)

```go
req := &types.ExecuteTemplateRequest{
    Query: "закажи в макдоналдсе карточку фри, оплати, введи адрес доставки, и напоминай когда курьер выедет с заказом выпить таблетки, и через два часа выпить еще одни таблетки",
    Language: "ru",
    Context: &types.UserContext{
        UserID: "user-123",
        Location: &types.UserLocation{
            Latitude:  55.7558,
            Longitude: 37.6173,
            Accuracy:  50,
        },
        Locale:   "ru-RU",
        Currency: "RUB",
        Region:   "RU",
    },
}

result, err := client.ExecuteTemplate(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Обработка многошагового workflow
if result.QueryType == "with_purchases_services" {
    fmt.Println("✅ Многошаговый сценарий обработан")
    
    // Работа с workflow (если доступен)
    if result.Workflow != nil {
        fmt.Println("\n📋 Workflow шаги:")
        steps := client.GetWorkflowSteps(result)
        for _, step := range steps {
            fmt.Printf("  Шаг %d: %s (%s) - статус: %s\n", 
                step.Step, step.Action, step.Domain, step.Status)
            if len(step.DependsOn) > 0 {
                fmt.Printf("    Зависит от: %v\n", step.DependsOn)
            }
        }
        
        // Получение следующего шага для выполнения
        nextStep := client.GetNextWorkflowStep(result)
        if nextStep != nil {
            fmt.Printf("\n➡️  Следующий шаг: %s (%s)\n", nextStep.Action, nextStep.Domain)
        }
        
        // Получение шагов по домену
        commerceSteps := client.GetWorkflowStepByDomain(result, "commerce")
        if len(commerceSteps) > 0 {
            fmt.Printf("\n🛒 Шаги commerce домена: %d\n", len(commerceSteps))
        }
    }
    
    // Обработка каждого домена
    for _, section := range result.Sections {
        switch section.DomainID {
        case "commerce":
            fmt.Println("\n🍔 Заказ еды:")
            for _, item := range section.Results {
                fmt.Printf("  - %s: %s\n", item.Title, item.Data["price"])
                // Выполнение заказа через action
                for _, action := range item.Actions {
                    if action.Type == "order_now" {
                        fmt.Printf("    → Действие: %s\n", action.Label)
                    }
                }
            }
            
        case "payment":
            fmt.Println("\n💳 Оплата:")
            for _, item := range section.Results {
                fmt.Printf("  - Сумма: %s\n", item.Data["amount"])
                for _, action := range item.Actions {
                    if action.Type == "process_payment" {
                        fmt.Printf("    → Действие: %s\n", action.Label)
                    }
                }
            }
            
        case "delivery":
            fmt.Println("\n🚚 Доставка:")
            for _, item := range section.Results {
                fmt.Printf("  - %s\n", item.Title)
                for _, action := range item.Actions {
                    fmt.Printf("    → Действие: %s\n", action.Label)
                }
            }
            
        case "notifications":
            fmt.Println("\n🔔 Напоминания:")
            for _, item := range section.Results {
                fmt.Printf("  - %s\n", item.Title)
                if item.Data != nil {
                    if reminderType, ok := item.Data["reminder_type"].(string); ok {
                        fmt.Printf("    Тип: %s\n", reminderType)
                    }
                    if delay, ok := item.Data["delay_hours"].(float64); ok {
                        fmt.Printf("    Задержка: %.0f часов\n", delay)
                    }
                }
                for _, action := range item.Actions {
                    fmt.Printf("    → Действие: %s\n", action.Label)
                }
            }
        }
    }
}
```

### С опциями выполнения

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
    Options: &types.ExecuteOptions{
        TimeoutMS:           30000,  // 30 секунд
        MaxResultsPerDomain: 5,
        ParallelExecution:   true,
        IncludeWebSearch:    true,
    },
}
```

### С кастомными метаданными

```go
metadata := types.NewRequestMetadata("2.0.0", "2.0.0") // Nexus Protocol v2.0.0
metadata.ClientID = "my-app"
metadata.ClientType = "web"
metadata.CustomHeaders = map[string]string{
    "x-feature-flag": "new-ui",
    "x-experiment-id": "exp-123",
}

req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Metadata: metadata,
}
```

## Получение статуса выполнения

```go
status, err := client.GetExecutionStatus("execution-id-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status)
fmt.Printf("Sections: %d\n", len(status.Sections))
```

## Streaming результатов

```go
resp, err := client.StreamTemplateResults("execution-id-123")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// Читаем Server-Sent Events
scanner := bufio.NewScanner(resp.Body)
for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "data: ") {
        // Парсим JSON из data: {...}
        data := line[6:]
        // Обработка данных
    }
}
```

## Получение конфигурации фронтенда

### Публичный endpoint (без аутентификации)

```go
// Получение активной конфигурации фронтенда
config, err := client.GetFrontendConfig(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Theme: %s\n", config.Theme)
fmt.Printf("Primary Color: %s\n", config.Colors["primary"])

// Использование конфигурации для настройки UI
if config.Branding != nil {
    logoURL := config.Branding["logo"]
    appName := config.Branding["name"]
    fmt.Printf("Logo: %s, Name: %s\n", logoURL, appName)
}

// Применение цветовой схемы
primaryColor := config.Colors["primary"]
secondaryColor := config.Colors["secondary"]
// ... применение в UI
```

## Обработка ошибок

### Базовая обработка

```go
result, err := client.ExecuteTemplate(req)
if err != nil {
    log.Printf("Ошибка: %v", err)
    return
}
```

### Детальная обработка

```go
result, err := client.ExecuteTemplate(req)
if err != nil {
    if errDetail, ok := err.(*types.ErrorDetail); ok {
        fmt.Printf("Code: %s\n", errDetail.Code)
        fmt.Printf("Type: %s\n", errDetail.Type)
        fmt.Printf("Message: %s\n", errDetail.Message)
        
        if errDetail.Field != "" {
            fmt.Printf("Field: %s\n", errDetail.Field)
        }
        
        if errDetail.Details != "" {
            fmt.Printf("Details: %s\n", errDetail.Details)
        }
    } else {
        log.Printf("Неожиданная ошибка: %v", err)
    }
    return
}
```

### Проверка типа ошибки

```go
if errDetail, ok := err.(*types.ErrorDetail); ok {
    switch {
    case errDetail.IsValidationError():
        fmt.Println("Ошибка валидации")
    case errDetail.IsAuthenticationError():
        fmt.Println("Ошибка аутентификации - проверьте токен")
    case errDetail.IsAuthorizationError():
        fmt.Println("Ошибка авторизации - недостаточно прав")
    case errDetail.IsRateLimitError():
        fmt.Println("Превышен лимит запросов")
    case errDetail.IsInternalError():
        fmt.Println("Внутренняя ошибка сервера")
    }
}
```

## Работа с результатами

### Обработка секций по доменам

```go
result, err := client.ExecuteTemplate(req)
if err != nil {
    log.Fatal(err)
}

for _, section := range result.Sections {
    fmt.Printf("Domain: %s\n", section.DomainID)
    fmt.Printf("Status: %s\n", section.Status)
    fmt.Printf("Results: %d\n", len(section.Results))
    
    for _, item := range section.Results {
        fmt.Printf("  - %s (relevance: %.2f)\n", 
            item.Title, item.Relevance)
    }
}
```

### Обработка веб-поиска

```go
if result.WebSearch != nil {
    fmt.Printf("Search Engine: %s\n", result.WebSearch.SearchEngine)
    fmt.Printf("Total Results: %d\n", result.WebSearch.TotalResults)
    
    for _, searchResult := range result.WebSearch.Results {
        fmt.Printf("  - %s\n", searchResult.Title)
        fmt.Printf("    URL: %s\n", searchResult.URL)
    }
}
```

### Обработка ранжирования

```go
if result.Ranking != nil {
    fmt.Printf("Algorithm: %s\n", result.Ranking.Algorithm)
    
    for _, item := range result.Ranking.Items {
        fmt.Printf("  Rank %d: %s (score: %.2f)\n", 
            item.Rank, item.ID, item.Score)
    }
}
```

## Enterprise возможности (v2.0.0) ✨

### Настройка enterprise параметров

```go
// Настройка приоритетов и кэширования
client.SetPriority("high")                    // low, normal, high, critical
client.SetCacheControl("cache-first")          // no-cache, cache-only, cache-first, network-first
client.SetCacheTTL(300)                       // TTL в секундах
client.SetRequestSource("batch")              // user, system, batch, webhook
client.SetExperiment("enterprise-rollout")     // A/B тестирование
client.SetFeatureFlag("advanced_analytics", "enabled")
```

### Расширенные фильтры поиска

```go
req := &types.ExecuteTemplateRequest{
    Query: "купить смартфон с хорошей камерой",
    Filters: &types.AdvancedFilters{
        Domains:        []string{"commerce", "reviews"},
        ExcludeDomains: []string{"adult"},
        MinRelevance:   0.8,
        MaxResults:     50,
        SortBy:         "relevance", // relevance, date, price, rating
        DateRange: &types.DateRange{
            From: time.Now().AddDate(0, 0, -30).Unix(), // Последние 30 дней
            To:   time.Now().Unix(),
        },
    },
}
```

### Локализация и региональные настройки

```go
req := &types.ExecuteTemplateRequest{
    Query: "купить ноутбук",
    Context: &types.UserContext{
        UserID:    "user-123",
        TenantID:  "enterprise-company-abc",
        Locale:    "ru-RU",              // ru-RU, en-US, etc.
        Timezone:  "Europe/Moscow",       // IANA timezone
        Currency:  "RUB",                 // RUB, USD, EUR
        Region:    "RU",                  // RU, US, EU
    },
}
```

### Enterprise метрики в ответах

```go
result, err := client.ExecuteTemplate(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Проверка enterprise метрик
if result.ResponseMetadata != nil {
    // Rate limiting
    if result.ResponseMetadata.RateLimitInfo != nil {
        fmt.Printf("Rate limit: %d/%d (reset: %d)\n",
            result.ResponseMetadata.RateLimitInfo.Remaining,
            result.ResponseMetadata.RateLimitInfo.Limit,
            result.ResponseMetadata.RateLimitInfo.ResetAt)
    }
    
    // Кэширование
    if result.ResponseMetadata.CacheInfo != nil {
        fmt.Printf("Cache: %s (TTL: %ds)\n",
            map[bool]string{true: "hit", false: "miss"}[result.ResponseMetadata.CacheInfo.CacheHit],
            result.ResponseMetadata.CacheInfo.CacheTTL)
    }
    
    // Квоты
    if result.ResponseMetadata.QuotaInfo != nil {
        fmt.Printf("Quota: %d/%d (%s)\n",
            result.ResponseMetadata.QuotaInfo.QuotaUsed,
            result.ResponseMetadata.QuotaInfo.QuotaLimit,
            result.ResponseMetadata.QuotaInfo.QuotaType)
    }
}

// Пагинация
if result.Pagination != nil {
    fmt.Printf("Page %d/%d (%d items)\n",
        result.Pagination.Page,
        result.Pagination.TotalPages,
        result.Pagination.TotalItems)
    
    if result.Pagination.HasNext {
        // Загрузить следующую страницу используя next_cursor
        nextReq := &types.ExecuteTemplateRequest{
            Query: req.Query,
            Filters: &types.AdvancedFilters{
                // Используйте next_cursor для следующей страницы
            },
        }
    }
}
```

### Batch операции

```go
// Создание batch запроса
batch := client.NewBatchBuilder().
    AddOperation("execute_template", &types.ExecuteTemplateRequest{
        Query: "купить iPhone 15",
        Context: &types.UserContext{TenantID: "enterprise-company-abc"},
    }).
    AddOperation("execute_template", &types.ExecuteTemplateRequest{
        Query: "забронировать отель в Париже",
        Context: &types.UserContext{TenantID: "enterprise-company-abc"},
    }).
    AddOperation("log_event", &types.LogEventRequest{
        EventType: "batch_operation",
        TenantID:  "enterprise-company-abc",
        Data:      map[string]interface{}{"batch_size": 2},
    }).
    SetOptions(&types.BatchOptions{
        Parallel:      true,  // Параллельное выполнение
        StopOnError:   false, // Продолжать при ошибках
        MaxConcurrency: 10,   // Максимальная параллельность
    })

// Выполнение batch
batchResult, err := batch.Execute(ctx, client)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Batch: %d/%d successful, %d failed\n",
    batchResult.Successful, batchResult.Total, batchResult.Failed)

// Обработка результатов
for _, res := range batchResult.Results {
    if res.Success {
        fmt.Printf("Operation %d: ✅ %d ms\n", res.OperationID, res.ExecutionTimeMS)
    } else {
        fmt.Printf("Operation %d: ❌ %s\n", res.OperationID, res.Error.Message)
    }
}
```

### Webhooks для асинхронных операций

```go
// Регистрация webhook
webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
    Config: &types.WebhookConfig{
        URL:    "https://my-app.company.com/webhooks/nexus",
        Events: []string{"template.completed", "template.failed", "batch.completed"},
        Secret: "webhook-secret-123",
        RetryPolicy: &types.WebhookRetryPolicy{
            MaxRetries:    3,
            InitialDelay:  1000,  // 1 секунда
            MaxDelay:      30000, // 30 секунд
            BackoffFactor: 2.0,
        },
        Active:      true,
        Description: "Enterprise webhook for async operations",
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Webhook registered: %s\n", webhookResp.WebhookID)

// Получение списка webhooks
webhooks, err := client.ListWebhooks(ctx, &types.ListWebhooksRequest{
    ActiveOnly: true,
    Limit:      10,
    Offset:     0,
})
if err != nil {
    log.Fatal(err)
}

for _, wh := range webhooks.Webhooks {
    fmt.Printf("Webhook %s: %s (%d/%d успехов/ошибок)\n",
        wh.ID, wh.Config.URL, wh.SuccessCount, wh.ErrorCount)
}

// Тестирование webhook
testResp, err := client.TestWebhook(ctx, &types.TestWebhookRequest{
    WebhookID: webhookResp.WebhookID,
    Event:     "template.completed",
    Data:      map[string]interface{}{"test": true},
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Webhook test: %s (%d ms, code %d)\n",
    testResp.Status, testResp.ResponseTimeMS, testResp.ResponseCode)

// Удаление webhook
deleteResp, err := client.DeleteWebhook(ctx, webhookResp.WebhookID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Webhook deleted: %s\n", deleteResp.WebhookID)
```

### Расширенная аналитика

```go
// Получение enterprise аналитики
stats, err := client.GetStats(ctx, &types.GetStatsRequest{
    TenantID: "enterprise-company-abc",
    Days:     30,
})
if err != nil {
    log.Fatal(err)
}

// Метрики конверсии
if stats.ConversionMetrics != nil {
    fmt.Printf("Search → Result: %.1f%%\n", stats.ConversionMetrics.SearchToResult*100)
    fmt.Printf("Result → Action: %.1f%%\n", stats.ConversionMetrics.ResultToAction*100)
    fmt.Printf("Template Success: %.1f%%\n", stats.ConversionMetrics.TemplateSuccess*100)
    fmt.Printf("User Retention: %.1f%%\n", stats.ConversionMetrics.UserRetention*100)
}

// Метрики производительности
if stats.PerformanceMetrics != nil {
    fmt.Printf("Avg Response Time: %.0f ms\n", stats.PerformanceMetrics.AvgResponseTimeMS)
    fmt.Printf("P95 Response Time: %.0f ms\n", stats.PerformanceMetrics.P95ResponseTimeMS)
    fmt.Printf("P99 Response Time: %.0f ms\n", stats.PerformanceMetrics.P99ResponseTimeMS)
    fmt.Printf("Error Rate: %.2f%%\n", stats.PerformanceMetrics.ErrorRate*100)
    fmt.Printf("Throughput: %d req/min\n", stats.PerformanceMetrics.ThroughputRPM)
}

// Разбивка по доменам
if stats.DomainBreakdown != nil {
    for domain, metrics := range stats.DomainBreakdown {
        fmt.Printf("%s: %d requests, %.1f%% success, %.0f ms avg\n",
            domain, metrics.RequestsCount, metrics.SuccessRate*100, metrics.AvgResponseTimeMS)
    }
}
```

### Детальный health check

```go
// Базовый health check
health, err := client.Health(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Health: %s (version: %s)\n", health.Status, health.Version)

// Enterprise readiness check
ready, err := client.Ready(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Readiness: %s\n", ready.Status)
fmt.Printf("Database: %s\n", ready.Checks.Database)
fmt.Printf("Redis: %s\n", ready.Checks.Redis)
fmt.Printf("AI Services: %s\n", ready.Checks.AIServices)

// Детальный статус компонентов
if ready.Components != nil {
    for name, component := range ready.Components {
        status := "✅"
        if component.Status != "healthy" {
            status = "⚠️"
        }
        fmt.Printf("%s %s: %s", status, name, component.Status)
        if component.LatencyMS > 0 {
            fmt.Printf(" (%d ms)", component.LatencyMS)
        }
        fmt.Println()
    }
}

// Информация о емкости
if ready.Capacity != nil {
    fmt.Printf("Current Load: %.1f%%\n", ready.Capacity.CurrentLoad*100)
    fmt.Printf("Max Capacity: %d req/sec\n", ready.Capacity.MaxCapacity)
    fmt.Printf("Available Capacity: %d req/sec\n", ready.Capacity.AvailableCapacity)
    fmt.Printf("Active Connections: %d\n", ready.Capacity.ActiveConnections)
}
```

## Admin API (v2.0.0) 🔧

Admin API предоставляет полный контроль над конфигурацией системы для администраторов.
Требует соответствующих прав доступа (superuser/admin роли).

### Получение Admin клиента

```go
// Получаем admin клиент
admin := client.Admin()
```

### Управление AI конфигурацией

```go
// Получить текущую конфигурацию AI
aiConfig, err := admin.GetAIConfig(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("AI Provider: %s, Model: %s\n", aiConfig.Provider, aiConfig.Model)

// Обновить конфигурацию AI
newConfig := &types.AIConfig{
    Provider:   "openai",
    Model:      "gpt-4-turbo",
    APIKey:     "encrypted_key",
    MaxTokens:  4000,
    Temperature: 0.7,
    TopP:       1.0,
    Timeout:    60,
    Enabled:    true,
}

err = admin.UpdateAIConfig(ctx, newConfig)
if err != nil {
    log.Fatal(err)
}
```

### Управление промптами

```go
// Получить список промптов для домена commerce
prompts, err := admin.ListPrompts(ctx, "commerce")
if err != nil {
    log.Fatal(err)
}

for _, prompt := range prompts {
    fmt.Printf("Prompt: %s (%s)\n", prompt.Name, prompt.Type)
}

// Создать новый промпт
newPrompt := &types.PromptConfig{
    Name:        "Commerce Search v2",
    Description: "Улучшенный промпт для поиска товаров",
    Domain:      "commerce",
    Type:        "system",
    Template:    "Ты помощник для поиска товаров. Запрос: {{query}}",
    Variables:   []string{"query"},
    Version:     1,
    Active:      true,
}

createdPrompt, err := admin.CreatePrompt(ctx, newPrompt)
if err != nil {
    log.Fatal(err)
}
```

### Управление доменами

```go
// Получить список всех доменов
domains, err := admin.ListDomains(ctx)
if err != nil {
    log.Fatal(err)
}

for _, domain := range domains {
    fmt.Printf("Domain: %s (%s) - %s\n", domain.Name, domain.Type, domain.Endpoint)
}

// Обновить ключевые слова домена
keywords := []string{"купить", "заказать", "товар", "цена", "доставка", "оплата"}
err = admin.UpdateDomainKeywords(ctx, "commerce", keywords)
if err != nil {
    log.Fatal(err)
}

// Обновить правила качества домена
qualityRules := []types.QualityRule{
    {
        Metric:      "relevance",
        Condition:   "min_relevance",
        Threshold:   0.7,
        Weight:      0.3,
        Description: "Релевантность должна быть выше 0.7",
    },
    {
        Metric:      "completeness",
        Condition:   "has_price",
        Threshold:   1.0,
        Weight:      0.25,
        Description: "Должен содержать информацию о цене",
    },
}

err = admin.UpdateDomainQualityRules(ctx, "commerce", qualityRules)
if err != nil {
    log.Fatal(err)
}
```

### Управление интеграциями

```go
// Получить список платежных интеграций
integrations, err := admin.ListIntegrations(ctx, "payment")
if err != nil {
    log.Fatal(err)
}

for _, integration := range integrations {
    fmt.Printf("Integration: %s (%s) - %s\n", integration.Name, integration.Provider, integration.Type)
}

// Создать новую интеграцию
newIntegration := &types.IntegrationConfig{
    Name:        "Stripe Payment",
    Type:        "payment",
    Provider:    "stripe",
    Enabled:     true,
    Config:      map[string]interface{}{"currency": "RUB"},
    Credentials: map[string]string{"api_key": "encrypted_key"},
    WebhookURL:  "https://api.nexus.dev/webhooks/stripe",
}

createdIntegration, err := admin.CreateIntegration(ctx, newIntegration)
if err != nil {
    log.Fatal(err)
}
```

### Управление frontend конфигурациями

```go
// Получить активную конфигурацию
activeConfig, err := admin.GetActiveFrontendConfig(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Active theme: %s\n", activeConfig.Theme)

// Создать новую конфигурацию
newConfig := &types.FrontendConfig{
    Name:   "Dark Theme v2",
    Theme:  "dark",
    Colors: map[string]string{
        "primary":   "#6200ea",
        "secondary": "#03dac6",
        "accent":    "#ff4081",
        "background": "#121212",
        "text":      "#ffffff",
    },
    Active: true,
}

createdConfig, err := admin.CreateFrontendConfig(ctx, newConfig)
if err != nil {
    log.Fatal(err)
}

// Установить как активную
err = admin.SetActiveFrontendConfig(ctx, createdConfig.ID)
if err != nil {
    log.Fatal(err)
}
```

### Инициализация доменов по умолчанию

```go
// Создать стандартные домены с настройками по умолчанию
err = admin.InitializeDefaultDomains(ctx)
if err != nil {
    log.Printf("Failed to initialize domains: %v", err)
} else {
    fmt.Println("Default domains initialized successfully")
}
```

## Проверка здоровья сервера

```go
health, err := client.Health(ctx)
if err != nil {
    log.Printf("Сервер недоступен: %v", err)
} else {
    fmt.Printf("Сервер доступен: %s (version: %s)\n", health.Status, health.Version)
}
```

## Переменные окружения

```go
import "os"

baseURL := os.Getenv("NEXUS_BASE_URL")
if baseURL == "" {
    baseURL = "http://localhost:8080"
}

token := os.Getenv("NEXUS_TOKEN")

client := nexus.NewClient(nexus.Config{
    BaseURL: baseURL,
    Token:   token,
})
```

## Лучшие практики

### 1. Всегда обрабатывайте ошибки

```go
result, err := client.ExecuteTemplate(req)
if err != nil {
    // Всегда обрабатывайте ошибки
    return fmt.Errorf("failed to execute template: %w", err)
}
```

### 2. Используйте контекст для таймаутов

```go
// Используйте Timeout в конфигурации клиента
cfg := nexus.Config{
    BaseURL: "https://api.nexus.dev",
    Timeout: 10 * time.Second,
}
```

### 3. Проверяйте метаданные ответа

```go
if result.ResponseMetadata != nil {
    fmt.Printf("Server version: %s\n", result.ResponseMetadata.ServerVersion)
    fmt.Printf("Processing time: %d ms\n", result.ResponseMetadata.ProcessingTimeMS)
}
```

### 4. Используйте правильные версии протокола

```go
cfg := nexus.Config{
    ProtocolVersion: "2.0.0", // Nexus Protocol v2.0.0 - указывайте версию явно
    ClientVersion:   "2.0.0",
}
```

## Примеры

Полные примеры находятся в директории `examples/`:

- `examples/basic/main.go` - базовое использование
- `examples/error_handling/main.go` - обработка ошибок

Запуск примеров:

```bash
# Базовый пример
make run-basic

# Пример обработки ошибок
make run-error

# Или напрямую
go run ./examples/basic
go run ./examples/error_handling
```

---


## Продвинутое использование SDK

Полное руководство по использованию advanced возможностей Nexus Protocol SDK v2.0.0.

## 🎯 Advanced возможности

### ✨ Расширенные возможности в v2.0.0

1. **Advanced метрики** - Rate limiting, кэширование, квоты
2. **Batch операции** - Параллельное выполнение множественных запросов
3. **Webhooks** - Асинхронная обработка результатов
4. **Расширенная аналитика** - Метрики конверсии и производительности
5. **Детальный health check** - Статус компонентов и емкость системы
6. **Расширенные фильтры** - Продвинутый поиск с фильтрами
7. **Пагинация** - Поддержка больших результатов
8. **Локализация** - Поддержка locale, timezone, currency

## 📚 Документация

- [USAGE.md](./USAGE.md) - Полное руководство по использованию SDK
- [README.md](./README.md) - Обзор SDK и быстрый старт
- [Advanced Examples](./examples/advanced/) - Примеры использования

## 🚀 Быстрый старт

### Установка

```bash
go get github.com/pro-deploy/nexus-protocol/sdk/go
```

### Базовый advanced клиент

```go
import (
    "context"
    "time"
    
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

cfg := client.Config{
    BaseURL:         "https://api.company.com",
    Token:           "advanced-jwt-token",
    ProtocolVersion: "2.0.0", // Nexus Protocol v2.0.0 с расширенными возможностями
    ClientVersion:   "2.0.0",
    ClientID:        "advanced-app",
    ClientType:      "api",
    RetryConfig: &client.RetryConfig{
        MaxRetries: 5,
        InitialDelay: 200 * time.Millisecond,
        MaxDelay: 10 * time.Second,
    },
}

client := client.NewClient(cfg)
ctx := context.Background()
```

### Настройка advanced параметров

```go
// Приоритеты и кэширование
client.SetPriority("high")
client.SetCacheControl("cache-first")
client.SetCacheTTL(300)

// A/B тестирование
client.SetExperiment("advanced-rollout")
client.SetFeatureFlag("advanced_analytics", "enabled")
```

## 📖 Примеры использования

### 1. Расширенный поиск с фильтрами

```go
req := &types.ExecuteTemplateRequest{
    Query: "купить смартфон с хорошей камерой",
    Context: &types.UserContext{
        UserID:   "user-123",
        TenantID: "advanced-company-abc",
        Locale:   "ru-RU",
        Currency: "RUB",
        Region:   "RU",
    },
    Filters: &types.AdvancedFilters{
        Domains:      []string{"commerce", "reviews"},
        MinRelevance: 0.8,
        MaxResults:  50,
        SortBy:      "relevance",
    },
}

result, err := client.ExecuteTemplate(ctx, req)
```

### 2. Batch операции

```go
batch := client.NewBatchBuilder().
    AddOperation("execute_template", templateReq1).
    AddOperation("execute_template", templateReq2).
    SetOptions(&types.BatchOptions{
        Parallel: true,
    })

result, err := batch.Execute(ctx, client)
```

### 3. Webhooks

```go
webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
    Config: &types.WebhookConfig{
        URL:    "https://app.company.com/webhooks",
        Events: []string{"template.completed"},
        Secret: "webhook-secret",
    },
})
```

## 🏗️ Deployment

Готовые конфигурации для развертывания:

- **Docker Compose**: `../../deployment/docker-compose.yml`
- **Kubernetes**: `../../deployment/kubernetes/`
- **Deployment Guide**: `../../deployment/DEPLOYMENT.md`

## 📊 Мониторинг

### Advanced метрики в ответах

```go
if result.ResponseMetadata != nil {
    // Rate limiting
    if result.ResponseMetadata.RateLimitInfo != nil {
        fmt.Printf("Rate limit: %d/%d\n",
            result.ResponseMetadata.RateLimitInfo.Remaining,
            result.ResponseMetadata.RateLimitInfo.Limit)
    }
    
    // Кэширование
    if result.ResponseMetadata.CacheInfo != nil {
        fmt.Printf("Cache: %s\n",
            map[bool]string{true: "hit", false: "miss"}[result.ResponseMetadata.CacheInfo.CacheHit])
    }
    
    // Квоты
    if result.ResponseMetadata.QuotaInfo != nil {
        fmt.Printf("Quota: %d/%d\n",
            result.ResponseMetadata.QuotaInfo.QuotaUsed,
            result.ResponseMetadata.QuotaInfo.QuotaLimit)
    }
}
```

### Health check

```go
ready, err := client.Ready(ctx)
if err != nil {
    log.Fatal(err)
}

// Детальный статус компонентов
for name, component := range ready.Components {
    fmt.Printf("%s: %s (%d ms)\n",
        name, component.Status, component.LatencyMS)
}

// Емкость системы
if ready.Capacity != nil {
    fmt.Printf("Load: %.1f%%\n", ready.Capacity.CurrentLoad*100)
}
```

## 💰 Бизнес-преимущества

### Средний бизнес (50-500 сотрудников)
- **Внедрение**: 1-3 дня вместо 2-6 месяцев
- **Конверсия**: +75% (30% → 67.5%)
- **Экономия**: $200K-500K/год

### Крупный бизнес (500+ сотрудников)
- **Multi-tenant**: полная изоляция данных
- **Advanced monitoring**: детальные health checks
- **Batch operations**: высокая производительность
- **Экономия**: $500K-2M/год

## 🔗 Полезные ссылки

- [API Reference](./README.md#api-reference)
- [Examples](./examples/)
- [Deployment Guide](../../deployment/DEPLOYMENT.md)
- [Protocol Documentation](../../protocol/)

## 📞 Поддержка

Для advanced клиентов:
- Email: support@nexus-protocol.com
- Slack: #nexus-advanced
- 24/7 техническая поддержка

---

**Nexus Protocol SDK v2.0.0** - Advanced Ready! 🚀

---


## Обработка ошибок в SDK

SDK автоматически парсит ошибки протокола и предоставляет удобные методы для их обработки.

## Базовая обработка

```go
result, err := nexusClient.ExecuteTemplate(ctx, req)
if err != nil {
    log.Printf("Ошибка: %v", err)
    return
}
```

## Детальная обработка

```go
result, err := nexusClient.ExecuteTemplate(ctx, req)
if err != nil {
    if errDetail, ok := err.(*types.ErrorDetail); ok {
        fmt.Printf("Code: %s\n", errDetail.Code)
        fmt.Printf("Type: %s\n", errDetail.Type)
        fmt.Printf("Message: %s\n", errDetail.Message)
        
        if errDetail.Field != "" {
            fmt.Printf("Field: %s\n", errDetail.Field)
        }
        
        if errDetail.Details != "" {
            fmt.Printf("Details: %s\n", errDetail.Details)
        }
    } else {
        log.Printf("Неожиданная ошибка: %v", err)
    }
    return
}
```

## Проверка типа ошибки

```go
if errDetail, ok := err.(*types.ErrorDetail); ok {
    switch {
    case errDetail.IsValidationError():
        fmt.Println("Ошибка валидации")
    case errDetail.IsAuthenticationError():
        fmt.Println("Ошибка аутентификации - проверьте токен")
    case errDetail.IsAuthorizationError():
        fmt.Println("Ошибка авторизации - недостаточно прав")
    case errDetail.IsRateLimitError():
        fmt.Println("Превышен лимит запросов")
    case errDetail.IsInternalError():
        fmt.Println("Внутренняя ошибка сервера")
    }
}
```

## Типы ошибок

### ValidationError

Ошибка валидации входных данных:

```go
if errDetail.IsValidationError() {
    // Обработка ошибки валидации
}
```

### AuthenticationError

Ошибка аутентификации (неверный токен):

```go
if errDetail.IsAuthenticationError() {
    // Обновить токен или запросить новый
}
```

### AuthorizationError

Ошибка авторизации (недостаточно прав):

```go
if errDetail.IsAuthorizationError() {
    // Проверить права доступа
}
```

### RateLimitError

Превышен лимит запросов:

```go
if errDetail.IsRateLimitError() {
    // Подождать и повторить запрос
}
```

### InternalError

Внутренняя ошибка сервера:

```go
if errDetail.IsInternalError() {
    // Логировать и повторить позже
}
```

## Retry при ошибках

SDK автоматически повторяет запросы при сетевых ошибках, если настроен RetryConfig:

```go
retryCfg := client.RetryConfig{
    MaxRetries:        3,
    InitialDelay:      100 * time.Millisecond,
    MaxDelay:          5 * time.Second,
    BackoffMultiplier: 2.0,
}

cfg := client.Config{
    BaseURL:     "https://api.nexus.dev",
    RetryConfig: &retryCfg,
}
```

---


## Примеры использования SDK

### Базовые примеры

### Пример 1: Простой запрос

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
    cfg := client.Config{
        BaseURL: "https://api.nexus.dev",
        Token:   "your-jwt-token",
    }
    
    client := client.NewClient(cfg)
    ctx := context.Background()
    
    req := &types.ExecuteTemplateRequest{
        Query:    "хочу борщ",
        Language: "ru",
    }
    
    result, err := client.ExecuteTemplate(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Execution ID: %s\n", result.ExecutionID)
}
```

### Пример 2: Запрос с контекстом

```go
req := &types.ExecuteTemplateRequest{
    Query:    "Найди где рядом продается кокакола",
    Language: "ru",
    Context: &types.UserContext{
        UserID: "user-123",
        Location: &types.UserLocation{
            Latitude:  55.7558,
            Longitude: 37.6173,
            Accuracy:  50,
        },
        Locale:   "ru-RU",
        Currency: "RUB",
        Region:   "RU",
    },
}

result, err := client.ExecuteTemplate(ctx, req)
```

## Enterprise примеры

### Batch операции

```go
batch := client.NewBatchBuilder().
    AddOperation("execute_template", &types.ExecuteTemplateRequest{
        Query: "купить iPhone",
        Context: &types.UserContext{UserID: "user-1"},
    }).
    AddOperation("execute_template", &types.ExecuteTemplateRequest{
        Query: "забронировать отель",
        Context: &types.UserContext{UserID: "user-1"},
    }).
    SetOptions(&types.BatchOptions{
        Parallel: true,
    })

batchResult, err := batch.Execute(ctx, client)
```

### Webhooks

```go
webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
    Config: &types.WebhookConfig{
        URL:    "https://myapp.com/webhook",
        Events: []string{"template.completed", "template.failed"},
        Secret: "webhook-secret",
    },
})
```

## Полные примеры

Полные примеры находятся в директории `examples/` SDK:

- `examples/basic/main.go` - базовое использование
- `examples/error_handling/main.go` - обработка ошибок
- `examples/iam/main.go` - аутентификация
- `examples/conversations/main.go` - беседы с AI
- `examples/analytics/main.go` - аналитика
- `examples/advanced/` - enterprise примеры

Запуск примеров:

```bash
cd sdk/go
make run-basic
make run-error
make run-iam
```

---


## Batch операции

Batch операции позволяют выполнять множественные запросы параллельно для повышения производительности.

## Создание batch запроса

```go
import (
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

batch := client.NewBatchBuilder().
    AddOperation("execute_template", &types.ExecuteTemplateRequest{
        Query: "купить iPhone 15",
        Context: &types.UserContext{UserID: "user-1"},
    }).
    AddOperation("execute_template", &types.ExecuteTemplateRequest{
        Query: "забронировать отель в Париже",
        Context: &types.UserContext{UserID: "user-1"},
    }).
    SetOptions(&types.BatchOptions{
        Parallel:      true,  // Параллельное выполнение
        StopOnError:   false, // Продолжать при ошибках
        MaxConcurrency: 10,   // Максимальная параллельность
    })
```

## Выполнение batch

```go
batchResult, err := batch.Execute(ctx, client)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Batch: %d/%d successful, %d failed\n",
    batchResult.Successful, batchResult.Total, batchResult.Failed)
```

## Обработка результатов

```go
for _, res := range batchResult.Results {
    if res.Success {
        fmt.Printf("Operation %d: ✅ %d ms\n", 
            res.OperationID, res.ExecutionTimeMS)
    } else {
        fmt.Printf("Operation %d: ❌ %s\n", 
            res.OperationID, res.Error.Message)
    }
}
```

## Опции batch

### Parallel

Параллельное выполнение операций:

```go
SetOptions(&types.BatchOptions{
    Parallel: true,
})
```

### StopOnError

Остановка при первой ошибке:

```go
SetOptions(&types.BatchOptions{
    StopOnError: true,
})
```

### MaxConcurrency

Максимальное количество параллельных операций:

```go
SetOptions(&types.BatchOptions{
    MaxConcurrency: 10,
})
```

## Преимущества

- ✅ Высокая производительность
- ✅ Параллельное выполнение
- ✅ Единая обработка ошибок
- ✅ Оптимизация использования ресурсов

---


## Webhooks

Webhooks позволяют получать асинхронные уведомления о событиях.

## Регистрация webhook

```go
webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
    Config: &types.WebhookConfig{
        URL:    "https://myapp.com/webhook",
        Events: []string{"template.completed", "template.failed", "batch.completed"},
        Secret: "webhook-secret-123",
        RetryPolicy: &types.WebhookRetryPolicy{
            MaxRetries:    3,
            InitialDelay:  1000,  // 1 секунда
            MaxDelay:      30000, // 30 секунд
            BackoffFactor: 2.0,
        },
        Active:      true,
        Description: "Webhook for async operations",
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Webhook registered: %s\n", webhookResp.WebhookID)
```

## Получение списка webhooks

```go
webhooks, err := client.ListWebhooks(ctx, &types.ListWebhooksRequest{
    ActiveOnly: true,
    Limit:      10,
    Offset:     0,
})
if err != nil {
    log.Fatal(err)
}

for _, wh := range webhooks.Webhooks {
    fmt.Printf("Webhook %s: %s (%d/%d успехов/ошибок)\n",
        wh.ID, wh.Config.URL, wh.SuccessCount, wh.ErrorCount)
}
```

## Тестирование webhook

```go
testResp, err := client.TestWebhook(ctx, &types.TestWebhookRequest{
    WebhookID: webhookResp.WebhookID,
    Event:     "template.completed",
    Data:      map[string]interface{}{"test": true},
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Webhook test: %s (%d ms, code %d)\n",
    testResp.Status, testResp.ResponseTimeMS, testResp.ResponseCode)
```

## Удаление webhook

```go
deleteResp, err := client.DeleteWebhook(ctx, webhookResp.WebhookID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Webhook deleted: %s\n", deleteResp.WebhookID)
```

## Поддерживаемые события

- `template.completed` - шаблон выполнен успешно
- `template.failed` - ошибка выполнения шаблона
- `batch.completed` - batch операция завершена
- `batch.failed` - ошибка batch операции

## Retry Policy

Webhook автоматически повторяет отправку при ошибках согласно настройкам RetryPolicy:

```go
RetryPolicy: &types.WebhookRetryPolicy{
    MaxRetries:    3,      // Максимум попыток
    InitialDelay:  1000,   // Начальная задержка (мс)
    MaxDelay:      30000,  // Максимальная задержка (мс)
    BackoffFactor: 2.0,    // Множитель задержки
}
```

---


## Аналитика

SDK предоставляет расширенные возможности аналитики для отслеживания использования и производительности.

## Логирование события

```go
logResp, err := client.LogEvent(ctx, &types.LogEventRequest{
    EventType: "user_action",
    UserID:    "user-123",
    Data: map[string]interface{}{
        "action": "viewed_page",
        "page":   "/products",
    },
})
```

## Получение событий

```go
eventsResp, err := client.GetEvents(ctx, &types.GetEventsRequest{
    EventType: "user_action",
    Limit:     10,
    Offset:    0,
})
```

## Получение статистики

```go
stats, err := client.GetStats(ctx, &types.GetStatsRequest{
    UserID: "user-123",
    Days:   7,
})
```

## Метрики конверсии

```go
if stats.ConversionMetrics != nil {
    fmt.Printf("Search → Result: %.1f%%\n", 
        stats.ConversionMetrics.SearchToResult*100)
    fmt.Printf("Result → Action: %.1f%%\n", 
        stats.ConversionMetrics.ResultToAction*100)
    fmt.Printf("Template Success: %.1f%%\n", 
        stats.ConversionMetrics.TemplateSuccess*100)
    fmt.Printf("User Retention: %.1f%%\n", 
        stats.ConversionMetrics.UserRetention*100)
}
```

## Метрики производительности

```go
if stats.PerformanceMetrics != nil {
    fmt.Printf("Avg Response Time: %.0f ms\n", 
        stats.PerformanceMetrics.AvgResponseTimeMS)
    fmt.Printf("P95 Response Time: %.0f ms\n", 
        stats.PerformanceMetrics.P95ResponseTimeMS)
    fmt.Printf("P99 Response Time: %.0f ms\n", 
        stats.PerformanceMetrics.P99ResponseTimeMS)
    fmt.Printf("Error Rate: %.2f%%\n", 
        stats.PerformanceMetrics.ErrorRate*100)
    fmt.Printf("Throughput: %d req/min\n", 
        stats.PerformanceMetrics.ThroughputRPM)
}
```

## Разбивка по доменам

```go
if stats.DomainBreakdown != nil {
    for domain, metrics := range stats.DomainBreakdown {
        fmt.Printf("%s: %d requests, %.1f%% success, %.0f ms avg\n",
            domain, metrics.RequestsCount, 
            metrics.SuccessRate*100, 
            metrics.AvgResponseTimeMS)
    }
}
```

---


## Admin API

Admin API предоставляет полный контроль над конфигурацией системы для администраторов. Требует соответствующих прав доступа (superuser/admin роли).

## Получение Admin клиента

```go
admin := client.Admin()
```

## Управление AI конфигурацией

### Получение конфигурации

```go
aiConfig, err := admin.GetAIConfig(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("AI Provider: %s, Model: %s\n", aiConfig.Provider, aiConfig.Model)
```

### Обновление конфигурации

```go
newConfig := &types.AIConfig{
    Provider:   "openai",
    Model:      "gpt-4-turbo",
    APIKey:     "encrypted_key",
    MaxTokens:  4000,
    Temperature: 0.7,
    TopP:       1.0,
    Timeout:    60,
    Enabled:    true,
}

err = admin.UpdateAIConfig(ctx, newConfig)
```

## Управление промптами

### Получение списка промптов

```go
prompts, err := admin.ListPrompts(ctx, "commerce")
if err != nil {
    log.Fatal(err)
}

for _, prompt := range prompts {
    fmt.Printf("Prompt: %s (%s)\n", prompt.Name, prompt.Type)
}
```

### Создание промпта

```go
newPrompt := &types.PromptConfig{
    Name:        "Commerce Search v2",
    Description: "Улучшенный промпт для поиска товаров",
    Domain:      "commerce",
    Type:        "system",
    Template:    "Ты помощник для поиска товаров. Запрос: {{query}}",
    Variables:   []string{"query"},
    Version:     1,
    Active:      true,
}

createdPrompt, err := admin.CreatePrompt(ctx, newPrompt)
```

## Управление доменами

### Получение списка доменов

```go
domains, err := admin.ListDomains(ctx)
if err != nil {
    log.Fatal(err)
}

for _, domain := range domains {
    fmt.Printf("Domain: %s (%s) - %s\n", 
        domain.Name, domain.Type, domain.Endpoint)
}
```

### Обновление ключевых слов домена

```go
keywords := []string{"купить", "заказать", "товар", "цена", "доставка", "оплата"}
err = admin.UpdateDomainKeywords(ctx, "commerce", keywords)
```

## Управление интеграциями

### Получение списка интеграций

```go
integrations, err := admin.ListIntegrations(ctx, "payment")
if err != nil {
    log.Fatal(err)
}

for _, integration := range integrations {
    fmt.Printf("Integration: %s (%s) - %s\n", 
        integration.Name, integration.Provider, integration.Type)
}
```

### Создание интеграции

```go
newIntegration := &types.IntegrationConfig{
    Name:        "Stripe Payment",
    Type:        "payment",
    Provider:    "stripe",
    Enabled:     true,
    Config:      map[string]interface{}{"currency": "RUB"},
    Credentials: map[string]string{"api_key": "encrypted_key"},
    WebhookURL:  "https://api.nexus.dev/webhooks/stripe",
}

createdIntegration, err := admin.CreateIntegration(ctx, newIntegration)
```

## Управление frontend конфигурациями

### Получение активной конфигурации

```go
activeConfig, err := admin.GetActiveFrontendConfig(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Active theme: %s\n", activeConfig.Theme)
```

### Создание конфигурации

```go
newConfig := &types.FrontendConfig{
    Name:   "Dark Theme v2",
    Theme:  "dark",
    Colors: map[string]string{
        "primary":   "#6200ea",
        "secondary": "#03dac6",
        "accent":    "#ff4081",
    },
    Active: true,
}

createdConfig, err := admin.CreateFrontendConfig(ctx, newConfig)
```

### Установка активной конфигурации

```go
err = admin.SetActiveFrontendConfig(ctx, createdConfig.ID)
```

---


## Client API Reference

Полная документация по API клиента SDK.

## NewClient

Создает новый клиент с указанной конфигурацией.

```go
func NewClient(config Config) *Client
```

### Пример

```go
cfg := client.Config{
    BaseURL: "https://api.nexus.dev",
    Token:   "jwt-token",
}
client := client.NewClient(cfg)
```

## ExecuteTemplate

Выполняет контекстно-зависимый шаблон.

```go
func (c *Client) ExecuteTemplate(ctx context.Context, req *types.ExecuteTemplateRequest) (*types.ExecuteTemplateResponse, error)
```

### Пример

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
}
result, err := client.ExecuteTemplate(ctx, req)
```

## GetExecutionStatus

Получает статус выполнения шаблона.

```go
func (c *Client) GetExecutionStatus(ctx context.Context, executionID string) (*types.ExecuteTemplateResponse, error)
```

### Пример

```go
status, err := client.GetExecutionStatus(ctx, "execution-id")
```

## ExecuteBatch

Выполняет пакет операций для высокой производительности (Enterprise).

```go
func (c *Client) ExecuteBatch(ctx context.Context, req *types.BatchRequest) (*types.BatchResponse, error)
```

## RegisterWebhook

Регистрирует webhook для асинхронной обработки (Enterprise).

```go
func (c *Client) RegisterWebhook(ctx context.Context, req *types.RegisterWebhookRequest) (*types.RegisterWebhookResponse, error)
```

## Health

Проверяет здоровье сервера.

```go
func (c *Client) Health(ctx context.Context) (*types.HealthResponse, error)
```

## Ready

Проверяет готовность сервера с enterprise метриками (Enterprise).

```go
func (c *Client) Ready(ctx context.Context) (*types.ReadinessResponse, error)
```

## GetFrontendConfig

Получает активную конфигурацию фронтенда (публичный endpoint).

```go
func (c *Client) GetFrontendConfig(ctx context.Context) (*types.FrontendConfig, error)
```

## SetToken

Устанавливает новый JWT токен.

```go
func (c *Client) SetToken(token string)
```

## SetPriority

Устанавливает приоритет запросов (Enterprise).

```go
func (c *Client) SetPriority(priority string)
```

Приоритеты: `low`, `normal`, `high`, `critical`

## SetCacheControl

Устанавливает контроль кэширования (Enterprise).

```go
func (c *Client) SetCacheControl(cacheControl string)
```

Опции: `no-cache`, `cache-only`, `cache-first`, `network-first`

## Admin

Получает admin клиент для управления конфигурацией.

```go
func (c *Client) Admin() *AdminClient
```

---


## Types Reference

Документация по типам данных SDK.

## ExecuteTemplateRequest

Запрос на выполнение шаблона.

```go
type ExecuteTemplateRequest struct {
    Query    string
    Language string
    Context  *UserContext
    Options  *ExecuteOptions
    Filters  *AdvancedFilters
    Metadata *RequestMetadata
}
```

## ExecuteTemplateResponse

Ответ на выполнение шаблона.

```go
type ExecuteTemplateResponse struct {
    ExecutionID string
    Status      string
    QueryType   string
    Sections    []Section
    Workflow    *Workflow
    Metadata    *ResponseMetadata
}
```

## UserContext

Контекст пользователя.

```go
type UserContext struct {
    UserID    string
    SessionID string
    TenantID  string
    Location  *UserLocation
    Locale    string
    Timezone  string
    Currency  string
    Region    string
}
```

## RequestMetadata

Метаданные запроса.

```go
type RequestMetadata struct {
    RequestID      string
    ProtocolVersion string
    ClientVersion  string
    ClientID       string
    ClientType     string
    Timestamp      int64
    CustomHeaders  map[string]string
}
```

## ResponseMetadata

Метаданные ответа.

```go
type ResponseMetadata struct {
    RequestID        string
    ProtocolVersion  string
    ServerVersion    string
    Timestamp        int64
    ProcessingTimeMS int64
    RateLimitInfo    *RateLimitInfo
    CacheInfo        *CacheInfo
    QuotaInfo        *QuotaInfo
}
```

## ErrorDetail

Детальная информация об ошибке.

```go
type ErrorDetail struct {
    Code    string
    Type    string
    Message string
    Field   string
    Details string
}
```

### Методы

- `IsValidationError() bool`
- `IsAuthenticationError() bool`
- `IsAuthorizationError() bool`
- `IsRateLimitError() bool`
- `IsInternalError() bool`

## BatchRequest

Запрос на выполнение batch операций.

```go
type BatchRequest struct {
    Operations []BatchOperation
    Options     *BatchOptions
}
```

## WebhookConfig

Конфигурация webhook.

```go
type WebhookConfig struct {
    URL         string
    Events      []string
    Secret      string
    RetryPolicy *WebhookRetryPolicy
    Active      bool
    Description string
}
```

---

