# PRD: AI-Assisted Setup & Configuration Architecture

> **Version**: 1.0  
> **Status**: Draft - For Review  
> **Last Updated**: 2026-01-29  
> **Parent**: `docs/prd/prd.md` Section 7 (Product Surfaces)  

---

## 1) Problem Statement: The Configuration Paradox

### The Challenge
Users have different comfort levels with configuration:

| User Type | Preference | Pain Point |
|-----------|------------|------------|
| **Power User** | Manual config files, CLI | Wants precision, version control, automation |
| **Visual User** | TUI/Web UI with forms | Wants discoverability, validation, previews |
| **Natural Language User** | AI-assisted setup | Wants to say "connect my Telegram" and have it work |

**Risk**: If AI-assisted setup is the ONLY path, we alienate power users. If it's NOT available, we lose casual users.

### The Solution: Three Coexisting Configuration Paths

All three methods must work **independently** and **respect each other's configurations**.

---

## 2) Configuration Method Architecture

### 2.1 The Three Paths

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Configuration Entry Points                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌──────────────┐    ┌──────────────┐    ┌──────────────┐         │
│   │   MANUAL     │    │    VISUAL    │    │ AI-ASSISTED  │         │
│   │              │    │              │    │              │         │
│   │ • Config     │    │ • TUI Forms  │    │ • Natural    │         │
│   │   files      │    │ • Web UI     │    │   language   │         │
│   │ • CLI flags  │    │ • Wizards    │    │ • Contextual │         │
│   │ • Env vars   │    │ • Validation │    │   help       │         │
│   └──────┬───────┘    └──────┬───────┘    └──────┬───────┘         │
│          │                   │                   │                  │
│          └───────────────────┴───────────────────┘                  │
│                              │                                      │
│                              ▼                                      │
│   ┌──────────────────────────────────────────────────────────┐     │
│   │              Unified Configuration Store                  │     │
│   │                                                           │     │
│   │  ~/.pryx/config.json (JSON5)                             │     │
│   │  ~/.pryx/vault/credentials.json (encrypted)              │     │
│   │  ~/.pryx/mcp/servers.json                                │     │
│   │  ~/.pryx/channels/telegram.json                          │     │
│   │                                                           │     │
│   └──────────────────────────────────────────────────────────┘     │
│                              │                                      │
│                              ▼                                      │
│   ┌──────────────────────────────────────────────────────────┐     │
│   │            Validation & Conflict Resolution               │     │
│   │                                                           │     │
│   │  • Schema validation (Zod)                               │     │
│   │  • Conflict detection (last-write-wins with backup)      │     │
│   │  • Audit logging (who changed what, when)                │     │
│   │  • Rollback capability (keep last 5 versions)            │     │
│   │                                                           │     │
│   └──────────────────────────────────────────────────────────┘     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.2 Core Principle: No Method is Superior

**CRITICAL**: All three methods are **first-class citizens**. None is a fallback or secondary option.

| Method | Best For | Not For |
|--------|----------|---------|
| **Manual** | CI/CD, version control, power users who know exactly what they want | First-time users |
| **Visual** | Discovering options, validation, seeing current state | Headless servers, automation |
| **AI-Assisted** | Quick setup, natural language preference, getting unstuck | Users who want explicit control over every detail |

### 2.3 Conflict Resolution Rules

When the same configuration is modified by different methods:

1. **Last-Write-Wins** (default): Most recent change takes effect
2. **User Confirmation** (for destructive changes): "AI wants to change your model provider from Claude to GPT-4. Approve?"
3. **Backup Always**: Keep previous version before any change
4. **Audit Trail**: Log who/what made the change

---

## 3) AI-Assisted Flows: Design Specification

### 3.1 When AI Assistance is Triggered

**Explicit triggers** (user initiates):
```
User: "Help me set up Telegram"
User: "Connect to OpenAI"
User: "What integrations are available?"
```

**Implicit triggers** (AI suggests):
```
User: "Send this to my team"
AI: "I notice you don't have any channels set up. Would you like help connecting Slack or Telegram?"
```

**NOT triggered when**:
- User is editing config files directly
- User is in TUI settings menu
- User explicitly says "manual mode" or "show me the config"

### 3.2 AI-Assisted Flow UX Pattern

```
User: "Connect my Telegram bot"

AI:
┌─────────────────────────────────────────────────────────┐
│ I'll help you connect Telegram! You have 3 options:     │
│                                                         │
│ [AI-ASSISTED - RECOMMENDED]                             │
│ I can walk you through getting a bot token from         │
│ @BotFather and connect it automatically.                │
│                                                         │
│ [MANUAL OPTIONS]                                        │
│ • TUI: Run `pryx integrations add telegram`             │
│ • Config: Edit ~/.pryx/channels/telegram.json           │
│ • CLI: pryx channel telegram --token <token>            │
│                                                         │
│ Which would you prefer? (Reply with 1, 2, or 3)         │
└─────────────────────────────────────────────────────────┘
```

### 3.3 AI Must Always Offer Manual Escape Hatch

**Rule**: Every AI-assisted flow must provide a manual alternative.

**Example - During AI Setup**:
```
AI: "Great! I see you have a Telegram bot token. I'll connect it now."
      
      [User can type at any time:]
      User: "wait, show me the config instead"
      
      AI: "Sure! Here's what I was about to save:
      
      ~/.pryx/channels/telegram.json:
      {
        "enabled": true,
        "botToken": "[REDACTED - stored in vault]",
        "allowFrom": [],
        "webhook": {
          "url": "https://webhook.pryx.dev/...",
          "secret": "[REDACTED]"
        }
      }
      
      Would you like me to:
      1. Save this configuration
      2. Show you how to edit it manually
      3. Cancel and do nothing"
```

### 3.4 AI-Assisted Flows Must Respect Existing Config

**Scenario 1**: User has manually configured Telegram
```
User: "Connect Telegram"
AI: "I see Telegram is already configured! Current status:
     • Bot: @mybot (connected)
     • Webhook: Active
     • Channels: 2 linked chats
     
     Would you like to:
     1. Test the connection
     2. View/edit configuration
     3. Add another bot"
```

**Scenario 2**: AI-assisted setup partially completed, user switches to manual
```
User: [halfway through AI Telegram setup]
User: "Actually, let me do this in the TUI"
AI: "No problem! I'll save your progress. You can continue in the TUI:
     `pryx integrations telegram --continue-setup`
     
     Or start fresh from scratch there."
```

---

## 4) Implementation: The 3-Phase Plan

### Phase 1: Foundation (Weeks 1-3)

**Goal**: Establish configuration infrastructure that supports all 3 methods

| Task | Component | Priority |
|------|-----------|----------|
| Vault encryption core (Argon2 + AES-256-GCM) | Host | P0 |
| Configuration schemas (Zod validation) | Runtime | P0 |
| Unified config store (JSON5 + encryption) | Runtime | P0 |
| Config file watcher (hot reload) | Runtime | P0 |
| Audit logging system | Runtime | P0 |
| Backup/rollback (keep last 5 versions) | Runtime | P0 |
| Provider config schema | Runtime | P1 |
| Channel config schema | Runtime | P1 |
| MCP config schema | Runtime | P1 |

**Deliverables**:
- All config changes are validated
- All config changes are logged
- All config changes create backups
- Manual editing of config files works perfectly

### Phase 2: Feature Implementation (Weeks 4-8)

**Goal**: Build the three configuration paths in parallel

#### Track A: Manual/CLI (Power Users)

| Task | Component | Priority |
|------|-----------|----------|
| CLI commands: `pryx config get/set` | CLI | P1 |
| CLI commands: `pryx provider add/remove` | CLI | P1 |
| CLI commands: `pryx channel add/remove` | CLI | P1 |
| CLI commands: `pryx mcp add/remove` | CLI | P1 |
| CLI commands: `pryx vault add/remove` | CLI | P1 |
| Environment variable support | Runtime | P1 |
| Config validation command | CLI | P1 |
| Config export/import | CLI | P2 |

#### Track B: Visual/TUI (Visual Users)

| Task | Component | Priority |
|------|-----------|----------|
| TUI Provider management screen | TUI | P1 |
| TUI Channel configuration forms | TUI | P1 |
| TUI MCP server browser | TUI | P1 |
| TUI Vault credential manager | TUI | P1 |
| Form validation with helpful errors | TUI | P1 |
| Connection testing UI | TUI | P1 |
| Web UI equivalent screens | Web | P2 |

#### Track C: AI-Assisted (Natural Language Users)

| Task | Component | Priority |
|------|-----------|----------|
| Natural language intent parser | Runtime | P1 |
| Provider setup dialogue flows | Runtime | P1 |
| Channel setup dialogue flows | Runtime | P1 |
| MCP setup dialogue flows | Runtime | P1 |
| Contextual help system | Runtime | P1 |
| Manual escape hatch ("show me the config") | Runtime | P1 |
| Progress persistence (save/resume) | Runtime | P2 |
| Setup verification & testing | Runtime | P2 |

**Deliverables**:
- All three methods work independently
- Changes from any method are visible to all others
- No conflicts or data loss between methods

### Phase 3: Integration & Polish (Weeks 9-12)

**Goal**: Seamless integration between methods, advanced features

| Task | Component | Priority |
|------|-----------|----------|
| Natural language cron parser | Runtime | P2 |
| Cron scheduler service | Runtime | P2 |
| Cron TUI management | TUI | P2 |
| Security audit command (`pryx doctor`) | CLI | P1 |
| Vault TUI with secure input | TUI | P1 |
| Environment variable import | CLI | P2 |
| AI-guided setup for all integrations | Runtime | P2 |
| Setup templates & recommendations | Runtime | P2 |
| Conflict resolution UI | TUI | P2 |
| Configuration sync validation | Runtime | P2 |

**Deliverables**:
- Users can switch between methods seamlessly
- AI suggests but never forces
- All edge cases handled (half-done setups, conflicts, rollbacks)

---

## 5) User Scenarios: All Three Methods in Action

### Scenario 1: Power User (Manual Method)

```
# Day 1: Initial setup via config file
$ cat ~/.pryx/config.json
{
  "providers": {
    "anthropic": {
      "apiKey": "${ANTHROPIC_API_KEY}",  // env var reference
      "defaultModel": "claude-3-opus"
    }
  }
}

$ pryx config validate
✓ Config is valid

# Day 2: Add Telegram via CLI
$ pryx channel add telegram --token $BOT_TOKEN
✓ Telegram bot @mybot connected

# Day 3: Check status
$ pryx status
Providers: anthropic ✓
Channels: telegram ✓
MCP Servers: filesystem ✓, shell ✓
```

### Scenario 2: Visual User (TUI Method)

```
$ pryx
[Opens TUI]

┌─────────────────────────────────────────┐
│  🔧 Pryx Control Center                │
│                                         │
│  [Settings] [Integrations] [Sessions]   │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  Integrations                   │   │
│  │                                 │   │
│  │  AI Providers:               │   │
│  │    ○ Anthropic [Connect]      │   │
│  │    ○ OpenAI [Connect]         │   │
│  │    ● Local (Ollama) [Configure]│   │
│  │                                 │   │
│  │  Channels:                   │   │
│  │    ● Telegram ✓               │   │
│  │    ○ Discord [Connect]        │   │
│  │    ○ Slack [Connect]          │   │
│  │                                 │   │
│  └─────────────────────────────────┘   │
│                                         │
└─────────────────────────────────────────┘

[User clicks "Connect" on Discord]
[Form opens: Bot Token, Permissions, Webhook URL]
[Validation happens in real-time]
[Test Connection button]
```

### Scenario 3: Natural Language User (AI-Assisted Method)

```
$ pryx

Pryx: Hello! How can I help you today?

User: I want to connect my Telegram bot

Pryx: I'll help you connect Telegram! I see a few options:

      1. 🚀 Quick Setup (Recommended)
         I'll guide you through getting a bot token and connect it.
      
      2. 📋 Step-by-Step
         Show me instructions and I'll help if you get stuck.
      
      3. ⚙️ Manual Configuration
         Open the configuration form to enter details yourself.
      
      Which would you prefer? (1, 2, or 3)

User: 1

Pryx: Great! First, do you already have a bot token from @BotFather?
      
      (If not, I can walk you through creating one)

User: yes, it's 123456:ABC-DEF...

Pryx: Perfect! I'll verify this token and set up the connection.
      
      [Testing connection...]
      ✓ Bot token is valid
      ✓ Bot name: @myawesomebot
      
      Would you like me to:
      • Set up webhook delivery (recommended)
      • Use polling mode (for local-only setups)
      
      Or type "show config" to see exactly what I'll save.

User: webhook is fine

Pryx: ✓ All set! Your Telegram bot is connected.
      
      Try it: Send /start to @myawesomebot on Telegram.
      
      You can manage this anytime via:
      • TUI: pryx → Integrations → Telegram
      • CLI: pryx channel telegram status
      • Config: ~/.pryx/channels/telegram.json
```

### Scenario 4: Mixed Method (Real-World)

```
# User starts with AI-assisted
User: "Connect to OpenAI"
AI: [guides through API key input]

# User switches to TUI mid-setup
User: "Actually, show me the TUI"
AI: "Opening TUI with your OpenAI settings..."
[TUI opens with partial OpenAI config filled in]

# User edits in TUI, saves
[TUI: User changes model from gpt-4 to gpt-3.5-turbo]

# Later, user checks via CLI
$ pryx config get providers.openai.model
gpt-3.5-turbo

# User edits via config file
$ vim ~/.pryx/config.json
[changes model to claude-3-sonnet via anthropic]

# AI recognizes the change
User: "What model am I using?"
AI: "You're currently using Claude 3 Sonnet via Anthropic. 
      (I see you switched from OpenAI in your config - good choice for this task!)"
```

---

## 6) Anti-Patterns to Avoid

### ❌ Anti-Pattern 1: AI as Gatekeeper
```
# BAD: AI prevents manual access
User: "Show me the Telegram config"
AI: "I can help you with that! Let me walk you through it..."
      [Forces AI-guided flow, no manual option]
```

**Fix**: Always offer manual alternative immediately.

### ❌ Anti-Pattern 2: AI Overrides Without Consent
```
# BAD: AI changes config without asking
User: "My Telegram isn't working"
AI: [automatically resets webhook and changes config]
      "I fixed it!"
```

**Fix**: Always ask before changing existing configuration.

### ❌ Anti-Pattern 3: Inconsistent State
```
# BAD: TUI shows different state than reality
TUI shows: Telegram disconnected
Config file: Telegram enabled and working
CLI shows: Telegram connected
```

**Fix**: Single source of truth (config files), all methods read from same source.

### ❌ Anti-Pattern 4: AI-Only Features
```
# BAD: Feature only accessible via AI
User: "How do I set up cron jobs?"
AI: "Just tell me what you want to schedule!"
      [No documentation on manual cron syntax]
```

**Fix**: All features must be accessible via all three methods.

---

## 7) Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Method usage distribution** | Roughly equal 33/33/33% | Track which method users choose for new setups |
| **Switch rate** | <10% need to switch methods | Users who start with AI but switch to manual |
| **Conflict incidents** | <1% of config changes | User complaints about AI overriding their settings |
| **Setup completion rate** | >90% for all methods | Successful integration setup via each method |
| **Time to setup** | AI <3 min, TUI <5 min, Manual <10 min | Average time to complete common setups |

---

## 8) Integration with Main PRD

### Updates Required to `docs/prd/prd.md`:

**Section 7.1 (CLI & TUI)** - Add:
```
All CLI commands have equivalent TUI and AI-assisted paths:
- Manual: `pryx mcp add <name> --url <url>`
- Visual: TUI → Integrations → MCP Servers → Add
- AI-Assisted: "Connect MCP server at <url>"
```

**Section 7.2 (Web UI)** - Add:
```
Configuration methods coexist:
- Users can switch between manual, visual, and AI-assisted modes at any time
- No method is privileged over another
- Changes from any method are immediately visible to all others
```

**Section 9 (Functional Requirements)** - Add:
```
FR-X: Configuration Method Parity
- All configuration operations must be available via:
  1. Manual (config files + CLI)
  2. Visual (TUI + Web UI forms)
  3. AI-Assisted (natural language)
- Users can switch methods mid-flow
- All methods share single source of truth
```

---

## 9) Summary

**Key Principles**:

1. **Three First-Class Methods**: Manual, Visual, AI-Assisted all equal
2. **No Gatekeeping**: AI assists but never blocks manual access
3. **Single Source of Truth**: Config files are canonical
4. **Seamless Switching**: Users can change methods anytime
5. **Respect User Choice**: Never override without explicit consent

**Implementation Order**:
1. **Phase 1**: Foundation (config store, validation, audit)
2. **Phase 2**: Parallel tracks (Manual, Visual, AI-Assisted)
3. **Phase 3**: Integration (seamless switching, advanced features)

This ensures **AI-assisted flows enhance the experience without replacing manual control**, addressing your concern about conflicts.

---

*Document Status*: Ready for review and integration into main PRD
