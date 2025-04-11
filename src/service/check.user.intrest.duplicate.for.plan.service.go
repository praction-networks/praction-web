package service

import (
	"context"
	"errors"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckUserIntrestDuplicate checks if a user's email or mobile exists in the database
// and whether it is verified.
func CheckUserIntrestForPlanDuplicate(user *models.AvailableUserRequest) (models.AvailableUserRequest, string, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("UserPlanIntrest")

	// Context with timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query to find a record with matching email or mobile
	query := bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"mobile": user.Mobile},
		},
	}

	var existingUser models.AvailableUserRequest
	err := collection.FindOne(ctx, query).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No matching document found
			logger.Warn("No matching user found in database", "Query", query)
			return models.AvailableUserRequest{}, "NotFound", nil
		}

		// Log and return unexpected MongoDB errors
		logger.Error("Error querying the database", "Error", err)
		return models.AvailableUserRequest{}, "", errors.New("database error: unable to fetch user data")
	}

	if existingUser.PlanID != user.PlanID {
		logger.Warn("Plen", "Query", query)
		return models.AvailableUserRequest{}, "NotFound", nil

	}
	// Check if the found record is verified
	if existingUser.IsVerified {
		logger.Info("Duplicate verified user found", "User", existingUser)
		return existingUser, "Verified", nil
	}

	// Not verified
	logger.Info("Duplicate user found but not verified", "User", existingUser)
	return existingUser, "NotVerified", nil
}
