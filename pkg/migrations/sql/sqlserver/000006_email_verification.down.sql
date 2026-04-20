IF EXISTS (SELECT 1 FROM sys.tables WHERE name = N'users' AND schema_id = SCHEMA_ID(N'dbo'))
BEGIN
    IF EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_users_email_verification_token' AND object_id = OBJECT_ID(N'dbo.users'))
        DROP INDEX idx_users_email_verification_token ON dbo.users;
    IF COL_LENGTH(N'dbo.users', N'email_verification_expires_at') IS NOT NULL
        ALTER TABLE dbo.users DROP COLUMN email_verification_expires_at;
    IF COL_LENGTH(N'dbo.users', N'email_verification_token') IS NOT NULL
        ALTER TABLE dbo.users DROP COLUMN email_verification_token;
    IF COL_LENGTH(N'dbo.users', N'email_verified_at') IS NOT NULL
        ALTER TABLE dbo.users DROP COLUMN email_verified_at;
END
