package providers

import (
	"github.com/zatrano/framework/core"
	"github.com/zatrano/framework/core/middleware/csrf"
)

// AppServiceProvider registers application-level services.
type AppServiceProvider struct{}

// Register registers bindings into the container.
func (p *AppServiceProvider) Register(app *core.Application) {
	// Application-specific bindings go here.
}

// Boot boots application services.
func (p *AppServiceProvider) Boot(app *core.Application) {
	app.Router().Use(csrf.Except("/api"))

	app.View().Share("appUrl", app.Config().GetString("app.url"))
	if m := app.Assets(); m != nil && len(m.All()) > 0 {
		app.View().Share("assetCss", m.URL("css/app.css"))
		app.View().Share("assetJs", m.URL("js/app.js"))
	}
}
