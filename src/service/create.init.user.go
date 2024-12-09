package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

// Initialize a validator instance
var validate = validator.New()

func CreateUser(ctx context.Context, username, password, email, firstName, lastName, role string) error {
	// Create a temporary User object to validate the input
	tempUser := models.User{
		Username:  username,
		Password:  password,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}

	// Validate the user
	if err := validate.Struct(tempUser); err != nil {
		logger.Error(fmt.Sprintf("Validation failed: %v", err))
		return fmt.Errorf("validation error: %w", err)
	}

	// Hash the password using Argon2
	salt, err := utils.GenerateSalt()
	if err != nil {
		logger.Error(fmt.Sprintf("Error generating salt: %v", err))
		return fmt.Errorf("error generating salt: %w", err)
	}
	hashedPassword := hashPassword(password, salt)

	saltstr := base64.RawStdEncoding.EncodeToString(salt)

	// Create a new User object
	user := models.User{
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		Salt:      saltstr,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}

	// Insert the user into the MongoDB database
	err = insertUserIntoDB(ctx, user)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to insert user into DB: %v", err))
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info(fmt.Sprintf("User %s created successfully.", username))
	return nil
}

// hashPassword hashes the user's password using Argon2
func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash) // Encode hash to Base64
}

// insertUserIntoDB inserts the new user into the MongoDB database
func insertUserIntoDB(ctx context.Context, user models.User) error {
	// Get the MongoDB client from the database package
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("User")

	// Insert the user document into the collection
	_, err := collection.InsertOne(ctx, user)
	if err != nil {

		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range mongoErr.WriteErrors {
				if writeErr.Code == 11000 {
					// Log the duplicate key error
					logger.Info(fmt.Sprintf("Duplicate key error: %v", writeErr.Message))
				}

			}
		}

		return fmt.Errorf("error inserting user into database: %w", err)

	}

	return nil
}
