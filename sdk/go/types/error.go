package types

// ErrorDetail содержит детальную информацию об ошибке
type ErrorDetail struct {
	Code     string            `json:"code"`
	Type     string            `json:"type"`
	Message  string            `json:"message"`
	Field    string            `json:"field,omitempty"`
	Details  string            `json:"details,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
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

