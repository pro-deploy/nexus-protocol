# Nexus Protocol Deployment

Готовые конфигурации для развертывания Nexus Protocol в различных окружениях.

## 📁 Структура

```
deployment/
├── docker-compose.yml          # Docker Compose для локальной разработки
├── DEPLOYMENT.md               # Полное руководство по развертыванию
├── README.md                   # Этот файл
└── kubernetes/
    ├── deployment.yaml         # Kubernetes Deployment для API
    ├── postgres.yaml           # PostgreSQL StatefulSet
    └── redis.yaml              # Redis Deployment
```

## 🚀 Быстрый старт

### Docker Compose (рекомендуется для разработки)

```bash
cd deployment
docker-compose up -d
```

### Kubernetes (production)

```bash
kubectl apply -f kubernetes/
```

## 📖 Документация

Полное руководство по развертыванию: [DEPLOYMENT.md](./DEPLOYMENT.md)

## 🔧 Поддерживаемые платформы

- ✅ Docker Compose
- ✅ Kubernetes
- ✅ AWS ECS Fargate
- ✅ GCP Cloud Run
- ✅ Azure Container Instances

## 📊 Мониторинг

Все конфигурации включают:
- Prometheus для метрик
- Grafana для визуализации
- Health checks
- Readiness probes

## 🔒 Безопасность

- TLS/SSL поддержка
- Secrets management
- Network policies
- Rate limiting

---

**Nexus Protocol v1.1.0** - Enterprise Ready! 🚀
