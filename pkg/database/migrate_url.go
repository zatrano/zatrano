package database

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zatrano/zatrano/pkg/config"
)

// MigrateURL returns a connection string understood by golang-migrate for cfg.DatabaseDriver.
// See: https://github.com/golang-migrate/migrate#databases
func MigrateURL(cfg *config.Config) (string, error) {
	u := strings.TrimSpace(cfg.DatabaseURL)
	if u == "" {
		return "", fmt.Errorf("database_url is empty")
	}
	switch cfg.NormalizedDatabaseDriver() {
	case "postgres":
		if !strings.HasPrefix(strings.ToLower(u), "postgres://") && !strings.HasPrefix(strings.ToLower(u), "postgresql://") {
			return "", fmt.Errorf("postgres database_url must start with postgres:// or postgresql:// (for golang-migrate)")
		}
		return u, nil
	case "mysql":
		if strings.HasPrefix(strings.ToLower(u), "mysql://") {
			return u, nil
		}
		return "mysql://" + u, nil
	case "sqlserver":
		if !strings.HasPrefix(strings.ToLower(u), "sqlserver://") {
			return "", fmt.Errorf("sqlserver database_url must start with sqlserver://")
		}
		return u, nil
	case "sqlite":
		lo := strings.ToLower(u)
		if strings.HasPrefix(lo, "sqlite3://") {
			return u, nil
		}
		if lo == ":memory:" || lo == "memory" {
			return "sqlite3://:memory:", nil
		}
		abs, err := filepath.Abs(u)
		if err != nil {
			return "", fmt.Errorf("sqlite path: %w", err)
		}
		abs = filepath.ToSlash(abs)
		if vol := filepath.VolumeName(abs); vol != "" && strings.HasPrefix(abs, vol) {
			// Windows: C:/x -> golang-migrate accepts file path in URL
		}
		return "sqlite3://" + abs, nil
	default:
		return "", fmt.Errorf("unsupported database_driver %q", cfg.DatabaseDriver)
	}
}
