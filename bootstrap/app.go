package bootstrap

import (
	"fmt"
	"log"

	"github.com/zatrano/zatrano/config"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
	"github.com/zatrano/zatrano/internal/zatrano/module"
	"github.com/zatrano/zatrano/internal/zatrano/view"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
)

// Application, ZATRANO uygulamasının tüm bileşenlerini ve yaşam döngüsünü yönetir.
type Application struct {
	Fiber     *fiber.App
	Kernel    kernel.IKernel
	providers []module.Provider
}

// New, yeni bir ZATRANO uygulaması örneği oluşturur ve temel yapılandırmayı yapar.
func New() *Application {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using default or environment variables.")
	}

	appKernel := kernel.New()
	
	engine := view.NewEngine()
	fiberApp := fiber.New(fiber.Config{
		Views:     engine,
		AppName:   config.Get("APP_NAME", "Zatrano"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			// Hata sayfalarını daha sonra oluşturacağız.
			return c.Status(code).SendString(fmt.Sprintf("Error %d: %s", code, err.Error()))
		},
	})

	return &Application{
		Fiber:     fiberApp,
		Kernel:    appKernel,
		providers: make([]module.Provider, 0),
	}
}

// AddProviders, uygulamaya bir veya daha fazla servis sağlayıcı (modül) ekler.
func (app *Application) AddProviders(providers []module.Provider) {
	app.providers = append(app.providers, providers...)
}

// GetProviders, eklenmiş sağlayıcıların bir kopyasını döndürür.
// Bu, main.go'nun seeder gibi özel durumlar için döngüleri ayırmasına olanak tanır.
func (app *Application) GetProviders() []module.Provider {
	return app.providers
}

// RegisterProviders, tüm sağlayıcıların Register metodunu çalıştırır.
func (app *Application) RegisterProviders() {
	for _, provider := range app.providers {
		provider.Register(app.Kernel)
	}
}

// Boot, uygulamanın çekirdek servislerini ve sağlayıcılarını başlatır.
// Bu metod, seeder'lar için gerekli değildir, bu yüzden ayrılmıştır.
func (app *Application) Boot() {
	// 1. Framework'ün çekirdek servislerini ve middleware'lerini kur.
	app.bootCoreServices()
	
	// 2. Tüm sağlayıcıların `Boot` metodlarını çalıştır.
	for _, provider := range app.providers {
		provider.Boot(app.Kernel, app.Fiber)
	}
}

// Run, web sunucusunu başlatır.
// Bu metod, seeder modunda çağrılmaz.
func (app *Application) Run() {
	// Statik dosyalar için public klasörünü ayarla.
	app.Fiber.Static("/", "./public")

	// Hiçbir rota eşleşmezse 404 Not Found hatası döndüren bir middleware ekle.
	app.Fiber.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})
	
	// Web sunucusunu dinlemeye başla.
	port := config.Get("APP_PORT", "3000")
	log.Printf("🚀 Zatrano server is running on port %s", port)
	log.Printf("   Press CTRL+C to stop")
	log.Fatal(app.Fiber.Listen(":" + port))
}


// bootCoreServices, framework'ün kendisi için gerekli olan temel servisleri
// ve middleware'leri kaydeder ve başlatır.
func (app *Application) bootCoreServices() {
	app.Fiber.Use(recover.New())
	
	if config.Get("APP_ENV") == "development" {
		app.Fiber.Use(logger.New())
	}

	sessionStore := session.New()
	
	app.Kernel.RegisterSingleton("session.store", func(k kernel.IKernel) (interface{}, error) {
		return sessionStore, nil
	})
	
	app.Fiber.Use(sessionStore)
	app.Fiber.Use(func(c *fiber.Ctx) error {
		c.Locals("session_store", sessionStore)
		return c.Next()
	})

	app.Fiber.Use(csrf.New(csrf.Config{
		Session:    sessionStore,
		ContextKey: "csrf",
	}))

	app.Fiber.Use(func(c *fiber.Ctx) error {
		// Bu, AuditingPlugin'in o anki isteğin context'ine erişmesini sağlar.
		app.Kernel.Register("http.context", func(k kernel.IKernel) (interface{}, error) {
			return c, nil
		})
		return c.Next()
	})
}