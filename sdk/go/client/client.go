package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

const (
	// DefaultProtocolVersion версия протокола по умолчанию
	// Соответствует версии Application Protocol v2.0.0
	DefaultProtocolVersion = "2.0.0"
	// DefaultClientVersion версия клиента по умолчанию
	DefaultClientVersion = "2.0.0"
	// DefaultTimeout таймаут по умолчанию
	DefaultTimeout = 30 * time.Second
)

// Client представляет клиент Nexus Protocol для взаимодействия с API.
// Все методы клиента поддерживают context.Context для отмены запросов и таймаутов.
type Client struct {
	baseURL         string
	token           string
	httpClient      *http.Client
	protocolVersion string
	clientVersion   string
	clientID        string
	clientType      string
	customHeaders   map[string]string
	retryConfig     RetryConfig
	logger          Logger
	interceptors    []Interceptor
	validator       *Validator
}

// Config содержит конфигурацию клиента.
// BaseURL - базовый URL API сервера (например, "https://api.nexus.dev").
// Token - JWT токен для аутентификации (опционально, можно установить позже через SetToken).
type Config struct {
	BaseURL         string
	Token           string
	Timeout         time.Duration
	ProtocolVersion string
	ClientVersion   string
	ClientID        string
	ClientType      string
	RetryConfig     *RetryConfig // Конфигурация retry (nil = использовать по умолчанию)
	Logger          Logger      // Логгер (nil = логирование отключено)
	Validator       *Validator  // Валидатор для JSON Schema (nil = валидация отключена)
}

// NewClient создает новый клиент Nexus Protocol с указанной конфигурацией.
// Если параметры не указаны, используются значения по умолчанию:
// - ProtocolVersion: "2.0.0" (Nexus Protocol v2.0.0)
// - ClientVersion: "2.0.0"
// - Timeout: 30 секунд
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}
	if config.ProtocolVersion == "" {
		config.ProtocolVersion = DefaultProtocolVersion
	}
	if config.ClientVersion == "" {
		config.ClientVersion = DefaultClientVersion
	}

	retryCfg := DefaultRetryConfig()
	if config.RetryConfig != nil {
		retryCfg = *config.RetryConfig
	}

	logger := config.Logger
	if logger == nil {
		logger = &NoOpLogger{}
	}

	return &Client{
		baseURL:         config.BaseURL,
		token:           config.Token,
		protocolVersion: config.ProtocolVersion,
		clientVersion:   config.ClientVersion,
		clientID:        config.ClientID,
		clientType:      config.ClientType,
		customHeaders:   make(map[string]string),
		retryConfig:     retryCfg,
		logger:          logger,
		interceptors:    make([]Interceptor, 0),
		validator:       config.Validator,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SetToken устанавливает JWT токен для аутентификации.
// Токен будет использоваться во всех последующих запросах.
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetCustomHeader устанавливает кастомный заголовок для всех последующих запросов.
// Заголовок будет включен в RequestMetadata.custom_headers.
func (c *Client) SetCustomHeader(key, value string) {
	if c.customHeaders == nil {
		c.customHeaders = make(map[string]string)
	}
	c.customHeaders[key] = value
}

// RemoveCustomHeader удаляет кастомный заголовок.
func (c *Client) RemoveCustomHeader(key string) {
	if c.customHeaders != nil {
		delete(c.customHeaders, key)
	}
}

// ClearCustomHeaders удаляет все кастомные заголовки.
func (c *Client) ClearCustomHeaders() {
	c.customHeaders = make(map[string]string)
}

// getCustomHeaders возвращает копию кастомных заголовков для включения в метаданные
func (c *Client) getCustomHeaders() map[string]string {
	if c.customHeaders == nil {
		return make(map[string]string)
	}
	// Возвращаем копию, чтобы избежать изменений извне
	result := make(map[string]string, len(c.customHeaders))
	for k, v := range c.customHeaders {
		result[k] = v
	}
	return result
}

// SetPriority устанавливает приоритет запроса
// Значения: "low", "normal", "high", "critical"
func (c *Client) SetPriority(priority string) {
	c.SetCustomHeader("x-priority", priority)
}

// SetRequestSource устанавливает источник запроса
// Значения: "user", "system", "batch", "webhook"
func (c *Client) SetRequestSource(source string) {
	c.SetCustomHeader("x-request-source", source)
}

// SetCacheControl устанавливает контроль кэширования
// Значения: "no-cache", "cache-only", "cache-first", "network-first"
func (c *Client) SetCacheControl(cacheControl string) {
	c.SetCustomHeader("x-cache-control", cacheControl)
}

// SetCacheTTL устанавливает TTL кэша в секундах
func (c *Client) SetCacheTTL(ttl int32) {
	c.SetCustomHeader("x-cache-ttl", fmt.Sprintf("%d", ttl))
}

// SetCacheKey устанавливает кастомный ключ кэша
func (c *Client) SetCacheKey(key string) {
	c.SetCustomHeader("x-cache-key", key)
}

// SetFeatureFlag устанавливает feature flag для A/B тестирования
func (c *Client) SetFeatureFlag(flag, value string) {
	c.SetCustomHeader(fmt.Sprintf("x-feature-%s", flag), value)
}

// SetExperiment устанавливает ID эксперимента
func (c *Client) SetExperiment(experimentID string) {
	c.SetCustomHeader("x-experiment-id", experimentID)
}

// createRequestMetadata создает RequestMetadata с настройками клиента
func (c *Client) createRequestMetadata() *types.RequestMetadata {
	metadata := types.NewRequestMetadata(c.protocolVersion, c.clientVersion)
	metadata.ClientID = c.clientID
	metadata.ClientType = c.clientType
	metadata.CustomHeaders = c.getCustomHeaders()
	return metadata
}

// doRequest выполняет HTTP запрос с поддержкой context, retry и rate limiting
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			// Вычисляем задержку для retry
			backoff := c.calculateBackoff(attempt - 1)
			
			c.logger.Debug("Retrying request",
				Field{Key: "attempt", Value: attempt},
				Field{Key: "backoff_ms", Value: backoff.Milliseconds()},
				Field{Key: "path", Value: path},
			)

			// Ждем перед повтором
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Создаем запрос
		var reqBody io.Reader
		if body != nil {
			// Валидация запроса (если валидатор настроен)
			if c.validator != nil && attempt == 0 {
				// Определяем схему по пути (можно улучшить)
				schemaName := c.getSchemaNameForPath(path)
				if schemaName != "" {
					if err := c.validator.ValidateRequest(schemaName, body); err != nil {
						return nil, fmt.Errorf("request validation failed: %w", err)
					}
				}
			}

			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonData)
		}

		url := c.baseURL + path
		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		// Применяем interceptors перед запросом
		if err := c.applyInterceptorsBefore(ctx, req); err != nil {
			return nil, fmt.Errorf("interceptor error: %w", err)
		}

		// Логируем запрос
		c.logger.Debug("Sending request",
			Field{Key: "method", Value: method},
			Field{Key: "path", Value: path},
			Field{Key: "attempt", Value: attempt + 1},
		)

		startTime := time.Now()
		resp, err := c.httpClient.Do(req)
		duration := time.Since(startTime)

		// Логируем ответ
		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}
		c.logger.Debug("Received response",
			Field{Key: "method", Value: method},
			Field{Key: "path", Value: path},
			Field{Key: "status_code", Value: statusCode},
			Field{Key: "duration_ms", Value: duration.Milliseconds()},
		)

		// Применяем interceptors после ответа
		if resp != nil {
			if err := c.applyInterceptorsAfter(ctx, req, resp); err != nil {
				if resp.Body != nil {
					resp.Body.Close()
				}
				return nil, fmt.Errorf("interceptor error: %w", err)
			}
		}

		// Обрабатываем ошибки
		if err != nil {
			lastErr = err
			if !c.shouldRetry(attempt+1, err, 0) {
				return nil, fmt.Errorf("request failed: %w", err)
			}
			continue
		}

		// Обрабатываем rate limiting (HTTP 429)
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := c.handleRateLimit(resp)
			if c.shouldRetry(attempt+1, nil, resp.StatusCode) {
				c.logger.Warn("Rate limited, waiting",
					Field{Key: "retry_after_sec", Value: retryAfter.Seconds()},
					Field{Key: "path", Value: path},
				)
				resp.Body.Close()
				lastResp = resp
				lastErr = fmt.Errorf("rate limited")
				
				// Ждем указанное время или используем backoff
				waitTime := retryAfter
				if waitTime == 0 {
					waitTime = c.calculateBackoff(attempt)
				}
				
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(waitTime):
				}
				continue
			}
		}

		// Проверяем другие retryable статусы
		if resp.StatusCode >= 400 {
			if c.shouldRetry(attempt+1, nil, resp.StatusCode) {
				lastResp = resp
				lastErr = fmt.Errorf("request failed with status %d", resp.StatusCode)
				resp.Body.Close()
				continue
			}
			// Не retryable ошибка - возвращаем сразу
			return resp, nil
		}

		// Успешный ответ - валидация ответа (если валидатор настроен)
		if c.validator != nil && resp.StatusCode < 400 {
			schemaName := c.getSchemaNameForPath(path)
			if schemaName != "" {
				// Читаем body для валидации
				bodyBytes, err := io.ReadAll(resp.Body)
				if err == nil {
					// Восстанавливаем body
					resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					
					// Валидируем (можно улучшить, добавив схему для ответа)
					// Пока пропускаем валидацию ответов, так как нужны разные схемы
				}
			}
		}

		return resp, nil
	}

	// Все попытки исчерпаны
	if lastResp != nil {
		return lastResp, lastErr
	}
	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryConfig.MaxRetries+1, lastErr)
}

// getSchemaNameForPath возвращает имя схемы для пути (базовая реализация)
func (c *Client) getSchemaNameForPath(_ string) string {
	// Можно улучшить, добавив маппинг путей к схемам
	// Пока возвращаем пустую строку (валидация отключена по умолчанию)
	return ""
}

// SetValidator устанавливает валидатор для клиента
func (c *Client) SetValidator(validator *Validator) {
	c.validator = validator
}

// handleRateLimit обрабатывает rate limiting и возвращает время ожидания
func (c *Client) handleRateLimit(resp *http.Response) time.Duration {
	// Пытаемся получить Retry-After заголовок
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter != "" {
		// Может быть число секунд или HTTP date
		if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
			return seconds
		}
		if date, err := http.ParseTime(retryAfter); err == nil {
			return time.Until(date)
		}
	}

	// Пытаемся получить из метаданных ошибки
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var errResp types.ErrorResponse
		if json.Unmarshal(body, &errResp) == nil {
			if resetAt, ok := errResp.Error.Metadata["reset_at"]; ok {
				if resetTime, err := time.Parse(time.RFC3339, resetAt); err == nil {
					return time.Until(resetTime)
				}
			}
		}
		// Восстанавливаем body для дальнейшей обработки
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return 0
}


// parseResponse парсит ответ и обрабатывает ошибки
func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Проверяем на ошибку
	if resp.StatusCode >= 400 {
		var errResp types.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return &errResp.Error
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Парсим успешный ответ
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Проверяем совместимость версий протокола, если в ответе есть ResponseMetadata
		// Пытаемся извлечь ResponseMetadata из результата
		if responseWithMetadata, ok := result.(interface{ GetMetadata() *types.ResponseMetadata }); ok {
			if metadata := responseWithMetadata.GetMetadata(); metadata != nil {
				if err := types.ValidateResponseMetadata(metadata); err != nil {
					c.logger.Warn("Invalid response metadata",
						Field{Key: "error", Value: err.Error()},
					)
				} else {
					compatible, err := types.IsCompatible(c.protocolVersion, metadata.ProtocolVersion)
					if err != nil {
						c.logger.Warn("Failed to check version compatibility",
							Field{Key: "error", Value: err.Error()},
						)
					} else if !compatible {
						return &types.ErrorDetail{
							Code:    "PROTOCOL_VERSION_MISMATCH",
							Type:    "PROTOCOL_VERSION_ERROR",
							Message: fmt.Sprintf("Protocol version mismatch: client %s is not compatible with server %s", c.protocolVersion, metadata.ProtocolVersion),
							Details: fmt.Sprintf("Client version %s is not compatible with server version %s. Major versions must match, and client minor version must not exceed server minor version.", c.protocolVersion, metadata.ProtocolVersion),
						}
					}
				}
			}
		}

		// Альтернативный способ: проверяем структуры с полем Metadata типа *ResponseMetadata
		var tempStruct struct {
			Metadata *types.ResponseMetadata `json:"metadata"`
		}
		if err := json.Unmarshal(body, &tempStruct); err == nil && tempStruct.Metadata != nil {
			if err := types.ValidateResponseMetadata(tempStruct.Metadata); err != nil {
				c.logger.Warn("Invalid response metadata",
					Field{Key: "error", Value: err.Error()},
				)
			} else {
				compatible, err := types.IsCompatible(c.protocolVersion, tempStruct.Metadata.ProtocolVersion)
				if err != nil {
					c.logger.Warn("Failed to check version compatibility",
						Field{Key: "error", Value: err.Error()},
					)
				} else if !compatible {
					return &types.ErrorDetail{
						Code:    "PROTOCOL_VERSION_MISMATCH",
						Type:    "PROTOCOL_VERSION_ERROR",
						Message: fmt.Sprintf("Protocol version mismatch: client %s is not compatible with server %s", c.protocolVersion, tempStruct.Metadata.ProtocolVersion),
						Details: fmt.Sprintf("Client version %s is not compatible with server version %s. Major versions must match, and client minor version must not exceed server minor version.", c.protocolVersion, tempStruct.Metadata.ProtocolVersion),
					}
				}
			}
		}
	}

	return nil
}

// Health проверяет здоровье сервера.
// Возвращает информацию о статусе сервера и его версии.
func (c *Client) Health(ctx context.Context) (*types.HealthResponse, error) {
	resp, err := c.doRequest(ctx, "GET", PathHealth, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.HealthResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Ready проверяет готовность сервера (readiness probe для Kubernetes).
// Возвращает детальную информацию о состоянии всех компонентов сервера.
//
// Пример использования:
//
//	ctx := context.Background()
//	ready, err := client.Ready(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if ready.Ready {
//		fmt.Println("Server is ready")
//	}

// Admin returns an admin client for administrative operations
func (c *Client) Admin() *AdminClient {
	return &AdminClient{client: c}
}

func (c *Client) Ready(ctx context.Context) (*types.ReadinessResponse, error) {
	resp, err := c.doRequest(ctx, "GET", PathReady, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.ReadinessResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}


