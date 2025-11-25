package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func TestRegisterUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1AuthRegister {
			t.Errorf("Expected path %s, got %s", PathAPIV1AuthRegister, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"data": {
				"user_id": "user-123",
				"message": "User registered successfully",
				"verification_required": true
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.RegisterUserRequest{
		Email:     "user@example.com",
		Password:  "password123",
		FirstName: "Иван",
		LastName:  "Иванов",
	}

	resp, err := client.RegisterUser(ctx, req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if resp.UserID != "user-123" {
		t.Errorf("Expected user_id 'user-123', got %s", resp.UserID)
	}

	if !resp.VerificationRequired {
		t.Error("Expected verification_required to be true")
	}
}

func TestLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1AuthLogin {
			t.Errorf("Expected path %s, got %s", PathAPIV1AuthLogin, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"access_token": "token-123",
				"refresh_token": "refresh-123",
				"token_type": "Bearer",
				"expires_in": 3600,
				"user": {
					"id": "user-123",
					"email": "user@example.com",
					"status": "active",
					"roles": ["user"],
					"created_at": 1640995200
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.AccessToken != "token-123" {
		t.Errorf("Expected access_token 'token-123', got %s", resp.AccessToken)
	}

	if client.token != "token-123" {
		t.Error("Expected token to be automatically set in client")
	}

	if resp.User == nil {
		t.Error("Expected user in response")
	}
}

func TestRefreshToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"access_token": "new-token-123",
				"token_type": "Bearer",
				"expires_in": 3600
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.RefreshTokenRequest{
		RefreshToken: "refresh-123",
	}

	resp, err := client.RefreshToken(ctx, req)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	if resp.AccessToken != "new-token-123" {
		t.Errorf("Expected access_token 'new-token-123', got %s", resp.AccessToken)
	}

	if client.token != "new-token-123" {
		t.Error("Expected token to be automatically updated in client")
	}
}

func TestGetUserProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1UsersProfile {
			t.Errorf("Expected path %s, got %s", PathAPIV1UsersProfile, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"id": "user-123",
				"email": "user@example.com",
				"username": "ivan",
				"first_name": "Иван",
				"last_name": "Иванов",
				"status": "active",
				"roles": ["user", "admin"],
				"created_at": 1640995200
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})
	ctx := context.Background()

	profile, err := client.GetUserProfile(ctx)
	if err != nil {
		t.Fatalf("GetUserProfile failed: %v", err)
	}

	if profile.ID != "user-123" {
		t.Errorf("Expected id 'user-123', got %s", profile.ID)
	}

	if profile.Email != "user@example.com" {
		t.Errorf("Expected email 'user@example.com', got %s", profile.Email)
	}

	if len(profile.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(profile.Roles))
	}
}

func TestUpdateUserProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"id": "user-123",
				"email": "user@example.com",
				"first_name": "Иван",
				"last_name": "Петров",
				"status": "active",
				"roles": ["user"],
				"created_at": 1640995200
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, Token: "test-token"})
	ctx := context.Background()

	req := &types.UpdateProfileRequest{
		FirstName: "Иван",
		LastName:  "Петров",
		Bio:       "Разработчик",
	}

	profile, err := client.UpdateUserProfile(ctx, req)
	if err != nil {
		t.Fatalf("UpdateUserProfile failed: %v", err)
	}

	if profile.LastName != "Петров" {
		t.Errorf("Expected last_name 'Петров', got %s", profile.LastName)
	}
}

func TestLogin_AuthenticationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"error": {
				"code": "AUTHENTICATION_FAILED",
				"type": "AUTHENTICATION_ERROR",
				"message": "Invalid credentials"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.LoginRequest{
		Email:    "user@example.com",
		Password: "wrong-password",
	}

	_, err := client.Login(ctx, req)
	if err == nil {
		t.Error("Expected authentication error, got nil")
		return
	}

	errDetail, ok := err.(*types.ErrorDetail)
	if !ok {
		t.Errorf("Expected ErrorDetail, got %T", err)
		return
	}

	if !errDetail.IsAuthenticationError() {
		t.Error("Expected IsAuthenticationError() to return true")
	}
}

