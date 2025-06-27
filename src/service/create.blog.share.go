package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBlogShare(ctx context.Context, BlogID primitive.ObjectID) error {

	//Check If Blog is Valid or Not

	err := checkBlogAndAddShare(ctx, BlogID)
	if err != nil {
		if errors.Is(err, ErrNoBlogFound) {
			logger.Info(fmt.Sprintf("No blog found with ID: %v", BlogID))
			return fmt.Errorf("blog with ID %v does not exist", BlogID)
		}
		if errors.Is(err, ErrDatabaseOperation) {
			logger.Error(fmt.Sprintf("Database error while checking blog availability: %v", err))
			return fmt.Errorf("internal database error occurred")
		}

		// Handle other unexpected errors
		logger.Error(fmt.Sprintf("Unexpected error: %v", err))
		return fmt.Errorf("unexpected error occurred: %v", err)
	}

	logger.Info(fmt.Sprintf("Blog Share incremented successfully for blog ID %s created successfully.", BlogID))
	return nil
}

func checkBlogAndAddShare(ctx context.Context, blogID primitive.ObjectID) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Create a filter to search for the blog by its ID
	filter := bson.D{{Key: "_id", Value: blogID}}

	// Update to increment the view count
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "shares", Value: 1}}}}

	// Attempt to find the blog
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Error updating share count for blog with ID %v: %v", blogID, err))
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	// Check if the blog exists
	if result.MatchedCount == 0 {
		logger.Info(fmt.Sprintf("No blog found with ID: %v", blogID))
		return ErrNoBlogFound
	}

	logger.Info(fmt.Sprintf("Successfully incremented share count for blog with ID %v", blogID))
	return nil
}
