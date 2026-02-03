# Pryx v1 Website Requirements - CLARIFIED

**Date:** 2026-02-03
**Priority:** HIGHEST - Production-ready UI needed for both websites

---

## 📋 Website Architecture Overview

### Website 1: pryx.dev (Cloudflare Worker) - Public-Facing

**Purpose:** Public-facing website for Pryx Cloud broker

**Scope:** Minimal, focused, production-ready UX

**Pages to Build:**

1. **Landing Page** ()
   - Product overview with clear value proposition
   - "Get started" links for both Tailnet and Public modes
   - "How it works" section explaining architecture
   - Clean, modern design

2. **Auth Pages**
   -  - Device code flow login
     - Show device code entry field
     - Call  endpoint
     - Show activation instructions
   
   -  - QR code pairing UI
     - Show pairing code field
     - Display QR code (generate from pairing code)
     - Call  and 
     - Real-time pairing status updates

3. **User Dashboard** ()
   - Minimal device view (NO sync yet)
   - Device configurations display
   - Simple view only - no editing
   - Link to localhost control panel for full management
   - Status indicators (online/offline per device)

4. **Super Admin Dashboard** ()
   - User management (list all users, view telemetry)
   - Error logging display (from /api/telemetry/ingest)
   - System monitoring (health status, active sessions)
   - Use existing API routes, just need UI

**What NOT to Build on pryx.dev:**
   ❌ Full configuration management (that's on devices)
   ❌ Session sync between devices (future feature)
   ❌ Chat interface (that's on localhost website)

---

### Website 2: Localhost Gateway Control Panel

**Purpose:** Full gateway control + chat room, running locally

**Scope:** Full-featured, comprehensive control, connected everything

**Features Required:**

1. **Gateway Control Panel**
   - Start/Stop gateway service
   - View real-time logs (tail with filters)
   - Status monitoring (health, uptime, connected clients)
   - Configuration management:
     - Providers (add/remove/edit API keys)
     - Channels (Telegram, Discord, Slack setup)
     - MCP (add/enable/disable/configure servers)
     - Skills (install/enable/disable/manage)

2. **Chat Room**
   - Local web chat interface (full chat UI)
   - Multi-input support (CLI, TUI, Channels, Web)
   - Message history per session
   - Real-time updates (streaming output)
   - Presence indicators (who's watching this session)

3. **Connected Everything**
   - CLI + TUI integration
   - Localhost web chat
   - Channel message mirroring (show messages from Telegram/Discord/Slack)
   - Session management (create/delete/switch)
   - Real-time updates across all clients

4. **Session Management UI**
   - Session list with metadata (title, created/updated)
   - Create new session
   - Delete session
   - Switch between sessions
   - Export/import sessions (optional for v1)

5. **Better than OpenClaw**
   - Similar architecture but improved UX/DX
   - Cleaner, more intuitive UI
   - Better onboarding flow
   - Superior error handling and recovery

---

## 🎯 Implementation Priority

### Phase 1: pryx.dev Public Site (HIGHEST PRIORITY)

**Estimated Time:** 6-10 hours

**Tasks:**

**1.1: Landing Page** (1-2 hours)
- [ ] Modern hero section with value prop
- [ ] "Get started" CTA buttons
- [ ] "How it works" diagram
- [ ] Features overview
- [ ] Mobile responsive design
- [ ] Dark mode support

**1.2: Auth Pages** (2-3 hours)
- [ ] Device code login UI
- [ ] QR code pairing UI
- [ ] Real-time pairing status updates
- [ ] Error handling and recovery
- [ ] Mobile-friendly QR code display

**1.3: User Dashboard** (2-3 hours)
- [ ] Device list view
- [ ] Device status indicators (online/offline)
- [ ] Device configuration display
- [ ] Link to localhost control panel
- [ ] Minimal, clean design

**1.4: Super Admin Dashboard** (1-2 hours)
- [ ] User management interface
- [ ] Error logging display (telemetry viewer)
- [ ] System monitoring dashboard
- [ ] Use existing API endpoints

**Deployment:** Deploy to Cloudflare Workers (wrangler deploy)

---

### Phase 2: Localhost Control Panel (HIGH PRIORITY)

**Estimated Time:** 12-18 hours

**Tasks:**

**2.1: Gateway Control Panel** (4-6 hours)
- [ ] Start/Stop service controls
- [ ] Real-time logs viewer (with filters/search)
- [ ] Status monitoring dashboard (health, uptime, clients)
- [ ] Provider management UI (add/remove/edit)
- [ ] Channel configuration UI (setup wizard)
- [ ] MCP management UI
- [ ] Skills management UI

**2.2: Web Chat Room** (4-6 hours)
- [ ] Full chat UI (messages, input, history)
- [ ] Multi-input indicator (show messages from all sources)
- [ ] Message history per session
- [ ] Real-time streaming output display
- [ ] Presence indicators (who's watching)
- [ ] Markdown rendering for code/structured content

**2.3: Session Management UI** (2-3 hours)
- [ ] Session list view
- [ ] Create/delete/switch sessions
- [ ] Session metadata display
- [ ] Active session highlighting
- [ ] Export/import (optional)

**2.4: Integration with CLI/TUI** (1-2 hours)
- [ ] Web chat ↔ TUI integration
- [ ] Web chat ↔ CLI integration
- [ ] Channel message mirroring
- [ ] Real-time sync across all clients

**Deployment:** Serve from Go runtime or separate static server

---

### Phase 3: Pryx Cloud Broker (Already in Progress)

**Status:** Covered in existing requirements (UX_DX_REQUIREMENTS.md)
- WebSocket hub for device discovery
- QR code pairing
- Session mirroring (future)
- Presence tracking

**Estimated Time:** 8-12 hours (from previous plan)

---

## 🔧 Technical Architecture

### pryx.dev (Cloudflare Workers)

**Stack:**
- Frontend: Astro + React (or Svelte for better DX)
- Backend: Cloudflare Workers (Hono API + static pages)
- Deployment: Wrangler CLI

**Structure:**
````
apps/web/
├── src/
│   ├── routes/
│   │   ├── index.tsx         # Landing page
│   │   ├── auth/
│   │   │   ├── index.tsx    # Device code login
│   │   │   └── qr.tsx       # QR code pairing
│   │   ├── dashboard.tsx    # User dashboard
│   │   └── admin.tsx        # Super admin
│   └── components/
│       ├── Layout.tsx
│       ├── Hero.tsx
│       ├── DeviceList.tsx
│       └── TelemetryViewer.tsx
├── public/                 # Static assets
└── worker.ts              # Cloudflare Worker entry
```

### Localhost Control Panel

**Stack:**
- Frontend: React (or Svelte) + Vite
- Backend: Go runtime (existing HTTP server + WebSocket)
- Deployment: Serve from Go runtime or separate static server

**Structure:**
```
apps/control-panel/          # New directory
├── src/
│   ├── routes/
│   │   ├── dashboard.tsx   # Gateway control + chat room
│   │   ├── settings/
│   │   │   ├── providers.tsx    # Provider management
│   │   │   ├── channels.tsx      # Channel configuration
│   │   │   ├── mcp.tsx          # MCP servers
│   │   │   └── skills.tsx       # Skills management
│   │   └── sessions.tsx     # Session management
│   └── components/
│       ├── ChatRoom.tsx
│       ├── LogViewer.tsx
│       ├── StatusMonitor.tsx
│       └── PresenceIndicator.tsx
└── package.json
```

---

## 📊 Success Criteria

### pryx.dev (Public Site)

**Critical (Must Have):**
- [ ] Landing page is production-ready (mobile, dark mode)
- [ ] Device code login works end-to-end
- [ ] QR code pairing works with status updates
- [ ] User dashboard shows devices correctly
- [ ] Super admin can view users and errors

**High (Should Have):**
- [ ] Real-time pairing status updates (WebSocket)
- [ ] Device status auto-refreshes
- [ ] Telemetry viewer is filterable
- [ ] Analytics/basic metrics on admin dashboard

**Medium (Nice to Have):**
- [ ] User settings page (theme, preferences)
- [ ] API documentation viewer
- [ ] Changelog/history pages

---

### Localhost Control Panel

**Critical (Must Have):**
- [ ] Gateway can start/stop from web UI
- [ ] Real-time logs viewer works
- [ ] Provider management functional
- [ ] Channel setup wizard works
- [ ] Web chat room receives CLI/TUI messages
- [ ] Channel messages appear in web chat
- [ ] Session create/delete/switch works

**High (Should Have):**
- [ ] Log filtering and search works
- [ ] Status monitoring shows real-time data
- [ ] MCP server management works
- [ ] Skills install/enable/disable works
- [ ] Message history persists

**Medium (Nice to Have):**
- [ ] Message export
- [ ] Session templates
- [ ] Dark mode toggle
- [ ] Keyboard shortcuts

---

## 📁 Key Files to Create

### pryx.dev (Cloudflare Workers)
1.  - Landing page
2.  - Device code login
3.  - QR code pairing
4.  - User dashboard
5.  - Super admin
6.  - Site layout
7.  - Hero section
8.  - Device list

### Localhost Control Panel
1.  - New directory
2.  - Main dashboard
3.  - Providers
4.  - Channels
5.  - MCP servers
6.  - Skills
7.  - Sessions
8.  - Web chat room
9.  - Logs viewer
10.  - Status monitor

---

## 🚦 Current Blockers

1. **Time** - 26-40 hours of work remaining for both websites
2. **Decision** - React vs Svelte for pryx.dev (affects DX)
3. **Integration** - How localhost control panel connects to Go runtime
4. **QR Code Generation** - Need library for pryx.dev

---

## 📊 Production Readiness Impact

**Current Score:** 70%
**With pryx.dev:** 75%
**With Localhost Control Panel:** 85%
**With Both Websites:** 90%

**Combined with Other Requirements:**
- Telemetry fixed: 90%
- Cloud broker complete: 95%

**Target:** 95%+ (production ready)

**Total Estimated Time:** 26-40 hours (~1-2 weeks)

---

## 📝 Notes

### Key Clarifications
1. ✅ User dashboard = minimal device view (NO sync yet)
2. ✅ Super admin = user management + error logs (already have API)
3. ✅ Localhost site = full control + chat (separate from pryx.dev)
4. ✅ Focus on better UX/DX than OpenClaw (cleaner, more intuitive)

### Architecture Decision
This separation of concerns is smart:
- **pryx.dev** = Public, minimal, marketing + basic auth
- **Localhost** = Private, full-featured, tooling + chat
- Both can exist independently and serve different purposes

### Production Path
1. Build pryx.dev first (faster, easier, public)
2. Build localhost control panel (needs more time)
3. Connect both via Pryx Cloud broker (already planned)
4. Deliver production-ready Pryx v1

---

**Status:** 🟡 REQUIREMENTS CLARIFIED - READY FOR IMPLEMENTATION
