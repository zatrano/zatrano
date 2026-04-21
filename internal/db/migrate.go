package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/zatrano/zatrano/pkg/config"
	zdb "github.com/zatrano/zatrano/pkg/database"
	zm "github.com/zatrano/zatrano/pkg/migrations"
)

// MigrateRequest controls where migrations are loaded from.
type MigrateRequest struct {
	// Dir is the filesystem migrations directory when using file source.
	Dir string
	// FileSource is true when the CLI passed an explicit --migrations path.
	FileSource bool
	// Steps for Up: 0 = all; for Down: 0 or negative treated as 1 in legacy path (see MigrateDown).
	Steps int
	// DatabaseURL, if set, overrides cfg.DatabaseURL for opening *sql.DB (PostgreSQL tenant search_path DSN).
	DatabaseURL string
}

func useFileMigrations(cfg *config.Config, req MigrateRequest) bool {
	if req.FileSource {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(cfg.MigrationsSource), "file")
}

// MigrateUp applies SQL migrations using cfg (database_driver + database_url) and request options.
func MigrateUp(cfg *config.Config, req MigrateRequest) (version uint, dirty bool, err error) {
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return 0, false, fmt.Errorf("database URL is empty (set DATABASE_URL)")
	}
	dsn := strings.TrimSpace(req.DatabaseURL)
	if dsn == "" {
		mu, err := zdb.MigrateURL(cfg)
		if err != nil {
			return 0, false, fmt.Errorf("migrate url: %w", err)
		}
		dsn = mu
	}
	if useFileMigrations(cfg, req) {
		dir := strings.TrimSpace(req.Dir)
		if dir == "" {
			dir = cfg.MigrationsDir
		}
		return migrateWithFileURL(dsn, dir, req.Steps)
	}
	return migrateEmbedded(cfg, dsn, req.Steps, false)
}

// MigrateDown rolls back migrations.
func MigrateDown(cfg *config.Config, req MigrateRequest) (version uint, dirty bool, err error) {
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return 0, false, fmt.Errorf("database URL is empty (set DATABASE_URL)")
	}
	dsn := strings.TrimSpace(req.DatabaseURL)
	if dsn == "" {
		mu, err := zdb.MigrateURL(cfg)
		if err != nil {
			return 0, false, fmt.Errorf("migrate url: %w", err)
		}
		dsn = mu
	}
	if useFileMigrations(cfg, req) {
		dir := strings.TrimSpace(req.Dir)
		if dir == "" {
			dir = cfg.MigrationsDir
		}
		return migrateDownWithFileURL(dsn, dir, req.Steps)
	}
	return migrateEmbedded(cfg, dsn, req.Steps, true)
}

func migrateWithFileURL(databaseURL, dir string, steps int) (version uint, dirty bool, err error) {
	src, err := fileSourceURL(dir)
	if err != nil {
		return 0, false, err
	}
	m, err := migrate.New(src, databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("migrate init: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if steps > 0 {
		if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
			v, d, _ := m.Version()
			return v, d, err
		}
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			v, d, _ := m.Version()
			return v, d, err
		}
	}
	return m.Version()
}

func migrateDownWithFileURL(databaseURL, dir string, steps int) (version uint, dirty bool, err error) {
	src, err := fileSourceURL(dir)
	if err != nil {
		return 0, false, err
	}
	m, err := migrate.New(src, databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("migrate init: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if steps <= 0 {
		steps = 1
	}
	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		v, d, _ := m.Version()
		return v, d, err
	}
	return m.Version()
}

func migrateEmbedded(cfg *config.Config, databaseURL string, steps int, down bool) (version uint, dirty bool, err error) {
	sub := "sql/" + zm.SQLSubdir(cfg.NormalizedDatabaseDriver())
	src, err := iofs.New(zm.SQL, sub)
	if err != nil {
		return 0, false, fmt.Errorf("embedded migrations (%s): %w", sub, err)
	}

	db, err := openSQLDBForMigrate(cfg, databaseURL)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = db.Close() }()

	dbDriver, err := databaseDriverForMigrate(cfg, db)
	if err != nil {
		return 0, false, err
	}

	m, err := migrate.NewWithInstance("iofs", src, cfg.NormalizedDatabaseDriver(), dbDriver)
	if err != nil {
		return 0, false, fmt.Errorf("migrate init (embedded): %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if down {
		if steps <= 0 {
			steps = 1
		}
		if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
			v, d, _ := m.Version()
			return v, d, err
		}
	} else {
		if steps > 0 {
			if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
				v, d, _ := m.Version()
				return v, d, err
			}
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				v, d, _ := m.Version()
				return v, d, err
			}
		}
	}
	return m.Version()
}

func openSQLDBForMigrate(cfg *config.Config, databaseURL string) (*sql.DB, error) {
	u := strings.TrimSpace(databaseURL)
	if u != "" && u != strings.TrimSpace(cfg.DatabaseURL) {
		if cfg.NormalizedDatabaseDriver() != "postgres" {
			return nil, fmt.Errorf("custom migrate DSN is only supported for database_driver postgres")
		}
		return sql.Open("pgx", u)
	}
	return zdb.OpenStdSQLForMigrate(cfg)
}

func databaseDriverForMigrate(cfg *config.Config, db *sql.DB) (migratedb.Driver, error) {
	switch cfg.NormalizedDatabaseDriver() {
	case "postgres":
		return postgres.WithInstance(db, &postgres.Config{})
	case "mysql":
		return mysql.WithInstance(db, &mysql.Config{})
	case "sqlite":
		return sqlite3.WithInstance(db, &sqlite3.Config{})
	default:
		return nil, fmt.Errorf("unsupported database_driver %q", cfg.DatabaseDriver)
	}
}

func fileSourceURL(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	abs = filepath.ToSlash(filepath.Clean(abs))
	if vol := filepath.VolumeName(abs); vol != "" {
		abs = "/" + abs
	}
	return "file://" + abs, nil
}
