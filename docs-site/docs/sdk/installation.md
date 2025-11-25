---
id: installation
title: Установка
sidebar_label: Установка
---

# Установка SDK

## Требования

- Go 1.18 или выше
- Доступ к Nexus Protocol API

## Установка через go get

```bash
go get github.com/pro-deploy/nexus-protocol/sdk/go
```

## Установка через go.mod

Добавьте в ваш `go.mod`:

```go
module your-module

go 1.18

require (
    github.com/pro-deploy/nexus-protocol/sdk/go v2.0.0
)
```

Затем выполните:

```bash
go mod download
go mod tidy
```

## Импорт

```go
import (
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)
```

## Проверка установки

Создайте простой тестовый файл:

```go
package main

import (
    "fmt"
    "github.com/pro-deploy/nexus-protocol/sdk/go/client"
)

func main() {
    cfg := client.Config{
        BaseURL: "http://localhost:8080",
    }
    client := client.NewClient(cfg)
    fmt.Println("SDK установлен успешно!")
}
```

Запустите:

```bash
go run main.go
```

## Зависимости

SDK использует следующие зависимости:

- `github.com/google/uuid` - генерация UUID
- `github.com/xeipuuv/gojsonschema` - валидация JSON Schema (опционально)

Все зависимости устанавливаются автоматически при установке SDK.

