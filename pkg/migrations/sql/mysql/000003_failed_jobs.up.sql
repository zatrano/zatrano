CREATE TABLE IF NOT EXISTS zatrano_failed_jobs (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    job_id      VARCHAR(36)  NOT NULL,
    queue       VARCHAR(255) NOT NULL DEFAULT 'default',
    job_name    VARCHAR(255) NOT NULL,
    payload     TEXT         NOT NULL,
    error       TEXT         NOT NULL,
    stack_trace TEXT,
    failed_at   DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6))
);

CREATE INDEX idx_failed_jobs_job_id ON zatrano_failed_jobs (job_id);
CREATE INDEX idx_failed_jobs_queue ON zatrano_failed_jobs (queue);
