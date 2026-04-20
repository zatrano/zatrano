package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genModelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Generate a standalone model and migration files",
	Long: `Creates a model scaffold under pkg/repository/models/ and PostgreSQL migration stubs under pkg/migrations/sql/postgres/ (copy or adapt for mysql/sqlite/sqlserver when needed).

Example:
  zatrano gen model user`,
	Args: cobra.ExactArgs(1),
	RunE: runGenModel,
}

func init() {
	genModelCmd.Flags().String("out", "pkg/repository/models", "base directory for generated model file relative to module-root")
	genModelCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genModelCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genModelCmd)
}

func runGenModel(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Model(moduleRoot, out, args[0], dry)
	if err != nil {
		return err
	}
	if dry {
		fmt.Println("dry-run — would write:")
	} else {
		fmt.Println("written:")
	}
	fmt.Println(strings.Join(paths, "\n"))
	return nil
}
