package webhook

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HealthChecker performs health checks on webhook channels
type HealthChecker struct {
	client *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Check performs a health check on a webhook channel
func (h *HealthChecker) Check(channel *Channel) HealthResult {
	result := HealthResult{
		ChannelID: channel.ID(),
		Timestamp: time.Now(),
	}

	// Check if channel is enabled
	config := channel.Config()
	if !config.Enabled {
		result.Status = "disabled"
		result.Message = "Channel is disabled"
		return result
	}

	// Check channel status
	if channel.Status() != "connected" {
		result.Status = "disconnected"
		result.Message = fmt.Sprintf("Channel status: %s", channel.Status())
		return result
	}

	// For outgoing webhooks, verify endpoint is reachable
	if config.TargetURL != "" {
		if err := h.checkEndpoint(config.TargetURL); err != nil {
			result.Status = "unhealthy"
			result.Message = fmt.Sprintf("Endpoint check failed: %v", err)
			return result
		}
	}

	// Check recent delivery success rate
	logs := channel.GetDeliveryLogs(10)
	if len(logs) > 0 {
		successCount := 0
		for _, log := range logs {
			if log.Status == DeliveryStatusDelivered {
				successCount++
			}
		}
		successRate := float64(successCount) / float64(len(logs))
		if successRate < 0.5 {
			result.Status = "degraded"
			result.Message = fmt.Sprintf("Low success rate: %.0f%%", successRate*100)
			return result
		}
	}

	result.Status = "healthy"
	result.Message = "All checks passed"
	return result
}

// checkEndpoint performs a lightweight check on the endpoint
func (h *HealthChecker) checkEndpoint(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Any response is considered healthy (even 404, 401, etc.)
	// We just want to know the endpoint is reachable
	return nil
}

// HealthResult represents the result of a health check
type HealthResult struct {
	ChannelID string
	Status    string
	Message   string
	Timestamp time.Time
}

// IsHealthy returns true if the health check passed
func (r HealthResult) IsHealthy() bool {
	return r.Status == "healthy"
}
