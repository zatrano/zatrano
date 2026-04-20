IF EXISTS (SELECT 1 FROM sys.tables WHERE name = N'users' AND schema_id = SCHEMA_ID(N'dbo'))
BEGIN
    IF COL_LENGTH(N'dbo.users', N'email_verified_at') IS NULL
        ALTER TABLE dbo.users ADD email_verified_at DATETIME2(7) NULL;
    IF COL_LENGTH(N'dbo.users', N'email_verification_token') IS NULL
        ALTER TABLE dbo.users ADD email_verification_token NVARCHAR(255) NULL;
    IF COL_LENGTH(N'dbo.users', N'email_verification_expires_at') IS NULL
        ALTER TABLE dbo.users ADD email_verification_expires_at DATETIME2(7) NULL;
    IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_users_email_verification_token' AND object_id = OBJECT_ID(N'dbo.users'))
        CREATE INDEX idx_users_email_verification_token ON dbo.users (email_verification_token);
END
