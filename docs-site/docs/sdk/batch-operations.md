---
id: batch-operations
title: Batch операции
sidebar_label: Batch операции
---

# Batch операции

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

