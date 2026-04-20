package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genFactoryCmd = &cobra.Command{
	Use:   "factory [name]",
	Short: "Generate a factory stub for test data generation",
	Long: `Creates a factory scaffold under pkg/factory/.

Example:
  zatrano gen factory user`,
	Args: cobra.ExactArgs(1),
	RunE: runGenFactory,
}

func init() {
	genFactoryCmd.Flags().String("out", "pkg/factory", "base directory for generated factory file relative to module-root")
	genFactoryCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genFactoryCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genFactoryCmd)
}

func runGenFactory(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Factory(moduleRoot, out, args[0], dry)
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
