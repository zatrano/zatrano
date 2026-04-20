package config

import (
	"fmt"
	"strings"
)

// Admin configures optional HTML admin area (dashboard, metrics, logs) and gen admin route prefix.
type Admin struct {
	// Enabled mounts /admin/* routes when true (requires View for HTML pages).
	Enabled bool `mapstructure:"enabled"`
	// PathPrefix is the URL prefix (default /admin).
	PathPrefix string `mapstructure:"path_prefix"`
	// Secret, if non-empty, requires matching X-Admin-Key header or ?admin_key= for every /admin request.
	Secret string `mapstructure:"secret"`
	// LogFile is an optional path to a log file tail for /admin/logs (e.g. storage/logs/app.log).
	LogFile string `mapstructure:"log_file"`
	// QueueNames lists Redis queue names to show depth for on /admin/metrics (ready + delayed).
	QueueNames []string `mapstructure:"queue_names"`
}

func (c *Config) applyAdminDefaults() {
	if strings.TrimSpace(c.Admin.PathPrefix) == "" {
		c.Admin.PathPrefix = "/admin"
	}
	c.Admin.PathPrefix = "/" + strings.Trim(strings.TrimSpace(c.Admin.PathPrefix), "/")
	if len(c.Admin.QueueNames) == 0 {
		c.Admin.QueueNames = []string{"default", "mails", "emails", "events", "notifications"}
	}
}

func (c *Config) validateAdmin() error {
	if !c.Admin.Enabled {
		return nil
	}
	if strings.EqualFold(strings.TrimSpace(c.Env), "prod") && strings.TrimSpace(c.Admin.Secret) == "" {
		return fmt.Errorf("admin.enabled in prod requires admin.secret (set X-Admin-Key / ?admin_key=)")
	}
	return nil
}
