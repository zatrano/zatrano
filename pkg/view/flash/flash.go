// Package flash implements a session-backed flash message system for
// server-rendered forms. Messages set in one HTTP handler are available to
// templates in the very next request, then automatically cleared.
//
// It also stores old-input values (field → value) so that forms can
// repopulate themselves after a failed validation.
//
// Usage in a handler (save flash then redirect):
//
//	flash.Set(c, flash.Success, "Record saved successfully.")
//	flash.SetOld(c, map[string]string{"email": "user@example.com"})
//	return c.Redirect("/dashboard")
//
// Usage in a template (via ViewData):
//
//	{{range .Flash}}
//	  <div class="alert alert-{{.Type}}">{{.Message}}</div>
//	{{end}}
//	{{old "email" .Old}}
package flash

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// Type represents the flash message category.
type Type string

const (
	Success Type = "success"
	Error   Type = "error"
	Warning Type = "warning"
	Info    Type = "info"
)

// Message is a single flash message entry.
type Message struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

const (
	sessionKeyFlash = "zatrano.flash"
	sessionKeyOld   = "zatrano.old"
)

// Manager wraps the Fiber session store for flash operations.
type Manager struct {
	store *session.Store
}

// New returns a Manager backed by the given session.Store.
func New(store *session.Store) *Manager {
	return &Manager{store: store}
}

// Set adds a flash message of the given type. Multiple messages of any type can
// coexist; they are all drained on the next Get call.
func (m *Manager) Set(c fiber.Ctx, t Type, msg string) error {
	sess, err := m.store.Get(c)
	if err != nil {
		return err
	}
	msgs := m.currentMessages(sess)
	msgs = append(msgs, Message{Type: t, Message: msg})
	return m.saveMessages(sess, msgs)
}

// Get returns all flash messages for the current request and clears them from
// the session so they are not shown again.
func (m *Manager) Get(c fiber.Ctx) ([]Message, error) {
	sess, err := m.store.Get(c)
	if err != nil {
		return nil, err
	}
	msgs := m.currentMessages(sess)
	sess.Delete(sessionKeyFlash)
	return msgs, sess.Save()
}

// SetOld stores the old-input map (field name → string value) in the session.
// Call this before redirecting after a failed form submission.
func (m *Manager) SetOld(c fiber.Ctx, input map[string]string) error {
	sess, err := m.store.Get(c)
	if err != nil {
		return err
	}
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	sess.Set(sessionKeyOld, string(data))
	return sess.Save()
}

// GetOld returns the old-input map and clears it from the session.
func (m *Manager) GetOld(c fiber.Ctx) (map[string]string, error) {
	sess, err := m.store.Get(c)
	if err != nil {
		return nil, err
	}
	raw, _ := sess.Get(sessionKeyOld).(string)
	sess.Delete(sessionKeyOld)
	if saveErr := sess.Save(); saveErr != nil {
		return nil, saveErr
	}
	if raw == "" {
		return map[string]string{}, nil
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return map[string]string{}, nil
	}
	return result, nil
}

// currentMessages reads existing flash messages from the session without clearing.
func (m *Manager) currentMessages(sess *session.Session) []Message {
	raw, _ := sess.Get(sessionKeyFlash).(string)
	if raw == "" {
		return nil
	}
	var msgs []Message
	if err := json.Unmarshal([]byte(raw), &msgs); err != nil {
		return nil
	}
	return msgs
}

func (m *Manager) saveMessages(sess *session.Session, msgs []Message) error {
	data, err := json.Marshal(msgs)
	if err != nil {
		return err
	}
	sess.Set(sessionKeyFlash, string(data))
	return sess.Save()
}

// ----------------------------------------------------------------------------
// Package-level helpers (require Manager in Locals)
// ----------------------------------------------------------------------------

const localsManager = "zatrano.flash.manager"

// Middleware stores the flash Manager on c.Locals so that Set/Get/SetOld/GetOld
// package-level helpers work without passing the manager explicitly.
func Middleware(m *Manager) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals(localsManager, m)
		return c.Next()
	}
}

// fromCtx extracts the Manager from Locals; returns nil if not set.
func fromCtx(c fiber.Ctx) *Manager {
	m, _ := c.Locals(localsManager).(*Manager)
	return m
}

// Set is a convenience wrapper around Manager.Set.
func Set(c fiber.Ctx, t Type, msg string) error {
	m := fromCtx(c)
	if m == nil {
		return nil
	}
	return m.Set(c, t, msg)
}

// Get is a convenience wrapper around Manager.Get.
func Get(c fiber.Ctx) ([]Message, error) {
	m := fromCtx(c)
	if m == nil {
		return nil, nil
	}
	return m.Get(c)
}

// SetOld is a convenience wrapper around Manager.SetOld.
func SetOld(c fiber.Ctx, input map[string]string) error {
	m := fromCtx(c)
	if m == nil {
		return nil
	}
	return m.SetOld(c, input)
}

// GetOld is a convenience wrapper around Manager.GetOld.
func GetOld(c fiber.Ctx) (map[string]string, error) {
	m := fromCtx(c)
	if m == nil {
		return map[string]string{}, nil
	}
	return m.GetOld(c)
}

// ----------------------------------------------------------------------------
// ViewData helper — builds the map that templates receive
// ----------------------------------------------------------------------------

// ViewData is the canonical template data map enriched with flash and old-input.
type ViewData map[string]any

// WithFlash loads flash messages and old input from the session into the map
// and returns the enriched map.  Keys added:
//
//   - "Flash"   []flash.Message
//   - "Old"     map[string]string
//   - "CSRF"    string (from the X-CSRF-Token header or form value)
func WithFlash(c fiber.Ctx, data ViewData) ViewData {
	if data == nil {
		data = ViewData{}
	}
	msgs, _ := Get(c)
	old, _ := GetOld(c)
	if msgs == nil {
		msgs = []Message{}
	}
	if old == nil {
		old = map[string]string{}
	}
	data["Flash"] = msgs
	data["Old"] = old
	// Expose CSRF token for form builder usage.
	csrfToken := c.Locals("fiber:csrf_token")
	if csrfToken == nil {
		csrfToken = c.Get("X-CSRF-Token")
	}
	data["CSRF"] = csrfToken
	return data
}

// OldValue looks up a field in the old-input map stored under "Old" in data.
// Returns empty string when not found.  Use in templates as {{old "email" .Old}}.
func OldValue(field string, old map[string]string) string {
	if old == nil {
		return ""
	}
	return old[field]
}
