# 🚀 План перехода на микросервисную архитектуру

## 📋 Обзор плана

**Цель:** Постепенная миграция от монолитной архитектуры к полноценной микросервисной системе с минимальными рисками и downtime.

**Текущая архитектура:** Монолит с 6 сервисами внутри одного процесса Go
**Целевая архитектура:** 7 независимых микросервисов + API Gateway

**Продолжительность:** 8-12 недель
**Команда:** 2-3 разработчика + DevOps
**Риски:** Минимальные (пошаговая миграция)

---

## 🎯 Этапы миграции

### **Этап 0: Подготовка (Неделя 1)**

#### **Цели:**
- ✅ Подготовить инфраструктуру
- ✅ Создать CI/CD pipelines
- ✅ Настроить мониторинг
- ✅ Протестировать микросервисную инфраструктуру

#### **Задачи:**

##### **1. Infrastructure Setup**
```bash
# Создать Kubernetes namespace
kubectl create namespace nexus-prod

# Настроить PostgreSQL и Redis кластеры
kubectl apply -f infrastructure/postgres-cluster.yml
kubectl apply -f infrastructure/redis-cluster.yml

# Настроить Keycloak
kubectl apply -f infrastructure/keycloak.yml

# Создать секреты
kubectl create secret generic nexus-secrets \
  --from-literal=jwt-secret='prod-secret-key' \
  --from-literal=db-password='prod-db-password' \
  --from-literal=redis-password='prod-redis-password' \
  --from-literal=openai-api-key='prod-openai-key'
```

##### **2. CI/CD Setup**
```yaml
# GitHub Actions workflows
.github/workflows/
├── deploy-monolith.yml     # Деплой монолита
├── deploy-ai-service.yml    # Деплой AI сервиса
├── deploy-api-gateway.yml   # Деплой Gateway
└── integration-tests.yml    # Интеграционные тесты
```

##### **3. Monitoring Setup**
```bash
# Prometheus + Grafana
kubectl apply -f monitoring/prometheus.yml
kubectl apply -f monitoring/grafana.yml

# ELK Stack для логов
kubectl apply -f monitoring/elasticsearch.yml
kubectl apply -f monitoring/logstash.yml
kubectl apply -f monitoring/kibana.yml

# Jaeger для tracing
kubectl apply -f monitoring/jaeger.yml
```

##### **4. Testing Infrastructure**
```bash
# Запустить микросервисную инфраструктуру в staging
cd server/go/docker/microservices
docker-compose -f docker-compose.microservices.yml up -d

# Протестировать service discovery
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # AI Service
curl http://localhost:8086/health  # Auth Service
```

**Критерии готовности:**
- [ ] Kubernetes кластер настроен
- [ ] PostgreSQL и Redis кластеры работают
- [ ] Keycloak realm настроен
- [ ] CI/CD pipelines созданы
- [ ] Мониторинг работает
- [ ] Микросервисная инфраструктура протестирована

---

### **Этап 1: AI Service (Неделя 2-3)**

#### **Почему AI Service первым?**
- Самый нагруженный сервис
- Независимая бизнес-логика
- Легко протестировать изоляцию

#### **Шаги миграции:**

##### **1. Code Extraction**
```bash
# Создать отдельный репозиторий или директорию
mkdir services/ai-service
cd services/ai-service

# Скопировать AI-related код
cp -r ../../server/go/internal/ai ./
cp -r ../../server/go/pkg ./
cp ../../server/go/go.mod ./
cp ../../server/go/cmd/main.go ./cmd/

# Адаптировать main.go для standalone AI service
# Убрать остальные сервисы, оставить только AI
```

##### **2. API Interface Definition**
```go
// internal/api/handlers/ai.go
type AIServiceInterface interface {
    ExecuteTemplate(ctx context.Context, req *types.ExecuteTemplateRequest) (*types.ExecuteTemplateResponse, error)
    GetTemplateStatus(ctx context.Context, executionID string) (*types.ExecuteTemplateResponse, error)
    StreamResults(ctx context.Context, executionID string, callback func(*types.DomainSection) error) error
}
```

##### **3. Database Migration**
```sql
-- Создать отдельные таблицы для AI service
CREATE SCHEMA ai_service;

-- Execution logs
CREATE TABLE ai_service.template_executions (
    id UUID PRIMARY KEY,
    user_id UUID,
    query TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP,
    completed_at TIMESTAMP,
    processing_time_ms INTEGER
);

-- Domain results
CREATE TABLE ai_service.domain_results (
    id UUID PRIMARY KEY,
    execution_id UUID REFERENCES ai_service.template_executions(id),
    domain_name VARCHAR(100),
    status VARCHAR(50),
    response_time_ms INTEGER,
    results JSONB
);
```

##### **4. Docker & Kubernetes**
```yaml
# services/ai-service/docker/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-service ./cmd

FROM alpine:latest
COPY --from=builder /app/ai-service .
EXPOSE 8080
CMD ["./ai-service"]
```

```yaml
# services/ai-service/k8s/deployment.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-service
  namespace: nexus-prod
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ai-service
  template:
    metadata:
      labels:
        app: ai-service
    spec:
      containers:
      - name: ai-service
        image: nexus-protocol/ai-service:v1.1.0
        ports:
        - containerPort: 8080
        env:
        - name: SERVICE_NAME
          value: "ai-service"
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: nexus-secrets
              key: openai-api-key
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
```

##### **5. API Gateway Configuration**
```go
// server/go/internal/api/gateway/routes.go
func (g *Gateway) routeToAIService(w http.ResponseWriter, r *http.Request) {
    // Проксировать запросы /api/v1/templates/* -> ai-service:8080
    targetURL := "http://ai-service.nexus-prod.svc.cluster.local:8080"

    // Forward request with authentication headers
    proxy := httputil.NewSingleHostReverseProxy(targetURL)
    proxy.ServeHTTP(w, r)
}
```

##### **6. Testing & Validation**
```bash
# Unit tests
cd services/ai-service
go test ./...

# Integration tests
curl -X POST http://localhost:8080/api/v1/templates/execute \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query": "test query", "language": "ru"}'

# Performance tests
ab -n 1000 -c 10 http://localhost:8080/api/v1/templates/execute
```

**Критерии завершения:**
- [ ] AI Service развернут отдельно
- [ ] API Gateway проксирует запросы
- [ ] Все тесты проходят
- [ ] Performance не ухудшилась
- [ ] Мониторинг показывает корректные метрики

---

### **Этап 2: Auth Service + API Gateway (Неделя 4-5)**

#### **Цели:**
- Выделить аутентификацию и авторизацию
- Создать полноценный API Gateway

#### **Шаги:**

##### **1. Auth Service Extraction**
```bash
mkdir services/auth-service
cd services/auth-service

# Скопировать auth-related код
cp -r ../../server/go/internal/auth ./
cp -r ../../server/go/pkg/config ./
cp -r ../../server/go/pkg/types ./

# Создать standalone main.go для auth service
```

##### **2. API Gateway Development**
```go
// services/api-gateway/main.go
func main() {
    // Routes
    r := mux.NewRouter()

    // Health checks
    r.HandleFunc("/health", healthHandler)
    r.HandleFunc("/ready", readinessHandler)

    // Auth endpoints - local processing
    r.HandleFunc("/api/v1/auth/login", loginHandler)
    r.HandleFunc("/api/v1/auth/register", registerHandler)
    r.HandleFunc("/api/v1/auth/refresh", refreshHandler)

    // Protected routes - proxy to services
    protected := r.PathPrefix("/api/v1").Subrouter()
    protected.Use(authMiddleware)

    // Route to AI Service
    protected.HandleFunc("/templates/{path:.*}", proxyToAIService)

    // Route to other services (when ready)
    protected.HandleFunc("/batch/{path:.*}", proxyToBatchService)
    protected.HandleFunc("/webhooks/{path:.*}", proxyToWebhookService)

    http.ListenAndServe(":8080", r)
}
```

##### **3. Service Mesh Setup**
```yaml
# Istio Service Mesh для routing и observability
kubectl apply -f infrastructure/istio-gateway.yml
kubectl apply -f infrastructure/istio-virtualservice.yml

# Traffic policies
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: nexus-gateway
spec:
  http:
  - match:
    - uri:
        prefix: "/api/v1/templates"
    route:
    - destination:
        host: ai-service
  - match:
    - uri:
        prefix: "/api/v1/auth"
    route:
    - destination:
        host: auth-service
```

##### **4. Authentication Flow**
```
Client → API Gateway → Auth Service → JWT Token
                    ↓
           Protected Routes → Service Mesh → Target Service
```

##### **5. Database Separation**
```sql
-- Auth service database
CREATE SCHEMA auth_service;

CREATE TABLE auth_service.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    roles TEXT[],
    created_at TIMESTAMP,
    last_login TIMESTAMP
);

-- Migrate existing users
INSERT INTO auth_service.users
SELECT id, email, password_hash, roles, created_at, last_login
FROM public.users;
```

**Критерии завершения:**
- [ ] Auth Service работает отдельно
- [ ] API Gateway маршрутизирует запросы
- [ ] JWT токены валидны между сервисами
- [ ] Service Mesh настроен
- [ ] Аутентификация не сломана

---

### **Этап 3: Batch Service (Неделя 6)**

#### **Цели:**
- Выделить параллельную обработку
- Настроить очереди и background jobs

#### **Шаги:**

##### **1. Batch Service Setup**
```bash
mkdir services/batch-service
cd services/batch-service

# Extract batch-related code
cp -r ../../server/go/internal/batch ./
# ... настройки для standalone
```

##### **2. Queue Infrastructure**
```yaml
# Redis-based queues для batch jobs
services:
  redis-queue:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis-queue-data:/data

  batch-worker:
    image: nexus-protocol/batch-service:v1.1.0
    command: ["./batch-service", "worker"]
    environment:
      - WORKER_MODE=true
      - QUEUE_NAME=batch-jobs
    depends_on:
      - redis-queue
```

##### **3. API Gateway Updates**
```go
// Add batch routes to gateway
protected.HandleFunc("/batch/{path:.*}", proxyToBatchService)

// Async processing
r.HandleFunc("/api/v1/batch/webhook", batchWebhookHandler)  // Для callback'ов
```

**Критерии завершения:**
- [ ] Batch operations работают асинхронно
- [ ] Queue система настроена
- [ ] Status tracking работает
- [ ] Performance улучшилась

---

### **Этап 4: Supporting Services (Неделя 7-8)**

#### **Цели:**
- Выделить Webhook, Analytics, Conversation сервисы
- Полная микросервисная архитектура

#### **Шаги:**

##### **1. Webhook Service**
```bash
mkdir services/webhook-service
# Extract webhook code, setup Redis для retry logic
```

##### **2. Analytics Service**
```bash
mkdir services/analytics-service
# Extract analytics, setup ClickHouse для метрик
```

##### **3. Conversation Service**
```bash
mkdir services/conversation-service
# Extract conversation logic, setup WebSocket support
```

##### **4. Shared Infrastructure Migration**
```sql
-- Migrate remaining data to service-specific schemas
-- Setup cross-service data sharing через APIs
```

**Критерии завершения:**
- [ ] Все сервисы выделены
- [ ] Cross-service communication работает
- [ ] Data consistency поддерживается
- [ ] Все APIs доступны через Gateway

---

### **Этап 5: Production Migration (Неделя 9-10)**

#### **Цели:**
- Полная миграция в production
- Zero-downtime deployment
- Rollback план

#### **Blue-Green Deployment:**
```bash
# Старая система (monolith)
kubectl apply -f monolith-deployment.yml

# Новая система (microservices)
kubectl apply -f microservices-deployment.yml

# Traffic switching
kubectl patch ingress nexus-ingress \
  -p '{"spec":{"rules":[{"host":"api.nexus.dev","http":{"paths":[{"path":"/","pathType":"Prefix","backend":{"service":{"name":"api-gateway-new","port":{"number":80}}}}]}}]}}'

# Verification
curl -f https://api.nexus.dev/health

# Rollback if needed
kubectl patch ingress nexus-ingress \
  -p '{"spec":{"rules":[{"host":"api.nexus.dev","http":{"paths":[{"path":"/","pathType":"Prefix","backend":{"service":{"name":"monolith","port":{"number":80}}}}]}}]}}'
```

#### **Data Migration:**
```bash
# Online data migration
pg_dump monolith_db > monolith_backup.sql

# Schema migration scripts
psql -f migrations/01_auth_schema.sql
psql -f migrations/02_ai_schema.sql
psql -f migrations/03_batch_schema.sql

# Data migration with validation
python migrations/migrate_users.py
python migrations/migrate_templates.py
```

#### **Monitoring & Alerting:**
```yaml
# Alert rules для микросервисов
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: nexus-microservices-alerts
spec:
  groups:
  - name: nexus
    rules:
    - alert: AIServiceDown
      expr: up{job="ai-service"} == 0
      for: 5m
      labels:
        severity: critical
    - alert: HighLatency
      expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2
      labels:
        severity: warning
```

**Критерии завершения:**
- [ ] Production traffic переключен
- [ ] Zero downtime достигнут
- [ ] Rollback план протестирован
- [ ] Monitoring и alerting работают

---

### **Этап 6: Optimization & Scaling (Неделя 11-12)**

#### **Цели:**
- Оптимизация производительности
- Auto-scaling настройка
- Cost optimization

#### **Performance Optimization:**
```yaml
# HPA для каждого сервиса
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ai-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ai-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### **Service Mesh Optimization:**
```yaml
# Circuit breakers
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: ai-service-circuit-breaker
spec:
  host: ai-service
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 10
        maxRequestsPerConnection: 10
    outlierDetection:
      consecutive5xxErrors: 3
      interval: 10s
      baseEjectionTime: 30s
```

#### **Database Optimization:**
```sql
-- Read replicas для analytics
CREATE PUBLICATION nexus_analytics FOR TABLE template_executions, user_actions;
-- Setup logical replication to analytics database
```

**Критерии завершения:**
- [ ] Auto-scaling работает
- [ ] Performance улучшилась на 30%
- [ ] Cost optimization достигнута
- [ ] Service mesh полностью настроен

---

## 📊 Метрики успеха

### **Technical Metrics:**
- **Latency:** < 200ms для 95% запросов
- **Availability:** > 99.9% uptime
- **Error Rate:** < 0.1% errors
- **Throughput:** 10,000+ req/min

### **Business Metrics:**
- **Time to Deploy:** < 15 минут
- **Rollback Time:** < 5 минут
- **Development Velocity:** +50% после миграции
- **Cost Efficiency:** -20% infrastructure costs

### **Team Metrics:**
- **Cross-team Dependencies:** Минимальные
- **Deployment Frequency:** Ежедневно
- **Lead Time:** < 1 час
- **Change Failure Rate:** < 5%

---

## 🚨 Риски и mitigation

### **High Risk:**
- **Data Consistency:** Shared database migration
  - *Mitigation:* Comprehensive testing, backup/restore procedures

- **Service Discovery:** Network communication between services
  - *Mitigation:* Service mesh (Istio), comprehensive testing

### **Medium Risk:**
- **Performance Degradation:** Network latency
  - *Mitigation:* Optimize service communication, caching

- **Increased Complexity:** Debugging distributed systems
  - *Mitigation:* Comprehensive logging, tracing, monitoring

### **Low Risk:**
- **Development Overhead:** More complex CI/CD
  - *Mitigation:* Automation, shared tooling

---

## 📋 Checklist готовности

### **Pre-Migration:**
- [ ] Infrastructure ready
- [ ] CI/CD pipelines tested
- [ ] Monitoring configured
- [ ] Team trained
- [ ] Rollback plan documented

### **Per-Service Migration:**
- [ ] Code extracted and tested
- [ ] Database schema migrated
- [ ] API contracts defined
- [ ] Integration tests passing
- [ ] Performance benchmarks met
- [ ] Monitoring alerts configured

### **Post-Migration:**
- [ ] Load testing completed
- [ ] Chaos engineering tested
- [ ] Documentation updated
- [ ] Team feedback collected
- [ ] Lessons learned documented

---

## 🎯 Success Criteria

**Миграция успешна если:**
- ✅ Все сервисы работают в production
- ✅ Performance не хуже монолита
- ✅ Time-to-market улучшился
- ✅ Команда может разрабатывать независимо
- ✅ Масштабирование работает автоматически
- ✅ Downtime = 0 минут

**Откат если:**
- ❌ Performance degradation > 20%
- ❌ Error rate increase > 5%
- ❌ Team velocity decreased
- ❌ Operational complexity too high

---

## 📞 Контакты и поддержка

**Migration Team:**
- Tech Lead: [Имя]
- DevOps: [Имя]
- QA: [Имя]

**Communication:**
- Slack: #nexus-migration
- Wiki: https://wiki.company.com/nexus-migration
- Standups: Ежедневно 10:00

**Emergency Contacts:**
- Production issues: PagerDuty
- Rollback procedures: GitHub Wiki

---

**Этот план обеспечивает controlled, low-risk миграцию к микросервисной архитектуре с максимальной надежностью и минимальными рисками!** 🚀
