---
id: basic-usage
title: Базовое использование
sidebar_label: Базовое использование
---

# Базовое использование SDK

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

