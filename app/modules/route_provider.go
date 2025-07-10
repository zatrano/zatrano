package modules

import (
	"log"
	"strings"
	"time"

	"github.com/zatrano/zatrano/app/routes"
	"github.com/zatrano/zatrano/config"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
	"github.com/zatrano/zatrano/internal/zatrano/module"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type RouteProvider struct{}

func (p *RouteProvider) Register(k kernel.IKernel) {
	// --- CSRF MIDDLEWARE ---
	k.RegisterSingleton("middleware.csrf", func(kern kernel.IKernel) (interface{}, error) {
		sessionStore, _ := kern.Get("session.store")
		return csrf.New(csrf.Config{
			Session:    sessionStore.(*session.Store),
			ContextKey: "csrf",
			Next: func(c *fiber.Ctx) bool {
				// /api/ ile başlayan veya /payment/callback olan yolları CSRF'den muaf tut.
				return strings.HasPrefix(c.Path(), "/api") ||
				       strings.HasPrefix(c.Path(), "/payment/callback")
			},
		}), nil
	})

	// --- RATE LIMITER MIDDLEWARE'LERİ ---

	// 1. API için Rate Limiter (Daha Yüksek Limit)
	k.RegisterSingleton("middleware.limiter.api", func(kern kernel.IKernel) (interface{}, error) {
		return limiter.New(limiter.Config{
			Max:        config.GetInt("API_RATE_LIMIT", 60), // Dakikada 60 istek
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too Many Requests"})
			},
		}), nil
	})

	// 2. WEB Formları için Genel Rate Limiter (Daha Düşük Limit)
	k.RegisterSingleton("middleware.limiter.web", func(kern kernel.IKernel) (interface{}, error) {
		return limiter.New(limiter.Config{
			Max:        config.GetInt("WEB_FORM_RATE_LIMIT", 15), // Dakikada 15 form gönderimi
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
			LimitReached: func(c *fiber.Ctx) error {
				// Form gönderimi olduğu için JSON yerine bir hata mesajı flash'leyip geri yönlendirmek daha iyi.
				// flash.Error(c, "You are submitting forms too frequently. Please try again later.")
				// return c.RedirectBack("/")
				// Şimdilik basit bir hata döndürelim:
				return c.Status(fiber.StatusTooManyRequests).SendString("Too many form submissions. Please try again in a minute.")
			},
			// Sadece POST, PUT, DELETE, PATCH isteklerini sınırla. GET istekleri serbest.
			Next: func(c *fiber.Ctx) bool {
				return c.Method() == "GET"
			},
		}), nil
	})

	// 3. Hassas İşlemler için Rate Limiter (Çok Düşük Limit)
	k.RegisterSingleton("middleware.limiter.auth", func(kern kernel.IKernel) (interface{}, error) {
		return limiter.New(limiter.Config{
			Max:        config.GetInt("AUTH_RATE_LIMIT", 5), // 15 dakikada 5 deneme
			Expiration: 15 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).SendString("Too many login attempts. Please try again later.")
			},
		}), nil
	})
}

func (p *RouteProvider) Boot(k kernel.IKernel, app *fiber.App) {
	log.Println("Booting Route Provider and applying middleware groups...")

	// Middleware'leri IoC'den al
	csrfMiddleware, _ := k.Get("middleware.csrf")
	apiLimiter, _ := k.Get("middleware.limiter.api")
	webLimiter, _ := k.Get("middleware.limiter.web")
	authLimiter, _ := k.Get("middleware.limiter.auth")

	// ==================================================================
	// GRUP 1: WEB ROTALARI
	// Middleware Zinciri: [Rate Limiter (Web)] -> [CSRF]
	// ==================================================================
	webRouter := app.Group("/", webLimiter.(fiber.Handler), csrfMiddleware.(fiber.Handler))
	
	routes.RegisterWebRoutes(webRouter, k)
	// Login gibi hassas rotalar için AYRI bir grup oluşturuyoruz.
	
	// ==================================================================
	// GRUP 2: AUTH ROTALARI (Giriş, Kayıt, Şifre Sıfırlama)
	// Middleware Zinciri: [Rate Limiter (Web)] -> [Rate Limiter (Auth)] -> [CSRF]
	// ==================================================================
	// webRouter'dan türetildiği için zaten webLimiter ve CSRF'ye sahip.
	// Ek olarak daha sıkı bir rate limit uyguluyoruz.
	authRouter := webRouter.With(authLimiter.(fiber.Handler))
	routes.RegisterAuthRoutes(authRouter, k)

	// ==================================================================
	// GRUP 3: API ROTALARI
	// Middleware Zinciri: [Rate Limiter (API)]
	// ==================================================================
	// Bu grup, doğrudan `app`'ten türetildiği için CSRF veya webLimiter'a sahip değil.
	apiRouter := app.Group("/api", apiLimiter.(fiber.Handler))
	routes.RegisterApiRoutes(apiRouter, k)
	
	// ... (CSRF'siz web rotaları grubu burada aynı kalabilir) ...
}