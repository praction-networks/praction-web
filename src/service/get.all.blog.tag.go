package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllBlogTag(ctx context.Context, params utils.PaginationParams) ([]models.BlogTag, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("BlogTag")

	// Build the query filters
	filter := bson.M{
		"isActive":  true,  // Only include active tags
		"isDeleted": false, // Exclude deleted tags
	}
	for key, value := range params.Filters {
		filter[key] = value
	}

	// Define find options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	// Query to fetch blog categories
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Info("No Blog Tag found in the database.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving Blog Tag from database: %v", err))
		return nil, fmt.Errorf("error fetching Blog Tag: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode all documents into the slice
	var blogTag []models.BlogTag
	if err = cursor.All(ctx, &blogTag); err != nil {
		logger.Error(fmt.Sprintf("Error decoding Blog Tag: %v", err))
		return nil, fmt.Errorf("error decoding Blog Tag: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved %d Blog Tag successfully.", len(blogTag)))
	return blogTag, nil
}

func DeleteBlogTagByID(ctx context.Context, id primitive.ObjectID) error {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("BlogTag")

	// Perform the delete operation
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error("Failed to delete blog tag", "id", id, "error", err)
		return fmt.Errorf("failed to delete blog tag: %w", err)
	}

	// Check if any document was deleted
	if result.DeletedCount == 0 {
		logger.Info("No blog tag found to delete", "id", id)
		return fmt.Errorf("no blog tag found with ID: %s", id.Hex())
	}

	logger.Info("Successfully deleted blog tag", "id", id)
	return nil
}
