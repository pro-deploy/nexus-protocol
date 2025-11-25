---
id: types
title: Types Reference
sidebar_label: Types
---

# Types Reference

Документация по типам данных SDK.

## ExecuteTemplateRequest

Запрос на выполнение шаблона.

```go
type ExecuteTemplateRequest struct {
    Query    string
    Language string
    Context  *UserContext
    Options  *ExecuteOptions
    Filters  *AdvancedFilters
    Metadata *RequestMetadata
}
```

## ExecuteTemplateResponse

Ответ на выполнение шаблона.

```go
type ExecuteTemplateResponse struct {
    ExecutionID string
    Status      string
    QueryType   string
    Sections    []Section
    Workflow    *Workflow
    Metadata    *ResponseMetadata
}
```

## UserContext

Контекст пользователя.

```go
type UserContext struct {
    UserID    string
    SessionID string
    TenantID  string
    Location  *UserLocation
    Locale    string
    Timezone  string
    Currency  string
    Region    string
}
```

## RequestMetadata

Метаданные запроса.

```go
type RequestMetadata struct {
    RequestID      string
    ProtocolVersion string
    ClientVersion  string
    ClientID       string
    ClientType     string
    Timestamp      int64
    CustomHeaders  map[string]string
}
```

## ResponseMetadata

Метаданные ответа.

```go
type ResponseMetadata struct {
    RequestID        string
    ProtocolVersion  string
    ServerVersion    string
    Timestamp        int64
    ProcessingTimeMS int64
    RateLimitInfo    *RateLimitInfo
    CacheInfo        *CacheInfo
    QuotaInfo        *QuotaInfo
}
```

## ErrorDetail

Детальная информация об ошибке.

```go
type ErrorDetail struct {
    Code    string
    Type    string
    Message string
    Field   string
    Details string
}
```

### Методы

- `IsValidationError() bool`
- `IsAuthenticationError() bool`
- `IsAuthorizationError() bool`
- `IsRateLimitError() bool`
- `IsInternalError() bool`

## BatchRequest

Запрос на выполнение batch операций.

```go
type BatchRequest struct {
    Operations []BatchOperation
    Options     *BatchOptions
}
```

## WebhookConfig

Конфигурация webhook.

```go
type WebhookConfig struct {
    URL         string
    Events      []string
    Secret      string
    RetryPolicy *WebhookRetryPolicy
    Active      bool
    Description string
}
```

