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
		Token:   "invalid-token",
	}

	nexusClient := nexus.NewClient(cfg)
	ctx := context.Background()

	// Пример обработки ошибок
	req := &types.ExecuteTemplateRequest{
		Query: "", // Пустой запрос для демонстрации ошибки валидации
	}

	result, err := nexusClient.ExecuteTemplate(ctx, req)
	if err != nil {
		// Проверяем тип ошибки
		if errDetail, ok := err.(*types.ErrorDetail); ok {
			fmt.Printf("Ошибка протокола:\n")
			fmt.Printf("  Code: %s\n", errDetail.Code)
			fmt.Printf("  Type: %s\n", errDetail.Type)
			fmt.Printf("  Message: %s\n", errDetail.Message)

			// Обрабатываем разные типы ошибок
			switch {
			case errDetail.IsValidationError():
				fmt.Println("  → Это ошибка валидации")
				if errDetail.Field != "" {
					fmt.Printf("  → Проблемное поле: %s\n", errDetail.Field)
				}
			case errDetail.IsAuthenticationError():
				fmt.Println("  → Это ошибка аутентификации")
				fmt.Println("  → Проверьте токен")
			case errDetail.IsAuthorizationError():
				fmt.Println("  → Это ошибка авторизации")
				fmt.Println("  → У вас недостаточно прав")
			case errDetail.IsRateLimitError():
				fmt.Println("  → Превышен лимит запросов")
			case errDetail.IsInternalError():
				fmt.Println("  → Внутренняя ошибка сервера")
			}
		} else {
			log.Printf("Неожиданная ошибка: %v", err)
		}
		return
	}

	fmt.Printf("Результат: %+v\n", result)
}

