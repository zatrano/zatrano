CREATE TABLE IF NOT EXISTS zatrano_feature_flags (
    id BIGSERIAL PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    rollout_percent SMALLINT NOT NULL DEFAULT 0 CHECK (rollout_percent >= 0 AND rollout_percent <= 100),
    allowed_roles JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_zatrano_feature_flags_key_lower ON zatrano_feature_flags ((lower(key)));
