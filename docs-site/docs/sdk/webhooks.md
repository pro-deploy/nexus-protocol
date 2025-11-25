---
id: webhooks
title: Webhooks
sidebar_label: Webhooks
---

# Webhooks

Webhooks позволяют получать асинхронные уведомления о событиях.

## Регистрация webhook

```go
webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
    Config: &types.WebhookConfig{
        URL:    "https://myapp.com/webhook",
        Events: []string{"template.completed", "template.failed", "batch.completed"},
        Secret: "webhook-secret-123",
        RetryPolicy: &types.WebhookRetryPolicy{
            MaxRetries:    3,
            InitialDelay:  1000,  // 1 секунда
            MaxDelay:      30000, // 30 секунд
            BackoffFactor: 2.0,
        },
        Active:      true,
        Description: "Webhook for async operations",
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Webhook registered: %s\n", webhookResp.WebhookID)
```

## Получение списка webhooks

```go
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
```

## Тестирование webhook

```go
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
```

## Удаление webhook

```go
deleteResp, err := client.DeleteWebhook(ctx, webhookResp.WebhookID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Webhook deleted: %s\n", deleteResp.WebhookID)
```

## Поддерживаемые события

- `template.completed` - шаблон выполнен успешно
- `template.failed` - ошибка выполнения шаблона
- `batch.completed` - batch операция завершена
- `batch.failed` - ошибка batch операции

## Retry Policy

Webhook автоматически повторяет отправку при ошибках согласно настройкам RetryPolicy:

```go
RetryPolicy: &types.WebhookRetryPolicy{
    MaxRetries:    3,      // Максимум попыток
    InitialDelay:  1000,   // Начальная задержка (мс)
    MaxDelay:      30000,  // Максимальная задержка (мс)
    BackoffFactor: 2.0,    // Множитель задержки
}
```

