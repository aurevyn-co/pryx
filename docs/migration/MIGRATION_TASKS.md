# Migration Tasks for Pryx → Rust

> Generated from gap analysis on 2026-02-15

## Phase 1: Foundation (HIGH Priority)

### MIG-001: Rename ZeroClaw to Pryx
- **Status**: pending
- **Priority**: HIGH
- **Description**: Replace all zeroclaw references with pryx in source code
- **Files**:
  - `newkind/zeroclaw-main/Cargo.toml`
  - `newkind/zeroclaw-main/src/**/*.rs`
  - `newkind/zeroclaw-main/README.md`
  - All doc files
- **Commands**:
  ```bash
  # Update Cargo.toml
  sed -i '' 's/zeroclaw/pryx/g' Cargo.toml
  
  # Update source files
  find src -name "*.rs" -exec sed -i '' 's/zeroclaw/pryx/g' {} \;
  
  # Update config paths
  # ~/.zeroclaw → ~/.pryx
  # ZEROCLAW_* → PRYX_*
  ```

### MIG-002: Add models.dev Provider Integration
- **Status**: pending
- **Priority**: HIGH
- **Description**: Integrate models.dev API for 84+ AI providers
- **Files**:
  - `src/providers/models_dev.rs` (new)
  - `src/providers/mod.rs` (modify)
- **Implementation**:
  - Fetch provider list from models.dev API
  - Create Provider trait implementations dynamically
  - Support OpenAI-compatible endpoints

### MIG-003: Add MCP Native Client
- **Status**: pending
- **Priority**: HIGH
- **Description**: Implement MCP client as a Tool trait implementation
- **Files**:
  - `src/tools/mcp.rs` (new)
  - `src/tools/mod.rs` (modify)
- **Implementation**:
  - Support stdio and HTTP transports
  - Tool discovery from MCP servers
  - Permission integration

### MIG-004: Add TUI Interface
- **Status**: pending
- **Priority**: HIGH
- **Description**: Port or create new TUI with ratatui
- **Files**:
  - `src/tui/mod.rs` (new)
  - `src/tui/app.rs` (new)
  - `src/tui/ui.rs` (new)
  - `src/tui/handlers.rs` (new)
- **Reference**: OpenCode's Bubble Tea TUI patterns

## Phase 2: Extended Features (MEDIUM Priority)

### MIG-005: Add Local Web Admin
- **Status**: pending
- **Priority**: MEDIUM
- **Description**: Port React web admin or create with axum + leptos/yew
- **Files**:
  - `src/web/mod.rs` (new)
  - `src/web/routes.rs` (new)
  - `src/web/static/` (new)

### MIG-006: Add Pryx Mesh Sync
- **Status**: pending
- **Priority**: MEDIUM
- **Description**: Implement multi-device sync protocol
- **Files**:
  - `src/mesh/mod.rs` (new)
  - `src/mesh/protocol.rs` (new)
  - `src/mesh/discovery.rs` (new)

### MIG-007: Add Session Timeline
- **Status**: pending
- **Priority**: MEDIUM
- **Description**: Track all actions, tool calls, approvals
- **Files**:
  - `src/session/mod.rs` (new)
  - `src/session/timeline.rs` (new)

### MIG-008: Add Cost Tracking
- **Status**: pending
- **Priority**: MEDIUM
- **Description**: Token usage and cost monitoring
- **Files**:
  - `src/observability/cost.rs` (new)

### MIG-009: Add OAuth Device Flow
- **Status**: pending
- **Priority**: MEDIUM
- **Description**: Implement RFC 8628 OAuth Device Flow
- **Files**:
  - `src/auth/mod.rs` (new)
  - `src/auth/oauth.rs` (new)

### MIG-010: Add Agent Spawning
- **Status**: pending
- **Priority**: MEDIUM
- **Description**: Multi-agent orchestration system
- **Files**:
  - `src/agent/spawn.rs` (new)

## Phase 3: Desktop Integration (LOW Priority)

### MIG-011: Add Tauri Desktop Wrapper
- **Status**: pending
- **Priority**: LOW
- **Description**: Create Tauri v2 wrapper for desktop app
- **Files**:
  - `apps/desktop/` (new)

### MIG-012: Add System Tray
- **Status**: pending
- **Priority**: LOW
- **Description**: Native system tray integration
- **Files**:
  - `apps/desktop/src/tray.rs` (new)

---

## Progress Tracking

| Task ID | Title | Status | Started | Completed |
|---------|-------|--------|---------|-----------|
| MIG-001 | Rename ZeroClaw to Pryx | pending | - | - |
| MIG-002 | models.dev integration | pending | - | - |
| MIG-003 | MCP native client | pending | - | - |
| MIG-004 | TUI interface | pending | - | - |
| MIG-005 | Local web admin | pending | - | - |
| MIG-006 | Pryx Mesh sync | pending | - | - |
| MIG-007 | Session timeline | pending | - | - |
| MIG-008 | Cost tracking | pending | - | - |
| MIG-009 | OAuth device flow | pending | - | - |
| MIG-010 | Agent spawning | pending | - | - |
| MIG-011 | Tauri desktop | pending | - | - |
| MIG-012 | System tray | pending | - | - |

---

## Notes

- All tasks should have corresponding test coverage
- Documentation must be updated for each task
- Breaking changes need migration guides
