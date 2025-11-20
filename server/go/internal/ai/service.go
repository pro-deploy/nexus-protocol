package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nexus-protocol/server/pkg/config"
	"github.com/nexus-protocol/server/pkg/types"
	"go.uber.org/zap"
)

// Service represents the AI service for context-aware templates
type Service struct {
	config *config.AIConfig
	logger *zap.Logger
	client *http.Client
}

// NewService creates a new AI service
func NewService(cfg *config.AIConfig, logger *zap.Logger) *Service {
	return &Service{
		config: cfg,
		logger: logger,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// ExecuteTemplate executes a context-aware template using AI
func (s *Service) ExecuteTemplate(ctx context.Context, req *types.ExecuteTemplateRequest) (*types.ExecuteTemplateResponse, error) {
	startTime := time.Now()

	executionID := uuid.New().String()
	intentID := uuid.New().String()

	s.logger.Info("Executing template",
		zap.String("execution_id", executionID),
		zap.String("query", req.Query),
		zap.String("language", req.Language))

	// Analyze intent and determine domains
	domains, err := s.analyzeIntent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("intent analysis failed: %w", err)
	}

	// Execute in parallel across domains
	results := s.executeParallel(ctx, req, domains)

	// Rank and filter results
	rankedResults := s.rankResults(results)

	// Create response
	response := &types.ExecuteTemplateResponse{
		ExecutionID:      executionID,
		IntentID:         intentID,
		Status:           s.determineOverallStatus(results),
		QueryType:        s.determineQueryType(req.Query),
		Sections:         s.buildSections(results),
		WebSearch:        nil, // TODO: implement web search
		Ranking:          s.buildRanking(rankedResults),
		ProcessingTimeMS: int32(time.Since(startTime).Milliseconds()),
		Metadata:         s.buildResponseMetadata(startTime),
	}

	return response, nil
}

// analyzeIntent analyzes the user query to determine relevant domains
func (s *Service) analyzeIntent(ctx context.Context, req *types.ExecuteTemplateRequest) ([]string, error) {
	// This is a simplified intent analysis
	// In production, this would use ML/NLP models

	query := strings.ToLower(req.Query)

	domains := []string{}

	// Domain detection based on keywords
	if s.containsAny(query, []string{"купить", "цена", "магазин", "товар", "заказать"}) {
		domains = append(domains, "commerce")
	}

	if s.containsAny(query, []string{"рецепт", "готовить", "еда", "кухня", "ингредиенты"}) {
		domains = append(domains, "recipes")
	}

	if s.containsAny(query, []string{"отель", "бронировать", "путешествие", "тур", "авиабилет"}) {
		domains = append(domains, "travel")
	}

	if s.containsAny(query, []string{"документ", "инструкция", "руководство", "справка"}) {
		domains = append(domains, "knowledge")
	}

	// Default to commerce if no specific domain detected
	if len(domains) == 0 {
		domains = append(domains, "commerce")
	}

	return domains, nil
}

// executeParallel executes the query across multiple domains in parallel
func (s *Service) executeParallel(ctx context.Context, req *types.ExecuteTemplateRequest, domains []string) map[string]*domainResult {
	results := make(map[string]*domainResult)
	resultsChan := make(chan *domainResult, len(domains))

	// Execute in parallel
	for _, domain := range domains {
		go func(domain string) {
			result := s.executeDomain(ctx, req, domain)
			resultsChan <- result
		}(domain)
	}

	// Collect results
	for i := 0; i < len(domains); i++ {
		result := <-resultsChan
		results[result.domain] = result
	}

	return results
}

// executeDomain executes the query for a specific domain
func (s *Service) executeDomain(ctx context.Context, req *types.ExecuteTemplateRequest, domain string) *domainResult {
	startTime := time.Now()

	// This is where you would integrate with actual AI models or APIs
	// For now, we'll simulate responses based on domain

	var results []types.ResultItem
	var err error

	switch domain {
	case "commerce":
		results, err = s.executeCommerceDomain(req)
	case "recipes":
		results, err = s.executeRecipesDomain(req)
	case "travel":
		results, err = s.executeTravelDomain(req)
	case "knowledge":
		results, err = s.executeKnowledgeDomain(req)
	default:
		results, err = s.executeGenericDomain(req, domain)
	}

	status := "success"
	if err != nil {
		status = "error"
		s.logger.Error("Domain execution failed",
			zap.String("domain", domain),
			zap.Error(err))
	}

	return &domainResult{
		domain:           domain,
		status:           status,
		results:          results,
		processingTimeMS: int32(time.Since(startTime).Milliseconds()),
		error:            err,
	}
}

// domainResult represents the result from a single domain execution
type domainResult struct {
	domain           string
	status           string
	results          []types.ResultItem
	processingTimeMS int32
	error            error
}

// Mock implementations for different domains
func (s *Service) executeCommerceDomain(req *types.ExecuteTemplateRequest) ([]types.ResultItem, error) {
	return []types.ResultItem{
		{
			ID:          uuid.New().String(),
			Type:        "product",
			Title:       "Рекомендуемый товар",
			Description: "На основе вашего запроса мы подобрали оптимальный вариант",
			Data: map[string]string{
				"price":        "1500 руб",
				"availability": "в наличии",
				"rating":       "4.5",
			},
			Relevance:  0.95,
			Confidence: 0.88,
			Actions: []types.Action{
				{
					Type:   "purchase",
					Label:  "Купить",
					Method: "POST",
					URL:    "/api/v1/commerce/purchase",
				},
			},
		},
	}, nil
}

func (s *Service) executeRecipesDomain(req *types.ExecuteTemplateRequest) ([]types.ResultItem, error) {
	return []types.ResultItem{
		{
			ID:          uuid.New().String(),
			Type:        "recipe",
			Title:       "Рецепт блюда",
			Description: "Подробный рецепт с ингредиентами и инструкцией",
			Data: map[string]string{
				"cooking_time": "45 мин",
				"difficulty":   "средний",
				"servings":     "4 порции",
			},
			Relevance:  0.92,
			Confidence: 0.85,
			Actions: []types.Action{
				{
					Type:   "view_recipe",
					Label:  "Посмотреть рецепт",
					Method: "GET",
					URL:    "/api/v1/recipes/details",
				},
			},
		},
	}, nil
}

func (s *Service) executeTravelDomain(req *types.ExecuteTemplateRequest) ([]types.ResultItem, error) {
	return []types.ResultItem{
		{
			ID:          uuid.New().String(),
			Type:        "hotel",
			Title:       "Рекомендация отеля",
			Description: "Идеальный вариант для вашего путешествия",
			Data: map[string]string{
				"location":        "Москва",
				"price_per_night": "5000 руб",
				"rating":          "4.7",
				"amenities":       "WiFi, бассейн, завтрак",
			},
			Relevance:  0.89,
			Confidence: 0.82,
			Actions: []types.Action{
				{
					Type:   "book_hotel",
					Label:  "Забронировать",
					Method: "POST",
					URL:    "/api/v1/travel/book",
				},
			},
		},
	}, nil
}

func (s *Service) executeKnowledgeDomain(req *types.ExecuteTemplateRequest) ([]types.ResultItem, error) {
	return []types.ResultItem{
		{
			ID:          uuid.New().String(),
			Type:        "document",
			Title:       "Справочная информация",
			Description: "Подробная информация по вашему запросу",
			Data: map[string]string{
				"source":       "официальная документация",
				"last_updated": "2024-01-15",
				"category":     "справочная информация",
			},
			Relevance:  0.87,
			Confidence: 0.80,
			Actions: []types.Action{
				{
					Type:   "view_document",
					Label:  "Посмотреть",
					Method: "GET",
					URL:    "/api/v1/knowledge/view",
				},
			},
		},
	}, nil
}

func (s *Service) executeGenericDomain(req *types.ExecuteTemplateRequest, domain string) ([]types.ResultItem, error) {
	return []types.ResultItem{
		{
			ID:          uuid.New().String(),
			Type:        "generic",
			Title:       "Результат поиска",
			Description: fmt.Sprintf("Результаты для домена: %s", domain),
			Data: map[string]string{
				"domain": domain,
				"query":  req.Query,
			},
			Relevance:  0.75,
			Confidence: 0.70,
		},
	}, nil
}

// Helper functions
func (s *Service) containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (s *Service) determineOverallStatus(results map[string]*domainResult) string {
	hasSuccess := false
	hasError := false

	for _, result := range results {
		if result.status == "success" {
			hasSuccess = true
		} else if result.status == "error" {
			hasError = true
		}
	}

	if hasError && hasSuccess {
		return "partial"
	} else if hasError {
		return "failed"
	}
	return "completed"
}

func (s *Service) determineQueryType(query string) string {
	query = strings.ToLower(query)

	if s.containsAny(query, []string{"купить", "заказать", "приобрести"}) {
		return "with_purchases_services"
	} else if s.containsAny(query, []string{"рецепт", "как", "инструкция"}) {
		return "information_only"
	}

	return "mixed"
}

func (s *Service) buildSections(results map[string]*domainResult) []types.DomainSection {
	sections := []types.DomainSection{}

	for domain, result := range results {
		section := types.DomainSection{
			DomainID:       domain,
			Title:          s.getDomainTitle(domain),
			Status:         result.status,
			ResponseTimeMS: result.processingTimeMS,
			Results:        result.results,
		}

		if result.error != nil {
			section.Error = result.error.Error()
		}

		sections = append(sections, section)
	}

	return sections
}

func (s *Service) getDomainTitle(domain string) string {
	titles := map[string]string{
		"commerce":  "Коммерческие предложения",
		"recipes":   "Рецепты и кулинария",
		"travel":    "Путешествия и туризм",
		"knowledge": "Справочная информация",
	}

	if title, exists := titles[domain]; exists {
		return title
	}
	return fmt.Sprintf("Результаты по домену: %s", domain)
}

func (s *Service) rankResults(results map[string]*domainResult) []*rankedItem {
	var allItems []*rankedItem

	for domain, result := range results {
		for _, item := range result.results {
			allItems = append(allItems, &rankedItem{
				domain: domain,
				item:   item,
				score:  item.Relevance*0.6 + item.Confidence*0.4, // Weighted score
			})
		}
	}

	// Sort by score (simple bubble sort for demo)
	for i := 0; i < len(allItems); i++ {
		for j := i + 1; j < len(allItems); j++ {
			if allItems[i].score < allItems[j].score {
				allItems[i], allItems[j] = allItems[j], allItems[i]
			}
		}
	}

	return allItems
}

type rankedItem struct {
	domain string
	item   types.ResultItem
	score  float32
}

func (s *Service) buildRanking(rankedItems []*rankedItem) *types.RankingResult {
	if len(rankedItems) == 0 {
		return nil
	}

	items := make([]types.RankedItem, len(rankedItems))
	for i, item := range rankedItems {
		items[i] = types.RankedItem{
			ID:    item.item.ID,
			Score: item.score,
			Rank:  int32(i + 1),
		}
	}

	return &types.RankingResult{
		Items:     items,
		Algorithm: "weighted_relevance_confidence",
	}
}

func (s *Service) buildResponseMetadata(startTime time.Time) *types.ExecutionMetadata {
	return &types.ExecutionMetadata{
		StartedAt:   startTime.Unix(),
		CompletedAt: time.Now().Unix(),
		TotalTimeMS: time.Now().Sub(startTime).Milliseconds(),
	}
}

// OpenAI integration (for future use)
func (s *Service) callOpenAI(ctx context.Context, prompt string) (string, error) {
	if s.config.Provider != "openai" {
		return "", fmt.Errorf("OpenAI provider not configured")
	}

	requestBody := map[string]interface{}{
		"model":       s.config.Model,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens":  s.config.MaxTokens,
		"temperature": s.config.Temperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/chat/completions", s.config.BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return response.Choices[0].Message.Content, nil
}
