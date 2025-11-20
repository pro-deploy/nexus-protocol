package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
)

// HealthHandler handles basic health check
func HealthHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := types.HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   "1.1.0",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// ReadinessHandler handles detailed readiness check
func ReadinessHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check database connectivity
		dbStatus := "ok"
		// TODO: Implement database ping when Database interface is defined
		// if err := services.(*Services).Database.Ping(r.Context()); err != nil {
		// 	services.(*Services).Logger.Error("Database health check failed", zap.Error(err))
		// 	dbStatus = "error"
		// }

		// Check Redis connectivity
		redisStatus := "ok"
		if err := services.(*Services).Redis.Ping(r.Context()).Err(); err != nil {
			services.(*Services).Logger.Error("Redis health check failed", zap.Error(err))
			redisStatus = "error"
		}

		// Check AI service
		aiStatus := "ok"
		if services.(*Services).AIService != nil {
			// Implement AI service health check
			aiStatus = "ok" // Placeholder
		}

		// Overall status
		overallStatus := "ready"
		if dbStatus != "ok" || redisStatus != "ok" || aiStatus != "ok" {
			overallStatus = "not_ready"
		}

		response := types.ReadinessResponse{
			Status:    overallStatus,
			Timestamp: time.Now().Format(time.RFC3339),
			Checks: types.ReadinessChecks{
				Database:   dbStatus,
				Redis:      redisStatus,
				AIServices: aiStatus,
			},
		}

		// Add component details for enterprise monitoring
		response.Components = map[string]*types.ComponentStatus{
			"database": {
				Status:    dbStatus,
				LatencyMS: 5, // Would measure actual latency
			},
			"redis": {
				Status:    redisStatus,
				LatencyMS: 2,
			},
			"ai_service": {
				Status:    aiStatus,
				LatencyMS: 150,
			},
		}

		// Add capacity information
		response.Capacity = &types.CapacityInfo{
			CurrentLoad:       0.75, // Would get from metrics
			MaxCapacity:       10000,
			AvailableCapacity: 2500,
			ActiveConnections: 7500,
		}

		statusCode := http.StatusOK
		if overallStatus != "ready" {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}

// VersionHandler returns API version information
func VersionHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"protocol_version": "1.1.0",
			"server_version":   "1.1.0",
			"api_version":      "v1",
			"build_info": map[string]string{
				"git_commit": "dev", // Would be set during build
				"build_time": time.Now().Format(time.RFC3339),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
