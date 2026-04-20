IF OBJECT_ID(N'dbo.zatrano_activity_logs', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.zatrano_activity_logs (
        id           BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        created_at   DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        user_id      NVARCHAR(255) NULL,
        subject_type NVARCHAR(255) NOT NULL,
        subject_id   NVARCHAR(255) NOT NULL,
        action       NVARCHAR(255) NOT NULL,
        changes      NVARCHAR(MAX) NULL,
        request_id   NVARCHAR(255) NULL,
        ip           NVARCHAR(255) NULL,
        metadata     NVARCHAR(MAX) NULL
    );
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'zatrano_activity_logs_subject_idx' AND object_id = OBJECT_ID(N'dbo.zatrano_activity_logs'))
    CREATE INDEX zatrano_activity_logs_subject_idx ON dbo.zatrano_activity_logs (subject_type, subject_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'zatrano_activity_logs_created_idx' AND object_id = OBJECT_ID(N'dbo.zatrano_activity_logs'))
    CREATE INDEX zatrano_activity_logs_created_idx ON dbo.zatrano_activity_logs (created_at DESC);

IF OBJECT_ID(N'dbo.zatrano_http_audit_logs', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.zatrano_http_audit_logs (
        id          BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        created_at  DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        user_id     NVARCHAR(255) NULL,
        method      NVARCHAR(16) NOT NULL,
        path        NVARCHAR(2048) NOT NULL,
        url_query   NVARCHAR(MAX) NULL,
        status      INT NOT NULL,
        duration_ms INT NOT NULL,
        request_id  NVARCHAR(255) NULL,
        ip          NVARCHAR(255) NULL
    );
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'zatrano_http_audit_logs_created_idx' AND object_id = OBJECT_ID(N'dbo.zatrano_http_audit_logs'))
    CREATE INDEX zatrano_http_audit_logs_created_idx ON dbo.zatrano_http_audit_logs (created_at DESC);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'zatrano_http_audit_logs_user_idx' AND object_id = OBJECT_ID(N'dbo.zatrano_http_audit_logs'))
    CREATE INDEX zatrano_http_audit_logs_user_idx ON dbo.zatrano_http_audit_logs (user_id);
