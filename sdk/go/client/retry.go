package client

import (
	"math"
	"net/http"
	"time"
)

// RetryConfig содержит конфигурацию для retry логики
type RetryConfig struct {
	MaxRetries      int           // Максимальное количество попыток (0 = без retry)
	InitialDelay    time.Duration // Начальная задержка
	MaxDelay        time.Duration // Максимальная задержка
	BackoffMultiplier float64     // Множитель для exponential backoff
	RetryableStatusCodes []int     // HTTP статусы, при которых нужно повторять запрос
}

// DefaultRetryConfig возвращает конфигурацию retry по умолчанию
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:        3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableStatusCodes: []int{
			http.StatusRequestTimeout,      // 408
			http.StatusTooManyRequests,     // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout,      // 504
		},
	}
}

// isRetryableError проверяет, можно ли повторить запрос при данной ошибке
func (c *Client) isRetryableError(err error, statusCode int) bool {
	// Проверяем статус код
	for _, code := range c.retryConfig.RetryableStatusCodes {
		if statusCode == code {
			return true
		}
	}

	// Проверяем на сетевые ошибки (timeout, connection refused и т.д.)
	if err != nil {
		// Сетевые ошибки обычно можно повторить
		return true
	}

	return false
}

// calculateBackoff вычисляет задержку для retry с exponential backoff
func (c *Client) calculateBackoff(attempt int) time.Duration {
	delay := float64(c.retryConfig.InitialDelay) * math.Pow(c.retryConfig.BackoffMultiplier, float64(attempt))
	
	if delay > float64(c.retryConfig.MaxDelay) {
		delay = float64(c.retryConfig.MaxDelay)
	}

	return time.Duration(delay)
}

// shouldRetry проверяет, нужно ли повторить запрос
// attempt - номер следующей попытки (1-based)
func (c *Client) shouldRetry(attempt int, err error, statusCode int) bool {
	if c.retryConfig.MaxRetries == 0 {
		return false
	}

	// attempt уже включает текущую попытку, поэтому проверяем > MaxRetries
	if attempt > c.retryConfig.MaxRetries {
		return false
	}

	return c.isRetryableError(err, statusCode)
}

