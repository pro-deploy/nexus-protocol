---
id: step-by-step
title: Пошаговая миграция
sidebar_label: Пошаговое руководство
---

# 📋 Пошаговая миграция на Nexus Protocol v2.0.0

Это подробное руководство проведет вас через процесс миграции шаг за шагом.

## 1️⃣ Шаг 1: Оценка текущей интеграции

### Анализ использования

Определите, какие части Nexus Protocol вы используете:

```bash
# Найдите все импорты SDK
grep -r "nexus-protocol" --include="*.go" --include="*.js" --include="*.py" --include="*.java" .

# Найдите API вызовы
grep -r "api.nexus.dev" --include="*.go" --include="*.js" --include="*.py" --include="*.java" .

# Найдите обработки ошибок
grep -r "ErrorDetail\|error\." --include="*.go" --include="*.js" --include="*.py" --include="*.java" .
```

### Создайте инвентарь

**Таблица 1: Текущие интеграции**

| Компонент | Версия | Функции | Критичность |
|-----------|--------|---------|-------------|
| Web Client | 1.2.1 | ExecuteTemplate, Health | Высокая |
| Mobile App | 1.1.3 | ExecuteTemplate | Высокая |
| Admin Panel | 1.0.8 | Basic CRUD | Средняя |

**Таблица 2: API endpoints**

| Endpoint | Метод | Частота | Важность |
|----------|-------|---------|----------|
| `/templates/execute` | POST | 1000/min | Критичная |
| `/health` | GET | 1/min | Высокая |
| `/version` | GET | 1/hour | Низкая |

### Риски и зависимости

**Высокий риск:**
- Пиковые нагрузки (1000+ RPS)
- Критичные бизнес-процессы
- Отсутствие feature flags

**Низкий риск:**
- Админ панели
- Мониторинг
- Логирование

## 2️⃣ Шаг 2: Обновление зависимостей

### SDK обновление

#### Go SDK

```bash
# Проверьте текущую версию
go list -m github.com/pro-deploy/nexus-protocol/sdk/go

# Обновите до v2.0.0
go get github.com/pro-deploy/nexus-protocol/sdk/go@v2.0.0

# Обновите go.mod
go mod tidy

# Проверьте обновление
go list -m github.com/pro-deploy/nexus-protocol/sdk/go
```

#### Node.js SDK

```bash
# Проверьте текущую версию
npm list nexus-protocol

# Обновите до v2.0.0
npm update nexus-protocol@2.0.0

# Или установите явно
npm install nexus-protocol@2.0.0

# Проверьте package.json
cat package.json | grep nexus-protocol
```

#### Python SDK

```bash
# Проверьте текущую версию
pip show nexus-protocol

# Обновите до v2.0.0
pip install --upgrade nexus-protocol==2.0.0

# Проверьте версию
pip show nexus-protocol
```

### Зависимости проверки

Убедитесь, что все зависимости совместимы:

```bash
# Go: проверьте на конфликты
go mod graph | grep nexus-protocol

# Node.js: проверьте на уязвимости
npm audit

# Python: проверьте зависимости
pip check
```

## 3️⃣ Шаг 3: Изменение кода

### Минимальные изменения (обратная совместимость)

#### Измените protocol_version

```go
// Было
cfg := client.Config{
    ProtocolVersion: "1.2.1",
    ClientVersion:   "1.2.1",
}

// Стало
cfg := client.Config{
    ProtocolVersion: "2.0.0", // Обновлено
    ClientVersion:   "2.0.0", // Обновлено
}
```

```javascript
// Было
const request = {
  metadata: {
    protocol_version: "1.2.1",
    client_version: "1.2.1"
  }
};

// Стало
const request = {
  metadata: {
    protocol_version: "2.0.0", // Обновлено
    client_version: "2.0.0"    // Обновлено
  }
};
```

#### Обновите обработку ResponseMetadata

```go
// v2.0.0 добавляет enterprise поля
response, err := client.ExecuteTemplate(ctx, req)
if err != nil {
    return err
}

// Новые поля в ResponseMetadata (опциональные)
if response.ResponseMetadata != nil {
    fmt.Printf("Rate limit: %d/%d\n",
        response.ResponseMetadata.RateLimitInfo.Remaining,
        response.ResponseMetadata.RateLimitInfo.Limit)

    if response.ResponseMetadata.CacheInfo.CacheHit {
        fmt.Println("Response from cache")
    }
}
```

### Расширенные изменения (новые возможности)

#### Добавьте контекст пользователя

```go
req := &types.ExecuteTemplateRequest{
    Query:    "хочу борщ",
    Language: "ru",
    Context: &types.UserContext{ // Новое в v2.0.0
        UserID:    "user-123",
        SessionID: "session-456",
        Location: &types.UserLocation{
            Latitude:  55.7558,
            Longitude: 37.6173,
            Accuracy:  50,
        },
        Locale:   "ru-RU",
        Currency: "RUB",
        Region:   "RU",
    },
}
```

#### Используйте batch операции

```go
// Вместо множественных отдельных запросов
batchReq := &types.BatchExecuteRequest{
    Requests: []*types.ExecuteTemplateRequest{
        {Query: "хочу борщ", Language: "ru"},
        {Query: "find pizza", Language: "en"},
    },
    Options: &types.BatchOptions{
        ParallelExecution: true,
        MaxConcurrency:    5,
    },
}

results, err := client.BatchExecute(ctx, batchReq)
// 10x быстрее чем отдельные запросы!
```

#### Настройте retry логику

```go
cfg := client.Config{
    BaseURL: "https://api.nexus.dev",
    RetryConfig: &client.RetryConfig{
        MaxRetries:        3,
        InitialDelay:      100 * time.Millisecond,
        MaxDelay:          5 * time.Second,
        BackoffMultiplier: 2.0,
    },
    Logger: client.NewSimpleLogger(client.LogLevelInfo),
}
```

### Обновление обработки ошибок

```go
// v2.0.0 имеет улучшенную типизацию ошибок
result, err := client.ExecuteTemplate(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *types.ErrorDetail:
        // Структурированная ошибка протокола
        log.Printf("Protocol error [%s]: %s", e.Code, e.Message)
        if e.Field != "" {
            log.Printf("Field: %s", e.Field)
        }

        // Специфическая обработка по кодам
        switch e.Code {
        case "RATE_LIMIT_EXCEEDED":
            // Ждем и повторяем
            time.Sleep(time.Duration(e.Metadata["reset_at"]) * time.Second)
        case "VALIDATION_FAILED":
            // Показываем пользователю
            showValidationError(e.Field, e.Message)
        }

    case *types.ValidationError:
        // Ошибка валидации входных данных
        log.Printf("Validation error: %s", e.Field)

    default:
        // Другие ошибки
        log.Printf("Unknown error: %v", err)
    }
    return
}
```

## 4️⃣ Шаг 4: Тестирование

### Создайте тестовую среду

```bash
# Создайте отдельную базу данных для тестирования
export NEXUS_DB_URL="postgres://test:test@localhost:5432/nexus_test"

# Используйте тестовый API endpoint
export NEXUS_BASE_URL="https://staging-api.nexus.dev"

# Включите verbose логирование
export NEXUS_LOG_LEVEL="debug"
```

### Unit тесты

```go
func TestMigrationV2Compatibility(t *testing.T) {
    // Тест обратной совместимости
    cfg := client.Config{
        ProtocolVersion: "2.0.0",
        ClientVersion:   "2.0.0",
    }

    client := client.NewClient(cfg)

    req := &types.ExecuteTemplateRequest{
        Query: "test query",
    }

    resp, err := client.ExecuteTemplate(context.Background(), req)
    assert.NoError(t, err)
    assert.Equal(t, "2.0.0", resp.ResponseMetadata.ProtocolVersion)
}
```

### Integration тесты

```javascript
describe('Nexus Protocol v2.0 Migration', () => {
  it('should work with v2.0.0 protocol', async () => {
    const client = new NexusClient({
      protocolVersion: '2.0.0',
      clientVersion: '2.0.0'
    });

    const response = await client.executeTemplate({
      query: 'test query'
    });

    expect(response.metadata.protocol_version).toBe('2.0.0');
    expect(response.metadata.server_version).toBe('2.0.0');
  });

  it('should handle enterprise features', async () => {
    const response = await client.executeTemplate({
      query: 'test query',
      context: {
        user_id: 'test-user',
        location: { latitude: 55.7558, longitude: 37.6173 }
      }
    });

    expect(response.data.domain_analysis).toBeDefined();
    expect(response.metadata.rate_limit_info).toBeDefined();
  });
});
```

### Performance тесты

```bash
# Тест производительности batch операций
ab -n 1000 -c 10 \
  -T 'application/json' \
  -H 'Authorization: Bearer <token>' \
  -p batch_payload.json \
  https://api.nexus.dev/api/v1/batch/execute

# Сравните с обычными запросами
ab -n 1000 -c 10 \
  -T 'application/json' \
  -H 'Authorization: Bearer <token>' \
  -p single_payload.json \
  https://api.nexus.dev/api/v1/templates/execute
```

### Load тесты

Используйте инструменты для нагрузочного тестирования:

```bash
# k6 для нагрузочного тестирования
k6 run migration-load-test.js

# Artillery для HTTP тестирования
artillery run migration-test.yml
```

### Регрессионное тестирование

```bash
# Запустите полный набор тестов
npm run test:regression

# Проверьте все endpoints
./scripts/test-all-endpoints.sh

# Валидация всех JSON схем
./scripts/validate-schemas.sh
```

## 5️⃣ Шаг 5: Продакшн деплой

### Blue-Green Deployment

```
Старый трафик → API v1.x ──┐
                             ├── Load Balancer
Новый трафик → API v2.0 ──┘
```

### Feature Flags

```go
// Используйте feature flags для постепенного rollout
type FeatureFlags struct {
    UseV2Protocol     bool
    EnableBatchOps    bool
    EnableAnalytics   bool
    EnableWebhooks    bool
}

func (f *FeatureFlags) IsEnabled(flag string) bool {
    // Проверка из конфига или сервиса feature flags
    return getFeatureFlag(flag)
}

// В коде
if flags.IsEnabled("v2_protocol") {
    cfg.ProtocolVersion = "2.0.0"
} else {
    cfg.ProtocolVersion = "1.2.1" // fallback
}
```

### Rollback план

**Критические метрики для мониторинга:**

```bash
# Error rate - не более 1%
curl -s https://api.nexus.dev/metrics | grep error_rate

# Response time - не более 500ms
curl -s https://api.nexus.dev/metrics | grep response_time

# Success rate - не менее 99.9%
curl -s https://api.nexus.dev/metrics | grep success_rate
```

**Rollback шаги:**

```bash
# 1. Остановите трафик на v2.0
kubectl set image deployment/api api=nexus-api:v1.2.1

# 2. Восстановите конфигурацию
kubectl apply -f config-v1.2.1.yaml

# 3. Проверьте восстановление
curl https://api.nexus.dev/health

# 4. Восстановите трафик
kubectl scale deployment/api-v1 --replicas=3
kubectl scale deployment/api-v2 --replicas=0
```

### Пост-деплой мониторинг

#### Dashboards

Создайте dashboards для отслеживания:

- **Migration Success Rate**: процент успешных миграций
- **Performance Comparison**: сравнение производительности v1 vs v2
- **Error Rate by Version**: ошибки по версиям протокола
- **Feature Adoption**: использование новых функций

#### Alerts

```yaml
# Prometheus alerts для миграции
groups:
  - name: migration
    rules:
      - alert: MigrationErrorRateHigh
        expr: rate(errors_total{version="2.0.0"}[5m]) > 0.01
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate after v2.0 migration"

      - alert: MigrationPerformanceDegraded
        expr: histogram_quantile(0.95, rate(response_time_bucket{version="2.0.0"}[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Response time degraded after migration"
```

## 🎯 Специальные случаи

### Миграция с custom интеграциями

Если у вас кастомные интеграции:

```go
// Создайте адаптер для плавной миграции
type ProtocolAdapter struct {
    v1Client *old.Client
    v2Client *nexus.Client
    useV2    bool
}

func (a *ProtocolAdapter) ExecuteTemplate(req interface{}) (interface{}, error) {
    if a.useV2 {
        // Конвертируем запрос в v2.0 формат
        v2Req := convertToV2Request(req)
        return a.v2Client.ExecuteTemplate(context.Background(), v2Req)
    } else {
        // Используем старую версию
        return a.v1Client.ExecuteTemplate(req)
    }
}
```

### Миграция баз данных

Если протокол затрагивает хранение данных:

```sql
-- Добавьте новые колонки с default значениями
ALTER TABLE requests ADD COLUMN protocol_version VARCHAR(20) DEFAULT '1.2.1';
ALTER TABLE requests ADD COLUMN enterprise_features JSONB DEFAULT '{}';

-- Создайте индексы для новых полей
CREATE INDEX idx_requests_protocol_version ON requests(protocol_version);
CREATE INDEX idx_requests_enterprise_features ON requests USING GIN(enterprise_features);

-- Миграция существующих данных
UPDATE requests SET protocol_version = '2.0.0' WHERE created_at > '2025-01-01';
```

### Миграция конфигураций

```yaml
# config-v1.2.yaml
api:
  version: "1.2.1"
  features:
    - basic_templates
    - error_handling

# config-v2.0.yaml
api:
  version: "2.0.0"
  features:
    - basic_templates
    - error_handling
    - batch_operations      # NEW
    - webhooks             # NEW
    - analytics            # NEW
    - enterprise_metrics   # NEW
```

## 🆘 Troubleshooting

### Распространенные проблемы

#### 1. Authentication errors

```
Ошибка: UNAUTHENTICATED
Решение: Проверьте JWT токен и его валидность
```

```bash
# Проверьте токен
curl -H "Authorization: Bearer <token>" https://api.nexus.dev/health

# Сгенерируйте новый токен
curl -X POST https://auth.nexus.dev/token \
  -d '{"username":"user","password":"pass"}'
```

#### 2. Protocol version mismatch

```
Ошибка: PROTOCOL_VERSION_MISMATCH
Решение: Обновите client_version или protocol_version
```

#### 3. Rate limiting

```
Ошибка: RATE_LIMIT_EXCEEDED
Решение: Проверьте limits и добавьте задержки
```

```go
// Добавьте exponential backoff
backoff := time.Second
for retries := 0; retries < 3; retries++ {
    resp, err := client.ExecuteTemplate(ctx, req)
    if err != nil && isRateLimitError(err) {
        time.Sleep(backoff)
        backoff *= 2
        continue
    }
    return resp, err
}
```

### Логи и отладка

```bash
# Включите debug логирование
export NEXUS_LOG_LEVEL=debug
export NEXUS_DEBUG=true

# Проверьте логи
kubectl logs -f deployment/nexus-api

# Используйте distributed tracing
curl https://api.nexus.dev/debug/trace/<request_id>
```

## 📞 Поддержка

### Быстрая помощь

- 📖 **[Документация](../)** - подробные гайды
- 💬 **[Slack Community](https://nexus-protocol.slack.com)** - живое общение
- 🐛 **[GitHub Issues](https://github.com/nexus-protocol/nexus-protocol/issues)** - багрепорты

### Enterprise поддержка

Для enterprise клиентов:

- 🚀 **Migration Assessment** - анализ вашей интеграции
- 👥 **Dedicated Engineer** - персональный инженер поддержки
- 📞 **24/7 Hotline** - круглосуточная поддержка
- 🎯 **Migration Workshop** - очные воркшопы

[Связаться с enterprise поддержкой](mailto:enterprise@nexus.dev)

---

## ✅ Checklist завершения миграции

### Pre-migration
- [ ] Оценка текущей интеграции завершена
- [ ] Риски идентифицированы и mitigated
- [ ] План rollback подготовлен
- [ ] Feature flags настроены

### Migration
- [ ] Зависимости обновлены
- [ ] Код изменен согласно гайдам
- [ ] Protocol version обновлена до 2.0.0
- [ ] Enterprise функции протестированы

### Testing
- [ ] Unit тесты проходят
- [ ] Integration тесты проходят
- [ ] Performance тесты в норме
- [ ] Load тесты успешны

### Production
- [ ] Blue-green deployment настроен
- [ ] Monitoring dashboards активны
- [ ] Alerts настроены
- [ ] Rollback процедуры документированы

### Post-migration
- [ ] Метрики стабильны
- [ ] Пользователи не жалуются
- [ ] Документация обновлена
- [ ] Команда обучена новым функциям

**Миграция завершена успешно! 🎉**
