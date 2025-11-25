package main

import (
	"context"
	"fmt"
	"log"
	"time"

	nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// Пример использования batch операций для массовой обработки запросов
func main() {
	cfg := nexus.Config{
		BaseURL:         "http://localhost:8080",
		Token:           "enterprise-jwt-token",
		ProtocolVersion: "2.0.0", // Nexus Protocol v2.0.0
		ClientVersion:   "2.0.0",
		ClientID:        "enterprise-batch-processor",
		ClientType:      "api",
	}

	client := nexus.NewClient(cfg)
	ctx := context.Background()

	fmt.Println("📦 Batch Operations Demo")
	fmt.Println("=======================")

	// Пример 1: Параллельная обработка множественных запросов
	demonstrateParallelBatch(ctx, client)

	// Пример 2: Batch с обработкой ошибок
	demonstrateBatchWithErrors(ctx, client)

	// Пример 3: Batch с ограничением параллельности
	demonstrateBatchWithConcurrency(ctx, client)

	// Пример 4: Комбинированные операции (templates + analytics)
	demonstrateMixedBatch(ctx, client)
}

func demonstrateParallelBatch(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n1️⃣ Параллельная обработка множественных запросов")

	// Создаем batch с 10 запросами
	batch := nexus.NewBatchBuilder()

	queries := []string{
		"купить iPhone 15",
		"забронировать отель в Москве",
		"найти ресторан с итальянской кухней",
		"купить билеты в театр",
		"найти автосервис рядом",
		"заказать доставку еды",
		"найти фитнес-клуб",
		"купить подарок на день рождения",
		"найти стоматолога",
		"забронировать столик в ресторане",
	}

	for i, query := range queries {
		batch.AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: query,
			Context: &types.UserContext{
				UserID:   fmt.Sprintf("user-%d", i+1),
				TenantID: "enterprise-company-abc",
			},
		})
	}

	// Настраиваем параллельное выполнение
	batch.SetOptions(&types.BatchOptions{
		Parallel:      true,
		StopOnError:   false,
		MaxConcurrency: 5, // Максимум 5 параллельных операций
	})

	start := time.Now()
	result, err := batch.Execute(ctx, client)
	if err != nil {
		log.Fatalf("Batch execution failed: %v", err)
	}

	duration := time.Since(start)

	fmt.Printf("✅ Batch выполнен за %v\n", duration)
	fmt.Printf("   - Всего операций: %d\n", result.Total)
	fmt.Printf("   - Успешных: %d\n", result.Successful)
	fmt.Printf("   - Неудачных: %d\n", result.Failed)
	fmt.Printf("   - Общее время выполнения: %d ms\n", result.TotalTimeMS)
	fmt.Printf("   - Среднее время на операцию: %.0f ms\n", float64(result.TotalTimeMS)/float64(result.Total))

	// Показываем результаты по операциям
	for i, res := range result.Results {
		status := "✅"
		if !res.Success {
			status = "❌"
		}
		fmt.Printf("   %d. %s Операция #%d - %d ms\n",
			i+1, status, res.OperationID, res.ExecutionTimeMS)
	}
}

func demonstrateBatchWithErrors(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n2️⃣ Batch с обработкой ошибок")

	batch := nexus.NewBatchBuilder().
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "купить ноутбук", // Валидный запрос
			Context: &types.UserContext{TenantID: "enterprise-company-abc"},
		}).
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "", // Невалидный запрос (пустой)
			Context: &types.UserContext{TenantID: "enterprise-company-abc"},
		}).
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "забронировать отель", // Валидный запрос
			Context: &types.UserContext{TenantID: "enterprise-company-abc"},
		}).
		SetOptions(&types.BatchOptions{
			Parallel:    true,
			StopOnError: false, // Продолжать при ошибках
		})

	result, err := batch.Execute(ctx, client)
	if err != nil {
		log.Fatalf("Batch execution failed: %v", err)
	}

	fmt.Printf("✅ Batch выполнен с обработкой ошибок\n")
	fmt.Printf("   - Успешных: %d\n", result.Successful)
	fmt.Printf("   - Неудачных: %d\n", result.Failed)

	// Обрабатываем ошибки
	for _, res := range result.Results {
		if !res.Success && res.Error != nil {
			fmt.Printf("   ❌ Операция #%d: %s (%s)\n",
				res.OperationID, res.Error.Message, res.Error.Code)
		}
	}
}

func demonstrateBatchWithConcurrency(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n3️⃣ Batch с ограничением параллельности")

	// Создаем batch с 20 операциями
	batch := nexus.NewBatchBuilder()

	for i := 0; i < 20; i++ {
		batch.AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: fmt.Sprintf("запрос #%d", i+1),
			Context: &types.UserContext{
				UserID:   fmt.Sprintf("user-%d", i+1),
				TenantID: "enterprise-company-abc",
			},
		})
	}

	// Ограничиваем параллельность до 3 операций одновременно
	batch.SetOptions(&types.BatchOptions{
		Parallel:      true,
		MaxConcurrency: 3,
		StopOnError:   false,
	})

	start := time.Now()
	result, err := batch.Execute(ctx, client)
	if err != nil {
		log.Fatalf("Batch execution failed: %v", err)
	}
	duration := time.Since(start)

	fmt.Printf("✅ Batch выполнен с ограничением параллельности\n")
	fmt.Printf("   - Всего операций: %d\n", result.Total)
	fmt.Printf("   - Максимальная параллельность: 3\n")
	fmt.Printf("   - Время выполнения: %v\n", duration)
	fmt.Printf("   - Пропускная способность: %.1f ops/sec\n",
		float64(result.Total)/duration.Seconds())
}

func demonstrateMixedBatch(ctx context.Context, client *nexus.Client) {
	fmt.Println("\n4️⃣ Комбинированные операции (templates + analytics)")

	batch := nexus.NewBatchBuilder().
		// Выполняем несколько шаблонов
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "купить смартфон",
			Context: &types.UserContext{TenantID: "enterprise-company-abc"},
		}).
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "забронировать отель",
			Context: &types.UserContext{TenantID: "enterprise-company-abc"},
		}).
		// Логируем события аналитики
		AddOperation("log_event", &types.LogEventRequest{
			EventType: "batch_operation_started",
			TenantID:  "enterprise-company-abc",
			Data: map[string]interface{}{
				"batch_size": 2,
				"timestamp":  time.Now().Unix(),
			},
		}).
		AddOperation("log_event", &types.LogEventRequest{
			EventType: "batch_operation_completed",
			TenantID:  "enterprise-company-abc",
			Data: map[string]interface{}{
				"batch_size": 2,
				"timestamp":  time.Now().Unix(),
			},
		}).
		SetOptions(&types.BatchOptions{
			Parallel:    true,
			StopOnError: false,
		})

	result, err := batch.Execute(ctx, client)
	if err != nil {
		log.Fatalf("Mixed batch execution failed: %v", err)
	}

	fmt.Printf("✅ Комбинированный batch выполнен\n")
	fmt.Printf("   - Всего операций: %d\n", result.Total)
	fmt.Printf("   - Успешных: %d\n", result.Successful)
	fmt.Printf("   - Типы операций: execute_template, log_event\n")

	// Группируем результаты по типам
	templateCount := 0
	eventCount := 0

	for _, res := range result.Results {
		if res.Success {
			// Определяем тип операции по результату
			if res.Data != nil {
				if _, ok := res.Data.(map[string]interface{})["execution_id"]; ok {
					templateCount++
				} else {
					eventCount++
				}
			}
		}
	}

	fmt.Printf("   - Выполнено шаблонов: %d\n", templateCount)
	fmt.Printf("   - Залогировано событий: %d\n", eventCount)
}
