package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genMiddlewareCmd = &cobra.Command{
	Use:   "middleware [name]",
	Short: "Generate a Fiber middleware stub",
	Long: `Creates a middleware scaffold under pkg/middleware/.

Example:
  zatrano gen middleware auth`,
	Args: cobra.ExactArgs(1),
	RunE: runGenMiddleware,
}

func init() {
	genMiddlewareCmd.Flags().String("out", "pkg/middleware", "base directory for generated middleware file relative to module-root")
	genMiddlewareCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genMiddlewareCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genMiddlewareCmd)
}

func runGenMiddleware(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Middleware(moduleRoot, out, args[0], dry)
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
