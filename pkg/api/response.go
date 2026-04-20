package api

import "github.com/gofiber/fiber/v3"

// Envelope wraps API responses in a standard JSON structure.
type Envelope struct {
	Data  any `json:"data"`
	Meta  any `json:"meta,omitempty"`
	Links any `json:"links,omitempty"`
}

// Wrap constructs a response envelope.
func Wrap(data any, meta any, links any) Envelope {
	return Envelope{Data: data, Meta: meta, Links: links}
}

// JSON writes a standardized API response.
func JSON(c fiber.Ctx, status int, data any, meta any, links any) error {
	return c.Status(status).JSON(Wrap(data, meta, links))
}

// Success writes a 200 OK standardized response.
func Success(c fiber.Ctx, data any) error {
	return JSON(c, fiber.StatusOK, data, nil, nil)
}

// Created writes a 201 Created standardized response.
func Created(c fiber.Ctx, data any) error {
	return JSON(c, fiber.StatusCreated, data, nil, nil)
}

// Transform slices of models into response resources.
func Transform[T any](items []T, fn func(T) any) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	return out
}
