package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zatrano/zatrano/pkg/api"
	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/core"
)

var apiKeyCmd = &cobra.Command{
	Use:   "api-key",
	Short: "Manage API keys for external client authentication",
}

var apiKeyCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new API key",
	Long: `Creates a new API key with optional scopes and expiration.

Examples:
  zatrano api-key create "My App" --scopes read,write --expires 2025-12-31
  zatrano api-key create "Read Only" --scopes read`,
	RunE: runAPIKeyCreate,
}

var apiKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active API keys",
	RunE:  runAPIKeyList,
}

var apiKeyRevokeCmd = &cobra.Command{
	Use:   "revoke [id]",
	Short: "Revoke an API key by ID",
	Long:  `Immediately expires an API key, preventing further use.`,
	RunE:  runAPIKeyRevoke,
}

func init() {
	apiKeyCreateCmd.Flags().StringSlice("scopes", []string{}, "comma-separated list of scopes")
	apiKeyCreateCmd.Flags().String("expires", "", "expiration date (YYYY-MM-DD)")

	apiKeyCmd.AddCommand(apiKeyCreateCmd, apiKeyListCmd, apiKeyRevokeCmd)
	rootCmd.AddCommand(apiKeyCmd)
}

func runAPIKeyCreate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires exactly one argument: name")
	}
	name := args[0]

	scopes, _ := cmd.Flags().GetStringSlice("scopes")
	expiresStr, _ := cmd.Flags().GetString("expires")

	var expiresAt *time.Time
	if expiresStr != "" {
		t, err := time.Parse("2006-01-02", expiresStr)
		if err != nil {
			return fmt.Errorf("invalid expires date format, use YYYY-MM-DD: %w", err)
		}
		expiresAt = &t
	}

	cfg, err := config.Load(config.LoadOptions{})
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	app, err := core.Bootstrap(cfg)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}
	defer app.Close()

	if app.DB == nil {
		return fmt.Errorf("database required for API keys")
	}

	manager := api.NewKeyManager(app.DB)
	key, plain, err := manager.Create(name, scopes, expiresAt)
	if err != nil {
		return fmt.Errorf("create api key: %w", err)
	}

	fmt.Printf("API Key created:\n")
	fmt.Printf("  ID:      %d\n", key.ID)
	fmt.Printf("  Name:    %s\n", key.Name)
	fmt.Printf("  Prefix:  %s\n", key.Prefix)
	fmt.Printf("  Scopes:  %s\n", strings.Join(key.Scopes, ", "))
	if key.ExpiresAt != nil {
		fmt.Printf("  Expires: %s\n", key.ExpiresAt.Format("2006-01-02"))
	}
	fmt.Printf("  Key:     %s.%s\n", key.Prefix, plain)
	fmt.Printf("\n⚠️  Store this key securely — it cannot be recovered!\n")
	return nil
}

func runAPIKeyList(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(config.LoadOptions{})
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	app, err := core.Bootstrap(cfg)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}
	defer app.Close()

	if app.DB == nil {
		return fmt.Errorf("database required for API keys")
	}

	manager := api.NewKeyManager(app.DB)
	keys, err := manager.List()
	if err != nil {
		return fmt.Errorf("list api keys: %w", err)
	}

	if len(keys) == 0 {
		fmt.Println("No active API keys found.")
		return nil
	}

	fmt.Printf("%-4s %-20s %-10s %-20s %-15s\n", "ID", "Name", "Prefix", "Scopes", "Expires")
	fmt.Println(strings.Repeat("-", 70))
	for _, key := range keys {
		expires := ""
		if key.ExpiresAt != nil {
			expires = key.ExpiresAt.Format("2006-01-02")
		}
		fmt.Printf("%-4d %-20s %-10s %-20s %-15s\n",
			key.ID, key.Name, key.Prefix, strings.Join(key.Scopes, ","), expires)
	}
	return nil
}

func runAPIKeyRevoke(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires exactly one argument: id")
	}

	id, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	cfg, err := config.Load(config.LoadOptions{})
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	app, err := core.Bootstrap(cfg)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}
	defer app.Close()

	if app.DB == nil {
		return fmt.Errorf("database required for API keys")
	}

	manager := api.NewKeyManager(app.DB)
	if err := manager.Revoke(uint(id)); err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}

	fmt.Printf("API key %d revoked.\n", id)
	return nil
}
