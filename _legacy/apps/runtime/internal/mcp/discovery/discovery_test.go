package discovery

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultCuratedRegistry(t *testing.T) {
	registry := DefaultCuratedRegistry()

	if registry.Version == "" {
		t.Error("expected registry version to be set")
	}

	if registry.Updated == "" {
		t.Error("expected registry updated date to be set")
	}

	if len(registry.Servers) == 0 {
		t.Error("expected registry to have servers")
	}
}

func TestCuratedRegistry_GetByID(t *testing.T) {
	registry := DefaultCuratedRegistry()

	server, ok := registry.GetByID("filesystem")
	if !ok {
		t.Error("expected to find filesystem server")
	}
	if server.ID != "filesystem" {
		t.Errorf("expected server ID to be filesystem, got %s", server.ID)
	}

	_, ok = registry.GetByID("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent server")
	}
}

func TestCuratedRegistry_GetByCategory(t *testing.T) {
	registry := DefaultCuratedRegistry()

	filesystemServers := registry.GetByCategory(CategoryFilesystem)
	if len(filesystemServers) == 0 {
		t.Error("expected to find filesystem servers")
	}

	for _, server := range filesystemServers {
		if server.Category != CategoryFilesystem {
			t.Errorf("expected category to be filesystem, got %s", server.Category)
		}
	}

	dbServers := registry.GetByCategory(CategoryDatabase)
	if len(dbServers) < 2 {
		t.Errorf("expected at least 2 database servers, got %d", len(dbServers))
	}
}

func TestNewDiscoveryService(t *testing.T) {
	ds := NewDiscoveryService()

	if ds == nil {
		t.Fatal("expected discovery service to be created")
	}

	if ds.registry == nil {
		t.Error("expected registry to be initialized")
	}
}

func TestDiscoveryService_GetCuratedServer(t *testing.T) {
	ds := NewDiscoveryService()

	server, ok := ds.GetCuratedServer("github")
	if !ok {
		t.Error("expected to find github server")
	}
	if server.ID != "github" {
		t.Errorf("expected server ID to be github, got %s", server.ID)
	}
}

func TestDiscoveryService_SearchCuratedServers(t *testing.T) {
	ds := NewDiscoveryService()

	tests := []struct {
		name     string
		filter   SearchFilter
		minCount int
		maxCount int
	}{
		{
			name:     "search by query",
			filter:   SearchFilter{Query: "github"},
			minCount: 1,
			maxCount: 10,
		},
		{
			name:     "search by category",
			filter:   SearchFilter{Category: CategoryWeb},
			minCount: 1,
			maxCount: 20,
		},
		{
			name:     "search verified only",
			filter:   SearchFilter{VerifiedOnly: true},
			minCount: 1,
			maxCount: 100,
		},
		{
			name:     "search by security level",
			filter:   SearchFilter{SecurityLevel: SecurityLevelA},
			minCount: 1,
			maxCount: 20,
		},
		{
			name:     "combined filters",
			filter:   SearchFilter{Category: CategoryDatabase, VerifiedOnly: true},
			minCount: 1,
			maxCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := ds.SearchCuratedServers(tt.filter)
			if len(results) < tt.minCount {
				t.Errorf("expected at least %d results, got %d", tt.minCount, len(results))
			}
			if len(results) > tt.maxCount {
				t.Errorf("expected at most %d results, got %d", tt.maxCount, len(results))
			}
		})
	}
}

func TestDiscoveryService_GetCategories(t *testing.T) {
	ds := NewDiscoveryService()

	categories := ds.GetCategories()

	if len(categories) == 0 {
		t.Error("expected categories to be returned")
	}

	if _, ok := categories[CategoryFilesystem]; !ok {
		t.Error("expected filesystem category")
	}
}

func TestDiscoveryService_ValidateCustomURL(t *testing.T) {
	ds := NewDiscoveryService()

	tests := []struct {
		name       string
		url        string
		shouldPass bool
		shouldWarn bool
	}{
		{
			name:       "valid https URL",
			url:        "https://example.com/mcp",
			shouldPass: true,
			shouldWarn: false,
		},
		{
			name:       "http URL warns",
			url:        "http://example.com/mcp",
			shouldPass: true,
			shouldWarn: true,
		},
		{
			name:       "localhost blocked",
			url:        "https://localhost:3000",
			shouldPass: false,
			shouldWarn: false,
		},
		{
			name:       "127.0.0.1 blocked",
			url:        "https://127.0.0.1:3000",
			shouldPass: false,
			shouldWarn: false,
		},
		{
			name:       "invalid URL format",
			url:        "://invalid-url",
			shouldPass: false,
			shouldWarn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ds.ValidateCustomURL(tt.url)

			if tt.shouldPass && !result.Valid {
				t.Errorf("expected URL to be valid, got errors: %v", result.Errors)
			}
			if !tt.shouldPass && result.Valid {
				t.Error("expected URL to be invalid")
			}
			if tt.shouldWarn && len(result.Warnings) == 0 {
				t.Error("expected warnings but got none")
			}
		})
	}
}

func TestDiscoveryService_CustomServers(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-servers.json")

	ds := NewDiscoveryServiceWithPath(configPath)

	entry, err := ds.AddCustomServer("My Server", "https://example.com/mcp")
	if err != nil {
		t.Fatalf("failed to add custom server: %v", err)
	}

	if entry.ID == "" {
		t.Error("expected entry ID to be set")
	}
	if entry.Name != "My Server" {
		t.Errorf("expected name to be 'My Server', got %s", entry.Name)
	}

	servers := ds.GetCustomServers()
	if len(servers) != 1 {
		t.Errorf("expected 1 custom server, got %d", len(servers))
	}

	_, err = ds.AddCustomServer("Duplicate", "https://example.com/mcp")
	if err == nil {
		t.Error("expected error when adding duplicate URL")
	}

	removed := ds.RemoveCustomServer(entry.ID)
	if !removed {
		t.Error("expected server to be removed")
	}

	servers = ds.GetCustomServers()
	if len(servers) != 0 {
		t.Errorf("expected 0 custom servers after removal, got %d", len(servers))
	}
}

func TestDiscoveryService_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-servers.json")

	ds1 := NewDiscoveryServiceWithPath(configPath)
	_, err := ds1.AddCustomServer("Test Server", "https://test.example.com")
	if err != nil {
		t.Fatalf("failed to add custom server: %v", err)
	}

	ds2 := NewDiscoveryServiceWithPath(configPath)
	servers := ds2.GetCustomServers()
	if len(servers) != 1 {
		t.Errorf("expected 1 server after loading from disk, got %d", len(servers))
	}
	if servers[0].Name != "Test Server" {
		t.Errorf("expected server name to be 'Test Server', got %s", servers[0].Name)
	}
}

func TestDiscoveryService_GetRecommendedConfig(t *testing.T) {
	ds := NewDiscoveryService()

	config, err := ds.GetRecommendedConfig("github")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(config) == 0 {
		t.Error("expected config for github server")
	}
	if _, ok := config["GITHUB_TOKEN"]; !ok {
		t.Error("expected GITHUB_TOKEN in config")
	}

	_, err = ds.GetRecommendedConfig("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent server")
	}
}

func TestDiscoveryService_GetSecurityWarnings(t *testing.T) {
	ds := NewDiscoveryService()

	warnings, err := ds.GetSecurityWarnings("shell")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(warnings) == 0 {
		t.Error("expected warnings for shell server")
	}

	hasExecWarning := false
	for _, w := range warnings {
		if strings.Contains(strings.ToLower(w), "shell") || strings.Contains(strings.ToLower(w), "command") {
			hasExecWarning = true
			break
		}
	}
	if !hasExecWarning {
		t.Error("expected warning about shell command execution")
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		host     string
		pattern  string
		expected bool
	}{
		{"example.com", "example.com", true},
		{"sub.example.com", "*.example.com", true},
		{"example.com", "*.example.com", false},
		{"example.org", "example.com", false},
	}

	for _, tt := range tests {
		result := matchesPattern(tt.host, tt.pattern)
		if result != tt.expected {
			t.Errorf("matchesPattern(%q, %q) = %v, expected %v", tt.host, tt.pattern, result, tt.expected)
		}
	}
}

func TestGetSecurityLevelColor(t *testing.T) {
	tests := []struct {
		level    SecurityLevel
		expected string
	}{
		{SecurityLevelA, "green"},
		{SecurityLevelB, "light_green"},
		{SecurityLevelC, "yellow"},
		{SecurityLevelD, "orange"},
		{SecurityLevelF, "red"},
		{"unknown", "gray"},
	}

	for _, tt := range tests {
		result := GetSecurityLevelColor(tt.level)
		if result != tt.expected {
			t.Errorf("GetSecurityLevelColor(%s) = %s, expected %s", tt.level, result, tt.expected)
		}
	}
}

func TestGetCategoryIcon(t *testing.T) {
	tests := []struct {
		category Category
		expected string
	}{
		{CategoryFilesystem, "folder"},
		{CategoryWeb, "globe"},
		{CategoryDatabase, "database"},
		{CategoryAI, "brain"},
		{CategoryUtility, "wrench"},
		{"unknown", "question"},
	}

	for _, tt := range tests {
		result := GetCategoryIcon(tt.category)
		if result != tt.expected {
			t.Errorf("GetCategoryIcon(%s) = %s, expected %s", tt.category, result, tt.expected)
		}
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID("Test Server")
	id2 := generateID("Test Server")

	if id1 == id2 {
		t.Error("expected generated IDs to be unique")
	}

	if !strings.Contains(strings.ToLower(id1), "test") {
		t.Error("expected ID to contain sanitized name")
	}

	id3 := generateID("Server With Spaces!")
	if strings.Contains(id3, " ") || strings.Contains(id3, "!") {
		t.Error("expected ID to not contain spaces or special chars")
	}
}

func TestCuratedServerToServerConfig(t *testing.T) {
	registry := DefaultCuratedRegistry()

	server, ok := registry.GetByID("postgresql")
	if !ok {
		t.Skip("postgresql server not found in registry")
	}

	if server.Transport != "stdio" {
		t.Errorf("expected transport to be stdio, got %s", server.Transport)
	}

	if len(server.Command) == 0 {
		t.Error("expected command to be set for stdio transport")
	}

	if len(server.EnvironmentRequired) == 0 {
		t.Error("expected environment variables to be required")
	}
}
