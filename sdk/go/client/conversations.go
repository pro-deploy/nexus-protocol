package client

import (
	"context"
	"fmt"

	"github.com/nexus-protocol/go-sdk/types"
)

// CreateConversation создает новую беседу с AI.
// Требует валидный JWT токен.
func (c *Client) CreateConversation(ctx context.Context, req *types.CreateConversationRequest) (*types.Conversation, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1Conversations, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.Conversation `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetConversation получает беседу по ID вместе с сообщениями.
// conversationID должен быть валидным UUID.
func (c *Client) GetConversation(ctx context.Context, conversationID string) (*types.Conversation, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1Conversations, conversationID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.Conversation `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// SendMessage отправляет сообщение в беседу и получает ответ от AI.
// Тип сообщения по умолчанию: "text", если не указан.
func (c *Client) SendMessage(ctx context.Context, conversationID string, req *types.SendMessageRequest) (*types.MessageResponse, error) {
	// Устанавливаем тип сообщения по умолчанию
	if req.MessageType == "" {
		req.MessageType = "text"
	}

	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	path := fmt.Sprintf("%s/%s/messages", PathAPIV1Conversations, conversationID)
	resp, err := c.doRequest(ctx, "POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.MessageResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

