package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/features"
)

// RequireFeature returns middleware that continues only when the flag is enabled
// for the current user (same resolution as features.FromFiber(reg).IsEnabled).
// Anonymous users fail role-gated or partial-rollout flags.
func RequireFeature(reg *features.Registry, key string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if reg == nil || !reg.Enabled() {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		ev := reg.FromFiber(c)
		if !ev.IsEnabled(c.Context(), key) {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}
		return c.Next()
	}
}
