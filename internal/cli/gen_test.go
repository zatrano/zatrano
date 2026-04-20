package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genTestCmd = &cobra.Command{
	Use:   "test [name]",
	Short: "Generate handler and service test stubs",
	Long: `Creates a test scaffold under tests/ with placeholder handler and service tests.

Example:
  zatrano gen test user`,
	Args: cobra.ExactArgs(1),
	RunE: runGenTest,
}

func init() {
	genTestCmd.Flags().String("out", "tests", "base directory for generated test files relative to module-root")
	genTestCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genTestCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genTestCmd)
}

func runGenTest(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Test(moduleRoot, out, args[0], dry)
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
