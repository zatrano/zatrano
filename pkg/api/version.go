package api

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/zatrano/zatrano/pkg/config"
)

// VersionManager handles API versioning with automatic route prefixing.
type VersionManager struct {
	version string
	prefix  string
}

// NewVersionManager creates a version manager from config.
func NewVersionManager(cfg *config.Config) *VersionManager {
	return &VersionManager{
		version: cfg.API.Version,
		prefix:  cfg.API.Prefix,
	}
}

// Group creates a versioned route group.
func (v *VersionManager) Group(app *fiber.App, path string) fiber.Router {
	fullPath := strings.TrimSuffix(v.prefix, "/") + "/" + strings.TrimPrefix(path, "/")
	return app.Group(fullPath)
}

// CurrentVersion returns the active API version.
func (v *VersionManager) CurrentVersion() string {
	return v.version
}

// Prefix returns the full API prefix.
func (v *VersionManager) Prefix() string {
	return v.prefix
}

// Middleware adds version headers to responses.
func (v *VersionManager) Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("X-API-Version", v.version)
		return c.Next()
	}
}
