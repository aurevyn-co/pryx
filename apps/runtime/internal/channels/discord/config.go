package discord

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultConfigDir  = ".pryx/config"
	defaultConfigFile = "discord.json"
)

// Config represents the Discord bot configuration
type Config struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	TokenRef           string               `json:"token_ref"`
	Token              string               `json:"token,omitempty"`
	ApplicationID      string               `json:"application_id,omitempty"`
	Intents            Intent               `json:"intents"`
	AllowedGuilds      []string             `json:"allowed_guilds"`
	AllowedChannels    []string             `json:"allowed_channels"`
	Commands           []ApplicationCommand `json:"commands"`
	DefaultPermissions *string              `json:"default_permissions,omitempty"`
	DMPermission       bool                 `json:"dm_permission"`
	Presence           *UpdateStatus        `json:"presence,omitempty"`
	ShardID            int                  `json:"shard_id,omitempty"`
	NumShards          int                  `json:"num_shards,omitempty"`
	LargeThreshold     int                  `json:"large_threshold,omitempty"`
	Enabled            bool                 `json:"enabled"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
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

	if c.Intents == 0 {
		c.Intents = DefaultIntents()
	}

	if c.LargeThreshold == 0 {
		c.LargeThreshold = 50
	}

	return nil
}

// IsGuildAllowed checks if a guild ID is in the whitelist
func (c *Config) IsGuildAllowed(guildID string) bool {
	if len(c.AllowedGuilds) == 0 {
		return true
	}

	for _, id := range c.AllowedGuilds {
		if id == guildID {
			return true
		}
	}

	return false
}

// IsChannelAllowed checks if a channel ID is in the whitelist
func (c *Config) IsChannelAllowed(channelID string) bool {
	if len(c.AllowedChannels) == 0 {
		return true
	}

	for _, id := range c.AllowedChannels {
		if id == channelID {
			return true
		}
	}

	return false
}

// SetDefaults sets default values for optional fields
func (c *Config) SetDefaults() {
	if c.Intents == 0 {
		c.Intents = DefaultIntents()
	}

	if c.LargeThreshold == 0 {
		c.LargeThreshold = 50
	}

	if c.NumShards == 0 {
		c.NumShards = 1
	}
}

// ConfigManager manages Discord bot configurations
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

// LoadAll loads all Discord configurations
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

	for i := range configs {
		configs[i].SetDefaults()
	}

	return configs, nil
}

// SaveAll saves all Discord configurations
func (cm *ConfigManager) SaveAll(configs []Config) error {
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configsToSave := make([]Config, len(configs))
	for i, config := range configs {
		configsToSave[i] = config
		configsToSave[i].Token = ""
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

	return nil, fmt.Errorf("discord config not found: %s", id)
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
		return fmt.Errorf("discord config not found: %s", id)
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

	if name, ok := updates["name"].(string); ok {
		config.Name = name
	}
	if tokenRef, ok := updates["token_ref"].(string); ok {
		config.TokenRef = tokenRef
	}
	if token, ok := updates["token"].(string); ok {
		config.Token = token
	}
	if applicationID, ok := updates["application_id"].(string); ok {
		config.ApplicationID = applicationID
	}
	if intents, ok := updates["intents"].(Intent); ok {
		config.Intents = intents
	}
	if allowedGuilds, ok := updates["allowed_guilds"].([]string); ok {
		config.AllowedGuilds = allowedGuilds
	}
	if allowedChannels, ok := updates["allowed_channels"].([]string); ok {
		config.AllowedChannels = allowedChannels
	}
	if commands, ok := updates["commands"].([]ApplicationCommand); ok {
		config.Commands = commands
	}
	if presence, ok := updates["presence"].(*UpdateStatus); ok {
		config.Presence = presence
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

// AddAllowedGuild adds a guild ID to the whitelist
func (cm *ConfigManager) AddAllowedGuild(configID string, guildID string) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	for _, id := range config.AllowedGuilds {
		if id == guildID {
			return nil
		}
	}

	config.AllowedGuilds = append(config.AllowedGuilds, guildID)
	return cm.Save(*config)
}

// RemoveAllowedGuild removes a guild ID from the whitelist
func (cm *ConfigManager) RemoveAllowedGuild(configID string, guildID string) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	filtered := make([]string, 0, len(config.AllowedGuilds))
	for _, id := range config.AllowedGuilds {
		if id != guildID {
			filtered = append(filtered, id)
		}
	}

	config.AllowedGuilds = filtered
	return cm.Save(*config)
}

// AddAllowedChannel adds a channel ID to the whitelist
func (cm *ConfigManager) AddAllowedChannel(configID string, channelID string) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	for _, id := range config.AllowedChannels {
		if id == channelID {
			return nil
		}
	}

	config.AllowedChannels = append(config.AllowedChannels, channelID)
	return cm.Save(*config)
}

// RemoveAllowedChannel removes a channel ID from the whitelist
func (cm *ConfigManager) RemoveAllowedChannel(configID string, channelID string) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	filtered := make([]string, 0, len(config.AllowedChannels))
	for _, id := range config.AllowedChannels {
		if id != channelID {
			filtered = append(filtered, id)
		}
	}

	config.AllowedChannels = filtered
	return cm.Save(*config)
}

// SetToken sets the bot token (not persisted, runtime only)
func (cm *ConfigManager) SetToken(configID string, token string) error {
	config, err := cm.Get(configID)
	if err != nil {
		return err
	}

	config.Token = token
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
	return fmt.Sprintf("discord-%d", time.Now().UnixNano())
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Intents:         DefaultIntents(),
		LargeThreshold:  50,
		NumShards:       1,
		DMPermission:    true,
		Enabled:         true,
		AllowedGuilds:   []string{},
		AllowedChannels: []string{},
		Commands:        []ApplicationCommand{},
	}
}

// NewBotConfig creates a new bot configuration with basic settings
func NewBotConfig(name, tokenRef string) Config {
	config := DefaultConfig()
	config.Name = name
	config.TokenRef = tokenRef
	return config
}

// NewGuildOnlyConfig creates a configuration restricted to specific guilds
func NewGuildOnlyConfig(name, tokenRef string, guildIDs []string) Config {
	config := DefaultConfig()
	config.Name = name
	config.TokenRef = tokenRef
	config.AllowedGuilds = guildIDs
	config.DMPermission = false
	return config
}
