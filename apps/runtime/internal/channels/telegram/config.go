package telegram

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultConfigDir       = ".pryx/config"
	defaultConfigFile      = "telegram.json"
	defaultPollingInterval = 30 * time.Second
)

// Config represents the Telegram bot configuration
type Config struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	TokenRef              string        `json:"token_ref"`       // Reference to token in vault
	Token                 string        `json:"token,omitempty"` // Token value (loaded from vault, not persisted)
	Mode                  string        `json:"mode"`            // "polling" or "webhook"
	WebhookURL            string        `json:"webhook_url,omitempty"`
	WebhookSecret         string        `json:"webhook_secret,omitempty"` // Secret for webhook validation
	PollingInterval       time.Duration `json:"polling_interval"`
	AllowedChats          []int64       `json:"allowed_chats"`   // Whitelist of chat IDs
	AllowedUpdates        []string      `json:"allowed_updates"` // Types of updates to receive
	MaxConnections        int           `json:"max_connections"` // For webhook mode
	DropPendingUpdates    bool          `json:"drop_pending_updates"`
	Commands              []BotCommand  `json:"commands"`   // Bot commands to register
	ParseMode             ParseMode     `json:"parse_mode"` // Default parse mode
	DisableWebPagePreview bool          `json:"disable_web_page_preview"`
	DisableNotification   bool          `json:"disable_notification"`
	Enabled               bool          `json:"enabled"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("config ID is required")
	}

	if c.Name == "" {
		return fmt.Errorf("config name is required")
	}

	if c.TokenRef == "" {
		return fmt.Errorf("token reference is required")
	}

	if c.Mode != "polling" && c.Mode != "webhook" {
		return fmt.Errorf("mode must be 'polling' or 'webhook', got: %s", c.Mode)
	}

	if c.Mode == "webhook" && c.WebhookURL == "" {
		return fmt.Errorf("webhook URL is required in webhook mode")
	}

	if c.PollingInterval <= 0 {
		c.PollingInterval = defaultPollingInterval
	}

	if c.MaxConnections <= 0 {
		c.MaxConnections = 40
	}

	return nil
}

// IsChatAllowed checks if a chat ID is in the whitelist
func (c *Config) IsChatAllowed(chatID int64) bool {
	if len(c.AllowedChats) == 0 {
		return true // No whitelist means all chats allowed
	}

	for _, id := range c.AllowedChats {
		if id == chatID {
			return true
		}
	}

	return false
}

// SetDefaults sets default values for optional fields
func (c *Config) SetDefaults() {
	if c.Mode == "" {
		c.Mode = "polling"
	}

	if c.PollingInterval <= 0 {
		c.PollingInterval = defaultPollingInterval
	}

	if c.MaxConnections <= 0 {
		c.MaxConnections = 40
	}

	if c.ParseMode == "" {
		c.ParseMode = ParseModeMarkdown
	}

	if len(c.AllowedUpdates) == 0 {
		c.AllowedUpdates = []string{
			"message",
			"edited_message",
			"callback_query",
			"my_chat_member",
		}
	}
}

// ConfigManager manages Telegram bot configurations
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new config manager
func NewConfigManager() *ConfigManager {
	home, _ := os.UserHomeDir()
	return &ConfigManager{
		configPath: filepath.Join(home, defaultConfigDir, defaultConfigFile),
	}
}

// NewConfigManagerWithPath creates a config manager with a custom path
func NewConfigManagerWithPath(path string) *ConfigManager {
	return &ConfigManager{
		configPath: path,
	}
}

// LoadAll loads all Telegram configurations
func (cm *ConfigManager) LoadAll() ([]Config, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configs []Config
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults for loaded configs
	for i := range configs {
		configs[i].SetDefaults()
	}

	return configs, nil
}

// SaveAll saves all Telegram configurations
func (cm *ConfigManager) SaveAll(configs []Config) error {
	// Ensure directory exists
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Remove token values before saving (security)
	configsToSave := make([]Config, len(configs))
	for i, config := range configs {
		configsToSave[i] = config
		configsToSave[i].Token = "" // Never persist tokens to disk
	}

	data, err := json.MarshalIndent(configsToSave, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Get retrieves a configuration by ID
func (cm *ConfigManager) Get(id string) (*Config, error) {
	configs, err := cm.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if config.ID == id {
			return &config, nil
		}
	}

	return nil, fmt.Errorf("telegram config not found: %s", id)
}

// Save saves a single configuration (creates or updates)
func (cm *ConfigManager) Save(config Config) error {
	configs, err := cm.LoadAll()
	if err != nil {
		return err
	}

	now := time.Now()
	config.UpdatedAt = now

	found := false
	for i, c := range configs {
		if c.ID == config.ID {
			// Preserve creation time
			config.CreatedAt = c.CreatedAt
			configs[i] = config
			found = true
			break
		}
	}

	if !found {
		config.CreatedAt = now
		configs = append(configs, config)
	}

	return cm.SaveAll(configs)
}

// Delete removes a configuration by ID
func (cm *ConfigManager) Delete(id string) error {
	configs, err := cm.LoadAll()
	if err != nil {
		return err
	}

	filtered := make([]Config, 0, len(configs))
	found := false
	for _, config := range configs {
		if config.ID != id {
			filtered = append(filtered, config)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("telegram config not found: %s", id)
	}

	return cm.SaveAll(filtered)
}

// Create creates a new configuration with generated ID
func (cm *ConfigManager) Create(config Config) (*Config, error) {
	if config.ID == "" {
		config.ID = generateID()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.SetDefaults()

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	if err := cm.Save(config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Update updates an existing configuration
func (cm *ConfigManager) Update(id string, updates map[string]interface{}) (*Config, error) {
	config, err := cm.Get(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		config.Name = name
	}
	if mode, ok := updates["mode"].(string); ok {
		config.Mode = mode
	}
	if webhookURL, ok := updates["webhook_url"].(string); ok {
		config.WebhookURL = webhookURL
	}
	if tokenRef, ok := updates["token_ref"].(string); ok {
		config.TokenRef = tokenRef
	}
	if token, ok := updates["token"].(string); ok {
		config.Token = token
	}
	if pollingInterval, ok := updates["polling_interval"].(time.Duration); ok {
		config.PollingInterval = pollingInterval
	}
	if allowedChats, ok := updates["allowed_chats"].([]int64); ok {
		config.AllowedChats = allowedChats
	}
	if allowedUpdates, ok := updates["allowed_updates"].([]string); ok {
		config.AllowedUpdates = allowedUpdates
	}
	if maxConnections, ok := updates["max_connections"].(int); ok {
		config.MaxConnections = maxConnections
	}
	if dropPending, ok := updates["drop_pending_updates"].(bool); ok {
		config.DropPendingUpdates = dropPending
	}
	if commands, ok := updates["commands"].([]BotCommand); ok {
		config.Commands = commands
	}
	if parseMode, ok := updates["parse_mode"].(ParseMode); ok {
		config.ParseMode = parseMode
	}
	if disablePreview, ok := updates["disable_web_page_preview"].(bool); ok {
		config.DisableWebPagePreview = disablePreview
	}
	if disableNotification, ok := updates["disable_notification"].(bool); ok {
		config.DisableNotification = disableNotification
	}
	if enabled, ok := updates["enabled"].(bool); ok {
		config.Enabled = enabled
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.UpdatedAt = time.Now()

	if err := cm.Save(*config); err != nil {
		return nil, err
	}

	return config, nil
}

// AddAllowedChat adds a chat ID to the whitelist
func (cm *ConfigManager) AddAllowedChat(configID string, chatID int64) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	// Check if already exists
	for _, id := range config.AllowedChats {
		if id == chatID {
			return nil // Already exists
		}
	}

	config.AllowedChats = append(config.AllowedChats, chatID)
	return cm.Save(*config)
}

// RemoveAllowedChat removes a chat ID from the whitelist
func (cm *ConfigManager) RemoveAllowedChat(configID string, chatID int64) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	filtered := make([]int64, 0, len(config.AllowedChats))
	for _, id := range config.AllowedChats {
		if id != chatID {
			filtered = append(filtered, id)
		}
	}

	config.AllowedChats = filtered
	return cm.Save(*config)
}

// SetToken sets the bot token (not persisted, runtime only)
func (cm *ConfigManager) SetToken(configID string, token string) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	config.Token = token
	// Don't save - token is not persisted
	return nil
}

// List returns all configurations
func (cm *ConfigManager) List() ([]Config, error) {
	return cm.LoadAll()
}

// ListEnabled returns only enabled configurations
func (cm *ConfigManager) ListEnabled() ([]Config, error) {
	configs, err := cm.LoadAll()
	if err != nil {
		return nil, err
	}

	enabled := make([]Config, 0)
	for _, config := range configs {
		if config.Enabled {
			enabled = append(enabled, config)
		}
	}

	return enabled, nil
}

// generateID generates a unique ID for a configuration
func generateID() string {
	return fmt.Sprintf("telegram-%d", time.Now().UnixNano())
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Mode:                  "polling",
		PollingInterval:       defaultPollingInterval,
		MaxConnections:        40,
		ParseMode:             ParseModeMarkdown,
		DisableWebPagePreview: false,
		DisableNotification:   false,
		Enabled:               true,
		AllowedUpdates: []string{
			"message",
			"edited_message",
			"callback_query",
			"my_chat_member",
		},
	}
}

// NewPollingConfig creates a new polling mode configuration
func NewPollingConfig(name, tokenRef string) Config {
	config := DefaultConfig()
	config.Name = name
	config.TokenRef = tokenRef
	config.Mode = "polling"
	return config
}

// NewWebhookConfig creates a new webhook mode configuration
func NewWebhookConfig(name, tokenRef, webhookURL string) Config {
	config := DefaultConfig()
	config.Name = name
	config.TokenRef = tokenRef
	config.Mode = "webhook"
	config.WebhookURL = webhookURL
	return config
}
