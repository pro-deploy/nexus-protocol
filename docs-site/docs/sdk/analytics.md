---
id: analytics
title: Аналитика
sidebar_label: Аналитика
---

# Аналитика

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

