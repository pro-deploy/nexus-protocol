package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nexus-protocol/server/pkg/config"
	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration
func RegisterHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch services.(*Services).Config.Auth.Provider {
		case "keycloak":
			handleKeycloakRegistration(services, w, r)
		case "jwt":
			fallthrough
		default:
			handleLocalRegistration(services, w, r)
		}
	}
}

// handleLocalRegistration handles local JWT-based registration
func handleLocalRegistration(services interface{}, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		TenantID  string `json:"tenant_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request format",
		}, http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validateRegistrationRequest(); err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		services.(*Services).Logger.Error("Password hashing failed", zap.Error(err))
		sendError(w, types.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Type:    "INTERNAL_ERROR",
			Message: "Registration failed",
		}, http.StatusInternalServerError)
		return
	}

	// Create user
	userID := uuid.New().String()
	user := &types.User{
		ID:        userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    "active",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
	}

		if err := services.(*Services).Auth.RegisterUser(r.Context(), user, string(hashedPassword), req.TenantID); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			sendError(w, types.ErrorDetail{
				Code:    "RESOURCE_CONFLICT",
				Type:    "CONFLICT",
				Message: "User with this email already exists",
			}, http.StatusConflict)
			return
		}

		services.(*Services).Logger.Error("User registration failed", zap.Error(err))
		sendError(w, types.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Type:    "INTERNAL_ERROR",
			Message: "Registration failed",
		}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"user_id": userID,
		"message": "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// handleKeycloakRegistration handles Keycloak-based registration
func handleKeycloakRegistration(services interface{}, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		TenantID  string `json:"tenant_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request format",
		}, http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validateRegistrationRequest(); err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Register user in Keycloak
	userID, err := registerUserInKeycloak(req, services.(*Services).Config)
	if err != nil {
		services.(*Services).Logger.Error("Keycloak user registration failed", zap.Error(err))
		sendError(w, types.ErrorDetail{
			Code:    "EXTERNAL_ERROR",
			Type:    "EXTERNAL_ERROR",
			Message: "Registration failed",
			Details: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// Optionally sync to local database
	if services.(*Services).Auth != nil {
		user := &types.User{
			ID:        userID,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Status:    "active",
			Roles:     []string{"user"},
			CreatedAt: time.Now(),
		}

		// Don't fail if local sync fails - Keycloak is source of truth
		if err := services.(*Services).Auth.SyncUserFromKeycloak(r.Context(), user, req.TenantID); err != nil {
			services.(*Services).Logger.Warn("Local user sync failed", zap.Error(err))
		}
	}

	response := map[string]interface{}{
		"success":    true,
		"user_id":    userID,
		"message":    "User registered successfully in Keycloak",
		"login_url": fmt.Sprintf("%s/realms/%s/account", services.(*Services).Config.Auth.KeycloakURL, services.(*Services).Config.Auth.KeycloakRealm),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// LoginHandler handles user login
func LoginHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch services.(*Services).Config.Auth.Provider {
		case "keycloak":
			handleKeycloakLogin(services, w, r)
		case "jwt":
			fallthrough
		default:
			handleLocalLogin(services, w, r)
		}
	}
}

// handleLocalLogin handles local JWT-based login
func handleLocalLogin(services interface{}, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request format",
		}, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: "Email and password are required",
		}, http.StatusBadRequest)
		return
	}

		// Authenticate user
		user, err := services.(*Services).Auth.AuthenticateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "AUTHENTICATION_FAILED",
			Type:    "AUTHENTICATION_ERROR",
			Message: "Invalid email or password",
		}, http.StatusUnauthorized)
		return
	}

		// Generate tokens
		accessToken, refreshToken, err := services.(*Services).Auth.GenerateTokens(user)
	if err != nil {
		services.(*Services).Logger.Error("Token generation failed", zap.Error(err))
		sendError(w, types.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Type:    "INTERNAL_ERROR",
			Message: "Login failed",
		}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(services.(*Services).Config.Auth.JWTExpiry.Seconds()),
		"user":          user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleKeycloakLogin handles Keycloak-based login
func handleKeycloakLogin(services interface{}, w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request format",
		}, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		sendError(w, types.ErrorDetail{
			Code:    "VALIDATION_FAILED",
			Type:    "VALIDATION_ERROR",
			Message: "Username and password are required",
		}, http.StatusBadRequest)
		return
	}

	// Get token from Keycloak
	tokenResponse, err := getKeycloakToken(req.Username, req.Password, services.(*Services).Config)
	if err != nil {
		sendError(w, types.ErrorDetail{
			Code:    "AUTHENTICATION_FAILED",
			Type:    "AUTHENTICATION_ERROR",
			Message: "Invalid username or password",
		}, http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success":            true,
		"access_token":       tokenResponse.AccessToken,
		"refresh_token":      tokenResponse.RefreshToken,
		"token_type":         "Bearer",
		"expires_in":         tokenResponse.ExpiresIn,
		"refresh_expires_in": tokenResponse.RefreshExpiresIn,
		"login_url":          fmt.Sprintf("%s/realms/%s/account", services.(*Services).Config.Auth.KeycloakURL, services.(*Services).Config.Auth.KeycloakRealm),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		// Refresh token
		accessToken, refreshToken, err := services.(*Services).Auth.RefreshToken(r.Context(), req.RefreshToken)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "AUTHENTICATION_FAILED",
				Type:    "AUTHENTICATION_ERROR",
				Message: "Invalid refresh token",
			}, http.StatusUnauthorized)
			return
		}

		response := map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    int(services.(*Services).Config.Auth.JWTExpiry.Seconds()),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// GetUserProfileHandler handles getting user profile
func GetUserProfileHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		user, err := services.(*Services).Auth.GetUserProfile(r.Context(), userID)
		if err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "NOT_FOUND",
				Type:    "NOT_FOUND",
				Message: "User not found",
			}, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": user,
		})
	}
}

// UpdateUserProfileHandler handles updating user profile
func UpdateUserProfileHandler(services interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r.Context())

		var req struct {
			FirstName string `json:"first_name,omitempty"`
			LastName  string `json:"last_name,omitempty"`
			Bio       string `json:"bio,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Type:    "VALIDATION_ERROR",
				Message: "Invalid request format",
			}, http.StatusBadRequest)
			return
		}

		err := services.(*Services).Auth.UpdateUserProfile(r.Context(), userID, req.FirstName, req.LastName, req.Bio)
		if err != nil {
			services.(*Services).Logger.Error("Profile update failed", zap.Error(err))
			sendError(w, types.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Type:    "INTERNAL_ERROR",
				Message: "Profile update failed",
			}, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Profile updated successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func validateRegistrationRequest() error {
	// Implementation would validate email format, password strength, etc.
	return nil
}

// KeycloakTokenResponse represents Keycloak token response
type KeycloakTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

// registerUserInKeycloak registers user in Keycloak
func registerUserInKeycloak(req struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TenantID  string `json:"tenant_id,omitempty"`
}, cfg *config.Config) (string, error) {
	// Keycloak Admin API endpoint for user creation
	adminURL := fmt.Sprintf("%s/admin/realms/%s/users", cfg.Auth.KeycloakURL, cfg.Auth.KeycloakRealm)

	userData := map[string]interface{}{
		"username":  req.Email,
		"email":     req.Email,
		"firstName": req.FirstName,
		"lastName":  req.LastName,
		"enabled":   true,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     req.Password,
				"temporary": false,
			},
		},
	}

	// Add groups/attributes for tenant support
	if req.TenantID != "" {
		userData["attributes"] = map[string]interface{}{
			"tenant_id": []string{req.TenantID},
		}
		userData["groups"] = []string{fmt.Sprintf("/tenants/%s", req.TenantID)}
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Get admin token first
	adminToken, err := getKeycloakAdminToken(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get admin token: %w", err)
	}

	httpReq, err := http.NewRequest("POST", adminURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+adminToken)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("user creation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("user creation failed with status %d", resp.StatusCode)
	}

	// Extract user ID from Location header
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no location header in response")
	}

	// Location format: http://keycloak/auth/admin/realms/nexus/users/{user-id}
	parts := strings.Split(location, "/")
	userID := parts[len(parts)-1]

	return userID, nil
}

// getKeycloakToken gets access token from Keycloak
func getKeycloakToken(username, password string, cfg *config.Config) (*KeycloakTokenResponse, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		cfg.Auth.KeycloakURL, cfg.Auth.KeycloakRealm)

	data := fmt.Sprintf("grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		cfg.Auth.KeycloakClientID, cfg.Auth.KeycloakSecret, username, password)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}

	var tokenResp KeycloakTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// getKeycloakAdminToken gets admin token for Keycloak Admin API
func getKeycloakAdminToken(cfg *config.Config) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", cfg.Auth.KeycloakURL)

	// This requires admin credentials - in production, use service account
	data := "grant_type=password&client_id=admin-cli&username=admin&password=admin"

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create admin token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("admin token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("admin token request failed with status %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode admin token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
