package types

// LogEventRequest представляет запрос логирования события
type LogEventRequest struct {
	EventType string                 `json:"event_type"`
	UserID    string                 `json:"user_id,omitempty"`
	TenantID  string                 `json:"tenant_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Metadata  *RequestMetadata       `json:"metadata,omitempty"`
}

// LogEventResponse представляет ответ логирования события
type LogEventResponse struct {
	EventID   string `json:"event_id"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
}

// GetEventsRequest представляет запрос получения событий
type GetEventsRequest struct {
	EventType string `json:"event_type,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Limit     int32  `json:"limit,omitempty"`
	Offset    int32  `json:"offset,omitempty"`
}

// GetEventsResponse представляет ответ получения событий
type GetEventsResponse struct {
	Events []AnalyticsEvent `json:"events"`
	Total  int32            `json:"total"`
	Limit  int32            `json:"limit"`
	Offset int32            `json:"offset"`
}

// AnalyticsEvent представляет событие аналитики
type AnalyticsEvent struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	UserID    string                 `json:"user_id,omitempty"`
	TenantID  string                 `json:"tenant_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp string                 `json:"timestamp"`
}

// GetStatsRequest представляет запрос получения статистики
type GetStatsRequest struct {
	UserID   string `json:"user_id,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Days     int32  `json:"days,omitempty"`
}

// ConversionMetrics содержит метрики конверсии
type ConversionMetrics struct {
	SearchToResult  float32 `json:"search_to_result,omitempty"` // конверсия поиска в результаты (0-1)
	ResultToAction  float32 `json:"result_to_action,omitempty"` // конверсия результатов в действия (0-1)
	TemplateSuccess float32 `json:"template_success,omitempty"` // успешность выполнения шаблонов (0-1)
	UserRetention   float32 `json:"user_retention,omitempty"`   // удержание пользователей (0-1)
}

// PerformanceMetrics содержит метрики производительности
type PerformanceMetrics struct {
	AvgResponseTimeMS float32 `json:"avg_response_time_ms,omitempty"` // среднее время ответа
	P95ResponseTimeMS float32 `json:"p95_response_time_ms,omitempty"` // 95-й перцентиль
	P99ResponseTimeMS float32 `json:"p99_response_time_ms,omitempty"` // 99-й перцентиль
	ErrorRate         float32 `json:"error_rate,omitempty"`           // процент ошибок (0-1)
	ThroughputRPM     int32   `json:"throughput_rpm,omitempty"`       // пропускная способность (запросов в минуту)
}

// DomainMetrics содержит метрики по домену
type DomainMetrics struct {
	RequestsCount       int32   `json:"requests_count,omitempty"`       // количество запросов
	SuccessRate         float32 `json:"success_rate,omitempty"`         // процент успешности (0-1)
	AvgResponseTimeMS   float32 `json:"avg_response_time_ms,omitempty"` // среднее время ответа
	ErrorCount          int32   `json:"error_count,omitempty"`          // количество ошибок
	CacheHitRate        float32 `json:"cache_hit_rate,omitempty"`      // процент попаданий в кэш (0-1)
	RelevanceScore      float32 `json:"relevance_score,omitempty"`      // средняя релевантность результатов (0-1)
}

// AnalyticsStats представляет статистику аналитики
type AnalyticsStats struct {
	PeriodDays        int32                      `json:"period_days"`
	TotalEvents       int32                      `json:"total_events"`
	TotalUsers        int32                      `json:"total_users"`
	ActiveUsers       int32                      `json:"active_users,omitempty"`
	EventsToday       int32                      `json:"events_today,omitempty"`
	TopEvents         []TopEvent                 `json:"top_events,omitempty"`
	UserActivity      []UserActivityDay          `json:"user_activity,omitempty"`
	ConversionMetrics *ConversionMetrics         `json:"conversion_metrics,omitempty"`
	PerformanceMetrics *PerformanceMetrics       `json:"performance_metrics,omitempty"`
	DomainBreakdown   map[string]*DomainMetrics  `json:"domain_breakdown,omitempty"`
}

// TopEvent представляет топ событие
type TopEvent struct {
	Event      string  `json:"event"`
	Count      int32   `json:"count"`
	Percentage float32 `json:"percentage"`
}

// UserActivityDay представляет активность пользователей за день
type UserActivityDay struct {
	Date        string `json:"date"`
	ActiveUsers int32  `json:"active_users"`
	TotalEvents int32  `json:"total_events"`
}

