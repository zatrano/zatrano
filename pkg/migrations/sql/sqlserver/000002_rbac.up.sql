IF OBJECT_ID(N'dbo.roles', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.roles (
        id          INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        name        NVARCHAR(100) NOT NULL,
        description NVARCHAR(255) NULL DEFAULT N'',
        created_at  DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        updated_at  DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
    CREATE UNIQUE INDEX uq_roles_name ON dbo.roles (name);
END

IF OBJECT_ID(N'dbo.permissions', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.permissions (
        id          INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
        name        NVARCHAR(100) NOT NULL,
        description NVARCHAR(255) NULL DEFAULT N'',
        created_at  DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME(),
        updated_at  DATETIME2(7) NOT NULL DEFAULT SYSUTCDATETIME()
    );
    CREATE UNIQUE INDEX uq_permissions_name ON dbo.permissions (name);
END

IF OBJECT_ID(N'dbo.role_permissions', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.role_permissions (
        role_id       INT NOT NULL,
        permission_id INT NOT NULL,
        CONSTRAINT PK_role_permissions PRIMARY KEY (role_id, permission_id),
        CONSTRAINT FK_rp_role FOREIGN KEY (role_id) REFERENCES dbo.roles (id) ON DELETE CASCADE,
        CONSTRAINT FK_rp_perm FOREIGN KEY (permission_id) REFERENCES dbo.permissions (id) ON DELETE CASCADE
    );
END

IF OBJECT_ID(N'dbo.zatrano_user_roles', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.zatrano_user_roles (
        user_id INT NOT NULL,
        role_id INT NOT NULL,
        CONSTRAINT PK_zatrano_user_roles PRIMARY KEY (user_id, role_id),
        CONSTRAINT FK_zur_role FOREIGN KEY (role_id) REFERENCES dbo.roles (id) ON DELETE CASCADE
    );
END

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_role_permissions_role' AND object_id = OBJECT_ID(N'dbo.role_permissions'))
    CREATE INDEX idx_role_permissions_role ON dbo.role_permissions (role_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_role_permissions_perm' AND object_id = OBJECT_ID(N'dbo.role_permissions'))
    CREATE INDEX idx_role_permissions_perm ON dbo.role_permissions (permission_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_user_roles_user' AND object_id = OBJECT_ID(N'dbo.zatrano_user_roles'))
    CREATE INDEX idx_user_roles_user ON dbo.zatrano_user_roles (user_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = N'idx_user_roles_role' AND object_id = OBJECT_ID(N'dbo.zatrano_user_roles'))
    CREATE INDEX idx_user_roles_role ON dbo.zatrano_user_roles (role_id);
