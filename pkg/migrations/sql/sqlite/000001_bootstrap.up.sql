CREATE TABLE IF NOT EXISTS zatrano_schema_migrations_meta (
    id         INTEGER PRIMARY KEY CHECK (id = 1) DEFAULT 1,
    note       TEXT NOT NULL DEFAULT 'ZATRANO placeholder — replace with your domain tables',
    created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);
