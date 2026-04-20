CREATE TABLE IF NOT EXISTS zatrano_failed_jobs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id      VARCHAR(36)  NOT NULL,
    queue       VARCHAR(255) NOT NULL DEFAULT 'default',
    job_name    VARCHAR(255) NOT NULL,
    payload     TEXT         NOT NULL,
    error       TEXT         NOT NULL,
    stack_trace TEXT,
    failed_at   TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX IF NOT EXISTS idx_failed_jobs_job_id ON zatrano_failed_jobs (job_id);
CREATE INDEX IF NOT EXISTS idx_failed_jobs_queue ON zatrano_failed_jobs (queue);
