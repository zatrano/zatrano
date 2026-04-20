IF OBJECT_ID(N'dbo.zatrano_feature_flags', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.zatrano_feature_flags (
        id BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        [key] NVARCHAR(255) NOT NULL,
        enabled BIT NOT NULL CONSTRAINT DF_zff_enabled DEFAULT 0,
        rollout_percent TINYINT NOT NULL CONSTRAINT DF_zff_rollout DEFAULT 0,
        allowed_roles NVARCHAR(MAX) NOT NULL CONSTRAINT DF_zff_roles DEFAULT N'[]',
        updated_at DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        CONSTRAINT CK_zff_rollout CHECK (rollout_percent >= 0 AND rollout_percent <= 100)
    );
    CREATE UNIQUE INDEX uq_zatrano_feature_flags_key ON dbo.zatrano_feature_flags ([key]);
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_zatrano_feature_flags_key_lower' AND object_id = OBJECT_ID(N'dbo.zatrano_feature_flags'))
    CREATE INDEX idx_zatrano_feature_flags_key_lower ON dbo.zatrano_feature_flags (LOWER([key]));
