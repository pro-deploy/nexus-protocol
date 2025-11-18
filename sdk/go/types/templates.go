package types

// ExecuteTemplateRequest представляет запрос на выполнение шаблона
type ExecuteTemplateRequest struct {
	Query    string                 `json:"query"`
	Language string                 `json:"language,omitempty"`
	Context  *UserContext            `json:"context,omitempty"`
	Options  *ExecuteOptions         `json:"options,omitempty"`
	Metadata *RequestMetadata        `json:"metadata,omitempty"`
}

// ExecuteTemplateResponse представляет ответ на выполнение шаблона
type ExecuteTemplateResponse struct {
	ExecutionID      string            `json:"execution_id"`
	IntentID         string            `json:"intent_id,omitempty"`
	Status           string            `json:"status"`
	QueryType        string            `json:"query_type,omitempty"`
	Sections         []DomainSection   `json:"sections,omitempty"`
	WebSearch        *WebSearchResult  `json:"web_search,omitempty"`
	Ranking          *RankingResult    `json:"ranking,omitempty"`
	Metadata         *ExecutionMetadata `json:"metadata,omitempty"`
	ProcessingTimeMS int32             `json:"processing_time_ms"`
	ResponseMetadata *ResponseMetadata  `json:"response_metadata,omitempty"`
}

// UserContext содержит контекст пользователя
type UserContext struct {
	UserID    string       `json:"user_id,omitempty"`
	SessionID string       `json:"session_id,omitempty"`
	TenantID  string       `json:"tenant_id,omitempty"`
	Location  *UserLocation `json:"location,omitempty"`
}

// UserLocation содержит информацию о местоположении
type UserLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

// ExecuteOptions содержит опции выполнения
type ExecuteOptions struct {
	TimeoutMS           int32 `json:"timeout_ms,omitempty"`
	MaxResultsPerDomain int32 `json:"max_results_per_domain,omitempty"`
	ParallelExecution   bool  `json:"parallel_execution,omitempty"`
	IncludeWebSearch    bool  `json:"include_web_search,omitempty"`
}

// DomainSection представляет секцию результатов домена
type DomainSection struct {
	DomainID      string       `json:"domain_id"`
	Title         string       `json:"title,omitempty"`
	Status        string       `json:"status"`
	Error         string       `json:"error,omitempty"`
	ResponseTimeMS int32        `json:"response_time_ms,omitempty"`
	Results       []ResultItem  `json:"results,omitempty"`
}

// ResultItem представляет элемент результата
type ResultItem struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	Relevance   float32           `json:"relevance"`
	Confidence  float32           `json:"confidence"`
	Actions     []Action          `json:"actions,omitempty"`
}

// Action представляет действие, доступное для результата
type Action struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	URL         string `json:"url,omitempty"`
	Method      string `json:"method,omitempty"`
	ConfirmText string `json:"confirm_text,omitempty"`
}

// WebSearchResult представляет результаты веб-поиска
type WebSearchResult struct {
	Results      []SearchResult `json:"results,omitempty"`
	SearchEngine string         `json:"search_engine,omitempty"`
	TotalResults int32          `json:"total_results,omitempty"`
}

// SearchResult представляет результат поиска
type SearchResult struct {
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	Snippet   string  `json:"snippet,omitempty"`
	Relevance float32 `json:"relevance,omitempty"`
}

// RankingResult представляет результаты ранжирования
type RankingResult struct {
	Items     []RankedItem `json:"items,omitempty"`
	Algorithm string       `json:"algorithm,omitempty"`
}

// RankedItem представляет ранжированный элемент
type RankedItem struct {
	ID    string  `json:"id"`
	Score float32 `json:"score"`
	Rank  int32   `json:"rank"`
}

// ExecutionMetadata содержит метаданные выполнения
type ExecutionMetadata struct {
	StartedAt       int64            `json:"started_at,omitempty"`
	CompletedAt     int64            `json:"completed_at,omitempty"`
	TotalTimeMS     int64            `json:"total_time_ms,omitempty"`
	DomainStats     map[string]int32 `json:"domain_stats,omitempty"`
	DomainsExecuted int32            `json:"domains_executed,omitempty"`
	ResultsCount    int32            `json:"results_count,omitempty"`
}

