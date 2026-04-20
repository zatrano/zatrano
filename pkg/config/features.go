package config

import (
	"fmt"
	"strings"
)

// Features configures optional feature flags (YAML and/or PostgreSQL).
type Features struct {
	// Enabled turns on the registry; when false, all flags evaluate to disabled.
	Enabled bool `mapstructure:"enabled"`
	// Source selects where definitions are read: config | db | both.
	// both: a row in zatrano_feature_flags overrides the static definition for the same key.
	Source string `mapstructure:"source"`
	// Definitions are static flags (used when source is config or both).
	Definitions []FeatureDefinition `mapstructure:"definitions"`
}

// FeatureDefinition is a single flag in YAML (config-only or defaults for "both").
type FeatureDefinition struct {
	Key              string   `mapstructure:"key"`
	Enabled          bool     `mapstructure:"enabled"`
	RolloutPercent   int      `mapstructure:"rollout_percent"`
	AllowedRoles     []string `mapstructure:"allowed_roles"`
}

func (c *Config) applyFeaturesDefaults() {
	f := &c.Features
	if strings.TrimSpace(f.Source) == "" {
		f.Source = "config"
	}
	f.Source = strings.ToLower(strings.TrimSpace(f.Source))
}

func (c *Config) validateFeatures() error {
	if !c.Features.Enabled {
		return nil
	}
	switch c.Features.Source {
	case "config", "db", "both":
	default:
		return fmt.Errorf("features.source must be config, db, or both (got %q)", c.Features.Source)
	}
	if c.Features.Source == "db" || c.Features.Source == "both" {
		if strings.TrimSpace(c.DatabaseURL) == "" {
			return fmt.Errorf("features.source %q requires database_url", c.Features.Source)
		}
	}
	for _, d := range c.Features.Definitions {
		if strings.TrimSpace(d.Key) == "" {
			return fmt.Errorf("features.definitions entry has empty key")
		}
		if d.RolloutPercent < 0 || d.RolloutPercent > 100 {
			return fmt.Errorf("features.definitions rollout_percent for %q must be 0..100", d.Key)
		}
	}
	return nil
}
