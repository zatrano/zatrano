package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var pgSchemaName = regexp.MustCompile(`^[a-z0-9][a-z0-9_]{0,62}$`)

func validatePGSchemaName(name string) error {
	s := strings.ToLower(strings.TrimSpace(name))
	if !pgSchemaName.MatchString(s) {
		return fmt.Errorf("invalid PostgreSQL schema name %q (use [a-z0-9][a-z0-9_]*)", name)
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

// MigrateUpWithSchema runs golang-migrate against databaseURL with search_path restricted to schema (and public for extensions).
func MigrateUpWithSchema(databaseURL, migrationsDir, schema string, steps int) (version uint, dirty bool, err error) {
	dsn, err := PostgresWithSearchPath(databaseURL, schema)
	if err != nil {
		return 0, false, err
	}
	return MigrateUp(dsn, migrationsDir, steps)
}

// MigrateDownWithSchema rolls back migrations in the tenant schema search_path.
func MigrateDownWithSchema(databaseURL, migrationsDir, schema string, steps int) (version uint, dirty bool, err error) {
	dsn, err := PostgresWithSearchPath(databaseURL, schema)
	if err != nil {
		return 0, false, err
	}
	return MigrateDown(dsn, migrationsDir, steps)
}

// CreateTenantSchema runs CREATE SCHEMA IF NOT EXISTS for PostgreSQL.
func CreateTenantSchema(databaseURL, schema string) error {
	if err := validatePGSchemaName(schema); err != nil {
		return err
	}
	db, err := sql.Open("pgx", databaseURL)
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
