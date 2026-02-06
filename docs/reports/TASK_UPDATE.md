# Architecture Refactoring: Rust Host Serves Local Web + API on Same Port

## Overview
Refactor Pryx architecture to follow OpenCode/OpenClaw pattern where Rust host serves everything on the same port, enabling local-first deployment with unified access.

## CONFIRMED DETAILS

### Static Port: 42424
- Memorable, unprivileged, not commonly used
- Serves HTTP, WebSocket, and static files

---

## Rust Host Responsibilities (apps/host/)

1. HTTP/WebSocket server on port 42424
2. Local web UI served from apps/local-web/
3. Admin API endpoints:
   - Health checks
   - Config read/write
   - Skills list/enable/disable
   - Providers list
   - MCP tools list
4. IPC bridge to Go runtime via JSON-RPC
5. Tauri updater for updates
6. Telemetry export to Cloudflare Workers
7. Error/logging (local + cloudflare)

---

## Go Runtime Responsibilities (apps/runtime/)

1. Agent orchestration
2. AI model providers (OpenAI, Anthropic, Google, etc. - 84+)
3. Local AI providers (Ollama)
4. Channels (Telegram, Slack, Discord, Webhooks)
5. Memory/RAG system
6. Vault/secret management (Argon2id, keychain)
7. MCP tool execution
8. Mesh/P2P networking (Tailscale + our system)

---

## Folder Structure

- `apps/local-web/` - **NEW**: Local admin web UI (moved from apps/web/)
- `apps/web/` - **STAYS**: Cloud-only deployment (pryx.dev)
- `apps/runtime/` - **STAYS**: Agent execution, channels, integrations
- `apps/host/` - **MODIFIED**: Add HTTP/WebSocket server, serve local-web

---

## Migration Plan (one by one)

### 1. Create apps/local-web/ folder
- Move local admin components from apps/web/
- Keep apps/web/ for cloud-only (pryx.dev)

### 2. Add HTTP/WebSocket server to Rust host
- Add dependencies (actix-web or axum, websocket, tower-http)
- Create server module: apps/host/src/server/
- Static port 42424
- Serve static files from local-web/

### 3. Migrate admin API to Rust
- Health endpoint
- Config read/write (IPC to Go for write)
- Skills list/enable/disable
- Providers list (IPC to Go)
- MCP tools list (IPC to Go)

### 4. Update TUI connection
- Connect to localhost:42424/ws instead of runtime port

### 5. Implement IPC layer
- Extend existing hostrpc/ module
- JSON-RPC between Rust and Go

### 6. Migrate telemetry to Rust
- Collection stays in Go
- Export to Cloudflare Workers from Rust

### 7. Update CLI commands
- Commands connect via IPC to Rust host
- Rust host forwards to Go runtime as needed

---

## Implementation Order (Subtasks)

### Phase 1: Foundation
1. [ ] Create `apps/local-web/` directory structure (Vite + React)
2. [ ] Move local components from `apps/web/` to `apps/local-web/`
3. [ ] Update `apps/web/` to be cloud-only (remove local admin)

### Phase 2: Rust Host Server
4. [ ] Add HTTP server dependencies to `apps/host/Cargo.toml`
5. [ ] Create `apps/host/src/server/` module
6. [ ] Implement static file serving for local web
7. [ ] Implement REST API endpoints in Rust
8. [ ] Implement WebSocket handler for TUI events
9. [ ] Update Tauri config to use port 42424

### Phase 3: Migration
10. [ ] Migrate admin API handlers from Go to Rust
11. [ ] Update TUI to connect to localhost:42424/ws instead of runtime
12. [ ] Enhance Go runtime IPC (hostrpc) for sidecar mode

### Phase 4: Integration
13. [ ] Update `pryx` CLI to start host (desktop mode)
14. [ ] Update `pryx-core` for headless/server mode
15. [ ] Update documentation (README, architecture docs)
16. [ ] Test full workflow: host → local web → API → runtime → agent

---

## File Changes Summary

| Action | Source → Destination |
|--------|---------------------|
| MOVE | `apps/web/src/components/Dashboard.tsx` → `apps/local-web/src/Dashboard.tsx` |
| MOVE | `apps/web/src/components/skills/` → `apps/local-web/src/components/skills/` |
| MOVE | `apps/web/src/pages/dashboard.astro` → `apps/local-web/src/pages/dashboard.astro` |
| MOVE | `apps/web/src/pages/skills.astro` → `apps/local-web/src/pages/skills.astro` |
| MOVE | `apps/web/src/middleware.ts` → `apps/local-web/src/middleware.ts` |
| MODIFY | `apps/host/Cargo.toml` (add HTTP/WebSocket deps) |
| CREATE | `apps/host/src/server/` (new server module) |
| MODIFY | `apps/host/src-tauri/src/main.rs` (add server startup) |
| MODIFY | `apps/tui/src/services/skills-api.ts` (change port to 42424) |
| MODIFY | `apps/tui/src/services/ws.ts` (connect to localhost:42424/ws) |
| CREATE | `docs/architecture/HOST_SERVER.md` (new architecture doc) |

---

## Success Criteria

- [ ] Port 42424 serves local web UI, HTTP API, and WebSocket
- [ ] apps/local-web/ exists and serves admin UI
- [ ] apps/web/ serves only cloud (pryx.dev)
- [ ] Go runtime runs as sidecar, handles agent work
- [ ] TUI connects to localhost:42424/ws
- [ ] All admin operations work via local web UI
- [ ] Telemetry exports to Cloudflare Workers
