package middleware

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nexus-protocol/server/pkg/config"
	"go.uber.org/zap"
)

// KeycloakClaims represents Keycloak JWT claims
type KeycloakClaims struct {
	Sub               string                 `json:"sub"`
	PreferredUsername string                 `json:"preferred_username"`
	Email             string                 `json:"email"`
	Name              string                 `json:"name"`
	GivenName         string                 `json:"given_name"`
	FamilyName        string                 `json:"family_name"`
	Roles             []string               `json:"roles,omitempty"`
	Groups            []string               `json:"groups,omitempty"`
	RealmAccess       map[string]interface{} `json:"realm_access,omitempty"`
	ResourceAccess    map[string]interface{} `json:"resource_access,omitempty"`
}

// Auth middleware validates tokens (JWT or Keycloak)
func Auth(cfg *config.Config, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string
			var claims interface{}

			switch cfg.Auth.Provider {
			case "keycloak":
				var err error
				userID, claims, err = validateKeycloakToken(r, cfg, logger)
				if err != nil {
					logger.Warn("Keycloak token validation failed", zap.Error(err))
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			case "jwt":
				fallthrough
			default:
				var err error
				userID, claims, err = validateJWToken(r, cfg, logger)
				if err != nil {
					logger.Warn("JWT token validation failed", zap.Error(err))
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// Add user context to request
			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "claims", claims)
			r = r.WithContext(ctx)

			// Continue with request
			next.ServeHTTP(w, r)
		})
	}
}

// validateKeycloakToken validates Keycloak JWT token
func validateKeycloakToken(r *http.Request, cfg *config.Config, logger *zap.Logger) (string, *KeycloakClaims, error) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, fmt.Errorf("missing Authorization header")
	}

	// Check Bearer token format
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", nil, fmt.Errorf("invalid Authorization header format")
	}

	tokenString := parts[1]

	// Validate token with Keycloak introspection endpoint
	claims, err := introspectKeycloakToken(tokenString, cfg)
	if err != nil {
		return "", nil, fmt.Errorf("token introspection failed: %w", err)
	}

	return claims.Sub, claims, nil
}

// introspectKeycloakToken validates token via Keycloak introspection
func introspectKeycloakToken(token string, cfg *config.Config) (*KeycloakClaims, error) {
	// Keycloak token introspection endpoint
	introspectURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
		cfg.Auth.KeycloakURL, cfg.Auth.KeycloakRealm)

	// Prepare form data
	data := fmt.Sprintf("token=%s&client_id=%s&client_secret=%s",
		token, cfg.Auth.KeycloakClientID, cfg.Auth.KeycloakSecret)

	// Create HTTP request
	req, err := http.NewRequest("POST", introspectURL, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Skip TLS verification for development (remove in production)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("introspection request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("introspection returned status %d", resp.StatusCode)
	}

	// Parse response
	var introspectionResp struct {
		Active         bool                   `json:"active"`
		Sub            string                 `json:"sub"`
		Username       string                 `json:"preferred_username"`
		Email          string                 `json:"email"`
		Name           string                 `json:"name"`
		GivenName      string                 `json:"given_name"`
		FamilyName     string                 `json:"family_name"`
		Groups         []string               `json:"groups"`
		RealmAccess    map[string]interface{} `json:"realm_access"`
		ResourceAccess map[string]interface{} `json:"resource_access"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&introspectionResp); err != nil {
		return nil, fmt.Errorf("failed to decode introspection response: %w", err)
	}

	if !introspectionResp.Active {
		return nil, fmt.Errorf("token is not active")
	}

	// Extract roles from realm_access
	var roles []string
	if introspectionResp.RealmAccess != nil {
		if rolesInterface, ok := introspectionResp.RealmAccess["roles"]; ok {
			if rolesSlice, ok := rolesInterface.([]interface{}); ok {
				for _, role := range rolesSlice {
					if roleStr, ok := role.(string); ok {
						roles = append(roles, roleStr)
					}
				}
			}
		}
	}

	claims := &KeycloakClaims{
		Sub:               introspectionResp.Sub,
		PreferredUsername: introspectionResp.Username,
		Email:             introspectionResp.Email,
		Name:              introspectionResp.Name,
		GivenName:         introspectionResp.GivenName,
		FamilyName:        introspectionResp.FamilyName,
		Roles:             roles,
		Groups:            introspectionResp.Groups,
		RealmAccess:       introspectionResp.RealmAccess,
		ResourceAccess:    introspectionResp.ResourceAccess,
	}

	return claims, nil
}

// validateJWToken validates local JWT token (fallback)
func validateJWToken(r *http.Request, cfg *config.Config, logger *zap.Logger) (string, interface{}, error) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, fmt.Errorf("missing Authorization header")
	}

	// Check Bearer token format
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", nil, fmt.Errorf("invalid Authorization header format")
	}

	tokenString := parts[1]

	// Parse and validate JWT token
	claims, err := parseJWTClaims(tokenString, cfg.Auth.JWTSecret)
	if err != nil {
		return "", nil, fmt.Errorf("invalid token: %w", err)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", nil, fmt.Errorf("missing user_id in token claims")
	}

	return userID, claims, nil
}

// parseJWTClaims extracts claims from JWT token
func parseJWTClaims(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Logging middleware logs HTTP requests
func Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Log request
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr))

			next.ServeHTTP(w, r)

			// Log response
			logger.Info("HTTP response",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", time.Since(start)))
		})
	}
}

// CORS middleware adds CORS headers
func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Recovery middleware recovers from panics
func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered", zap.Any("error", err))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
