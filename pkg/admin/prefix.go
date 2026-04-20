package admin

import (
	"strings"

	"github.com/zatrano/zatrano/pkg/config"
)

// URLPrefix returns the normalized admin.path_prefix (default /admin, no trailing slash).
func URLPrefix(cfg *config.Config) string {
	if cfg == nil {
		return "/admin"
	}
	p := strings.TrimRight(strings.TrimSpace(cfg.Admin.PathPrefix), "/")
	if p == "" {
		return "/admin"
	}
	return p
}
