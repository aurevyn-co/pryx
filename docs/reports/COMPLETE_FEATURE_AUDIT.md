# Pryx Complete Feature Audit & Testing Checklist

**Version:** 1.0.0  
**Date:** February 2, 2026  
**Scope:** Full user journey from installation to active usage  
**Status:** Production Readiness Assessment

---

## Executive Summary

This document provides a comprehensive audit of all Pryx features following the complete user journey:

```
Installation → First Run (Auth) → Provider Setup → MCP/Skills/Channels → Chat
```

### Critical Blockers Status

| Feature | Status | Priority |
|---------|--------|----------|
| Skills Install | ✅ IMPLEMENTED | P0 |
| Skills Uninstall | ✅ IMPLEMENTED | P0 |
| OAuth PKCE | ✅ IMPLEMENTED | P0 |
| Rate Limiting | ✅ IMPLEMENTED | P1 |
| Token Refresh | ✅ ALREADY EXISTS | P0 |

---

## Phase 1: Installation & First Run

### 1.1 Installation Methods

#### Test Case: macOS Homebrew Install
**Steps:**
1. Run `brew install irfndi/pryx/pryx`
2. Verify binary installed to `/usr/local/bin/pryx`
3. Run `pryx --version`

**Expected:**
- Binary accessible in PATH
- Version displayed correctly
- No dependency errors

**Test Status:** ⬜ NOT TESTED

#### Test Case: Linux Install Script (Ubuntu/Debian)
**Steps:**
1. Run `curl -fsSL https://get.pryx.ai/install.sh | bash`
2. Verify `~/.pryx/` directory created
3. Check binary in `~/.local/bin/` or `/usr/local/bin/`

**Expected:**
- Directory structure created:
  - `~/.pryx/`
  - `~/.pryx/logs/`
  - `~/.pryx/skills/`
  - `~/.pryx/channels/`
- Binary executable

**Test Status:** ⬜ NOT TESTED

#### Test Case: Linux RPM Install (Fedora/RHEL)
**Steps:**
1. Download RPM from releases
2. Run `sudo dnf install pryx-*.rpm`
3. Verify installation

**Expected:**
- Package installs without conflicts
- Binary in `/usr/bin/pryx`

**Test Status:** ⬜ NOT TESTED

#### Test Case: Windows Winget Install
**Steps:**
1. Run `winget install pryx`
2. Verify installation in Program Files
3. Test `pryx` command in PowerShell

**Expected:**
- Application appears in Start Menu
- Binary accessible from command line
- No DLL errors

**Test Status:** ⬜ NOT TESTED

#### Test Case: Direct Binary Download
**Steps:**
1. Download from GitHub releases
2. Extract archive
3. Run `./pryx --version`

**Expected:**
- Runs without installation
- Self-contained binary

**Test Status:** ⬜ NOT TESTED

### 1.2 First Run Initialization

#### Test Case: Fresh Install - Directory Structure
**Steps:**
1. Delete `~/.pryx/` if exists
2. Run `pryx` for first time
3. Check created directories

**Expected Directories:**
```
~/.pryx/
├── config.json          # Configuration file
├── pryx.db             # SQLite database
├── logs/               # Log files
│   ├── runtime.log
│   └── tui.log
├── skills/             # Installed skills
└── channels/           # Channel configs
```

**Test Status:** ⬜ NOT TESTED

#### Test Case: Database Initialization
**Steps:**
1. First run
2. Check `~/.pryx/pryx.db` exists
3. Verify tables created

**Expected:**
- Database file exists
- Schema properly initialized
- Migrations applied

**Test Status:** ⬜ NOT TESTED

#### Test Case: Config File Creation
**Steps:**
1. First run
2. Check `~/.pryx/config.json`
3. Verify default values

**Expected Defaults:**
```json
{
  "model_provider": "",
  "model_name": "",
  "max_messages_per_session": 1000,
  "enable_telemetry": true,
  "enable_memory_profiling": false
}
```

**Test Status:** ⬜ NOT TESTED

---

## Phase 2: Authentication (CLI & TUI)

### 2.1 Cloud Login - TUI Flow

#### Test Case: First Run Shows Login Screen
**Steps:**
1. Fresh install
2. Start TUI with `pryx`
3. Observe initial screen

**Expected:**
- Step 0: "Pryx Cloud Login" displayed
- User code and verification URL shown
- Instructions to open browser
- "Press S to skip (offline)" option

**Test Status:** ⬜ NOT TESTED

#### Test Case: Successful OAuth Device Flow
**Steps:**
1. Start TUI
2. Press Enter to start login
3. Open browser at verification URL
4. Enter user code
5. Authorize in browser
6. Return to TUI, press Enter to poll

**Expected:**
- Device code generated
- User code displayed (format: XXXX-XXXX)
- Verification URL opens successfully
- Poll completes successfully
- Advances to Step 1 (Provider selection)
- Token stored in OS keychain

**Test Status:** ⬜ NOT TESTED

#### Test Case: OAuth with PKCE (Security)
**Steps:**
1. Start login flow
2. Verify PKCE parameters in backend

**Expected:**
- Code verifier generated (128 chars)
- Code challenge generated (S256)
- Challenge sent to auth endpoint
- Verifier used in token exchange

**Test Status:** ⬜ NOT TESTED (Implemented, needs verification)

#### Test Case: Cancel During Login
**Steps:**
1. Start login
2. Press Ctrl+C or wait for timeout

**Expected:**
- Graceful cancellation
- Returns to offline option
- No hanging processes

**Test Status:** ⬜ NOT TESTED

#### Test Case: Skip Cloud Login (Offline Mode)
**Steps:**
1. Start TUI
2. Press 'S' to skip

**Expected:**
- Skips to Step 1 (Provider setup)
- Can configure local providers (Ollama)
- Cloud features disabled

**Test Status:** ⬜ NOT TESTED

#### Test Case: Invalid/Expired User Code
**Steps:**
1. Start login
2. Wait for expiration (default 10 min)
3. Try to complete in browser

**Expected:**
- Clear error message
- Option to restart login
- No crash

**Test Status:** ⬜ NOT TESTED

### 2.2 Cloud Login - CLI Flow

#### Test Case: CLI Login Command
**Steps:**
1. Run `pryx-core login`
2. Follow prompts

**Expected:**
- Same flow as TUI
- Displays URL and code
- Polls for completion
- Stores token in keychain

**Test Status:** ⬜ NOT TESTED

### 2.3 Token Management

#### Test Case: Token Storage
**Steps:**
1. Complete login
2. Verify token in keychain

**Expected:**
- macOS: In Keychain Access
- Linux: In secret service/keyring
- Windows: In Credential Manager
- Token not in plaintext files

**Test Status:** ⬜ NOT TESTED

#### Test Case: Token Refresh (OAuth Providers)
**Steps:**
1. Login with OAuth provider (e.g., Google)
2. Wait for token to approach expiry
3. Use provider

**Expected:**
- Automatic token refresh before expiry
- New access token obtained
- Refresh token rotation (if applicable)
- No user interruption

**Test Status:** ⬜ NOT TESTED (Already implemented)

#### Test Case: Logout
**Steps:**
1. Run `pryx-core config remove cloud_access_token`
2. Or use TUI settings

**Expected:**
- Token removed from keychain
- Cloud features disabled
- Can login again

**Test Status:** ⬜ NOT TESTED

---

## Phase 3: AI Provider Setup

### 3.1 Provider Configuration - TUI

#### Test Case: Provider Selection Screen
**Steps:**
1. Complete login (or skip)
2. Step 1: Provider selection

**Expected:**
- List of 84+ providers displayed
- Provider names and API key requirements shown
- Navigation with ↑↓ arrows
- Select with Enter

**Test Status:** ⬜ NOT TESTED

#### Test Case: Provider with API Key (OpenAI)
**Steps:**
1. Select "OpenAI"
2. Step 2: Choose model
3. Step 3: Enter API key

**Expected:**
- Models fetched from models.dev
- API key input masked
- Key validated before saving
- Stored in OS keychain

**Test Status:** ⬜ NOT TESTED

#### Test Case: Local Provider (Ollama)
**Steps:**
1. Select "Ollama"
2. Optional: enter custom endpoint
3. No API key required

**Expected:**
- Works offline
- Default endpoint: localhost:11434
- Can specify custom endpoint
- Models fetched from Ollama

**Test Status:** ⬜ NOT TESTED

#### Test Case: OAuth Provider (Google)
**Steps:**
1. Select "Google"
2. Browser opens for OAuth
3. Complete authorization
4. Token stored

**Expected:**
- OAuth flow initiated
- Browser opens automatically
- Token saved to keychain
- Models available after auth

**Test Status:** ⬜ NOT TESTED

#### Test Case: Invalid API Key
**Steps:**
1. Select provider
2. Enter invalid key
3. Submit

**Expected:**
- Validation error shown
- Prompt to re-enter
- Helpful error message

**Test Status:** ⬜ NOT TESTED

#### Test Case: Provider Test Connection
**Steps:**
1. Configure provider
2. Run `pryx-core provider test <name>`

**Expected:**
- Connection test performed
- Success/failure message
- Details on failure

**Test Status:** ⬜ NOT TESTED

### 3.2 Provider Configuration - CLI

#### Test Case: Add Provider
**Steps:**
```bash
pryx-core provider add openai
pryx-core provider set-key openai
# Enter API key
```

**Expected:**
- Provider added to config
- Key stored securely
- Can switch to it

**Test Status:** ⬜ NOT TESTED

#### Test Case: List Providers
**Steps:**
`pryx-core provider list`

**Expected:**
- Shows all configured providers
- Indicates active provider
- Shows key status (configured/missing)

**Test Status:** ⬜ NOT TESTED

#### Test Case: Switch Provider
**Steps:**
`pryx-core provider use anthropic`

**Expected:**
- Active provider changed
- Used for new chats
- Persisted in config

**Test Status:** ⬜ NOT TESTED

#### Test Case: Remove Provider
**Steps:**
`pryx-core provider remove openai`

**Expected:**
- Provider removed from config
- API key deleted from keychain
- Confirmation prompt

**Test Status:** ⬜ NOT TESTED

---

## Phase 4: MCP, Skills & Channels Setup

### 4.1 MCP Server Management

#### Test Case: List MCP Servers
**Steps:**
`pryx-core mcp list`

**Expected:**
- Shows available MCP servers
- Status (enabled/disabled)
- Can add/remove

**Test Status:** ⬜ NOT TESTED

#### Test Case: Add MCP Server
**Steps:**
```bash
pryx-core mcp add filesystem
# Configure allowed directories
```

**Expected:**
- Server added to config
- Security validation performed
- Can enable/disable

**Test Status:** ⬜ NOT TESTED

#### Test Case: MCP Security Validation
**Steps:**
1. Add untrusted MCP server
2. Observe security warnings

**Expected:**
- Risk rating displayed (A-F)
- Warnings for HTTP/non-verified
- User confirmation required

**Test Status:** ⬜ NOT TESTED

### 4.2 Skills Management

#### Test Case: List Skills
**Steps:**
`pryx-core skills list`

**Expected:**
- Shows bundled, managed, workspace skills
- Status indicators (✓ enabled, ⚠ issues)
- Eligibility check

**Test Status:** ⬜ NOT TESTED

#### Test Case: Install Skill from Bundled
**Steps:**
```bash
pryx-core skills install web-search --from bundled/web-search
```

**Expected:**
- ✅ IMPLEMENTED
- Copies from bundled to managed
- Runs installers if defined
- Success message

**Test Status:** ⬜ NOT TESTED (Implemented)

#### Test Case: Install Skill from Path
**Steps:**
```bash
pryx-core skills install my-skill --from /path/to/skill
```

**Expected:**
- ✅ IMPLEMENTED
- Copies SKILL.md to managed directory
- Validates skill structure
- Can enable after install

**Test Status:** ⬜ NOT TESTED (Implemented)

#### Test Case: Uninstall Skill
**Steps:**
```bash
pryx-core skills uninstall web-search
```

**Expected:**
- ✅ IMPLEMENTED
- Removes from managed directory
- Disables if enabled
- Confirmation for non-managed

**Test Status:** ⬜ NOT TESTED (Implemented)

#### Test Case: Enable/Disable Skill
**Steps:**
```bash
pryx-core skills enable web-search
pryx-core skills disable web-search
```

**Expected:**
- Toggles enabled state
- Persisted in enabled.json
- Affects agent behavior

**Test Status:** ⬜ NOT TESTED

#### Test Case: Check Skill Health
**Steps:**
`pryx-core skills check`

**Expected:**
- Validates all skills
- Reports missing fields
- Checks required binaries/env
- Summary of issues

**Test Status:** ⬜ NOT TESTED

### 4.3 Channel Integration

#### Test Case: Telegram Bot Setup
**Steps:**
1. Create bot with @BotFather
2. `pryx-core channel add telegram mybot --token <token>`
3. `pryx-core channel enable mybot`
4. Test with `/start`

**Expected:**
- Bot responds to messages
- Webhook/polling configured
- Session tracking works
- Messages relayed to agent

**Test Status:** ⬜ NOT TESTED

#### Test Case: Discord Bot Setup
**Steps:**
1. Create Discord app
2. Get bot token
3. `pryx-core channel add discord mybot --token <token>`
4. Invite bot to server
5. Test with `/chat hello`

**Expected:**
- Bot appears online
- Responds to slash commands
- Messages processed
- Session tracking

**Test Status:** ⬜ NOT TESTED

#### Test Case: Slack App Setup
**Steps:**
1. Create Slack app
2. Configure bot scopes
3. `pryx-core channel add slack myapp --token <token>`
4. Test in channel

**Expected:**
- Bot responds to mentions
- DMs work
- Session tracking
- Webhook verification

**Test Status:** ⬜ NOT TESTED

#### Test Case: Channel Status Check
**Steps:**
`pryx-core channel status`

**Expected:**
- Shows all channels
- Connection status
- Last activity
- Error states

**Test Status:** ⬜ NOT TESTED

---

## Phase 5: Chat & Messaging

### 5.1 TUI Chat Interface

#### Test Case: Start New Chat
**Steps:**
1. Open TUI
2. Press `/` for command palette
3. Select "Chat" or press `1`
4. Type message
5. Press Enter

**Expected:**
- Message sent to runtime
- Agent processes request
- Response displayed
- Session created

**Test Status:** ⬜ NOT TESTED

#### Test Case: Multi-turn Conversation
**Steps:**
1. Start chat
2. Send: "My name is Alice"
3. Send: "What's my name?"

**Expected:**
- Context maintained
- Agent remembers "Alice"
- Session history preserved

**Test Status:** ⬜ NOT TESTED

#### Test Case: View Session History
**Steps:**
1. Press `/`
2. Select "Sessions" or press `2`
3. Browse sessions

**Expected:**
- List of past sessions
- Timestamps
- Can resume session

**Test Status:** ⬜ NOT TESTED

#### Test Case: Resume Session
**Steps:**
1. Go to Sessions view
2. Select previous session
3. Press Enter

**Expected:**
- Chat view opens
- Previous context loaded
- Can continue conversation

**Test Status:** ⬜ NOT TESTED

#### Test Case: Tool Execution in Chat
**Steps:**
1. Enable skill with tools
2. Chat: "Search for X"
3. Observe tool execution

**Expected:**
- Tool called automatically
- Approval shown (if required)
- Results integrated
- Cost tracked

**Test Status:** ⬜ NOT TESTED

### 5.2 Channel Chat

#### Test Case: Telegram Chat
**Steps:**
1. Message bot in Telegram
2. Bot responds via Pryx

**Expected:**
- Message received
- Agent processes
- Response sent
- Session tracked

**Test Status:** ⬜ NOT TESTED

#### Test Case: Discord Chat
**Steps:**
1. Use `/chat` command in Discord
2. Or mention bot

**Expected:**
- Command recognized
- Agent responds
- Embeds formatted
- Session tracked

**Test Status:** ⬜ NOT TESTED

#### Test Case: Slack Chat
**Steps:**
1. Mention bot in channel
2. Or DM the bot

**Expected:**
- Message processed
- Thread replies
- Session tracked
- File attachments handled

**Test Status:** ⬜ NOT TESTED

### 5.3 Session Management

#### Test Case: List Sessions
**Steps:**
`pryx-core session list`

**Expected:**
- Shows all sessions
- IDs, titles, timestamps
- Can export/delete

**Test Status:** ⬜ NOT TESTED

#### Test Case: Export Session
**Steps:**
`pryx-core session export <id> --format json`

**Expected:**
- Exports to file
- JSON/markdown format
- Includes all messages

**Test Status:** ⬜ NOT TESTED

#### Test Case: Delete Session
**Steps:**
`pryx-core session delete <id>`

**Expected:**
- Removes from database
- Confirmation prompt
- Irreversible

**Test Status:** ⬜ NOT TESTED

---

## Phase 6: Settings & Configuration

### 6.1 Configuration Management

#### Test Case: View All Config
**Steps:**
`pryx-core config list`

**Expected:**
- All settings displayed
- Current values shown
- Sources indicated

**Test Status:** ⬜ NOT TESTED

#### Test Case: Get/Set Config
**Steps:**
```bash
pryx-core config get max_messages_per_session
pryx-core config set max_messages_per_session 500
```

**Expected:**
- Value retrieved
- Value updated
- Persisted to config.json

**Test Status:** ⬜ NOT TESTED

### 6.2 Cost Tracking

#### Test Case: View Cost Summary
**Steps:**
`pryx-core cost summary`

**Expected:**
- Total cost displayed
- Breakdown by provider
- Token usage stats

**Test Status:** ⬜ NOT TESTED

#### Test Case: Daily Cost Breakdown
**Steps:**
`pryx-core cost daily 7`

**Expected:**
- Last 7 days breakdown
- Per-day costs
- Trends shown

**Test Status:** ⬜ NOT TESTED

#### Test Case: Set Budget
**Steps:**
`pryx-core cost budget set 100`

**Expected:**
- Budget configured
- Warnings at 80%
- Blocks at 100% (optional)

**Test Status:** ⬜ NOT TESTED

### 6.3 System Diagnostics

#### Test Case: Run Doctor
**Steps:**
`pryx-core doctor`

**Expected:**
- Database check
- Keychain check
- Provider connectivity
- Channel health
- MCP status
- Summary report

**Test Status:** ⬜ NOT TESTED

---

## Edge Cases & Error Handling

### Error Scenarios

#### Test Case: Runtime Not Running
**Steps:**
1. Kill runtime process
2. Try to use TUI

**Expected:**
- "Disconnected" shown
- Auto-reconnect attempts
- Clear error message
- Option to restart

**Test Status:** ⬜ NOT TESTED

#### Test Case: Network Failure
**Steps:**
1. Disconnect network
2. Try to chat

**Expected:**
- Offline mode detection
- Clear error message
- Retry option
- Local providers still work

**Test Status:** ⬜ NOT TESTED

#### Test Case: Invalid Provider Config
**Steps:**
1. Corrupt config.json
2. Start Pryx

**Expected:**
- Detects corruption
- Offers reset
- Backup old config
- Clean state

**Test Status:** ⬜ NOT TESTED

#### Test Case: Keychain Locked
**Steps:**
1. Lock OS keychain
2. Try to access API key

**Expected:**
- Prompt to unlock
- Fallback to file (if configured)
- Clear error
- Instructions to fix

**Test Status:** ⬜ NOT TESTED

#### Test Case: Rate Limit Hit
**Steps:**
1. Send many requests quickly
2. Observe rate limiting

**Expected:**
- 429 Too Many Requests
- Clear message
- Retry-after header
- No crash

**Test Status:** ⬜ NOT TESTED (Implemented)

---

## Security Audit

### Security Features

| Feature | Status | Notes |
|---------|--------|-------|
| OAuth PKCE | ✅ IMPLEMENTED | RFC 7636 compliance |
| Token Refresh | ✅ EXISTS | Automatic rotation |
| Rate Limiting | ✅ IMPLEMENTED | 10 req/sec default |
| OS Keychain | ✅ EXISTS | Secure secret storage |
| Input Validation | ✅ EXISTS | 200+ validators |
| MCP Sandboxing | ✅ EXISTS | Risk ratings A-F |
| Audit Logging | ✅ EXISTS | PII redaction |
| Policy Engine | ✅ EXISTS | Tool approval |
| Argon2id | ✅ EXISTS | Password hashing |

---

## Testing Checklist Summary

### Must Test Before Release (P0)

- [ ] Installation on macOS (Homebrew)
- [ ] Installation on Linux (script)
- [ ] First run creates directories
- [ ] OAuth device flow (TUI)
- [ ] OAuth PKCE implementation
- [ ] Provider setup (OpenAI with API key)
- [ ] Provider setup (Ollama local)
- [ ] Skills install/uninstall
- [ ] MCP server add/enable
- [ ] Telegram channel setup
- [ ] Basic chat in TUI
- [ ] Multi-turn conversation
- [ ] Session persistence
- [ ] Rate limiting works

### Should Test (P1)

- [ ] Installation on Windows
- [ ] OAuth with Google
- [ ] Token refresh
- [ ] Discord channel setup
- [ ] Slack channel setup
- [ ] Tool execution in chat
- [ ] Cost tracking
- [ ] Session export
- [ ] Doctor command
- [ ] Error handling scenarios

### Nice to Test (P2)

- [ ] RPM package install
- [ ] Direct binary download
- [ ] Multiple provider switching
- [ ] Advanced MCP configurations
- [ ] Channel webhook verification
- [ ] Budget alerts
- [ ] Performance under load

---

## Production Readiness Score

| Category | Score | Status |
|----------|-------|--------|
| **Core Features** | 90% | ✅ Good |
| **Security** | 95% | ✅ Excellent |
| **User Experience** | 75% | ⚠️ Needs Testing |
| **Documentation** | 60% | ⚠️ Incomplete |
| **Testing Coverage** | 40% | ❌ Poor |
| **Overall** | 72% | ⚠️ Not Ready |

---

## Blockers for Production

### Critical (Must Fix)

1. **Testing** - Most features not manually tested
2. **Installation** - Scripts not validated on clean systems
3. **Documentation** - User guides incomplete

### High Priority (Should Fix)

1. **Error Messages** - Need more user-friendly errors
2. **Onboarding** - First-run experience needs polish
3. **Recovery** - Better handling of corrupted state

### Medium Priority (Nice to Have)

1. **Performance** - Optimize startup time
2. **Telemetry** - Add usage analytics (opt-in)
3. **Auto-update** - Check for updates

---

## Next Steps

1. **Execute Test Plan** - Run all ⬜ NOT TESTED cases
2. **Fix Issues** - Address any bugs found
3. **Documentation** - Complete user guides
4. **Release** - Tag v1.0.0 when all P0 tests pass

---

*Document created by Sisyphus Agent*  
*Last updated: February 2, 2026*
