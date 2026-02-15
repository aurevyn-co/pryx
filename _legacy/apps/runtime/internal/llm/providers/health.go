package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ConnectionStatus represents the health status of a provider connection
type ConnectionStatus string

const (
	StatusHealthy     ConnectionStatus = "healthy"
	StatusDegraded    ConnectionStatus = "degraded"
	StatusError       ConnectionStatus = "error"
	StatusUnreachable ConnectionStatus = "unreachable"
)

// ProviderHealth represents the health check result for a provider
type ProviderHealth struct {
	ProviderID   string           `json:"provider_id"`
	Status       ConnectionStatus `json:"status"`
	LastChecked  time.Time        `json:"last_checked"`
	LastError    string           `json:"last_error,omitempty"`
	ModelsCount  int              `json:"models_count"`
	ResponseTime time.Duration    `json:"response_time_ms"`
	APIKeyValid  bool             `json:"api_key_valid"`
}

// HealthChecker performs health checks on LLM providers
type HealthChecker struct {
	results map[string]*ProviderHealth
	mu      sync.RWMutex
	client  *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		results: make(map[string]*ProviderHealth),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckProvider performs a health check on a specific provider
func (h *HealthChecker) CheckProvider(ctx context.Context, providerID, apiKey, baseURL string) (*ProviderHealth, error) {
	start := time.Now()

	health := &ProviderHealth{
		ProviderID:  providerID,
		LastChecked: time.Now(),
	}

	// Check based on provider type
	switch providerID {
	case "openai":
		h.checkOpenAI(ctx, health, apiKey, baseURL)
	case "anthropic":
		h.checkAnthropic(ctx, health, apiKey, baseURL)
	case "google":
		h.checkGoogle(ctx, health, apiKey, baseURL)
	case "ollama":
		h.checkOllama(ctx, health, apiKey, baseURL)
	case "openrouter":
		h.checkOpenRouter(ctx, health, apiKey, baseURL)
	default:
		health.Status = StatusError
		health.LastError = fmt.Sprintf("unsupported provider: %s", providerID)
	}

	health.ResponseTime = time.Since(start)

	// Store result
	h.mu.Lock()
	h.results[providerID] = health
	h.mu.Unlock()

	return health, nil
}

// checkOpenAI checks OpenAI API health
func (h *HealthChecker) checkOpenAI(ctx context.Context, health *ProviderHealth, apiKey, baseURL string) {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/models", nil)
	if err != nil {
		health.Status = StatusError
		health.LastError = err.Error()
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		health.Status = StatusUnreachable
		health.LastError = err.Error()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		health.Status = StatusError
		health.LastError = "invalid API key"
		health.APIKeyValid = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		health.Status = StatusDegraded
		health.LastError = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		health.Status = StatusDegraded
		health.LastError = err.Error()
		return
	}

	health.Status = StatusHealthy
	health.APIKeyValid = true
	health.ModelsCount = len(result.Data)
}

// checkAnthropic checks Anthropic API health
func (h *HealthChecker) checkAnthropic(ctx context.Context, health *ProviderHealth, apiKey, baseURL string) {
	// Anthropic doesn't have a simple models endpoint, check with minimal request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.anthropic.com/v1/models", nil)
	if err != nil {
		health.Status = StatusError
		health.LastError = err.Error()
		return
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := h.client.Do(req)
	if err != nil {
		health.Status = StatusUnreachable
		health.LastError = err.Error()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		health.Status = StatusError
		health.LastError = "invalid API key"
		health.APIKeyValid = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		health.Status = StatusDegraded
		health.LastError = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return
	}

	health.Status = StatusHealthy
	health.APIKeyValid = true
	health.ModelsCount = 3 // Claude 3 models
}

// checkGoogle checks Google AI API health
func (h *HealthChecker) checkGoogle(ctx context.Context, health *ProviderHealth, apiKey, baseURL string) {
	// Google uses API key in query param for some endpoints
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		health.Status = StatusError
		health.LastError = err.Error()
		return
	}

	resp, err := h.client.Do(req)
	if err != nil {
		health.Status = StatusUnreachable
		health.LastError = err.Error()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		health.Status = StatusError
		health.LastError = "invalid API key"
		health.APIKeyValid = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		health.Status = StatusDegraded
		health.LastError = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		health.Status = StatusDegraded
		health.LastError = err.Error()
		return
	}

	health.Status = StatusHealthy
	health.APIKeyValid = true
	health.ModelsCount = len(result.Models)
}

// checkOllama checks Ollama local server health
func (h *HealthChecker) checkOllama(ctx context.Context, health *ProviderHealth, apiKey, baseURL string) {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/tags", nil)
	if err != nil {
		health.Status = StatusError
		health.LastError = err.Error()
		return
	}

	resp, err := h.client.Do(req)
	if err != nil {
		health.Status = StatusUnreachable
		health.LastError = err.Error()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		health.Status = StatusDegraded
		health.LastError = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		health.Status = StatusDegraded
		health.LastError = err.Error()
		return
	}

	health.Status = StatusHealthy
	health.APIKeyValid = true // Ollama doesn't require API key
	health.ModelsCount = len(result.Models)
}

// checkOpenRouter checks OpenRouter API health
func (h *HealthChecker) checkOpenRouter(ctx context.Context, health *ProviderHealth, apiKey, baseURL string) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		health.Status = StatusError
		health.LastError = err.Error()
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		health.Status = StatusUnreachable
		health.LastError = err.Error()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		health.Status = StatusError
		health.LastError = "invalid API key"
		health.APIKeyValid = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		health.Status = StatusDegraded
		health.LastError = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		health.Status = StatusDegraded
		health.LastError = err.Error()
		return
	}

	health.Status = StatusHealthy
	health.APIKeyValid = true
	health.ModelsCount = len(result.Data)
}

// GetHealth returns the cached health status for a provider
func (h *HealthChecker) GetHealth(providerID string) (*ProviderHealth, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	health, ok := h.results[providerID]
	return health, ok
}

// GetAllHealth returns health status for all checked providers
func (h *HealthChecker) GetAllHealth() map[string]*ProviderHealth {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]*ProviderHealth)
	for k, v := range h.results {
		result[k] = v
	}
	return result
}

// IsStale checks if the health check is older than the given duration
func (h *ProviderHealth) IsStale(maxAge time.Duration) bool {
	return time.Since(h.LastChecked) > maxAge
}
