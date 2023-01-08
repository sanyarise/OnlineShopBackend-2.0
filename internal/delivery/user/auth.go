package user

import (
	"OnlineShopBackend/internal/models"
	"crypto/sha1"
	"fmt"
	"unicode"
)

const salt = "sjdhkashdsw823rgfeg"

func ValidationCheck(user models.User) error {
	if user.Email == "" && user.Firstname == "" && user.Lastname == "" {
		return fmt.Errorf("empty filed")
	}
	if len(user.Password) < 5 {
		return fmt.Errorf("password is too short")
	}
	for _, char := range user.Password {
		if !unicode.IsDigit(char) && !unicode.Is(unicode.Latin, char) {
			return fmt.Errorf("password should contain lathin letter or numbers only")
		}
	}
	return nil
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

