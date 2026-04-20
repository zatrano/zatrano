package config

// View holds view engine configuration.
type View struct {
	// Root is the template root directory (default: "views").
	Root string `mapstructure:"root"`
	// Extension is the template file extension (default: ".html").
	Extension string `mapstructure:"extension"`
	// ComponentsDir is the sub-directory for reusable components (default: "components").
	ComponentsDir string `mapstructure:"components_dir"`
	// LayoutsDir is the sub-directory for layout templates (default: "layouts").
	LayoutsDir string `mapstructure:"layouts_dir"`
	// DevMode disables template caching (default: false; mirrors log_development when unset).
	DevMode bool `mapstructure:"dev_mode"`

	// Asset holds static asset pipeline configuration.
	Asset ViewAsset `mapstructure:"asset"`
}

// ViewAsset holds asset pipeline configuration embedded in View.
type ViewAsset struct {
	// PublicDir is the filesystem path to the public static assets directory (default: "public").
	PublicDir string `mapstructure:"public_dir"`
	// PublicURL is the URL prefix for static assets (default: "/public").
	PublicURL string `mapstructure:"public_url"`
	// ViteManifest is the path to the Vite/esbuild manifest.json (optional).
	// e.g. "public/build/.vite/manifest.json"
	ViteManifest string `mapstructure:"vite_manifest"`
	// ViteDevURL is the Vite development server URL (e.g. "http://localhost:5173").
	// When non-empty and View.DevMode is true, asset URLs proxy to this server.
	ViteDevURL string `mapstructure:"vite_dev_url"`
}
