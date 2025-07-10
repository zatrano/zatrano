package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/app/handlers"
	"github.com/zatrano/zatrano/internal/zatrano/kernel"
)

// RegisterPanelRoutes, sadece giriş yapmış kullanıcıların erişebileceği panel rotalarını kaydeder.
func RegisterPanelRoutes(router fiber.Router, k kernel.IKernel) {
	panelHandler, _ := k.Get("handler.panel")
	
	router.Get("/", panelHandler.(*handlers.PanelHandler).Dashboard)
	router.Get("/settings", panelHandler.(*handlers.PanelHandler).Settings)
}