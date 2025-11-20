package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nexus-protocol/server/internal/ai"
	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
)

// Service handles batch operations
type Service struct {
	aiService *ai.Service
	logger    *zap.Logger
	maxConcurrency int
	maxBatchSize   int
}

// BatchJob represents a single batch job
type BatchJob struct {
	ID           string                          `json:"id"`
	Status       string                          `json:"status"`
	Operations   []*BatchOperation               `json:"operations"`
	Results      []*types.ExecuteTemplateResponse `json:"results,omitempty"`
	CreatedAt    time.Time                       `json:"created_at"`
	CompletedAt  *time.Time                      `json:"completed_at,omitempty"`
	Error        string                          `json:"error,omitempty"`
	TotalTimeMS  int64                           `json:"total_time_ms,omitempty"`
}

// BatchOperation represents a single operation in a batch
type BatchOperation struct {
	ID       string                      `json:"id"`
	Type     string                      `json:"type"`
	Request  *types.ExecuteTemplateRequest `json:"request"`
	Response *types.ExecuteTemplateResponse `json:"response,omitempty"`
	Error    string                      `json:"error,omitempty"`
	Status   string                      `json:"status"`
	StartTime *time.Time                 `json:"start_time,omitempty"`
	EndTime   *time.Time                 `json:"end_time,omitempty"`
}

// NewService creates a new batch service
func NewService(aiService *ai.Service, logger *zap.Logger, maxConcurrency, maxBatchSize int) *Service {
	return &Service{
		aiService:      aiService,
		logger:         logger,
		maxConcurrency: maxConcurrency,
		maxBatchSize:   maxBatchSize,
	}
}

// ExecuteBatch executes multiple operations in batch
func (s *Service) ExecuteBatch(ctx context.Context, requests []*types.ExecuteTemplateRequest, parallel bool) (*BatchJob, error) {
	if len(requests) > s.maxBatchSize {
		return nil, fmt.Errorf("batch size exceeds maximum allowed (%d)", s.maxBatchSize)
	}

	jobID := uuid.New().String()
	job := &BatchJob{
		ID:         jobID,
		Status:     "running",
		CreatedAt:  time.Now(),
		Operations: make([]*BatchOperation, len(requests)),
		Results:    make([]*types.ExecuteTemplateResponse, 0, len(requests)),
	}

	s.logger.Info("Starting batch job",
		zap.String("job_id", jobID),
		zap.Int("operation_count", len(requests)),
		zap.Bool("parallel", parallel))

	// Create operations
	for i, req := range requests {
		job.Operations[i] = &BatchOperation{
			ID:      uuid.New().String(),
			Type:    "execute_template",
			Request: req,
			Status:  "pending",
		}
	}

	// Execute operations
	if parallel {
		s.executeParallel(ctx, job)
	} else {
		s.executeSequential(ctx, job)
	}

	// Mark job as completed
	job.Status = "completed"
	now := time.Now()
	job.CompletedAt = &now
	job.TotalTimeMS = now.Sub(job.CreatedAt).Milliseconds()

	s.logger.Info("Batch job completed",
		zap.String("job_id", jobID),
		zap.Int64("total_time_ms", job.TotalTimeMS),
		zap.Int("successful_operations", len(job.Results)))

	return job, nil
}

// executeParallel executes operations in parallel with concurrency control
func (s *Service) executeParallel(ctx context.Context, job *BatchJob) {
	semaphore := make(chan struct{}, s.maxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, operation := range job.Operations {
		wg.Add(1)
		go func(op *BatchOperation, index int) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			s.executeOperation(ctx, op, index, &mu, job)
		}(operation, i)
	}

	wg.Wait()
}

// executeSequential executes operations one by one
func (s *Service) executeSequential(ctx context.Context, job *BatchJob) {
	var mu sync.Mutex

	for i, operation := range job.Operations {
		s.executeOperation(ctx, operation, i, &mu, job)
	}
}

// executeOperation executes a single operation
func (s *Service) executeOperation(ctx context.Context, operation *BatchOperation, index int, mu *sync.Mutex, job *BatchJob) {
	startTime := time.Now()
	operation.StartTime = &startTime
	operation.Status = "running"

	s.logger.Debug("Executing batch operation",
		zap.String("operation_id", operation.ID),
		zap.String("type", operation.Type),
		zap.String("query", operation.Request.Query))

	// Execute the operation
	response, err := s.aiService.ExecuteTemplate(ctx, operation.Request)

	endTime := time.Now()
	operation.EndTime = &endTime

	mu.Lock()
	defer mu.Unlock()

	if err != nil {
		operation.Status = "failed"
		operation.Error = err.Error()
		s.logger.Error("Batch operation failed",
			zap.String("operation_id", operation.ID),
			zap.Error(err))
		return
	}

	operation.Status = "completed"
	operation.Response = response
	job.Results = append(job.Results, response)

	s.logger.Debug("Batch operation completed",
		zap.String("operation_id", operation.ID),
		zap.Int64("duration_ms", endTime.Sub(startTime).Milliseconds()))
}

// GetBatchStatus returns the status of a batch job
func (s *Service) GetBatchStatus(ctx context.Context, jobID string) (*BatchJob, error) {
	// In a real implementation, this would retrieve from database/cache
	// For now, return a mock response
	return &BatchJob{
		ID:        jobID,
		Status:    "completed",
		CreatedAt: time.Now().Add(-time.Minute),
		Results:   []*types.ExecuteTemplateResponse{},
	}, nil
}

// CancelBatch cancels a running batch job
func (s *Service) CancelBatch(ctx context.Context, jobID string) error {
	// In a real implementation, this would cancel running operations
	// For now, just log the cancellation
	s.logger.Info("Batch job cancelled", zap.String("job_id", jobID))
	return nil
}

// GetBatchStats returns statistics about batch operations
func (s *Service) GetBatchStats(ctx context.Context) (*BatchStats, error) {
	// Mock statistics - in real implementation would query database
	return &BatchStats{
		TotalJobs:         150,
		RunningJobs:       5,
		CompletedJobs:     140,
		FailedJobs:        5,
		AverageJobTimeMS:  2500,
		TotalOperations:   2500,
		SuccessfulOperations: 2400,
		FailedOperations:  100,
	}, nil
}

// BatchStats represents batch operation statistics
type BatchStats struct {
	TotalJobs            int64 `json:"total_jobs"`
	RunningJobs          int64 `json:"running_jobs"`
	CompletedJobs        int64 `json:"completed_jobs"`
	FailedJobs           int64 `json:"failed_jobs"`
	AverageJobTimeMS     int64 `json:"average_job_time_ms"`
	TotalOperations      int64 `json:"total_operations"`
	SuccessfulOperations int64 `json:"successful_operations"`
	FailedOperations     int64 `json:"failed_operations"`
}

// ValidateBatchRequest validates a batch request
func (s *Service) ValidateBatchRequest(requests []*types.ExecuteTemplateRequest) error {
	if len(requests) == 0 {
		return fmt.Errorf("batch must contain at least one operation")
	}

	if len(requests) > s.maxBatchSize {
		return fmt.Errorf("batch size %d exceeds maximum allowed %d", len(requests), s.maxBatchSize)
	}

	for i, req := range requests {
		if req == nil {
			return fmt.Errorf("operation %d: request cannot be nil", i)
		}
		if req.Query == "" {
			return fmt.Errorf("operation %d: query cannot be empty", i)
		}
		if req.Language != "ru" && req.Language != "en" {
			return fmt.Errorf("operation %d: invalid language '%s'", i, req.Language)
		}
	}

	return nil
}

// GetBatchOperations returns operations for a batch job
func (s *Service) GetBatchOperations(ctx context.Context, jobID string, limit, offset int) ([]*BatchOperation, error) {
	// Mock implementation - in real system would query database
	operations := []*BatchOperation{
		{
			ID:     uuid.New().String(),
			Type:   "execute_template",
			Status: "completed",
			StartTime: &time.Time{},
			EndTime:   &time.Time{},
		},
	}

	return operations, nil
}
