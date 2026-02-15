# Migration Tasks for Pryx → Rust

> **Version**: 5.0 (Final - Triple-Checked)  
> **Updated**: 2026-02-15  
> **Total Tasks**: 34 (33 open + 1 closed duplicate)  
> **Status**: Ready for implementation  
> **Testing**: See [TESTING_REQUIREMENTS.md](./TESTING_REQUIREMENTS.md)

---

## Quick Stats

| Priority | Count | Tasks |
|----------|-------|-------|
| **P0 (HIGH)** | 10 | Foundation + Core gaps |
| **P1 (MEDIUM)** | 13 | Extended features |
| **P2 (LOW)** | 11 | Nice-to-have |
| **TOTAL OPEN** | **34** | |

---

## Pre-Requisites (Do First) - P0

### MIG-000: Directory Restructure ✅ DONE
- **ID**: `pryx-1bhd`
- **Status**: COMPLETED
- **Description**: Move pryx-rust to root, old codebases to _legacy/
- **Completed**: 2026-02-15
- **Result**:
  - `apps/` → `_legacy/apps/`
  - `packages/` → `_legacy/packages/`
  - `newkind/pryx-rust/src/` → `src/`
  - `newkind/pryx-rust/tests/` → `tests/`
  - `newkind/pryx-rust/examples/` → `examples/`
  - `newkind/pryx-rust/Cargo.toml` → `Cargo.toml`
  - `cargo build --release` works ✅
  - `cargo test` passes (842/843 - 1 pre-existing env-dependent failure)

### MIG-001T: Test Coverage for Existing Code
- **ID**: `pryx-4dy6`
- **Description**: Verify 80%+ coverage on existing 843 tests
- **Modules**: providers, channels, tools, memory, security, agent, gateway

---

## Phase 1: HIGH Priority (P0) - 8 remaining

| ID | Task | Foundation | Gap |
|----|------|------------|-----|
| `pryx-pw1k.1` | MIG-002: models.dev Catalog | providers/compatible.rs | Dynamic catalog |
| `pryx-9vpv` | MIG-003: MCP Native Client | None | Protocol client |
| `pryx-pw1k.3` | MIG-004: TUI Interface | None | 24 components to port |
| `pryx-nht0` | MIG-006: Mesh Sync | security/pairing.rs | Multi-device |
| `pryx-ak6j` | MIG-010: Agent Spawning | agent/loop_.rs | Sub-agents + handoff + capability |
| `pryx-h7wm` | MIG-021: OS Keychain | security/secrets.rs | Keychain storage |
| `pryx-yktr` | MIG-023: Human-in-Loop | security/policy.rs | Approval prompts |
| `pryx-q4kc` | MIG-029: Event Bus | observability (one-way) | Full pub/sub |

---

## Phase 2: MEDIUM Priority (P1) - 13 tasks

| ID | Task | Source |
|----|------|--------|
| `pryx-ogcq` | MIG-005: Local Web Admin | NEW |
| `pryx-y623` | MIG-007: Session Timeline | NEW |
| `pryx-ckp6` | MIG-008: Cost Tracking | NEW |
| `pryx-tcvm` | MIG-009: OAuth Device Flow | NEW |
| `pryx-gnft` | MIG-013: OAuth Token Refresh | NEW |
| `pryx-60pq` | MIG-018: Clipboard Tool | NEW |
| `pryx-324p` | MIG-019: Screen/Terminal Tool | NEW |
| `pryx-rxi1` | MIG-020: LSP Integration | OpenCode |
| `pryx-vqf6` | MIG-022: Audit Logging | OLD Pryx |
| `pryx-6bk1` | MIG-026: Multi-Agent Routing | OpenClaw |
| `pryx-a1hh` | MIG-027: PKCE Support | NEW |
| `pryx-i8ez` | MIG-030: Input Validation | OLD Pryx |
| `pryx-xric` | MIG-033: Universal Agent Client | OLD Pryx (universal/) |

---

## Phase 3: LOW Priority (P2) - 11 tasks

| ID | Task | Source |
|----|------|--------|
| `pryx-7rgp` | MIG-011: Tauri Desktop + hostrpc | OLD Pryx |
| `pryx-8j4v` | MIG-012: System Tray | NEW |
| `pryx-d1mn` | MIG-014: Email Channel | OpenClaw |
| `pryx-8i2w` | MIG-015: Signal Channel | OpenClaw |
| `pryx-hewu` | MIG-016: Teams Channel | OpenClaw |
| `pryx-mdpt` | MIG-017: Voice Wake + Talk | OpenClaw |
| `pryx-26ih` | MIG-024: Device Management | NEW |
| `pryx-j8ln` | MIG-025: WebSocket Coordination | NEW |
| `pryx-mwot` | MIG-028: Auto-Update | NEW |
| `pryx-twvw` | MIG-031: Patch/Diff Tool | OpenCode |
| `pryx-0tg1` | MIG-032: Web Chat Channel | OpenClaw |

---

## Existing Foundation (Don't Rebuild!)

| Module | Status | Files |
|--------|--------|-------|
| **Providers** | ✅ | OpenAI, Anthropic, OpenRouter, Ollama, Compatible (22+) |
| **Channels** | ✅ | CLI, Telegram, Discord, Slack, Matrix, WhatsApp, iMessage, Webhook |
| **Tools** | ✅ | Shell, File R/W, Memory (3), Browser, Composio |
| **Memory** | ✅ | SQLite + FTS5 + Vector + Embeddings |
| **Security** | ✅ | Pairing (6-digit), Policy (sandbox), Secrets (file) |
| **Agent** | ✅ | Single agent loop |
| **Gateway** | ✅ | Axum HTTP server |
| **Daemon** | ✅ | Long-running supervisor |
| **Cron** | ✅ | Task scheduler |
| **Tunnel** | ✅ | Cloudflare, Tailscale, ngrok, custom |
| **Skills** | ✅ | TOML manifests |
| **Observability** | ✅ | Observer trait (one-way) |

---

## Dependencies

```
MIG-000 (Restructure) ─┬─> All other tasks
MIG-001T (Tests) ──────┘

MIG-029 (Event Bus) ───┬─> MIG-010 (Agent Spawning)
                       ├─> MIG-006 (Mesh Sync)
                       ├─> MIG-026 (Multi-Agent Routing)
                       └─> MIG-033 (Universal Agent Client)

MIG-009 (OAuth) ───────┬─> MIG-013 (Token Refresh)
                       └─> MIG-027 (PKCE)

MIG-010 (Agent Spawn) ─> MIG-026 (Routing)
                       └> MIG-033 (Universal Client)

MIG-011 (Desktop) ─────> MIG-012 (Tray)
```

---

## Triple-Check Verification

### OLD Pryx Modules (41 total)
- ✅ 20 already exist in Rust
- ✅ 20 covered by tasks
- ✅ 4 not needed (federation, marketplace, nlp, trust)
- ✅ 1 added as MIG-033 (universal/)

### OpenClaw Features
- ✅ All reviewed and covered

### OpenCode Features
- ✅ All reviewed and covered

### TUI Components (24 total)
- ✅ All listed in MIG-004 description

---

## Related Documents

- [TESTING_REQUIREMENTS.md](./TESTING_REQUIREMENTS.md) - Mandatory testing standards
- [MIGRATION_GAP_ANALYSIS.md](./MIGRATION_GAP_ANALYSIS.md) - Detailed gap analysis

---

*Document updated: 2026-02-15 - Triple-checked final version with 34 tasks*
