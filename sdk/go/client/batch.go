package client

import (
	"context"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// ExecuteBatch выполняет пакет операций согласно протоколу v2.0.0.
// Поддерживает выполнение нескольких template операций в одном запросе.
// Полезно для enterprise сценариев с множественными операциями.
//
// Пример использования:
//
//	req := &types.BatchRequest{
//		Requests: []*types.ExecuteTemplateRequest{
//			{
//				Query:    "хочу борщ",
//				Language: "ru",
//			},
//			{
//				Query:    "найди ресторан",
//				Language: "ru",
//			},
//		},
//		BatchOptions: &types.ExecuteOptions{
//			ParallelExecution: true,
//		},
//	}
//
//	result, err := client.ExecuteBatch(ctx, req)
func (c *Client) ExecuteBatch(ctx context.Context, req *types.BatchRequest) (*types.BatchResponse, error) {
	// Если метаданные не указаны, создаем их автоматически
	if req.Metadata == nil {
		req.Metadata = c.createRequestMetadata()
	}

	resp, err := c.doRequest(ctx, "POST", PathAPIV1BatchExecute, req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data     types.BatchResponse `json:"data"`
		Metadata *types.ResponseMetadata `json:"metadata,omitempty"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	// Устанавливаем ResponseMetadata если он есть
	if result.Metadata != nil {
		result.Data.ResponseMetadata = result.Metadata
	}

	return &result.Data, nil
}

// BatchBuilder помогает строить batch запросы согласно протоколу v2.0.0
type BatchBuilder struct {
	requests     []*types.ExecuteTemplateRequest
	batchOptions *types.ExecuteOptions
}

// NewBatchBuilder создает новый BatchBuilder
func NewBatchBuilder() *BatchBuilder {
	return &BatchBuilder{
		requests: make([]*types.ExecuteTemplateRequest, 0),
	}
}

// AddRequest добавляет ExecuteTemplateRequest в batch
func (b *BatchBuilder) AddRequest(req *types.ExecuteTemplateRequest) *BatchBuilder {
	b.requests = append(b.requests, req)
	return b
}

// SetBatchOptions устанавливает опции batch выполнения
func (b *BatchBuilder) SetBatchOptions(options *types.ExecuteOptions) *BatchBuilder {
	b.batchOptions = options
	return b
}

// Build создает BatchRequest согласно протоколу v2.0.0
func (b *BatchBuilder) Build() *types.BatchRequest {
	return &types.BatchRequest{
		Requests:     b.requests,
		BatchOptions: b.batchOptions,
	}
}

// Execute выполняет batch через клиент
func (b *BatchBuilder) Execute(ctx context.Context, client *Client) (*types.BatchResponse, error) {
	req := b.Build()
	return client.ExecuteBatch(ctx, req)
}
