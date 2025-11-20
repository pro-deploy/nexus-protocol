package client

// API пути для версии 1
const (
	// Health endpoints
	PathHealth = "/health"
	PathReady  = "/ready"

	// Templates endpoints
	PathAPIV1TemplatesExecute = "/api/v1/templates/execute"
	PathAPIV1TemplatesStatus  = "/api/v1/templates/status"
	PathAPIV1TemplatesStream  = "/api/v1/templates/stream"

	// Batch endpoints
	PathAPIV1BatchExecute = "/api/v1/batch/execute"

	// Webhooks endpoints
	PathAPIV1Webhooks = "/api/v1/webhooks"

	// IAM endpoints
	PathAPIV1AuthRegister = "/api/v1/auth/register"
	PathAPIV1AuthLogin    = "/api/v1/auth/login"
	PathAPIV1AuthRefresh  = "/api/v1/auth/refresh"
	PathAPIV1UsersProfile = "/api/v1/users/profile"

	// Conversations endpoints
	PathAPIV1Conversations = "/api/v1/conversations"

	// Analytics endpoints
	PathAPIV1AnalyticsEvents = "/api/v1/analytics/events"
	PathAPIV1AnalyticsStats  = "/api/v1/analytics/stats"
)

