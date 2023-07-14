package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ClientError(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"message": msg,
	})
}
