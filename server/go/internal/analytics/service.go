package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
)

// Service handles analytics and metrics collection
type Service struct {
	logger *zap.Logger
	// In production, add database/cache clients
}

// Event represents an analytics event
type Event struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	SessionID   string                 `json:"session_id,omitempty"`
	EventType   string                 `json:"event_type"`
	EventData   map[string]interface{} `json:"event_data"`
	Timestamp   time.Time              `json:"timestamp"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	TenantID    string                 `json:"tenant_id,omitempty"`
}

// ConversionMetrics represents conversion metrics
type ConversionMetrics struct {
	SearchToResult      float64 `json:"search_to_result"`
	ResultToAction      float64 `json:"result_to_action"`
	TemplateSuccess     float64 `json:"template_success"`
	AverageSessionTime  int64   `json:"average_session_time_seconds"`
	BounceRate          float64 `json:"bounce_rate"`
	ReturnVisitorRate   float64 `json:"return_visitor_rate"`
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	AverageResponseTime    int64 `json:"average_response_time_ms"`
	MedianResponseTime     int64 `json:"median_response_time_ms"`
	P95ResponseTime        int64 `json:"p95_response_time_ms"`
	P99ResponseTime        int64 `json:"p99_response_time_ms"`
	RequestsPerSecond      float64 `json:"requests_per_second"`
	ErrorRate              float64 `json:"error_rate"`
	ThroughputRPM          int64   `json:"throughput_rpm"`
	CacheHitRate           float64 `json:"cache_hit_rate"`
}

// DomainMetrics represents metrics for specific domains
type DomainMetrics struct {
	DomainName         string  `json:"domain_name"`
	RequestCount       int64   `json:"request_count"`
	SuccessRate        float64 `json:"success_rate"`
	AverageResponseTime int64  `json:"average_response_time_ms"`
	UniqueUsers        int64   `json:"unique_users"`
	PopularQueries     []string `json:"popular_queries"`
}

// NewService creates a new analytics service
func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// LogEvent logs an analytics event
func (s *Service) LogEvent(ctx context.Context, event *Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	s.logger.Info("Analytics event logged",
		zap.String("event_id", event.ID),
		zap.String("user_id", event.UserID),
		zap.String("event_type", event.EventType),
		zap.Time("timestamp", event.Timestamp))

	// In production, this would save to database/clickhouse/analytics platform
	// For now, just log the event

	return nil
}

// LogTemplateExecution logs template execution events
func (s *Service) LogTemplateExecution(ctx context.Context, userID, requestID string, request *types.ExecuteTemplateRequest, response *types.ExecuteTemplateResponse) error {
	event := &Event{
		UserID:    userID,
		RequestID: requestID,
		EventType: "template.execution",
		EventData: map[string]interface{}{
			"query":             request.Query,
			"language":          request.Language,
			"execution_id":      response.ExecutionID,
			"status":           response.Status,
			"processing_time_ms": response.ProcessingTimeMS,
			"sections_count":   len(response.Sections),
			"query_type":       response.QueryType,
		},
		Timestamp: time.Now(),
	}

	// Add domain-specific data
	if len(response.Sections) > 0 {
		domains := make([]string, len(response.Sections))
		for i, section := range response.Sections {
			domains[i] = section.DomainID
		}
		event.EventData["domains"] = domains
	}

	return s.LogEvent(ctx, event)
}

// LogUserAction logs user actions (clicks, purchases, etc.)
func (s *Service) LogUserAction(ctx context.Context, userID, requestID, actionType string, actionData map[string]interface{}) error {
	event := &Event{
		UserID:    userID,
		RequestID: requestID,
		EventType: "user.action",
		EventData: map[string]interface{}{
			"action_type": actionType,
		},
		Timestamp: time.Now(),
	}

	// Merge action data
	for k, v := range actionData {
		event.EventData[k] = v
	}

	return s.LogEvent(ctx, event)
}

// LogError logs error events
func (s *Service) LogError(ctx context.Context, userID, requestID string, errorType, errorMessage string, metadata map[string]interface{}) error {
	event := &Event{
		UserID:    userID,
		RequestID: requestID,
		EventType: "error",
		EventData: map[string]interface{}{
			"error_type":    errorType,
			"error_message": errorMessage,
			"metadata":      metadata,
		},
		Timestamp: time.Now(),
	}

	return s.LogEvent(ctx, event)
}

// GetConversionMetrics returns conversion metrics
func (s *Service) GetConversionMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) (*ConversionMetrics, error) {
	// Mock implementation - in production would query analytics database
	return &ConversionMetrics{
		SearchToResult:     0.85,
		ResultToAction:     0.65,
		TemplateSuccess:    0.92,
		AverageSessionTime: 180,
		BounceRate:         0.15,
		ReturnVisitorRate:  0.35,
	}, nil
}

// GetPerformanceMetrics returns performance metrics
func (s *Service) GetPerformanceMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) (*PerformanceMetrics, error) {
	// Mock implementation
	return &PerformanceMetrics{
		AverageResponseTime: 245,
		MedianResponseTime:  200,
		P95ResponseTime:     500,
		P99ResponseTime:     1000,
		RequestsPerSecond:   50.5,
		ErrorRate:           0.025,
		ThroughputRPM:       30000,
		CacheHitRate:        0.75,
	}, nil
}

// GetDomainMetrics returns metrics for specific domains
func (s *Service) GetDomainMetrics(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*DomainMetrics, error) {
	// Mock implementation
	return []*DomainMetrics{
		{
			DomainName:          "commerce",
			RequestCount:        1500,
			SuccessRate:         0.95,
			AverageResponseTime: 300,
			UniqueUsers:         450,
			PopularQueries:      []string{"купить телефон", "заказать пиццу", "найди ресторан"},
		},
		{
			DomainName:          "recipes",
			RequestCount:        800,
			SuccessRate:         0.98,
			AverageResponseTime: 200,
			UniqueUsers:         280,
			PopularQueries:      []string{"рецепт борща", "как приготовить пасту", "быстрые рецепты"},
		},
		{
			DomainName:          "travel",
			RequestCount:        600,
			SuccessRate:         0.92,
			AverageResponseTime: 400,
			UniqueUsers:         200,
			PopularQueries:      []string{"отели в москве", "билеты в париж", "туры в египет"},
		},
	}, nil
}

// GetUserBehaviorAnalytics returns user behavior analytics
func (s *Service) GetUserBehaviorAnalytics(ctx context.Context, userID string, days int) (*UserBehaviorAnalytics, error) {
	// Mock implementation
	return &UserBehaviorAnalytics{
		UserID:               userID,
		TotalSessions:        25,
		TotalRequests:        150,
		AverageSessionTime:   420,
		MostUsedDomains:      []string{"commerce", "recipes", "travel"},
		PreferredLanguage:    "ru",
		LastActivity:         time.Now().Add(-2 * time.Hour),
		ConversionRate:       0.75,
		AverageResponseTime:  280,
		ErrorRate:            0.02,
	}, nil
}

// UserBehaviorAnalytics represents user behavior analytics
type UserBehaviorAnalytics struct {
	UserID              string    `json:"user_id"`
	TotalSessions       int64     `json:"total_sessions"`
	TotalRequests       int64     `json:"total_requests"`
	AverageSessionTime  int64     `json:"average_session_time_seconds"`
	MostUsedDomains     []string  `json:"most_used_domains"`
	PreferredLanguage   string    `json:"preferred_language"`
	LastActivity        time.Time `json:"last_activity"`
	ConversionRate      float64   `json:"conversion_rate"`
	AverageResponseTime int64     `json:"average_response_time_ms"`
	ErrorRate           float64   `json:"error_rate"`
}

// GetRealtimeMetrics returns real-time metrics
func (s *Service) GetRealtimeMetrics(ctx context.Context) (*RealtimeMetrics, error) {
	// Mock implementation - would use Redis/cache for real-time data
	return &RealtimeMetrics{
		ActiveUsers:        1250,
		RequestsPerSecond:  45.5,
		AverageLatency:     220,
		ErrorRate:          0.015,
		CacheHitRate:       0.78,
		QueueSize:          5,
		Timestamp:          time.Now(),
	}, nil
}

// RealtimeMetrics represents real-time system metrics
type RealtimeMetrics struct {
	ActiveUsers       int64     `json:"active_users"`
	RequestsPerSecond float64   `json:"requests_per_second"`
	AverageLatency    int64     `json:"average_latency_ms"`
	ErrorRate         float64    `json:"error_rate"`
	CacheHitRate      float64    `json:"cache_hit_rate"`
	QueueSize         int64     `json:"queue_size"`
	Timestamp         time.Time `json:"timestamp"`
}

// ExportAnalytics exports analytics data
func (s *Service) ExportAnalytics(ctx context.Context, tenantID string, format string, startDate, endDate time.Time) ([]byte, error) {
	// Mock implementation - would generate CSV/JSON export
	s.logger.Info("Exporting analytics data",
		zap.String("tenant_id", tenantID),
		zap.String("format", format),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate))

	return []byte(`{"message": "Analytics export would be generated here"}`), nil
}

// CleanOldData cleans old analytics data
func (s *Service) CleanOldData(ctx context.Context, daysToKeep int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)

	s.logger.Info("Cleaning old analytics data",
		zap.Int("days_to_keep", daysToKeep),
		zap.Time("cutoff_date", cutoffDate))

	// In production, this would delete old records from database
	return nil
}
