CREATE TABLE IF NOT EXISTS zatrano_activity_logs (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at   TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    user_id      TEXT,
    subject_type TEXT NOT NULL,
    subject_id   TEXT NOT NULL,
    action       TEXT NOT NULL,
    changes      TEXT,
    request_id   TEXT,
    ip           TEXT,
    metadata     TEXT
);

CREATE INDEX IF NOT EXISTS zatrano_activity_logs_subject_idx
    ON zatrano_activity_logs (subject_type, subject_id);
CREATE INDEX IF NOT EXISTS zatrano_activity_logs_created_idx
    ON zatrano_activity_logs (created_at);

CREATE TABLE IF NOT EXISTS zatrano_http_audit_logs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at  TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    user_id     TEXT,
    method      TEXT NOT NULL,
    path        TEXT NOT NULL,
    url_query   TEXT,
    status      INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    request_id  TEXT,
    ip          TEXT
);

CREATE INDEX IF NOT EXISTS zatrano_http_audit_logs_created_idx
    ON zatrano_http_audit_logs (created_at);
CREATE INDEX IF NOT EXISTS zatrano_http_audit_logs_user_idx
    ON zatrano_http_audit_logs (user_id);
