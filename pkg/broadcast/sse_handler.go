package broadcast

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/zatrano/zatrano/pkg/config"
)

// SSEHandler streams JSON lines (each hub message is one SSE `data:` frame) for lightweight one-way push.
func SSEHandler(h *Hub, cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		if h == nil || !cfg.Broadcast.SSEEnabled {
			return fiber.NewError(fiber.StatusNotFound, "sse disabled")
		}
		chName := strings.TrimSpace(c.Params("channel"))
		if chName == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing channel")
		}
		class, _ := ClassifyChannel(chName)
		tok := readToken(c, cfg)
		var userID string
		if tok != "" {
			if mc, err := ParseAccessToken(cfg, tok); err == nil {
				userID = Subject(mc)
			}
		}
		if class != ChannelPublic && userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "token required")
		}
		if strings.HasPrefix(chName, "private-user-") {
			suffix := strings.TrimPrefix(chName, "private-user-")
			if suffix != "" && userID != suffix {
				return fiber.NewError(fiber.StatusForbidden, "forbidden")
			}
		}
		if strings.HasPrefix(chName, "presence-user-") {
			suffix := strings.TrimPrefix(chName, "presence-user-")
			if suffix != "" && userID != suffix {
				return fiber.NewError(fiber.StatusForbidden, "forbidden")
			}
		}

		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("X-Accel-Buffering", "no")

		subID := strings.ReplaceAll(uuid.NewString(), "-", "")
		events, cleanup := h.RegisterSSE(chName, subID, userID)

		return c.SendStreamWriter(func(w *bufio.Writer) {
			defer cleanup()
			_, _ = fmt.Fprintf(w, ": connected %s\n\n", time.Now().Format(time.RFC3339))
			_ = w.Flush()
			ctx := c.Context()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-events:
					if !ok {
						return
					}
					_, _ = fmt.Fprintf(w, "data: %s\n\n", string(msg))
					_ = w.Flush()
				}
			}
		})
	}
}
