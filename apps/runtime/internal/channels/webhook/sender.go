package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Sender struct {
	config WebhookConfig
	client *http.Client
}

func NewSender(config WebhookConfig) *Sender {
	return &Sender{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (s *Sender) Send(ctx context.Context, payload []byte) (*DeliveryLog, error) {
	if s.config.TargetURL == "" {
		return nil, fmt.Errorf("no target URL configured")
	}

	log := &DeliveryLog{
		ID:        generateID(),
		ChannelID: s.config.ID,
		Status:    DeliveryStatusPending,
		CreatedAt: time.Now(),
	}

	retryConfig := s.config.RetryConfig
	if retryConfig.MaxRetries == 0 {
		retryConfig = DefaultRetry()
	}

	var lastErr error
	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		log.Attempt = attempt + 1
		log.Status = DeliveryStatusRetrying

		if attempt > 0 {
			delay := calculateBackoff(attempt, retryConfig)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				log.Status = DeliveryStatusFailed
				log.Error = "context cancelled during retry"
				return log, ctx.Err()
			}
		}

		respCode, err := s.sendRequest(ctx, payload)
		log.ResponseCode = respCode

		if err == nil && respCode < 400 {
			log.Status = DeliveryStatusDelivered
			return log, nil
		}

		if err != nil {
			lastErr = err
			log.Error = err.Error()
		} else {
			lastErr = fmt.Errorf("HTTP %d", respCode)
			log.Error = lastErr.Error()
		}

		if !shouldRetry(respCode) {
			break
		}
	}

	log.Status = DeliveryStatusFailed
	return log, fmt.Errorf("failed after %d attempts: %w", retryConfig.MaxRetries+1, lastErr)
}

func (s *Sender) sendRequest(ctx context.Context, payload []byte) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.TargetURL, bytes.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Pryx-Webhook/1.0")

	for name, value := range s.config.Headers {
		req.Header.Set(name, value)
	}

	if s.config.Secret != "" {
		mac := hmac.New(sha256.New, []byte(s.config.Secret))
		mac.Write(payload)
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Webhook-Signature", "sha256="+signature)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode, nil
}

func (s *Sender) SendJSON(ctx context.Context, data interface{}) (*DeliveryLog, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return s.Send(ctx, payload)
}

func calculateBackoff(attempt int, config RetryConfig) time.Duration {
	delay := config.BaseDelay * (1 << uint(attempt-1))
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	return delay
}

func shouldRetry(statusCode int) bool {
	if statusCode >= 500 {
		return true
	}
	if statusCode == 429 {
		return true
	}
	if statusCode == 0 {
		return true
	}
	return false
}
