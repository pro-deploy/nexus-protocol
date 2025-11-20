package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/nexus-protocol/server/pkg/config"
	"go.uber.org/zap"
)

// NewClient creates a new Redis client
func NewClient(cfg *config.RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	logger.Info("Initializing Redis connection",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB))

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	if err := client.Ping(client.Context()).Err(); err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	logger.Info("Redis connection established")
	return client, nil
}
