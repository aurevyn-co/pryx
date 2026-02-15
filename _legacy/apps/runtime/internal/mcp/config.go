package mcp

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type ServersConfig struct {
	Servers map[string]ServerConfig `json:"servers"`
}

type ServerConfig struct {
	Transport       string            `json:"transport"`
	URL             string            `json:"url,omitempty"`
	Command         []string          `json:"command,omitempty"`
	Cwd             string            `json:"cwd,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	ProtocolVersion string            `json:"protocol_version,omitempty"`
	Auth            *AuthConfig       `json:"auth,omitempty"`
}

type AuthConfig struct {
	Type     string `json:"type"`
	TokenRef string `json:"token_ref,omitempty"`
}

func DefaultServersConfigPaths() []string {
	var paths []string

	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, ".pryx", "mcp", "servers.json"))
	}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".pryx", "mcp", "servers.json"))
	}

	return paths
}

func LoadServersConfigFromFirstExisting(paths []string) (*ServersConfig, string, error) {
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, "", err
		}

		cfg := &ServersConfig{}
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, "", err
		}
		if cfg.Servers == nil {
			cfg.Servers = map[string]ServerConfig{}
		}
		normalizeServerTransports(cfg)
		return cfg, p, nil
	}

	return &ServersConfig{Servers: map[string]ServerConfig{}}, "", nil
}

func normalizeServerTransports(cfg *ServersConfig) {
	for name, server := range cfg.Servers {
		transport := strings.TrimSpace(server.Transport)
		if transport == "" {
			if server.URL != "" {
				server.Transport = "http"
			} else if len(server.Command) > 0 {
				server.Transport = "stdio"
			}
		} else if transport == "stdio" && server.URL != "" && len(server.Command) == 0 {
			server.Transport = "http"
		}
		cfg.Servers[name] = server
	}
}
