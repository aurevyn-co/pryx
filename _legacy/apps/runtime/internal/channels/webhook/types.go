package webhook

import (
	"time"
)

// RetryConfig defines retry behavior for outgoing webhooks
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// DefaultRetryConfig returns sensible defaults
func DefaultRetry() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  time.Second,
		MaxDelay:   time.Minute,
	}
}

// DeliveryStatus represents the status of a webhook delivery
type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "pending"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
	DeliveryStatusRetrying  DeliveryStatus = "retrying"
)

// DeliveryLog tracks a webhook delivery attempt
type DeliveryLog struct {
	ID           string
	ChannelID    string
	MessageID    string
	Status       DeliveryStatus
	Attempt      int
	Error        string
	ResponseCode int
	CreatedAt    time.Time
}

// SignatureFormat defines how webhook signatures are formatted
type SignatureFormat string

const (
	SignatureFormatGeneric SignatureFormat = "generic"
	SignatureFormatGitHub  SignatureFormat = "github"
	SignatureFormatStripe  SignatureFormat = "stripe"
)

// IncomingWebhook represents a received webhook message
type IncomingWebhook struct {
	ID        string
	ChannelID string
	Payload   []byte
	Headers   map[string]string
	Timestamp time.Time
}

// OutgoingWebhook represents an outgoing webhook request
type OutgoingWebhook struct {
	ID        string
	ChannelID string
	URL       string
	Payload   []byte
	Headers   map[string]string
	Attempt   int
	CreatedAt time.Time
}

// FilterRule defines filtering criteria for incoming webhooks
type FilterRule struct {
	Path    string
	Header  string
	Content string
}

// RateLimitConfig defines rate limiting for incoming webhooks
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// DefaultRateLimit returns sensible defaults
func DefaultRateLimit() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
	}
}
