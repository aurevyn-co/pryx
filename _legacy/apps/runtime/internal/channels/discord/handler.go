package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"pryx-core/internal/bus"
	"pryx-core/internal/channels"
)

// Handler processes incoming Discord Gateway events
type Handler struct {
	config   *Config
	client   *Client
	eventBus *bus.Bus
	commands map[string]CommandHandler
	appID    string
}

// CommandHandler is a function that handles a slash command
type CommandHandler func(ctx context.Context, interaction *Interaction) error

// NewHandler creates a new Discord event handler
func NewHandler(config *Config, client *Client, eventBus *bus.Bus) *Handler {
	h := &Handler{
		config:   config,
		client:   client,
		eventBus: eventBus,
		commands: make(map[string]CommandHandler),
	}

	h.registerDefaultCommands()
	return h
}

// registerDefaultCommands registers the built-in slash commands
func (h *Handler) registerDefaultCommands() {
	h.commands["chat"] = h.handleChat
	h.commands["status"] = h.handleStatus
	h.commands["help"] = h.handleHelp
}

// RegisterCommand registers a custom command handler
func (h *Handler) RegisterCommand(name string, handler CommandHandler) {
	h.commands[name] = handler
}

// HandleEvent processes a Gateway event
func (h *Handler) HandleEvent(ctx context.Context, event *GatewayEvent) error {
	switch event.Type {
	case "READY":
		return h.handleReady(ctx, event.Data)
	case "MESSAGE_CREATE":
		return h.handleMessageCreate(ctx, event.Data)
	case "INTERACTION_CREATE":
		return h.handleInteractionCreate(ctx, event.Data)
	case "GUILD_CREATE":
		return h.handleGuildCreate(ctx, event.Data)
	case "GUILD_DELETE":
		return h.handleGuildDelete(ctx, event.Data)
	default:
		// Ignore other events
		return nil
	}
}

// handleReady processes the Ready event
func (h *Handler) handleReady(ctx context.Context, data json.RawMessage) error {
	var ready ReadyData
	if err := json.Unmarshal(data, &ready); err != nil {
		return fmt.Errorf("failed to unmarshal ready: %w", err)
	}

	h.appID = ready.Application.ID

	// Register slash commands
	if err := h.registerSlashCommands(ctx); err != nil {
		// Log but don't fail - commands are optional
		h.publishError(fmt.Sprintf("Failed to register commands: %v", err))
	}

	// Publish status event
	h.publishStatus("ready", map[string]interface{}{
		"user_id":    ready.User.ID,
		"username":   ready.User.Username,
		"guilds":     len(ready.Guilds),
		"session_id": ready.SessionID,
	})

	return nil
}

// handleMessageCreate processes MESSAGE_CREATE events
func (h *Handler) handleMessageCreate(ctx context.Context, data json.RawMessage) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Ignore messages from bots
	if msg.Author != nil && msg.Author.Bot {
		return nil
	}

	// Check guild whitelist
	if msg.GuildID != nil && !h.config.IsGuildAllowed(*msg.GuildID) {
		return nil
	}

	// Check channel whitelist
	if !h.config.IsChannelAllowed(msg.ChannelID) {
		return nil
	}

	// Check if message mentions the bot
	isMention := false
	if msg.Author != nil {
		for _, mention := range msg.Mentions {
			if mention.Bot {
				isMention = true
				break
			}
		}
	}

	// Check for DM
	isDM := msg.GuildID == nil || *msg.GuildID == ""

	// Only process DMs or mentions
	if !isDM && !isMention {
		return nil
	}

	// Publish message to event bus
	h.publishMessage(&msg, isDM, isMention)

	return nil
}

// handleInteractionCreate processes INTERACTION_CREATE events
func (h *Handler) handleInteractionCreate(ctx context.Context, data json.RawMessage) error {
	var interaction Interaction
	if err := json.Unmarshal(data, &interaction); err != nil {
		return fmt.Errorf("failed to unmarshal interaction: %w", err)
	}

	// Only handle application commands
	if interaction.Type != InteractionTypeApplicationCommand {
		return nil
	}

	// Check guild whitelist
	if interaction.GuildID != "" && !h.config.IsGuildAllowed(interaction.GuildID) {
		return nil
	}

	// Check channel whitelist
	if interaction.ChannelID != "" && !h.config.IsChannelAllowed(interaction.ChannelID) {
		return nil
	}

	// Get command name
	if interaction.Data == nil {
		return nil
	}

	cmdName := interaction.Data.Name

	// Handle subcommands
	if len(interaction.Data.Options) > 0 && interaction.Data.Options[0].Type == ApplicationCommandOptionTypeSubCommand {
		cmdName = interaction.Data.Options[0].Name
	}

	// Find and execute command handler
	if handler, ok := h.commands[cmdName]; ok {
		if err := handler(ctx, &interaction); err != nil {
			// Send error response
			h.sendErrorResponse(ctx, &interaction, err.Error())
			return err
		}
	} else {
		// Unknown command
		h.sendErrorResponse(ctx, &interaction, "Unknown command")
	}

	return nil
}

// handleGuildCreate processes GUILD_CREATE events
func (h *Handler) handleGuildCreate(ctx context.Context, data json.RawMessage) error {
	var guild Guild
	if err := json.Unmarshal(data, &guild); err != nil {
		return fmt.Errorf("failed to unmarshal guild: %w", err)
	}

	// Check if guild is allowed
	if !h.config.IsGuildAllowed(guild.ID) {
		// Could leave guild here if desired
		return nil
	}

	h.publishStatus("guild_join", map[string]interface{}{
		"guild_id":   guild.ID,
		"guild_name": guild.Name,
	})

	return nil
}

// handleGuildDelete processes GUILD_DELETE events
func (h *Handler) handleGuildDelete(ctx context.Context, data json.RawMessage) error {
	var guild Guild
	if err := json.Unmarshal(data, &guild); err != nil {
		return fmt.Errorf("failed to unmarshal guild: %w", err)
	}

	h.publishStatus("guild_leave", map[string]interface{}{
		"guild_id":   guild.ID,
		"guild_name": guild.Name,
	})

	return nil
}

// handleChat handles the /pryx chat command
func (h *Handler) handleChat(ctx context.Context, interaction *Interaction) error {
	response := InteractionResponse{
		Type: InteractionResponseTypeChannelMessageWithSource,
		Data: &InteractionResponseData{
			Embeds: []Embed{
				*InfoEmbed("Chat Mode Activated", "I'm ready to chat! Send me messages and I'll respond.").
					SetFooter("Powered by Pryx", "").
					BuildPtr(),
			},
		},
	}

	if err := h.client.RespondToInteraction(ctx, interaction.ID, interaction.Token, &response); err != nil {
		return fmt.Errorf("failed to respond: %w", err)
	}

	// Publish chat request event
	h.publishChatRequest(interaction)

	return nil
}

// handleStatus handles the /pryx status command
func (h *Handler) handleStatus(ctx context.Context, interaction *Interaction) error {
	// Get bot info
	botInfo, err := h.client.GetMe(ctx)
	if err != nil {
		return h.sendErrorResponse(ctx, interaction, "Failed to get bot status")
	}

	embed := NewEmbed().
		SetTitle("Bot Status").
		SetColor(ColorBlue).
		AddField("Name", botInfo.Username, true).
		AddField("ID", botInfo.ID, true).
		AddField("Status", "Online", true).
		SetFooter("Pryx Discord Bot", "").
		Build()

	response := InteractionResponse{
		Type: InteractionResponseTypeChannelMessageWithSource,
		Data: &InteractionResponseData{
			Embeds: []Embed{embed},
		},
	}

	return h.client.RespondToInteraction(ctx, interaction.ID, interaction.Token, &response)
}

// handleHelp handles the /pryx help command
func (h *Handler) handleHelp(ctx context.Context, interaction *Interaction) error {
	embed := NewEmbed().
		SetTitle("Pryx Bot Commands").
		SetDescription("Here are the available commands:").
		SetColor(ColorBlue).
		AddField("/pryx chat", "Start a conversation with Pryx", false).
		AddField("/pryx status", "Check bot status and information", false).
		AddField("/pryx help", "Show this help message", false).
		AddField("Mention", "Mention me in a message to chat", false).
		SetFooter("Powered by Pryx", "").
		Build()

	response := InteractionResponse{
		Type: InteractionResponseTypeChannelMessageWithSource,
		Data: &InteractionResponseData{
			Embeds: []Embed{embed},
		},
	}

	return h.client.RespondToInteraction(ctx, interaction.ID, interaction.Token, &response)
}

// registerSlashCommands registers the bot's slash commands
func (h *Handler) registerSlashCommands(ctx context.Context) error {
	if h.appID == "" {
		return fmt.Errorf("application ID not set")
	}

	// Create main command with subcommands
	mainCmd := &ApplicationCommand{
		Name:        "pryx",
		Description: "Pryx AI Assistant",
		Options: []ApplicationCommandOption{
			{
				Type:        ApplicationCommandOptionTypeSubCommand,
				Name:        "chat",
				Description: "Start a conversation with Pryx",
			},
			{
				Type:        ApplicationCommandOptionTypeSubCommand,
				Name:        "status",
				Description: "Check bot status and information",
			},
			{
				Type:        ApplicationCommandOptionTypeSubCommand,
				Name:        "help",
				Description: "Show help information",
			},
		},
		DMPermission: boolPtr(true),
	}

	// Register global command
	_, err := h.client.CreateSlashCommand(ctx, h.appID, mainCmd)
	return err
}

// sendErrorResponse sends an error response to an interaction
func (h *Handler) sendErrorResponse(ctx context.Context, interaction *Interaction, errMsg string) error {
	embed := *ErrorEmbed("Error", errMsg).BuildPtr()

	response := InteractionResponse{
		Type: InteractionResponseTypeChannelMessageWithSource,
		Data: &InteractionResponseData{
			Embeds: []Embed{embed},
			Flags:  1 << 6, // Ephemeral flag
		},
	}

	return h.client.RespondToInteraction(ctx, interaction.ID, interaction.Token, &response)
}

// publishMessage publishes a message to the event bus
func (h *Handler) publishMessage(msg *Message, isDM, isMention bool) {
	if h.eventBus == nil {
		return
	}

	content := msg.Content
	// Remove bot mentions from content
	if isMention && msg.Author != nil {
		for _, mention := range msg.Mentions {
			if mention.Bot {
				mentionStr := fmt.Sprintf("<@%s>", mention.ID)
				content = strings.ReplaceAll(content, mentionStr, "")
				mentionStr = fmt.Sprintf("<@!%s>", mention.ID)
				content = strings.ReplaceAll(content, mentionStr, "")
			}
		}
		content = strings.TrimSpace(content)
	}

	channelMsg := channels.Message{
		ID:        msg.ID,
		Content:   content,
		Source:    h.config.ID,
		ChannelID: msg.ChannelID,
		SenderID:  "",
		CreatedAt: msg.Timestamp,
		Metadata:  h.extractMetadata(msg, isDM, isMention),
	}

	if msg.Author != nil {
		channelMsg.SenderID = msg.Author.ID
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventChannelMessage, "", channelMsg))
}

// publishChatRequest publishes a chat request event
func (h *Handler) publishChatRequest(interaction *Interaction) {
	if h.eventBus == nil {
		return
	}

	userID := ""
	username := ""

	if interaction.User != nil {
		userID = interaction.User.ID
		username = interaction.User.Username
	} else if interaction.Member != nil && interaction.Member.User != nil {
		userID = interaction.Member.User.ID
		username = interaction.Member.User.Username
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventChatRequest, "", map[string]interface{}{
		"channel_id": h.config.ID,
		"guild_id":   interaction.GuildID,
		"channel":    interaction.ChannelID,
		"user_id":    userID,
		"username":   username,
	}))
}

// publishStatus publishes a status event
func (h *Handler) publishStatus(status string, data map[string]interface{}) {
	if h.eventBus == nil {
		return
	}

	payload := map[string]interface{}{
		"channel_id":   h.config.ID,
		"channel_type": "discord",
		"status":       status,
	}

	for k, v := range data {
		payload[k] = v
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventChannelStatus, "", payload))
}

// publishError publishes an error event
func (h *Handler) publishError(errMsg string) {
	if h.eventBus == nil {
		return
	}

	h.eventBus.Publish(bus.NewEvent(bus.EventErrorOccurred, "", map[string]interface{}{
		"channel_id": h.config.ID,
		"error":      errMsg,
		"type":       "discord",
	}))
}

// extractMetadata extracts metadata from a message
func (h *Handler) extractMetadata(msg *Message, isDM, isMention bool) map[string]string {
	metadata := make(map[string]string)

	if msg.Author != nil {
		metadata["username"] = msg.Author.Username
		metadata["discriminator"] = msg.Author.Discriminator
		metadata["user_id"] = msg.Author.ID
	}

	if msg.GuildID != nil {
		metadata["guild_id"] = *msg.GuildID
	}

	metadata["channel_id"] = msg.ChannelID
	metadata["is_dm"] = fmt.Sprintf("%v", isDM)
	metadata["is_mention"] = fmt.Sprintf("%v", isMention)

	if msg.MessageReference != nil {
		metadata["reply_to"] = msg.MessageReference.MessageID
	}

	return metadata
}

// SendResponse sends a message response to a channel
func (h *Handler) SendResponse(ctx context.Context, channelID string, content string) error {
	_, err := h.client.SendMessage(ctx, channelID, content)
	return err
}

// SendEmbedResponse sends an embed response to a channel
func (h *Handler) SendEmbedResponse(ctx context.Context, channelID string, embed *Embed) error {
	_, err := h.client.SendEmbed(ctx, channelID, embed)
	return err
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}
