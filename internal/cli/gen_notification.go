package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/gen"
)

var genNotificationCmd = &cobra.Command{
	Use:   "notification [name]",
	Short: "Generate a notification stub",
	Long: `Creates a notification scaffold under modules/notifications/.

Example:
  zatrano gen notification welcome`,
	Args: cobra.ExactArgs(1),
	RunE: runGenNotification,
}

func init() {
	genNotificationCmd.Flags().String("out", "modules/notifications", "base directory for generated notification file relative to module-root")
	genNotificationCmd.Flags().String("module-root", ".", "directory containing go.mod")
	genNotificationCmd.Flags().Bool("dry-run", false, "print paths only, do not write files")
	genCmd.AddCommand(genNotificationCmd)
}

func runGenNotification(cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetString("out")
	moduleRoot, _ := cmd.Flags().GetString("module-root")
	dry, _ := cmd.Flags().GetBool("dry-run")

	paths, err := gen.Notification(moduleRoot, out, args[0], dry)
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
