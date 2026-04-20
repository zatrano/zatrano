CREATE TABLE IF NOT EXISTS zatrano_activity_logs (
    id           BIGSERIAL PRIMARY KEY,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id      TEXT NULL,
    subject_type TEXT NOT NULL,
    subject_id   TEXT NOT NULL,
    action       TEXT NOT NULL,
    changes      JSONB NULL,
    request_id   TEXT NULL,
    ip           TEXT NULL,
    metadata     JSONB NULL
);

CREATE INDEX IF NOT EXISTS zatrano_activity_logs_subject_idx
    ON zatrano_activity_logs (subject_type, subject_id);
CREATE INDEX IF NOT EXISTS zatrano_activity_logs_created_idx
    ON zatrano_activity_logs (created_at DESC);

CREATE TABLE IF NOT EXISTS zatrano_http_audit_logs (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id     TEXT NULL,
    method      TEXT NOT NULL,
    path        TEXT NOT NULL,
    url_query   TEXT NULL,
    status      INT NOT NULL,
    duration_ms INT NOT NULL,
    request_id  TEXT NULL,
    ip          TEXT NULL
);

CREATE INDEX IF NOT EXISTS zatrano_http_audit_logs_created_idx
    ON zatrano_http_audit_logs (created_at DESC);
CREATE INDEX IF NOT EXISTS zatrano_http_audit_logs_user_idx
    ON zatrano_http_audit_logs (user_id);
