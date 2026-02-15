package security

import (
	"strings"
	"testing"
)

func TestRedactString_APIKeys(t *testing.T) {
	redactor := NewRedactor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "OpenAI API key",
			input:    "Using key: sk-abc123def456ghi789jkl012mno345pqr678stu901vwx234yz",
			expected: "Using key: [REDACTED_CREDENTIAL]",
		},
		{
			name:     "Bearer token",
			input:    "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected: "Authorization: Bearer [REDACTED_CREDENTIAL]",
		},
		{
			name:     "Password in string",
			input:    "password=mysecretpass123",
			expected: "password=[REDACTED_CREDENTIAL]",
		},
		{
			name:     "API key with quotes",
			input:    `api_key="supersecrettoken12345"`,
			expected: `api_key=[REDACTED_CREDENTIAL]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactString(tt.input)
			if result != tt.expected {
				t.Errorf("RedactString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRedactString_PII(t *testing.T) {
	redactor := NewRedactor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Email address",
			input:    "Contact me at john.doe@example.com please",
			expected: "Contact me at [REDACTED_PII] please",
		},
		{
			name:     "Phone number",
			input:    "Call me at +1-555-123-4567",
			expected: "Call me at +[REDACTED_PII]",
		},
		{
			name:     "Credit card",
			input:    "Card: 1234-5678-9012-3456",
			expected: "Card: [REDACTED_PII]",
		},
		{
			name:     "SSN",
			input:    "SSN: 123-45-6789",
			expected: "SSN: [REDACTED_PII]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactor.RedactString(tt.input)
			if result != tt.expected {
				t.Errorf("RedactString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRedactMap(t *testing.T) {
	redactor := NewRedactor()

	input := map[string]interface{}{
		"username":   "john_doe",
		"password":   "secret123",
		"api_key":    "sk-abc123",
		"email":      "john@example.com",
		"safe_field": "this is fine",
	}

	result := redactor.RedactMap(input)

	if result["username"] != "john_doe" {
		t.Errorf("username should not be redacted, got %v", result["username"])
	}

	if result["password"] != "[REDACTED]" {
		t.Errorf("password should be redacted, got %v", result["password"])
	}

	if result["api_key"] != "[REDACTED]" {
		t.Errorf("api_key should be redacted, got %v", result["api_key"])
	}

	if result["safe_field"] != "this is fine" {
		t.Errorf("safe_field should not be changed, got %v", result["safe_field"])
	}
}

func TestRedactMap_Nested(t *testing.T) {
	redactor := NewRedactor()

	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name":     "John",
			"password": "secret",
		},
		"config": map[string]interface{}{
			"token": "abc123",
		},
	}

	result := redactor.RedactMap(input)

	userMap := result["user"].(map[string]interface{})
	if userMap["password"] != "[REDACTED]" {
		t.Errorf("nested password should be redacted, got %v", userMap["password"])
	}

	configMap := result["config"].(map[string]interface{})
	if configMap["token"] != "[REDACTED]" {
		t.Errorf("nested token should be redacted, got %v", configMap["token"])
	}
}

func TestIsSensitiveKey(t *testing.T) {
	redactor := NewRedactor()

	sensitiveKeys := []string{
		"password",
		"api_key",
		"SECRET_TOKEN",
		"auth_header",
		"private_key",
	}

	nonSensitiveKeys := []string{
		"username",
		"count",
		"enabled",
		"timestamp",
	}

	for _, key := range sensitiveKeys {
		if !redactor.isSensitiveKey(key) {
			t.Errorf("key %q should be sensitive", key)
		}
	}

	for _, key := range nonSensitiveKeys {
		if redactor.isSensitiveKey(key) {
			t.Errorf("key %q should NOT be sensitive", key)
		}
	}
}

func TestRedactor_SetEnabled(t *testing.T) {
	redactor := NewRedactor()

	redactor.SetEnabled(false)
	if redactor.IsEnabled() {
		t.Error("redactor should be disabled")
	}

	input := "password=secret123"
	result := redactor.RedactString(input)
	if result != input {
		t.Error("when disabled, redaction should not occur")
	}

	redactor.SetEnabled(true)
	if !redactor.IsEnabled() {
		t.Error("redactor should be enabled")
	}
}

func TestGlobalRedactor(t *testing.T) {
	input := "api_key=thisisaverylongsecretkey123456"
	result := RedactString(input)

	if result == input {
		t.Error("global RedactString should redact")
	}

	if !strings.Contains(result, "REDACTED") {
		t.Errorf("redacted result should contain REDACTED, got %s", result)
	}
}

func BenchmarkRedactString(b *testing.B) {
	redactor := NewRedactor()
	input := "User email is test@example.com and API key is sk-abc123def456"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		redactor.RedactString(input)
	}
}
