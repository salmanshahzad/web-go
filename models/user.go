package models

import (
	"github.com/matthewhartstonge/argon2"
)

var argon = argon2.DefaultConfig()

type User struct {
	ID       uint
	Username string `gorm:"uniqueIndex"`
	Password string
}

func (u *User) HashPassword() error {
	encoded, err := argon.HashEncoded([]byte(u.Password))
	if err != nil {
		return err
	}
	u.Password = string(encoded)
	return nil
}

func (u *User) VerifyPassword(password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(u.Password))
}
