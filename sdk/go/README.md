# Nexus Protocol Go SDK

Go SDK для работы с Nexus Application Protocol.

## Установка

```bash
go get github.com/nexus-protocol/go-sdk
```

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/nexus-protocol/go-sdk/client"
    "github.com/nexus-protocol/go-sdk/types"
)

func main() {
    // Создаем клиент
    cfg := client.Config{
        BaseURL:         "http://localhost:8080",
        Token:           "your-jwt-token",
        ProtocolVersion: "1.0.0",
        ClientVersion:   "1.0.0",
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

## Основные возможности

### Создание клиента

```go
import (
    "context"
    "time"
    
    "github.com/nexus-protocol/go-sdk/client"
)

// Базовый клиент
cfg := client.Config{
    BaseURL:         "https://api.nexus.dev",
    Token:           "jwt-token",
    Timeout:         30 * time.Second,
    ProtocolVersion: "1.0.0",
    ClientVersion:   "1.0.0",
    ClientID:        "my-app",
    ClientType:      "web",
}

nexusClient := client.NewClient(cfg)
ctx := context.Background()

// Клиент с retry и логированием
retryCfg := client.RetryConfig{
    MaxRetries:        3,
    InitialDelay:      100 * time.Millisecond,
    MaxDelay:          5 * time.Second,
    BackoffMultiplier: 2.0,
}

logger := client.NewSimpleLogger(client.LogLevelInfo)

cfgWithRetry := client.Config{
    BaseURL:     "https://api.nexus.dev",
    Token:       "jwt-token",
    RetryConfig: &retryCfg,
    Logger:      logger,
}

nexusClientWithRetry := client.NewClient(cfgWithRetry)
```

### Выполнение шаблона

```go
ctx := context.Background()

req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
    Context: &types.UserContext{
        UserID:    "user-123",
        SessionID: "session-456",
    },
    Options: &types.ExecuteOptions{
        TimeoutMS:           30000,
        MaxResultsPerDomain: 5,
        ParallelExecution:   true,
        IncludeWebSearch:    true,
    },
}

result, err := nexusClient.ExecuteTemplate(ctx, req)
if err != nil {
    // Обработка ошибки
    if errDetail, ok := err.(*types.ErrorDetail); ok {
        fmt.Printf("Error: %s (%s)\n", errDetail.Message, errDetail.Code)
    }
    return
}

fmt.Printf("Execution ID: %s\n", result.ExecutionID)
fmt.Printf("Status: %s\n", result.Status)
```

### Получение статуса выполнения

```go
ctx := context.Background()
status, err := nexusClient.GetExecutionStatus(ctx, "execution-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status)
```

### Conversations (Беседы с AI)

```go
// Создание беседы
conversation, err := client.CreateConversation(ctx, &types.CreateConversationRequest{
    Title:        "Обсуждение рецептов",
    BotID:        "bot-123",
    SystemPrompt: "Ты помощник по кулинарии",
})

// Отправка сообщения
messageResp, err := client.SendMessage(ctx, conversation.ID, &types.SendMessageRequest{
    Content:     "Расскажи рецепт борща",
    MessageType: "text",
})

// Получение беседы с историей
fullConversation, err := client.GetConversation(ctx, conversation.ID)
```

### Analytics (Аналитика)

```go
// Логирование события
logResp, err := client.LogEvent(ctx, &types.LogEventRequest{
    EventType: "user_action",
    UserID:    "user-123",
    Data: map[string]interface{}{
        "action": "viewed_page",
    },
})

// Получение событий
eventsResp, err := client.GetEvents(ctx, &types.GetEventsRequest{
    EventType: "user_action",
    Limit:     10,
    Offset:    0,
})

// Получение статистики
stats, err := client.GetStats(ctx, &types.GetStatsRequest{
    UserID: "user-123",
    Days:   7,
})
```

### IAM (Аутентификация и авторизация)

```go
// Регистрация
registerResp, err := client.RegisterUser(ctx, &types.RegisterUserRequest{
    Email:     "user@example.com",
    Password:  "password123",
    FirstName: "Иван",
    LastName:  "Иванов",
})

// Вход (токен устанавливается автоматически)
loginResp, err := client.Login(ctx, &types.LoginRequest{
    Email:    "user@example.com",
    Password: "password123",
})

// Получение профиля
profile, err := client.GetUserProfile(ctx)

// Обновление профиля
updatedProfile, err := client.UpdateUserProfile(ctx, &types.UpdateProfileRequest{
    FirstName: "Иван",
    LastName:  "Петров",
})

// Обновление токена
refreshResp, err := client.RefreshToken(ctx, &types.RefreshTokenRequest{
    RefreshToken: loginResp.RefreshToken,
})
```

### Обработка ошибок

SDK автоматически парсит ошибки протокола:

```go
result, err := nexusClient.ExecuteTemplate(req)
if err != nil {
    if errDetail, ok := err.(*types.ErrorDetail); ok {
        switch {
        case errDetail.IsValidationError():
            fmt.Println("Ошибка валидации:", errDetail.Message)
        case errDetail.IsAuthenticationError():
            fmt.Println("Ошибка аутентификации")
        case errDetail.IsAuthorizationError():
            fmt.Println("Ошибка авторизации")
        case errDetail.IsRateLimitError():
            fmt.Println("Превышен лимит запросов")
        }
    }
}
```

### Автоматические метаданные

SDK автоматически создает метаданные запроса:

```go
// Метаданные создаются автоматически
req := &types.ExecuteTemplateRequest{
    Query: "хочу борщ",
    // Metadata будет создан автоматически
}

// Или создайте вручную
metadata := types.NewRequestMetadata("1.0.0", "1.0.0")
metadata.ClientID = "my-app"
metadata.ClientType = "web"

req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Metadata: metadata,
}
```

## API Reference

### Client

#### `NewClient(config Config) *Client`

Создает новый клиент с указанной конфигурацией.

#### `ExecuteTemplate(req *ExecuteTemplateRequest) (*ExecuteTemplateResponse, error)`

Выполняет контекстно-зависимый шаблон.

#### `GetExecutionStatus(executionID string) (*ExecuteTemplateResponse, error)`

Получает статус выполнения шаблона.

#### `StreamTemplateResults(executionID string) (*http.Response, error)`

Получает поток результатов выполнения (Server-Sent Events).

#### `Health() error`

Проверяет здоровье сервера.

### Types

#### `RequestMetadata`

Метаданные запроса с полями:
- `RequestID` - UUID запроса
- `ProtocolVersion` - версия протокола
- `ClientVersion` - версия клиента
- `ClientID` - идентификатор клиента
- `ClientType` - тип клиента (web, mobile, sdk, api, desktop)
- `Timestamp` - временная метка
- `CustomHeaders` - кастомные заголовки

#### `ErrorDetail`

Детальная информация об ошибке с методами:
- `IsValidationError()` - проверка типа ошибки
- `IsAuthenticationError()` - проверка типа ошибки
- `IsAuthorizationError()` - проверка типа ошибки
- `IsRateLimitError()` - проверка типа ошибки
- `IsInternalError()` - проверка типа ошибки

### Retry и Rate Limiting

SDK автоматически повторяет запросы при сетевых ошибках и обрабатывает rate limiting:

```go
// Настройка retry
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

// При HTTP 429 (rate limit) SDK автоматически ждет и повторяет запрос
// Используется заголовок Retry-After или exponential backoff
```

### Логирование

```go
import "github.com/nexus-protocol/go-sdk/client"

// Создание логгера
logger := client.NewSimpleLogger(client.LogLevelDebug)

cfg := client.Config{
    BaseURL: "https://api.nexus.dev",
    Logger:  logger,
}

// Или установить позже
client.SetLogger(logger)
```

### Interceptors (Middleware)

```go
// Создание interceptor для измерения времени
type TimingInterceptor struct{}

func (t *TimingInterceptor) BeforeRequest(ctx context.Context, req *http.Request) error {
    // Логика перед запросом
    return nil
}

func (t *TimingInterceptor) AfterResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
    // Логика после ответа
    return nil
}

// Добавление interceptor
client.AddInterceptor(&TimingInterceptor{})
```

### Метрики

```go
// Создание коллектора метрик
metricsCollector := client.NewSimpleMetricsCollector()

// Создание interceptor для метрик
metricsInterceptor := client.NewMetricsInterceptor(metricsCollector)

// Добавление interceptor
client.AddInterceptor(metricsInterceptor)

// После выполнения запросов получаем статистику
stats := metricsCollector.GetStats()
fmt.Printf("Requests: %v\n", stats["requests"])
fmt.Printf("Errors: %v\n", stats["errors"])
fmt.Printf("Avg durations: %v\n", stats["avg_durations"])
```

### Валидация по JSON Schema

```go
// Создание валидатора
validator := client.NewValidator()

// Загрузка схемы
err := validator.LoadSchema("execute-template", "schemas/message-schema.json")
if err != nil {
    log.Fatal(err)
}

// Использование валидатора в клиенте
cfg := client.Config{
    BaseURL:  "https://api.nexus.dev",
    Validator: validator,
}

client := client.NewClient(cfg)

// Или установить позже
client.SetValidator(validator)
```

## Примеры

Примеры использования находятся в директории `examples/`:

- `basic/main.go` - базовое использование
- `error_handling/main.go` - обработка ошибок
- `iam/main.go` - аутентификация и управление пользователями
- `conversations/main.go` - беседы с AI
- `analytics/main.go` - аналитика и события
- `retry/main.go` - retry логика и rate limiting ✨
- `interceptors/main.go` - использование interceptors ✨
- `metrics/main.go` - сбор метрик ✨

Запуск примеров:

```bash
make run-basic         # Базовый пример
make run-error         # Обработка ошибок
make run-iam           # IAM пример
make run-conversations # Conversations пример
make run-analytics     # Analytics пример
make run-retry         # Retry пример
make run-interceptors  # Interceptors пример
make run-metrics       # Metrics пример
```

## Зависимости

- `github.com/google/uuid` - генерация UUID
- `github.com/xeipuuv/gojsonschema` - валидация JSON Schema (опционально)

## Лицензия

MIT License

