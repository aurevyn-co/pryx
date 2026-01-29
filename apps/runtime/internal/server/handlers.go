package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"pryx-core/internal/skills"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleMCPTools(w http.ResponseWriter, r *http.Request) {
	refresh := strings.TrimSpace(r.URL.Query().Get("refresh")) == "1"
	tools, err := s.mcp.ListToolsFlat(r.Context(), refresh)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": tools,
	})
}

type mcpCallRequest struct {
	SessionID string                 `json:"session_id"`
	Tool      string                 `json:"tool"`
	Arguments map[string]interface{} `json:"arguments"`
}

func (s *Server) handleMCPCall(w http.ResponseWriter, r *http.Request) {
	req := mcpCallRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid json body",
		})
		return
	}
	if strings.TrimSpace(req.Tool) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "missing tool",
		})
		return
	}
	if req.Arguments == nil {
		req.Arguments = map[string]interface{}{}
	}

	res, err := s.mcp.CallTool(r.Context(), strings.TrimSpace(req.SessionID), req.Tool, req.Arguments)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleSkillsList(w http.ResponseWriter, r *http.Request) {
	reg := s.skills
	if reg == nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"skills": []skills.Skill{},
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"skills": reg.List(),
	})
}

func (s *Server) handleSkillsInfo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "missing id",
		})
		return
	}
	reg := s.skills
	if reg == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "not found",
		})
		return
	}
	skill, ok := reg.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "not found",
		})
		return
	}
	_ = json.NewEncoder(w).Encode(skill)
}

func (s *Server) handleSkillsBody(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "missing id",
		})
		return
	}
	reg := s.skills
	if reg == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "not found",
		})
		return
	}
	skill, ok := reg.Get(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "not found",
		})
		return
	}
	body, err := skill.Body()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"body": body,
	})
}

func (s *Server) handleProvidersList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if s.catalog != nil {
		var providers []map[string]interface{}
		for id, info := range s.catalog.Providers {
			requiresKey := len(info.Env) > 0
			providers = append(providers, map[string]interface{}{
				"id":               id,
				"name":             info.Name,
				"requires_api_key": requiresKey,
			})
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"providers": providers})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"providers": []map[string]interface{}{
			{"id": "openai", "name": "OpenAI", "requires_api_key": true},
			{"id": "anthropic", "name": "Anthropic", "requires_api_key": true},
			{"id": "google", "name": "Google AI", "requires_api_key": true},
			{"id": "ollama", "name": "Ollama (Local)", "requires_api_key": false},
		},
	})
}

func (s *Server) handleProviderModels(w http.ResponseWriter, r *http.Request) {
	providerID := strings.TrimSpace(chi.URLParam(r, "id"))
	if providerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing provider id"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if s.catalog != nil {
		models := s.catalog.GetProviderModels(providerID)
		var result []map[string]interface{}
		for _, m := range models {
			result = append(result, map[string]interface{}{
				"id":   m.ID,
				"name": m.Name,
			})
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"models": result})
		return
	}

	staticModels := map[string][]map[string]interface{}{
		"openai": {
			{"id": "gpt-4", "name": "GPT-4"},
			{"id": "gpt-4-turbo", "name": "GPT-4 Turbo"},
			{"id": "gpt-3.5-turbo", "name": "GPT-3.5 Turbo"},
		},
		"anthropic": {
			{"id": "claude-3-opus", "name": "Claude 3 Opus"},
			{"id": "claude-3-sonnet", "name": "Claude 3 Sonnet"},
			{"id": "claude-3-haiku", "name": "Claude 3 Haiku"},
		},
		"google": {
			{"id": "gemini-pro", "name": "Gemini Pro"},
			{"id": "gemini-ultra", "name": "Gemini Ultra"},
		},
		"ollama": {
			{"id": "llama3", "name": "Llama 3"},
			{"id": "llama2", "name": "Llama 2"},
			{"id": "mistral", "name": "Mistral"},
		},
	}

	if providerModels, ok := staticModels[providerID]; ok {
		json.NewEncoder(w).Encode(map[string]interface{}{"models": providerModels})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "provider not found"})
	}
}

func (s *Server) handleModelsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if s.catalog != nil {
		var result []map[string]interface{}
		for _, m := range s.catalog.Models {
			result = append(result, map[string]interface{}{
				"id":       m.ID,
				"name":     m.Name,
				"provider": m.Provider,
			})
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"models": result})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"models": []map[string]interface{}{
			{"id": "gpt-4", "name": "GPT-4", "provider": "openai"},
			{"id": "gpt-4-turbo", "name": "GPT-4 Turbo", "provider": "openai"},
			{"id": "gpt-3.5-turbo", "name": "GPT-3.5 Turbo", "provider": "openai"},
			{"id": "claude-3-opus", "name": "Claude 3 Opus", "provider": "anthropic"},
			{"id": "claude-3-sonnet", "name": "Claude 3 Sonnet", "provider": "anthropic"},
			{"id": "claude-3-haiku", "name": "Claude 3 Haiku", "provider": "anthropic"},
		},
	})
}
