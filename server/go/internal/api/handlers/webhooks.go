package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nexus-protocol/server/pkg/types"
)

// RegisterWebhookHandler handles webhook registration
func RegisterWebhookHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		var req struct {
			URL    string   `json:"url"`
			Events []string `json:"events"`
			Secret string   `json:"secret,omitempty"`
			Active *bool    `json:"active,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.URL == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Webhook URL is required",
			}, http.StatusBadRequest)
			return
		}

		if len(req.Events) == 0 {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "At least one event must be specified",
			}, http.StatusBadRequest)
			return
		}

		// Set defaults
		active := true
		if req.Active != nil {
			active = *req.Active
		}

		webhook, err := services.(*Services).Webhooks.RegisterWebhook(r.Context(), userID, req.URL, req.Events, req.Secret)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to register webhook",
				Details: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		// Set active status if provided
		if !active {
			webhook.Active = false
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhook": webhook,
			"message": "Webhook registered successfully",
		})
	}
}

// ListWebhooksHandler handles webhook listing
func ListWebhooksHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		webhooks, err := services.(*Services).Webhooks.ListWebhooks(r.Context(), userID)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to list webhooks",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhooks": webhooks,
			"count":    len(webhooks),
		})
	}
}

// UpdateWebhookHandler handles webhook updates
func UpdateWebhookHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		webhookID := vars["webhookId"]

		if webhookID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Webhook ID is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		webhook, err := services.(*Services).Webhooks.UpdateWebhook(r.Context(), webhookID, userID, updates)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to update webhook",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhook": webhook,
			"message": "Webhook updated successfully",
		})
	}
}

// DeleteWebhookHandler handles webhook deletion
func DeleteWebhookHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		webhookID := vars["webhookId"]

		if webhookID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Webhook ID is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		if err := services.(*Services).Webhooks.DeleteWebhook(r.Context(), webhookID, userID); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to delete webhook",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhook_id": webhookID,
			"message":    "Webhook deleted successfully",
		})
	}
}

// TestWebhookHandler handles webhook testing
func TestWebhookHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		webhookID := vars["webhookId"]

		if webhookID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Webhook ID is required",
			}, http.StatusBadRequest)
			return
		}

		if err := services.(*Services).Webhooks.TestWebhook(r.Context(), webhookID); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to test webhook",
				Details: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhook_id": webhookID,
			"message":    "Webhook test sent successfully",
		})
	}
}

// GetWebhookDeliveriesHandler handles webhook delivery history
func GetWebhookDeliveriesHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		webhookID := vars["webhookId"]

		if webhookID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Webhook ID is required",
			}, http.StatusBadRequest)
			return
		}

		// Parse pagination parameters
		limit := 50 // default
		offset := 0 // default

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}

		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			}
		}

		deliveries, err := services.(*Services).Webhooks.GetWebhookDeliveries(r.Context(), webhookID, limit, offset)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get webhook deliveries",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhook_id": webhookID,
			"deliveries": deliveries,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
				"count":  len(deliveries),
			},
		})
	}
}

// GetWebhookStatsHandler handles webhook statistics
func GetWebhookStatsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		stats, err := services.(*Services).Webhooks.GetWebhookStats(r.Context(), userID)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get webhook statistics",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}
