package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Command generates a Cobra CLI command scaffold under internal/cli/.
func Command(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid command name %q (use letters, digits, _ or -)", rawName)
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
	if err := os.WriteFile(path, []byte(tmplCommand(name, pascal)), 0o644); err != nil {
		return nil, err
	}
	return []string{path}, nil
}

func tmplCommand(useName, pascal string) string {
	return fmt.Sprintf(`package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var %sCmd = &cobra.Command{
	Use:   %q,
	Short: "TODO: implement %s command",
	Args:  cobra.ExactArgs(0),
	RunE:  run%s,
}

func init() {
	rootCmd.AddCommand(%sCmd)
}

func run%s(cmd *cobra.Command, args []string) error {
	return fmt.Errorf(%q)
}
`, pascal, useName, useName, pascal, pascal, pascal, "not implemented yet")
}
