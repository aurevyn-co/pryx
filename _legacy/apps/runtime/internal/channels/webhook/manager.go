package webhook

import (
	"context"
	"sync"

	"pryx-core/internal/bus"
)

// Manager manages webhook channels
type Manager struct {
	channels map[string]*Channel
	config   *ConfigManager
	bus      *bus.Bus
	mu       sync.RWMutex
}

// NewManager creates a new webhook manager
func NewManager(eventBus *bus.Bus) *Manager {
	return &Manager{
		channels: make(map[string]*Channel),
		config:   NewConfigManager(),
		bus:      eventBus,
	}
}

// LoadChannels loads all webhook channels from config
func (m *Manager) LoadChannels() error {
	configs, err := m.config.LoadAll()
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		channel := NewChannel(cfg, m.bus)
		m.channels[cfg.ID] = channel

		// Auto-connect if configured with a port
		if cfg.Port > 0 {
			ctx := context.Background()
			go channel.Connect(ctx)
		}
	}

	return nil
}

// GetChannel returns a webhook channel by ID
func (m *Manager) GetChannel(id string) *Channel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.channels[id]
}

// ListChannels returns all webhook channels
func (m *Manager) ListChannels() []*Channel {
	m.mu.RLock()
	defer m.mu.RUnlock()

	channels := make([]*Channel, 0, len(m.channels))
	for _, ch := range m.channels {
		channels = append(channels, ch)
	}
	return channels
}

// CreateChannel creates a new webhook channel
func (m *Manager) CreateChannel(config WebhookConfig) (*Channel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg, err := m.config.Create(config)
	if err != nil {
		return nil, err
	}

	channel := NewChannel(*cfg, m.bus)
	m.channels[cfg.ID] = channel

	// Auto-connect if configured with a port
	if cfg.Port > 0 && cfg.Enabled {
		ctx := context.Background()
		go channel.Connect(ctx)
	}

	return channel, nil
}

// UpdateChannel updates an existing webhook channel
func (m *Manager) UpdateChannel(id string, config WebhookConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	channel, exists := m.channels[id]
	if !exists {
		return nil
	}

	// Disconnect old channel
	ctx := context.Background()
	channel.Disconnect(ctx)

	// Update config
	config.ID = id
	if err := m.config.Save(config); err != nil {
		return err
	}

	// Create new channel with updated config
	newChannel := NewChannel(config, m.bus)
	m.channels[id] = newChannel

	// Reconnect if enabled
	if config.Enabled && config.Port > 0 {
		go newChannel.Connect(ctx)
	}

	return nil
}

// DeleteChannel deletes a webhook channel
func (m *Manager) DeleteChannel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	channel, exists := m.channels[id]
	if !exists {
		return nil
	}

	// Disconnect channel
	ctx := context.Background()
	channel.Disconnect(ctx)

	// Remove from config
	if err := m.config.Delete(id); err != nil {
		return err
	}

	// Remove from manager
	delete(m.channels, id)

	return nil
}

// Shutdown gracefully shuts down all channels
func (m *Manager) Shutdown() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()
	for _, channel := range m.channels {
		channel.Disconnect(ctx)
	}

	return nil
}
