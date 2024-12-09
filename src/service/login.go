package service

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/argon2"
)

// AuthenticateUser checks the username and password and returns the user if valid
func AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("User")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		// Return a nil pointer to User, and the error
		return nil, errors.New("user not found")
	}

	// Validate the password
	if !validatePassword(password, user.Password, user.Salt) {
		// Return a nil pointer to User, and the error
		return nil, errors.New("invalid password")
	}

	// Return a pointer to the user struct
	return &user, nil
}

// validatePassword checks if the provided password matches the stored hash
func validatePassword(password, hashedPassword, salt string) bool {
	// Decode the salt from base64
	saltBytes, _ := base64.RawStdEncoding.DecodeString(salt)
	// Hash the password using Argon2
	hash := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	// Compare the hash to the stored hashed password

	return base64.RawStdEncoding.EncodeToString(hash) == hashedPassword
}
