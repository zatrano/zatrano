package admin

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/core"
)

// Register mounts dashboard, metrics, and log viewer when admin.enabled is true.
func Register(a *core.App, app *fiber.App) {
	if a == nil || !a.Config.Admin.Enabled {
		return
	}
	prefix := strings.TrimRight(strings.TrimSpace(a.Config.Admin.PathPrefix), "/")
	if prefix == "" {
		prefix = "/admin"
	}
	g := app.Group(prefix, Middleware(a.Config))
	g.Get("/", handleDashboard(a))
	g.Get("/metrics", handleMetrics(a))
	g.Get("/logs", handleLogs(a))
	g.Get("/logs/download", handleLogsDownload(a))
}
