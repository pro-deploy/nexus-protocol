package types

import (
	"time"
)

// RequestMetadata contains metadata for all requests
type RequestMetadata struct {
	RequestID      string            `json:"request_id" validate:"required,uuid4"`
	ProtocolVersion string           `json:"protocol_version" validate:"required"`
	ClientVersion   string           `json:"client_version" validate:"required"`
	ClientID        string           `json:"client_id,omitempty"`
	ClientType      string           `json:"client_type,omitempty" validate:"omitempty,oneof=web mobile sdk api desktop"`
	Timestamp       int64            `json:"timestamp" validate:"required,min=0"`
	CustomHeaders   map[string]string `json:"custom_headers,omitempty"`
}

// ResponseMetadata contains metadata for all responses
type ResponseMetadata struct {
	RequestID        string         `json:"request_id" validate:"required,uuid4"`
	ProtocolVersion  string         `json:"protocol_version" validate:"required"`
	ServerVersion    string         `json:"server_version" validate:"required"`
	Timestamp        int64          `json:"timestamp" validate:"required,min=0"`
	ProcessingTimeMS int32          `json:"processing_time_ms" validate:"min=0"`

	// Enterprise features
	RateLimitInfo *RateLimitInfo `json:"rate_limit_info,omitempty"`
	CacheInfo     *CacheInfo     `json:"cache_info,omitempty"`
	QuotaInfo     *QuotaInfo     `json:"quota_info,omitempty"`
}

// RateLimitInfo contains rate limiting information
type RateLimitInfo struct {
	Limit    int32 `json:"limit" validate:"min=0"`
	Remaining int32 `json:"remaining" validate:"min=0"`
	ResetAt  int64 `json:"reset_at" validate:"min=0"`
}

// CacheInfo contains caching information
type CacheInfo struct {
	CacheHit bool   `json:"cache_hit"`
	CacheKey string `json:"cache_key,omitempty"`
	CacheTTL int32  `json:"cache_ttl,omitempty" validate:"min=0"`
}

// QuotaInfo contains quota information
type QuotaInfo struct {
	QuotaUsed  int64  `json:"quota_used" validate:"min=0"`
	QuotaLimit int64  `json:"quota_limit" validate:"min=0"`
	QuotaType  string `json:"quota_type" validate:"oneof=requests data storage bandwidth"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username,omitempty" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Status    string    `json:"status" db:"status" validate:"oneof=active inactive suspended"`
	Roles     []string  `json:"roles" db:"roles"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// UserContext contains user context information
type UserContext struct {
	UserID    string      `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	SessionID string      `json:"session_id,omitempty"`
	TenantID  string      `json:"tenant_id,omitempty" validate:"omitempty,uuid4"`
	Roles     []string    `json:"roles,omitempty"`
	Location  *Location   `json:"location,omitempty"`
	Locale    string      `json:"locale,omitempty"`
	Timezone  string      `json:"timezone,omitempty"`
	Currency  string      `json:"currency,omitempty"`
	Region    string      `json:"region,omitempty"`
}

// Location represents geographic location
type Location struct {
	Latitude  float64 `json:"latitude" validate:"min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"min=-180,max=180"`
	Accuracy  float64 `json:"accuracy,omitempty" validate:"min=0"`
}

// ExecuteTemplateRequest represents a template execution request
type ExecuteTemplateRequest struct {
	Query    string       `json:"query" validate:"required,min=1,max=1000"`
	Language string       `json:"language" validate:"required,oneof=ru en"`
	Context  *UserContext `json:"context,omitempty"`
	Options  *ExecuteOptions `json:"options,omitempty"`
	Metadata *RequestMetadata `json:"metadata,omitempty"`

	// Enterprise features
	Filters *AdvancedFilters `json:"filters,omitempty"`
}

// ExecuteOptions contains execution options
type ExecuteOptions struct {
	TimeoutMS           int32 `json:"timeout_ms,omitempty" validate:"min=0,max=120000"`
	MaxResultsPerDomain int32 `json:"max_results_per_domain,omitempty" validate:"min=0,max=50"`
	ParallelExecution   bool  `json:"parallel_execution,omitempty"`
	IncludeWebSearch    bool  `json:"include_web_search,omitempty"`
}

// AdvancedFilters contains advanced filtering options
type AdvancedFilters struct {
	Domains      []string `json:"domains,omitempty"`
	MinRelevance float32  `json:"min_relevance,omitempty" validate:"min=0,max=1"`
	MaxResults   int32    `json:"max_results,omitempty" validate:"min=0,max=100"`
	SortBy       string   `json:"sort_by,omitempty" validate:"oneof=relevance date rating"`
}

// ExecuteTemplateResponse represents a template execution response
type ExecuteTemplateResponse struct {
	ExecutionID       string                  `json:"execution_id" validate:"required,uuid4"`
	IntentID          string                  `json:"intent_id,omitempty" validate:"omitempty,uuid4"`
	Status            string                  `json:"status" validate:"required,oneof=in_progress completed partial failed timeout"`
	QueryType         string                  `json:"query_type,omitempty" validate:"omitempty,oneof=information_only with_purchases_services mixed"`
	Sections          []DomainSection         `json:"sections,omitempty"`
	WebSearch         *WebSearchResult        `json:"web_search,omitempty"`
	Ranking           *RankingResult          `json:"ranking,omitempty"`
	Metadata          *ExecutionMetadata      `json:"metadata,omitempty"`
	ProcessingTimeMS  int32                   `json:"processing_time_ms,omitempty"`
	ResponseMetadata  *ResponseMetadata       `json:"response_metadata,omitempty"`
	Pagination        *PaginationInfo         `json:"pagination,omitempty"`
}

// DomainSection represents results from a specific domain
type DomainSection struct {
	DomainID       string       `json:"domain_id" validate:"required"`
	Title          string       `json:"title" validate:"required"`
	Status         string       `json:"status" validate:"required,oneof=success error timeout partial"`
	Error          string       `json:"error,omitempty"`
	ResponseTimeMS int32        `json:"response_time_ms,omitempty" validate:"min=0"`
	Results        []ResultItem `json:"results,omitempty"`
}

// ResultItem represents a single result item
type ResultItem struct {
	ID          string            `json:"id" validate:"required"`
	Type        string            `json:"type" validate:"required"`
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	Relevance   float32           `json:"relevance,omitempty" validate:"min=0,max=1"`
	Confidence  float32           `json:"confidence,omitempty" validate:"min=0,max=1"`
	Actions     []Action          `json:"actions,omitempty"`
}

// Action represents an available action for a result
type Action struct {
	Type        string `json:"type" validate:"required"`
	Label       string `json:"label" validate:"required"`
	URL         string `json:"url,omitempty"`
	Method      string `json:"method,omitempty" validate:"omitempty,oneof=GET POST PUT DELETE"`
	ConfirmText string `json:"confirm_text,omitempty"`
}

// WebSearchResult contains web search results
type WebSearchResult struct {
	Results      []SearchResult `json:"results,omitempty"`
	SearchEngine string         `json:"search_engine,omitempty"`
	TotalResults int32          `json:"total_results,omitempty"`
}

// SearchResult represents a web search result
type SearchResult struct {
	Title    string  `json:"title" validate:"required"`
	URL      string  `json:"url" validate:"required,url"`
	Snippet  string  `json:"snippet,omitempty"`
	Relevance float32 `json:"relevance,omitempty" validate:"min=0,max=1"`
}

// RankingResult contains ranking results
type RankingResult struct {
	Items    []RankedItem `json:"items,omitempty"`
	Algorithm string      `json:"algorithm,omitempty"`
}

// RankedItem represents a ranked item
type RankedItem struct {
	ID    string  `json:"id" validate:"required"`
	Score float32 `json:"score" validate:"min=0"`
	Rank  int32   `json:"rank" validate:"min=1"`
}

// ExecutionMetadata contains execution metadata
type ExecutionMetadata struct {
	StartedAt      int64            `json:"started_at" validate:"min=0"`
	CompletedAt    int64            `json:"completed_at,omitempty" validate:"min=0"`
	TotalTimeMS    int64            `json:"total_time_ms,omitempty" validate:"min=0"`
	DomainStats    map[string]int32 `json:"domain_stats,omitempty"`
	DomainsExecuted int32           `json:"domains_executed,omitempty" validate:"min=0"`
	ResultsCount   int32            `json:"results_count,omitempty" validate:"min=0"`
}

// PaginationInfo contains pagination information
type PaginationInfo struct {
	Page      int32 `json:"page" validate:"min=1"`
	PageSize  int32 `json:"page_size" validate:"min=1,max=100"`
	TotalItems int32 `json:"total_items" validate:"min=0"`
	HasNext   bool  `json:"has_next"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code     string            `json:"error_code" validate:"required"`
	Type     string            `json:"error_type" validate:"required,oneof=VALIDATION_ERROR AUTHENTICATION_ERROR AUTHORIZATION_ERROR NOT_FOUND CONFLICT RATE_LIMIT_ERROR INTERNAL_ERROR EXTERNAL_ERROR PROTOCOL_VERSION_ERROR"`
	Message  string            `json:"message" validate:"required,max=1000"`
	Field    string            `json:"field,omitempty"`
	Details  string            `json:"details,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// ReadinessResponse represents readiness check response
type ReadinessResponse struct {
	Status    string                     `json:"status"`
	Timestamp string                     `json:"timestamp"`
	Checks    ReadinessChecks            `json:"checks"`
	Components map[string]*ComponentStatus `json:"components,omitempty"`
	Capacity  *CapacityInfo              `json:"capacity,omitempty"`
}

// ReadinessChecks contains readiness check results
type ReadinessChecks struct {
	Database   string `json:"database"`
	Redis      string `json:"redis"`
	AIServices string `json:"ai_services"`
}

// ComponentStatus represents component health status
type ComponentStatus struct {
	Status     string `json:"status"`
	LatencyMS  int64  `json:"latency_ms,omitempty"`
	Message    string `json:"message,omitempty"`
}

// CapacityInfo contains system capacity information
type CapacityInfo struct {
	CurrentLoad     float32 `json:"current_load"`
	MaxCapacity     int64   `json:"max_capacity"`
	AvailableCapacity int64 `json:"available_capacity"`
	ActiveConnections int64 `json:"active_connections"`
}
