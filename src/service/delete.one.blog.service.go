package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteOneBlog(ctx context.Context, id string) error {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("Blog")

	// Convert the string ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Invalid blog ID format", "error", err)
		return fmt.Errorf("invalid blog ID format: %w", err)
	}

	// Build the query filter
	filter := bson.M{"_id": objectID}

	// Attempt to delete the document
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error("Error deleting blog from the database", "error", err, "blogID", id)
		return fmt.Errorf("error deleting blog: %w", err)
	}

	// Check if any document was deleted
	if result.DeletedCount == 0 {
		logger.Warn("No blog was deleted", "blogID", id)
		return fmt.Errorf("no blog found with the given ID: %s", id)
	}

	logger.Info("Blog deleted successfully", "blogID", id)
	return nil
}
