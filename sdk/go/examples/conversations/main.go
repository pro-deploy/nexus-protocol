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

	// Создание беседы
	fmt.Println("Создание новой беседы...")
	createReq := &types.CreateConversationRequest{
		Title:       "Обсуждение рецептов",
		BotID:       "bot-123",
		SystemPrompt: "Ты помощник по кулинарии",
	}

	conversation, err := client.CreateConversation(ctx, createReq)
	if err != nil {
		log.Fatalf("Ошибка создания беседы: %v", err)
	}

	fmt.Printf("✓ Беседа создана:\n")
	fmt.Printf("  ID: %s\n", conversation.ID)
	fmt.Printf("  Title: %s\n", conversation.Title)
	fmt.Printf("  Status: %s\n", conversation.Status)

	// Отправка сообщения
	fmt.Println("\nОтправка сообщения...")
	messageReq := &types.SendMessageRequest{
		Content:     "Расскажи рецепт борща",
		MessageType: "text",
	}

	messageResp, err := client.SendMessage(ctx, conversation.ID, messageReq)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	fmt.Printf("✓ Сообщение отправлено\n")
	if messageResp.UserMessage != nil {
		fmt.Printf("  User message: %s\n", messageResp.UserMessage.Content)
	}
	if messageResp.AIResponse != nil {
		fmt.Printf("  AI response: %s\n", messageResp.AIResponse.Content)
	}
	fmt.Printf("  Total messages: %d\n", messageResp.TotalMessages)

	// Получение беседы с историей
	fmt.Println("\nПолучение беседы с историей...")
	fullConversation, err := client.GetConversation(ctx, conversation.ID)
	if err != nil {
		log.Fatalf("Ошибка получения беседы: %v", err)
	}

	fmt.Printf("✓ Беседа получена:\n")
	fmt.Printf("  Message count: %d\n", fullConversation.MessageCount)
	if len(fullConversation.Messages) > 0 {
		fmt.Println("  Messages:")
		for i, msg := range fullConversation.Messages {
			fmt.Printf("    %d. [%s] %s: %s\n",
				i+1,
				msg.SenderType,
				msg.Type,
				msg.Content[:min(50, len(msg.Content))],
			)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

