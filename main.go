package main

import (
	"github.com/zatrano/zatrano/bootstrap"
	"github.com/zatrano/zatrano/config"
)

func main() {
	// 1. Yeni bir uygulama örneği oluştur.
	// New() metodu, .env'i yükler, Fiber ve Kernel'i kurar.
	app := bootstrap.New()

	// 2. Uygulama modüllerini (sağlayıcıları) ekle.
	// config/modules.go'dan gelen listeyi kullanıyoruz.
	app.AddProviders(config.GetModules())
	
	// 3. Uygulamayı çalıştır.
	// Run() metodu, Register/Boot döngülerini yönetir ve sunucuyu başlatır.
	app.Run()
}