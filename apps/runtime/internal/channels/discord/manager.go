package discord

import (
	"context"
	"fmt"
	"sync"

	"pryx-core/internal/bus"
	"pryx-core/internal/channels"
)

// Manager manages Discord channel instances
type Manager struct {
	mu          sync.RWMutex
	channels    map[string]*DiscordChannel
	configMgr   *ConfigManager
	eventBus    *bus.Bus
	healthCheck *HealthChecker
}

// NewManager creates a new Discord channel manager
func NewManager(eventBus *bus.Bus) *Manager {
	return &Manager{
		channels:    make(map[string]*DiscordChannel),
		configMgr:   NewConfigManager(),
		eventBus:    eventBus,
		healthCheck: NewHealthChecker(),
	}
}

// LoadConfigs loads all Discord configurations and creates channels
func (m *Manager) LoadConfigs(ctx context.Context) error {
	configs, err := m.configMgr.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load configs: %w", err)
	}

	for _, config := range configs {
		if !config.Enabled {
			continue
		}

		if err := m.CreateChannel(ctx, config); err != nil {
			m.publishError(config.ID, fmt.Sprintf("failed to create channel: %v", err))
		}
	}

	return nil
}

// CreateChannel creates and registers a new Discord channel
func (m *Manager) CreateChannel(ctx context.Context, config Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.channels[config.ID]; exists {
		return fmt.Errorf("channel %s already exists", config.ID)
	}

	// Get token from config (in real usage, load from vault)
	token := config.Token
	if token == "" {
		return fmt.Errorf("bot token is required")
	}

	channel := NewDiscordChannel(config.ID, token, m.eventBus)

	m.channels[config.ID] = channel

	// Connect the channel
	if err := channel.Connect(ctx); err != nil {
		delete(m.channels, config.ID)
		return fmt.Errorf("failed to connect channel: %w", err)
	}

	m.publishStatus(config.ID, "created", nil)

	return nil
}

// GetChannel retrieves a channel by ID
func (m *Manager) GetChannel(id string) (*DiscordChannel, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ch, ok := m.channels[id]
	return ch, ok
}

// GetAllChannels returns all managed channels
func (m *Manager) GetAllChannels() []*DiscordChannel {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*DiscordChannel, 0, len(m.channels))
	for _, ch := range m.channels {
		result = append(result, ch)
	}
	return result
}

// RemoveChannel removes and disconnects a channel
func (m *Manager) RemoveChannel(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	channel, exists := m.channels[id]
	if !exists {
		return fmt.Errorf("channel %s not found", id)
	}

	if err := channel.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect channel: %w", err)
	}

	delete(m.channels, id)
	m.publishStatus(id, "removed", nil)

	return nil
}

// Shutdown disconnects all channels
func (m *Manager) Shutdown(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, channel := range m.channels {
		if err := channel.Disconnect(ctx); err != nil {
			m.publishError(id, fmt.Sprintf("failed to disconnect: %v", err))
		}
	}

	m.channels = make(map[string]*DiscordChannel)
}

// GetHealthStatus returns health status for all channels
func (m *Manager) GetHealthStatus() map[string]HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]HealthStatus)
	for id, channel := range m.channels {
		result[id] = channel.Health()
	}
	return result
}

// CheckHealth performs health checks on all channels
func (m *Manager) CheckHealth(ctx context.Context) map[string]HealthResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]HealthResult)
	for id, channel := range m.channels {
		result[id] = m.healthCheck.Check(ctx, channel)
	}
	return result
}

// SendMessage sends a message through a specific channel
func (m *Manager) SendMessage(ctx context.Context, channelID string, channelID2 string, text string) error {
	channel, exists := m.GetChannel(channelID)
	if !exists {
		return fmt.Errorf("channel %s not found", channelID)
	}

	msg := channels.Message{
		ChannelID: channelID2,
		Content:   text,
	}

	return channel.Send(ctx, msg)
}

// RegisterCommands registers slash commands for a channel
func (m *Manager) RegisterCommands(ctx context.Context, channelID string, cmds []ApplicationCommand) error {
	channel, exists := m.GetChannel(channelID)
	if !exists {
		return fmt.Errorf("channel %s not found", channelID)
	}

	return channel.RegisterCommands(ctx, cmds)
}

// GetChannelCount returns the number of managed channels
func (m *Manager) GetChannelCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.channels)
}

// IsChannelConnected checks if a channel is connected
func (m *Manager) IsChannelConnected(channelID string) bool {
	channel, exists := m.GetChannel(channelID)
	if !exists {
		return false
	}
	return channel.Status() == channels.StatusConnected
}

// publishStatus publishes a status event
func (m *Manager) publishStatus(channelID, status string, data interface{}) {
	if m.eventBus == nil {
		return
	}

	payload := map[string]interface{}{
		"channel_id": channelID,
		"status":     status,
		"type":       "discord",
	}

	if data != nil {
		payload["data"] = data
	}

	m.eventBus.Publish(bus.NewEvent(bus.EventChannelStatus, "", payload))
}

// publishError publishes an error event
func (m *Manager) publishError(channelID, errMsg string) {
	if m.eventBus == nil {
		return
	}

	m.eventBus.Publish(bus.NewEvent(bus.EventErrorOccurred, "", map[string]interface{}{
		"channel_id": channelID,
		"error":      errMsg,
		"type":       "discord",
	}))
}
