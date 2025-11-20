package client

// API пути для версии 2.0.0
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
	PathAPIV1FrontendConfig  = "/api/v1/frontend/config"

	// Admin endpoints (v2.0.0 enterprise features)
	PathAPIV1AdminAIConfig        = "/api/v1/admin/ai/config"
	PathAPIV1AdminPrompts         = "/api/v1/admin/prompts"
	PathAPIV1AdminDomains         = "/api/v1/admin/domains"
	PathAPIV1AdminIntegrations    = "/api/v1/admin/integrations"
	PathAPIV1AdminFrontendConfigs = "/api/v1/admin/frontend/configs"
	PathAPIV1AdminVersion         = "/api/v1/admin/version"
)

