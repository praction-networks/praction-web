package utils

import (
	"crypto/rand"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		logger.Fatal("Failed to hash password:", err)
	}
	return string(hashedPassword)
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16) // 16 bytes salt
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
