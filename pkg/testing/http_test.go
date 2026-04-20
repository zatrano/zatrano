package testing

import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestHTTPClient_Get(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello"})
	})

	client := NewHTTPClient(app)
	resp := client.Get("/test")

	if resp.Status() != 200 {
		t.Errorf("expected status 200, got %d", resp.Status())
	}
	if !strings.Contains(resp.Body(), "hello") {
		t.Errorf("expected body to contain 'hello', got %s", resp.Body())
	}
}

func TestHTTPClient_WithToken(t *testing.T) {
	app := fiber.New()
	app.Get("/protected", func(c fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "Bearer test-token" {
			return c.JSON(fiber.Map{"authorized": true})
		}
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	})

	client := NewHTTPClient(app).WithToken("test-token")
	resp := client.Get("/protected")

	if resp.Status() != 200 {
		t.Errorf("expected status 200, got %d", resp.Status())
	}
	if !strings.Contains(resp.Body(), "authorized") {
		t.Errorf("expected body to contain 'authorized', got %s", resp.Body())
	}
}

func TestHTTPClient_Post(t *testing.T) {
	app := fiber.New()
	app.Post("/echo", func(c fiber.Ctx) error {
		body := c.Body()
		return c.SendString(string(body))
	})

	client := NewHTTPClient(app)
	resp := client.Post("/echo", fiber.Map{"key": "value"})

	if resp.Status() != 200 {
		t.Errorf("expected status 200, got %d", resp.Status())
	}
	if !strings.Contains(resp.Body(), "value") {
		t.Errorf("expected body to contain 'value', got %s", resp.Body())
	}
}
