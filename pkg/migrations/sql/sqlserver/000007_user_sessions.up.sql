IF OBJECT_ID(N'dbo.user_sessions', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.user_sessions (
        id             INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        user_id        INT NOT NULL,
        session_token  NVARCHAR(255) NOT NULL,
        refresh_token  NVARCHAR(255) NULL,
        ip_address     NVARCHAR(45) NULL,
        user_agent     NVARCHAR(MAX) NULL,
        device_info    NVARCHAR(MAX) NULL,
        expires_at     DATETIME2(7) NOT NULL,
        last_activity  DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        created_at     DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
    CREATE UNIQUE INDEX uq_user_sessions_session_token ON dbo.user_sessions (session_token);
    CREATE UNIQUE INDEX uq_user_sessions_refresh_token ON dbo.user_sessions (refresh_token) WHERE refresh_token IS NOT NULL;
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_user_sessions_user_id' AND object_id = OBJECT_ID(N'dbo.user_sessions'))
    CREATE INDEX idx_user_sessions_user_id ON dbo.user_sessions (user_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_user_sessions_session_token' AND object_id = OBJECT_ID(N'dbo.user_sessions'))
    CREATE INDEX idx_user_sessions_session_token ON dbo.user_sessions (session_token);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_user_sessions_refresh_token' AND object_id = OBJECT_ID(N'dbo.user_sessions'))
    CREATE INDEX idx_user_sessions_refresh_token ON dbo.user_sessions (refresh_token);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_user_sessions_expires_at' AND object_id = OBJECT_ID(N'dbo.user_sessions'))
    CREATE INDEX idx_user_sessions_expires_at ON dbo.user_sessions (expires_at);
