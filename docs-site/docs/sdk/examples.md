---
id: examples
title: Примеры
sidebar_label: Примеры
---

# Примеры использования SDK

## Базовые примеры

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

