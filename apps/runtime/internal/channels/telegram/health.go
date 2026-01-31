package telegram

import (
	"context"
	"fmt"
	"time"
)

// HealthChecker performs health checks on Telegram channels
type HealthChecker struct {
	timeout time.Duration
}

// HealthResult represents the result of a health check
type HealthResult struct {
	ChannelID    string    `json:"channel_id"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	LastError    string    `json:"last_error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	ResponseTime int64     `json:"response_time_ms"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		timeout: 10 * time.Second,
	}
}

// Check performs a comprehensive health check on a Telegram channel
func (h *HealthChecker) Check(ctx context.Context, channel *Channel) HealthResult {
	start := time.Now()
	result := HealthResult{
		ChannelID: channel.ID(),
		Timestamp: start,
	}

	// Check channel status
	if channel.Status() != "connected" {
		result.Status = "disconnected"
		result.Message = fmt.Sprintf("Channel status: %s", channel.Status())
		return result
	}

	// Create context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// Test 1: Verify bot token (GetMe)
	botInfo, err := channel.GetBotInfo(checkCtx)
	if err != nil {
		result.Status = "unhealthy"
		result.Message = "Failed to validate bot token"
		result.LastError = err.Error()
		result.ResponseTime = time.Since(start).Milliseconds()
		return result
	}

	// Test 2: Check webhook status (if in webhook mode)
	config := channel.GetConfig()
	if config.Mode == "webhook" {
		webhookInfo, err := channel.client.GetWebhookInfo(checkCtx)
		if err != nil {
			result.Status = "degraded"
			result.Message = "Bot connected but webhook check failed"
			result.LastError = err.Error()
			result.ResponseTime = time.Since(start).Milliseconds()
			return result
		}

		if webhookInfo.URL == "" {
			result.Status = "degraded"
			result.Message = "Bot connected but webhook not set"
			result.ResponseTime = time.Since(start).Milliseconds()
			return result
		}

		if webhookInfo.LastErrorMessage != "" {
			result.Status = "degraded"
			result.Message = fmt.Sprintf("Webhook has errors: %s", webhookInfo.LastErrorMessage)
			result.LastError = webhookInfo.LastErrorMessage
			result.ResponseTime = time.Since(start).Milliseconds()
			return result
		}
	}

	// Test 3: Verify send capability (dry run - just validate we can call API)
	// We don't actually send a message, just verify the client works
	_, err = channel.client.GetMe(checkCtx)
	if err != nil {
		result.Status = "degraded"
		result.Message = "Bot connected but API call failed"
		result.LastError = err.Error()
		result.ResponseTime = time.Since(start).Milliseconds()
		return result
	}

	result.Status = "healthy"
	result.Message = fmt.Sprintf("Bot @%s is operational", botInfo.Username)
	result.ResponseTime = time.Since(start).Milliseconds()
	return result
}

// IsHealthy returns true if the health check passed
func (r HealthResult) IsHealthy() bool {
	return r.Status == "healthy"
}

// IsDegraded returns true if the channel is operational but has issues
func (r HealthResult) IsDegraded() bool {
	return r.Status == "degraded"
}

// IsUnhealthy returns true if the channel is not operational
func (r HealthResult) IsUnhealthy() bool {
	return r.Status == "unhealthy" || r.Status == "disconnected"
}
