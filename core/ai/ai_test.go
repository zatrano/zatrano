package ai_test

import (
	"strings"
	"testing"

	"github.com/zatrano/framework/core/ai"
)

func TestAIChat(t *testing.T) {
	m := ai.New()
	resp, err := m.Chat(ai.ChatRequest{
		Messages: []ai.Message{{Role: "user", Content: "ping"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Message.Content, "ping") {
		t.Fatalf("%v", resp.Message.Content)
	}
	if resp.Usage.TotalTokens < 1 {
		t.Fatal("usage")
	}
}
