---
id: client-api
title: Client API Reference
sidebar_label: Client API
---

# Client API Reference

Полная документация по API клиента SDK.

## NewClient

Создает новый клиент с указанной конфигурацией.

```go
func NewClient(config Config) *Client
```

### Пример

```go
cfg := client.Config{
    BaseURL: "https://api.nexus.dev",
    Token:   "jwt-token",
}
client := client.NewClient(cfg)
```

## ExecuteTemplate

Выполняет контекстно-зависимый шаблон.

```go
func (c *Client) ExecuteTemplate(ctx context.Context, req *types.ExecuteTemplateRequest) (*types.ExecuteTemplateResponse, error)
```

### Пример

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
}
result, err := client.ExecuteTemplate(ctx, req)
```

## GetExecutionStatus

Получает статус выполнения шаблона.

```go
func (c *Client) GetExecutionStatus(ctx context.Context, executionID string) (*types.ExecuteTemplateResponse, error)
```

### Пример

```go
status, err := client.GetExecutionStatus(ctx, "execution-id")
```

## ExecuteBatch

Выполняет пакет операций для высокой производительности (Enterprise).

```go
func (c *Client) ExecuteBatch(ctx context.Context, req *types.BatchRequest) (*types.BatchResponse, error)
```

## RegisterWebhook

Регистрирует webhook для асинхронной обработки (Enterprise).

```go
func (c *Client) RegisterWebhook(ctx context.Context, req *types.RegisterWebhookRequest) (*types.RegisterWebhookResponse, error)
```

## Health

Проверяет здоровье сервера.

```go
func (c *Client) Health(ctx context.Context) (*types.HealthResponse, error)
```

## Ready

Проверяет готовность сервера с enterprise метриками (Enterprise).

```go
func (c *Client) Ready(ctx context.Context) (*types.ReadinessResponse, error)
```

## GetFrontendConfig

Получает активную конфигурацию фронтенда (публичный endpoint).

```go
func (c *Client) GetFrontendConfig(ctx context.Context) (*types.FrontendConfig, error)
```

## SetToken

Устанавливает новый JWT токен.

```go
func (c *Client) SetToken(token string)
```

## SetPriority

Устанавливает приоритет запросов (Enterprise).

```go
func (c *Client) SetPriority(priority string)
```

Приоритеты: `low`, `normal`, `high`, `critical`

## SetCacheControl

Устанавливает контроль кэширования (Enterprise).

```go
func (c *Client) SetCacheControl(cacheControl string)
```

Опции: `no-cache`, `cache-only`, `cache-first`, `network-first`

## Admin

Получает admin клиент для управления конфигурацией.

```go
func (c *Client) Admin() *AdminClient
```

