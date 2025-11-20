package client

import (
	"context"

	"github.com/nexus-protocol/go-sdk/types"
)

// ExecuteBatch выполняет пакет операций.
// Поддерживает выполнение нескольких операций в одном запросе.
// Полезно для enterprise сценариев с множественными операциями.
//
// Пример использования:
//
//	req := &types.BatchRequest{
//		Operations: []types.BatchOperation{
//			{
//				ID:   1,
//				Type: "execute_template",
//				Request: &types.ExecuteTemplateRequest{
//					Query: "хочу борщ",
//				},
//			},
//			{
//				ID:   2,
//				Type: "log_event",
//				Request: &types.LogEventRequest{
//					EventType: "batch_operation",
//					Data:      map[string]interface{}{"batch_size": 2},
//				},
//			},
//		},
//		Options: &types.BatchOptions{
//			Parallel: true,
//			StopOnError: false,
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
		Data types.BatchResponse `json:"data"`
	}

	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// BatchBuilder помогает строить batch запросы
type BatchBuilder struct {
	operations []types.BatchOperation
	options    *types.BatchOptions
	nextID     int
}

// NewBatchBuilder создает новый BatchBuilder
func NewBatchBuilder() *BatchBuilder {
	return &BatchBuilder{
		operations: make([]types.BatchOperation, 0),
		nextID:     1,
	}
}

// AddOperation добавляет операцию в batch
func (b *BatchBuilder) AddOperation(operationType string, request interface{}) *BatchBuilder {
	b.operations = append(b.operations, types.BatchOperation{
		ID:      b.nextID,
		Type:    operationType,
		Request: request,
	})
	b.nextID++
	return b
}

// SetOptions устанавливает опции batch выполнения
func (b *BatchBuilder) SetOptions(options *types.BatchOptions) *BatchBuilder {
	b.options = options
	return b
}

// Build создает BatchRequest
func (b *BatchBuilder) Build() *types.BatchRequest {
	return &types.BatchRequest{
		Operations: b.operations,
		Options:    b.options,
	}
}

// Execute выполняет batch через клиент
func (b *BatchBuilder) Execute(ctx context.Context, client *Client) (*types.BatchResponse, error) {
	req := b.Build()
	return client.ExecuteBatch(ctx, req)
}
