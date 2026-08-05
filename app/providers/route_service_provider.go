package providers

import (
	"github.com/zatrano/framework/core"
	"github.com/zatrano/framework/routes"
)

// RouteServiceProvider registers application routes.
type RouteServiceProvider struct{}

// Register registers route-related services.
func (p *RouteServiceProvider) Register(app *core.Application) {}

// Boot loads route files.
func (p *RouteServiceProvider) Boot(app *core.Application) {
	routes.Web(app)
	routes.API(app)
	routes.Health(app)
}
