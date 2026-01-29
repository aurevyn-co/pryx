package factory

import (
	"os"
	"testing"

	"pryx-core/internal/llm/providers"
)

func TestNewProvider_OpenAI(t *testing.T) {
	p, err := NewProvider("openai", "test-key", "")
	if err != nil {
		t.Fatalf("Failed to create OpenAI provider: %v", err)
	}
	if _, ok := p.(*providers.OpenAIProvider); !ok {
		t.Errorf("Expected providers.OpenAIProvider type")
	}
}

func TestNewProvider_Anthropic(t *testing.T) {
	p, err := NewProvider("anthropic", "test-key", "")
	if err != nil {
		t.Fatalf("Failed to create Anthropic provider: %v", err)
	}
	if _, ok := p.(*providers.AnthropicProvider); !ok {
		t.Errorf("Expected providers.AnthropicProvider type")
	}
}

func TestNewProvider_OpenRouter(t *testing.T) {
	p, err := NewProvider("openrouter", "test-key", "")
	if err != nil {
		t.Fatalf("Failed to create OpenRouter provider: %v", err)
	}
	if _, ok := p.(*providers.OpenAIProvider); !ok {
		t.Errorf("Expected providers.OpenAIProvider type (OpenRouter uses OpenAI client)")
	}
}

func TestNewProvider_Ollama(t *testing.T) {
	p, err := NewProvider("ollama", "", "http://localhost:11434")
	if err != nil {
		t.Fatalf("Failed to create Ollama provider: %v", err)
	}
	if _, ok := p.(*providers.OpenAIProvider); !ok {
		t.Errorf("Expected providers.OpenAIProvider type")
	}
}

func TestNewProvider_CustomBaseURL(t *testing.T) {
	os.Setenv("OPENAI_BASE_URL", "https://custom.api/v1")
	defer os.Unsetenv("OPENAI_BASE_URL")

	p, err := NewProvider("openai", "test-key", "")
	if err != nil {
		t.Fatalf("Failed to create custom provider: %v", err)
	}

	if _, ok := p.(*providers.OpenAIProvider); !ok {
		t.Errorf("Expected providers.OpenAIProvider type")
	}
}
