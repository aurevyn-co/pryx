package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleSessionsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessions := []map[string]interface{}{
		{
			"id":         "test-session-1",
			"name":       "Test Session",
			"created_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"updated_at": time.Now().Format(time.RFC3339),
		},
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions": sessions,
	})
}

func (s *Server) handleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	sessionID := "session-" + time.Now().Format("20060102150405")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         sessionID,
		"name":       req.Name,
		"created_at": time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleSessionGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "session id is required"})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         sessionID,
		"name":       "Test Session",
		"status":     "active",
		"created_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleSessionDelete(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "session id is required"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
