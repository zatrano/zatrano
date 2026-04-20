package migrations

import "strings"

// SQLSubdir returns the embedded folder name for cfg.DatabaseDriver (postgres, mysql, sqlite, sqlserver).
func SQLSubdir(driver string) string {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "mysql":
		return "mysql"
	case "sqlite":
		return "sqlite"
	case "sqlserver":
		return "sqlserver"
	default:
		return "postgres"
	}
}
