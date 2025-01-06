package service

import (
	"context"
	"fmt"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateOneBlog(ctx context.Context, id string, updateBlogData *models.BlogUpdate) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Blog")

	// Convert the string ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid blog ID format", "blogID", id, "error", err)
		return fmt.Errorf("invalid blog ID format: %w", err)
	}

	// Build the query filter
	filter := bson.M{"_id": objectID}

	// Update the `UpdatedAt` field
	updateBlogData.UpdatedAt = time.Now()

	// Prepare the update document
	update := bson.M{"$set": updateBlogData}

	// Attempt to update the document
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error updating blog in the database", "blogID", id, "error", err)
		return fmt.Errorf("error updating blog: %w", err)
	}

	// Check if any document was matched and updated
	if result.MatchedCount == 0 {
		logger.Warn("No blog found with the given ID for updating", "blogID", id)
		return fmt.Errorf("no blog found with the given ID: %s", id)
	}

	if result.ModifiedCount == 0 {
		logger.Warn("Blog update request did not modify any fields", "blogID", id)
		return fmt.Errorf("blog update did not modify any fields")
	}

	logger.Info("Blog updated successfully", "blogID", id)
	return nil
}
