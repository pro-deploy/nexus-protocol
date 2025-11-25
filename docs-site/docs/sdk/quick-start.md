---
id: quick-start
title: Быстрый старт
sidebar_label: Быстрый старт
---

# Быстрый старт

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

