package webhook

import (
	"context"
	"fmt"

	"pryx-core/internal/bus"
	"pryx-core/internal/channels"
)

// Channel implements the channels.Channel interface for webhooks
type Channel struct {
	config   WebhookConfig
	receiver *Receiver
	sender   *Sender
	bus      *bus.Bus
	status   channels.Status
	logs     *LogStore
}

// NewChannel creates a new webhook channel
func NewChannel(config WebhookConfig, eventBus *bus.Bus) *Channel {
	return &Channel{
		config:   config,
		receiver: NewReceiver(config),
		sender:   NewSender(config),
		bus:      eventBus,
		status:   channels.StatusDisconnected,
		logs:     NewLogStore(),
	}
}

// ID returns the channel ID
func (c *Channel) ID() string {
	return c.config.ID
}

// Type returns the channel type
func (c *Channel) Type() string {
	return "webhook"
}

// Connect initializes the channel
func (c *Channel) Connect(ctx context.Context) error {
	if !c.config.Enabled {
		return fmt.Errorf("channel is disabled")
	}
	c.status = channels.StatusConnected
	return nil
}

// Disconnect closes the channel
func (c *Channel) Disconnect(ctx context.Context) error {
	c.status = channels.StatusDisconnected
	return nil
}

// Send sends a message via webhook
func (c *Channel) Send(ctx context.Context, msg channels.Message) error {
	if c.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}

	payload := []byte(msg.Content)
	log, err := c.sender.Send(ctx, payload)
	if err != nil {
		c.status = channels.StatusError
		c.logs.Add(log)
		return err
	}

	c.logs.Add(log)

	if c.bus != nil {
		c.bus.Publish(bus.NewEvent(bus.EventTraceEvent, "", map[string]interface{}{
			"kind":        "webhook.sent",
			"channel_id":  c.config.ID,
			"message_id":  msg.ID,
			"delivery_id": log.ID,
			"status":      log.Status,
		}))
	}

	return nil
}

// Status returns the channel status
func (c *Channel) Status() channels.Status {
	return c.status
}

// Health checks the channel health
func (c *Channel) Health() error {
	if c.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}
	if !c.config.Enabled {
		return fmt.Errorf("channel is disabled")
	}
	return nil
}

// Receive processes an incoming webhook
func (c *Channel) Receive(msg *IncomingWebhook) error {
	if c.status != channels.StatusConnected {
		return fmt.Errorf("channel not connected")
	}

	channelMsg := channels.Message{
		ID:        msg.ID,
		Content:   string(msg.Payload),
		Source:    c.config.ID,
		ChannelID: msg.ChannelID,
		SenderID:  "webhook",
		Metadata:  msg.Headers,
		CreatedAt: msg.Timestamp,
	}

	if c.bus != nil {
		c.bus.Publish(bus.NewEvent(bus.EventChannelMessage, "", channelMsg))
	}

	return nil
}

// Config returns the channel configuration
func (c *Channel) Config() WebhookConfig {
	return c.config
}

// GetDeliveryLogs returns recent delivery logs
func (c *Channel) GetDeliveryLogs(limit int) []DeliveryLog {
	return c.logs.GetRecent(limit)
}
