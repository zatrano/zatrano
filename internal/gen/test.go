package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Test generates handler and service test stubs.
func Test(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid test name %q (use letters, digits, _ or -)", rawName)
	}
	snake := name
	outDir := filepath.Join(moduleRoot, baseDir)
	path := filepath.Join(outDir, snake+"_test.go")
	if dryRun {
		return []string{path}, nil
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(tmplTest(snake, snakeToPascal(name))), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplTest(snake, pascal string) string {
	return fmt.Sprintf(`package tests

import (
	"testing"

	"github.com/zatrano/zatrano/pkg/testing"
)

func Test%sHandler(t *testing.T) {
	// TODO: Initialize your Fiber app and HTTP client
	// client := testing.NewHTTPClient(app)
	// client.WithToken("your-jwt-token").Get("/api/v1/%s").AssertStatus(200)

	t.Run("stub", func(t *testing.T) {
		// TODO: build a Fiber context and assert handler behavior.
	})
}

func Test%sService(t *testing.T) {
	// TODO: Initialize test suite for database rollback
	// suite := testing.NewTestSuite(db)
	// suite.SetupTest(t)
	// defer suite.TeardownTest()

	t.Run("stub", func(t *testing.T) {
		// TODO: construct the service and verify business logic.
	})
}
`, pascal, snake, pascal)
}
