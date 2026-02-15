package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultTimeout     = 30 * time.Second
	defaultAPIEndpoint = "https://api.telegram.org/bot"
)

// APIError represents an error returned by the Telegram Bot API
type APIError struct {
	Code        int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  *struct {
		MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
		RetryAfter      int   `json:"retry_after,omitempty"`
	} `json:"parameters,omitempty"`
}

func (e *APIError) Error() string {
	if e.Parameters != nil && e.Parameters.RetryAfter > 0 {
		return fmt.Sprintf("telegram API error %d: %s (retry after %d seconds)", e.Code, e.Description, e.Parameters.RetryAfter)
	}
	return fmt.Sprintf("telegram API error %d: %s", e.Code, e.Description)
}

// IsRetryable returns true if the error is retryable (rate limit, network error)
func (e *APIError) IsRetryable() bool {
	return e.Code == 429 || // Too Many Requests
		(e.Code >= 500 && e.Code < 600) // Server errors
}

// GetRetryAfter returns the number of seconds to wait before retrying
func (e *APIError) GetRetryAfter() int {
	if e.Parameters != nil {
		return e.Parameters.RetryAfter
	}
	return 0
}

// Response represents a generic response from the Telegram Bot API
type Response struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
	Parameters  *struct {
		MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
		RetryAfter      int   `json:"retry_after,omitempty"`
	} `json:"parameters,omitempty"`
}

// Client is a Telegram Bot API HTTP client
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
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

// NewClient creates a new Telegram Bot API client
func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		token:   token,
		baseURL: defaultAPIEndpoint + token + "/",
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// makeRequest performs an HTTP request to the Telegram API
func (c *Client) makeRequest(ctx context.Context, method string, params url.Values) (*Response, error) {
	apiURL := c.baseURL + method

	var body io.Reader
	if params != nil {
		body = bytes.NewBufferString(params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.OK {
		return nil, &APIError{
			Code:        apiResp.ErrorCode,
			Description: apiResp.Description,
			Parameters:  apiResp.Parameters,
		}
	}

	return &apiResp, nil
}

// makeMultipartRequest performs a multipart/form-data request (for file uploads)
func (c *Client) makeMultipartRequest(ctx context.Context, method string, params map[string]string, fileField string, file InputFile) (*Response, error) {
	apiURL := c.baseURL + method

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add params
	for key, value := range params {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	// Add file if provided
	if file.FileID == "" && (file.FilePath != "" || file.URL != "" || file.Data != nil) {
		var fileReader io.Reader
		var fileName string

		switch {
		case file.Data != nil:
			fileReader = bytes.NewReader(file.Data)
			fileName = file.FileName
			if fileName == "" {
				fileName = "file"
			}
		case file.FilePath != "":
			f, err := os.Open(file.FilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()
			fileReader = f
			fileName = filepath.Base(file.FilePath)
		case file.URL != "":
			// For URLs, Telegram downloads the file itself
			if err := writer.WriteField(fileField, file.URL); err != nil {
				return nil, fmt.Errorf("failed to write URL field: %w", err)
			}
		}

		if fileReader != nil {
			part, err := writer.CreateFormFile(fileField, fileName)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := io.Copy(part, fileReader); err != nil {
				return nil, fmt.Errorf("failed to copy file data: %w", err)
			}
		}
	} else if file.FileID != "" {
		// Existing file on Telegram servers
		if err := writer.WriteField(fileField, file.FileID); err != nil {
			return nil, fmt.Errorf("failed to write file ID: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.OK {
		return nil, &APIError{
			Code:        apiResp.ErrorCode,
			Description: apiResp.Description,
			Parameters:  apiResp.Parameters,
		}
	}

	return &apiResp, nil
}

// GetMe returns basic information about the bot
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	resp, err := c.makeRequest(ctx, "getMe", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp.Result, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetUpdates gets pending updates from Telegram
// offset: Identifier of the first update to be returned
// limit: Limits the number of updates to be retrieved (1-100)
// timeout: Timeout in seconds for long polling
// allowedUpdates: List of update types to receive
func (c *Client) GetUpdates(ctx context.Context, offset, limit, timeout int, allowedUpdates []string) ([]Update, error) {
	params := url.Values{}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if timeout > 0 {
		params.Set("timeout", strconv.Itoa(timeout))
	}
	if len(allowedUpdates) > 0 {
		allowedUpdatesJSON, _ := json.Marshal(allowedUpdates)
		params.Set("allowed_updates", string(allowedUpdatesJSON))
	}

	resp, err := c.makeRequest(ctx, "getUpdates", params)
	if err != nil {
		return nil, err
	}

	var updates []Update
	if err := json.Unmarshal(resp.Result, &updates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates: %w", err)
	}

	return updates, nil
}

// SendMessage sends a text message to a chat
func (c *Client) SendMessage(ctx context.Context, chatID int64, text string, opts ...SendMessageOption) (*Message, error) {
	params := url.Values{
		"chat_id": {strconv.FormatInt(chatID, 10)},
		"text":    {text},
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeRequest(ctx, "sendMessage", params)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(resp.Result, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// SendMessageOption is a functional option for SendMessage
type SendMessageOption func(url.Values)

// WithParseMode sets the parse mode for the message
func WithParseMode(mode ParseMode) SendMessageOption {
	return func(v url.Values) {
		v.Set("parse_mode", string(mode))
	}
}

// WithDisableWebPagePreview disables link previews
func WithDisableWebPagePreview(disable bool) SendMessageOption {
	return func(v url.Values) {
		v.Set("disable_web_page_preview", strconv.FormatBool(disable))
	}
}

// WithDisableNotification sends message silently
func WithDisableNotification(disable bool) SendMessageOption {
	return func(v url.Values) {
		v.Set("disable_notification", strconv.FormatBool(disable))
	}
}

// WithReplyToMessageID sets the message to reply to
func WithReplyToMessageID(messageID int) SendMessageOption {
	return func(v url.Values) {
		v.Set("reply_to_message_id", strconv.Itoa(messageID))
	}
}

// WithReplyMarkup sets the reply markup
func WithReplyMarkup(markup interface{}) SendMessageOption {
	return func(v url.Values) {
		markupJSON, _ := json.Marshal(markup)
		v.Set("reply_markup", string(markupJSON))
	}
}

// SendPhoto sends a photo to a chat
// photo can be a file_id, URL, or local file path
func (c *Client) SendPhoto(ctx context.Context, chatID int64, photo InputFile, caption string, opts ...SendPhotoOption) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeMultipartRequest(ctx, "sendPhoto", params, "photo", photo)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(resp.Result, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// SendPhotoOption is a functional option for SendPhoto
type SendPhotoOption func(map[string]string)

// WithPhotoParseMode sets the parse mode for the caption
func WithPhotoParseMode(mode ParseMode) SendPhotoOption {
	return func(m map[string]string) {
		m["parse_mode"] = string(mode)
	}
}

// WithPhotoDisableNotification sends photo silently
func WithPhotoDisableNotification(disable bool) SendPhotoOption {
	return func(m map[string]string) {
		m["disable_notification"] = strconv.FormatBool(disable)
	}
}

// WithPhotoReplyToMessageID sets the message to reply to
func WithPhotoReplyToMessageID(messageID int) SendPhotoOption {
	return func(m map[string]string) {
		m["reply_to_message_id"] = strconv.Itoa(messageID)
	}
}

// SendDocument sends a document to a chat
func (c *Client) SendDocument(ctx context.Context, chatID int64, document InputFile, caption string, opts ...SendDocumentOption) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeMultipartRequest(ctx, "sendDocument", params, "document", document)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(resp.Result, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// SendDocumentOption is a functional option for SendDocument
type SendDocumentOption func(map[string]string)

// WithDocumentParseMode sets the parse mode for the caption
func WithDocumentParseMode(mode ParseMode) SendDocumentOption {
	return func(m map[string]string) {
		m["parse_mode"] = string(mode)
	}
}

// WithDocumentDisableNotification sends document silently
func WithDocumentDisableNotification(disable bool) SendDocumentOption {
	return func(m map[string]string) {
		m["disable_notification"] = strconv.FormatBool(disable)
	}
}

// WithDocumentReplyToMessageID sets the message to reply to
func WithDocumentReplyToMessageID(messageID int) SendDocumentOption {
	return func(m map[string]string) {
		m["reply_to_message_id"] = strconv.Itoa(messageID)
	}
}

// WithThumbnail sets the thumbnail for the document
func WithThumbnail(thumbnail InputFile) SendDocumentOption {
	return func(m map[string]string) {
		// Note: Thumbnail upload would require additional handling
		// For now, we only support file_id or URL
		if thumbnail.FileID != "" {
			m["thumbnail"] = thumbnail.FileID
		} else if thumbnail.URL != "" {
			m["thumbnail"] = thumbnail.URL
		}
	}
}

// SetWebhook sets the webhook URL for receiving updates
func (c *Client) SetWebhook(ctx context.Context, webhookURL string, opts ...SetWebhookOption) (bool, error) {
	params := map[string]string{
		"url": webhookURL,
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeMultipartRequest(ctx, "setWebhook", params, "certificate", InputFile{})
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

// SetWebhookOption is a functional option for SetWebhook
type SetWebhookOption func(map[string]string)

// WithMaxConnections sets the maximum allowed number of simultaneous HTTPS connections
func WithMaxConnections(max int) SetWebhookOption {
	return func(m map[string]string) {
		m["max_connections"] = strconv.Itoa(max)
	}
}

// WithAllowedUpdates sets the list of update types to receive
func WithAllowedUpdates(updates []string) SetWebhookOption {
	return func(m map[string]string) {
		updatesJSON, _ := json.Marshal(updates)
		m["allowed_updates"] = string(updatesJSON)
	}
}

// WithDropPendingUpdates drops all pending updates
func WithDropPendingUpdates(drop bool) SetWebhookOption {
	return func(m map[string]string) {
		m["drop_pending_updates"] = strconv.FormatBool(drop)
	}
}

// WithCertificate sets the certificate for the webhook
func WithCertificate(cert InputFile) SetWebhookOption {
	return func(m map[string]string) {
		// This would require special handling in makeMultipartRequest
		// For now, it's a placeholder
	}
}

// DeleteWebhook removes the webhook integration
func (c *Client) DeleteWebhook(ctx context.Context, dropPendingUpdates bool) (bool, error) {
	params := url.Values{
		"drop_pending_updates": {strconv.FormatBool(dropPendingUpdates)},
	}

	resp, err := c.makeRequest(ctx, "deleteWebhook", params)
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

// GetWebhookInfo returns current webhook status
func (c *Client) GetWebhookInfo(ctx context.Context) (*WebhookInfo, error) {
	resp, err := c.makeRequest(ctx, "getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	var info WebhookInfo
	if err := json.Unmarshal(resp.Result, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook info: %w", err)
	}

	return &info, nil
}

// SendChatAction sends a chat action (typing, uploading_photo, etc.)
func (c *Client) SendChatAction(ctx context.Context, chatID int64, action ChatAction) (bool, error) {
	params := url.Values{
		"chat_id": {strconv.FormatInt(chatID, 10)},
		"action":  {string(action)},
	}

	resp, err := c.makeRequest(ctx, "sendChatAction", params)
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

// GetFile gets basic info about a file and prepare it for downloading
func (c *Client) GetFile(ctx context.Context, fileID string) (*File, error) {
	params := url.Values{
		"file_id": {fileID},
	}

	resp, err := c.makeRequest(ctx, "getFile", params)
	if err != nil {
		return nil, err
	}

	var file File
	if err := json.Unmarshal(resp.Result, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %w", err)
	}

	return &file, nil
}

// File represents a file ready to be downloaded
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// GetFileURL returns the full URL to download a file
func (c *Client) GetFileURL(file *File) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", c.token, file.FilePath)
}

// ForwardMessage forwards a message
func (c *Client) ForwardMessage(ctx context.Context, chatID, fromChatID int64, messageID int, opts ...ForwardMessageOption) (*Message, error) {
	params := url.Values{
		"chat_id":      {strconv.FormatInt(chatID, 10)},
		"from_chat_id": {strconv.FormatInt(fromChatID, 10)},
		"message_id":   {strconv.Itoa(messageID)},
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeRequest(ctx, "forwardMessage", params)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(resp.Result, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// ForwardMessageOption is a functional option for ForwardMessage
type ForwardMessageOption func(url.Values)

// WithDisableNotificationForward forwards message silently
func WithDisableNotificationForward(disable bool) ForwardMessageOption {
	return func(v url.Values) {
		v.Set("disable_notification", strconv.FormatBool(disable))
	}
}

// WithProtectContent protects the forwarded message from forwarding and saving
func WithProtectContent(protect bool) ForwardMessageOption {
	return func(v url.Values) {
		v.Set("protect_content", strconv.FormatBool(protect))
	}
}

// DeleteMessage deletes a message
func (c *Client) DeleteMessage(ctx context.Context, chatID int64, messageID int) (bool, error) {
	params := url.Values{
		"chat_id":    {strconv.FormatInt(chatID, 10)},
		"message_id": {strconv.Itoa(messageID)},
	}

	resp, err := c.makeRequest(ctx, "deleteMessage", params)
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

// EditMessageText edits text and game messages
func (c *Client) EditMessageText(ctx context.Context, chatID int64, messageID int, text string, opts ...EditMessageOption) (*Message, error) {
	params := url.Values{
		"chat_id":    {strconv.FormatInt(chatID, 10)},
		"message_id": {strconv.Itoa(messageID)},
		"text":       {text},
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeRequest(ctx, "editMessageText", params)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(resp.Result, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// EditMessageOption is a functional option for EditMessageText
type EditMessageOption func(url.Values)

// WithEditParseMode sets the parse mode for the edited message
func WithEditParseMode(mode ParseMode) EditMessageOption {
	return func(v url.Values) {
		v.Set("parse_mode", string(mode))
	}
}

// WithDisableWebPagePreviewEdit disables link previews in edited message
func WithDisableWebPagePreviewEdit(disable bool) EditMessageOption {
	return func(v url.Values) {
		v.Set("disable_web_page_preview", strconv.FormatBool(disable))
	}
}

// WithEditReplyMarkup sets the reply markup for the edited message
func WithEditReplyMarkup(markup interface{}) EditMessageOption {
	return func(v url.Values) {
		markupJSON, _ := json.Marshal(markup)
		v.Set("reply_markup", string(markupJSON))
	}
}

// SetMyCommands changes the list of the bot's commands
func (c *Client) SetMyCommands(ctx context.Context, commands []BotCommand, opts ...SetMyCommandsOption) (bool, error) {
	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		return false, fmt.Errorf("failed to marshal commands: %w", err)
	}

	params := url.Values{
		"commands": {string(commandsJSON)},
	}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeRequest(ctx, "setMyCommands", params)
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

// SetMyCommandsOption is a functional option for SetMyCommands
type SetMyCommandsOption func(url.Values)

// WithScope sets the scope of users for which the commands are relevant
func WithScope(scope interface{}) SetMyCommandsOption {
	return func(v url.Values) {
		scopeJSON, _ := json.Marshal(scope)
		v.Set("scope", string(scopeJSON))
	}
}

// WithLanguageCode sets the language code for the commands
func WithLanguageCode(code string) SetMyCommandsOption {
	return func(v url.Values) {
		v.Set("language_code", code)
	}
}

// GetMyCommands gets the current list of the bot's commands
func (c *Client) GetMyCommands(ctx context.Context, opts ...GetMyCommandsOption) ([]BotCommand, error) {
	params := url.Values{}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeRequest(ctx, "getMyCommands", params)
	if err != nil {
		return nil, err
	}

	var commands []BotCommand
	if err := json.Unmarshal(resp.Result, &commands); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
	}

	return commands, nil
}

// GetMyCommandsOption is a functional option for GetMyCommands
type GetMyCommandsOption func(url.Values)

// GetMyCommandsWithScope sets the scope for getting commands
func GetMyCommandsWithScope(scope interface{}) GetMyCommandsOption {
	return func(v url.Values) {
		scopeJSON, _ := json.Marshal(scope)
		v.Set("scope", string(scopeJSON))
	}
}

// GetMyCommandsWithLanguageCode sets the language code for getting commands
func GetMyCommandsWithLanguageCode(code string) GetMyCommandsOption {
	return func(v url.Values) {
		v.Set("language_code", code)
	}
}

// DeleteMyCommands deletes the list of the bot's commands
func (c *Client) DeleteMyCommands(ctx context.Context, opts ...DeleteMyCommandsOption) (bool, error) {
	params := url.Values{}

	for _, opt := range opts {
		opt(params)
	}

	resp, err := c.makeRequest(ctx, "deleteMyCommands", params)
	if err != nil {
		return false, err
	}

	var result bool
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

// DeleteMyCommandsOption is a functional option for DeleteMyCommands
type DeleteMyCommandsOption func(url.Values)

// DeleteMyCommandsWithScope sets the scope for deleting commands
func DeleteMyCommandsWithScope(scope interface{}) DeleteMyCommandsOption {
	return func(v url.Values) {
		scopeJSON, _ := json.Marshal(scope)
		v.Set("scope", string(scopeJSON))
	}
}

// DeleteMyCommandsWithLanguageCode sets the language code for deleting commands
func DeleteMyCommandsWithLanguageCode(code string) DeleteMyCommandsOption {
	return func(v url.Values) {
		v.Set("language_code", code)
	}
}

// ValidateToken validates the bot token by calling getMe
func (c *Client) ValidateToken(ctx context.Context) (*User, error) {
	return c.GetMe(ctx)
}
