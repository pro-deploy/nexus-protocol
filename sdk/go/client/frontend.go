package client

import (
	"context"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// GetFrontendConfig получает активную конфигурацию фронтенда.
// Это публичный endpoint, который не требует аутентификации.
//
// Application Protocol: Ответ приходит в формате Application Protocol:
//
//	{
//	  "metadata": { ... ResponseMetadata ... },
//	  "data": { ... FrontendConfig ... }
//	}
func (c *Client) GetFrontendConfig(ctx context.Context) (*types.FrontendConfig, error) {
	resp, err := c.doRequest(ctx, "GET", PathAPIV1FrontendConfig, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data     types.FrontendConfig    `json:"data"`
		Metadata *types.ResponseMetadata `json:"metadata"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
