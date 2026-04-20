IF OBJECT_ID(N'dbo.zatrano_failed_jobs', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.zatrano_failed_jobs (
        id          BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        job_id      NVARCHAR(36)  NOT NULL,
        queue       NVARCHAR(255) NOT NULL CONSTRAINT DF_zfj_queue DEFAULT N'default',
        job_name    NVARCHAR(255) NOT NULL,
        payload     NVARCHAR(MAX) NOT NULL,
        error       NVARCHAR(MAX) NOT NULL,
        stack_trace NVARCHAR(MAX) NULL,
        failed_at   DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_failed_jobs_job_id' AND object_id = OBJECT_ID(N'dbo.zatrano_failed_jobs'))
    CREATE INDEX idx_failed_jobs_job_id ON dbo.zatrano_failed_jobs (job_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_failed_jobs_queue' AND object_id = OBJECT_ID(N'dbo.zatrano_failed_jobs'))
    CREATE INDEX idx_failed_jobs_queue ON dbo.zatrano_failed_jobs (queue);
