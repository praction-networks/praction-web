package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/argon2"
)

func CreateUser(ctx context.Context, user models.Admin) error {

	// Hash the password using Argon2
	salt, err := utils.GenerateSalt()
	if err != nil {
		logger.Error(fmt.Sprintf("Error generating salt: %v", err))
		return fmt.Errorf("error generating salt: %w", err)
	}
	hashedPassword := hashPassword(user.Password, salt)

	saltstr := base64.RawStdEncoding.EncodeToString(salt)

	// Create a new User object
	userToAdd := models.Admin{
		ID:        primitive.NewObjectID(),
		Username:  user.Username,
		Password:  hashedPassword,
		Salt:      saltstr,
		Mobile:    user.Mobile,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	// Insert the user into the MongoDB database
	err = insertUserIntoDB(ctx, userToAdd)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert user into DB: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("User %s created successfully.", userToAdd.ID))
	return nil
}

// hashPassword hashes the user's password using Argon2
func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash) // Encode hash to Base64
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertUserIntoDB(ctx context.Context, user models.Admin) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("User")

	// Insert the user document
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		if IsDuplicateKeyError(err) {
			logger.Warn("Duplicate request for User creation", err)
			return err
		}
		// Handle other types of errors (internal server errors)
		return err
	}

	return nil
}

// Helper function to check if the error is a duplicate key error
func IsDuplicateKeyError(err error) bool {
	// MongoDB returns duplicate key errors with error code 11000
	if mongoErr, ok := err.(mongo.WriteException); ok {
		for _, writeErr := range mongoErr.WriteErrors {
			if writeErr.Code == 11000 { // Duplicate key error
				return true
			}
		}
	}
	return false
}
