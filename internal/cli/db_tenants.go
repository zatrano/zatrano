package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zatrano/zatrano/internal/db"
	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/tenant"
)

var dbTenantsCmd = &cobra.Command{
	Use:   "tenants",
	Short: "Multi-tenant PostgreSQL schema helpers",
	Long:  "Create tenant schemas and run golang-migrate with search_path scoped to a tenant (see README → Multi-tenancy).",
}

var dbTenantsMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply SQL migrations in a tenant schema (search_path)",
	RunE:  runTenantsMigrate,
}

var dbTenantsRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Roll back migrations in a tenant schema",
	RunE:  runTenantsRollback,
}

var dbTenantsCreateSchemaCmd = &cobra.Command{
	Use:   "create-schema",
	Short: "CREATE SCHEMA IF NOT EXISTS for a tenant",
	RunE:  runTenantsCreateSchema,
}

func init() {
	dbTenantsMigrateCmd.Flags().String("tenant", "", "tenant key (required)")
	dbTenantsMigrateCmd.Flags().String("env", "", "environment name; default ENV or dev")
	dbTenantsMigrateCmd.Flags().String("config-dir", "config", "directory containing {env}.yaml")
	dbTenantsMigrateCmd.Flags().Bool("no-dotenv", false, "do not load .env")
	dbTenantsMigrateCmd.Flags().String("migrations", "", "override migrations directory")
	dbTenantsMigrateCmd.Flags().Int("steps", 0, "number of up migrations (0 = all)")

	dbTenantsRollbackCmd.Flags().String("tenant", "", "tenant key (required)")
	dbTenantsRollbackCmd.Flags().String("env", "", "environment name; default ENV or dev")
	dbTenantsRollbackCmd.Flags().String("config-dir", "config", "directory containing {env}.yaml")
	dbTenantsRollbackCmd.Flags().Bool("no-dotenv", false, "do not load .env")
	dbTenantsRollbackCmd.Flags().String("migrations", "", "override migrations directory")
	dbTenantsRollbackCmd.Flags().Int("steps", 1, "number of migrations to roll back")

	dbTenantsCreateSchemaCmd.Flags().String("tenant", "", "tenant key (required)")
	dbTenantsCreateSchemaCmd.Flags().String("env", "", "environment name; default ENV or dev")
	dbTenantsCreateSchemaCmd.Flags().String("config-dir", "config", "directory containing {env}.yaml")
	dbTenantsCreateSchemaCmd.Flags().Bool("no-dotenv", false, "do not load .env")

	dbTenantsCmd.AddCommand(dbTenantsMigrateCmd, dbTenantsRollbackCmd, dbTenantsCreateSchemaCmd)
	dbCmd.AddCommand(dbTenantsCmd)
}

func runTenantsMigrate(cmd *cobra.Command, _ []string) error {
	cfg, err := loadCfgFlags(cmd)
	if err != nil {
		return err
	}
	key, _ := cmd.Flags().GetString("tenant")
	schema, err := tenantSchemaFromConfig(cfg, key)
	if err != nil {
		return err
	}
	dir := cfg.MigrationsDir
	fileSrc := false
	if f, _ := cmd.Flags().GetString("migrations"); strings.TrimSpace(f) != "" {
		fileSrc = true
		dir = f
	}
	steps, _ := cmd.Flags().GetInt("steps")
	ver, dirty, err := db.MigrateUpWithSchema(cfg, dir, schema, steps, fileSrc)
	if err != nil {
		return fmt.Errorf("tenant migrate: %w", err)
	}
	fmt.Printf("ok tenant=%q schema=%q version=%d dirty=%v\n", key, schema, ver, dirty)
	return nil
}

func runTenantsRollback(cmd *cobra.Command, _ []string) error {
	cfg, err := loadCfgFlags(cmd)
	if err != nil {
		return err
	}
	key, _ := cmd.Flags().GetString("tenant")
	schema, err := tenantSchemaFromConfig(cfg, key)
	if err != nil {
		return err
	}
	dir := cfg.MigrationsDir
	fileSrc := false
	if f, _ := cmd.Flags().GetString("migrations"); strings.TrimSpace(f) != "" {
		fileSrc = true
		dir = f
	}
	steps, _ := cmd.Flags().GetInt("steps")
	ver, dirty, err := db.MigrateDownWithSchema(cfg, dir, schema, steps, fileSrc)
	if err != nil {
		return fmt.Errorf("tenant rollback: %w", err)
	}
	fmt.Printf("ok tenant=%q schema=%q version=%d dirty=%v\n", key, schema, ver, dirty)
	return nil
}

func runTenantsCreateSchema(cmd *cobra.Command, _ []string) error {
	cfg, err := loadCfgFlags(cmd)
	if err != nil {
		return err
	}
	key, _ := cmd.Flags().GetString("tenant")
	schema, err := tenantSchemaFromConfig(cfg, key)
	if err != nil {
		return err
	}
	if err := db.CreateTenantSchema(cfg, schema); err != nil {
		return err
	}
	fmt.Printf("ok schema %q created (if not exists)\n", schema)
	return nil
}

func tenantSchemaFromConfig(cfg *config.Config, key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", fmt.Errorf("--tenant is required")
	}
	return tenant.SchemaName(cfg.Tenant.SchemaPrefix, key)
}
