# Руководство по использованию Go SDK

## Быстрый старт

### 1. Установка

```bash
go get github.com/nexus-protocol/go-sdk
```

### 2. Базовое использование

```go
package main

import (
    "fmt"
    "log"
    
    nexus "github.com/nexus-protocol/go-sdk/client"
    "github.com/nexus-protocol/go-sdk/types"
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
    ProtocolVersion: "1.0.0",
    ClientVersion:   "1.0.0",
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
metadata := types.NewRequestMetadata("1.0.0", "1.0.0")
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

## Проверка здоровья сервера

```go
if err := client.Health(); err != nil {
    log.Printf("Сервер недоступен: %v", err)
} else {
    fmt.Println("Сервер доступен")
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
    ProtocolVersion: "1.0.0", // Указывайте версию явно
    ClientVersion:   "1.0.0",
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

