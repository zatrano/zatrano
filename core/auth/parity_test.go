package auth_test

import (
	"errors"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zatrano/framework/core/auth"
	"github.com/zatrano/framework/core/encryption"
	"github.com/zatrano/framework/core/hashing"
	"github.com/zatrano/framework/core/http"
	"github.com/zatrano/framework/core/session"
	"github.com/zatrano/framework/core/totp"
)

func TestValidateOnceLoginUsingIDAndViaRemember(t *testing.T) {
	hash, err := hashing.Hash("secret")
	if err != nil {
		t.Fatal(err)
	}
	provider := newMemoryUserProvider()
	user, err := provider.Create(map[string]any{
		"email": "once@zatrano.test", "password": hash, "name": "Once",
	})
	if err != nil {
		t.Fatal(err)
	}
	manager := auth.NewManager("web")
	manager.Extend("web", auth.NewGuard("web", provider))

	ok, err := manager.Validate(map[string]string{"email": "once@zatrano.test", "password": "secret"})
	if err != nil || !ok {
		t.Fatalf("validate=%v err=%v", ok, err)
	}
	ok, err = manager.Validate(map[string]string{"email": "once@zatrano.test", "password": "bad"})
	if err != nil || ok {
		t.Fatalf("expected invalid credentials")
	}

	req := http.NewRequest(httptest.NewRequest(stdhttp.MethodGet, "/", nil))
	req.SetSession(&memSession{data: map[string]any{}})
	if err := manager.Once(req, user); err != nil {
		t.Fatal(err)
	}
	if manager.User(req) == nil {
		t.Fatal("once user missing")
	}
	if req.Session().Get("auth_user_id") != nil {
		t.Fatal("once must not write session")
	}

	req2 := http.NewRequest(httptest.NewRequest(stdhttp.MethodGet, "/", nil))
	req2.SetSession(&memSession{data: map[string]any{}})
	if err := manager.LoginUsingID(req2, user.AuthID()); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(manager.ID(req2)) != fmt.Sprint(user.AuthID()) {
		t.Fatal("login using id failed")
	}

	req3 := http.NewRequest(httptest.NewRequest(stdhttp.MethodGet, "/", nil))
	req3.SetSession(&memSession{data: map[string]any{}})
	if err := manager.Login(req3, user, true); err != nil {
		t.Fatal(err)
	}
	var cookieValue string
	for _, c := range req3.Cookies().Apply() {
		if c.Name == "remember_web" {
			cookieValue = c.Value
		}
	}
	raw := httptest.NewRequest(stdhttp.MethodGet, "/", nil)
	raw.AddCookie(&stdhttp.Cookie{Name: "remember_web", Value: cookieValue})
	req4 := http.NewRequest(raw)
	req4.SetSession(&memSession{data: map[string]any{}})
	if manager.User(req4) == nil || !manager.ViaRemember(req4) {
		t.Fatal("via remember expected")
	}
}

func TestLogoutOtherDevicesDestroysForeignSessions(t *testing.T) {
	dir := t.TempDir()
	sessMgr := session.NewManager(dir, 120)
	hash, _ := hashing.Hash("secret")
	provider := newMemoryUserProvider()
	user, _ := provider.Create(map[string]any{"email": "devices@zatrano.test", "password": hash})
	manager := auth.NewManager("web")
	manager.Extend("web", auth.NewGuard("web", provider))
	manager.SetSessionManager(sessMgr)

	other, err := sessMgr.Start("")
	if err != nil {
		t.Fatal(err)
	}
	other.Put("auth_user_id", fmt.Sprint(user.AuthID()))
	if err := sessMgr.Save(other); err != nil {
		t.Fatal(err)
	}

	current, err := sessMgr.Start("")
	if err != nil {
		t.Fatal(err)
	}
	req := http.NewRequest(httptest.NewRequest(stdhttp.MethodPost, "/", nil))
	req.SetSession(current)
	if err := manager.Login(req, user); err != nil {
		t.Fatal(err)
	}
	_ = sessMgr.Save(current)

	if err := manager.LogoutOtherDevices(req, "secret"); err != nil {
		t.Fatal(err)
	}
	// foreign session file should be gone
	bag, err := sessMgr.Start(other.ID())
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(bag.Get("auth_user_id")) == fmt.Sprint(user.AuthID()) {
		t.Fatal("other session should have been destroyed")
	}
}

func TestLockoutAfterFailedAttempts(t *testing.T) {
	hash, _ := hashing.Hash("secret")
	provider := newMemoryUserProvider()
	_, _ = provider.Create(map[string]any{"email": "lock@zatrano.test", "password": hash})
	manager := auth.NewManager("web")
	manager.Extend("web", auth.NewGuard("web", provider))
	manager.SetLockout(3, time.Minute)

	req := http.NewRequest(httptest.NewRequest(stdhttp.MethodPost, "/", nil))
	req.SetSession(&memSession{data: map[string]any{}})
	for i := 0; i < 3; i++ {
		_, err := manager.Attempt(req, map[string]string{"email": "lock@zatrano.test", "password": "bad"})
		if i == 2 && !errors.Is(err, auth.ErrLockout) {
			t.Fatalf("expected lockout, got %v", err)
		}
	}
}

func TestTwoFactorChallengeFlow(t *testing.T) {
	hash, _ := hashing.Hash("secret")
	provider := newMemoryUserProvider()
	user, _ := provider.Create(map[string]any{"email": "mfa@zatrano.test", "password": hash, "name": "MFA"})
	manager := auth.NewManager("web")
	manager.Extend("web", auth.NewGuard("web", provider))

	req := http.NewRequest(httptest.NewRequest(stdhttp.MethodPost, "/", nil))
	req.SetSession(&memSession{data: map[string]any{}})
	if err := manager.Login(req, user); err != nil {
		t.Fatal(err)
	}
	secret, _, _, err := manager.EnableTwoFactor(user)
	if err != nil {
		t.Fatal(err)
	}
	code, err := totp.Code(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.ConfirmTwoFactor(user, code); err != nil {
		t.Fatal(err)
	}
	_ = manager.Logout(req)

	ok, err := manager.Attempt(req, map[string]string{"email": "mfa@zatrano.test", "password": "secret"})
	if ok || !errors.Is(err, auth.ErrTwoFactorRequired) {
		t.Fatalf("expected two-factor required, ok=%v err=%v", ok, err)
	}
	code, err = totp.Code(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	ok, err = manager.ChallengeTwoFactor(req, code)
	if err != nil || !ok {
		t.Fatalf("challenge failed: %v", err)
	}
	if !manager.Check(req) {
		t.Fatal("expected authenticated after challenge")
	}
}

func TestShouldUseAndOnceBasic(t *testing.T) {
	hash, _ := hashing.Hash("secret")
	provider := newMemoryUserProvider()
	_, _ = provider.Create(map[string]any{"email": "basic@zatrano.test", "password": hash})
	manager := auth.NewManager("web")
	manager.Extend("web", auth.NewGuard("web", provider))
	manager.Extend("api", auth.NewGuard("api", provider))
	manager.ShouldUse("api")
	if manager.GetDefaultDriver() != "api" {
		t.Fatalf("default=%q", manager.GetDefaultDriver())
	}
	manager.SetDefaultDriver("web")
	if manager.GetDefaultDriver() != "web" {
		t.Fatal("set default driver failed")
	}

	raw := httptest.NewRequest(stdhttp.MethodGet, "/", nil)
	raw.SetBasicAuth("basic@zatrano.test", "secret")
	req := http.NewRequest(raw)
	req.SetSession(&memSession{data: map[string]any{}})
	if !manager.OnceBasic(req) {
		t.Fatal("once basic failed")
	}
	if manager.User(req) == nil {
		t.Fatal("expected user from once basic")
	}
	if req.Session().Get("auth_user_id") != nil {
		t.Fatal("once basic must not write session")
	}
}

func TestTwoFactorSecretsAreEncrypted(t *testing.T) {
	hash, _ := hashing.Hash("secret")
	provider := newMemoryUserProvider()
	user, _ := provider.Create(map[string]any{"email": "enc@zatrano.test", "password": hash})
	manager := auth.NewManager("web")
	manager.Extend("web", auth.NewGuard("web", provider))
	crypt, err := encryption.New("zatrano-dev-key")
	if err != nil {
		t.Fatal(err)
	}
	manager.SetEncrypter(crypt)

	req := http.NewRequest(httptest.NewRequest(stdhttp.MethodGet, "/", nil))
	req.SetSession(&memSession{data: map[string]any{}})
	_ = manager.Login(req, user)
	secret, _, _, err := manager.EnableTwoFactor(user)
	if err != nil {
		t.Fatal(err)
	}
	stored := fmt.Sprint(user.(*auth.GenericUser).Get("two_factor_secret"))
	if stored == "" || stored == secret {
		t.Fatalf("expected encrypted storage, secret=%q stored=%q", secret, stored)
	}
	code, err := totp.Code(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.ConfirmTwoFactor(user, code); err != nil {
		t.Fatal(err)
	}
}
