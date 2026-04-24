package admin

import (
	"path/filepath"
	"strings"
	"unicode/utf8"

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

func adminInitials(appName string) string {
	s := strings.TrimSpace(appName)
	if s == "" {
		return "ZA"
	}
	s = strings.ToUpper(s)
	r, w := utf8.DecodeRuneInString(s)
	if w == 0 {
		return "ZA"
	}
	if len(s) <= w {
		return string(r)
	}
	r2, _ := utf8.DecodeRuneInString(s[w:])
	if r2 == utf8.RuneError {
		return string(r)
	}
	return string(r) + string(r2)
}

func renderAdmin(a *core.App, c fiber.Ctx, view string, extra fiber.Map) error {
	if a.View == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "admin HTML için view yapılandırması gerekir (view.*)",
		})
	}
	data := a.View.ViewData(c, extra)
	data["AppName"] = a.Config.AppName
	data["AdminInitials"] = adminInitials(a.Config.AppName)
	data["AdminSubtitle"] = "Operasyon paneli"
	return c.Render(view, data)
}

func handleDashboard(a *core.App) fiber.Handler {
	return func(c fiber.Ctx) error {
		return renderAdmin(a, c, "admin/dashboard", fiber.Map{
			"Title":       "Yönetim",
			"AdminPrefix": adminURLPrefix(a.Config),
			"AdminNav":    "dashboard",
		})
	}
}

func handleMetrics(a *core.App) fiber.Handler {
	return func(c fiber.Ctx) error {
		snap := collectMetrics(c.Context(), a)
		return renderAdmin(a, c, "admin/metrics", fiber.Map{
			"Title":       "Metrikler",
			"AdminPrefix": adminURLPrefix(a.Config),
			"AdminNav":    "metrics",
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
		return renderAdmin(a, c, "admin/logs", fiber.Map{
			"Title":       "Günlükler",
			"AdminPrefix": adminURLPrefix(a.Config),
			"AdminNav":    "logs",
			"LogPath":     path,
			"LogQuery":    q,
			"LogLines":    lines,
			"LogError":    errMsg,
			"LogLineCap":  logTailMax,
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
