// +build integration

package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nexus-protocol/go-sdk/types"
)

// Integration тесты требуют запущенный сервер
// Запуск: go test -tags=integration ./client/...

func TestIntegration_Health(t *testing.T) {
	baseURL := os.Getenv("NEXUS_API_URL")
	if baseURL == "" {
		t.Skip("NEXUS_API_URL not set, skipping integration test")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		Timeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	health, err := client.Health(ctx)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	if health.Status == "" {
		t.Error("Expected status in health response")
	}
}

func TestIntegration_ExecuteTemplate(t *testing.T) {
	baseURL := os.Getenv("NEXUS_API_URL")
	token := os.Getenv("NEXUS_API_TOKEN")
	
	if baseURL == "" {
		t.Skip("NEXUS_API_URL not set, skipping integration test")
	}
	if token == "" {
		t.Skip("NEXUS_API_TOKEN not set, skipping integration test")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		Token:   token,
		Timeout: 30 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &types.ExecuteTemplateRequest{
		Query:    "тестовый запрос",
		Language: "ru",
	}

	result, err := client.ExecuteTemplate(ctx, req)
	if err != nil {
		t.Fatalf("ExecuteTemplate failed: %v", err)
	}

	if result.ExecutionID == "" {
		t.Error("Expected execution ID in response")
	}
}

func TestIntegration_Login(t *testing.T) {
	baseURL := os.Getenv("NEXUS_API_URL")
	email := os.Getenv("NEXUS_TEST_EMAIL")
	password := os.Getenv("NEXUS_TEST_PASSWORD")
	
	if baseURL == "" {
		t.Skip("NEXUS_API_URL not set, skipping integration test")
	}
	if email == "" || password == "" {
		t.Skip("NEXUS_TEST_EMAIL or NEXUS_TEST_PASSWORD not set, skipping integration test")
	}

	client := NewClient(Config{
		BaseURL: baseURL,
		Timeout: 10 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &types.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("Expected access token in response")
	}

	if client.token != resp.AccessToken {
		t.Error("Expected token to be set in client")
	}
}

