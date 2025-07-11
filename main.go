package main

import (
	"flag"
	"log"
	"os"

	"github.com/zatrano/zatrano/app/database/migrations"
	"github.com/zatrano/zatrano/app/database/seeders"
	"github.com/zatrano/zatrano/bootstrap"
	"github.com/zatrano/zatrano/config"
	"gorm.io/gorm"
)

func main() {
	// 1. Komut satırı bayraklarını tanımla
	var runMigration = flag.Bool("migrate", false, "Run database migrations")
	var runSeeder = flag.Bool("seed", false, "Run database seeders")
	flag.Parse() // Komut satırı argümanlarını işle

	// 2. Uygulamayı ve temel servisleri başlat
	// Bu, her mod için ortaktır çünkü hem migration hem de seeder
	// veritabanı bağlantısına ihtiyaç duyar.
	app := bootstrap.New()
	// Sadece veritabanı gibi temel provider'ları kaydedelim.
	// Rota, handler gibi provider'ları kaydetmeye gerek yok.
	// Bu, config/modules.go'yu daha akıllı hale getirmeyi gerektirebilir.
	// Şimdilik hepsini kaydedelim.
	app.AddProviders(config.GetModules())
	app.RegisterProviders()

	// 3. Mod seçimi
	if *runMigration {
		// --- MIGRATION MODU ---
		log.Println("Starting migration process...")
		db, err := app.Kernel.Get("db.pool") // Ana bağlantı havuzunu al
		if err != nil {
			log.Fatalf("Could not get database pool for migration: %v", err)
		}

		if err := migrations.RunMigrations(db.(*gorm.DB)); err != nil {
			log.Fatalf("Migration failed.")
		}

		log.Println("Migration process finished.")
		os.Exit(0) // Başarıyla tamamlandı, programı sonlandır.

	} else if *runSeeder {
		// --- SEEDER MODU ---
		log.Println("Starting seeder process...")
		db, err := app.Kernel.Get("db.pool") // Ana bağlantı havuzunu al
		if err != nil {
			log.Fatalf("Could not get database pool for seeding: %v", err)
		}

		databaseSeeder := seeders.NewDatabaseSeeder()
		if err := databaseSeeder.Run(db.(*gorm.DB)); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}

		log.Println("Seeder process finished.")
		os.Exit(0) // Başarıyla tamamlandı, programı sonlandır.

	} else {
		// --- NORMAL WEB SUNUCUSU MODU ---
		log.Println("Starting web server...")
		app.Boot() // Middleware'leri ve Provider'ların Boot metodlarını çalıştır
		app.Run()  // Sunucuyu dinle
	}
}
