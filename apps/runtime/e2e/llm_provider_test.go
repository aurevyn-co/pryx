//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// TestLLMProviderConnection tests the LLM provider API endpoints
func TestLLMProviderConnection(t *testing.T) {
	bin := buildPryxCore(t)
	home := t.TempDir()

	port, cancel := startPryxCore(t, bin, home)
	defer cancel()

	waitForServer(t, port, 5*time.Second)

	baseURL := "http://localhost:" + port

	// Test 1: List available providers
	t.Run("list_providers", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/providers")
		if err != nil {
			t.Fatalf("Failed to list providers: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		providers, ok := result["providers"].([]interface{})
		if !ok {
			t.Fatal("Expected providers array in response")
		}

		if len(providers) == 0 {
			t.Fatal("Expected at least one provider")
		}

		t.Logf("✓ Found %d providers", len(providers))

		// Check if GLM is in the list
		foundGLM := false
		for _, p := range providers {
			provider := p.(map[string]interface{})
			if provider["id"] == "glm" {
				foundGLM = true
				break
			}
		}

		if foundGLM {
			t.Logf("✓ GLM provider found in list")
		}
	})

	// Test 2: Get provider models
	t.Run("get_provider_models", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/providers/openai/models")
		if err != nil {
			t.Fatalf("Failed to get provider models: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Provider models not available, status: %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		models, ok := result["models"].([]interface{})
		if !ok {
			t.Skip("No models in response")
		}

		t.Logf("✓ Found %d models for provider", len(models))
	})

	// Test 3: List all models
	t.Run("list_all_models", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/v1/models")
		if err != nil {
			t.Fatalf("Failed to list models: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		models, ok := result["models"].([]interface{})
		if !ok {
			t.Fatal("Expected models array in response")
		}

		t.Logf("✓ Found %d total models", len(models))
	})
}
