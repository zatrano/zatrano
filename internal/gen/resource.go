package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Resource generates an API resource transformer stub.
func Resource(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid resource name %q (use letters, digits, _ or -)", rawName)
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
	if err := os.WriteFile(path, []byte(tmplResource(pascal)), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplResource(pascal string) string {
	return fmt.Sprintf(`package resources

// %sResource formats JSON responses for your resource.
type %sResource struct {
	ID   uint   `+"`json:\"id\"`"+`
	Name string `+"`json:\"name,omitempty\"`"+`
	// TODO: add additional fields and hide sensitive values here.
}

// Build%sResource maps a model to a JSON resource.
func Build%sResource(model any) *%sResource {
	// TODO: replace with concrete model mapping logic.
	return &%sResource{}
}
`, pascal, pascal, pascal, pascal, pascal, pascal)
}
