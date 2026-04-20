package storage

import (
	"log"
)

// Register initializes and registers storage drivers based on configuration.
func Register(cfg *Config) *Manager {
	manager := NewManager()

	switch cfg.Driver {
	case "s3", "minio", "r2":
		// Register S3-compatible driver
		driver, err := NewS3Driver(
			cfg.S3.Region,
			cfg.S3.Bucket,
			cfg.S3.Endpoint,
			cfg.S3.AccessKey,
			cfg.S3.SecretKey,
			cfg.S3.PublicURL,
			cfg.S3.UsePathStyle,
		)
		if err != nil {
			log.Fatalf("Failed to initialize S3 storage driver: %v", err)
		}
		manager.Register("s3", driver)
		manager.SetDefault("s3")

	case "local":
		fallthrough
	default:
		// Register local driver
		publicDriver := NewLocalDriver(
			cfg.Local.Root,
			cfg.Local.PublicURL,
			true,
		)
		manager.Register("local", publicDriver)

		// Register private driver if enabled
		if cfg.Private.Enabled {
			privateDriver := NewLocalDriver(
				cfg.Private.Root,
				"",
				false,
			)
			manager.Register("private", privateDriver)
		}

		manager.SetDefault("local")
	}

	return manager
}

// RegisterFromMap creates a manager from a configuration map.
// Useful for Dependency Injection or factory patterns.
func RegisterFromMap(configMap map[string]any) (*Manager, error) {
	cfg := &Config{}

	// Parse driver type
	if driver, ok := configMap["driver"].(string); ok {
		cfg.Driver = driver
	}

	// Parse local config
	if local, ok := configMap["local"].(map[string]any); ok {
		if root, ok := local["root"].(string); ok {
			cfg.Local.Root = root
		}
		if publicRoot, ok := local["public_root"].(string); ok {
			cfg.Local.PublicRoot = publicRoot
		}
		if publicURL, ok := local["public_url"].(string); ok {
			cfg.Local.PublicURL = publicURL
		}
	}

	// Parse S3 config
	if s3, ok := configMap["s3"].(map[string]any); ok {
		if region, ok := s3["region"].(string); ok {
			cfg.S3.Region = region
		}
		if bucket, ok := s3["bucket"].(string); ok {
			cfg.S3.Bucket = bucket
		}
		if endpoint, ok := s3["endpoint"].(string); ok {
			cfg.S3.Endpoint = endpoint
		}
		if accessKey, ok := s3["access_key"].(string); ok {
			cfg.S3.AccessKey = accessKey
		}
		if secretKey, ok := s3["secret_key"].(string); ok {
			cfg.S3.SecretKey = secretKey
		}
		if publicURL, ok := s3["public_url"].(string); ok {
			cfg.S3.PublicURL = publicURL
		}
		if usePathStyle, ok := s3["use_path_style"].(bool); ok {
			cfg.S3.UsePathStyle = usePathStyle
		}
	}

	// Parse private config
	if private, ok := configMap["private"].(map[string]any); ok {
		if enabled, ok := private["enabled"].(bool); ok {
			cfg.Private.Enabled = enabled
		}
		if root, ok := private["root"].(string); ok {
			cfg.Private.Root = root
		}
	}

	// Set defaults
	if cfg.Local.Root == "" {
		cfg.Local.Root = "storage/app"
	}
	if cfg.Local.PublicRoot == "" {
		cfg.Local.PublicRoot = "public/storage"
	}
	if cfg.Local.PublicURL == "" {
		cfg.Local.PublicURL = "/storage"
	}
	if cfg.Private.Root == "" {
		cfg.Private.Root = "storage/private"
	}

	return Register(cfg), nil
}

// Helper to create a manager with sensible defaults for testing/development.
func DefaultManager() *Manager {
	cfg := &Config{
		Driver: "local",
	}
	cfg.Local.Root = "storage/app"
	cfg.Local.PublicRoot = "public/storage"
	cfg.Local.PublicURL = "/storage"

	return Register(cfg)
}
