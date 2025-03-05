package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOneUser retrieves a single user by their ID or username
func GetOneUser(ctx context.Context, userIDOrUsername string) (*models.ResponseAdmin, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("User")

	// Log the incoming request
	logger.Info("Attempting to find user with ID/Username:", "Use_ID", userIDOrUsername)

	// Query filter (search by either _id or username)
	var filter bson.M
	if isValidObjectID(userIDOrUsername) {
		objectID, err := primitive.ObjectIDFromHex(userIDOrUsername)

		if err != nil {
			logger.Error("Error converting string to ObjectId: %v", err)
			return nil, fmt.Errorf("invalid ObjectId format")
		}

		// If the user ID is a valid ObjectID, use it as the filter
		filter = bson.M{"_id": objectID}
		logger.Info("Using _id to query the database for user with _id:", "_id", userIDOrUsername)
	} else {
		// Otherwise, search by username
		filter = bson.M{"username": userIDOrUsername}
		logger.Info("Using username to query the database for user with username:", "usernaem", userIDOrUsername)
	}

	// Create an empty user struct to hold the result
	var user models.Admin

	// Query the database for the user
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Log that no document was found
			logger.Warn("User not found with ID/Username:", userIDOrUsername)
			// Return an error if the user is not found
			return nil, fmt.Errorf("user not found")
		}
		// Log any other database errors
		logger.Error("Failed to retrieve user with ID/Username:", userIDOrUsername, "Error:", err)
		// Return any other errors (e.g., database issues)
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	returnUser := models.ResponseAdmin{
		ID:        user.ID,
		Username:  user.Username,
		Mobile:    user.Mobile,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	// Log success and return the found user
	logger.Info("User found with ID/Username:", "UserID", userIDOrUsername)
	return &returnUser, nil
}

// Helper function to check if the user ID is a valid MongoDB ObjectID
func isValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
