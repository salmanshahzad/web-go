package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/salmanshahzad/web-go/database"
	"github.com/salmanshahzad/web-go/middleware"
	"github.com/salmanshahzad/web-go/utils"
)

func NewUserRouter() *fiber.App {
	r := fiber.New()
	r.Post("/", handleCreateUser)
	auth := r.Group("/", middleware.Auth())
	auth.Get("/", handleGetUser)
	auth.Put("/username", handleEditUsername)
	auth.Put("/password", handleEditPassword)
	auth.Delete("/", handleDeleteUser)
	return r
}

func handleGetUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
	})
}

func handleCreateUser(c *fiber.Ctx) error {
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
		return utils.UnprocessableEntity(c, "Username and password are required")
	}

	userCount, err := database.Db.CountUsersWithUsername(database.Ctx, payload.Username)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return utils.UnprocessableEntity(c, "Username already exists")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	user := database.CreateUserParams{
		Username: payload.Username,
		Password: hashedPassword,
	}
	userId, err := database.Db.CreateUser(database.Ctx, user)
	if err != nil {
		return err
	}

	if err := utils.CreateSession(c, userId); err != nil {
		return err
	}
	log.Println("Created user", payload.Username)
	return c.SendStatus(fiber.StatusCreated)
}

func handleEditUsername(c *fiber.Ctx) error {
	type Payload struct {
		Username string
	}
	payload := new(Payload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	payload.Username = strings.TrimSpace(payload.Username)
	if len(payload.Username) == 0 {
		return utils.UnprocessableEntity(c, "Username is required")
	}

	userCount, err := database.Db.CountUsersWithUsername(database.Ctx, payload.Username)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return utils.UnprocessableEntity(c, "Username already exists")
	}

	user := c.Locals("user").(*database.User)
	params := database.UpdateUsernameParams{
		ID:       user.ID,
		Username: payload.Username,
	}
	if err := database.Db.UpdateUsername(database.Ctx, params); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func handleEditPassword(c *fiber.Ctx) error {
	type Payload struct {
		Password string
	}
	payload := new(Payload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	payload.Password = strings.TrimSpace(payload.Password)
	if len(payload.Password) == 0 {
		return utils.UnprocessableEntity(c, "Password is required")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	user := c.Locals("user").(*database.User)
	params := database.UpdatePasswordParams{
		ID:       user.ID,
		Password: hashedPassword,
	}
	if err := database.Db.UpdatePassword(database.Ctx, params); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func handleDeleteUser(c *fiber.Ctx) error {
	if err := utils.DeleteSession(c); err != nil {
		return err
	}
	user := c.Locals("user").(*database.User)
	if err := database.Db.DeleteUser(database.Ctx, user.ID); err != nil {
		return err
	}

	log.Println("Deleted user", user.Username)
	return c.SendStatus(fiber.StatusNoContent)
}
