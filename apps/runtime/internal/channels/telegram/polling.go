package telegram

import (
	"context"
	"fmt"
	"time"

	"pryx-core/internal/bus"
)

// Poller implements long-polling for Telegram updates
type Poller struct {
	config   *Config
	client   *Client
	handler  *Handler
	eventBus *bus.Bus
	offset   int
	running  bool
}

// NewPoller creates a new polling instance
func NewPoller(config *Config, client *Client, handler *Handler, eventBus *bus.Bus) *Poller {
	return &Poller{
		config:   config,
		client:   client,
		handler:  handler,
		eventBus: eventBus,
		offset:   0,
		running:  false,
	}
}

// Start begins the polling loop
func (p *Poller) Start(ctx context.Context) {
	if p.running {
		return
	}
	p.running = true
	defer func() { p.running = false }()

	// Initial delay before starting
	ticker := time.NewTicker(p.config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				p.publishError(fmt.Sprintf("polling error: %v", err))
			}
		}
	}
}

// poll fetches updates from Telegram
func (p *Poller) poll(ctx context.Context) error {
	// Calculate timeout for long polling (should be less than ticker interval)
	timeout := int(p.config.PollingInterval.Seconds()) - 5
	if timeout < 5 {
		timeout = 5
	}

	updates, err := p.client.GetUpdates(ctx, p.offset, 100, timeout, p.config.AllowedUpdates)
	if err != nil {
		return err
	}

	for _, update := range updates {
		// Update offset to acknowledge this update
		if update.UpdateID >= p.offset {
			p.offset = update.UpdateID + 1
		}

		// Process update asynchronously
		go func(u Update) {
			if err := p.handler.HandleUpdate(ctx, &u); err != nil {
				p.publishError(fmt.Sprintf("handle update error: %v", err))
			}
		}(update)
	}

	return nil
}

// Stop stops the polling loop
func (p *Poller) Stop() {
	p.running = false
}

// IsRunning returns whether the poller is running
func (p *Poller) IsRunning() bool {
	return p.running
}

// GetOffset returns the current update offset
func (p *Poller) GetOffset() int {
	return p.offset
}

// SetOffset sets the update offset
func (p *Poller) SetOffset(offset int) {
	p.offset = offset
}

// publishError publishes an error event
func (p *Poller) publishError(errMsg string) {
	if p.eventBus == nil {
		return
	}

	p.eventBus.Publish(bus.NewEvent(bus.EventErrorOccurred, "", map[string]interface{}{
		"channel_id": p.config.ID,
		"error":      errMsg,
	}))
}
