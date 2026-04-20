package broadcast

import (
	"encoding/json"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"github.com/zatrano/zatrano/pkg/config"
)

// websocketUpgrader builds a FastHTTP upgrader with optional origin checks.
func websocketUpgrader(cfg *config.Config) websocket.FastHTTPUpgrader {
	allow := cfg.Broadcast.AllowOrigins
	return websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			if len(allow) == 0 {
				return true
			}
			origin := string(ctx.Request.Header.Peek("Origin"))
			if origin == "" {
				return true
			}
			for _, o := range allow {
				if o == origin {
					return true
				}
			}
			for _, o := range cfg.HTTP.CORSAllowOrigins {
				if o == origin || o == "*" {
					return true
				}
			}
			return false
		},
	}
}

func readToken(c fiber.Ctx, cfg *config.Config) string {
	q := strings.TrimSpace(cfg.Broadcast.JWTQueryParam)
	if q == "" {
		q = "access_token"
	}
	if t := strings.TrimSpace(c.Query(q)); t != "" {
		return t
	}
	return strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
}

// WebSocketHandler upgrades to WebSocket and speaks a Pusher-compatible protocol subset.
func WebSocketHandler(h *Hub, cfg *config.Config) fiber.Handler {
	up := websocketUpgrader(cfg)
	return func(c fiber.Ctx) error {
		if h == nil {
			return fiber.NewError(fiber.StatusNotFound, "broadcast disabled")
		}
		var claims jwt.MapClaims
		if tok := readToken(c, cfg); tok != "" {
			mc, err := ParseAccessToken(cfg, tok)
			if err == nil {
				claims = mc
			}
		}
		userID := Subject(claims)

		err := up.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			sid := newSocketID()
			client := &wsConn{
				id:       sid,
				userID:   userID,
				channels: make(map[string]struct{}),
				socket:   conn,
			}
			welcome, werr := marshalConnectionEstablished(sid)
			if werr != nil {
				_ = conn.Close()
				return
			}
			client.writeMu.Lock()
			_ = conn.WriteMessage(websocket.TextMessage, welcome)
			client.writeMu.Unlock()

			for {
				_, payload, rerr := conn.ReadMessage()
				if rerr != nil {
					h.DetachWebSocket(client)
					return
				}
				var env Envelope
				if json.Unmarshal(payload, &env) != nil {
					continue
				}
				switch env.Event {
				case EventPing:
					if b, err := json.Marshal(Envelope{Event: EventPong, Data: json.RawMessage(`{}`)}); err == nil {
						client.writeMu.Lock()
						_ = conn.WriteMessage(websocket.TextMessage, b)
						client.writeMu.Unlock()
					}
				case EventSubscribe:
					sd, err := parseSubscribeData(env.Data)
					if err != nil {
						if b, e2 := marshalPusherError("invalid subscribe"); e2 == nil {
							client.writeMu.Lock()
							_ = conn.WriteMessage(websocket.TextMessage, b)
							client.writeMu.Unlock()
						}
						continue
					}
					ack, err := h.SubscribeWebSocket(client, sd.Channel, sd.ChannelData)
					if err != nil {
						if b, e2 := marshalPusherError(err.Error()); e2 == nil {
							client.writeMu.Lock()
							_ = conn.WriteMessage(websocket.TextMessage, b)
							client.writeMu.Unlock()
						}
						continue
					}
					client.writeMu.Lock()
					_ = conn.WriteMessage(websocket.TextMessage, ack)
					client.writeMu.Unlock()
				case EventUnsubscribe:
					ch := strings.TrimSpace(env.Channel)
					if ch != "" {
						h.UnsubscribeWebSocket(client, ch)
					}
				default:
					// ignore unknown client events
				}
			}
		})
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return nil
	}
}
