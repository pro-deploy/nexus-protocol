package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pro-deploy/nexus-protocol/sdk/go/client"
	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
	// Создание клиента
	config := client.Config{
		BaseURL: "http://localhost:8080", // URL вашего Nexus сервера
		Token:   "your-jwt-token",        // Ваш JWT токен
	}

	nexusClient := client.NewClient(config)
	ctx := context.Background()

	fmt.Println("🚀 Пример работы с интеграциями через Admin API")

	// Получаем Admin клиент
	adminClient := nexusClient.Admin()

	// 1. Получение списка всех интеграций
	fmt.Println("\n📋 Получение списка интеграций...")
	integrations, err := adminClient.ListIntegrations(ctx, "")
	if err != nil {
		log.Printf("Failed to list integrations: %v", err)
	} else {
		fmt.Printf("Найдено интеграций: %d\n", len(integrations))
		for _, integration := range integrations {
			fmt.Printf("- %s (%s) - %s [%s]\n",
				integration.Name,
				integration.ID,
				integration.Type,
				integration.Provider)
			if !integration.Enabled {
				fmt.Printf("  ⚠️  Отключена\n")
			}
		}
	}

	// 2. Фильтрация по типу
	fmt.Println("\n🔍 Фильтрация интеграций по типу 'data_source'...")
	dataSources, err := adminClient.ListIntegrations(ctx, "data_source")
	if err != nil {
		log.Printf("Failed to list data sources: %v", err)
	} else {
		fmt.Printf("Найдено источников данных: %d\n", len(dataSources))
		for _, ds := range dataSources {
			fmt.Printf("- %s: %s\n", ds.Name, ds.Provider)
		}
	}

	// 3. Получение конкретной интеграции
	if len(integrations) > 0 {
		fmt.Println("\n📥 Получение деталей интеграции...")
		integration, err := adminClient.GetIntegration(ctx, integrations[0].ID)
		if err != nil {
			log.Printf("Failed to get integration: %v", err)
		} else {
			fmt.Printf("Интеграция: %s\n", integration.Name)
			fmt.Printf("  ID: %s\n", integration.ID)
			fmt.Printf("  Тип: %s\n", integration.Type)
			fmt.Printf("  Провайдер: %s\n", integration.Provider)
			fmt.Printf("  Включена: %v\n", integration.Enabled)
			if len(integration.Metadata) > 0 {
				fmt.Printf("  Метаданные: %+v\n", integration.Metadata)
			}
		}
	}

	// 4. Создание новой интеграции
	fmt.Println("\n➕ Создание новой интеграции...")
	newIntegration := &types.IntegrationConfig{
		ID:       "example-weather-api",
		Name:     "Example Weather API",
		Type:     "data_source",
		Provider: "openweather",
		Enabled:  true,
		Config: map[string]interface{}{
			"base_url": "https://api.openweathermap.org/data/2.5",
			"timeout":  30,
		},
		Credentials: map[string]string{
			"api_key": "your_api_key_here",
		},
		Metadata: map[string]string{
			"version": "1.0.0",
			"region":  "global",
		},
	}

	created, err := adminClient.CreateIntegration(ctx, newIntegration)
	if err != nil {
		log.Printf("Failed to create integration: %v", err)
	} else {
		fmt.Printf("✅ Создана интеграция: %s (ID: %s)\n", created.Name, created.ID)
	}

	// 5. Обновление интеграции
	if created != nil {
		fmt.Println("\n✏️  Обновление интеграции...")
		created.Metadata["updated"] = "true"
		updated, err := adminClient.UpdateIntegration(ctx, created.ID, created)
		if err != nil {
			log.Printf("Failed to update integration: %v", err)
		} else {
			fmt.Printf("✅ Обновлена интеграция: %s\n", updated.Name)
		}
	}

	// 6. Использование интеграций через домен
	fmt.Println("\n🎯 Использование интеграций через домен 'integrations'...")
	fmt.Println("   (Используйте ExecuteTemplate с доменом 'integrations')")
	
	// Пример запроса через домен integrations
	templateReq := &types.ExecuteTemplateRequest{
		Query:    "получить данные из weather-api",
		Language: "ru",
		Context: &types.UserContext{
			UserID: "example-user",
		},
	}

	// Выполнение через обычный API (домен integrations обработает запрос)
	fmt.Println("   Запрос отправлен в домен 'integrations'")
	fmt.Printf("   Query: %s\n", templateReq.Query)
	fmt.Println("   Домен integrations автоматически выберет нужный источник данных")

	// 7. Удаление тестовой интеграции (опционально)
	if created != nil {
		fmt.Println("\n🗑️  Удаление тестовой интеграции...")
		if err := adminClient.DeleteIntegration(ctx, created.ID); err != nil {
			log.Printf("Failed to delete integration: %v", err)
		} else {
			fmt.Printf("✅ Удалена интеграция: %s\n", created.ID)
		}
	}

	fmt.Println("\n✅ Пример завершен!")
	fmt.Println("\n💡 Важно:")
	fmt.Println("   - MCP - это внутренний протокол сервера")
	fmt.Println("   - Пользователи SDK работают только с публичным Admin API")
	fmt.Println("   - Доступ к данным интеграций - через домен 'integrations'")
	fmt.Println("   - Управление интеграциями - через Admin API")
}