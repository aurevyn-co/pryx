package telegram

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"pryx-core/internal/bus"
	"pryx-core/internal/channels"
)

// HealthStatus represents the health status of the channel
type HealthStatus struct {
	Healthy   bool   `json:"healthy"`
	Message   string `json:"message,omitempty"`
	LastError string `json:"last_error,omitempty"`
}

// Channel implements the channels.Channel interface for Telegram
type Channel struct {
	id       string
	config   *Config
	client   *Client
	handler  *Handler
	eventBus *bus.Bus

	// Internal state
	status   channels.Status
	statusMu sync.RWMutex
	health   HealthStatus
	healthMu sync.RWMutex

	// Mode-specific components
	poller  *Poller
	webhook *WebhookReceiver

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewChannel creates a new Telegram channel
func NewChannel(config Config, eventBus *bus.Bus) (*Channel, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	config.SetDefaults()

	// Create API client
	client := NewClient(config.Token)

	// Create handler
	handler := NewHandler(&config, client, eventBus)

	ch := &Channel{
		id:       config.ID,
		config:   &config,
		client:   client,
		handler:  handler,
		eventBus: eventBus,
		status:   channels.StatusDisconnected,
		health: HealthStatus{
			Healthy: false,
			Message: "Not connected",
		},
	}

	return ch, nil
}

// ID returns the channel ID
func (ch *Channel) ID() string {
	return ch.id
}

// Type returns the channel type
func (ch *Channel) Type() string {
	return "telegram"
}

// Connect establishes connection to Telegram
func (ch *Channel) Connect(ctx context.Context) error {
	ch.statusMu.Lock()
	defer ch.statusMu.Unlock()

	if ch.status == channels.StatusConnected || ch.status == channels.StatusConnecting {
		return fmt.Errorf("channel already connected or connecting")
	}

	ch.status = channels.StatusConnecting

	// Validate token by calling getMe
	botInfo, err := ch.client.GetMe(ctx)
	if err != nil {
		ch.status = channels.StatusError
		ch.updateHealth(false, "Failed to validate token", err.Error())
		return fmt.Errorf("failed to validate token: %w", err)
	}

	// Register bot commands if configured
	if len(ch.config.Commands) > 0 {
		if _, err := ch.client.SetMyCommands(ctx, ch.config.Commands); err != nil {
			// Log but don't fail - commands are optional
			ch.publishError(fmt.Sprintf("Failed to set commands: %v", err))
		}
	}

	// Create context for the channel
	ch.ctx, ch.cancel = context.WithCancel(context.Background())

	// Setup mode-specific receiver
	switch ch.config.Mode {
	case "webhook":
		if err := ch.setupWebhook(ctx); err != nil {
			ch.status = channels.StatusError
			ch.updateHealth(false, "Failed to setup webhook", err.Error())
			return err
		}
	case "polling":
		ch.setupPolling()
	default:
		ch.status = channels.StatusError
		return fmt.Errorf("unknown mode: %s", ch.config.Mode)
	}

	ch.status = channels.StatusConnected
	ch.updateHealth(true, fmt.Sprintf("Connected as @%s", botInfo.Username), "")

	// Publish status event
	ch.publishStatus("connected", botInfo)

	return nil
}

// Disconnect closes the connection to Telegram
func (ch *Channel) Disconnect(ctx context.Context) error {
	ch.statusMu.Lock()
	defer ch.statusMu.Unlock()

	if ch.status != channels.StatusConnected {
		return nil
	}

	ch.status = channels.StatusDisconnected
	ch.updateHealth(false, "Disconnected", "")

	// Cancel context to stop goroutines
	if ch.cancel != nil {
		ch.cancel()
	}

	// Wait for goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
		ch.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines finished
	case <-ctx.Done():
		// Timeout
	}

	// Delete webhook if in webhook mode
	if ch.config.Mode == "webhook" {
		if _, err := ch.client.DeleteWebhook(ctx, false); err != nil {
			ch.publishError(fmt.Sprintf("Failed to delete webhook: %v", err))
		}
	}

	ch.poller = nil
	ch.webhook = nil

	ch.publishStatus("disconnected", nil)

	return nil
}

// Send sends a message to a Telegram chat
func (ch *Channel) Send(ctx context.Context, msg channels.Message) error {
	if ch.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}

	chatID, err := strconv.ParseInt(msg.ChannelID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	_, err = ch.client.SendMessage(ctx, chatID, msg.Content,
		WithParseMode(ch.config.ParseMode),
		WithDisableWebPagePreview(ch.config.DisableWebPagePreview),
		WithDisableNotification(ch.config.DisableNotification))

	return err
}

// Status returns the current connection status
func (ch *Channel) Status() channels.Status {
	ch.statusMu.RLock()
	defer ch.statusMu.RUnlock()
	return ch.status
}

// Health returns the health status of the channel
func (ch *Channel) Health() HealthStatus {
	ch.healthMu.RLock()
	defer ch.healthMu.RUnlock()
	return ch.health
}

// updateHealth updates the health status
func (ch *Channel) updateHealth(healthy bool, message, lastError string) {
	ch.healthMu.Lock()
	defer ch.healthMu.Unlock()
	ch.health = HealthStatus{
		Healthy:   healthy,
		Message:   message,
		LastError: lastError,
	}
}

// setupWebhook configures webhook mode
func (ch *Channel) setupWebhook(ctx context.Context) error {
	// Set webhook on Telegram
	_, err := ch.client.SetWebhook(ctx, ch.config.WebhookURL,
		WithMaxConnections(ch.config.MaxConnections),
		WithAllowedUpdates(ch.config.AllowedUpdates),
		WithDropPendingUpdates(ch.config.DropPendingUpdates))

	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	// Create webhook receiver
	ch.webhook = NewWebhookReceiver(ch.config, ch.handler, ch.eventBus)

	return nil
}

// setupPolling configures polling mode
func (ch *Channel) setupPolling() {
	ch.poller = NewPoller(ch.config, ch.client, ch.handler, ch.eventBus)

	ch.wg.Add(1)
	go func() {
		defer ch.wg.Done()
		ch.poller.Start(ch.ctx)
	}()
}

// GetWebhookHandler returns the HTTP handler for webhook mode
func (ch *Channel) GetWebhookHandler() http.Handler {
	if ch.webhook != nil {
		return ch.webhook
	}
	return nil
}

// SendPhoto sends a photo to a chat
func (ch *Channel) SendPhoto(ctx context.Context, chatID int64, photo InputFile, caption string) error {
	if ch.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}

	_, err := ch.client.SendPhoto(ctx, chatID, photo, caption,
		WithPhotoParseMode(ch.config.ParseMode),
		WithPhotoDisableNotification(ch.config.DisableNotification))

	return err
}

// SendDocument sends a document to a chat
func (ch *Channel) SendDocument(ctx context.Context, chatID int64, document InputFile, caption string) error {
	if ch.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}

	_, err := ch.client.SendDocument(ctx, chatID, document, caption,
		WithDocumentParseMode(ch.config.ParseMode),
		WithDocumentDisableNotification(ch.config.DisableNotification))

	return err
}

// SetCommands updates the bot's command list
func (ch *Channel) SetCommands(ctx context.Context, commands []BotCommand) error {
	if ch.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}

	_, err := ch.client.SetMyCommands(ctx, commands)
	return err
}

// GetBotInfo returns information about the bot
func (ch *Channel) GetBotInfo(ctx context.Context) (*User, error) {
	return ch.client.GetMe(ctx)
}

// IsChatAllowed checks if a chat ID is whitelisted
func (ch *Channel) IsChatAllowed(chatID int64) bool {
	return ch.config.IsChatAllowed(chatID)
}

// publishStatus publishes a status change event
func (ch *Channel) publishStatus(status string, botInfo *User) {
	if ch.eventBus == nil {
		return
	}

	payload := map[string]interface{}{
		"channel_id":   ch.id,
		"channel_type": ch.Type(),
		"status":       status,
		"mode":         ch.config.Mode,
	}

	if botInfo != nil {
		payload["bot_username"] = botInfo.Username
		payload["bot_name"] = botInfo.FirstName
	}

	ch.eventBus.Publish(bus.NewEvent(bus.EventChannelStatus, "", payload))
}

// publishError publishes an error event
func (ch *Channel) publishError(errMsg string) {
	if ch.eventBus == nil {
		return
	}

	ch.eventBus.Publish(bus.NewEvent(bus.EventErrorOccurred, "", map[string]interface{}{
		"channel_id": ch.id,
		"error":      errMsg,
	}))
}

// GetConfig returns the channel configuration
func (ch *Channel) GetConfig() Config {
	return *ch.config
}

// UpdateConfig updates the channel configuration
func (ch *Channel) UpdateConfig(config Config) error {
	if ch.status == channels.StatusConnected {
		return fmt.Errorf("cannot update config while connected")
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	ch.config = &config
	ch.handler = NewHandler(ch.config, ch.client, ch.eventBus)

	return nil
}
