package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Seeder generates a SQL seed file under db/seeds/.
func Seeder(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid seeder name %q (use letters, digits, _ or -)", rawName)
	}
	table := strings.ReplaceAll(name, "_", "-") + "s"
	outDir := filepath.Join(moduleRoot, baseDir)
	path := filepath.Join(outDir, name+"_seed.sql")
	if dryRun {
		return []string{path}, nil
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(tmplSeeder(table)), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplSeeder(table string) string {
	return fmt.Sprintf(`-- SQL seed for %s
-- Add initial data for your application.

INSERT INTO %s (name, email)
VALUES ('Example Name', 'example@example.com');
`, table, table)
}
