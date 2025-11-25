# Руководство по использованию Go SDK

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

