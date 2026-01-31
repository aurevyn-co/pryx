# Database Architecture Recommendations

## Executive Summary

**Recommendation: Keep SQLite as primary database, add Redis selectively for specific use cases.**

SQLite is sufficient for Pryx's current and near-term needs. Redis should be added only for specific high-performance or real-time features.

---

## Current State Analysis

### What Pryx Currently Stores (SQLite)

1. **Sessions** (`sessions` table)
   - Session ID, title, timestamps
   - Low write volume (one row per session)

2. **Messages** (`messages` table)
   - Message content, role, session relationship
   - Medium write volume (grows with chat history)
   - Indexed by session_id for fast retrieval

3. **Audit Log** (`audit_log` table)
   - All tool executions, API calls, user actions
   - High write volume (every operation logged)
   - Multiple indexes for querying

4. **Configuration** (YAML file)
   - Static settings, provider configs
   - Read-heavy, rarely written

5. **Catalog Cache** (JSON files)
   - models.dev data, OpenRouter cache
   - Refreshed periodically (24h TTL)

### Current Performance Characteristics

- **SQLite Connection Pool**: 25 max open, 25 max idle
- **Connection Lifetime**: 5 minutes
- **Local-First**: Single-user desktop application
- **Write Pattern**: Bursty (chat sessions), not sustained high throughput

---

## Redis Use Case Analysis

### 1. Provider Catalog Caching

**Current**: File-based JSON cache (`~/.pryx/cache/models.json`)

**Assessment**: ❌ **NOT NEEDED**

- Catalog is 1.4MB JSON (84 providers, 1417 models)
- Loaded once at startup, refreshed every 24h
- File I/O is sufficient for this use case
- Redis would add complexity without benefit

**Recommendation**: Keep file-based caching

---

### 2. Session State

**Current**: SQLite + in-memory Bus

**Assessment**: ⚠️ **MAYBE - For Multi-Device Mesh**

Current session state flow:
```
User Input → TUI → Runtime → SQLite → Bus → (Mesh broadcast)
```

SQLite handles this well for single-device use. However, for **mesh coordination** (multi-device sync):

- Session state needs to sync across devices
- WebSocket-based mesh already implemented
- Redis could help with:
  - Session presence (which device has active session)
  - Temporary session locks during sync
  - Cross-device session state cache

**Recommendation**: Add Redis **only when mesh scaling becomes an issue**

Priority: P2 (future enhancement)

---

### 3. Rate Limiting

**Current**: In-memory rate limiter (`internal/channels/ratelimit.go`)

**Assessment**: ❌ **NOT NEEDED**

Current implementation:
```go
type RateLimiter struct {
    tokens   float64
    rate     float64
    capacity float64
    last     time.Time
}
```

- Per-channel rate limiting (not global)
- In-memory is sufficient for single runtime instance
- No distributed rate limiting needed yet

**Recommendation**: Keep in-memory rate limiting

---

### 4. Pub/Sub for Real-Time Features

**Current**: In-memory Bus (`internal/bus/bus.go`)

**Assessment**: ⚠️ **CONSIDER FOR FUTURE**

Current Bus implementation:
- In-memory pub/sub with goroutines
- 100-event buffer per subscriber
- Drops events if subscriber is slow

**When Redis Pub/Sub would help:**

1. **Multi-instance runtime** (if Pryx ever runs distributed)
2. **External integrations** (webhooks, third-party services)
3. **Real-time analytics** (if added later)

**Current mesh coordination** uses WebSocket directly, not pub/sub.

**Recommendation**: Keep in-memory Bus for now, consider Redis for:
- External webhook delivery (reliability)
- Analytics pipeline (if built)

Priority: P2 (future enhancement)

---

### 5. MCP Tool Result Caching

**Current**: In-memory cache with TTL (`internal/mcp/manager.go`)

**Assessment**: ❌ **NOT NEEDED**

```go
type cachedTools struct {
    fetchedAt time.Time
    tools     []Tool
}
```

- Tool lists cached per MCP server
- Short TTL (5 minutes typical)
- In-memory is sufficient

**Recommendation**: Keep in-memory caching

---

### 6. Cost Tracking & Aggregation

**Current**: SQLite audit_log table

**Assessment**: ⚠️ **MAYBE - For Analytics**

Current cost tracking:
- Per-request cost stored in audit_log
- Aggregated on-demand for reports
- SQLite handles this fine for single-user

**When Redis would help:**

- Real-time cost dashboards (faster aggregation)
- Budget alerts (need fast counters)
- High-frequency cost tracking (1000+ requests/minute)

**Recommendation**: Add Redis counters **if** building real-time cost monitoring

Priority: P2 (future enhancement)

---

## SQLite Production Readiness

### ✅ SQLite is Sufficient For:

1. **Single-user desktop applications** (Pryx's primary use case)
2. **Write volumes < 1000 TPS** (Pryx: ~10-100 TPS during active use)
3. **Data sizes < 100GB** (Pryx: ~10-100MB typical)
4. **Local-first architecture** (no network dependency)

### ⚠️ SQLite Limitations:

1. **Concurrent writes**: Only one writer at a time
   - Mitigation: WAL mode (already enabled)
   - Impact: Minimal for single-user app

2. **No built-in replication**: For mesh sync
   - Mitigation: Application-level sync via WebSocket
   - Impact: Handled by mesh manager

3. **Query performance**: Complex aggregations can be slow
   - Mitigation: Proper indexing (already done)
   - Impact: Minimal for Pryx query patterns

---

## Recommended Architecture

### Phase 1: Current (SQLite Only)

```
┌─────────────────────────────────────────┐
│           Pryx Runtime                  │
│  ┌─────────────┐  ┌──────────────────┐  │
│  │   SQLite    │  │   In-Memory Bus  │  │
│  │  (Primary)  │  │   (Pub/Sub)      │  │
│  └─────────────┘  └──────────────────┘  │
│         │                  │            │
│  ┌─────────────┐  ┌──────────────────┐  │
│  │ File Cache  │  │ In-Memory Cache  │  │
│  │ (Catalogs)  │  │ (MCP Tools, etc) │  │
│  └─────────────┘  └──────────────────┘  │
└─────────────────────────────────────────┘
```

### Phase 2: Future (SQLite + Redis Selective)

```
┌─────────────────────────────────────────┐
│           Pryx Runtime                  │
│  ┌─────────────┐  ┌──────────────────┐  │
│  │   SQLite    │  │   In-Memory Bus  │  │
│  │  (Primary)  │  │   (Pub/Sub)      │  │
│  └─────────────┘  └──────────────────┘  │
│         │                  │            │
│  ┌─────────────┐  ┌──────────────────┐  │
│  │ File Cache  │  │ In-Memory Cache  │  │
│  │ (Catalogs)  │  │ (MCP Tools, etc) │  │
│  └─────────────┘  └──────────────────┘  │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │  Redis (Optional, Future)        │   │
│  │  - Session presence (mesh)       │   │
│  │  - Cost counters (analytics)     │   │
│  │  - Webhook delivery queue        │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

---

## Decision Matrix

| Feature | SQLite | Redis | Recommendation |
|---------|--------|-------|----------------|
| Session storage | ✅ | ❌ | SQLite |
| Message history | ✅ | ❌ | SQLite |
| Audit log | ✅ | ❌ | SQLite |
| Config | ✅ | ❌ | SQLite (file) |
| Catalog cache | ✅ | ❌ | File-based |
| MCP tool cache | ✅ | ❌ | In-memory |
| Rate limiting | ✅ | ❌ | In-memory |
| Pub/Sub (Bus) | ✅ | ❌ | In-memory |
| Session presence (mesh) | ⚠️ | ✅ | Redis (future) |
| Cost counters | ⚠️ | ✅ | Redis (future) |
| Webhook queue | ⚠️ | ✅ | Redis (future) |

---

## Implementation Guidelines

### If Adding Redis Later:

1. **Make it optional**: Pryx should work without Redis
2. **Use Redis for**: Ephemeral data, counters, queues
3. **Keep SQLite for**: Persistent data, audit trail, sessions
4. **Configuration**:
   ```yaml
   redis:
     enabled: false  # Default
     url: "redis://localhost:6379"
     use_for:
       - session_presence
       - cost_counters
       - webhook_queue
   ```

### Migration Path:

1. **Phase 1**: Implement Redis as optional dependency
2. **Phase 2**: Add feature flags for Redis-backed features
3. **Phase 3**: Enable Redis for specific deployments (enterprise, multi-device)

---

## Conclusion

**SQLite is the right choice for Pryx today.**

- Simplifies deployment (no external dependencies)
- Sufficient performance for single-user use
- Local-first architecture maintained
- Easy backup/restore (single file)

**Add Redis only when:**
- Mesh coordination needs scaling
- Real-time analytics dashboard built
- External webhook reliability required
- Multi-instance runtime deployment

**Estimated timeline**: 6-12 months before Redis becomes beneficial
