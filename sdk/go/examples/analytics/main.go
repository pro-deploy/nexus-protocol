package main

import (
	"context"
	"fmt"
	"log"

	nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
	cfg := nexus.Config{
		BaseURL: "http://localhost:8080",
		Token:   "your-jwt-token",
	}

	client := nexus.NewClient(cfg)
	ctx := context.Background()

	// Логирование события
	fmt.Println("Логирование события...")
	logReq := &types.LogEventRequest{
		EventType: "user_action",
		UserID:    "user-123",
		Data: map[string]interface{}{
			"action": "viewed_page",
			"page":   "/recipes",
		},
	}

	logResp, err := client.LogEvent(ctx, logReq)
	if err != nil {
		log.Fatalf("Ошибка логирования события: %v", err)
	}

	fmt.Printf("✓ Событие залогировано:\n")
	fmt.Printf("  Event ID: %s\n", logResp.EventID)
	fmt.Printf("  Timestamp: %s\n", logResp.Timestamp)

	// Получение событий
	fmt.Println("\nПолучение событий...")
	getEventsReq := &types.GetEventsRequest{
		EventType: "user_action",
		UserID:    "user-123",
		Limit:     10,
		Offset:    0,
	}

	eventsResp, err := client.GetEvents(ctx, getEventsReq)
	if err != nil {
		log.Fatalf("Ошибка получения событий: %v", err)
	}

	fmt.Printf("✓ События получены:\n")
	fmt.Printf("  Total: %d\n", eventsResp.Total)
	fmt.Printf("  Returned: %d\n", len(eventsResp.Events))
	for i, event := range eventsResp.Events {
		if i >= 5 {
			break
		}
		fmt.Printf("    %d. %s at %s\n", i+1, event.EventType, event.Timestamp)
	}

	// Получение статистики
	fmt.Println("\nПолучение статистики...")
	statsReq := &types.GetStatsRequest{
		UserID: "user-123",
		Days:   7,
	}

	stats, err := client.GetStats(ctx, statsReq)
	if err != nil {
		log.Fatalf("Ошибка получения статистики: %v", err)
	}

	fmt.Printf("✓ Статистика получена:\n")
	fmt.Printf("  Period: %d days\n", stats.PeriodDays)
	fmt.Printf("  Total events: %d\n", stats.TotalEvents)
	fmt.Printf("  Total users: %d\n", stats.TotalUsers)
	fmt.Printf("  Active users: %d\n", stats.ActiveUsers)
	fmt.Printf("  Events today: %d\n", stats.EventsToday)

	if len(stats.TopEvents) > 0 {
		fmt.Println("  Top events:")
		for i, event := range stats.TopEvents {
			if i >= 5 {
				break
			}
			fmt.Printf("    %d. %s: %d (%.1f%%)\n",
				i+1,
				event.Event,
				event.Count,
				event.Percentage*100,
			)
		}
	}
}

