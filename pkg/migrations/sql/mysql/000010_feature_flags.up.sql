CREATE TABLE IF NOT EXISTS zatrano_feature_flags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(255) NOT NULL,
    enabled TINYINT(1) NOT NULL DEFAULT 0,
    rollout_percent TINYINT NOT NULL DEFAULT 0,
    allowed_roles JSON NOT NULL DEFAULT (JSON_ARRAY()),
    updated_at DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT chk_rollout CHECK (rollout_percent >= 0 AND rollout_percent <= 100),
    UNIQUE KEY uq_zatrano_feature_flags_key (`key`)
);

CREATE INDEX idx_zatrano_feature_flags_key_lower ON zatrano_feature_flags ((LOWER(`key`)));
