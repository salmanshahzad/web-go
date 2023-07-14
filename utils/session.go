package utils

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"salmanshahzad.com/web-go/database"
	"salmanshahzad.com/web-go/models"
)

const cookieName = "sessionId"

func CreateSession(c *fiber.Ctx, userId uint) error {
	sid := uuid.New().String()
	return createSession(c, sid, userId)
}

func DeleteSession(c *fiber.Ctx) error {
	if _, err := database.Rdb.Del(database.Ctx, c.Cookies(cookieName)).Result(); err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Expires:  time.Now(),
		Name:     cookieName,
	})
	return nil
}

func GetSession(c *fiber.Ctx, resave bool) (*models.User, error) {
	sid := c.Cookies(cookieName)
	if sid == "" {
		return nil, nil
	}

	userId, err := database.Rdb.Get(database.Ctx, sid).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	err = database.Db.First(user, "id = ?", userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if resave {
		if err := createSession(c, sid, user.ID); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func createSession(c *fiber.Ctx, sid string, userId uint) error {
	expiration := time.Until(time.Now().Add(7 * 24 * time.Hour))
	if _, err := database.Rdb.SetEx(database.Ctx, sid, userId, expiration).Result(); err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Expires:  time.Now().Add(expiration),
		Name:     cookieName,
		Value:    sid,
	})
	return nil
}
