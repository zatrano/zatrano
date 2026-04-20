IF OBJECT_ID(N'dbo.zatrano_schema_migrations_meta', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.zatrano_schema_migrations_meta (
        id         SMALLINT NOT NULL CONSTRAINT PK_zatrano_schema_migrations_meta PRIMARY KEY DEFAULT 1,
        note       NVARCHAR(500) NOT NULL CONSTRAINT DF_zatrano_schema_meta_note DEFAULT N'ZATRANO placeholder — replace with your domain tables',
        created_at DATETIME2(7) NOT NULL CONSTRAINT DF_zatrano_schema_meta_created DEFAULT SYSUTCDATETIME(),
        CONSTRAINT CK_zatrano_single CHECK (id = 1)
    );
END
