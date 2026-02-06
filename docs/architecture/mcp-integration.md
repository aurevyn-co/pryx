# MCP (Model Context Protocol) Integration

This document describes how Pryx integrates with the Model Context Protocol (MCP).

## Overview

Pryx implements a complete MCP client and server, enabling bidirectional communication with MCP tools/servers.

## Architecture

```
┌───────────────────────────────────────────────────────────┐
│                                                 │
│   Pryx Runtime (Go)                              │
│   ┌─────────────────────────┐                    │
│   │  MCP Manager        │                    │
│   │  - Discovery Service │                    │
│   │  - Config Manager    │                    │
│   │  - Tool Executor    │                    │
│   └─────────────────────────┘                    │
│                                                 │
│   ┌─────────────────────────┐                    │
│   │  MCP Servers        │                    │
│   │  - Bundled Tools (Browser, Clipboard, etc.) │
│   │  - External Servers (Custom MCP servers)  │
│   │  - JSON-RPC Transport                    │
│   └─────────────────────────┘                    │
│                                                 │
└───────────────────────────────────────────────────────────┘
                          ↕                   ↕
                  JSON-RPC over stdio (host ↔ runtime)
```

## Components

### 1. MCP Manager

Located in `apps/runtime/internal/mcp/manager.go`

- **Lifecycle Management**: Start, stop, restart, reload servers
- **Connection Tracking**: Monitor server health, auto-reconnect
- **Resource Management**: Manage shared resources (files, sockets)
- **Event Publishing**: Publish MCP events to event bus

### 2. Discovery Service

Located in `apps/runtime/internal/mcp/discovery/`

- **Curated Servers**: Pre-configured MCP servers from community
- **Categories**: Organized servers by category (Productivity, Development, AI, etc.)
- **Search**: Full-text search across all MCP servers
- **Validation**: Verify server manifests and configurations
- **Installation**: One-click install of discovered servers

### 3. Config Manager

Located in `apps/runtime/internal/mcp/config.go`

- **Server Configuration**: Add/manage MCP servers
- **Tool Filtering**: Configure which tools to enable/disable
- **Authentication**: Store auth tokens securely in vault/keychain
- **Environment Variables**: Configure MCP server environments
- **Settings File**: `~/.pryx/mcp/servers.json` (persisted config)

### 4. Tool Executor

Located in `apps/runtime/internal/mcp/executor.go`

- **Tool Calls**: Execute MCP tools with parameters
- **Streaming**: Handle streaming responses
- **Error Handling**: Proper error propagation
- **Timeout Management**: Configurable tool call timeouts
- **Progress Reporting**: Real-time progress updates

### 5. Client (for MCP Servers)

Located in `apps/runtime/internal/mcp/client.go`

- **JSON-RPC 2.0**: Transport for MCP communication
- **Client Registration**: Register with MCP servers
- **Tool Implementation**: Server-side MCP tools for other clients
- **Event Subscription**: Receive tool execution events
- **Resource Serving**: Serve resources to other clients

### 6. Bundled Tools

Located in `apps/runtime/internal/mcp/bundled/`

#### 6.1. Browser Tool
- Read web pages, extract content
- Take screenshots
- Navigate and click elements
- Execute JavaScript in browser context

#### 6.2. Clipboard Tool
- Read and write system clipboard
- Support multiple formats (text, HTML, images)

#### 6.3. Filesystem Tool
- Virtual filesystem for tool execution
- File read/write operations
- Safe sandboxing

#### 6.4. Shell Tool
- Execute shell commands
- Capture stdout/stderr
- Working directory management

#### 6.5. Screen Tool
- Terminal emulation
- Resize and scroll
- ANSI color support

#### 6.6. Log Tool
- Access application logs
- Tail logs in real-time
- Filter and search logs

### 7. Event Types

```go
const (
    EventMCPServerStarted    = "mcp:server_started"
    EventMCPServerStopped     = "mcp:server_stopped"
    EventMCPToolCalled       = "mcp:tool_called"
    EventMCPToolCompleted    = "mcp:tool_completed"
    EventMCPToolErrored       = "mcp:tool_errored"
    EventMCPResourceCreated  = "mcp:resource_created"
)
)
```

## Configuration

### Servers Configuration

Stored in `~/.pryx/mcp/servers.json`:

```json
{
  "servers": [
    {
      "id": "filesystem",
      "name": "Filesystem",
      "enabled": true,
      "command": "/usr/local/bin/rye",
      "transport": "stdio",
      "env": {"HOME": "/home/user"}
    },
    {
      "id": "sqlite",
      "name": "SQLite",
      "enabled": true,
      "command": "sqlite3",
      "transport": "stdio"
    }
  ]
}
```

## API Endpoints

```
# MCP Server Management
GET    /api/v1/mcp/servers              # List all MCP servers
POST   /api/v1/mcp/servers              # Add new MCP server
PUT    /api/v1/mcp/servers/{id}      # Update server config
DELETE /api/v1/mcp/servers/{id}      # Remove server

# MCP Discovery
GET    /api/v1/mcp/discovery/curated          # Get curated server list
GET    /api/v1/mcp/discovery/categories       # Get categories
POST   /api/v1/mcp/discovery/validate        # Validate server URL
POST   /api/v1/mcp/discovery/custom          # Add custom server
GET    /api/v1/mcp/discovery/custom          # List custom servers
DELETE /api/v1/mcp/discovery/custom/{id}   # Remove custom server

# MCP Tools
GET    /api/v1/mcp/tools                  # List available MCP tools
POST   /api/v1/mcp/tools/call           # Execute MCP tool
GET    /api/v1/mcp/tools/{id}           # Get tool schema
POST   /api/v1/mcp/tools/subscribe         # Subscribe to tool events

# Bundled Tools
POST   /api/v1/mcp/browser/execute        # Browser tool execution
POST   /api/v1/mcp/clipboard/read         # Clipboard read
POST   /api/v1/mcp/clipboard/write        # Clipboard write
POST   /api/v1/mcp/filesystem/operation   # Filesystem operations
POST   /api/v1/mcp/shell/execute          # Shell command execution
POST   /api/v1/mcp/screen/execute           # Terminal screen operations
POST   /api/v1/mcp/logs/tail             # Log tailing
```

## Tool Execution Flow

1. Runtime receives tool execution request via HTTP API or JSON-RPC
2. Executor validates tool input against schema
3. Tool server processes request
4. Response streamed back to runtime
5. Results published to event bus
6. Agent can use tool results in subsequent actions

## Security

- **Transport Security**: JSON-RPC 2.0 over stdio for host↔runtime communication
- **Sandboxing**: All bundled tools run in isolated environments
- **Access Control**: Tool execution scoped to user permissions
- **Resource Limits**: CPU, memory, and disk quotas per tool
- **Input Validation**: Schema validation for all tool inputs

## Best Practices

1. **Async Operations**: All tool calls are non-blocking with timeout handling
2. **Error Recovery**: Failed tool calls are retried with exponential backoff
3. **Stream Handling**: Use streaming for long-running operations
4. **Resource Cleanup**: Properly close resources after tool execution
5. **Event-Driven**: All MCP operations publish events for observability
6. **Configuration**: MCP servers can be hot-reloaded without restart

## Discovery Integration

Pryx uses a federated discovery system:
- **Curated Servers**: Community-maintained list from MCP discovery service
- **Custom Servers**: User-added servers (e.g., internal tools)
- **Search**: Full-text and category-based search across all registered servers
- **Validation**: Verify server compatibility before connection

## Future Enhancements

- [ ] Add MCP server clustering and load balancing
- [ ] Implement tool chaining (output of one tool as input to another)
- [ ] Add tool permissions system (fine-grained access control)
- [ ] Add MCP prompt templates for common workflows
- [ ] Implement tool caching for performance
- [ ] Add batch execution support for multiple tools
- [ ] Add streaming response preview for long operations
- [ ] Add tool marketplace (discover and install MCP servers)
- [ ] Implement federated discovery across multiple discovery services
