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

func UserReferrelVerifyOTPAndUpdate(userOTP models.UserOTPVerify) (models.UserRefrence, string, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("UserReferal")

	// Context with timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"referredBy.email": userOTP.Email},
			{"referredBy.mobile": userOTP.Mobile},
		},
	}

	var existingUser models.UserRefrence
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No matching user found
			logger.Warn("No matching user found in the database", "Query", query)
			return models.UserRefrence{}, "User-not-found", nil
		}

		// Handle unexpected MongoDB errors
		logger.Error("Error querying the database", "Error", err)
		return models.UserRefrence{}, "", fmt.Errorf("database error: unable to fetch user data")
	}

	// If user is already verified
	if existingUser.IsVerified {
		logger.Info("User is already verified", "User", existingUser)
		return models.UserRefrence{}, "User-already-verified", nil
	}
	// Set the OTP expiration time limit (30 minutes)
	otpValidityDuration := 30 * time.Minute

	// Check if OTP has expired
	if time.Now().After(existingUser.OTPExpireTime.Add(otpValidityDuration)) {
		logger.Warn("OTP has expired", "Email", userOTP.Email, "Mobile", userOTP.Mobile)
		return models.UserRefrence{}, "OTP-expired", nil
	}
	// Validate OTP
	if existingUser.OTP != userOTP.OTP {
		logger.Warn("OTP does not match", "Email", userOTP.Email, "Mobile", userOTP.Mobile)
		return models.UserRefrence{}, "OTP-mismatch", nil
	}

	logger.Info("OTP matched successfully, updating user verification status", "UUID", existingUser.UUID)

	// Update user's verification status
	update := bson.M{
		"$set": bson.M{
			"updatedAt":    time.Now(),
			"VerifiedAt":   time.Now(),
			"isVerified":   true, // Mark the user as verified
			"refrelCoupon": utils.GenerateRefrelCoupon(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"uuid": existingUser.UUID}, update)

	if err != nil {
		logger.Error(fmt.Sprintf("Error updating user verification status for %s: %v", existingUser.ReferedBy.Name, err))
		return models.UserRefrence{}, "", fmt.Errorf("error updating user verification status: %w", err)
	}

	// Fetch the updated user data to include the RefrelCoupon
	err = collection.FindOne(ctx, bson.M{"uuid": existingUser.UUID}).Decode(&existingUser)
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching updated user data: %v", err))
		return models.UserRefrence{}, "", fmt.Errorf("error fetching updated user data: %w", err)
	}

	logger.Info("User verification updated successfully", "UUID", existingUser.UUID)
	return existingUser, "Verification-successful", nil
}
