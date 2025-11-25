package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

// AdminClient предоставляет методы для управления конфигурацией системы (только для администраторов)
type AdminClient struct {
	client *Client
}

// GetAIConfig получает текущую конфигурацию AI
func (ac *AdminClient) GetAIConfig(ctx context.Context) (*types.AIConfig, error) {
	resp, err := ac.client.doRequest(ctx, http.MethodGet, PathAPIV1AdminAIConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI config: %w", err)
	}
	defer resp.Body.Close()

	var config types.AIConfig
	if err := ac.client.parseResponse(resp, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateAIConfig обновляет конфигурацию AI
func (ac *AdminClient) UpdateAIConfig(ctx context.Context, config *types.AIConfig) error {
	resp, err := ac.client.doRequest(ctx, http.MethodPut, PathAPIV1AdminAIConfig, config)
	if err != nil {
		return fmt.Errorf("failed to update AI config: %w", err)
	}
	defer resp.Body.Close()

	if err := ac.client.parseResponse(resp, nil); err != nil {
		return err
	}
	return nil
}

// ListPrompts получает список всех промптов
func (ac *AdminClient) ListPrompts(ctx context.Context, domain string) ([]*types.PromptConfig, error) {
	path := PathAPIV1AdminPrompts
	if domain != "" {
		path += "?domain=" + domain
	}

	resp, err := ac.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}
	defer resp.Body.Close()

	var prompts []*types.PromptConfig
	if err := ac.client.parseResponse(resp, &prompts); err != nil {
		return nil, err
	}
	return prompts, nil
}

// GetPrompt получает промпт по ID
func (ac *AdminClient) GetPrompt(ctx context.Context, id string) (*types.PromptConfig, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminPrompts, id)
	resp, err := ac.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt %s: %w", id, err)
	}
	defer resp.Body.Close()

	var prompt types.PromptConfig
	if err := ac.client.parseResponse(resp, &prompt); err != nil {
		return nil, err
	}
	return &prompt, nil
}

// CreatePrompt создает новый промпт
func (ac *AdminClient) CreatePrompt(ctx context.Context, prompt *types.PromptConfig) (*types.PromptConfig, error) {
	resp, err := ac.client.doRequest(ctx, http.MethodPost, PathAPIV1AdminPrompts, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}
	defer resp.Body.Close()

	var created types.PromptConfig
	if err := ac.client.parseResponse(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdatePrompt обновляет существующий промпт
func (ac *AdminClient) UpdatePrompt(ctx context.Context, id string, prompt *types.PromptConfig) (*types.PromptConfig, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminPrompts, id)
	resp, err := ac.client.doRequest(ctx, http.MethodPut, path, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to update prompt %s: %w", id, err)
	}
	defer resp.Body.Close()

	var updated types.PromptConfig
	if err := ac.client.parseResponse(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeletePrompt удаляет промпт
func (ac *AdminClient) DeletePrompt(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminPrompts, id)
	resp, err := ac.client.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete prompt %s: %w", id, err)
	}
	defer resp.Body.Close()

	if err := ac.client.parseResponse(resp, nil); err != nil {
		return err
	}
	return nil
}

// ListDomains получает список всех доменов
func (ac *AdminClient) ListDomains(ctx context.Context) ([]*types.DomainConfig, error) {
	resp, err := ac.client.doRequest(ctx, http.MethodGet, PathAPIV1AdminDomains, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}
	defer resp.Body.Close()

	var domains []*types.DomainConfig
	if err := ac.client.parseResponse(resp, &domains); err != nil {
		return nil, err
	}
	return domains, nil
}

// GetDomain получает домен по ID
func (ac *AdminClient) GetDomain(ctx context.Context, id string) (*types.DomainConfig, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminDomains, id)
	resp, err := ac.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain %s: %w", id, err)
	}
	defer resp.Body.Close()

	var domain types.DomainConfig
	if err := ac.client.parseResponse(resp, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

// CreateDomain создает новый домен
func (ac *AdminClient) CreateDomain(ctx context.Context, domain *types.DomainConfig) (*types.DomainConfig, error) {
	resp, err := ac.client.doRequest(ctx, http.MethodPost, PathAPIV1AdminDomains, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}
	defer resp.Body.Close()

	var created types.DomainConfig
	if err := ac.client.parseResponse(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateDomain обновляет домен
func (ac *AdminClient) UpdateDomain(ctx context.Context, id string, domain *types.DomainConfig) (*types.DomainConfig, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminDomains, id)
	resp, err := ac.client.doRequest(ctx, http.MethodPut, path, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain %s: %w", id, err)
	}
	defer resp.Body.Close()

	var updated types.DomainConfig
	if err := ac.client.parseResponse(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteDomain удаляет домен
func (ac *AdminClient) DeleteDomain(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminDomains, id)
	resp, err := ac.client.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete domain %s: %w", id, err)
	}
	defer resp.Body.Close()

	if err := ac.client.parseResponse(resp, nil); err != nil {
		return err
	}
	return nil
}

// InitializeDefaultDomains инициализирует домены по умолчанию
func (ac *AdminClient) InitializeDefaultDomains(ctx context.Context) error {
	resp, err := ac.client.doRequest(ctx, http.MethodPost, PathAPIV1AdminDomains+"/initialize-default", nil)
	if err != nil {
		return fmt.Errorf("failed to initialize default domains: %w", err)
	}
	defer resp.Body.Close()

	if err := ac.client.parseResponse(resp, nil); err != nil {
		return err
	}
	return nil
}

// ListIntegrations получает список интеграций
func (ac *AdminClient) ListIntegrations(ctx context.Context, integrationType string) ([]*types.IntegrationConfig, error) {
	path := PathAPIV1AdminIntegrations
	if integrationType != "" {
		path += "?type=" + integrationType
	}
	resp, err := ac.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list integrations: %w", err)
	}
	defer resp.Body.Close()

	var integrations []*types.IntegrationConfig
	if err := ac.client.parseResponse(resp, &integrations); err != nil {
		return nil, err
	}
	return integrations, nil
}

// GetIntegration получает интеграцию по ID
func (ac *AdminClient) GetIntegration(ctx context.Context, id string) (*types.IntegrationConfig, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminIntegrations, id)
	resp, err := ac.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration %s: %w", id, err)
	}
	defer resp.Body.Close()

	var integration types.IntegrationConfig
	if err := ac.client.parseResponse(resp, &integration); err != nil {
		return nil, err
	}
	return &integration, nil
}

// CreateIntegration создает новую интеграцию
func (ac *AdminClient) CreateIntegration(ctx context.Context, config *types.IntegrationConfig) (*types.IntegrationConfig, error) {
	resp, err := ac.client.doRequest(ctx, http.MethodPost, PathAPIV1AdminIntegrations, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}
	defer resp.Body.Close()

	var created types.IntegrationConfig
	if err := ac.client.parseResponse(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateIntegration обновляет интеграцию
func (ac *AdminClient) UpdateIntegration(ctx context.Context, id string, config *types.IntegrationConfig) (*types.IntegrationConfig, error) {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminIntegrations, id)
	resp, err := ac.client.doRequest(ctx, http.MethodPut, path, config)
	if err != nil {
		return nil, fmt.Errorf("failed to update integration %s: %w", id, err)
	}
	defer resp.Body.Close()

	var updated types.IntegrationConfig
	if err := ac.client.parseResponse(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteIntegration удаляет интеграцию
func (ac *AdminClient) DeleteIntegration(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s", PathAPIV1AdminIntegrations, id)
	resp, err := ac.client.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete integration %s: %w", id, err)
	}
	defer resp.Body.Close()

	if err := ac.client.parseResponse(resp, nil); err != nil {
		return err
	}
	return nil
}

// GetVersion получает информацию о версии системы
func (ac *AdminClient) GetVersion(ctx context.Context) (map[string]string, error) {
	resp, err := ac.client.doRequest(ctx, http.MethodGet, PathAPIV1AdminVersion, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get version info: %w", err)
	}
	defer resp.Body.Close()

	var version map[string]string
	if err := ac.client.parseResponse(resp, &version); err != nil {
		return nil, err
	}
	return version, nil
}
