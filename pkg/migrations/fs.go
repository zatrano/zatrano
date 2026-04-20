package migrations

import "embed"

// SQL holds versioned golang-migrate files per database driver (subfolders of sql/).
//
//go:embed sql/postgres/*.sql sql/mysql/*.sql sql/sqlite/*.sql sql/sqlserver/*.sql
var SQL embed.FS
