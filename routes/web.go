package routes

import (
	"strings"

	"github.com/zatrano/framework/app/http/controllers/web"
	"github.com/zatrano/framework/core"
	"github.com/zatrano/framework/core/http"
	"github.com/zatrano/framework/core/localization"
	"github.com/zatrano/framework/core/routing"
)

// Web registers application web routes.
func Web(app *core.Application) {
	router := app.Router()

	routing.Controller(router, &web.HomeController{}, func(r *routing.Router, c *web.HomeController) {
		r.Get("/", c.Index).As("home")
	})

	router.Post("/locale", func(req *http.Request) *http.Response {
		locale := strings.TrimSpace(strings.ToLower(req.Input("locale")))
		langPath := app.BasePath("lang")
		if localization.Published(langPath) && localization.HasLocale(langPath, locale) {
			if sess := req.Session(); sess != nil {
				sess.Put("locale", locale)
			}
			if app.Translator() != nil {
				app.Translator().SetLocale(locale)
				_ = app.Translator().Load(locale)
			}
		}
		return http.Redirect("/")
	}).As("locale")
}
