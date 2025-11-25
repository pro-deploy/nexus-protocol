package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// TimingInterceptor измеряет время выполнения запросов
type TimingInterceptor struct{}

func (t *TimingInterceptor) BeforeRequest(ctx context.Context, req *http.Request) error {
	req.Header.Set("X-Request-Start-Time", fmt.Sprintf("%d", time.Now().UnixNano()))
	return nil
}

func (t *TimingInterceptor) AfterResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	startTimeStr := req.Header.Get("X-Request-Start-Time")
	if startTimeStr != "" {
		var startTime int64
		fmt.Sscanf(startTimeStr, "%d", &startTime)
		duration := time.Since(time.Unix(0, startTime))
		fmt.Printf("Request to %s took %v\n", req.URL.Path, duration)
	}
	return nil
}

// AuthInterceptor добавляет дополнительные заголовки
type AuthInterceptor struct{}

func (a *AuthInterceptor) BeforeRequest(ctx context.Context, req *http.Request) error {
	req.Header.Set("X-Client-Version", "2.0.0")
	req.Header.Set("X-Request-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	return nil
}

func (a *AuthInterceptor) AfterResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	// Можно обработать ответ, например, обновить токен
	return nil
}

func main() {
	cfg := nexus.Config{
		BaseURL: "http://localhost:8080",
		Token:   "your-jwt-token",
	}

	client := nexus.NewClient(cfg)

	// Добавляем interceptors
	client.AddInterceptor(&TimingInterceptor{})
	client.AddInterceptor(&AuthInterceptor{})

	ctx := context.Background()

	fmt.Println("Выполнение запроса с interceptors...")

	req := &types.ExecuteTemplateRequest{
		Query:    "хочу борщ",
		Language: "ru",
	}

	result, err := client.ExecuteTemplate(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Printf("✓ Запрос выполнен\n")
	fmt.Printf("  Execution ID: %s\n", result.ExecutionID)
}

