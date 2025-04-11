package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetUserPassword allows a user to change their password
func SetUserPassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Initialize the MongoDB collection (replace "users" with your collection name)
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("User")

	// Find the user by their ID
	var user models.Admin
	err := collection.FindOne(ctx, bson.M{"username": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Validate the password
	if !validatePassword(currentPassword, user.Password, user.Salt) {
		// Return a nil pointer to User, and the error
		return errors.New("current Password is not matched")
	}
	// Hash the password using Argon2
	salt, err := utils.GenerateSalt()
	if err != nil {
		logger.Error(fmt.Sprintf("Error generating salt: %v", err))
		return fmt.Errorf("error generating salt: %w", err)
	}
	hashedPassword := hashPassword(newPassword, salt)

	// Update the user's password in the database
	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}
	_, err = collection.UpdateOne(ctx, bson.M{"username": userID}, update)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Return success
	return nil
}
