package types

// WebhookConfig представляет конфигурацию webhook
type WebhookConfig struct {
	URL          string            `json:"url"`                     // URL для отправки webhook
	Events       []string          `json:"events"`                  // список событий для подписки
	Secret       string            `json:"secret"`                  // секрет для подписи webhook
	RetryPolicy  *WebhookRetryPolicy `json:"retry_policy,omitempty"` // политика повторных отправок
	Headers      map[string]string `json:"headers,omitempty"`       // дополнительные заголовки
	Active       bool              `json:"active,omitempty"`        // активна ли подписка
	Description  string            `json:"description,omitempty"`   // описание webhook
}

// WebhookRetryPolicy определяет политику повторных отправок
type WebhookRetryPolicy struct {
	MaxRetries    int32 `json:"max_retries,omitempty"`    // максимальное количество повторных попыток
	InitialDelay  int32 `json:"initial_delay,omitempty"`  // начальная задержка в миллисекундах
	MaxDelay      int32 `json:"max_delay,omitempty"`      // максимальная задержка в миллисекундах
	BackoffFactor float32 `json:"backoff_factor,omitempty"` // коэффициент увеличения задержки
}

// WebhookEvent представляет событие webhook
type WebhookEvent struct {
	ID        string                 `json:"id"`                   // уникальный ID события
	Event     string                 `json:"event"`               // тип события
	Timestamp int64                  `json:"timestamp"`           // время события
	Data      map[string]interface{} `json:"data,omitempty"`      // данные события
	Signature string                 `json:"signature,omitempty"` // подпись для верификации
}

// RegisterWebhookRequest представляет запрос на регистрацию webhook
type RegisterWebhookRequest struct {
	Config   *WebhookConfig   `json:"config"`
	Metadata *RequestMetadata `json:"metadata,omitempty"`
}

// RegisterWebhookResponse представляет ответ на регистрацию webhook
type RegisterWebhookResponse struct {
	WebhookID       string `json:"webhook_id"`
	Status          string `json:"status"`
	Message         string `json:"message,omitempty"`
	ResponseMetadata *ResponseMetadata `json:"response_metadata,omitempty"`
}

// ListWebhooksRequest представляет запрос на получение списка webhook
type ListWebhooksRequest struct {
	ActiveOnly bool             `json:"active_only,omitempty"` // только активные webhook
	Limit      int32            `json:"limit,omitempty"`       // лимит результатов
	Offset     int32            `json:"offset,omitempty"`      // смещение
	Metadata   *RequestMetadata `json:"metadata,omitempty"`
}

// WebhookInfo представляет информацию о webhook
type WebhookInfo struct {
	ID          string            `json:"id"`
	Config      *WebhookConfig    `json:"config"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
	LastUsedAt  int64             `json:"last_used_at,omitempty"`
	SuccessCount int32            `json:"success_count,omitempty"` // количество успешных отправок
	ErrorCount   int32            `json:"error_count,omitempty"`   // количество ошибок
}

// ListWebhooksResponse представляет ответ на запрос списка webhook
type ListWebhooksResponse struct {
	Webhooks        []WebhookInfo     `json:"webhooks"`
	Total           int32             `json:"total"`
	Limit           int32             `json:"limit"`
	Offset          int32             `json:"offset"`
	ResponseMetadata *ResponseMetadata `json:"response_metadata,omitempty"`
}

// DeleteWebhookRequest представляет запрос на удаление webhook
type DeleteWebhookRequest struct {
	WebhookID string          `json:"webhook_id"`
	Metadata  *RequestMetadata `json:"metadata,omitempty"`
}

// DeleteWebhookResponse представляет ответ на удаление webhook
type DeleteWebhookResponse struct {
	WebhookID       string `json:"webhook_id"`
	Status          string `json:"status"`
	Message         string `json:"message,omitempty"`
	ResponseMetadata *ResponseMetadata `json:"response_metadata,omitempty"`
}

// TestWebhookRequest представляет запрос на тестирование webhook
type TestWebhookRequest struct {
	WebhookID string          `json:"webhook_id"`
	Event     string          `json:"event,omitempty"`     // тестовое событие
	Data      map[string]interface{} `json:"data,omitempty"` // тестовые данные
	Metadata  *RequestMetadata `json:"metadata,omitempty"`
}

// TestWebhookResponse представляет ответ на тестирование webhook
type TestWebhookResponse struct {
	WebhookID       string `json:"webhook_id"`
	Status          string `json:"status"`
	ResponseCode    int32  `json:"response_code,omitempty"`   // HTTP код ответа
	ResponseTimeMS  int32  `json:"response_time_ms,omitempty"` // время ответа
	Error           string `json:"error,omitempty"`            // ошибка если была
	ResponseMetadata *ResponseMetadata `json:"response_metadata,omitempty"`
}
