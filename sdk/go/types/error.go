package types

// ErrorDetail содержит детальную информацию об ошибке
// Соответствует спецификации Nexus Protocol v2.0.0
type ErrorDetail struct {
	Code     string            `json:"code"`      // Машинно-читаемый код ошибки (UPPER_SNAKE_CASE)
	Type     string            `json:"type"`      // Категория ошибки (VALIDATION_ERROR, AUTHENTICATION_ERROR, etc.)
	Message  string            `json:"message"`   // Человеко-читаемое сообщение об ошибке
	Field    string            `json:"field,omitempty"`    // Поле, вызвавшее ошибку (для валидационных ошибок)
	Details  string            `json:"details,omitempty"`  // Детальная информация об ошибке
	Metadata map[string]string `json:"metadata,omitempty"` // Дополнительные метаданные ошибки
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// Error реализует интерфейс error
func (e *ErrorDetail) Error() string {
	return e.Message
}

// IsValidationError проверяет, является ли ошибка ошибкой валидации
func (e *ErrorDetail) IsValidationError() bool {
	return e.Type == "VALIDATION_ERROR"
}

// IsAuthenticationError проверяет, является ли ошибка ошибкой аутентификации
func (e *ErrorDetail) IsAuthenticationError() bool {
	return e.Type == "AUTHENTICATION_ERROR"
}

// IsAuthorizationError проверяет, является ли ошибка ошибкой авторизации
func (e *ErrorDetail) IsAuthorizationError() bool {
	return e.Type == "AUTHORIZATION_ERROR"
}

// IsRateLimitError проверяет, является ли ошибка ошибкой rate limit
func (e *ErrorDetail) IsRateLimitError() bool {
	return e.Type == "RATE_LIMIT_ERROR"
}

// IsInternalError проверяет, является ли ошибка внутренней ошибкой
func (e *ErrorDetail) IsInternalError() bool {
	return e.Type == "INTERNAL_ERROR"
}

// IsNotFoundError проверяет, является ли ошибка ошибкой "не найдено"
func (e *ErrorDetail) IsNotFoundError() bool {
	return e.Type == "NOT_FOUND"
}

// IsConflictError проверяет, является ли ошибка ошибкой конфликта
func (e *ErrorDetail) IsConflictError() bool {
	return e.Type == "CONFLICT"
}

// IsExternalError проверяет, является ли ошибка ошибкой внешнего сервиса
func (e *ErrorDetail) IsExternalError() bool {
	return e.Type == "EXTERNAL_ERROR"
}

// IsProtocolVersionError проверяет, является ли ошибка ошибкой несовместимости версий протокола
func (e *ErrorDetail) IsProtocolVersionError() bool {
	return e.Type == "PROTOCOL_VERSION_ERROR"
}

