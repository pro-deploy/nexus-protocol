package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// ExecuteTemplate выполняет контекстно-зависимый шаблон.
//
// Application Protocol: Запрос отправляется с автоматически добавленными метаданными
// в формате RequestMetadata. Ответ приходит в формате Application Protocol:
// {
//   "metadata": { ... ResponseMetadata ... },
//   "data": { ... ExecuteTemplateResponse ... }
// }
//
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
		req.Metadata = c.createRequestMetadata()
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

// GetWorkflowSteps возвращает шаги workflow из ответа, отсортированные по номеру шага
func (c *Client) GetWorkflowSteps(result *types.ExecuteTemplateResponse) []types.WorkflowStep {
	if result.Workflow == nil || len(result.Workflow.Steps) == 0 {
		return nil
	}
	
	// Сортируем шаги по номеру
	steps := make([]types.WorkflowStep, len(result.Workflow.Steps))
	copy(steps, result.Workflow.Steps)
	
	// Простая сортировка по номеру шага
	for i := 0; i < len(steps)-1; i++ {
		for j := i + 1; j < len(steps); j++ {
			if steps[i].Step > steps[j].Step {
				steps[i], steps[j] = steps[j], steps[i]
			}
		}
	}
	
	return steps
}

// GetNextWorkflowStep возвращает следующий шаг workflow, который готов к выполнению
// (все зависимости выполнены)
func (c *Client) GetNextWorkflowStep(result *types.ExecuteTemplateResponse) *types.WorkflowStep {
	if result.Workflow == nil {
		return nil
	}
	
	steps := c.GetWorkflowSteps(result)
	completedResults := make(map[string]bool)
	
	for _, step := range steps {
		// Проверяем, выполнены ли все зависимости
		allDepsCompleted := true
		for _, depID := range step.DependsOn {
			if !completedResults[depID] {
				allDepsCompleted = false
				break
			}
		}
		
		// Если зависимости выполнены и шаг еще не выполнен
		if allDepsCompleted && step.Status == "pending" {
			return &step
		}
		
		// Если шаг выполнен, добавляем его результат в список выполненных
		if step.Status == "completed" && step.ResultID != "" {
			completedResults[step.ResultID] = true
		}
	}
	
	return nil
}

// GetWorkflowStepByDomain возвращает шаги workflow для указанного домена
func (c *Client) GetWorkflowStepByDomain(result *types.ExecuteTemplateResponse, domain string) []types.WorkflowStep {
	if result.Workflow == nil {
		return nil
	}
	
	var domainSteps []types.WorkflowStep
	for _, step := range result.Workflow.Steps {
		if step.Domain == domain {
			domainSteps = append(domainSteps, step)
		}
	}
	
	return domainSteps
}

