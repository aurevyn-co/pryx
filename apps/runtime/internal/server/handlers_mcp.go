package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"pryx-core/internal/mcp/discovery"
	"pryx-core/internal/validation"
)

// handleMCPDiscoveryCurated returns the list of curated MCP servers
// Supports filtering by category, author, security level, and search query
func (s *Server) handleMCPDiscoveryCurated(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	filter := discovery.SearchFilter{
		Query:        r.URL.Query().Get("q"),
		Author:       r.URL.Query().Get("author"),
		VerifiedOnly: r.URL.Query().Get("verified") == "true",
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filter.Category = discovery.Category(category)
	}

	if securityLevel := r.URL.Query().Get("security_level"); securityLevel != "" {
		filter.SecurityLevel = discovery.SecurityLevel(securityLevel)
	}

	servers := s.mcpDiscovery.SearchCuratedServers(filter)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"servers": servers,
		"count":   len(servers),
	})
}

// handleMCPDiscoveryCategories returns all available categories with counts
func (s *Server) handleMCPDiscoveryCategories(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	categories := s.mcpDiscovery.GetCategories()

	var result []map[string]interface{}
	for cat, count := range categories {
		result = append(result, map[string]interface{}{
			"id":    cat,
			"name":  string(cat),
			"icon":  discovery.GetCategoryIcon(cat),
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": result,
	})
}

// handleMCPDiscoveryServer returns details of a specific curated server
func (s *Server) handleMCPDiscoveryServer(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	id := strings.TrimSpace(chi.URLParam(r, "id"))

	validator := validation.NewValidator()
	if err := validator.ValidateID("id", id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	server, ok := s.mcpDiscovery.GetCuratedServer(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "server not found"})
		return
	}

	warnings, _ := s.mcpDiscovery.GetSecurityWarnings(id)
	config, _ := s.mcpDiscovery.GetRecommendedConfig(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"server":               server,
		"security_warnings":    warnings,
		"recommended_config":   config,
		"security_level_color": discovery.GetSecurityLevelColor(server.SecurityLevel),
	})
}

// handleMCPDiscoveryValidateURL validates a custom MCP server URL
func (s *Server) handleMCPDiscoveryValidateURL(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "url is required"})
		return
	}

	result := s.mcpDiscovery.ValidateCustomURL(req.URL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleMCPDiscoveryAddCustom adds a custom MCP server
func (s *Server) handleMCPDiscoveryAddCustom(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	validator := validation.NewValidator()
	if err := validator.ValidateRequired("name", req.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := validator.ValidateRequired("url", req.URL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	entry, err := s.mcpDiscovery.AddCustomServer(req.Name, req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// handleMCPDiscoveryCustomServers returns all custom servers
func (s *Server) handleMCPDiscoveryCustomServers(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	servers := s.mcpDiscovery.GetCustomServers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"servers": servers,
		"count":   len(servers),
	})
}

// handleMCPDiscoveryRemoveCustom removes a custom server
func (s *Server) handleMCPDiscoveryRemoveCustom(w http.ResponseWriter, r *http.Request) {
	if s.mcpDiscovery == nil {
		s.mcpDiscovery = discovery.NewDiscoveryService()
	}

	id := strings.TrimSpace(chi.URLParam(r, "id"))

	validator := validation.NewValidator()
	if err := validator.ValidateID("id", id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	removed := s.mcpDiscovery.RemoveCustomServer(id)
	if !removed {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "server not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
