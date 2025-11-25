---
id: grpc-api
title: gRPC API
sidebar_label: gRPC API
---

# gRPC API

Nexus Protocol поддерживает gRPC для высокопроизводительного взаимодействия с низкой латентностью и эффективной сериализацией.

## 🌐 Базовая информация

### Адрес подключения
```
api.nexus.dev:50051
```

### Протокол
- **Transport**: HTTP/2
- **Serialization**: Protocol Buffers 3
- **Compression**: gzip (опционально)

### Поддерживаемые версии
- **Protocol Version**: 2.0.0
- **gRPC Version**: 1.50+

## 📋 Спецификация

Полная спецификация gRPC API определена в файле [nexus.proto](../../api/grpc/nexus.proto).

### Генерация клиента

```bash
# Go
protoc --go_out=. --go-grpc_out=. nexus.proto

# Python
python -m grpc_tools.protoc --python_out=. --grpc_python_out=. nexus.proto

# Node.js
grpc_tools_node_protoc --js_out=import_style=commonjs,binary:. --grpc_out=. nexus.proto
```

## 🔐 Аутентификация

gRPC поддерживает несколько методов аутентификации:

### 1. JWT в metadata (рекомендуемый)

```go
import (
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

func createAuthenticatedContext(token string) context.Context {
    md := metadata.New(map[string]string{
        "authorization": "Bearer " + token,
    })
    return metadata.NewOutgoingContext(context.Background(), md)
}

conn, err := grpc.Dial("api.nexus.dev:50051", grpc.WithTransportCredentials(credentials.NewTLS(nil)))
client := pb.NewNexusClient(conn)

ctx := createAuthenticatedContext("your-jwt-token")
```

### 2. mTLS (для enterprise)

```go
import (
    "crypto/tls"
    "google.golang.org/grpc/credentials"
)

cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
if err != nil {
    log.Fatal(err)
}

creds := credentials.NewTLS(&tls.Config{
    Certificates: []tls.Certificate{cert},
    RootCAs:      caCertPool,
})

conn, err := grpc.Dial("api.nexus.dev:50051", grpc.WithTransportCredentials(creds))
```

## 🚀 Основные сервисы

### NexusService

Основной сервис для работы с шаблонами.

```go
// Создание клиента
conn, err := grpc.Dial("api.nexus.dev:50051", grpc.WithTransportCredentials(credentials.NewTLS(nil)))
if err != nil {
