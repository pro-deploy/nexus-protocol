---
id: error-handling
title: Обработка ошибок
sidebar_label: Обработка ошибок
---

# Обработка ошибок

SDK автоматически парсит ошибки протокола и предоставляет удобные методы для их обработки.

## Базовая обработка

```go
result, err := nexusClient.ExecuteTemplate(ctx, req)
if err != nil {
    log.Printf("Ошибка: %v", err)
    return
}
```

## Детальная обработка

```go
result, err := nexusClient.ExecuteTemplate(ctx, req)
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

## Проверка типа ошибки

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

## Типы ошибок

### ValidationError

Ошибка валидации входных данных:

```go
if errDetail.IsValidationError() {
    // Обработка ошибки валидации
}
```

### AuthenticationError

Ошибка аутентификации (неверный токен):

```go
if errDetail.IsAuthenticationError() {
    // Обновить токен или запросить новый
}
```

### AuthorizationError

Ошибка авторизации (недостаточно прав):

```go
if errDetail.IsAuthorizationError() {
    // Проверить права доступа
}
```

### RateLimitError

Превышен лимит запросов:

```go
if errDetail.IsRateLimitError() {
    // Подождать и повторить запрос
}
```

### InternalError

Внутренняя ошибка сервера:

```go
if errDetail.IsInternalError() {
    // Логировать и повторить позже
}
```

## Retry при ошибках

SDK автоматически повторяет запросы при сетевых ошибках, если настроен RetryConfig:

```go
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
```

