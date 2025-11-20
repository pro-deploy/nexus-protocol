package main

import (
	"context"
	"fmt"
	"log"
	"time"

	nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// Пример использования webhooks для асинхронной обработки результатов
func main() {
	cfg := nexus.Config{
		BaseURL:         "http://localhost:8080",
		Token:           "enterprise-jwt-token",
		ProtocolVersion: "1.1.0",
		ClientVersion:   "2.0.0",
		ClientID:        "enterprise-webhook-manager",
		ClientType:      "api",
	}

	client := nexus.NewClient(cfg)
	ctx := context.Background()

	fmt.Println("🪝 Webhooks Demo")
	fmt.Println("=================")

	// Пример 1: Регистрация webhook
	webhookID := demonstrateWebhookRegistration(ctx, client)

	// Пример 2: Получение списка webhooks
	demonstrateListWebhooks(ctx, client)

	// Пример 3: Тестирование webhook
	if webhookID != "" {
		demonstrateTestWebhook(ctx, client, webhookID)
	}

	// Пример 4: Мониторинг webhook статистики
	if webhookID != "" {
		demonstrateWebhookMonitoring(ctx, client, webhookID)
	}

	// Пример 5: Удаление webhook
	if webhookID != "" {
		demonstrateWebhookDeletion(ctx, client, webhookID)
	}
}

func demonstrateWebhookRegistration(ctx context.Context, client *nexus.Client) string {
	fmt.Println("\n1️⃣ Регистрация webhook")

	// Регистрируем webhook для получения уведомлений о завершении операций
	webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
		Config: &types.WebhookConfig{
			URL:    "https://enterprise-app.company.com/webhooks/nexus",
			Events: []string{
				"template.completed",
				"template.failed",
				"batch.completed",
				"batch.failed",
			},
			Secret: "enterprise-webhook-secret-2024",
			RetryPolicy: &types.WebhookRetryPolicy{
				MaxRetries:    3,
				InitialDelay: 1000,  // 1 секунда
				MaxDelay:      30000, // 30 секунд
				BackoffFactor: 2.0,
			},
			Headers: map[string]string{
				"X-API-Key":    "webhook-api-key",
				"X-Client-ID":  "enterprise-app",
				"X-Timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
			},
			Active:      true,
			Description: "Enterprise webhook for async operations monitoring",
		},
	})

	if err != nil {
		log.Fatalf("Webhook registration failed: %v", err)
	}

	fmt.Printf("✅ Webhook зарегистрирован:\n")
	fmt.Printf("   - Webhook ID: %s\n", webhookResp.WebhookID)
	fmt.Printf("   - Status: %s\n", webhookResp.Status)
	if webhookResp.Message != "" {
		fmt.Printf("   - Message: %s\n", webhookResp.Message)
	}

	return webhookResp.WebhookID
}

func demonstrateListWebhooks(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n2️⃣ Получение списка webhooks")

	// Получаем все активные webhooks
	webhooks, err := client.ListWebhooks(ctx, &types.ListWebhooksRequest{
		ActiveOnly: true,
		Limit:       10,
		Offset:      0,
	})

	if err != nil {
		log.Fatalf("List webhooks failed: %v", err)
	}

	fmt.Printf("✅ Найдено webhooks: %d\n", webhooks.Total)
	fmt.Printf("   - Лимит: %d\n", webhooks.Limit)
	fmt.Printf("   - Смещение: %d\n", webhooks.Offset)

	for i, wh := range webhooks.Webhooks {
		fmt.Printf("\n   Webhook #%d:\n", i+1)
		fmt.Printf("     - ID: %s\n", wh.ID)
		fmt.Printf("     - URL: %s\n", wh.Config.URL)
		fmt.Printf("     - Events: %v\n", wh.Config.Events)
		fmt.Printf("     - Active: %v\n", wh.Config.Active)
		fmt.Printf("     - Created: %s\n", time.Unix(wh.CreatedAt, 0).Format(time.RFC3339))
		if wh.LastUsedAt > 0 {
			fmt.Printf("     - Last used: %s\n", time.Unix(wh.LastUsedAt, 0).Format(time.RFC3339))
		}
		fmt.Printf("     - Success count: %d\n", wh.SuccessCount)
		fmt.Printf("     - Error count: %d\n", wh.ErrorCount)
		if wh.SuccessCount+wh.ErrorCount > 0 {
			successRate := float64(wh.SuccessCount) / float64(wh.SuccessCount+wh.ErrorCount) * 100
			fmt.Printf("     - Success rate: %.1f%%\n", successRate)
		}
	}
}

func demonstrateTestWebhook(ctx context.Context, client *nexus.Client, webhookID string) {
	fmt.Println("\n3️⃣ Тестирование webhook")

	// Отправляем тестовое событие
	testResp, err := client.TestWebhook(ctx, &types.TestWebhookRequest{
		WebhookID: webhookID,
		Event:     "template.completed",
		Data: map[string]interface{}{
			"execution_id": "test-exec-123",
			"status":       "completed",
			"timestamp":    time.Now().Unix(),
			"test":         true,
		},
	})

	if err != nil {
		log.Fatalf("Webhook test failed: %v", err)
	}

	fmt.Printf("✅ Webhook тест выполнен:\n")
	fmt.Printf("   - Webhook ID: %s\n", testResp.WebhookID)
	fmt.Printf("   - Status: %s\n", testResp.Status)
	fmt.Printf("   - Response code: %d\n", testResp.ResponseCode)
	fmt.Printf("   - Response time: %d ms\n", testResp.ResponseTimeMS)

	if testResp.Error != "" {
		fmt.Printf("   - Error: %s\n", testResp.Error)
	} else {
		fmt.Printf("   - ✅ Webhook успешно получил тестовое событие\n")
	}
}

func demonstrateWebhookMonitoring(ctx context.Context, client *nexus.Client, webhookID string) {
	fmt.Println("\n4️⃣ Мониторинг webhook статистики")

	// Получаем обновленный список для мониторинга
	webhooks, err := client.ListWebhooks(ctx, &types.ListWebhooksRequest{
		ActiveOnly: true,
		Limit:      100,
	})

	if err != nil {
		log.Fatalf("List webhooks failed: %v", err)
	}

	// Находим наш webhook
	var targetWebhook *types.WebhookInfo
	for _, wh := range webhooks.Webhooks {
		if wh.ID == webhookID {
			targetWebhook = &wh
			break
		}
	}

	if targetWebhook == nil {
		fmt.Printf("⚠️ Webhook %s не найден\n", webhookID)
		return
	}

	fmt.Printf("✅ Статистика webhook %s:\n", webhookID)
	fmt.Printf("   - Всего отправок: %d\n", targetWebhook.SuccessCount+targetWebhook.ErrorCount)
	fmt.Printf("   - Успешных: %d\n", targetWebhook.SuccessCount)
	fmt.Printf("   - Ошибок: %d\n", targetWebhook.ErrorCount)

	if targetWebhook.SuccessCount+targetWebhook.ErrorCount > 0 {
		successRate := float64(targetWebhook.SuccessCount) / float64(targetWebhook.SuccessCount+targetWebhook.ErrorCount) * 100
		fmt.Printf("   - Success rate: %.1f%%\n", successRate)

		if successRate < 95.0 {
			fmt.Printf("   - ⚠️ Низкий success rate! Проверьте webhook endpoint\n")
		} else {
			fmt.Printf("   - ✅ Отличный success rate\n")
		}
	}

	if targetWebhook.LastUsedAt > 0 {
		lastUsed := time.Unix(targetWebhook.LastUsedAt, 0)
		timeSinceLastUse := time.Since(lastUsed)
		fmt.Printf("   - Последнее использование: %s (%v назад)\n",
			lastUsed.Format(time.RFC3339), timeSinceLastUse)

		if timeSinceLastUse > 24*time.Hour {
			fmt.Printf("   - ⚠️ Webhook не использовался более 24 часов\n")
		}
	}
}

func demonstrateWebhookDeletion(ctx context.Context, client *nexus.Client, webhookID string) {
	fmt.Println("\n5️⃣ Удаление webhook")

	deleteResp, err := client.DeleteWebhook(ctx, webhookID)
	if err != nil {
		log.Fatalf("Webhook deletion failed: %v", err)
	}

	fmt.Printf("✅ Webhook удален:\n")
	fmt.Printf("   - Webhook ID: %s\n", deleteResp.WebhookID)
	fmt.Printf("   - Status: %s\n", deleteResp.Status)
	if deleteResp.Message != "" {
		fmt.Printf("   - Message: %s\n", deleteResp.Message)
	}
}
