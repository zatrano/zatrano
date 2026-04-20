package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/zatrano/zatrano/pkg/config"
)

var pgSchemaName = regexp.MustCompile(`^[a-z0-9][a-z0-9_]{0,62}$`)

func validatePGSchemaName(name string) error {
	s := strings.ToLower(strings.TrimSpace(name))
	if !pgSchemaName.MatchString(s) {
		return fmt.Errorf("invalid PostgreSQL schema name %q (use [a-z0-9][a-z0-9_]*)", name)
	}
	return nil
}

func requirePostgresTenant(cfg *config.Config, cmd string) error {
	if cfg.NormalizedDatabaseDriver() != "postgres" {
		return fmt.Errorf("%s requires database_driver postgres (current %q)", cmd, cfg.NormalizedDatabaseDriver())
	}
	return nil
}

// PostgresWithSearchPath returns a copy of databaseURL with libpq options setting search_path for migrations and pooled connections.
func PostgresWithSearchPath(databaseURL, schema string) (string, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return "", fmt.Errorf("database URL is empty")
	}
	if err := validatePGSchemaName(schema); err != nil {
		return "", err
	}
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("parse database url: %w", err)
	}
	q := u.Query()
	q.Set("options", "-csearch_path="+schema+",public")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// MigrateUpWithSchema runs golang-migrate against cfg with search_path restricted to schema (and public for extensions).
// When migrationsFromCLI is true, file-based migrations are used from migrationsDir; otherwise behaviour follows cfg.MigrationsSource (embed vs file).
func MigrateUpWithSchema(cfg *config.Config, migrationsDir, schema string, steps int, migrationsFromCLI bool) (version uint, dirty bool, err error) {
	if err := requirePostgresTenant(cfg, "db tenants migrate"); err != nil {
		return 0, false, err
	}
	dsn, err := PostgresWithSearchPath(cfg.DatabaseURL, schema)
	if err != nil {
		return 0, false, err
	}
	return MigrateUp(cfg, MigrateRequest{
		Dir:          migrationsDir,
		FileSource:   migrationsFromCLI,
		Steps:        steps,
		DatabaseURL:  dsn,
	})
}

// MigrateDownWithSchema rolls back migrations in the tenant schema search_path.
func MigrateDownWithSchema(cfg *config.Config, migrationsDir, schema string, steps int, migrationsFromCLI bool) (version uint, dirty bool, err error) {
	if err := requirePostgresTenant(cfg, "db tenants rollback"); err != nil {
		return 0, false, err
	}
	dsn, err := PostgresWithSearchPath(cfg.DatabaseURL, schema)
	if err != nil {
		return 0, false, err
	}
	return MigrateDown(cfg, MigrateRequest{
		Dir:          migrationsDir,
		FileSource:   migrationsFromCLI,
		Steps:        steps,
		DatabaseURL:  dsn,
	})
}

// CreateTenantSchema runs CREATE SCHEMA IF NOT EXISTS for PostgreSQL.
func CreateTenantSchema(cfg *config.Config, schema string) error {
	if err := requirePostgresTenant(cfg, "db tenants create-schema"); err != nil {
		return err
	}
	if err := validatePGSchemaName(schema); err != nil {
		return err
	}
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()
	q := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %q`, schema)
	if _, err := db.Exec(q); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}
