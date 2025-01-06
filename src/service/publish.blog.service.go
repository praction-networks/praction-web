package service

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PublishBlog updates the blog's state to published if the publish flag is true.
func PublishBlog(ctx context.Context, id string, publish bool) error {
	if !publish {
		logger.Warn("Publish flag is false; no changes will be made", "blogID", id)
		return nil
	}

	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Convert the string ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid blog ID format", "blogID", id, "error", err)
		return fmt.Errorf("invalid blog ID format: %w", err)
	}

	// Build the query filter
	filter := bson.M{"_id": objectID,
		"isApproved": true}

	// Prepare the update document
	update := bson.M{
		"$set": bson.M{
			"isActive":  true,
			"status":    "published",
			"updatedAt": time.Now(),
		},
	}

	// Perform the update operation
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error updating blog for publishing", "blogID", id, "error", err)
		return fmt.Errorf("error updating blog: %w", err)
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		logger.Warn("No blog found with the given ID for publishing", "blogID", id)
		return fmt.Errorf("no blog found with the given ID: %s", id)
	}

	if result.ModifiedCount == 0 {
		logger.Warn("Blog publish request did not modify any fields", "blogID", id)
		return fmt.Errorf("blog publish did not modify any fields")
	}

	logger.Info("Blog published successfully", "blogID", id)
	return nil
}
