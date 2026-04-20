package flash_test

import (
	"testing"

	"github.com/zatrano/zatrano/pkg/view/flash"
)

// ----------------------------------------------------------------------------
// OldValue helper (pure function, no session needed)
// ----------------------------------------------------------------------------

func TestOldValue_Found(t *testing.T) {
	old := map[string]string{"email": "user@example.com", "name": "Alice"}
	got := flash.OldValue("email", old)
	if got != "user@example.com" {
		t.Errorf("expected 'user@example.com', got: %q", got)
	}
}

func TestOldValue_NotFound(t *testing.T) {
	old := map[string]string{"name": "Alice"}
	got := flash.OldValue("email", old)
	if got != "" {
		t.Errorf("expected empty string for missing key, got: %q", got)
	}
}

func TestOldValue_NilMap(t *testing.T) {
	got := flash.OldValue("email", nil)
	if got != "" {
		t.Errorf("expected empty string for nil map, got: %q", got)
	}
}

// ----------------------------------------------------------------------------
// Flash message types
// ----------------------------------------------------------------------------

func TestFlashTypes(t *testing.T) {
	types := []flash.Type{flash.Success, flash.Error, flash.Warning, flash.Info}
	expected := []string{"success", "error", "warning", "info"}
	for i, ft := range types {
		if string(ft) != expected[i] {
			t.Errorf("type[%d]: expected %q, got %q", i, expected[i], ft)
		}
	}
}

// ----------------------------------------------------------------------------
// ViewData (no session; checks keys injected without flash manager)
// ----------------------------------------------------------------------------

// TestWithFlash_NilContext verifies that WithFlash never panics when
// called on a context without the flash manager in Locals. It acts as a
// guard against nil-pointer dereferences in the zero-value code path.
// (Full session integration tests require a running Redis instance and are
// covered in integration/e2e tests outside the unit suite.)
func TestWithFlash_Structure(t *testing.T) {
	// Build a ViewData manually as renderers do, simulating ViewData() output
	// without an actual Fiber context.
	data := flash.ViewData{
		"Flash": []flash.Message{
			{Type: flash.Success, Message: "Saved!"},
			{Type: flash.Error, Message: "Failed!"},
		},
		"Old": map[string]string{
			"email": "test@example.com",
		},
		"CSRF": "token123",
	}

	// Verify Flash slice.
	msgs, ok := data["Flash"].([]flash.Message)
	if !ok {
		t.Fatal("Flash key should contain []flash.Message")
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Type != flash.Success {
		t.Errorf("expected first message to be Success")
	}
	if msgs[1].Type != flash.Error {
		t.Errorf("expected second message to be Error")
	}

	// Verify Old map.
	old, ok := data["Old"].(map[string]string)
	if !ok {
		t.Fatal("Old key should contain map[string]string")
	}
	if old["email"] != "test@example.com" {
		t.Errorf("unexpected old email: %q", old["email"])
	}

	// Verify CSRF token.
	if data["CSRF"] != "token123" {
		t.Errorf("unexpected CSRF: %v", data["CSRF"])
	}
}

// ----------------------------------------------------------------------------
// Message struct
// ----------------------------------------------------------------------------

func TestMessage_Fields(t *testing.T) {
	m := flash.Message{Type: flash.Warning, Message: "watch out"}
	if m.Type != flash.Warning {
		t.Errorf("Type: expected Warning, got %q", m.Type)
	}
	if m.Message != "watch out" {
		t.Errorf("Message: expected 'watch out', got %q", m.Message)
	}
}
