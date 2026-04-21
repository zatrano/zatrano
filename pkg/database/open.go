package database

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zatrano/zatrano/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenGORM opens a *gorm.DB using cfg.DatabaseDriver and cfg.DatabaseURL.
// Returns (nil, nil) when database_url is empty.
func OpenGORM(cfg *config.Config, gormLog logger.Interface) (*gorm.DB, error) {
	u := strings.TrimSpace(cfg.DatabaseURL)
	if u == "" {
		return nil, nil
	}
	gc := &gorm.Config{Logger: gormLog}
	switch cfg.NormalizedDatabaseDriver() {
	case "postgres":
		db, err := gorm.Open(postgres.Open(u), gc)
		if err != nil {
			return nil, fmt.Errorf("postgres: %w", err)
		}
		return db, nil
	case "mysql":
		dsn := mysqlDSNForGORM(u)
		db, err := gorm.Open(mysql.Open(dsn), gc)
		if err != nil {
			return nil, fmt.Errorf("mysql: %w", err)
		}
		return db, nil
	case "sqlite":
		path := sqlitePathForGORM(u)
		db, err := gorm.Open(sqlite.Open(path), gc)
		if err != nil {
			return nil, fmt.Errorf("sqlite: %w", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported database_driver %q", cfg.DatabaseDriver)
	}
}

func mysqlDSNForGORM(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(strings.ToLower(s), "mysql://") {
		return strings.TrimPrefix(s, "mysql://")
	}
	return s
}

func sqlitePathForGORM(raw string) string {
	s := strings.TrimSpace(raw)
	lo := strings.ToLower(s)
	switch {
	case strings.HasPrefix(lo, "sqlite://"):
		return strings.TrimPrefix(s, "sqlite://")
	case strings.HasPrefix(lo, "file:"):
		if u, err := url.Parse(s); err == nil && u.Path != "" {
			return u.Path
		}
		return strings.TrimPrefix(s, "file:")
	default:
		return s
	}
}
