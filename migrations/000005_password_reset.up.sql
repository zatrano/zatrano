-- +migrate Up
CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_password_reset_tokens_email (email),
    INDEX idx_password_reset_tokens_token (token),
    INDEX idx_password_reset_tokens_expires_at (expires_at)
);

-- +migrate Down
DROP TABLE password_reset_tokens;