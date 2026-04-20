package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Manage storage symlinks and file operations",
}

var storageLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create symlink from storage/app/public to public/storage",
	Long: `Creates a symbolic link from storage/app/public to public/storage.

This allows the framework to serve uploaded files through the web server.

Examples:
  zatrano storage:link
  zatrano storage:link --force   # Re-create if exists`,
	RunE: runStorageLink,
}

var storageClearCmd = &cobra.Command{
	Use:   "clear [disk]",
	Short: "Clear all files from a storage disk",
	Long: `Removes all files from the specified storage disk.

Use with caution! This will permanently delete all files.

Examples:
  zatrano storage:clear     # Clear default disk
  zatrano storage:clear temp   # Clear 'temp' disk`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStorageClear,
}

func init() {
	storageLinkCmd.Flags().Bool("force", false, "force creation if symlink already exists")
	storageLinkCmd.Flags().String("storage-path", "storage/app/public", "path to storage directory")
	storageLinkCmd.Flags().String("public-path", "public/storage", "path to public directory")

	storageClearCmd.Flags().Bool("force", false, "skip confirmation prompt")

	storageCmd.AddCommand(storageLinkCmd, storageClearCmd)
	rootCmd.AddCommand(storageCmd)
}

func runStorageLink(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	storagePath, _ := cmd.Flags().GetString("storage-path")
	publicPath, _ := cmd.Flags().GetString("public-path")

	// Resolve paths
	absStoragePath, err := filepath.Abs(storagePath)
	if err != nil {
		return fmt.Errorf("resolve storage path: %w", err)
	}

	absPublicPath, err := filepath.Abs(publicPath)
	if err != nil {
		return fmt.Errorf("resolve public path: %w", err)
	}

	// Check if storage path exists
	if _, err := os.Stat(absStoragePath); err != nil {
		if os.IsNotExist(err) {
			// Create the storage directory
			if err := os.MkdirAll(absStoragePath, 0o755); err != nil {
				return fmt.Errorf("create storage directory: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Created: %s\n", storagePath)
		} else {
			return fmt.Errorf("stat storage path: %w", err)
		}
	}

	// Create public directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(absPublicPath), 0o755); err != nil {
		return fmt.Errorf("create public directory: %w", err)
	}

	// Check if symlink already exists
	if _, err := os.Lstat(absPublicPath); err == nil {
		if !force {
			return fmt.Errorf("symlink already exists at %s (use --force to overwrite)", publicPath)
		}
		// Remove existing symlink
		if err := os.Remove(absPublicPath); err != nil {
			return fmt.Errorf("remove existing symlink: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Removed: %s\n", publicPath)
	}

	// Create symlink
	if err := os.Symlink(absStoragePath, absPublicPath); err != nil {
		return fmt.Errorf("create symlink: %w", err)
	}

	fmt.Printf("Symlink created:\n")
	fmt.Printf("  %s → %s\n", publicPath, storagePath)
	return nil
}

func runStorageClear(cmd *cobra.Command, args []string) error {
	disk := "default"
	if len(args) > 0 {
		disk = args[0]
	}

	force, _ := cmd.Flags().GetBool("force")

	if !force {
		fmt.Printf("⚠️  This will permanently delete all files from the '%s' disk.\n", disk)
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// TODO: Implement actual clearing with storage manager
	// For now, just show what would be cleared
	fmt.Printf("Would clear all files from '%s' disk\n", disk)
	fmt.Println("(Requires storage manager integration)")

	return nil
}
