package discord

import (
	"encoding/json"
	"time"
)

// GatewayOpCode represents the operation code in Gateway payloads
type GatewayOpCode int

const (
	GatewayOpDispatch            GatewayOpCode = 0
	GatewayOpHeartbeat           GatewayOpCode = 1
	GatewayOpIdentify            GatewayOpCode = 2
	GatewayOpPresenceUpdate      GatewayOpCode = 3
	GatewayOpVoiceStateUpdate    GatewayOpCode = 4
	GatewayOpResume              GatewayOpCode = 6
	GatewayOpReconnect           GatewayOpCode = 7
	GatewayOpRequestGuildMembers GatewayOpCode = 8
	GatewayOpInvalidSession      GatewayOpCode = 9
	GatewayOpHello               GatewayOpCode = 10
	GatewayOpHeartbeatACK        GatewayOpCode = 11
)

// GatewayPayload represents a payload sent over the Discord Gateway
type GatewayPayload struct {
	Op       int             `json:"op"`
	Data     json.RawMessage `json:"d"`
	Sequence *int            `json:"s,omitempty"`
	Type     string          `json:"t,omitempty"`
}

// GatewayEvent represents a dispatched event from the Gateway
type GatewayEvent struct {
	Op       int
	Type     string
	Data     json.RawMessage
	Sequence int
}

// IdentifyData represents the data sent in an Identify payload
type IdentifyData struct {
	Token          string        `json:"token"`
	Properties     Properties    `json:"properties"`
	Compress       bool          `json:"compress,omitempty"`
	LargeThreshold int           `json:"large_threshold,omitempty"`
	Shard          *[2]int       `json:"shard,omitempty"`
	Presence       *UpdateStatus `json:"presence,omitempty"`
	Intents        Intent        `json:"intents"`
}

// Properties represents the client properties sent during identification
type Properties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

// UpdateStatus represents the initial presence data
type UpdateStatus struct {
	Since      *int       `json:"since"`
	Activities []Activity `json:"activities"`
	Status     string     `json:"status"`
	AFK        bool       `json:"afk"`
}

// Activity represents a Discord activity (rich presence)
type Activity struct {
	Name          string              `json:"name"`
	Type          ActivityType        `json:"type"`
	URL           *string             `json:"url,omitempty"`
	CreatedAt     int64               `json:"created_at"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID string              `json:"application_id,omitempty"`
	Details       *string             `json:"details,omitempty"`
	State         *string             `json:"state,omitempty"`
	Emoji         *Emoji              `json:"emoji,omitempty"`
}

// ActivityType represents the type of activity
type ActivityType int

const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

// ActivityTimestamps represents timestamps for an activity
type ActivityTimestamps struct {
	Start *int64 `json:"start,omitempty"`
	End   *int64 `json:"end,omitempty"`
}

// HelloData represents the data in a Hello payload
type HelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// ReadyData represents the data in a Ready event
type ReadyData struct {
	Version     int         `json:"v"`
	User        User        `json:"user"`
	Guilds      []Guild     `json:"guilds"`
	SessionID   string      `json:"session_id"`
	Shard       *[2]int     `json:"shard,omitempty"`
	Application Application `json:"application"`
}

// ResumeData represents the data sent in a Resume payload
type ResumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

// User represents a Discord user
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
	System        bool   `json:"system,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Banner        string `json:"banner,omitempty"`
	AccentColor   int    `json:"accent_color,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int    `json:"flags,omitempty"`
	PremiumType   int    `json:"premium_type,omitempty"`
	PublicFlags   int    `json:"public_flags,omitempty"`
}

// Guild represents a Discord guild (server)
type Guild struct {
	ID                          string    `json:"id"`
	Name                        string    `json:"name"`
	Icon                        string    `json:"icon,omitempty"`
	Description                 string    `json:"description,omitempty"`
	Splash                      string    `json:"splash,omitempty"`
	OwnerID                     string    `json:"owner_id"`
	AFKChannelID                string    `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int       `json:"afk_timeout"`
	WidgetEnabled               bool      `json:"widget_enabled,omitempty"`
	WidgetChannelID             string    `json:"widget_channel_id,omitempty"`
	VerificationLevel           int       `json:"verification_level"`
	DefaultMessageNotifications int       `json:"default_message_notifications"`
	ExplicitContentFilter       int       `json:"explicit_content_filter"`
	Roles                       []Role    `json:"roles"`
	Emojis                      []Emoji   `json:"emojis"`
	Features                    []string  `json:"features"`
	MFALevel                    int       `json:"mfa_level"`
	ApplicationID               string    `json:"application_id,omitempty"`
	SystemChannelID             string    `json:"system_channel_id,omitempty"`
	RulesChannelID              string    `json:"rules_channel_id,omitempty"`
	JoinedAt                    time.Time `json:"joined_at,omitempty"`
	Large                       bool      `json:"large,omitempty"`
	Unavailable                 bool      `json:"unavailable,omitempty"`
	MemberCount                 int       `json:"member_count,omitempty"`
}

// Role represents a Discord role
type Role struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Color        int    `json:"color"`
	Hoist        bool   `json:"hoist"`
	Icon         string `json:"icon,omitempty"`
	UnicodeEmoji string `json:"unicode_emoji,omitempty"`
	Position     int    `json:"position"`
	Permissions  string `json:"permissions"`
	Managed      bool   `json:"managed"`
	Mentionable  bool   `json:"mentionable"`
}

// Emoji represents a Discord emoji
type Emoji struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles,omitempty"`
	User          *User    `json:"user,omitempty"`
	RequireColons bool     `json:"require_colons,omitempty"`
	Managed       bool     `json:"managed,omitempty"`
	Animated      bool     `json:"animated,omitempty"`
	Available     bool     `json:"available,omitempty"`
}

// Application represents a Discord application
type Application struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Icon                string `json:"icon,omitempty"`
	Description         string `json:"description,omitempty"`
	BotPublic           bool   `json:"bot_public"`
	BotRequireCodeGrant bool   `json:"bot_require_code_grant"`
	TermsOfServiceURL   string `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL    string `json:"privacy_policy_url,omitempty"`
	Owner               *User  `json:"owner,omitempty"`
}

// Channel represents a Discord channel
type Channel struct {
	ID                         string                `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildID                    string                `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      string                `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              string                `json:"last_message_id,omitempty"`
	Bitrate                    int                   `json:"bitrate,omitempty"`
	UserLimit                  int                   `json:"user_limit,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	Recipients                 []User                `json:"recipients,omitempty"`
	Icon                       string                `json:"icon,omitempty"`
	OwnerID                    string                `json:"owner_id,omitempty"`
	ApplicationID              string                `json:"application_id,omitempty"`
	ParentID                   string                `json:"parent_id,omitempty"`
	LastPinTimestamp           *time.Time            `json:"last_pin_timestamp,omitempty"`
	RTCRegion                  string                `json:"rtc_region,omitempty"`
	VideoQualityMode           int                   `json:"video_quality_mode,omitempty"`
	MessageCount               int                   `json:"message_count,omitempty"`
	MemberCount                int                   `json:"member_count,omitempty"`
	ThreadMetadata             *ThreadMetadata       `json:"thread_metadata,omitempty"`
	Member                     *ThreadMember         `json:"member,omitempty"`
	DefaultAutoArchiveDuration int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                string                `json:"permissions,omitempty"`
}

// ChannelType represents the type of a channel
type ChannelType int

const (
	ChannelTypeGuildText          ChannelType = 0
	ChannelTypeDM                 ChannelType = 1
	ChannelTypeGuildVoice         ChannelType = 2
	ChannelTypeGroupDM            ChannelType = 3
	ChannelTypeGuildCategory      ChannelType = 4
	ChannelTypeGuildNews          ChannelType = 5
	ChannelTypeGuildStore         ChannelType = 6
	ChannelTypeGuildNewsThread    ChannelType = 10
	ChannelTypeGuildPublicThread  ChannelType = 11
	ChannelTypeGuildPrivateThread ChannelType = 12
	ChannelTypeGuildStageVoice    ChannelType = 13
	ChannelTypeGuildDirectory     ChannelType = 14
	ChannelTypeGuildForum         ChannelType = 15
)

// PermissionOverwrite represents a permission overwrite for a channel
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  int    `json:"type"`
	Allow string `json:"allow"`
	Deny  string `json:"deny"`
}

// ThreadMetadata represents metadata for a thread
type ThreadMetadata struct {
	Archived            bool       `json:"archived"`
	AutoArchiveDuration int        `json:"auto_archive_duration"`
	ArchiveTimestamp    time.Time  `json:"archive_timestamp"`
	Locked              bool       `json:"locked,omitempty"`
	Invitable           bool       `json:"invitable,omitempty"`
	CreateTimestamp     *time.Time `json:"create_timestamp,omitempty"`
}

// ThreadMember represents a member of a thread
type ThreadMember struct {
	ID            string    `json:"id,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         int       `json:"flags"`
}

// Message represents a Discord message
type Message struct {
	ID                string              `json:"id"`
	ChannelID         string              `json:"channel_id"`
	GuildID           *string             `json:"guild_id,omitempty"`
	Author            *User               `json:"author,omitempty"`
	Member            *GuildMember        `json:"member,omitempty"`
	Content           string              `json:"content"`
	Timestamp         time.Time           `json:"timestamp"`
	EditedTimestamp   *time.Time          `json:"edited_timestamp,omitempty"`
	TTS               bool                `json:"tts"`
	MentionEveryone   bool                `json:"mention_everyone"`
	Mentions          []User              `json:"mentions"`
	MentionRoles      []string            `json:"mention_roles"`
	MentionChannels   []ChannelMention    `json:"mention_channels,omitempty"`
	Attachments       []Attachment        `json:"attachments"`
	Embeds            []Embed             `json:"embeds"`
	Reactions         []Reaction          `json:"reactions,omitempty"`
	Nonce             interface{}         `json:"nonce,omitempty"`
	Pinned            bool                `json:"pinned"`
	WebhookID         string              `json:"webhook_id,omitempty"`
	Type              MessageType         `json:"type"`
	Activity          *MessageActivity    `json:"activity,omitempty"`
	Application       *Application        `json:"application,omitempty"`
	MessageReference  *MessageReference   `json:"message_reference,omitempty"`
	Flags             int                 `json:"flags,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message,omitempty"`
	Interaction       *MessageInteraction `json:"interaction,omitempty"`
	Thread            *Channel            `json:"thread,omitempty"`
	Components        []MessageComponent  `json:"components,omitempty"`
	StickerItems      []StickerItem       `json:"sticker_items,omitempty"`
}

// MessageType represents the type of a message
type MessageType int

const (
	MessageTypeDefault                                 MessageType = 0
	MessageTypeRecipientAdd                            MessageType = 1
	MessageTypeRecipientRemove                         MessageType = 2
	MessageTypeCall                                    MessageType = 3
	MessageTypeChannelNameChange                       MessageType = 4
	MessageTypeChannelIconChange                       MessageType = 5
	MessageTypeChannelPinnedMessage                    MessageType = 6
	MessageTypeUserJoin                                MessageType = 7
	MessageTypeGuildBoost                              MessageType = 8
	MessageTypeGuildBoostTier1                         MessageType = 9
	MessageTypeGuildBoostTier2                         MessageType = 10
	MessageTypeGuildBoostTier3                         MessageType = 11
	MessageTypeChannelFollowAdd                        MessageType = 12
	MessageTypeGuildDiscoveryDisqualified              MessageType = 14
	MessageTypeGuildDiscoveryRequalified               MessageType = 15
	MessageTypeGuildDiscoveryGracePeriodInitialWarning MessageType = 16
	MessageTypeGuildDiscoveryGracePeriodFinalWarning   MessageType = 17
	MessageTypeThreadCreated                           MessageType = 18
	MessageTypeReply                                   MessageType = 19
	MessageTypeChatInputCommand                        MessageType = 20
	MessageTypeThreadStarterMessage                    MessageType = 21
	MessageTypeGuildInviteReminder                     MessageType = 22
	MessageTypeContextMenuCommand                      MessageType = 23
)

// GuildMember represents a member of a guild
type GuildMember struct {
	User                       *User      `json:"user,omitempty"`
	Nick                       string     `json:"nick,omitempty"`
	Avatar                     string     `json:"avatar,omitempty"`
	Roles                      []string   `json:"roles"`
	JoinedAt                   time.Time  `json:"joined_at"`
	PremiumSince               *time.Time `json:"premium_since,omitempty"`
	Deaf                       bool       `json:"deaf"`
	Mute                       bool       `json:"mute"`
	Pending                    bool       `json:"pending,omitempty"`
	Permissions                string     `json:"permissions,omitempty"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until,omitempty"`
}

// ChannelMention represents a channel mention in a message
type ChannelMention struct {
	ID      string      `json:"id"`
	GuildID string      `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}

// Attachment represents a file attachment in a message
type Attachment struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Height      *int   `json:"height,omitempty"`
	Width       *int   `json:"width,omitempty"`
	Ephemeral   bool   `json:"ephemeral,omitempty"`
}

// Embed represents an embed in a message
type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       int             `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField    `json:"fields,omitempty"`
}

// EmbedFooter represents the footer of an embed
type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedImage represents an image in an embed
type EmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// EmbedThumbnail represents a thumbnail in an embed
type EmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// EmbedVideo represents a video in an embed
type EmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// EmbedProvider represents a provider in an embed
type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// EmbedAuthor represents the author of an embed
type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedField represents a field in an embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Reaction represents a reaction to a message
type Reaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

// MessageActivity represents an activity in a message
type MessageActivity struct {
	Type    int    `json:"type"`
	PartyID string `json:"party_id,omitempty"`
}

// MessageReference represents a reference to another message
type MessageReference struct {
	MessageID       string `json:"message_id,omitempty"`
	ChannelID       string `json:"channel_id,omitempty"`
	GuildID         string `json:"guild_id,omitempty"`
	FailIfNotExists *bool  `json:"fail_if_not_exists,omitempty"`
}

// MessageInteraction represents an interaction in a message
type MessageInteraction struct {
	ID   string          `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User User            `json:"user"`
}

// InteractionType represents the type of an interaction
type InteractionType int

const (
	InteractionTypePing                           InteractionType = 1
	InteractionTypeApplicationCommand             InteractionType = 2
	InteractionTypeMessageComponent               InteractionType = 3
	InteractionTypeApplicationCommandAutocomplete InteractionType = 4
	InteractionTypeModalSubmit                    InteractionType = 5
)

// StickerItem represents a sticker in a message
type StickerItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FormatType int    `json:"format_type"`
}

// MessageComponent represents a component in a message
type MessageComponent struct {
	Type        int                `json:"type"`
	Components  []MessageComponent `json:"components,omitempty"`
	Style       int                `json:"style,omitempty"`
	Label       string             `json:"label,omitempty"`
	Emoji       *ComponentEmoji    `json:"emoji,omitempty"`
	CustomID    string             `json:"custom_id,omitempty"`
	URL         string             `json:"url,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   *int               `json:"min_values,omitempty"`
	MaxValues   *int               `json:"max_values,omitempty"`
}

// ComponentEmoji represents an emoji in a component
type ComponentEmoji struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Animated bool   `json:"animated,omitempty"`
}

// SelectMenuOption represents an option in a select menu
type SelectMenuOption struct {
	Label       string          `json:"label"`
	Value       string          `json:"value"`
	Description string          `json:"description,omitempty"`
	Emoji       *ComponentEmoji `json:"emoji,omitempty"`
	Default     bool            `json:"default,omitempty"`
}

// Intent represents Gateway intents for receiving events
type Intent int

const (
	IntentGuilds                      Intent = 1 << 0
	IntentGuildMembers                Intent = 1 << 1
	IntentGuildModeration             Intent = 1 << 2
	IntentGuildEmojisAndStickers      Intent = 1 << 3
	IntentGuildIntegrations           Intent = 1 << 4
	IntentGuildWebhooks               Intent = 1 << 5
	IntentGuildInvites                Intent = 1 << 6
	IntentGuildVoiceStates            Intent = 1 << 7
	IntentGuildPresences              Intent = 1 << 8
	IntentGuildMessages               Intent = 1 << 9
	IntentGuildMessageReactions       Intent = 1 << 10
	IntentGuildMessageTyping          Intent = 1 << 11
	IntentDirectMessages              Intent = 1 << 12
	IntentDirectMessageReactions      Intent = 1 << 13
	IntentDirectMessageTyping         Intent = 1 << 14
	IntentMessageContent              Intent = 1 << 15
	IntentGuildScheduledEvents        Intent = 1 << 16
	IntentAutoModerationConfiguration Intent = 1 << 20
	IntentAutoModerationExecution     Intent = 1 << 21
)

// DefaultIntents returns the default intents for a bot
func DefaultIntents() Intent {
	return IntentGuilds |
		IntentGuildMessages |
		IntentDirectMessages |
		IntentMessageContent |
		IntentGuildMessageReactions
}

// ApplicationCommand represents a slash command
type ApplicationCommand struct {
	ID                       string                     `json:"id,omitempty"`
	Type                     ApplicationCommandType     `json:"type,omitempty"`
	ApplicationID            string                     `json:"application_id,omitempty"`
	GuildID                  string                     `json:"guild_id,omitempty"`
	Name                     string                     `json:"name"`
	NameLocalizations        map[string]string          `json:"name_localizations,omitempty"`
	Description              string                     `json:"description"`
	DescriptionLocalizations map[string]string          `json:"description_localizations,omitempty"`
	Options                  []ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                    `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                      `json:"dm_permission,omitempty"`
	DefaultPermission        bool                       `json:"default_permission,omitempty"`
	Version                  string                     `json:"version,omitempty"`
}

// ApplicationCommandType represents the type of an application command
type ApplicationCommandType int

const (
	ApplicationCommandTypeChatInput ApplicationCommandType = 1
	ApplicationCommandTypeUser      ApplicationCommandType = 2
	ApplicationCommandTypeMessage   ApplicationCommandType = 3
)

// ApplicationCommandOption represents an option for an application command
type ApplicationCommandOption struct {
	Type                     ApplicationCommandOptionType     `json:"type"`
	Name                     string                           `json:"name"`
	NameLocalizations        map[string]string                `json:"name_localizations,omitempty"`
	Description              string                           `json:"description"`
	DescriptionLocalizations map[string]string                `json:"description_localizations,omitempty"`
	Required                 bool                             `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  []ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []ChannelType                    `json:"channel_types,omitempty"`
	MinValue                 *float64                         `json:"min_value,omitempty"`
	MaxValue                 *float64                         `json:"max_value,omitempty"`
	MinLength                *int                             `json:"min_length,omitempty"`
	MaxLength                *int                             `json:"max_length,omitempty"`
	Autocomplete             bool                             `json:"autocomplete,omitempty"`
}

// ApplicationCommandOptionType represents the type of an option
type ApplicationCommandOptionType int

const (
	ApplicationCommandOptionTypeSubCommand      ApplicationCommandOptionType = 1
	ApplicationCommandOptionTypeSubCommandGroup ApplicationCommandOptionType = 2
	ApplicationCommandOptionTypeString          ApplicationCommandOptionType = 3
	ApplicationCommandOptionTypeInteger         ApplicationCommandOptionType = 4
	ApplicationCommandOptionTypeBoolean         ApplicationCommandOptionType = 5
	ApplicationCommandOptionTypeUser            ApplicationCommandOptionType = 6
	ApplicationCommandOptionTypeChannel         ApplicationCommandOptionType = 7
	ApplicationCommandOptionTypeRole            ApplicationCommandOptionType = 8
	ApplicationCommandOptionTypeMentionable     ApplicationCommandOptionType = 9
	ApplicationCommandOptionTypeNumber          ApplicationCommandOptionType = 10
	ApplicationCommandOptionTypeAttachment      ApplicationCommandOptionType = 11
)

// ApplicationCommandOptionChoice represents a choice for an option
type ApplicationCommandOptionChoice struct {
	Name              string            `json:"name"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Value             interface{}       `json:"value"`
}

// Interaction represents an interaction (slash command, button click, etc.)
type Interaction struct {
	ID             string           `json:"id"`
	ApplicationID  string           `json:"application_id"`
	Type           InteractionType  `json:"type"`
	Data           *InteractionData `json:"data,omitempty"`
	GuildID        string           `json:"guild_id,omitempty"`
	ChannelID      string           `json:"channel_id,omitempty"`
	Member         *GuildMember     `json:"member,omitempty"`
	User           *User            `json:"user,omitempty"`
	Token          string           `json:"token"`
	Version        int              `json:"version"`
	Message        *Message         `json:"message,omitempty"`
	AppPermissions string           `json:"app_permissions,omitempty"`
	Locale         string           `json:"locale,omitempty"`
	GuildLocale    string           `json:"guild_locale,omitempty"`
}

// InteractionData represents the data for an interaction
type InteractionData struct {
	ID            string                                    `json:"id"`
	Name          string                                    `json:"name"`
	Type          ApplicationCommandType                    `json:"type"`
	Resolved      *ResolvedData                             `json:"resolved,omitempty"`
	Options       []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	GuildID       string                                    `json:"guild_id,omitempty"`
	TargetID      string                                    `json:"target_id,omitempty"`
	CustomID      string                                    `json:"custom_id,omitempty"`
	ComponentType int                                       `json:"component_type,omitempty"`
	Values        []string                                  `json:"values,omitempty"`
}

// ResolvedData represents resolved data for an interaction
type ResolvedData struct {
	Users       map[string]User        `json:"users,omitempty"`
	Members     map[string]GuildMember `json:"members,omitempty"`
	Roles       map[string]Role        `json:"roles,omitempty"`
	Channels    map[string]Channel     `json:"channels,omitempty"`
	Messages    map[string]Message     `json:"messages,omitempty"`
	Attachments map[string]Attachment  `json:"attachments,omitempty"`
}

// ApplicationCommandInteractionDataOption represents an option in an interaction
type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Type    ApplicationCommandOptionType              `json:"type"`
	Value   interface{}                               `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                      `json:"focused,omitempty"`
}

// InteractionResponse represents a response to an interaction
type InteractionResponse struct {
	Type InteractionResponseType  `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

// InteractionResponseType represents the type of an interaction response
type InteractionResponseType int

const (
	InteractionResponseTypePong                                 InteractionResponseType = 1
	InteractionResponseTypeChannelMessageWithSource             InteractionResponseType = 4
	InteractionResponseTypeDeferredChannelMessageWithSource     InteractionResponseType = 5
	InteractionResponseTypeDeferredUpdateMessage                InteractionResponseType = 6
	InteractionResponseTypeUpdateMessage                        InteractionResponseType = 7
	InteractionResponseTypeApplicationCommandAutocompleteResult InteractionResponseType = 8
	InteractionResponseTypeModal                                InteractionResponseType = 9
)

// InteractionResponseData represents the data for an interaction response
type InteractionResponseData struct {
	TTS             bool                             `json:"tts,omitempty"`
	Content         string                           `json:"content,omitempty"`
	Embeds          []Embed                          `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions                 `json:"allowed_mentions,omitempty"`
	Flags           int                              `json:"flags,omitempty"`
	Components      []MessageComponent               `json:"components,omitempty"`
	Attachments     []Attachment                     `json:"attachments,omitempty"`
	Choices         []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	CustomID        string                           `json:"custom_id,omitempty"`
	Title           string                           `json:"title,omitempty"`
}

// AllowedMentions represents allowed mentions in a message
type AllowedMentions struct {
	Parse       []string `json:"parse,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Users       []string `json:"users,omitempty"`
	RepliedUser bool     `json:"replied_user,omitempty"`
}

// Webhook represents a Discord webhook
type Webhook struct {
	ID            string   `json:"id"`
	Type          int      `json:"type"`
	GuildID       string   `json:"guild_id,omitempty"`
	ChannelID     string   `json:"channel_id"`
	User          *User    `json:"user,omitempty"`
	Name          string   `json:"name,omitempty"`
	Avatar        string   `json:"avatar,omitempty"`
	Token         string   `json:"token,omitempty"`
	ApplicationID string   `json:"application_id,omitempty"`
	SourceGuild   *Guild   `json:"source_guild,omitempty"`
	SourceChannel *Channel `json:"source_channel,omitempty"`
	URL           string   `json:"url,omitempty"`
}
