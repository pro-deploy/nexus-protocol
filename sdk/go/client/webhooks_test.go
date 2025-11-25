package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func TestRegisterWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/webhooks" {
			t.Errorf("Expected POST /api/v1/webhooks, got %s %s", r.Method, r.URL.Path)
		}

		var req types.RegisterWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Проверяем содержимое запроса
		if req.Config == nil {
			t.Error("Expected webhook config")
		}

		if req.Config.URL != "https://example.com/webhook" {
			t.Errorf("Expected URL 'https://example.com/webhook', got %s", req.Config.URL)
		}

		if len(req.Config.Events) != 2 {
			t.Errorf("Expected 2 events, got %d", len(req.Config.Events))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := types.RegisterWebhookResponse{
			WebhookID: "webhook-123",
			Status:    "registered",
			Message:   "Webhook registered successfully",
			ResponseMetadata: &types.ResponseMetadata{
				RequestID:      "req-123",
				ProtocolVersion: "2.0.0",
				ServerVersion:   "2.0.0",
				Timestamp:       1640995200,
			},
		}
		json.NewEncoder(w).Encode(map[string]types.RegisterWebhookResponse{"data": response})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})

	req := &types.RegisterWebhookRequest{
		Config: &types.WebhookConfig{
			URL:    "https://example.com/webhook",
			Events: []string{"template.completed", "template.failed"},
			Secret: "webhook-secret",
			RetryPolicy: &types.WebhookRetryPolicy{
				MaxRetries: 3,
				InitialDelay: 1000,
			},
		},
	}

	result, err := client.RegisterWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("RegisterWebhook failed: %v", err)
	}

	if result.WebhookID != "webhook-123" {
		t.Errorf("Expected webhook ID 'webhook-123', got %s", result.WebhookID)
	}

	if result.Status != "registered" {
		t.Errorf("Expected status 'registered', got %s", result.Status)
	}
}

func TestListWebhooks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/api/v1/webhooks" {
			t.Errorf("Expected GET /api/v1/webhooks, got %s %s", r.Method, r.URL.Path)
		}

		// Проверяем query параметры
		if r.URL.Query().Get("active_only") != "true" {
			t.Error("Expected active_only=true query parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := types.ListWebhooksResponse{
			Webhooks: []types.WebhookInfo{
				{
					ID: "webhook-1",
					Config: &types.WebhookConfig{
						URL:    "https://app1.com/webhook",
						Events: []string{"template.completed"},
						Active: true,
					},
					CreatedAt: 1640995200,
					SuccessCount: 100,
					ErrorCount:   5,
				},
				{
					ID: "webhook-2",
					Config: &types.WebhookConfig{
						URL:    "https://app2.com/webhook",
						Events: []string{"batch.completed"},
						Active: true,
					},
					CreatedAt: 1640995300,
					SuccessCount: 50,
					ErrorCount:   2,
				},
			},
			Total:  2,
			Limit:  10,
			Offset: 0,
			ResponseMetadata: &types.ResponseMetadata{
				RequestID:      "req-456",
				ProtocolVersion: "2.0.0",
				ServerVersion:   "2.0.0",
				Timestamp:       1640995400,
			},
		}
		json.NewEncoder(w).Encode(map[string]types.ListWebhooksResponse{"data": response})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})

	req := &types.ListWebhooksRequest{
		ActiveOnly: true,
		Limit:      10,
	}

	result, err := client.ListWebhooks(context.Background(), req)
	if err != nil {
		t.Fatalf("ListWebhooks failed: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}

	if len(result.Webhooks) != 2 {
		t.Errorf("Expected 2 webhooks, got %d", len(result.Webhooks))
	}

	if result.Webhooks[0].ID != "webhook-1" {
		t.Errorf("Expected first webhook ID 'webhook-1', got %s", result.Webhooks[0].ID)
	}

	if result.Webhooks[1].SuccessCount != 50 {
		t.Errorf("Expected second webhook success count 50, got %d", result.Webhooks[1].SuccessCount)
	}
}

func TestDeleteWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		expectedPath := "/api/v1/webhooks/webhook-123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := types.DeleteWebhookResponse{
			WebhookID: "webhook-123",
			Status:    "deleted",
			Message:   "Webhook deleted successfully",
			ResponseMetadata: &types.ResponseMetadata{
				RequestID:      "req-789",
				ProtocolVersion: "2.0.0",
				ServerVersion:   "2.0.0",
				Timestamp:       1640995500,
			},
		}
		json.NewEncoder(w).Encode(map[string]types.DeleteWebhookResponse{"data": response})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})

	result, err := client.DeleteWebhook(context.Background(), "webhook-123")
	if err != nil {
		t.Fatalf("DeleteWebhook failed: %v", err)
	}

	if result.WebhookID != "webhook-123" {
		t.Errorf("Expected webhook ID 'webhook-123', got %s", result.WebhookID)
	}

	if result.Status != "deleted" {
		t.Errorf("Expected status 'deleted', got %s", result.Status)
	}
}

func TestTestWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		expectedPath := "/api/v1/webhooks/webhook-123/test"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := types.TestWebhookResponse{
			WebhookID:   "webhook-123",
			Status:      "sent",
			ResponseCode: 200,
			ResponseTimeMS: 150,
			ResponseMetadata: &types.ResponseMetadata{
				RequestID:      "req-999",
				ProtocolVersion: "2.0.0",
				ServerVersion:   "2.0.0",
				Timestamp:       1640995600,
			},
		}
		json.NewEncoder(w).Encode(map[string]types.TestWebhookResponse{"data": response})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})

	req := &types.TestWebhookRequest{
		WebhookID: "webhook-123",
		Event:     "template.completed",
		Data:      map[string]interface{}{"test": true},
	}

	result, err := client.TestWebhook(context.Background(), req)
	if err != nil {
		t.Fatalf("TestWebhook failed: %v", err)
	}

	if result.WebhookID != "webhook-123" {
		t.Errorf("Expected webhook ID 'webhook-123', got %s", result.WebhookID)
	}

	if result.Status != "sent" {
		t.Errorf("Expected status 'sent', got %s", result.Status)
	}

	if result.ResponseCode != 200 {
		t.Errorf("Expected response code 200, got %d", result.ResponseCode)
	}

	if result.ResponseTimeMS != 150 {
		t.Errorf("Expected response time 150ms, got %d", result.ResponseTimeMS)
	}
}
