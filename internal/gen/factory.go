package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Factory generates a factory stub for test data generation.
func Factory(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid factory name %q (use letters, digits, _ or -)", rawName)
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
	if err := os.WriteFile(path, []byte(tmplFactory(pascal)), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplFactory(pascal string) string {
	return fmt.Sprintf(`package factory

import (
	"github.com/brianvoe/gofakeit/v6"
)

// %sFactory generates test data for your resource using gofakeit.
type %sFactory struct {
	faker *gofakeit.Faker
}

// New%sFactory constructs a new factory with a seeded faker for consistent test data.
func New%sFactory() *%sFactory {
	return &%sFactory{
		faker: gofakeit.New(0), // seed 0 for consistent results
	}
}

// Make returns a placeholder value for the generated resource.
// Override this method to return concrete model instances with fake data.
func (f *%sFactory) Make() any {
	// TODO: return a concrete model instance with fake data.
	// Example:
	// return &models.%s{
	//     Name: f.faker.Name(),
	//     Email: f.faker.Email(),
	//     CreatedAt: f.faker.Date(),
	// }
	return nil
}

// MakeMany generates multiple instances.
func (f *%sFactory) MakeMany(count int) []any {
	results := make([]any, count)
	for i := 0; i < count; i++ {
		results[i] = f.Make()
	}
	return results
}
`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}
