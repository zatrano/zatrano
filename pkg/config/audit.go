package config

import (
	"fmt"
	"strings"
)

// Audit configures model activity logging and HTTP request auditing.
type Audit struct {
	// Enabled is the master switch for audit features (middleware + helpers).
	Enabled bool `mapstructure:"enabled"`
	// ModelEnabled registers GORM callbacks for registered subjects (requires migrations).
	ModelEnabled bool `mapstructure:"model_enabled"`
	// HttpEnabled turns on middleware.AuditLog when Enabled is true.
	HttpEnabled bool `mapstructure:"http_enabled"`
	// HttpDriver is where HTTP audit lines go: db | file.
	HttpDriver string `mapstructure:"http_driver"`
	// HttpFilePath is the append-only JSONL path when HttpDriver is file.
	HttpFilePath string `mapstructure:"http_file_path"`
}

func (c *Config) applyAuditDefaults() {
	a := &c.Audit
	if strings.TrimSpace(a.HttpDriver) == "" {
		a.HttpDriver = "db"
	}
}

func (c *Config) validateAudit() error {
	if !c.Audit.Enabled {
		return nil
	}
	if c.Audit.HttpEnabled {
		switch strings.ToLower(strings.TrimSpace(c.Audit.HttpDriver)) {
		case "db", "file":
		default:
			return fmt.Errorf("audit.http_driver must be db or file (got %q)", c.Audit.HttpDriver)
		}
		if strings.EqualFold(c.Audit.HttpDriver, "file") && strings.TrimSpace(c.Audit.HttpFilePath) == "" {
			return fmt.Errorf("audit.http_file_path is required when audit.http_driver is file")
		}
	}
	if c.Audit.ModelEnabled && strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("audit.model_enabled requires database_url")
	}
	if c.Audit.HttpEnabled && strings.EqualFold(c.Audit.HttpDriver, "db") && strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("audit.http_enabled with driver db requires database_url")
	}
	return nil
}
