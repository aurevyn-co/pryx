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
`
