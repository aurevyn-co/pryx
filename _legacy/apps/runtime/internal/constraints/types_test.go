package constraints

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestActionConstants tests action constants
func TestActionConstants(t *testing.T) {
	assert.Equal(t, Action("allow"), ActionAllow)
	assert.Equal(t, Action("deny"), ActionDeny)
	assert.Equal(t, Action("fallback"), ActionFallback)
	assert.Equal(t, Action("ask"), ActionAsk)
}

// TestResolutionStructure tests Resolution struct
func TestResolutionStructure(t *testing.T) {
	res := Resolution{
		Action:           ActionAllow,
		TargetModel:      "gpt-4o",
		Reason:           "Model within budget",
		EstimatedCostUSD: 0.05,
	}

	assert.Equal(t, ActionAllow, res.Action)
	assert.Equal(t, "gpt-4o", res.TargetModel)
	assert.Equal(t, "Model within budget", res.Reason)
	assert.Equal(t, 0.05, res.EstimatedCostUSD)
}

// TestRequestStructure tests Request struct
func TestRequestStructure(t *testing.T) {
	req := Request{
		Model:          "gpt-4o",
		ProviderID:     "openai",
		PromptTokens:   1000,
		OutputTokens:   500,
		ThinkingTokens: 0,
		Tools:          []string{"weather", "calculator"},
		Images:         false,
		MaxCostUSD:     0.10,
	}

	assert.Equal(t, "gpt-4o", req.Model)
	assert.Equal(t, "openai", req.ProviderID)
	assert.Equal(t, 1000, req.PromptTokens)
	assert.Equal(t, 500, req.OutputTokens)
	assert.Len(t, req.Tools, 2)
	assert.False(t, req.Images)
	assert.Equal(t, 0.10, req.MaxCostUSD)
}

// TestCostEstimateStructure tests CostEstimate struct
func TestCostEstimateStructure(t *testing.T) {
	estimate := CostEstimate{
		InputTokensCost:    0.01,
		OutputTokensCost:   0.005,
		ThinkingTokensCost: 0.002,
		FixedCost:          0.001,
		TotalUSD:           0.018,
	}

	assert.Equal(t, 0.01, estimate.InputTokensCost)
	assert.Equal(t, 0.005, estimate.OutputTokensCost)
	assert.Equal(t, 0.002, estimate.ThinkingTokensCost)
	assert.Equal(t, 0.001, estimate.FixedCost)
	assert.Equal(t, 0.018, estimate.TotalUSD)
}

// TestCostEstimateCalculation tests cost calculation logic
func TestCostEstimateCalculation(t *testing.T) {
	estimate := CostEstimate{
		InputTokensCost:    0.01,
		OutputTokensCost:   0.005,
		ThinkingTokensCost: 0.002,
		FixedCost:          0.001,
		TotalUSD:           0.018,
	}

	// Verify total is sum of components
	expectedTotal := estimate.InputTokensCost + estimate.OutputTokensCost +
		estimate.ThinkingTokensCost + estimate.FixedCost
	assert.InDelta(t, expectedTotal, estimate.TotalUSD, 0.0001)
}

// TestRequestWithImages tests Request with images
func TestRequestWithImages(t *testing.T) {
	req := Request{
		Model:        "gpt-4o",
		ProviderID:   "openai",
		PromptTokens: 2000,
		OutputTokens: 1000,
		Images:       true,
	}

	assert.True(t, req.Images)
	assert.Equal(t, 2000, req.PromptTokens)
}

// TestRequestEmptyTools tests Request with no tools
func TestRequestEmptyTools(t *testing.T) {
	req := Request{
		Model:        "gpt-3.5-turbo",
		ProviderID:   "openai",
		PromptTokens: 500,
		OutputTokens: 250,
		Tools:        []string{},
		Images:       false,
	}

	assert.Empty(t, req.Tools)
	assert.False(t, req.Images)
}

// TestResolutionWithDifferentActions tests Resolution with different actions
func TestResolutionWithDifferentActions(t *testing.T) {
	tests := []struct {
		name   string
		action Action
	}{
		{"Allow action", ActionAllow},
		{"Deny action", ActionDeny},
		{"Fallback action", ActionFallback},
		{"Ask action", ActionAsk},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Resolution{Action: tt.action}
			assert.Equal(t, tt.action, res.Action)
		})
	}
}

// TestResolutionNoTargetModel tests Resolution without target model
func TestResolutionNoTargetModel(t *testing.T) {
	res := Resolution{
		Action: ActionDeny,
		Reason: "Cost exceeds budget",
	}

	assert.Empty(t, res.TargetModel)
	assert.Equal(t, "Cost exceeds budget", res.Reason)
}

// TestRequestZeroTokens tests Request with zero tokens
func TestRequestZeroTokens(t *testing.T) {
	req := Request{
		Model:        "test-model",
		ProviderID:   "test",
		PromptTokens: 0,
		OutputTokens: 0,
	}

	assert.Equal(t, 0, req.PromptTokens)
	assert.Equal(t, 0, req.OutputTokens)
}

// TestCostEstimateZeroCosts tests CostEstimate with zero costs
func TestCostEstimateZeroCosts(t *testing.T) {
	estimate := CostEstimate{
		InputTokensCost:    0,
		OutputTokensCost:   0,
		ThinkingTokensCost: 0,
		FixedCost:          0,
		TotalUSD:           0,
	}

	assert.Equal(t, 0.0, estimate.TotalUSD)
	assert.Equal(t, 0.0, estimate.InputTokensCost)
	assert.Equal(t, 0.0, estimate.OutputTokensCost)
}
