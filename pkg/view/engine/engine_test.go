package engine_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zatrano/zatrano/pkg/view/engine"
)

// setupViews creates a temporary views directory tree for testing and returns
// the root path and a cleanup function.
func setupViews(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	return root
}

func newEngine(root string, devMode bool) *engine.Engine {
	return engine.New(engine.Config{
		Root:          root,
		Extension:     ".html",
		ComponentsDir: "components",
		LayoutsDir:    "layouts",
		DevMode:       devMode,
	})
}

// ----------------------------------------------------------------------------
// Standalone template rendering
// ----------------------------------------------------------------------------

func TestRender_Standalone(t *testing.T) {
	root := setupViews(t, map[string]string{
		"home/index.html": `<h1>Hello {{.Name}}</h1>`,
	})
	e := newEngine(root, true)

	var buf bytes.Buffer
	if err := e.Render(&buf, "home/index", map[string]any{"Name": "World"}); err != nil {
		t.Fatalf("Render: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "Hello World") {
		t.Errorf("expected 'Hello World' in output, got: %s", got)
	}
}

// ----------------------------------------------------------------------------
// Layout inheritance
// ----------------------------------------------------------------------------

func TestRender_LayoutInheritance(t *testing.T) {
	root := setupViews(t, map[string]string{
		"layouts/app.html": `<html><head><title>{{block "title" .}}Default{{end}}</title></head><body>{{block "content" .}}{{end}}</body></html>`,
		"home/index.html": `{{extends "layouts/app"}}
{{block "title"}}Custom Title{{end}}
{{block "content"}}<p>Hello {{.Name}}</p>{{end}}`,
	})
	e := newEngine(root, true)

	out, err := e.RenderString("home/index", map[string]any{"Name": "Zatrano"})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, "Custom Title") {
		t.Errorf("expected 'Custom Title', got: %s", out)
	}
	if !strings.Contains(out, "Hello Zatrano") {
		t.Errorf("expected 'Hello Zatrano', got: %s", out)
	}
	if strings.Contains(out, "Default") {
		t.Errorf("default block content should be replaced, got: %s", out)
	}
}

func TestRender_LayoutDefaultBlock(t *testing.T) {
	root := setupViews(t, map[string]string{
		"layouts/app.html": `<title>{{block "title" .}}Default Title{{end}}</title><body>{{block "content" .}}{{end}}</body>`,
		"home/index.html": `{{extends "layouts/app"}}
{{block "content"}}<p>Content only</p>{{end}}`,
	})
	e := newEngine(root, true)

	out, err := e.RenderString("home/index", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, "Default Title") {
		t.Errorf("expected default title block to remain, got: %s", out)
	}
	if !strings.Contains(out, "Content only") {
		t.Errorf("expected 'Content only', got: %s", out)
	}
}

// ----------------------------------------------------------------------------
// Component loading
// ----------------------------------------------------------------------------

func TestRender_Component(t *testing.T) {
	root := setupViews(t, map[string]string{
		"components/alert.html": `{{define "components/alert"}}<div class="alert-{{.Type}}">{{.Message}}</div>{{end}}`,
		"home/index.html":       `{{template "components/alert" (dict "Type" "success" "Message" "Done!")}}`,
	})
	e := newEngine(root, true)

	out, err := e.RenderString("home/index", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `alert-success`) {
		t.Errorf("expected 'alert-success', got: %s", out)
	}
	if !strings.Contains(out, "Done!") {
		t.Errorf("expected 'Done!', got: %s", out)
	}
}

func TestRender_ComponentInLayout(t *testing.T) {
	root := setupViews(t, map[string]string{
		"components/badge.html": `{{define "components/badge"}}<span class="badge">{{.Label}}</span>{{end}}`,
		"layouts/app.html":      `<body>{{block "content" .}}{{end}}</body>`,
		"home/index.html": `{{extends "layouts/app"}}
{{block "content"}}{{template "components/badge" (dict "Label" "New")}}{{end}}`,
	})
	e := newEngine(root, true)

	out, err := e.RenderString("home/index", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `<span class="badge">New</span>`) {
		t.Errorf("expected badge component in output, got: %s", out)
	}
}

// ----------------------------------------------------------------------------
// Template FuncMap helpers
// ----------------------------------------------------------------------------

func TestFuncMap_Dict(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{$m := dict "key" "value"}}{{index $m "key"}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "value" {
		t.Errorf("expected 'value', got: %q", out)
	}
}

func TestFuncMap_OldHelper(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{old "email" .Old}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{
		"Old": map[string]string{"email": "user@example.com"},
	})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "user@example.com" {
		t.Errorf("expected 'user@example.com', got: %q", out)
	}
}

func TestFuncMap_OldHelper_Missing(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `"{{old "missing" .Old}}"`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{
		"Old": map[string]string{},
	})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `""`) {
		t.Errorf("expected empty old value, got: %q", out)
	}
}

func TestFuncMap_CsrfField(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{csrf_field .CSRF}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{"CSRF": "tok123"})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `name="_csrf"`) {
		t.Errorf("expected hidden csrf input, got: %s", out)
	}
	if !strings.Contains(out, `value="tok123"`) {
		t.Errorf("expected csrf value, got: %s", out)
	}
}

func TestFuncMap_CsrfField_Empty(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `"{{csrf_field ""}}"`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	// Empty token should produce empty output.
	if strings.Contains(out, `name="_csrf"`) {
		t.Errorf("expected no csrf input for empty token, got: %s", out)
	}
}

func TestFuncMap_Safe(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{safe .HTML}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{"HTML": "<b>bold</b>"})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, "<b>bold</b>") {
		t.Errorf("expected raw HTML, got: %s", out)
	}
}

func TestFuncMap_Upper_Lower(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{upper .A}}-{{lower .B}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{"A": "hello", "B": "WORLD"})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "HELLO-world" {
		t.Errorf("expected 'HELLO-world', got: %q", out)
	}
}

func TestFuncMap_Default(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{default "fallback" .Val}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{"Val": ""})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "fallback" {
		t.Errorf("expected 'fallback', got: %q", out)
	}
}

func TestFuncMap_Nl2br(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{nl2br .Text}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{"Text": "line1\nline2"})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, "<br>") {
		t.Errorf("expected <br> in output, got: %s", out)
	}
}

func TestFuncMap_Concat(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{concat "hello" " " "world"}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "hello world" {
		t.Errorf("expected 'hello world', got: %q", out)
	}
}

func TestFuncMap_Arithmetic(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		expected string
	}{
		{"add", `{{add 3 4}}`, "7"},
		{"sub", `{{sub 10 3}}`, "7"},
		{"mul", `{{mul 3 4}}`, "12"},
		{"div", `{{div 10 2}}`, "5"},
		{"mod", `{{mod 10 3}}`, "1"},
		{"div_zero", `{{div 5 0}}`, "0"},
		{"mod_zero", `{{mod 5 0}}`, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := setupViews(t, map[string]string{"test.html": tt.tmpl})
			e := newEngine(root, true)
			out, err := e.RenderString("test", nil)
			if err != nil {
				t.Fatalf("RenderString: %v", err)
			}
			if strings.TrimSpace(out) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strings.TrimSpace(out))
			}
		})
	}
}

func TestFuncMap_HasKey(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{if hasKey .M "exists"}}yes{{else}}no{{end}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", map[string]any{
		"M": map[string]string{"exists": "1"},
	})
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "yes" {
		t.Errorf("expected 'yes', got: %q", out)
	}
}

func TestFuncMap_Iterate(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{range iterate 3}}{{.}}{{end}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "012" {
		t.Errorf("expected '012', got: %q", out)
	}
}

// ----------------------------------------------------------------------------
// Form builder helpers
// ----------------------------------------------------------------------------

func TestFuncMap_FormOpen(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{form_open "/users" "POST"}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `action="/users"`) {
		t.Errorf("expected action, got: %s", out)
	}
	if !strings.Contains(out, `method="POST"`) {
		t.Errorf("expected method, got: %s", out)
	}
}

func TestFuncMap_FormClose(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{form_close}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `</form>`) {
		t.Errorf("expected </form>, got: %s", out)
	}
}

func TestFuncMap_Input(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{input "email" "email" "user@example.com"}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `type="email"`) {
		t.Errorf("expected type=email, got: %s", out)
	}
	if !strings.Contains(out, `value="user@example.com"`) {
		t.Errorf("expected value, got: %s", out)
	}
}

func TestFuncMap_Textarea(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{textarea "bio" "hello"}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `<textarea`) {
		t.Errorf("expected textarea, got: %s", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected value, got: %s", out)
	}
}

func TestFuncMap_Checkbox(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{checkbox "active" "1" true}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `type="checkbox"`) {
		t.Errorf("expected checkbox type, got: %s", out)
	}
	if !strings.Contains(out, "checked") {
		t.Errorf("expected checked attr, got: %s", out)
	}
}

// ----------------------------------------------------------------------------
// Template caching
// ----------------------------------------------------------------------------

func TestCache_HitInNonDevMode(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `version1`,
	})
	e := newEngine(root, false) // caching enabled

	out1, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("first render: %v", err)
	}

	// Overwrite the file — cache should still serve old content.
	if err := os.WriteFile(filepath.Join(root, "test.html"), []byte("version2"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	out2, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("second render: %v", err)
	}

	if out1 != out2 {
		t.Errorf("cache miss in non-dev mode: first=%q second=%q", out1, out2)
	}
}

func TestClearCache(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `version1`,
	})
	e := newEngine(root, false)

	_, _ = e.RenderString("test", nil)

	// Overwrite + clear cache.
	if err := os.WriteFile(filepath.Join(root, "test.html"), []byte("version2"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	e.ClearCache()

	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("render after clear: %v", err)
	}
	if !strings.Contains(out, "version2") {
		t.Errorf("expected version2 after cache clear, got: %q", out)
	}
}

// ----------------------------------------------------------------------------
// Error cases
// ----------------------------------------------------------------------------

func TestRender_MissingTemplate(t *testing.T) {
	root := setupViews(t, map[string]string{})
	e := newEngine(root, true)

	err := e.Render(nil, "nonexistent/view", nil)
	if err == nil {
		t.Error("expected error for missing template, got nil")
	}
}

func TestRender_MissingLayout(t *testing.T) {
	root := setupViews(t, map[string]string{
		"child.html": `{{extends "layouts/missing"}}{{block "content"}}hi{{end}}`,
	})
	e := newEngine(root, true)

	err := e.Render(nil, "child", nil)
	if err == nil {
		t.Error("expected error for missing layout, got nil")
	}
}

func TestFuncMap_Dict_OddArgs(t *testing.T) {
	root := setupViews(t, map[string]string{
		// dict with odd number of args should produce a template error
		"test.html": `{{dict "a"}}`,
	})
	e := newEngine(root, true)
	_, err := e.RenderString("test", nil)
	if err == nil {
		t.Error("expected error for dict with odd args, got nil")
	}
}

// ----------------------------------------------------------------------------
// RenderBytes
// ----------------------------------------------------------------------------

func TestRenderBytes(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `bytes test`,
	})
	e := newEngine(root, true)
	b, err := e.RenderBytes("test", nil)
	if err != nil {
		t.Fatalf("RenderBytes: %v", err)
	}
	if string(b) != "bytes test" {
		t.Errorf("expected 'bytes test', got: %q", string(b))
	}
}

func TestFuncMap_Slice(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{$s := slice "a" "b" "c"}}{{index $s 1}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "b" {
		t.Errorf("expected 'b', got: %q", out)
	}
}

func TestFuncMap_Arr(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{$a := arr "val" "Label"}}{{index $a 0}}-{{index $a 1}}`,
	})
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if strings.TrimSpace(out) != "val-Label" {
		t.Errorf("expected 'val-Label', got: %q", out)
	}
}

func TestFuncMap_Select(t *testing.T) {
	root := setupViews(t, map[string]string{
		"test.html": `{{select "role" "admin" (slice (arr "admin" "Admin") (arr "user" "User"))}}`,
	})
	// select takes [][2]string but the template uses []any with [2]string elements.
	// The helper casts internally — just verify the output structure.
	e := newEngine(root, true)
	out, err := e.RenderString("test", nil)
	if err != nil {
		t.Fatalf("RenderString: %v", err)
	}
	if !strings.Contains(out, `<select`) {
		t.Errorf("expected select element, got: %s", out)
	}
	if !strings.Contains(out, `name="role"`) {
		t.Errorf("expected name=role, got: %s", out)
	}
}
