package service

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserOTP(ctx context.Context, userIntrest models.UserInterest) error {

	// Check if the user exists in the database
	err := updateOTPInDB(ctx, userIntrest)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update OTP for User Interest: %v", err))
		return fmt.Errorf("failed to update OTP: %w", err)
	}

	logger.Info(fmt.Sprintf("OTP updated successfully for user %s.", userIntrest.Name))
	return nil
}

// updateOTPInDB updates the OTP and verification status for the user based on mobile or email
func updateOTPInDB(ctx context.Context, userIntrest models.UserInterest) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("UserIntrest")

	// Define the query for checking if the user exists based on email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"email": userIntrest.Email},
			{"mobile": userIntrest.Mobile},
		},
	}

	// Check if the user already exists
	var existingUser models.UserInterest
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No user found, return an error
			logger.Warn(fmt.Sprintf("No matching user found for mobile %s or email %s", userIntrest.Mobile, userIntrest.Email))
			return fmt.Errorf("no matching user found")
		}
		// Log any other errors
		logger.Error(fmt.Sprintf("Error checking user existence: %v", err))
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// If user exists, update the OTP and set IsVerified as false (if not already verified)
	update := bson.M{
		"$set": bson.M{
			"updatedAt":  time.Now(),
			"IsVarified": false, // Keep it false since this is an OTP update
		},
	}

	// Update the OTP field for the matching user
	_, err = collection.UpdateOne(ctx, bson.M{"uuid": existingUser.UUID}, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Error updating OTP for user %s: %v", userIntrest.Name, err))
		return fmt.Errorf("error updating OTP: %w", err)
	}

	logger.Info(fmt.Sprintf("Updated OTP for user %s successfully.", userIntrest.Name))
	return nil
}
