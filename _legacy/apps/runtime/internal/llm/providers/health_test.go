package providers

import (
	"context"
	"testing"
	"time"
)

func TestNewHealthChecker(t *testing.T) {
	hc := NewHealthChecker()
	if hc == nil {
		t.Fatal("NewHealthChecker() returned nil")
	}
	if hc.results == nil {
		t.Error("HealthChecker.results is nil")
	}
	if hc.client == nil {
		t.Error("HealthChecker.client is nil")
	}
}

func TestHealthChecker_CheckProvider_Unsupported(t *testing.T) {
	hc := NewHealthChecker()
	ctx := context.Background()

	health, err := hc.CheckProvider(ctx, "unknown", "test-key", "")
	if err != nil {
		t.Fatalf("CheckProvider failed: %v", err)
	}

	if health.Status != StatusError {
		t.Errorf("Expected status %s, got %s", StatusError, health.Status)
	}

	if health.ProviderID != "unknown" {
		t.Errorf("Expected providerID 'unknown', got %s", health.ProviderID)
	}
}

func TestHealthChecker_GetHealth(t *testing.T) {
	hc := NewHealthChecker()

	// Should return false for unknown provider
	_, ok := hc.GetHealth("unknown")
	if ok {
		t.Error("GetHealth should return false for unknown provider")
	}

	// Add a health result
	hc.mu.Lock()
	hc.results["test"] = &ProviderHealth{
		ProviderID:  "test",
		Status:      StatusHealthy,
		LastChecked: time.Now(),
	}
	hc.mu.Unlock()

	// Should return true now
	health, ok := hc.GetHealth("test")
	if !ok {
		t.Error("GetHealth should return true for existing provider")
	}
	if health.ProviderID != "test" {
		t.Errorf("Expected providerID 'test', got %s", health.ProviderID)
	}
}

func TestHealthChecker_GetAllHealth(t *testing.T) {
	hc := NewHealthChecker()

	// Add multiple results
	hc.mu.Lock()
	hc.results["openai"] = &ProviderHealth{
		ProviderID:  "openai",
		Status:      StatusHealthy,
		LastChecked: time.Now(),
	}
	hc.results["anthropic"] = &ProviderHealth{
		ProviderID:  "anthropic",
		Status:      StatusHealthy,
		LastChecked: time.Now(),
	}
	hc.mu.Unlock()

	all := hc.GetAllHealth()
	if len(all) != 2 {
		t.Errorf("Expected 2 health results, got %d", len(all))
	}

	if _, ok := all["openai"]; !ok {
		t.Error("Expected 'openai' in results")
	}
	if _, ok := all["anthropic"]; !ok {
		t.Error("Expected 'anthropic' in results")
	}
}

func TestProviderHealth_IsStale(t *testing.T) {
	health := &ProviderHealth{
		ProviderID:  "test",
		Status:      StatusHealthy,
		LastChecked: time.Now().Add(-10 * time.Minute),
	}

	if !health.IsStale(5 * time.Minute) {
		t.Error("Expected IsStale to return true for 10-minute-old check with 5-minute threshold")
	}

	health.LastChecked = time.Now()
	if health.IsStale(5 * time.Minute) {
		t.Error("Expected IsStale to return false for fresh check")
	}
}
