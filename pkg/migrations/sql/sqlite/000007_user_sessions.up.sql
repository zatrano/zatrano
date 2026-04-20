CREATE TABLE IF NOT EXISTS user_sessions (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id        INTEGER NOT NULL,
    session_token  VARCHAR(255) NOT NULL UNIQUE,
    refresh_token  VARCHAR(255) UNIQUE,
    ip_address     VARCHAR(45),
    user_agent     TEXT,
    device_info    TEXT,
    expires_at     TEXT NOT NULL,
    last_activity  TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    created_at     TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_session_token ON user_sessions (session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token ON user_sessions (refresh_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions (expires_at);
