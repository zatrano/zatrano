package broadcast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type presenceChannelData struct {
	UserID   string         `json:"user_id"`
	UserInfo map[string]any `json:"user_info"`
}

// SubscribeWebSocket registers client on channel (Pusher-shaped ack bytes).
func (h *Hub) SubscribeWebSocket(client *wsConn, channelName, channelData string) ([]byte, error) {
	class, chName := ClassifyChannel(channelName)
	if chName == "" {
		return nil, fmt.Errorf("empty channel")
	}
	if class != ChannelPublic && client.userID == "" {
		return nil, fmt.Errorf("jwt required for channel %s", chName)
	}
	if strings.HasPrefix(chName, "private-user-") {
		suffix := strings.TrimPrefix(chName, "private-user-")
		if suffix != "" && client.userID != suffix {
			return nil, fmt.Errorf("forbidden")
		}
	}
	if strings.HasPrefix(chName, "presence-user-") {
		suffix := strings.TrimPrefix(chName, "presence-user-")
		if suffix != "" && client.userID != suffix {
			return nil, fmt.Errorf("forbidden")
		}
	}

	userInfo := map[string]any{}
	if class == ChannelPresence && strings.TrimSpace(channelData) != "" {
		var pcd presenceChannelData
		if err := json.Unmarshal([]byte(channelData), &pcd); err != nil {
			return nil, fmt.Errorf("invalid channel_data")
		}
		if pcd.UserID != "" && pcd.UserID != client.userID {
			return nil, fmt.Errorf("user_id mismatch")
		}
		userInfo = pcd.UserInfo
	}

	ch := h.getOrCreateChannel(chName)
	ch.mu.Lock()
	ch.ws[client.id] = client
	var pres *presencePayload
	if class == ChannelPresence {
		member := map[string]any{
			"user_id":   client.userID,
			"user_info": userInfo,
		}
		raw, _ := json.Marshal(member)
		if client.userID != "" {
			ch.presence[client.userID] = json.RawMessage(raw)
		}
		pres = &presencePayload{}
		pres.Presence.Count = len(ch.presence)
		for uid := range ch.presence {
			pres.Presence.IDs = append(pres.Presence.IDs, uid)
		}
		if pres.Presence.Hash == nil {
			pres.Presence.Hash = make(map[string]string)
		}
		for uid, r := range ch.presence {
			pres.Presence.Hash[uid] = string(r)
		}
	}
	ch.mu.Unlock()

	client.channels[chName] = struct{}{}

	if class == ChannelPresence && client.userID != "" {
		member := map[string]any{
			"user_id":   client.userID,
			"user_info": userInfo,
		}
		raw, _ := json.Marshal(member)
		if b, err := marshalMemberEvent(EventMemberAdded, chName, json.RawMessage(raw)); err == nil {
			h.broadcastRawExcept(chName, b, client.id)
		}
	}

	return marshalSubscriptionSucceeded(chName, pres)
}

// UnsubscribeWebSocket removes client from a single channel.
func (h *Hub) UnsubscribeWebSocket(client *wsConn, channelName string) {
	if _, ok := client.channels[channelName]; !ok {
		return
	}
	delete(client.channels, channelName)
	h.removeWSFromChannel(channelName, client.id)
}

// DetachWebSocket removes client from all subscribed channels.
func (h *Hub) DetachWebSocket(client *wsConn) {
	for ch := range client.channels {
		h.removeWSFromChannel(ch, client.id)
	}
	client.channels = make(map[string]struct{})
}

func newSocketID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
