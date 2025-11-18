package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nexus-protocol/go-sdk/types"
)

// ExecuteTemplate выполняет контекстно-зависимый шаблон.
// Если метаданные не указаны в запросе, они создаются автоматически.
// Язык по умолчанию: "ru", если не указан.
//
// Пример использования:
//
//	ctx := context.Background()
//	req := &types.ExecuteTemplateRequest{
//		Query:    "хочу борщ",
//		Language: "ru",
//	}
//	result, err := client.ExecuteTemplate(ctx, req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Execution ID: %s\n", result.ExecutionID)
func (c *Client) ExecuteTemplate(ctx context.Context, req *types.ExecuteTemplateRequest) (*types.ExecuteTemplateResponse, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = types.NewRequestMetadata(c.protocolVersion, c.clientVersion)
		req.Metadata.ClientID = c.clientID
		req.Metadata.ClientType = c.clientType
	}

	// Устанавливаем язык по умолчанию
	if req.Language == "" {
		req.Language = "ru"
	}

	// Устанавливаем опции по умолчанию
	if req.Options == nil {
		req.Options = &types.ExecuteOptions{
			TimeoutMS:           30000,
			MaxResultsPerDomain: 5,
			ParallelExecution:   true,
			IncludeWebSearch:    true,
		}
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1TemplatesExecute, req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data     types.ExecuteTemplateResponse `json:"data"`
		Metadata *types.ResponseMetadata        `json:"metadata"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	// Сохраняем метаданные ответа
	result.Data.ResponseMetadata = result.Metadata

	return &result.Data, nil
}

// GetExecutionStatus получает статус выполнения шаблона по execution ID.
// executionID должен быть валидным UUID.
func (c *Client) GetExecutionStatus(ctx context.Context, executionID string) (*types.ExecuteTemplateResponse, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1TemplatesStatus, executionID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data     types.ExecuteTemplateResponse `json:"data"`
		Metadata *types.ResponseMetadata        `json:"metadata"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	result.Data.ResponseMetadata = result.Metadata
	return &result.Data, nil
}

// StreamTemplateResults получает поток результатов выполнения в реальном времени (Server-Sent Events).
// Возвращает http.Response, который нужно закрыть после использования.
// Для чтения событий используйте bufio.Scanner или аналогичные инструменты.
func (c *Client) StreamTemplateResults(ctx context.Context, executionID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1TemplatesStream, executionID)
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, c.parseResponse(resp, nil)
	}

	return resp, nil
}

