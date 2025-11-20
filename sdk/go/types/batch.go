package types

// BatchOperation представляет одну операцию в batch запросе
type BatchOperation struct {
	ID      int         `json:"id"`      // уникальный ID операции в batch
	Type    string      `json:"type"`    // тип операции (execute_template, log_event, etc.)
	Request interface{} `json:"request"` // тело запроса операции
}

// BatchOptions содержит опции выполнения batch операций
type BatchOptions struct {
	StopOnError bool `json:"stop_on_error,omitempty"` // остановиться при первой ошибке
	Parallel    bool `json:"parallel,omitempty"`      // выполнять параллельно
	MaxConcurrency int32 `json:"max_concurrency,omitempty"` // максимальная параллельность
}

// BatchRequest представляет batch запрос
type BatchRequest struct {
	Operations []BatchOperation `json:"operations"`
	Options    *BatchOptions    `json:"options,omitempty"`
	Metadata   *RequestMetadata `json:"metadata,omitempty"`
}

// BatchResult представляет результат выполнения одной операции в batch
type BatchResult struct {
	OperationID int         `json:"operation_id"` // ID операции из запроса
	Success     bool        `json:"success"`      // успешность выполнения
	Data        interface{} `json:"data,omitempty"`      // результат операции (если success=true)
	Error       *ErrorDetail `json:"error,omitempty"`     // ошибка операции (если success=false)
	ExecutionTimeMS int32   `json:"execution_time_ms,omitempty"` // время выполнения операции
}

// BatchResponse представляет ответ на batch запрос
type BatchResponse struct {
	Results         []BatchResult     `json:"results"`
	Total           int32             `json:"total"`                      // всего операций
	Successful      int32             `json:"successful"`                 // успешных операций
	Failed          int32             `json:"failed"`                     // неудачных операций
	TotalTimeMS     int32             `json:"total_time_ms,omitempty"`    // общее время выполнения
	ResponseMetadata *ResponseMetadata `json:"response_metadata,omitempty"`
}
