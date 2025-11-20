# Nexus Protocol Deployment Guide

Руководство по развертыванию Nexus Protocol для enterprise окружений.

## 📋 Содержание

- [Docker Compose](#docker-compose)
- [Kubernetes](#kubernetes)
- [AWS ECS Fargate](#aws-ecs-fargate)
- [GCP Cloud Run](#gcp-cloud-run)
- [Production Checklist](#production-checklist)

## 🐳 Docker Compose

### Быстрый старт

```bash
# Клонируем репозиторий
git clone https://github.com/nexus-protocol/nexus-protocol.git
cd nexus-protocol/deployment

# Создаем файл с секретами
cat > .env << EOF
JWT_SECRET=your-super-secret-jwt-key-change-in-production
GRAFANA_PASSWORD=admin
EOF

# Запускаем все сервисы
docker-compose up -d

# Проверяем статус
docker-compose ps

# Просмотр логов
docker-compose logs -f nexus-api
```

### Проверка работоспособности

```bash
# Health check
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready

# Metrics
curl http://localhost:9090/metrics
```

### Остановка

```bash
# Остановить все сервисы
docker-compose down

# Остановить и удалить volumes
docker-compose down -v
```

## ☸️ Kubernetes

### Предварительные требования

- Kubernetes cluster (версия 1.24+)
- kubectl настроен
- Доступ к namespace `nexus`

### Развертывание

```bash
# Создаем namespace
kubectl create namespace nexus

# Создаем secrets
kubectl create secret generic nexus-secrets \
  --from-literal=jwt-secret='your-jwt-secret' \
  --from-literal=database-url='postgresql://user:pass@nexus-postgres:5432/nexus_db' \
  --from-literal=redis-url='redis://nexus-redis:6379/0' \
  --from-literal=postgres-user='nexus' \
  --from-literal=postgres-password='nexus_password' \
  -n nexus

# Создаем ConfigMap
kubectl create configmap nexus-config \
  --from-file=config/app.yaml \
  -n nexus

# Развертываем PostgreSQL
kubectl apply -f kubernetes/postgres.yaml

# Развертываем Redis
kubectl apply -f kubernetes/redis.yaml

# Развертываем API
kubectl apply -f kubernetes/deployment.yaml

# Проверяем статус
kubectl get pods -n nexus
kubectl get services -n nexus
```

### Масштабирование

```bash
# Ручное масштабирование
kubectl scale deployment nexus-api --replicas=5 -n nexus

# HPA автоматически масштабирует от 3 до 10 реплик
kubectl get hpa -n nexus
```

### Мониторинг

```bash
# Просмотр логов
kubectl logs -f deployment/nexus-api -n nexus

# Просмотр метрик
kubectl port-forward svc/nexus-api 9090:9090 -n nexus
curl http://localhost:9090/metrics
```

## ☁️ AWS ECS Fargate

### Создание ECS кластера

```bash
# Создаем кластер
aws ecs create-cluster --cluster-name nexus-protocol

# Создаем task definition
aws ecs register-task-definition --cli-input-json file://aws/ecs-task-definition.json

# Создаем service
aws ecs create-service \
  --cluster nexus-protocol \
  --service-name nexus-api \
  --task-definition nexus-api:1 \
  --desired-count 3 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=ENABLED}"
```

### Конфигурация

См. `aws/ecs-task-definition.json` для полной конфигурации.

## ☁️ GCP Cloud Run

### Развертывание

```bash
# Сборка образа
gcloud builds submit --tag gcr.io/PROJECT_ID/nexus-api:1.1.0

# Развертывание
gcloud run deploy nexus-api \
  --image gcr.io/PROJECT_ID/nexus-api:1.1.0 \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --memory 2Gi \
  --cpu 2 \
  --min-instances 1 \
  --max-instances 10 \
  --set-env-vars "PROTOCOL_VERSION=1.1.0,ENABLE_ENTERPRISE_FEATURES=true"
```

## ✅ Production Checklist

### Безопасность

- [ ] Изменить все дефолтные пароли и секреты
- [ ] Настроить TLS/SSL сертификаты
- [ ] Включить firewall правила
- [ ] Настроить rate limiting
- [ ] Включить audit logging
- [ ] Настроить backup стратегию

### Производительность

- [ ] Настроить connection pooling для БД
- [ ] Настроить Redis для кэширования
- [ ] Включить CDN для статических ресурсов
- [ ] Настроить load balancing
- [ ] Оптимизировать database queries
- [ ] Настроить auto-scaling

### Мониторинг

- [ ] Настроить Prometheus для метрик
- [ ] Настроить Grafana dashboards
- [ ] Настроить alerting rules
- [ ] Включить distributed tracing
- [ ] Настроить log aggregation
- [ ] Настроить uptime monitoring

### Резервное копирование

- [ ] Настроить автоматические backups БД
- [ ] Настроить backup retention policy
- [ ] Протестировать restore процедуру
- [ ] Настроить disaster recovery plan

### Документация

- [ ] Обновить API документацию
- [ ] Создать runbook для операторов
- [ ] Документировать deployment процедуры
- [ ] Создать troubleshooting guide

## 🔧 Конфигурация

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PROTOCOL_VERSION` | Версия протокола | `1.1.0` |
| `ENABLE_ENTERPRISE_FEATURES` | Включить enterprise фичи | `true` |
| `RATE_LIMIT_ENABLED` | Включить rate limiting | `true` |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Лимит запросов в минуту | `1000` |
| `CACHE_ENABLED` | Включить кэширование | `true` |
| `CACHE_TTL_SECONDS` | TTL кэша в секундах | `300` |
| `WEBHOOK_ENABLED` | Включить webhooks | `true` |
| `BATCH_OPERATIONS_ENABLED` | Включить batch операции | `true` |
| `MAX_BATCH_SIZE` | Максимальный размер batch | `100` |
| `MAX_BATCH_CONCURRENCY` | Максимальная параллельность | `10` |

## 📊 Мониторинг

### Prometheus Metrics

- `nexus_requests_total` - общее количество запросов
- `nexus_requests_duration_seconds` - время выполнения запросов
- `nexus_rate_limit_remaining` - оставшиеся запросы в rate limit
- `nexus_cache_hits_total` - попадания в кэш
- `nexus_cache_misses_total` - промахи кэша
- `nexus_batch_operations_total` - количество batch операций
- `nexus_webhook_deliveries_total` - доставки webhooks

### Grafana Dashboards

Импортируйте готовые dashboards из `grafana/dashboards/`:
- Nexus API Overview
- Enterprise Metrics
- Performance Monitoring
- Error Tracking

## 🚨 Troubleshooting

### Проблемы с подключением к БД

```bash
# Проверка подключения
kubectl exec -it deployment/nexus-api -n nexus -- \
  psql $DATABASE_URL -c "SELECT 1"
```

### Проблемы с Redis

```bash
# Проверка Redis
kubectl exec -it deployment/nexus-redis -n nexus -- \
  redis-cli ping
```

### Высокая нагрузка

```bash
# Просмотр метрик
kubectl top pods -n nexus

# Масштабирование
kubectl scale deployment nexus-api --replicas=10 -n nexus
```

## 📞 Поддержка

Для enterprise клиентов доступна 24/7 поддержка:
- Email: support@nexus-protocol.com
- Slack: #nexus-enterprise
- Phone: +1-800-NEXUS-01

---

**Nexus Protocol v1.1.0** - готов к production deployment! 🚀
