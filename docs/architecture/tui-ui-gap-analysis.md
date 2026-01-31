# TUI/CLI UI Gap Analysis

## Overview

This document analyzes which runtime features have TUI/CLI UI and which are missing, prioritized by user impact.

---

## Current TUI Views (Implemented)

### 1. Chat View ✅
**File**: `apps/tui/src/components/Chat.tsx`

**Features**:
- Message input/output
- Streaming responses
- Message history display
- Keyboard shortcuts (1)

**Runtime APIs Used**:
- WebSocket (`/ws`)
- Session management (via store)

---

### 2. Sessions View ✅
**File**: `apps/tui/src/components/SessionExplorer.tsx`

**Features**:
- List all sessions
- Switch between sessions
- Create new sessions
- Delete sessions
- Keyboard shortcuts (2)

**Runtime APIs Used**:
- Session store (SQLite)
- Message store

---

### 3. Channels View ✅
**File**: `apps/tui/src/components/Channels.tsx`

**Features**:
- View channel integrations
- Configure webhooks
- Telegram/Discord/Slack settings
- Keyboard shortcuts (3)

**Runtime APIs Used**:
- `internal/channels`

---

### 4. Skills View ✅
**File**: `apps/tui/src/components/Skills.tsx`

**Features**:
- Browse available skills
- Enable/disable skills
- View skill details
- Keyboard shortcuts (4)

**Runtime APIs Used**:
- `GET /skills`
- `GET /skills/{id}`
- `GET /skills/{id}/body`

---

### 5. Settings View ✅
**File**: `apps/tui/src/components/Settings.tsx`

**Features**:
- Provider configuration
- Model selection
- API key management
- General settings
- Keyboard shortcuts (5)

**Runtime APIs Used**:
- Config service
- Provider APIs

---

### 6. Command Palette ✅
**File**: `apps/tui/src/components/SearchableCommandPalette.tsx`

**Features**:
- Quick navigation (1-5)
- Command search
- Keyboard shortcuts (?)
- Categories: Navigation, Chat, System

---

### 7. Setup/Onboarding ✅
**File**: `apps/tui/src/components/SetupRequired.tsx`, `OnboardingWizard.tsx`

**Features**:
- Provider setup
- API key input
- Model selection
- First-time configuration

---

## Missing TUI/CLI UI (Runtime Features Without UI)

### 🔴 P0 - Critical Missing UI

#### 1. Provider Management /connect Command
**Runtime Feature**: Provider configuration, connection testing

**Current State**: 
- Settings view has basic provider config
- No `/connect` command in TUI
- Users can't see provider status or test connections

**What's Missing**:
- Provider connection status indicator
- `/connect` command to add/configure providers
- Provider health check UI
- Model selection per provider
- Connection testing with error feedback

**Runtime APIs Available**:
- `GET /api/v1/providers`
- `GET /api/v1/providers/{id}/models`
- `GET /api/v1/models`
- Catalog service (84 providers, 1417 models)

**Priority**: **P0** - Blocking basic functionality

---

#### 2. MCP Tools Management
**Runtime Feature**: MCP server management, tool execution

**Current State**:
- Runtime has full MCP support (`internal/mcp`)
- TUI has no MCP management UI
- Users can't see available tools
- Can't enable/disable MCP servers

**What's Missing**:
- List MCP servers
- View available tools per server
- Enable/disable MCP servers
- Tool execution history
- Tool approval UI (for dangerous operations)

**Runtime APIs Available**:
- `GET /mcp/tools`
- `POST /mcp/tools/call`
- MCP Manager with caching

**Priority**: **P0** - Core feature not accessible

---

#### 3. Audit Log Viewer
**Runtime Feature**: Complete audit trail of all operations

**Current State**:
- Runtime logs everything to `audit_log` table
- TUI has no way to view audit history
- Users can't see what tools were called

**What's Missing**:
- Audit log browser
- Filter by date/action/tool
- Export audit log
- Cost breakdown per operation
- Tool execution details

**Runtime APIs Available**:
- `internal/audit` repository
- SQLite audit_log table

**Priority**: **P0** - Security/compliance feature

---

### 🟡 P1 - High Priority Missing UI

#### 4. Cost Tracking Dashboard
**Runtime Feature**: Cost tracking and optimization

**Current State**:
- Runtime tracks costs (`internal/cost`)
- TUI has no cost visibility
- Users can't see spending

**What's Missing**:
- Daily/weekly/monthly cost breakdown
- Cost per provider
- Cost per session
- Budget alerts
- Cost optimization suggestions

**Runtime APIs Available**:
- `internal/cost/service.go`
- Audit log with cost data

**Priority**: **P1** - Important for user awareness

---

#### 5. Agent Spawning / Sub-agents
**Runtime Feature**: Spawn sub-agents for concurrent tasks

**Current State**:
- Runtime supports agent spawning (`internal/agent/spawn`)
- TUI has no sub-agent UI
- Users can't spawn or monitor sub-agents

**What's Missing**:
- Sub-agent list
- Spawn new agent UI
- Monitor agent progress
- View agent results
- Agent tree visualization

**Runtime APIs Available**:
- `internal/agent/spawn`
- Session forking

**Priority**: **P1** - Advanced feature needs UI

---

#### 6. Policy Engine / Approvals
**Runtime Feature**: Policy-based operation gating

**Current State**:
- Runtime has policy engine (`internal/policy`)
- TUI has no approval UI
- Dangerous operations can't be approved

**What's Missing**:
- Approval request UI
- Policy configuration
- Approval history
- Scope management (workspace, host)
- Dangerous operation warnings

**Runtime APIs Available**:
- `internal/policy`

**Priority**: **P1** - Security feature

---

#### 7. Mesh / Multi-Device Status
**Runtime Feature**: Multi-device coordination

**Current State**:
- Runtime has mesh manager (`internal/mesh`)
- TUI has no mesh UI
- Users can't see connected devices

**What's Missing**:
- Device list
- Connection status
- Session sync status
- Device management (rename, remove)
- Mesh activity log

**Runtime APIs Available**:
- `internal/mesh`
- WebSocket mesh protocol

**Priority**: **P1** - Multi-device users need this

---

### 🟢 P2 - Medium Priority Missing UI

#### 8. Memory Management
**Runtime Feature**: Conversation memory management

**Current State**:
- Runtime has memory manager (`internal/memory`)
- TUI has no memory UI
- Users can't manage context window

**What's Missing**:
- Memory usage display
- Context window visualization
- Manual memory cleanup
- Memory optimization suggestions
- Token count display

**Runtime APIs Available**:
- `internal/memory`

**Priority**: **P2** - Nice to have for power users

---

#### 9. Telemetry / Health Dashboard
**Runtime Feature**: System health and telemetry

**Current State**:
- Runtime exports telemetry (`internal/telemetry`)
- TUI has basic connection status only
- No detailed health metrics

**What's Missing**:
- System health dashboard
- Telemetry export status
- Performance metrics
- Error logs viewer
- Diagnostics (like `pryx doctor`)

**Runtime APIs Available**:
- `GET /health`
- `internal/telemetry`
- `internal/doctor`

**Priority**: **P2** - Debugging/support feature

---

#### 10. Constraints / Model Catalog Browser
**Runtime Feature**: Model constraints and capabilities

**Current State**:
- Runtime has constraints engine (`internal/constraints`)
- Catalog loads 1417 models from models.dev
- TUI has no model browser

**What's Missing**:
- Model catalog browser
- Filter by capabilities (vision, tools, etc.)
- Pricing comparison
- Model recommendations
- Constraint visualization

**Runtime APIs Available**:
- `internal/constraints`
- `internal/models` (catalog)

**Priority**: **P2** - Power user feature

---

## TUI↔Runtime Communication Architecture

### Current Architecture

```
┌─────────────────────────────────────────────────────────┐
│                        TUI                              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │  Chat   │ │Sessions │ │Channels │ │Settings │       │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │
│       │           │           │           │             │
│       └───────────┴───────────┴───────────┘             │
│                   │                                     │
│              ┌────┴────┐                                │
│              │ WebSocket│ ←── Real-time events          │
│              └────┬────┘                                │
└───────────────────┼─────────────────────────────────────┘
                    │ HTTP/WS
┌───────────────────┼─────────────────────────────────────┐
│              ┌────┴────┐                                │
│              │  Server  │                                │
│              └────┬────┘                                │
│       ┌───────────┼───────────┐                         │
│       ▼           ▼           ▼                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                   │
│  │   Bus   │ │  Store  │ │  Skills │                   │
│  └─────────┘ └─────────┘ └─────────┘                   │
│                                                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │   MCP   │ │  Agent  │ │  Mesh   │ │  Policy │       │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │
│                                                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │  Audit  │ │  Cost   │ │Telemetry│ │ Memory  │       │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │
└─────────────────────────────────────────────────────────┘
```

### Communication Patterns

1. **HTTP REST API**: Configuration, skills, providers
2. **WebSocket**: Real-time chat, events, streaming
3. **File System**: Config, cache, logs

---

## Implementation Priority Matrix

| Feature | User Impact | Implementation Complexity | Priority |
|---------|-------------|---------------------------|----------|
| Provider Management | High | Medium | **P0** |
| MCP Tools UI | High | Medium | **P0** |
| Audit Log Viewer | High | Low | **P0** |
| Cost Dashboard | Medium | Low | **P1** |
| Agent Spawning | Medium | High | **P1** |
| Policy/Approvals | High | Medium | **P1** |
| Mesh Status | Low | Medium | **P1** |
| Memory Management | Low | Medium | **P2** |
| Telemetry Dashboard | Low | Low | **P2** |
| Model Catalog Browser | Low | Medium | **P2** |

---

## Recommended Implementation Order

### Phase 1: Critical (P0)
1. **Provider Management** - Fix `/connect` command, add provider status
2. **MCP Tools UI** - List servers, view tools, enable/disable
3. **Audit Log Viewer** - Basic log browser with filters

### Phase 2: Important (P1)
4. **Cost Dashboard** - Daily breakdown, per-provider costs
5. **Policy/Approvals** - Approval UI for dangerous operations
6. **Agent Spawning** - Sub-agent list and spawn UI

### Phase 3: Enhancement (P2)
7. **Mesh Status** - Device list and sync status
8. **Memory Management** - Context window visualization
9. **Telemetry Dashboard** - Health metrics and logs
10. **Model Catalog Browser** - Browse 1417 models

---

## Technical Notes

### Adding New Views

To add a new view to the TUI:

1. **Create component** in `apps/tui/src/components/`
2. **Add to View type** in `App.tsx`:
   ```typescript
   type View = "chat" | "sessions" | ... | "newview";
   ```
3. **Add route** in `App.tsx`:
   ```typescript
   <Match when={view() === "newview"}>
     <NewView />
   </Match>
   ```
4. **Add command** in command palette:
   ```typescript
   {
     id: "newview",
     name: "New View",
     category: "Navigation",
     shortcut: "6",
     action: () => setView("newview"),
   }
   ```
5. **Add keyboard shortcut** (1-9)

### Adding CLI Commands

Runtime CLI commands are in `apps/runtime/cmd/pryx-core/`:

```go
// Example: pryx-core provider list
var providerListCmd = &cobra.Command{
    Use:   "list",
    Short: "List configured providers",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

---

## Summary

**Current State**: 5 main views implemented (Chat, Sessions, Channels, Skills, Settings)

**Critical Gaps (P0)**: 
- Provider management (`/connect` command)
- MCP tools management
- Audit log viewer

**Total Missing**: 10 major features need UI

**Estimated Effort**: 
- P0: 2-3 weeks
- P1: 3-4 weeks  
- P2: 2-3 weeks

**Recommendation**: Focus on P0 items first to unblock core functionality.
