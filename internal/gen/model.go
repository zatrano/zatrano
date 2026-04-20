package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Model generates a standalone model file and matching SQL migration files.
func Model(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid model name %q (use letters, digits, _ or -)", rawName)
	}
	pascal := snakeToPascal(name)
	plural := strings.ReplaceAll(name, "_", "-") + "s"
	modelDir := filepath.Join(moduleRoot, baseDir)
	migDir := filepath.Join(moduleRoot, "migrations")
	stamp := time.Now().UTC().Format("20060102150405")
	modelPath := filepath.Join(modelDir, name+".go")
	upPath := filepath.Join(migDir, stamp+"_"+name+".up.sql")
	downPath := filepath.Join(migDir, stamp+"_"+name+".down.sql")
	paths := []string{modelPath, upPath, downPath}
	if dryRun {
		return paths, nil
	}
	if err := os.MkdirAll(modelDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(migDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(modelPath, []byte(tmplModel(pascal)), 0o644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(upPath, []byte(tmplModelMigrationUp(plural)), 0o644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(downPath, []byte(tmplModelMigrationDown(plural)), 0o644); err != nil {
		return nil, err
	}
	return paths, nil
}

func tmplModel(pascal string) string {
	return fmt.Sprintf(`package models

import "time"

// %s is the model scaffold for your resource.
type %s struct {
	ID        uint      `+"`gorm:\"primaryKey\" json:\"id\"`"+`
	CreatedAt time.Time `+"`json:\"created_at\"`"+`
	UpdatedAt time.Time `+"`json:\"updated_at\"`"+`
	// TODO: add your fields here.
}
`, pascal, pascal)
}

func tmplModelMigrationUp(plural string) string {
	return fmt.Sprintf(`-- Migration for %s

CREATE TABLE IF NOT EXISTS %s (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`, plural, plural)
}

func tmplModelMigrationDown(plural string) string {
	return fmt.Sprintf(`-- Rollback migration for %s

DROP TABLE IF EXISTS %s;
`, plural, plural)
}
