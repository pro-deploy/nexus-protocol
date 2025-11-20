package client

import (
	"context"

	"github.com/nexus-protocol/go-sdk/types"
)

// RegisterUser регистрирует нового пользователя в системе.
// Требует email, password, first_name и last_name.
//
// Пример использования:
//
//	ctx := context.Background()
//	req := &types.RegisterUserRequest{
//		Email:     "user@example.com",
//		Password:  "secure-password",
//		FirstName: "Иван",
//		LastName:  "Иванов",
//	}
//	resp, err := client.RegisterUser(ctx, req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User ID: %s\n", resp.UserID)
func (c *Client) RegisterUser(ctx context.Context, req *types.RegisterUserRequest) (*types.RegisterUserResponse, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1AuthRegister, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.RegisterUserResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// Login выполняет аутентификацию пользователя и возвращает JWT токены.
// Access token автоматически устанавливается в клиент для последующих запросов.
//
// Пример использования:
//
//	ctx := context.Background()
//	req := &types.LoginRequest{
//		Email:    "user@example.com",
//		Password: "password",
//	}
//	resp, err := client.Login(ctx, req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Access token: %s\n", resp.AccessToken)
func (c *Client) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1AuthLogin, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.LoginResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	// Автоматически устанавливаем токен при успешном входе
	if result.Data.AccessToken != "" {
		c.SetToken(result.Data.AccessToken)
	}

	return &result.Data, nil
}

// RefreshToken обновляет access token используя refresh token.
// Новый access token автоматически устанавливается в клиент.
func (c *Client) RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.RefreshTokenResponse, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1AuthRefresh, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.RefreshTokenResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	// Автоматически обновляем токен
	if result.Data.AccessToken != "" {
		c.SetToken(result.Data.AccessToken)
	}

	return &result.Data, nil
}

// GetUserProfile получает профиль текущего аутентифицированного пользователя.
// Требует валидный JWT токен.
func (c *Client) GetUserProfile(ctx context.Context) (*types.UserProfile, error) {
	resp, err := c.doRequest(ctx, "GET", PathAPIV1UsersProfile, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.UserProfile `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// UpdateUserProfile обновляет профиль текущего аутентифицированного пользователя.
// Можно обновить first_name, last_name и bio.
func (c *Client) UpdateUserProfile(ctx context.Context, req *types.UpdateProfileRequest) (*types.UserProfile, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "PUT", PathAPIV1UsersProfile, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.UserProfile `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

