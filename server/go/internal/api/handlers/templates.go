package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nexus-protocol/server/internal/ai"
	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
)

// ExecuteTemplateHandler handles template execution requests
func ExecuteTemplateHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Parse request
		var req types.ExecuteTemplateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
				Details: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		// Validate request
		if err := validateExecuteTemplateRequest(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		// Create request metadata if not provided
		if req.Metadata == nil {
			req.Metadata = &types.RequestMetadata{
				RequestID:       uuid.New().String(),
				ProtocolVersion: "1.1.0",
				ClientVersion:   "1.0.0", // Would get from request headers
				Timestamp:       time.Now().Unix(),
			}
		}

		// Get user context from request
		userID := getUserIDFromContext(r.Context())

		// Log analytics event
		if err := services.(*Services).Analytics.LogTemplateExecution(r.Context(), userID, req.Metadata.RequestID, &req, nil); err != nil {
			services.(*Services).Logger.Warn("Failed to log analytics event", zap.Error(err))
		}

		// Execute template
		aiService := services.(*Services).AIService.(*ai.Service)
		result, err := aiService.ExecuteTemplate(r.Context(), &req)
		if err != nil {
			services.(*Services).Logger.Error("Template execution failed", zap.Error(err))

			// Log error analytics
			if logErr := services.(*Services).Analytics.LogError(r.Context(), userID, req.Metadata.RequestID, "TEMPLATE_EXECUTION_ERROR", err.Error(), nil); logErr != nil {
				services.(*Services).Logger.Warn("Failed to log error analytics", zap.Error(logErr))
			}

			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Template execution failed",
				Details: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		// Log successful execution
		if logErr := services.(*Services).Analytics.LogTemplateExecution(r.Context(), userID, req.Metadata.RequestID, &req, result); logErr != nil {
			services.(*Services).Logger.Warn("Failed to log execution analytics", zap.Error(logErr))
		}

		// Trigger webhook if configured
		if services.(*Services).Webhooks != nil {
			go func() {
				if err := services.(*Services).Webhooks.TriggerWebhook(r.Context(), "", "template.completed", result, userID, req.Metadata.RequestID); err != nil {
					services.(*Services).Logger.Warn("Failed to trigger webhook", zap.Error(err))
				}
			}()
		}

		// Create response metadata
		responseMetadata := &types.ResponseMetadata{
			RequestID:        req.Metadata.RequestID,
			ProtocolVersion:  "1.1.0",
			ServerVersion:    "1.1.0",
			Timestamp:        time.Now().Unix(),
			ProcessingTimeMS: int32(time.Since(startTime).Milliseconds()),
		}

		// Add enterprise metadata
		if services.(*Services).Config.RateLimit.Enabled {
			responseMetadata.RateLimitInfo = &types.RateLimitInfo{
				Limit:    int32(services.(*Services).Config.RateLimit.RequestsPerMin),
				Remaining: 950, // Would get from Redis
				ResetAt:  time.Now().Add(time.Minute).Unix(),
			}
		}

		if services.(*Services).Config.Cache.Enabled {
			responseMetadata.CacheInfo = &types.CacheInfo{
				CacheHit: false, // Would check cache
				CacheTTL:  int32(services.(*Services).Config.Cache.TTL.Seconds()),
			}
		}

		result.ResponseMetadata = responseMetadata

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":     result,
			"metadata": responseMetadata,
		})
	}
}

// GetTemplateStatusHandler handles template status requests
func GetTemplateStatusHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		executionID := vars["executionId"]

		if executionID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Execution ID is required",
			}, http.StatusBadRequest)
			return
		}

		// For now, return mock status (would check actual execution status)
		status := map[string]interface{}{
			"execution_id": executionID,
			"status":       "completed",
			"progress":     100,
			"message":      "Execution completed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

// StreamTemplateResultsHandler handles streaming template results
func StreamTemplateResultsHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		executionID := vars["executionId"]

		if executionID == "" {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Execution ID is required",
			}, http.StatusBadRequest)
			return
		}

		// Set headers for Server-Sent Events
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Get flusher
		flusher, ok := w.(http.Flusher)
		if !ok {
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Streaming not supported",
			}, http.StatusInternalServerError)
			return
		}

		// Simulate streaming results (in production, would stream from actual execution)
		results := []types.DomainSection{
			{
				DomainID:       "commerce",
				Title:          "Коммерческие предложения",
				Status:         "completed",
				ResponseTimeMS: 250,
				Results: []types.ResultItem{
					{
						ID:          uuid.New().String(),
						Type:        "product",
						Title:       "Рекомендуемый товар",
						Description: "На основе вашего запроса",
						Relevance:   0.95,
					},
				},
			},
		}

		for _, result := range results {
			data := map[string]interface{}{
				"type": "domain_result",
				"data": result,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
			time.Sleep(100 * time.Millisecond) // Simulate delay
		}

		// Send completion event
		completionData := map[string]interface{}{
			"type": "completed",
			"execution_id": executionID,
		}

		jsonData, _ := json.Marshal(completionData)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}
}

func validateExecuteTemplateRequest(req *types.ExecuteTemplateRequest) error {
	if req.Query == "" {
		return fmt.Errorf("query is required")
	}
	if len(req.Query) > 1000 {
		return fmt.Errorf("query too long (max 1000 characters)")
	}
	if req.Language != "ru" && req.Language != "en" {
		return fmt.Errorf("invalid language (must be 'ru' or 'en')")
	}
	return nil
}