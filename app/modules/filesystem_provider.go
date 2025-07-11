package modules

import (
	"github.com/zatrano/zatrano/config"
	"github.com/zatrano/zatrano/internal/zatrano/filesystem"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"

	"github.com/gofiber/fiber/v2"
)

type FilesystemProvider struct{}

func (p *FilesystemProvider) Register(k kernel.IKernel) {
	k.RegisterSingleton("storage", func(kern kernel.IKernel) (interface{}, error) {
		// Konfigürasyon dosyasından ayarları al
		fsConfig := config.GetFilesystemConfig()

		// Yeni bir filesystem yöneticisi oluştur
		return filesystem.NewFilesystem(fsConfig)
	})
}

func (p *FilesystemProvider) Boot(k kernel.IKernel, app *fiber.App) {
	// 'public' diskinin dosyalarını sunmak için bir statik rota ekle.
	// Bu, `public/storage` klasörünü `/storage` URL'ine bağlar.
	app.Static("/storage", "./public/storage")
}
