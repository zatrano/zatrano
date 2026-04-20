package config

import (
	"fmt"
	"strings"
)

// Search configures PostgreSQL FTS helpers and optional external engines (Meilisearch / Typesense).
type Search struct {
	// Enabled turns on pkg/search client wiring in bootstrap (external drivers only).
	Enabled bool `mapstructure:"enabled"`
	// Driver selects the external index backend: meilisearch | typesense | "" (postgres-only, scopes only).
	Driver string `mapstructure:"driver"`
	// DefaultIndexPrefix is prepended to logical index names (e.g. "app_" + "products").
	DefaultIndexPrefix string `mapstructure:"default_index_prefix"`
	// MeilisearchURL is the HTTP origin (e.g. http://127.0.0.1:7700).
	MeilisearchURL string `mapstructure:"meilisearch_url"`
	// MeilisearchAPIKey is sent as Bearer (may be empty on local dev).
	MeilisearchAPIKey string `mapstructure:"meilisearch_api_key"`
	// TypesenseURL is the node origin (e.g. http://127.0.0.1:8108).
	TypesenseURL string `mapstructure:"typesense_url"`
	// TypesenseAPIKey is sent as X-TYPESENSE-API-KEY.
	TypesenseAPIKey string `mapstructure:"typesense_api_key"`
	// PostgresFTSLanguage is the default regconfig for plainto_tsquery scopes (e.g. english, simple, turkish).
	PostgresFTSLanguage string `mapstructure:"postgres_fts_language"`
}

func (c *Config) applySearchDefaults() {
	s := &c.Search
	if strings.TrimSpace(s.PostgresFTSLanguage) == "" {
		s.PostgresFTSLanguage = "simple"
	}
	if strings.TrimSpace(s.DefaultIndexPrefix) == "" {
		s.DefaultIndexPrefix = "zatrano_"
	}
}

func (c *Config) validateSearch() error {
	if !c.Search.Enabled {
		return nil
	}
	d := strings.ToLower(strings.TrimSpace(c.Search.Driver))
	switch d {
	case "meilisearch":
		if strings.TrimSpace(c.Search.MeilisearchURL) == "" {
			return fmt.Errorf("search.meilisearch_url is required when search.driver is meilisearch")
		}
	case "typesense":
		if strings.TrimSpace(c.Search.TypesenseURL) == "" {
			return fmt.Errorf("search.typesense_url is required when search.driver is typesense")
		}
		if strings.TrimSpace(c.Search.TypesenseAPIKey) == "" {
			return fmt.Errorf("search.typesense_api_key is required when search.driver is typesense")
		}
	case "":
		return fmt.Errorf("search.driver is required when search.enabled is true (meilisearch or typesense)")
	default:
		return fmt.Errorf("search.driver must be meilisearch or typesense (got %q)", c.Search.Driver)
	}
	return nil
}
