// Package security provides telemetry PII redaction
package security

import (
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TelemetryRedactor handles PII redaction for telemetry data
type TelemetryRedactor struct {
	redactor   *Redactor
	enabled    bool
	redactKeys map[string]bool
	allowKeys  map[string]bool
}

// NewTelemetryRedactor creates a new telemetry redactor
func NewTelemetryRedactor() *TelemetryRedactor {
	return &TelemetryRedactor{
		redactor:   NewRedactor(),
		enabled:    true,
		redactKeys: defaultRedactKeys(),
		allowKeys:  defaultAllowKeys(),
	}
}

// defaultRedactKeys returns keys that should always be redacted
func defaultRedactKeys() map[string]bool {
	return map[string]bool{
		"password":      true,
		"secret":        true,
		"token":         true,
		"api_key":       true,
		"apikey":        true,
		"auth":          true,
		"authorization": true,
		"bearer":        true,
		"credential":    true,
		"private_key":   true,
		"client_secret": true,
		"access_token":  true,
		"refresh_token": true,
		"pin":           true,
		"passphrase":    true,
		"mnemonic":      true,
		"seed":          true,
		"email":         true,
		"phone":         true,
		"ssn":           true,
		"credit_card":   true,
		"address":       true,
		"name":          true,
		"full_name":     true,
		"first_name":    true,
		"last_name":     true,
		"user_id":       true,
		"username":      true,
		"session_id":    true,
		"ip_address":    true,
		"mac_address":   true,
		"hostname":      true,
		"device_id":     true,
	}
}

// defaultAllowKeys returns keys that are safe to log
func defaultAllowKeys() map[string]bool {
	return map[string]bool{
		"llm.provider":      true,
		"llm.model":         true,
		"llm.tokens.input":  true,
		"llm.tokens.output": true,
		"llm.tokens.total":  true,
		"tool.name":         true,
		"tool.arg.command":  true,
		"tool.arg.path":     true,
		"session.id":        true,
		"channel.type":      true,
		"channel.id":        true,
		"mcp.server":        true,
		"mcp.tool":          true,
		"cost.amount":       true,
		"cost.currency":     true,
		"duration_ms":       true,
		"action":            true,
		"status":            true,
		"error":             true,
		"version":           true,
		"platform":          true,
	}
}

// RedactAttribute redacts a single attribute
func (t *TelemetryRedactor) RedactAttribute(attr attribute.KeyValue) attribute.KeyValue {
	if !t.enabled {
		return attr
	}

	key := string(attr.Key)

	// Check if key should be redacted
	if t.shouldRedactKey(key) {
		return attribute.String(key, "[REDACTED]")
	}

	// Check if value contains sensitive data
	switch attr.Value.Type() {
	case attribute.STRING:
		redacted := t.redactor.RedactString(attr.Value.AsString())
		if redacted != attr.Value.AsString() {
			return attribute.String(key, redacted)
		}
	}

	return attr
}

// RedactAttributes redacts multiple attributes
func (t *TelemetryRedactor) RedactAttributes(attrs []attribute.KeyValue) []attribute.KeyValue {
	if !t.enabled {
		return attrs
	}

	result := make([]attribute.KeyValue, len(attrs))
	for i, attr := range attrs {
		result[i] = t.RedactAttribute(attr)
	}
	return result
}

// RedactSpanAttributes redacts all attributes on a span
func (t *TelemetryRedactor) RedactSpanAttributes(span trace.Span, attrs []attribute.KeyValue) {
	if !t.enabled || span == nil {
		return
	}

	redacted := t.RedactAttributes(attrs)
	span.SetAttributes(redacted...)
}

// shouldRedactKey checks if a key should be redacted
func (t *TelemetryRedactor) shouldRedactKey(key string) bool {
	keyLower := strings.ToLower(key)

	// Check explicit redact list
	if t.redactKeys[keyLower] {
		return true
	}

	// Check for sensitive patterns in key name
	sensitivePatterns := []string{
		"password", "secret", "token", "api_key", "apikey",
		"auth", "credential", "private", "bearer",
	}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(keyLower, pattern) {
			// Unless explicitly allowed
			if !t.allowKeys[keyLower] {
				return true
			}
		}
	}

	return false
}

// RedactString redacts a string value
func (t *TelemetryRedactor) RedactString(value string) string {
	if !t.enabled {
		return value
	}
	return t.redactor.RedactString(value)
}

// SetEnabled enables or disables redaction
func (t *TelemetryRedactor) SetEnabled(enabled bool) {
	t.enabled = enabled
}

// AddRedactKey adds a key to the redaction list
func (t *TelemetryRedactor) AddRedactKey(key string) {
	t.redactKeys[strings.ToLower(key)] = true
}

// AddAllowKey adds a key to the allow list (overrides redaction)
func (t *TelemetryRedactor) AddAllowKey(key string) {
	t.allowKeys[strings.ToLower(key)] = true
}

// TelemetrySpan wraps a span with automatic redaction
type TelemetrySpan struct {
	trace.Span
	redactor *TelemetryRedactor
}

// SetAttributes sets attributes with automatic redaction
func (s *TelemetrySpan) SetAttributes(attrs ...attribute.KeyValue) {
	redacted := s.redactor.RedactAttributes(attrs)
	s.Span.SetAttributes(redacted...)
}

// SetAttributesWithRedaction sets attributes with custom redaction function
func (s *TelemetrySpan) SetAttributesWithRedaction(attrs []attribute.KeyValue, redactFn func(attribute.KeyValue) attribute.KeyValue) {
	result := make([]attribute.KeyValue, len(attrs))
	for i, attr := range attrs {
		result[i] = redactFn(attr)
	}
	s.Span.SetAttributes(result...)
}

// WrapSpan wraps a span with redaction capabilities
func (t *TelemetryRedactor) WrapSpan(span trace.Span) *TelemetrySpan {
	return &TelemetrySpan{
		Span:     span,
		redactor: t,
	}
}

// Global telemetry redactor
var globalTelemetryRedactor = NewTelemetryRedactor()

// RedactTelemetryAttribute redacts using global redactor
func RedactTelemetryAttribute(attr attribute.KeyValue) attribute.KeyValue {
	return globalTelemetryRedactor.RedactAttribute(attr)
}

// RedactTelemetryAttributes redacts using global redactor
func RedactTelemetryAttributes(attrs []attribute.KeyValue) []attribute.KeyValue {
	return globalTelemetryRedactor.RedactAttributes(attrs)
}

// RedactSpanString redacts a string for telemetry
func RedactSpanString(value string) string {
	return globalTelemetryRedactor.RedactString(value)
}

// PII Patterns for common sensitive data
var (
	// Email pattern
	EmailPattern = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)

	// Phone pattern
	PhonePattern = regexp.MustCompile(`\b\+?\d{1,3}[-.\s]?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}\b`)

	// Credit card pattern
	CreditCardPattern = regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`)

	// SSN pattern
	SSNPattern = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
)

// ContainsPII checks if a string contains potential PII
func ContainsPII(value string) bool {
	if EmailPattern.MatchString(value) {
		return true
	}
	if PhonePattern.MatchString(value) {
		return true
	}
	if CreditCardPattern.MatchString(value) {
		return true
	}
	if SSNPattern.MatchString(value) {
		return true
	}
	return false
}

// SanitizeForTelemetry sanitizes a map for telemetry export
func SanitizeForTelemetry(data map[string]interface{}) map[string]interface{} {
	return globalTelemetryRedactor.redactor.RedactMap(data)
}
