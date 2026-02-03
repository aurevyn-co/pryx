package store

const schema = `
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL,
	role TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS audit_log (
	id TEXT PRIMARY KEY,
	timestamp DATETIME NOT NULL,
	session_id TEXT,
	surface TEXT,
	tool TEXT,
	action TEXT NOT NULL,
	description TEXT,
	payload TEXT,
	cost TEXT,
	duration INTEGER,
	user_id TEXT,
	success BOOLEAN NOT NULL DEFAULT 1,
	error_msg TEXT,
	metadata TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_sessions_updated_at ON sessions(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_log(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_audit_session_id ON audit_log(session_id);
CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_surface ON audit_log(surface);
CREATE INDEX IF NOT EXISTS idx_audit_tool ON audit_log(tool);

-- Memory entries table for RAG system
CREATE TABLE IF NOT EXISTS memory_entries (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    date TEXT,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    access_count INTEGER DEFAULT 0,
    last_accessed DATETIME
);

-- Memory sources (which files/tools contributed to this memory)
CREATE TABLE IF NOT EXISTS memory_sources (
    id TEXT PRIMARY KEY,
    entry_id TEXT,
    source_type TEXT,
    source_path TEXT,
    FOREIGN KEY (entry_id) REFERENCES memory_entries(id) ON DELETE CASCADE
);

-- Vector embeddings (placeholder for future embedding provider plugin)
CREATE TABLE IF NOT EXISTS memory_vectors (
    entry_id TEXT PRIMARY KEY,
    embedding BLOB,
    FOREIGN KEY (entry_id) REFERENCES memory_entries(id) ON DELETE CASCADE
);

-- Indexes for memory tables
CREATE INDEX IF NOT EXISTS idx_memory_type ON memory_entries(type);
CREATE INDEX IF NOT EXISTS idx_memory_date ON memory_entries(date);
CREATE INDEX IF NOT EXISTS idx_memory_created_at ON memory_entries(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_memory_sources_entry_id ON memory_sources(entry_id);

-- Mesh pairing sessions (for QR code pairing)
CREATE TABLE IF NOT EXISTS mesh_pairing_sessions (
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    device_id TEXT NOT NULL,
    device_name TEXT NOT NULL,
    server_url TEXT NOT NULL,
    nonce TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mesh_pairing_code ON mesh_pairing_sessions(code);
CREATE INDEX IF NOT EXISTS idx_mesh_pairing_expires ON mesh_pairing_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_mesh_pairing_status ON mesh_pairing_sessions(status);

-- Mesh devices (paired devices)
CREATE TABLE IF NOT EXISTS mesh_devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    public_key TEXT NOT NULL,
    paired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    metadata TEXT
);

CREATE INDEX IF NOT EXISTS idx_mesh_devices_active ON mesh_devices(is_active);
CREATE INDEX IF NOT EXISTS idx_mesh_devices_paired ON mesh_devices(paired_at DESC);

-- Mesh sync events
CREATE TABLE IF NOT EXISTS mesh_sync_events (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    source_device_id TEXT NOT NULL,
    target_device_id TEXT,
    payload TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mesh_events_source ON mesh_sync_events(source_device_id);
CREATE INDEX IF NOT EXISTS idx_mesh_events_created ON mesh_sync_events(created_at DESC);
`
