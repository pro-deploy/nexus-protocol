package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRetry_NetworkError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Симулируем сетевую ошибку, закрывая соединение
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	cfg := Config{
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:      3,
			InitialDelay:    10 * time.Millisecond,
			MaxDelay:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
		},
	}

	client := NewClient(cfg)
	ctx := context.Background()

	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("Request failed after retries: %v", err)
	}
	defer resp.Body.Close()

	if attempts < 3 {
		t.Errorf("Expected at least 3 attempts, got %d", attempts)
	}
}

func TestRetry_ServerError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	cfg := Config{
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:      2,
			InitialDelay:    10 * time.Millisecond,
			MaxDelay:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
			RetryableStatusCodes: []int{http.StatusInternalServerError},
		},
	}

	client := NewClient(cfg)
	ctx := context.Background()

	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if attempts < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", attempts)
	}
}

func TestRetry_MaxRetriesExceeded(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := Config{
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:      2,
			InitialDelay:    10 * time.Millisecond,
			MaxDelay:        100 * time.Millisecond,
			BackoffMultiplier: 2.0,
			RetryableStatusCodes: []int{http.StatusInternalServerError},
		},
	}

	client := NewClient(cfg)
	ctx := context.Background()

	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err == nil && resp == nil {
		t.Error("Expected error or response after max retries, got nil")
	}

	// Проверяем, что был последний ответ с ошибкой
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", resp.StatusCode)
		}
	}

	expectedAttempts := cfg.RetryConfig.MaxRetries + 1
	if attempts != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attempts)
	}
}

func TestRateLimit_RetryAfter(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	cfg := Config{
		BaseURL: server.URL,
		RetryConfig: &RetryConfig{
			MaxRetries:      3,
			InitialDelay:    50 * time.Millisecond,
			MaxDelay:        500 * time.Millisecond,
			BackoffMultiplier: 2.0,
			RetryableStatusCodes: []int{http.StatusTooManyRequests},
		},
	}

	client := NewClient(cfg)
	ctx := context.Background()

	start := time.Now()
	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	duration := time.Since(start)
	defer resp.Body.Close()

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Должна быть задержка из-за Retry-After
	if duration < 1*time.Second {
		t.Errorf("Expected delay of at least 1 second, got %v", duration)
	}

	if attempts < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", attempts)
	}
}

func TestCalculateBackoff(t *testing.T) {
	cfg := RetryConfig{
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
	}

	client := &Client{retryConfig: cfg}

	// Первая попытка (attempt = 0)
	backoff := client.calculateBackoff(0)
	if backoff != 100*time.Millisecond {
		t.Errorf("Expected 100ms, got %v", backoff)
	}

	// Вторая попытка (attempt = 1)
	backoff = client.calculateBackoff(1)
	if backoff != 200*time.Millisecond {
		t.Errorf("Expected 200ms, got %v", backoff)
	}

	// Третья попытка (attempt = 2)
	backoff = client.calculateBackoff(2)
	if backoff != 400*time.Millisecond {
		t.Errorf("Expected 400ms, got %v", backoff)
	}

	// Проверяем ограничение MaxDelay
	cfg.MaxDelay = 300 * time.Millisecond
	client.retryConfig = cfg
	backoff = client.calculateBackoff(2)
	if backoff > 300*time.Millisecond {
		t.Errorf("Expected backoff <= 300ms, got %v", backoff)
	}
}

