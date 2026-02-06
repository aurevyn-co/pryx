package skills

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type EnabledConfig struct {
	EnabledSkills map[string]bool `yaml:"enabled_skills" json:"enabled_skills"`
}

func EnabledConfigPath() string {
	if p := strings.TrimSpace(os.Getenv("PRYX_SKILLS_CONFIG_PATH")); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pryx", "skills.yaml")
}

func LoadEnabledConfig(path string) (*EnabledConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &EnabledConfig{EnabledSkills: map[string]bool{}}, nil
		}
		return nil, err
	}

	cfg := EnabledConfig{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.EnabledSkills == nil {
		cfg.EnabledSkills = map[string]bool{}
	}
	return &cfg, nil
}

func SaveEnabledConfig(path string, cfg *EnabledConfig) error {
	if cfg.EnabledSkills == nil {
		cfg.EnabledSkills = map[string]bool{}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
