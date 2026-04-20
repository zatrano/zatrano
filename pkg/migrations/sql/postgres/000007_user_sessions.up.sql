CREATE TABLE IF NOT EXISTS user_sessions (
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER NOT NULL,
    session_token  VARCHAR(255) NOT NULL UNIQUE,
    refresh_token  VARCHAR(255) NULL UNIQUE,
    ip_address     INET,
    user_agent     TEXT,
    device_info    JSONB,
    expires_at     TIMESTAMPTZ NOT NULL,
    last_activity  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id         ON user_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_session_token   ON user_sessions (session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token   ON user_sessions (refresh_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at      ON user_sessions (expires_at);
