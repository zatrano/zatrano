package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genCommandCmd = &cobra.Command{
	Use:   "command [name]",
	Short: "Generate a stub CLI command",
	Long: `Creates a Cobra command scaffold under internal/cli/.

Example:
  zatrano gen command notify`,
	Args: cobra.ExactArgs(1),
	RunE: runGenCommand,
}

func init() {
	genCommandCmd.Flags().String("out", "internal/cli", "base directory for generated CLI command file relative to module-root")
	genCommandCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genCommandCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genCommandCmd)
}

func runGenCommand(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Command(moduleRoot, out, args[0], dry)
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
