---
id: analytics-index
title: Аналитика и метрики
sidebar_label: Обзор аналитики
slug: /analytics
---

# 📊 Аналитика и метрики Nexus Protocol

Nexus Protocol предоставляет комплексные инструменты для мониторинга, аналитики и оптимизации производительности ваших интеграций.

## 🎯 Обзор возможностей

### Core Analytics (v2.0.0)

<span className="feature-badge feature-badge--enterprise">ENTERPRISE</span>

- **📈 Performance Metrics**: Детальный мониторинг производительности
- **🔍 Usage Analytics**: Анализ использования API
- **💰 Business Intelligence**: Метрики конверсии и ROI
- **🚨 Real-time Alerts**: Автоматические оповещения
- **📊 Custom Dashboards**: Персонализированные дашборды
- **🔗 Integration Monitoring**: Отслеживание интеграций

## 📋 Доступные метрики

### Performance Metrics

| Метрика | Описание | Частота | Важность |
|---------|----------|---------|----------|
| `response_time` | Время ответа API | Real-time | Критичная |
| `error_rate` | Процент ошибок | 1 min | Высокая |
| `throughput` | RPS (запросов в секунду) | 1 min | Высокая |
| `success_rate` | Процент успешных запросов | 1 min | Критичная |
| `latency_p95` | 95-й перцентиль латентности | 5 min | Высокая |

### Business Metrics

| Метрика | Описание | Частота | Важность |
|---------|----------|---------|----------|
| `conversion_rate` | Конверсия запросов в действия | Hourly | Высокая |
| `user_engagement` | Уровень вовлеченности | Daily | Средняя |
| `domain_usage` | Использование доменов | Hourly | Средняя |
| `geographic_distribution` | Географическое распределение | Daily | Низкая |

### System Metrics

| Метрика | Описание | Частота | Важность |
|---------|----------|---------|----------|
| `cpu_usage` | Использование CPU | 30 sec | Средняя |
| `memory_usage` | Использование памяти | 30 sec | Средняя |
| `disk_usage` | Использование диска | 5 min | Низкая |
| `network_io` | Сетевой трафик | 1 min | Средняя |

## 🚀 Быстрый старт с аналитикой

### 1. Включение метрик

```bash
# В конфигурации сервера
export NEXUS_METRICS_ENABLED=true
export NEXUS_METRICS_ENDPOINT=/metrics
export NEXUS_ANALYTICS_ENABLED=true
```

### 2. Получение базовых метрик

```bash
# Prometheus формат
curl http://localhost:8080/metrics

# JSON формат
curl http://localhost:8080/analytics/summary
```

### 3. Создание первого дашборда

```json
{
  "dashboard": {
    "name": "API Performance",
    "widgets": [
      {
        "type": "line_chart",
        "metric": "response_time",
        "title": "Response Time Trend",
        "period": "1h"
      },
      {
        "type": "bar_chart",
        "metric": "error_rate",
        "title": "Error Rate by Endpoint",
        "period": "24h"
      }
    ]
  }
}
```

## 📊 Дашборды

### Предустановленные дашборды

#### 1. API Performance Dashboard

![API Performance Dashboard](https://via.placeholder.com/800x400/0066cc/ffffff?text=API+Performance+Dashboard)

**Метрики:**
- Response time по endpoints
- Error rate trends
- Throughput graphs
- Latency percentiles

#### 2. Business Intelligence Dashboard

![Business Intelligence Dashboard](https://via.placeholder.com/800x400/00cc66/ffffff?text=Business+Intelligence+Dashboard)

**Метрики:**
- Conversion rates
- Domain usage statistics
- Geographic distribution
- User engagement metrics

#### 3. System Health Dashboard

![System Health Dashboard](https://via.placeholder.com/800x400/ff6600/ffffff?text=System+Health+Dashboard)

**Метрики:**
- CPU/Memory usage
- Disk space monitoring
- Network I/O
- Service availability

### Кастомные дашборды

```json
// Создание кастомного дашборда
POST /analytics/dashboards
{
  "name": "My Custom Dashboard",
  "description": "Dashboard for my specific use case",
  "widgets": [
    {
      "type": "time_series",
      "metric": "custom_metric",
      "title": "My Custom Metric",
      "config": {
        "period": "7d",
        "aggregation": "avg",
        "filters": {
          "client_id": "my-app"
        }
      }
    }
  ],
  "permissions": {
    "read": ["admin", "developer"],
    "write": ["admin"]
  }
}
```

## 📈 Real-time Analytics API

### Получение метрик

```bash
# Текущие метрики
GET /analytics/metrics/current

# Исторические данные
GET /analytics/metrics/history?period=24h&metrics=response_time,error_rate

# Агрегированные данные
GET /analytics/metrics/aggregate?period=7d&group_by=hour&metrics=throughput
```

### Примеры ответов

#### Current Metrics

```json
{
  "timestamp": "2025-01-18T10:00:00Z",
  "metrics": {
    "response_time": {
      "value": 245.5,
      "unit": "ms",
      "trend": "stable"
    },
    "error_rate": {
      "value": 0.005,
      "unit": "percent",
      "trend": "decreasing"
    },
    "throughput": {
      "value": 1250,
      "unit": "rps",
      "trend": "increasing"
    }
  }
}
```

#### Historical Data

```json
{
  "period": "24h",
  "interval": "1h",
  "data": [
    {
      "timestamp": "2025-01-17T10:00:00Z",
      "response_time": 245.5,
      "error_rate": 0.005,
      "throughput": 1250
    },
    {
      "timestamp": "2025-01-17T11:00:00Z",
      "response_time": 267.8,
      "error_rate": 0.003,
      "throughput": 1180
    }
  ]
}
```

## 🚨 Alerts и оповещения

### Настройка алертов

```json
POST /analytics/alerts
{
  "name": "High Error Rate Alert",
  "description": "Alert when error rate exceeds 1%",
  "condition": {
    "metric": "error_rate",
    "operator": ">",
    "threshold": 0.01,
    "period": "5m"
  },
  "channels": [
    {
      "type": "email",
      "recipients": ["admin@nexus.dev", "dev-team@nexus.dev"]
    },
    {
      "type": "slack",
      "webhook": "https://hooks.slack.com/...",
      "channel": "#alerts"
    }
  ],
  "cooldown": "10m"
}
```

### Типы алертов

#### Performance Alerts
- Response time degradation
- Error rate spikes
- Throughput drops
- Latency percentile violations

#### Business Alerts
- Conversion rate drops
- Domain failures
- Geographic outages

#### System Alerts
- High resource usage
- Service unavailability
- Disk space warnings

## 🔍 Детальный анализ

### Query Analytics

```sql
-- Анализ медленных запросов
SELECT
  request_id,
  endpoint,
  response_time,
  error_code,
  client_id,
  timestamp
FROM api_requests
WHERE response_time > 1000
  AND timestamp > NOW() - INTERVAL '24 hours'
ORDER BY response_time DESC
LIMIT 100;

-- Анализ ошибок по типам
SELECT
  error_code,
  COUNT(*) as count,
  AVG(response_time) as avg_response_time
FROM api_requests
WHERE error_code IS NOT NULL
  AND timestamp > NOW() - INTERVAL '7 days'
GROUP BY error_code
ORDER BY count DESC;
```

### User Behavior Analytics

```sql
-- Анализ паттернов использования
SELECT
  client_id,
  COUNT(*) as total_requests,
  AVG(response_time) as avg_response_time,
  SUM(CASE WHEN error_code IS NULL THEN 1 ELSE 0 END) / COUNT(*)::float as success_rate
FROM api_requests
WHERE timestamp > NOW() - INTERVAL '30 days'
GROUP BY client_id
ORDER BY total_requests DESC;

-- Географический анализ
SELECT
  country,
  region,
  COUNT(*) as requests,
  AVG(response_time) as avg_response_time
FROM api_requests r
JOIN client_locations l ON r.client_id = l.client_id
WHERE timestamp > NOW() - INTERVAL '7 days'
GROUP BY country, region
ORDER BY requests DESC;
```

## 📊 Интеграция с внешними системами

### Prometheus Integration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'nexus-api'
    static_configs:
      - targets: ['api.nexus.dev:8080']
    metrics_path: '/metrics'
    params:
      format: ['prometheus']
```

### Grafana Dashboards

```json
{
  "dashboard": {
    "title": "Nexus Protocol Analytics",
    "panels": [
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(nexus_response_time_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

### ELK Stack Integration

```json
// Logstash конфигурация
input {
  http {
    host => "0.0.0.0"
    port => 8080
    codec => json
  }
}

filter {
  json {
    source => "message"
  }

  mutate {
    add_field => {
      "service" => "nexus-api"
      "version" => "%{[metadata][protocol_version]}"
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "nexus-analytics-%{+YYYY.MM.dd}"
  }
}
```

## 🎯 Business Intelligence

### Conversion Funnels

```sql
-- Анализ конверсии запросов
WITH request_flow AS (
  SELECT
    session_id,
    MIN(CASE WHEN endpoint = '/templates/execute' THEN timestamp END) as execute_time,
    MIN(CASE WHEN endpoint = '/commerce/purchase' THEN timestamp END) as purchase_time,
    MIN(CASE WHEN endpoint = '/delivery/order' THEN timestamp END) as delivery_time
  FROM api_requests
  WHERE timestamp > NOW() - INTERVAL '30 days'
  GROUP BY session_id
)
SELECT
  COUNT(*) as total_sessions,
  COUNT(execute_time) as executed_templates,
  COUNT(purchase_time) as made_purchases,
  COUNT(delivery_time) as ordered_delivery,
  ROUND(COUNT(purchase_time)::float / COUNT(execute_time) * 100, 2) as conversion_rate
FROM request_flow;
```

### ROI Calculation

```sql
-- Расчет ROI по клиентам
SELECT
  c.client_id,
  c.name,
  COUNT(r.request_id) as total_requests,
  SUM(c.revenue) as total_revenue,
  SUM(c.costs) as total_costs,
  ROUND((SUM(c.revenue) - SUM(c.costs)) / SUM(c.costs) * 100, 2) as roi_percentage
FROM clients c
LEFT JOIN api_requests r ON c.client_id = r.client_id
WHERE r.timestamp > NOW() - INTERVAL '90 days'
GROUP BY c.client_id, c.name
ORDER BY roi_percentage DESC;
```

## 📱 Real-time Dashboards

### Live Metrics Component

```jsx
import React, { useState, useEffect } from 'react';

const LiveMetrics = () => {
  const [metrics, setMetrics] = useState({});

  useEffect(() => {
    const ws = new WebSocket('wss://api.nexus.dev/analytics/live');

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setMetrics(data);
    };

    return () => ws.close();
  }, []);

  return (
    <div className="live-metrics">
      <div className="metric-card">
        <h3>RPS</h3>
        <div className="metric-value">{metrics.throughput || 0}</div>
        <div className="metric-trend">+5.2%</div>
      </div>

      <div className="metric-card">
        <h3>Response Time</h3>
        <div className="metric-value">{metrics.response_time || 0}ms</div>
        <div className="metric-trend">-2.1%</div>
      </div>

      <div className="metric-card">
        <h3>Error Rate</h3>
        <div className="metric-value">{(metrics.error_rate || 0) * 100}%</div>
        <div className="metric-trend">0.0%</div>
      </div>
    </div>
  );
};
```

## 🔧 API для аналитики

### Analytics Endpoints

#### GET /analytics/summary
Общая сводка метрик

#### GET /analytics/metrics/\{metric_name\}
Детальные метрики по имени

#### GET /analytics/dashboards
Список дашбордов

#### POST /analytics/dashboards
Создание дашборда

#### GET /analytics/reports/\{report_id\}
Генерация отчетов

#### POST /analytics/alerts
Создание алерта

#### WebSocket /analytics/live
Real-time метрики

### Authentication

```bash
# Для analytics API требуется специальный токен
curl -H "Authorization: Bearer <analytics-token>" \
     https://api.nexus.dev/analytics/summary
```

## 📈 Advanced Analytics

### Predictive Analytics

```python
# ML для предсказания нагрузки
import pandas as pd
from sklearn.ensemble import RandomForestRegressor

# Загрузка исторических данных
data = pd.read_csv('api_metrics.csv')

# Features для предсказания
features = ['hour', 'day_of_week', 'is_holiday', 'previous_load']
X = data[features]
y = data['predicted_load']

# Обучение модели
model = RandomForestRegressor(n_estimators=100)
model.fit(X, y)

# Предсказание
next_hour_load = model.predict([[14, 1, False, current_load]])
```

### Anomaly Detection

```python
# Обнаружение аномалий
from sklearn.ensemble import IsolationForest

# Обучение на нормальных данных
model = IsolationForest(contamination=0.1)
model.fit(normal_data)

# Обнаружение аномалий
anomaly_score = model.predict(current_metrics)
if anomaly_score == -1:
    alert_anomaly(current_metrics)
```

## 🎯 Best Practices

### 1. Metrics Collection
- Собирайте все ключевые метрики
- Используйте подходящие aggregation интервалы
- Храните исторические данные

### 2. Alert Configuration
- Настройте meaningful thresholds
- Избегайте alert fatigue
- Используйте разные severity levels

### 3. Dashboard Design
- Фокус на важных метриках
- Используйте подходящие визуализации
- Регулярно обновляйте дашборды

### 4. Performance Optimization
- Используйте кэширование метрик
- Оптимизируйте запросы к данным
- Масштабируйте инфраструктуру

## 🔗 Связанные ресурсы

- [API Documentation](../protocol/rest-api) - REST API endpoints
- [SDK Documentation](../sdk/analytics) - Analytics SDK
- [Monitoring Guide](../deployment/monitoring) - Настройка мониторинга
- [Grafana Integration](https://grafana.com/docs/grafana/latest/) - Grafana docs

---

## 📞 Поддержка

- 📧 **Analytics Support**: analytics@nexus.dev
- 📖 **[Analytics API Docs](../sdk/analytics)**
- 🎯 **[Enterprise Analytics](mailto:enterprise@nexus.dev)** - для enterprise клиентов
