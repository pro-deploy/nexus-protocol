package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nexus-protocol/server/internal/analytics"
	"github.com/nexus-protocol/server/pkg/types"
)

// LogAnalyticsEventHandler handles analytics event logging
func LogAnalyticsEventHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		var req struct {
			EventType string                 `json:"event_type"`
			EventData map[string]interface{} `json:"event_data,omitempty"`
			RequestID string                 `json:"request_id,omitempty"`
			SessionID string                 `json:"session_id,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		if req.EventType == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Event type is required",
			}, http.StatusBadRequest)
			return
		}

		event := &analytics.Event{
			ID:        uuid.New().String(),
			UserID:    userID,
			EventType: req.EventType,
			EventData: req.EventData,
			Timestamp: time.Now(),
			UserAgent: r.Header.Get("User-Agent"),
			IPAddress: getClientIP(r),
		}

		if err := services.(*Services).Analytics.LogEvent(r.Context(), event); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to log analytics event",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"event_id": event.ID,
			"message":  "Analytics event logged successfully",
		})
	}
}

// GetAnalyticsStatsHandler handles analytics statistics requests
func GetAnalyticsStatsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Check permissions for analytics access
		if err := services.(*Services).Auth.ValidateUserPermissions(userID, []string{"admin"}, "analytics", "read"); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHORIZATION_FAILED",
				Type:    "AUTHORIZATION_ERROR",
				Message: "Insufficient permissions to access analytics",
			}, http.StatusForbidden)
			return
		}

		// Parse date range
		startDate := time.Now().AddDate(0, 0, -30) // Last 30 days by default
		endDate := time.Now()

		if startStr := r.URL.Query().Get("start_date"); startStr != "" {
			if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
				startDate = parsed
			}
		}

		if endStr := r.URL.Query().Get("end_date"); endStr != "" {
			if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
				endDate = parsed
			}
		}

		// Get conversion metrics
		conversionMetrics, err := services.(*Services).Analytics.GetConversionMetrics(r.Context(), "", startDate, endDate)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get conversion metrics",
			}, http.StatusInternalServerError)
			return
		}

		// Get performance metrics
		performanceMetrics, err := services.(*Services).Analytics.GetPerformanceMetrics(r.Context(), "", startDate, endDate)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get performance metrics",
			}, http.StatusInternalServerError)
			return
		}

		// Get domain metrics
		domainMetrics, err := services.(*Services).Analytics.GetDomainMetrics(r.Context(), "", startDate, endDate)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get domain metrics",
			}, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"period": map[string]interface{}{
				"start_date": startDate.Format("2006-01-02"),
				"end_date":   endDate.Format("2006-01-02"),
			},
			"conversion_metrics":  conversionMetrics,
			"performance_metrics": performanceMetrics,
			"domain_metrics":      domainMetrics,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// GetUserAnalyticsHandler handles user-specific analytics
func GetUserAnalyticsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Parse days parameter
		days := 30 // default
		if daysStr := r.URL.Query().Get("days"); daysStr != "" {
			if parsed, err := strconv.Atoi(daysStr); err == nil && parsed > 0 && parsed <= 365 {
				days = parsed
			}
		}

		userAnalytics, err := services.(*Services).Analytics.GetUserBehaviorAnalytics(r.Context(), userID, days)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get user analytics",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userAnalytics)
	}
}

// GetRealtimeMetricsHandler handles real-time metrics requests
func GetRealtimeMetricsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Check permissions for real-time metrics
		if err := services.(*Services).Auth.ValidateUserPermissions(userID, []string{"admin"}, "analytics", "realtime"); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHORIZATION_FAILED",
				Type:    "AUTHORIZATION_ERROR",
				Message: "Insufficient permissions to access real-time metrics",
			}, http.StatusForbidden)
			return
		}

		metrics, err := services.(*Services).Analytics.GetRealtimeMetrics(r.Context())
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get real-time metrics",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metrics)
	}
}

// ExportAnalyticsHandler handles analytics data export
func ExportAnalyticsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Check permissions for export
		if err := services.(*Services).Auth.ValidateUserPermissions(userID, []string{"admin"}, "analytics", "export"); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHORIZATION_FAILED",
				Type:    "AUTHORIZATION_ERROR",
				Message: "Insufficient permissions to export analytics",
			}, http.StatusForbidden)
			return
		}

		// Parse parameters
		exportFormat := r.URL.Query().Get("format")
		if exportFormat == "" {
			exportFormat = "json"
		}

		startDate := time.Now().AddDate(0, 0, -30)
		endDate := time.Now()

		if startStr := r.URL.Query().Get("start_date"); startStr != "" {
			if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
				startDate = parsed
			}
		}

		if endStr := r.URL.Query().Get("end_date"); endStr != "" {
			if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
				endDate = parsed
			}
		}

		data, err := services.(*Services).Analytics.ExportAnalytics(r.Context(), "", exportFormat, startDate, endDate)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to export analytics data",
			}, http.StatusInternalServerError)
			return
		}

		// Set appropriate headers based on format
		switch exportFormat {
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment; filename=analytics.csv")
		case "json":
			fallthrough
		default:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", "attachment; filename=analytics.json")
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// CleanAnalyticsDataHandler handles analytics data cleanup
func CleanAnalyticsDataHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Only admins can clean data
		if err := services.(*Services).Auth.ValidateUserPermissions(userID, []string{"admin"}, "analytics", "clean"); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHORIZATION_FAILED",
				Type:    "AUTHORIZATION_ERROR",
				Message: "Insufficient permissions to clean analytics data",
			}, http.StatusForbidden)
			return
		}

		daysToKeep := 90 // default
		if daysStr := r.URL.Query().Get("days_to_keep"); daysStr != "" {
			if parsed, err := strconv.Atoi(daysStr); err == nil && parsed > 0 {
				daysToKeep = parsed
			}
		}

		if err := services.(*Services).Analytics.CleanOldData(r.Context(), daysToKeep); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to clean analytics data",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"days_kept": daysToKeep,
			"message":   "Analytics data cleaned successfully",
		})
	}
}

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP if multiple
		if idx := strings.Index(xForwardedFor, ","); idx > 0 {
			return strings.TrimSpace(xForwardedFor[:idx])
		}
		return strings.TrimSpace(xForwardedFor)
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// Fall back to RemoteAddr
	return strings.TrimSpace(r.RemoteAddr)
}
