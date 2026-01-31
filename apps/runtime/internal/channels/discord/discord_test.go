package discord

import (
	"testing"
)

// TestEmbedBuilder tests the embed builder
func TestEmbedBuilder(t *testing.T) {
	// Test basic embed creation
	embed := NewEmbed().
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetColor(ColorBlue).
		Build()

	if embed.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", embed.Title)
	}

	if embed.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got '%s'", embed.Description)
	}

	if embed.Color != ColorBlue {
		t.Errorf("Expected color %d, got %d", ColorBlue, embed.Color)
	}
}

// TestEmbedBuilderFields tests adding fields
func TestEmbedBuilderFields(t *testing.T) {
	embed := NewEmbed().
		SetTitle("Field Test").
		AddField("Name1", "Value1", false).
		AddInlineField("Name2", "Value2").
		Build()

	if len(embed.Fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(embed.Fields))
	}

	if embed.Fields[0].Name != "Name1" {
		t.Errorf("Expected field name 'Name1', got '%s'", embed.Fields[0].Name)
	}

	if !embed.Fields[1].Inline {
		t.Error("Expected second field to be inline")
	}
}

// TestEmbedBuilderValidation tests embed validation
func TestEmbedBuilderValidation(t *testing.T) {
	// Test title limit (256 chars)
	longTitle := make([]byte, 300)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	embed := NewEmbed().SetTitle(string(longTitle))
	if err := embed.Validate(); err == nil {
		t.Error("Expected validation error for long title")
	}

	// Test valid embed
	validEmbed := NewEmbed().
		SetTitle("Valid Title").
		SetDescription("Valid Description").
		AddField("Field", "Value", false)

	if err := validEmbed.Validate(); err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}
}

// TestHexToInt tests hex color conversion
func TestHexToInt(t *testing.T) {
	tests := []struct {
		hex      string
		expected int
	}{
		{"#FF0000", 0xFF0000},
		{"#00FF00", 0x00FF00},
		{"#0000FF", 0x0000FF},
		{"#FFFFFF", 0xFFFFFF},
		{"#000000", 0x000000},
		{"FF0000", 0xFF0000},
	}

	for _, test := range tests {
		result := hexToInt(test.hex)
		if result != test.expected {
			t.Errorf("hexToInt(%s) = %d, expected %d", test.hex, result, test.expected)
		}
	}
}

// TestDefaultIntents tests default intents
func TestDefaultIntents(t *testing.T) {
	intents := DefaultIntents()

	// Should include basic intents
	if intents&IntentGuilds == 0 {
		t.Error("Default intents should include Guilds")
	}

	if intents&IntentGuildMessages == 0 {
		t.Error("Default intents should include GuildMessages")
	}

	if intents&IntentDirectMessages == 0 {
		t.Error("Default intents should include DirectMessages")
	}

	if intents&IntentMessageContent == 0 {
		t.Error("Default intents should include MessageContent")
	}
}

// TestConfigValidation tests config validation
func TestConfigValidation(t *testing.T) {
	// Valid config
	validConfig := Config{
		ID:       "test-id",
		Name:     "Test Bot",
		TokenRef: "vault://token",
		Intents:  DefaultIntents(),
	}

	if err := validConfig.Validate(); err != nil {
		t.Errorf("Valid config should not error: %v", err)
	}

	// Missing ID
	invalidConfig := Config{
		Name:     "Test Bot",
		TokenRef: "vault://token",
	}

	if err := invalidConfig.Validate(); err == nil {
		t.Error("Config without ID should error")
	}

	// Missing name
	invalidConfig2 := Config{
		ID:       "test-id",
		TokenRef: "vault://token",
	}

	if err := invalidConfig2.Validate(); err == nil {
		t.Error("Config without name should error")
	}

	// Missing token ref
	invalidConfig3 := Config{
		ID:   "test-id",
		Name: "Test Bot",
	}

	if err := invalidConfig3.Validate(); err == nil {
		t.Error("Config without token_ref should error")
	}
}

// TestConfigWhitelist tests whitelist functionality
func TestConfigWhitelist(t *testing.T) {
	config := Config{
		ID:            "test-id",
		Name:          "Test Bot",
		TokenRef:      "vault://token",
		AllowedGuilds: []string{"guild1", "guild2"},
	}

	// Test guild whitelist
	if !config.IsGuildAllowed("guild1") {
		t.Error("Should allow guild1")
	}

	if !config.IsGuildAllowed("guild2") {
		t.Error("Should allow guild2")
	}

	if config.IsGuildAllowed("guild3") {
		t.Error("Should not allow guild3")
	}

	// Test empty whitelist (allow all)
	config2 := Config{
		ID:            "test-id",
		Name:          "Test Bot",
		TokenRef:      "vault://token",
		AllowedGuilds: []string{},
	}

	if !config2.IsGuildAllowed("any-guild") {
		t.Error("Empty whitelist should allow all guilds")
	}
}

// TestPrebuiltEmbeds tests pre-built embed templates
func TestPrebuiltEmbeds(t *testing.T) {
	// Success embed
	success := SuccessEmbed("Success", "Operation completed").Build()
	if success.Color != ColorGreen {
		t.Errorf("Success embed should be green, got %d", success.Color)
	}

	// Error embed
	errEmbed := ErrorEmbed("Error", "Something went wrong").Build()
	if errEmbed.Color != ColorRed {
		t.Errorf("Error embed should be red, got %d", errEmbed.Color)
	}

	// Info embed
	info := InfoEmbed("Info", "Here is some information").Build()
	if info.Color != ColorBlue {
		t.Errorf("Info embed should be blue, got %d", info.Color)
	}

	// Warning embed
	warning := WarningEmbed("Warning", "Be careful").Build()
	if warning.Color != ColorGold {
		t.Errorf("Warning embed should be gold, got %d", warning.Color)
	}
}

// BenchmarkEmbedBuilder benchmarks embed creation
func BenchmarkEmbedBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewEmbed().
			SetTitle("Benchmark Title").
			SetDescription("Benchmark Description").
			SetColor(ColorBlue).
			AddField("Field1", "Value1", false).
			AddInlineField("Field2", "Value2").
			Build()
	}
}
