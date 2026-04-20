package config

import (
	"fmt"
	"strings"
)

// Tenant configures multi-tenancy resolution and data isolation hints.
type Tenant struct {
	// Enabled registers ResolveTenant middleware and tenant helpers.
	Enabled bool `mapstructure:"enabled"`
	// Mode is how the tenant key is read: "header" or "subdomain".
	Mode string `mapstructure:"mode"`
	// HeaderName is the HTTP header for Mode=header (default X-Tenant-ID).
	HeaderName string `mapstructure:"header_name"`
	// SubdomainSuffix must match the end of Host for Mode=subdomain (e.g. ".app.test").
	// The tenant key is the label before this suffix (e.g. acme.app.test → acme when suffix is .app.test).
	SubdomainSuffix string `mapstructure:"subdomain_suffix"`
	// Required returns 400 when no tenant could be resolved (when Enabled).
	Required bool `mapstructure:"required"`
	// Isolation selects row-level (tenant_id column) vs PostgreSQL schema per tenant.
	// Row: use repository.NewTenantAware and tenant.FromContext.
	// Schema: set search_path via tenant.GormSession; run migrations with db tenants migrate.
	Isolation string `mapstructure:"isolation"` // row | schema
	// RowColumn is the SQL column for row isolation WHERE (default tenant_id).
	RowColumn string `mapstructure:"row_column"`
	// SchemaPrefix is prepended to sanitized keys for CREATE SCHEMA / search_path (default tenant_).
	SchemaPrefix string `mapstructure:"schema_prefix"`
}

func (c *Config) applyTenantDefaults() {
	t := &c.Tenant
	if strings.TrimSpace(t.Mode) == "" {
		t.Mode = "header"
	}
	if strings.TrimSpace(t.HeaderName) == "" {
		t.HeaderName = "X-Tenant-ID"
	}
	if strings.TrimSpace(t.RowColumn) == "" {
		t.RowColumn = "tenant_id"
	}
	if strings.TrimSpace(t.SchemaPrefix) == "" {
		t.SchemaPrefix = "tenant_"
	}
	if strings.TrimSpace(t.Isolation) == "" {
		t.Isolation = "row"
	}
}

func (c *Config) validateTenant() error {
	if !c.Tenant.Enabled {
		return nil
	}
	switch strings.ToLower(strings.TrimSpace(c.Tenant.Mode)) {
	case "header", "subdomain":
	default:
		return fmt.Errorf("tenant.mode must be header or subdomain (got %q)", c.Tenant.Mode)
	}
	switch strings.ToLower(strings.TrimSpace(c.Tenant.Isolation)) {
	case "row", "schema":
	default:
		return fmt.Errorf("tenant.isolation must be row or schema (got %q)", c.Tenant.Isolation)
	}
	if strings.EqualFold(c.Tenant.Mode, "subdomain") && strings.TrimSpace(c.Tenant.SubdomainSuffix) == "" {
		return fmt.Errorf("tenant.subdomain_suffix is required when tenant.mode is subdomain")
	}
	if strings.EqualFold(strings.TrimSpace(c.Tenant.Isolation), "schema") && c.NormalizedDatabaseDriver() != "postgres" {
		return fmt.Errorf("tenant.isolation schema requires database_driver postgres (got %q)", c.NormalizedDatabaseDriver())
	}
	return nil
}
