package bootstrap

import (
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
// Bu fonksiyon, uygulama yaşam döngüsünde sadece bir kez, main.go tarafından çağrılır.
func New() *Application {
	// 1. .env dosyasını yükle. Hata varsa sadece uyar, uygulamayı durdurma.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using default or environment variables.")
	}

	// 2. Çekirdek bileşenleri oluştur: IoC Konteyneri ve Fiber App.
	appKernel := kernel.New()
	
	engine := view.NewEngine()
	fiberApp := fiber.New(fiber.Config{
		Views:     engine,
		AppName:   config.Get("APP_NAME", "Zatrano"),
		// Hata yönetimi için özel bir handler.
		// Bu, 404 gibi hataları veya handler'lardan dönen hataları yakalar.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			// Hata sayfasını render et.
			return view.Render(c, "pages/errors/default", fiber.Map{
				"Title":   fmt.Sprintf("Error %d", code),
				"Code":    code,
				"Message": err.Error(),
			}, "layouts/error")
		},
	})

	return &Application{
		Fiber:     fiberApp,
		Kernel:    appKernel,
		providers: make([]module.Provider, 0),
	}
}

// AddProvider, uygulamaya bir servis sağlayıcı (modül) ekler.
func (app *Application) AddProvider(provider module.Provider) {
	app.providers = append(app.providers, provider)
}

// AddProviders, birden fazla sağlayıcıyı tek seferde ekler.
func (app *Application) AddProviders(providers []module.GProvider) {
	app.providers = append(app.providers, providers...)
}

// Run, uygulamanın tüm yaşam döngüsünü başlatır. Bu, main.go'dan çağrılan son fonksiyondur.
func (app *Application) Run() {
	// 1. Framework'ün çekirdek servislerini ve middleware'lerini kur.
	app.bootCoreServices()

	// 2. Geliştiricinin config/modules.go'da tanımladığı tüm sağlayıcıların
	// `Register` metodlarını çalıştırarak servisleri konteynere kaydet.
	for _, provider := range app.providers {
		provider.Register(app.Kernel)
	}
	
	// 3. Tüm servisler kaydedildikten sonra, sağlayıcıların `Boot` metodlarını
	// çalıştırarak uygulamayı başlat (örn: rotaları tanımla).
	for _, provider := range app.providers {
		provider.Boot(app.Kernel, app.Fiber)
	}

	// 4. Statik dosyalar için public klasörünü ayarla.
	// Bu, /css/style.css gibi isteklere ./public/css/style.css dosyasını sunar.
	app.Fiber.Static("/", "./public")

	// 5. Hiçbir rota eşleşmezse 404 Not Found hatası döndüren bir middleware ekle.
	// Bu, her zaman en sonda olmalıdır.
	app.Fiber.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})
	
	// 6. Web sunucusunu dinlemeye başla.
	port := config.Get("APP_PORT", "3000")
	log.Printf("🚀 Zatrano server is running on port %s", port)
	log.Printf("   Press CTRL+C to stop")
	log.Fatal(app.Fiber.Listen(":" + port))
}

// bootCoreServices, framework'ün kendisi için gerekli olan temel servisleri
// ve middleware'leri kaydeder ve başlatır.
func (app *Application) bootCoreServices() {
	// Temel Middleware'ler
	app.Fiber.Use(recover.New()) // Panikleri yakalayıp 500 hatasına çevirir.
	
	// Sadece geliştirme ortamında detaylı loglama yap.
	if config.Get("APP_ENV") == "development" {
		app.Fiber.Use(logger.New())
	}

	// Session'ı kur
	sessionStore := session.New(session.Config{
		// Production için bu ayarlar .env'den okunmalıdır.
		// CookieSecure:  config.GetBool("SESSION_SECURE", true),
		// CookieHTTPOnly: true,
		// Expiration: time.Hour * 24,
	})
	
	// Session store'u hem IoC konteynerine kaydet hem de context'e koy.
	// Bu, uygulamanın her yerinden erişilebilir olmasını sağlar.
	app.Kernel.RegisterSingleton("session.store", func(k kernel.IKernel) (interface{}, error) {
		return sessionStore, nil
	})
	
	// Middleware olarak session'ı Fiber'a ekle.
	// Bu, her istekte session'ı otomatik olarak başlatır.
	// Ayrıca, session'ı context'e de ekleriz ki diğer middleware'ler ve handler'lar
	// bu session'a erişebilsin.
	// Bu, örneğin flash mesajlar veya form verileri için kullanılabilir.
	// Fiber'ın middleware zincirine session'ı eklerken, session'ı
	// IoC konteynerine de kaydediyoruz ki diğer bileşenler bu session'a erişebilsin.
	app.Fiber.Use(sessionStore)
	app.Fiber.Use(func(c *fiber.Ctx) error {
		c.Locals("session_store", sessionStore)
		return c.Next()
	})

	// CSRF'yi kur (Session'dan sonra gelmelidir!)
	app.Fiber.Use(csrf.New(csrf.Config{
		Session:    sessionStore,
		ContextKey: "csrf", // Bu anahtar ile `c.Locals("csrf")` üzerinden token'a erişilir.
		// TokenLookup: "header:X-CSRF-Token,form:_csrf",
	}))

	// Otomatik denetim (auditing) için her isteğin context'ini
	// IoC konteynerine kaydeden bir middleware.
	// NOT: Bu yaklaşım, her isteğin kendi "scoped" kernel'ine sahip olduğu
	// daha gelişmiş bir yapıyla daha güvenli hale getirilebilir.
	// Şimdilik, her isteğin DB oturumunun doğru context'i almasını sağlar.
	app.Fiber.Use(func(c *fiber.Ctx) error {
		app.Kernel.Register("http.context", func(k kernel.IKernel) (interface{}, error) {
			return c, nil
		})
		return c.Next()
	})
}