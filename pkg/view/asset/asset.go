// Package asset provides helpers for resolving versioned static asset URLs.
//
// It supports two resolution strategies:
//  1. Vite / esbuild manifest.json — reads hashed filenames from the build manifest.
//  2. MD5 file-hash fallback — appends ?v=<hash> to any file in PublicDir.
//
// Example (vite build produces manifest.json):
//
//	manager := asset.New(asset.Config{
//	    PublicDir:    "public",
//	    PublicURL:    "/public",
//	    ViteManifest: "public/build/.vite/manifest.json",
//	})
//	url := manager.URL("app.css")   // → /public/build/app-a1b2c3.css
//
// Hot-Module Replacement (HMR) during development:
//
//	manager := asset.New(asset.Config{
//	    PublicDir:    "public",
//	    PublicURL:    "/public",
//	    ViteDevURL:   "http://localhost:5173",  // Vite dev server
//	    DevMode:      true,
//	})
//	url := manager.URL("src/app.ts")  // → http://localhost:5173/src/app.ts
package asset

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config holds the asset manager configuration.
type Config struct {
	// PublicDir is the filesystem path to the directory serving static assets (default: "public").
	PublicDir string
	// PublicURL is the URL prefix used to build asset URLs (default: "/public").
	PublicURL string
	// ViteManifest is the path to the Vite/esbuild manifest.json (optional).
	// When set, URL() resolves hashed filenames from the manifest.
	ViteManifest string
	// ViteDevURL is the URL of the Vite development server (e.g. "http://localhost:5173").
	// When DevMode is true and ViteDevURL is non-empty, URL() proxies to the Vite dev server.
	ViteDevURL string
	// DevMode disables the manifest cache and falls back to the Vite dev server when ViteDevURL is set.
	DevMode bool
}

// ViteEntry is a single entry in a Vite/esbuild manifest.json.
type ViteEntry struct {
	File    string   `json:"file"`
	Src     string   `json:"src"`
	CSS     []string `json:"css"`
	Assets  []string `json:"assets"`
	Imports []string `json:"imports"`
	IsEntry bool     `json:"isEntry"`
}

// Manager resolves static asset URLs with cache-busting.
type Manager struct {
	cfg           Config
	mu            sync.RWMutex
	manifest      map[string]ViteEntry
	manifestMTime time.Time
	hashCache     map[string]string // filepath → short MD5 hash
}

// New creates a Manager with sensible defaults applied.
func New(cfg Config) *Manager {
	if cfg.PublicDir == "" {
		cfg.PublicDir = "public"
	}
	if cfg.PublicURL == "" {
		cfg.PublicURL = "/public"
	}
	m := &Manager{
		cfg:       cfg,
		hashCache: make(map[string]string),
	}
	m.loadManifest()
	return m
}

// URL returns the versioned URL for the given asset path (relative to PublicDir).
// Resolution order:
//  1. Vite dev server (DevMode + ViteDevURL set)
//  2. Vite/esbuild manifest
//  3. MD5 file-hash query string
//  4. Plain URL
func (m *Manager) URL(path string) string {
	path = strings.TrimPrefix(path, "/")

	// 1. Vite HMR dev server.
	if m.cfg.DevMode && m.cfg.ViteDevURL != "" {
		base := strings.TrimSuffix(m.cfg.ViteDevURL, "/")
		return base + "/" + path
	}

	// 2. Manifest.
	m.loadManifest() // no-op if mtime unchanged
	m.mu.RLock()
	entry, ok := m.manifest[path]
	m.mu.RUnlock()
	if ok {
		return m.cfg.PublicURL + "/" + entry.File
	}

	// 3. MD5 hash.
	fullPath := filepath.Join(m.cfg.PublicDir, path)
	if hash, found := m.fileHash(fullPath); found {
		return m.cfg.PublicURL + "/" + path + "?v=" + hash
	}

	// 4. Plain URL.
	return m.cfg.PublicURL + "/" + path
}

// CSS returns a versioned URL for the given CSS asset path.
func (m *Manager) CSS(path string) string { return m.URL(path) }

// JS returns a versioned URL for the given JS asset path.
func (m *Manager) JS(path string) string { return m.URL(path) }

// LinkTag returns a <link rel="stylesheet"> HTML tag for the given CSS path.
func (m *Manager) LinkTag(path string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<link rel="stylesheet" href="%s">`,
		template.HTMLEscapeString(m.URL(path)),
	))
}

// ScriptTag returns a <script src="..." defer></script> HTML tag.
func (m *Manager) ScriptTag(path string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<script src="%s" defer></script>`,
		template.HTMLEscapeString(m.URL(path)),
	))
}

// ViteHead returns the full <script> / <link> tags needed for a Vite entry point,
// including any associated CSS files discovered from the manifest.
// In DevMode it injects the Vite client HMR script instead.
func (m *Manager) ViteHead(entryPath string) template.HTML {
	if m.cfg.DevMode && m.cfg.ViteDevURL != "" {
		base := strings.TrimSuffix(m.cfg.ViteDevURL, "/")
		return template.HTML(fmt.Sprintf(
			`<script type="module" src="%s/@vite/client"></script>`+"\n"+
				`<script type="module" src="%s/%s"></script>`,
			base, base, strings.TrimPrefix(entryPath, "/"),
		))
	}

	m.loadManifest()
	m.mu.RLock()
	entry, ok := m.manifest[entryPath]
	m.mu.RUnlock()

	if !ok {
		// Fallback to plain script tag.
		return m.ScriptTag(entryPath)
	}

	var sb strings.Builder
	// CSS files first.
	for _, css := range entry.CSS {
		sb.WriteString(fmt.Sprintf(`<link rel="stylesheet" href="%s/%s">`+"\n",
			m.cfg.PublicURL, css))
	}
	// JS entry.
	sb.WriteString(fmt.Sprintf(`<script type="module" src="%s/%s" defer></script>`,
		m.cfg.PublicURL, entry.File))

	return template.HTML(sb.String())
}

// ManifestEntry returns the raw ViteEntry for the given source path, and whether
// the entry was found.
func (m *Manager) ManifestEntry(src string) (ViteEntry, bool) {
	m.loadManifest()
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.manifest[src]
	return e, ok
}

// loadManifest reads the manifest.json file when it has been modified since the
// last read. It is safe to call concurrently; it is a no-op when no manifest is
// configured.
func (m *Manager) loadManifest() {
	if m.cfg.ViteManifest == "" {
		return
	}
	fi, err := os.Stat(m.cfg.ViteManifest)
	if err != nil {
		return
	}

	m.mu.RLock()
	unchanged := !fi.ModTime().After(m.manifestMTime)
	m.mu.RUnlock()
	if unchanged {
		return
	}

	data, err := os.ReadFile(m.cfg.ViteManifest)
	if err != nil {
		return
	}
	var manifest map[string]ViteEntry
	if json.Unmarshal(data, &manifest) != nil {
		return
	}

	m.mu.Lock()
	m.manifest = manifest
	m.manifestMTime = fi.ModTime()
	m.mu.Unlock()
}

// fileHash returns a short (8-char) MD5 hex hash of the file at path.
// Results are cached in memory; the cache is never invalidated — restart the
// server to pick up changed files in non-dev mode.
func (m *Manager) fileHash(path string) (string, bool) {
	m.mu.RLock()
	if h, ok := m.hashCache[path]; ok {
		m.mu.RUnlock()
		return h, true
	}
	m.mu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	sum := md5.Sum(data) //nolint:gosec
	hash := fmt.Sprintf("%x", sum[:4])

	m.mu.Lock()
	m.hashCache[path] = hash
	m.mu.Unlock()
	return hash, true
}

// ClearCache discards cached hashes and the manifest so they will be reloaded
// on the next request. Useful in tests or after a hot deploy.
func (m *Manager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hashCache = make(map[string]string)
	m.manifest = nil
	m.manifestMTime = time.Time{}
}

// TemplateFuncs returns a template.FuncMap wiring the Manager's helpers into
// Go html/template so they can be used as {{asset "app.css"}}.
func (m *Manager) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"asset":      m.URL,
		"assetCSS":   m.CSS,
		"assetJS":    m.JS,
		"assetLink":  m.LinkTag,
		"assetScript": m.ScriptTag,
		"viteHead":   m.ViteHead,
	}
}
