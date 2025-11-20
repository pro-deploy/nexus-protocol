package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nexus-protocol/server/internal/analytics"
	"github.com/nexus-protocol/server/internal/auth"
	"github.com/nexus-protocol/server/internal/batch"
	"github.com/nexus-protocol/server/internal/conversation"
	"github.com/nexus-protocol/server/internal/webhook"
	"github.com/nexus-protocol/server/pkg/config"
	"github.com/nexus-protocol/server/pkg/types"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// Services contains all service dependencies
type Services struct {
	Config       *config.Config
	Logger       *zap.Logger
	Database     interface{} // *db.DB
	Redis        *redis.Client
	AIService    interface{} // *ai.Service
	Auth         *auth.Service
	Templates    interface{} // *ai.Service (same as AI)
	Batch        *batch.Service
	Webhooks     *webhook.Service
	Analytics    *analytics.Service
	Conversations *conversation.Service
}

// sendError sends a standardized error response
func sendError(w http.ResponseWriter, err types.ErrorDetail, statusCode int) {
	response := types.ErrorResponse{
		Error: err,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Error-Code", err.Code)
	w.Header().Set("X-Error-Type", err.Type)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// getUserIDFromContext extracts user ID from JWT token in context
func getUserIDFromContext(_ context.Context) string {
	// This would extract user ID from JWT token stored in context by auth middleware
	// For now, return a placeholder
	return "user-123"
}

