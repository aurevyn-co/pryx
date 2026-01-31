// Package security provides security utilities for Pryx runtime
package security

import (
	"regexp"
	"strings"
	"sync"
)

// Sensitive patterns for redaction
var (
	// API keys, tokens, secrets patterns
	apiKeyPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key[:=]\s*)["']?[a-zA-Z0-9_\-]{16,}["']?`),
		regexp.MustCompile(`(?i)(bearer\s+)[a-zA-Z0-9_\-\.]{20,}`),
		regexp.MustCompile(`(?i)(token[:=]\s*)["']?[a-zA-Z0-9_\-]{16,}["']?`),
		regexp.MustCompile(`(?i)(secret[:=]\s*)["']?[a-zA-Z0-9_\-]{8,}["']?`),
		regexp.MustCompile(`(?i)(password[:=]\s*)["']?[^"'\s]{4,}["']?`),
		regexp.MustCompile(`(?i)(authorization[:=]\s*basic\s+)[a-zA-Z0-9+/=]{10,}`),
		regexp.MustCompile(`(?i)(sk-[a-zA-Z0-9]{20,})`),         // OpenAI key pattern
		regexp.MustCompile(`(?i)(sk-ant-[a-zA-Z0-9_-]{20,})`),   // Anthropic key pattern
		regexp.MustCompile(`(?i)([a-f0-9]{32,})`),               // Hex tokens
		regexp.MustCompile(`(?i)(gh[pousr]_[A-Za-z0-9_]{36,})`), // GitHub tokens
	}

	// PII patterns
	piiPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),      // Email
		regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),                                    // SSN
		regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),               // Credit card
		regexp.MustCompile(`\b\+?\d{1,3}[-.\s]?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}\b`), // Phone
		regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`),                        // IP address
	}

	// Sensitive keys in maps/objects
	sensitiveKeys = []string{
		"password", "secret", "token", "key", "auth", "credential",
		"private", "apikey", "api_key", "bearer", "authorization",
		"access_token", "refresh_token", "client_secret", "client_id",
		"pin", "passphrase", "seed", "mnemonic", "wallet",
	}
)

// Redactor handles redaction of sensitive information
type Redactor struct {
	mu             sync.RWMutex
	customPatterns []*regexp.Regexp
	enabled        bool
}

// NewRedactor creates a new redactor
func NewRedactor() *Redactor {
	return &Redactor{
		enabled: true,
	}
}

// RedactString redacts sensitive information from a string
func (r *Redactor) RedactString(input string) string {
	if !r.enabled || input == "" {
		return input
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := input

	// Redact API keys and tokens - preserve prefix, redact value
	for _, pattern := range apiKeyPatterns {
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			// Find where the actual credential starts
			loc := pattern.FindStringIndex(match)
			if loc == nil {
				return match
			}
			// Try to find prefix separator
			for _, sep := range []string{"=", ":", " "} {
				if idx := strings.Index(match, sep); idx != -1 {
					return match[:idx+len(sep)] + "[REDACTED_CREDENTIAL]"
				}
			}
			return "[REDACTED_CREDENTIAL]"
		})
	}

	// Redact PII
	for _, pattern := range piiPatterns {
		result = pattern.ReplaceAllString(result, "[REDACTED_PII]")
	}

	// Apply custom patterns
	for _, pattern := range r.customPatterns {
		result = pattern.ReplaceAllString(result, "[REDACTED_CUSTOM]")
	}

	return result
}

// RedactMap redacts sensitive values from a map
func (r *Redactor) RedactMap(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	result := make(map[string]interface{}, len(data))
	for key, value := range data {
		if r.isSensitiveKey(key) {
			result[key] = "[REDACTED]"
		} else {
			result[key] = r.redactValue(value)
		}
	}
	return result
}

// redactValue recursively redacts values
func (r *Redactor) redactValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return r.RedactString(v)
	case map[string]interface{}:
		return r.RedactMap(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = r.redactValue(item)
		}
		return result
	default:
		return value
	}
}

// isSensitiveKey checks if a key might contain sensitive information
func (r *Redactor) isSensitiveKey(key string) bool {
	keyLower := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(keyLower, sensitive) {
			return true
		}
	}
	return false
}

// AddPattern adds a custom redaction pattern
func (r *Redactor) AddPattern(pattern *regexp.Regexp) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.customPatterns = append(r.customPatterns, pattern)
}

// SetEnabled enables or disables redaction
func (r *Redactor) SetEnabled(enabled bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.enabled = enabled
}

// IsEnabled returns whether redaction is enabled
func (r *Redactor) IsEnabled() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.enabled
}

// Global redactor instance
var globalRedactor = NewRedactor()

// RedactString redacts using the global redactor
func RedactString(input string) string {
	return globalRedactor.RedactString(input)
}

// RedactMap redacts using the global redactor
func RedactMap(data map[string]interface{}) map[string]interface{} {
	return globalRedactor.RedactMap(data)
}

// AddGlobalPattern adds a pattern to the global redactor
func AddGlobalPattern(pattern *regexp.Regexp) {
	globalRedactor.AddPattern(pattern)
}

// SetGlobalEnabled sets enabled state on global redactor
func SetGlobalEnabled(enabled bool) {
	globalRedactor.SetEnabled(enabled)
}
