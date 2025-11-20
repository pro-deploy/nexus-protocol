package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nexus-protocol/server/pkg/types"
)

// ExecuteBatchHandler handles batch execution requests
func ExecuteBatchHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var req struct {
			Requests []*types.ExecuteTemplateRequest `json:"requests"`
			Options  *struct {
				Parallel bool `json:"parallel,omitempty"`
			} `json:"options,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		// Validate batch request
		if err := services.(*Services).Batch.ValidateBatchRequest(req.Requests); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		// Determine execution mode
		parallel := false
		if req.Options != nil {
			parallel = req.Options.Parallel
		}

		// Execute batch
		job, err := services.(*Services).Batch.ExecuteBatch(r.Context(), req.Requests, parallel)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Batch execution failed",
				Details: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"batch_id": job.ID,
			"status":   job.Status,
			"operations_count": len(job.Operations),
			"message": "Batch execution started",
		})
	}
}

// GetBatchStatusHandler handles batch status requests
func GetBatchStatusHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		batchID := vars["batchId"]

		if batchID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Batch ID is required",
			}, http.StatusBadRequest)
			return
		}

		// Get batch status
		job, err := services.(*Services).Batch.GetBatchStatus(r.Context(), batchID)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "NOT_FOUND",
				Type:    "NOT_FOUND",
				Message: "Batch not found",
			}, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(job)
	}
}

// GetBatchStatsHandler handles batch statistics requests
func GetBatchStatsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		// Check permissions (enterprise feature)
		if err := services.(*Services).Auth.ValidateUserPermissions(userID, []string{"admin"}, "batch", "read_stats"); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHORIZATION_FAILED",
				Type:    "AUTHORIZATION_ERROR",
				Message: "Insufficient permissions to view batch statistics",
			}, http.StatusForbidden)
			return
		}

		stats, err := services.(*Services).Batch.GetBatchStats(r.Context())
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get batch statistics",
			}, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}

// CancelBatchHandler handles batch cancellation requests
func CancelBatchHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		batchID := vars["batchId"]

		if batchID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Batch ID is required",
			}, http.StatusBadRequest)
			return
		}

		userID := getUserIDFromContext(r.Context())

		// Check permissions
		if err := services.(*Services).Auth.ValidateUserPermissions(userID, []string{"admin"}, "batch", "cancel"); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHORIZATION_FAILED",
				Type:    "AUTHORIZATION_ERROR",
				Message: "Insufficient permissions to cancel batch",
			}, http.StatusForbidden)
			return
		}

		if err := services.(*Services).Batch.CancelBatch(r.Context(), batchID); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to cancel batch",
			}, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"batch_id": batchID,
			"status":   "cancelled",
			"message":  "Batch cancellation requested",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// GetBatchOperationsHandler handles requests for batch operations
func GetBatchOperationsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		batchID := vars["batchId"]

		if batchID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Batch ID is required",
			}, http.StatusBadRequest)
			return
		}

		// Parse pagination parameters
		limit := 50  // default
		offset := 0  // default

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

		operations, err := services.(*Services).Batch.GetBatchOperations(r.Context(), batchID, limit, offset)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Failed to get batch operations",
			}, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"batch_id":   batchID,
			"operations": operations,
			"pagination": map[string]interface{}{
				"limit":  limit,
				"offset": offset,
				"count":  len(operations),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
