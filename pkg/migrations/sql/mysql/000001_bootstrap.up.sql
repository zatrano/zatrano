CREATE TABLE IF NOT EXISTS zatrano_schema_migrations_meta (
    id         TINYINT PRIMARY KEY DEFAULT 1,
    note       VARCHAR(500) NOT NULL DEFAULT 'ZATRANO placeholder — replace with your domain tables',
    created_at DATETIME(6) NOT NULL DEFAULT (CURRENT_TIMESTAMP(6)),
    CONSTRAINT zatrano_single CHECK (id = 1)
);
