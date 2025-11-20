# Nexus Protocol - Микросервисная архитектура

## 🏗️ Архитектура

Nexus Protocol использует **полноценную микросервисную архитектуру** с API Gateway и независимыми сервисами.

```
┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │────│  Auth Service   │
│    (Port 8080)  │    │   (Port 8086)   │
└─────────────────┘    └─────────────────┘
          │                       │
          ├───────────────────────┼───────────────────────┐
          │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   AI Service    │    │ Batch Service   │    │Webhook Service  │
│   (Port 8081)   │    │  (Port 8082)    │    │  (Port 8083)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│Conversation Svc │    │Analytics Svc   │    │   Keycloak      │
│  (Port 8085)    │    │  (Port 8084)    │    │   (Port 8081)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
                       ┌─────────────────┐
                       │  Shared Infra   │
                       │ PostgreSQL+Redis│
                       └─────────────────┘
```

## 🚀 Запуск микросервисов

### Docker Compose

```bash
cd server/go/docker/microservices
docker-compose -f docker-compose.microservices.yml up -d
```

### Kubernetes

```bash
cd server/go/docker/microservices
kubectl apply -f k8s-deployment.yml
```

## 📋 Сервисы

### 1. API Gateway (Port 8080)
- **Роль**: Единая точка входа, маршрутизация запросов
- **Функции**: Аутентификация, rate limiting, логирование
- **Зависимости**: Все остальные сервисы

### 2. AI Service (Port 8081)
- **Роль**: Context-Aware Templates, AI обработка
- **Функции**: Анализ запросов, генерация ответов по доменам
- **Зависимости**: OpenAI API, PostgreSQL, Redis

### 3. Batch Service (Port 8082)
- **Роль**: Параллельная обработка множественных операций
- **Функции**: Очереди, concurrency control, статус tracking
- **Зависимости**: AI Service, PostgreSQL, Redis

### 4. Webhook Service (Port 8083)
- **Роль**: Асинхронные уведомления
- **Функции**: Retry logic, delivery tracking, security
- **Зависимости**: PostgreSQL, Redis

### 5. Analytics Service (Port 8084)
- **Роль**: Метрики и аналитика
- **Функции**: Сбор событий, статистика, отчеты
- **Зависимости**: PostgreSQL, Redis

### 6. Conversation Service (Port 8085)
- **Роль**: AI беседы с памятью
- **Функции**: Управление диалогами, контекст, typing indicators
- **Зависимости**: AI Service, PostgreSQL, Redis

### 7. Auth Service (Port 8086)
- **Роль**: Аутентификация и авторизация
- **Функции**: JWT tokens, Keycloak integration, user management
- **Зависимости**: Keycloak, PostgreSQL, Redis

## 🔄 Взаимодействие сервисов

### Синхронное взаимодействие
```
Client → API Gateway → [Auth Service] → Target Service
```

### Асинхронное взаимодействие
```
AI Service → Webhook Service → External API
Batch Service → Analytics Service → Metrics Storage
```

### Service Discovery
- **Kubernetes**: Service discovery через DNS
- **Docker**: Через service names
- **Production**: Consul или Eureka

## 🗄️ Shared Infrastructure

### PostgreSQL
- **Роль**: Основное хранилище данных
- **Базы**: nexus_db (основная), keycloak_db (аутентификация)

### Redis
- **Роль**: Кэш, сессии, rate limiting, очереди
- **Кластер**: Для production масштабирования

### Keycloak
- **Роль**: Identity & Access Management
- **Realm**: nexus (преднастроенный)
- **Интеграция**: OpenID Connect, SAML

## 📊 Масштабирование

### Horizontal Scaling
```yaml
ai-service:
  replicas: 3  # AI требует много ресурсов

batch-service:
  replicas: 2  # Средняя нагрузка

webhook-service:
  replicas: 1  # Низкая нагрузка
```

### Vertical Scaling
```yaml
ai-service:
  resources:
    requests:
      memory: "1Gi"
      cpu: "500m"
    limits:
      memory: "2Gi"
      cpu: "2000m"
```

### Auto-scaling
```yaml
hpa:
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

## 🔒 Безопасность

### Service-to-Service
- **mTLS**: Mutual TLS между сервисами
- **Service Accounts**: Для Kubernetes
- **API Keys**: Для внешних вызовов

### External APIs
- **OAuth2**: Для OpenAI, внешних сервисов
- **JWT Tokens**: Для межсервисного общения
- **Rate Limiting**: На уровне API Gateway

## 📈 Мониторинг

### Метрики по сервисам
- **API Gateway**: Request latency, error rates, throughput
- **AI Service**: Token usage, model performance, domain stats
- **Batch Service**: Queue size, processing times, success rates
- **Webhook Service**: Delivery rates, retry counts, failures

### Логирование
- **Structured Logging**: JSON format для всех сервисов
- **Centralized**: ELK stack или Loki
- **Tracing**: Jaeger или Zipkin

## 🚀 Преимущества архитектуры

### ✅ Независимое масштабирование
- Каждый сервис можно масштабировать отдельно
- Оптимальное использование ресурсов

### ✅ Fault Isolation
- Сбой одного сервиса не влияет на другие
- Graceful degradation

### ✅ Technology Diversity
- Разные языки/фреймворки для разных задач
- Оптимальный выбор инструментов

### ✅ Independent Deployment
- CI/CD pipelines для каждого сервиса
- Blue-green deployments
- Canary releases

### ✅ Team Autonomy
- Разные команды за разные сервисы
- Независимый релиз цикл

## 🔧 Переход от монолита

Текущая реализация - **гибридный подход** (монолит + микросервисы внутри).

### Этапы перехода:

1. **Текущее состояние**: Все сервисы в одном процессе
2. **Первый этап**: Выделить AI Service (самый нагруженный)
3. **Второй этап**: Auth Service + API Gateway
4. **Третий этап**: Batch + Webhook + Analytics
5. **Финальный этап**: Полная микросервисная архитектура

### Миграция данных:
- **Database**: Shared PostgreSQL (мягкая миграция)
- **Cache**: Shared Redis (мягкая миграция)
- **API**: Backward compatibility (API Gateway)

## 📋 API Endpoints

### API Gateway (Port 8080)
```
GET    /health
GET    /ready
POST   /api/v1/auth/*
POST   /api/v1/templates/*
POST   /api/v1/batch/*
POST   /api/v1/webhooks/*
POST   /api/v1/analytics/*
POST   /api/v1/conversations/*
```

### AI Service (Port 8081)
```
POST   /execute
GET    /status/{id}
GET    /stream/{id}
```

### Batch Service (Port 8082)
```
POST   /execute
GET    /status/{id}
GET    /stats
POST   /cancel/{id}
```

### Webhook Service (Port 8083)
```
POST   /register
GET    /list
PUT    /{id}
DELETE /{id}
POST   /{id}/test
GET    /{id}/deliveries
GET    /stats
```

## 🎯 Готовность к production

- ✅ **Service Discovery**: Kubernetes DNS
- ✅ **Load Balancing**: Kubernetes Services
- ✅ **Health Checks**: Readiness/Liveness probes
- ✅ **Configuration**: ConfigMaps/Secrets
- ✅ **Logging**: Structured JSON logs
- ✅ **Monitoring**: Prometheus metrics
- ✅ **Security**: mTLS, RBAC, Network Policies

**Микросервисная архитектура полностью готова к production!** 🚀
