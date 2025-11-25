package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
	// Создаем клиент
	cfg := nexus.Config{
		BaseURL:         getEnv("NEXUS_BASE_URL", "http://localhost:8080"),
		Token:           getEnv("NEXUS_TOKEN", ""),
		ProtocolVersion: "2.0.0", // Nexus Protocol v2.0.0
		ClientVersion:   "2.0.0",
		ClientID:        "go-sdk-example",
		ClientType:      "sdk",
	}

	nexusClient := nexus.NewClient(cfg)
	ctx := context.Background()

	// Проверяем здоровье сервера
	fmt.Println("Проверка здоровья сервера...")
	health, err := nexusClient.Health(ctx)
	if err != nil {
		log.Fatalf("Сервер недоступен: %v", err)
	}
	fmt.Printf("✓ Сервер доступен (версия: %s)\n", health.Version)

	// Выполняем шаблон
	fmt.Println("\nВыполнение шаблона...")
	req := &types.ExecuteTemplateRequest{
		Query:    "хочу борщ",
		Language: "ru",
		Context: &types.UserContext{
			UserID:    "user-123",
			SessionID: "session-456",
		},
		Options: &types.ExecuteOptions{
			TimeoutMS:           30000,
			MaxResultsPerDomain: 5,
			ParallelExecution:   true,
			IncludeWebSearch:    true,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := nexusClient.ExecuteTemplate(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка выполнения шаблона: %v", err)
	}

	fmt.Printf("✓ Шаблон выполнен\n")
	fmt.Printf("  Execution ID: %s\n", result.ExecutionID)
	fmt.Printf("  Status: %s\n", result.Status)
	fmt.Printf("  Processing time: %d ms\n", result.ProcessingTimeMS)

	if result.ResponseMetadata != nil {
		fmt.Printf("  Server version: %s\n", result.ResponseMetadata.ServerVersion)
		fmt.Printf("  Protocol version: %s\n", result.ResponseMetadata.ProtocolVersion)
	}

	// Получаем статус выполнения
	if result.ExecutionID != "" {
		fmt.Println("\nПолучение статуса выполнения...")
		status, err := nexusClient.GetExecutionStatus(ctx, result.ExecutionID)
		if err != nil {
			log.Printf("Ошибка получения статуса: %v", err)
		} else {
			fmt.Printf("✓ Статус получен: %s\n", status.Status)
		}
	}

	// Выводим результаты по доменам
	if len(result.Sections) > 0 {
		fmt.Println("\nРезультаты по доменам:")
		for _, section := range result.Sections {
			fmt.Printf("  - %s (%s): %d результатов\n",
				section.DomainID,
				section.Status,
				len(section.Results),
			)
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

