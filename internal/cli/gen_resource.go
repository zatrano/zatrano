package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genResourceCmd = &cobra.Command{
	Use:   "resource [name]",
	Short: "Generate an API resource transformer stub",
	Long: `Creates a resource transformer scaffold under pkg/resources/.

Example:
  zatrano gen resource user`,
	Args: cobra.ExactArgs(1),
	RunE: runGenResource,
}

func init() {
	genResourceCmd.Flags().String("out", "pkg/resources", "base directory for generated resource file relative to module-root")
	genResourceCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genResourceCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genResourceCmd)
}

func runGenResource(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Resource(moduleRoot, out, args[0], dry)
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
