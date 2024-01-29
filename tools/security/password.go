package security

import (
	"golang.org/x/crypto/bcrypt"
)

var PasswordHashStrength = 10

// HashPassword generates a hash using the bcrypt.GenerateFromPassword.
func HashPassword(password string) (string, error) {
	// salt를 사용하여 bcrypt 해싱
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword compares the hash.
func ComparePassword(hash, password string) bool {
	if len(password) == 0 || len(hash) == 0 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
