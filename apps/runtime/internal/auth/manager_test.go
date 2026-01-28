package auth_test

import (
	"context"
	"testing"
	"time"

	"pryx-core/internal/auth"
	"pryx-core/internal/bus"
	"pryx-core/internal/config"
	"pryx-core/internal/keychain"
)

func TestNewManager(t *testing.T) {
	cfg := &config.AuthConfig{
		OAuthProviders: map[string]*config.OAuthProvider{},
	}
	kc := &mockKeychain{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := auth.NewManager(cfg, kc, store, bus)

	if manager == nil {
		t.Error("Expected non-nil manager")
	}
}

func TestInitiateDeviceFlow(t *testing.T) {
	ctx := context.Background()
	cfg := &config.AuthConfig{
		OAuthProviders: map[string]*config.OAuthProvider{
			"test_provider": {
				Name:     "Test Provider",
				ClientID: "test_client",
				ClientSecret: "test_secret",
				AuthURL:    "https://auth.example.com/authorize",
				TokenURL:     "https://auth.example.com/token",
				Scopes:      []string{"read", "write"},
			},
		},
	}

	kc := &mockKeychain{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := auth.NewManager(cfg, kc, store, bus)

	redirectURI := "pryx://callback/test"

	state, err := manager.InitiateDeviceFlow(ctx, "test_provider", redirectURI)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if state.ProviderID != "test_provider" {
		t.Errorf("Expected test_provider, got: %s", state.ProviderID)
	}

	if state.ClientID != "test_client" {
		t.Errorf("Expected test_client, got: %s", state.ClientID)
	}
}

func TestSetManualToken(t *testing.T) {
	ctx := context.Background()
	cfg := &config.AuthConfig{
		OAuthProviders: map[string]*config.OAuthProvider{},
	}

	kc := &mockKeychain{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := auth.NewManager(cfg, kc, store, bus)

	err := manager.SetManualToken(ctx, "test_provider", "manual_token_value")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Mock types
type mockKeychain struct {
	Store map[string]string
}

type mockStore struct {
	Store map[string]*auth.OAuthState
}

func (m *mockKeychain) Get(ctx context.Context, key string) (string, error) {
	return m.Store[key], nil
}

func (m *mockKeychain) Set(ctx context.Context, key string, value string) error {
	m.Store[key] = value
	return nil
}

func (m *mockKeychain) Delete(ctx context.Context, key string) error {
	delete(m.Store, key)
	return nil
}

func (m *mockStore) GetOAuthState(ctx context.Context, state string) (*auth.OAuthState, error) {
	return m.Store[state], nil
}

func (m *mockBus) NewEvent(eventType bus.EventType, id string, data map[string]interface{}) bus.Event {
	return bus.Event{Type: eventType, ID: id, Data: data}
}

func (m *mockBus) Publish(event bus.Event) {
	// No-op for tests
}
