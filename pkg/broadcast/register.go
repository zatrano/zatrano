package broadcast

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/zatrano/zatrano/pkg/config"
)

// Register mounts WebSocket and optional SSE routes under cfg.Broadcast.PathPrefix.
func Register(h *Hub, cfg *config.Config, app *fiber.App) {
	if h == nil || !cfg.Broadcast.Enabled {
		return
	}
	prefix := strings.TrimSpace(cfg.Broadcast.PathPrefix)
	if prefix == "" {
		prefix = "/broadcast"
	}
	prefix = strings.TrimRight(prefix, "/")
	if prefix == "" {
		prefix = "/broadcast"
	}

	g := app.Group(prefix)
	g.Get("/ws", WebSocketHandler(h, cfg))
	if cfg.Broadcast.SSEEnabled {
		g.Get("/sse/:channel", SSEHandler(h, cfg))
	}
}
