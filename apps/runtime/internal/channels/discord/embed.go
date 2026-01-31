package discord

import (
	"fmt"
	"time"
)

// EmbedBuilder helps construct rich embed messages
type EmbedBuilder struct {
	embed Embed
}

// NewEmbed creates a new embed builder
func NewEmbed() *EmbedBuilder {
	return &EmbedBuilder{
		embed: Embed{
			Type: "rich",
		},
	}
}

// SetTitle sets the title of the embed
func (b *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	b.embed.Title = title
	return b
}

// SetDescription sets the description of the embed
func (b *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	b.embed.Description = description
	return b
}

// SetURL sets the URL of the embed
func (b *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	b.embed.URL = url
	return b
}

// SetTimestamp sets the timestamp of the embed
func (b *EmbedBuilder) SetTimestamp(t time.Time) *EmbedBuilder {
	b.embed.Timestamp = &t
	return b
}

// SetColor sets the color of the embed (Discord color integer)
func (b *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	b.embed.Color = color
	return b
}

// SetColorHex sets the color using a hex code (e.g., "#FF0000")
func (b *EmbedBuilder) SetColorHex(hex string) *EmbedBuilder {
	color := hexToInt(hex)
	b.embed.Color = color
	return b
}

// AddField adds a field to the embed
func (b *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	field := EmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	b.embed.Fields = append(b.embed.Fields, field)
	return b
}

// AddInlineField adds an inline field to the embed
func (b *EmbedBuilder) AddInlineField(name, value string) *EmbedBuilder {
	return b.AddField(name, value, true)
}

// AddField adds a regular (non-inline) field to the embed
func (b *EmbedBuilder) AddRegularField(name, value string) *EmbedBuilder {
	return b.AddField(name, value, false)
}

// SetFooter sets the footer of the embed
func (b *EmbedBuilder) SetFooter(text, iconURL string) *EmbedBuilder {
	b.embed.Footer = &EmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return b
}

// SetFooterText sets only the footer text
func (b *EmbedBuilder) SetFooterText(text string) *EmbedBuilder {
	if b.embed.Footer == nil {
		b.embed.Footer = &EmbedFooter{}
	}
	b.embed.Footer.Text = text
	return b
}

// SetFooterIcon sets only the footer icon
func (b *EmbedBuilder) SetFooterIcon(iconURL string) *EmbedBuilder {
	if b.embed.Footer == nil {
		b.embed.Footer = &EmbedFooter{}
	}
	b.embed.Footer.IconURL = iconURL
	return b
}

// SetImage sets the image of the embed
func (b *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	b.embed.Image = &EmbedImage{
		URL: url,
	}
	return b
}

// SetThumbnail sets the thumbnail of the embed
func (b *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	b.embed.Thumbnail = &EmbedThumbnail{
		URL: url,
	}
	return b
}

// SetVideo sets the video of the embed
func (b *EmbedBuilder) SetVideo(url string) *EmbedBuilder {
	b.embed.Video = &EmbedVideo{
		URL: url,
	}
	return b
}

// SetAuthor sets the author of the embed
func (b *EmbedBuilder) SetAuthor(name, url, iconURL string) *EmbedBuilder {
	b.embed.Author = &EmbedAuthor{
		Name:    name,
		URL:     url,
		IconURL: iconURL,
	}
	return b
}

// SetAuthorName sets only the author name
func (b *EmbedBuilder) SetAuthorName(name string) *EmbedBuilder {
	if b.embed.Author == nil {
		b.embed.Author = &EmbedAuthor{}
	}
	b.embed.Author.Name = name
	return b
}

// SetAuthorURL sets only the author URL
func (b *EmbedBuilder) SetAuthorURL(url string) *EmbedBuilder {
	if b.embed.Author == nil {
		b.embed.Author = &EmbedAuthor{}
	}
	b.embed.Author.URL = url
	return b
}

// SetAuthorIcon sets only the author icon
func (b *EmbedBuilder) SetAuthorIcon(iconURL string) *EmbedBuilder {
	if b.embed.Author == nil {
		b.embed.Author = &EmbedAuthor{}
	}
	b.embed.Author.IconURL = iconURL
	return b
}

// SetProvider sets the provider of the embed
func (b *EmbedBuilder) SetProvider(name, url string) *EmbedBuilder {
	b.embed.Provider = &EmbedProvider{
		Name: name,
		URL:  url,
	}
	return b
}

// Build returns the constructed embed
func (b *EmbedBuilder) Build() Embed {
	return b.embed
}

// BuildPtr returns a pointer to the constructed embed
func (b *EmbedBuilder) BuildPtr() *Embed {
	return &b.embed
}

// Validate checks if the embed is valid according to Discord limits
func (b *EmbedBuilder) Validate() error {
	return validateEmbed(&b.embed)
}

// IsEmpty returns true if the embed has no content
func (b *EmbedBuilder) IsEmpty() bool {
	return b.embed.Title == "" &&
		b.embed.Description == "" &&
		b.embed.URL == "" &&
		len(b.embed.Fields) == 0 &&
		b.embed.Image == nil &&
		b.embed.Thumbnail == nil &&
		b.embed.Video == nil &&
		b.embed.Footer == nil &&
		b.embed.Author == nil
}

// hexToInt converts a hex color string to an integer
func hexToInt(hex string) int {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	var result int
	for _, c := range hex {
		result *= 16
		switch {
		case c >= '0' && c <= '9':
			result += int(c - '0')
		case c >= 'a' && c <= 'f':
			result += int(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			result += int(c - 'A' + 10)
		}
	}

	return result
}

// validateEmbed validates an embed according to Discord limits
func validateEmbed(embed *Embed) error {
	// Title limit: 256 characters
	if len(embed.Title) > 256 {
		return fmt.Errorf("embed title exceeds 256 characters")
	}

	// Description limit: 4096 characters
	if len(embed.Description) > 4096 {
		return fmt.Errorf("embed description exceeds 4096 characters")
	}

	// Fields limit: 25 fields
	if len(embed.Fields) > 25 {
		return fmt.Errorf("embed exceeds 25 fields")
	}

	// Field name limit: 256 characters
	// Field value limit: 1024 characters
	for i, field := range embed.Fields {
		if len(field.Name) > 256 {
			return fmt.Errorf("field %d name exceeds 256 characters", i)
		}
		if len(field.Value) > 1024 {
			return fmt.Errorf("field %d value exceeds 1024 characters", i)
		}
	}

	// Footer text limit: 2048 characters
	if embed.Footer != nil && len(embed.Footer.Text) > 2048 {
		return fmt.Errorf("embed footer text exceeds 2048 characters")
	}

	// Author name limit: 256 characters
	if embed.Author != nil && len(embed.Author.Name) > 256 {
		return fmt.Errorf("embed author name exceeds 256 characters")
	}

	return nil
}

// Common Discord colors
const (
	ColorDefault           = 0
	ColorAqua              = 1752220
	ColorGreen             = 3066993
	ColorBlue              = 3447003
	ColorPurple            = 10181046
	ColorGold              = 15844367
	ColorOrange            = 15105570
	ColorRed               = 15158332
	ColorGrey              = 9807270
	ColorDarkerGrey        = 8359053
	ColorNavy              = 3426654
	ColorDarkAqua          = 1146986
	ColorDarkGreen         = 2067276
	ColorDarkBlue          = 2123412
	ColorDarkPurple        = 7419530
	ColorDarkGold          = 12745742
	ColorDarkOrange        = 11027200
	ColorDarkRed           = 10038562
	ColorDarkGrey          = 9936031
	ColorLightGrey         = 12370112
	ColorDarkNavy          = 2899536
	ColorLuminousVividPink = 16580705
	ColorDarkVividPink     = 12320855
)

// SuccessEmbed creates a success-themed embed
func SuccessEmbed(title, description string) *EmbedBuilder {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorGreen)
}

// ErrorEmbed creates an error-themed embed
func ErrorEmbed(title, description string) *EmbedBuilder {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorRed)
}

// InfoEmbed creates an info-themed embed
func InfoEmbed(title, description string) *EmbedBuilder {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorBlue)
}

// WarningEmbed creates a warning-themed embed
func WarningEmbed(title, description string) *EmbedBuilder {
	return NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(ColorGold)
}
