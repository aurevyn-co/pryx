package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Receiver struct {
	config      WebhookConfig
	rateLimiter *RateLimiter
}

func NewReceiver(config WebhookConfig) *Receiver {
	return &Receiver{
		config:      config,
		rateLimiter: NewRateLimiter(DefaultRateLimit()),
	}
}

func (r *Receiver) Handle(req *http.Request) (*IncomingWebhook, error) {
	if !r.rateLimiter.Allow(r.config.ID) {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	if req.Method != http.MethodPost {
		return nil, fmt.Errorf("method not allowed: %s", req.Method)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	defer req.Body.Close()

	if r.config.Secret != "" {
		if err := r.verifySignature(req, body); err != nil {
			return nil, fmt.Errorf("signature verification failed: %w", err)
		}
	}

	headers := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	return &IncomingWebhook{
		ID:        generateID(),
		ChannelID: r.config.ID,
		Payload:   body,
		Headers:   headers,
		Timestamp: time.Now(),
	}, nil
}

func (r *Receiver) verifySignature(req *http.Request, body []byte) error {
	formats := []SignatureFormat{
		SignatureFormatStripe,
		SignatureFormatGitHub,
		SignatureFormatGeneric,
	}

	for _, format := range formats {
		if err := r.verifySignatureFormat(req, body, format); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no valid signature found")
}

func (r *Receiver) verifySignatureFormat(req *http.Request, body []byte, format SignatureFormat) error {
	switch format {
	case SignatureFormatStripe:
		return r.verifyStripeSignature(req, body)
	case SignatureFormatGitHub:
		return r.verifyGitHubSignature(req, body)
	case SignatureFormatGeneric:
		return r.verifyGenericSignature(req, body)
	default:
		return fmt.Errorf("unknown signature format: %s", format)
	}
}

func (r *Receiver) verifyStripeSignature(req *http.Request, body []byte) error {
	sigHeader := req.Header.Get("Stripe-Signature")
	if sigHeader == "" {
		return fmt.Errorf("no Stripe-Signature header")
	}

	parts := strings.Split(sigHeader, ",")
	var timestamp, signature string
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signature = kv[1]
		}
	}

	if timestamp == "" || signature == "" {
		return fmt.Errorf("invalid Stripe-Signature format")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp")
	}
	if time.Since(time.Unix(ts, 0)) > 5*time.Minute {
		return fmt.Errorf("timestamp too old")
	}

	payload := timestamp + "." + string(body)
	mac := hmac.New(sha256.New, []byte(r.config.Secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func (r *Receiver) verifyGitHubSignature(req *http.Request, body []byte) error {
	sigHeader := req.Header.Get("X-Hub-Signature-256")
	if sigHeader == "" {
		return fmt.Errorf("no X-Hub-Signature-256 header")
	}

	sigHeader = strings.TrimPrefix(sigHeader, "sha256=")

	mac := hmac.New(sha256.New, []byte(r.config.Secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sigHeader), []byte(expected)) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func (r *Receiver) verifyGenericSignature(req *http.Request, body []byte) error {
	sigHeader := req.Header.Get("X-Webhook-Signature")
	if sigHeader == "" {
		return fmt.Errorf("no X-Webhook-Signature header")
	}

	sigHeader = strings.TrimPrefix(sigHeader, "sha256=")

	mac := hmac.New(sha256.New, []byte(r.config.Secret))
	mac.Write(body)
	expectedHex := hex.EncodeToString(mac.Sum(nil))

	if hmac.Equal([]byte(sigHeader), []byte(expectedHex)) {
		return nil
	}

	expectedB64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if hmac.Equal([]byte(sigHeader), []byte(expectedB64)) {
		return nil
	}

	return fmt.Errorf("signature mismatch")
}

func generateID() string {
	return fmt.Sprintf("wh-%d", time.Now().UnixNano())
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	config  RateLimitConfig
	buckets map[string]*tokenBucket
	mu      sync.RWMutex
}

type tokenBucket struct {
	tokens     float64
	lastUpdate time.Time
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:  config,
		buckets: make(map[string]*tokenBucket),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		rl.buckets[key] = &tokenBucket{
			tokens:     float64(rl.config.BurstSize - 1),
			lastUpdate: time.Now(),
		}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(bucket.lastUpdate).Minutes()
	bucket.tokens = min(float64(rl.config.BurstSize), bucket.tokens+elapsed*float64(rl.config.RequestsPerMinute))
	bucket.lastUpdate = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
