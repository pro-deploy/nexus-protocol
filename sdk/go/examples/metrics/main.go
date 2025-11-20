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
	// Создаем коллектор метрик
	metricsCollector := nexus.NewSimpleMetricsCollector()

	// Создаем interceptor для метрик
	metricsInterceptor := nexus.NewMetricsInterceptor(metricsCollector)

	cfg := nexus.Config{
		BaseURL: "http://localhost:8080",
		Token:   "your-jwt-token",
	}

	client := nexus.NewClient(cfg)

	// Добавляем interceptor для метрик
	client.AddInterceptor(metricsInterceptor)

	ctx := context.Background()

	fmt.Println("Выполнение запросов с метриками...")

	// Выполняем несколько запросов
	for i := 0; i < 5; i++ {
		req := &types.ExecuteTemplateRequest{
			Query:    fmt.Sprintf("запрос %d", i+1),
			Language: "ru",
		}

		_, err := client.ExecuteTemplate(ctx, req)
		if err != nil {
			log.Printf("Ошибка запроса %d: %v", i+1, err)
		} else {
			fmt.Printf("✓ Запрос %d выполнен\n", i+1)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Получаем статистику
	stats := metricsCollector.GetStats()
	
	fmt.Println("\n📊 Статистика метрик:")
	fmt.Printf("  Запросы: %v\n", stats["requests"])
	fmt.Printf("  Ошибки: %v\n", stats["errors"])
	fmt.Printf("  Средние длительности: %v\n", stats["avg_durations"])
}

