// Package engine provides a Go html/template-based view engine with layout
// inheritance ({{extends "layouts/app"}}), section/block system, component
// partials, and a rich set of template helpers (form builder, flash, old input,
// asset versioning, Vite/esbuild manifest support).
package engine

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config holds the view engine configuration.
type Config struct {
	// Root is the directory that contains view templates (default: "views").
	Root string
	// Extension is the file extension for templates (default: ".html").
	Extension string
	// ComponentsDir is the sub-directory for components relative to Root (default: "components").
	ComponentsDir string
	// LayoutsDir is the sub-directory for layouts relative to Root (default: "layouts").
	LayoutsDir string
	// PublicDir is the directory for static assets (default: "public").
	PublicDir string
	// PublicURL is the URL prefix for static assets (default: "/public").
	PublicURL string
	// ViteManifest is the path to the Vite/esbuild manifest.json file (optional).
	// When set, {{asset "app.js"}} will resolve hashed file names from the manifest.
	ViteManifest string
	// DevMode disables template caching so changes are picked up without restart.
	DevMode bool
	// Reload is an alias for DevMode kept for clarity.
	Reload bool
	// FuncMap adds extra template functions on top of the built-in helpers.
	FuncMap template.FuncMap
}

// Engine is the main view engine.
type Engine struct {
	cfg      Config
	mu       sync.RWMutex
	cache    map[string]*template.Template
	manifest map[string]viteEntry // Vite/esbuild manifest cache
	mtime    time.Time            // last manifest load time

	// fileHashCache caches MD5-based version hashes for files not in the manifest.
	fileHashCache map[string]string
}

type viteEntry struct {
	File    string   `json:"file"`
	CSS     []string `json:"css"`
	Imports []string `json:"imports"`
}

// New returns an initialised Engine with sensible defaults applied.
func New(cfg Config) *Engine {
	if cfg.Root == "" {
		cfg.Root = "views"
	}
	if cfg.Extension == "" {
		cfg.Extension = ".html"
	}
	if cfg.ComponentsDir == "" {
		cfg.ComponentsDir = "components"
	}
	if cfg.LayoutsDir == "" {
		cfg.LayoutsDir = "layouts"
	}
	if cfg.PublicDir == "" {
		cfg.PublicDir = "public"
	}
	if cfg.PublicURL == "" {
		cfg.PublicURL = "/public"
	}
	e := &Engine{
		cfg:           cfg,
		cache:         make(map[string]*template.Template),
		fileHashCache: make(map[string]string),
	}
	e.loadManifest()
	return e
}

// ----------------------------------------------------------------------------
// Public API
// ----------------------------------------------------------------------------

// Render executes template name with data into w. The template name is relative
// to the Root directory (e.g. "home/index" resolves to "views/home/index.html").
//
// Layout inheritance: if the template starts with {{extends "layouts/app"}} the
// engine renders the layout with the child's blocks injected.
func (e *Engine) Render(w io.Writer, name string, data any, layouts ...string) error {
	t, err := e.load(name)
	if err != nil {
		return err
	}
	// Execute into a buffer so we can inspect for extends directives.
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	if w != nil {
		_, err = w.Write(buf.Bytes())
	}
	return err
}

// RenderBytes executes a template and returns the output as []byte.
func (e *Engine) RenderBytes(name string, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.Render(&buf, name, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// RenderString executes a template and returns the output as a string.
func (e *Engine) RenderString(name string, data any) (string, error) {
	b, err := e.RenderBytes(name, data)
	return string(b), err
}

// ----------------------------------------------------------------------------
// Template loading & caching
// ----------------------------------------------------------------------------

// load returns a parsed *template.Template for the given view name, using the
// cache when DevMode is false.
func (e *Engine) load(name string) (*template.Template, error) {
	if !e.cfg.DevMode && !e.cfg.Reload {
		e.mu.RLock()
		if t, ok := e.cache[name]; ok {
			e.mu.RUnlock()
			return t, nil
		}
		e.mu.RUnlock()
	}

	t, err := e.parse(name)
	if err != nil {
		return nil, err
	}

	if !e.cfg.DevMode && !e.cfg.Reload {
		e.mu.Lock()
		e.cache[name] = t
		e.mu.Unlock()
	}
	return t, nil
}

// parse reads the template file, detects layout directives, and builds the
// final template tree including the layout + all components.
func (e *Engine) parse(name string) (*template.Template, error) {
	viewPath := e.resolvePath(name)
	src, err := os.ReadFile(viewPath)
	if err != nil {
		return nil, fmt.Errorf("view %q: %w", name, err)
	}

	// Detect {{extends "..."}} directive on the first non-blank line.
	layoutName, body := extractExtends(string(src))

	fm := e.buildFuncMap()

	if layoutName == "" {
		// Standalone template: load all components, then parse.
		return e.parseWithComponents(name, body, fm)
	}

	// Load layout source.
	layoutPath := e.resolvePath(layoutName)
	layoutSrc, err := os.ReadFile(layoutPath)
	if err != nil {
		return nil, fmt.Errorf("layout %q: %w", layoutName, err)
	}

	// Parse child to collect {{block}} definitions.
	childBlocks := parseBlocks(body)

	// Merge child blocks into layout source.
	merged := injectBlocks(string(layoutSrc), childBlocks)

	return e.parseWithComponents(layoutName+"::"+name, merged, fm)
}

// parseWithComponents loads all component files and parses them together with
// the main source so that {{template "components/alert" .}} works in views.
func (e *Engine) parseWithComponents(name, src string, fm template.FuncMap) (*template.Template, error) {
	t := template.New(name).Funcs(fm)

	// Parse the main source first.
	if _, err := t.Parse(src); err != nil {
		return nil, fmt.Errorf("parse %q: %w", name, err)
	}

	// Walk components directory and register each as a named template.
	compDir := filepath.Join(e.cfg.Root, e.cfg.ComponentsDir)
	if fi, err := os.Stat(compDir); err == nil && fi.IsDir() {
		err = filepath.WalkDir(compDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if !strings.HasSuffix(d.Name(), e.cfg.Extension) {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			// Template name is relative to Root, e.g. "components/alert".
			rel, _ := filepath.Rel(e.cfg.Root, path)
			rel = strings.TrimSuffix(filepath.ToSlash(rel), e.cfg.Extension)
			if _, parseErr := t.New(rel).Parse(string(data)); parseErr != nil {
				return fmt.Errorf("component %q: %w", rel, parseErr)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

// resolvePath returns the absolute filesystem path for a view/layout name.
func (e *Engine) resolvePath(name string) string {
	name = filepath.ToSlash(name)
	if !strings.HasSuffix(name, e.cfg.Extension) {
		name += e.cfg.Extension
	}
	return filepath.Join(e.cfg.Root, name)
}

// ----------------------------------------------------------------------------
// Extends / Block parsing helpers
// ----------------------------------------------------------------------------

// extractExtends detects the first {{extends "..."}} directive and returns the
// layout name and the remaining template body (with the directive stripped).
func extractExtends(src string) (layout, body string) {
	for _, line := range strings.SplitN(src, "\n", 5) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `{{extends "`) || strings.HasPrefix(trimmed, "{{extends `") {
			// Extract quoted argument.
			inner := strings.TrimPrefix(trimmed, "{{extends ")
			inner = strings.TrimSuffix(inner, "}}")
			inner = strings.Trim(strings.TrimSpace(inner), `"`+"`")
			layout = inner
			body = strings.Replace(src, line, "", 1)
			return
		}
		if trimmed != "" {
			break // first non-blank line is not an extends directive
		}
	}
	return "", src
}

// parseBlocks scans template source for {{block "name"}}...{{end}} sections
// and returns a map of block name → block content (the inner body).
func parseBlocks(src string) map[string]string {
	blocks := make(map[string]string)
	i := 0
	for i < len(src) {
		start := strings.Index(src[i:], `{{block "`)
		if start == -1 {
			break
		}
		start += i
		// Find block name.
		nameStart := start + len(`{{block "`)
		nameEnd := strings.Index(src[nameStart:], `"`)
		if nameEnd == -1 {
			break
		}
		blockName := src[nameStart : nameStart+nameEnd]
		// Find the body after the closing }}.
		bodyStart := strings.Index(src[nameStart+nameEnd:], "}}") + nameStart + nameEnd + 2
		// Find matching {{end}}.
		end := strings.Index(src[bodyStart:], "{{end}}")
		if end == -1 {
			break
		}
		bodyContent := src[bodyStart : bodyStart+end]
		blocks[blockName] = bodyContent
		i = bodyStart + end + len("{{end}}")
	}
	return blocks
}

// injectBlocks replaces {{block "name"}}default{{end}} occurrences in layoutSrc
// with the child's block content (if defined) or keeps the default.
func injectBlocks(layoutSrc string, childBlocks map[string]string) string {
	result := layoutSrc
	for name, content := range childBlocks {
		placeholder := `{{block "` + name + `"`
		start := strings.Index(result, placeholder)
		for start != -1 {
			// Find closing tag + {{end}}.
			bodyStart := strings.Index(result[start:], "}}") + start + 2
			end := strings.Index(result[bodyStart:], "{{end}}")
			if end == -1 {
				break
			}
			endPos := bodyStart + end + len("{{end}}")
			result = result[:start] + content + result[endPos:]
			start = strings.Index(result, placeholder)
		}
	}
	return result
}

// ----------------------------------------------------------------------------
// Template FuncMap
// ----------------------------------------------------------------------------

// buildFuncMap returns the complete template.FuncMap for the engine, combining
// built-in helpers with any user-supplied functions from Config.FuncMap.
func (e *Engine) buildFuncMap() template.FuncMap {
	fm := template.FuncMap{
		// Asset helpers
		"asset":    e.assetHelper,
		"assetCSS": e.assetCSSHelper,
		"assetJS":  e.assetJSHelper,

		// HTML escape safety
		"safe":     func(s string) template.HTML { return template.HTML(s) },
		"safeURL":  func(s string) template.URL { return template.URL(s) },
		"safeAttr": func(s string) template.HTMLAttr { return template.HTMLAttr(s) },

		// Form helpers — these are no-ops at the FuncMap level; they are
		// intended to be called from templates and rely on .CSRF being passed
		// through the data map (see ViewData helpers in flash package).
		"form_open":  formOpen,
		"form_close": formClose,
		"input":      inputHelper,
		"textarea":   textareaHelper,
		"select": func(name, selected string, options any, attrs ...string) template.HTML {
			return selectHelper(name, selected, options, attrs...)
		},
		"checkbox": checkboxHelper,

		// csrf_field emits a hidden <input> with the CSRF token.
		// Usage: {{csrf_field .CSRF}}
		"csrf_field": func(token string) template.HTML {
			if token == "" {
				return ""
			}
			return template.HTML(fmt.Sprintf(
				`<input type="hidden" name="_csrf" value="%s">`,
				template.HTMLEscapeString(token),
			))
		},

		// String utilities
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title, //nolint:staticcheck
		"trim":  strings.TrimSpace,

		// Old-input helper: {{old "email" .Old}}
		"old": func(field string, old map[string]string) string {
			if old == nil {
				return ""
			}
			return old[field]
		},

		// Misc
		"nl2br": func(s string) template.HTML {
			return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(s), "\n", "<br>"))
		},
		"default": func(def, val any) any {
			if val == nil || val == "" || val == 0 || val == false {
				return def
			}
			return val
		},
		"json": func(v any) (string, error) {
			b, err := json.Marshal(v)
			return string(b), err
		},
		"iterate": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},

		// dict creates a map[string]any from key-value pairs.
		// Usage in templates: {{template "components/form-input" (dict "Name" "email" "Label" "Email")}}
		"dict": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, fmt.Errorf("dict requires an even number of arguments (key-value pairs)")
			}
			m := make(map[string]any, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				key, ok := pairs[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict key at position %d must be a string", i)
				}
				m[key] = pairs[i+1]
			}
			return m, nil
		},

		// concat joins strings without a separator.
		"concat": func(parts ...string) string { return strings.Join(parts, "") },

		// Arithmetic helpers used in templates (e.g. pagination).
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},

		// slice builds a []any from variadic arguments — useful for passing
		// option lists to form helpers without Go code.
		// Example: {{select "role" .Role (slice (arr "admin" "Admin") (arr "user" "User"))}}
		"slice": func(items ...any) []any { return items },

		// arr builds a [2]string from two string arguments for select options.
		// Example: {{arr "value" "Label"}}
		"arr": func(value, label string) [2]string { return [2]string{value, label} },

		// hasKey checks whether a map[string]string or map[string]any has a key.
		"hasKey": func(m any, key string) bool {
			switch v := m.(type) {
			case map[string]string:
				_, ok := v[key]
				return ok
			case map[string]any:
				_, ok := v[key]
				return ok
			}
			return false
		},
		"eq":  func(a, b any) bool { return fmt.Sprint(a) == fmt.Sprint(b) },
		"neq": func(a, b any) bool { return fmt.Sprint(a) != fmt.Sprint(b) },
	}

	// Merge user-supplied functions (they can override built-ins).
	for k, v := range e.cfg.FuncMap {
		fm[k] = v
	}
	return fm
}

// ----------------------------------------------------------------------------
// Asset helpers
// ----------------------------------------------------------------------------

// assetHelper returns a versioned URL for the given asset path.
// Resolution order:
//  1. Vite/esbuild manifest (if configured and entry exists)
//  2. MD5 hash of the file from PublicDir
//  3. Plain URL with no version query string
func (e *Engine) assetHelper(path string) string {
	path = strings.TrimPrefix(path, "/")

	// 1. Try manifest.
	if entry, ok := e.manifestEntry(path); ok {
		return e.cfg.PublicURL + "/" + entry.File
	}

	// 2. Try MD5 hash.
	fullPath := filepath.Join(e.cfg.PublicDir, path)
	if hash, ok := e.fileHash(fullPath); ok {
		return e.cfg.PublicURL + "/" + path + "?v=" + hash
	}

	// 3. Plain URL.
	return e.cfg.PublicURL + "/" + path
}

// assetCSSHelper emits a <link rel="stylesheet"> tag.
func (e *Engine) assetCSSHelper(path string) template.HTML {
	url := e.assetHelper(path)
	return template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, template.HTMLEscapeString(url)))
}

// assetJSHelper emits a <script src="..."> tag.
func (e *Engine) assetJSHelper(path string) template.HTML {
	url := e.assetHelper(path)
	return template.HTML(fmt.Sprintf(`<script src="%s" defer></script>`, template.HTMLEscapeString(url)))
}

// manifestEntry returns the Vite manifest entry for the given source path.
func (e *Engine) manifestEntry(path string) (viteEntry, bool) {
	if e.cfg.ViteManifest == "" {
		return viteEntry{}, false
	}
	e.loadManifest()
	e.mu.RLock()
	entry, ok := e.manifest[path]
	e.mu.RUnlock()
	return entry, ok
}

// loadManifest (re)loads the Vite/esbuild manifest.json when it has changed on disk.
func (e *Engine) loadManifest() {
	if e.cfg.ViteManifest == "" {
		return
	}
	fi, err := os.Stat(e.cfg.ViteManifest)
	if err != nil {
		return
	}

	e.mu.RLock()
	unchanged := !fi.ModTime().After(e.mtime)
	e.mu.RUnlock()
	if unchanged {
		return
	}

	data, err := os.ReadFile(e.cfg.ViteManifest)
	if err != nil {
		return
	}

	var m map[string]viteEntry
	if json.Unmarshal(data, &m) != nil {
		return
	}

	e.mu.Lock()
	e.manifest = m
	e.mtime = fi.ModTime()
	e.mu.Unlock()
}

// fileHash returns a short MD5 hash of a file for cache-busting.
func (e *Engine) fileHash(path string) (string, bool) {
	e.mu.RLock()
	if h, ok := e.fileHashCache[path]; ok {
		e.mu.RUnlock()
		return h, true
	}
	e.mu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	sum := md5.Sum(data) //nolint:gosec
	hash := fmt.Sprintf("%x", sum[:4])

	e.mu.Lock()
	e.fileHashCache[path] = hash
	e.mu.Unlock()
	return hash, true
}

// ClearCache discards all cached templates and file hashes. Call on deploy or
// when templates change in DevMode = false environments.
func (e *Engine) ClearCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cache = make(map[string]*template.Template)
	e.fileHashCache = make(map[string]string)
	e.manifest = nil
	e.mtime = time.Time{}
}

// ----------------------------------------------------------------------------
// Form builder helpers (standalone functions used in FuncMap)
// ----------------------------------------------------------------------------

// formOpen renders an HTML <form> opening tag. Pass csrf token via data.CSRF.
func formOpen(action, method string, attrs ...string) template.HTML {
	if method == "" {
		method = "POST"
	}
	extra := ""
	if len(attrs) > 0 {
		extra = " " + strings.Join(attrs, " ")
	}
	return template.HTML(fmt.Sprintf(`<form action="%s" method="%s"%s>`,
		template.HTMLEscapeString(action),
		template.HTMLEscapeString(strings.ToUpper(method)),
		extra,
	))
}

func formClose() template.HTML { return `</form>` }

// inputHelper renders an <input> element.
// Usage: {{input "text" "email" .FormData.Email "class=\"form-control\""}}
func inputHelper(typ, name, value string, attrs ...string) template.HTML {
	extra := ""
	if len(attrs) > 0 {
		extra = " " + strings.Join(attrs, " ")
	}
	return template.HTML(fmt.Sprintf(
		`<input type="%s" name="%s" value="%s"%s>`,
		template.HTMLEscapeString(typ),
		template.HTMLEscapeString(name),
		template.HTMLEscapeString(value),
		extra,
	))
}

// textareaHelper renders a <textarea> element.
func textareaHelper(name, value string, attrs ...string) template.HTML {
	extra := ""
	if len(attrs) > 0 {
		extra = " " + strings.Join(attrs, " ")
	}
	return template.HTML(fmt.Sprintf(
		`<textarea name="%s"%s>%s</textarea>`,
		template.HTMLEscapeString(name),
		extra,
		template.HTMLEscapeString(value),
	))
}

// selectHelper renders a <select> element.
// options accepts [][2]string or []any (from the {{slice}} helper) where each
// element is a [2]string{value, label} pair.
func selectHelper(name, selected string, options any, attrs ...string) template.HTML {
	extra := ""
	if len(attrs) > 0 {
		extra = " " + strings.Join(attrs, " ")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<select name="%s"%s>`, template.HTMLEscapeString(name), extra))

	renderOpt := func(value, label string) {
		sel := ""
		if value == selected {
			sel = ` selected`
		}
		sb.WriteString(fmt.Sprintf(
			`<option value="%s"%s>%s</option>`,
			template.HTMLEscapeString(value),
			sel,
			template.HTMLEscapeString(label),
		))
	}

	switch opts := options.(type) {
	case [][2]string:
		for _, opt := range opts {
			renderOpt(opt[0], opt[1])
		}
	case []any:
		for _, item := range opts {
			switch pair := item.(type) {
			case [2]string:
				renderOpt(pair[0], pair[1])
			case []string:
				if len(pair) >= 2 {
					renderOpt(pair[0], pair[1])
				}
			}
		}
	}

	sb.WriteString(`</select>`)
	return template.HTML(sb.String())
}

// checkboxHelper renders an <input type="checkbox"> element.
func checkboxHelper(name, value string, checked bool, attrs ...string) template.HTML {
	extra := ""
	if len(attrs) > 0 {
		extra = " " + strings.Join(attrs, " ")
	}
	chk := ""
	if checked {
		chk = " checked"
	}
	return template.HTML(fmt.Sprintf(
		`<input type="checkbox" name="%s" value="%s"%s%s>`,
		template.HTMLEscapeString(name),
		template.HTMLEscapeString(value),
		chk,
		extra,
	))
}
