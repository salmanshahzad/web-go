package utils

import "github.com/matthewhartstonge/argon2"

var argon = argon2.DefaultConfig()

func HashPassword(password string) (string, error) {
	hashed, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(password, hashed string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(hashed))
}
