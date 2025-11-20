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

// AnalyticsStats представляет статистику аналитики
type AnalyticsStats struct {
	PeriodDays  int32              `json:"period_days"`
	TotalEvents int32              `json:"total_events"`
	TotalUsers  int32              `json:"total_users"`
	ActiveUsers int32              `json:"active_users,omitempty"`
	EventsToday int32              `json:"events_today,omitempty"`
	TopEvents   []TopEvent         `json:"top_events,omitempty"`
	UserActivity []UserActivityDay `json:"user_activity,omitempty"`
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

