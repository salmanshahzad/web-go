package utils

import (
	"github.com/gofiber/fiber/v2"
)

func Unauthorized(c *fiber.Ctx, msg string) error {
	return sendMessage(c, fiber.StatusUnauthorized, msg)
}

func UnprocessableEntity(c *fiber.Ctx, msg string) error {
	return sendMessage(c, fiber.StatusUnprocessableEntity, msg)
}

func sendMessage(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"message": msg,
	})
}
