package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"pryx-core/internal/bus"
)

// WebhookReceiver handles incoming webhook requests from Telegram
type WebhookReceiver struct {
	config   *Config
	handler  *Handler
	eventBus EventPublisher
}

// EventPublisher is the interface for publishing events
type EventPublisher interface {
	Publish(event bus.Event)
}

// NewWebhookReceiver creates a new webhook receiver
func NewWebhookReceiver(config *Config, handler *Handler, eventBus EventPublisher) *WebhookReceiver {
	return &WebhookReceiver{
		config:   config,
		handler:  handler,
		eventBus: eventBus,
	}
}

// ServeHTTP implements the http.Handler interface
func (wr *WebhookReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify secret token if configured
	if wr.config.WebhookSecret != "" {
		secret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if secret != wr.config.WebhookSecret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Read and parse the update
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process the update asynchronously
	go wr.handler.HandleUpdate(r.Context(), &update)

	// Respond immediately to Telegram
	w.WriteHeader(http.StatusOK)
}

// VerifySecret verifies the webhook secret from a request
func (wr *WebhookReceiver) VerifySecret(r *http.Request) bool {
	if wr.config.WebhookSecret == "" {
		return true
	}
	return r.Header.Get("X-Telegram-Bot-Api-Secret-Token") == wr.config.WebhookSecret
}

// ParseUpdate parses an update from the request body
func (wr *WebhookReceiver) ParseUpdate(r *http.Request) (*Update, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		return nil, err
	}

	return &update, nil
}

// HandleUpdate processes an update
func (wr *WebhookReceiver) HandleUpdate(update *Update) error {
	return wr.handler.HandleUpdate(nil, update)
}

// WebhookResponse is the response structure for webhook setup
type WebhookResponse struct {
	OK          bool   `json:"ok"`
	Result      bool   `json:"result,omitempty"`
	Description string `json:"description,omitempty"`
}

// WriteResponse writes a JSON response
func WriteResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ReadBody reads the request body and returns it as bytes
func ReadBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

// ResetBody resets the request body for re-reading
func ResetBody(r *http.Request, body []byte) {
	r.Body = io.NopCloser(bytes.NewBuffer(body))
}
