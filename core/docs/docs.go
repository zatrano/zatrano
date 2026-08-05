package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zatrano/framework/core/http"
	"github.com/zatrano/framework/core/markdown"
	"github.com/zatrano/framework/core/routing"
)

// Page is a documentation markdown page.
type Page struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

// Repository loads markdown docs from a directory.
type Repository struct {
	root string
}

// New creates a docs repository rooted at dir.
func New(root string) *Repository {
	return &Repository{root: root}
}

// List returns available doc pages (without full content).
func (r *Repository) List() ([]Page, error) {
	entries, err := os.ReadDir(r.root)
	if err != nil {
		if os.IsNotExist(err) {
			return []Page{}, nil
		}
		return nil, err
	}
	pages := make([]Page, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		slug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		title := humanize(slug)
		pages = append(pages, Page{
			Slug:  slug,
			Title: title,
			Path:  filepath.Join(r.root, entry.Name()),
		})
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].Slug < pages[j].Slug })
	return pages, nil
}

// Get loads a page by slug.
func (r *Repository) Get(slug string) (*Page, error) {
	slug = strings.Trim(slug, "/")
	if slug == "" {
		slug = "index"
	}
	path := filepath.Join(r.root, slug+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	title := slug
	if lines := strings.SplitN(string(raw), "\n", 2); len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "# ") {
		title = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[0]), "# "))
	}
	return &Page{
		Slug:    slug,
		Title:   title,
		Path:    path,
		Content: string(raw),
	}, nil
}

// HTML returns rendered HTML for a page.
func (r *Repository) HTML(slug string) (string, *Page, error) {
	page, err := r.Get(slug)
	if err != nil {
		return "", nil, err
	}
	return markdown.ToHTML(page.Content), page, nil
}

// IndexHandler lists documentation pages.
func (r *Repository) IndexHandler() routing.HandlerFunc {
	return func(req *http.Request) *http.Response {
		pages, err := r.List()
		if err != nil {
			return http.JSON(map[string]any{"message": err.Error()}).Status(500)
		}
		if req.WantsJSON() {
			return http.JSON(map[string]any{"data": pages})
		}
		var b strings.Builder
		b.WriteString(`<!doctype html><html><head><meta charset="utf-8"><title>Docs</title>
<style>body{font-family:ui-sans-serif,system-ui;background:#0b1220;color:#e8eef8;padding:2rem;max-width:800px;margin:0 auto}
h1{color:#3dd6c6}a{color:#3dd6c6}li{margin:.4rem 0}</style></head><body>`)
		b.WriteString(`<h1>Documentation</h1><ul>`)
		for _, page := range pages {
			b.WriteString(fmt.Sprintf(`<li><a href="/documentation/%s">%s</a></li>`, page.Slug, page.Title))
		}
		b.WriteString(`</ul></body></html>`)
		return http.HTML(b.String())
	}
}

// ShowHandler renders a markdown doc page.
func (r *Repository) ShowHandler() routing.HandlerFunc {
	return func(req *http.Request) *http.Response {
		slug := req.Route("slug", "index")
		html, page, err := r.HTML(slug)
		if err != nil {
			if req.WantsJSON() {
				return http.JSON(map[string]any{"message": "Document not found"}).Status(404)
			}
			return http.HTML("<h1>Not Found</h1>").Status(404)
		}
		if req.WantsJSON() {
			return http.JSON(map[string]any{"slug": page.Slug, "title": page.Title, "content": page.Content})
		}
		body := fmt.Sprintf(`<!doctype html><html><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:ui-sans-serif,system-ui;background:#0b1220;color:#e8eef8;padding:2rem;max-width:800px;margin:0 auto;line-height:1.6}
h1,h2,h3{color:#3dd6c6}a{color:#3dd6c6}code{background:#121a2b;padding:.1rem .35rem;border-radius:4px}
pre{background:#121a2b;padding:1rem;border-radius:8px;overflow:auto}</style></head>
<body><p><a href="/documentation">&larr; Docs</a></p>%s</body></html>`, page.Title, html)
		return http.HTML(body)
	}
}

func humanize(slug string) string {
	parts := strings.Split(strings.ReplaceAll(slug, "_", "-"), "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
