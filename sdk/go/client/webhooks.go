package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/nexus-protocol/go-sdk/types"
)

// RegisterWebhook регистрирует новый webhook для получения уведомлений об асинхронных операциях.
//
// Пример использования:
//
//	config := &types.WebhookConfig{
//		URL:    "https://my-app.com/webhook",
//		Events: []string{"template.completed", "template.failed"},
//		Secret: "webhook-secret-123",
//		RetryPolicy: &types.WebhookRetryPolicy{
//			MaxRetries:   3,
//			InitialDelay: 1000,
//		},
//	}
//
//	resp, err := client.RegisterWebhook(ctx, &types.RegisterWebhookRequest{
//		Config: config,
//	})
func (c *Client) RegisterWebhook(ctx context.Context, req *types.RegisterWebhookRequest) (*types.RegisterWebhookResponse, error) {
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1Webhooks, req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data types.RegisterWebhookResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// ListWebhooks получает список зарегистрированных webhooks.
func (c *Client) ListWebhooks(ctx context.Context, req *types.ListWebhooksRequest) (*types.ListWebhooksResponse, error) {
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	// Строим query параметры
	params := make([]string, 0)
	if req.ActiveOnly {
		params = append(params, "active_only=true")
	}
	if req.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", req.Limit))
	}
	if req.Offset > 0 {
		params = append(params, fmt.Sprintf("offset=%d", req.Offset))
	}

	path := PathAPIV1Webhooks
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data types.ListWebhooksResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// DeleteWebhook удаляет webhook по ID.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) (*types.DeleteWebhookResponse, error) {
	req := &types.DeleteWebhookRequest{
		WebhookID: webhookID,
		Metadata:  c.createRequestMetadata(),
	}

	resp, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("%s/%s", PathAPIV1Webhooks, webhookID), req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data types.DeleteWebhookResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// TestWebhook отправляет тестовое событие на webhook.
func (c *Client) TestWebhook(ctx context.Context, req *types.TestWebhookRequest) (*types.TestWebhookResponse, error) {
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", fmt.Sprintf("%s/%s/test", PathAPIV1Webhooks, req.WebhookID), req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data types.TestWebhookResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
