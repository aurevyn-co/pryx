# Skills System

This document describes the Skills/Extensibility system in Pryx.

## Overview

Pryx supports external skills through an HTTP-based bridge, enabling skills to be written in any language (JavaScript, Python, Rust, etc.) and integrated as MCP tools.

## Architecture

```
┌──────────────────────────────────────────────────────┐
│   Pryx Runtime (Go)                              │
│   ┌─────────────────────────────────────────┐      │
│   │  Skills Bridge          │ HTTP Client  │
│   │                       │                │        │
│   │  - Discovery          │           │        │
│   │  - Installer          │   HTTP APIs  │
│   │  - Parser             │           │        │
│   │  - Registry           │           │        │
│   │  - Remote Execution    │           │        │
│   └─────────────────────────┘      │        │
└───────────────────────────────────────────────────────┘
                          ↕                    ↕
                  External Skills (HTTP)
                  ┌─────────────────────┐      ┌─────────────────────────┐
                  │  ┌──────────────┐    │      │     ┌──────────────┐    │
                  │  │  skill.json     │    │      │     │  skill.go       │    │
                  │  │  install.json    │    │      │     │  skill-server.go │    │
                  │  │  index.json      │    │      │     │  (MCP server)   │
                  │  └──────────────┘    │      │     │                  │
                  └─────────────────────────┘      └────────────────────────────┘
                          ↕                   ↕
```

## Core Components

### 1. Skills Bridge

Located in `apps/runtime/internal/skills/bridge.go`

- **Purpose**: HTTP client for communicating with external skill servers
- **Authentication**: Supports Bearer and API key authentication
- **Request/Response**: Structured HTTP communication
- **Error Handling**: Proper HTTP error propagation
- **Timeout Configurable**: 30 second default timeout

### 2. Skill Discovery

Located in `apps/runtime/internal/skills/discover.go`

- **Endpoint Discovery**: Automatic discovery from skill.json manifest
- **Package Detection**: Identifies skill packages (pako, zod, pkgjs)
- **Registry System**: Maps external skill packages to types
- **Installer System**: Download and install skill dependencies

### 3. Installer System

Located in `apps/runtime/internal/skills/installer/`

- **Dependency Resolution**: Handles npm-style dependencies
- **Download Caching**: Cache downloaded packages
- **Rollback Support**: Rollback failed installations
- **Installer Registry**: Plugin system for installers (npx, cargo, etc.)

### 4. Parser

Located in `apps/runtime/internal/skills/parser/`

- **Manifest Parser**: Parse skill.json manifests
- **Schema Validation**: Validate against known schemas
- **Type Safety**: Generate TypeScript type definitions

### 5. Registry System

Located in `apps/runtime/internal/skills/registry.go`

- **Package Registry**: Maps package names to metadata
- **Version Resolution**: Select compatible versions
- **Health Checks**: Ping registry servers
- **Mirror Support**: Fallback mirror URLs

### 6. Remote Execution

Located in `apps/runtime/internal/skills/remote.go`

- **Job Queue**: Queue for skill execution
- **Result Streaming**: Stream execution results
- **Cancellation**: Cancel running jobs
- **Timeout Handling**: Job timeouts

## Skill Manifest Format

### skill.json

```json
{
  "name": "My Skill",
  "version": "1.0.0",
  "description": "A helpful skill that does things",
  "author": "Author Name",
  "license": "MIT",
  "repository": "https://github.com/author/repo",
  "homepage": "https://myskill.com",
  "entrypoint": "index.js",
  "packages": {
    "pako": "2.0.0",
    "zod": "^3.0.0"
  },
  "api": {
    "type": "stdio",
    "version": "1.0.0"
  }
}
```

### install.json

```json
{
  "runtime": "node",
  "packages": ["pako"],
  "script": "node install"
}
```

## API Endpoints

```
# Skill Discovery
GET    /api/v1/skills/discover          # Discover skills from remote URL
GET    /api/v1/skills/packages           # Get supported packages
GET    /api/v1/skills/installers         # List available installers
POST   /api/v1/skills/install           # Install skill from URL

# Skill Bridge
GET    /api/v1/skills                       # List installed skills
POST   /api/v1/skills                       # Add skill from local path
PUT    /api/v1/skills/{id}                # Update skill config
DELETE /api/v1/skills/{id}                # Remove skill

# Remote Execution
POST   /api/v1/skills/execute               # Execute skill job
GET    /api/v1/skills/execute/{id}/status   # Get job status
GET    /api/v1/skills/execute/{id}/stream   # Stream job results
DELETE /api/v1/skills/execute/{id}              # Cancel job
```

## Supported Packages

### Data Compression
- `pako` - Fast zlib compressor
- `zod` - TypeScript-first schema validation

### Data Validation
- `pkgjs` - Parameter parsing library

## Security

- **Manifest Validation**: Validate skill.json before installation
- **Sandboxing**: Skills run in isolated environments
- **Input Validation**: Schema validation for skill inputs
- **Rate Limiting**: Control skill execution rate
- **Resource Quotas**: CPU, memory, disk limits per skill

## Integration with MCP

External skills can be published as MCP tools:

1. Skill registers as MCP server
2. Pryx discovers and connects to skill
3. Skill exposes tools via MCP tools/list
4. Tool execution streamed back to runtime
5. Results integrated into agent workflows

## Best Practices

1. **Skill Packaging**: Include proper manifest (skill.json) and installation instructions
2. **Dependency Pinning**: Use exact versions for stability
3. **Error Handling**: Skills should handle errors gracefully
4. **Resource Cleanup**: Skills should clean up resources after execution
5. **Timeout Management**: Skills should respect configured timeouts
6. **Logging**: Skills should log important operations
7. **Idempotency**: Operations should be idempotent where possible

## Future Enhancements

- [ ] Skill marketplace integration (discover & install from registry)
- [ ] Skill version management (automatic updates)
- [ ] Skill permissions system (fine-grained access control)
- [ ] Skill sandbox improvements (resource quotas, rate limiting)
- [ ] Skill dependency resolution (conflict detection)
- [ ] Skill health monitoring (ping, status checks)
- [ ] Skill analytics (usage statistics, performance metrics)
- [ ] Skill templates (boilerplate skills for common patterns)
- [ ] Multi-skill orchestration (skills calling other skills)
- [ ] Batch skill execution (execute multiple skills in one request)
- [ ] Streaming skill execution (long-running processes)
- [ ] Skill hot-reloading (reload skills without restart)
- [ ] Skill development toolkit (CLI tool for skill development)
- [ ] Skill testing framework (automated skill testing)
