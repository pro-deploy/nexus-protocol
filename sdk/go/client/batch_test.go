package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func TestExecuteBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/batch/execute" {
			t.Errorf("Expected POST /api/v1/batch/execute, got %s %s", r.Method, r.URL.Path)
		}

		var req types.BatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Проверяем, что запрос содержит операции
		if len(req.Operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(req.Operations))
		}

		// Возвращаем успешный ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := types.BatchResponse{
			Results: []types.BatchResult{
				{
					OperationID:   0,
					Success:       true,
					Data:          map[string]interface{}{"result": "success"},
					ExecutionTimeMS: 100,
				},
				{
					OperationID:   1,
					Success:       true,
					Data:          map[string]interface{}{"result": "success"},
					ExecutionTimeMS: 150,
				},
			},
			Total:     2,
			Successful: 2,
			Failed:    0,
			TotalTimeMS: 250,
			ResponseMetadata: &types.ResponseMetadata{
				RequestID:      "req-123",
				ProtocolVersion: "2.0.0",
				ServerVersion:   "2.0.0",
				Timestamp:       1640995200,
				ProcessingTimeMS: 250,
			},
		}
		json.NewEncoder(w).Encode(map[string]types.BatchResponse{"data": response})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})

	// Создаем batch с двумя операциями
	batch := NewBatchBuilder().
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "test query 1",
		}).
		AddOperation("execute_template", &types.ExecuteTemplateRequest{
			Query: "test query 2",
		}).
		SetOptions(&types.BatchOptions{
			Parallel: true,
		})

	result, err := batch.Execute(context.Background(), client)
	if err != nil {
		t.Fatalf("ExecuteBatch failed: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}

	if result.Successful != 2 {
		t.Errorf("Expected successful 2, got %d", result.Successful)
	}

	if result.Failed != 0 {
		t.Errorf("Expected failed 0, got %d", result.Failed)
	}

	if result.TotalTimeMS != 250 {
		t.Errorf("Expected total time 250ms, got %d", result.TotalTimeMS)
	}
}

func TestExecuteBatch_PartialFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := types.BatchResponse{
			Results: []types.BatchResult{
				{
					OperationID:   0,
					Success:       true,
					Data:          map[string]interface{}{"result": "success"},
					ExecutionTimeMS: 100,
				},
				{
					OperationID:   1,
					Success:       false,
					Error:         &types.ErrorDetail{Code: "VALIDATION_FAILED", Type: "VALIDATION_ERROR", Message: "Invalid query"},
					ExecutionTimeMS: 50,
				},
			},
			Total:     2,
			Successful: 1,
			Failed:    1,
			TotalTimeMS: 150,
		}
		json.NewEncoder(w).Encode(map[string]types.BatchResponse{"data": response})
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})

	batch := NewBatchBuilder().
		AddOperation("execute_template", &types.ExecuteTemplateRequest{Query: "valid query"}).
		AddOperation("execute_template", &types.ExecuteTemplateRequest{Query: "invalid query"})

	result, err := batch.Execute(context.Background(), client)
	if err != nil {
		t.Fatalf("ExecuteBatch failed: %v", err)
	}

	if result.Successful != 1 {
		t.Errorf("Expected successful 1, got %d", result.Successful)
	}

	if result.Failed != 1 {
		t.Errorf("Expected failed 1, got %d", result.Failed)
	}

	if result.Results[1].Success {
		t.Error("Expected second operation to fail")
	}

	if result.Results[1].Error == nil {
		t.Error("Expected error for failed operation")
	}
}

func TestBatchBuilder(t *testing.T) {
	builder := NewBatchBuilder()

	// Добавляем операции
	builder.AddOperation("execute_template", &types.ExecuteTemplateRequest{Query: "test1"})
	builder.AddOperation("log_event", &types.LogEventRequest{EventType: "test"})

	// Устанавливаем опции
	builder.SetOptions(&types.BatchOptions{Parallel: true, StopOnError: false})

	req := builder.Build()

	if len(req.Operations) != 2 {
		t.Errorf("Expected 2 operations, got %d", len(req.Operations))
	}

	if req.Options == nil || !req.Options.Parallel {
		t.Error("Expected parallel execution option to be set")
	}

	if req.Options.StopOnError {
		t.Error("Expected StopOnError to be false")
	}
}
