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

func GetAllBlogCategory(ctx context.Context, params utils.PaginationParams) ([]models.BlogCategory, error) {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("BlogCategory")

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
			logger.Info("No Blog Category found in the database.")
			return nil, nil
		}
		logger.Error(fmt.Sprintf("Error retrieving Blog Category from database: %v", err))
		return nil, fmt.Errorf("error fetching Blog Category: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode all documents into the slice
	var blogCategories []models.BlogCategory
	if err = cursor.All(ctx, &blogCategories); err != nil {
		logger.Error(fmt.Sprintf("Error decoding Blog Category: %v", err))
		return nil, fmt.Errorf("error decoding Blog Category: %w", err)
	}

	logger.Info(fmt.Sprintf("Retrieved %d Blog Categories successfully.", len(blogCategories)))
	return blogCategories, nil
}

func DeleteBlogCategoryByID(ctx context.Context, id primitive.ObjectID) error {
	client := database.GetClient()
	collection := client.Database("uvfiberweb").Collection("BlogCategory")

	// Perform the delete operation
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error("Failed to delete blog category", "id", id, "error", err)
		return fmt.Errorf("failed to delete blog category: %w", err)
	}

	// Check if any document was deleted
	if result.DeletedCount == 0 {
		logger.Info("No blog category found to delete", "id", id)
		return fmt.Errorf("no blog category found with ID: %s", id.Hex())
	}

	logger.Info("Successfully deleted blog category", "id", id)
	return nil
}
