// Package view integrates the zatrano view engine with Fiber.
// It provides a Renderer that can be registered as Fiber's Views engine, wires
// in flash messages, old input, CSRF token, and the asset manager helpers.
package view

import (
	"html/template"
	"io"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"

	"github.com/zatrano/zatrano/pkg/features"
	"github.com/zatrano/zatrano/pkg/view/asset"
	"github.com/zatrano/zatrano/pkg/view/engine"
	"github.com/zatrano/zatrano/pkg/view/flash"
)

// Renderer implements fiber.Views and wraps the zatrano view engine.
// Register it in Fiber via fiber.Config{Views: renderer}.
type Renderer struct {
	engine   *engine.Engine
	assets   *asset.Manager
	session  *session.Store
	flash    *flash.Manager
	features *features.Registry
}

// Config holds the combined configuration for the Renderer.
type Config struct {
	Engine engine.Config
	Asset  asset.Config
	// Features is set at bootstrap (not loaded from YAML); optional.
	Features *features.Registry `mapstructure:"-" yaml:"-"`
}

// New returns a Renderer.  Pass a session.Store when you want flash and
// old-input support; pass nil to disable those features.
func New(cfg Config, store *session.Store) *Renderer {
	eng := engine.New(cfg.Engine)

	// Merge asset template functions into the engine's FuncMap.
	mgr := asset.New(cfg.Asset)
	if cfg.Engine.FuncMap == nil {
		cfg.Engine.FuncMap = template.FuncMap{}
	}
	for k, v := range mgr.TemplateFuncs() {
		cfg.Engine.FuncMap[k] = v
	}
	for k, v := range features.TemplateFuncMap() {
		cfg.Engine.FuncMap[k] = v
	}
	// Re-create engine with the merged FuncMap.
	eng = engine.New(cfg.Engine)

	r := &Renderer{
		engine:   eng,
		assets:   mgr,
		features: cfg.Features,
	}
	if store != nil {
		r.session = store
		r.flash = flash.New(store)
	}
	return r
}

// Load implements fiber.Views.
func (r *Renderer) Load() error { return nil }

// Render implements fiber.Views. The binding map is expected to be
// map[string]any or fiber.Map; flash/old/CSRF are injected automatically
// when a session store is configured.
func (r *Renderer) Render(w io.Writer, name string, binding any, layout ...string) error {
	data, _ := binding.(map[string]any)
	if data == nil {
		data = make(map[string]any)
	}

	return r.engine.Render(w, name, data)
}

// Engine returns the underlying engine so callers can call RenderString / RenderBytes directly.
func (r *Renderer) Engine() *engine.Engine { return r.engine }

// Assets returns the asset manager.
func (r *Renderer) Assets() *asset.Manager { return r.assets }

// Flash returns the flash manager (may be nil if no session store was provided).
func (r *Renderer) Flash() *flash.Manager { return r.flash }

// ClearCache discards all cached templates and asset hashes.
func (r *Renderer) ClearCache() {
	r.engine.ClearCache()
	r.assets.ClearCache()
}

// ----------------------------------------------------------------------------
// Fiber middleware — enriches c.Locals with flash / old-input / CSRF
// ----------------------------------------------------------------------------

// Middleware returns a Fiber handler that injects flash messages, old input, and
// the CSRF token into c.Locals so templates can access them via the data map.
//
// Usage:
//
//	app.Use(renderer.Middleware())
func (r *Renderer) Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		if r.flash != nil {
			c.Locals("zatrano.flash.manager", r.flash)
		}
		return c.Next()
	}
}

// ViewData builds a map[string]any pre-populated with flash messages, old input
// and the CSRF token for use as template binding data.  Additional keys can be
// merged via extra.
//
// Usage in a handler:
//
//	data := renderer.ViewData(c, fiber.Map{"Title": "Login"})
//	return c.Render("auth/login", data)
func (r *Renderer) ViewData(c fiber.Ctx, extra ...fiber.Map) map[string]any {
	data := make(map[string]any)

	// Merge caller-supplied keys first.
	for _, m := range extra {
		for k, v := range m {
			data[k] = v
		}
	}

	// Flash + old input.
	if r.flash != nil {
		msgs, _ := r.flash.Get(c)
		old, _ := r.flash.GetOld(c)
		if msgs == nil {
			msgs = []flash.Message{}
		}
		if old == nil {
			old = map[string]string{}
		}
		data["Flash"] = msgs
		data["Old"] = old

		// OldFn is a template-callable function: {{call .OldFn "email"}}
		data["OldFn"] = func(field string) string {
			return old[field]
		}
	}

	// CSRF token — try Fiber locals first (set by gofiber csrf middleware), then header.
	if tok := csrf.TokenFromContext(c); tok != "" {
		data["CSRF"] = tok
	} else if h := c.Get("X-CSRF-Token"); h != "" {
		data["CSRF"] = h
	}

	// App globals that templates commonly need.
	data["Path"] = c.Path()
	data["Method"] = c.Method()

	if r.features != nil && r.features.Enabled() {
		data[features.TemplateBindKey] = r.features.FromFiber(c)
	}

	return data
}
