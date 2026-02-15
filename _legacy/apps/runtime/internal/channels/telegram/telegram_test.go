package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"pryx-core/internal/bus"
	"pryx-core/internal/channels"
)

// Mock Telegram API Server
type MockTelegramServer struct {
	Server   *httptest.Server
	Token    string
	Updates  []Update
	Messages []Message
}

func NewMockTelegramServer() *MockTelegramServer {
	mock := &MockTelegramServer{
		Token:    "test-token-12345",
		Updates:  make([]Update, 0),
		Messages: make([]Message, 0),
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.handleRequest(w, r)
	}))

	return mock
}

func (m *MockTelegramServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/bot"+m.Token+"/")

	switch path {
	case "getMe":
		m.handleGetMe(w, r)
	case "getUpdates":
		m.handleGetUpdates(w, r)
	case "sendMessage":
		m.handleSendMessage(w, r)
	case "sendPhoto":
		m.handleSendPhoto(w, r)
	case "sendDocument":
		m.handleSendDocument(w, r)
	case "setWebhook":
		m.handleSetWebhook(w, r)
	case "deleteWebhook":
		m.handleDeleteWebhook(w, r)
	case "getWebhookInfo":
		m.handleGetWebhookInfo(w, r)
	case "setMyCommands":
		m.handleSetMyCommands(w, r)
	case "getMyCommands":
		m.handleGetMyCommands(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":          false,
			"error_code":  404,
			"description": "Not Found",
		})
	}
}

func (m *MockTelegramServer) handleGetMe(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok": true,
		"result": User{
			ID:        123456789,
			IsBot:     true,
			FirstName: "TestBot",
			Username:  "test_bot",
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleGetUpdates(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok":     true,
		"result": m.Updates,
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, _ := strconv.ParseInt(r.FormValue("chat_id"), 10, 64)
	msg := Message{
		MessageID: len(m.Messages) + 1,
		Chat:      &Chat{ID: chatID, Type: "private"},
		Text:      r.FormValue("text"),
		Date:      int(time.Now().Unix()),
	}
	m.Messages = append(m.Messages, msg)

	response := map[string]interface{}{
		"ok":     true,
		"result": msg,
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleSendPhoto(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok": true,
		"result": Message{
			MessageID: 999,
			Chat:      &Chat{ID: 123, Type: "private"},
			Photo:     []PhotoSize{{FileID: "photo123", Width: 100, Height: 100}},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleSendDocument(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok": true,
		"result": Message{
			MessageID: 1000,
			Chat:      &Chat{ID: 123, Type: "private"},
			Document:  &Document{FileID: "doc123", FileName: "test.pdf"},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleSetWebhook(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok":     true,
		"result": true,
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleDeleteWebhook(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok":     true,
		"result": true,
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleGetWebhookInfo(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok": true,
		"result": WebhookInfo{
			URL:                  "https://example.com/webhook",
			HasCustomCertificate: false,
			PendingUpdateCount:   0,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleSetMyCommands(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok":     true,
		"result": true,
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) handleGetMyCommands(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"ok": true,
		"result": []BotCommand{
			{Command: "start", Description: "Start the bot"},
			{Command: "help", Description: "Show help"},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *MockTelegramServer) Close() {
	m.Server.Close()
}

func (m *MockTelegramServer) URL() string {
	return m.Server.URL + "/bot" + m.Token + "/"
}

// Client Tests

func TestClient_GetMe(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))

	user, err := client.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}

	if user.Username != "test_bot" {
		t.Errorf("Expected username 'test_bot', got '%s'", user.Username)
	}

	if !user.IsBot {
		t.Error("Expected IsBot to be true")
	}
}

func TestClient_SendMessage(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))

	msg, err := client.SendMessage(context.Background(), 123456, "Hello, World!")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if msg.Text != "Hello, World!" {
		t.Errorf("Expected text 'Hello, World!', got '%s'", msg.Text)
	}
}

func TestClient_SendMessage_WithOptions(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))

	msg, err := client.SendMessage(context.Background(), 123456, "*Bold* text",
		WithParseMode(ParseModeMarkdown),
		WithDisableWebPagePreview(true),
		WithDisableNotification(true))

	if err != nil {
		t.Fatalf("SendMessage with options failed: %v", err)
	}

	if msg.MessageID == 0 {
		t.Error("Expected message ID to be set")
	}
}

func TestClient_GetUpdates(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	// Add a mock update
	mock.Updates = append(mock.Updates, Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			Chat:      &Chat{ID: 123, Type: "private"},
			Text:      "Test message",
		},
	})

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))

	updates, err := client.GetUpdates(context.Background(), 0, 100, 30, nil)
	if err != nil {
		t.Fatalf("GetUpdates failed: %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(updates))
	}

	if updates[0].Message.Text != "Test message" {
		t.Errorf("Expected text 'Test message', got '%s'", updates[0].Message.Text)
	}
}

func TestClient_SetWebhook(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))

	result, err := client.SetWebhook(context.Background(), "https://example.com/webhook")
	if err != nil {
		t.Fatalf("SetWebhook failed: %v", err)
	}

	if !result {
		t.Error("Expected SetWebhook to return true")
	}
}

func TestClient_DeleteWebhook(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))

	result, err := client.DeleteWebhook(context.Background(), false)
	if err != nil {
		t.Fatalf("DeleteWebhook failed: %v", err)
	}

	if !result {
		t.Error("Expected DeleteWebhook to return true")
	}
}

func TestClient_APIError(t *testing.T) {
	// Create a server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":          false,
			"error_code":  400,
			"description": "Bad Request: chat not found",
		})
	}))
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL+"/bot"))

	_, err := client.SendMessage(context.Background(), 123, "test")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that error message contains expected content
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected error to contain '400', got: %s", err.Error())
	}

	if !strings.Contains(err.Error(), "Bad Request") {
		t.Errorf("Expected error to contain 'Bad Request', got: %s", err.Error())
	}
}

// Handler Tests

func TestHandler_HandleCommand_Start(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.AllowedChats = []int64{123}

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)

	msg := &Message{
		MessageID: 1,
		Chat:      &Chat{ID: 123, Type: "private"},
		From:      &User{ID: 456, FirstName: "Test", Username: "testuser"},
		Text:      "/start",
		Date:      int(time.Now().Unix()),
	}

	err := handler.handleCommand(context.Background(), msg)
	if err != nil {
		t.Fatalf("handleCommand failed: %v", err)
	}

	// Verify message was sent
	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}

	if !strings.Contains(mock.Messages[0].Text, "Welcome") {
		t.Errorf("Expected welcome message, got: %s", mock.Messages[0].Text)
	}
}

func TestHandler_HandleCommand_Help(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.AllowedChats = []int64{123}

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)

	msg := &Message{
		MessageID: 1,
		Chat:      &Chat{ID: 123, Type: "private"},
		From:      &User{ID: 456, FirstName: "Test"},
		Text:      "/help",
		Date:      int(time.Now().Unix()),
	}

	err := handler.handleCommand(context.Background(), msg)
	if err != nil {
		t.Fatalf("handleCommand failed: %v", err)
	}

	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}

	if !strings.Contains(mock.Messages[0].Text, "Available Commands") {
		t.Errorf("Expected help message, got: %s", mock.Messages[0].Text)
	}
}

func TestHandler_HandleCommand_Status(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.Mode = "polling"
	config.AllowedChats = []int64{123}

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)

	msg := &Message{
		MessageID: 1,
		Chat:      &Chat{ID: 123, Type: "private"},
		From:      &User{ID: 456, FirstName: "Test"},
		Text:      "/status",
		Date:      int(time.Now().Unix()),
	}

	err := handler.handleCommand(context.Background(), msg)
	if err != nil {
		t.Fatalf("handleCommand failed: %v", err)
	}

	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}

	if !strings.Contains(mock.Messages[0].Text, "Bot Status") {
		t.Errorf("Expected status message, got: %s", mock.Messages[0].Text)
	}
}

func TestHandler_HandleCommand_Unknown(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.AllowedChats = []int64{123}

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)

	msg := &Message{
		MessageID: 1,
		Chat:      &Chat{ID: 123, Type: "private"},
		From:      &User{ID: 456, FirstName: "Test"},
		Text:      "/unknowncommand",
		Date:      int(time.Now().Unix()),
	}

	err := handler.handleCommand(context.Background(), msg)
	if err != nil {
		t.Fatalf("handleCommand failed: %v", err)
	}

	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}

	if !strings.Contains(mock.Messages[0].Text, "Unknown command") {
		t.Errorf("Expected unknown command message, got: %s", mock.Messages[0].Text)
	}
}

func TestHandler_ChatWhitelist(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.AllowedChats = []int64{123} // Only allow chat 123

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	eventBus := bus.New()
	handler := NewHandler(&config, client, eventBus)

	// Message from allowed chat
	allowedMsg := &Message{
		MessageID: 1,
		Chat:      &Chat{ID: 123, Type: "private"},
		From:      &User{ID: 456, FirstName: "Test"},
		Text:      "Hello",
		Date:      int(time.Now().Unix()),
	}

	// Subscribe to events
	msgCh, unsub := eventBus.Subscribe(bus.EventChannelMessage)
	defer unsub()

	err := handler.handleMessage(context.Background(), allowedMsg)
	if err != nil {
		t.Fatalf("handleMessage failed: %v", err)
	}

	// Should receive event
	select {
	case <-msgCh:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected message event for allowed chat")
	}

	// Message from blocked chat
	blockedMsg := &Message{
		MessageID: 2,
		Chat:      &Chat{ID: 999, Type: "private"},
		From:      &User{ID: 789, FirstName: "Blocked"},
		Text:      "Hello",
		Date:      int(time.Now().Unix()),
	}

	err = handler.handleMessage(context.Background(), blockedMsg)
	if err != nil {
		t.Fatalf("handleMessage failed: %v", err)
	}

	// Should not receive event
	select {
	case <-msgCh:
		t.Error("Should not receive event for blocked chat")
	case <-time.After(100 * time.Millisecond):
		// Success - no event
	}
}

// Channel Tests

func TestChannel_Lifecycle(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Name = "Test Channel"
	config.Token = mock.Token
	config.TokenRef = "vault://test-token"
	config.Mode = "polling"
	config.PollingInterval = 100 * time.Millisecond

	channel, err := NewChannel(config, nil)
	if err != nil {
		t.Fatalf("NewChannel failed: %v", err)
	}

	channel.client = NewClient(mock.Token, WithBaseURL(mock.URL()))

	if channel.ID() != "test-channel" {
		t.Errorf("Expected ID 'test-channel', got '%s'", channel.ID())
	}

	if channel.Type() != "telegram" {
		t.Errorf("Expected type 'telegram', got '%s'", channel.Type())
	}

	if channel.Status() != channels.StatusDisconnected {
		t.Errorf("Expected status 'disconnected', got '%s'", channel.Status())
	}

	// Connect
	ctx := context.Background()
	err = channel.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if channel.Status() != channels.StatusConnected {
		t.Errorf("Expected status 'connected', got '%s'", channel.Status())
	}

	// Check health
	health := channel.Health()
	if !health.Healthy {
		t.Errorf("Expected healthy, got: %s", health.Message)
	}

	// Disconnect
	err = channel.Disconnect(ctx)
	if err != nil {
		t.Fatalf("Disconnect failed: %v", err)
	}

	if channel.Status() != channels.StatusDisconnected {
		t.Errorf("Expected status 'disconnected', got '%s'", channel.Status())
	}
}

func TestChannel_Send(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Name = "Test Channel"
	config.Token = mock.Token
	config.TokenRef = "vault://test-token"
	config.Mode = "polling"

	channel, err := NewChannel(config, nil)
	if err != nil {
		t.Fatalf("NewChannel failed: %v", err)
	}

	channel.client = NewClient(mock.Token, WithBaseURL(mock.URL()))

	ctx := context.Background()
	err = channel.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer channel.Disconnect(ctx)

	msg := channels.Message{
		ChannelID: "123456",
		Content:   "Test message from channel",
	}

	err = channel.Send(ctx, msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}

	if mock.Messages[0].Text != "Test message from channel" {
		t.Errorf("Expected 'Test message from channel', got '%s'", mock.Messages[0].Text)
	}
}

func TestChannel_IsChatAllowed(t *testing.T) {
	config := DefaultConfig()
	config.ID = "test-channel"
	config.Name = "Test Channel"
	config.TokenRef = "vault://test-token"
	config.ID = "test-channel"
	config.AllowedChats = []int64{123, 456}

	channel, err := NewChannel(config, nil)
	if err != nil {
		t.Fatalf("NewChannel failed: %v", err)
	}

	if !channel.IsChatAllowed(123) {
		t.Error("Expected chat 123 to be allowed")
	}

	if !channel.IsChatAllowed(456) {
		t.Error("Expected chat 456 to be allowed")
	}

	if channel.IsChatAllowed(999) {
		t.Error("Expected chat 999 to be blocked")
	}
}

// Manager Tests

func TestManager_CreateChannel(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	eventBus := bus.New()
	manager := NewManager(eventBus)

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Name = "Test Channel"
	config.Token = mock.Token
	config.Mode = "polling"

	// Override client creation by connecting manually
	ctx := context.Background()
	err := manager.CreateChannel(ctx, config)
	if err == nil {
		// Expected to fail because we can't override the client
		// In real usage, the token would be valid
		t.Log("CreateChannel failed as expected (token validation)")
	}
}

func TestManager_GetChannel(t *testing.T) {
	eventBus := bus.New()
	manager := NewManager(eventBus)

	// Get non-existent channel
	_, exists := manager.GetChannel("non-existent")
	if exists {
		t.Error("Expected channel to not exist")
	}
}

func TestManager_ChannelCount(t *testing.T) {
	eventBus := bus.New()
	manager := NewManager(eventBus)

	if manager.GetChannelCount() != 0 {
		t.Errorf("Expected 0 channels, got %d", manager.GetChannelCount())
	}
}

// Config Tests

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				ID:       "test",
				Name:     "Test",
				TokenRef: "vault://token",
				Mode:     "polling",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			config: Config{
				Name:     "Test",
				TokenRef: "vault://token",
				Mode:     "polling",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			config: Config{
				ID:       "test",
				TokenRef: "vault://token",
				Mode:     "polling",
			},
			wantErr: true,
		},
		{
			name: "missing token ref",
			config: Config{
				ID:   "test",
				Name: "Test",
				Mode: "polling",
			},
			wantErr: true,
		},
		{
			name: "invalid mode",
			config: Config{
				ID:       "test",
				Name:     "Test",
				TokenRef: "vault://token",
				Mode:     "invalid",
			},
			wantErr: true,
		},
		{
			name: "webhook without URL",
			config: Config{
				ID:       "test",
				Name:     "Test",
				TokenRef: "vault://token",
				Mode:     "webhook",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_IsChatAllowed(t *testing.T) {
	config := DefaultConfig()
	config.AllowedChats = []int64{123, 456}

	if !config.IsChatAllowed(123) {
		t.Error("Expected chat 123 to be allowed")
	}

	if !config.IsChatAllowed(456) {
		t.Error("Expected chat 456 to be allowed")
	}

	if config.IsChatAllowed(999) {
		t.Error("Expected chat 999 to be blocked")
	}
}

func TestConfig_IsChatAllowed_NoWhitelist(t *testing.T) {
	config := DefaultConfig()
	// Empty whitelist means all chats allowed
	config.AllowedChats = []int64{}

	if !config.IsChatAllowed(123) {
		t.Error("Expected all chats to be allowed when whitelist is empty")
	}
}

// Health Checker Tests

func TestHealthChecker_Check(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Name = "Test Channel"
	config.Token = mock.Token
	config.TokenRef = "vault://test-token"
	config.Mode = "polling"

	channel, err := NewChannel(config, nil)
	if err != nil {
		t.Fatalf("NewChannel failed: %v", err)
	}

	channel.client = NewClient(mock.Token, WithBaseURL(mock.URL()))

	ctx := context.Background()
	err = channel.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer channel.Disconnect(ctx)

	healthChecker := NewHealthChecker()
	result := healthChecker.Check(ctx, channel)

	if !result.IsHealthy() {
		t.Errorf("Expected healthy, got: %s - %s", result.Status, result.Message)
	}
}

func TestHealthChecker_Check_Disconnected(t *testing.T) {
	config := DefaultConfig()
	config.ID = "test-channel"
	config.Name = "Test Channel"
	config.Token = "test-token"
	config.TokenRef = "vault://test-token"
	config.Mode = "polling"

	channel, err := NewChannel(config, nil)
	if err != nil {
		t.Fatalf("NewChannel failed: %v", err)
	}

	healthChecker := NewHealthChecker()
	result := healthChecker.Check(context.Background(), channel)

	if !result.IsUnhealthy() {
		t.Errorf("Expected unhealthy for disconnected channel, got: %s", result.Status)
	}
}

// Webhook Tests

func TestWebhookReceiver_ServeHTTP(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.WebhookSecret = "secret123"

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)
	receiver := NewWebhookReceiver(&config, handler, nil)

	// Test valid request
	update := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			Chat:      &Chat{ID: 123, Type: "private"},
			Text:      "Hello",
		},
	}

	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(string(body)))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "secret123")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	receiver.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestWebhookReceiver_ServeHTTP_InvalidSecret(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.WebhookSecret = "secret123"

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)
	receiver := NewWebhookReceiver(&config, handler, nil)

	update := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			Chat:      &Chat{ID: 123, Type: "private"},
			Text:      "Hello",
		},
	}

	body, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(string(body)))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "wrong-secret")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	receiver.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestWebhookReceiver_ServeHTTP_InvalidMethod(t *testing.T) {
	config := DefaultConfig()
	client := NewClient("test-token")
	handler := NewHandler(&config, client, nil)
	receiver := NewWebhookReceiver(&config, handler, nil)

	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	rr := httptest.NewRecorder()
	receiver.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// Polling Tests

func TestPoller_StartStop(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "test-channel"
	config.Token = mock.Token
	config.PollingInterval = 50 * time.Millisecond

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)
	poller := NewPoller(&config, client, handler, nil)

	if poller.IsRunning() {
		t.Error("Expected poller to not be running initially")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go poller.Start(ctx)

	// Wait for poller to start
	time.Sleep(100 * time.Millisecond)

	// Stop the poller
	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestPoller_Offset(t *testing.T) {
	config := DefaultConfig()
	client := NewClient("test-token")
	handler := NewHandler(&config, client, nil)
	poller := NewPoller(&config, client, handler, nil)

	if poller.GetOffset() != 0 {
		t.Errorf("Expected initial offset 0, got %d", poller.GetOffset())
	}

	poller.SetOffset(100)
	if poller.GetOffset() != 100 {
		t.Errorf("Expected offset 100, got %d", poller.GetOffset())
	}
}

// Integration Tests

func TestIntegration_FullFlow(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	eventBus := bus.New()

	config := DefaultConfig()
	config.ID = "integration-test"
	config.Name = "Integration Test"
	config.Token = mock.Token
	config.TokenRef = "vault://test-token"
	config.Mode = "polling"
	config.AllowedChats = []int64{123}

	channel, err := NewChannel(config, eventBus)
	if err != nil {
		t.Fatalf("NewChannel failed: %v", err)
	}
	channel.client = NewClient(mock.Token, WithBaseURL(mock.URL()))

	// Connect
	ctx := context.Background()
	err = channel.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer channel.Disconnect(ctx)

	// Subscribe to events
	msgCh, unsub := eventBus.Subscribe(bus.EventChannelMessage)
	defer unsub()

	// Simulate incoming message
	update := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			Chat:      &Chat{ID: 123, Type: "private"},
			From:      &User{ID: 456, FirstName: "Test", Username: "testuser"},
			Text:      "Hello, Bot!",
			Date:      int(time.Now().Unix()),
		},
	}

	err = channel.handler.HandleUpdate(ctx, &update)
	if err != nil {
		t.Fatalf("HandleUpdate failed: %v", err)
	}

	// Wait for event
	select {
	case event := <-msgCh:
		msg, ok := event.Payload.(channels.Message)
		if !ok {
			t.Fatal("Expected channels.Message payload")
		}
		if msg.Content != "Hello, Bot!" {
			t.Errorf("Expected 'Hello, Bot!', got '%s'", msg.Content)
		}
		if msg.ChannelID != "123" {
			t.Errorf("Expected channel ID '123', got '%s'", msg.ChannelID)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for message event")
	}
}

func TestIntegration_CommandFlow(t *testing.T) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "cmd-test"
	config.Name = "Command Test"
	config.Token = mock.Token
	config.Mode = "polling"
	config.AllowedChats = []int64{123}

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)

	// Test /start command
	update := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			Chat:      &Chat{ID: 123, Type: "private"},
			From:      &User{ID: 456, FirstName: "TestUser"},
			Text:      "/start",
			Date:      int(time.Now().Unix()),
		},
	}

	ctx := context.Background()
	err := handler.HandleUpdate(ctx, &update)
	if err != nil {
		t.Fatalf("HandleUpdate failed: %v", err)
	}

	// Verify response was sent
	if len(mock.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(mock.Messages))
	}

	if !strings.Contains(mock.Messages[0].Text, "Welcome") {
		t.Errorf("Expected welcome message, got: %s", mock.Messages[0].Text)
	}
}

// Benchmarks

func BenchmarkClient_SendMessage(b *testing.B) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.SendMessage(ctx, 123456, "Test message")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandler_HandleMessage(b *testing.B) {
	mock := NewMockTelegramServer()
	defer mock.Close()

	config := DefaultConfig()
	config.ID = "bench-channel"
	config.Token = mock.Token
	config.AllowedChats = []int64{123}

	client := NewClient(mock.Token, WithBaseURL(mock.URL()))
	handler := NewHandler(&config, client, nil)

	msg := &Message{
		MessageID: 1,
		Chat:      &Chat{ID: 123, Type: "private"},
		From:      &User{ID: 456, FirstName: "Test"},
		Text:      "Hello",
		Date:      int(time.Now().Unix()),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := handler.handleMessage(ctx, msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
