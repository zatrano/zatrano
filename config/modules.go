package config

import (
	"github.com/zatrano/zatrano/app/modules"
	"github.com/zatrano/zatrano/internal/zatrano/module"
)

// GetModules, uygulamada yüklenecek tüm servis sağlayıcıların listesini döndürür.
// Modül ekleme/çıkarma işlemleri sadece bu dosyadan yapılır.
// Sıralama önemlidir! Bir modül diğerine bağımlıysa, önce bağımlı olunan yüklenmelidir.
func GetModules() []module.Provider {
	return []module.Provider{
		// İlk olarak veritabanı gibi temel modüller
		&modules.DatabaseProvider{},
		
		// Sonra bu temelleri kullanan modüller
		&modules.RepositoryProvider{},
		&modules.ServiceProvider{},
		&modules.HandlerProvider{},
		
		// Rota sağlayıcısı genellikle en sonda yer alır,
		// çünkü handler'ların kaydedilmiş olmasını bekler.
		&modules.RouteProvider{},
	}
}