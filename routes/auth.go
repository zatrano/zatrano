package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/app/handlers"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
)

// RegisterAuthRoutes, kullanıcı kimlik doğrulama işlemlerini (login, logout, register) kaydeder.
func RegisterAuthRoutes(router fiber.Router, k kernel.IKernel) {
	authHandler, _ := k.Get("handler.auth")

	// Misafir middleware'i eklenebilir (giriş yapmışsa anasayfaya yönlendir)
	router.Get("/login", authHandler.(*handlers.AuthHandler).ShowLoginForm)
	router.Post("/login", authHandler.(*handlers.AuthHandler).Login)
	router.Post("/logout", authHandler.(*handlers.AuthHandler).Logout)
}