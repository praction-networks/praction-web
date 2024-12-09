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

func GetOneBlog(ctx context.Context, id string) (*models.Blog, error) {

	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Convert the string ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid blog ID format: %v", err))
		return nil, fmt.Errorf("invalid blog ID format: %w", err)
	}
	// Build the query filter
	filter := bson.M{
		"_id":        objectID,
		"isApproved": true,
		"isActive":   true,
		"isDeleted":  false,
		"status":     "published",
	}

	// Query the collection
	var blog models.Blog
	err = collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info(fmt.Sprintf("No blog found with ID: %s", id))
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving blog from the database: %v", err))
		return nil, fmt.Errorf("error fetching blog details: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved blog with ID: %s successfully.", id))
	return &blog, nil
}
