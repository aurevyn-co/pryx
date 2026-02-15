// Package llm provides LLM provider integration for Pryx.
// It supports multiple providers including OpenAI, Anthropic, and OpenRouter.
package llm

// Role represents the role of a message sender in a conversation.
type Role string

// Message roles for chat conversations.
const (
	// RoleUser represents a message from the user.
	RoleUser Role = "user"
	// RoleAssistant represents a message from the AI assistant.
	RoleAssistant Role = "assistant"
	// RoleSystem represents a system message that sets context/behavior.
	RoleSystem Role = "system"
)

// Message represents a single message in a chat conversation.
type Message struct {
	// Role is the sender's role (user, assistant, or system).
	Role Role `json:"role"`
	// Content is the message text content.
	Content string `json:"content"`
}

// ChatRequest represents a request to an LLM for chat completion.
type ChatRequest struct {
	// Model is the model identifier (e.g., "gpt-4", "claude-3-opus").
	Model string `json:"model"`
	// Messages is the conversation history.
	Messages []Message `json:"messages"`
	// MaxTokens limits the number of tokens in the response.
	MaxTokens int `json:"max_tokens,omitempty"`
	// Temperature controls randomness (0.0-2.0, lower is more deterministic).
	Temperature float64 `json:"temperature,omitempty"`
	// Stream indicates whether to stream the response.
	Stream bool `json:"stream,omitempty"`
}

// ChatResponse represents a response from an LLM chat completion.
type ChatResponse struct {
	// Content is the generated response text.
	Content string `json:"content"`
	// Role is always RoleAssistant for responses.
	Role Role `json:"role"`
	// FinishReason indicates why the generation stopped (e.g., "stop", "length").
	FinishReason string `json:"finish_reason,omitempty"`
	// Usage contains token count information.
	Usage Usage `json:"usage"`
}

// Usage contains token usage statistics for an LLM request.
type Usage struct {
	// PromptTokens is the number of tokens in the input messages.
	PromptTokens int `json:"prompt_tokens"`
	// CompletionTokens is the number of tokens in the generated response.
	CompletionTokens int `json:"completion_tokens"`
	// TotalTokens is the sum of prompt and completion tokens.
	TotalTokens int `json:"total_tokens"`
}

// StreamChunk represents a single chunk in a streaming response.
type StreamChunk struct {
	// Content is the incremental text content (delta).
	Content string `json:"content"`
	// Done indicates if this is the final chunk.
	Done bool `json:"done"`
	// Err contains any error that occurred during streaming (not serialized).
	Err error `json:"-"`
}
