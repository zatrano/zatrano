package modules

import (
	"fmt"
	"log"
	"time"

	"github.com/zatrano/zatrano/config"
	"github.com/zatrano/zatrano/internal/zatrano/database"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseProvider, veritabanı bağlantısını ve oturumlarını yönetir.
// Tek sorumluluğu, IoC konteynerine doğru veritabanı servislerini kaydetmektir.
type DatabaseProvider struct{}

// Register, IoC konteynerine iki farklı veritabanı servisi kaydeder:
// 1. "db.pool": Uygulama boyunca yaşayacak olan ana veritabanı bağlantı havuzu (Singleton).
// 2. "db": Her HTTP isteği için oluşturulan, o isteğe özel GORM oturumu (Transient).
func (p *DatabaseProvider) Register(k kernel.IKernel) {

	// 1. Ana veritabanı bağlantı havuzunu (DB Pool) singleton olarak kaydet.
	// Bu, sadece uygulama ilk başladığında bir kez çalışır.
	k.RegisterSingleton("db.pool", func(kern kernel.IKernel) (interface{}, error) {
		host := config.Get("DB_HOST", "127.0.0.1")
		port := config.GetInt("DB_PORT", 5432)
		user := config.Get("DB_USERNAME", "postgres")
		password := config.Get("DB_PASSWORD", "password")
		dbname := config.Get("DB_DATABASE", "zatrano")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)

		logLevel := logger.Silent
		if config.Get("APP_ENV") == "development" {
			logLevel = logger.Info
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})

		if err != nil {
			return nil, fmt.Errorf("FATAL: failed to connect to database pool: %w", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("FATAL: failed to get underlying sql.DB: %w", err)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Println("Database connection pool established successfully.")
		return db, nil
	})

	// 2. Her istek için yeni bir GORM oturumu oluşturan "transient" bir fabrika kaydet.
	k.Register("db", func(kern kernel.IKernel) (interface{}, error) {
		pool, err := kern.Get("db.pool")
		if err != nil {
			return nil, err
		}
		db := pool.(*gorm.DB)

		ctx, ok := kern.Get("http.context")
		if !ok || ctx == nil {
			// HTTP context'i yoksa (CLI komutu gibi), eklentisiz devam et.
			return db.Session(&gorm.Session{}), nil
		}

		// HTTP context'i varsa, Auditing eklentisini tak.
		tx := db.Session(&gorm.Session{})
		tx.Use(&database.AuditingPlugin{Context: ctx.(*fiber.Ctx)})

		return tx, nil
	})
}

// Boot metodu artık hiçbir işlem yapmıyor.
// Arayüze uyum sağlamak için mevcuttur, ancak içi bilinçli olarak boştur.
// Migration işlemleri artık sadece `go run main.go --migrate` komutuyla, merkezi
// olarak `app/database/migrations` paketinden çalıştırılır.
func (p *DatabaseProvider) Boot(k kernel.IKernel, app *fiber.App) {
	// BOŞ BIRAKILDI
}
