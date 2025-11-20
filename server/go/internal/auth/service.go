package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nexus-protocol/server/pkg/config"
	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication and authorization
type Service struct {
	config *config.Config
	logger *zap.Logger
	db     *sql.DB // In production, use proper database interface
}

// Claims represents JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	TenantID string   `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

// NewService creates a new auth service
func NewService(cfg *config.Config, logger *zap.Logger) *Service {
	return &Service{
		config: cfg,
		logger: logger,
	}
}

// RegisterUser registers a new user
func (s *Service) RegisterUser(ctx context.Context, user *types.User, hashedPassword string, tenantID string) error {
	// In production, save to database
	s.logger.Info("User registered",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email))

	// Mock implementation - in production would insert into database
	return nil
}

// AuthenticateUser authenticates user credentials
func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (*types.User, error) {
	// Mock authentication - in production would query database
	if email == "test@example.com" && password == "password123" {
		return &types.User{
			ID:        uuid.New().String(),
			Email:     email,
			FirstName: "Test",
			LastName:  "User",
			Status:    "active",
			Roles:     []string{"user"},
			CreatedAt: time.Now(),
		}, nil
	}

	return nil, fmt.Errorf("invalid credentials")
}

// GenerateTokens generates access and refresh tokens
func (s *Service) GenerateTokens(user *types.User) (string, string, error) {
	// Create access token
	accessClaims := Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Roles:    user.Roles,
		TenantID: "", // Would get from user context
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "nexus-protocol",
			Subject:   user.ID,
			ID:        uuid.New().String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token
	refreshClaims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "nexus-protocol",
			Subject:   user.ID,
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateToken validates JWT token and returns claims
func (s *Service) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return map[string]interface{}{
			"user_id":   claims.UserID,
			"email":     claims.Email,
			"roles":     claims.Roles,
			"tenant_id": claims.TenantID,
			"exp":       claims.ExpiresAt.Unix(),
			"iat":       claims.IssuedAt.Unix(),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshToken refreshes access token using refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid refresh token claims")
	}

	// Get user from database (mock implementation)
	user := &types.User{
		ID:    userID,
		Email: claims["email"].(string),
		Roles: []string{"user"}, // Would get from database
	}

	// Generate new tokens
	return s.GenerateTokens(user)
}

// GetUserProfile gets user profile
func (s *Service) GetUserProfile(ctx context.Context, userID string) (*types.User, error) {
	// Mock implementation - in production would query database
	return &types.User{
		ID:        userID,
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    "active",
		Roles:     []string{"user"},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		LastLogin: &time.Time{},
	}, nil
}

// UpdateUserProfile updates user profile
func (s *Service) UpdateUserProfile(ctx context.Context, userID, firstName, lastName, bio string) error {
	// Mock implementation - in production would update database
	s.logger.Info("User profile updated",
		zap.String("user_id", userID),
		zap.String("first_name", firstName),
		zap.String("last_name", lastName))
	return nil
}

// SyncUserFromKeycloak syncs user data from Keycloak (for hybrid mode)
func (s *Service) SyncUserFromKeycloak(ctx context.Context, user *types.User, tenantID string) error {
	// This method is called when using Keycloak auth but local database sync is needed
	s.logger.Info("User synced from Keycloak",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("tenant_id", tenantID))
	return nil
}

// ValidateUserPermissions validates if user has required permissions
func (s *Service) ValidateUserPermissions(userID string, requiredRoles []string, resource string, action string) error {
	// Get user roles (mock implementation)
	userRoles := []string{"user"} // Would get from database/cache

	// Check if user has required roles
	for _, requiredRole := range requiredRoles {
		found := false
		for _, userRole := range userRoles {
			if userRole == requiredRole {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("insufficient permissions: missing role %s", requiredRole)
		}
	}

	s.logger.Debug("User permissions validated",
		zap.String("user_id", userID),
		zap.Strings("required_roles", requiredRoles),
		zap.String("resource", resource),
		zap.String("action", action))

	return nil
}

// GetUserTenants gets tenants for a user
func (s *Service) GetUserTenants(ctx context.Context, userID string) ([]string, error) {
	// Mock implementation - in production would query database
	return []string{"tenant-1", "tenant-2"}, nil
}

// HashPassword hashes a password
func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword checks if password matches hash
func (s *Service) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RevokeUserTokens revokes all tokens for a user (logout from all devices)
func (s *Service) RevokeUserTokens(ctx context.Context, userID string) error {
	// In production, this would add user to token blacklist or increment token version
	s.logger.Info("User tokens revoked", zap.String("user_id", userID))
	return nil
}

// GetUserSessions gets active sessions for a user
func (s *Service) GetUserSessions(ctx context.Context, userID string) ([]*UserSession, error) {
	// Mock implementation
	return []*UserSession{
		{
			ID:         uuid.New().String(),
			UserID:     userID,
			UserAgent:  "Mozilla/5.0...",
			IPAddress:  "192.168.1.1",
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			LastActivity: time.Now().Add(-5 * time.Minute),
		},
	}, nil
}

// UserSession represents a user session
type UserSession struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
}

// TerminateSession terminates a specific session
func (s *Service) TerminateSession(ctx context.Context, userID, sessionID string) error {
	s.logger.Info("Session terminated",
		zap.String("user_id", userID),
		zap.String("session_id", sessionID))
	return nil
}
