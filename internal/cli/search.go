package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zatrano/zatrano/pkg/search"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Full-text helpers and external search index import",
}

var searchImportCmd = &cobra.Command{
	Use:   "import <model>",
	Short: "Bulk-index registered models into Meilisearch or Typesense",
	Long: `Runs the importer registered with search.RegisterImporter for the given model name (case-insensitive).

Requires database_url, search.enabled, and search.driver (meilisearch or typesense). Create matching indexes/collections in the engine before importing.`,
	Args: cobra.ExactArgs(1),
	RunE: runSearchImport,
}

func init() {
	searchImportCmd.Flags().String("env", "", "environment name; default ENV or dev")
	searchImportCmd.Flags().String("config-dir", "config", "directory containing {env}.yaml")
	searchImportCmd.Flags().Bool("no-dotenv", false, "do not load .env")

	searchCmd.AddCommand(searchImportCmd)
	rootCmd.AddCommand(searchCmd)
}

func runSearchImport(cmd *cobra.Command, args []string) error {
	cfg, err := loadCfgFlags(cmd)
	if err != nil {
		return err
	}
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return fmt.Errorf("database_url is required for search import")
	}
	drv, err := search.NewDriverForCLI(cfg)
	if err != nil {
		return err
	}
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	model := strings.TrimSpace(args[0])
	ctx := context.Background()
	if err := search.RunImport(ctx, db, drv, model); err != nil {
		return err
	}
	fmt.Printf("ok: search import finished for %q\n", model)
	return nil
}
