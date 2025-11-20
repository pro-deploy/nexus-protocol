package types

// RegisterUserRequest представляет запрос регистрации пользователя
type RegisterUserRequest struct {
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	TenantID  string          `json:"tenant_id,omitempty"`
	Metadata  *RequestMetadata `json:"metadata,omitempty"`
}

// RegisterUserResponse представляет ответ регистрации
type RegisterUserResponse struct {
	UserID             string `json:"user_id"`
	Message            string `json:"message"`
	VerificationRequired bool `json:"verification_required"`
}

// LoginRequest представляет запрос входа
type LoginRequest struct {
	Email    string           `json:"email"`
	Password string           `json:"password"`
	Metadata *RequestMetadata `json:"metadata,omitempty"`
}

// LoginResponse представляет ответ входа
type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	ExpiresIn    int32      `json:"expires_in"`
	User         *UserProfile `json:"user"`
}

// RefreshTokenRequest представляет запрос обновления токена
type RefreshTokenRequest struct {
	RefreshToken string          `json:"refresh_token"`
	Metadata     *RequestMetadata `json:"metadata,omitempty"`
}

// RefreshTokenResponse представляет ответ обновления токена
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
}

// UserProfile представляет профиль пользователя
type UserProfile struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	Username   string   `json:"username,omitempty"`
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	Status     string   `json:"status"`
	Roles      []string `json:"roles"`
	CreatedAt  int64    `json:"created_at"`
	LastLoginAt int64   `json:"last_login_at,omitempty"`
}

// UpdateProfileRequest представляет запрос обновления профиля
type UpdateProfileRequest struct {
	FirstName string           `json:"first_name,omitempty"`
	LastName  string           `json:"last_name,omitempty"`
	Bio       string           `json:"bio,omitempty"`
	Metadata  *RequestMetadata `json:"metadata,omitempty"`
}

