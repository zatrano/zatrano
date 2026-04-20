package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/tenant"
)

// LocalsTenant is the Fiber Locals key for *tenant.Info (same tenant as request context).
const LocalsTenant = "zatrano.tenant"

// ResolveTenant resolves the current tenant from header or subdomain and stores it in
// Fiber Locals and the request context.Context (for GORM / repository.NewTenantAware).
func ResolveTenant(cfg *config.Config) fiber.Handler {
	if cfg == nil || !cfg.Tenant.Enabled {
		return func(c fiber.Ctx) error { return c.Next() }
	}
	mode := strings.ToLower(strings.TrimSpace(cfg.Tenant.Mode))
	isolation := strings.ToLower(strings.TrimSpace(cfg.Tenant.Isolation))
	prefix := cfg.Tenant.SchemaPrefix

	return func(c fiber.Ctx) error {
		key := resolveTenantKey(c, cfg, mode)
		if strings.TrimSpace(key) == "" {
			if cfg.Tenant.Required {
				return fiber.NewError(fiber.StatusBadRequest, "tenant is required")
			}
			return c.Next()
		}

		info := tenant.Info{
			Key:       key,
			NumericID: tenant.ParseNumericKey(key),
		}
		if isolation == "schema" {
			sch, err := tenant.SchemaName(prefix, key)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			info.Schema = sch
		}

		c.Locals(LocalsTenant, &info)
		c.SetContext(tenant.WithContext(c.Context(), info))
		return c.Next()
	}
}

func resolveTenantKey(c fiber.Ctx, cfg *config.Config, mode string) string {
	switch mode {
	case "subdomain":
		host := strings.ToLower(c.Hostname())
		suf := strings.ToLower(strings.TrimSpace(cfg.Tenant.SubdomainSuffix))
		if suf == "" || !strings.HasSuffix(host, suf) {
			return ""
		}
		left := strings.TrimSuffix(host, suf)
		left = strings.TrimSuffix(left, ".")
		return strings.TrimSpace(left)
	default: // header
		h := strings.TrimSpace(cfg.Tenant.HeaderName)
		if h == "" {
			h = "X-Tenant-ID"
		}
		return strings.TrimSpace(c.Get(h))
	}
}
