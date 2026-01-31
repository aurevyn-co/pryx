package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// URLValidationResult contains the result of validating a custom URL
type URLValidationResult struct {
	Valid         bool     `json:"valid"`
	HTTPS         bool     `json:"https"`
	DomainAllowed bool     `json:"domain_allowed"`
	Errors        []string `json:"errors,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
	NormalizedURL string   `json:"normalized_url,omitempty"`
}

// SearchFilter contains criteria for filtering curated servers
type SearchFilter struct {
	Category      Category      `json:"category,omitempty"`
	Author        string        `json:"author,omitempty"`
	SecurityLevel SecurityLevel `json:"security_level,omitempty"`
	VerifiedOnly  bool          `json:"verified_only,omitempty"`
	Query         string        `json:"query,omitempty"`
}

// DiscoveryService provides curated MCP server discovery and custom URL validation
type DiscoveryService struct {
	mu            sync.RWMutex
	registry      *CuratedRegistry
	customServers []CustomServerEntry
	allowlist     []string
	blocklist     []string
	configPath    string
}

// CustomServerEntry represents a user-added custom MCP server
type CustomServerEntry struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	URL           string              `json:"url"`
	AddedAt       time.Time           `json:"added_at"`
	Validated     bool                `json:"validated"`
	SecurityCheck URLValidationResult `json:"security_check"`
}

// NewDiscoveryService creates a new discovery service with the default curated registry
func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{
		registry:  DefaultCuratedRegistry(),
		allowlist: defaultAllowlist(),
		blocklist: defaultBlocklist(),
	}
}

// NewDiscoveryServiceWithPath creates a discovery service with custom config path
func NewDiscoveryServiceWithPath(configPath string) *DiscoveryService {
	ds := &DiscoveryService{
		registry:   DefaultCuratedRegistry(),
		allowlist:  defaultAllowlist(),
		blocklist:  defaultBlocklist(),
		configPath: configPath,
	}

	if configPath != "" {
		ds.loadCustomServers()
	}

	return ds
}

// GetCuratedRegistry returns the curated registry
func (ds *DiscoveryService) GetCuratedRegistry() *CuratedRegistry {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.registry
}

// SetCuratedRegistry updates the curated registry
func (ds *DiscoveryService) SetCuratedRegistry(registry *CuratedRegistry) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.registry = registry
}

// GetCuratedServer returns a curated server by ID
func (ds *DiscoveryService) GetCuratedServer(id string) (CuratedServer, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.registry.GetByID(id)
}

// SearchCuratedServers searches curated servers by query and filters
func (ds *DiscoveryService) SearchCuratedServers(filter SearchFilter) []CuratedServer {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var results []CuratedServer
	queryLower := strings.ToLower(filter.Query)

	for _, server := range ds.registry.Servers {
		if filter.Category != "" && server.Category != filter.Category {
			continue
		}

		if filter.Author != "" && !strings.EqualFold(server.Author, filter.Author) {
			continue
		}

		if filter.SecurityLevel != "" && server.SecurityLevel != filter.SecurityLevel {
			continue
		}

		if filter.VerifiedOnly && !server.Verified {
			continue
		}

		if filter.Query != "" {
			match := false
			if strings.Contains(strings.ToLower(server.Name), queryLower) {
				match = true
			}
			if strings.Contains(strings.ToLower(server.Description), queryLower) {
				match = true
			}
			for _, tag := range server.Tags {
				if strings.Contains(strings.ToLower(tag), queryLower) {
					match = true
					break
				}
			}
			for _, tool := range server.Tools {
				if strings.Contains(strings.ToLower(tool.Name), queryLower) {
					match = true
					break
				}
				if strings.Contains(strings.ToLower(tool.Description), queryLower) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		results = append(results, server)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].SecurityLevel != results[j].SecurityLevel {
			return results[i].SecurityLevel < results[j].SecurityLevel
		}
		return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
	})

	return results
}

// GetCategories returns all available categories with counts
func (ds *DiscoveryService) GetCategories() map[Category]int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	counts := make(map[Category]int)
	for _, server := range ds.registry.Servers {
		counts[server.Category]++
	}
	return counts
}

// ValidateCustomURL validates a custom MCP server URL
func (ds *DiscoveryService) ValidateCustomURL(urlStr string) URLValidationResult {
	result := URLValidationResult{
		Valid:         true,
		HTTPS:         true,
		DomainAllowed: true,
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid URL format: "+err.Error())
		return result
	}

	if parsed.Scheme != "https" {
		result.HTTPS = false
		if parsed.Scheme == "http" {
			result.Warnings = append(result.Warnings, "HTTP is insecure. HTTPS is strongly recommended.")
		} else {
			result.Valid = false
			result.Errors = append(result.Errors, "URL must use HTTPS protocol")
		}
	}

	if parsed.Hostname() == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "URL must have a valid hostname")
		return result
	}

	host := strings.ToLower(parsed.Hostname())

	for _, blocked := range ds.blocklist {
		if strings.Contains(host, blocked) {
			result.DomainAllowed = false
			result.Valid = false
			result.Errors = append(result.Errors, "Domain is blocked: "+blocked)
			return result
		}
	}

	if len(ds.allowlist) > 0 {
		allowed := false
		for _, allowedPattern := range ds.allowlist {
			if matchesPattern(host, allowedPattern) {
				allowed = true
				break
			}
		}
		if !allowed {
			result.DomainAllowed = false
			result.Warnings = append(result.Warnings, "Domain not in allowlist. Additional verification recommended.")
		}
	}

	result.NormalizedURL = parsed.String()

	return result
}

func (ds *DiscoveryService) AddCustomServer(name, urlStr string) (CustomServerEntry, error) {
	validation := ds.ValidateCustomURL(urlStr)
	if !validation.Valid {
		return CustomServerEntry{}, errors.New("URL validation failed: " + strings.Join(validation.Errors, ", "))
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()

	for _, existing := range ds.customServers {
		if existing.URL == urlStr {
			return CustomServerEntry{}, errors.New("server already exists: " + urlStr)
		}
	}

	entry := CustomServerEntry{
		ID:            generateID(name),
		Name:          name,
		URL:           validation.NormalizedURL,
		AddedAt:       time.Now(),
		Validated:     true,
		SecurityCheck: validation,
	}

	ds.customServers = append(ds.customServers, entry)

	if ds.configPath != "" {
		ds.saveCustomServers()
	}

	return entry, nil
}

// GetCustomServers returns all custom servers
func (ds *DiscoveryService) GetCustomServers() []CustomServerEntry {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	result := make([]CustomServerEntry, len(ds.customServers))
	copy(result, ds.customServers)
	return result
}

func (ds *DiscoveryService) RemoveCustomServer(id string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for i, server := range ds.customServers {
		if server.ID == id {
			ds.customServers = append(ds.customServers[:i], ds.customServers[i+1:]...)
			if ds.configPath != "" {
				ds.saveCustomServers()
			}
			return true
		}
	}
	return false
}

func (ds *DiscoveryService) loadCustomServers() error {
	data, err := os.ReadFile(ds.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var servers []CustomServerEntry
	if err := json.Unmarshal(data, &servers); err != nil {
		return err
	}

	ds.customServers = servers
	return nil
}

func (ds *DiscoveryService) saveCustomServers() error {
	dir := filepath.Dir(ds.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(ds.customServers, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ds.configPath, data, 0644)
}

func (ds *DiscoveryService) GetRecommendedConfig(serverID string) (map[string]string, error) {
	server, ok := ds.GetCuratedServer(serverID)
	if !ok {
		return nil, fmt.Errorf("server not found: %s", serverID)
	}

	return server.RecommendedConfig, nil
}

func (ds *DiscoveryService) GetSecurityWarnings(serverID string) ([]string, error) {
	server, ok := ds.GetCuratedServer(serverID)
	if !ok {
		return nil, fmt.Errorf("server not found: %s", serverID)
	}

	var warnings []string

	if server.SecurityLevel == SecurityLevelC || server.SecurityLevel == SecurityLevelD {
		warnings = append(warnings, fmt.Sprintf("Security Level %s: Use with caution", server.SecurityLevel))
	}

	warnings = append(warnings, server.SecurityWarnings...)

	return warnings, nil
}

func defaultAllowlist() []string {
	return []string{}
}

func defaultBlocklist() []string {
	return []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
		"[::]",
	}
}

func matchesPattern(host, pattern string) bool {
	pattern = strings.ToLower(pattern)
	host = strings.ToLower(host)

	if host == pattern {
		return true
	}

	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[2:]
		return strings.HasSuffix(host, suffix) && host != suffix
	}

	if strings.HasPrefix(pattern, "regex:") {
		regex := pattern[6:]
		re, err := regexp.Compile(regex)
		if err != nil {
			return false
		}
		return re.MatchString(host)
	}

	return false
}

func generateID(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	id := re.ReplaceAllString(name, "-")
	id = strings.ToLower(strings.Trim(id, "-"))

	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s-%d", id, timestamp)
}

func GetSecurityLevelColor(level SecurityLevel) string {
	switch level {
	case SecurityLevelA:
		return "green"
	case SecurityLevelB:
		return "light_green"
	case SecurityLevelC:
		return "yellow"
	case SecurityLevelD:
		return "orange"
	case SecurityLevelF:
		return "red"
	default:
		return "gray"
	}
}

func GetCategoryIcon(category Category) string {
	switch category {
	case CategoryFilesystem:
		return "folder"
	case CategoryWeb:
		return "globe"
	case CategoryDatabase:
		return "database"
	case CategoryAI:
		return "brain"
	case CategoryUtility:
		return "wrench"
	default:
		return "question"
	}
}
