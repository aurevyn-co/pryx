CREATE TABLE sessions_new (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token_hash TEXT UNIQUE,
    expires_at TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

INSERT INTO sessions_new (id, user_id, expires_at, created_at) 
SELECT id, user_id, expires_at, created_at FROM sessions;

DROP TABLE sessions;
ALTER TABLE sessions_new RENAME TO sessions;

CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
