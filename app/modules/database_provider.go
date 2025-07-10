package modules

import (
	"fmt"
	"log"
	"time"
	
	"github.com/zatrano/zatrano/app/models" // Uygulamanın modelleri
	"github.com/zatrano/zatrano/config"
	"github.com/zatrano/zatrano/internal/zatrano/database" // Framework'ün db eklentileri
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
	"github.com/zatrano/zatrano/internal/zatrano/module"
	
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseProvider, veritabanı bağlantısını, oturumlarını ve migration'ları yönetir.
type DatabaseProvider struct{}

// Register, IoC konteynerine iki farklı veritabanı servisi kaydeder:
// 1. "db.pool": Uygulama boyunca yaşayacak olan ana veritabanı bağlantı havuzu (Singleton).
// 2. "db": Her HTTP isteği için oluşturulan, o isteğe özel GORM oturumu (Transient).
func (p *DatabaseProvider) Register(k kernel.IKernel) {

	// 1. Ana veritabanı bağlantı havuzunu (DB Pool) singleton olarak kaydet.
	// Bu, sadece uygulama ilk başladığında bir kez çalışır.
	k.RegisterSingleton("db.pool", func(kern kernel.IKernel) (interface{}, error) {
		// Veritabanı bağlantı bilgilerini .env'den (config aracılığıyla) al
		host := config.Get("DB_HOST", "127.0.0.1")
		port := config.GetInt("DB_PORT", 5432)
		user := config.Get("DB_USERNAME", "postgres")
		password := config.Get("DB_PASSWORD", "password")
		dbname := config.Get("DB_DATABASE", "zatrano")
		
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)

		logLevel := logger.Silent
		if config.Get("APP_ENV") == "development" {
			// Geliştirme ortamında çalışan tüm SQL sorgularını logla.
			logLevel = logger.Info
		}

		// GORM ile veritabanına bağlan
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})

		if err != nil {
			// Bağlantı hatası durumunda uygulama başlatılamaz.
			return nil, fmt.Errorf("FATAL: failed to connect to database pool: %w", err)
		}
		
		// Bağlantı havuzu ayarları
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
	// Bu, her `kernel.Get("db")` çağrıldığında çalışır.
	k.Register("db", func(kern kernel.IKernel) (interface{}, error) {
		// Ana bağlantı havuzunu konteynerden al.
		pool, err := kern.Get("db.pool")
		if err != nil {
			return nil, err
		}
		
		db := pool.(*gorm.DB)
		
		// O anki isteğin context'ini al. Bu, bir middleware tarafından sağlanmalı.
		ctx, ok := kern.Get("http.context")
		if !ok || ctx == nil {
			// Eğer bir HTTP context'i yoksa (örn: seeder veya CLI komutu çalışıyorsa),
			// eklentisiz, saf bir GORM oturumu döndür.
			return db.Session(&gorm.Session{}), nil
		}

		// Eğer bir HTTP context'i varsa, o context'e özel yeni bir GORM oturumu başlat
		// ve bu oturuma "Auditing" eklentisini tak.
		tx := db.Session(&gorm.Session{})
		tx.Use(&database.AuditingPlugin{Context: ctx.(*fiber.Ctx)})
		
		return tx, nil
	})
}

// Boot, veritabanı migration'larını çalıştırmak için kullanılır.
func (p *DatabaseProvider) Boot(k kernel.IKernel, app *fiber.App) {
	// Sadece geliştirme ortamında otomatik migration yap.
	if config.Get("APP_ENV") != "development" {
		return
	}

	// Migration için ana bağlantı havuzunu kullan.
	db, err := k.Get("db.pool")
	if err != nil {
		log.Fatalf("Migration failed: could not get db.pool from kernel: %v", err)
	}

	log.Println("Running auto-migrations...")
	
	// Uygulamadaki tüm modelleri buraya ekle.
	// GORM, tabloları oluşturacak veya güncelleyecektir.
	err = db.(*gorm.DB).AutoMigrate(
		&models.User{},
		// &models.Address{},
		// &models.Country{},
		// &models.Session{},
	)
	
	if err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}
	
	log.Println("Auto-migrations completed successfully.")
}