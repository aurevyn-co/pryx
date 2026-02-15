package agent

import (
	"fmt"
	"strings"
)

// PreActionValidation performs validation before taking any action
type PreActionValidation struct {
	AvailableTools []string
}

// ValidationResult contains the outcome of pre-action validation
type ValidationResult struct {
	ShouldProceed bool
	Reason        string
	Confidence    ConfidenceLevel
}

// ConfidenceLevel represents how confident we are about an action
type ConfidenceLevel int

const (
	ConfidenceLow ConfidenceLevel = iota
	ConfidenceMedium
	ConfidenceHigh
)

func (c ConfidenceLevel) String() string {
	switch c {
	case ConfidenceHigh:
		return "HIGH"
	case ConfidenceMedium:
		return "MEDIUM"
	case ConfidenceLow:
		return "LOW"
	default:
		return "UNKNOWN"
	}
}

// NewPreActionValidation creates a new validator with available tools
func NewPreActionValidation(tools []string) *PreActionValidation {
	return &PreActionValidation{
		AvailableTools: tools,
	}
}

// ValidateToolUse checks if a tool should be used for a given request
func (p *PreActionValidation) ValidateToolUse(request string, toolName string) ValidationResult {
	// Check if tool exists
	if !p.isToolAvailable(toolName) {
		return ValidationResult{
			ShouldProceed: false,
			Reason:        fmt.Sprintf("Tool '%s' is not in available tools list", toolName),
			Confidence:    ConfidenceLow,
		}
	}

	// Check if tool is appropriate for the request
	if !p.isToolAppropriate(request, toolName) {
		return ValidationResult{
			ShouldProceed: false,
			Reason:        fmt.Sprintf("Tool '%s' does not appear appropriate for this request", toolName),
			Confidence:    ConfidenceLow,
		}
	}

	return ValidationResult{
		ShouldProceed: true,
		Reason:        fmt.Sprintf("Tool '%s' is available and appropriate", toolName),
		Confidence:    ConfidenceHigh,
	}
}

// isToolAvailable checks if a tool is in the available tools list
func (p *PreActionValidation) isToolAvailable(toolName string) bool {
	toolName = strings.ToLower(toolName)
	for _, tool := range p.AvailableTools {
		if strings.ToLower(tool) == toolName {
			return true
		}
	}
	return false
}

// isToolAppropriate uses heuristics to determine if a tool is appropriate
func (p *PreActionValidation) isToolAppropriate(request, toolName string) bool {
	request = strings.ToLower(request)
	toolName = strings.ToLower(toolName)

	switch toolName {
	case "filesystem", "file":
		fileKeywords := []string{"file", "read", "write", "directory", "folder", "path", "content", "open", "save"}
		return containsAny(request, fileKeywords)
	case "shell", "bash", "terminal":
		shellKeywords := []string{"run", "command", "execute", "script", "terminal", "shell", "bash", "sh "}
		return containsAny(request, shellKeywords)
	case "browser", "web", "fetch":
		browserKeywords := []string{"web", "url", "http", "website", "page", "browser", "internet", "search", "scrap"}
		return containsAny(request, browserKeywords)
	case "clipboard":
		clipboardKeywords := []string{"clipboard", "copy", "paste"}
		return containsAny(request, clipboardKeywords)
	default:
		return true
	}
}

// AssessConfidence evaluates confidence level for a given request
func (p *PreActionValidation) AssessConfidence(request string, hasContext bool) ConfidenceLevel {
	request = strings.ToLower(request)

	// High confidence indicators
	highConfidencePatterns := []string{
		"what is", "how to", "explain", "define",
		"list", "show me", "tell me", "what are",
	}

	// Low confidence indicators
	lowConfidencePatterns := []string{
		"fix", "debug", "solve", "implement",
		"create", "build", "design", "optimize",
	}

	// Check for context-dependent requests
	contextRequiredPatterns := []string{
		"this file", "the code", "here", "above",
		"previous", "last message", "context",
	}

	// If request requires context but none provided
	if !hasContext && containsAny(request, contextRequiredPatterns) {
		return ConfidenceLow
	}

	// If matches high confidence patterns
	if containsAny(request, highConfidencePatterns) {
		return ConfidenceHigh
	}

	// If matches low confidence patterns
	if containsAny(request, lowConfidencePatterns) {
		return ConfidenceLow
	}

	return ConfidenceMedium
}

// containsAny checks if the text contains any of the keywords
func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// ValidateRequest performs full validation on a user request
func (p *PreActionValidation) ValidateRequest(request string, proposedTool string, hasContext bool) ValidationResult {
	// First assess confidence
	confidence := p.AssessConfidence(request, hasContext)

	// If low confidence, suggest asking for clarification
	if confidence == ConfidenceLow {
		return ValidationResult{
			ShouldProceed: false,
			Reason:        "Low confidence in understanding the request - ask for clarification",
			Confidence:    confidence,
		}
	}

	// If a tool is proposed, validate it
	if proposedTool != "" {
		toolValidation := p.ValidateToolUse(request, proposedTool)
		if !toolValidation.ShouldProceed {
			return toolValidation
		}
		// Use the lower of the two confidence levels
		if confidence < toolValidation.Confidence {
			toolValidation.Confidence = confidence
		}
		return toolValidation
	}

	return ValidationResult{
		ShouldProceed: true,
		Reason:        "Request validated - no tool required",
		Confidence:    confidence,
	}
}
