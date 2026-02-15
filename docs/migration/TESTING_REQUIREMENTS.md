# Testing Requirements for All Migration Tasks

> **Version**: 1.0  
> **Updated**: 2026-02-15  
> **Applies to**: All MIG-* tasks

---

## Mandatory Testing Standards

Every migration task MUST include comprehensive testing before being marked complete.

### Test Types Required

| Type | Requirement | Location |
|------|-------------|----------|
| **Unit Tests** | All public functions | `src/*/tests.rs` or inline `#[test]` |
| **Integration Tests** | Module interactions | `tests/` directory |
| **Async Tests** | Async operations | `#[tokio::test]` |
| **Fixtures** | Test data | `tests/fixtures/` |
| **E2E Tests** | Critical paths | `tests/e2e/` |

### Coverage Requirements

| Module Type | Minimum Coverage |
|-------------|------------------|
| Core logic | 90% |
| Providers | 80% |
| Channels | 80% |
| Tools | 85% |
| Security | 95% |
| UI/TUI | 70% |

---

## Test Structure

### 1. Unit Tests (Inline)

```rust
// src/providers/models_dev.rs

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_fetch_provider_catalog() {
        // Test catalog fetch
    }

    #[test]
    fn test_provider_from_catalog_entry() {
        // Test provider creation
    }

    #[tokio::test]
    async fn test_async_catalog_fetch() {
        // Test async operations
    }
}
```

### 2. Integration Tests

```rust
// tests/providers_integration_test.rs

use pryx::providers::*;

#[tokio::test]
async fn test_models_dev_provider_integration() {
    // Test with real (mocked) API
}

#[test]
fn test_provider_factory_with_models_dev() {
    // Test provider factory integration
}
```

### 3. Fixtures

```
tests/fixtures/
├── providers/
│   ├── catalog_response.json
│   └── provider_configs.toml
├── channels/
│   ├── telegram_message.json
│   └── discord_event.json
└── memory/
    └── sample_memory.db
```

### 4. E2E Tests

```rust
// tests/e2e/agent_flow_test.rs

#[tokio::test]
#[ignore] // Run with --ignored flag for E2E
async fn test_full_agent_flow_with_models_dev() {
    // Test complete agent flow
}
```

---

## Task-Specific Testing Requirements

### MIG-002: models.dev Provider Integration
- [ ] Unit: Catalog fetch parsing
- [ ] Unit: Provider creation from catalog
- [ ] Integration: Provider factory integration
- [ ] E2E: Full chat with models.dev provider

### MIG-003: MCP Native Client
- [ ] Unit: Protocol message parsing
- [ ] Unit: Transport layer (stdio, HTTP, SSE)
- [ ] Integration: Tool discovery
- [ ] Integration: Permission enforcement
- [ ] E2E: Full MCP tool execution

### MIG-004: TUI Interface
- [ ] Unit: Component rendering
- [ ] Unit: Event handling
- [ ] Integration: State management
- [ ] E2E: Full TUI session

### MIG-006: Pryx Mesh Sync
- [ ] Unit: Protocol messages
- [ ] Unit: Device pairing
- [ ] Integration: WebSocket communication
- [ ] Integration: Event sync
- [ ] E2E: Multi-device sync

### MIG-010: Agent Spawning
- [ ] Unit: Sub-agent creation
- [ ] Unit: Resource scoping
- [ ] Integration: Message bus
- [ ] E2E: Parallel agent execution

### MIG-021: OS Keychain Integration
- [ ] Unit: Keychain operations (mocked)
- [ ] Integration: macOS Keychain
- [ ] Integration: Linux Secret Service
- [ ] E2E: Secret storage and retrieval

### MIG-023: Human-in-the-Loop Approvals
- [ ] Unit: Permission catalog
- [ ] Unit: Policy resolution
- [ ] Integration: TUI prompts
- [ ] E2E: Approval flow

### MIG-005: Local Web Admin
- [ ] Unit: Route handlers
- [ ] Integration: HTTP endpoints
- [ ] E2E: Full admin session

### MIG-007: Session Timeline
- [ ] Unit: Timeline entry creation
- [ ] Unit: Query operations
- [ ] Integration: Storage persistence
- [ ] E2E: Full session tracking

### MIG-008: Cost Tracking
- [ ] Unit: Token counting
- [ ] Unit: Cost calculation
- [ ] Integration: Provider integration
- [ ] E2E: Budget enforcement

### MIG-009: OAuth Device Flow
- [ ] Unit: Code generation
- [ ] Unit: Token exchange
- [ ] Integration: Provider OAuth
- [ ] E2E: Full OAuth flow

### MIG-013: OAuth Token Refresh
- [ ] Unit: Refresh logic
- [ ] Integration: Token storage
- [ ] E2E: Token expiry and refresh

### MIG-018: Clipboard Tool
- [ ] Unit: Read/write operations
- [ ] Integration: Format support
- [ ] E2E: Clipboard roundtrip

### MIG-019: Screen/Terminal Tool
- [ ] Unit: Screen capture
- [ ] Unit: ANSI parsing
- [ ] Integration: Terminal emulation
- [ ] E2E: Screen interaction

### MIG-020: LSP Integration
- [ ] Unit: LSP message parsing
- [ ] Integration: Server connection
- [ ] Integration: Go to definition
- [ ] E2E: Full LSP session

### MIG-022: Audit Logging
- [ ] Unit: Log entry creation
- [ ] Unit: Query operations
- [ ] Integration: Storage
- [ ] E2E: Full audit trail

### MIG-026: Multi-Agent Routing
- [ ] Unit: Agent selection
- [ ] Unit: Load balancing
- [ ] Integration: Routing logic
- [ ] E2E: Multi-agent dispatch

### MIG-027: PKCE Support
- [ ] Unit: Verifier generation
- [ ] Unit: Challenge generation
- [ ] Integration: OAuth flow
- [ ] E2E: PKCE auth

### MIG-011: Tauri Desktop
- [ ] Unit: Window management
- [ ] Integration: IPC
- [ ] E2E: Desktop app launch

### MIG-012: System Tray
- [ ] Unit: Tray operations
- [ ] Integration: Menu handling
- [ ] E2E: Tray interaction

### MIG-014: Email Channel
- [ ] Unit: IMAP operations
- [ ] Unit: SMTP operations
- [ ] Integration: Email parsing
- [ ] E2E: Send/receive email

### MIG-015: Signal Channel
- [ ] Unit: Signal protocol
- [ ] Integration: Message handling
- [ ] E2E: Signal conversation

### MIG-016: Teams Channel
- [ ] Unit: Teams API
- [ ] Integration: Message handling
- [ ] E2E: Teams conversation

### MIG-017: Voice Wake + Talk
- [ ] Unit: Wake word detection
- [ ] Unit: Voice processing
- [ ] Integration: Audio pipeline
- [ ] E2E: Voice session

### MIG-024: Device Management
- [ ] Unit: Device CRUD
- [ ] Integration: Storage
- [ ] E2E: Device lifecycle

### MIG-025: WebSocket Coordination
- [ ] Unit: Message protocol
- [ ] Integration: Connection handling
- [ ] E2E: Real-time sync

### MIG-028: Auto-Update
- [ ] Unit: Version check
- [ ] Unit: Download
- [ ] Integration: Update apply
- [ ] E2E: Full update flow

---

## Running Tests

```bash
# All tests
cargo test

# Unit tests only
cargo test --lib

# Integration tests only
cargo test --test '*'

# E2E tests (ignored by default)
cargo test -- --ignored

# With coverage
cargo tarpaulin --out Html

# Specific module
cargo test --lib providers::
```

---

## Test Checklist for Task Completion

Before marking any MIG-* task as complete:

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] E2E tests pass (if applicable)
- [ ] Coverage meets minimum requirements
- [ ] Fixtures added for test data
- [ ] Test documentation updated

---

## Existing Test Coverage

Current pryx-rust has:
- **675** unit tests (`#[test]`)
- **168** async tests (`#[tokio::test]`)
- **2** integration tests (`tests/` dir)

All existing tests MUST continue to pass after any changes.

---

*Document created: 2026-02-15*
