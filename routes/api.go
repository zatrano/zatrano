package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
)

// RegisterApiRoutes, tüm API endpoint'lerini kaydeder.
func RegisterApiRoutes(router fiber.Router, k kernel.IKernel) {
	// Örnek bir public API endpoint'i
	router.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "API is running",
			"time":   time.Now(),
		})
	})
	
	// Örnek bir veri endpoint'i
	router.Get("/users", func(c *fiber.Ctx) error {
		// userHandler, _ := k.Get("handler.user.api")
		// return userHandler.GetAll(c)
		
		// Şimdilik sahte veri döndürelim
		return c.JSON(fiber.Map{
			"data": []fiber.Map{
				{"id": 1, "name": "API User 1"},
				{"id": 2, "name": "API User 2"},
			},
		})
	})
}