package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ChannelConfig represents a simplified channel configuration for CLI
type ChannelConfig struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"` // telegram, discord, slack, webhook
	Name      string            `json:"name"`
	Enabled   bool              `json:"enabled"`
	Config    map[string]string `json:"config"` // Type-specific config
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

func runChannel(args []string) int {
	if len(args) < 1 {
		channelUsage()
		return 2
	}

	cmd := args[0]

	switch cmd {
	case "list", "ls":
		return runChannelList(args[1:])
	case "add":
		return runChannelAdd(args[1:])
	case "update":
		return runChannelUpdate(args[1:])
	case "remove", "rm", "delete":
		return runChannelRemove(args[1:])
	case "enable":
		return runChannelEnable(args[1:])
	case "disable":
		return runChannelDisable(args[1:])
	case "test":
		return runChannelTest(args[1:])
	case "status":
		return runChannelStatus(args[1:])
	case "sync":
		return runChannelSync(args[1:])
	case "help", "-h", "--help":
		channelUsage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		channelUsage()
		return 2
	}
}

func runChannelList(args []string) int {
	jsonOutput := false
	detailed := false

	for _, arg := range args {
		if arg == "--json" || arg == "-j" {
			jsonOutput = true
		}
		if arg == "--verbose" || arg == "-v" {
			detailed = true
		}
	}

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	if jsonOutput {
		data, err := json.MarshalIndent(channels, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal channels: %v\n", err)
			return 1
		}
		fmt.Println(string(data))
	} else {
		fmt.Printf("Channels (%d)\n", len(channels))
		fmt.Println(strings.Repeat("=", 50))

		if len(channels) == 0 {
			fmt.Println("No channels configured.")
			fmt.Println("Use 'pryx-core channel add <type> <name>' to add a channel.")
		} else {
			for _, ch := range channels {
				status := "disabled"
				if ch.Enabled {
					status = "enabled"
				}
				fmt.Printf("• %s [%s] - %s\n", ch.Name, ch.Type, status)

				if detailed {
					fmt.Printf("  ID: %s\n", ch.ID)
					if len(ch.Config) > 0 {
						fmt.Printf("  Config:\n")
						for k, v := range ch.Config {
							// Don't print sensitive values
							if strings.Contains(k, "token") || strings.Contains(k, "secret") {
								fmt.Printf("    %s: ***\n", k)
							} else {
								fmt.Printf("    %s: %s\n", k, v)
							}
						}
					}
					fmt.Printf("  Created: %s\n", ch.CreatedAt)
					fmt.Printf("  Updated: %s\n", ch.UpdatedAt)
					fmt.Println()
				}
			}
		}
	}

	return 0
}

func runChannelAdd(args []string) int {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: channel type and name required\n")
		fmt.Fprintf(os.Stderr, "Usage: pryx-core channel add <type> <name> [--<key> <value>...]\n")
		return 2
	}

	channelType := args[0]
	name := args[1]
	configValues := make(map[string]string)

	// Parse additional config values
	i := 2
	for i < len(args) {
		if strings.HasPrefix(args[i], "--") {
			key := strings.TrimPrefix(args[i], "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				configValues[key] = args[i+1]
				i += 2
			} else {
				configValues[key] = ""
				i += 1
			}
		} else {
			i++
		}
	}

	// Validate channel type
	validTypes := map[string]bool{
		"telegram": true,
		"discord":  true,
		"slack":    true,
		"webhook":  true,
	}

	if !validTypes[channelType] {
		fmt.Fprintf(os.Stderr, "Error: invalid channel type: %s\n", channelType)
		fmt.Fprintf(os.Stderr, "Valid types: telegram, discord, slack, webhook\n")
		return 1
	}

	// Load existing channels
	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	// Check for duplicate name
	for _, ch := range channels {
		if ch.Name == name {
			fmt.Fprintf(os.Stderr, "Error: channel name already exists: %s\n", name)
			return 1
		}
	}

	// Create new channel
	now := getTimestamp()
	newChannel := ChannelConfig{
		ID:        fmt.Sprintf("%s-%s", channelType, now),
		Type:      channelType,
		Name:      name,
		Enabled:   false, // Disabled by default until tested
		Config:    configValues,
		CreatedAt: now,
		UpdatedAt: now,
	}

	channels = append(channels, newChannel)

	// Save channels
	if err := saveChannels(channels); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save channels: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Added channel: %s (type: %s)\n", name, channelType)
	fmt.Printf("  ID: %s\n", newChannel.ID)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Set authentication (if required):")
	fmt.Printf("     pryx-core channel enable %s\n", name)
	fmt.Println("  2. Test the connection:")
	fmt.Printf("     pryx-core channel test %s\n", name)
	if missing := validateChannelConfig(newChannel); len(missing) > 0 {
		fmt.Println()
		fmt.Printf("⚠ Missing required config: %s\n", strings.Join(missing, ", "))
		printChannelConfigHelp(channelType, name)
	}

	return 0
}

func runChannelRemove(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: channel name or ID required\n")
		return 2
	}

	name := args[0]

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	found := false
	filtered := make([]ChannelConfig, 0, len(channels))
	for _, ch := range channels {
		if ch.Name == name || ch.ID == name {
			found = true
			fmt.Printf("Removing channel: %s (%s)\n", ch.Name, ch.ID)
		} else {
			filtered = append(filtered, ch)
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	if err := saveChannels(filtered); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save channels: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Removed channel: %s\n", name)
	return 0
}

func runChannelUpdate(args []string) int {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: channel name or ID and at least one config value required\n")
		fmt.Fprintf(os.Stderr, "Usage: pryx-core channel update <name> [--<key> <value>...] [--unset <key>...]\n")
		return 2
	}

	name := args[0]
	configValues := make(map[string]string)
	unsetKeys := make(map[string]bool)

	i := 1
	for i < len(args) {
		arg := args[i]
		switch {
		case arg == "--unset":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
				fmt.Fprintf(os.Stderr, "Error: --unset requires a key\n")
				return 2
			}
			unsetKeys[args[i+1]] = true
			i += 2
		case strings.HasPrefix(arg, "--"):
			key := strings.TrimPrefix(arg, "--")
			if key == "" {
				fmt.Fprintf(os.Stderr, "Error: invalid key\n")
				return 2
			}
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				configValues[key] = args[i+1]
				i += 2
			} else {
				configValues[key] = ""
				i += 1
			}
		default:
			i++
		}
	}

	if len(configValues) == 0 && len(unsetKeys) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no updates provided\n")
		return 2
	}

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	updated := false
	for i, ch := range channels {
		if ch.Name == name || ch.ID == name {
			if ch.Config == nil {
				ch.Config = map[string]string{}
			}
			for k, v := range configValues {
				ch.Config[k] = v
			}
			for k := range unsetKeys {
				delete(ch.Config, k)
			}
			channels[i].Config = ch.Config
			channels[i].UpdatedAt = getTimestamp()
			updated = true
			fmt.Printf("✓ Updated channel: %s\n", ch.Name)
		}
	}

	if !updated {
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	if err := saveChannels(channels); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save channels: %v\n", err)
		return 1
	}

	return 0
}

func runChannelEnable(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: channel name or ID required\n")
		return 2
	}

	name := args[0]

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	found := false
	for i, ch := range channels {
		if ch.Name == name || ch.ID == name {
			found = true
			if missing := validateChannelConfig(ch); len(missing) > 0 {
				fmt.Fprintf(os.Stderr, "Error: channel configuration incomplete for %s (%s)\n", ch.Name, ch.Type)
				fmt.Fprintf(os.Stderr, "Missing: %s\n", strings.Join(missing, ", "))
				printChannelConfigHelp(ch.Type, ch.Name)
				return 1
			}
			channels[i].Enabled = true
			channels[i].UpdatedAt = getTimestamp()
			fmt.Printf("✓ Enabled channel: %s\n", ch.Name)
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	if err := saveChannels(channels); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save channels: %v\n", err)
		return 1
	}

	return 0
}

func runChannelDisable(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: channel name or ID required\n")
		return 2
	}

	name := args[0]

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	found := false
	for i, ch := range channels {
		if ch.Name == name || ch.ID == name {
			found = true
			channels[i].Enabled = false
			channels[i].UpdatedAt = getTimestamp()
			fmt.Printf("✓ Disabled channel: %s\n", ch.Name)
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	if err := saveChannels(channels); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save channels: %v\n", err)
		return 1
	}

	return 0
}

func runChannelTest(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: channel name or ID required\n")
		return 2
	}

	name := args[0]

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	var target *ChannelConfig
	for _, ch := range channels {
		if ch.Name == name || ch.ID == name {
			target = &ch
			break
		}
	}

	if target == nil {
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	fmt.Printf("Testing channel: %s (%s)\n", target.Name, target.Type)
	fmt.Println(strings.Repeat("=", 40))

	missing := validateChannelConfig(*target)
	if len(missing) > 0 {
		fmt.Printf("✗ Missing required config: %s\n", strings.Join(missing, ", "))
		printChannelConfigHelp(target.Type, target.Name)
		return 1
	}
	fmt.Printf("✓ Required configuration present\n")

	fmt.Println()
	fmt.Println("Note: Full connection testing requires runtime to be running")
	fmt.Println("Start runtime with: pryx-core")
	fmt.Println("Or test in TUI for interactive verification")

	return 0
}

func runChannelStatus(args []string) int {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	if name != "" {
		// Show status for specific channel
		for _, ch := range channels {
			if ch.Name == name || ch.ID == name {
				fmt.Printf("Channel: %s\n", ch.Name)
				fmt.Println(strings.Repeat("=", 40))
				fmt.Printf("ID:      %s\n", ch.ID)
				fmt.Printf("Type:    %s\n", ch.Type)
				fmt.Printf("Status:  ")
				if ch.Enabled {
					fmt.Printf("enabled\n")
				} else {
					fmt.Printf("disabled\n")
				}
				fmt.Printf("Created: %s\n", ch.CreatedAt)
				fmt.Printf("Updated: %s\n", ch.UpdatedAt)

				if len(ch.Config) > 0 {
					fmt.Println("\nConfiguration:")
					for k, v := range ch.Config {
						fmt.Printf("  %s: ", k)
						if strings.Contains(k, "token") || strings.Contains(k, "secret") || strings.Contains(k, "password") {
							fmt.Println("***")
						} else {
							fmt.Println(v)
						}
					}
				}

				return 0
			}
		}
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	// Show status for all channels
	fmt.Printf("Channel Status\n")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()

	enabledCount := 0
	for _, ch := range channels {
		status := "disabled"
		if ch.Enabled {
			status = "enabled"
			enabledCount++
		}
		fmt.Printf("• %s [%s] - %s\n", ch.Name, ch.Type, status)
	}

	fmt.Println()
	fmt.Printf("Total: %d channels (%d enabled, %d disabled)\n",
		len(channels), enabledCount, len(channels)-enabledCount)

	return 0
}

func runChannelSync(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: channel name or ID required\n")
		return 2
	}

	name := args[0]

	channels, err := loadChannels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load channels: %v\n", err)
		return 1
	}

	var target *ChannelConfig
	for _, ch := range channels {
		if ch.Name == name || ch.ID == name {
			target = &ch
			break
		}
	}

	if target == nil {
		fmt.Fprintf(os.Stderr, "Error: channel not found: %s\n", name)
		return 1
	}

	fmt.Printf("Syncing channel: %s (%s)\n", target.Name, target.Type)

	// Channel-specific sync logic
	switch target.Type {
	case "discord":
		fmt.Println("Syncing Discord slash commands...")
		fmt.Println("(Requires runtime to be running)")
		fmt.Println("Start runtime with: pryx-core")
	case "slack":
		fmt.Println("Syncing Slack app configuration...")
		fmt.Println("(Requires runtime to be running)")
		fmt.Println("Start runtime with: pryx-core")
	default:
		fmt.Printf("Sync not required for %s channels\n", target.Type)
	}

	return 0
}

func channelUsage() {
	fmt.Println("pryx-core channel - Manage communication channels")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  list [--json] [--verbose]        List all channels")
	fmt.Println("  add <type> <name> [--key val]    Add a new channel")
	fmt.Println("  update <name> [--key val]       Update channel configuration")
	fmt.Println("  remove <name>                    Remove a channel")
	fmt.Println("  enable <name>                   Enable a channel")
	fmt.Println("  disable <name>                  Disable a channel")
	fmt.Println("  test <name>                     Test channel connection")
	fmt.Println("  status [name]                   Show channel status")
	fmt.Println("  sync <name>                     Sync channel configuration")
	fmt.Println("")
	fmt.Println("Channel types:")
	fmt.Println("  telegram                         Telegram bot")
	fmt.Println("  discord                          Discord bot")
	fmt.Println("  slack                            Slack app")
	fmt.Println("  webhook                          Webhook endpoint")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  pryx-core channel add telegram my-bot --token YOUR_TOKEN")
	fmt.Println("  pryx-core channel add discord my-bot --token YOUR_TOKEN")
	fmt.Println("  pryx-core channel add slack my-bot --bot-token xoxb-... --app-token xapp-...")
	fmt.Println("  pryx-core channel add webhook my-hook --url https://example.com/webhook")
	fmt.Println("  pryx-core channel add webhook my-local --port 8080 --path /webhooks/pryx")
	fmt.Println("  pryx-core channel update my-bot --token NEW_TOKEN")
	fmt.Println("  pryx-core channel update my-hook --url https://example.com/webhook")
	fmt.Println("  pryx-core channel update my-bot --unset token")
	fmt.Println("  pryx-core channel enable my-bot")
	fmt.Println("  pryx-core channel test my-bot")
}

func validateChannelConfig(ch ChannelConfig) []string {
	var missing []string
	switch ch.Type {
	case "telegram", "discord":
		if ch.Config["token"] == "" && ch.Config["token_ref"] == "" {
			missing = append(missing, "token")
		}
	case "slack":
		if ch.Config["bot_token"] == "" {
			missing = append(missing, "bot_token")
		}
		if ch.Config["app_token"] == "" {
			missing = append(missing, "app_token")
		}
	case "webhook":
		if ch.Config["url"] == "" {
			if portStr := strings.TrimSpace(ch.Config["port"]); portStr != "" {
				if port, err := strconv.Atoi(portStr); err != nil || port <= 0 {
					missing = append(missing, "port")
				}
			} else {
				missing = append(missing, "url or port")
			}
		}
	default:
		missing = append(missing, "valid channel type")
	}
	return missing
}

func printChannelConfigHelp(channelType, name string) {
	fmt.Println("How to fix:")
	switch channelType {
	case "telegram":
		fmt.Printf("  pryx-core channel remove %s\n", name)
		fmt.Printf("  pryx-core channel add telegram %s --token YOUR_TOKEN\n", name)
	case "discord":
		fmt.Printf("  pryx-core channel remove %s\n", name)
		fmt.Printf("  pryx-core channel add discord %s --token YOUR_TOKEN\n", name)
	case "slack":
		fmt.Printf("  pryx-core channel remove %s\n", name)
		fmt.Printf("  pryx-core channel add slack %s --bot-token xoxb-... --app-token xapp-...\n", name)
	case "webhook":
		fmt.Printf("  pryx-core channel remove %s\n", name)
		fmt.Printf("  pryx-core channel add webhook %s --url https://example.com/webhook\n", name)
		fmt.Printf("  pryx-core channel add webhook %s --port 8080 --path /webhooks/pryx\n", name)
	default:
		fmt.Println("  Check channel configuration in ~/.pryx/channels.json")
	}
	fmt.Println("  Or edit ~/.pryx/channels.json directly if you prefer.")
}

func loadChannels() ([]ChannelConfig, error) {
	path := getChannelsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ChannelConfig{}, nil
		}
		return nil, err
	}

	var channels []ChannelConfig
	if err := json.Unmarshal(data, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

func saveChannels(channels []ChannelConfig) error {
	path := getChannelsPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(channels, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func getChannelsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pryx", "channels.json")
}

func getTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
