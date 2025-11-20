package protocol

import (
	"github.com/nexus-protocol/go-sdk/types"
)

// RequestMessage представляет запрос в формате Application Protocol.
// Application Protocol требует, чтобы все запросы были обернуты в структуру
// с полями metadata и data.
//
// Формат:
//   {
//     "metadata": { ... },
//     "data": { ... }
//   }
type RequestMessage struct {
	Metadata *types.RequestMetadata `json:"metadata"`
	Data     interface{}            `json:"data"`
}

// ResponseMessage представляет ответ в формате Application Protocol.
// Application Protocol требует, чтобы все ответы были обернуты в структуру
// с полями metadata и data.
//
// Формат:
//   {
//     "metadata": { ... },
//     "data": { ... }
//   }
type ResponseMessage struct {
	Metadata *types.ResponseMetadata `json:"metadata"`
	Data     interface{}             `json:"data"`
}

// NewRequestMessage создает новый RequestMessage в формате Application Protocol.
// Это правильный способ создания запросов согласно Application Protocol.
func NewRequestMessage(metadata *types.RequestMetadata, data interface{}) *RequestMessage {
	return &RequestMessage{
		Metadata: metadata,
		Data:     data,
	}
}

// NewResponseMessage создает новый ResponseMessage в формате Application Protocol.
func NewResponseMessage(metadata *types.ResponseMetadata, data interface{}) *ResponseMessage {
	return &ResponseMessage{
		Metadata: metadata,
		Data:     data,
	}
}

