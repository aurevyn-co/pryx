# LLM Provider Configuration

This guide documents how to configure and manage AI providers in Pryx.

## Overview

Pryx supports 84+ AI providers via [models.dev](https://models.dev). Providers are configured through the runtime API or TUI.

## Supported Providers

### OpenAI
```json
{
  "id": "openai",
  "name": "OpenAI",
  "models": ["gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"],
  "auth_type": "api_key"
}
```

### Anthropic
```json
{
  "id": "anthropic",
  "name": "Anthropic",
  "models": ["claude-3-opus", "claude-2", "claude-instant-1"],
  "auth_type": "api_key"
}
```

### Google AI
```json
{
  "id": "google",
  "name": "Google AI",
  "models": ["gemini-pro", "gemini-1.5-pro", "gemini-1.5-flash"],
  "auth_type": "oauth2"
}
```

### xAI
```json
{
  "id": "xai",
  "name": "xAI",
  "models": ["grok-beta", "grok-beta-vision"],
  "auth_type": "api_key"
}
```

### OpenRouter
```json
{
  "id": "openrouter",
  "name": "OpenRouter",
  "models": "dynamic",
  "auth_type": "api_key"
}
```

### Ollama
```json
{
  "id": "ollama",
  "name": "Ollama",
  "models": "dynamic",
  "auth_type": "none",
  "base_url": "http://localhost:11434"
}
```

### DeepSeek
```json
{
  "id": "deepseek",
  "name": "DeepSeek",
  "models": ["deepseek-chat", "deepseek-coder"],
  "auth_type": "api_key"
}
```

## Configuration Methods

### API Key Providers

For providers using API keys (OpenAI, Anthropic, xAI, DeepSeek, etc.):

```bash
# Set API key via CLI
pryx provider add openai
pryx provider set-key openai
# Enter your API key when prompted

# Or via TUI
# Press `/` → "Providers" → Select provider → Set Key
# Or via HTTP API
POST /api/v1/providers/openai/key
{
  "key": "sk-..."
}
```

### OAuth2 Providers

For providers using OAuth2 (Google AI):

```bash
# Start OAuth flow via CLI
pryx provider add google
pryx provider oauth google
# Follow browser prompts

# Or via TUI
# Press `/` → "Providers" → Select provider → OAuth
```

OAuth2 Flow:
1. User initiates OAuth via CLI/TUI
2. Runtime starts local callback server
3. Provider displays authorization URL
4. User approves in browser
5. Runtime receives callback and exchanges auth code
6. Runtime stores tokens securely in vault/keychain
7. Tokens available to LLM orchestration components

### Local LLM Providers

For providers running locally (Ollama):

```bash
# Add Ollama provider
pryx provider add ollama

# Configure Ollama endpoint
pryx provider set-config ollama --base-url http://localhost:11434
```

## API Endpoints

```
GET    /api/v1/providers                  # List configured providers
POST   /api/v1/providers/{id}            # Add new provider
PUT    /api/v1/providers/{id}            # Update provider config
DELETE /api/v1/providers/{id}            # Remove provider
GET    /api/v1/providers/{id}/models       # Get models for provider
POST   /api/v1/providers/{id}/key          # Set API key
GET    /api/v1/providers/{id}/key          # Check if key is set
DELETE /api/v1/providers/{id}/key         # Delete API key
POST   /api/v1/providers/{id}/oauth        # Start OAuth flow
```

## Provider Configuration Schema

```json
{
  "id": "openai",
  "name": "My OpenAI",
  "enabled": true,
  "api_key_ref": "vault:openai",
  "model": "gpt-4",
  "base_url": "https://api.openai.com/v1",
  "created_at": "2026-02-04T12:00:00Z",
  "updated_at": "2026-02-04T12:00:00Z"
}
```

## Security

- **API Keys**: Stored encrypted in vault/keychain with scope-based access
- **OAuth Tokens**: Stored in vault with encrypted refresh tokens
- **Key Rotation**: Manual trigger via CLI/API
- **Audit Logging**: All key access logged in audit system
- **Scope Validation**: Keys only available to authorized components

## Model Management

```bash
# List available models
pryx provider models openai

# Set default model
pryx provider use-model openai gpt-4

# Get current model
pryx provider current-model
```

## Environment Variables

- `PRYX_PROVIDER`: Current active provider
- `PRYX_MODEL`: Current active model
- `PRYX_API_KEY_REF`: Vault reference to API key

## Testing

```bash
# Test provider connection
pryx provider test openai

# List models
pryx provider models openai
```

## Troubleshooting

### Provider not connecting
1. Check internet connection
2. Verify API key validity
3. Check provider status: `pryx provider status`
4. View audit logs: `pryx audit --provider`

### OAuth flow not completing
1. Verify callback server is running
2. Check firewall settings
3. Ensure redirect URI matches configured URL
4. Verify client ID in provider config

### API key errors
1. Verify key format for provider
2. Check key permissions in provider console
3. Contact provider support if issue persists
