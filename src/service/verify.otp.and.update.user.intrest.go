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

func UserInterestVerifyOTPAndUpdate(userOTP models.UserOTPVerify) (models.UserInterest, string, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("UserIntrest")

	// Context with timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"email": userOTP.Email},
			{"mobile": userOTP.Mobile},
		},
	}

	var existingUser models.UserInterest
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No matching user found
			logger.Warn("No matching user found in the database", "Query", query)
			return models.UserInterest{}, "User-not-found", nil
		}

		// Handle unexpected MongoDB errors
		logger.Error("Error querying the database", "Error", err)
		return models.UserInterest{}, "", fmt.Errorf("database error: unable to fetch user data")
	}

	// If user is already verified
	if existingUser.IsVerified {
		logger.Info("User is already verified", "User", existingUser)
		return models.UserInterest{}, "User-already-verified", nil
	}

	otpValidityDuration := 30 * time.Minute

	// Check if OTP has expired
	if time.Now().After(existingUser.OTPExpireTime.Add(otpValidityDuration)) {
		logger.Warn("OTP has expired", "Email", userOTP.Email, "Mobile", userOTP.Mobile)
		return models.UserInterest{}, "OTP-expired", nil
	}

	// Validate OTP
	if existingUser.OTP != userOTP.OTP {
		logger.Warn("OTP does not match", "Email", userOTP.Email, "Mobile", userOTP.Mobile)
		return models.UserInterest{}, "OTP-mismatch", nil
	}

	logger.Info("OTP matched successfully, updating user verification status", "User", existingUser)

	// Update user's verification status
	update := bson.M{
		"$set": bson.M{
			"updatedAt":     time.Now(),
			"verifiedAt":    time.Now(),
			"isVerified":    true, // Mark the user as verified
			"otpExpireTime": nil,
			"otp":           nil,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"uuid": existingUser.UUID}, update)

	if err != nil {
		logger.Error(fmt.Sprintf("Error updating user verification status for %s: %v", existingUser.Name, err))
		return models.UserInterest{}, "", fmt.Errorf("error updating user verification status: %w", err)
	}

	logger.Info("User verification updated successfully", "User", existingUser)
	return existingUser, "Verification-successful", nil
}
