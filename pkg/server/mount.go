package server

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/zatrano/zatrano/pkg/admin"
	"github.com/zatrano/zatrano/pkg/broadcast"
	"github.com/zatrano/zatrano/pkg/core"
	zgraphql "github.com/zatrano/zatrano/pkg/graphql"
	"github.com/zatrano/zatrano/pkg/health"
	"github.com/zatrano/zatrano/pkg/i18n"
	"github.com/zatrano/zatrano/pkg/meta"
	zmw "github.com/zatrano/zatrano/pkg/middleware"
	"github.com/zatrano/zatrano/pkg/oauth"
	"github.com/zatrano/zatrano/pkg/openapi"
	"github.com/zatrano/zatrano/pkg/security"
)

func bodyLimitJSON(n int) string {
	if n <= 0 {
		return "default"
	}
	return strconv.Itoa(n)
}

// MountOptions configures optional route registration for application modules.
type MountOptions struct {
	// RegisterRoutes runs after registerAPI and registerModules, before OAuth and OpenAPI.
	RegisterRoutes func(a *core.App, app *fiber.App)
}

// Mount wires global middlewares and routes for the HTTP server.
func Mount(a *core.App, app *fiber.App, opts MountOptions) {
	app.Use(recover.New())
	app.Use(requestid.New())
	if a.Config.Tenant.Enabled {
		app.Use(zmw.ResolveTenant(a.Config))
	}
	registerI18nMiddleware(a, app)
	registerHTTPMiddleware(a, app)
	app.Use(helmet.New())
	app.Use(compress.New())

	security.RegisterSessionAndCSRF(a, app)

	app.Use(zmw.AccessLog(a.Log))
	if a.Config.Audit.Enabled && a.Config.Audit.HttpEnabled {
		app.Use(zmw.AuditLog(a.Config, a.Audit, a.Log))
	}

	// View middleware: injects flash manager into Locals for every request.
	if a.View != nil {
		app.Use(a.View.Middleware())
	}
	if a.Features != nil && a.Config.Features.Enabled {
		app.Use(a.Features.LocalsMiddleware())
	}

	app.Get("/", func(c fiber.Ctx) error {
		ep := fiber.Map{
			"health":   "/health",
			"ready":    "/ready",
			"status":   "/status",
			"openapi":  "/openapi.yaml",
			"docs":     "/docs",
			"api_ping": "/api/v1/public/ping",
			"api_me":   "/api/v1/private/me",
		}
		if a.Config.OAuth.Enabled {
			ep["oauth_google"] = "/auth/oauth/google/login"
			ep["oauth_github"] = "/auth/oauth/github/login"
		}
		if a.Broadcast != nil && a.Config.Broadcast.Enabled {
			p := strings.TrimRight(strings.TrimSpace(a.Config.Broadcast.PathPrefix), "/")
			if p == "" {
				p = "/broadcast"
			}
			ep["broadcast_ws"] = p + "/ws"
			if a.Config.Broadcast.SSEEnabled {
				ep["broadcast_sse"] = p + "/sse/{channel}"
			}
		}
		if a.Config.GraphQL.Enabled {
			gp := strings.TrimSpace(a.Config.GraphQL.Path)
			if gp == "" {
				gp = "/graphql"
			}
			ep["graphql"] = gp
			if a.Config.GraphQL.Playground {
				pp := strings.TrimSpace(a.Config.GraphQL.PlaygroundPath)
				if pp == "" {
					pp = "/playground"
				}
				ep["graphql_playground"] = pp
			}
		}
		h := a.Config.HTTP
		ic := a.Config.I18n
		idx := fiber.Map{
			"name":      a.Config.AppName,
			"env":       a.Config.Env,
			"version":   meta.Version,
			"endpoints": ep,
			// Set by requestid middleware; included in JSON error bodies for support correlation.
			"error_includes_request_id": true,
			"http": fiber.Map{
				"cors_enabled":       h.CORSEnabled,
				"rate_limit_enabled": h.RateLimitEnabled,
				"request_timeout":    h.RequestTimeout.String(),
				"body_limit":         bodyLimitJSON(h.BodyLimit),
			},
		}
		if ic.Enabled {
			idx["i18n"] = fiber.Map{
				"enabled":           true,
				"default_locale":    ic.DefaultLocale,
				"supported_locales": ic.SupportedLocales,
				"active_locale":     i18n.Locale(c),
			}
		} else {
			idx["i18n"] = fiber.Map{"enabled": false}
		}
		return c.JSON(idx)
	})

	health.Register(a, app)
	registerAPI(a, app)
	zgraphql.Register(a, app)
	broadcast.Register(a.Broadcast, a.Config, app)
	registerModules(a, app)
	admin.Register(a, app)
	if opts.RegisterRoutes != nil {
		opts.RegisterRoutes(a, app)
	}
	oauth.Register(a, app)
	openapi.Register(a, app)

	if p := a.Config.StaticPath; p != "" {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			prefix := a.Config.StaticURLPrefix
			if prefix == "" {
				prefix = "/static"
			}
			app.Use(prefix, static.New(p))
		}
	}
}
