package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler manages HTTP endpoints for webhooks
type Handler struct {
	manager *Manager
}

// NewHandler creates a new webhook handler
func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// RegisterRoutes registers webhook HTTP routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/webhooks/{channelId}", h.handleIncomingWebhook)
	r.Get("/webhooks", h.handleListWebhooks)
	r.Post("/webhooks", h.handleCreateWebhook)
	r.Get("/webhooks/{channelId}/logs", h.handleGetLogs)
	r.Get("/webhooks/{channelId}/health", h.handleHealth)
}

// handleIncomingWebhook processes incoming webhook requests
func (h *Handler) handleIncomingWebhook(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	if channelID == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channel := h.manager.GetChannel(channelID)
	if channel == nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	receiver := NewReceiver(channel.Config())
	msg, err := receiver.Handle(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := channel.Receive(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"id":     msg.ID,
	})
}

// handleListWebhooks lists all configured webhooks
func (h *Handler) handleListWebhooks(w http.ResponseWriter, r *http.Request) {
	channels := h.manager.ListChannels()

	configs := make([]WebhookConfig, 0, len(channels))
	for _, ch := range channels {
		configs = append(configs, ch.Config())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configs)
}

// handleCreateWebhook creates a new webhook channel
func (h *Handler) handleCreateWebhook(w http.ResponseWriter, r *http.Request) {
	var config WebhookConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	channel, err := h.manager.CreateChannel(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(channel.Config())
}

// handleGetLogs returns delivery logs for a webhook
func (h *Handler) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	if channelID == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channel := h.manager.GetChannel(channelID)
	if channel == nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	logs := channel.GetDeliveryLogs(100)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// handleHealth checks the health of a webhook channel
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelId")
	if channelID == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channel := h.manager.GetChannel(channelID)
	if channel == nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	if err := channel.Health(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
