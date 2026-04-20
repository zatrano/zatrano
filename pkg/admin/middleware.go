package admin

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/config"
)

// Middleware enforces admin.secret when set (X-Admin-Key header or ?admin_key=).
func Middleware(cfg *config.Config) fiber.Handler {
	secret := strings.TrimSpace(cfg.Admin.Secret)
	return func(c fiber.Ctx) error {
		if secret == "" {
			return c.Next()
		}
		hdr := strings.TrimSpace(c.Get("X-Admin-Key"))
		q := strings.TrimSpace(c.Query("admin_key"))
		if hdr != secret && q != secret {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}
