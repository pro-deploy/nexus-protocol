package types

import (
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

// ResponseMetadata содержит метаданные ответа
type ResponseMetadata struct {
	RequestID        string `json:"request_id"`
	ProtocolVersion  string `json:"protocol_version"`
	ServerVersion    string `json:"server_version"`
	Timestamp        int64  `json:"timestamp"`
	ProcessingTimeMS int32  `json:"processing_time_ms"`
}

