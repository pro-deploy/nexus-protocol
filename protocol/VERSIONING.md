# Версионирование Nexus Protocol

## Обзор

Nexus Application Protocol использует **Semantic Versioning** (SemVer) для управления версиями протокола. Версионирование обеспечивает:

- ✅ Предсказуемое развитие протокола
- ✅ Обратную совместимость
- ✅ Понятные правила обновления
- ✅ Безопасную миграцию между версиями

## Формат версии

### Semantic Versioning

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
│     │    │        │           │
│     │    │        │           └─ Build metadata (опционально)
│     │    │        └───────────── Prerelease identifier (опционально)
│     │    └────────────────────── Patch: обратно совместимые исправления
│     └────────────────────────── Minor: обратно совместимые новые функции
└──────────────────────────────── Major: несовместимые изменения
```

### Примеры версий

- `1.0.0` - первая стабильная версия
- `1.1.0` - добавлена новая функция (обратно совместимая)
- `1.1.1` - исправлена ошибка (обратно совместимая)
- `2.0.0` - несовместимые изменения
- `1.0.0-alpha.1` - альфа-версия
- `1.0.0-rc.1+build.123` - release candidate

## Текущая версия

**Protocol Version:** `1.0.0`

## Правила версионирования

### MAJOR версия (несовместимые изменения)

Увеличивается при:
- ❌ Изменении формата сообщений (RequestMetadata, ResponseMetadata)
- ❌ Удалении полей из сообщений
- ❌ Изменении кодов ошибок
- ❌ Изменении обязательности полей
- ❌ Изменении поведения существующих функций

**Примеры:**
- `1.0.0` → `2.0.0`: Изменена структура RequestMetadata
- `1.0.0` → `2.0.0`: Удалено поле `client_id` из RequestMetadata
- `1.0.0` → `2.0.0`: Изменен формат ErrorDetail

### MINOR версия (обратно совместимые новые функции)

Увеличивается при:
- ✅ Добавлении новых опциональных полей
- ✅ Добавлении новых типов сообщений
- ✅ Добавлении новых кодов ошибок
- ✅ Расширении существующих enum значений

**Примеры:**
- `1.0.0` → `1.1.0`: Добавлено опциональное поле `custom_headers` в RequestMetadata
- `1.0.0` → `1.1.0`: Добавлен новый тип сообщения `subscribe`
- `1.0.0` → `1.1.0`: Добавлен новый код ошибки `TIMEOUT`

### PATCH версия (обратно совместимые исправления)

Увеличивается при:
- ✅ Исправлении ошибок в документации
- ✅ Уточнении описаний полей
- ✅ Исправлении примеров
- ✅ Улучшении валидации (без изменения поведения)

**Примеры:**
- `1.0.0` → `1.0.1`: Исправлена опечатка в документации
- `1.0.0` → `1.0.1`: Уточнено описание поля `timestamp`
- `1.0.0` → `1.0.1`: Улучшена валидация UUID

## Совместимость версий

### Правила совместимости

Версии совместимы если:

1. **Major версии совпадают** - несовместимые изменения только в Major
2. **Minor версия клиента ≤ Minor версии сервера** - клиент не может быть новее сервера

### Матрица совместимости

| Client Version | Server Version | Compatible | Notes |
|----------------|----------------|------------|-------|
| 1.0.0         | 1.0.0         | ✅        | Полная совместимость |
| 1.0.0         | 1.1.0         | ✅        | Server имеет дополнительные функции |
| 1.1.0         | 1.0.0         | ❌        | Client ожидает функций, которых нет |
| 1.0.0         | 2.0.0         | ❌        | Несовместимые изменения |
| 2.0.0         | 1.0.0         | ❌        | Client несовместим |

### Проверка совместимости

```javascript
function isCompatible(clientVersion, serverVersion) {
  const [clientMajor, clientMinor] = clientVersion.split('.').map(Number);
  const [serverMajor, serverMinor] = serverVersion.split('.').map(Number);
  
  // Major версии должны совпадать
  if (clientMajor !== serverMajor) {
    return false;
  }
  
  // Minor версия клиента не должна быть больше сервера
  if (clientMinor > serverMinor) {
    return false;
  }
  
  return true;
}
```

## Version Negotiation

### Request Metadata

Клиент указывает версию протокола в RequestMetadata:

```json
{
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "1.0.0",
    "client_version": "1.0.0",
    "timestamp": 1640995200
  }
}
```

### Response Metadata

Сервер возвращает версию протокола в ResponseMetadata:

```json
{
  "metadata": {
    "request_id": "req-123",
    "protocol_version": "1.0.0",
    "server_version": "1.0.2",
    "timestamp": 1640995235,
    "processing_time_ms": 3500
  }
}
```

### Проверка на сервере

```go
func validateProtocolVersion(clientVersion string) error {
    serverVersion := "1.0.0"
    
    if !isCompatible(clientVersion, serverVersion) {
        return &ProtocolVersionError{
            ClientVersion: clientVersion,
            ServerVersion: serverVersion,
            Message: "Protocol version mismatch",
        }
    }
    
    return nil
}
```

## Миграция между версиями

### Миграция 1.0.0 → 1.1.0 (Minor)

**Изменения:**
- Добавлено опциональное поле `custom_headers` в RequestMetadata
- Добавлен новый тип сообщения `subscribe`

**Миграция:**
1. Обновить клиент до версии 1.1.0
2. Использовать новые функции (опционально)
3. Обратная совместимость гарантирована

### Миграция 1.x.x → 2.0.0 (Major)

**Изменения:**
- Изменена структура RequestMetadata
- Удалено поле `client_id`
- Изменен формат ErrorDetail

**Миграция:**
1. **Audit:** Проверить использование устаревших полей
2. **Update:** Обновить код для новой структуры
3. **Test:** Протестировать в staging среде
4. **Deploy:** Развернуть с rollback plan

## Deprecation Policy

### Процесс устаревания

1. **Announcement:** Поле/функция помечается как deprecated в релизе
2. **Documentation:** Обновляется документация с предупреждением
3. **Alternative:** Предоставляется альтернативный способ
4. **Removal:** Удаляется через 2 major версии

### Пример

```json
// v1.0.0
{
  "metadata": {
    "request_id": "req-123",
    "client_id": "web-app"  // Актуально
  }
}

// v1.1.0 - deprecated
{
  "metadata": {
    "request_id": "req-123",
    "client_id": "web-app",  // deprecated, используйте client_type
    "client_type": "web"     // Новое поле
  }
}

// v2.0.0 - удалено
{
  "metadata": {
    "request_id": "req-123",
    "client_type": "web"     // client_id удален
  }
}
```

## Version Endpoints

### GET /api/v1/version

**Response:**
```json
{
  "protocol_version": "1.0.0",
  "server_version": "1.0.2",
  "api_version": "v1",
  "build_info": {
    "git_commit": "abc123def456",
    "build_time": "2025-01-18T10:00:00Z"
  }
}
```

## Best Practices

### 1. Всегда указывайте protocol_version

```json
{
  "metadata": {
    "protocol_version": "1.0.0"  // Обязательно
  }
}
```

### 2. Проверяйте совместимость на клиенте

```javascript
const clientVersion = "1.0.0";
const serverVersion = response.metadata.protocol_version;

if (!isCompatible(clientVersion, serverVersion)) {
  console.error("Protocol version mismatch");
  // Обработка несовместимости
}
```

### 3. Используйте Semantic Versioning

```javascript
// ✅ Правильно
const version = "1.0.0";

// ❌ Неправильно
const version = "1.0";
const version = "v1.0.0";
const version = "1.0.0-beta";
```

## См. также

- [Формат сообщений](./MESSAGE_FORMAT.md) - структура сообщений
- [Метаданные](./METADATA.md) - protocol_version в метаданных
- [Обработка ошибок](./ERROR_HANDLING.md) - PROTOCOL_VERSION_ERROR
