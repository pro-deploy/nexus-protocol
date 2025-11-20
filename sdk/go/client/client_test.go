package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nexus-protocol/go-sdk/types"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		BaseURL: "http://localhost:8080",
		Token:   "test-token",
	}

	client := NewClient(cfg)

	if client.baseURL != "http://localhost:8080" {
		t.Errorf("Expected baseURL %s, got %s", "http://localhost:8080", client.baseURL)
	}

	if client.token != "test-token" {
		t.Errorf("Expected token %s, got %s", "test-token", client.token)
	}

	if client.protocolVersion != DefaultProtocolVersion {
		t.Errorf("Expected protocolVersion %s, got %s", DefaultProtocolVersion, client.protocolVersion)
	}
}

func TestNewClient_Defaults(t *testing.T) {
	cfg := Config{
		BaseURL: "http://localhost:8080",
	}

	client := NewClient(cfg)

	if client.protocolVersion != DefaultProtocolVersion {
		t.Errorf("Expected default protocolVersion %s, got %s", DefaultProtocolVersion, client.protocolVersion)
	}

	if client.clientVersion != DefaultClientVersion {
		t.Errorf("Expected default clientVersion %s, got %s", DefaultClientVersion, client.clientVersion)
	}

	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("Expected default timeout %v, got %v", DefaultTimeout, client.httpClient.Timeout)
	}
}

func TestSetToken(t *testing.T) {
	client := NewClient(Config{BaseURL: "http://localhost:8080"})
	client.SetToken("new-token")

	if client.token != "new-token" {
		t.Errorf("Expected token %s, got %s", "new-token", client.token)
	}
}

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathHealth {
			t.Errorf("Expected path %s, got %s", PathHealth, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","timestamp":"2025-01-18T10:00:00Z","version":"1.0.0"}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	health, err := client.Health(ctx)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", health.Status)
	}

	if health.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", health.Version)
	}
}

func TestHealth_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	_, err := client.Health(ctx)
	if err == nil {
		t.Error("Expected error for 500 status, got nil")
	}
}

func TestExecuteTemplate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1TemplatesExecute {
			t.Errorf("Expected path %s, got %s", PathAPIV1TemplatesExecute, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"execution_id": "exec-123",
				"status": "completed",
				"processing_time_ms": 100
			},
			"metadata": {
				"request_id": "req-123",
				"protocol_version": "1.0.0",
				"server_version": "1.0.0",
				"timestamp": 1640995200,
				"processing_time_ms": 100
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.ExecuteTemplateRequest{
		Query:    "test query",
		Language: "ru",
	}

	result, err := client.ExecuteTemplate(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteTemplate failed: %v", err)
	}

	if result.ExecutionID != "exec-123" {
		t.Errorf("Expected execution_id 'exec-123', got %s", result.ExecutionID)
	}

	if result.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", result.Status)
	}
}

func TestExecuteTemplate_ValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"error_code": "VALIDATION_FAILED",
				"error_type": "VALIDATION_ERROR",
				"message": "Query cannot be empty",
				"field": "query"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.ExecuteTemplateRequest{
		Query: "",
	}

	_, err := client.ExecuteTemplate(ctx, req)
	if err == nil {
		t.Error("Expected validation error, got nil")
		return
	}

	errDetail, ok := err.(*types.ErrorDetail)
	if !ok {
		t.Errorf("Expected ErrorDetail, got %T", err)
		return
	}

	if errDetail.Code != "VALIDATION_FAILED" {
		t.Errorf("Expected code 'VALIDATION_FAILED', got %s", errDetail.Code)
	}

	if !errDetail.IsValidationError() {
		t.Error("Expected IsValidationError() to return true")
	}
}

func TestContext_Cancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req := &types.ExecuteTemplateRequest{
		Query: "test",
	}

	_, err := client.ExecuteTemplate(ctx, req)
	if err == nil {
		t.Error("Expected context timeout error, got nil")
	}
}

