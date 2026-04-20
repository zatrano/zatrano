CREATE TABLE IF NOT EXISTS api_keys (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    `key`      VARCHAR(64) NOT NULL,
    prefix     VARCHAR(8) NOT NULL,
    scopes     JSON DEFAULT (JSON_ARRAY()),
    expires_at DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    updated_at DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY uq_api_keys_key (`key`)
);

CREATE INDEX idx_api_keys_prefix ON api_keys (prefix);
CREATE INDEX idx_api_keys_expires_at ON api_keys (expires_at);

CREATE TABLE IF NOT EXISTS notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NULL,
    type VARCHAR(255) NOT NULL DEFAULT 'notification',
    subject VARCHAR(255) NULL,
    body TEXT NULL,
    data TEXT NULL,
    read_at DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    updated_at DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)) ON UPDATE CURRENT_TIMESTAMP(6)
);

CREATE INDEX idx_notifications_user_id ON notifications (user_id);
CREATE INDEX idx_notifications_read_at ON notifications (read_at);
CREATE INDEX idx_notifications_created_at ON notifications (created_at);
