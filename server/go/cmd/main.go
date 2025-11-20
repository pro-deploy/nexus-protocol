package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	redisClient "github.com/go-redis/redis/v8"
	"github.com/nexus-protocol/server/internal/ai"
	"github.com/nexus-protocol/server/internal/analytics"
	"github.com/nexus-protocol/server/internal/api"
	"github.com/nexus-protocol/server/internal/api/handlers"
	"github.com/nexus-protocol/server/internal/auth"
	"github.com/nexus-protocol/server/internal/batch"
	"github.com/nexus-protocol/server/internal/conversation"
	"github.com/nexus-protocol/server/internal/db"
	"github.com/nexus-protocol/server/internal/redis"
	"github.com/nexus-protocol/server/internal/webhook"
	"github.com/nexus-protocol/server/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger := setupLogger()
	defer logger.Sync()

	logger.Info("🚀 Starting Nexus Protocol Server v1.1.0")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize services
	services, err := initializeServices(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize services", zap.Error(err))
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      api.NewRouter(services, logger),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start metrics server if enabled
	if cfg.Metrics.Enabled {
		go startMetricsServer(cfg, logger)
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", zap.Int("port", cfg.Server.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func setupLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		log.Fatal("Failed to create logger", err)
	}

	return logger
}

func initializeServices(cfg *config.Config, logger *zap.Logger) (*handlers.Services, error) {
	// Initialize database
	db, err := initializeDatabase(&cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis
	redisClient, err := initializeRedis(&cfg.Redis, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize AI service
	aiService, err := initializeAIService(&cfg.AI, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI service: %w", err)
	}

	// Initialize services
	authService := auth.NewService(cfg, logger)
	batchService := batch.NewService(aiService, logger, cfg.RateLimit.BurstSize, 100)
	webhookService := webhook.NewService(logger)
	analyticsService := analytics.NewService(logger)
	conversationService := conversation.NewService(logger)

	// Create services struct
	services := &handlers.Services{
		Config:        cfg,
		Logger:        logger,
		Database:      db,
		Redis:         redisClient,
		AIService:     aiService,
		Auth:          authService,
		Templates:     aiService, // Templates service is part of AI service
		Batch:         batchService,
		Webhooks:      webhookService,
		Analytics:     analyticsService,
		Conversations: conversationService,
	}

	return services, nil
}

func initializeDatabase(cfg *config.DatabaseConfig, logger *zap.Logger) (*db.DB, error) {
	return db.NewDB(cfg, logger)
}

func initializeRedis(cfg *config.RedisConfig, logger *zap.Logger) (*redisClient.Client, error) {
	return redis.NewClient(cfg, logger)
}

func initializeAIService(cfg *config.AIConfig, logger *zap.Logger) (*ai.Service, error) {
	return ai.NewService(cfg, logger), nil
}

func startMetricsServer(cfg *config.Config, logger *zap.Logger) {
	metricsMux := http.NewServeMux()

	// Add Prometheus metrics endpoint
	metricsMux.Handle("/metrics", promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
		Handler: metricsMux,
	}

	logger.Info("Metrics server starting", zap.Int("port", cfg.Metrics.Port))
	if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Metrics server failed", zap.Error(err))
	}
}
