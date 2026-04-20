IF OBJECT_ID(N'dbo.api_keys', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.api_keys (
        id         INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        name       NVARCHAR(100) NOT NULL,
        [key]      NVARCHAR(64) NOT NULL,
        prefix     NVARCHAR(8) NOT NULL,
        scopes     NVARCHAR(MAX) NOT NULL CONSTRAINT DF_api_keys_scopes DEFAULT N'[]',
        expires_at DATETIME2(7) NULL,
        created_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        updated_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
    CREATE UNIQUE INDEX uq_api_keys_key ON dbo.api_keys ([key]);
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_api_keys_prefix' AND object_id = OBJECT_ID(N'dbo.api_keys'))
    CREATE INDEX idx_api_keys_prefix ON dbo.api_keys (prefix);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_api_keys_expires_at' AND object_id = OBJECT_ID(N'dbo.api_keys'))
    CREATE INDEX idx_api_keys_expires_at ON dbo.api_keys (expires_at);

IF OBJECT_ID(N'dbo.notifications', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.notifications (
        id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        user_id INT NULL,
        type NVARCHAR(255) NOT NULL CONSTRAINT DF_notifications_type DEFAULT N'notification',
        subject NVARCHAR(255) NULL,
        body NVARCHAR(MAX) NULL,
        data NVARCHAR(MAX) NULL,
        read_at DATETIME2(7) NULL,
        created_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        updated_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_notifications_user_id' AND object_id = OBJECT_ID(N'dbo.notifications'))
    CREATE INDEX idx_notifications_user_id ON dbo.notifications (user_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_notifications_read_at' AND object_id = OBJECT_ID(N'dbo.notifications'))
    CREATE INDEX idx_notifications_read_at ON dbo.notifications (read_at);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_notifications_created_at' AND object_id = OBJECT_ID(N'dbo.notifications'))
    CREATE INDEX idx_notifications_created_at ON dbo.notifications (created_at);
