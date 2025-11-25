package main

import (
	"context"
	"fmt"
	"log"
	"time"

	nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
	// Конфигурация для расширенных возможностей (v2.0.0)
	cfg := nexus.Config{
		BaseURL:         "http://localhost:8080",
		Token:           "jwt-token",
		ProtocolVersion: "2.0.0", // Nexus Protocol v2.0.0 с расширенными возможностями
		ClientVersion:   "2.0.0",
		ClientID:        "advanced-app",
		ClientType:      "api",
		RetryConfig: &nexus.RetryConfig{
			MaxRetries:        5,
			InitialDelay:      200 * time.Millisecond,
			MaxDelay:          10 * time.Second,
			BackoffMultiplier: 2.0,
		},
	}

	client := nexus.NewClient(cfg)
	ctx := context.Background()

	fmt.Println("🚀 Nexus Protocol Advanced Features Demo (v2.0.0)")
	fmt.Println("=================================")

	// 1. Настройка enterprise параметров
	demonstrateEnterpriseSetup(client)

	// 2. Расширенный поиск с фильтрами и пагинацией
	demonstrateAdvancedSearch(ctx, client)

	// 3. Batch операции для высокой производительности
	demonstrateBatchOperations(ctx, client)

	// 4. Webhooks для асинхронной обработки
	demonstrateWebhooks(ctx, client)

	// 5. Enterprise аналитика и метрики
	demonstrateAnalytics(ctx, client)

	// 6. Детальный health check
	demonstrateHealthCheck(ctx, client)

	fmt.Println("\n✅ Enterprise demo завершен!")
}

func demonstrateEnterpriseSetup(client *nexus.Client) {
	fmt.Println("\n📋 1. Настройка enterprise параметров")

	// Настройка приоритетов и кэширования
	client.SetPriority("high")
	client.SetCacheControl("cache-first")
	client.SetCacheTTL(300) // 5 минут кэша
	client.SetExperiment("enterprise-demo")
	client.SetFeatureFlag("advanced_filters", "enabled")
	client.SetFeatureFlag("batch_operations", "enabled")

	fmt.Println("✅ Настроены enterprise параметры:")
	fmt.Println("   - Приоритет: high")
	fmt.Println("   - Кэширование: cache-first (TTL: 300s)")
	fmt.Println("   - Эксперимент: enterprise-demo")
	fmt.Println("   - Фичи: advanced_filters, batch_operations")
}

func demonstrateAdvancedSearch(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n🔍 2. Расширенный поиск с enterprise фильтрами")

	// Запрос с расширенными фильтрами и контекстом
	req := &types.ExecuteTemplateRequest{
		Query: "умный дом с голосовым управлением в Москве",
		Language: "ru",
		Context: &types.UserContext{
			UserID:    "enterprise-user-123",
			SessionID: "enterprise-session-456",
			TenantID:  "enterprise-company-abc",
			Location: &types.UserLocation{
				Latitude:  55.7558,
				Longitude: 37.6173,
				Accuracy:  10.0,
			},
			Locale:   "ru-RU",
			Timezone: "Europe/Moscow",
			Currency: "RUB",
			Region:   "RU",
		},
		Options: &types.ExecuteOptions{
			TimeoutMS:           60000,
			MaxResultsPerDomain: 20,
			ParallelExecution:   true,
			IncludeWebSearch:    true,
		},
		Filters: &types.AdvancedFilters{
			Domains:       []string{"commerce", "smart_home", "reviews"},
			MinRelevance:  0.8,
			MaxResults:    50,
			SortBy:        "relevance",
			DateRange: &types.DateRange{
				From: time.Now().AddDate(0, 0, -30).Unix(), // Последние 30 дней
				To:   time.Now().Unix(),
			},
		},
	}

	result, err := client.ExecuteTemplate(ctx, req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return
	}

	fmt.Printf("✅ Поиск выполнен:\n")
	fmt.Printf("   - Execution ID: %s\n", result.ExecutionID)
	fmt.Printf("   - Status: %s\n", result.Status)
	fmt.Printf("   - Processing time: %d ms\n", result.ProcessingTimeMS)
	fmt.Printf("   - Total results: %d\n", len(result.Sections))

	// Показываем enterprise метрики из ResponseMetadata
	if result.ResponseMetadata != nil {
		fmt.Println("   - Enterprise metrics:")
		if result.ResponseMetadata.RateLimitInfo != nil {
			fmt.Printf("     * Rate limit: %d/%d (reset: %d)\n",
				result.ResponseMetadata.RateLimitInfo.Remaining,
				result.ResponseMetadata.RateLimitInfo.Limit,
				result.ResponseMetadata.RateLimitInfo.ResetAt)
		}
		if result.ResponseMetadata.CacheInfo != nil {
			fmt.Printf("     * Cache: %s (TTL: %ds)\n",
				map[bool]string{true: "hit", false: "miss"}[result.ResponseMetadata.CacheInfo.CacheHit],
				result.ResponseMetadata.CacheInfo.CacheTTL)
		}
		if result.ResponseMetadata.QuotaInfo != nil {
			fmt.Printf("     * Quota: %d/%d (%s)\n",
				result.ResponseMetadata.QuotaInfo.QuotaUsed,
				result.ResponseMetadata.QuotaInfo.QuotaLimit,
				result.ResponseMetadata.QuotaInfo.QuotaType)
		}
	}

	// Показываем пагинацию
	if result.Pagination != nil {
		fmt.Printf("   - Pagination: page %d/%d (%d items)\n",
			result.Pagination.Page,
			result.Pagination.TotalPages,
			result.Pagination.TotalItems)
	}
}

func demonstrateBatchOperations(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n📦 3. Batch операции для высокой производительности")

	// Создаем batch с несколькими операциями
	batch := nexus.NewBatchBuilder().
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "купить iPhone 15",
			Language: "ru",
			Context: &types.UserContext{
				UserID: "batch-user-1",
				TenantID: "enterprise-company-abc",
			},
		}).
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "забронировать отель в Париже",
			Language: "ru",
			Context: &types.UserContext{
				UserID: "batch-user-2",
				TenantID: "enterprise-company-abc",
			},
		}).
		AddOperation("log_event", &types.LogEventRequest{
			EventType: "batch_operation_demo",
			UserID:    "batch-user-1",
			TenantID:  "enterprise-company-abc",
			Data: map[string]interface{}{
				"operation": "batch_demo",
				"timestamp": time.Now().Unix(),
			},
		}).
		SetOptions(&types.BatchOptions{
			Parallel:    true,
			StopOnError: false,
		})

	// Выполняем batch
	batchResult, err := batch.Execute(ctx, client)
	if err != nil {
		log.Printf("Ошибка выполнения batch: %v", err)
		return
	}

	fmt.Printf("✅ Batch выполнен:\n")
	fmt.Printf("   - Всего операций: %d\n", batchResult.Total)
	fmt.Printf("   - Успешных: %d\n", batchResult.Successful)
	fmt.Printf("   - Неудачных: %d\n", batchResult.Failed)
	fmt.Printf("   - Общее время: %d ms\n", batchResult.TotalTimeMS)

	// Показываем результаты по операциям
	for i, res := range batchResult.Results {
		status := "✅"
		if !res.Success {
			status = "❌"
		}
		fmt.Printf("   %d. %s Операция #%d - %d ms\n",
			i+1, status, res.OperationID, res.ExecutionTimeMS)
	}
}

func demonstrateWebhooks(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n🪝 4. Webhooks для асинхронной обработки")

	// Регистрируем webhook
	webhookResp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
		Config: &types.WebhookConfig{
			URL:    "https://enterprise-app.company.com/webhooks/nexus",
			Events: []string{"template.completed", "template.failed", "batch.completed"},
			Secret: "enterprise-webhook-secret-2024",
			RetryPolicy: &types.WebhookRetryPolicy{
				MaxRetries:    3,
				InitialDelay:  1000,
				MaxDelay:      30000,
				BackoffFactor: 2.0,
			},
			Active:     true,
			Description: "Enterprise webhook for async operations",
		},
	})

	if err != nil {
		log.Printf("Ошибка регистрации webhook: %v", err)
		return
	}

	fmt.Printf("✅ Webhook зарегистрирован:\n")
	fmt.Printf("   - Webhook ID: %s\n", webhookResp.WebhookID)
	fmt.Printf("   - Status: %s\n", webhookResp.Status)

	// Получаем список webhooks
	webhooks, err := client.ListWebhooks(ctx, &types.ListWebhooksRequest{
		ActiveOnly: true,
		Limit:      10,
	})

	if err != nil {
		log.Printf("Ошибка получения webhooks: %v", err)
		return
	}

	fmt.Printf("   - Активных webhooks: %d\n", len(webhooks.Webhooks))
	for _, wh := range webhooks.Webhooks {
		fmt.Printf("     * %s: %s (%d/%d успехов/ошибок)\n",
			wh.ID, wh.Config.URL, wh.SuccessCount, wh.ErrorCount)
	}
}

func demonstrateAnalytics(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n📊 5. Enterprise аналитика и метрики")

	// Получаем расширенную аналитику
	stats, err := client.GetStats(ctx, &types.GetStatsRequest{
		TenantID: "enterprise-company-abc",
		Days:     7,
	})

	if err != nil {
		log.Printf("Ошибка получения статистики: %v", err)
		return
	}

	fmt.Printf("✅ Enterprise аналитика за %d дней:\n", stats.PeriodDays)
	fmt.Printf("   - Всего событий: %d\n", stats.TotalEvents)
	fmt.Printf("   - Всего пользователей: %d\n", stats.TotalUsers)
	fmt.Printf("   - Активных пользователей: %d\n", stats.ActiveUsers)

	// Показываем метрики конверсии
	if stats.ConversionMetrics != nil {
		fmt.Println("   - Метрики конверсии:")
		fmt.Printf("     * Поиск → Результат: %.1f%%\n", stats.ConversionMetrics.SearchToResult*100)
		fmt.Printf("     * Результат → Действие: %.1f%%\n", stats.ConversionMetrics.ResultToAction*100)
		fmt.Printf("     * Успешность шаблонов: %.1f%%\n", stats.ConversionMetrics.TemplateSuccess*100)
		fmt.Printf("     * Удержание пользователей: %.1f%%\n", stats.ConversionMetrics.UserRetention*100)
	}

	// Показываем метрики производительности
	if stats.PerformanceMetrics != nil {
		fmt.Println("   - Метрики производительности:")
		fmt.Printf("     * Среднее время ответа: %.0f ms\n", stats.PerformanceMetrics.AvgResponseTimeMS)
		fmt.Printf("     * 95-й перцентиль: %.0f ms\n", stats.PerformanceMetrics.P95ResponseTimeMS)
		fmt.Printf("     * 99-й перцентиль: %.0f ms\n", stats.PerformanceMetrics.P99ResponseTimeMS)
		fmt.Printf("     * Процент ошибок: %.2f%%\n", stats.PerformanceMetrics.ErrorRate*100)
		fmt.Printf("     * Пропускная способность: %d req/min\n", stats.PerformanceMetrics.ThroughputRPM)
	}

	// Показываем разбивку по доменам
	if len(stats.DomainBreakdown) > 0 {
		fmt.Println("   - Разбивка по доменам:")
		for domain, metrics := range stats.DomainBreakdown {
			fmt.Printf("     * %s: %d запросов, %.1f%% успех, %.0f ms среднее\n",
				domain,
				metrics.RequestsCount,
				metrics.SuccessRate*100,
				metrics.AvgResponseTimeMS)
		}
	}
}

func demonstrateHealthCheck(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n🏥 6. Детальный enterprise health check")

	// Проверяем здоровье системы
	health, err := client.Health(ctx)
	if err != nil {
		log.Printf("Ошибка проверки здоровья: %v", err)
		return
	}

	fmt.Printf("✅ Health check: %s (version: %s)\n", health.Status, health.Version)

	// Проверяем готовность с enterprise метриками
	ready, err := client.Ready(ctx)
	if err != nil {
		log.Printf("Ошибка проверки готовности: %v", err)
		return
	}

	fmt.Printf("✅ Readiness check: %s\n", ready.Status)
	fmt.Printf("   - Database: %s\n", ready.Checks.Database)
	fmt.Printf("   - Redis: %s\n", ready.Checks.Redis)
	fmt.Printf("   - AI Services: %s\n", ready.Checks.AIServices)

	// Показываем детальную информацию о компонентах
	if len(ready.Components) > 0 {
		fmt.Println("   - Детальный статус компонентов:")
		for name, component := range ready.Components {
			status := "✅"
			if component.Status != "healthy" {
				status = "⚠️"
			}
			fmt.Printf("     * %s %s: %s", status, name, component.Status)
			if component.LatencyMS > 0 {
				fmt.Printf(" (%d ms)", component.LatencyMS)
			}
			if component.Message != "" {
				fmt.Printf(" - %s", component.Message)
			}
			fmt.Println()
		}
	}

	// Показываем информацию о емкости
	if ready.Capacity != nil {
		fmt.Println("   - Информация о емкости:")
		fmt.Printf("     * Текущая нагрузка: %.1f%%\n", ready.Capacity.CurrentLoad*100)
		fmt.Printf("     * Максимальная емкость: %d req/sec\n", ready.Capacity.MaxCapacity)
		fmt.Printf("     * Доступная емкость: %d req/sec\n", ready.Capacity.AvailableCapacity)
		fmt.Printf("     * Размер очереди: %d\n", ready.Capacity.QueueSize)
		fmt.Printf("     * Активные соединения: %d\n", ready.Capacity.ActiveConnections)
	}
}
