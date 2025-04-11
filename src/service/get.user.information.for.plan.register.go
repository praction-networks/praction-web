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

func GetUserDetailsFromDBForPlan(userDara *models.UserOTPPlanResend) (models.AvailableUserRequest, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("UserPlanIntrest")

	// Context with timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email or mobile
	query := bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"email": userDara.Email},
					{"mobile": userDara.Mobile},
				},
			},
			{"planID": userDara.PlanID},
		},
	}

	var existingUser models.AvailableUserRequest
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No matching user found
			logger.Warn("No matching user found in the database", "Query", query)
			return models.AvailableUserRequest{}, fmt.Errorf("user not found in database")
		}

		// Handle unexpected MongoDB errors
		logger.Error("Error querying the database", "Error", err)
		return models.AvailableUserRequest{}, fmt.Errorf("database error: unable to fetch user data")
	}

	newExpireTime := time.Now().Add(30 * time.Minute) // Set new expiration time (e.g., 5 minutes from now)

	update := bson.M{
		"$set": bson.M{
			"otpExpireTime": newExpireTime,
		},
	}

	_, updateErr := collection.UpdateOne(ctx, query, update)
	if updateErr != nil {
		logger.Error("Failed to update OTP expiration time", "Error", updateErr)
		return models.AvailableUserRequest{}, fmt.Errorf("database error: unable to update OTP expiration time")
	}

	logger.Info("OTP expiration time updated successfully", "User", existingUser.Email, "NewExpireTime", newExpireTime)

	return existingUser, nil

}
