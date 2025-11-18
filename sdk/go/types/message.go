package types

// RequestMessage представляет запрос в формате протокола
type RequestMessage struct {
	Metadata *RequestMetadata `json:"metadata"`
	Data     interface{}      `json:"data"`
}

// ResponseMessage представляет ответ в формате протокола
type ResponseMessage struct {
	Metadata *ResponseMetadata `json:"metadata"`
	Data     interface{}       `json:"data"`
}

// NewRequestMessage создает новый RequestMessage
func NewRequestMessage(metadata *RequestMetadata, data interface{}) *RequestMessage {
	return &RequestMessage{
		Metadata: metadata,
		Data:     data,
	}
}

