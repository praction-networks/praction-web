package service

import (
	"context"
	"fmt"

	"github.com/praction-networks/quantum-ISP365/webapp/src/database"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetImageByName checks if an image with the given name exists in the database.
func GetImageByName(ctx context.Context, name string) error {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Image")

	filter := bson.M{"name": name}

	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Image does not exist, log and return nil (no error)
			logger.Info("Image not found", "name", name)
			return nil
		}
		// Log and return any other database error
		logger.Error("Error checking image existence", "name", name, "error", err)
		return err
	}

	// Image exists, log and return a validation error
	logger.Info("Image found", "name", name)
	return fmt.Errorf("image already exists with this name, please use a different name to upload the image")
}

func GetAllImageService(ctx context.Context, params utils.PaginationParams) ([]models.Image, error) {
	client := database.GetClient()
	collection := client.Database("practionweb").Collection("Image")

	// Build the query filter
	filter := bson.M{
		"isActive": true, // Base filter
	}

	// Dynamically add filters from params.Filters
	for key, value := range params.Filters {
		if key == "tag" {
			switch v := value.(type) {
			case []string: // Handle multiple tags
				filter[key] = bson.M{"$in": v}
			case string: // Handle single tag
				filter[key] = v
			default:
				logger.Error(fmt.Sprintf("Unexpected tag filter type: %T", v))
			}
		} else {
			filter[key] = value
		}
	}
	// Define options for pagination and sorting
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: params.SortField, Value: params.SortOrder}})
	findOptions.SetSkip(int64((params.Page - 1) * params.PageSize))
	findOptions.SetLimit(int64(params.PageSize))

	// Query the collection
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("Error retrieving Images from the database: %v", err))
		return nil, fmt.Errorf("error fetching images details: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice
	var images []models.Image
	if err = cursor.All(ctx, &images); err != nil {
		logger.Error(fmt.Sprintf("Error decoding images: %v", err))
		return nil, fmt.Errorf("error decoding images: %w", err)
	}

	// Log details
	logger.Info(fmt.Sprintf(
		"Retrieved %d image(s) with page=%d, pageSize=%d, sortField=%s, sortOrder=%d, filters=%v",
		len(images), params.Page, params.PageSize, params.SortField, params.SortOrder, params.Filters,
	))

	logger.Info("Images", images)

	// Return an empty slice if no results found
	if len(images) == 0 {
		return []models.Image{}, nil
	}

	return images, nil
}
