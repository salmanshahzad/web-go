package api

import (
	"github.com/gofiber/fiber/v2"
)

func NewHealthRouter() *fiber.App {
	r := fiber.New()
	r.Get("/", handleHealthCheck)
	return r
}

func handleHealthCheck(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
