package broadcast

import (
	"encoding/json"
	"sync"

	"github.com/fasthttp/websocket"
	"go.uber.org/zap"
)

// Hub is a channel-oriented broadcast router (in-memory, single process).
// Use Publish / PublishJSON from handlers or workers; WebSocket and SSE subscribers receive Pusher-shaped envelopes.
type Hub struct {
	log *zap.Logger
	mu  sync.RWMutex
	// channel name -> subscribers
	channels map[string]*channelHub
}

type channelHub struct {
	mu sync.RWMutex
	// socketID -> ws client
	ws map[string]*wsConn
	// subscriber id -> sse writer
	sse map[string]*sseConn
	// presence: userID -> minimal member info (last subscribe wins)
	presence map[string]json.RawMessage
}

type wsConn struct {
	id       string
	userID   string
	userInfo map[string]any
	socket   *websocket.Conn
	writeMu  sync.Mutex
	channels map[string]struct{}
	hub      *Hub
}

type sseConn struct {
	id      string
	userID  string
	channel string
	send    chan []byte
}

// NewHub constructs an empty hub.
func NewHub(log *zap.Logger) *Hub {
	if log == nil {
		log = zap.NewNop()
	}
	return &Hub{
		log:      log,
		channels: make(map[string]*channelHub),
	}
}

func (h *Hub) getOrCreateChannel(name string) *channelHub {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch, ok := h.channels[name]
	if !ok {
		ch = &channelHub{
			ws:       make(map[string]*wsConn),
			sse:      make(map[string]*sseConn),
			presence: make(map[string]json.RawMessage),
		}
		h.channels[name] = ch
	}
	return ch
}

func (h *Hub) removeWSFromChannel(chName, socketID string) {
	h.mu.RLock()
	ch, ok := h.channels[chName]
	h.mu.RUnlock()
	if !ok {
		return
	}
	var removedUser string
	var had bool
	ch.mu.Lock()
	if c, ok := ch.ws[socketID]; ok {
		had = true
		removedUser = c.userID
		delete(ch.ws, socketID)
		if c.userID != "" {
			delete(ch.presence, c.userID)
		}
	}
	empty := len(ch.ws) == 0 && len(ch.sse) == 0
	ch.mu.Unlock()

	if had && removedUser != "" {
		if b, err := marshalMemberEvent(EventMemberRemoved, chName, json.RawMessage(mustObjectJSON("user_id", removedUser))); err == nil {
			h.broadcastRawExcept(chName, b, socketID)
		}
	}

	if empty {
		h.mu.Lock()
		delete(h.channels, chName)
		h.mu.Unlock()
	}
}

func (h *Hub) broadcastRawExcept(chName string, payload []byte, exceptSocketID string) {
	snap := h.snapshotWS(chName, exceptSocketID)
	for _, c := range snap {
		c.writeMu.Lock()
		_ = c.socket.WriteMessage(websocket.TextMessage, payload)
		c.writeMu.Unlock()
	}
}

func (h *Hub) snapshotWS(chName, exceptSocketID string) []*wsConn {
	h.mu.RLock()
	ch, ok := h.channels[chName]
	h.mu.RUnlock()
	if !ok {
		return nil
	}
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	out := make([]*wsConn, 0, len(ch.ws))
	for id, c := range ch.ws {
		if id == exceptSocketID {
			continue
		}
		out = append(out, c)
	}
	return out
}

// PublishJSON sends {event, channel, data} to WebSocket and SSE subscribers on channel.
func (h *Hub) PublishJSON(channel, event string, data any) error {
	b, err := marshalServerEvent(channel, event, data)
	if err != nil {
		return err
	}
	h.publishRaw(channel, b)
	return nil
}

func (h *Hub) publishRaw(channel string, payload []byte) {
	h.mu.RLock()
	ch, ok := h.channels[channel]
	h.mu.RUnlock()
	if !ok {
		return
	}
	ch.mu.RLock()
	wsClients := make([]*wsConn, 0, len(ch.ws))
	for _, c := range ch.ws {
		wsClients = append(wsClients, c)
	}
	sseClients := make([]*sseConn, 0, len(ch.sse))
	for _, s := range ch.sse {
		sseClients = append(sseClients, s)
	}
	ch.mu.RUnlock()

	for _, c := range wsClients {
		c.writeMu.Lock()
		_ = c.socket.WriteMessage(websocket.TextMessage, payload)
		c.writeMu.Unlock()
	}
	for _, s := range sseClients {
		select {
		case s.send <- append([]byte(nil), payload...):
		default:
			// slow consumer — drop
		}
	}
}

// OnlineOn returns user IDs currently present on a presence channel (best-effort, in-memory).
func (h *Hub) OnlineOn(channel string) []string {
	h.mu.RLock()
	ch, ok := h.channels[channel]
	h.mu.RUnlock()
	if !ok {
		return nil
	}
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	out := make([]string, 0, len(ch.presence))
	for id := range ch.presence {
		out = append(out, id)
	}
	return out
}

// RegisterSSE attaches a Server-Sent stream to channel. Call cleanup when the stream ends (client disconnect).
func (h *Hub) RegisterSSE(channel, subscriberID, userID string) (chan []byte, func()) {
	ch := h.getOrCreateChannel(channel)
	s := &sseConn{
		id:      subscriberID,
		userID:  userID,
		channel: channel,
		send:    make(chan []byte, 16),
	}
	ch.mu.Lock()
	ch.sse[subscriberID] = s
	ch.mu.Unlock()

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			ch.mu.Lock()
			delete(ch.sse, subscriberID)
			empty := len(ch.ws) == 0 && len(ch.sse) == 0
			ch.mu.Unlock()
			close(s.send)
			if empty {
				h.mu.Lock()
				delete(h.channels, channel)
				h.mu.Unlock()
			}
		})
	}

	return s.send, cleanup
}
