package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateUser updates the user details in the database
func UpdateUser(ctx context.Context, userID string, updatedAdmin *models.UpdateAdmin) (*models.Admin, error) {
	// Initialize the MongoDB collection (replace "users" with your collection name)
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("User")

	// Create a filter to find the user by their ID
	filter := bson.M{"_id": userID}

	// Create an update document with the fields to update
	update := bson.M{
		"$set": bson.M{
			"username":   updatedAdmin.Username,
			"mobile":     updatedAdmin.Mobile,
			"email":      updatedAdmin.Email,
			"first_name": updatedAdmin.FirstName,
			"last_name":  updatedAdmin.LastName,
			"role":       updatedAdmin.Role,
		},
	}

	// Perform the update
	result := collection.FindOneAndUpdate(ctx, filter, update)

	// Check if the update was successful
	var updatedUser models.Admin
	err := result.Decode(&updatedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return an error if the user is not found
			return nil, fmt.Errorf("user not found")
		}
		// Return other errors (e.g., database issues)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Return the updated user
	return &updatedUser, nil
}
