# Nexus Protocol Go SDK - Advanced Guide

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
