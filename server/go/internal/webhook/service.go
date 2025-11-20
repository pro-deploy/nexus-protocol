package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service handles webhook operations
type Service struct {
	logger         *zap.Logger
	httpClient     *http.Client
	maxRetries     int
	defaultTimeout time.Duration
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID          string            `json:"id" db:"id"`
	UserID      string            `json:"user_id" db:"user_id"`
	URL         string            `json:"url" db:"url"`
	Events      []string          `json:"events" db:"events"`
	Secret      string            `json:"secret,omitempty" db:"secret"`
	Active      bool              `json:"active" db:"active"`
	Headers     map[string]string `json:"headers,omitempty" db:"headers"`
	RetryPolicy *RetryPolicy      `json:"retry_policy,omitempty" db:"retry_policy"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

// RetryPolicy defines retry behavior for failed webhooks
type RetryPolicy struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxDelay      time.Duration `json:"max_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID         string    `json:"id" db:"id"`
	WebhookID  string    `json:"webhook_id" db:"webhook_id"`
	EventType  string    `json:"event_type" db:"event_type"`
	Payload    string    `json:"payload" db:"payload"`
	Status     string    `json:"status" db:"status"` // pending, success, failed, retry
	StatusCode int       `json:"status_code,omitempty" db:"status_code"`
	Error      string    `json:"error,omitempty" db:"error"`
	Attempt    int       `json:"attempt" db:"attempt"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// EventPayload represents the payload sent to webhooks
type EventPayload struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
	UserID    string      `json:"user_id,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// NewService creates a new webhook service
func NewService(logger *zap.Logger) *Service {
	return &Service{
		logger:         logger,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
		maxRetries:     3,
		defaultTimeout: 30 * time.Second,
	}
}

// RegisterWebhook registers a new webhook
func (s *Service) RegisterWebhook(ctx context.Context, userID string, url string, events []string, secret string) (*Webhook, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook URL cannot be empty")
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("webhook must have at least one event")
	}

	if secret == "" {
		secret = s.generateSecret()
	}

	webhook := &Webhook{
		ID:          uuid.New().String(),
		UserID:      userID,
		URL:         url,
		Events:      events,
		Secret:      secret,
		Active:      true,
		Headers:     make(map[string]string),
		RetryPolicy: &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  1 * time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// In a real implementation, save to database
	s.logger.Info("Webhook registered",
		zap.String("webhook_id", webhook.ID),
		zap.String("user_id", userID),
		zap.String("url", url),
		zap.Strings("events", events))

	return webhook, nil
}

// UpdateWebhook updates an existing webhook
func (s *Service) UpdateWebhook(ctx context.Context, webhookID, userID string, updates map[string]interface{}) (*Webhook, error) {
	// In a real implementation, retrieve and update from database
	webhook := &Webhook{
		ID:     webhookID,
		UserID: userID,
		URL:    "https://example.com/webhook", // Mock
		Events: []string{"template.completed"},
		Active: true,
	}

	// Apply updates
	if url, ok := updates["url"].(string); ok {
		webhook.URL = url
	}
	if events, ok := updates["events"].([]string); ok {
		webhook.Events = events
	}
	if active, ok := updates["active"].(bool); ok {
		webhook.Active = active
	}

	webhook.UpdatedAt = time.Now()

	s.logger.Info("Webhook updated", zap.String("webhook_id", webhookID))
	return webhook, nil
}

// DeleteWebhook deletes a webhook
func (s *Service) DeleteWebhook(ctx context.Context, webhookID, userID string) error {
	// In a real implementation, delete from database
	s.logger.Info("Webhook deleted",
		zap.String("webhook_id", webhookID),
		zap.String("user_id", userID))
	return nil
}

// ListWebhooks returns webhooks for a user
func (s *Service) ListWebhooks(ctx context.Context, userID string) ([]*Webhook, error) {
	// Mock implementation - in real system would query database
	webhooks := []*Webhook{
		{
			ID:        uuid.New().String(),
			UserID:    userID,
			URL:       "https://app.example.com/webhooks",
			Events:    []string{"template.completed", "batch.finished"},
			Active:    true,
			CreatedAt: time.Now().Add(-time.Hour),
			UpdatedAt: time.Now().Add(-time.Hour),
		},
	}

	return webhooks, nil
}

// TriggerWebhook triggers a webhook for an event
func (s *Service) TriggerWebhook(ctx context.Context, webhookID string, eventType string, data interface{}, userID string, requestID string) error {
	// Get webhook configuration
	webhook, err := s.getWebhook(ctx, webhookID)
	if err != nil {
		return fmt.Errorf("failed to get webhook: %w", err)
	}

	if !webhook.Active {
		return fmt.Errorf("webhook is not active")
	}

	// Check if webhook handles this event type
	if !s.webhookHandlesEvent(webhook, eventType) {
		return nil // Silently ignore events that webhook doesn't handle
	}

	// Create event payload
	payload := EventPayload{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().Unix(),
		Data:      data,
		UserID:    userID,
		RequestID: requestID,
	}

	// Deliver webhook asynchronously
	go s.deliverWebhook(webhook, payload)

	return nil
}

// deliverWebhook delivers a webhook with retry logic
func (s *Service) deliverWebhook(webhook *Webhook, payload EventPayload) {
	deliveryID := uuid.New().String()

	s.logger.Info("Delivering webhook",
		zap.String("delivery_id", deliveryID),
		zap.String("webhook_id", webhook.ID),
		zap.String("event_type", payload.EventType),
		zap.String("url", webhook.URL))

	// Create delivery record
	delivery := &WebhookDelivery{
		ID:        deliveryID,
		WebhookID: webhook.ID,
		EventType: payload.EventType,
		Payload:   s.payloadToString(payload),
		Status:    "pending",
		Attempt:   1,
		CreatedAt: time.Now(),
	}

	maxRetries := s.maxRetries
	if webhook.RetryPolicy != nil {
		maxRetries = webhook.RetryPolicy.MaxRetries
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		delivery.Attempt = attempt

		err := s.attemptDelivery(webhook, payload, delivery)
		if err == nil {
			delivery.Status = "success"
			delivery.DeliveredAt = &time.Time{}
			*delivery.DeliveredAt = time.Now()
			break
		}

		lastErr = err
		delivery.Error = err.Error()
		delivery.Status = "retry"

		s.logger.Warn("Webhook delivery failed, retrying",
			zap.String("delivery_id", deliveryID),
			zap.Int("attempt", attempt),
			zap.Error(err))

		if attempt < maxRetries {
			delay := s.calculateDelay(attempt, webhook.RetryPolicy)
			time.Sleep(delay)
		}
	}

	if delivery.Status != "success" {
		delivery.Status = "failed"
		s.logger.Error("Webhook delivery failed permanently",
			zap.String("delivery_id", deliveryID),
			zap.Error(lastErr))
	}

	// Save delivery record (in real implementation)
	_ = delivery
}

// attemptDelivery attempts to deliver a webhook once
func (s *Service) attemptDelivery(webhook *Webhook, payload EventPayload, delivery *WebhookDelivery) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Nexus-Webhook/1.0")
	req.Header.Set("X-Webhook-ID", webhook.ID)
	req.Header.Set("X-Webhook-Event", payload.EventType)
	req.Header.Set("X-Webhook-Delivery", delivery.ID)

	// Add custom headers
	for key, value := range webhook.Headers {
		req.Header.Set(key, value)
	}

	// Add signature if secret is configured
	if webhook.Secret != "" {
		signature := s.generateSignature(jsonPayload, webhook.Secret)
		req.Header.Set("X-Hub-Signature-256", "sha256="+signature)
	}

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	delivery.StatusCode = resp.StatusCode

	// Consider 2xx status codes as success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("webhook returned status %d", resp.StatusCode)
}

// GetWebhookDeliveries returns delivery history for a webhook
func (s *Service) GetWebhookDeliveries(ctx context.Context, webhookID string, limit, offset int) ([]*WebhookDelivery, error) {
	// Mock implementation
	deliveries := []*WebhookDelivery{
		{
			ID:         uuid.New().String(),
			WebhookID:  webhookID,
			EventType:  "template.completed",
			Status:     "success",
			StatusCode: 200,
			Attempt:    1,
			DeliveredAt: &time.Time{},
			CreatedAt:  time.Now().Add(-time.Minute),
		},
	}

	return deliveries, nil
}

// TestWebhook tests a webhook by sending a test event
func (s *Service) TestWebhook(ctx context.Context, webhookID string) error {
	testPayload := EventPayload{
		EventID:   uuid.New().String(),
		EventType: "webhook.test",
		Timestamp: time.Now().Unix(),
		Data: map[string]string{
			"message": "This is a test webhook delivery",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	return s.TriggerWebhook(ctx, webhookID, "webhook.test", testPayload.Data, "test-user", "test-request")
}

// GetWebhookStats returns statistics for webhooks
func (s *Service) GetWebhookStats(ctx context.Context, userID string) (*WebhookStats, error) {
	// Mock statistics
	return &WebhookStats{
		TotalWebhooks:      5,
		ActiveWebhooks:     4,
		TotalDeliveries:    1250,
		SuccessfulDeliveries: 1180,
		FailedDeliveries:   70,
		AverageResponseTime: 250,
	}, nil
}

// WebhookStats represents webhook statistics
type WebhookStats struct {
	TotalWebhooks        int64 `json:"total_webhooks"`
	ActiveWebhooks       int64 `json:"active_webhooks"`
	TotalDeliveries      int64 `json:"total_deliveries"`
	SuccessfulDeliveries int64 `json:"successful_deliveries"`
	FailedDeliveries     int64 `json:"failed_deliveries"`
	AverageResponseTime  int64 `json:"average_response_time_ms"`
}

// Helper functions

func (s *Service) getWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	// Mock implementation - in real system would query database
	return &Webhook{
		ID:     webhookID,
		URL:    "https://app.example.com/webhooks",
		Events: []string{"template.completed", "batch.finished", "webhook.test"},
		Active: true,
		Secret: "webhook-secret-123",
		RetryPolicy: &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  1 * time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
		},
	}, nil
}

func (s *Service) webhookHandlesEvent(webhook *Webhook, eventType string) bool {
	for _, event := range webhook.Events {
		if event == eventType {
			return true
		}
	}
	return false
}

func (s *Service) payloadToString(payload EventPayload) string {
	data, _ := json.Marshal(payload)
	return string(data)
}

func (s *Service) generateSecret() string {
	return uuid.New().String() + uuid.New().String()
}

func (s *Service) generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Service) calculateDelay(attempt int, policy *RetryPolicy) time.Duration {
	if policy == nil {
		// Default exponential backoff
		delay := time.Duration(attempt) * time.Second
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}
		return delay
	}

	delay := float64(policy.InitialDelay) * pow(policy.BackoffFactor, attempt-1)
	if delay > float64(policy.MaxDelay) {
		delay = float64(policy.MaxDelay)
	}

	return time.Duration(delay)
}

func pow(base float64, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
