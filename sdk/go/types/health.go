package types

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ReadinessResponse представляет ответ readiness probe
type ReadinessResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    ReadinessChecks   `json:"checks"`
}

// ReadinessChecks содержит проверки готовности
type ReadinessChecks struct {
	Database   string `json:"database,omitempty"`
	Redis      string `json:"redis,omitempty"`
	AIServices string `json:"ai_services,omitempty"`
}

