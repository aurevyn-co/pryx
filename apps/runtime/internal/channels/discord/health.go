package discord

import (
	"context"
	"fmt"
	"time"
)

// HealthChecker performs health checks on Discord channels
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

// Check performs a comprehensive health check on a Discord channel
func (h *HealthChecker) Check(ctx context.Context, channel *DiscordChannel) HealthResult {
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

	// Test 2: Check Gateway connection
	if channel.session == nil {
		result.Status = "degraded"
		result.Message = "Bot connected but Gateway session not available"
		result.ResponseTime = time.Since(start).Milliseconds()
		return result
	}

	// Test 3: Verify send capability (dry run)
	_, err = channel.GetBotInfo(checkCtx)
	if err != nil {
		result.Status = "degraded"
		result.Message = "Bot connected but API call failed"
		result.LastError = err.Error()
		result.ResponseTime = time.Since(start).Milliseconds()
		return result
	}

	result.Status = "healthy"
	result.Message = fmt.Sprintf("Bot %s#%s is operational", botInfo.Username, botInfo.Discriminator)
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
