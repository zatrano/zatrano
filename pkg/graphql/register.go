package graphql

import (
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

	"github.com/zatrano/zatrano/pkg/core"
	"github.com/zatrano/zatrano/pkg/graphql/graph"
)

// Register mounts GraphQL and optional GraphiQL playground when graphql.enabled is true.
func Register(a *core.App, app *fiber.App) {
	if a == nil || a.Config == nil || !a.Config.GraphQL.Enabled {
		return
	}
	cfg := a.Config.GraphQL
	path := stringsTrimPath(cfg.Path)
	if path == "" {
		path = "/graphql"
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{App: a},
	}))

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithLoaders(r.Context(), NewLoaders(a))
		srv.ServeHTTP(w, r.WithContext(ctx))
	})

	app.All(path, adaptor.HTTPHandler(h))

	if cfg.Playground {
		pp := stringsTrimPath(cfg.PlaygroundPath)
		if pp == "" {
			pp = "/playground"
		}
		pg := playground.Handler("GraphQL playground", path)
		app.Get(pp, adaptor.HTTPHandler(pg))
	}
}

func stringsTrimPath(s string) string {
	s = strings.TrimSpace(s)
	for len(s) > 1 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
