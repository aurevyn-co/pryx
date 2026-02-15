package webhook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const configDir = ".pryx/config"
const configFile = "webhooks.json"

type ConfigManager struct {
	configPath string
}

func NewConfigManager() *ConfigManager {
	home, _ := os.UserHomeDir()
	return &ConfigManager{
		configPath: filepath.Join(home, configDir, configFile),
	}
}

func (cm *ConfigManager) LoadAll() ([]WebhookConfig, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []WebhookConfig{}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var configs []WebhookConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return configs, nil
}

func (cm *ConfigManager) SaveAll(configs []WebhookConfig) error {
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (cm *ConfigManager) Get(id string) (*WebhookConfig, error) {
	configs, err := cm.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if config.ID == id {
			return &config, nil
		}
	}

	return nil, fmt.Errorf("webhook config not found: %s", id)
}

func (cm *ConfigManager) Save(config WebhookConfig) error {
	configs, err := cm.LoadAll()
	if err != nil {
		return err
	}

	now := time.Now()
	config.UpdatedAt = now

	found := false
	for i, c := range configs {
		if c.ID == config.ID {
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

func (cm *ConfigManager) Delete(id string) error {
	configs, err := cm.LoadAll()
	if err != nil {
		return err
	}

	filtered := make([]WebhookConfig, 0, len(configs))
	for _, config := range configs {
		if config.ID != id {
			filtered = append(filtered, config)
		}
	}

	if len(filtered) == len(configs) {
		return fmt.Errorf("webhook config not found: %s", id)
	}

	return cm.SaveAll(filtered)
}

func (cm *ConfigManager) Validate(config WebhookConfig) error {
	if config.ID == "" {
		return fmt.Errorf("ID is required")
	}

	if config.Name == "" {
		return fmt.Errorf("name is required")
	}

	if config.TargetURL == "" && config.Port == 0 {
		return fmt.Errorf("either TargetURL or Port must be specified")
	}

	if config.RetryConfig.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	return nil
}

func (cm *ConfigManager) Create(config WebhookConfig) (*WebhookConfig, error) {
	if config.ID == "" {
		config.ID = generateID()
	}

	if err := cm.Validate(config); err != nil {
		return nil, err
	}

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	if config.RetryConfig.MaxRetries == 0 {
		config.RetryConfig = DefaultRetry()
	}

	if err := cm.Save(config); err != nil {
		return nil, err
	}

	return &config, nil
}
