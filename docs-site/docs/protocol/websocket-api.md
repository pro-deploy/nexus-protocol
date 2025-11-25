---
id: websocket-api
title: WebSocket API
sidebar_label: WebSocket API
---

# WebSocket API

Nexus Protocol поддерживает WebSocket для двусторонней связи в реальном времени с низкой латентностью и поддержкой streaming.

## 🌐 Базовая информация

### URL подключения

#### Development
```
ws://localhost:8080/ws?token=<jwt_token>
```

#### Production
```
wss://api.nexus.dev/ws?token=<jwt_token>
```

### Subprotocol
```
nexus-json
```

### WebSocket Version
- **RFC 6455** (WebSocket Protocol)
- **Protocol Version**: 2.0.0

## 🔌 Подключение

### JavaScript/TypeScript

```javascript
const ws = new WebSocket('wss://api.nexus.dev/ws?token=<jwt_token>', ['nexus-json']);

ws.onopen = (event) => {
  console.log('Connected to Nexus WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleMessage(message);
};

ws.onclose = (event) => {
  console.log('Disconnected:', event.code, event.reason);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};
```

### Python

```python
import websocket
import json
import threading

def on_message(ws, message):
    data = json.loads(message)
    handle_message(data)

def on_error(ws, error):
    print(f"Error: {error}")

def on_close(ws, close_status_code, close_msg):
    print(f"Closed: {close_status_code}, {close_msg}")

def on_open(ws):
    print("Connected")

ws = websocket.WebSocketApp(
    "wss://api.nexus.dev/ws?token=<jwt_token>",
    subprotocols=["nexus-json"],
    on_open=on_open,
    on_message=on_message,
    on_error=on_error,
    on_close=on_close
)

ws.run_forever()
```

### Go

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/url"
    "time"

    "github.com/gorilla/websocket"
)

type WSClient struct {
    conn *websocket.Conn
    ctx  context.Context
    cancel context.CancelFunc
}

func NewWSClient(token string) (*WSClient, error) {
    u := url.URL{
        Scheme: "wss",
        Host:   "api.nexus.dev",
        Path:   "/ws",
        RawQuery: "token=" + token,
    }

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(context.Background())

    return &WSClient{
        conn:   conn,
        ctx:    ctx,
        cancel: cancel,
    }, nil
}

func (c *WSClient) Start() {
    go c.readLoop()
}

func (c *WSClient) readLoop() {
    defer c.conn.Close()
    defer c.cancel()

    for {
        select {
        case <-c.ctx.Done():
            return
        default:
            _, message, err := c.conn.ReadMessage()
            if err != nil {
                log.Printf("Read error: %v", err)
                return
            }

            var msg map[string]interface{}
            if err := json.Unmarshal(message, &msg); err != nil {
                log.Printf("Parse error: %v", err)
                continue
            }

            c.handleMessage(msg)
        }
    }
}

func (c *WSClient) handleMessage(msg map[string]interface{}) {
    msgType := msg["type"].(string)
    switch msgType {
    case "template_result":
        // Обработка результата
    case "error":
        // Обработка ошибки
    case "heartbeat":
        // Обработка heartbeat
    }
}

func (c *WSClient) Send(message interface{}) error {
    data, err := json.Marshal(message)
    if err != nil {
        return err
    }

    return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *WSClient) Close() error {
    c.cancel()
    return c.conn.Close()
}
```

## 📋 Формат сообщений

### Структура сообщения

Все сообщения следуют единому формату:

```json
{
  "type": "message_type",
  "request_id": "uuid",
  "success": true|false,
  "metadata": {
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "timestamp": 1640995200
  },
  "payload": {
    // Данные сообщения
  },
  "error": {
    // Информация об ошибке (если success: false)
  }
}
```

### Типы сообщений

#### 1. Запросы (Client → Server)

##### execute_template

```json
{
  "type": "execute_template",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "web-app",
    "client_type": "web"
  },
  "payload": {
    "query": "хочу борщ",
    "language": "ru",
    "context": {
      "user_id": "user-123",
      "location": {
        "latitude": 55.7558,
        "longitude": 37.6173
      },
      "locale": "ru-RU",
      "currency": "RUB"
    },
    "options": {
      "timeout_ms": 30000,
      "max_results_per_domain": 5
    }
  },
  "timestamp": 1640995200
}
```

##### get_execution_status

```json
{
  "type": "get_execution_status",
  "request_id": "550e8400-e29b-41d4-a716-446655440001",
  "metadata": {
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "payload": {
    "execution_id": "exec-123"
  },
  "timestamp": 1640995200
}
```

##### batch_execute (Enterprise)

```json
{
  "type": "batch_execute",
  "request_id": "550e8400-e29b-41d4-a716-446655440002",
  "metadata": {
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "payload": {
    "requests": [
      {
        "query": "хочу борщ",
        "language": "ru"
      },
      {
        "query": "find pizza",
        "language": "en"
      }
    ],
    "options": {
      "parallel_execution": true,
      "max_concurrency": 5
    }
  },
  "timestamp": 1640995200
}
```

#### 2. Ответы (Server → Client)

##### template_result

```json
{
  "type": "template_result",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": true,
  "metadata": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 3500
  },
  "payload": {
    "execution_id": "exec-456",
    "status": "completed",
    "query_type": "information_only",
    "sections": [
      {
        "domain_id": "recipes",
        "title": "Рецепты и кулинария",
        "status": "success",
        "results": [
          {
            "id": "recipe-123",
            "type": "recipe",
            "title": "Борщ украинский",
            "description": "Классический рецепт украинского борща",
            "data": {
              "ingredients": ["свекла", "картофель", "капуста"],
              "cooking_time": "2 часа",
              "difficulty": "средне"
            },
            "relevance": 0.95,
            "confidence": 0.88
          }
        ]
      }
    ],
    "domain_analysis": {
      "selected_domains": [
        {
          "domain_id": "recipes",
          "name": "Рецепты и кулинария",
          "confidence": 0.87,
          "reason": "Высокая уверенность: найдены ключевые слова рецепта"
        }
      ]
    }
  },
  "timestamp": 1640995235
}
```

##### execution_status

```json
{
  "type": "execution_status",
  "request_id": "550e8400-e29b-41d4-a716-446655440001",
  "success": true,
  "metadata": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0"
  },
  "payload": {
    "execution_id": "exec-123",
    "status": "processing",
    "progress": 75,
    "stage": "ai_processing",
    "created_at": "2025-01-18T10:00:00Z",
    "updated_at": "2025-01-18T10:00:12Z"
  },
  "timestamp": 1640995212
}
```

##### batch_result (Enterprise)

```json
{
  "type": "batch_result",
  "request_id": "550e8400-e29b-41d4-a716-446655440002",
  "success": true,
  "metadata": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 4200
  },
  "payload": {
    "batch_id": "batch-789",
    "total_requests": 2,
    "completed_requests": 2,
    "failed_requests": 0,
    "results": [
      {
        "request_index": 0,
        "execution_id": "exec-456",
        "status": "completed",
        "sections": [...]
      },
      {
        "request_index": 1,
        "execution_id": "exec-457",
        "status": "completed",
        "sections": [...]
      }
    ]
  },
  "timestamp": 1640995240
}
```

#### 3. События и уведомления

##### template_progress

```json
{
  "type": "template_progress",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": true,
  "metadata": {
    "protocol_version": "2.0.0"
  },
  "payload": {
    "execution_id": "exec-456",
    "progress": 75,
    "stage": "ai_processing",
    "message": "Анализ доменов завершен, обработка рецептов",
    "current_domain": "recipes",
    "remaining_time_ms": 5000
  },
  "timestamp": 1640995220
}
```

##### error

```json
{
  "type": "error",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": false,
  "metadata": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0"
  },
  "error": {
    "code": "VALIDATION_FAILED",
    "type": "VALIDATION_ERROR",
    "message": "Query cannot be empty",
    "field": "query",
    "details": "The query field is required for template execution"
  },
  "timestamp": 1640995205
}
```

##### heartbeat

```json
{
  "type": "heartbeat",
  "request_id": "550e8400-e29b-41d4-a716-446655440003",
  "success": true,
  "metadata": {
    "protocol_version": "2.0.0"
  },
  "payload": {
    "timestamp": 1640995260,
    "server_time": "2025-01-18T10:01:00Z"
  },
  "timestamp": 1640995260
}
```

## ⚡ Real-time Features

### Streaming Results

WebSocket поддерживает streaming результатов для длительных операций:

```javascript
// Запрос с streaming
ws.send(JSON.stringify({
  type: "execute_template",
  request_id: "req-123",
  payload: {
    query: "длинный анализ данных",
    options: {
      enable_streaming: true,
      stream_progress: true
    }
  },
  timestamp: Date.now()
}));

// Получение прогресса
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  if (message.type === 'template_progress') {
    updateProgress(message.payload.progress);
  }

  if (message.type === 'template_result') {
    displayResult(message.payload);
  }
};
```

### Bidirectional Communication

```javascript
// Клиент может отправлять несколько запросов параллельно
const request1 = {
  type: "execute_template",
  request_id: "req-1",
  payload: { query: "рецепт борща" }
};

const request2 = {
  type: "execute_template",
  request_id: "req-2",
  payload: { query: "рецепт салата" }
};

ws.send(JSON.stringify(request1));
ws.send(JSON.stringify(request2));

// Сервер ответит на каждый request_id отдельно
```

## 🔄 Connection Management

### Heartbeat

```javascript
// Клиент отправляет heartbeat каждые 30 секунд
setInterval(() => {
  ws.send(JSON.stringify({
    type: "heartbeat",
    request_id: generateUUID(),
    timestamp: Date.now()
  }));
}, 30000);

// Сервер отвечает на heartbeat
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.type === 'heartbeat') {
    // Сервер жив
    lastHeartbeat = Date.now();
  }
};
```

### Reconnection

```javascript
class WSReconnectingClient {
  constructor(url, token) {
    this.url = url;
    this.token = token;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000; // ms
  }

  connect() {
    this.ws = new WebSocket(this.url, ['nexus-json']);

    this.ws.onopen = () => {
      console.log('Connected');
      this.reconnectAttempts = 0;
    };

    this.ws.onclose = (event) => {
      console.log('Disconnected:', event.code);
      this.handleReconnect();
    };

    this.ws.onerror = (error) => {
      console.error('Error:', error);
    };
  }

  handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);

      setTimeout(() => {
        console.log(`Reconnecting (attempt ${this.reconnectAttempts})`);
        this.connect();
      }, delay);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }
}
```

## 🛡️ Error Handling

### WebSocket Errors

```javascript
ws.onerror = (error) => {
  console.error('WebSocket error:', error);

  // Попытка переподключения
  setTimeout(() => {
    reconnect();
  }, 1000);
};

ws.onclose = (event) => {
  console.log(`Connection closed: ${event.code} - ${event.reason}`);

  switch (event.code) {
    case 1000: // Normal closure
      console.log('Normal closure');
      break;
    case 1001: // Going away
      console.log('Server going away');
      break;
    case 1006: // Abnormal closure
      console.log('Abnormal closure - reconnecting');
      reconnect();
      break;
    case 1011: // Internal server error
      console.log('Server error');
      break;
    default:
      console.log('Unknown close code:', event.code);
  }
};
```

### Protocol Errors

```javascript
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  if (!message.success && message.error) {
    handleProtocolError(message.error);
  }
};

function handleProtocolError(error) {
  console.error(`Protocol error: ${error.code} - ${error.message}`);

  switch (error.code) {
    case 'VALIDATION_FAILED':
      // Показать ошибку валидации
      showValidationError(error.field, error.message);
      break;
    case 'AUTHENTICATION_FAILED':
      // Перенаправить на логин
      redirectToLogin();
      break;
    case 'RATE_LIMIT_EXCEEDED':
      // Показать сообщение о лимите
      showRateLimitMessage(error.details);
      break;
    default:
      // Общая ошибка
      showGenericError(error.message);
  }
}
```

## 📊 Monitoring

### Connection Metrics

```javascript
class WSMonitor {
  constructor(ws) {
    this.ws = ws;
    this.metrics = {
      messagesSent: 0,
      messagesReceived: 0,
      bytesSent: 0,
      bytesReceived: 0,
      connectionTime: Date.now(),
      lastActivity: Date.now()
    };

    this.attachListeners();
  }

  attachListeners() {
    const originalSend = this.ws.send;
    this.ws.send = (data) => {
      this.metrics.messagesSent++;
      this.metrics.bytesSent += data.length;
      this.metrics.lastActivity = Date.now();
      return originalSend.call(this.ws, data);
    };

    this.ws.addEventListener('message', (event) => {
      this.metrics.messagesReceived++;
      this.metrics.bytesReceived += event.data.length;
      this.metrics.lastActivity = Date.now();
    });
  }

  getMetrics() {
    return {
      ...this.metrics,
      uptime: Date.now() - this.metrics.connectionTime,
      idleTime: Date.now() - this.metrics.lastActivity
    };
  }
}
```

## 🔧 Best Practices

### 1. Message Size Limits

```javascript
const MAX_MESSAGE_SIZE = 1024 * 1024; // 1MB

function sendMessage(ws, message) {
  const data = JSON.stringify(message);

  if (data.length > MAX_MESSAGE_SIZE) {
    throw new Error('Message too large');
  }

  ws.send(data);
}
```

### 2. Rate Limiting

```javascript
class WSRateLimiter {
  constructor(maxRequestsPerSecond = 10) {
    this.maxRequests = maxRequestsPerSecond;
    this.requests = [];
  }

  canSend() {
    const now = Date.now();
    // Удаляем старые запросы (старше 1 секунды)
    this.requests = this.requests.filter(time => now - time < 1000);

    return this.requests.length < this.maxRequests;
  }

  recordRequest() {
    this.requests.push(Date.now());
  }

  send(ws, message) {
    if (!this.canSend()) {
      throw new Error('Rate limit exceeded');
    }

    this.recordRequest();
    ws.send(JSON.stringify(message));
  }
}
```

### 3. Message Queue

```javascript
class WSMessageQueue {
  constructor(ws, options = {}) {
    this.ws = ws;
    this.queue = [];
    this.processing = false;
    this.maxRetries = options.maxRetries || 3;
    this.retryDelay = options.retryDelay || 1000;
  }

  send(message, options = {}) {
    return new Promise((resolve, reject) => {
      this.queue.push({
        message,
        resolve,
        reject,
        retries: 0,
        options
      });

      this.process();
    });
  }

  async process() {
    if (this.processing || this.queue.length === 0) {
      return;
    }

    this.processing = true;

    while (this.queue.length > 0) {
      const item = this.queue.shift();

      try {
        if (this.ws.readyState === WebSocket.OPEN) {
          this.ws.send(JSON.stringify(item.message));
          item.resolve();
        } else {
          throw new Error('WebSocket not connected');
        }
      } catch (error) {
        if (item.retries < this.maxRetries) {
          item.retries++;
          setTimeout(() => {
            this.queue.unshift(item);
            this.process();
          }, this.retryDelay * item.retries);
        } else {
          item.reject(error);
        }
        break;
      }
    }

    this.processing = false;
  }
}
```

### 4. Connection Pool

```javascript
class WSConnectionPool {
  constructor(url, options = {}) {
    this.url = url;
    this.poolSize = options.poolSize || 3;
    this.connections = [];
    this.currentIndex = 0;
  }

  async init() {
    for (let i = 0; i < this.poolSize; i++) {
      const ws = new WebSocket(this.url, ['nexus-json']);
      await new Promise((resolve, reject) => {
        ws.onopen = resolve;
        ws.onerror = reject;
      });
      this.connections.push(ws);
    }
  }

  getConnection() {
    const connection = this.connections[this.currentIndex];
    this.currentIndex = (this.currentIndex + 1) % this.connections.length;
    return connection;
  }

  send(message) {
    const ws = this.getConnection();
    ws.send(JSON.stringify(message));
  }
}
```
```

### Ответ

```json
{
  "type": "execute_template_response",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "protocol_version": "2.0.0",
    "server_version": "2.0.0",
    "processing_time_ms": 3500
  },
  "payload": {
    "execution_id": "exec-123",
    "status": "completed"
  },
  "timestamp": 1640995235
}
```

## Пример использования

```javascript
const ws = new WebSocket('wss://api.nexus.dev/ws?token=<jwt_token>', ['nexus-json']);

ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'execute_template',
    request_id: 'req-123',
    metadata: {
      protocol_version: '2.0.0',
      client_version: '2.0.0'
    },
    payload: {
      query: 'хочу борщ',
      language: 'ru'
    },
    timestamp: Date.now() / 1000
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket closed');
};
```

## Типы сообщений

- `execute_template` - выполнение шаблона
- `execute_template_response` - ответ на выполнение шаблона
- `status_update` - обновление статуса выполнения
- `error` - сообщение об ошибке

## Спецификация

Полная спецификация доступна в файле [protocol.json](../../api/websocket/protocol.json).

