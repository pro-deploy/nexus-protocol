package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nexus-protocol/go-sdk/types"
)

// LogEvent логирует событие аналитики для отслеживания.
// Требует валидный JWT токен.
func (c *Client) LogEvent(ctx context.Context, req *types.LogEventRequest) (*types.LogEventResponse, error) {
	resp, err := c.doRequest(ctx, "POST", PathAPIV1AnalyticsEvents, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.LogEventResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetEvents получает события аналитики с фильтрацией по типу события, пользователю и пагинацией.
// Поддерживает фильтрацию по event_type, user_id, limit и offset.
func (c *Client) GetEvents(ctx context.Context, req *types.GetEventsRequest) (*types.GetEventsResponse, error) {
	// Строим query параметры
	params := url.Values{}
	if req.EventType != "" {
		params.Add("event_type", req.EventType)
	}
	if req.UserID != "" {
		params.Add("user_id", req.UserID)
	}
	if req.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", req.Limit))
	}
	if req.Offset > 0 {
		params.Add("offset", fmt.Sprintf("%d", req.Offset))
	}

	path := PathAPIV1AnalyticsEvents
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.GetEventsResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetStats получает комплексную статистику аналитики.
// Поддерживает фильтрацию по user_id, tenant_id и периоду (days).
// По умолчанию период: 7 дней.
func (c *Client) GetStats(ctx context.Context, req *types.GetStatsRequest) (*types.AnalyticsStats, error) {
	// Строим query параметры
	params := url.Values{}
	if req.UserID != "" {
		params.Add("user_id", req.UserID)
	}
	if req.TenantID != "" {
		params.Add("tenant_id", req.TenantID)
	}
	if req.Days > 0 {
		params.Add("days", fmt.Sprintf("%d", req.Days))
	}

	path := PathAPIV1AnalyticsStats
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data types.AnalyticsStats `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

