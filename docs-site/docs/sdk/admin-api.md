---
id: admin-api
title: Admin API
sidebar_label: Admin API
---

# Admin API

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

