package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllUsers retrieves all users from the database
func GetAllUsers(ctx context.Context) ([]models.ResponseAdmin, error) {
	// Initialize the MongoDB collection (replace "users" with your collection name)
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("User")

	// Query all users
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.ResponseAdmin
	for cursor.Next(ctx) {
		var user models.Admin
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		responseUser := models.ResponseAdmin{
			ID:        user.ID,
			Username:  user.Username,
			Mobile:    user.Mobile,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
		}
		users = append(users, responseUser)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}
