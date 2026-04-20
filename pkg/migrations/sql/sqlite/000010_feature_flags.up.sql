CREATE TABLE IF NOT EXISTS zatrano_feature_flags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    "key" TEXT NOT NULL UNIQUE,
    enabled INTEGER NOT NULL DEFAULT 0,
    rollout_percent INTEGER NOT NULL DEFAULT 0 CHECK (rollout_percent >= 0 AND rollout_percent <= 100),
    allowed_roles TEXT NOT NULL DEFAULT '[]',
    updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX IF NOT EXISTS idx_zatrano_feature_flags_key_lower ON zatrano_feature_flags (lower("key"));
