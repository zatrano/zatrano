CREATE TABLE IF NOT EXISTS user_sessions (
    id             INT AUTO_INCREMENT PRIMARY KEY,
    user_id        INT NOT NULL,
    session_token  VARCHAR(255) NOT NULL,
    refresh_token  VARCHAR(255) NULL,
    ip_address     VARCHAR(45) NULL,
    user_agent     TEXT NULL,
    device_info    JSON NULL,
    expires_at     DATETIME(6) NOT NULL,
    last_activity  DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    created_at     DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    UNIQUE KEY uq_user_sessions_session_token (session_token),
    UNIQUE KEY uq_user_sessions_refresh_token (refresh_token)
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_session_token ON user_sessions (session_token);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions (refresh_token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions (expires_at);
