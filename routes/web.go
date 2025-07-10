package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/app/handlers"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
)

// RegisterWebRoutes, genel, kimlik doğrulaması gerektirmeyen web rotalarını kaydeder.
func RegisterWebRoutes(router fiber.Router, k kernel.IKernel) {
	homeHandler, _ := k.Get("handler.home")
	
	router.Get("/", homeHandler.(*handlers.HomeHandler).Index)
	router.Get("/about", homeHandler.(*handlers.HomeHandler).About)
}