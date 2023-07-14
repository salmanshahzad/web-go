package api

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"salmanshahzad.com/web-go/database"
	"salmanshahzad.com/web-go/models"
	"salmanshahzad.com/web-go/utils"
)

func NewSessionRouter() *fiber.App {
	r := fiber.New()
	r.Post("/", handleSignIn)
	r.Delete("/", handleSignOut)
	return r
}

func handleSignIn(c *fiber.Ctx) error {
	type Payload struct {
		Username string
		Password string
	}
	payload := new(Payload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	payload.Username = strings.TrimSpace(payload.Username)
	payload.Password = strings.TrimSpace(payload.Password)

	if len(payload.Username) == 0 || len(payload.Password) == 0 {
		return utils.ClientError(c, "Username and password are required")
	}

	user := new(models.User)
	if err := database.Db.First(user, "username = ?", payload.Username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Incorrect username or password",
			})
		}
		return err
	}

	validPassword, err := user.VerifyPassword(payload.Password)
	if err != nil {
		return err
	}
	if !validPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect username or password",
		})
	}

	if err := utils.CreateSession(c, user.ID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func handleSignOut(c *fiber.Ctx) error {
	if err := utils.DeleteSession(c); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
