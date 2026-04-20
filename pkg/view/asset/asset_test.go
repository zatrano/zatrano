package asset_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zatrano/zatrano/pkg/view/asset"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}

// ----------------------------------------------------------------------------
// Plain URL fallback
// ----------------------------------------------------------------------------

func TestURL_Plain(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir: t.TempDir(),
		PublicURL: "/public",
	})
	url := m.URL("css/app.css")
	// No manifest, no file on disk → plain URL
	if url != "/public/css/app.css" {
		t.Errorf("expected plain URL, got: %s", url)
	}
}

func TestURL_LeadingSlash(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir: t.TempDir(),
		PublicURL: "/public",
	})
	// Leading slash should be trimmed.
	url := m.URL("/css/app.css")
	if url != "/public/css/app.css" {
		t.Errorf("expected /public/css/app.css, got: %s", url)
	}
}

// ----------------------------------------------------------------------------
// MD5 hash fallback
// ----------------------------------------------------------------------------

func TestURL_FileHash(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "css/app.css", "body{color:red}")

	m := asset.New(asset.Config{
		PublicDir: dir,
		PublicURL: "/public",
	})
	url := m.URL("css/app.css")
	if !strings.HasPrefix(url, "/public/css/app.css?v=") {
		t.Errorf("expected hash query string, got: %s", url)
	}
	// Hash should be 8 hex chars (4 bytes).
	parts := strings.SplitN(url, "?v=", 2)
	if len(parts) != 2 || len(parts[1]) != 8 {
		t.Errorf("expected 8-char hash, got: %q (full URL: %s)", parts, url)
	}
}

func TestURL_FileHash_Cached(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "js/app.js", "console.log(1)")

	m := asset.New(asset.Config{
		PublicDir: dir,
		PublicURL: "/pub",
	})
	url1 := m.URL("js/app.js")
	url2 := m.URL("js/app.js") // should hit cache
	if url1 != url2 {
		t.Errorf("hash should be stable across calls: %s vs %s", url1, url2)
	}
}

// ----------------------------------------------------------------------------
// Vite manifest
// ----------------------------------------------------------------------------

func writeManifest(t *testing.T, dir string, entries map[string]asset.ViteEntry) string {
	t.Helper()
	data, _ := json.Marshal(entries)
	path := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return path
}

func TestURL_ViteManifest(t *testing.T) {
	dir := t.TempDir()
	manifest := writeManifest(t, dir, map[string]asset.ViteEntry{
		"src/main.ts": {File: "assets/main-a1b2c3d4.js"},
	})

	m := asset.New(asset.Config{
		PublicDir:    dir,
		PublicURL:    "/build",
		ViteManifest: manifest,
	})
	url := m.URL("src/main.ts")
	if url != "/build/assets/main-a1b2c3d4.js" {
		t.Errorf("expected hashed URL from manifest, got: %s", url)
	}
}

func TestURL_ViteManifest_MissingEntry_FallsBack(t *testing.T) {
	dir := t.TempDir()
	manifest := writeManifest(t, dir, map[string]asset.ViteEntry{
		"src/other.ts": {File: "assets/other-xyz.js"},
	})

	m := asset.New(asset.Config{
		PublicDir:    dir,
		PublicURL:    "/build",
		ViteManifest: manifest,
	})
	// Entry not in manifest → plain URL
	url := m.URL("src/main.ts")
	if url != "/build/src/main.ts" {
		t.Errorf("expected plain fallback URL, got: %s", url)
	}
}

func TestManifestEntry(t *testing.T) {
	dir := t.TempDir()
	manifest := writeManifest(t, dir, map[string]asset.ViteEntry{
		"src/app.css": {File: "assets/app-abc.css", CSS: []string{"assets/app-abc.css"}},
	})

	m := asset.New(asset.Config{
		PublicDir:    dir,
		PublicURL:    "/build",
		ViteManifest: manifest,
	})
	entry, ok := m.ManifestEntry("src/app.css")
	if !ok {
		t.Fatal("expected manifest entry to be found")
	}
	if entry.File != "assets/app-abc.css" {
		t.Errorf("unexpected file: %s", entry.File)
	}
}

// ----------------------------------------------------------------------------
// Vite dev mode (HMR)
// ----------------------------------------------------------------------------

func TestURL_DevMode_ViteDevURL(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir:  t.TempDir(),
		PublicURL:  "/build",
		ViteDevURL: "http://localhost:5173",
		DevMode:    true,
	})
	url := m.URL("src/main.ts")
	if url != "http://localhost:5173/src/main.ts" {
		t.Errorf("expected Vite dev URL, got: %s", url)
	}
}

func TestURL_DevMode_NoViteDevURL_UsesManifest(t *testing.T) {
	dir := t.TempDir()
	manifest := writeManifest(t, dir, map[string]asset.ViteEntry{
		"src/main.ts": {File: "assets/main-prod.js"},
	})
	m := asset.New(asset.Config{
		PublicDir:    dir,
		PublicURL:    "/build",
		ViteManifest: manifest,
		DevMode:      true, // dev mode but no ViteDevURL → fall through to manifest
	})
	url := m.URL("src/main.ts")
	if url != "/build/assets/main-prod.js" {
		t.Errorf("expected manifest URL when no ViteDevURL, got: %s", url)
	}
}

// ----------------------------------------------------------------------------
// HTML tag helpers
// ----------------------------------------------------------------------------

func TestLinkTag(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir: t.TempDir(),
		PublicURL: "/public",
	})
	tag := string(m.LinkTag("css/app.css"))
	if !strings.Contains(tag, `rel="stylesheet"`) {
		t.Errorf("expected rel=stylesheet, got: %s", tag)
	}
	if !strings.Contains(tag, "/public/css/app.css") {
		t.Errorf("expected href, got: %s", tag)
	}
}

func TestScriptTag(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir: t.TempDir(),
		PublicURL: "/public",
	})
	tag := string(m.ScriptTag("js/app.js"))
	if !strings.Contains(tag, `<script`) {
		t.Errorf("expected script tag, got: %s", tag)
	}
	if !strings.Contains(tag, "defer") {
		t.Errorf("expected defer attribute, got: %s", tag)
	}
}

func TestViteHead_DevMode(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir:  t.TempDir(),
		PublicURL:  "/build",
		ViteDevURL: "http://localhost:5173",
		DevMode:    true,
	})
	head := string(m.ViteHead("src/main.ts"))
	if !strings.Contains(head, "@vite/client") {
		t.Errorf("expected vite client script, got: %s", head)
	}
	if !strings.Contains(head, "src/main.ts") {
		t.Errorf("expected entry script, got: %s", head)
	}
	if !strings.Contains(head, `type="module"`) {
		t.Errorf("expected type=module, got: %s", head)
	}
}

func TestViteHead_Prod(t *testing.T) {
	dir := t.TempDir()
	manifest := writeManifest(t, dir, map[string]asset.ViteEntry{
		"src/main.ts": {
			File: "assets/main-xyz.js",
			CSS:  []string{"assets/main-xyz.css"},
		},
	})
	m := asset.New(asset.Config{
		PublicDir:    dir,
		PublicURL:    "/build",
		ViteManifest: manifest,
		DevMode:      false,
	})
	head := string(m.ViteHead("src/main.ts"))
	if !strings.Contains(head, "assets/main-xyz.js") {
		t.Errorf("expected hashed JS, got: %s", head)
	}
	if !strings.Contains(head, "assets/main-xyz.css") {
		t.Errorf("expected hashed CSS link, got: %s", head)
	}
	if strings.Contains(head, "@vite/client") {
		t.Errorf("should not include vite client in prod, got: %s", head)
	}
}

// ----------------------------------------------------------------------------
// TemplateFuncs
// ----------------------------------------------------------------------------

func TestTemplateFuncs_Keys(t *testing.T) {
	m := asset.New(asset.Config{
		PublicDir: t.TempDir(),
		PublicURL: "/public",
	})
	funcs := m.TemplateFuncs()
	required := []string{"asset", "assetCSS", "assetJS", "assetLink", "assetScript", "viteHead"}
	for _, key := range required {
		if _, ok := funcs[key]; !ok {
			t.Errorf("missing template func: %s", key)
		}
	}
}

// ----------------------------------------------------------------------------
// ClearCache
// ----------------------------------------------------------------------------

func TestClearCache(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "css/app.css", "v1")

	m := asset.New(asset.Config{
		PublicDir: dir,
		PublicURL: "/public",
	})

	url1 := m.URL("css/app.css")
	if !strings.Contains(url1, "?v=") {
		t.Fatalf("expected hash in first URL, got: %s", url1)
	}

	// Change file content and clear cache.
	writeTempFile(t, dir, "css/app.css", "v2 with different content")
	m.ClearCache()

	url2 := m.URL("css/app.css")
	if url1 == url2 {
		t.Errorf("expected different hash after cache clear, both: %s", url1)
	}
}
