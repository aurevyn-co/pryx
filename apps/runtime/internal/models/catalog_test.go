package models

import (
	"testing"
	"time"
)

func TestCatalog_IsStale(t *testing.T) {
	fresh := &Catalog{
		CachedAt:  time.Now(),
		Models:    make(map[string]ModelInfo),
		Providers: make(map[string]ProviderInfo),
	}
	if fresh.IsStale() {
		t.Error("Expected fresh catalog to not be stale")
	}

	stale := &Catalog{
		CachedAt:  time.Now().Add(-25 * time.Hour),
		Models:    make(map[string]ModelInfo),
		Providers: make(map[string]ProviderInfo),
	}
	if !stale.IsStale() {
		t.Error("Expected old catalog to be stale")
	}
}

func TestCatalog_GetProviderModels(t *testing.T) {
	catalog := &Catalog{
		Models: map[string]ModelInfo{
			"gpt-4": {
				ID:       "gpt-4",
				Name:     "GPT-4",
				Provider: "openai",
			},
			"claude-3": {
				ID:       "claude-3",
				Name:     "Claude 3",
				Provider: "anthropic",
			},
			"gpt-3.5": {
				ID:       "gpt-3.5",
				Name:     "GPT-3.5",
				Provider: "openai",
			},
		},
		Providers: make(map[string]ProviderInfo),
		CachedAt:  time.Now(),
	}

	openaiModels := catalog.GetProviderModels("openai")
	if len(openaiModels) != 2 {
		t.Errorf("Expected 2 OpenAI models, got %d", len(openaiModels))
	}

	anthropicModels := catalog.GetProviderModels("anthropic")
	if len(anthropicModels) != 1 {
		t.Errorf("Expected 1 Anthropic model, got %d", len(anthropicModels))
	}

	nonexistent := catalog.GetProviderModels("nonexistent")
	if len(nonexistent) != 0 {
		t.Errorf("Expected 0 models for nonexistent provider, got %d", len(nonexistent))
	}
}

func TestCatalog_GetModel(t *testing.T) {
	catalog := &Catalog{
		Models: map[string]ModelInfo{
			"gpt-4": {
				ID:   "gpt-4",
				Name: "GPT-4",
			},
		},
		Providers: make(map[string]ProviderInfo),
		CachedAt:  time.Now(),
	}

	model, ok := catalog.GetModel("gpt-4")
	if !ok {
		t.Error("Expected to find GPT-4 model")
	}
	if model.Name != "GPT-4" {
		t.Errorf("Expected model name 'GPT-4', got '%s'", model.Name)
	}

	_, ok = catalog.GetModel("nonexistent")
	if ok {
		t.Error("Expected to not find nonexistent model")
	}
}

func TestCatalog_GetProvider(t *testing.T) {
	catalog := &Catalog{
		Models: make(map[string]ModelInfo),
		Providers: map[string]ProviderInfo{
			"openai": {Name: "OpenAI"},
		},
		CachedAt: time.Now(),
	}

	provider, ok := catalog.GetProvider("openai")
	if !ok {
		t.Error("Expected to find OpenAI provider")
	}
	if provider.Name != "OpenAI" {
		t.Errorf("Expected provider name 'OpenAI', got '%s'", provider.Name)
	}

	_, ok = catalog.GetProvider("nonexistent")
	if ok {
		t.Error("Expected to not find nonexistent provider")
	}
}

func TestModelInfo_SupportsTools(t *testing.T) {
	withTools := ModelInfo{ToolCall: true}
	if !withTools.SupportsTools() {
		t.Error("Expected model with ToolCall=true to support tools")
	}

	withoutTools := ModelInfo{ToolCall: false}
	if withoutTools.SupportsTools() {
		t.Error("Expected model with ToolCall=false to not support tools")
	}
}

func TestModelInfo_SupportsVision(t *testing.T) {
	withVision := ModelInfo{
		Modalities: struct {
			Input  []string `json:"input"`
			Output []string `json:"output"`
		}{
			Input: []string{"text", "image"},
		},
	}
	if !withVision.SupportsVision() {
		t.Error("Expected model with image input to support vision")
	}

	withoutVision := ModelInfo{
		Modalities: struct {
			Input  []string `json:"input"`
			Output []string `json:"output"`
		}{
			Input: []string{"text"},
		},
	}
	if withoutVision.SupportsVision() {
		t.Error("Expected model without image input to not support vision")
	}
}

func TestModelInfo_CalculateCost(t *testing.T) {
	model := ModelInfo{}
	model.Cost.Input = 2.50
	model.Cost.Output = 10.00

	cost := model.CalculateCost(1_000_000, 500_000)
	expected := 2.50 + 5.00
	if cost != expected {
		t.Errorf("Expected cost %.2f, got %.2f", expected, cost)
	}
}

func TestService_NewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService should return a non-nil service")
	}
	if service.cachePath == "" {
		t.Error("Service should have a cache path set")
	}
}

func TestGetSupportedProviders(t *testing.T) {
	providers := GetSupportedProviders()
	if len(providers) == 0 {
		t.Error("Expected non-empty list of supported providers")
	}

	hasOpenAI := false
	for _, p := range providers {
		if p == "openai" {
			hasOpenAI = true
			break
		}
	}
	if !hasOpenAI {
		t.Error("Expected 'openai' in supported providers")
	}
}

func TestRawProviderDataStructure(t *testing.T) {
	// Test that RawProviderData can hold the expected structure
	provider := RawProviderData{
		ID:   "test-provider",
		Name: "Test Provider",
		NPM:  "@test/sdk",
		Env:  []string{"TEST_API_KEY"},
		Doc:  "https://test.com/docs",
		API:  "https://api.test.com",
		Models: map[string]ModelInfo{
			"model-1": {
				ID:   "model-1",
				Name: "Model 1",
			},
		},
	}

	if provider.ID != "test-provider" {
		t.Error("Provider ID mismatch")
	}
	if len(provider.Models) != 1 {
		t.Error("Expected 1 model")
	}
}
