# Pryx Implementation: 3-Phase Roadmap

> **Version**: 1.0  
> **Date**: 2026-01-29  
> **Status**: Planning Complete, Ready for Implementation  
> **Related**: `docs/prd/prd.md`, `docs/prd/prd-v2.md`, `docs/prd/ai-assisted-setup.md`

---

## Executive Summary

This document consolidates the **implementation roadmap** for Pryx's core features across three phases, ensuring **AI-assisted flows coexist with manual setup** without conflicts.

**Key Achievement**: AI-assisted setup is designed as an **optional enhancement**, not a replacement. All three configuration methods (Manual, Visual, AI-Assisted) are first-class citizens.

---

## The Three Configuration Methods

```
┌─────────────────────────────────────────────────────────────────────┐
│              PRYX CONFIGURATION ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   All paths lead to the same configuration store                    │
│                                                                     │
│    ┌─────────────┐      ┌─────────────┐      ┌─────────────┐       │
│    │   MANUAL    │      │   VISUAL    │      │  AI-ASSISTED│       │
│    │             │      │             │      │             │       │
│    │ Config files│◄────►│  TUI Forms  │◄────►│  Natural    │       │
│    │ CLI commands│      │  Web UI     │      │  Language   │       │
│    │ Env vars    │      │  Wizards    │      │  Dialogues  │       │
│    └──────┬──────┘      └──────┬──────┘      └──────┬──────┘       │
│           │                    │                    │              │
│           └────────────────────┴────────────────────┘              │
│                              │                                      │
│                              ▼                                      │
│   ┌──────────────────────────────────────────────────────────┐     │
│   │         Unified Config Store (Single Source of Truth)     │     │
│   │                                                           │     │
│   │  • ~/.pryx/config.json (JSON5)                           │     │
│   │  • ~/.pryx/vault/credentials.json (encrypted)            │     │
│   │  • ~/.pryx/mcp/servers.json                              │     │
│   │  • ~/.pryx/channels/*.json                               │     │
│   │                                                           │     │
│   └──────────────────────────────────────────────────────────┘     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Method Comparison

| Aspect | Manual | Visual | AI-Assisted |
|--------|--------|--------|-------------|
| **Best For** | Power users, CI/CD, automation | Discoverability, validation, first-time | Natural language preference, quick setup |
| **Speed** | Fastest for experts | Medium | Fastest for beginners |
| **Precision** | Exact control | Guided with validation | Conversational |
| **Examples** | Edit config.json, CLI flags | TUI forms, wizards | "Connect my Telegram" |

**Conflict Resolution**: Last-write-wins with backup + audit logging

---

## Phase 1: Foundation (Weeks 1-3)

**Theme**: Build the infrastructure that enables all three configuration methods

### 1.1 Vault & Security (Critical)

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-npno | Host | Vault encryption core (Argon2 + AES-256-GCM) | P0 |
| pryx-w8k2 | Host | Master password derivation with key stretching | P0 |
| pryx-x9m3 | Runtime | Credential storage schema (encrypted JSON) | P0 |
| pryx-y4p5 | Runtime | Memory-only decryption (clear after use) | P0 |
| pryx-z1q7 | Runtime | Access control (read-only, write-own, full) | P1 |
| pryx-a2r8 | Runtime | Audit logging (who accessed what, when) | P1 |

**Deliverable**: Secure credential vault that works with all three methods

### 1.2 Configuration Infrastructure

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-b3s9 | Runtime | Config schemas (Zod validation) | P0 |
| pryx-c4t0 | Runtime | Unified config store (JSON5 support) | P0 |
| pryx-d5u1 | Runtime | Config file watcher (hot reload) | P0 |
| pryx-e6v2 | Runtime | Backup/rollback (keep last 5 versions) | P0 |
| pryx-f7w3 | Runtime | Provider config schema (OpenAI, Anthropic, etc.) | P1 |
| pryx-g8x4 | Runtime | Channel config schema (Telegram, Discord, etc.) | P1 |
| pryx-h9y5 | Runtime | MCP config schema (servers, tools) | P1 |

**Deliverable**: Config system that supports manual editing, TUI forms, and AI parsing

### Phase 1 Success Criteria
- [ ] Can edit config files manually with validation
- [ ] Changes are automatically detected and applied
- [ ] All changes create backups
- [ ] Audit log tracks every modification
- [ ] Schema validation prevents invalid configs

---

## Phase 2: Feature Implementation (Weeks 4-8)

**Theme**: Build the three configuration paths in parallel

### Track A: Manual/CLI (Power Users)

**Goal**: Complete CLI coverage for all configuration operations

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-i0z6 | CLI | `pryx config get/set` commands | P1 |
| pryx-j1a7 | CLI | `pryx provider add/remove/list` commands | P1 |
| pryx-k2b8 | CLI | `pryx channel add/remove/list` commands | P1 |
| pryx-l3c9 | CLI | `pryx mcp add/remove/test` commands | P1 |
| pryx-m4d0 | CLI | `pryx vault add/remove/list` commands | P1 |
| pryx-n5e1 | CLI | Environment variable support ($PROVIDER_API_KEY) | P1 |
| pryx-o6f2 | CLI | Config validation command | P1 |
| pryx-p7g3 | CLI | Config export/import | P2 |

**Example Usage**:
```bash
# Add provider
pryx provider add anthropic --api-key $ANTHROPIC_API_KEY

# Add channel
pryx channel add telegram --token $BOT_TOKEN

# Add MCP server
pryx mcp add filesystem --transport stdio --command "npx -y @modelcontextprotocol/server-filesystem"

# Verify config
pryx config validate
```

### Track B: Visual/TUI (Visual Users)

**Goal**: Rich TUI forms with validation and testing

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-q8h4 | TUI | Provider management screen (add, edit, test) | P1 |
| pryx-r9i5 | TUI | Channel configuration forms (Telegram, Discord, etc.) | P1 |
| pryx-s0j6 | TUI | MCP server browser (curated list + custom URL) | P1 |
| pryx-t1k7 | TUI | Vault credential manager (secure input) | P1 |
| pryx-u2l8 | TUI | Form validation with helpful error messages | P1 |
| pryx-v3m9 | TUI | Connection testing UI (test before saving) | P1 |
| pryx-w4n0 | TUI | Configuration diff viewer (see changes before applying) | P2 |
| pryx-x5o1 | Web | Web UI equivalent screens (for headless servers) | P2 |

**Example Flow**:
```
$ pryx
[Opens TUI]

┌─────────────────────────────────────────┐
│  🔧 Pryx Control Center                │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  AI Providers                   │   │
│  │                                 │   │
│  │  ○ Anthropic [Connect]         │   │
│  │  ○ OpenAI [Connect]            │   │
│  │  ● Local (Ollama) ✓            │   │
│  │                                 │   │
│  │  [Connect New Provider]         │   │
│  └─────────────────────────────────┘   │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  Channels                       │   │
│  │                                 │   │
│  │  ● Telegram ✓                   │   │
│  │    Bot: @mybot                  │   │
│  │    Status: Connected            │   │
│  │                                 │   │
│  │  ○ Discord [Connect]           │   │
│  │  ○ Slack [Connect]             │   │
│  │                                 │   │
│  │  [Add Channel]                  │   │
│  └─────────────────────────────────┘   │
│                                         │
└─────────────────────────────────────────┘
```

### Track C: AI-Assisted (Natural Language Users)

**Goal**: Conversational setup that never blocks manual access

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-y6p2 | Runtime | Natural language intent parser | P1 |
| pryx-z7q3 | Runtime | Provider setup dialogue flows | P1 |
| pryx-a8r4 | Runtime | Channel setup dialogue flows | P1 |
| pryx-b9s5 | Runtime | MCP setup dialogue flows | P1 |
| pryx-c0t6 | Runtime | Contextual help system | P1 |
| pryx-d1u7 | Runtime | **Manual escape hatch** ("show me the config") | P1 |
| pryx-e2v8 | Runtime | Progress persistence (save/resume setup) | P2 |
| pryx-f3w9 | Runtime | Setup verification & testing | P2 |

**Critical Design Principle**: AI-assisted flows must always offer manual alternatives.

**Example Dialogue**:
```
User: "Connect my Telegram bot"

AI: I'll help you connect Telegram! You have 3 options:

     1. 🚀 Quick Setup (Recommended)
        I'll guide you through getting a bot token from @BotFather.
     
     2. 📋 Step-by-Step
        Show me instructions and I'll help if you get stuck.
     
     3. ⚙️ Manual Configuration
        Open the configuration form to enter details yourself.
     
     Which would you prefer? (1, 2, or 3)

User: 1

AI: Great! First, do you have a bot token from @BotFather?
     
     (If not, I can walk you through creating one)

User: yes, it's 123456:ABC-DEF...

AI: Perfect! I'll verify this token and set up the connection.
     
     [Testing connection...]
     ✓ Bot token is valid
     ✓ Bot name: @myawesomebot
     
     Would you like me to:
     • Set up webhook delivery (recommended)
     • Use polling mode (for local-only setups)
     
     Or type "show config" to see exactly what I'll save.
```

### Phase 2 Success Criteria
- [ ] All three methods work independently
- [ ] Changes from any method are visible to all others
- [ ] No data loss or conflicts between methods
- [ ] Users can switch methods mid-setup
- [ ] AI always offers manual escape hatch

---

## Phase 3: Advanced Features (Weeks 9-12)

**Theme**: Seamless integration between methods, advanced automation

### 3.1 Cron Jobs & Scheduled Tasks

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-g4x0 | Runtime | Natural language cron parser ("every weekday at 8pm") | P2 |
| pryx-h5y1 | Runtime | Cron scheduler service with persistence | P2 |
| pryx-i6z2 | Runtime | Task isolation (isolated agent sessions) | P2 |
| pryx-j7a3 | Runtime | Delivery to channels (Telegram, etc.) | P2 |
| pryx-k8b4 | TUI | Cron job management dashboard | P2 |
| pryx-l9c5 | TUI | Task history and logs viewer | P2 |
| pryx-m0d6 | CLI | `pryx cron add/remove/list` commands | P2 |
| pryx-n1e7 | Runtime | Retry policies and failure handling | P2 |

**Example**:
```
User: "Every weekday at 8 PM, summarize my finance data"

AI: I'll create a scheduled task for you:
     
     Task: Daily Finance Summary
     Schedule: 0 20 * * 1-5 (Mon-Fri at 8pm)
     Action: Summarize data from vault:banking
     Deliver to: Telegram
     
     Estimated cost: ~$0.05 per run
     
     Create this task? (yes/no/show config)
```

### 3.2 Security & Audit

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-o2f8 | CLI | `pryx doctor` security audit command | P1 |
| pryx-p3g9 | Runtime | Filesystem permission checks | P1 |
| pryx-q4h0 | Runtime | Secrets detection in configs | P1 |
| pryx-r5i1 | Runtime | Vault security validation | P1 |
| pryx-s6j2 | TUI | Security dashboard with recommendations | P2 |
| pryx-t7k3 | CLI | Environment variable import (`pryx vault import-env`) | P2 |

### 3.3 Integration & Polish

| Task ID | Component | Description | Priority |
|---------|-----------|-------------|----------|
| pryx-u8l4 | Runtime | Conflict resolution UI (when AI wants to override) | P2 |
| pryx-v9m5 | Runtime | Configuration sync validation | P2 |
| pryx-w0n6 | Runtime | Setup templates & recommendations | P2 |
| pryx-x1o7 | TUI | Main dashboard with unified view | P2 |

### Phase 3 Success Criteria
- [ ] Can create cron jobs via all three methods
- [ ] Security audit runs without errors
- [ ] AI suggestions don't override manual config without consent
- [ ] Seamless switching between methods works in all edge cases
- [ ] Performance: All operations complete in <3 seconds

---

## Implementation Timeline Summary

```
Month 1 (Weeks 1-4): Foundation
├── Week 1-2: Vault encryption, config schemas
├── Week 3: Unified config store, validation
└── Week 4: Provider/Channel/MCP schemas

Month 2-3 (Weeks 5-8): Feature Implementation  
├── Week 5-6: Parallel development
│   ├── Track A: CLI commands
│   ├── Track B: TUI screens
│   └── Track C: AI dialogue flows
├── Week 7: Testing, integration
└── Week 8: Polish, documentation

Month 3 (Weeks 9-12): Advanced Features
├── Week 9-10: Cron jobs, scheduler
├── Week 11: Security audit, vault TUI
└── Week 12: Integration, final testing
```

---

## Key Anti-Patterns Avoided

### ❌ AI as Gatekeeper
```
# BAD: AI blocks manual access
User: "Show me the Telegram config"
AI: "Let me walk you through it..." [forces AI-guided flow]
```

✅ **Fix**: Always offer manual option immediately.

### ❌ AI Overrides Without Consent
```
# BAD: AI changes config automatically
User: "My Telegram isn't working"
AI: [silently resets webhook] "Fixed it!"
```

✅ **Fix**: Always ask before changing existing configuration.

### ❌ Inconsistent State Between Methods
```
# BAD: Different methods show different states
TUI shows: Telegram disconnected
Config file: Telegram enabled
CLI shows: Telegram connected
```

✅ **Fix**: Single source of truth (config files), all methods read from same source.

### ❌ AI-Only Features
```
# BAD: Feature only accessible via AI
User: "How do I set up cron jobs?"
AI: "Just tell me what to schedule!" [no manual docs]
```

✅ **Fix**: All features must be accessible via all three methods.

---

## Success Metrics

| Metric | Target | Phase |
|--------|--------|-------|
| Config change success rate | >99% | Phase 1 |
| Method usage distribution | Roughly equal 33/33/33% | Phase 2 |
| Setup completion rate | >90% for all methods | Phase 2 |
| AI-to-manual switch rate | <10% | Phase 2 |
| Conflict incidents | <1% of changes | Phase 3 |
| Cron job success rate | >95% | Phase 3 |
| Security audit pass rate | 100% | Phase 3 |

---

## Integration with Existing Documents

This roadmap integrates with:

1. **`docs/prd/prd.md`** (v1 PRD)
   - Section 7: Product Surfaces (all three methods)
   - Section 8: Architecture (Session Bus, Gateway)
   - Section 9: Functional Requirements

2. **`docs/prd/prd-v2.md`** (v2 Roadmap)
   - Section 4: Local AI integration
   - Section 5: Ecosystem & Channels
   - Section 7: Autonomous workflows (cron jobs)

3. **`docs/prd/ai-assisted-setup.md`** (This feature)
   - Detailed specification of 3-method architecture
   - UX patterns for AI-assisted flows
   - Anti-patterns and conflict resolution

4. **`docs/prd/scheduled-tasks.md`**
   - Cron job implementation details
   - Natural language parsing

---

## Next Steps

1. **Review & Approve**: Stakeholder review of 3-phase plan
2. **Start Phase 1**: Begin vault encryption implementation
3. **Parallel Development**: Set up teams for Tracks A, B, C in Phase 2
4. **Weekly Sync**: Cross-track coordination meetings
5. **Milestone Reviews**: End-of-phase demos and validation

---

## Summary

**What We Built**:
- ✅ **3-Phase Implementation Plan** with clear deliverables
- ✅ **Three Configuration Methods** designed to coexist
- ✅ **Conflict Resolution** ensuring no method dominates
- ✅ **AI-Assisted as Optional** enhancement, not replacement
- ✅ **Complete Feature Coverage** across all phases

**Key Achievement**: AI-assisted setup enhances the experience **without**:
- Blocking manual configuration
- Overriding settings without consent
- Creating inconsistent state
- Removing power user capabilities

**Ready to Start**: Phase 1 can begin immediately with vault encryption implementation.

---

*Document Status*: ✅ Complete and ready for implementation
