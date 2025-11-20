package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RequestMetadata содержит метаданные запроса
type RequestMetadata struct {
	RequestID      string            `json:"request_id"`
	ProtocolVersion string           `json:"protocol_version"`
	ClientVersion   string           `json:"client_version"`
	ClientID        string           `json:"client_id,omitempty"`
	ClientType      string           `json:"client_type,omitempty"`
	Timestamp       int64            `json:"timestamp"`
	CustomHeaders   map[string]string `json:"custom_headers,omitempty"`
}

// NewRequestMetadata создает новый RequestMetadata с автогенерацией полей
func NewRequestMetadata(protocolVersion, clientVersion string) *RequestMetadata {
	return &RequestMetadata{
		RequestID:       uuid.New().String(),
		ProtocolVersion: protocolVersion,
		ClientVersion:   clientVersion,
		Timestamp:       time.Now().Unix(),
		CustomHeaders:   make(map[string]string),
	}
}

// RateLimitInfo содержит информацию о rate limiting
type RateLimitInfo struct {
	Limit    int32 `json:"limit,omitempty"`     // лимит запросов
	Remaining int32 `json:"remaining,omitempty"` // оставшиеся запросы
	ResetAt  int64 `json:"reset_at,omitempty"`  // время сброса лимита (Unix timestamp)
}

// CacheInfo содержит информацию о кэшировании
type CacheInfo struct {
	CacheHit bool   `json:"cache_hit,omitempty"` // был ли кэш
	CacheKey string `json:"cache_key,omitempty"` // ключ кэша
	CacheTTL int32  `json:"cache_ttl,omitempty"` // TTL кэша в секундах
}

// QuotaInfo содержит информацию о квотах
type QuotaInfo struct {
	QuotaUsed  int64  `json:"quota_used,omitempty"`  // использовано квоты
	QuotaLimit int64  `json:"quota_limit,omitempty"` // лимит квоты
	QuotaType  string `json:"quota_type,omitempty"`  // тип квоты (requests, data, etc.)
}

// ResponseMetadata содержит метаданные ответа
type ResponseMetadata struct {
	RequestID        string         `json:"request_id"`
	ProtocolVersion  string         `json:"protocol_version"`
	ServerVersion    string         `json:"server_version"`
	Timestamp        int64          `json:"timestamp"`
	ProcessingTimeMS int32          `json:"processing_time_ms"`
	RateLimitInfo    *RateLimitInfo `json:"rate_limit_info,omitempty"`
	CacheInfo        *CacheInfo     `json:"cache_info,omitempty"`
	QuotaInfo        *QuotaInfo     `json:"quota_info,omitempty"`
}

var (
	// versionPattern соответствует Semantic Versioning: MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$`)
	// uuidPattern соответствует UUID v4
	uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// ValidateVersion проверяет формат версии по Semantic Versioning
func ValidateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	if !versionPattern.MatchString(version) {
		return fmt.Errorf("invalid version format: %s (expected MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD])", version)
	}
	return nil
}

// ValidateUUID проверяет формат UUID v4
func ValidateUUID(uuidStr string) error {
	if uuidStr == "" {
		return fmt.Errorf("UUID cannot be empty")
	}
	if !uuidPattern.MatchString(strings.ToLower(uuidStr)) {
		return fmt.Errorf("invalid UUID format: %s (expected UUID v4)", uuidStr)
	}
	return nil
}

// IsCompatible проверяет совместимость версий протокола согласно Nexus Protocol v2.0.0
// Версии совместимы если:
// 1. Major версии совпадают (несовместимые изменения)
// 2. Minor версия клиента <= Minor версии сервера (новые функции в сервере)
func IsCompatible(clientVersion, serverVersion string) (bool, error) {
	if err := ValidateVersion(clientVersion); err != nil {
		return false, fmt.Errorf("invalid client version: %w", err)
	}
	if err := ValidateVersion(serverVersion); err != nil {
		return false, fmt.Errorf("invalid server version: %w", err)
	}

	// Парсим версии (игнорируем prerelease и build metadata)
	clientParts := strings.Split(strings.Split(clientVersion, "+")[0], "-")[0]
	serverParts := strings.Split(strings.Split(serverVersion, "+")[0], "-")[0]

	clientNums := strings.Split(clientParts, ".")
	serverNums := strings.Split(serverParts, ".")

	if len(clientNums) < 3 || len(serverNums) < 3 {
		return false, fmt.Errorf("invalid version format")
	}

	clientMajor, err := strconv.Atoi(clientNums[0])
	if err != nil {
		return false, fmt.Errorf("invalid client major version: %w", err)
	}

	clientMinor, err := strconv.Atoi(clientNums[1])
	if err != nil {
		return false, fmt.Errorf("invalid client minor version: %w", err)
	}

	serverMajor, err := strconv.Atoi(serverNums[0])
	if err != nil {
		return false, fmt.Errorf("invalid server major version: %w", err)
	}

	serverMinor, err := strconv.Atoi(serverNums[1])
	if err != nil {
		return false, fmt.Errorf("invalid server minor version: %w", err)
	}

	// Major версии должны совпадать для совместимости
	if clientMajor != serverMajor {
		return false, fmt.Errorf("incompatible protocol versions: client major %d != server major %d", clientMajor, serverMajor)
	}

	// Minor версия клиента не должна быть больше сервера
	if clientMinor > serverMinor {
		return false, fmt.Errorf("client protocol version %s requires server %d.%d.x or higher, but server supports %s",
			clientVersion, clientMajor, clientMinor, serverVersion)
	}

	return true, nil
}

// ValidateRequestMetadata валидирует RequestMetadata
func ValidateRequestMetadata(metadata *RequestMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}
	if err := ValidateUUID(metadata.RequestID); err != nil {
		return fmt.Errorf("invalid request_id: %w", err)
	}
	if err := ValidateVersion(metadata.ProtocolVersion); err != nil {
		return fmt.Errorf("invalid protocol_version: %w", err)
	}
	if err := ValidateVersion(metadata.ClientVersion); err != nil {
		return fmt.Errorf("invalid client_version: %w", err)
	}
	if metadata.Timestamp <= 0 {
		return fmt.Errorf("timestamp must be positive")
	}
	return nil
}

// ValidateResponseMetadata валидирует ResponseMetadata
func ValidateResponseMetadata(metadata *ResponseMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}
	if err := ValidateUUID(metadata.RequestID); err != nil {
		return fmt.Errorf("invalid request_id: %w", err)
	}
	if err := ValidateVersion(metadata.ProtocolVersion); err != nil {
		return fmt.Errorf("invalid protocol_version: %w", err)
	}
	if err := ValidateVersion(metadata.ServerVersion); err != nil {
		return fmt.Errorf("invalid server_version: %w", err)
	}
	if metadata.Timestamp <= 0 {
		return fmt.Errorf("timestamp must be positive")
	}
	if metadata.ProcessingTimeMS < 0 {
		return fmt.Errorf("processing_time_ms must be non-negative")
	}
	return nil
}

