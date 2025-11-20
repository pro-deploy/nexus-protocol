package api

import (
	"github.com/gorilla/mux"
	"github.com/nexus-protocol/server/internal/api/handlers"
	"github.com/nexus-protocol/server/pkg/middleware"
	"go.uber.org/zap"
)

// NewRouter creates the main HTTP router with all routes
func NewRouter(services *handlers.Services, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.Logging(logger))
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery(logger))

	// Rate limiting middleware
	if services.Config.RateLimit.Enabled {
		r.Use(middleware.RateLimit(services.Redis, services.Config.RateLimit))
	}

	// Authentication middleware for protected routes
	authMiddleware := middleware.Auth(services.Config, logger)

	// Health check endpoints (no auth required)
	r.HandleFunc("/health", handlers.HealthHandler(services)).Methods("GET")
	r.HandleFunc("/ready", handlers.ReadinessHandler(services)).Methods("GET")

	// Version endpoint
	r.HandleFunc("/api/v1/version", handlers.VersionHandler(services)).Methods("GET")

	// Authentication routes (no auth required)
	authRouter := r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/register", handlers.RegisterHandler(services)).Methods("POST")
	authRouter.HandleFunc("/login", handlers.LoginHandler(services)).Methods("POST")
	authRouter.HandleFunc("/refresh", handlers.RefreshTokenHandler(services)).Methods("POST")

	// Protected routes (require authentication)
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(authMiddleware)

	// User profile routes
	apiRouter.HandleFunc("/users/profile", handlers.GetUserProfileHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/users/profile", handlers.UpdateUserProfileHandler(services)).Methods("PUT")

	// Template execution routes
	apiRouter.HandleFunc("/templates/execute", handlers.ExecuteTemplateHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/templates/status/{executionId}", handlers.GetTemplateStatusHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/templates/stream/{executionId}", handlers.StreamTemplateResultsHandler(services)).Methods("GET")

	// Batch operations routes
	apiRouter.HandleFunc("/batch/execute", handlers.ExecuteBatchHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/batch/status/{batchId}", handlers.GetBatchStatusHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/batch/stats", handlers.GetBatchStatsHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/batch/{batchId}/cancel", handlers.CancelBatchHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/batch/{batchId}/operations", handlers.GetBatchOperationsHandler(services)).Methods("GET")

	// Webhook routes
	apiRouter.HandleFunc("/webhooks", handlers.RegisterWebhookHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/webhooks", handlers.ListWebhooksHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/webhooks/{webhookId}", handlers.UpdateWebhookHandler(services)).Methods("PUT")
	apiRouter.HandleFunc("/webhooks/{webhookId}", handlers.DeleteWebhookHandler(services)).Methods("DELETE")
	apiRouter.HandleFunc("/webhooks/{webhookId}/test", handlers.TestWebhookHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/webhooks/{webhookId}/deliveries", handlers.GetWebhookDeliveriesHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/webhooks/stats", handlers.GetWebhookStatsHandler(services)).Methods("GET")

	// Conversation routes
	apiRouter.HandleFunc("/conversations", handlers.CreateConversationHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/conversations", handlers.ListUserConversationsHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/conversations/{conversationId}", handlers.GetConversationHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/conversations/{conversationId}", handlers.ArchiveConversationHandler(services)).Methods("DELETE")
	apiRouter.HandleFunc("/conversations/{conversationId}/messages", handlers.SendMessageHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/conversations/{conversationId}/history", handlers.GetConversationHistoryHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/conversations/{conversationId}/typing", handlers.UpdateTypingStatusHandler(services)).Methods("POST")

	// Analytics routes
	apiRouter.HandleFunc("/analytics/events", handlers.LogAnalyticsEventHandler(services)).Methods("POST")
	apiRouter.HandleFunc("/analytics/stats", handlers.GetAnalyticsStatsHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/analytics/user", handlers.GetUserAnalyticsHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/analytics/realtime", handlers.GetRealtimeMetricsHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/analytics/export", handlers.ExportAnalyticsHandler(services)).Methods("GET")
	apiRouter.HandleFunc("/analytics/clean", handlers.CleanAnalyticsDataHandler(services)).Methods("POST")

	return r
}
