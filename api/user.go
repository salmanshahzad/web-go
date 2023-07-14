package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"salmanshahzad.com/web-go/database"
	"salmanshahzad.com/web-go/middleware"
	"salmanshahzad.com/web-go/models"
	"salmanshahzad.com/web-go/utils"
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
	user := c.Locals("user").(*models.User)
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
		return utils.ClientError(c, "Username and password are required")
	}

	var userCount int64
	if err := database.Db.Model(&models.User{}).Where("username = ?", payload.Username).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return utils.ClientError(c, "Username already exists")
	}

	user := models.User{
		Username: payload.Username,
		Password: payload.Password,
	}
	if err := user.HashPassword(); err != nil {
		return err
	}
	if err := database.Db.Create(&user).Error; err != nil {
		return err
	}

	if err := utils.CreateSession(c, user.ID); err != nil {
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
		return utils.ClientError(c, "Username is required")
	}

	var userCount int64
	if err := database.Db.Model(&models.User{}).Where("username = ?", payload.Username).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return utils.ClientError(c, "Username already exists")
	}

	user := c.Locals("user").(*models.User)
	user.Username = payload.Username
	if err := database.Db.Save(user).Error; err != nil {
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
		return utils.ClientError(c, "Password is required")
	}

	user := c.Locals("user").(*models.User)
	user.Password = payload.Password
	if err := user.HashPassword(); err != nil {
		return err
	}
	if err := database.Db.Save(user).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func handleDeleteUser(c *fiber.Ctx) error {
	if err := utils.DeleteSession(c); err != nil {
		return err
	}
	user := c.Locals("user").(*models.User)
	if err := database.Db.Delete(user).Error; err != nil {
		return err
	}

	log.Println("Deleted user", user.Username)
	return c.SendStatus(fiber.StatusNoContent)
}
