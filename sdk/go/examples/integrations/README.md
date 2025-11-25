# Пример работы с интеграциями

Этот пример показывает, как использовать **Admin API** для управления интеграциями в Nexus Protocol.

## ⚠️ Важно: SDK и внутренние протоколы

**SDK содержит только публичные API**, которые определены в OpenAPI спецификации.

- ✅ **Admin API** - публичный API для управления интеграциями
- ❌ **MCP (Model Context Protocol)** - внутренний протокол сервера, не входит в SDK

## Запуск примера

1. **Запустите Nexus сервер**:
   ```bash
   cd /path/to/server
   go run cmd/main.go
   ```

2. **Запустите пример**:
   ```bash
   cd sdk/go/examples/integrations
   go run main.go
   ```

## Что делает пример

### 1. Получение списка интеграций
```go
integrations, err := adminClient.ListIntegrations(ctx, "")
```
Показывает все зарегистрированные интеграции.

### 2. Фильтрация по типу
```go
dataSources, err := adminClient.ListIntegrations(ctx, "data_source")
```
Получает только интеграции определенного типа.

### 3. Получение деталей интеграции
```go
integration, err := adminClient.GetIntegration(ctx, integrationID)
```
Получает полную информацию об интеграции.

### 4. Создание интеграции
```go
created, err := adminClient.CreateIntegration(ctx, config)
```
Создает новую интеграцию с указанной конфигурацией.

### 5. Обновление интеграции
```go
updated, err := adminClient.UpdateIntegration(ctx, id, config)
```
Обновляет существующую интеграцию.

### 6. Использование через домен
```go
// Запрос через домен 'integrations'
templateReq := &types.ExecuteTemplateRequest{
    Query: "получить данные из weather-api",
    Language: "ru",
}
```
Доступ к данным интеграций происходит через домен `integrations` в обычном API.

### 7. Удаление интеграции
```go
err := adminClient.DeleteIntegration(ctx, id)
```
Удаляет интеграцию.

## Архитектура

```
┌─────────────┐
│   SDK       │
│  (Client)   │
└──────┬──────┘
       │ Admin API (публичный)
       ▼
┌─────────────┐
│   Server    │
│  (Admin)    │
└──────┬──────┘
       │
       ▼
┌─────────────┐     ┌─────────────┐
│  MCP Server │────▶│  Sources    │
│ (внутренний)│     │  (GitHub,   │
│             │     │   Weather)  │
└─────────────┘     └─────────────┘
```

**SDK** работает только с **публичным Admin API**.  
**MCP** - это внутренний протокол, который используется сервером для работы с источниками данных.

## Использование интеграций

### Управление (через Admin API)
```go
adminClient := nexusClient.Admin()

// Создать интеграцию
config := &types.IntegrationConfig{
    ID:       "github-api",
    Name:     "GitHub API",
    Type:     "data_source",
    Provider: "github",
    Enabled:  true,
    Config: map[string]interface{}{
        "base_url": "https://api.github.com",
    },
    Credentials: map[string]string{
        "api_key": "your_token",
    },
}

created, err := adminClient.CreateIntegration(ctx, config)
```

### Доступ к данным (через домен)
```go
// Запрос через обычный API
result, err := nexusClient.ExecuteTemplate(ctx, &types.ExecuteTemplateRequest{
    Query:    "получить репозитории пользователя octocat",
    Language: "ru",
})

// Система автоматически:
// 1. Определит домен 'integrations'
// 2. Выберет нужный источник (github-api)
// 3. Получит данные через MCP
// 4. Вернет результат
```

## Типы интеграций

- `payment` - платежные системы
- `delivery` - доставка
- `notifications` - уведомления
- `analytics` - аналитика
- `data_source` - источники данных (GitHub, Weather, etc.)
- `custom` - кастомные интеграции

## Ошибки и отладка

### Возможные ошибки

1. **401 Unauthorized**
   ```
   Failed to list integrations: unauthorized
   ```
   **Решение**: Проверьте JWT токен в конфигурации клиента

2. **403 Forbidden**
   ```
   Failed to create integration: forbidden
   ```
   **Решение**: Убедитесь, что у пользователя есть права администратора

3. **404 Not Found**
   ```
   Failed to get integration: not found
   ```
   **Решение**: Проверьте ID интеграции

## Связанные файлы

- `sdk/go/client/admin.go` - Admin API клиент
- `sdk/go/types/templates.go` - типы данных (IntegrationConfig)
- `server/go/internal/admin/service.go` - реализация Admin API
- `server/go/internal/integrations/` - MCP сервер (внутренний)

## Принципы SDK

✅ **Включаем в SDK:**
- Публичные API эндпоинты
- Типы данных для публичных API
- Утилиты для работы с протоколом

❌ **НЕ включаем в SDK:**
- Клиенты для внутренних протоколов (MCP)
- Знание о внутренней архитектуре
- Прямой доступ к внутренним сервисам