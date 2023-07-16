package api

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/salmanshahzad/web-go/database"
	"github.com/salmanshahzad/web-go/utils"
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

	user, err := database.Db.GetUserByUsername(database.Ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Incorrect username or password",
			})
		}
		return err
	}

	validPassword, err := utils.VerifyPassword(payload.Password, user.Password)
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
