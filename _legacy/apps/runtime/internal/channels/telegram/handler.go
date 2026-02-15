package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"pryx-core/internal/bus"
	"pryx-core/internal/channels"
)

// Handler processes incoming Telegram updates and commands
type Handler struct {
	config   *Config
	client   *Client
	eventBus *bus.Bus
	commands map[string]CommandFunc
}

// CommandFunc is a function that handles a bot command
type CommandFunc func(ctx context.Context, msg *Message, args string) error

// NewHandler creates a new Telegram update handler
func NewHandler(config *Config, client *Client, eventBus *bus.Bus) *Handler {
	h := &Handler{
		config:   config,
		client:   client,
		eventBus: eventBus,
		commands: make(map[string]CommandFunc),
	}

	h.registerDefaultCommands()
	return h
}

// registerDefaultCommands registers the built-in bot commands
func (h *Handler) registerDefaultCommands() {
	h.commands["start"] = h.handleStart
	h.commands["help"] = h.handleHelp
	h.commands["chat"] = h.handleChat
	h.commands["status"] = h.handleStatus
}

// RegisterCommand registers a custom command handler
func (h *Handler) RegisterCommand(name string, handler CommandFunc) {
	h.commands[name] = handler
}

// HandleUpdate processes a single Telegram update
func (h *Handler) HandleUpdate(ctx context.Context, update *Update) error {
	// Handle different update types
	switch {
	case update.Message != nil:
		return h.handleMessage(ctx, update.Message)
	case update.EditedMessage != nil:
		return h.handleEditedMessage(ctx, update.EditedMessage)
	case update.CallbackQuery != nil:
		return h.handleCallbackQuery(ctx, update.CallbackQuery)
	case update.MyChatMember != nil:
		return h.handleChatMemberUpdate(ctx, update.MyChatMember)
	default:
		// Ignore other update types for now
		return nil
	}
}

// handleMessage processes incoming messages
func (h *Handler) handleMessage(ctx context.Context, msg *Message) error {
	if msg.Chat == nil {
		return nil
	}

	// Check chat whitelist
	if !h.config.IsChatAllowed(msg.Chat.ID) {
		h.publishError(fmt.Sprintf("chat %d not in whitelist", msg.Chat.ID))
		return nil
	}

	// Check if it's a command
	if msg.Text != "" && strings.HasPrefix(msg.Text, "/") {
		return h.handleCommand(ctx, msg)
	}

	// Handle regular message - publish to event bus
	h.publishMessage(msg)
	return nil
}

// handleCommand processes bot commands
func (h *Handler) handleCommand(ctx context.Context, msg *Message) error {
	// Parse command and arguments
	parts := strings.Fields(msg.Text)
	if len(parts) == 0 {
		return nil
	}

	// Extract command name (remove @botname if present)
	cmdText := parts[0]
	cmdText = strings.TrimPrefix(cmdText, "/")

	// Remove bot username suffix if present
	if idx := strings.Index(cmdText, "@"); idx != -1 {
		cmdText = cmdText[:idx]
	}

	cmdText = strings.ToLower(cmdText)
	args := ""
	if len(parts) > 1 {
		args = strings.Join(parts[1:], " ")
	}

	// Find and execute command handler
	if handler, ok := h.commands[cmdText]; ok {
		return handler(ctx, msg, args)
	}

	// Unknown command
	return h.sendUnknownCommand(ctx, msg.Chat.ID, cmdText)
}

// handleStart handles the /start command
func (h *Handler) handleStart(ctx context.Context, msg *Message, args string) error {
	var userName string
	if msg.From != nil {
		if msg.From.FirstName != "" {
			userName = msg.From.FirstName
		} else if msg.From.Username != "" {
			userName = msg.From.Username
		}
	}

	welcomeText := fmt.Sprintf(
		"👋 *Welcome%s!*\n\n"+
			"I'm your Pryx assistant bot. I can help you:\n\n"+
			"• 💬 Chat with you in conversation mode\n"+
			"• 📊 Show system status\n"+
			"• 🔄 Process commands\n\n"+
			"Use /help to see available commands.",
		func() string {
			if userName != "" {
				return ", " + userName
			}
			return ""
		}(),
	)

	_, err := h.client.SendMessage(ctx, msg.Chat.ID, welcomeText, WithParseMode(ParseModeMarkdown))
	return err
}

// handleHelp handles the /help command
func (h *Handler) handleHelp(ctx context.Context, msg *Message, args string) error {
	helpText := "*Available Commands*\n\n" +
		"/start - Start the bot and show welcome message\n" +
		"/help - Show this help message\n" +
		"/chat - Start conversation mode\n" +
		"/status - Show bot status\n\n" +
		"You can also send me regular messages and I'll process them."

	_, err := h.client.SendMessage(ctx, msg.Chat.ID, helpText, WithParseMode(ParseModeMarkdown))
	return err
}

// handleChat handles the /chat command
func (h *Handler) handleChat(ctx context.Context, msg *Message, args string) error {
	response := "💬 *Conversation Mode Activated*\n\n" +
		"I'm ready to chat! Send me any message and I'll respond.\n\n" +
		"Your messages will be processed by the Pryx agent system."

	_, err := h.client.SendMessage(ctx, msg.Chat.ID, response, WithParseMode(ParseModeMarkdown))

	if err == nil {
		// Publish chat request event to start conversation
		h.publishChatRequest(msg)
	}

	return err
}

// handleStatus handles the /status command
func (h *Handler) handleStatus(ctx context.Context, msg *Message, args string) error {
	// Get bot info
	botInfo, err := h.client.GetMe(ctx)
	if err != nil {
		return h.sendError(ctx, msg.Chat.ID, "Failed to get bot status")
	}

	// Build status message
	statusText := fmt.Sprintf(
		"*Bot Status* 🤖\n\n"+
			"*Name:* %s\n"+
			"*Username:* @%s\n"+
			"*Mode:* %s\n"+
			"*Chat ID:* %d\n"+
			"*Parse Mode:* %s\n",
		botInfo.FirstName,
		botInfo.Username,
		strings.Title(h.config.Mode),
		msg.Chat.ID,
		h.config.ParseMode,
	)

	if len(h.config.AllowedChats) > 0 {
		statusText += fmt.Sprintf("\n*Whitelisted Chats:* %d", len(h.config.AllowedChats))
	}

	_, err = h.client.SendMessage(ctx, msg.Chat.ID, statusText, WithParseMode(ParseModeMarkdown))
	return err
}

// handleEditedMessage processes edited messages
func (h *Handler) handleEditedMessage(ctx context.Context, msg *Message) error {
	// For now, treat edited messages the same as new messages
	// In the future, could track message history
	return h.handleMessage(ctx, msg)
}

// handleCallbackQuery processes callback queries (inline keyboard buttons)
func (h *Handler) handleCallbackQuery(ctx context.Context, query *CallbackQuery) error {
	if query.Message == nil {
		return nil
	}

	// Acknowledge the callback
	_, err := h.client.SendMessage(ctx, query.Message.Chat.ID,
		fmt.Sprintf("You selected: %s", query.Data),
		WithReplyToMessageID(query.Message.MessageID))

	return err
}

// handleChatMemberUpdate processes chat member updates
func (h *Handler) handleChatMemberUpdate(ctx context.Context, update *ChatMemberUpdated) error {
	// Log chat member changes
	if h.eventBus != nil {
		h.eventBus.Publish(bus.NewEvent(bus.EventChannelStatus, "", map[string]interface{}{
			"channel_id": h.config.ID,
			"chat_id":    update.Chat.ID,
			"old_status": update.OldChatMember.Status,
			"new_status": update.NewChatMember.Status,
			"user":       update.From,
		}))
	}
	return nil
}

// sendUnknownCommand sends a response for unknown commands
func (h *Handler) sendUnknownCommand(ctx context.Context, chatID int64, cmd string) error {
	msg := fmt.Sprintf("❓ Unknown command: `/%s`\n\nUse /help to see available commands.", cmd)
	_, err := h.client.SendMessage(ctx, chatID, msg, WithParseMode(ParseModeMarkdown))
	return err
}

// sendError sends an error message to the chat
func (h *Handler) sendError(ctx context.Context, chatID int64, errMsg string) error {
	msg := fmt.Sprintf("❌ *Error:* %s", errMsg)
	_, err := h.client.SendMessage(ctx, chatID, msg, WithParseMode(ParseModeMarkdown))
	return err
}

// publishMessage publishes a message to the event bus
func (h *Handler) publishMessage(msg *Message) {
	if h.eventBus == nil {
		return
	}

	content := msg.Text
	if content == "" {
		content = h.formatMediaContent(msg)
	}

	channelMsg := channels.Message{
		ID:        strconv.Itoa(msg.MessageID),
		Content:   content,
		Source:    h.config.ID,
		ChannelID: strconv.FormatInt(msg.Chat.ID, 10),
		SenderID:  h.formatSenderID(msg.From),
		CreatedAt: time.Unix(int64(msg.Date), 0),
		Metadata:  h.extractMetadata(msg),
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventChannelMessage, "", channelMsg))
}

// publishChatRequest publishes a chat request event
func (h *Handler) publishChatRequest(msg *Message) {
	if h.eventBus == nil {
		return
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventChatRequest, "", map[string]interface{}{
		"channel_id": h.config.ID,
		"chat_id":    msg.Chat.ID,
		"user_id":    h.formatSenderID(msg.From),
		"username":   h.getUsername(msg.From),
		"message":    msg.Text,
	}))
}

// publishError publishes an error event
func (h *Handler) publishError(errMsg string) {
	if h.eventBus == nil {
		return
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventErrorOccurred, "", map[string]interface{}{
		"channel_id": h.config.ID,
		"error":      errMsg,
	}))
}

// formatMediaContent formats media messages for display
func (h *Handler) formatMediaContent(msg *Message) string {
	switch {
	case msg.Photo != nil && len(msg.Photo) > 0:
		if msg.Caption != "" {
			return fmt.Sprintf("[Photo] %s", msg.Caption)
		}
		return "[Photo]"
	case msg.Document != nil:
		if msg.Caption != "" {
			return fmt.Sprintf("[Document: %s] %s", msg.Document.FileName, msg.Caption)
		}
		return fmt.Sprintf("[Document: %s]", msg.Document.FileName)
	case msg.Voice != nil:
		return "[Voice Message]"
	case msg.Audio != nil:
		if msg.Caption != "" {
			return fmt.Sprintf("[Audio: %s] %s", msg.Audio.Title, msg.Caption)
		}
		return fmt.Sprintf("[Audio: %s]", msg.Audio.Title)
	case msg.Video != nil:
		if msg.Caption != "" {
			return fmt.Sprintf("[Video] %s", msg.Caption)
		}
		return "[Video]"
	case msg.Sticker != nil:
		return fmt.Sprintf("[Sticker: %s]", msg.Sticker.Emoji)
	case msg.Animation != nil:
		return "[GIF/Animation]"
	case msg.Location != nil:
		return fmt.Sprintf("[Location: %.6f, %.6f]", msg.Location.Latitude, msg.Location.Longitude)
	case msg.Contact != nil:
		return fmt.Sprintf("[Contact: %s %s]", msg.Contact.FirstName, msg.Contact.LastName)
	default:
		return "[Media]"
	}
}

// formatSenderID formats the sender ID
func (h *Handler) formatSenderID(user *User) string {
	if user == nil {
		return "unknown"
	}
	return strconv.FormatInt(user.ID, 10)
}

// getUsername gets the username or name from user
func (h *Handler) getUsername(user *User) string {
	if user == nil {
		return "unknown"
	}
	if user.Username != "" {
		return user.Username
	}
	if user.FirstName != "" {
		return user.FirstName
	}
	return strconv.FormatInt(user.ID, 10)
}

// extractMetadata extracts metadata from a message
func (h *Handler) extractMetadata(msg *Message) map[string]string {
	metadata := make(map[string]string)

	if msg.From != nil {
		metadata["username"] = msg.From.Username
		metadata["first_name"] = msg.From.FirstName
		metadata["last_name"] = msg.From.LastName
		metadata["language_code"] = msg.From.LanguageCode
	}

	if msg.Chat != nil {
		metadata["chat_type"] = msg.Chat.Type
		metadata["chat_title"] = msg.Chat.Title
	}

	if msg.ReplyToMessage != nil {
		metadata["reply_to_message_id"] = strconv.Itoa(msg.ReplyToMessage.MessageID)
	}

	return metadata
}

// SendResponse sends a response message to a chat
func (h *Handler) SendResponse(ctx context.Context, chatID int64, text string) error {
	_, err := h.client.SendMessage(ctx, chatID, text,
		WithParseMode(h.config.ParseMode),
		WithDisableWebPagePreview(h.config.DisableWebPagePreview),
		WithDisableNotification(h.config.DisableNotification))
	return err
}

// SendResponseWithMarkup sends a response with inline keyboard markup
func (h *Handler) SendResponseWithMarkup(ctx context.Context, chatID int64, text string, markup *InlineKeyboardMarkup) error {
	_, err := h.client.SendMessage(ctx, chatID, text,
		WithParseMode(h.config.ParseMode),
		WithDisableWebPagePreview(h.config.DisableWebPagePreview),
		WithDisableNotification(h.config.DisableNotification),
		WithReplyMarkup(markup))
	return err
}
