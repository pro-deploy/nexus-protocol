package main

import (
	"context"
	"fmt"
	"log"

	nexus "github.com/nexus-protocol/go-sdk/client"
	"github.com/nexus-protocol/go-sdk/types"
)

func main() {
	cfg := nexus.Config{
		BaseURL: "http://localhost:8080",
	}

	client := nexus.NewClient(cfg)
	ctx := context.Background()

	// Пример регистрации пользователя
	fmt.Println("Регистрация пользователя...")
	registerReq := &types.RegisterUserRequest{
		Email:     "user@example.com",
		Password:  "password123",
		FirstName: "Иван",
		LastName:  "Иванов",
	}

	registerResp, err := client.RegisterUser(ctx, registerReq)
	if err != nil {
		log.Printf("Ошибка регистрации: %v", err)
	} else {
		fmt.Printf("✓ Пользователь зарегистрирован: %s\n", registerResp.UserID)
		if registerResp.VerificationRequired {
			fmt.Println("  Требуется верификация email")
		}
	}

	// Пример входа
	fmt.Println("\nВход в систему...")
	loginReq := &types.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	}

	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		log.Printf("Ошибка входа: %v", err)
		return
	}

	fmt.Printf("✓ Вход выполнен успешно\n")
	fmt.Printf("  Access Token: %s...\n", loginResp.AccessToken[:20])
	fmt.Printf("  Token Type: %s\n", loginResp.TokenType)
	fmt.Printf("  Expires In: %d секунд\n", loginResp.ExpiresIn)

	if loginResp.User != nil {
		fmt.Printf("  User: %s %s (%s)\n",
			loginResp.User.FirstName,
			loginResp.User.LastName,
			loginResp.User.Email,
		)
	}

	// Получение профиля
	fmt.Println("\nПолучение профиля...")
	profile, err := client.GetUserProfile(ctx)
	if err != nil {
		log.Printf("Ошибка получения профиля: %v", err)
	} else {
		fmt.Printf("✓ Профиль получен:\n")
		fmt.Printf("  ID: %s\n", profile.ID)
		fmt.Printf("  Email: %s\n", profile.Email)
		fmt.Printf("  Status: %s\n", profile.Status)
		fmt.Printf("  Roles: %v\n", profile.Roles)
	}

	// Обновление профиля
	fmt.Println("\nОбновление профиля...")
	updateReq := &types.UpdateProfileRequest{
		FirstName: "Иван",
		LastName:  "Петров",
		Bio:       "Разработчик Go",
	}

	updatedProfile, err := client.UpdateUserProfile(ctx, updateReq)
	if err != nil {
		log.Printf("Ошибка обновления профиля: %v", err)
	} else {
		fmt.Printf("✓ Профиль обновлен:\n")
		fmt.Printf("  Имя: %s %s\n", updatedProfile.FirstName, updatedProfile.LastName)
	}

	// Обновление токена
	fmt.Println("\nОбновление токена...")
	refreshReq := &types.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}

	refreshResp, err := client.RefreshToken(ctx, refreshReq)
	if err != nil {
		log.Printf("Ошибка обновления токена: %v", err)
	} else {
		fmt.Printf("✓ Токен обновлен\n")
		fmt.Printf("  Новый Access Token: %s...\n", refreshResp.AccessToken[:20])
		fmt.Printf("  Expires In: %d секунд\n", refreshResp.ExpiresIn)
	}
}

