package types

// AIConfig представляет конфигурацию AI провайдера
type AIConfig struct {
	Provider     string            `json:"provider"` // openai, anthropic, local, custom
	Model        string            `json:"model"`    // gpt-4, claude-3, llama-2 и т.д.
	APIKey       string            `json:"api_key"`  // зашифрованный ключ
	BaseURL      string            `json:"base_url"` // для custom провайдеров
	MaxTokens    int               `json:"max_tokens"`
	Temperature  float64           `json:"temperature"`
	TopP         float64           `json:"top_p"`
	Timeout      int               `json:"timeout"`       // в секундах
	CustomParams map[string]string `json:"custom_params"` // дополнительные параметры
	Enabled      bool              `json:"enabled"`
}

// DomainCapability описывает возможности домена
type DomainCapability struct {
	Type        string            `json:"type"` // search, execute, analyze, etc.
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// DomainMLModel содержит конфигурацию ML модели домена
type DomainMLModel struct {
	Type       string            `json:"type"` // classification, regression, nlp
	Version    string            `json:"version,omitempty"`
	Accuracy   float32           `json:"accuracy,omitempty"`  // 0.0-1.0
	Threshold  float32           `json:"threshold,omitempty"` // порог уверенности 0.0-1.0
	Parameters map[string]string `json:"parameters,omitempty"`
}

// DomainConfig представляет расширенную конфигурацию домена/микросервиса
type DomainConfig struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Type         string             `json:"type"` // commerce, recipes, travel, knowledge, health, finance, custom
	Enabled      bool               `json:"enabled"`
	Endpoint     string             `json:"endpoint"`    // URL микросервиса
	AuthType     string             `json:"auth_type"`   // none, api_key, oauth2, jwt
	AuthConfig   map[string]string  `json:"auth_config"` // конфигурация авторизации
	Timeout      int                `json:"timeout"`     // в секундах
	RetryCount   int                `json:"retry_count"`
	Priority     int                `json:"priority"`      // приоритет выполнения (0-100)
	Keywords     []string           `json:"keywords"`      // ключевые слова для распознавания
	Capabilities []DomainCapability `json:"capabilities"`  // возможности домена
	MLModel      *DomainMLModel     `json:"ml_model"`      // ML модель домена
	QualityRules []QualityRule      `json:"quality_rules"` // правила оценки качества
	Metadata     map[string]string  `json:"metadata"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}

// IntegrationConfig представляет конфигурацию интеграции
type IntegrationConfig struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // payment, delivery, notifications, analytics, custom
	Provider    string                 `json:"provider"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
	Credentials map[string]string      `json:"credentials"`
	WebhookURL  string                 `json:"webhook_url"`
	Metadata    map[string]string      `json:"metadata"`
}

// PromptConfig представляет конфигурацию промпта
type PromptConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Domain      string            `json:"domain"`    // commerce, recipes, travel, knowledge
	Type        string            `json:"type"`      // system, user, assistant
	Template    string            `json:"template"`  // шаблон промпта
	Variables   []string          `json:"variables"` // переменные в шаблоне
	Version     int               `json:"version"`
	Active      bool              `json:"active"`
	Metadata    map[string]string `json:"metadata"`
}

// QualityRule определяет правило оценки качества для домена
type QualityRule struct {
	Metric      string  `json:"metric"`    // relevance, completeness, accuracy, helpfulness
	Condition   string  `json:"condition"` // min_relevance, has_price, etc.
	Threshold   float32 `json:"threshold"` // пороговое значение
	Weight      float32 `json:"weight"`    // вес правила 0.0-1.0
	Description string  `json:"description"`
}

// ExecuteTemplateRequest представляет запрос на выполнение шаблона
type ExecuteTemplateRequest struct {
	Query    string           `json:"query"`
	Language string           `json:"language,omitempty"`
	Context  *UserContext     `json:"context,omitempty"`
	Options  *ExecuteOptions  `json:"options,omitempty"`
	Filters  *AdvancedFilters `json:"filters,omitempty"` // расширенные фильтры
	Metadata *RequestMetadata `json:"metadata,omitempty"`
}

// ExecuteTemplateResponse представляет ответ на выполнение шаблона
type ExecuteTemplateResponse struct {
	ExecutionID      string                `json:"execution_id"`
	IntentID         string                `json:"intent_id,omitempty"`
	Status           string                `json:"status"`
	QueryType        string                `json:"query_type,omitempty"`
	Sections         []DomainSection       `json:"sections,omitempty"`
	WebSearch        *WebSearchResult      `json:"web_search,omitempty"`
	Ranking          *RankingResult        `json:"ranking,omitempty"`
	Metadata         *ExecutionMetadata    `json:"metadata,omitempty"`
	ProcessingTimeMS int32                 `json:"processing_time_ms"`
	ResponseMetadata *ResponseMetadata     `json:"response_metadata,omitempty"`
	Pagination       *PaginationInfo       `json:"pagination,omitempty"`      // информация о пагинации
	Workflow         *Workflow             `json:"workflow,omitempty"`        // workflow для многошаговых сценариев
	DomainAnalysis   *DomainAnalysisResult `json:"domain_analysis,omitempty"` // анализ выбора доменов
}

// UserContext содержит контекст пользователя
type UserContext struct {
	UserID    string        `json:"user_id,omitempty"`
	SessionID string        `json:"session_id,omitempty"`
	TenantID  string        `json:"tenant_id,omitempty"`
	Location  *UserLocation `json:"location,omitempty"`
	Locale    string        `json:"locale,omitempty"`   // локаль пользователя (ru-RU, en-US)
	Timezone  string        `json:"timezone,omitempty"` // часовой пояс (Europe/Moscow)
	Currency  string        `json:"currency,omitempty"` // валюта (RUB, USD, EUR)
	Region    string        `json:"region,omitempty"`   // регион (RU, US, EU)
}

// UserLocation содержит информацию о местоположении
type UserLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

// PaginationInfo содержит информацию о пагинации
type PaginationInfo struct {
	Page        int32  `json:"page,omitempty"`         // текущая страница (1-based)
	PageSize    int32  `json:"page_size,omitempty"`    // размер страницы
	TotalPages  int32  `json:"total_pages,omitempty"`  // всего страниц
	TotalItems  int64  `json:"total_items,omitempty"`  // всего элементов
	HasNext     bool   `json:"has_next,omitempty"`     // есть ли следующая страница
	HasPrevious bool   `json:"has_previous,omitempty"` // есть ли предыдущая страница
	NextCursor  string `json:"next_cursor,omitempty"`  // курсор для следующей страницы
	PrevCursor  string `json:"prev_cursor,omitempty"`  // курсор для предыдущей страницы
}

// ExecuteOptions содержит опции выполнения
type ExecuteOptions struct {
	TimeoutMS           int32 `json:"timeout_ms,omitempty"`
	MaxResultsPerDomain int32 `json:"max_results_per_domain,omitempty"`
	ParallelExecution   bool  `json:"parallel_execution,omitempty"`
	IncludeWebSearch    bool  `json:"include_web_search,omitempty"`
}

// AdvancedFilters содержит расширенные фильтры результатов
type AdvancedFilters struct {
	Domains        []string   `json:"domains,omitempty"`         // список доменов для включения
	ExcludeDomains []string   `json:"exclude_domains,omitempty"` // список доменов для исключения
	MinRelevance   float32    `json:"min_relevance,omitempty"`   // минимальная релевантность (0-1)
	MaxResults     int32      `json:"max_results,omitempty"`     // максимальное количество результатов
	SortBy         string     `json:"sort_by,omitempty"`         // сортировка: relevance, date, price
	DateRange      *DateRange `json:"date_range,omitempty"`      // диапазон дат
}

// DateRange содержит диапазон дат
type DateRange struct {
	From int64 `json:"from,omitempty"` // timestamp начала
	To   int64 `json:"to,omitempty"`   // timestamp конца
}

// DomainSection представляет секцию результатов домена
type DomainSection struct {
	DomainID       string       `json:"domain_id"`
	Title          string       `json:"title,omitempty"`
	Status         string       `json:"status"`
	Error          string       `json:"error,omitempty"`
	ResponseTimeMS int32        `json:"response_time_ms,omitempty"`
	Results        []ResultItem `json:"results,omitempty"`
}

// ResultItem представляет элемент результата
type ResultItem struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"` // изменено на interface{} для поддержки сложных структур (stores, addresses и т.д.)
	Relevance   float32                `json:"relevance"`
	Confidence  float32                `json:"confidence"`
	Actions     []Action               `json:"actions,omitempty"`
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

// Workflow представляет многошаговый сценарий с зависимостями между шагами
type Workflow struct {
	Steps []WorkflowStep `json:"steps,omitempty"`
}

// WorkflowStep представляет один шаг в workflow
type WorkflowStep struct {
	Step      int32    `json:"step"`                 // номер шага (1-based)
	Action    string   `json:"action"`               // тип действия (order_food, process_payment и т.д.)
	Domain    string   `json:"domain"`               // домен (commerce, payment, delivery, notifications)
	Status    string   `json:"status"`               // статус: pending, in_progress, completed, failed
	ResultID  string   `json:"result_id,omitempty"`  // ID результата этого шага
	DependsOn []string `json:"depends_on,omitempty"` // список ID результатов, от которых зависит этот шаг
}

// DomainAnalysisResult содержит анализ выбора доменов
type DomainAnalysisResult struct {
	SelectedDomains   []DomainSelection `json:"selected_domains,omitempty"`
	RejectedDomains   []DomainSelection `json:"rejected_domains,omitempty"`
	Confidence        float32           `json:"confidence"`
	AnalysisAlgorithm string            `json:"analysis_algorithm,omitempty"`
}

// DomainSelection представляет выбор домена с анализом
type DomainSelection struct {
	DomainID     string             `json:"domain_id"`
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	Confidence   float32            `json:"confidence"`
	Relevance    float32            `json:"relevance"`
	Reason       string             `json:"reason,omitempty"`
	Priority     int32              `json:"priority"`
	Metadata     map[string]string  `json:"metadata,omitempty"`
	Capabilities []DomainCapability `json:"capabilities,omitempty"`
}

// ResponseQualityAnalysis содержит анализ качества ответа домена
type ResponseQualityAnalysis struct {
	DomainID          string            `json:"domain_id"`
	OverallQuality    float32           `json:"overall_quality"`
	RelevanceScore    float32           `json:"relevance_score"`
	CompletenessScore float32           `json:"completeness_score"`
	AccuracyScore     float32           `json:"accuracy_score"`
	HelpfulnessScore  float32           `json:"helpfulness_score"`
	Issues            []QualityIssue    `json:"issues,omitempty"`
	Suggestions       []string          `json:"suggestions,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// QualityIssue описывает проблему с качеством ответа
type QualityIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// DomainRoutingDecision содержит решение о маршрутизации запроса
type DomainRoutingDecision struct {
	DomainID     string            `json:"domain_id"`
	Action       string            `json:"action"`
	Priority     int32             `json:"priority"`
	Reason       string            `json:"reason"`
	Alternatives []string          `json:"alternatives,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}
