package modules

import (
	"log"
	"strings"
	"time"

	// --- DOĞRU IMPORT YOLLARI ---
	// Projenin kendi içindeki paketlere, go.mod'da tanımlanan tam yolla erişiyoruz.
	"github.com/zatrano/zatrano/routes"
	"github.com/zatrano/zatrano/config"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
	"github.com/zatrano/zatrano/internal/zatrano/module"
	
	// --- HARİCİ PAKETLER ---
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// RouteProvider, tüm uygulama rotalarını kaydeder ve yapılandırır.
type RouteProvider struct{}

// Register, rota katmanının ihtiyaç duyacağı middleware'leri IoC konteynerine kaydeder.
func (p *RouteProvider) Register(k kernel.IKernel) {
	// --- CSRF MIDDLEWARE ---
	k.RegisterSingleton("middleware.csrf", func(kern kernel.IKernel) (interface{}, error) {
		sessionStore, err := kern.Get("session.store")
		if err != nil {
			return nil, err
		}
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

	// API için Rate Limiter
	k.RegisterSingleton("middleware.limiter.api", func(kern kernel.IKernel) (interface{}, error) {
		return limiter.New(limiter.Config{
			Max:        config.GetInt("API_RATE_LIMIT", 60),
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too Many Requests"})
			},
		}), nil
	})

	// WEB Formları için Genel Rate Limiter
	k.RegisterSingleton("middleware.limiter.web", func(kern kernel.IKernel) (interface{}, error) {
		return limiter.New(limiter.Config{
			Max:        config.GetInt("WEB_FORM_RATE_LIMIT", 15),
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).SendString("Too many form submissions. Please try again in a minute.")
			},
			Next: func(c *fiber.Ctx) bool {
				// Sadece state değiştiren metodları sınırla, GET'leri serbest bırak.
				switch c.Method() {
				case "POST", "PUT", "DELETE", "PATCH":
					return false // Limiter çalışsın
				default:
					return true // Limiter'ı atla
				}
			},
		}), nil
	})

	// Hassas İşlemler için Rate Limiter (Login, vb.)
	k.RegisterSingleton("middleware.limiter.auth", func(kern kernel.IKernel) (interface{}, error) {
		return limiter.New(limiter.Config{
			Max:        config.GetInt("AUTH_RATE_LIMIT", 5),
			Expiration: 15 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).SendString("Too many login attempts. Please try again later.")
			},
		}), nil
	})
	
	// NOT: Diğer middleware'ler (auth, guest) de burada kaydedilmelidir.
	// k.Register("middleware.auth", ...)
}

// Boot, rota gruplarını oluşturur, middleware'leri atar ve rota dosyalarını çağırır.
func (p *RouteProvider) Boot(k kernel.IKernel, app *fiber.App) {
	log.Println("Booting Route Provider and applying middleware groups...")

	// Middleware'leri IoC'den al
	csrfMiddleware, _ := k.Get("middleware.csrf")
	apiLimiter, _ := k.Get("middleware.limiter.api")
	webLimiter, _ := k.Get("middleware.limiter.web")
	authLimiter, _ := k.Get("middleware.limiter.auth")

	// ==================================================================
	// GRUP 1: API ROTALARI
	// Middleware Zinciri: [Rate Limiter (API)]
	// ==================================================================
	apiRouter := app.Group("/api", apiLimiter.(fiber.Handler))
	routes.RegisterApiRoutes(apiRouter, k)

	// ==================================================================
	// GRUP 2: WEB ROTALARI (Genel)
	// Middleware Zinciri: [Rate Limiter (Web)] -> [CSRF]
	// ==================================================================
	webRouter := app.Group("/", webLimiter.(fiber.Handler), csrfMiddleware.(fiber.Handler))
	
	routes.RegisterWebRoutes(webRouter, k)
	
	// ==================================================================
	// GRUP 3: AUTH ROTALARI (Daha Sıkı Kurallar)
	// Middleware Zinciri: (webRouter'dan miras) [Rate Limiter (Web)] -> [CSRF] -> [Rate Limiter (Auth)]
	// ==================================================================
	// `webRouter.Group` yerine `webRouter.With` kullanmak, aynı ön ekte (`/`) kalıp
	// sadece ek middleware uygulamamızı sağlar.
	authRouter := webRouter.With(authLimiter.(fiber.Handler))
	routes.RegisterAuthRoutes(authRouter, k)

	// ==================================================================
	// GRUP 4: PANEL ROTALARI (Kimlik Doğrulama Gerekli)
	// Middleware Zinciri: (webRouter'dan miras) [Rate Limiter (Web)] -> [CSRF] -> [Auth Middleware]
	// ==================================================================
	// authMiddleware, _ := k.Get("middleware.auth")
	// panelRouter := webRouter.Group("/panel", authMiddleware.(fiber.Handler))
	// routes.RegisterPanelRoutes(panelRouter, k)

	// ==================================================================
	// GRUP 5: CSRF'SİZ WEB ROTALARI (Webhook'lar vb.)
	// ==================================================================
	// Doğrudan ana `app` üzerinden tanımlandığı için hiçbir grup middleware'ini almaz.
	app.Post("/payment/callback", func(c *fiber.Ctx) error {
		log.Println("Received payment callback:", string(c.Body()))
		return c.SendStatus(200)
	})

	log.Println("All routes have been registered.")
}