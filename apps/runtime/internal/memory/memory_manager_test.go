package memory_test

import (
	"context"
	"testing"
	"time"

	"pryx-core/internal/memory"
)

func TestNewManager(t *testing.T) {
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	if manager == nil {
		t.Error("Expected non-nil manager")
	}
}

func TestGetMemoryUsage_Basic(t *testing.T) {
	ctx := context.Background()
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	// Mock session
	session := &Session{
		ID:          "test-session",
		MessagesCount: 100,
		MaxTokens:   128000,
	}

	store.GetSessionFunc = func(ctx context.Context, sessionID string) (*Session, error) {
		return session, nil
	}
	store.GetSessionConstraintsFunc = func(ctx context.Context, sessionID string) (*constraints.SessionConstraints, error) {
		return &constraints.SessionConstraints{
			MaxTokens: 128000,
		}, nil
	}

	usage, err := manager.GetMemoryUsage(ctx, "test-session")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if usage.UsagePercent != 78.125 {
		t.Errorf("Expected 78.125%% usage, got: %.2f%%", usage.UsagePercent)
	}

	if usage.WarningLevel != "info" {
		t.Errorf("Expected info warning level, got: %s", usage.WarningLevel)
	}
}

func TestGetMemoryUsage_OverLimit(t *testing.T) {
	ctx := context.Background()
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	session := &Session{
		ID:          "test-session",
		MessagesCount: 128000,
		MaxTokens:   128000,
	}

	store.GetSessionFunc = func(ctx context.Context, sessionID string) (*Session, error) {
		return session, nil
	}
	store.GetSessionConstraintsFunc = func(ctx context.Context, sessionID string) (*constraints.SessionConstraints, error) {
		return &constraints.SessionConstraints{
			MaxTokens: 128000,
		}, nil
	}

	usage, err := manager.GetMemoryUsage(ctx, "test-session")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if usage.UsagePercent != 100 {
		t.Errorf("Expected 100%% usage, got: %.2f%%", usage.UsagePercent)
	}

	if usage.WarningLevel != "critical" {
		t.Errorf("Expected critical warning level, got: %s", usage.WarningLevel)
	}
}

func TestSummarizeSession_Basic(t *testing.T) {
	ctx := context.Background()
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	session := &Session{
		ID:          "test-session",
		Title:       "Test Session",
		CreatedAt:   time.Now(),
		Messages:    []Message{},
	}

	archivedCount := 0

	// Set up mock for compression
	store.ListMessagesFunc = func(ctx context.Context, sessionID string, offset, limit int, before time.Time, filter interface{}) ([]Message, error) {
		return append(make([]Message, 20), Message{
			ID:         fmt.Sprintf("msg-%d", i),
			Tokens:     10,
			CreatedAt: time.Now(),
		}), nil
	}
	store.ArchiveMessageFunc = func(ctx context.Context, msgID string) error {
		return nil
	}
	store.UpdateSessionFunc = func(ctx context.Context, sessionID string, updates map[string]interface{}) error {
		if compressedTokens, ok := updates["compressed_tokens"].(int64); ok {
			session.CompressedTokens = compressedTokens
		}
		return nil
	}

	result, err := manager.SummarizeSession(ctx, "test-session")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.CompressedCount != 10 {
		t.Errorf("Expected 10 compressed messages, got: %d", result.CompressedCount)
	}
}

func TestArchiveSession(t *testing.T) {
	ctx := context.Background()
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	session := &Session{
		ID:          "test-session",
		Title:       "Test Session",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		Messages:    []Message{},
	}

	store.UpdateSessionFunc = func(ctx context.Context, sessionID string, updates map[string]interface{}) error {
		return nil
	}

	archivedCount := 0

	result, err := manager.ArchiveSession(ctx, "test-session")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.ArchivedCount != 1 {
		t.Errorf("Expected 1 archived session, got: %d", result.ArchivedCount)
	}

	if len(result.ArchivedSessions) != 1 {
		t.Errorf("Expected 1 archived session ID, got %d", len(result.ArchivedSessions))
	}
}

func TestCreateChildSession(t *testing.T) {
	ctx := context.Background()
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	store.CreateSessionFunc = func(ctx context.Context, parentID string, title string) (string, error) {
		return "child-session", nil
	}
	store.GetSessionFunc = func(ctx context.Context, sessionID string) (*Session, error) {
		return &Session{
			ID:          "child-session",
			ParentSessionID: "test-session",
			Title:       "Child Session",
			CreatedAt:   time.Now(),
		}, nil
	}

	sessionID, err := manager.CreateChildSession(ctx, "test-session", "Child Session")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if sessionID != "child-session" {
		t.Errorf("Expected child-session ID, got: %s", sessionID)
	}
}

func TestGetSessionMemory(t *testing.T) {
	ctx := context.Background()
	cfg := constraints.Manager{}
	store := &mockStore{}
	bus := &mockBus{}

	manager := memory.NewManager(&cfg, store, bus)

	session := &Session{
		ID:          "test-session",
		Title:       "Test Session",
		CreatedAt:   time.Now(),
		Messages:    []Message{
			{ID: "msg-1", Tokens: 10, CreatedAt: time.Now()},
			{ID: "msg-2", Tokens: 10, CreatedAt: time.Now()},
			{ID: "msg-3", Tokens: 10, CreatedAt: time.Now()},
		},
	}

	store.ListMessagesFunc = func(ctx context.Context, sessionID string, offset, limit int, before time.Time, filter interface{}) ([]Message, error) {
		return session.Messages, nil
	}

	memory, err := manager.GetSessionMemory(ctx, "test-session")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if memory.TotalTokens != 30 {
		t.Errorf("Expected 30 total tokens, got: %d", memory.TotalTokens)
	}
}

// Mock types for testing
type mockStore struct{}
type mockBus struct{}

func (m *mockStore) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	return &Session{ID: sessionID}, nil
}

func (m *mockStore) GetSessionConstraints(ctx context.Context, sessionID string) (*constraints.SessionConstraints, error) {
	return &constraints.SessionConstraints{MaxTokens: 128000}, nil
}

func (m *mockStore) ListMessages(ctx context.Context, sessionID string, offset, limit int, before time.Time, filter interface{}) ([]memory.Message, error) {
	return []memory.Message{}, nil
}

func (m *mockStore) UpdateSession(ctx context.Context, sessionID string, updates map[string]interface{}) error {
	return nil
}

func (m *mockStore) ArchiveMessage(ctx context.Context, msgID string) error {
	return nil
}

func (m *mockBus) Publish(event bus.Event) {
	// No-op for testing
}

func (m *mockBus) NewEvent(eventType string, id string, data map[string]interface{}) bus.Event {
	return bus.Event{Type: eventType, ID: id, Data: data}
}
