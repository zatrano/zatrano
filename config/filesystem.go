package config

import "github.com/zatrano/zatrano/internal/zatrano/filesystem"

// FilesystemConfig, tüm dosya sistemi ayarlarını içerir.
type FilesystemConfig struct {
	// Varsayılan olarak kullanılacak disk. .env'den okunabilir.
	Default string

	// Mevcut tüm disklerin yapılandırması.
	Disks map[string]filesystem.DiskConfig
}

// GetFilesystemConfig, dosya sistemi yapılandırmasını döndürür.
func GetFilesystemConfig() FilesystemConfig {
	return FilesystemConfig{
		Default: Get("FILESYSTEM_DISK", "public"),

		Disks: map[string]filesystem.DiskConfig{

			// 'public' diski: Dışarıdan erişilebilen dosyalar için (avatarlar, resimler).
			// `public/storage` altına kaydeder ve `http://app.url/storage/...` üzerinden erişilir.
			"public": {
				"driver": "local",
				"root":   "./public/storage", // Fiziksel olarak kaydedileceği yer
				"url":    Get("APP_URL", "http://localhost:3000") + "/storage",
			},

			// 'local' diski: Dışarıdan erişilememesi gereken dosyalar için (loglar, özel dosyalar).
			"local": {
				"driver": "local",
				"root":   "./storage/app",
			},

			// 's3' diski: Amazon S3 için yapılandırma (gelecekte eklenebilir).
			"s3": {
				"driver": "s3",
				"key":    Get("AWS_ACCESS_KEY_ID"),
				"secret": Get("AWS_SECRET_ACCESS_KEY"),
				"region": Get("AWS_DEFAULT_REGION"),
				"bucket": Get("AWS_BUCKET"),
			},
		},
	}
}
