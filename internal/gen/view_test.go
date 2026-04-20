package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zatrano/zatrano/internal/gen"
)

// ----------------------------------------------------------------------------
// gen.View — dry-run
// ----------------------------------------------------------------------------

func TestView_DryRun_DefaultFiles(t *testing.T) {
	root := t.TempDir()
	paths, err := gen.View(root, "post", gen.ViewOptions{DryRun: true})
	if err != nil {
		t.Fatalf("View dry-run: %v", err)
	}
	// Default: index + show only.
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d: %v", len(paths), paths)
	}
	hasIndex := false
	hasShow := false
	for _, p := range paths {
		base := filepath.Base(p)
		if base == "index.html" {
			hasIndex = true
		}
		if base == "show.html" {
			hasShow = true
		}
	}
	if !hasIndex {
		t.Error("expected index.html in paths")
	}
	if !hasShow {
		t.Error("expected show.html in paths")
	}
	// Dry-run: no files should be written.
	if _, err := os.Stat(filepath.Join(root, "post", "index.html")); err == nil {
		t.Error("dry-run should not write files")
	}
}

func TestView_DryRun_WithForm(t *testing.T) {
	root := t.TempDir()
	paths, err := gen.View(root, "article", gen.ViewOptions{
		WithForm: true,
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("View dry-run with-form: %v", err)
	}
	// Should have index, show, create, edit.
	if len(paths) != 4 {
		t.Errorf("expected 4 paths, got %d: %v", len(paths), paths)
	}
	bases := make(map[string]bool)
	for _, p := range paths {
		bases[filepath.Base(p)] = true
	}
	for _, want := range []string{"index.html", "show.html", "create.html", "edit.html"} {
		if !bases[want] {
			t.Errorf("expected %s in paths", want)
		}
	}
}

// ----------------------------------------------------------------------------
// gen.View — actual file writing
// ----------------------------------------------------------------------------

func TestView_WritesFiles(t *testing.T) {
	root := t.TempDir()
	_, err := gen.View(root, "user", gen.ViewOptions{
		Layout:   "layouts/app",
		WithForm: true,
		DryRun:   false,
	})
	if err != nil {
		t.Fatalf("View: %v", err)
	}

	for _, name := range []string{"index.html", "show.html", "create.html", "edit.html"} {
		path := filepath.Join(root, "user", name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", path)
		}
	}
}

func TestView_IndexContent(t *testing.T) {
	root := t.TempDir()
	_, err := gen.View(root, "product", gen.ViewOptions{
		Layout: "layouts/app",
		DryRun: false,
	})
	if err != nil {
		t.Fatalf("View: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "product", "index.html"))
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `{{extends "layouts/app"}}`) {
		t.Errorf("expected extends directive, got:\n%s", content)
	}
	if !strings.Contains(content, `{{block "content"}}`) {
		t.Errorf("expected content block, got:\n%s", content)
	}
	// Should reference the module name in links.
	if !strings.Contains(content, "/product") {
		t.Errorf("expected module name in links, got:\n%s", content)
	}
}

func TestView_CreateContent_HasCSRF(t *testing.T) {
	root := t.TempDir()
	_, err := gen.View(root, "order", gen.ViewOptions{
		Layout:   "layouts/app",
		WithForm: true,
		DryRun:   false,
	})
	if err != nil {
		t.Fatalf("View: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "order", "create.html"))
	if err != nil {
		t.Fatalf("read create.html: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "csrf_field") {
		t.Errorf("expected csrf_field in create form, got:\n%s", content)
	}
	if !strings.Contains(content, "form_open") {
		t.Errorf("expected form_open in create form, got:\n%s", content)
	}
	if !strings.Contains(content, "form_close") {
		t.Errorf("expected form_close in create form, got:\n%s", content)
	}
}

func TestView_EditContent_HasMethodOverride(t *testing.T) {
	root := t.TempDir()
	_, err := gen.View(root, "ticket", gen.ViewOptions{
		Layout:   "layouts/app",
		WithForm: true,
		DryRun:   false,
	})
	if err != nil {
		t.Fatalf("View: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "ticket", "edit.html"))
	if err != nil {
		t.Fatalf("read edit.html: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `_method`) {
		t.Errorf("expected HTTP method override in edit form, got:\n%s", content)
	}
	if !strings.Contains(content, `old "name"`) {
		t.Errorf("expected old input helper in edit form, got:\n%s", content)
	}
}

// ----------------------------------------------------------------------------
// gen.View — custom layout
// ----------------------------------------------------------------------------

func TestView_CustomLayout(t *testing.T) {
	root := t.TempDir()
	_, err := gen.View(root, "report", gen.ViewOptions{
		Layout: "layouts/admin",
		DryRun: false,
	})
	if err != nil {
		t.Fatalf("View: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(root, "report", "index.html"))
	if !strings.Contains(string(data), `{{extends "layouts/admin"}}`) {
		t.Errorf("expected custom layout 'layouts/admin', got:\n%s", string(data))
	}
}

// ----------------------------------------------------------------------------
// gen.View — invalid name
// ----------------------------------------------------------------------------

func TestView_InvalidName(t *testing.T) {
	root := t.TempDir()
	_, err := gen.View(root, "   ", gen.ViewOptions{})
	if err == nil {
		t.Error("expected error for blank name, got nil")
	}
}

func TestView_NameNormalization(t *testing.T) {
	root := t.TempDir()
	// Dashes should be converted to underscores.
	paths, err := gen.View(root, "blog-post", gen.ViewOptions{DryRun: true})
	if err != nil {
		t.Fatalf("View: %v", err)
	}
	for _, p := range paths {
		if strings.Contains(filepath.ToSlash(p), "blog-post") {
			t.Errorf("expected dashes normalized to underscores, got: %s", p)
		}
	}
	_ = paths
}
