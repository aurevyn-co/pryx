package cost

import (
	"testing"

	"pryx-core/internal/llm"

	"github.com/stretchr/testify/assert"
)

// TestCostCalculatorCalculateFromUsage tests cost calculation from LLM usage
func TestCostCalculatorCalculateFromUsage(t *testing.T) {
	pricingManager := NewPricingManager()
	calculator := NewCostCalculator(pricingManager)

	tests := []struct {
		name       string
		modelID    string
		usage      llm.Usage
		expectCost bool
	}{
		{
			name:    "GPT-4o small usage",
			modelID: "gpt-4o",
			usage: llm.Usage{
				PromptTokens:     1000,
				CompletionTokens: 500,
				TotalTokens:      1500,
			},
			expectCost: true,
		},
		{
			name:    "GPT-4o-mini usage",
			modelID: "gpt-4o-mini",
			usage: llm.Usage{
				PromptTokens:     10000,
				CompletionTokens: 5000,
				TotalTokens:      15000,
			},
			expectCost: true,
		},
		{
			name:    "Unknown model",
			modelID: "unknown-model",
			usage: llm.Usage{
				PromptTokens:     1000,
				CompletionTokens: 500,
				TotalTokens:      1500,
			},
			expectCost: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			costInfo, err := calculator.CalculateFromUsage(tt.modelID, tt.usage)
			assert.NoError(t, err)

			if tt.expectCost {
				assert.Greater(t, costInfo.TotalCost, 0.0, "cost should be positive")
				assert.Equal(t, int64(tt.usage.PromptTokens), costInfo.InputTokens)
				assert.Equal(t, int64(tt.usage.CompletionTokens), costInfo.OutputTokens)
				assert.Equal(t, int64(tt.usage.TotalTokens), costInfo.TotalTokens)
			} else {
				assert.Equal(t, 0.0, costInfo.TotalCost, "unknown model should have zero cost")
			}
		})
	}
}

// TestCostCalculatorCalculateSessionCost tests session cost aggregation
func TestCostCalculatorCalculateSessionCost(t *testing.T) {
	pricingManager := NewPricingManager()
	calculator := NewCostCalculator(pricingManager)

	requests := []CostInfo{
		{
			InputTokens:  1000,
			OutputTokens: 500,
			TotalTokens:  1500,
			InputCost:    0.01,
			OutputCost:   0.005,
			TotalCost:    0.015,
		},
		{
			InputTokens:  2000,
			OutputTokens: 1000,
			TotalTokens:  3000,
			InputCost:    0.02,
			OutputCost:   0.01,
			TotalCost:    0.03,
		},
		{
			InputTokens:  500,
			OutputTokens: 250,
			TotalTokens:  750,
			InputCost:    0.005,
			OutputCost:   0.0025,
			TotalCost:    0.0075,
		},
	}

	summary := calculator.CalculateSessionCost(requests)

	assert.Equal(t, int64(3500), summary.TotalInputTokens)
	assert.Equal(t, int64(1750), summary.TotalOutputTokens)
	assert.Equal(t, int64(5250), summary.TotalTokens)
	assert.InDelta(t, 0.035, summary.TotalInputCost, 0.0001)
	assert.InDelta(t, 0.0175, summary.TotalOutputCost, 0.0001)
	assert.InDelta(t, 0.0525, summary.TotalCost, 0.0001)
	assert.Equal(t, 3, summary.RequestCount)
	assert.InDelta(t, 0.0175, summary.AverageCostPerReq, 0.0001)
}

// TestCostCalculatorEmptySession tests empty session cost
func TestCostCalculatorEmptySession(t *testing.T) {
	pricingManager := NewPricingManager()
	calculator := NewCostCalculator(pricingManager)

	requests := []CostInfo{}
	summary := calculator.CalculateSessionCost(requests)

	assert.Equal(t, int64(0), summary.TotalInputTokens)
	assert.Equal(t, int64(0), summary.TotalOutputTokens)
	assert.Equal(t, int64(0), summary.TotalTokens)
	assert.Equal(t, 0.0, summary.TotalInputCost)
	assert.Equal(t, 0.0, summary.TotalOutputCost)
	assert.Equal(t, 0.0, summary.TotalCost)
	assert.Equal(t, 0, summary.RequestCount)
	assert.Equal(t, 0.0, summary.AverageCostPerReq)
}

// TestPricingManagerGetPricing tests pricing lookup
func TestPricingManagerGetPricing(t *testing.T) {
	manager := NewPricingManager()

	tests := []struct {
		modelID     string
		shouldExist bool
	}{
		{"gpt-4o", true},
		{"gpt-4o-mini", true},
		{"gpt-3.5-turbo", true},
		{"unknown-model", false},
	}

	for _, tt := range tests {
		t.Run(tt.modelID, func(t *testing.T) {
			pricing, exists := manager.GetPricing(tt.modelID)
			if tt.shouldExist {
				assert.True(t, exists, "model should exist")
				assert.NotEmpty(t, pricing.ModelID)
				assert.Greater(t, pricing.InputPricePer1K, 0.0)
				assert.Greater(t, pricing.OutputPricePer1K, 0.0)
			} else {
				assert.False(t, exists, "model should not exist")
			}
		})
	}
}

// TestPricingManagerAllModels tests that all expected models have pricing
func TestPricingManagerAllModels(t *testing.T) {
	manager := NewPricingManager()

	expectedModels := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-3.5-turbo",
		"gemini-1.5-pro",
		"gemini-1.5-flash",
	}

	for _, model := range expectedModels {
		t.Run(model, func(t *testing.T) {
			pricing, exists := manager.GetPricing(model)
			assert.True(t, exists, "model %s should exist in pricing", model)
			assert.NotEmpty(t, pricing.Provider)
		})
	}
}

// TestCostInfoStructure tests CostInfo structure
func TestCostInfoStructure(t *testing.T) {
	info := CostInfo{
		InputTokens:  100,
		OutputTokens: 50,
		TotalTokens:  150,
		InputCost:    0.001,
		OutputCost:   0.0005,
		TotalCost:    0.0015,
		Model:        "test-model",
	}

	assert.Equal(t, int64(100), info.InputTokens)
	assert.Equal(t, int64(50), info.OutputTokens)
	assert.Equal(t, int64(150), info.TotalTokens)
	assert.Equal(t, 0.001, info.InputCost)
	assert.Equal(t, 0.0005, info.OutputCost)
	assert.Equal(t, 0.0015, info.TotalCost)
	assert.Equal(t, "test-model", info.Model)
}

// TestBudgetConfigStructure tests BudgetConfig structure
func TestBudgetConfigStructure(t *testing.T) {
	config := BudgetConfig{
		DailyBudget:      10.0,
		MonthlyBudget:    200.0,
		WarningThreshold: 0.8,
	}

	assert.Equal(t, 10.0, config.DailyBudget)
	assert.Equal(t, 200.0, config.MonthlyBudget)
	assert.Equal(t, 0.8, config.WarningThreshold)
}

// TestBudgetStatusStructure tests BudgetStatus structure
func TestBudgetStatusStructure(t *testing.T) {
	status := BudgetStatus{
		DailySpent:       5.0,
		DailyRemaining:   5.0,
		DailyPercent:     50.0,
		MonthlySpent:     100.0,
		MonthlyRemaining: 100.0,
		MonthlyPercent:   50.0,
		IsOverBudget:     false,
		Warnings:         []string{},
	}

	assert.Equal(t, 5.0, status.DailySpent)
	assert.Equal(t, 5.0, status.DailyRemaining)
	assert.Equal(t, 50.0, status.DailyPercent)
	assert.False(t, status.IsOverBudget)
}

// TestCostOptimizationStructure tests CostOptimization structure
func TestCostOptimizationStructure(t *testing.T) {
	opt := CostOptimization{
		Type:            "model_switch",
		SavingsEstimate: 0.5,
		Description:     "Switch to a cheaper model",
		Priority:        1,
	}

	assert.Equal(t, "model_switch", opt.Type)
	assert.Equal(t, 0.5, opt.SavingsEstimate)
	assert.Equal(t, "Switch to a cheaper model", opt.Description)
	assert.Equal(t, 1, opt.Priority)
}
