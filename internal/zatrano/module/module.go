package module

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zatrano/zatrano/internal/zatrano/kernel" // Proje adını `zatrano-app` olarak kullandık
)

type Provider interface {
	Register(k kernel.IKernel)
	Boot(k kernel.IKernel, app *fiber.App)
}