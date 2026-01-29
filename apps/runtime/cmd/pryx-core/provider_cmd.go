package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pryx-core/internal/config"
	"pryx-core/internal/keychain"
)

// ProviderInfo holds information about a provider
type ProviderInfo struct {
	Name        string
	DisplayName string
	Description string
	DefaultURL  string
}

// SupportedProviders is the list of supported LLM providers
var SupportedProviders = []ProviderInfo{
	{Name: "openai", DisplayName: "OpenAI", Description: "GPT-4, GPT-3.5", DefaultURL: "https://api.openai.com/v1"},
	{Name: "anthropic", DisplayName: "Anthropic", Description: "Claude 3 models", DefaultURL: "https://api.anthropic.com/v1"},
	{Name: "glm", DisplayName: "GLM (Zhipu)", Description: "GLM-4, ChatGLM", DefaultURL: "https://open.bigmodel.cn/api/paas/v4"},
	{Name: "openrouter", DisplayName: "OpenRouter", Description: "Multi-provider access", DefaultURL: "https://openrouter.ai/api/v1"},
	{Name: "together", DisplayName: "Together AI", Description: "Open source models", DefaultURL: "https://api.together.xyz/v1"},
	{Name: "groq", DisplayName: "Groq", Description: "Fast inference", DefaultURL: "https://api.groq.com/openai/v1"},
	{Name: "xai", DisplayName: "xAI", Description: "Grok models", DefaultURL: "https://api.x.ai/v1"},
	{Name: "mistral", DisplayName: "Mistral AI", Description: "Mistral models", DefaultURL: "https://api.mistral.ai/v1"},
	{Name: "cohere", DisplayName: "Cohere", Description: "Command models", DefaultURL: "https://api.cohere.com/v1"},
	{Name: "google", DisplayName: "Google AI", Description: "Gemini models", DefaultURL: "https://generativelanguage.googleapis.com/v1"},
	{Name: "ollama", DisplayName: "Ollama", Description: "Local models", DefaultURL: "http://localhost:11434"},
}

func runProvider(args []string) int {
	if len(args) < 1 {
		usageProvider()
		return 1
	}

	command := args[0]
	path := config.DefaultPath()
	cfg := config.Load()

	if fileCfg, err := config.LoadFromFile(path); err == nil {
		cfg = fileCfg
	}

	kc := keychain.New("pryx")

	switch command {
	case "list":
		return providerList(cfg, kc)
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: pryx-core provider add <name>")
			return 1
		}
		return providerAdd(args[1], cfg, path, kc)
	case "set-key":
		if len(args) < 2 {
			fmt.Println("Usage: pryx-core provider set-key <name>")
			return 1
		}
		return providerSetKey(args[1], kc)
	case "remove":
		if len(args) < 2 {
			fmt.Println("Usage: pryx-core provider remove <name>")
			return 1
		}
		return providerRemove(args[1], cfg, path, kc)
	case "use":
		if len(args) < 2 {
			fmt.Println("Usage: pryx-core provider use <name>")
			return 1
		}
		return providerUse(args[1], cfg, path, kc)
	case "test":
		if len(args) < 2 {
			fmt.Println("Usage: pryx-core provider test <name>")
			return 1
		}
		return providerTest(args[1], cfg, kc)
	default:
		usageProvider()
		return 1
	}
}

func usageProvider() {
	fmt.Println("Usage:")
	fmt.Println("  pryx-core provider list                    List all configured providers")
	fmt.Println("  pryx-core provider add <name>              Add new provider interactively")
	fmt.Println("  pryx-core provider set-key <name>          Set API key for provider")
	fmt.Println("  pryx-core provider remove <name>           Remove provider config")
	fmt.Println("  pryx-core provider use <name>              Set as active/default provider")
	fmt.Println("  pryx-core provider test <name>             Test connection to provider")
	fmt.Println("")
	fmt.Println("Supported providers:")
	for _, p := range SupportedProviders {
		fmt.Printf("  %-12s %s\n", p.Name, p.Description)
	}
}

func providerList(cfg *config.Config, kc *keychain.Keychain) int {
	fmt.Println("Configured Providers:")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-15s %-10s %-20s %s\n", "NAME", "STATUS", "API KEY", "ACTIVE")
	fmt.Println(strings.Repeat("-", 60))

	activeProvider := cfg.ModelProvider

	for _, p := range SupportedProviders {
		name := p.Name
		status := "not configured"
		keyStatus := "not set"
		isActive := ""

		// Check if provider has API key
		if key, err := kc.GetProviderKey(name); err == nil && key != "" {
			keyStatus = "set"
			status = "configured"
		}

		// Special handling for ollama (local, no key needed)
		if name == "ollama" {
			keyStatus = "n/a"
			if cfg.OllamaEndpoint != "" {
				status = "configured"
			}
		}

		if name == activeProvider {
			isActive = "*"
		}

		fmt.Printf("%-15s %-10s %-20s %s\n", name, status, keyStatus, isActive)
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("* = active provider")
	return 0
}

func providerAdd(name string, cfg *config.Config, path string, kc *keychain.Keychain) int {
	// Validate provider name
	providerInfo, ok := getProviderInfo(name)
	if !ok {
		fmt.Printf("Unknown provider: %s\n", name)
		fmt.Println("Run 'pryx-core provider list' to see supported providers.")
		return 1
	}

	fmt.Printf("Adding provider: %s (%s)\n", providerInfo.DisplayName, providerInfo.Description)
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	// For ollama, we just need the endpoint
	if name == "ollama" {
		fmt.Printf("Ollama endpoint [%s]: ", providerInfo.DefaultURL)
		endpoint, _ := reader.ReadString('\n')
		endpoint = strings.TrimSpace(endpoint)
		if endpoint == "" {
			endpoint = providerInfo.DefaultURL
		}
		cfg.OllamaEndpoint = endpoint
		cfg.ModelProvider = "ollama"
		cfg.ModelName = "llama3"

		if err := cfg.Save(path); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return 1
		}

		fmt.Println("Ollama provider configured successfully!")
		fmt.Printf("  Endpoint: %s\n", endpoint)
		fmt.Println("  Model: llama3 (change with 'pryx-core config set model_name <model>')")
		return 0
	}

	// For cloud providers, ask for API key
	fmt.Print("Enter API key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		fmt.Println("API key cannot be empty.")
		return 1
	}

	// Store key in keychain
	if err := kc.SetProviderKey(name, apiKey); err != nil {
		fmt.Printf("Failed to store API key: %v\n", err)
		return 1
	}

	// Ask if this should be the active provider
	fmt.Print("Set as active provider? [Y/n]: ")
	setActive, _ := reader.ReadString('\n')
	setActive = strings.TrimSpace(strings.ToLower(setActive))

	if setActive == "" || setActive == "y" || setActive == "yes" {
		cfg.ModelProvider = name
		// Set a reasonable default model
		cfg.ModelName = getDefaultModelForProvider(name)
	}

	if err := cfg.Save(path); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return 1
	}

	fmt.Printf("Provider '%s' added successfully!\n", name)
	if cfg.ModelProvider == name {
		fmt.Println("Set as active provider.")
	}

	return 0
}

func providerSetKey(name string, kc *keychain.Keychain) int {
	// Validate provider name
	if _, ok := getProviderInfo(name); !ok {
		fmt.Printf("Unknown provider: %s\n", name)
		return 1
	}

	if name == "ollama" {
		fmt.Println("Ollama provider does not require an API key (local deployment).")
		return 1
	}

	fmt.Printf("Setting API key for provider: %s\n", name)
	fmt.Print("Enter API key: ")

	reader := bufio.NewReader(os.Stdin)
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		fmt.Println("API key cannot be empty.")
		return 1
	}

	if err := kc.SetProviderKey(name, apiKey); err != nil {
		fmt.Printf("Failed to store API key: %v\n", err)
		return 1
	}

	fmt.Printf("API key for '%s' stored securely in keychain.\n", name)
	return 0
}

func providerRemove(name string, cfg *config.Config, path string, kc *keychain.Keychain) int {
	// Validate provider name
	if _, ok := getProviderInfo(name); !ok {
		fmt.Printf("Unknown provider: %s\n", name)
		return 1
	}

	// Confirm removal
	fmt.Printf("Are you sure you want to remove provider '%s'? [y/N]: ", name)
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("Cancelled.")
		return 0
	}

	// Remove API key from keychain
	if err := kc.DeleteProviderKey(name); err != nil {
		fmt.Printf("Warning: Failed to remove API key from keychain: %v\n", err)
	}

	// If this was the active provider, reset to ollama
	if cfg.ModelProvider == name {
		cfg.ModelProvider = "ollama"
		cfg.ModelName = "llama3"
		if err := cfg.Save(path); err != nil {
			fmt.Printf("Failed to update config: %v\n", err)
			return 1
		}
		fmt.Println("Active provider reset to 'ollama'.")
	}

	fmt.Printf("Provider '%s' removed successfully.\n", name)
	return 0
}

func providerUse(name string, cfg *config.Config, path string, kc *keychain.Keychain) int {
	// Validate provider name
	providerInfo, ok := getProviderInfo(name)
	if !ok {
		fmt.Printf("Unknown provider: %s\n", name)
		return 1
	}

	// Check if provider is configured
	if name != "ollama" {
		if key, err := kc.GetProviderKey(name); err != nil || key == "" {
			fmt.Printf("Provider '%s' is not configured. Run 'pryx-core provider add %s' first.\n", name, name)
			return 1
		}
	}

	// Update config
	cfg.ModelProvider = name
	if cfg.ModelName == "" {
		cfg.ModelName = getDefaultModelForProvider(name)
	}

	if err := cfg.Save(path); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
		return 1
	}

	fmt.Printf("Now using provider: %s (%s)\n", providerInfo.DisplayName, providerInfo.Description)
	fmt.Printf("Model: %s\n", cfg.ModelName)
	return 0
}

func providerTest(name string, cfg *config.Config, kc *keychain.Keychain) int {
	// Validate provider name
	providerInfo, ok := getProviderInfo(name)
	if !ok {
		fmt.Printf("Unknown provider: %s\n", name)
		return 1
	}

	fmt.Printf("Testing connection to %s...\n", providerInfo.DisplayName)

	// Check if configured
	if name != "ollama" {
		if key, err := kc.GetProviderKey(name); err != nil || key == "" {
			fmt.Printf("❌ Provider '%s' is not configured (no API key found).\n", name)
			fmt.Printf("   Run 'pryx-core provider set-key %s' to configure.\n", name)
			return 1
		}
	}

	// For ollama, check if endpoint is reachable
	if name == "ollama" {
		// Simple HTTP check to ollama endpoint
		endpoint := cfg.OllamaEndpoint
		if endpoint == "" {
			endpoint = "http://localhost:11434"
		}
		fmt.Printf("Checking Ollama at %s...\n", endpoint)
		// Note: In a real implementation, we'd make an HTTP request here
		fmt.Println("✓ Ollama endpoint configured (manual verification needed)")
		return 0
	}

	// For cloud providers, we'd make a test API call
	// For now, just verify key exists
	fmt.Printf("✓ API key found for %s\n", providerInfo.DisplayName)
	fmt.Println("✓ Provider configuration valid")
	fmt.Println("")
	fmt.Println("Note: Full connection test will be implemented with LLM client integration.")

	return 0
}

func getProviderInfo(name string) (ProviderInfo, bool) {
	name = strings.ToLower(name)
	for _, p := range SupportedProviders {
		if p.Name == name {
			return p, true
		}
	}
	return ProviderInfo{}, false
}

func getDefaultModelForProvider(provider string) string {
	switch provider {
	case "openai":
		return "gpt-4"
	case "anthropic":
		return "claude-3-opus"
	case "glm":
		return "glm-4-flash"
	case "openrouter":
		return "anthropic/claude-3-opus"
	case "together":
		return "meta-llama/Llama-3-70b-chat-hf"
	case "groq":
		return "llama3-70b-8192"
	case "xai":
		return "grok-beta"
	case "mistral":
		return "mistral-large"
	case "cohere":
		return "command-r-plus"
	case "google":
		return "gemini-pro"
	case "ollama":
		return "llama3"
	default:
		return ""
	}
}
