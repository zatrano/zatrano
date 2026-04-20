IF OBJECT_ID(N'dbo.refresh_tokens', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.refresh_tokens (
        id         INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        user_id    INT NOT NULL,
        token      NVARCHAR(255) NOT NULL,
        expires_at DATETIME2(7) NOT NULL,
        revoked    BIT NOT NULL CONSTRAINT DF_refresh_tokens_revoked DEFAULT 0,
        created_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
    CREATE UNIQUE INDEX uq_refresh_tokens_token ON dbo.refresh_tokens (token);
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_refresh_tokens_user_id' AND object_id = OBJECT_ID(N'dbo.refresh_tokens'))
    CREATE INDEX idx_refresh_tokens_user_id ON dbo.refresh_tokens (user_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_refresh_tokens_token' AND object_id = OBJECT_ID(N'dbo.refresh_tokens'))
    CREATE INDEX idx_refresh_tokens_token ON dbo.refresh_tokens (token);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_refresh_tokens_expires_at' AND object_id = OBJECT_ID(N'dbo.refresh_tokens'))
    CREATE INDEX idx_refresh_tokens_expires_at ON dbo.refresh_tokens (expires_at);
