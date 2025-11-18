# Версионирование Nexus Application Protocol

## Обзор

Nexus Application Protocol использует **Semantic Versioning** (SemVer) для управления версиями протокола. Версионирование обеспечивает предсказуемое развитие и обратную совместимость.

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

### Текущая версия

**Protocol Version:** `1.0.0`

## Правила версионирования

### MAJOR версия (несовместимые изменения)

Увеличивается при:
- ❌ Изменении формата сообщений (RequestMetadata, ResponseMetadata)
- ❌ Удалении полей из сообщений
- ❌ Изменении кодов ошибок
- ❌ Изменении обязательности полей

**Пример:** `1.0.0` → `2.0.0`

### MINOR версия (обратно совместимые новые функции)

Увеличивается при:
- ✅ Добавлении новых опциональных полей
- ✅ Добавлении новых типов сообщений
- ✅ Добавлении новых кодов ошибок

**Пример:** `1.0.0` → `1.1.0`

### PATCH версия (обратно совместимые исправления)

Увеличивается при:
- ✅ Исправлении ошибок в документации
- ✅ Уточнении описаний полей
- ✅ Улучшении валидации

**Пример:** `1.0.0` → `1.0.1`

## Совместимость версий

### Правила совместимости

Версии совместимы если:

1. **Major версии совпадают**
2. **Minor версия клиента ≤ Minor версии сервера**

### Матрица совместимости

| Client Version | Server Version | Compatible |
|----------------|----------------|------------|
| 1.0.0         | 1.0.0         | ✅        |
| 1.0.0         | 1.1.0         | ✅        |
| 1.1.0         | 1.0.0         | ❌        |
| 1.0.0         | 2.0.0         | ❌        |

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

Клиент указывает версию протокола:

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

Сервер возвращает версию протокола:

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

### Обработка несовместимости

Если версии несовместимы, сервер возвращает ошибку:

```json
{
  "error": {
    "code": "PROTOCOL_VERSION_MISMATCH",
    "type": "PROTOCOL_VERSION_ERROR",
    "message": "Protocol version mismatch",
    "details": "Client version 1.1.0 is not compatible with server version 1.0.0"
  }
}
```

## Миграция между версиями

### Minor версия (1.0.0 → 1.1.0)

**Изменения:**
- Добавлены опциональные поля
- Добавлены новые типы сообщений

**Миграция:**
1. Обновить клиент до версии 1.1.0
2. Использовать новые функции (опционально)
3. Обратная совместимость гарантирована

### Major версия (1.x.x → 2.0.0)

**Изменения:**
- Изменена структура сообщений
- Удалены устаревшие поля

**Миграция:**
1. **Audit:** Проверить использование устаревших полей
2. **Update:** Обновить код для новой структуры
3. **Test:** Протестировать в staging среде
4. **Deploy:** Развернуть с rollback plan

## Deprecation Policy

### Процесс устаревания

1. **Announcement:** Поле помечается как deprecated
2. **Documentation:** Обновляется документация
3. **Alternative:** Предоставляется альтернатива
4. **Removal:** Удаляется через 2 major версии

### Пример

```json
// v1.0.0
{
  "metadata": {
    "client_id": "web-app"  // Актуально
  }
}

// v1.1.0 - deprecated
{
  "metadata": {
    "client_id": "web-app",  // deprecated
    "client_type": "web"     // Новое поле
  }
}

// v2.0.0 - удалено
{
  "metadata": {
    "client_type": "web"     // client_id удален
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

### 2. Проверяйте совместимость

```javascript
const clientVersion = "1.0.0";
const serverVersion = response.metadata.protocol_version;

if (!isCompatible(clientVersion, serverVersion)) {
  console.error("Protocol version mismatch");
}
```

### 3. Используйте Semantic Versioning

```javascript
// ✅ Правильно
const version = "1.0.0";

// ❌ Неправильно
const version = "1.0";
const version = "v1.0.0";
```

## См. также

- [Метаданные](../protocol/METADATA.md) - protocol_version в метаданных
- [Обработка ошибок](../protocol/ERROR_HANDLING.md) - PROTOCOL_VERSION_ERROR
- [Формат сообщений](../protocol/MESSAGE_FORMAT.md) - структура сообщений