package service

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserDetailsFromDBRefeeral(userData *models.UserOTPResend) (models.UserRefrence, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("UserReferal")

	// Context with timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"referredBy.email": userData.Email},
			{"referredBy.mobile": userData.Mobile},
		},
	}

	var existingUser models.UserRefrence
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No matching user found
			logger.Warn("No matching user found in the database", "Query", query)
			return models.UserRefrence{}, fmt.Errorf("user not found in database")
		}

		// Handle unexpected MongoDB errors
		logger.Error("Error querying the database", "Error", err)
		return models.UserRefrence{}, fmt.Errorf("database error: unable to fetch user data")
	}

	// Generate a new OTP (Assuming you have a utility function for OTP generation)
	newOTP := utils.GenerateRandomOTP(6)

	// Set OTP expiration time (30 minutes from now)
	otpExpireTime := time.Now().Add(30 * time.Minute)

	// Update the user's OTP and OTPExpireTime in the database
	update := bson.M{
		"$set": bson.M{
			"otp":           newOTP,
			"OTPExpireTime": otpExpireTime,
		},
	}

	// Perform the update
	_, err = collection.UpdateOne(ctx, bson.M{"uuid": existingUser.UUID}, update)

	if err != nil {
		logger.Error(fmt.Sprintf("Error updating OTP for user %s: %v", existingUser.ReferedBy.Name, err))
		return models.UserRefrence{}, fmt.Errorf("error updating OTP: %w", err)
	}

	// Fetch the updated user from the database to reflect the new OTP
	err = collection.FindOne(ctx, bson.M{"uuid": existingUser.UUID}).Decode(&existingUser)
	if err != nil {
		logger.Error("Error fetching updated user details", "Error", err)
		return models.UserRefrence{}, fmt.Errorf("error fetching updated user details: %w", err)
	}

	// Return the updated user details along with the new OTP
	return existingUser, nil

}
