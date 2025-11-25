package types

// BatchRequest представляет batch запрос согласно протоколу v2.0.0
// Соответствует gRPC BatchRequest и OpenAPI спецификации
type BatchRequest struct {
	BatchID     string                  `json:"batch_id,omitempty"`     // Уникальный ID батча (генерируется сервером, если не указан)
	Requests    []*ExecuteTemplateRequest `json:"requests"`              // Запросы в батче
	BatchOptions *ExecuteOptions         `json:"batch_options,omitempty"` // Общие опции для батча
	Metadata    *RequestMetadata        `json:"metadata,omitempty"`     // Метаданные запроса
}

// BatchResponse представляет ответ на batch запрос согласно протоколу v2.0.0
// Соответствует gRPC BatchResponse и OpenAPI спецификации
type BatchResponse struct {
	BatchID         string                `json:"batch_id"`              // ID батча
	Responses       []*ExecuteTemplateResponse `json:"responses"`         // Ответы на запросы
	BatchMetadata   *BatchMetadata        `json:"batch_metadata"`        // Метаданные батча
	ResponseMetadata *ResponseMetadata    `json:"response_metadata,omitempty"` // Метаданные ответа
}

// BatchMetadata содержит метаданные выполнения batch операций
type BatchMetadata struct {
	TotalRequests        int32 `json:"total_requests"`                  // Общее количество запросов
	SuccessfulRequests   int32 `json:"successful_requests"`             // Успешные запросы
	FailedRequests       int32 `json:"failed_requests"`                  // Неудачные запросы
	StartedAt            int64 `json:"started_at"`                      // Время начала (Unix timestamp)
	CompletedAt          int64 `json:"completed_at,omitempty"`          // Время завершения (Unix timestamp)
	TotalProcessingTimeMS int32 `json:"total_processing_time_ms,omitempty"` // Общее время обработки в миллисекундах
}
