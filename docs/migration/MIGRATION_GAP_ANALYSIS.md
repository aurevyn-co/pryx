# Pryx → Rust Migration Gap Analysis

> **Version**: 1.0-draft  
> **Date**: 2026-02-15  
> **Status**: Draft - Work in Progress  

---

## Executive Summary

This document analyzes the gaps between the current Pryx polyglot codebase (Rust/Tauri + Go/TypeScript) and the target Rust-based ZeroClaw codebase, with additional feature insights from OpenClaw and OpenCode.

**Key Findings:**
- ZeroClaw provides a solid foundation with 22+ providers, 8 channels, and trait-based architecture
- Pryx has unique features (Mesh, TUI, Desktop App, Local Web) that need migration planning
- OpenCode's sqlite session management and error patterns should be adopted
- OpenClaw's extensive channel ecosystem and multi-agent routing are reference implementations

---

## 1. Architecture Comparison

### 1.1 Current Pryx Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Desktop Host (Rust + Tauri)                   │
│  Port: 42424                                                         │
│  • HTTP server (axum)                                                │
│  • WebSocket for real-time TUI communication                         │
│  • Local web UI admin panel (apps/local-web/)                        │
│  • Go runtime sidecar management                                    │
│  • Native dialogs & system tray                                      │
└─────────────────────────────────────────────────────────────────────┘
                              │
                    Sidecar (Go Runtime)
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Runtime (Go)                                     │
│  • Agent execution & orchestration                                   │
│  • HTTP API + WebSocket server                                      │
│  • 84+ AI providers (models.dev)                                    │
│  • Channels (Telegram, Discord, Slack, Webhooks)                     │
│  • MCP integration                                                   │
│  • Memory & RAG                                                      │
│  • Vault (Argon2id encryption)                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 1.2 Target ZeroClaw Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        ZeroClaw (100% Rust)                          │
│  ~3.4MB binary · <10ms startup · 1,017 tests                        │
├─────────────────────────────────────────────────────────────────────┤
│  • Every subsystem is a TRAIT - pluggable architecture               │
│  • 22+ AI providers (Provider trait)                                │
│  • 8 channels (Channel trait)                                       │
│  • SQLite + Vector memory (Memory trait)                            │
│  • Shell, File, Browser tools (Tool trait)                          │
│  • Security: pairing, sandbox, allowlists (SecurityPolicy trait)    │
│  • Gateway/daemon mode                                              │
│  • Skills system via TOML manifests                                 │
│  • 50+ integrations registry                                        │
└─────────────────────────────────────────────────────────────────────┘
```

### 1.3 Key Architectural Differences

| Aspect | Pryx (Current) | ZeroClaw (Target) | Gap/Action |
|--------|----------------|-------------------|------------|
| **Language** | Polyglot (Rust/Go/TS) | 100% Rust | ✅ Rust migration complete |
| **Runtime** | Tauri sidecar + Go | Single binary | ✅ Simpler deployment |
| **UI** | TUI (Solid/TS) + Local Web (React) | CLI only | ⚠️ Need to add TUI/Web |
| **Providers** | 84+ via models.dev | 22+ built-in | ⚠️ Need models.dev integration |
| **Channels** | 4 built-in | 8 built-in | ✅ More channels available |
| **Memory** | SQLite + RAG | SQLite + FTS5 + Vector | ✅ Better memory system |
| **MCP** | Native client | Skills/Integrations | ⚠️ Need MCP adapter |

---

## 2. Feature Gap Analysis

### 2.1 Features in ZeroClaw (Target) ✅

| Feature | Implementation | Status |
|---------|---------------|--------|
| Provider trait system | `src/providers/traits.rs` | ✅ Ready |
| 22+ AI providers | OpenAI, Anthropic, OpenRouter, Ollama, etc. | ✅ Ready |
| Channel trait system | `src/channels/traits.rs` | ✅ Ready |
| CLI channel | `src/channels/cli.rs` | ✅ Ready |
| Telegram channel | `src/channels/telegram.rs` | ✅ Ready |
| Discord channel | `src/channels/discord.rs` | ✅ Ready |
| Slack channel | `src/channels/slack.rs` | ✅ Ready |
| Matrix channel | `src/channels/matrix.rs` | ✅ Ready |
| WhatsApp channel | `src/channels/whatsapp.rs` | ✅ Ready |
| iMessage channel | `src/channels/imessage.rs` | ✅ Ready |
| SQLite memory | `src/memory/sqlite.rs` | ✅ Ready |
| Vector search | `src/memory/vector.rs` | ✅ Ready |
| FTS5 keyword search | `src/memory/sqlite.rs` | ✅ Ready |
| Shell tool | `src/tools/shell.rs` | ✅ Ready |
| File read/write | `src/tools/file_*.rs` | ✅ Ready |
| Memory tools | `src/tools/memory_*.rs` | ✅ Ready |
| Browser tool | `src/tools/browser.rs` | ✅ Ready |
| Security pairing | `src/security/pairing.rs` | ✅ Ready |
| Security policy | `src/security/policy.rs` | ✅ Ready |
| Gateway server | `src/gateway/mod.rs` | ✅ Ready |
| Daemon mode | `src/daemon/mod.rs` | ✅ Ready |
| Cron scheduler | `src/cron/scheduler.rs` | ✅ Ready |
| Skills system | `src/skills/mod.rs` | ✅ Ready |
| Integrations | `src/integrations/registry.rs` | ✅ Ready |
| Tunnels | `src/tunnel/*.rs` | ✅ Ready |
| AIEOS identity | `src/config/schema.rs` | ✅ Ready |
| Doctor command | `src/doctor/mod.rs` | ✅ Ready |

### 2.2 Features in Pryx (Missing in ZeroClaw) ❌

| Feature | Pryx Location | Gap Analysis | Priority |
|---------|--------------|--------------|----------|
| **TUI (Terminal UI)** | `apps/tui/` | ZeroClaw has CLI only | HIGH |
| **Local Web Admin** | `apps/local-web/` | ZeroClaw has gateway API only | HIGH |
| **Desktop App (Tauri)** | `apps/desktop/` | ZeroClaw is CLI-first | MEDIUM |
| **Pryx Mesh** | `pkg/mesh/` | Multi-device sync | HIGH |
| **84+ Providers** | Go runtime + models.dev | ZeroClaw has 22+ | MEDIUM |
| **MCP Native Client** | Go runtime | ZeroClaw has Skills/Integrations | HIGH |
| **Vault System** | Go runtime | ZeroClaw has encrypted secrets | LOW |
| **Cost Tracking** | Go runtime | Need observability trait | MEDIUM |
| **Session Timeline** | Go runtime | Need session management | MEDIUM |
| **OTLP Telemetry** | Go runtime | ZeroClaw has Observer trait | LOW |
| **OAuth Device Flow** | Cloudflare workers | Need auth system | MEDIUM |
| **Policy Engine** | Go runtime | ZeroClaw has SecurityPolicy | LOW |
| **Agent Spawning** | Go runtime | Need agent orchestration | HIGH |

### 2.3 Features from OpenClaw (Reference Implementation)

| Feature | OpenClaw Implementation | Relevance to Pryx |
|---------|------------------------|-------------------|
| **Multi-agent routing** | Gateway config routing | Adopt for agent orchestration |
| **DM pairing policy** | `dmPolicy="pairing"` | ✅ ZeroClaw has this |
| **Voice Wake + Talk Mode** | macOS/iOS/Android nodes | Future consideration |
| **Live Canvas** | A2UI system | Future consideration |
| **Browser control** | Dedicated Chrome/Chromium | ✅ ZeroClaw has browser tool |
| **Group routing** | Mention gating, reply tags | Adopt for channel management |
| **Session model** | `main` for direct, group isolation | Adopt session patterns |
| **Gmail Pub/Sub** | Webhook triggers | Adopt webhook patterns |

### 2.4 Features from OpenCode (Patterns to Adopt)

| Feature | OpenCode Implementation | Relevance to Pryx |
|---------|------------------------|-------------------|
| **SQLite sessions** | Single DB for all state | ✅ ZeroClaw uses SQLite |
| **Error handling** | Structured error types | Adopt error patterns |
| **Tool calling** | TypeBox schemas | Adopt for tool validation |
| **LSP integration** | Language Server Protocol | Adopt for code intelligence |
| **Auto compact** | Context window management | Adopt for long sessions |
| **Custom commands** | Markdown files | Adopt for skills |
| **External editor** | Editor integration | Adopt for TUI |

---

## 3. Migration Tasks

### 3.1 Phase 1: Core Migration (HIGH Priority)

```yaml
tasks:
  - id: MIG-001
    title: "Rename ZeroClaw to Pryx"
    description: "Replace all zeroclaw references with pryx in source code"
    files:
      - Cargo.toml
      - src/**/*.rs
      - README.md
      - docs/**
    status: pending
    
  - id: MIG-002
    title: "Add models.dev provider integration"
    description: "Integrate models.dev API for 84+ providers"
    files:
      - src/providers/models_dev.rs (new)
      - src/providers/mod.rs
    status: pending
    
  - id: MIG-003
    title: "Add MCP native client"
    description: "Implement MCP client as a Tool trait implementation"
    files:
      - src/tools/mcp.rs (new)
      - src/tools/mod.rs
    status: pending
    
  - id: MIG-004
    title: "Add TUI interface"
    description: "Port Solid/TUI or create new Rust TUI with ratatui"
    files:
      - src/tui/ (new)
      - src/main.rs
    status: pending
    
  - id: MIG-005
    title: "Add Local Web Admin"
    description: "Port React web admin or create new with axum + leptos/yew"
    files:
      - src/web/ (new)
      - src/gateway/mod.rs
    status: pending
```

### 3.2 Phase 2: Extended Features (MEDIUM Priority)

```yaml
tasks:
  - id: MIG-006
    title: "Add Pryx Mesh sync"
    description: "Implement multi-device sync protocol"
    files:
      - src/mesh/ (new)
    status: pending
    
  - id: MIG-007
    title: "Add session timeline"
    description: "Track all actions, tool calls, approvals"
    files:
      - src/session/ (new)
    status: pending
    
  - id: MIG-008
    title: "Add cost tracking"
    description: "Token usage and cost monitoring"
    files:
      - src/observability/cost.rs (new)
    status: pending
    
  - id: MIG-009
    title: "Add OAuth device flow"
    description: "Implement RFC 8628 OAuth Device Flow"
    files:
      - src/auth/ (new)
    status: pending
    
  - id: MIG-010
    title: "Add agent spawning"
    description: "Multi-agent orchestration system"
    files:
      - src/agent/spawn.rs (new)
    status: pending
```

### 3.3 Phase 3: Desktop Integration (LOW Priority)

```yaml
tasks:
  - id: MIG-011
    title: "Add Tauri desktop wrapper"
    description: "Create Tauri v2 wrapper for desktop app"
    files:
      - apps/desktop/ (new)
    status: pending
    
  - id: MIG-012
    title: "Add system tray"
    description: "Native system tray integration"
    files:
      - apps/desktop/src/tray.rs (new)
    status: pending
```

---

## 4. Compatibility Matrix

### 4.1 Channel Support

| Channel | Pryx (Go) | ZeroClaw (Rust) | Gap |
|---------|-----------|-----------------|-----|
| CLI | ✅ | ✅ | None |
| Telegram | ✅ | ✅ | None |
| Discord | ✅ | ✅ | None |
| Slack | ✅ | ✅ | None |
| WhatsApp | ❌ | ✅ | ZeroClaw has more |
| Matrix | ❌ | ✅ | ZeroClaw has more |
| iMessage | ❌ | ✅ | ZeroClaw has more |
| Webhook | ✅ | ✅ | None |

### 4.2 Provider Support

| Provider Type | Pryx (Go) | ZeroClaw (Rust) | Gap |
|---------------|-----------|-----------------|-----|
| OpenAI | ✅ | ✅ | None |
| Anthropic | ✅ | ✅ | None |
| Google/Gemini | ✅ | ⚠️ | Via OpenRouter |
| xAI/Grok | ✅ | ⚠️ | Via OpenRouter |
| Ollama | ✅ | ✅ | None |
| OpenRouter | ✅ | ✅ | None |
| Groq | ✅ | ⚠️ | Via OpenRouter |
| Mistral | ✅ | ⚠️ | Via OpenRouter |
| Custom endpoints | ✅ | ✅ | custom:https:// |

### 4.3 Feature Parity Checklist

- [x] Agent loop with tools
- [x] SQLite memory
- [x] Vector search
- [x] FTS5 keyword search
- [x] Shell tool
- [x] File tools
- [x] Browser tool
- [x] Gateway server
- [x] Daemon mode
- [x] Cron scheduler
- [x] Security pairing
- [x] Channel allowlists
- [x] Tunnel support
- [ ] TUI interface
- [ ] Local web admin
- [ ] 84+ providers via models.dev
- [ ] MCP native client
- [ ] Mesh sync
- [ ] Session timeline
- [ ] Cost tracking
- [ ] Agent spawning
- [ ] Desktop app

---

## 5. Recommended Migration Order

### Phase 1: Foundation (Week 1-2)
1. Rename zeroclaw → pryx (MIG-001)
2. Verify all tests pass
3. Update documentation

### Phase 2: Core Features (Week 3-4)
1. Add models.dev integration (MIG-002)
2. Add MCP client (MIG-003)
3. Add TUI (MIG-004)

### Phase 3: UI/UX (Week 5-6)
1. Add Local Web Admin (MIG-005)
2. Port desktop app if needed (MIG-011)

### Phase 4: Extended Features (Week 7-8)
1. Add Mesh sync (MIG-006)
2. Add session timeline (MIG-007)
3. Add cost tracking (MIG-008)
4. Add agent spawning (MIG-010)

---

## 6. Files to Modify/Create

### 6.1 Rename Operations

```bash
# Binary name
s/zeroclaw/pryx/g Cargo.toml

# Config directory
~/.zeroclaw → ~/.pryx

# Environment variables
ZEROCLAW_* → PRYX_*

# Source references
grep -r "zeroclaw" src/ | xargs sed -i 's/zeroclaw/pryx/g'
```

### 6.2 New Files to Create

```
src/
├── providers/
│   └── models_dev.rs      # models.dev API integration
├── tools/
│   └── mcp.rs             # MCP native client
├── tui/
│   ├── mod.rs             # TUI module
│   ├── app.rs             # App state
│   ├── ui.rs              # UI components
│   └── handlers.rs        # Event handlers
├── web/
│   ├── mod.rs             # Web admin module
│   ├── routes.rs          # HTTP routes
│   └── static/            # Static assets
├── mesh/
│   ├── mod.rs             # Mesh sync module
│   ├── protocol.rs        # Sync protocol
│   └── discovery.rs       # Device discovery
├── session/
│   ├── mod.rs             # Session management
│   └── timeline.rs        # Action timeline
└── agent/
    └── spawn.rs           # Agent spawning
```

---

## 7. Testing Strategy

### 7.1 Unit Tests
- All existing ZeroClaw tests (1,017 tests) must pass
- New modules must have >80% coverage

### 7.2 Integration Tests
- Channel connectivity tests
- Provider integration tests
- Memory system tests
- Security tests

### 7.3 E2E Tests
- Full agent loop tests
- Multi-channel tests
- Mesh sync tests

---

## 8. Documentation Updates

### Files to Update
- [ ] README.md - Update branding and features
- [ ] CHANGELOG.md - Add migration notes
- [ ] docs/prd/prd.md - Update architecture section
- [ ] docs/architecture/SYSTEM_ARCHITECTURE.md - New Rust architecture
- [ ] docs/guides/*.md - Update all guides

### New Documentation
- [ ] docs/migration/README.md - Migration guide
- [ ] docs/migration/MIGRATION_CHECKLIST.md - Step-by-step checklist
- [ ] docs/architecture/RUST_ARCHITECTURE.md - New architecture docs

---

## 9. Appendix

### A. Reference Repositories

| Repo | URL | Purpose |
|------|-----|---------|
| ZeroClaw | `.temp_refs/zeroclaw-main/` | Target codebase |
| OpenClaw | `.temp_refs/openclaw/` | Reference implementation |
| OpenCode | `.temp_refs/opencode/` | SQLite session patterns |

### B. Key Commands

```bash
# Build
cargo build --release

# Test
cargo test

# Run
cargo run --release

# Install
cargo install --path .
```

### C. Configuration Migration

```toml
# Old (Pryx Go)
# ~/.pryx/config.json

# New (Pryx Rust)
# ~/.pryx/config.toml
[pryx]
default_provider = "openrouter"
default_model = "anthropic/claude-sonnet-4"

[memory]
backend = "sqlite"
embedding_provider = "openai"

[gateway]
require_pairing = true
port = 8080

[mesh]
enabled = false
```

---

*Document generated by Pryx migration analysis on 2026-02-15*
