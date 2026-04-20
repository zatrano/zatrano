IF EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_zatrano_feature_flags_key_lower' AND object_id = OBJECT_ID(N'dbo.zatrano_feature_flags'))
    DROP INDEX idx_zatrano_feature_flags_key_lower ON dbo.zatrano_feature_flags;
DROP TABLE IF EXISTS dbo.zatrano_feature_flags;
