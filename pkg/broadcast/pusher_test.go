package broadcast

import (
	"encoding/json"
	"testing"
)

func TestClassifyChannel(t *testing.T) {
	c, n := ClassifyChannel("public-news")
	if c != ChannelPublic || n != "public-news" {
		t.Fatalf("public: got %v %q", c, n)
	}
	c, n = ClassifyChannel("private-user-1")
	if c != ChannelPrivate || n != "private-user-1" {
		t.Fatalf("private: got %v %q", c, n)
	}
	c, n = ClassifyChannel("presence-room")
	if c != ChannelPresence || n != "presence-room" {
		t.Fatalf("presence: got %v %q", c, n)
	}
}

func TestMarshalConnectionEstablished(t *testing.T) {
	b, err := marshalConnectionEstablished("abc123")
	if err != nil {
		t.Fatal(err)
	}
	var env Envelope
	if err := json.Unmarshal(b, &env); err != nil {
		t.Fatal(err)
	}
	if env.Event != EventConnectionEstablished {
		t.Fatalf("event %q", env.Event)
	}
	var innerJSON string
	if err := json.Unmarshal(env.Data, &innerJSON); err != nil {
		t.Fatal(err)
	}
	var inner map[string]string
	if err := json.Unmarshal([]byte(innerJSON), &inner); err != nil {
		t.Fatal(err)
	}
	if inner["socket_id"] != "abc123" {
		t.Fatalf("socket_id %v", inner)
	}
}
