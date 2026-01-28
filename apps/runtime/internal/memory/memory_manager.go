package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pryx-core/internal/bus"
	"pryx-core/internal/store"
	"pryx-core/internal/constraints"
)

const (
	WarnThresholdPercent  = 80
	SummarizeThresholdPercent = 90
	CompressionRatio           = 0.2 // Compress oldest 20%
	MaxContextTokens      = 128000 // Example: GPT-4 context limit
	SessionArchiveDays     = 7    // Archive sessions older than 7 days
)

type MemoryUsage struct {
	UsedTokens      int     `json:"used_tokens"`
	MaxTokens       int     `json:"max_tokens"`
	UsagePercent     float64 `json:"usage_percent"`
	WarningLevel    string  `json:"warning_level"`
}

type SessionMemory struct {
	SessionID       string    `json:"session_id"`
	ParentSessionID string    `json:"parent_session_id,omitempty"`
	Title           string    `json:"title"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	MessagesCount   int       `json:"messages_count"`
	TotalTokens     int       `json:"total_tokens"`
	CompressedTokens int       `json:"compressed_tokens"`
	Archived        bool      `json:"archived"`
}

type SummaryRequest struct {
	OldestMessages    int `json:"oldest_messages"`
	CompressionRatio float64 `json:"compression_ratio"`
}

type CompressionResult struct {
	CompressedCount int        `json:"compressed_count"`
	NewTotalTokens  int        `json:"new_total_tokens"`
	SavedTokens      int       `json:"saved_tokens"`
}

type ArchiveResult struct {
	ArchivedCount int `json:"archived_count"`
	ArchivedSessions []string `json:"archived_sessions"`
}

type Manager struct {
	config       *constraints.Manager
	store        *store.Store
	bus           *bus.Bus
}

func NewManager(cfg *constraints.Manager, store *store.Store, bus *bus.Bus) *Manager {
	return &Manager{
		config: cfg,
		store:  store,
		bus:  bus,
	}
}

func (m *Manager) GetMemoryUsage(ctx context.Context, sessionID string) (MemoryUsage, error) {
	session, err := m.store.GetSession(ctx, sessionID)
	if err != nil {
		return MemoryUsage{}, err
	}

	// Get constraints for this session
	sessionConstraints, err := m.config.GetSessionConstraints(ctx, sessionID)
	if err != nil {
		return MemoryUsage{
			UsedTokens:    session.MessagesCount,
			MaxTokens:     sessionConstraints.MaxTokens,
			UsagePercent:   0,
			WarningLevel:  "",
		}
	}

	// Calculate token usage
	usagePercent := 0.0
	if sessionConstraints.MaxTokens > 0 {
		usagePercent = float64(session.MessagesCount) / float64(sessionConstraints.MaxTokens) * 100.0
	}

	// Determine warning level
	warningLevel := ""
	if usagePercent >= 100 {
		warningLevel = "critical"
	} else if usagePercent >= float64(SummarizeThresholdPercent) {
		warningLevel = "warn"
	} else if usagePercent >= float64(WarnThresholdPercent) {
		warningLevel = "info"
	}

	return MemoryUsage{
		UsedTokens:    session.MessagesCount,
		MaxTokens:     sessionConstraints.MaxTokens,
		UsagePercent:   usagePercent,
		WarningLevel:    warningLevel,
	}
}

func (m *Manager) CheckAndWarn(ctx context.Context, sessionID string) error {
	usage, err := m.GetMemoryUsage(ctx, sessionID)
	if err != nil {
		return err
	}

	// Publish warning if needed
	if usage.UsagePercent >= float64(WarnThresholdPercent) {
		if m.bus != nil {
			event := bus.NewEvent(bus.EventMemoryWarning, sessionID, map[string]interface{}{
				"usage_percent":   usage.UsagePercent,
				"used_tokens":     usage.UsedTokens,
				"max_tokens":      usage.MaxTokens,
				"warning_level":    usage.WarningLevel,
			})
			m.bus.Publish(event)
		}
	}

	if usage.UsagePercent >= float64(SummarizeThresholdPercent) {
		if m.bus != nil {
			event := bus.NewEvent(bus.EventMemorySummarizeRequest, sessionID, map[string]interface{}{
				"oldest_messages": usage.UsedTokens / 2,
				"compression_ratio": CompressionRatio,
			})
			m.bus.Publish(event)
		}
	}

	return nil
}

func (m *Manager) SummarizeSession(ctx context.Context, sessionID string) (*CompressionResult, error) {
	messages, err := m.store.ListMessages(ctx, sessionID, 0, 50, nil) // Get last 50 messages
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return &CompressionResult{
			CompressedCount: 0,
			NewTotalTokens:  0,
			SavedTokens:      0,
		}, nil
	}

	// Compress oldest messages (20% of messages)
	compressCount := int(float64(len(messages)) * CompressionRatio)
	if compressCount == 0 {
		compressCount = 1
	}

	// Create summary of compressed messages
	var summary string
	totalTokens := 0
	for i, msg := range messages {
		if i >= compressCount {
			totalTokens += msg.Tokens
		}
	}

	summary = fmt.Sprintf("Compressed %d messages (%d tokens)", compressCount, totalTokens)

	// Archive original messages
	for i, msg := range messages {
		if i < compressCount {
			_ = m.store.ArchiveMessage(ctx, msg.ID)
		}
	}

	// Update session metadata
	_ = m.store.UpdateSession(ctx, sessionID, map[string]interface{}{
		"compressed_at":      time.Now(),
		"compressed_tokens": totalTokens,
	})

	if m.bus != nil {
		event := bus.NewEvent(bus.EventMemorySummarized, sessionID, map[string]interface{}{
			"compressed_count": compressCount,
			"saved_tokens":      totalTokens,
			"summary":          summary,
		})
		m.bus.Publish(event)
	}

	return &CompressionResult{
		CompressedCount: compressCount,
		NewTotalTokens:  0,
		SavedTokens:      totalTokens,
	}, nil
}

func (m *Manager) ArchiveSession(ctx context.Context, sessionID string) (*ArchiveResult, error) {
	session, err := m.store.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Mark session as archived
	err = m.store.UpdateSession(ctx, sessionID, map[string]interface{}{
		"archived": true,
		"archived_at": time.Now(),
	})

	if err != nil {
		return nil, err
	}

	archivedSessions := []string{sessionID}

	// Archive all child sessions
	childSessions, err := m.store.ListChildSessions(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	for _, child := range childSessions {
		_ = m.store.UpdateSession(ctx, child.ID, map[string]interface{}{
			"archived": true,
			"archived_at": time.Now(),
		})
		archivedSessions = append(archivedSessions, child.ID)
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventSessionArchived, sessionID, map[string]interface{}{
			"archived_count": len(archivedSessions) + 1,
			"archived_sessions": archivedSessions,
		})
		m.bus.Publish(event)
	}

	return &ArchiveResult{
		ArchivedCount: len(archivedSessions) + 1,
		ArchivedSessions: archivedSessions,
	}, nil
}

func (m *Manager) UnarchiveSession(ctx context.Context, sessionID string) error {
	return m.store.UpdateSession(ctx, sessionID, map[string]interface{}{
		"archived": false,
		"archived_at": nil,
	})
}

func (m *Manager) CreateChildSession(ctx context.Context, parentSessionID, title string) (string, error) {
	sessionID, err := m.store.CreateSession(ctx, parentSessionID, title)
	if err != nil {
		return "", err
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventSessionCreated, sessionID, map[string]interface{}{
			"parent_session_id": parentSessionID,
			"title":            title,
		})
		m.bus.Publish(event)
	}

	return sessionID, nil
}

func (m *Manager) GetSessionMemory(ctx context.Context, sessionID string) (SessionMemory, error) {
	session, err := m.store.GetSession(ctx, sessionID)
	if err != nil {
		return SessionMemory{}, err
	}

	messages, err := m.store.ListMessages(ctx, sessionID, 0, 100, nil)
	if err != nil {
		return SessionMemory{SessionID: sessionID}, nil
	}

	totalTokens := 0
	for _, msg := range messages {
		totalTokens += msg.Tokens
	}

	childSessions, err := m.store.ListChildSessions(ctx, sessionID)
	if err != nil {
		return SessionMemory{
			SessionID:       sessionID,
			Title:           session.Title,
			CreatedAt:       session.CreatedAt,
			UpdatedAt:       session.UpdatedAt,
			MessagesCount:   len(messages),
			TotalTokens:     totalTokens,
			CompressedTokens: session.CompressedTokens,
			Archived:        session.Archived,
		}, nil
	}
}

func (m *Manager) GetAllSessionsMemory(ctx context.Context) ([]SessionMemory, error) {
	sessions, err := m.store.ListSessions(ctx, nil)
	if err != nil {
		return nil, err
	}

	var sessionMemories []SessionMemory
	for _, session := range sessions {
		messages, err := m.store.ListMessages(ctx, session.ID, 0, 50, nil)
		if err != nil {
			continue
		}

		totalTokens := 0
		for _, msg := range messages {
			totalTokens += msg.Tokens
		}

		sessionMemories = append(sessionMemories, SessionMemory{
			SessionID:       session.ID,
			Title:           session.Title,
			CreatedAt:       session.CreatedAt,
			UpdatedAt:       session.UpdatedAt,
			MessagesCount:   len(messages),
			TotalTokens:     totalTokens,
			CompressedTokens: session.CompressedTokens,
			Archived:        session.Archived,
		})
	}

	return sessionMemories, nil
}

func (m *Manager) AutoManageMemory(ctx context.Context, sessionID string) error {
	usage, err := m.GetMemoryUsage(ctx, sessionID)
	if err != nil {
		return err
	}

	// Auto-summarize at threshold
	if usage.UsagePercent >= float64(SummarizeThresholdPercent) {
		_, err := m.SummarizeSession(ctx, sessionID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) CleanupOldSessions(ctx context.Context) (int, error) {
	sessions, err := m.store.ListSessions(ctx, nil)
	if err != nil {
		return 0, err
	}

	archiveThreshold := time.Now().AddDate(-SessionArchiveDays * 24 * time.Hour)
	archivedCount := 0

	for _, session := range sessions {
		// Don't archive already archived sessions
		if session.Archived {
			continue
		}

		// Archive sessions older than threshold
		if session.UpdatedAt.Before(archiveThreshold) {
			_, err := m.ArchiveSession(ctx, session.ID)
			if err == nil {
				archivedCount++
			}
		}
	}

	// Publish cleanup event
	if m.bus != nil {
		event := bus.NewEvent(bus.EventSessionsCleaned, "", map[string]interface{}{
			"archived_count": archivedCount,
		})
		m.bus.Publish(event)
	}

	return archivedCount, nil
}

func (m *Manager) QueryRAG(ctx context.Context, sessionID, query string) (map[string]interface{}, error) {
	// Placeholder for RAG integration
	// TODO: Implement actual RAG backend integration
	
	response := map[string]interface{}{
		"query":    query,
		"session_id": sessionID,
		"results":  []map[string]interface{}{
			"source": "context_memory",
			"content": fmt.Sprintf("Found %d relevant messages for query", 10),
		},
	},
	}

	return response, nil
}

func compressMessages(messages []Message, compressCount int) ([]Message, int) {
	// Create compressed summary message
	compressedContent := ""
	for i, msg := range messages {
		if i < compressCount {
			compressedContent += msg.Content + "\n\n"
		}
	}

	summaryMsg := Message{
		Role:    "system",
		Content: fmt.Sprintf("[Compressed %d messages]", compressCount),
		Tokens: 0, // System messages don't count
		CreatedAt: time.Now(),
	}

	return append([]Message{summaryMsg}, messages[:compressCount]...), compressCount
}

func estimateTokens(text string) int {
	// Simple estimation: ~4 chars per token
	return len(text) / 4
}
