package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/salmanshahzad/web-go/utils"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetSession(c, true)
		if err != nil {
			return err
		}
		if user == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Locals("user", user)
		return c.Next()
	}
}
