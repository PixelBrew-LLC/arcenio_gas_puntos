package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var numericPINRegex = regexp.MustCompile(`^\d+$`)

// HashPassword hashea un PIN/password con bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifica si el PIN coincide con su hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsValidPIN verifica que el PIN contenga solo números
func IsValidPIN(pin string) bool {
	return len(pin) >= 4 && numericPINRegex.MatchString(pin)
}
