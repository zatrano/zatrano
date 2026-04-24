package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Options for generating a new application.
type Options struct {
	Dir         string
	AppName     string
	Module      string
	ZatranoPath string // if set, go.mod gets a replace directive for local ZATRANO development
}

type fileT struct {
	rel  string
	tmpl string
	// raw: when true, write the file as-is; Zatrano view syntax is not run through text/template.Parse.
	raw bool
}

// Run writes the starter layout to Dir.
func Run(opts Options) error {
	if strings.TrimSpace(opts.Dir) == "" {
		return fmt.Errorf("output directory is required")
	}
	if strings.TrimSpace(opts.AppName) == "" {
		return fmt.Errorf("app name is required")
	}
	if strings.TrimSpace(opts.Module) == "" {
		return fmt.Errorf("module path is required (e.g. github.com/acme/myapp)")
	}
	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return err
	}

	rep := strings.TrimSpace(opts.ZatranoPath)
	if rep != "" {
		rep = filepath.ToSlash(rep)
		if strings.ContainsAny(rep, " \t") {
			rep = `"` + rep + `"`
		}
	}
	data := map[string]any{
		"Module":      opts.Module,
		"AppName":     opts.AppName,
		"ReplacePath": rep,
	}

	files := []fileT{
		{filepath.Join("go.mod"), tplGoMod, false},
		{filepath.Join("internal", "routes", "register.go"), tplRoutesRegister, false},
		{filepath.Join("cmd", opts.AppName, "main.go"), tplMain, false},
		{filepath.Join("config", "examples", "dev.yaml"), tplDevYAML, false},
		{filepath.Join("locales", "en.json"), tplLocalesEn, false},
		{filepath.Join("locales", "tr.json"), tplLocalesTr, false},
		{filepath.Join("api", "openapi.yaml"), tplOpenAPI, false},
		{filepath.Join("migrations", "000001_init.up.sql"), tplMigrationUp, false},
		{filepath.Join("migrations", "000001_init.down.sql"), tplMigrationDown, false},
		{filepath.Join("db", "seeds", ".gitkeep"), "", false},
		// View / template system (raw writes; not parsed as scaffold text/template)
		{filepath.Join("views", "layouts", "app.html"), tplViewLayoutApp, true},
		{filepath.Join("views", "layouts", "auth.html"), tplViewLayoutAuth, true},
		{filepath.Join("views", "components", "alert.html"), tplViewComponentAlert, true},
		{filepath.Join("views", "components", "button.html"), tplViewComponentButton, true},
		{filepath.Join("views", "components", "form-input.html"), tplViewComponentFormInput, true},
		{filepath.Join("views", "components", "form-select.html"), tplViewComponentFormSelect, true},
		{filepath.Join("views", "components", "form-textarea.html"), tplViewComponentFormTextarea, true},
		{filepath.Join("views", "components", "csrf.html"), tplViewComponentCSRF, true},
		{filepath.Join("views", "components", "pagination.html"), tplViewComponentPagination, true},
		{filepath.Join("views", "partials", "flash-messages.html"), tplViewPartialFlash, true},
		{filepath.Join("views", "home", "index.html"), tplViewHomeIndex, true},
		// Static asset placeholders
		{filepath.Join("public", "css", ".gitkeep"), "", false},
		{filepath.Join("public", "js", ".gitkeep"), "", false},
		{"README.md", tplReadme, false},
	}

	for _, f := range files {
		out := filepath.Join(opts.Dir, f.rel)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		if f.tmpl == "" {
			if err := os.WriteFile(out, []byte{}, 0o644); err != nil {
				return err
			}
			continue
		}
		if f.raw {
			if err := os.WriteFile(out, []byte(f.tmpl), 0o644); err != nil {
				return err
			}
			continue
		}
		t, err := template.New(f.rel).Parse(f.tmpl)
		if err != nil {
			return fmt.Errorf("parse %s: %w", f.rel, err)
		}
		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("render %s: %w", f.rel, err)
		}
		if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}

const tplGoMod = `module {{.Module}}

go 1.25.0

require github.com/zatrano/zatrano v0.0.0

{{- if .ReplacePath}}

replace github.com/zatrano/zatrano => {{.ReplacePath}}
{{- end}}
`

const tplMain = `package main

import (
	"log"

	"{{.Module}}/internal/routes"
	"github.com/zatrano/zatrano/pkg/zatrano"
)

func main() {
	if err := zatrano.Start(zatrano.StartOptions{
		RegisterRoutes: routes.Register,
	}); err != nil {
		log.Fatal(err)
	}
}
`

const tplRoutesRegister = `package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/core"
	// zatrano:wire:imports:start
	// zatrano:wire:imports:end
)

// Register mounts application modules (updated by zatrano gen module / gen crud).
func Register(a *core.App, app *fiber.App) {
	// Welcome page — renders views/home/index.html via the view engine.
	app.Get("/home", func(c fiber.Ctx) error {
		data := a.View.ViewData(c, fiber.Map{
			"AppName": a.Config.AppName,
		})
		return c.Render("home/index", data)
	})

	// zatrano:wire:register:start
	// zatrano:wire:register:end
}
`

const tplDevYAML = `env: dev
app_name: {{.AppName}}

http_addr: ":8080"
http_read_timeout: 30s

# http:
#   cors_enabled: true
#   cors_allow_origins: ["http://localhost:5173"]

# i18n:
#   enabled: true
#   default_locale: en
#   supported_locales: [en, tr]
#   locales_dir: locales

database_url: ""
database_required: false

redis_url: ""
redis_required: false

log_level: info
log_development: true

migrations_source: file
migrations_dir: migrations
seeds_dir: db/seeds
openapi_path: api/openapi.yaml

security:
  session_enabled: true
  csrf_enabled: true
  csrf_skip_prefixes:
    - /api/
  jwt_secret: "change-me-in-dev-only"
  jwt_issuer: zatrano
  jwt_expiry: 60m
  cookie_secure: false
  demo_token_endpoint: true

static_path: public
static_url_prefix: /static

# View / template engine
view:
  root: views
  extension: .html
  components_dir: components
  layouts_dir: layouts
  dev_mode: true
  asset:
    public_dir: public
    public_url: /public
    # vite_manifest: public/build/.vite/manifest.json
    # vite_dev_url: http://localhost:5173
`

const tplLocalesEn = `{
  "app": {
    "welcome": "Welcome to {{.AppName}}"
  },
  "validation": {
    "required": "This field is required",
    "email": "Must be a valid email address",
    "min": "Must be at least {{"{{.Param}}"}} characters",
    "max": "Must be at most {{"{{.Param}}"}} characters",
    "gte": "Must be greater than or equal to {{"{{.Param}}"}}",
    "lte": "Must be less than or equal to {{"{{.Param}}"}}",
    "len": "Must be exactly {{"{{.Param}}"}} characters",
    "url": "Must be a valid URL",
    "uuid": "Must be a valid UUID",
    "oneof": "Must be one of: {{"{{.Param}}"}}",
    "numeric": "Must be numeric",
    "alpha": "Must contain only letters",
    "alphanum": "Must contain only letters and numbers"
  }
}
`

const tplLocalesTr = `{
  "app": {
    "welcome": "{{.AppName}} uygulamasına hoş geldiniz"
  },
  "validation": {
    "required": "Bu alan zorunludur",
    "email": "Geçerli bir e-posta adresi olmalıdır",
    "min": "En az {{"{{.Param}}"}} karakter olmalıdır",
    "max": "En fazla {{"{{.Param}}"}} karakter olmalıdır",
    "gte": "{{"{{.Param}}"}} değerinden büyük veya eşit olmalıdır",
    "lte": "{{"{{.Param}}"}} değerinden küçük veya eşit olmalıdır",
    "len": "Tam olarak {{"{{.Param}}"}} karakter olmalıdır",
    "url": "Geçerli bir URL olmalıdır",
    "uuid": "Geçerli bir UUID olmalıdır",
    "oneof": "Şunlardan biri olmalıdır: {{"{{.Param}}"}}",
    "numeric": "Sayısal bir değer olmalıdır",
    "alpha": "Sadece harf içermelidir",
    "alphanum": "Sadece harf ve rakam içermelidir"
  }
}
`

const tplOpenAPI = `openapi: 3.0.3
info:
  title: {{.AppName}} API
  version: 0.1.0
paths:
  /api/v1/public/ping:
    get:
      summary: Ping
      responses:
        "200":
          description: OK
`

const tplMigrationUp = `-- {{.AppName}} initial migration
CREATE TABLE IF NOT EXISTS app_hello (
    id bigserial PRIMARY KEY,
    message text NOT NULL DEFAULT 'hello from {{.AppName}}',
    created_at timestamptz NOT NULL DEFAULT now()
);

`

const tplMigrationDown = `DROP TABLE IF EXISTS app_hello;
`

const tplReadme = `# {{.AppName}}

Generated by [ZATRANO](https://github.com/zatrano/zatrano).

## Run

` + "```bash" + `
cp config/examples/dev.yaml config/dev.yaml
go mod tidy
go run ./cmd/{{.AppName}}
` + "```" + `

Set ` + "`DATABASE_URL`" + ` and ` + "`REDIS_URL`" + ` when you enable Postgres/Redis.

## Modules

` + "```bash" + `
zatrano gen module my_feature
go fmt ./...
` + "```" + `

This updates ` + "`internal/routes/register.go`" + ` and runs ` + "`go fmt`" + ` on it. Use ` + "`--skip-wire`" + ` to only generate files under ` + "`modules/`" + `, then ` + "`zatrano gen wire <name>`" + ` when ready.

## Checks (optional)

` + "```bash" + `
zatrano verify
` + "```" + `

Runs ` + "`go vet`" + `, ` + "`go test`" + `, and merged OpenAPI validation (install ` + "`zatrano`" + ` on PATH first).

## Migrations & seeds

` + "```bash" + `
go install github.com/zatrano/zatrano/cmd/zatrano@latest
zatrano db migrate
zatrano db seed   # after adding .sql files under db/seeds/
` + "```" + `
`

// ─── View template constants ─────────────────────────────────────────────────

const tplViewLayoutApp = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{block "title" .}}{{if .AppName}}{{.AppName}}{{else}}APP{{end}}{{end}}</title>
  {{if .CSRF}}<meta name="csrf-token" content="{{.CSRF}}">{{end}}
  {{assetLink "css/app.css"}}
  {{block "head" .}}{{end}}
</head>
<body class="{{block "body_class" .}}{{end}}">

  {{block "header" .}}
  <header>
    <nav>
      <a href="/">Home</a>
    </nav>
  </header>
  {{end}}

  {{template "partials/flash-messages" .}}

  <main>
    {{block "content" .}}{{end}}
  </main>

  {{block "footer" .}}
  <footer>
    <p>Powered by ZATRANO</p>
  </footer>
  {{end}}

  {{assetScript "js/app.js"}}
  {{block "scripts" .}}{{end}}
  <script>window.__CSRF__ = "{{.CSRF}}";</script>
</body>
</html>
`

const tplViewLayoutAuth = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{block "title" .}}Auth{{end}}</title>
  {{if .CSRF}}<meta name="csrf-token" content="{{.CSRF}}">{{end}}
  {{assetLink "css/app.css"}}
  {{block "head" .}}{{end}}
</head>
<body class="auth-layout {{block "body_class" .}}{{end}}">
  <main class="auth-main">
    {{template "partials/flash-messages" .}}
    {{block "content" .}}{{end}}
  </main>
  {{assetScript "js/app.js"}}
  {{block "scripts" .}}{{end}}
</body>
</html>
`

const tplViewComponentAlert = `{{define "components/alert"}}
<div role="alert" class="alert alert-{{.Type}}">
  <span class="alert-message">{{.Message}}</span>
</div>
{{end}}
`

const tplViewComponentButton = `{{define "components/button"}}
<button
  type="{{if .Type}}{{.Type}}{{else}}button{{end}}"
  class="btn btn-{{if .Variant}}{{.Variant}}{{else}}primary{{end}}{{if .Class}} {{.Class}}{{end}}"
  {{if .Attrs}}{{safe .Attrs}}{{end}}>
  {{.Label}}
</button>
{{end}}
`

const tplViewComponentFormInput = `{{define "components/form-input"}}
<div class="form-group{{if .Error}} has-error{{end}}{{if .Class}} {{.Class}}{{end}}">
  {{if .Label}}
  <label for="input-{{.Name}}" class="form-label">
    {{.Label}}{{if .Required}}<span class="required" aria-hidden="true">*</span>{{end}}
  </label>
  {{end}}
  <input
    type="{{if .Type}}{{.Type}}{{else}}text{{end}}"
    id="input-{{.Name}}"
    name="{{.Name}}"
    value="{{.Value}}"
    {{if .Placeholder}}placeholder="{{.Placeholder}}"{{end}}
    {{if .Required}}required{{end}}
    class="form-control{{if .Error}} is-invalid{{end}}"
    {{if .Attrs}}{{safe .Attrs}}{{end}}
  >
  {{if .Error}}
  <p class="form-error" role="alert">{{.Error}}</p>
  {{else if .Hint}}
  <p class="form-hint">{{.Hint}}</p>
  {{end}}
</div>
{{end}}
`

const tplViewComponentFormSelect = `{{define "components/form-select"}}
<div class="form-group{{if .Error}} has-error{{end}}{{if .Class}} {{.Class}}{{end}}">
  {{if .Label}}
  <label for="select-{{.Name}}" class="form-label">
    {{.Label}}{{if .Required}}<span class="required" aria-hidden="true">*</span>{{end}}
  </label>
  {{end}}
  <select
    id="select-{{.Name}}"
    name="{{.Name}}"
    {{if .Required}}required{{end}}
    class="form-control form-select{{if .Error}} is-invalid{{end}}"
    {{if .Attrs}}{{safe .Attrs}}{{end}}
  >
    {{range .Options}}
    <option value="{{index . 0}}"{{if eq (index . 0) $.Value}} selected{{end}}>{{index . 1}}</option>
    {{end}}
  </select>
  {{if .Error}}
  <p class="form-error" role="alert">{{.Error}}</p>
  {{else if .Hint}}
  <p class="form-hint">{{.Hint}}</p>
  {{end}}
</div>
{{end}}
`

const tplViewComponentFormTextarea = `{{define "components/form-textarea"}}
<div class="form-group{{if .Error}} has-error{{end}}{{if .Class}} {{.Class}}{{end}}">
  {{if .Label}}
  <label for="textarea-{{.Name}}" class="form-label">
    {{.Label}}{{if .Required}}<span class="required" aria-hidden="true">*</span>{{end}}
  </label>
  {{end}}
  <textarea
    id="textarea-{{.Name}}"
    name="{{.Name}}"
    rows="{{if .Rows}}{{.Rows}}{{else}}4{{end}}"
    {{if .Placeholder}}placeholder="{{.Placeholder}}"{{end}}
    {{if .Required}}required{{end}}
    class="form-control{{if .Error}} is-invalid{{end}}"
    {{if .Attrs}}{{safe .Attrs}}{{end}}
  >{{.Value}}</textarea>
  {{if .Error}}
  <p class="form-error" role="alert">{{.Error}}</p>
  {{else if .Hint}}
  <p class="form-hint">{{.Hint}}</p>
  {{end}}
</div>
{{end}}
`

const tplViewComponentCSRF = `{{define "components/csrf"}}
<input type="hidden" name="_csrf" value="{{.CSRF}}">
{{end}}
`

const tplViewComponentPagination = `{{define "components/pagination"}}
{{if gt .TotalPages 1}}
<nav class="pagination" aria-label="Pagination">
  {{if gt .CurrentPage 1}}
  <a href="?page={{sub .CurrentPage 1}}" class="pagination-prev">&laquo; Prev</a>
  {{end}}
  {{range iterate .TotalPages}}
  {{$page := add . 1}}
  <a href="?page={{$page}}"
     class="pagination-page{{if eq $page $.CurrentPage}} active{{end}}"
     {{if eq $page $.CurrentPage}}aria-current="page"{{end}}>
    {{$page}}
  </a>
  {{end}}
  {{if lt .CurrentPage .TotalPages}}
  <a href="?page={{add .CurrentPage 1}}" class="pagination-next">Next &raquo;</a>
  {{end}}
</nav>
{{end}}
{{end}}
`

const tplViewPartialFlash = `{{define "partials/flash-messages"}}
{{if .Flash}}
<div class="flash-container" aria-live="polite">
  {{range .Flash}}
    {{template "components/alert" .}}
  {{end}}
</div>
{{end}}
{{end}}
`

const tplViewHomeIndex = `{{extends "layouts/app"}}

{{block "title"}}Welcome — {{.AppName}}{{end}}

{{block "content"}}
<div style="max-width:640px;margin:4rem auto;text-align:center;font-family:system-ui,sans-serif">
  <h1 style="font-size:2.5rem;margin-bottom:.5rem">🚀 {{.AppName}}</h1>
  <p style="color:#666;margin-bottom:2rem">Generated by <strong>ZATRANO</strong></p>
  <div style="background:#f4f4f5;border-radius:8px;padding:1.5rem;text-align:left">
    <p style="margin:0 0 .5rem;font-weight:600">Next steps:</p>
    <pre style="margin:0;font-size:.85rem">zatrano gen module my_feature --with-form
zatrano gen view my_feature --with-form
go fmt ./...</pre>
  </div>
</div>
{{end}}
`
