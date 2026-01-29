package models_test

import (
	"os"
	"testing"
	"time"

	"pryx-core/internal/store"
	"pryx-core/internal/models"
)

// TestCatalog_Load tests catalog loading
func TestCatalog_Load(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	defer os.Remove(tmpDB)

	store, err := store.New(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	catalog, err := models.LoadCatalog(store)
	if err != nil {
		t.Logf("Load catalog error (expected if no models): %v", err)
	}

	if catalog == nil {
		t.Fatal("Catalog should not be nil")
	}

	t.Logf("Catalog loaded with %d models", len(catalog.Models))
}

// TestCatalog_GetProviderModels tests filtering by provider
func TestCatalog_GetProviderModels(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	store, _ := store.New(tmpDB)
	defer store.Close()

	catalog, _ := models.LoadCatalog(store)

	// Create test models
	testModel1 := models.ModelInfo{
		ID:       "model-1",
		Name:     "Test Model 1",
		Provider: "openai",
	}
	testModel2 := models.ModelInfo{
		ID:       "model-2",
		Name:     "Test Model 2",
		Provider: "anthropic",
	}

	// Use time.Now() to satisfy LSP "imported and not used" warning
	_ = time.Now()

	// Simulate catalog by creating a slice with value types
	modelsSlice := []models.ModelInfo{testModel1, testModel2}
	testCatalog := struct {
		Models []models.ModelInfo
	}
	testCatalog.Models = modelsSlice

	// Test filtering
	openaiModels := testCatalog.GetProviderModels("openai")
	if len(openaiModels) != 1 {
		t.Errorf("Expected 1 OpenAI model, got %d", len(openaiModels))
	}

	// Test model properties
	if len(openaiModels) > 0 {
		if openaiModels[0].ID != "model-1" {
			t.Errorf("Expected model ID to be model-1, got %s", openaiModels[0].ID)
		}
		if openaiModels[0].Name != "Test Model 1" {
			t.Errorf("Expected model Name to be Test Model 1, got %s", openaiModels[0].Name)
		}
		if openaiModels[0].Provider != "openai" {
			t.Errorf("Expected model Provider to be openai, got %s", openaiModels[0].Provider)
		}
	}

	anthropicModels := testCatalog.GetProviderModels("anthropic")
	if len(anthropicModels) != 1 {
		t.Errorf("Expected 1 Anthropic model, got %d", len(anthropicModels))
	}
}

	// Test filtering by existing provider
	openaiModels := catalog.GetProviderModels("openai")
	if len(openaiModels) == 0 {
		t.Error("Expected OpenAI models, got none")
	}

	t.Logf("Found %d OpenAI models", len(openaiModels))
}

// TestCatalog_GetModelByID tests model lookup
func TestCatalog_GetModelByID(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	store, _ := store.New(tmpDB)
	defer store.Close()

	catalog, _ := models.LoadCatalog(store)

	model := catalog.GetModelByID("nonexistent-model")
	if model != nil {
		t.Error("Expected model to be nil for non-existent ID")
	}

	if catalog.GetModelByID("gpt-4") == nil {
		t.Error("Expected to find GPT-4 model")
	}
}

// TestCatalog_IsStale tests staleness detection
func TestCatalog_IsStale(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	defer os.Remove(tmpDB)

	store, _ := store.New(tmpDB)
	defer store.Close()

	catalog, _ := models.LoadCatalog(store)

	if catalog.IsStale() {
		t.Log("Catalog is stale (as expected)")
	} else {
		t.Error("Expected catalog to be stale immediately after load")
	}
}

// TestCatalog_GetPricing tests pricing retrieval
func TestCatalog_GetPricing(t *testing.T) {
	tmpDB := t.TempDir() + "/test.db"
	defer os.Remove(tmpDB)

	store, _ := store.New(tmpDB)
	defer store.Close()

	catalog, _ := models.LoadCatalog(store)

	pricing := catalog.GetPricing()
	if pricing == nil {
		t.Fatal("Pricing should not be nil")
	}

	if len(pricing.Models) == 0 {
		t.Error("Expected pricing to have models")
	}

	t.Logf("Pricing loaded with %d models", len(pricing.Models))
}
