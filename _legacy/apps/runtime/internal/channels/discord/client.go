package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	apiBaseURL     = "https://discord.com/api/v10"
	defaultTimeout = 30 * time.Second
)

// APIError represents an error returned by the Discord API
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after,omitempty"`
	Global     bool   `json:"global,omitempty"`
}

func (e *APIError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("discord API error %d: %s (retry after %d seconds)", e.Code, e.Message, e.RetryAfter)
	}
	return fmt.Sprintf("discord API error %d: %s", e.Code, e.Message)
}

// IsRateLimit returns true if the error is a rate limit
func (e *APIError) IsRateLimit() bool {
	return e.Code == 429 || e.RetryAfter > 0
}

// GetRetryAfter returns the number of seconds to wait before retrying
func (e *APIError) GetRetryAfter() time.Duration {
	return time.Duration(e.RetryAfter) * time.Second
}

// Client is a Discord REST API client
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client

	// Rate limiting
	rateLimits map[string]*RateLimitInfo
	rateMu     sync.RWMutex
}

// RateLimitInfo stores rate limit information for an endpoint
type RateLimitInfo struct {
	Limit      int
	Remaining  int
	Reset      time.Time
	ResetAfter time.Duration
}

// ClientOption is a functional option for configuring the Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new Discord REST API client
func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		token:      token,
		baseURL:    apiBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		rateLimits: make(map[string]*RateLimitInfo),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// makeRequest performs an HTTP request to the Discord API
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	// Check rate limits
	if err := c.checkRateLimit(endpoint); err != nil {
		return nil, err
	}

	apiURL := c.baseURL + endpoint

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "PryxBot (https://github.com/aurevyn-co/pryx, 1.0)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Update rate limit info
	c.updateRateLimit(endpoint, resp)

	return resp, nil
}

// checkRateLimit checks if we should wait before making a request
func (c *Client) checkRateLimit(endpoint string) error {
	c.rateMu.RLock()
	info, exists := c.rateLimits[endpoint]
	c.rateMu.RUnlock()

	if !exists {
		return nil
	}

	if info.Remaining <= 0 && time.Now().Before(info.Reset) {
		waitTime := time.Until(info.Reset)
		if waitTime > 0 {
			return fmt.Errorf("rate limited: wait %v", waitTime)
		}
	}

	return nil
}

// updateRateLimit updates rate limit info from response headers
func (c *Client) updateRateLimit(endpoint string, resp *http.Response) {
	limit := resp.Header.Get("X-RateLimit-Limit")
	remaining := resp.Header.Get("X-RateLimit-Remaining")
	reset := resp.Header.Get("X-RateLimit-Reset")
	resetAfter := resp.Header.Get("X-RateLimit-Reset-After")

	if limit == "" && remaining == "" {
		return
	}

	info := &RateLimitInfo{}

	if limit != "" {
		info.Limit, _ = strconv.Atoi(limit)
	}
	if remaining != "" {
		info.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset != "" {
		if resetUnix, err := strconv.ParseInt(reset, 10, 64); err == nil {
			info.Reset = time.Unix(resetUnix, 0)
		}
	}
	if resetAfter != "" {
		if after, err := strconv.ParseFloat(resetAfter, 64); err == nil {
			info.ResetAfter = time.Duration(after * float64(time.Second))
		}
	}

	c.rateMu.Lock()
	c.rateLimits[endpoint] = info
	c.rateMu.Unlock()
}

// handleResponse processes the API response
func (c *Client) handleResponse(resp *http.Response) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Handle rate limiting
	if resp.StatusCode == 429 {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			if after, err := strconv.Atoi(retryAfter); err == nil {
				apiErr.RetryAfter = after
			}
		}
		apiErr.Global = resp.Header.Get("X-RateLimit-Global") == "true"
	}

	return &apiErr
}

// GetMe returns the current bot user
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/users/@me", nil)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

// GetUser returns a user by ID
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/users/"+userID, nil)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

// GetChannel returns a channel by ID
func (c *Client) GetChannel(ctx context.Context, channelID string) (*Channel, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/channels/"+channelID, nil)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var channel Channel
	if err := json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, fmt.Errorf("failed to decode channel: %w", err)
	}

	return &channel, nil
}

// SendMessage sends a text message to a channel
func (c *Client) SendMessage(ctx context.Context, channelID, content string) (*Message, error) {
	payload := map[string]interface{}{
		"content": content,
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, "/channels/"+channelID+"/messages", payload)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// SendMessageComplex sends a message with full options
func (c *Client) SendMessageComplex(ctx context.Context, channelID string, data *MessageSend) (*Message, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/channels/"+channelID+"/messages", data)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// SendEmbed sends an embed message to a channel
func (c *Client) SendEmbed(ctx context.Context, channelID string, embed *Embed) (*Message, error) {
	payload := map[string]interface{}{
		"embeds": []Embed{*embed},
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, "/channels/"+channelID+"/messages", payload)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// SendFile sends a file to a channel
func (c *Client) SendFile(ctx context.Context, channelID, filename string, data []byte, content string) (*Message, error) {
	apiURL := c.baseURL + "/channels/" + channelID + "/messages"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// Add content if provided
	if content != "" {
		if err := writer.WriteField("content", content); err != nil {
			return nil, fmt.Errorf("failed to write content field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "PryxBot (https://github.com/aurevyn-co/pryx, 1.0)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// EditMessage edits an existing message
func (c *Client) EditMessage(ctx context.Context, channelID, messageID string, content string, embeds []Embed) (*Message, error) {
	payload := map[string]interface{}{}
	if content != "" {
		payload["content"] = content
	}
	if len(embeds) > 0 {
		payload["embeds"] = embeds
	}

	resp, err := c.makeRequest(ctx, http.MethodPatch, "/channels/"+channelID+"/messages/"+messageID, payload)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// DeleteMessage deletes a message
func (c *Client) DeleteMessage(ctx context.Context, channelID, messageID string) error {
	resp, err := c.makeRequest(ctx, http.MethodDelete, "/channels/"+channelID+"/messages/"+messageID, nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

// CreateSlashCommand creates a global slash command
func (c *Client) CreateSlashCommand(ctx context.Context, appID string, cmd *ApplicationCommand) (*ApplicationCommand, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/applications/"+appID+"/commands", cmd)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var createdCmd ApplicationCommand
	if err := json.NewDecoder(resp.Body).Decode(&createdCmd); err != nil {
		return nil, fmt.Errorf("failed to decode command: %w", err)
	}

	return &createdCmd, nil
}

// CreateGuildSlashCommand creates a guild-specific slash command
func (c *Client) CreateGuildSlashCommand(ctx context.Context, appID, guildID string, cmd *ApplicationCommand) (*ApplicationCommand, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/applications/"+appID+"/guilds/"+guildID+"/commands", cmd)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var createdCmd ApplicationCommand
	if err := json.NewDecoder(resp.Body).Decode(&createdCmd); err != nil {
		return nil, fmt.Errorf("failed to decode command: %w", err)
	}

	return &createdCmd, nil
}

// GetSlashCommands returns all global slash commands for an application
func (c *Client) GetSlashCommands(ctx context.Context, appID string) ([]ApplicationCommand, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/applications/"+appID+"/commands", nil)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var cmds []ApplicationCommand
	if err := json.NewDecoder(resp.Body).Decode(&cmds); err != nil {
		return nil, fmt.Errorf("failed to decode commands: %w", err)
	}

	return cmds, nil
}

// DeleteSlashCommand deletes a global slash command
func (c *Client) DeleteSlashCommand(ctx context.Context, appID, cmdID string) error {
	resp, err := c.makeRequest(ctx, http.MethodDelete, "/applications/"+appID+"/commands/"+cmdID, nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

// DeleteGuildSlashCommand deletes a guild-specific slash command
func (c *Client) DeleteGuildSlashCommand(ctx context.Context, appID, guildID, cmdID string) error {
	resp, err := c.makeRequest(ctx, http.MethodDelete, "/applications/"+appID+"/guilds/"+guildID+"/commands/"+cmdID, nil)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

// RespondToInteraction responds to an interaction (slash command)
func (c *Client) RespondToInteraction(ctx context.Context, interactionID, token string, response *InteractionResponse) error {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/interactions/"+interactionID+"/"+token+"/callback", response)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

// EditInteractionResponse edits the original response to an interaction
func (c *Client) EditInteractionResponse(ctx context.Context, appID, token string, response *InteractionResponseData) (*Message, error) {
	resp, err := c.makeRequest(ctx, http.MethodPatch, "/webhooks/"+appID+"/"+token+"/messages/@original", response)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// FollowupInteraction sends a followup message to an interaction
func (c *Client) FollowupInteraction(ctx context.Context, appID, token string, data *InteractionResponseData) (*Message, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/webhooks/"+appID+"/"+token, data)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return &msg, nil
}

// BulkOverwriteCommands overwrites all global commands for an application
func (c *Client) BulkOverwriteCommands(ctx context.Context, appID string, cmds []ApplicationCommand) ([]ApplicationCommand, error) {
	resp, err := c.makeRequest(ctx, http.MethodPut, "/applications/"+appID+"/commands", cmds)
	if err != nil {
		return nil, err
	}

	if err := c.handleResponse(resp); err != nil {
		return nil, err
	}

	var createdCmds []ApplicationCommand
	if err := json.NewDecoder(resp.Body).Decode(&createdCmds); err != nil {
		return nil, fmt.Errorf("failed to decode commands: %w", err)
	}

	return createdCmds, nil
}

// MessageSend represents data for sending a message
type MessageSend struct {
	Content         string             `json:"content,omitempty"`
	Embeds          []Embed            `json:"embeds,omitempty"`
	TTS             bool               `json:"tts,omitempty"`
	AllowedMentions *AllowedMentions   `json:"allowed_mentions,omitempty"`
	Components      []MessageComponent `json:"components,omitempty"`
	StickerIDs      []string           `json:"sticker_ids,omitempty"`
	Flags           int                `json:"flags,omitempty"`
}
