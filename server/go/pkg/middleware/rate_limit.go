package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nexus-protocol/server/pkg/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RateLimit middleware implements rate limiting
func RateLimit(redisClient *redis.Client, cfg config.RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client identifier (IP address or user ID)
			clientID := getClientIdentifier(r)

			// Check rate limit
			key := "rate_limit:" + clientID
			count, err := redisClient.Incr(r.Context(), key).Result()
			if err != nil {
				// If Redis fails, allow request but log error
				zap.L().Error("Redis rate limit check failed", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			// Set expiration on first request
			if count == 1 {
				redisClient.Expire(r.Context(), key, time.Minute)
			}

			// Check if limit exceeded
			if count > int64(cfg.RequestsPerMin) {
				// Get reset time
				ttl, _ := redisClient.TTL(r.Context(), key).Result()
				resetAt := time.Now().Add(ttl).Unix()

				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.RequestsPerMin))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt))
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))

				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add rate limit headers
			remaining := cfg.RequestsPerMin - int(count)
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.RequestsPerMin))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIdentifier(r *http.Request) string {
	// Try to get user ID from context (authenticated requests)
	if userID := r.Context().Value("user_id"); userID != nil {
		return userID.(string)
	}

	// Fall back to IP address
	return r.RemoteAddr
}
