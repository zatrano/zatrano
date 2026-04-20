package broadcast

import (
	"encoding/json"
	"strings"
)

// Pusher-compatible JSON envelopes (subset for laravel-echo / pusher-js style clients).

const (
	EventConnectionEstablished = "pusher:connection_established"
	EventSubscribe             = "pusher:subscribe"
	EventUnsubscribe           = "pusher:unsubscribe"
	EventPing                  = "pusher:ping"
	EventPong                  = "pusher:pong"
	EventSubscriptionSucceeded = "pusher_internal:subscription_succeeded"
	EventMemberAdded           = "pusher_internal:member_added"
	EventMemberRemoved         = "pusher_internal:member_removed"
	EventError                 = "pusher:error"
)

// Envelope is a wire message from/to clients.
type Envelope struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type subscribeData struct {
	Channel     string `json:"channel"`
	ChannelData string `json:"channel_data,omitempty"`
	Auth        string `json:"auth,omitempty"`
}

func parseSubscribeData(raw json.RawMessage) (subscribeData, error) {
	var sd subscribeData
	if len(raw) == 0 {
		return sd, errInvalidSubscribe
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return sd, errInvalidSubscribe
		}
		if err := json.Unmarshal([]byte(s), &sd); err != nil {
			return sd, errInvalidSubscribe
		}
		return sd, nil
	}
	if err := json.Unmarshal(raw, &sd); err != nil {
		return sd, errInvalidSubscribe
	}
	return sd, nil
}

type presencePayload struct {
	Presence struct {
		Count int               `json:"count"`
		IDs   []string          `json:"ids"`
		Hash  map[string]string `json:"hash,omitempty"`
	} `json:"presence"`
}

func marshalConnectionEstablished(socketID string) ([]byte, error) {
	inner, err := json.Marshal(map[string]string{"socket_id": socketID})
	if err != nil {
		return nil, err
	}
	return json.Marshal(Envelope{
		Event: EventConnectionEstablished,
		Data:  json.RawMessage(quoteJSON(string(inner))),
	})
}

func marshalSubscriptionSucceeded(channel string, presence *presencePayload) ([]byte, error) {
	var data []byte
	var err error
	if presence != nil {
		data, err = json.Marshal(presence)
	} else {
		data = []byte("{}")
	}
	if err != nil {
		return nil, err
	}
	return json.Marshal(Envelope{
		Event:   EventSubscriptionSucceeded,
		Channel: channel,
		Data:    data,
	})
}

func marshalMemberEvent(event, channel string, member json.RawMessage) ([]byte, error) {
	return json.Marshal(Envelope{
		Event:   event,
		Channel: channel,
		Data:    member,
	})
}

func marshalServerEvent(channel, event string, data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Envelope{
		Event:   event,
		Channel: channel,
		Data:    json.RawMessage(b),
	})
}

func marshalPusherError(message string) ([]byte, error) {
	return json.Marshal(Envelope{
		Event: EventError,
		Data:  json.RawMessage(mustObjectJSON("message", message)),
	})
}

func quoteJSON(s string) json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

func mustObjectJSON(k, v string) json.RawMessage {
	m := map[string]string{k: v}
	b, _ := json.Marshal(m)
	return json.RawMessage(b)
}

// ChannelClass describes subscription rules for a channel name.
type ChannelClass int

const (
	ChannelPublic ChannelClass = iota
	ChannelPrivate
	ChannelPresence
)

func ClassifyChannel(name string) (ChannelClass, string) {
	name = strings.TrimSpace(name)
	switch {
	case strings.HasPrefix(name, "presence-"):
		return ChannelPresence, name
	case strings.HasPrefix(name, "private-"):
		return ChannelPrivate, name
	default:
		return ChannelPublic, name
	}
}
