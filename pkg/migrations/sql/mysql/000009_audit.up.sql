CREATE TABLE IF NOT EXISTS zatrano_activity_logs (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    created_at   DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    user_id      VARCHAR(255) NULL,
    subject_type VARCHAR(255) NOT NULL,
    subject_id   VARCHAR(255) NOT NULL,
    action       VARCHAR(255) NOT NULL,
    changes      JSON NULL,
    request_id   VARCHAR(255) NULL,
    ip           VARCHAR(255) NULL,
    metadata     JSON NULL
);

CREATE INDEX zatrano_activity_logs_subject_idx ON zatrano_activity_logs (subject_type, subject_id);
CREATE INDEX zatrano_activity_logs_created_idx ON zatrano_activity_logs (created_at);

CREATE TABLE IF NOT EXISTS zatrano_http_audit_logs (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    user_id     VARCHAR(255) NULL,
    method      VARCHAR(16) NOT NULL,
    path        VARCHAR(2048) NOT NULL,
    url_query   TEXT NULL,
    status      INT NOT NULL,
    duration_ms INT NOT NULL,
    request_id  VARCHAR(255) NULL,
    ip          VARCHAR(255) NULL
);

CREATE INDEX zatrano_http_audit_logs_created_idx ON zatrano_http_audit_logs (created_at);
CREATE INDEX zatrano_http_audit_logs_user_idx ON zatrano_http_audit_logs (user_id);
