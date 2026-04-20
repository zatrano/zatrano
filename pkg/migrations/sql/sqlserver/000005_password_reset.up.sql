IF OBJECT_ID(N'dbo.password_reset_tokens', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.password_reset_tokens (
        id         INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        email      NVARCHAR(255) NOT NULL,
        token      NVARCHAR(255) NOT NULL,
        expires_at DATETIME2(7) NOT NULL,
        created_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
    CREATE UNIQUE INDEX uq_password_reset_tokens_token ON dbo.password_reset_tokens (token);
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_password_reset_tokens_email' AND object_id = OBJECT_ID(N'dbo.password_reset_tokens'))
    CREATE INDEX idx_password_reset_tokens_email ON dbo.password_reset_tokens (email);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_password_reset_tokens_token' AND object_id = OBJECT_ID(N'dbo.password_reset_tokens'))
    CREATE INDEX idx_password_reset_tokens_token ON dbo.password_reset_tokens (token);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_password_reset_tokens_expires_at' AND object_id = OBJECT_ID(N'dbo.password_reset_tokens'))
    CREATE INDEX idx_password_reset_tokens_expires_at ON dbo.password_reset_tokens (expires_at);
