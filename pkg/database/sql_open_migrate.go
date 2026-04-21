package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zatrano/zatrano/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

// OpenStdSQLForMigrate opens a *sql.DB for golang-migrate (WithInstance), matching cfg.DatabaseDriver.
func OpenStdSQLForMigrate(cfg *config.Config) (*sql.DB, error) {
	u := strings.TrimSpace(cfg.DatabaseURL)
	if u == "" {
		return nil, fmt.Errorf("database_url is empty")
	}
	switch cfg.NormalizedDatabaseDriver() {
	case "postgres":
		return sql.Open("pgx", u)
	case "mysql":
		dsn := mysqlDSNForGORM(u)
		if !strings.Contains(strings.ToLower(dsn), "multistatements") {
			if strings.Contains(dsn, "?") {
				dsn += "&multiStatements=true"
			} else {
				dsn += "?multiStatements=true"
			}
		}
		return sql.Open("mysql", dsn)
	case "sqlite":
		return sql.Open("sqlite3", sqlitePathForGORM(u))
	default:
		return nil, fmt.Errorf("unsupported database_driver %q", cfg.DatabaseDriver)
	}
}
