package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genSeederCmd = &cobra.Command{
	Use:   "seeder [name]",
	Short: "Generate a SQL seeder file",
	Long: `Creates a SQL seed file under db/seeds/.

Example:
  zatrano gen seeder users`,
	Args: cobra.ExactArgs(1),
	RunE: runGenSeeder,
}

func init() {
	genSeederCmd.Flags().String("out", "db/seeds", "base directory for generated seed files relative to module-root")
	genSeederCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genSeederCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genSeederCmd)
}

func runGenSeeder(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Seeder(moduleRoot, out, args[0], dry)
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
