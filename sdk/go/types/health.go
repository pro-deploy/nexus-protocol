package types

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ComponentStatus содержит статус компонента
type ComponentStatus struct {
	Status    string `json:"status"`      // healthy, degraded, unhealthy
	LatencyMS int32  `json:"latency_ms,omitempty"` // задержка в миллисекундах
	Message   string `json:"message,omitempty"`   // дополнительное сообщение
}

// ExternalServiceStatus содержит статус внешнего сервиса
type ExternalServiceStatus struct {
	Name      string `json:"name"`                 // имя сервиса
	Status    string `json:"status"`              // healthy, degraded, unhealthy
	LatencyMS int32  `json:"latency_ms,omitempty"` // задержка в миллисекундах
	Endpoint  string `json:"endpoint,omitempty"`   // endpoint сервиса
}

// CapacityInfo содержит информацию о емкости системы
type CapacityInfo struct {
	CurrentLoad       float32 `json:"current_load,omitempty"`        // текущая нагрузка (0-1)
	MaxCapacity       int64   `json:"max_capacity,omitempty"`        // максимальная емкость (запросов/сек)
	AvailableCapacity int64   `json:"available_capacity,omitempty"`  // доступная емкость (запросов/сек)
	QueueSize         int32   `json:"queue_size,omitempty"`          // размер очереди
	ActiveConnections int32   `json:"active_connections,omitempty"` // активные соединения
}

// ReadinessChecks содержит проверки готовности
type ReadinessChecks struct {
	Database   string `json:"database,omitempty"`
	Redis      string `json:"redis,omitempty"`
	AIServices string `json:"ai_services,omitempty"`
}

// ReadinessResponse представляет ответ readiness probe
type ReadinessResponse struct {
	Status    string                       `json:"status"`
	Timestamp string                       `json:"timestamp"`
	Checks    ReadinessChecks              `json:"checks"`
	Components map[string]*ComponentStatus `json:"components,omitempty"` // детальный статус компонентов
	Capacity  *CapacityInfo                `json:"capacity,omitempty"`    // информация о емкости
}

