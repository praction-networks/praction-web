package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
)

// GenerateJWT generates a JWT token for the authenticated user
func GenerateJWT(user *models.User) (string, error) {

	secret := GetJWTSECRET()

	if secret == "" {
		logger.Fatal("Failed to Retive JWT SECRET")
	}
	// Define JWT claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(90 * 24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key (e.g., "mysecretkey")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
