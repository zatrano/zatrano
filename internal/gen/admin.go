package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Admin generates HTML admin list scaffolding under modules/<name>/ and
// views/admin/<name>/index.html.
// Requires modules/<name>/ (run `zatrano gen module <name>` first).
// Expects models.%[Pascal] in package models (e.g. zatrano gen model <name> under models/).
func Admin(moduleRoot, baseDir, rawName string, dryRun bool) ([]string, error) {
	name := normalizeName(rawName)
	if name == "" {
		return nil, fmt.Errorf("invalid model / module name %q", rawName)
	}
	modDir := filepath.Join(moduleRoot, baseDir, name)
	if fi, err := os.Stat(modDir); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("modules klasörü yok: %s (önce zatrano gen module %s çalıştırın)", filepath.ToSlash(modDir), rawName)
	}
	modPath, err := ModuleImportPath(moduleRoot)
	if err != nil {
		return nil, err
	}
	pascal := snakeToPascal(name)
	table := adminTableName(name)
	modelImport := modPath + "/models"

	viewPath := filepath.Join(moduleRoot, "views", "admin", name, "index.html")

	files := []struct {
		path string
		body string
	}{
		{filepath.Join(modDir, "admin_handlers.go"), tmplAdminHandlers(name, pascal, modelImport)},
		{filepath.Join(modDir, "admin_register.go"), tmplAdminRegister(name, pascal)},
		{viewPath, tmplAdminView(name, pascal, table)},
	}
	var written []string
	for _, f := range files {
		written = append(written, f.path)
		if dryRun {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(f.path, []byte(f.body), 0o644); err != nil {
			return nil, err
		}
	}
	return written, nil
}

func adminTableName(name string) string {
	base := strings.ReplaceAll(name, "-", "_")
	return base + "s"
}

func tmplAdminHandlers(pkg, pascal, modelImport string) string {
	return fmt.Sprintf(`package %[1]s

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/admin"
	"github.com/zatrano/zatrano/pkg/core"
	"%[3]s"
)

// %[2]sAdminHandler is a scaffolded admin CRUD surface (HTML list + TODO mutating routes).
type %[2]sAdminHandler struct {
	app *core.App
}

// New%[2]sAdminHandler constructs the handler.
func New%[2]sAdminHandler(app *core.App) *%[2]sAdminHandler {
	return &%[2]sAdminHandler{app: app}
}

// Index lists rows with optional id search (?q=) and pagination (?page=).
func (h *%[2]sAdminHandler) Index(c fiber.Ctx) error {
	a := h.app
	if a == nil || a.DB == nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString("veritabanı yok")
	}
	if a.View == nil {
		return c.Status(fiber.StatusServiceUnavailable).SendString("view yapılandırması gerekir")
	}
	q := strings.TrimSpace(c.Query("q"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 20
	offset := (page - 1) * perPage

	base := a.DB.WithContext(c.Context()).Model(&models.%[2]s{})
	if q != "" {
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			base = base.Where("id = ?", id)
		}
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	var items []models.%[2]s
	listQ := a.DB.WithContext(c.Context()).Model(&models.%[2]s{}).Order("id DESC").Limit(perPage).Offset(offset)
	if q != "" {
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			listQ = listQ.Where("id = ?", id)
		}
	}
	if err := listQ.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	prefix := admin.URLPrefix(a.Config)
	data := a.View.ViewData(c, fiber.Map{
		"Title":        %[2]q + " — Admin",
		"AdminPrefix":  prefix,
		"ResourceName": %[2]q,
		"Items":        items,
		"Total":        total,
		"Query":        q,
		"CurrentPage":  page,
		"TotalPages":   totalPages,
		"PerPage":      perPage,
		"ListBaseURL":  prefix + "/%[1]s",
	})
	data["AppName"] = a.Config.AppName
	return c.Render("admin/%[1]s/index", data)
}

// New renders a placeholder create form.
func (h *%[2]sAdminHandler) New(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("TODO: create form")
}

// Create accepts POST /admin/<pkg>/ .
func (h *%[2]sAdminHandler) Create(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("TODO: persist create")
}

// Edit renders edit placeholder.
func (h *%[2]sAdminHandler) Edit(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("TODO: edit form")
}

// Update handles PUT/PATCH.
func (h *%[2]sAdminHandler) Update(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("TODO: persist update")
}

// Destroy removes a row.
func (h *%[2]sAdminHandler) Destroy(c fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("TODO: destroy")
}
`, pkg, pascal, modelImport)
}

func tmplAdminRegister(pkg, pascal string) string {
	return fmt.Sprintf(`package %[1]s

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/admin"
	"github.com/zatrano/zatrano/pkg/core"
)

// RegisterAdmin mounts HTML admin routes for %[1]s under admin.path_prefix (config).
func RegisterAdmin(a *core.App, app *fiber.App) {
	if a == nil || !a.Config.Admin.Enabled {
		return
	}
	p := strings.TrimRight(strings.TrimSpace(a.Config.Admin.PathPrefix), "/")
	if p == "" {
		p = "/admin"
	}
	h := New%[2]sAdminHandler(a)
	g := app.Group(p+"/%[1]s", admin.Middleware(a.Config))
	g.Get("/", h.Index)
	g.Get("/new", h.New)
	g.Post("/", h.Create)
	g.Get("/:id/edit", h.Edit)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Destroy)
}
`, pkg, pascal)
}

func tmplAdminView(pkg, pascal, table string) string {
	return fmt.Sprintf(`{{extends "layouts/app"}}
{{define "title"}}{{.Title}} — {{.AppName}}{{end}}

{{define "nav"}}
<a class="navbar-link" href="{{.AdminPrefix}}">Özet</a>
<a class="navbar-link" href="{{.AdminPrefix}}/metrics">Metrikler</a>
<a class="navbar-link" href="{{.AdminPrefix}}/logs">Günlükler</a>
<span class="navbar-link active">%[2]s</span>
{{end}}

{{define "content"}}
<section class="admin-crud">
  <h1>%[2]s</h1>
  <p class="text-muted">GORM tablosu (varsayılan): <code>%[3]s</code> · sayfa başına {{.PerPage}}</p>

  <form method="get" action="{{.ListBaseURL}}" class="admin-log-filter">
    <label for="q">Ara (id)</label>
    <input id="q" name="q" type="search" value="{{.Query}}" placeholder="ör. 12">
    <button type="submit">Ara</button>
    {{if .Query}}<a href="{{.ListBaseURL}}">Temizle</a>{{end}}
  </form>

  <table class="data-table">
    <thead><tr><th>ID</th><th>Oluşturulma</th><th>Güncellenme</th></tr></thead>
    <tbody>
      {{range .Items}}
      <tr>
        <td>{{.ID}}</td>
        <td>{{.CreatedAt}}</td>
        <td>{{.UpdatedAt}}</td>
      </tr>
      {{else}}
      <tr><td colspan="3">Kayıt yok</td></tr>
      {{end}}
    </tbody>
  </table>

  {{template "components/pagination" (dict "CurrentPage" .CurrentPage "TotalPages" .TotalPages "BaseURL" .ListBaseURL)}}
</section>
{{end}}
`, pkg, pascal, table)
}
