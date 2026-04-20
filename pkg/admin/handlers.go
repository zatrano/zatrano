package admin

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/config"
	"github.com/zatrano/zatrano/pkg/core"
)

const logTailMax = 2000

func adminURLPrefix(cfg *config.Config) string {
	p := strings.TrimRight(strings.TrimSpace(cfg.Admin.PathPrefix), "/")
	if p == "" {
		return "/admin"
	}
	return p
}

func renderPage(a *core.App, c fiber.Ctx, view string, extra fiber.Map) error {
	if a.View == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "admin HTML için view yapılandırması gerekir (view.*)",
		})
	}
	data := a.View.ViewData(c, extra)
	data["AppName"] = a.Config.AppName
	return c.Render(view, data)
}

func handleDashboard(a *core.App) fiber.Handler {
	return func(c fiber.Ctx) error {
		return renderPage(a, c, "admin/dashboard", fiber.Map{
			"Title":       "Yönetim",
			"AdminPrefix": adminURLPrefix(a.Config),
		})
	}
}

func handleMetrics(a *core.App) fiber.Handler {
	return func(c fiber.Ctx) error {
		snap := collectMetrics(c.Context(), a)
		return renderPage(a, c, "admin/metrics", fiber.Map{
			"Title":       "Metrikler",
			"AdminPrefix": adminURLPrefix(a.Config),
			"Metrics":     snap,
		})
	}
}

func handleLogs(a *core.App) fiber.Handler {
	return func(c fiber.Ctx) error {
		path := strings.TrimSpace(a.Config.Admin.LogFile)
		q := strings.TrimSpace(c.Query("q"))
		var lines []string
		var errMsg string
		if path == "" {
			errMsg = "admin.log_file yapılandırılmadı."
		} else {
			raw, err := tailLogFile(path, logTailMax)
			if err != nil {
				errMsg = err.Error()
			} else {
				lines = filterLogLines(raw, q)
			}
		}
		return renderPage(a, c, "admin/logs", fiber.Map{
			"Title":       "Günlükler",
			"AdminPrefix": adminURLPrefix(a.Config),
			"LogPath":     path,
			"LogQuery":   q,
			"LogLines":   lines,
			"LogError":   errMsg,
			"LogLineCap": logTailMax,
		})
	}
}

func handleLogsDownload(a *core.App) fiber.Handler {
	return func(c fiber.Ctx) error {
		path := strings.TrimSpace(a.Config.Admin.LogFile)
		if path == "" {
			return c.SendStatus(fiber.StatusNotFound)
		}
		name := filepath.Base(path)
		if name == "" || name == "." {
			name = "app.log"
		}
		return c.Download(path, name)
	}
}

func filterLogLines(lines []string, q string) []string {
	if q == "" {
		return lines
	}
	q = strings.ToLower(q)
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		if strings.Contains(strings.ToLower(ln), q) {
			out = append(out, ln)
		}
	}
	return out
}
