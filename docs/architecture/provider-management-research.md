# Provider Management Research: opencode/moltbot Patterns

## Overview

Research into how opencode (and moltbot) handle provider management in their CLI/TUI applications.

---

## OpenCode Provider Management Patterns

### 1. Command Structure

OpenCode uses a dedicated command for provider management:

```typescript
// From i18n translation files
"command.category.provider": "Provider",
"command.provider.connect": "Connect Provider",
```

**Commands Available**:
- `/connect` - Connect a new provider
- Provider-specific connection flows
- Disconnect providers
- View connected providers

---

### 2. Provider Connection Flow

OpenCode implements a sophisticated multi-step connection flow:

#### A. Provider Selection Dialog
```typescript
"dialog.provider.search.placeholder": "Search providers",
"dialog.provider.empty": "No providers found",
"dialog.provider.group.popular": "Popular",
"dialog.provider.group.other": "Others",
"dialog.provider.tag.recommended": "Recommended",
```

**Features**:
- Searchable provider list
- Grouped by popularity
- Tags for recommended providers
- Provider notes/descriptions

#### B. Provider Categories
OpenCode categorizes providers for better UX:

**Popular Providers**:
- Anthropic (Claude)
- OpenAI (GPT)
- Google (Gemini)
- xAI (Grok)
- Meta (Llama)

**Connection Methods**:
- Claude Pro/Max OAuth
- ChatGPT Pro/Plus OAuth
- Copilot OAuth
- API Key (universal)

#### C. Connection Methods

**1. API Key Method**:
```typescript
"provider.connect.method.apiKey": "API Key",
"provider.connect.apiKey.description": "Enter your {{provider}} API key to connect your account",
"provider.connect.apiKey.label": "{{provider}} API Key",
"provider.connect.apiKey.placeholder": "API Key",
"provider.connect.apiKey.required": "API key is required",
```

**2. OAuth Methods**:
```typescript
"provider.connect.oauth.code.visit.link": "this link",
"provider.connect.oauth.code.label": "{{method}} Authorization Code",
"provider.connect.oauth.auto.confirmationCode": "Confirmation Code",
```

**OAuth Flows**:
- Auto OAuth (redirect-based)
- Manual OAuth (code-based)
- Provider-specific OAuth (Claude Pro/Max)

#### D. Connection Status States

```typescript
"provider.connect.status.inProgress": "Authorizing...",
"provider.connect.status.waiting": "Waiting for authorization...",
"provider.connect.status.failed": "Authorization failed: {{error}}",
```

---

### 3. Settings Integration

OpenCode has a dedicated Providers section in Settings:

```typescript
"settings.providers.title": "Providers",
"settings.providers.description": "Provider settings can be configured here",
"settings.providers.section.connected": "Connected Providers",
"settings.providers.connected.empty": "No connected providers",
"settings.providers.section.popular": "Popular Providers",
"settings.providers.tag.environment": "Environment",
"settings.providers.tag.config": "Config",
"settings.providers.tag.custom": "Custom",
```

**Settings Features**:
- View all connected providers
- Disconnect providers
- See provider source (environment, config, custom)
- Provider tags for organization

---

### 4. Provider-Specific Features

#### Provider Notes
OpenCode shows provider-specific notes:

```typescript
"dialog.provider.anthropic.note": "Connect with Claude Pro/Max or API key",
"dialog.provider.openai.note": "Connect with ChatGPT Pro/Plus or API key",
"dialog.provider.copilot.note": "Connect with Copilot or API key",
"dialog.provider.opencode.note": "Curated models including Claude, GPT, Gemini",
```

#### Special Providers

**OpenCode Zen** (Unified API):
```typescript
"provider.connect.opencodeZen.line2": "Access models like Claude, GPT, Gemini with a single API key",
"provider.connect.opencodeZen.visit.link": "opencode.ai/zen",
```

---

### 5. Error Handling

OpenCode provides detailed error messages:

```typescript
"error.chain.modelNotFound": "Model not found: {{provider}}/{{model}}",
"error.chain.providerAuthFailed": "Provider authentication failed ({{provider}}): {{message}}",
"error.chain.providerInitFailed": 'Failed to initialize provider "{{provider}}". Check credentials and configuration.',
```

---

### 6. Toast Notifications

Success/Failure feedback:

```typescript
"provider.connect.toast.connected.title": "{{provider}} Connected",
"provider.connect.toast.connected.description": "{{provider}} models are now available",
"provider.disconnect.toast.disconnected.title": "{{provider}} Disconnected",
"provider.disconnect.toast.disconnected.description": "{{provider}} models are no longer available",
```

---

## Moltbot Provider Management Patterns

### 1. CLI-First Approach

Moltbot uses CLI commands for provider management:

```typescript
// From moltbot CLI
.option("--no-usage", "Skip model provider usage/quota snapshots")
.description("Show provider capabilities (intents/scopes + supported features)")
```

**CLI Commands**:
- `moltbot providers list` - List available providers
- `moltbot providers connect <provider>` - Connect a provider
- `moltbot providers capabilities` - Show provider capabilities

---

### 2. Provider Configuration Schema

Moltbot uses a comprehensive provider configuration:

```typescript
// Provider configuration structure
providers: {
  [providerID: string]: {
    models: Array<{
      id: string;
      name: string;
      capabilities: string[];
    }>;
    auth: {
      type: "apiKey" | "oauth";
      required: boolean;
    };
  }
}
```

---

### 3. Connection Status Tracking

Moltbot tracks connection status for distributed systems:

```typescript
// Node connection status
interface NodeStatus {
  connected: boolean;
  connectedAtMs?: number;
  reconnectAttempts: number;
  lastDisconnect?: {
    code: number;
    reason: string;
    at: number;
  };
}
```

---

### 4. Provider Capabilities

Moltbot exposes provider capabilities:

```typescript
// Provider capability detection
function getProviderCapabilities(provider: string): {
  intents: string[];
  scopes: string[];
  features: string[];
  supportsStreaming: boolean;
  supportsVision: boolean;
  supportsTools: boolean;
}
```

---

## Key Patterns Summary

### 1. **Multi-Method Authentication**
- Support both API Key and OAuth
- Provider-specific OAuth flows
- Fallback authentication methods

### 2. **Progressive Disclosure**
- Searchable provider list
- Group providers (popular vs others)
- Show details only when needed

### 3. **Visual Feedback**
- Connection status indicators
- Toast notifications
- Error messages with context

### 4. **Settings Integration**
- Centralized provider management
- Tags for organization
- Easy disconnect/reconnect

### 5. **Error Handling**
- Specific error messages per provider
- Suggestions for fixing issues
- Clear status states

---

## Recommendations for Pryx

### 1. Implement `/connect` Command

Add a TUI command palette command:

```typescript
{
  id: "connect-provider",
  name: "Connect Provider",
  description: "Add a new LLM provider",
  category: "Provider",
  shortcut: "c",
  action: () => openProviderDialog(),
}
```

### 2. Provider Selection Dialog

Create a searchable list with:
- Provider icons
- Connection status
- Popular providers highlighted
- Search by name

### 3. Connection Methods

Support multiple auth methods:
- **API Key**: Direct input with validation
- **OAuth**: Browser-based flow with callback
- **Environment**: Read from env vars

### 4. Provider Status in Header

Show connection status in TUI header:
```
[OpenAI: ✓] [Anthropic: ✗] [GLM: ✓]
```

### 5. Settings View Enhancement

Add to Settings view:
- Connected providers list
- Provider health check
- Model selection per provider
- Disconnect button

### 6. Error Handling

Implement user-friendly errors:
- "Failed to connect to OpenAI: Invalid API key"
- "Anthropic: Rate limited, retrying..."
- "GLM: Connection timeout"

---

## Implementation Priority

### Phase 1: Basic Connection
1. `/connect` command
2. Provider list (from catalog)
3. API key input
4. Connection test

### Phase 2: Enhanced UX
1. OAuth support
2. Provider status indicators
3. Settings integration
4. Error handling

### Phase 3: Advanced Features
1. Provider health monitoring
2. Auto-reconnect
3. Model selection UI
4. Provider comparison

---

## UI Mockup

```
┌─────────────────────────────────────────────────────────┐
│  Connect Provider                    [Search: ______]   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ★ Popular Providers                                    │
│  ┌─────────────────────────────────────────────────┐   │
│  │ 🤖 OpenAI                    [Connect]          │   │
│  │    GPT-4, GPT-3.5, DALL-E                       │   │
│  ├─────────────────────────────────────────────────┤   │
│  │ 🧠 Anthropic                 [Connected ✓]      │   │
│  │    Claude 3.5 Sonnet, Claude 3 Opus             │   │
│  ├─────────────────────────────────────────────────┤   │
│  │ 🔍 Google                    [Connect]          │   │
│  │    Gemini Pro, Gemini Ultra                     │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  All Providers (84)                                     │
│  ┌─────────────────────────────────────────────────┐   │
│  │ • Alibaba Cloud              [Connect]          │   │
│  │ • AWS Bedrock                [Connect]          │   │
│  │ • Azure OpenAI               [Connect]          │   │
│  │ • ...                                             │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## Conclusion

Both opencode and moltbot provide excellent patterns for provider management:

- **OpenCode**: Focuses on user-friendly OAuth flows and visual feedback
- **Moltbot**: Emphasizes CLI commands and distributed system monitoring

**For Pryx**: Combine OpenCode's UX patterns with moltbot's CLI approach for a hybrid TUI/CLI experience.
