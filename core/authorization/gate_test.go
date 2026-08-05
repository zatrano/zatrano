package authorization_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/zatrano/framework/core/auth"
	"github.com/zatrano/framework/core/authorization"
	"github.com/zatrano/framework/core/http"
)

type fakeUser struct{ id any }

func (u fakeUser) AuthID() any          { return u.id }
func (u fakeUser) AuthPassword() string { return "" }

func TestGateAndPolicy(t *testing.T) {
	gate := authorization.New()
	gate.Define("always", func(user auth.Authenticatable, arguments ...any) bool {
		return user != nil
	})
	gate.Policy("post", authorization.NewPolicy().Define("update", func(user auth.Authenticatable, arguments ...any) bool {
		return user.AuthID() == arguments[0]
	}))

	user := fakeUser{id: 7}
	if !gate.Allows(user, "always") {
		t.Fatal("expected always allow")
	}
	if !gate.Check(user, "always") {
		t.Fatal("Check should alias Allows")
	}
	if !gate.Allows(user, "post.update", 7) {
		t.Fatal("expected policy allow")
	}
	if gate.Allows(user, "post.update", 8) {
		t.Fatal("expected policy deny")
	}
	if !gate.Has("always") || !gate.Has("post.update") || gate.Has("missing") {
		t.Fatal("Has ability registry")
	}
}

func TestBeforeAfterAndForUser(t *testing.T) {
	gate := authorization.New()
	gate.Define("edit", func(user auth.Authenticatable, arguments ...any) bool {
		return false
	})

	trueVal := true
	falseVal := false
	gate.Before(func(user auth.Authenticatable, ability string, arguments ...any) *bool {
		if ability == "super" {
			return &trueVal
		}
		return nil
	})
	gate.After(func(user auth.Authenticatable, ability string, result bool, arguments ...any) *bool {
		if ability == "edit" && len(arguments) > 0 && arguments[0] == "force" {
			return &trueVal
		}
		if ability == "super" {
			return &falseVal
		}
		return nil
	})

	user := fakeUser{id: 1}
	if gate.Allows(user, "super") {
		t.Fatal("after should deny super (overrides before allow)")
	}
	if gate.Allows(user, "edit") {
		t.Fatal("edit should deny by default")
	}
	if !gate.Allows(user, "edit", "force") {
		t.Fatal("after should force-allow edit")
	}

	pending := gate.ForUser(user)
	if !pending.Any([]string{"edit", "missing"}, "force") {
		t.Fatal("Any should allow edit with force")
	}
	if !pending.None([]string{"missing", "nope"}) {
		t.Fatal("None should be true when all denied")
	}
	if err := pending.Authorize("edit", "force"); err != nil {
		t.Fatal(err)
	}
	if err := pending.Authorize("edit"); err == nil {
		t.Fatal("expected authorize error")
	} else if _, ok := err.(authorization.AuthorizationException); !ok {
		t.Fatalf("expected AuthorizationException, got %T", err)
	}
}

func TestAnyNoneAndAuthorizeResponse(t *testing.T) {
	gate := authorization.New()
	gate.Define("read", func(user auth.Authenticatable, arguments ...any) bool { return true })
	gate.Define("write", func(user auth.Authenticatable, arguments ...any) bool { return false })

	user := fakeUser{id: 1}
	if !gate.Any(user, []string{"write", "read"}) {
		t.Fatal("Any")
	}
	if !gate.None(user, []string{"write", "delete"}) {
		t.Fatal("None")
	}
	err := gate.Authorize(user, "write")
	resp := authorization.ResponseFor(err)
	if resp == nil || resp.StatusCode() != 403 {
		t.Fatalf("response=%v", resp)
	}
}

func TestMiddlewareWithArgs(t *testing.T) {
	gate := authorization.New()
	gate.Define("post.update", func(user auth.Authenticatable, arguments ...any) bool {
		return len(arguments) == 1 && arguments[0] == "42"
	})
	gate.Define("read", func(user auth.Authenticatable, arguments ...any) bool { return true })

	mgr := auth.NewManager("web")
	mgr.Extend("web", auth.NewGuard("web", nil))

	user := fakeUser{id: 9}
	raw := httptest.NewRequest(stdhttp.MethodGet, "/posts/42", nil)
	req := http.NewRequest(raw)
	req.Set("auth.user", user)
	req.SetRouteParams(map[string]string{"id": "42"})

	called := false
	mw := authorization.Middleware(gate, mgr, "post.update", func(r *http.Request) any {
		return r.Route("id")
	})
	resp := mw(func(r *http.Request) *http.Response {
		called = true
		return http.JSON(map[string]any{"ok": true})
	})(req)
	if !called || resp.StatusCode() != 200 {
		t.Fatalf("expected allow status=%d called=%v", resp.StatusCode(), called)
	}

	deny := authorization.Middleware(gate, mgr, "post.update", "99")(func(r *http.Request) *http.Response {
		t.Fatal("should not run")
		return nil
	})(req)
	if deny.StatusCode() != 403 {
		t.Fatalf("expected 403, got %d", deny.StatusCode())
	}

	anyMW := authorization.MiddlewareAny(gate, mgr, []string{"missing", "read"})
	anyResp := anyMW(func(r *http.Request) *http.Response {
		return http.JSON(map[string]any{"ok": true})
	})(req)
	if anyResp.StatusCode() != 200 {
		t.Fatalf("MiddlewareAny status=%d", anyResp.StatusCode())
	}
}
