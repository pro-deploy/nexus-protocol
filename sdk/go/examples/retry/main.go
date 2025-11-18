package main

import (
	"context"
	"fmt"
	"log"
	"time"

	nexus "github.com/nexus-protocol/go-sdk/client"
	"github.com/nexus-protocol/go-sdk/types"
)

func main() {
	// Создаем клиент с retry конфигурацией
	retryCfg := nexus.RetryConfig{
		MaxRetries:        3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableStatusCodes: []int{
			408, // Request Timeout
			429, // Too Many Requests
			500, // Internal Server Error
			502, // Bad Gateway
			503, // Service Unavailable
			504, // Gateway Timeout
		},
	}

	// Создаем логгер для отслеживания retry
	logger := nexus.NewSimpleLogger(nexus.LogLevelDebug)

	cfg := nexus.Config{
		BaseURL:     "http://localhost:8080",
		Token:       "your-jwt-token",
		RetryConfig: &retryCfg,
		Logger:      logger,
	}

	client := nexus.NewClient(cfg)
	ctx := context.Background()

	fmt.Println("Выполнение запроса с retry логикой...")
	fmt.Println("(Если сервер вернет ошибку 500, запрос будет повторен автоматически)")

	req := &types.ExecuteTemplateRequest{
		Query:    "хочу борщ",
		Language: "ru",
	}

	result, err := client.ExecuteTemplate(ctx, req)
	if err != nil {
		log.Printf("Ошибка после всех попыток: %v", err)
		return
	}

	fmt.Printf("✓ Запрос выполнен успешно\n")
	fmt.Printf("  Execution ID: %s\n", result.ExecutionID)
	fmt.Printf("  Status: %s\n", result.Status)
}

