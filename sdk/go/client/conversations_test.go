package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func TestCreateConversation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1Conversations {
			t.Errorf("Expected path %s, got %s", PathAPIV1Conversations, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"data": {
				"id": "conv-123",
				"user_id": "user-123",
				"title": "Test Conversation",
				"status": "active",
				"message_count": 0,
				"created_at": "2025-01-18T10:00:00Z"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.CreateConversationRequest{
		Title: "Test Conversation",
	}

	conv, err := client.CreateConversation(ctx, req)
	if err != nil {
		t.Fatalf("CreateConversation failed: %v", err)
	}

	if conv.ID != "conv-123" {
		t.Errorf("Expected ID 'conv-123', got %s", conv.ID)
	}

	if conv.Title != "Test Conversation" {
		t.Errorf("Expected title 'Test Conversation', got %s", conv.Title)
	}
}

func TestGetConversation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"id": "conv-123",
				"user_id": "user-123",
				"status": "active",
				"message_count": 2,
				"created_at": "2025-01-18T10:00:00Z",
				"messages": [
					{
						"id": "msg-1",
						"conversation_id": "conv-123",
						"sender_type": "user",
						"type": "text",
						"content": "Hello",
						"created_at": "2025-01-18T10:00:00Z"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	conv, err := client.GetConversation(ctx, "conv-123")
	if err != nil {
		t.Fatalf("GetConversation failed: %v", err)
	}

	if conv.ID != "conv-123" {
		t.Errorf("Expected ID 'conv-123', got %s", conv.ID)
	}

	if conv.MessageCount != 2 {
		t.Errorf("Expected message_count 2, got %d", conv.MessageCount)
	}

	if len(conv.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(conv.Messages))
	}
}

func TestSendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"conversation_id": "conv-123",
				"user_message": {
					"id": "msg-1",
					"content": "Hello",
					"sender_type": "user"
				},
				"ai_response": {
					"id": "msg-2",
					"content": "Hi there!",
					"sender_type": "assistant"
				},
				"total_messages": 2
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.SendMessageRequest{
		Content: "Hello",
	}

	resp, err := client.SendMessage(ctx, "conv-123", req)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if resp.ConversationID != "conv-123" {
		t.Errorf("Expected conversation_id 'conv-123', got %s", resp.ConversationID)
	}

	if resp.UserMessage == nil {
		t.Error("Expected user_message, got nil")
	}

	if resp.AIResponse == nil {
		t.Error("Expected ai_response, got nil")
	}
}

