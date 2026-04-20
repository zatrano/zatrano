package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Middleware generates a Fiber middleware stub.
func Middleware(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid middleware name %q (use letters, digits, _ or -)", rawName)
	}
	pascal := snakeToPascal(name)
	outDir := filepath.Join(moduleRoot, baseDir)
	path := filepath.Join(outDir, name+".go")
	if dryRun {
		return []string{path}, nil
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(tmplMiddleware(pascal)), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplMiddleware(pascal string) string {
	return fmt.Sprintf(`package middleware

import "github.com/gofiber/fiber/v3"

// %sMiddleware is a starter middleware stub.
func %sMiddleware(c fiber.Ctx) error {
	// TODO: add middleware logic (auth, CORS, headers, rate limiting, etc.)
	return c.Next()
}
`, pascal, pascal)
}
